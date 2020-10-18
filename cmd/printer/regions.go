package printer

import (
	"github.com/vultr/govultr/v2"
)

func Regions(avail []govultr.Region, meta *govultr.Meta) {
	col := columns{"ID", "CITY", "COUNTRY", "CONTINENT", "OPTIONS"}
	display(col)
	for _, r := range avail {
		display(columns{r.ID, r.City, r.Country, r.Continent, r.Options})
	}

	Meta(meta)
	flush()
}

func RegionAvailability(avail *govultr.PlanAvailability) {
	display(columns{"AVAILABLE PLANS"})

	for _, r := range avail.AvailablePlans {
		display(columns{r})
	}

	flush()
}
