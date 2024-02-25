package applications

import (
	"strconv"

	"github.com/vultr/govultr/v3"

	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// ApplicationsPrinter represents the plans data from the API
type ApplicationsPrinter struct {
	Applications []govultr.Application `json:"applications"`
	Meta         *govultr.Meta         `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (a *ApplicationsPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML provides the YAML formatted byte data
func (a *ApplicationsPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns provides the plan columns for the printer
func (a *ApplicationsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"SHORT NAME",
		"DEPLOY NAME",
		"TYPE",
		"VENDOR",
		"IMAGE ID",
	}}
}

// Data provides the plan data for the printer
func (a *ApplicationsPrinter) Data() [][]string {
	var data [][]string

	if len(a.Applications) == 0 {
		data = append(data, []string{"---", "---", "---", "---", "---", "---", "---"})
		return data
	}

	for i := range a.Applications {
		data = append(data, []string{
			strconv.Itoa(a.Applications[i].ID),
			a.Applications[i].Name,
			a.Applications[i].ShortName,
			a.Applications[i].DeployName,
			a.Applications[i].Type,
			a.Applications[i].Vendor,
			a.Applications[i].ImageID,
		})
	}

	return data
}

// Paging validates and forms the paging data for output
func (a *ApplicationsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(a.Meta).Compose()
}
