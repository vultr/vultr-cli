package printer

import (
	"github.com/vultr/govultr"
)

func Regions(vultrOS []govultr.Region) {
	col := columns{"REGION ID", "NAME", "COUNTRY", "CONTINENT", "STATE", "DDOS", "BLOCK STORAGE", "REGION CODE"}
	display(col)
	for _, r := range vultrOS {
		display(columns{r.RegionID, r.Name, r.Country, r.Continent, r.State, r.Ddos, r.BlockStorage, r.RegionCode})
	}
	flush()
}

func RegionAvailability(avail []int) {
	col := columns{"AVAILABILITY"}
	display(col)
	for _, a := range avail {
		display(columns{a})
	}
	flush()
}
