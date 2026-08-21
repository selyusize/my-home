package git

// Provider identifies a git hosting backend.
type Provider string

const (
	ProviderUnknown Provider = ""
	ProviderGitHub  Provider = "github"
	ProviderGitLab  Provider = "gitlab"
)
