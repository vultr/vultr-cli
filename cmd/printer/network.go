package printer

import "github.com/vultr/govultr"

func NetworkList(network []govultr.Network, meta *govultr.Meta) {
	col := columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"}
	display(col)
	for _, n := range network {
		display(columns{n.NetworkID, n.Region, n.Description, n.V4Subnet, n.V4SubnetMask, n.DateCreated})
	}

	Meta(meta)
	flush()
}

func Network(network *govultr.Network) {
	display(columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"})
	display(columns{network.NetworkID, network.Region, network.Description, network.V4Subnet, network.V4SubnetMask, network.DateCreated})

	flush()
}
