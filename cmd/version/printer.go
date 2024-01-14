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
	js, err := json.MarshalIndent(v, "", printer.JSONIndent)
	if err != nil {
		panic("move this into byte")
	}

	return js
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
func (v *VersionPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (v *VersionPrinter) Data() [][]string {
	return [][]string{0: {v.Version}}
}

// Paging ...
func (v *VersionPrinter) Paging() [][]string {
	return nil
}
