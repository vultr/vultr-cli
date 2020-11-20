package printer

import "github.com/vultr/govultr/v2"

func ServerBandwidth(bandwidth *govultr.Bandwidth) {
	col := columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"}
	display(col)
	for k, b := range bandwidth.Bandwidth {
		display(columns{k, b.IncomingBytes, b.OutgoingBytes})
	}
	flush()
}

func ServerIPV4(ip []govultr.IPv4, meta *govultr.Meta) {
	col := columns{"IP", "NETMASK", "GATEWAY", "TYPE", "REVERSE"}
	display(col)
	for _, i := range ip {
		display(columns{i.IP, i.Netmask, i.Gateway, i.Type, i.Reverse})
	}

	Meta(meta)
	flush()
}

func ServerIPV6(ip []govultr.IPv6, meta *govultr.Meta) {
	col := columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"}
	display(col)
	for _, i := range ip {
		display(columns{i.IP, i.Network, i.NetworkSize, i.Type})
	}

	Meta(meta)
	flush()
}

// func Server(server *govultr.Instance) {
// 	display(columns{"ID", "IP", "LABEL", "OS", "STATUS", "Region", "CPU", "RAM", "DISK", "BANDWIDTH"})
// 	display(columns{server.ID, server.MainIP, server.Label, server.Os, server.Status, server.Region, server.VCPUCount, server.Ram, server.Disk, server.AllowedBandwidth})

// 	flush()
// }

func ServerList(server []govultr.Instance, meta *govultr.Meta) {
	col := columns{"ID", "IP", "LABEL", "OS", "STATUS", "Region", "CPU", "RAM", "DISK", "BANDWIDTH"}
	display(col)
	for _, s := range server {
		display(columns{s.ID, s.MainIP, s.Label, s.Os, s.Status, s.Region, s.VCPUCount, s.RAM, s.Disk, s.AllowedBandwidth})
	}

	Meta(meta)
	flush()
}

func Server(server *govultr.Instance) {
	col := columns{"SERVER INFO"}
	display(col)
	display(columns{"ID", server.ID})
	display(columns{"Os", server.Os})
	display(columns{"RAM", server.RAM})
	display(columns{"DISK", server.Disk})
	display(columns{"MAIN IP", server.MainIP})
	display(columns{"VCPU CONT", server.VCPUCount})
	display(columns{"REGION", server.Region})
	display(columns{"DATE CREATED", server.DateCreated})
	display(columns{"STATUS", server.Status})
	display(columns{"ALLOWED BANDWIDTH", server.AllowedBandwidth})
	display(columns{"NETMASK V4", server.NetmaskV4})
	display(columns{"GATEWAY V4", server.GatewayV4})
	display(columns{"POWER STATUS", server.PowerStatus})
	display(columns{"SERVER STATE", server.ServerStatus})
	display(columns{"PLAN", server.Plan})
	display(columns{"LABEL", server.Label})
	display(columns{"INTERNAL IP", server.InternalIP})
	display(columns{"KVM URL", server.KVM})
	display(columns{"TAG", server.Tag})
	display(columns{"OsID", server.OsID})
	display(columns{"AppID", server.AppID})
	display(columns{"FIREWALL GROUP ID", server.FirewallGroupID})
	display(columns{"V6 MAIN IP", server.V6MainIP})
	display(columns{"V6 NETWORK", server.V6Network})
	display(columns{"V6 NETWORK SIZE", server.V6NetworkSize})
	display(columns{"FEATURES", server.Features})

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
	col := columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME"}
	display(col)
	for _, a := range app {
		display(columns{a.ID, a.Name, a.ShortName, a.DeployName})
	}
	flush()
}

func BackupsGet(b *govultr.BackupSchedule) {
	col := columns{"ENABLED", "CRON TYPE", "NEXT RUN", "HOUR", "DOW", "DOM"}
	display(col)
	display(columns{b.Enabled, b.Type, b.NextScheduleTimeUTC, b.Hour, b.Dow, b.Dom})
	flush()
}

func IsoStatus(iso *govultr.Iso) {
	col := columns{"ISO ID", "STATE"}
	display(col)
	display(columns{iso.IsoID, iso.State})
	flush()
}

func PlansList(plans []int) {
	col := columns{"PLAN ID"}
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
