package printer

import "github.com/vultr/govultr/v2"

func VPCList(vpc []govultr.VPC, meta *govultr.Meta) {
	col := columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"}
	display(col)
	for _, n := range vpc {
		display(columns{n.ID, n.Region, n.Description, n.V4Subnet, n.V4SubnetMask, n.DateCreated})
	}

	Meta(meta)
	flush()
}

func VPC(vpc *govultr.VPC) {
	display(columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"})
	display(columns{vpc.ID, vpc.Region, vpc.Description, vpc.V4Subnet, vpc.V4SubnetMask, vpc.DateCreated})

	flush()
}
