package db

// Provider identifies a database backend.
type Provider string

const (
	ProviderUnknown  Provider = ""
	ProviderPostgres Provider = "postgres"
)
