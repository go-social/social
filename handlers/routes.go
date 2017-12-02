package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

type OAuthLoginFunc func(w http.ResponseWriter, r *http.Request, creds []social.Credentials, user *social.User, err error)

func Routes(oauthLoginFn OAuthLoginFunc) http.Handler {
	r := chi.NewRouter()

	r.Get("/", ListProviders)

	r.Route("/{provider}", func(r chi.Router) {
		r.Use(ProviderCtx)

		r.Get("/", OAuth) // open

		r.Group(func(r chi.Router) {
			// secure, via jwt state token
			r.Use(jwtauth.Verify(providers.TokenAuth, tokenFromQuery("state")))
			r.Use(jwtauth.Authenticator)
			r.Get("/callback", OAuthCallback(oauthLoginFn))
		})
	})

	// TODO: this needs to be secured as well, as
	// its a callback router
	r.Get("/loopback/{route}", Loopback)

	return r
}

func ListProviders(w http.ResponseWriter, r *http.Request) {
	plist := []string{}
	for id, _ := range providers.Registry {
		plist = append(plist, id)
	}

	// TODO: make payloads for all request / response objects, and render.Render()
	render.JSON(w, r, plist)

}

// Loopback redirects the client to another path on our router
func Loopback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)

	route := chi.URLParam(r, "route")
	providerID, ok := claims["provider"]

	if !ok {
		render.Status(r, 403) // TODO: defined payload..
		render.JSON(w, r, "invalid provider id")
		return
	}

	switch route {
	case "googleapi":
		redirectURL := fmt.Sprintf("/auth/%s/callback?%s", providerID, r.URL.Query().Encode())
		http.Redirect(w, r, redirectURL, 302)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func tokenFromQuery(param string) func(r *http.Request) string {
	// Get token from query param
	return func(r *http.Request) string {
		return r.URL.Query().Get(param)
	}
}
