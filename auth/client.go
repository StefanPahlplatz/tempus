package auth

import (
	auth "github.com/StefanPahlplatz/tempus/auth/claim"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// NewClient returns a gRPC client for interacting with auth service.
// After calling it, run a defer close on the close function
func NewClient() (auth.AuthServiceClient, func() error, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, errors.Wrap(err, "did not connect")
	}
	return auth.NewAuthServiceClient(conn), conn.Close, nil
}
