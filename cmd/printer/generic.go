package printer

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var _ ResourceOutput = &Generic{}

// Generic ...
type Generic struct {
	Message string
}

// JSON ...
func (g *Generic) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return prettyJSON
}

// YAML ...
func (g *Generic) YAML() []byte {
	yam, err := yaml.Marshal(g)
	if err != nil {
		panic(err.Error())
	}
	return yam
}

// Columns ...
func (g *Generic) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"message"}}
}

// Data ...
func (g *Generic) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {g.Message}}
}

// Paging ...
func (g *Generic) Paging() map[int][]interface{} {
	return nil
}
