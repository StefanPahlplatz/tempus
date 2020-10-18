package main

import (
	auth2 "github.com/StefanPahlplatz/tempus/auth"
	auth "github.com/StefanPahlplatz/tempus/auth/claim"
	"github.com/StefanPahlplatz/tempus/environments"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	// ServiceName is how this app is identified in logs and error handlers
	ServiceName string = "auth"
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

// Listen for incoming requests, then validate, sanitize, and route them.
func main() {
	lis, err := net.Listen("tcp", auth2.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", auth2.ServerPort, err)
	}

	s := auth.Server{}

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, &s)

	logger.Infof("Initialized environment %s", config.Name)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to service gRPC server over port 8000: %v", err)
	}
}
