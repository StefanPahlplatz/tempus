package client

import (
	authgrpc "github.com/StefanPahlplatz/tempus/auth/grpc"
	"google.golang.org/grpc"
)

func NewAuthClient(connString string) (authgrpc.AuthServiceClient, func() error, error) {
	conn, err := grpc.Dial(connString, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return authgrpc.NewAuthServiceClient(conn), conn.Close, nil
}
