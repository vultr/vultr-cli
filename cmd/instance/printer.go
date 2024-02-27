package instance

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// InstancesPrinter ...
type InstancesPrinter struct {
	Instances []govultr.Instance `json:"instances"`
	Meta      *govultr.Meta      `json:"meta"`
}

// JSON ...
func (i *InstancesPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *InstancesPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *InstancesPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"IP",
		"LABEL",
		"OS",
		"STATUS",
		"REGION",
		"CPU",
		"RAM",
		"DISK",
		"BANDWIDTH",
		"TAGS",
	}}
}

// Data ...
func (i *InstancesPrinter) Data() [][]string {
	if len(i.Instances) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for j := range i.Instances {
		data = append(data, []string{
			i.Instances[j].ID,
			i.Instances[j].MainIP,
			i.Instances[j].Label,
			i.Instances[j].Os,
			i.Instances[j].Status,
			i.Instances[j].Region,
			strconv.Itoa(i.Instances[j].VCPUCount),
			strconv.Itoa(i.Instances[j].RAM),
			strconv.Itoa(i.Instances[j].Disk),
			strconv.Itoa(i.Instances[j].AllowedBandwidth),
			printer.ArrayOfStringsToString(i.Instances[j].Tags),
		})
	}
	return data
}

// Paging ...
func (i *InstancesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(i.Meta).Compose()
}

// ======================================

// InstancePrinter ...
type InstancePrinter struct {
	Instance *govultr.Instance `json:"instance"`
}

// JSON ...
func (i *InstancePrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *InstancePrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *InstancePrinter) Columns() [][]string {
	return [][]string{0: {"INSTANCE INFO"}}
}

// Data ...
func (i *InstancePrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", i.Instance.ID},
		[]string{"OS", i.Instance.Os},
		[]string{"OS ID", strconv.Itoa(i.Instance.OsID)},
		[]string{"APP ID", strconv.Itoa(i.Instance.AppID)},
		[]string{"RAM", strconv.Itoa(i.Instance.RAM)},
		[]string{"DISK", strconv.Itoa(i.Instance.Disk)},
		[]string{"MAIN IP", i.Instance.MainIP},
		[]string{"VCPU COUNT", strconv.Itoa(i.Instance.VCPUCount)},
		[]string{"REGION", i.Instance.Region},
		[]string{"DATE CREATED", i.Instance.DateCreated},
		[]string{"STATUS", i.Instance.Status},
		[]string{"ALLOWED BANDWIDTH", strconv.Itoa(i.Instance.AllowedBandwidth)},
		[]string{"NETMASK V4", i.Instance.NetmaskV4},
		[]string{"GATEWAY V4", i.Instance.GatewayV4},
		[]string{"POWER STATUS", i.Instance.PowerStatus},
		[]string{"SERVER STATE", i.Instance.ServerStatus},
		[]string{"PLAN", i.Instance.Plan},
		[]string{"LABEL", i.Instance.Label},
		[]string{"INTERNAL IP", i.Instance.InternalIP},
		[]string{"KVM URL", i.Instance.KVM},
		[]string{"FIREWALL GROUP ID", i.Instance.FirewallGroupID},
		[]string{"V6 MAIN IP", i.Instance.V6MainIP},
		[]string{"V6 NETWORK", i.Instance.V6Network},
		[]string{"V6 NETWORK SIZE", strconv.Itoa(i.Instance.V6NetworkSize)},
		[]string{"FEATURES", printer.ArrayOfStringsToString(i.Instance.Features)},
		[]string{"TAGS", printer.ArrayOfStringsToString(i.Instance.Tags)},
	)

	return data
}

// Paging ...
func (i *InstancePrinter) Paging() [][]string {
	return nil
}

// ======================================

// BandwidthPrinter ...
type BandwidthPrinter struct {
	Bandwidth govultr.Bandwidth `json:"bandwidth"`
}

// JSON ...
func (b *BandwidthPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BandwidthPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BandwidthPrinter) Columns() [][]string {
	return [][]string{0: {
		"DATE",
		"INCOMING BYTES",
		"OUTGOING BYTES",
	}}
}

// Data ...
func (b *BandwidthPrinter) Data() [][]string {
	var data [][]string
	for i := range b.Bandwidth.Bandwidth {
		data = append(data, []string{
			i,
			strconv.Itoa(b.Bandwidth.Bandwidth[i].IncomingBytes),
			strconv.Itoa(b.Bandwidth.Bandwidth[i].OutgoingBytes),
		})
	}

	return data
}

// Paging ...
func (b *BandwidthPrinter) Paging() [][]string {
	return nil
}

// ======================================

// BackupPrinter ...
type BackupPrinter struct {
	Backup govultr.BackupSchedule `json:"backup_schedule"`
}

// JSON ...
func (b *BackupPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BackupPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BackupPrinter) Columns() [][]string {
	return [][]string{0: {
		"ENABLED",
		"CRON TYPE",
		"NEXT RUN",
		"HOUR",
		"DOW",
		"DOM",
	}}
}

// Data ...
func (b *BackupPrinter) Data() [][]string {
	return [][]string{0: {
		strconv.FormatBool(*b.Backup.Enabled),
		b.Backup.Type,
		b.Backup.NextScheduleTimeUTC,
		strconv.Itoa(b.Backup.Hour),
		strconv.Itoa(b.Backup.Dow),
		strconv.Itoa(b.Backup.Dom),
	}}
}

// Paging ...
func (b *BackupPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ISOPrinter ...
type ISOPrinter struct {
	ISO govultr.Iso `json:"iso_status"`
}

// JSON ...
func (i *ISOPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *ISOPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *ISOPrinter) Columns() [][]string {
	return [][]string{0: {
		"ISO ID",
		"STATE",
	}}
}

// Data ...
func (i *ISOPrinter) Data() [][]string {
	return [][]string{0: {
		i.ISO.IsoID,
		i.ISO.State,
	}}
}

// Paging ...
func (i *ISOPrinter) Paging() [][]string {
	return nil
}

// ======================================

// OSsPrinter ...
type OSsPrinter struct {
	OSs []govultr.OS `json:"operating_systems"`
}

// JSON ...
func (o *OSsPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *OSsPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *OSsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"ARCH",
		"FAMILY",
	}}
}

// Data ...
func (o *OSsPrinter) Data() [][]string {
	if len(o.OSs) == 0 {
		return [][]string{0: {"---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range o.OSs {
		data = append(data, []string{
			strconv.Itoa(o.OSs[i].ID),
			o.OSs[i].Name,
			o.OSs[i].Arch,
			o.OSs[i].Family,
		})
	}

	return data
}

// Paging ...
func (o *OSsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AppsPrinter ...
type AppsPrinter struct {
	Apps []govultr.Application `json:"applications"`
}

// JSON ...
func (a *AppsPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AppsPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AppsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"SHORT NAME",
		"DEPLOY NAME",
		"TYPE",
		"VENDOR",
		"IMAGE ID",
	}}
}

// Data ...
func (a *AppsPrinter) Data() [][]string {
	if len(a.Apps) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range a.Apps {
		data = append(data, []string{
			strconv.Itoa(a.Apps[i].ID),
			a.Apps[i].Name,
			a.Apps[i].ShortName,
			a.Apps[i].DeployName,
			a.Apps[i].Type,
			a.Apps[i].Vendor,
			a.Apps[i].ImageID,
		})
	}

	return data
}

// Paging ...
func (a *AppsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// PlansPrinter ...
type PlansPrinter struct {
	Plans []string `json:"plans"`
}

// JSON ...
func (p *PlansPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PlansPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PlansPrinter) Columns() [][]string {
	return [][]string{0: {
		"PLAN NAME",
	}}
}

// Data ...
func (p *PlansPrinter) Data() [][]string {
	if len(p.Plans) == 0 {
		return [][]string{0: {"---"}}
	}

	var data [][]string
	for i := range p.Plans {
		data = append(data, []string{
			p.Plans[i],
		})
	}

	return data
}

// Paging ...
func (p *PlansPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ReverseIPsPrinter ...
type ReverseIPsPrinter struct {
	ReverseIPs []govultr.ReverseIP `json:"reverse_ips"`
}

// JSON ...
func (r *ReverseIPsPrinter) JSON() []byte {
	return printer.MarshalObject(r, "json")
}

// YAML ...
func (r *ReverseIPsPrinter) YAML() []byte {
	return printer.MarshalObject(r, "yaml")
}

// Columns ...
func (r *ReverseIPsPrinter) Columns() [][]string {
	return [][]string{0: {
		"IP",
		"REVERSE",
	}}
}

// Data ...
func (r *ReverseIPsPrinter) Data() [][]string {
	if len(r.ReverseIPs) == 0 {
		return [][]string{0: {"---", "---"}}
	}

	var data [][]string
	for j := range r.ReverseIPs {
		data = append(data, []string{
			r.ReverseIPs[j].IP,
			r.ReverseIPs[j].Reverse,
		})
	}

	return data
}

// Paging ...
func (r *ReverseIPsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// VPC2sPrinter ...
type VPC2sPrinter struct {
	VPC2s []govultr.VPC2Info `json:"vpcs"`
	Meta  *govultr.Meta      `json:"meta"`
}

// JSON ...
func (v *VPC2sPrinter) JSON() []byte {
	return printer.MarshalObject(v, "json")
}

// YAML ...
func (v *VPC2sPrinter) YAML() []byte {
	return printer.MarshalObject(v, "yaml")
}

// Columns ...
func (v *VPC2sPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"MAC ADDRESS",
		"IP ADDRESS",
	}}
}

// Data ...
func (v *VPC2sPrinter) Data() [][]string {
	var data [][]string

	if len(v.VPC2s) == 0 {
		return [][]string{0: {"---", "---", "---"}}
	}

	for i := range v.VPC2s {
		data = append(data, []string{
			v.VPC2s[i].ID,
			v.VPC2s[i].MacAddress,
			v.VPC2s[i].IPAddress,
		})
	}

	return data
}

// Paging ...
func (v *VPC2sPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(v.Meta).Compose()
}
