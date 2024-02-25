package dns

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// DNSRecordsPrinter ...
type DNSRecordsPrinter struct {
	Records []govultr.DomainRecord `json:"records"`
	Meta    *govultr.Meta          `json:"meta"`
}

// JSON ...
func (d *DNSRecordsPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DNSRecordsPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DNSRecordsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"TYPE",
		"NAME",
		"DATA",
		"PRIORITY",
		"TTL",
	}}
}

// Data ...
func (d *DNSRecordsPrinter) Data() [][]string {
	if len(d.Records) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range d.Records {
		data = append(data, []string{
			d.Records[i].ID,
			d.Records[i].Type,
			d.Records[i].Name,
			d.Records[i].Data,
			strconv.Itoa(d.Records[i].Priority),
			strconv.Itoa(d.Records[i].TTL),
		})
	}

	return data
}

// Paging ...
func (d *DNSRecordsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(d.Meta).Compose()
}

// ======================================

// DNSRecordPrinter ...
type DNSRecordPrinter struct {
	Record govultr.DomainRecord `json:"records"`
}

// JSON ...
func (d *DNSRecordPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DNSRecordPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DNSRecordPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"TYPE",
		"NAME",
		"DATA",
		"PRIORITY",
		"TTL",
	}}
}

// Data ...
func (d *DNSRecordPrinter) Data() [][]string {
	return [][]string{0: {
		d.Record.ID,
		d.Record.Type,
		d.Record.Name,
		d.Record.Data,
		strconv.Itoa(d.Record.Priority),
		strconv.Itoa(d.Record.TTL),
	}}
}

// Paging ...
func (d *DNSRecordPrinter) Paging() [][]string {
	return nil
}

// ======================================

// DNSDomainsPrinter ...
type DNSDomainsPrinter struct {
	Domains []govultr.Domain `json:"domains"`
	Meta    *govultr.Meta    `json:"meta"`
}

// JSON ...
func (d *DNSDomainsPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DNSDomainsPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DNSDomainsPrinter) Columns() [][]string {
	return [][]string{0: {
		"DOMAIN",
		"DATE CREATED",
		"DNSSEC",
	}}
}

// Data ...
func (d *DNSDomainsPrinter) Data() [][]string {
	if len(d.Domains) == 0 {
		return [][]string{0: {"---", "---", "---"}}
	}

	var data [][]string
	for i := range d.Domains {
		data = append(data, []string{
			d.Domains[i].Domain,
			d.Domains[i].DateCreated,
			d.Domains[i].DNSSec,
		})
	}

	return data
}

// Paging ...
func (d *DNSDomainsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(d.Meta).Compose()
}

// ======================================

// DNSDomainPrinter ...
type DNSDomainPrinter struct {
	Domain govultr.Domain `json:"domain"`
}

// JSON ...
func (d *DNSDomainPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DNSDomainPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DNSDomainPrinter) Columns() [][]string {
	return [][]string{0: {
		"DOMAIN",
		"DATE CREATED",
		"DNS SEC",
	}}
}

// Data ...
func (d *DNSDomainPrinter) Data() [][]string {
	return [][]string{0: {
		d.Domain.Domain,
		d.Domain.DateCreated,
		d.Domain.DNSSec,
	}}
}

// Paging ...
func (d *DNSDomainPrinter) Paging() [][]string {
	return nil
}

// ======================================

// DNSSOAPrinter ...
type DNSSOAPrinter struct {
	SOA govultr.Soa `json:"dns_soa"`
}

// JSON ...
func (d *DNSSOAPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DNSSOAPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DNSSOAPrinter) Columns() [][]string {
	return [][]string{0: {
		"NS PRIMARY",
		"EMAIL",
	}}
}

// Data ...
func (d *DNSSOAPrinter) Data() [][]string {
	return [][]string{0: {
		d.SOA.NSPrimary,
		d.SOA.Email,
	}}
}

// Paging ...
func (d *DNSSOAPrinter) Paging() [][]string {
	return nil
}

// ======================================

// DNSSECPrinter ...
type DNSSECPrinter struct {
	SEC []string `json:"dns_sec"`
}

// JSON ...
func (d *DNSSECPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DNSSECPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DNSSECPrinter) Columns() [][]string {
	return [][]string{0: {
		"DNSSEC INFO",
	}}
}

// Data ...
func (d *DNSSECPrinter) Data() [][]string {
	if len(d.SEC) == 0 {
		return [][]string{0: {"---"}}
	}

	var data [][]string
	for i := range d.SEC {
		data = append(data, []string{d.SEC[i]})
	}

	return data
}

// Paging ...
func (d *DNSSECPrinter) Paging() [][]string {
	return nil
}
