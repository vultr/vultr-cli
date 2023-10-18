package printer

import "github.com/vultr/govultr/v3"

func ReservedIPList(reservedIP []govultr.ReservedIP, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "REGION", "IP TYPE", "SUBNET", "SUBNET SIZE", "LABEL", "ATTACHED TO"})

	if len(reservedIP) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range reservedIP {
		display(columns{
			reservedIP[i].ID,
			reservedIP[i].Region,
			reservedIP[i].IPType,
			reservedIP[i].Subnet,
			reservedIP[i].SubnetSize,
			reservedIP[i].Label,
			reservedIP[i].InstanceID,
		})
	}

	Meta(meta)
}

func ReservedIP(reservedIP *govultr.ReservedIP) {
	defer flush()

	display(columns{"ID", "REGION", "IP TYPE", "SUBNET", "SUBNET SIZE", "LABEL", "ATTACHED TO"})
	display(columns{
		reservedIP.ID,
		reservedIP.Region,
		reservedIP.IPType,
		reservedIP.Subnet,
		reservedIP.SubnetSize,
		reservedIP.Label,
		reservedIP.InstanceID,
	})
}
