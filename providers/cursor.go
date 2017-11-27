package providers

// Query cursor
type Cursor struct {
	Next *Query
	Prev *Query
}

func NewCursor(query Query, prevID string, nextID string) *Cursor {
	c := &Cursor{}
	p := query
	n := query

	c.Prev = &p
	c.Prev.SinceID = prevID
	c.Prev.UntilID = ""

	c.Next = &n
	c.Next.SinceID = ""
	c.Next.UntilID = nextID

	return c
}
