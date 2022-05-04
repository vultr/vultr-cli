package printer

import (
	"github.com/vultr/govultr/v2"
)

func BareMetal(b *govultr.BareMetalServer) {
	col := columns{"ID", "IP", "TAG", "MAC ADDRESS", "LABEL", "OS", "STATUS", "REGION", "CPU", "RAM", "DISK", "FEATURES", "TAGS"}
	display(col)

	display(columns{b.ID, b.MainIP, b.Tag, b.MacAddress, b.Label, b.Os, b.Status, b.Region, b.CPUCount, b.RAM, b.Disk, b.Features, b.Tags})

	flush()
}

func BareMetalList(bms []govultr.BareMetalServer, meta *govultr.Meta) {
	col := columns{"ID", "IP", "TAG", "MAC ADDRESS", "LABEL", "OS", "STATUS", "REGION", "CPU", "RAM", "DISK", "FEATURES", "TAGS"}
	display(col)
	for _, b := range bms {
		display(columns{b.ID, b.MainIP, b.Tag, b.MacAddress, b.Label, b.Os, b.Status, b.Region, b.CPUCount, b.RAM, b.Disk, b.Features, b.Tags})
	}

	Meta(meta)

	flush()
}

func BareMetalBandwidth(bw *govultr.Bandwidth) {
	display(columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"})
	for k, b := range bw.Bandwidth {
		display(columns{k, b.IncomingBytes, b.OutgoingBytes})
	}
	flush()
}

func BareMetalIPV4Info(info []govultr.IPv4, meta *govultr.Meta) {
	display(columns{"IP", "NETMASK", "GATEWAY", "TYPE"})
	for _, i := range info {
		display(columns{i.IP, i.Netmask, i.Gateway, i.Type})
	}

	Meta(meta)
	flush()
}

func BareMetalIPV6Info(info []govultr.IPv6, meta *govultr.Meta) {
	display(columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"})
	for _, i := range info {
		display(columns{i.IP, i.Network, i.NetworkSize, i.Type})
	}

	Meta(meta)
	flush()
}

func BareMetalVNCUrl(vnc *govultr.VNCUrl) {
	display(columns{"VNC URL"})
	display(columns{vnc.URL})
	flush()
}
