package providers

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultNumResults = 20
	MaxNumResults     = 50
)

var (
	NoQuery = Query{}
)

// TODO: use https://github.com/google/go-querystring with
// `url` struct tags to easier parse/build query strings
type Query struct {
	Search   SearchParts
	Filter   string // Second-pass keywords filter
	Username string // Query by a specific username
	UserID   string

	Limit   int
	Sort    string // recent,popular (default: recent)
	SinceID string // TODO: rename these to NextID and PrevID?
	UntilID string
	Perm    string // read or write, default: read

	Params url.Values
}

func NewQuery(args url.Values) Query {
	q := Query{
		Limit: DefaultNumResults,
		Sort:  "recent",
	}

	if args == nil {
		return q
	}

	q.Search = NewSearchParts(args.Get("q"))
	q.Filter = args.Get("filter")
	q.Username = args.Get("username")
	q.UserID = args.Get("userid")

	if q.Username != "" {
		q.Search.Usernames = append(q.Search.Usernames, q.Username)
	}

	limit, err := strconv.Atoi(args.Get("limit"))
	if err != nil || limit < 1 || limit > MaxNumResults {
		q.Limit = DefaultNumResults
	} else {
		q.Limit = limit
	}

	sort := args.Get("sort")
	if sort == "popular" {
		q.Sort = sort
	}

	q.SinceID = args.Get("since_id")
	q.UntilID = args.Get("until_id")
	q.Perm = args.Get("perm")

	q.Params = args

	return q
}

func (q Query) ToURLArgs() url.Values {
	args := q.Params
	if q.Search.String() != "" {
		args.Set("q", q.Search.String())
	}
	if q.Filter != "" {
		args.Set("filter", q.Filter)
	}
	if q.Username != "" {
		args.Set("username", q.Username)
	}
	if q.UserID != "" {
		args.Set("userid", q.UserID)
	}
	args.Set("limit", strconv.Itoa(q.Limit))

	if q.Sort != "recent" {
		args.Set("sort", q.Sort)
	}

	if q.SinceID != "" {
		args.Set("since_id", q.SinceID)
	}
	if q.UntilID != "" {
		args.Set("until_id", q.UntilID)
	}
	if q.Perm != "" {
		args.Set("perm", q.Perm)
	}
	return args
}

// Search query parts
type SearchParts struct {
	Usernames, Tags, Words []string
}

func NewSearchParts(q string) SearchParts {
	qs := strings.TrimSpace(q)
	qp := SearchParts{}
	parts := strings.Split(qs, " ")

	for _, k := range parts {
		if len(k) == 0 {
			continue
		}
		switch k[0:1] {
		case "@":
			qp.Usernames = append(qp.Usernames, k[1:])
		case "#":
			qp.Tags = append(qp.Tags, k[1:])
		default:
			qp.Words = append(qp.Words, k)
		}
	}
	return qp
}

func (sq SearchParts) String() (s string) {
	if len(sq.Usernames) > 0 {
		s += sq.buildPart(sq.Usernames, "@") + " "
	}
	if len(sq.Tags) > 0 {
		s += sq.buildPart(sq.Tags, "#") + " "
	}
	if len(sq.Words) > 0 {
		s += sq.buildPart(sq.Words, "")
	}
	return strings.TrimSpace(s)
}

// Return only the tags and words in the search query
func (sq SearchParts) Keywords(prefix ...bool) (s string) {
	addPrefix := false
	if len(prefix) > 0 {
		addPrefix = prefix[0]
	}
	if len(sq.Tags) > 0 {
		p := ""
		if addPrefix {
			p = "#"
		}
		s += sq.buildPart(sq.Tags, p) + " "
	}
	if len(sq.Words) > 0 {
		s += sq.buildPart(sq.Words, "")
	}
	return strings.TrimSpace(s)
}

// Returns the first username found
func (sq SearchParts) Username() (s string) {
	if len(sq.Usernames) > 0 {
		s = sq.Usernames[0]
	}
	return
}

func (sq SearchParts) buildPart(parts []string, prefix string) (s string) {
	if len(parts) == 0 {
		return
	}
	for i, k := range parts {
		if prefix != "" {
			s += prefix + k
		} else {
			s += k
		}
		if i < len(parts)-1 {
			s += " "
		}
	}
	return
}
