package version

import (
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// VersionPrinter represents the version data
type VersionPrinter struct {
	Version string `json:"version"`
}

// JSON ...
func (v *VersionPrinter) JSON() []byte {
	return printer.MarshalObject(v, "json")
}

// YAML ...
func (v *VersionPrinter) YAML() []byte {
	return printer.MarshalObject(v, "yaml")
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
