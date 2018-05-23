package config

const (
	Version     = "1.0.0"
	ProjectName = "Sirius"
	AppName     = "Actor"
	TraceKey    = "request_id"

	// Header name
	HeaderRequestID = "X-Request-ID"

	CountryJSON = "countries.json"

	// Resource of error descriptions
	ErrEndpointNotFound    = "endpoint your requested not found"
	ErrDatabaseNil         = "database object is nil"
	ErrDatabaseConnectFail = "failed connected to database"
	ErrPingDatabaseFail    = "ping database connection failed"

	ErrCastingData = "failed casting data type"
)
