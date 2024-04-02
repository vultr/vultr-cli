package plans

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// PlansPrinter represents the plans data from the API
type PlansPrinter struct {
	Plans []govultr.Plan `json:"plans"`
	Meta  *govultr.Meta  `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (p *PlansPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML provides the YAML formatted byte data
func (p *PlansPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns provides the plan columns for the printer
func (p *PlansPrinter) Columns() [][]string {
	return [][]string{0: {
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
func (p *PlansPrinter) Data() [][]string {
	data := [][]string{}

	if len(p.Plans) == 0 {
		data = append(data, []string{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		return data
	}

	for i := range p.Plans {
		data = append(data, []string{
			p.Plans[i].ID,
			strconv.Itoa(p.Plans[i].VCPUCount),
			strconv.Itoa(p.Plans[i].RAM),
			strconv.Itoa(p.Plans[i].Disk),
			strconv.Itoa(p.Plans[i].DiskCount),
			strconv.Itoa(p.Plans[i].Bandwidth),
			strconv.FormatFloat(float64(p.Plans[i].MonthlyCost), 'f', utils.FloatPrecision, 32),
			p.Plans[i].Type,
			strconv.Itoa(p.Plans[i].GPUVRAM),
			p.Plans[i].GPUType,
			printer.ArrayOfStringsToString(p.Plans[i].Locations),
		})
	}

	return data
}

// Paging validates and forms the paging data for output
func (p *PlansPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(p.Meta).Compose()
}

// ======================================

// MetalPlansPrinter represents the bare metal plans data from the API
type MetalPlansPrinter struct {
	Plans []govultr.BareMetalPlan `json:"plans_metal"`
	Meta  *govultr.Meta           `json:"meta"`
}

// JSON provides the JSON formatted byte data
func (m *MetalPlansPrinter) JSON() []byte {
	return printer.MarshalObject(m, "json")
}

// YAML provides the YAML formatted byte data
func (m *MetalPlansPrinter) YAML() []byte {
	return printer.MarshalObject(m, "yaml")
}

// Columns provides the plan columns for the printer
func (m *MetalPlansPrinter) Columns() [][]string {
	return [][]string{0: {
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
func (m *MetalPlansPrinter) Data() [][]string {
	data := [][]string{}

	if len(m.Plans) == 0 {
		data = append(data, []string{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		return data
	}

	for i := range m.Plans {
		data = append(data, []string{
			m.Plans[i].ID,
			strconv.Itoa(m.Plans[i].CPUCount),
			m.Plans[i].CPUModel,
			strconv.Itoa(m.Plans[i].CPUThreads),
			strconv.Itoa(m.Plans[i].RAM),
			strconv.Itoa(m.Plans[i].Disk),
			strconv.Itoa(m.Plans[i].DiskCount),
			strconv.Itoa(m.Plans[i].Bandwidth),
			strconv.FormatFloat(float64(m.Plans[i].MonthlyCost), 'f', utils.FloatPrecision, 32),
			m.Plans[i].Type,
			printer.ArrayOfStringsToString(m.Plans[i].Locations),
		})
	}

	return data
}

// Paging validates and forms the paging data for output
func (m *MetalPlansPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(m.Meta).Compose()
}
