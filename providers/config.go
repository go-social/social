package providers

import (
	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

type ProviderConfig struct {
	AppID         string `toml:"app_id"`
	AppSecret     string `toml:"app_secret"`
	OAuthCallback string `toml:"oauth_callback"`
}

type ProviderConfigs map[string]ProviderConfig

func Configure(confs ProviderConfigs, tokenAuth *jwtauth.JWTAuth) {
	for id, conf := range confs {
		if p, ok := Registry[id]; ok {
			p.Configure(conf.AppID, conf.AppSecret, conf.OAuthCallback)
		}
	}
	TokenAuth = tokenAuth
}
