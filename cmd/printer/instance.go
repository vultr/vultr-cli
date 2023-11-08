package printer

import (
	"github.com/vultr/govultr/v3"
)

func InstanceBandwidth(bandwidth *govultr.Bandwidth) {
	display(columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"})
	for k, b := range bandwidth.Bandwidth {
		display(columns{k, b.IncomingBytes, b.OutgoingBytes})
	}
	flush()
}

func InstanceIPV4(ip []govultr.IPv4, meta *govultr.Meta) {
	defer flush()

	display(columns{"IP", "NETMASK", "GATEWAY", "TYPE", "REVERSE"})

	if len(ip) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range ip {
		display(columns{
			ip[i].IP,
			ip[i].Netmask,
			ip[i].Gateway,
			ip[i].Type,
			ip[i].Reverse,
		})
	}

	Meta(meta)
}

func InstanceIPV6(ip []govultr.IPv6, meta *govultr.Meta) {
	defer flush()

	display(columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"})

	if len(ip) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range ip {
		display(columns{
			ip[i].IP,
			ip[i].Network,
			ip[i].NetworkSize,
			ip[i].Type,
		})
	}

	Meta(meta)
}

func InstanceList(instance []govultr.Instance, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "IP", "LABEL", "OS", "STATUS", "REGION", "CPU", "RAM", "DISK", "BANDWIDTH", "TAGS"})

	if len(instance) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range instance {
		display(columns{
			instance[i].ID,
			instance[i].MainIP,
			instance[i].Label,
			instance[i].Os,
			instance[i].Status,
			instance[i].Region,
			instance[i].VCPUCount,
			instance[i].RAM,
			instance[i].Disk,
			instance[i].AllowedBandwidth,
			arrayOfStringsToString(instance[i].Tags),
		})
	}

	Meta(meta)
}

func Instance(instance *govultr.Instance) {
	defer flush()

	display(columns{"INSTANCE INFO"})

	display(columns{"ID", instance.ID})
	display(columns{"Os", instance.Os})
	display(columns{"RAM", instance.RAM})
	display(columns{"DISK", instance.Disk})
	display(columns{"MAIN IP", instance.MainIP})
	display(columns{"VCPU COUNT", instance.VCPUCount})
	display(columns{"REGION", instance.Region})
	display(columns{"DATE CREATED", instance.DateCreated})
	display(columns{"STATUS", instance.Status})
	display(columns{"ALLOWED BANDWIDTH", instance.AllowedBandwidth})
	display(columns{"NETMASK V4", instance.NetmaskV4})
	display(columns{"GATEWAY V4", instance.GatewayV4})
	display(columns{"POWER STATUS", instance.PowerStatus})
	display(columns{"SERVER STATE", instance.ServerStatus})
	display(columns{"PLAN", instance.Plan})
	display(columns{"LABEL", instance.Label})
	display(columns{"INTERNAL IP", instance.InternalIP})
	display(columns{"KVM URL", instance.KVM})
	display(columns{"TAG", instance.Tag}) //nolint:all
	display(columns{"OsID", instance.OsID})
	display(columns{"AppID", instance.AppID})
	display(columns{"FIREWALL GROUP ID", instance.FirewallGroupID})
	display(columns{"V6 MAIN IP", instance.V6MainIP})
	display(columns{"V6 NETWORK", instance.V6Network})
	display(columns{"V6 NETWORK SIZE", instance.V6NetworkSize})
	display(columns{"FEATURES", arrayOfStringsToString(instance.Features)})
	display(columns{"TAGS", arrayOfStringsToString(instance.Tags)})
}

func OsList(os []govultr.OS) {
	defer flush()

	display(columns{"ID", "NAME", "ARCH", "FAMILY"})

	if len(os) == 0 {
		display(columns{"---", "---", "---", "---"})
		return
	}

	for i := range os {
		display(columns{
			os[i].ID,
			os[i].Name,
			os[i].Arch,
			os[i].Family,
		})
	}
}

func AppList(app []govultr.Application) {
	defer flush()

	display(columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME", "TYPE", "VENDOR", "IMAGE ID"})

	if len(app) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		return
	}

	for i := range app {
		display(columns{
			app[i].ID,
			app[i].Name,
			app[i].ShortName,
			app[i].DeployName,
			app[i].Type,
			app[i].Vendor,
			app[i].ImageID,
		})
	}
}

func BackupsGet(b *govultr.BackupSchedule) {
	defer flush()

	display(columns{"ENABLED", "CRON TYPE", "NEXT RUN", "HOUR", "DOW", "DOM"})
	display(columns{*b.Enabled, b.Type, b.NextScheduleTimeUTC, b.Hour, b.Dow, b.Dom})
}

func IsoStatus(iso *govultr.Iso) {
	defer flush()

	display(columns{"ISO ID", "STATE"})
	display(columns{iso.IsoID, iso.State})
}

func PlansList(plans []string) {
	defer flush()

	display(columns{"PLAN NAME"})

	if len(plans) == 0 {
		display(columns{"---"})
		return
	}

	for i := range plans {
		display(columns{plans[i]})
	}
}

func ReverseIpv6(rip []govultr.ReverseIP) {
	display(columns{"IP", "REVERSE"})
	for _, r := range rip {
		display(columns{r.IP, r.Reverse})
	}
	flush()
}

// InstanceVPC2List Generate a printer display of all VPC 2.0 networks attached to a given instance
func InstanceVPC2List(vpc2s []govultr.VPC2Info, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "MAC ADDRESS", "IP ADDRESS"})

	if len(vpc2s) == 0 {
		display(columns{"---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range vpc2s {
		display(columns{
			vpc2s[i].ID,
			vpc2s[i].MacAddress,
			vpc2s[i].IPAddress,
		})
	}

	Meta(meta)
}
