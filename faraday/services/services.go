package services

// Configuration for back-end services

const (
	// Public security means a user may be logged out or logged in
	Public = iota
	// Authenticated security means a user must be logged in
	Authenticated = iota
	// Admin security means a user must be both logged in and have sudo flag
	Admin = iota
)

// ServiceDirectory allows access to a backend service using its subdomain
type ServiceDirectory map[string]Service

// Service is an app on Staffjoy that runs on a subdomain
type Service struct {
	Security      int    // Public, Authenticated, or Admin
	RestrictDev   bool   // If True, service is suppressed in stage and prod
	BackendDomain string // Backend service to query
	NoCacheHTML   bool   // If True, injects a header for HTML responses telling the browser not to cache HTML

}

// StaffjoyServices is a map of subdomains -> specs
// Sudomain is <string> + Env["rootDomain"]
// e.g. "login" service on prod is "login" + "staffjoy.com""
//
// KEEP THIS LIST IN ALPHABETICAL ORDER please
var StaffjoyServices = ServiceDirectory{
	"auth": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "localhost:8001",
	},
	"faraday": {
		// Debug site for faraday
		Security:      Admin,
		RestrictDev:   true,
		BackendDomain: "httpbin.org",
	},
	"www": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "www-service",
	},
}
