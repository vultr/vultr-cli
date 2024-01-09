package version

import (
	"encoding/json"

	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// VersionPrinter represents the version data
type VersionPrinter struct {
	Version string `json:"version"`
}

// JSON ...
func (v *VersionPrinter) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(v, "", printer.JSONIndent)
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

// YAML ...
func (v *VersionPrinter) YAML() []byte {
	yam, err := yaml.Marshal(v)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (v *VersionPrinter) Columns() map[int][]interface{} {
	return nil
}

// Data ...
func (v *VersionPrinter) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {v.Version}}
}

// Paging ...
func (v *VersionPrinter) Paging() map[int][]interface{} {
	return nil
}
