package facebook

import (
	"fmt"
	"strings"

	"github.com/go-social/social"
	"github.com/go-social/social/providers"
)

type Mapper struct{}

func (m *Mapper) BuildPosts(fbPosts []FbPost) social.Posts {
	var posts social.Posts
	for _, fbPost := range fbPosts {
		post := m.BuildPost(fbPost)
		if post != nil {
			posts.Add(post)
		}
	}
	return posts
}

func (m *Mapper) BuildPost(fbPost FbPost) *social.Post {
	var post *social.Post

	idSlice := strings.Split(fbPost.ID, "_")
	if len(idSlice) < 2 {
		return nil
	}
	userID, postID := idSlice[0], idSlice[1]

	post.URL = "https://facebook.com/" + userID + "/posts/" + postID // TODO: need this...? cuz there is fbPost.Link ..?
	post.ID = fbPost.ID
	post.Provider = ProviderID
	post.NumShares = int32(fbPost.Shares.Count)

	post.Author = social.User{
		ID:         fbPost.From.ID,
		Name:       fbPost.From.Name,
		ProfileURL: "https://facebook.com/profile.php?id=" + fbPost.From.ID,
		AvatarURL:  fmt.Sprintf("https://graph.facebook.com/%v/picture?type=large", fbPost.From.ID),
	}

	publishedAt, _ := providers.GetUTCTimeForLayout(fbPost.CreatedTime, timeLayout)
	updatedAt, _ := providers.GetUTCTimeForLayout(fbPost.UpdatedTime, timeLayout)

	post.PublishedAt = &publishedAt
	post.UpdatedAt = &updatedAt

	return post
}
