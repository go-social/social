package social

import (
	"context"
	"net/http"
	"time"
)

type OAuth interface {
	ProviderID() string

	// Authorization URL to start OAuth process and retrieve a request token
	AuthCodeURL(r *http.Request, claims map[string]interface{}) (string, error)

	// Request an access token from a request token to complete authentication
	Exchange(ctx context.Context, r *http.Request) ([]Credentials, error)
}

type Credentials interface {
	ProviderID() string
	ProviderUserID() string
	AccessToken() string
	AccessTokenSecret() string
	RefreshToken() string
	ExpiresAt() *time.Time
	Permission() Permission
	SetPermission(string)
}

type User struct {
	Provider     string     `json:"provider" url:"provider"`
	ID           string     `json:"id" url:"id"`
	Username     string     `json:"username" url:"username"`
	Name         string     `json:"name" url:"name"`
	Email        string     `json:"email" url:"email"`
	ProfileURL   string     `json:"profile_url" url:"profile_url"`
	AvatarURL    string     `json:"avatar_url" url:"avatar_url"`
	NumPosts     int32      `json:"num_posts" url:"num_posts"`
	NumFollowers int32      `json:"num_followers" url:"num_folowers"`
	NumFollowing int32      `json:"num_following" url:"num_following"`
	Lang         string     `json:"lang" url:"lang"`
	Location     string     `json:"location" url:"location"`
	Timezone     string     `json:"timezone" url:"timezone"`
	Private      bool       `json:"private" url:"private"`
	LastSyncAt   *time.Time `json:"last_sync_at" url:"last_sync_at"`
}
