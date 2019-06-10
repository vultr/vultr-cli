package printer

import (
	"github.com/vultr/govultr"
)

func BareMetal(b *govultr.BareMetalServer) {
	display(columns{"SERVER INFO", ""})
	display(columns{"ID", b.BareMetalServerID})
	display(columns{"OS", b.Os})
	display(columns{"RAM", b.RAM})
	display(columns{"Disk", b.Disk})
	display(columns{"Main IP", b.MainIP})
	display(columns{"CPU Count", b.CPUCount})
	display(columns{"Location", b.Location})
	display(columns{"Region ID", b.RegionID})
	display(columns{"Date Created", b.DateCreated})
	display(columns{"Status", b.Status})
	display(columns{"Netmask V4", b.NetmaskV4})
	display(columns{"Gateway V4", b.GatewayV4})
	display(columns{"Plan", b.BareMetalPlanID})
	display(columns{"V6 Networks", b.V6Networks})
	display(columns{"Label", b.Label})
	display(columns{"Tag", b.Tag})
	display(columns{"OSID", b.OsID})
	display(columns{"AppID", b.AppID})

	flush()
}

func BareMetalList(bms []govultr.BareMetalServer) {
	col := columns{"ID", "IP", "TAG", "LABEL", "OS", "STATUS", "REGION", "CPU", "RAM", "DISK"}
	display(col)
	for _, b := range bms {
		display(columns{b.BareMetalServerID, b.MainIP, b.Tag, b.Label, b.Os, b.Status, b.RegionID, b.CPUCount, b.RAM, b.Disk})
	}
	flush()
}

func BareMetalAppInfo(app *govultr.AppInfo) {
	display(columns{"APP INFO"})
	display(columns{app.AppInfo})
	flush()
}

func BareMetalBandwidth(bw []map[string]string) {
	display(columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"})
	for _, b := range bw {
		display(columns{b["date"], b["incoming"], b["outgoing"]})
	}
	flush()
}

func BareMetalIPV4Info(info []govultr.BareMetalServerIPV4) {
	display(columns{"IP", "NETMASK", "GATEWAY", "TYPE"})
	for _, i := range info {
		display(columns{i.IP, i.Netmask, i.Gateway, i.Type})
	}
	flush()
}

func BareMetalIPV6Info(info []govultr.BareMetalServerIPV6) {
	display(columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"})
	for _, i := range info {
		display(columns{i.IP, i.Network, i.NetworkSize, i.Type})
	}
	flush()
}
