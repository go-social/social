package twitter

import (
	"context"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

type OAuth struct {
	client oauth.Client
}

func NewOAuth() social.OAuth {
	client := oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}
	client.Credentials.Token = AppID
	client.Credentials.Secret = AppSecret
	return &OAuth{client}
}

func (oa *OAuth) ProviderID() string {
	return ProviderID
}

func (oa *OAuth) AuthCodeURL(r *http.Request, claims map[string]interface{}) (string, error) {
	// NOTE: Twitter API permissions are set per application, not per token.
	// Default to read only permission
	if _, ok := claims["perm"]; !ok {
		claims["perm"] = social.PermissionRead.String()
	}

	_, stateToken, err := providers.TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}
	callbackURL := OAuthCallback + "?state=" + stateToken

	tempCred, err := oa.client.RequestTemporaryCredentials(http.DefaultClient, callbackURL, nil)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	if _, ok := claims["force_login"]; ok {
		v.Set("force_login", "true")
		v.Set("screen_name", "")
	}

	authURL := oa.client.AuthorizationURL(tempCred, v)
	return authURL, nil
}

func (oa *OAuth) Exchange(ctx context.Context, r *http.Request) ([]social.Credentials, error) {
	callbackArgs := r.URL.Query()
	reqToken := callbackArgs.Get("oauth_token")
	verifier := callbackArgs.Get("oauth_verifier")
	tempCred := &oauth.Credentials{reqToken, ""}

	access, _, err := oa.client.RequestToken(http.DefaultClient, tempCred, verifier)
	if err != nil {
		return nil, providerError(err)
	}

	// NOTE: twitter does not expire oauth tokens:
	// https://dev.twitter.com/oauth/overview/faq
	creds := []social.Credentials{
		&providers.OAuth1Creds{
			AuthProviderID:        ProviderID,
			AuthAccessToken:       access.Token,
			AuthAccessTokenSecret: access.Secret,
		},
	}
	return creds, nil
}
