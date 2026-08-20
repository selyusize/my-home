package github

import (
	"errors"
	"time"
)

var (
	ErrNotAuthenticated = errors.New("github: not authenticated")
	ErrMissingOAuthConfig = errors.New("github: oauth client id/secret are required")
)

// Config holds OAuth application settings.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// User is a simplified GitHub profile.
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	Blog      string `json:"blog"`
	AvatarURL string `json:"avatarUrl"`
	HTMLURL   string `json:"htmlUrl"`
	PublicRepos int `json:"publicRepos"`
	Followers     int `json:"followers"`
	Following     int `json:"following"`
	CreatedAt     time.Time `json:"createdAt"`
}

// Repository is a simplified GitHub repository.
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
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	Archived      bool      `json:"archived"`
	Stars         int       `json:"stars"`
	Forks         int       `json:"forks"`
	OpenIssues    int       `json:"openIssues"`
	OwnerLogin    string    `json:"ownerLogin"`
	OwnerAvatar   string    `json:"ownerAvatar"`
	PushedAt      time.Time `json:"pushedAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// ActivityEvent is a simplified GitHub activity event.
type ActivityEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     string    `json:"actor"`
	RepoName  string    `json:"repoName"`
	RepoURL   string    `json:"repoUrl"`
	Public    bool      `json:"public"`
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
