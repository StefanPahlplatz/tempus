package auth

import (
	"context"
	"log"
)

type Server struct {
}

func (s *Server) Authenticate(ctx context.Context, authReq *AuthenticateRequest) (*Claim, error) {
	log.Printf("Received grpc request")
	return &Claim{
		Role: "admin",
	}, nil
}
