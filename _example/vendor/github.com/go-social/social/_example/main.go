package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-social/social"
	authHandlers "github.com/go-social/social/handlers"
	"github.com/go-social/social/providers"
	_ "github.com/go-social/social/providers/facebook"
	_ "github.com/go-social/social/providers/twitter"
)

func main() {
	// JWT token authentication for service access
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	// Configure social providers
	pcfg := providers.ProviderConfigs{
		"twitter": {
			AppID:         "x",
			AppSecret:     "y",
			OAuthCallback: "http://localhost:1515/auth/twitter/callback",
		},
		"facebook": {
			AppID:         "x",
			AppSecret:     "y",
			OAuthCallback: "http://localhost:1515/auth/facebook/callback",
		},
		"google": {
			AppID:         "x",
			AppSecret:     "y",
			OAuthCallback: "http://localhost:1515/auth/google/callback",
		},
	}
	providers.Configure(pcfg, tokenAuth)

	// HTTP service
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Mount("/auth", authHandlers.Routes(oauthLoginHandler))

	// Start the server on port 0.0.0.0:1515
	http.ListenAndServe(":1515", r)
}

func oauthLoginHandler(r *http.Request, creds []social.Credentials, user *social.User, err error) {
	fmt.Println("oauth login sequence complete")

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("success!")
	fmt.Println("creds:", creds)
	fmt.Println("user:", user)

	cred := creds[0] // pick first one in case there are multiple (ie. fb)

	provider, err := providers.NewSession(context.Background(), cred.ProviderID(), cred)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	profile, err := provider.GetUser(providers.NoQuery)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("provider.GetUser():", profile)
}
