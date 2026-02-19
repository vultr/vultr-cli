package vpc

import (
	"fmt"
	"strconv"
	"strings"

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
			[]string{"ID", ng.NATGateways[i].ID},
			[]string{"VPC ID", ng.NATGateways[i].VPCID},
			[]string{"DATE CREATED", ng.NATGateways[i].DateCreated},
			[]string{"STATUS", ng.NATGateways[i].Status},
			[]string{"LABEL", ng.NATGateways[i].Label},
			[]string{"TAG", ng.NATGateways[i].Tag},
			[]string{"PUBLIC IPS", strings.Join(ng.NATGateways[i].PublicIPs, ",")},
			[]string{"PUBLIC IPSV6", strings.Join(ng.NATGateways[i].PublicIPsV6, ",")},
			[]string{"PRIVATE IPS", strings.Join(ng.NATGateways[i].PrivateIPs, ",")},
			[]string{"BILLING CHARGES", fmt.Sprintf("%.2f", ng.NATGateways[i].Billing.Charges)},
			[]string{"BILLING MONTHLY", fmt.Sprintf("%.2f", ng.NATGateways[i].Billing.Monthly)},
		)

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
	return printer.MarshalObject(ng, "json")
}

// YAML ...
func (ng *NATGatewayPrinter) YAML() []byte {
	return printer.MarshalObject(ng, "yaml")
}

// Columns ...
func (ng *NATGatewayPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (ng *NATGatewayPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", ng.NATGateway.ID},
		[]string{"VPC ID", ng.NATGateway.VPCID},
		[]string{"DATE CREATED", ng.NATGateway.DateCreated},
		[]string{"STATUS", ng.NATGateway.Status},
		[]string{"LABEL", ng.NATGateway.Label},
		[]string{"TAG", ng.NATGateway.Tag},
		[]string{"PUBLIC IPS", strings.Join(ng.NATGateway.PublicIPs, ",")},
		[]string{"PUBLIC IPSV6", strings.Join(ng.NATGateway.PublicIPsV6, ",")},
		[]string{"PRIVATE IPS", strings.Join(ng.NATGateway.PrivateIPs, ",")},
		[]string{"BILLING CHARGES", fmt.Sprintf("%.2f", ng.NATGateway.Billing.Charges)},
		[]string{"BILLING MONTHLY", fmt.Sprintf("%.2f", ng.NATGateway.Billing.Monthly)},
	)

	return data
}

// Paging ...
func (ng *NATGatewayPrinter) Paging() [][]string {
	return nil
}
