package auth

const (
	// ServerPort tells the gRPC internal what port to listen on
	ServerPort = ":8001"

	// Endpoint defines the DNS of the account internal for clients
	// to access the internal in Kubernetes.
	Endpoint = "localhost" + ServerPort

	cookieName = "staffjoy-faraday"
	cookie
	uuidKey       = "uuid"
	supportKey    = "support"
	expirationKey = "exp"
	// for GRPC
	currentUserMetadata = "faraday-current-user-uuid"
	// header set for internal user id
	currentUserHeader = "Grpc-Metadata-Faraday-Current-User-Uuid"

	// AuthorizationHeader is the http request header
	// key used for accessing the internal authorization.
	AuthorizationHeader = "Authorization"

	// AuthorizationMetadata is the grpce metadadata key used
	// for accessing the internal authorization
	AuthorizationMetadata = "authorization"

	// AuthorizationAnonymousWeb is set as the Authorization header to denote that
	// a request is being made bu an unauthenticated web user
	AuthorizationAnonymousWeb = "faraday-anonymous"

	// AuthorizationAuthenticatedUser is set as the Authorization header to denote that
	// a request is being made by an authenticated web user
	AuthorizationAuthenticatedUser = "faraday-authenticated"

	// AuthorizationSupportUser is set as the Authorization header to denote that
	// a request is being made by a Staffjoy team me
	AuthorizationSupportUser = "faraday-support"

	// AuthorizationWWWService is set as the Authorization header to denote that
	// a request is being made by the www login / signup system
	AuthorizationWWWService = "www-service"

	// AuthorizationCompanyService is set as the Authorization header to denote
	// that a request is being made by the company api/internal
	AuthorizationCompanyService = "company-service"

	// AuthorizationSuperpowersService is set as the Authorization header to
	// denote that a request is being made by the dev-only superpowers service
	AuthorizationSuperpowersService = "superpowers-service"

	// AuthorizationWhoamiService is set as the Authorization heade to denote that
	// a request is being made by the whoami microservice
	AuthorizationWhoamiService = "whoami-service"

	// AuthorizationBotService is set as the Authorization header to denote that
	// a request is being made by the bot microservice
	AuthorizationBotService = "bot-service"

	// AuthorizationAccountService is set as the Authorization header to denote that
	// a request is being made by the account service
	AuthorizationAccountService = "account-service"

	// AuthorizationICalService is set as the Authorization header to denote that
	// a request is being made by the ical service
	AuthorizationICalService = "ical-service"
)
