package main

import (
	"context"
	"github.com/StefanPahlplatz/tempus/auth"
	"github.com/StefanPahlplatz/tempus/auth/grpc"
)

type authServiceController struct {
	grpc.UnimplementedAuthServiceServer
	authService auth.Service
}

func NewAuthServiceController(authService auth.Service) grpc.AuthServiceServer {
	return &authServiceController{
		authService: authService,
	}
}

func (a *authServiceController) Login(context.Context, *grpc.LoginRequest) (*grpc.LoginResponse, error) {
	return &grpc.LoginResponse{}, nil
}
