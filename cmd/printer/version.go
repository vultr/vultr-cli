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
	data := map[int][]interface{}{}
	data[0] = []interface{}{"Version"}
	return data
}

func (v *Version) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	data[0] = []interface{}{v.Version}
	return data
}

func (v Version) Paging() map[int][]interface{} {
	return nil
}
