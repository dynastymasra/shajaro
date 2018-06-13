package config

const (
	Version     = "1.0.0"
	ProjectName = "Shajaro"
	AppName     = "Actor"
	TraceKey    = "request_id"
	RedirectURI = "https://auth.shajaro.com"

	// Header name
	HeaderRequestID  = "X-Request-ID"
	ScopeHeader      = "X-Authenticated-Scope"
	AuthUserIDHeader = "X-Authenticated-Userid"
	ConsumerIDHeader = "X-Consumer-Id"

	CountryJSON = "countries.json"

	ActorScopes = "actor.read actor.update actor.delete"
	ActorRead   = "actor.read"
	ActorUpdate = "actor.update"
	ActorDelete = "actor.delete"

	// Resource of error descriptions
	ErrEndpointNotFound    = "endpoint your requested not found"
	ErrDatabaseNil         = "database object is nil"
	ErrDatabaseConnectFail = "failed connected to database"
	ErrPingDatabaseFail    = "ping database connection failed"
	ErrFailedSaveNewUser   = "failed to save new user"
	ErrFailedLogin         = "failed login user, check your email and password"

	ErrCastingData = "failed casting data type"
)
