package operatingSystems

import (
	"encoding/json"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// ApplicationsPrinter represents the plans data from the API
type OSPrinter struct {
	OperatingSystems []govultr.OS  `json:"os"`
	Meta             *govultr.Meta `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (o *OSPrinter) JSON() []byte {
	json, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return json
}

// YAML provides the YAML formatted byte data
func (o *OSPrinter) YAML() []byte {
	yml, err := yaml.Marshal(o)
	if err != nil {
		panic(err.Error())
	}
	return yml
}

// Columns provides the plan columns for the printer
func (o *OSPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"ID",
		"NAME",
		"ARCH",
		"FAMILY",
	}}
}

// Data provides the plan data for the printer
func (o *OSPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}

	if len(o.OperatingSystems) == 0 {
		data[0] = []interface{}{"---", "---", "---", "---"}
		return data
	}

	for k, v := range o.OperatingSystems {
		data[k] = []interface{}{
			v.ID,
			v.Name,
			v.Arch,
			v.Family,
		}
	}

	return data
}

// Paging validates and forms the paging data for output
func (o *OSPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(o.Meta.Total, &o.Meta.Links.Next, &o.Meta.Links.Prev).Compose()
}
