package regions

import (
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// RegionsPrinter represents the regions data from the API and contains the
// methods to format and print the data via the ResourceOutput interface
type RegionsPrinter struct {
	Regions []govultr.Region `json:"regions"`
	Meta    *govultr.Meta    `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (r *RegionsPrinter) JSON() []byte {
	return printer.MarshalObject(r, "json")
}

// YAML provides the YAML formatted byte data
func (r *RegionsPrinter) YAML() []byte {
	return printer.MarshalObject(r, "yaml")
}

// Columns provides the columns for the printer
func (r *RegionsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"CITY",
		"COUNTRY",
		"CONTINENT",
		"OPTIONS",
	}}
}

// Data provides the data for the printer
func (r *RegionsPrinter) Data() [][]string {
	data := [][]string{}

	if len(r.Regions) == 0 {
		data = append(data, []string{"---", "---", "---", "---", "---"})
		return data
	}

	for i := range r.Regions {
		data = append(data, []string{
			r.Regions[i].ID,
			r.Regions[i].City,
			r.Regions[i].Country,
			r.Regions[i].Continent,
			printer.ArrayOfStringsToString(r.Regions[i].Options),
		})
	}

	return data
}

// Paging validates and forms the paging data for output
func (r *RegionsPrinter) Paging() [][]string {
	return printer.NewPaging(r.Meta.Total, &r.Meta.Links.Next, &r.Meta.Links.Prev).Compose()
}

// ======================================

// RegionsAvailabilityPrinter represents the plan availability data for a
// region from the API and contains the methods to format and print the data
// via the ResourceOutput interface
type RegionsAvailabilityPrinter struct {
	Plans *govultr.PlanAvailability `json:"available_plans"`
}

// JSON provides the JSON formatted byte data
func (r *RegionsAvailabilityPrinter) JSON() []byte {
	return printer.MarshalObject(r, "json")
}

// YAML provides the YAML formatted byte data
func (r *RegionsAvailabilityPrinter) YAML() []byte {
	return printer.MarshalObject(r, "yaml")
}

// Columns provides the available plans columns for the printer
func (r *RegionsAvailabilityPrinter) Columns() [][]string {
	return [][]string{0: {
		"AVAILABLE PLANS",
	}}
}

// Data provides the region availability plan data for the printer
func (r *RegionsAvailabilityPrinter) Data() [][]string {
	data := [][]string{}

	if len(r.Plans.AvailablePlans) == 0 {
		data = append(data, []string{"---"})
		return data
	}

	for i := range r.Plans.AvailablePlans {
		data = append(data, []string{
			r.Plans.AvailablePlans[i],
		})
	}

	return data
}

// Paging validates and forms the paging data for output
func (r *RegionsAvailabilityPrinter) Paging() [][]string {
	return nil
}
