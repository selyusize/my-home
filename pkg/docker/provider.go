package docker

// Provider identifies a container engine backend.
type Provider string

const (
	ProviderUnknown Provider = ""
	ProviderEngine  Provider = "engine"
)
