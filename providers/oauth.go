package providers

import (
	"time"

	"github.com/go-social/social"
	"golang.org/x/oauth2"
)

func NewOAuth(providerID string) (social.OAuth, error) {
	r, ok := Registry[providerID]
	if !ok {
		return nil, ErrUnknownProviderID
	}
	return r.NewOAuth(), nil
}

// OAuth1Creds is a normalized social.Credentials implementation to
// use with social providers using oauth1 (ie. twitter)
type OAuth1Creds struct {
	CredProviderID        string
	CredProviderUserID    string
	CredAccessToken       string
	CredAccessTokenSecret string
	CredRefreshToken      string
	CredExpiresAt         *time.Time
	CredPermission        social.Permission
}

var _ social.Credentials = &OAuth1Creds{}

func (c *OAuth1Creds) ProviderID() string {
	return c.CredProviderID
}

func (c *OAuth1Creds) ProviderUserID() string {
	return c.CredProviderUserID
}

func (c *OAuth1Creds) AccessToken() string {
	return c.CredAccessToken
}

func (c *OAuth1Creds) AccessTokenSecret() string {
	return c.CredAccessTokenSecret
}

func (c *OAuth1Creds) RefreshToken() string {
	return c.CredRefreshToken
}

func (c *OAuth1Creds) Permission() social.Permission {
	return c.CredPermission
}

func (c *OAuth1Creds) SetPermission(perm string) {
	c.CredPermission = social.PermissionFromString(perm)
}

func (c *OAuth1Creds) ExpiresAt() *time.Time {
	return c.CredExpiresAt
}

// OAuth2Creds is a normalized social.Credentials implementation to
// use with social providers using oauth2 (ie. facebook, google, ..)
type OAuth2Creds struct {
	*oauth2.Token
	CredProviderID     string
	CredProviderUserID string
	CredPermission     social.Permission
}

var _ social.Credentials = &OAuth2Creds{}

func (c *OAuth2Creds) ProviderID() string {
	return c.CredProviderID
}

func (c *OAuth2Creds) ProviderUserID() string {
	return c.CredProviderUserID
}

func (c *OAuth2Creds) AccessToken() string {
	return c.Token.AccessToken
}

func (c *OAuth2Creds) AccessTokenSecret() string {
	return ""
}

func (c *OAuth2Creds) RefreshToken() string {
	return c.Token.RefreshToken
}

func (c *OAuth2Creds) Permission() social.Permission {
	return c.CredPermission
}

func (c *OAuth2Creds) SetPermission(perm string) {
	c.CredPermission = social.PermissionFromString(perm)
}

func (c *OAuth2Creds) ExpiresAt() *time.Time {
	return &c.Token.Expiry
}
