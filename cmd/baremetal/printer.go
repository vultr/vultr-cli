package baremetal

import (
	"strconv"

	"github.com/vultr/govultr/v3"

	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// BareMetalsPrinter ...
type BareMetalsPrinter struct {
	BareMetals []govultr.BareMetalServer `json:"bare_metals"`
	Meta       *govultr.Meta             `json:"meta"`
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
	return printer.NewPagingFromMeta(b.Meta).Compose()
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
	return [][]string{0: {"BARE METAL INFO"}}
}

// Data ...
func (b *BareMetalPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", b.BareMetal.ID},
		[]string{"IP", b.BareMetal.MainIP},
		[]string{"USER SCHEME", b.BareMetal.UserScheme},
	)

	if b.BareMetal.DefaultPassword != "" {
		data = append(data, []string{"PASSWORD", b.BareMetal.DefaultPassword})
	} else {
		data = append(data, []string{"PASSWORD", "UNAVAILABLE"})
	}

	data = append(data,
		[]string{"MAC ADDRESS", strconv.Itoa(b.BareMetal.MacAddress)},
		[]string{"LABEL", b.BareMetal.Label},
		[]string{"OS", b.BareMetal.Os},
		[]string{"SNAPSHOT ID", b.BareMetal.SnapshotID},
		[]string{"STATUS", b.BareMetal.Status},
		[]string{"REGION", b.BareMetal.Region},
		[]string{"CPU", strconv.Itoa(b.BareMetal.CPUCount)},
		[]string{"RAM", b.BareMetal.RAM},
		[]string{"DISK", b.BareMetal.Disk},
		[]string{"FEATURES", printer.ArrayOfStringsToString(b.BareMetal.Features)},
		[]string{"TAGS", printer.ArrayOfStringsToString(b.BareMetal.Tags)},
	)

	return data
}

// Paging ...
func (b *BareMetalPrinter) Paging() [][]string {
	return nil
}

// ======================================

// BareMetalVNCPrinter ...
type BareMetalVNCPrinter struct {
	VNC govultr.VNCUrl `json:"vnc"`
}

// JSON ...
func (b *BareMetalVNCPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BareMetalVNCPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BareMetalVNCPrinter) Columns() [][]string {
	return [][]string{0: {
		"URL",
	}}
}

// Data ...
func (b *BareMetalVNCPrinter) Data() [][]string {
	return [][]string{0: {
		b.VNC.URL,
	}}
}

// Paging ...
func (b *BareMetalVNCPrinter) Paging() [][]string {
	return nil
}

// ======================================

// BareMetalBandwidthPrinter ...
type BareMetalBandwidthPrinter struct {
	Bandwidth govultr.Bandwidth `json:"all_bandwidth"`
}

// JSON ...
func (b *BareMetalBandwidthPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BareMetalBandwidthPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BareMetalBandwidthPrinter) Columns() [][]string {
	return [][]string{0: {
		"DATE",
		"INCOMING BYTES",
		"OUTGOING BYTES",
	}}
}

// Data ...
func (b *BareMetalBandwidthPrinter) Data() [][]string {
	var data [][]string
	for k := range b.Bandwidth.Bandwidth {
		data = append(data, []string{
			k,
			strconv.Itoa(b.Bandwidth.Bandwidth[k].IncomingBytes),
			strconv.Itoa(b.Bandwidth.Bandwidth[k].OutgoingBytes),
		})
	}

	return data
}

// Paging ...
func (b *BareMetalBandwidthPrinter) Paging() [][]string {
	return nil
}

// ======================================

// BareMetalVPC2sPrinter ...
type BareMetalVPC2sPrinter struct {
	VPC2s []govultr.VPC2Info `json:"vpcs"`
}

// JSON ...
func (b *BareMetalVPC2sPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BareMetalVPC2sPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BareMetalVPC2sPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"MAC ADDRESS",
		"IP ADDRESS",
	}}
}

// Data ...
func (b *BareMetalVPC2sPrinter) Data() [][]string {
	var data [][]string

	if len(b.VPC2s) == 0 {
		return [][]string{0: {"---", "---", "---"}}
	}

	for i := range b.VPC2s {
		data = append(data, []string{
			b.VPC2s[i].ID,
			b.VPC2s[i].MacAddress,
			b.VPC2s[i].IPAddress,
		})
	}

	return data
}

// Paging ...
func (b *BareMetalVPC2sPrinter) Paging() [][]string {
	return nil
}
