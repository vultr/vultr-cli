package printer

import (
	"encoding/json"

	"github.com/go-yaml/yaml"
)

var _ ResourceOutput = &Version{}

type Version struct {
	Version string
}

func (v *Version) Json() []byte {
	prettyJSON, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (v *Version) Yaml() []byte {
	yam, err := yaml.Marshal(v)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (v *Version) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"version"}}

}

func (v *Version) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {v.Version}}

}

func (v Version) Paging() map[int][]interface{} {
	return nil
}
