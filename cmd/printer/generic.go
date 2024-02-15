package printer

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var _ ResourceOutput = &Generic{}

// Generic ...
type Generic struct {
	Message string `json:"message"`
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
func (g *Generic) Columns() [][]string {
	return [][]string{0: {"MESSAGE"}}
}

// Data ...
func (g *Generic) Data() [][]string {
	return [][]string{0: {g.Message}}
}

// Paging ...
func (g *Generic) Paging() [][]string {
	return nil
}
