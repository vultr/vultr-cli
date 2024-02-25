package vpc2

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// VPC2sPrinter ...
type VPC2sPrinter struct {
	VPC2s []govultr.VPC2 `json:"vpcs"`
	Meta  *govultr.Meta  `json:"meta"`
}

// JSON ...
func (s *VPC2sPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *VPC2sPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *VPC2sPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"REGION",
		"DESCRIPTION",
		"IP BLOCK",
		"PREFIX LENGTH",
	}}
}

// Data ...
func (s *VPC2sPrinter) Data() [][]string {
	if len(s.VPC2s) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range s.VPC2s {
		data = append(data, []string{
			s.VPC2s[i].ID,
			s.VPC2s[i].DateCreated,
			s.VPC2s[i].Region,
			s.VPC2s[i].Description,
			s.VPC2s[i].IPBlock,
			strconv.Itoa(s.VPC2s[i].PrefixLength),
		})
	}

	return data
}

// Paging ...
func (s *VPC2sPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(s.Meta).Compose()
}

// ======================================

// VPC2Printer ...
type VPC2Printer struct {
	VPC2 *govultr.VPC2 `json:"vpc"`
}

// JSON ...
func (s *VPC2Printer) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *VPC2Printer) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *VPC2Printer) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"REGION",
		"DESCRIPTION",
		"IP BLOCK",
		"PREFIX LENGTH",
	}}
}

// Data ...
func (s *VPC2Printer) Data() [][]string {
	return [][]string{0: {
		s.VPC2.ID,
		s.VPC2.DateCreated,
		s.VPC2.Region,
		s.VPC2.Description,
		s.VPC2.IPBlock,
		strconv.Itoa(s.VPC2.PrefixLength),
	}}
}

// Paging ...
func (s *VPC2Printer) Paging() [][]string {
	return nil
}

// ======================================

// VPC2NodesPrinter ...
type VPC2NodesPrinter struct {
	Nodes []govultr.VPC2Node `json:"nodes"`
	Meta  *govultr.Meta      `json:"meta"`
}

// JSON ...
func (s *VPC2NodesPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *VPC2NodesPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *VPC2NodesPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"IP ADDRESS",
		"MAC ADDRESS",
		"DESCRIPTION",
		"TYPE",
		"NODE STATUS",
	}}
}

// Data ...
func (s *VPC2NodesPrinter) Data() [][]string {
	if len(s.Nodes) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range s.Nodes {
		data = append(data, []string{
			s.Nodes[i].ID,
			s.Nodes[i].IPAddress,
			strconv.Itoa(s.Nodes[i].MACAddress),
			s.Nodes[i].Description,
			s.Nodes[i].Type,
			s.Nodes[i].NodeStatus,
		})
	}

	return data
}

// Paging ...
func (s *VPC2NodesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(s.Meta).Compose()
}
