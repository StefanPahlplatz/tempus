package auth

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func test() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen on port 8000: %v", err)
	}

	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to service gRPC server over port 8000: %v", err)
	}
}
