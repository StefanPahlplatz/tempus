package auth

const (
	// ServerPort tells the gRPC server what port to listen on
	ServerPort = ":8001"

	// Endpoint defines the DNS of the account server for clients
	// to access the server in Kubernetes.
	Endpoint = "auth-service" + ServerPort
)
