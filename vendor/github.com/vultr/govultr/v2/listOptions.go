package govultr

// ListOptions
type ListOptions struct {
	PerPage int    `url:"per_page,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
}
