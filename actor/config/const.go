package config

const (
	Version     = "1.0.0"
	ProjectName = "Shajaro"
	AppName     = "Actor"
	TraceKey    = "trace_id"
	RedirectURI = "https://auth.shajaro.com"

	// Header name
	HeaderTraceID    = "X-Trace-ID"
	ScopeHeader      = "X-Authenticated-Scope"
	AuthUserIDHeader = "X-Authenticated-Userid"
	ConsumerIDHeader = "X-Consumer-Id"

	CountryJSON = "countries.json"

	ActorScopes = "actor.read actor.update actor.delete"
	ActorRead   = "actor.read"
	ActorUpdate = "actor.update"
	ActorDelete = "actor.delete"

	// Resource of error descriptions
	ErrEndpointNotFound = "endpoint your requested not found"
	ErrDatabaseNil      = "database object is nil"

	ErrCastingData = "failed casting data type"
)
