package main

import (
	"github.com/StefanPahlplatz/tempus/auth/core"
	authgrpc "github.com/StefanPahlplatz/tempus/auth/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	// configure our core service
	userService := core.NewService()

	// configure our gRPC service controller
	userServiceController := NewAuthServiceController(userService)

	// start a gRPC server
	server := grpc.NewServer()
	authgrpc.RegisterAuthServiceServer(server, userServiceController)

	con, err := net.Listen("tcp", os.Getenv("GRPC_ADDR"))
	if err != nil {
		panic(err)
	}

	log.Printf("Starting gRPC user service on %s...\n", con.Addr().String())
	err = server.Serve(con)
	if err != nil {
		panic(err)
	}
}
