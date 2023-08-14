package printer

import "github.com/vultr/govultr/v3"

func InstanceBandwidth(bandwidth *govultr.Bandwidth) {
	col := columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"}
	display(col)
	for k, b := range bandwidth.Bandwidth {
		display(columns{k, b.IncomingBytes, b.OutgoingBytes})
	}
	flush()
}

func InstanceIPV4(ip []govultr.IPv4, meta *govultr.Meta) {
	col := columns{"IP", "NETMASK", "GATEWAY", "TYPE", "REVERSE"}
	display(col)
	for _, i := range ip {
		display(columns{i.IP, i.Netmask, i.Gateway, i.Type, i.Reverse})
	}

	Meta(meta)
	flush()
}

func InstanceIPV6(ip []govultr.IPv6, meta *govultr.Meta) {
	col := columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"}
	display(col)
	for _, i := range ip {
		display(columns{i.IP, i.Network, i.NetworkSize, i.Type})
	}

	Meta(meta)
	flush()
}

func InstanceList(instance []govultr.Instance, meta *govultr.Meta) {
	col := columns{"ID", "IP", "LABEL", "OS", "STATUS", "Region", "CPU", "RAM", "DISK", "BANDWIDTH", "TAGS"}
	display(col)
	for _, s := range instance {
		display(columns{s.ID, s.MainIP, s.Label, s.Os, s.Status, s.Region, s.VCPUCount, s.RAM, s.Disk, s.AllowedBandwidth, s.Tags})
	}

	Meta(meta)
	flush()
}

func Instance(instance *govultr.Instance) {
	col := columns{"INSTANCE INFO"}
	display(col)
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
	display(columns{"FEATURES", instance.Features})
	display(columns{"TAGS", instance.Tags})

	flush()
}

func OsList(os []govultr.OS) {
	col := columns{"ID", "NAME", "ARCH", "FAMILY"}
	display(col)
	for _, o := range os {
		display(columns{o.ID, o.Name, o.Arch, o.Family})
	}
	flush()
}

func AppList(app []govultr.Application) {
	col := columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME", "TYPE", "VENDOR", "IMAGE ID"}
	display(col)
	for _, a := range app {
		display(columns{a.ID, a.Name, a.ShortName, a.DeployName, a.Type, a.Vendor, a.ImageID})
	}
	flush()
}

func BackupsGet(b *govultr.BackupSchedule) {
	col := columns{"ENABLED", "CRON TYPE", "NEXT RUN", "HOUR", "DOW", "DOM"}
	display(col)
	display(columns{*b.Enabled, b.Type, b.NextScheduleTimeUTC, b.Hour, b.Dow, b.Dom})
	flush()
}

func IsoStatus(iso *govultr.Iso) {
	col := columns{"ISO ID", "STATE"}
	display(col)
	display(columns{iso.IsoID, iso.State})
	flush()
}

func PlansList(plans []string) {
	col := columns{"PLAN NAME"}
	display(col)
	for _, p := range plans {
		display(columns{p})
	}
	flush()
}

func ReverseIpv6(rip []govultr.ReverseIP) {
	col := columns{"IP", "REVERSE"}
	display(col)
	for _, r := range rip {
		display(columns{r.IP, r.Reverse})
	}
	flush()
}

// InstanceVPC2List Generate a printer display of all VPC 2.0 networks attached to a given instance
func InstanceVPC2List(vpc2s []govultr.VPC2Info, meta *govultr.Meta) {
	display(columns{"ID", "MAC ADDRESS", "IP ADDRESS"})
	for _, r := range vpc2s {
		display(columns{r.ID, r.MacAddress, r.IPAddress})
	}

	Meta(meta)
	flush()
}
