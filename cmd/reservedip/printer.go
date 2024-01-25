package reservedip

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// ReservedIPsPrinter...
type ReservedIPsPrinter struct {
	IPs  []govultr.ReservedIP `json:"reserved_ips"`
	Meta *govultr.Meta        `json:"meta"`
}

// JSON ...
func (r *ReservedIPsPrinter) JSON() []byte {
	return printer.MarshalObject(r, "json")
}

// YAML ...
func (r *ReservedIPsPrinter) YAML() []byte {
	return printer.MarshalObject(r, "yaml")
}

// Columns ...
func (r *ReservedIPsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"IP TYPE",
		"SUBNET",
		"SUBNET SIZE",
		"LABEL",
		"ATTACHED TO",
	}}
}

// Data ...
func (r *ReservedIPsPrinter) Data() [][]string {
	if len(r.IPs) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}

	}

	var data [][]string
	for i := range r.IPs {
		data = append(data, []string{
			r.IPs[i].ID,
			r.IPs[i].Region,
			r.IPs[i].IPType,
			r.IPs[i].Subnet,
			strconv.Itoa(r.IPs[i].SubnetSize),
			r.IPs[i].Label,
			r.IPs[i].InstanceID,
		})
	}

	return data
}

// Paging ...
func (r *ReservedIPsPrinter) Paging() [][]string {
	return printer.NewPaging(r.Meta.Total, &r.Meta.Links.Next, &r.Meta.Links.Prev).Compose()
}

// ======================================

// ReservedIPPrinter...
type ReservedIPPrinter struct {
	IP *govultr.ReservedIP `json:"reserved_ip"`
}

// JSON ...
func (r *ReservedIPPrinter) JSON() []byte {
	return printer.MarshalObject(r, "json")
}

// YAML ...
func (r *ReservedIPPrinter) YAML() []byte {
	return printer.MarshalObject(r, "yaml")
}

// Columns ...
func (r *ReservedIPPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"IP TYPE",
		"SUBNET",
		"SUBNET SIZE",
		"LABEL",
		"ATTACHED TO",
	}}
}

// Data ...
func (r *ReservedIPPrinter) Data() [][]string {
	return [][]string{0: {
		r.IP.ID,
		r.IP.Region,
		r.IP.IPType,
		r.IP.Subnet,
		strconv.Itoa(r.IP.SubnetSize),
		r.IP.Label,
		r.IP.InstanceID,
	}}
}

// Paging ...
func (r *ReservedIPPrinter) Paging() [][]string {
	return nil
}

// ======================================
