package twitter

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

type Mapper struct{}

func (m Mapper) BuildPosts(tweets []anaconda.Tweet) social.Posts {
	var posts social.Posts
	for _, tweet := range tweets {
		post := m.BuildPost(tweet)
		if post != nil {
			posts.Add(post)
		}
	}
	return posts
}

func (m Mapper) BuildPost(tweet anaconda.Tweet) *social.Post {
	postURL := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IdStr)

	post := &social.Post{
		Raw:      tweet,
		ID:       tweet.IdStr,
		Provider: ProviderID,
		URL:      postURL,
		Author: social.User{
			ID:           tweet.User.IdStr,
			Name:         tweet.User.Name,
			Username:     tweet.User.ScreenName,
			ProfileURL:   "https://twitter.com/" + tweet.User.ScreenName,
			AvatarURL:    getLargeProfileImageURL(tweet.User.ProfileImageURL),
			NumFollowers: int32(tweet.User.FollowersCount),
			NumFollowing: int32(tweet.User.FriendsCount),
		},
		Contents:  tweet.Text,
		NumShares: int32(tweet.RetweetCount),
		NumLikes:  int32(tweet.FavoriteCount),
	}

	publishedAt, _ := providers.GetUTCTimeForLayout(tweet.CreatedAt, TimeLayout)
	post.PublishedAt = &publishedAt

	return post
}

type UserMapper struct{}

func (m UserMapper) BuildUsers(us []anaconda.User) []*social.User {
	var users []*social.User
	for _, u := range us {
		users = append(users, m.BuildUser(u))
	}
	return users
}

func (m UserMapper) BuildUser(u anaconda.User) *social.User {
	return &social.User{
		Provider:     ProviderID,
		ID:           u.IdStr,
		Username:     u.ScreenName,
		Name:         u.Name,
		ProfileURL:   fmt.Sprintf("https://twitter.com/%s", u.ScreenName),
		AvatarURL:    getLargeProfileImageURL(u.ProfileImageURL),
		NumPosts:     int32(u.StatusesCount),
		NumFollowers: int32(u.FollowersCount),
		NumFollowing: int32(u.FriendsCount),
		Lang:         u.Lang,
		Location:     u.Location,
		Timezone:     u.TimeZone,
		Private:      u.Protected,
	}
}

func getLargeProfileImageURL(url string) string {
	parts := strings.Split(url, "?")
	ext := path.Ext(parts[0])
	r, _ := regexp.Compile("(?i)_normal" + ext)
	return r.ReplaceAllString(url, ext)
}
