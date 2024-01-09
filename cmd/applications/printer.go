package applications

import (
	"encoding/json"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// ApplicationsPrinter represents the plans data from the API
type ApplicationsPrinter struct {
	Applications []govultr.Application `json:"applications"`
	Meta         *govultr.Meta         `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (a *ApplicationsPrinter) JSON() []byte {
	json, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return json
}

// YAML provides the YAML formatted byte data
func (a *ApplicationsPrinter) YAML() []byte {
	yml, err := yaml.Marshal(a)
	if err != nil {
		panic(err.Error())
	}
	return yml
}

// Columns provides the plan columns for the printer
func (a *ApplicationsPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
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
func (a *ApplicationsPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}

	if len(a.Applications) == 0 {
		data[0] = []interface{}{"---", "---", "---", "---", "---", "---", "---"}
		return data
	}

	for k, v := range a.Applications {
		data[k] = []interface{}{
			v.ID,
			v.Name,
			v.ShortName,
			v.DeployName,
			v.Type,
			v.Vendor,
			v.ImageID,
		}
	}

	return data
}

// Paging validates and forms the paging data for output
func (a *ApplicationsPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(a.Meta.Total, &a.Meta.Links.Next, &a.Meta.Links.Prev).Compose()
}
