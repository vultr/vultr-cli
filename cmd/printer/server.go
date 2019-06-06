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
