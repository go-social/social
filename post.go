package social

import "time"

// Post is a normalized `post` object across multiple providers. As normalized
// as possible. The original data is available in Raw
type Post struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	URL      string `json:"url"`

	Author    User   `json:"author"`
	Contents  string `json:"contents"`
	NumShares int32  `json:"num_shares"`
	NumLikes  int32  `json:"num_likes"`

	Tags  []string `json:"tags"`
	Links []string `json:"links"`

	PublishedAt *time.Time `json:"published_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`

	Raw interface{}
}

type Posts []*Post

func (ps *Posts) Add(post ...*Post) {
	*ps = append(*ps, post...)
}
