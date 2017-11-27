package twitter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

const (
	ProviderID = `twitter`
	TimeLayout = `Mon Jan 02 15:04:05 -0700 2006`
)

var (
	AppID         string
	AppSecret     string
	OAuthCallback string
)

type Provider struct {
	creds social.Credentials
	api   *anaconda.TwitterApi
}

func New(ctx context.Context, creds social.Credentials) (providers.ProviderSession, error) {
	api := anaconda.NewTwitterApi(creds.AccessToken(), creds.AccessTokenSecret())
	api.ReturnRateLimitError(true)
	api.DisableThrottling()
	api.HttpClient = http.DefaultClient
	return &Provider{creds: creds, api: api}, nil
}

func Configure(appID string, appSecret string, oauthCallback string) {
	AppID = appID
	AppSecret = appSecret
	OAuthCallback = oauthCallback

	anaconda.SetConsumerKey(appID)
	anaconda.SetConsumerSecret(appSecret)
}

func init() {
	providers.Register(ProviderID, &providers.Provider{
		Configure: Configure,
		New:       New,
		NewOAuth:  NewOAuth,
	})
}

func (p *Provider) ID() string {
	return ProviderID
}

// Post a tweet to twitter
func (p *Provider) Post(ctx context.Context, msg string, shareLink string) (*social.Post, error) {
	// Append the share link to the message
	if shareLink != "" && strings.Index(msg, shareLink) < 0 {
		msg = fmt.Sprintf("%s %s", strings.TrimSpace(msg), shareLink)
	}

	// Send tweet
	tweet, err := p.api.PostTweet(msg, url.Values{})
	if err != nil {
		perr := providerError(err)
		return nil, perr
	}

	newPost := &social.Post{
		Provider: p.ID(),
		ID:       tweet.IdStr,
		URL:      fmt.Sprintf("https://twitter.com/statuses/%s", tweet.IdStr),
	}

	return newPost, nil
}

// Search Twitter via their REST API
func (p *Provider) Search(query providers.Query) (social.Posts, *providers.Cursor, error) {
	var tweets []anaconda.Tweet
	var err error

	args := url.Values{}
	args.Add("count", strconv.Itoa(query.Limit))
	args.Add("include_entities", "true")

	// remove retweets from a user's timeline
	args.Set("include_rts", "false")

	// remove replies from a user's timeline
	args.Set("exclude_replies", "true")

	if query.Sort == "popular" {
		args.Add("result_type", "mixed")
	} else {
		// note, twitter supports: mixed, recent, popular
		args.Add("result_type", "recent")
	}

	if query.UntilID != "" {
		args.Add("max_id", query.UntilID)
	}
	if query.SinceID != "" {
		args.Add("since_id", query.SinceID)
	}

	if query.Search.Username() != "" {
		// See: https://dev.twitter.com/rest/reference/get/statuses/user_timeline
		args.Set("screen_name", query.Search.Username())

		// turn off trim_user, give us the entire object
		args.Set("trim_user", "false")

		// Query twitter's rest api
		tweets, err = p.api.GetUserTimeline(args)

	} else {
		q := query.Search.Keywords(true)

		var resp anaconda.SearchResponse
		resp, err = p.api.GetSearch(q, args)
		tweets = resp.Statuses
	}

	if err != nil {
		perr := providerError(err)
		return nil, nil, perr
	}

	posts := (Mapper{}).BuildPosts(tweets)
	prev, next := getCursorIDs(posts)
	cursor := providers.NewCursor(query, prev, next)

	return posts, cursor, nil
}

func (p *Provider) GetFeed(query providers.Query) (social.Posts, *providers.Cursor, error) {
	var tweets []anaconda.Tweet
	args := url.Values{}

	args.Add("count", strconv.Itoa(query.Limit))
	if query.UntilID != "" {
		args.Add("max_id", query.UntilID)
	}
	if query.SinceID != "" {
		args.Add("since_id", query.SinceID)
	}

	tweets, err := p.api.GetHomeTimeline(args)
	if err != nil {
		perr := providerError(err)
		return nil, nil, perr
	}

	posts := (Mapper{}).BuildPosts(tweets)
	prev, next := getCursorIDs(posts)
	cursor := providers.NewCursor(query, prev, next)

	return posts, cursor, providerError(err)
}

func (p *Provider) GetPosts(query providers.Query) (social.Posts, *providers.Cursor, error) {
	var tweets []anaconda.Tweet
	args := url.Values{}

	args.Add("count", strconv.Itoa(query.Limit))
	if query.UntilID != "" {
		args.Add("max_id", query.UntilID)
	}
	if query.SinceID != "" {
		args.Add("since_id", query.SinceID)
	}

	tweets, err := p.api.GetUserTimeline(args)
	if err != nil {
		perr := providerError(err)
		return nil, nil, perr
	}

	posts := (Mapper{}).BuildPosts(tweets)
	prev, next := getCursorIDs(posts)
	cursor := providers.NewCursor(query, prev, next)

	return posts, cursor, providerError(err)
}

func (p *Provider) GetUser(query providers.Query) (*social.User, error) {
	var u anaconda.User
	var err error

	if query.UserID != "" {
		userid, _ := strconv.Atoi(query.UserID)
		u, err = p.api.GetUsersShowById(int64(userid), nil)
	} else if query.Username == "" {
		u, err = p.api.GetSelf(nil)
	} else {
		u, err = p.api.GetUsersShow(query.Username, nil)
		// TODO: when twitter returns 404, then the user doesn't exist,
		// therefore we should respond with a user not found error
	}

	if err != nil {
		perr := providerError(err)
		return nil, perr
	}

	user := (UserMapper{}).BuildUser(u)
	return user, nil
}

// Get a user's friends (aka following)
// Network docs: https://dev.twitter.com/rest/reference/get/friends/list
func (p *Provider) GetFriends(query providers.Query) ([]*social.User, *providers.Cursor, error) {
	v := url.Values{}
	v.Add("count", strconv.Itoa(query.Limit))
	v.Add("include_user_entities", "true")

	if query.UserID != "" {
		v.Add("user_id", query.UserID)
	} else if query.Username != "" {
		v.Add("screen_name", query.Username)
	} else {
		return nil, nil, errors.New("social: UserID or Username not specified in query")
	}
	if query.UntilID != "" {
		v.Set("cursor", query.UntilID)
	}
	if query.SinceID != "" {
		v.Set("cursor", query.SinceID)
	}

	q, err := p.api.GetFriendsList(v)
	if err != nil {
		perr := providerError(err)
		return nil, nil, perr
	}

	users := (UserMapper{}).BuildUsers(q.Users)
	cursor := providers.NewCursor(query, q.Previous_cursor_str, q.Next_cursor_str)

	return users, cursor, nil
}

// Get a user's followers
// Network docs: https://dev.twitter.com/rest/reference/get/followers/list
func (p *Provider) GetFollowers(query providers.Query) ([]*social.User, *providers.Cursor, error) {
	v := url.Values{}
	v.Add("count", strconv.Itoa(query.Limit))
	v.Add("include_user_entities", "true")

	if query.UserID != "" {
		v.Add("user_id", query.UserID)
	} else if query.Username != "" {
		v.Add("screen_name", query.Username)
	} else {
		return nil, nil, errors.New("UserID or Username not specified in query")
	}
	if query.UntilID != "" {
		v.Set("cursor", query.UntilID)
	}
	if query.SinceID != "" {
		v.Set("cursor", query.SinceID)
	}

	q, err := p.api.GetFollowersList(v)
	if err != nil {
		perr := providerError(err)
		return nil, nil, perr
	}

	users := (UserMapper{}).BuildUsers(q.Users)
	cursor := providers.NewCursor(query, q.Previous_cursor_str, q.Next_cursor_str)

	return users, cursor, nil
}

func getCursorIDs(posts social.Posts) (prevID, nextID string) {
	if len(posts) == 0 {
		return
	}
	prevID = posts[0].ID
	nextID = posts[len(posts)-1].ID
	return
}
