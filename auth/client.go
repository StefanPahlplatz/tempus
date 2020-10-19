package auth

import (
	protos "github.com/StefanPahlplatz/tempus/auth/protos"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// NewClient returns a gRPC client for interacting with auth service.
// After calling it, run a defer close on the close function
func NewClient() (protos.AuthServiceClient, func() error, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, errors.Wrap(err, "did not connect")
	}
	return protos.NewAuthServiceClient(conn), conn.Close, nil
}
