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
