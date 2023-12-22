package printer

import "github.com/vultr/govultr/v3"

func NetworkList(network []govultr.Network, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"})

	if len(network) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range network {
		display(columns{
			network[i].NetworkID,
			network[i].Region,
			network[i].Description,
			network[i].V4Subnet,
			network[i].V4SubnetMask,
			network[i].DateCreated,
		})
	}

	Meta(meta)
}

func Network(network *govultr.Network) {
	defer flush()

	display(columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"})
	display(columns{network.NetworkID, network.Region, network.Description, network.V4Subnet, network.V4SubnetMask, network.DateCreated})
}
