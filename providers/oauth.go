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
	AuthProviderID        string
	AuthAccessToken       string
	AuthAccessTokenSecret string
	AuthRefreshToken      string
	AuthExpiresAt         *time.Time
	AuthPermission        social.Permission
}

var _ social.Credentials = &OAuth1Creds{}

func (c *OAuth1Creds) ProviderID() string {
	return c.AuthProviderID
}

func (c *OAuth1Creds) AccessToken() string {
	return c.AuthAccessToken
}

func (c *OAuth1Creds) AccessTokenSecret() string {
	return c.AuthAccessTokenSecret
}

func (c *OAuth1Creds) RefreshToken() string {
	return c.AuthRefreshToken
}

func (c *OAuth1Creds) Permission() social.Permission {
	return c.AuthPermission
}

func (c *OAuth1Creds) SetPermission(perm string) {
	c.AuthPermission = social.PermissionFromString(perm)
}

func (c *OAuth1Creds) ExpiresAt() *time.Time {
	return c.AuthExpiresAt
}

// OAuth2Creds is a normalized social.Credentials implementation to
// use with social providers using oauth2 (ie. facebook, google, ..)
type OAuth2Creds struct {
	*oauth2.Token
	AuthProviderID string
	AuthPermission social.Permission
}

var _ social.Credentials = &OAuth2Creds{}

func (c *OAuth2Creds) ProviderID() string {
	return c.AuthProviderID
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
	return c.AuthPermission
}

func (c *OAuth2Creds) SetPermission(perm string) {
	c.AuthPermission = social.PermissionFromString(perm)
}

func (c *OAuth2Creds) ExpiresAt() *time.Time {
	return &c.Token.Expiry
}
