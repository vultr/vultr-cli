package printer

// Message provides generic output for the CLI
type Message struct {
	Message string `json:"message"`
}

// JSON ...
func (m *Message) JSON() []byte {
	return MarshalObject(m, "json")
}

// YAML ...
func (m *Message) YAML() []byte {
	return MarshalObject(m, "yaml")
}

// Columns ...
func (m *Message) Columns() [][]string {
	return [][]string{0: {"MESSAGE"}}
}

// Data ...
func (m *Message) Data() [][]string {
	return [][]string{0: {m.Message}}
}

// Paging ...
func (m *Message) Paging() [][]string {
	return nil
}
