package printer

import (
	"github.com/vultr/govultr/v3"
)

func Regions(avail []govultr.Region, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "CITY", "COUNTRY", "CONTINENT", "OPTIONS"})

	if len(avail) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range avail {
		display(columns{
			avail[i].ID,
			avail[i].City,
			avail[i].Country,
			avail[i].Continent,
			avail[i].Options,
		})
	}

	Meta(meta)
}

func RegionAvailability(avail *govultr.PlanAvailability) {
	defer flush()

	display(columns{"AVAILABLE PLANS"})

	if len(avail.AvailablePlans) == 0 {
		display(columns{"---"})
		return
	}

	for i := range avail.AvailablePlans {
		display(columns{avail.AvailablePlans[i]})
	}
}
