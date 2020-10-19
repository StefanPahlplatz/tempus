package auth

import (
	"fmt"
	"github.com/StefanPahlplatz/tempus/environments"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/metadata"
)

var (
	signingSecret string
	shortSession  = time.Duration(12 * time.Hour)
	longSession   = time.Duration(30 * 24 * time.Hour)
	config        environments.Config
)

func init() {
	signingSecret = os.Getenv("SIGNING_SECRET")

	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
}

// SetInternalHeaders is used by Faraday to sanitize incoming external requests
// and convert them to internal requests with authorization information
func SetInternalHeaders(externalReq *http.Request, internalHeaders http.Header) {
	ProxyHeaders(externalReq.Header, internalHeaders)

	// default to anonymous web then prove otherwise
	//authorization := AuthorizationAnonymousWeb
	//uuid, support, err := getSession(externalReq)
	//if err != nil {
	//	internalHeaders.Set(AuthorizationHeader, authorization)
	//	return
	//}
	//
	//if support {
	//	authorization = AuthorizationSupportUser
	//} else {
	//	authorization = AuthorizationAuthenticatedUser
	//}
	//internalHeaders.Set(currentUserHeader, uuid)
}

// ProxyHeaders copies http headers
func ProxyHeaders(from, to http.Header) {
	// Range over the headres
	for k, v := range from {
		// TODO - filter restricted headers

		// Multiple header values may exist per key
		for _, x := range v {
			to.Add(k, x)
		}

	}
}

// GetCurrentUserUUIDFromMetadata allows backend gRPC services with
// authorization methods of AuthenticatedUser or SupportUser to access
// the uuid of the user making the request
func GetCurrentUserUUIDFromMetadata(data metadata.MD) (uuid string, err error) {
	res, ok := data[currentUserMetadata]
	if !ok || len(res) == 0 {
		err = fmt.Errorf("User not authenticated")
		return
	}
	uuid = res[0]
	return
}

// GetCurrentUserUUIDFromHeader allows backend http services with
// authorization methods of AuthenticatedUser or SupportUser to access
// the uuid of the user making the request
func GetCurrentUserUUIDFromHeader(data http.Header) (uuid string, err error) {
	res, ok := data[currentUserHeader]
	if !ok || len(res) == 0 {
		err = fmt.Errorf("User not authenticated")
		return
	}
	uuid = res[0]
	return
}
