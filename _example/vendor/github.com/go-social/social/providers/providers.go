package providers

import (
	"context"

	"github.com/go-social/social"
)

type Provider struct {
	Configure func(appID string, appSecret string, oauthCallback string)
	New       func(ctx context.Context, creds social.Credentials) (ProviderSession, error)
	NewOAuth  func() social.OAuth
}

type ProviderSession interface {
	// ID of the Provider
	ID() string

	// Post a message to the provider and return the new Post object created.
	Post(ctx context.Context, msg string, link string) (*social.Post, error)

	// Search content on a provider network
	Search(query Query) (social.Posts, *Cursor, error)

	// Get a user's feed/wall
	GetFeed(query Query) (social.Posts, *Cursor, error) // Feed

	// Get a user's own posts
	GetPosts(query Query) (social.Posts, *Cursor, error) // Posts

	// Get the user social profile object
	GetUser(query Query) (*social.User, error)

	// Get a user's friends list (aka following)
	GetFriends(query Query) ([]*social.User, *Cursor, error)

	// Get a user's followers list
	GetFollowers(query Query) ([]*social.User, *Cursor, error)
}

func NewSession(ctx context.Context, providerID string, creds social.Credentials) (ProviderSession, error) {
	r, ok := Registry[providerID]
	if !ok {
		return nil, ErrUnknownProviderID
	}
	return r.New(ctx, creds)
}

var Registry = make(map[string]*Provider)

func Register(providerID string, provider *Provider) {
	Registry[providerID] = provider
}
