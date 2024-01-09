package regions

import (
	"encoding/json"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// RegionsPrinter represents the regions data from the API and contains the
// methods to format and print the data via the ResourceOutput interface
type RegionsPrinter struct {
	Regions []govultr.Region `json:"regions"`
	Meta    *govultr.Meta    `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (r *RegionsPrinter) JSON() []byte {
	json, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return json
}

// YAML provides the YAML formatted byte data
func (r *RegionsPrinter) YAML() []byte {
	yml, err := yaml.Marshal(r)
	if err != nil {
		panic(err.Error())
	}
	return yml
}

// Columns provides the columns for the printer
func (r *RegionsPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"ID",
		"CITY",
		"COUNTRY",
		"CONTINENT",
		"OPTIONS",
	}}
}

// Data provides the data for the printer
func (r *RegionsPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}

	if len(r.Regions) == 0 {
		data[0] = []interface{}{"---", "---", "---", "---", "---"}
		return data
	}

	for k, v := range r.Regions {
		data[k] = []interface{}{
			v.ID,
			v.City,
			v.Country,
			v.Continent,
			v.Options,
		}
	}

	return data
}

// Paging validates and forms the paging data for output
func (r *RegionsPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(r.Meta.Total, &r.Meta.Links.Next, &r.Meta.Links.Prev).Compose()
}

// RegionsAvailabilityPrinter represents the plan availability data for a
// region from the API and contains the methods to format and print the data
// via the ResourceOutput interface
type RegionsAvailabilityPrinter struct {
	// TODO: test json marshalling on this
	Plans *govultr.PlanAvailability `json:"available_plans"`
	Meta  *govultr.Meta             `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (r *RegionsAvailabilityPrinter) JSON() []byte {
	json, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return json
}

// YAML provides the YAML formatted byte data
func (r *RegionsAvailabilityPrinter) YAML() []byte {
	yml, err := yaml.Marshal(r)
	if err != nil {
		panic(err.Error())
	}
	return yml
}

// Columns provides the available plans columns for the printer
func (r *RegionsAvailabilityPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"AVAILABLE PLANS",
	}}
}

// Data provides the region availability plan data for the printer
func (r *RegionsAvailabilityPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}

	if len(r.Plans.AvailablePlans) == 0 {
		data[0] = []interface{}{"---"}
		return data
	}

	for k, v := range r.Plans.AvailablePlans {
		data[k] = []interface{}{
			v,
		}
	}

	return data
}

// Paging validates and forms the paging data for output
func (r *RegionsAvailabilityPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(r.Meta.Total, &r.Meta.Links.Next, &r.Meta.Links.Prev).Compose()
}
