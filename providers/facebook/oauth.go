package facebook

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-social/social"
	"github.com/go-social/social/providers"
	fb "github.com/huandu/facebook"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	// https://developers.facebook.com/docs/facebook-login/permissions/v2.11
	loginScope = []string{
		"public_profile",
		"email",
		"user_location",
	}
	readScope = []string{
		"user_friends",
		"user_posts",
		"user_status",
		"user_likes",
		"user_photos",
		"user_videos",
	}
	writeScope = []string{
		"manage_pages",
		"publish_pages",
		"publish_actions",
	}
)

type OAuth struct {
	*oauth2.Config
}

func NewOAuth() social.OAuth {
	config := &oauth2.Config{
		ClientID:     AppID,
		ClientSecret: AppSecret,
		RedirectURL:  OAuthCallback,
		Endpoint:     facebook.Endpoint,
	}

	return &OAuth{config}
}

func (oa *OAuth) ProviderID() string {
	return ProviderID
}

func (oa *OAuth) AuthCodeURL(r *http.Request, claims map[string]interface{}) (string, error) {
	scope := loginScope

	perm, _ := claims["perm"].(string)
	switch social.PermissionFromString(perm) {
	case social.PermissionRead:
		scope = append(scope, readScope...)
	case social.PermissionWrite:
		scope = append(scope, writeScope...)
	case social.PermissionReadWrite:
		scope = append(scope, append(readScope, writeScope...)...)
	}

	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("scope", strings.Join(scope, ",")),
	}
	if _, ok := claims["force_login"]; ok {
		opts = append(opts, oauth2.SetAuthURLParam("auth_type", "reauthenticate"))
	}

	_, stateToken, err := providers.TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return oa.Config.AuthCodeURL(stateToken, opts...), nil
}

func (oa *OAuth) Exchange(ctx context.Context, r *http.Request) ([]social.Credentials, error) {
	callbackArgs := r.URL.Query()
	code := callbackArgs.Get("code")
	cbError := callbackArgs.Get("error")
	cbErrorMsg := callbackArgs.Get("error_reason")
	cbErrorDesc := callbackArgs.Get("error_description")

	if cbError != "" {
		msg := fmt.Sprintf("Error:%v,  ErrorReason:%v,  ErrorDescription:%v", cbError, cbErrorMsg, cbErrorDesc)
		// TODO should probably use a different error code
		return nil, providers.ErrAuthFailed.Err(errors.New(msg))
	}
	if code == "" {
		msg := "empty code in facebook callback"
		return nil, providers.ErrAuthFailed.Err(errors.New(msg))
	}

	userToken, err := oa.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	creds := []social.Credentials{
		&providers.OAuth2Creds{Token: userToken},
	}

	client := oa.Config.Client(ctx, userToken)
	api := &fb.Session{
		Version:    FacebookApiVersion,
		HttpClient: client,
	}

	// Fetch tokens for FB pages.
	resp, err := api.Get("/me/accounts", getFbParams(url.Values{}))
	if err != nil {
		return creds, nil
	}

	var accounts FbResponseAccounts
	if err := resp.Decode(&accounts); err != nil {
		return creds, nil
	}

	for _, account := range accounts.Data {
		// Re-use user's TokenType, RefreshToken and Expiry.
		pageToken := &oauth2.Token{
			AccessToken:  account.AccessToken,
			TokenType:    userToken.TokenType,
			RefreshToken: userToken.RefreshToken,
			Expiry:       userToken.Expiry,
		}

		creds = append(creds, &providers.OAuth2Creds{Token: pageToken})
	}

	return creds, nil
}
