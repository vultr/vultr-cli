package printer

import (
	"github.com/vultr/govultr/v2"
)

var _ ResourceOutput = &Plans{}

type Plans struct {
	Plan []govultr.Plan
	Meta *govultr.Meta
}

func (p *Plans) Json() {

}

func (p *Plans) Yaml() {

}

func (p *Plans) Columns() map[int][]interface{} {
	data := map[int][]interface{}{}
	data[0] = []interface{}{"ID", "VCPU COUNT", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"}
	return data
}

func (p *Plans) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, p := range p.Plan {
		data[k] = []interface{}{p.ID, p.VCPUCount, p.RAM, p.Disk, p.DiskCount, p.Bandwidth, p.MonthlyCost, p.Type, p.Locations}
	}
	return data
}

func (p *Plans) Paging() map[int][]interface{} {
	data := map[int][]interface{}{}
	data[0] = []interface{}{"======================================"}
	data[1] = []interface{}{"TOTAL", "NEXT PAGE", "PREV PAGE"}
	data[2] = []interface{}{p.Meta.Total, p.Meta.Links.Next, p.Meta.Links.Prev}
	return data
}


var _ ResourceOutput = &BaremetalPlans{}
type BaremetalPlans struct {
	Plan []govultr.BareMetalPlan
	Meta *govultr.Meta
}

func (b *BaremetalPlans) Json() {

}

func (b *BaremetalPlans) Yaml(){

}

func (b *BaremetalPlans) Columns() map[int][]interface{} {
	data := map[int][]interface{}{}
	data[0] = []interface{}{"ID", "CPU COUNT", "CPU MODEL", "CPU THREADS", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"}
	return data
}

func (b *BaremetalPlans) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, p := range b.Plan {
		data[k] = []interface{}{p.ID, p.CPUCount, p.CPUModel, p.CPUThreads, p.RAM, p.Disk, p.DiskCount, p.Bandwidth, p.MonthlyCost, p.Type, p.Locations}
	}
	return data
}

func (b *BaremetalPlans) Paging() map[int][]interface{} {
	data := map[int][]interface{}{}
	data[0] = []interface{}{"======================================"}
	data[1] = []interface{}{"TOTAL", "NEXT PAGE", "PREV PAGE"}
	data[2] = []interface{}{b.Meta.Total, b.Meta.Links.Next, b.Meta.Links.Prev}
	return data
}