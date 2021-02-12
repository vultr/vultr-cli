package printer

import (
	"github.com/vultr/govultr/v2"
)

func Plan(plan []govultr.Plan, meta *govultr.Meta) {
	col := columns{"ID", "VCPU COUNT", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"}
	display(col)
	for _, p := range plan {
		display(columns{p.ID, p.VCPUCount, p.RAM, p.Disk, p.DiskCount, p.Bandwidth, p.MonthlyCost, p.Type, p.Locations})
	}

	Meta(meta)
	flush()
}

func PlanBareMetal(plan []govultr.BareMetalPlan, meta *govultr.Meta) {
	col := columns{"ID", "CPU COUNT", "CPU MODEL", "CPU THREADS", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"}
	display(col)
	for _, p := range plan {
		display(columns{p.ID, p.CPUCount, p.CPUModel, p.CPUThreads, p.RAM, p.Disk, p.DiskCount, p.Bandwidth, p.MonthlyCost, p.Type, p.Locations})
	}

	Meta(meta)
	flush()
}
