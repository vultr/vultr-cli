package printer

import (
	"github.com/vultr/govultr/v3"
)

func Plan(plan []govultr.Plan, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "VCPU COUNT", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "GPU VRAM", "GPU TYPE", "REGIONS"})

	if len(plan) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range plan {
		display(columns{
			plan[i].ID,
			plan[i].VCPUCount,
			plan[i].RAM,
			plan[i].Disk,
			plan[i].DiskCount,
			plan[i].Bandwidth,
			plan[i].MonthlyCost,
			plan[i].Type,
			plan[i].GPUVRAM,
			plan[i].GPUType,
			plan[i].Locations,
		})
	}

	Meta(meta)
}

func PlanBareMetal(plan []govultr.BareMetalPlan, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "CPU COUNT", "CPU MODEL", "CPU THREADS", "RAM", "DISK", "DISK COUNT", "BANDWIDTH GB", "PRICE PER MONTH", "TYPE", "REGIONS"})

	if len(plan) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range plan {
		display(columns{
			plan[i].ID,
			plan[i].CPUCount,
			plan[i].CPUModel,
			plan[i].CPUThreads,
			plan[i].RAM,
			plan[i].Disk,
			plan[i].DiskCount,
			plan[i].Bandwidth,
			plan[i].MonthlyCost,
			plan[i].Type,
			plan[i].Locations,
		})
	}

	Meta(meta)
}
