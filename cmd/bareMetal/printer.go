package baremetal

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// BareMetalsPrinter ...
type BareMetalsPrinter struct {
	BareMetals []govultr.BareMetalServer `json:"bare_metals"`
	Meta       *govultr.Meta
}

// JSON ...
func (b *BareMetalsPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BareMetalsPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BareMetalsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"IP",
		"TAG",
		"MAC ADDRESS",
		"LABEL",
		"OS",
		"STATUS",
		"REGION",
		"CPU",
		"RAM",
		"DISK",
		"FEATURES",
		"TAGS",
	}}
}

// Data ...
func (b *BareMetalsPrinter) Data() [][]string {
	var data [][]string
	for i := range b.BareMetals {
		data = append(data, []string{
			b.BareMetals[i].ID,
			b.BareMetals[i].MainIP,
			b.BareMetals[i].Tag, //nolint: staticcheck
			strconv.Itoa(b.BareMetals[i].MacAddress),
			b.BareMetals[i].Label,
			b.BareMetals[i].Os,
			b.BareMetals[i].Status,
			b.BareMetals[i].Region,
			strconv.Itoa(b.BareMetals[i].CPUCount),
			b.BareMetals[i].RAM,
			b.BareMetals[i].Disk,
			printer.ArrayOfStringsToString(b.BareMetals[i].Features),
			printer.ArrayOfStringsToString(b.BareMetals[i].Tags),
		})
	}
	return data
}

// Paging ...
func (b *BareMetalsPrinter) Paging() [][]string {
	return printer.NewPaging(b.Meta.Total, &b.Meta.Links.Next, &b.Meta.Links.Prev).Compose()
}

// ======================================

// BareMetalPrinter ...
type BareMetalPrinter struct {
	BareMetal govultr.BareMetalServer `json:"bare_metal"`
}

// JSON ...
func (b *BareMetalPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BareMetalPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BareMetalPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"IP",
		"TAG",
		"MAC ADDRESS",
		"LABEL",
		"OS",
		"STATUS",
		"REGION",
		"CPU",
		"RAM",
		"DISK",
		"FEATURES",
		"TAGS",
	}}
}

// Data ...
func (b *BareMetalPrinter) Data() [][]string {
	return [][]string{0: {
		b.BareMetal.ID,
		b.BareMetal.MainIP,
		b.BareMetal.Tag, //nolint: staticcheck
		strconv.Itoa(b.BareMetal.MacAddress),
		b.BareMetal.Label,
		b.BareMetal.Os,
		b.BareMetal.Status,
		b.BareMetal.Region,
		strconv.Itoa(b.BareMetal.CPUCount),
		b.BareMetal.RAM,
		b.BareMetal.Disk,
		printer.ArrayOfStringsToString(b.BareMetal.Features),
		printer.ArrayOfStringsToString(b.BareMetal.Tags),
	}}
}

// Paging ...
func (b *BareMetalPrinter) Paging() [][]string {
	return nil
}
