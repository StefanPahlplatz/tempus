package server

import (
	"context"
	"github.com/StefanPahlplatz/tempus/auth"
	protos "github.com/StefanPahlplatz/tempus/auth/protos"
	"github.com/sirupsen/logrus"
)

type Server struct {
	protos.UnimplementedAuthServiceServer
	logger *logrus.Entry
}

func NewServer(logger *logrus.Entry) *Server {
	return &Server{logger: logger}
}

func (s *Server) Authenticate(ctx context.Context, in *protos.AuthenticateRequest) (*protos.Claim, error) {
	s.logger.Info("Handle Authenticate for email:", in.GetEmail())

	return &protos.Claim{Role: auth.AuthorizationAuthenticatedUser}, nil
}
