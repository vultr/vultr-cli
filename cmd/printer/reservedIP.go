package printer

import "github.com/vultr/govultr"

func ReservedIPList(reservedIP []govultr.ReservedIP) {
	col := columns{"ID", "REGION", "IP TYPE", "SUBNET", "SUBNET SIZE", "LABEL", "ATTACHED TO"}
	display(col)
	for _, r := range reservedIP {
		display(columns{r.ReservedIPID, r.RegionID, r.IPType, r.Subnet, r.SubnetSize, r.Label, r.AttachedID})
	}
	flush()
}
