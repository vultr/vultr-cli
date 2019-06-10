package printer

import "github.com/vultr/govultr"

func ServerBandwidth(bandwidth []map[string]string) {
	col := columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"}
	display(col)
	for _, b := range bandwidth {
		display(columns{b["date"], b["incoming"], b["outgoing"]})
	}
	flush()
}

func ServerIPV4(ip []govultr.IPV4) {
	col := columns{"IP", "NETMASK", "GATEWAY", "TYPE", "REVERSE"}
	display(col)
	for _, i := range ip {
		display(columns{i.IP, i.Netmask, i.Gateway, i.Type, i.Reverse})
	}
	flush()
}

func ServerIPV6(ip []govultr.IPV6) {
	col := columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"}
	display(col)
	for _, i := range ip {
		display(columns{i.IP, i.Network, i.NetworkSize, i.Type})
	}
	flush()
}

func ServerList(server []govultr.Server) {
	col := columns{"ID", "IP", "LABEL", "OS", "STATUS", "Region", "CPU", "RAM", "DISK", "BANDWIDTH", "COST"}
	display(col)
	for _, s := range server {
		display(columns{s.VpsID, s.MainIP, s.Label, s.Os, s.Status, s.RegionID, s.VPSCpus, s.RAM, s.Disk, s.CurrentBandwidth, s.Cost})
	}
	flush()
}

func ServerInfo(server *govultr.Server) {
	col := columns{"SERVER INFO"}
	display(col)
	display(columns{"Instance ID", server.VpsID})
	display(columns{"Os", server.Os})
	display(columns{"Ram", server.RAM})
	display(columns{"Disk", server.Disk})
	display(columns{"Main IP", server.MainIP})
	display(columns{"VCPUS", server.VpsID})
	display(columns{"RegionID", server.RegionID})
	display(columns{"Date Created", server.Created})
	display(columns{"Pending Charges", server.PendingCharges})
	display(columns{"Status", server.Status})
	display(columns{"Monthly Cost", server.Cost})
	display(columns{"Current Bandwidth", server.CurrentBandwidth})
	display(columns{"Allowed Bandwidth", server.Cost})
	display(columns{"Netmask V4", server.NetmaskV4})
	display(columns{"Gateway V4", server.GatewayV4})
	display(columns{"Power Status", server.Status})
	display(columns{"Server State", server.ServerState})
	display(columns{"Plan", server.PlanID})
	display(columns{"Label", server.Label})
	display(columns{"Internal IP", server.InternalIP})
	display(columns{"KVM URL", server.KVMUrl})
	display(columns{"Auto Backup", server.AutoBackups})
	display(columns{"Tag", server.Tag})
	display(columns{"OsID", server.OsID})
	display(columns{"AppID", server.AppID})
	display(columns{"Firewall Group ID", server.FirewallGroupID})
	display(columns{"V6 Networks", server.V6Networks})
	flush()
}

func OsList(os []govultr.OS) {
	col := columns{"ID", "NAME", "ARCH", "FAMILY"}
	display(col)
	for _, o := range os {
		display(columns{o.OsID, o.Name, o.Arch, o.Family})
	}
	flush()
}

func AppList(app []govultr.Application) {
	col := columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME"}
	display(col)
	for _, a := range app {
		display(columns{a.AppID, a.Name, a.ShortName, a.DeployName})
	}
	flush()
}

func ServerAppInfo(app *govultr.AppInfo) {
	col := columns{"APP INFO"}
	display(col)
	display(columns{app.AppInfo})
	flush()
}

func BackupsGet(b *govultr.BackupSchedule) {
	col := columns{"ENABLED", "CRON TYPE", "NEXT RUN", "HOUR", "DOW", "DOM"}
	display(col)
	display(columns{b.Enabled, b.CronType, b.NextRun, b.Hour, b.Dow, b.Dom})
	flush()
}

func IsoStatus(iso *govultr.ServerIso) {
	col := columns{"ISO ID", "STATE"}
	display(col)
	display(columns{iso.IsoID, iso.State})
	flush()
}
