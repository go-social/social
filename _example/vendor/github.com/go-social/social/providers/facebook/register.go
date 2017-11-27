package facebook

import "github.com/go-social/social/providers"

func Configure(appID string, appSecret string, oauthCallback string) {
	AppID = appID
	AppSecret = appSecret
	OAuthCallback = oauthCallback
}

func init() {
	providers.Register(ProviderID, &providers.Provider{
		Configure: Configure,
		New:       New,
		NewOAuth:  NewOAuth,
	})
}
