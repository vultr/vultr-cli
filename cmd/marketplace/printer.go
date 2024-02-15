package marketplace

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// VariablesPrinter ...
type VariablesPrinter struct {
	Variables []govultr.MarketplaceAppVariable `json:"variables"`
}

// JSON ...
func (v *VariablesPrinter) JSON() []byte {
	return printer.MarshalObject(v, "json")
}

// YAML ...
func (v *VariablesPrinter) YAML() []byte {
	return printer.MarshalObject(v, "yaml")
}

// Columns ...
func (v *VariablesPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
		"DESCRIPTION",
		"REQUIRED",
	}}
}

// Data ...
func (v *VariablesPrinter) Data() [][]string {
	if len(v.Variables) == 0 {
		return [][]string{0: {"---", "---", "---"}}
	}

	var data [][]string
	for i := range v.Variables {
		data = append(data, []string{
			v.Variables[i].Name,
			v.Variables[i].Description,
			strconv.FormatBool(*v.Variables[i].Required),
		})
	}

	return data
}

// Paging ...
func (v *VariablesPrinter) Paging() [][]string {
	return nil
}
