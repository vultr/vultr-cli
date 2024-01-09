package plans

import (
	"encoding/json"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// PlansPrinter represents the plans data from the API
type PlansPrinter struct {
	Plans []govultr.Plan `json:"plans"`
	Meta  *govultr.Meta  `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (p *PlansPrinter) JSON() []byte {
	json, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return json
}

// YAML provides the YAML formatted byte data
func (p *PlansPrinter) YAML() []byte {
	yml, err := yaml.Marshal(p)
	if err != nil {
		panic(err.Error())
	}
	return yml
}

// Columns provides the plan columns for the printer
func (p *PlansPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"ID",
		"VCPU COUNT",
		"RAM",
		"DISK",
		"DISK COUNT",
		"BANDWIDTH GB",
		"PRICE PER MONTH",
		"TYPE",
		"GPU VRAM",
		"GPU TYPE",
		"REGIONS",
	}}
}

// Data provides the plan data for the printer
func (p *PlansPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}

	if len(p.Plans) == 0 {
		data[0] = []interface{}{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"}
		return data
	}

	for k, v := range p.Plans {
		data[k] = []interface{}{
			v.ID,
			v.VCPUCount,
			v.RAM,
			v.Disk,
			v.DiskCount,
			v.Bandwidth,
			v.MonthlyCost,
			v.Type,
			v.GPUVRAM,
			v.GPUType,
			v.Locations,
		}
	}

	return data
}

// Paging validates and forms the paging data for output
func (p *PlansPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(p.Meta.Total, &p.Meta.Links.Next, &p.Meta.Links.Prev).Compose()
}

// MetalPlansPrinter represents the bare metal plans data from the API
type MetalPlansPrinter struct {
	Plans []govultr.BareMetalPlan `json:"plans_metal"`
	Meta  *govultr.Meta           `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (m *MetalPlansPrinter) JSON() []byte {
	json, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	return json
}

// YAML provides the YAML formatted byte data
func (m *MetalPlansPrinter) YAML() []byte {
	yml, err := yaml.Marshal(m)
	if err != nil {
		panic(err.Error())
	}
	return yml
}

// Columns provides the plan columns for the printer
func (m *MetalPlansPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"ID",
		"CPU COUNT",
		"CPU MODEL",
		"CPU THREADS",
		"RAM",
		"DISK",
		"DISK COUNT",
		"BANDWIDTH GB",
		"PRICE PER MONTH",
		"TYPE",
		"REGIONS",
	}}
}

// Data provides the plan data for the printer
func (m *MetalPlansPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}

	if len(data) == 0 {
		data[0] = []interface{}{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"}
		return data
	}

	for k, v := range m.Plans {
		data[k] = []interface{}{
			v.ID,
			v.CPUCount,
			v.CPUModel,
			v.CPUThreads,
			v.RAM,
			v.Disk,
			v.DiskCount,
			v.Bandwidth,
			v.MonthlyCost,
			v.Type,
			v.Locations,
		}
	}

	return data
}

// Paging validates and forms the paging data for output
func (m *MetalPlansPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(m.Meta.Total, &m.Meta.Links.Next, &m.Meta.Links.Prev).Compose()
}
