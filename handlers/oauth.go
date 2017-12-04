package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

func ProviderCtx(oauthErrorFn ErrorHandlerFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			providerID := strings.ToLower(chi.URLParam(r, "provider"))

			oauth, err := providers.NewOAuth(providerID)
			if err != nil {
				oauthErrorFn(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), ProviderIDCtxKey, providerID)
			ctx = context.WithValue(ctx, ProviderOAuthCtxKey, oauth)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OAuth(oauthErrorFn ErrorHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oauth := ctx.Value(ProviderOAuthCtxKey).(social.OAuth)

		state := jwtauth.Claims{
			"sub":      "OAuthCallback",
			"provider": oauth.ProviderID(),
			"perm":     social.PermissionFromString(r.URL.Query().Get("perm")).String(),
		}

		// Give users 15 minutes to authenticate
		state.SetIssuedNow()
		state.SetExpiryIn(15 * time.Minute)

		authURL, err := oauth.AuthCodeURL(r, state)
		if err != nil {
			oauthErrorFn(w, r, err)
			return
		}

		http.Redirect(w, r, authURL, 302)
	}
}

func OAuthCallback(oauthCallbackFn CallbackHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var mcreds []social.Credentials
		var providerUser *social.User

		ctx := r.Context()
		oauth := ctx.Value(ProviderOAuthCtxKey).(social.OAuth)

		defer func() {
			oauthCallbackFn(w, r, mcreds, providerUser, err)
		}()

		mcreds, err = oauth.Exchange(ctx, r)
		if err != nil {
			return
		}
		creds := mcreds[0]

		p, err := providers.NewSession(ctx, oauth.ProviderID(), creds)
		if err != nil {
			return
		}

		providerUser, err = p.GetUser(providers.NoQuery)
		if err != nil {
			return
		}
	}
}
