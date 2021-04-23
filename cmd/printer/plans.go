package printer

import (
	"encoding/json"

	"github.com/go-yaml/yaml"
	"github.com/vultr/govultr/v2"
)

var _ ResourceOutput = &Plans{}

type Plans struct {
	Plan []govultr.Plan
	Meta *govultr.Meta
}

// why do we have this struct? when it's literally the same thing up above? go to bed your tired
type plansBase struct {
	Plans []govultr.Plan `json:"plans"`
	Meta  *govultr.Meta  `json:"meta"`
}

func (p *Plans) Json() []byte {
	prettyJSON, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (p *Plans) Yaml() []byte {
	yam, err := yaml.Marshal(p)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (p *Plans) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "VCPU COUNT", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"}}
}

func (p *Plans) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, p := range p.Plan {
		data[k] = []interface{}{p.ID, p.VCPUCount, p.RAM, p.Disk, p.DiskCount, p.Bandwidth, p.MonthlyCost, p.Type, p.Locations}
	}
	return data
}

func (p *Plans) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {p.Meta.Total, p.Meta.Links.Next, p.Meta.Links.Prev},
	}
}

var _ ResourceOutput = &BaremetalPlans{}

type BaremetalPlans struct {
	Plan []govultr.BareMetalPlan
	Meta *govultr.Meta
}

func (b *BaremetalPlans) Json() []byte {
	prettyJSON, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (b *BaremetalPlans) Yaml() []byte {
	yam, err := yaml.Marshal(b)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (b *BaremetalPlans) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "CPU COUNT", "CPU MODEL", "CPU THREADS", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"}}
}

func (b *BaremetalPlans) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, p := range b.Plan {
		data[k] = []interface{}{p.ID, p.CPUCount, p.CPUModel, p.CPUThreads, p.RAM, p.Disk, p.DiskCount, p.Bandwidth, p.MonthlyCost, p.Type, p.Locations}
	}
	return data
}

func (b *BaremetalPlans) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {b.Meta.Total, b.Meta.Links.Next, b.Meta.Links.Prev},
	}
}
