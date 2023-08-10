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

// VPC2 Generate a printer display of a given VPC 2.0 network
func VPC2(vpc2 *govultr.VPC2) {
	display(columns{"ID", "DATE CREATED", "REGION", "DESCRIPTION", "IP BLOCK", "PREFIX LENGTH"})
	display(columns{vpc2.ID, vpc2.DateCreated, vpc2.Region, vpc2.Description, vpc2.IPBlock, vpc2.PrefixLength})
	flush()
}
