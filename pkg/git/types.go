package git

import "time"

// Config holds OAuth and instance settings shared by git providers.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	// BaseURL is an optional instance URL (used by self-hosted GitLab).
	BaseURL string
}

// User is a simplified git host profile.
type User struct {
	ID          int64     `json:"id"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Bio         string    `json:"bio"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Blog        string    `json:"blog"`
	AvatarURL   string    `json:"avatarUrl"`
	HTMLURL     string    `json:"htmlUrl"`
	PublicRepos int       `json:"publicRepos"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ContributionDay is one cell of a year-long contribution heatmap.
type ContributionDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
	Level int    `json:"level"`
}

// ContributionCalendar is a GitHub-style contribution heatmap.
type ContributionCalendar struct {
	Total int               `json:"total"`
	Days  []ContributionDay `json:"days"`
}

// Repository is a simplified git host repository.
type Repository struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"fullName"`
	Description   string    `json:"description"`
	HTMLURL       string    `json:"htmlUrl"`
	CloneURL      string    `json:"cloneUrl"`
	SSHURL        string    `json:"sshUrl"`
	DefaultBranch string    `json:"defaultBranch"`
	Language      string    `json:"language"`
	IsPrivate     bool      `json:"private"`
	IsFork        bool      `json:"fork"`
	IsArchived    bool      `json:"archived"`
	Stars         int       `json:"stars"`
	Forks         int       `json:"forks"`
	OpenIssues    int       `json:"openIssues"`
	OwnerLogin    string    `json:"ownerLogin"`
	OwnerAvatar   string    `json:"ownerAvatar"`
	PushedAt      time.Time `json:"pushedAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// Branch is a simplified git host branch.
type Branch struct {
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	IsProtected bool   `json:"protected"`
}

// Commit is a simplified git host commit.
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	AvatarURL string    `json:"avatarUrl"`
	HTMLURL   string    `json:"htmlUrl"`
	Date      time.Time `json:"date"`
}

// ActivityEvent is a simplified git host activity event.
type ActivityEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     string    `json:"actor"`
	RepoName  string    `json:"repoName"`
	RepoURL   string    `json:"repoUrl"`
	IsPublic  bool      `json:"public"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListOptions controls pagination.
type ListOptions struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

// PageInfo describes pagination metadata.
type PageInfo struct {
	NextPage  int `json:"nextPage"`
	PrevPage  int `json:"prevPage"`
	FirstPage int `json:"firstPage"`
	LastPage  int `json:"lastPage"`
}

// NormalizeListOptions applies default page bounds shared by providers.
func NormalizeListOptions(opts ListOptions) ListOptions {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 30
	}
	if opts.PerPage > 100 {
		opts.PerPage = 100
	}
	return opts
}
