package printer

import (
	"github.com/vultr/govultr"
)

func Plan(plan []govultr.Plan) {
	col := columns{"VPSPLANID", "NAME", "VCPU COUNT", "RAM", "DISK", "BANDWIDTH", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS", "DEPRECATED"}
	display(col)
	for _, p := range plan {
		display(columns{p.PlanID, p.Name, p.VCPUs, p.RAM, p.Disk, p.Bandwidth, p.BandwidthGB, p.Price, p.PlanType, p.Regions, p.Deprecated})
	}
	flush()
}

func PlanBareMetal(plan []govultr.BareMetalPlan) {
	col := columns{"METALPLANID", "NAME", "VCPU COUNT", "RAM", "DISK", "BANDWIDTH TB", "PRICE PER MONTH", "TYPE", "REGIONS", "DEPRECATED"}
	display(col)
	for _, p := range plan {
		display(columns{p.PlanID, p.Name, p.CPUs, p.RAM, p.Disk, p.BandwidthTB, p.Price, p.PlanType, p.Regions, p.Deprecated})
	}
	flush()
}
