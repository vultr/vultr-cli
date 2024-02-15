package operatingsystems

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// OSPrinter represents the plans data from the API
type OSPrinter struct {
	OperatingSystems []govultr.OS  `json:"os"`
	Meta             *govultr.Meta `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (o *OSPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML provides the YAML formatted byte data
func (o *OSPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns provides the plan columns for the printer
func (o *OSPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"ARCH",
		"FAMILY",
	}}
}

// Data provides the plan data for the printer
func (o *OSPrinter) Data() [][]string {
	var data [][]string

	if len(o.OperatingSystems) == 0 {
		data = append(data, []string{"---", "---", "---", "---"})
		return data
	}

	for i := range o.OperatingSystems {
		data = append(data, []string{
			strconv.Itoa(o.OperatingSystems[i].ID),
			o.OperatingSystems[i].Name,
			o.OperatingSystems[i].Arch,
			o.OperatingSystems[i].Family,
		})
	}

	return data
}

// Paging validates and forms the paging data for output
func (o *OSPrinter) Paging() [][]string {
	return printer.NewPaging(o.Meta.Total, &o.Meta.Links.Next, &o.Meta.Links.Prev).Compose()
}
