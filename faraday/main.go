// Faraday proxies all requests to Staffjoy
package main

import (
	"bytes"
	"fmt"
	"github.com/StefanPahlplatz/tempus/auth"
	protos "github.com/StefanPahlplatz/tempus/auth/protos"
	"github.com/StefanPahlplatz/tempus/middlewares"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/StefanPahlplatz/tempus/environments"
	"github.com/StefanPahlplatz/tempus/errorpages"
	"github.com/StefanPahlplatz/tempus/faraday/services"
	"github.com/StefanPahlplatz/tempus/healthcheck"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	// ServiceName is how this app is identified in logs and error handlers
	ServiceName          string = "faraday"
	userID               string = "userID"
	userSudo             string = "userSudo"
	requestAuthenticated string = "requestAuthenticated"
	requestedService     string = "service"
)

var (
	logger       *logrus.Entry
	config       environments.Config
	signingToken = os.Getenv("SIGNING_SECRET")
	bannedUsers  = map[string]string{ // Use a map for constant tempus lookups. Value doesn't matter
		// Hypothetically these should be universally unique, so we don't have to limit by env
		"d7b9dbed-9719-4856-5f19-23da2d0e3dec": "hidden",
	}
)

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)
}

// Listen for incoming requests, then validate, sanitize, and route them.
func main() {
	logger.Infof("Initialized environment %s", config.Name)

	r := NewRouter(config, logger)
	// Set up http internal
	// Note - we do this without the Negroni convenience func so that
	// we can add in TLS support in the future too.
	s := &http.Server{
		Addr:           ":8000",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// TODO - add in a logging system and have it do a fatal call here
	logger.Panicf("main: error while serving: %v", s.ListenAndServe())
}

// NewRouter returns a router composed of internal and external parts
func NewRouter(config environments.Config, logger *logrus.Entry) http.Handler {
	// Create a new router. We use Gorilla instead of stdlib because it handles
	// memory clean up for the 'context' package correctly
	externalRouter := mux.NewRouter()
	internalRouter := mux.NewRouter().PathPrefix("/").Subrouter().StrictSlash(true)

	// Make this available always, e.g. for kubernetes health checks
	externalRouter.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	externalRouter.HandleFunc(MobileConfigPath, MobileConfigHandler)

	sentryPublicDSN, err := environments.GetPublicSentryDSN(config.GetSentryDSN())
	if err != nil {
		logger.Fatalf("Cannot get sentry info - %s", err)
	}

	// only apply security to the internal routes
	externalRouter.PathPrefix("/").Handler(negroni.New(
		middlewares.NewRecovery(ServiceName, config, sentryPublicDSN),
		NewSecurityMiddleware(config),
		NewServiceMiddleware(config, services.StaffjoyServices),
		negroni.Wrap(internalRouter),
	))
	internalRouter.PathPrefix("/").HandlerFunc(proxyHandler)

	return externalRouter
}

// HTTP function that handles proxying after all of the middlewares
func proxyHandler(res http.ResponseWriter, req *http.Request) {
	// Get the service as set in the middleware.
	service := req.Context().Value(requestedService).(services.Service)

	// No security on backend right now :-(
	destination := "http://" + service.BackendDomain + req.URL.RequestURI()
	logger.Debugf("Proxying to %s", destination)

	// Get the body to pass to the new request.
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("Could not read request body - %s", err))
	}

	// Create the internal request to send to send via grpc.
	internalReq, err := http.NewRequest(req.Method, destination, bytes.NewReader(b))
	if err != nil {
		panic(fmt.Sprintf("Unable to create request - %s", err))
	}

	auth.SetInternalHeaders(req, internalReq.Header)

	authClient, closeClient, err := auth.NewClient()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to auth internal - %s", err))
	}
	defer closeClient()

	a, err := authClient.Authenticate(req.Context(), &protos.AuthenticateRequest{
		Email:    "",
		Password: "",
	})
	if err != nil {
		logger.Panicf("Unable to authenticate - %s", err)
	}
	logger.Infof("auth response: %v", a)

	currentUserUUID, err := auth.GetCurrentUserUUIDFromHeader(internalReq.Header)
	if err == nil {
		// authenticated request
		if _, isBanned := bannedUsers[currentUserUUID]; isBanned {
			logger.Warningf("Banned user accessing service - user %s", currentUserUUID)
			errorpages.Forbidden(res)
			return
		}
	}

	// Right here - check response Authorization and see if it's ok
	// with the requested service

	// Check perimeter authorization
	switch a.GetRole() {
	case auth.AuthorizationAnonymousWeb:
		if service.Security != services.Public {
			// send to login
			scheme := "https"
			if config.Name == "development" || config.Name == "test" {
				scheme = "http"
			}
			redirectDest := &url.URL{Host: "www." + config.ExternalApex, Scheme: scheme, Path: "/login/"}

			url := req.Host + req.URL.EscapedPath()

			http.Redirect(res, req, redirectDest.String()+"?return_to="+url, http.StatusTemporaryRedirect)
			return
		}
	case auth.AuthorizationAuthenticatedUser:
		if service.Security == services.Admin {
			errorpages.Forbidden(res)
			return
		}
	case auth.AuthorizationSupportUser:
		// no restrictions
	default:
		logger.Panicf("proxyHandler: unknown authorization header: %s", a.GetRole())
	}

	client := http.Client{
		// RETURN a redirect, do not FOLLOW it (which ends up causing relative redirect issues)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	internalRes, err := client.Do(internalReq)
	if err != nil {
		logger.Warningf("Unable to query backend - %s", err)
		errorpages.GatewayTimeout(res)
		return
	}
	// Copy headers from service to user
	auth.ProxyHeaders(internalRes.Header, res.Header())

	if service.NoCacheHTML {
		if strings.Contains(strings.Join(res.Header()["Content-Type"], ""), "text/html") {
			// insert header to prevent caching
			res.Header().Set("Cache-Control", "no-cache")
		}
	}

	res.WriteHeader(internalRes.StatusCode)
	io.Copy(res, internalRes.Body)
	internalRes.Body.Close()

}
