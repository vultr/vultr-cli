package printer

import (
	"github.com/vultr/govultr/v3"
)

// VPC2List Generates a printer display of all VPC 2.0 networks on the account
func VPC2List(vpc2s []govultr.VPC2, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "REGION", "DESCRIPTION", "IP BLOCK", "PREFIX LENGTH"})

	if len(vpc2s) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range vpc2s {
		display(columns{
			vpc2s[i].ID,
			vpc2s[i].DateCreated,
			vpc2s[i].Region,
			vpc2s[i].Description,
			vpc2s[i].IPBlock,
			vpc2s[i].PrefixLength,
		})
	}

	Meta(meta)
}

// VPC2 Generates a printer display of a given VPC 2.0 network
func VPC2(vpc2 *govultr.VPC2) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "REGION", "DESCRIPTION", "IP BLOCK", "PREFIX LENGTH"})
	display(columns{vpc2.ID, vpc2.DateCreated, vpc2.Region, vpc2.Description, vpc2.IPBlock, vpc2.PrefixLength})
}

// VPC2ListNodes Generates a printer display of all nodes attached to a VPC 2.0 network
func VPC2ListNodes(vpc2Nodes []govultr.VPC2Node, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "IP ADDRESS", "MAC ADDRESS", "DESCRIPTION", "TYPE", "NODE STATUS"})

	if len(vpc2Nodes) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range vpc2Nodes {
		display(columns{
			vpc2Nodes[i].ID,
			vpc2Nodes[i].IPAddress,
			vpc2Nodes[i].MACAddress,
			vpc2Nodes[i].Description,
			vpc2Nodes[i].Type,
			vpc2Nodes[i].NodeStatus,
		})
	}

	Meta(meta)
}
