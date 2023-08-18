package printer

import (
	"github.com/vultr/govultr/v3"
)

// VPC2List Generates a printer display of all VPC 2.0 networks on the account
func VPC2List(vpc2s []govultr.VPC2, meta *govultr.Meta) {
	display(columns{"ID", "DATE CREATED", "REGION", "DESCRIPTION", "IP BLOCK", "PREFIX LENGTH"})
	for d := range vpc2s {
		display(columns{vpc2s[d].ID, vpc2s[d].DateCreated, vpc2s[d].Region, vpc2s[d].Description, vpc2s[d].IPBlock, vpc2s[d].PrefixLength})
	}

	Meta(meta)
	flush()
}

// VPC2 Generates a printer display of a given VPC 2.0 network
func VPC2(vpc2 *govultr.VPC2) {
	display(columns{"ID", "DATE CREATED", "REGION", "DESCRIPTION", "IP BLOCK", "PREFIX LENGTH"})
	display(columns{vpc2.ID, vpc2.DateCreated, vpc2.Region, vpc2.Description, vpc2.IPBlock, vpc2.PrefixLength})
	flush()
}

// VPC2ListNodes Generates a printer display of all nodes attached to a VPC 2.0 network
func VPC2ListNodes(vpc2Nodes []govultr.VPC2Node, meta *govultr.Meta) {
	display(columns{"ID", "IP ADDRESS", "MAC ADDRESS", "DESCRIPTION", "TYPE", "NODE STATUS"})
	for d := range vpc2Nodes {
		display(columns{vpc2Nodes[d].ID, vpc2Nodes[d].IPAddress, vpc2Nodes[d].MACAddress, vpc2Nodes[d].Description, vpc2Nodes[d].Type, vpc2Nodes[d].NodeStatus})
	}

	Meta(meta)
	flush()
}
