package facebook

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-social/social"
	"github.com/go-social/social/providers"
	fb "github.com/huandu/facebook"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const (
	ProviderID = `facebook`
)

var (
	postFields = strings.Join([]string{
		"actions",
		"application",
		"attachments",
		"caption",
		"created_time",
		"description",
		"from",
		"full_picture",
		"icon",
		"id",
		"likes",
		"link",
		"message",
		"message_tags",
		"name",
		"object_id",
		"place",
		"privacy",
		"properties",
		"shares",
		"source",
		"status_type",
		"type",
		"updated_time",
	}, ",")

	basicFields = strings.Join([]string{
		"about",
		"id",
		"link",
		"location",
		"name",
		"picture.type(large)",
	}, ",")

	userFields = strings.Join([]string{
		basicFields,
		"email",
		"timezone",
	}, ",")
)

const timeLayout = `2006-01-02T15:04:05-0700`

var (
	AppID         string
	AppSecret     string
	OAuthCallback string

	// Facebook API version
	// See: https://developers.facebook.com/docs/apps/changelog for updates
	FacebookApiVersion = "v2.11"
)

type Provider struct {
	creds social.Credentials
	api   *fb.Session
}

func New(ctx context.Context, creds social.Credentials) (providers.ProviderSession, error) {
	conf := &oauth2.Config{
		ClientID:     AppID,
		ClientSecret: AppSecret,
		RedirectURL:  OAuthCallback,
		Endpoint:     facebook.Endpoint,
	}

	token := &oauth2.Token{
		AccessToken:  creds.AccessToken(),
		RefreshToken: creds.RefreshToken(),
	}
	if expiresAt := creds.ExpiresAt(); expiresAt != nil {
		token.Expiry = *expiresAt
	}

	api := &fb.Session{
		Version:    FacebookApiVersion,
		HttpClient: conf.Client(ctx, token),
	}

	if err := api.Validate(); err != nil {
		return nil, providerError(err)
	}

	return &Provider{api: api, creds: creds}, nil
}

func (p *Provider) ID() string {
	return ProviderID
}

func (p *Provider) Post(ctx context.Context, msg string, shareLink string) (*social.Post, error) {
	return nil, providers.ErrNotImplemented
}

func (p *Provider) Search(query providers.Query) (social.Posts, *providers.Cursor, error) {
	return nil, nil, providers.ErrUnsupported
}

func (p *Provider) GetFeed(query providers.Query) (social.Posts, *providers.Cursor, error) {
	return nil, nil, providers.ErrNotImplemented
}

func (p *Provider) GetPosts(query providers.Query) (social.Posts, *providers.Cursor, error) {
	return nil, nil, providers.ErrNotImplemented
}

func (p *Provider) GetUser(query providers.Query) (*social.User, error) {
	var resp fb.Result
	var err error

	username := query.Username
	if username == "" {
		username = "me"
	}

	args := url.Values{}
	args.Set("fields", userFields)
	resp, err = p.api.Get(username, getFbParams(args))
	if err != nil && providerError(err) == errNodeTypePage {
		// Try again, this time this will query for basic information only.
		args := url.Values{}
		args.Set("fields", basicFields)
		resp, err = p.api.Get(username, getFbParams(args))
	}
	if err != nil {
		return nil, providerError(err)
	}

	fbResponse, err := getActualResponseAndError(resp, err)
	if err != nil {
		return nil, providerError(err)
	}

	user := &social.User{
		Provider:     ProviderID,
		ID:           fbResponse.ID,
		Name:         fbResponse.Name,
		Username:     fbResponse.Name,
		ProfileURL:   fbResponse.Link,
		AvatarURL:    fbResponse.Picture.Data.URL,
		Location:     fbResponse.Location.Name,
		Email:        fbResponse.Email,
		Timezone:     strconv.FormatFloat(float64(fbResponse.Timezone), 'f', 2, 32),
		NumFollowers: int32(fbResponse.Friends.Summary.TotalCount),
	}

	return user, nil
}

func (p *Provider) GetFriends(query providers.Query) ([]*social.User, *providers.Cursor, error) {
	return nil, nil, providers.ErrNotImplemented
}

func (p *Provider) GetFollowers(query providers.Query) ([]*social.User, *providers.Cursor, error) {
	return nil, nil, providers.ErrNotImplemented
}

func getFbParams(args url.Values) fb.Params {
	params := fb.Params{}
	for key := range args {
		params[key] = args.Get(key)
	}
	return params
}

func getActualResponseAndError(res map[string]interface{}, err error) (FbResponse, error) {
	getFbResponse := func() FbResponse {
		var fbResponse FbResponse
		b, _ := json.Marshal(res)
		json.Unmarshal(b, &fbResponse)
		return fbResponse
	}

	if res != nil && len(res) > 0 {
		// check to see if error is present in the response
		fbResponse := getFbResponse()
		if fbResponse.FbError.Message != "" || fbResponse.FbError.Type != "" {
			// This means that fb returned an error code
			return FbResponse{}, &(fbResponse.FbError)
		} else {
			return fbResponse, nil
		}
	}
	if err != nil {
		// This means fb didnt return an error but some other error happened
		return FbResponse{}, err
	}
	fbResponse := getFbResponse()
	return fbResponse, nil
}
