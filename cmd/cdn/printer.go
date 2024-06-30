package cdn

import (
	"fmt"
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// PullZonesPrinter ...
type PullZonesPrinter struct {
	PullZones []govultr.CDNZone `json:"pull_zones"`
	Meta      *govultr.Meta     `json:"meta"`
}

// JSON ...
func (p *PullZonesPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PullZonesPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PullZonesPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PullZonesPrinter) Data() [][]string {
	if len(p.PullZones) == 0 {
		return [][]string{0: {"No active CDN pull zones"}}
	}

	var data [][]string
	for i := range p.PullZones {
		data = append(data,
			[]string{"---------------------------"},
			[]string{"ID", p.PullZones[i].ID},
			[]string{"STATUS", p.PullZones[i].Status},
			[]string{"LABEL", p.PullZones[i].Label},
			[]string{"ORIGIN SCHEME", p.PullZones[i].OriginScheme},
			[]string{"ORIGIN DOMAIN", p.PullZones[i].OriginDomain},
			[]string{"URL", p.PullZones[i].CDNURL},
			[]string{"DATE CREATED", p.PullZones[i].DateCreated},
			[]string{"DATE PURGED", p.PullZones[i].DatePurged},
			[]string{"CACHE SIZE", strconv.Itoa(p.PullZones[i].CacheSize)},
			[]string{"REQUESTS", strconv.Itoa(p.PullZones[i].Requests)},
			[]string{"BYTES IN", strconv.Itoa(p.PullZones[i].BytesIn)},
			[]string{"BYTES OUT", strconv.Itoa(p.PullZones[i].BytesOut)},
			[]string{"PACKETS PER SECOND", strconv.Itoa(p.PullZones[i].PacketsPerSec)},
			[]string{" "},
			[]string{"OPTIONS"},
			[]string{"CORS", strconv.FormatBool(p.PullZones[i].CORS)},
			[]string{"GZIP", strconv.FormatBool(p.PullZones[i].GZIP)},
			[]string{"BLOCK AI", strconv.FormatBool(p.PullZones[i].BlockAI)},
			[]string{"BLOCK BAD BOTS", strconv.FormatBool(p.PullZones[i].BlockBadBots)},
			[]string{" "},
			[]string{"REGIONS"},
		)

		for j := range p.PullZones[i].Regions {
			data = append(data, []string{fmt.Sprintf(" - %s", p.PullZones[i].Regions[j])})
		}
	}

	return data
}

// Paging ...
func (p *PullZonesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(p.Meta).Compose()
}

// ======================================

// PullZonePrinter ...
type PullZonePrinter struct {
	PullZone *govultr.CDNZone `json:"pull_zone"`
}

// JSON ...
func (p *PullZonePrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PullZonePrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PullZonePrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PullZonePrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", p.PullZone.ID},
		[]string{"STATUS", p.PullZone.Status},
		[]string{"LABEL", p.PullZone.Label},
		[]string{"ORIGIN SCHEME", p.PullZone.OriginScheme},
		[]string{"ORIGIN DOMAIN", p.PullZone.OriginDomain},
		[]string{"URL", p.PullZone.CDNURL},
		[]string{"DATE CREATED", p.PullZone.DateCreated},
		[]string{"DATE PURGED", p.PullZone.DatePurged},
		[]string{"CACHE SIZE", strconv.Itoa(p.PullZone.CacheSize)},
		[]string{"REQUESTS", strconv.Itoa(p.PullZone.Requests)},
		[]string{"BYTES IN", strconv.Itoa(p.PullZone.BytesIn)},
		[]string{"BYTES OUT", strconv.Itoa(p.PullZone.BytesOut)},
		[]string{"PACKETS PER SECOND", strconv.Itoa(p.PullZone.PacketsPerSec)},
		[]string{" "},
		[]string{"OPTIONS"},
		[]string{"CORS", strconv.FormatBool(p.PullZone.CORS)},
		[]string{"GZIP", strconv.FormatBool(p.PullZone.GZIP)},
		[]string{"BLOCK AI", strconv.FormatBool(p.PullZone.BlockAI)},
		[]string{"BLOCK BAD BOTS", strconv.FormatBool(p.PullZone.BlockBadBots)},
		[]string{" "},
		[]string{"REGIONS"},
	)

	for i := range p.PullZone.Regions {
		data = append(data, []string{fmt.Sprintf(" - %s", p.PullZone.Regions[i])})
	}

	return data
}

// Paging ...
func (p *PullZonePrinter) Paging() [][]string {
	return nil
}

// ======================================

// PushZonesPrinter ...
type PushZonesPrinter struct {
	PushZones []govultr.CDNZone `json:"push_zones"`
	Meta      *govultr.Meta     `json:"meta"`
}

// JSON ...
func (p *PushZonesPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PushZonesPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PushZonesPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PushZonesPrinter) Data() [][]string {
	if len(p.PushZones) == 0 {
		return [][]string{0: {"No active CDN push zones"}}
	}

	var data [][]string
	for i := range p.PushZones {
		data = append(data,
			[]string{"---------------------------"},
			[]string{"ID", p.PushZones[i].ID},
			[]string{"STATUS", p.PushZones[i].Status},
			[]string{"LABEL", p.PushZones[i].Label},
			[]string{"URL", p.PushZones[i].CDNURL},
			[]string{"DATE CREATED", p.PushZones[i].DateCreated},
			[]string{"CACHE SIZE", strconv.Itoa(p.PushZones[i].CacheSize)},
			[]string{"REQUESTS", strconv.Itoa(p.PushZones[i].Requests)},
			[]string{"BYTES IN", strconv.Itoa(p.PushZones[i].BytesIn)},
			[]string{"BYTES OUT", strconv.Itoa(p.PushZones[i].BytesOut)},
			[]string{"PACKETS PER SECOND", strconv.Itoa(p.PushZones[i].PacketsPerSec)},
			[]string{" "},
			[]string{"OPTIONS"},
			[]string{"CORS", strconv.FormatBool(p.PushZones[i].CORS)},
			[]string{"GZIP", strconv.FormatBool(p.PushZones[i].GZIP)},
			[]string{"BLOCK AI", strconv.FormatBool(p.PushZones[i].BlockAI)},
			[]string{"BLOCK BAD BOTS", strconv.FormatBool(p.PushZones[i].BlockBadBots)},
			[]string{" "},
			[]string{"REGIONS"},
		)

		for j := range p.PushZones[i].Regions {
			data = append(data, []string{fmt.Sprintf(" - %s", p.PushZones[i].Regions[j])})
		}
	}

	return data
}

// Paging ...
func (p *PushZonesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(p.Meta).Compose()
}

// ======================================

// PushZonePrinter ...
type PushZonePrinter struct {
	PushZone *govultr.CDNZone `json:"push_zone"`
}

// JSON ...
func (p *PushZonePrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PushZonePrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PushZonePrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PushZonePrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", p.PushZone.ID},
		[]string{"STATUS", p.PushZone.Status},
		[]string{"LABEL", p.PushZone.Label},
		[]string{"ORIGIN SCHEME", p.PushZone.OriginScheme},
		[]string{"ORIGIN DOMAIN", p.PushZone.OriginDomain},
		[]string{"URL", p.PushZone.CDNURL},
		[]string{"DATE CREATED", p.PushZone.DateCreated},
		[]string{"DATE PURGED", p.PushZone.DatePurged},
		[]string{"CACHE SIZE", strconv.Itoa(p.PushZone.CacheSize)},
		[]string{"REQUESTS", strconv.Itoa(p.PushZone.Requests)},
		[]string{"BYTES IN", strconv.Itoa(p.PushZone.BytesIn)},
		[]string{"BYTES OUT", strconv.Itoa(p.PushZone.BytesOut)},
		[]string{"PACKETS PER SECOND", strconv.Itoa(p.PushZone.PacketsPerSec)},
		[]string{" "},
		[]string{"OPTIONS"},
		[]string{"CORS", strconv.FormatBool(p.PushZone.CORS)},
		[]string{"GZIP", strconv.FormatBool(p.PushZone.GZIP)},
		[]string{"BLOCK AI", strconv.FormatBool(p.PushZone.BlockAI)},
		[]string{"BLOCK BAD BOTS", strconv.FormatBool(p.PushZone.BlockBadBots)},
		[]string{" "},
		[]string{"REGIONS"},
	)

	for i := range p.PushZone.Regions {
		data = append(data, []string{fmt.Sprintf(" - %s", p.PushZone.Regions[i])})
	}

	return data
}

// Paging ...
func (p *PushZonePrinter) Paging() [][]string {
	return nil
}

// ======================================

// PushZoneFilesPrinter ...
type PushZoneFilesPrinter struct {
	FileData *govultr.CDNZoneFileData `json:"file_data"`
}

// JSON ...
func (p *PushZoneFilesPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PushZoneFilesPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PushZoneFilesPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PushZoneFilesPrinter) Data() [][]string {
	if len(p.FileData.Files) == 0 {
		return [][]string{0: {"No CDN push zone file data"}}
	}

	var data [][]string
	for i := range p.FileData.Files {
		data = append(data,
			[]string{"NAME", p.FileData.Files[i].Name},
			[]string{"SIZE", strconv.Itoa(p.FileData.Files[i].Size)},
			[]string{"DATE MODIFIED", p.FileData.Files[i].DateModified},
		)
	}

	data = append(data,
		[]string{" "},
		[]string{"COUNT", strconv.Itoa(p.FileData.Count)},
		[]string{"TOTAL SIZE", strconv.Itoa(p.FileData.Size)},
	)

	return data
}

// Paging ...
func (p *PushZoneFilesPrinter) Paging() [][]string {
	return nil
}

// ======================================

// PushZoneFilePrinter ...
type PushZoneFilePrinter struct {
	File *govultr.CDNZoneFile `json:"file"`
}

// JSON ...
func (p *PushZoneFilePrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PushZoneFilePrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PushZoneFilePrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PushZoneFilePrinter) Data() [][]string {
	return [][]string{
		[]string{"NAME", p.File.Name},
		[]string{"SIZE", strconv.Itoa(p.File.Size)},
		[]string{"DATE MODIFIED", p.File.DateModified},
		// TODO: missing from govultr
		// []string{"MIME", strconv.Itoa(p.File.MIME)},
		// []string{"CONTENT", p.File.Content},
	}
}

// Paging ...
func (p *PushZoneFilePrinter) Paging() [][]string {
	return nil
}

// ======================================

// PushZoneEndpointPrinter ...
type PushZoneEndpointPrinter struct {
	Endpoint *govultr.CDNZoneEndpoint `json:"upload_endpoint"`
}

// JSON ...
func (p *PushZoneEndpointPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PushZoneEndpointPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PushZoneEndpointPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PushZoneEndpointPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"URL", p.Endpoint.URL},
		[]string{"ACL", p.Endpoint.Inputs.ACL},
		[]string{"KEY", p.Endpoint.Inputs.Key},
		[]string{"X-AMZ-CREDENTIAL", p.Endpoint.Inputs.Credential},
		[]string{"X-AMZ-ALGORITHM", p.Endpoint.Inputs.Algorithm},
		[]string{"X-AMZ-SIGNATURE", p.Endpoint.Inputs.Signature},
		[]string{"POLICY", p.Endpoint.Inputs.Policy},
	)

	return data
}

// Paging ...
func (p *PushZoneEndpointPrinter) Paging() [][]string {
	return nil
}

// ======================================
