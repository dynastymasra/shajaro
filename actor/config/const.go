package config

const (
	Version     = "1.0.0-dev"
	ProjectName = "Shajaro"
	AppName     = "Actor"
	TraceKey    = "request_id"
	RedirectURI = "https://auth.shajaro.com"

	// Header name
	HeaderRequestID = "X-Request-ID"

	CountryJSON = "countries.json"

	// Resource of error descriptions
	ErrEndpointNotFound    = "endpoint your requested not found"
	ErrDatabaseNil         = "database object is nil"
	ErrDatabaseConnectFail = "failed connected to database"
	ErrPingDatabaseFail    = "ping database connection failed"
	ErrFailedSaveNewUser   = "failed to save new user"
	ErrFailedLogin         = "failed login user, check your email and password"

	ErrCastingData = "failed casting data type"
)
