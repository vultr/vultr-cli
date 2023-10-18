package printer

import "github.com/vultr/govultr/v3"

func VPCList(vpc []govultr.VPC, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"})

	if len(vpc) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range vpc {
		display(columns{
			vpc[i].ID,
			vpc[i].Region,
			vpc[i].Description,
			vpc[i].V4Subnet,
			vpc[i].V4SubnetMask,
			vpc[i].DateCreated,
		})
	}

	Meta(meta)
}

func VPC(vpc *govultr.VPC) {
	defer flush()

	display(columns{"ID", "REGION", "DESCRIPTION", "V4 SUBNET", "V4 SUBNET MASK", "DATE CREATED"})
	display(columns{vpc.ID, vpc.Region, vpc.Description, vpc.V4Subnet, vpc.V4SubnetMask, vpc.DateCreated})
}
