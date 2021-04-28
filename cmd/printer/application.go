package printer

import (
	"encoding/json"

	"github.com/go-yaml/yaml"
	"github.com/vultr/govultr/v2"
)

var _ ResourceOutput = &Applications{}

type Applications struct {
	Applications []govultr.Application `json:"applications"`
	Meta         *govultr.Meta
}

func (a *Applications) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (a *Applications) Yaml() []byte {
	yam, err := yaml.Marshal(a)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (a *Applications) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "NAME", "SHORT NAME", "DEPLOY NAME"}}
}

func (a *Applications) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, a := range a.Applications {
		data[k] = []interface{}{a.ID, a.Name, a.ShortName, a.DeployName}
	}
	return data
}

func (a *Applications) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {a.Meta.Total, a.Meta.Links.Next, a.Meta.Links.Prev},
	}
}
