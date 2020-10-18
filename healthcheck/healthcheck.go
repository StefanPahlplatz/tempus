// healthcheck is a library that provides a basic health check handler for Staffjoy applications.
// We generally host this endpoint at "/health" on port 80
//
// Usage:
// r.HandleFunc(healthcheck.HEALTHPATH healthcheck.Handler)

package healthcheck

import (
	"encoding/json"
	"github.com/StefanPahlplatz/tempus/environments"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	// HEALTHPATH is the standard healthcheck path in our app
	ServiceName string = "healthcheck"
	HEALTHPATH  string = "/health"
)

var (
	logger *logrus.Entry
	config environments.Config
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

// Handler returns a basic JSON
func Handler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	// We shouldn't have any errors
	msg, _ := json.Marshal(map[string]string{"status": "ok"})
	_, err := res.Write(msg)
	if err != nil {
		logger.Fatalf("healthcheck.handler: unable to write response: %v", err)
	}
	return
}
