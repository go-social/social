package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

func ProviderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providerID := chi.URLParam(r, "provider")

		oauth, err := providers.NewOAuth(providerID)
		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, err) // TODO: use payloads..
			// TODO: if render.SetContentType to JSON..
			// we need to be able to override the content-type
			// on a per response or router basis.. check that out.
			return
		}

		ctx := context.WithValue(r.Context(), ProviderOAuthCtxKey, oauth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OAuth(w http.ResponseWriter, r *http.Request) {
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
		render.Status(r, 401)
		render.JSON(w, r, err) // TODO: render.Render()
		return
	}

	http.Redirect(w, r, authURL, 302)
}

func OAuthCallback(oauthLoginFn OAuthLoginFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var mcreds []social.Credentials
		var providerUser *social.User

		ctx := r.Context()
		oauth := ctx.Value(ProviderOAuthCtxKey).(social.OAuth)

		defer func() {
			oauthLoginFn(w, r, mcreds, providerUser, err)
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
