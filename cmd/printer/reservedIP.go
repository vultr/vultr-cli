package printer

import "github.com/vultr/govultr/v2"

func ReservedIPList(reservedIP []govultr.ReservedIP, meta *govultr.Meta) {
	col := columns{"ID", "REGION", "IP TYPE", "SUBNET", "SUBNET SIZE", "LABEL", "ATTACHED TO"}
	display(col)
	for _, r := range reservedIP {
		display(columns{r.ID, r.Region, r.IPType, r.Subnet, r.SubnetSize, r.Label, r.InstanceID})
	}

	Meta(meta)
	flush()
}

func ReservedIP(reservedIP *govultr.ReservedIP) {
	col := columns{"ID", "REGION", "IP TYPE", "SUBNET", "SUBNET SIZE", "LABEL", "ATTACHED TO"}
	display(col)
	display(columns{reservedIP.ID, reservedIP.Region, reservedIP.IPType, reservedIP.Subnet, reservedIP.SubnetSize, reservedIP.Label, reservedIP.InstanceID})

	flush()
}
