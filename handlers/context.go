package handlers

var (
	ProviderIDCtxKey    = &contextKey{"ProviderID"}
	ProviderOAuthCtxKey = &contextKey{"ProviderOAuth"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "context value " + k.name
}
