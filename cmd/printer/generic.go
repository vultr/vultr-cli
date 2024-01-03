package printer

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var _ ResourceOutput = &Generic{}

type Generic struct {
	Message string
}

func (g *Generic) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return prettyJSON
}

func (g *Generic) Yaml() []byte {
	yam, err := yaml.Marshal(g)
	if err != nil {
		panic(err.Error())
	}
	return yam
}

func (g *Generic) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"message"}}
}

func (g *Generic) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {g.Message}}
}

func (g *Generic) Paging() map[int][]interface{} {
	return nil
}
