package printer

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var _ ResourceOutput = &Version{}

type Version struct {
	Version string
}

// JSON ...
func (v *Version) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

// YAML ...
func (v *Version) YAML() []byte {
	yam, err := yaml.Marshal(v)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (v *Version) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"version"}}

}

// Data ...
func (v *Version) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {v.Version}}

}

// Paging ...
func (v Version) Paging() map[int][]interface{} {
	return nil
}
