package config

const (
	Version     = "1.0.0"
	ProjectName = "Gapura"
	AppName     = "Actor"
	TraceKey    = "request_id"
	RedirectURI = "https://www.sirius.co.id"

	// Header name
	HeaderRequestID = "X-Request-ID"

	CountryJSON = "countries.json"

	// Resource of error descriptions
	ErrEndpointNotFound    = "endpoint your requested not found"
	ErrDatabaseNil         = "database object is nil"
	ErrDatabaseConnectFail = "failed connected to database"
	ErrPingDatabaseFail    = "ping database connection failed"
	ErrFailedSaveNewUser   = "failed to save new user"

	ErrCastingData = "failed casting data type"
)
