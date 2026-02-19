package vpc

import (
	"strconv"

	"github.com/vultr/govultr/v3"

	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// VPCsPrinter ...
type VPCsPrinter struct {
	VPCs []govultr.VPC `json:"vpcs"`
	Meta *govultr.Meta `json:"meta"`
}

// JSON ...
func (s *VPCsPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *VPCsPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *VPCsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"DESCRIPTION",
		"V4 SUBNET",
		"V4 SUBNET MASK",
		"DATE CREATED",
	}}
}

// Data ...
func (s *VPCsPrinter) Data() [][]string {
	if len(s.VPCs) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range s.VPCs {
		data = append(data, []string{
			s.VPCs[i].ID,
			s.VPCs[i].Region,
			s.VPCs[i].Description,
			s.VPCs[i].V4Subnet,
			strconv.Itoa(s.VPCs[i].V4SubnetMask),
			s.VPCs[i].DateCreated,
		})
	}

	return data
}

// Paging ...
func (s *VPCsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(s.Meta).Compose()
}

// ======================================

// VPCPrinter ...
type VPCPrinter struct {
	VPC *govultr.VPC `json:"vpc"`
}

// JSON ...
func (s *VPCPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *VPCPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *VPCPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"DESCRIPTION",
		"V4 SUBNET",
		"V4 SUBNET MASK",
		"DATE CREATED",
	}}
}

// Data ...
func (s *VPCPrinter) Data() [][]string {
	return [][]string{0: {
		s.VPC.ID,
		s.VPC.Region,
		s.VPC.Description,
		s.VPC.V4Subnet,
		strconv.Itoa(s.VPC.V4SubnetMask),
		s.VPC.DateCreated,
	}}
}

// Paging ...
func (s *VPCPrinter) Paging() [][]string {
	return nil
}

// ======================================

// NATGatewaysPrinter ...
type NATGatewaysPrinter struct {
	NATGateways []govultr.NATGateway `json:"nat_gateways"`
	Meta        *govultr.Meta        `json:"meta"`
}

// JSON ...
func (ng *NATGatewaysPrinter) JSON() []byte {
	return printer.MarshalObject(ng, "json")
}

// YAML ...
func (ng *NATGatewaysPrinter) YAML() []byte {
	return printer.MarshalObject(ng, "yaml")
}

// Columns ...
func (ng *NATGatewaysPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (ng *NATGatewaysPrinter) Data() [][]string {
	if len(ng.NATGateways) == 0 {
		return [][]string{0: {"No NAT Gateways"}}
	}

	var data [][]string
	for i := range ng.NATGateways {
		data = append(data,
			[]string{"USERNAME", ng.NATGateways[i].Username},
			[]string{"PASSWORD", ng.NATGateways[i].Password},
		)

		if ng.NATGateways[i].Encryption != "" {
			data = append(data, []string{"ENCRYPTION", ng.NATGateways[i].Encryption})
		}

		if ng.NATGateways[i].AccessControl != nil {
			data = append(data,
				[]string{"ACCESS CONTROL"},
				[]string{"ACL CATEGORIES", printer.ArrayOfStringsToString(ng.NATGateways[i].AccessControl.ACLCategories)},
				[]string{"ACL CHANNELS", printer.ArrayOfStringsToString(ng.NATGateways[i].AccessControl.ACLChannels)},
				[]string{"ACL COMMANDS", printer.ArrayOfStringsToString(ng.NATGateways[i].AccessControl.ACLCommands)},
				[]string{"ACL KEYS", printer.ArrayOfStringsToString(ng.NATGateways[i].AccessControl.ACLKeys)},
			)
		}

		if ng.NATGateways[i].Permission != "" {
			data = append(data, []string{"PERMISSION", ng.NATGateways[i].Permission})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (ng *NATGatewaysPrinter) Paging() [][]string {
	paging := &printer.Total{Total: ng.Meta.Total}
	return paging.Compose()
}

// ======================================

// NATGatewayPrinter ...
type NATGatewayPrinter struct {
	NATGateway *govultr.NATGateway `json:"nat_gateway"`
}

// JSON ...
func (ng *NATGatewayPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (ng *NATGatewayPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (ng *NATGatewayPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (ng *NATGatewayPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"USERNAME", ng.NATGateway.Username},
		[]string{"PASSWORD", ng.NATGateway.Password},
	)

	if ng.NATGateway.Encryption != "" {
		data = append(data, []string{"ENCRYPTION", ng.NATGateway.Encryption})
	}

	if ng.NATGateway.AccessControl != nil {
		data = append(data,
			[]string{"ACCESS CONTROL"},
			[]string{"ACL CATEGORIES", printer.ArrayOfStringsToString(ng.NATGateway.AccessControl.ACLCategories)},
			[]string{"ACL CHANNELS", printer.ArrayOfStringsToString(ng.NATGateway.AccessControl.ACLChannels)},
			[]string{"ACL COMMANDS", printer.ArrayOfStringsToString(ng.NATGateway.AccessControl.ACLCommands)},
			[]string{"ACL KEYS", printer.ArrayOfStringsToString(ng.NATGateway.AccessControl.ACLKeys)},
		)
	}

	if ng.NATGateway.Permission != "" {
		data = append(data, []string{"PERMISSION", ng.NATGateway.Permission})
	}

	if ng.NATGateway.AccessKey != "" {
		data = append(data, []string{"ACCESS KEY", ng.NATGateway.AccessKey})
	}

	if ng.NATGateway.AccessCert != "" {
		data = append(data, []string{"ACCESS CERT", ng.NATGateway.AccessCert})
	}

	return data
}

// Paging ...
func (ng *NATGatewayPrinter) Paging() [][]string {
	return nil
}
