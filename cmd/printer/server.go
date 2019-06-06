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
