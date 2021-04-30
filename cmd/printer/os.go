package printer

import (
	"encoding/json"
	"github.com/go-yaml/yaml"
	"github.com/vultr/govultr/v2"
)

var _ ResourceOutput = &OS{}

type OS struct {
	OS   []govultr.OS
	Meta *govultr.Meta
}

func (o *OS) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (o *OS) Yaml() []byte {
	yam, err := yaml.Marshal(o)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (o *OS) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "NAME", "ARCH", "FAMILY"}}
}
func (o *OS) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, os := range o.OS {
		data[k] = []interface{}{os.ID, os.Name, os.Arch, os.Family}
	}
	return data
}

func (o *OS) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {o.Meta.Total, o.Meta.Links.Next, o.Meta.Links.Prev},
	}
}
