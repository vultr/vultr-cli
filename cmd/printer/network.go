package printer

import "github.com/vultr/govultr"

func NetworkList(network []govultr.Network) {
	col := columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"}
	display(col)
	for _, n := range network {
		display(columns{n.NetworkID, n.RegionID, n.Description, n.V4Subnet, n.V4SubnetMask, n.DateCreated})
	}
	flush()
}
