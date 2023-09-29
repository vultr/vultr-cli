package printer

import (
	"github.com/vultr/govultr/v3"
)

func BareMetal(b *govultr.BareMetalServer) {
	defer flush()

	display(columns{"ID", "IP", "TAG", "MAC ADDRESS", "LABEL", "OS", "STATUS", "REGION", "CPU", "RAM", "DISK", "FEATURES", "TAGS"})
	display(columns{b.ID, b.MainIP, b.Tag, b.MacAddress, b.Label, b.Os, b.Status, b.Region, b.CPUCount, b.RAM, b.Disk, b.Features, b.Tags}) //nolint:all
}

func BareMetalList(bms []govultr.BareMetalServer, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "IP", "TAG", "MAC ADDRESS", "LABEL", "OS", "STATUS", "REGION", "CPU", "RAM", "DISK", "FEATURES", "TAGS"})

	if len(bms) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range bms {
		display(columns{
			bms[i].ID,
			bms[i].MainIP,
			bms[i].MacAddress,
			bms[i].Label,
			bms[i].Os,
			bms[i].Status,
			bms[i].Region,
			bms[i].CPUCount,
			bms[i].RAM,
			bms[i].Disk,
			bms[i].Features,
			bms[i].Tags,
		})
	}

	Meta(meta)
}

func BareMetalBandwidth(bw *govultr.Bandwidth) {
	defer flush()

	display(columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"})
	if len(bw.Bandwidth) == 0 {
		display(columns{"---", "---", "---"})
		return
	}

	for i := range bw.Bandwidth {
		display(columns{
			i,
			bw.Bandwidth[i].IncomingBytes,
			bw.Bandwidth[i].OutgoingBytes,
		})
	}
}

func BareMetalIPV4Info(info []govultr.IPv4, meta *govultr.Meta) {
	defer flush()

	display(columns{"IP", "NETMASK", "GATEWAY", "TYPE"})

	if len(info) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range info {
		display(columns{
			info[i].IP,
			info[i].Netmask,
			info[i].Gateway,
			info[i].Type,
		})
	}

	Meta(meta)
}

func BareMetalIPV6Info(info []govultr.IPv6, meta *govultr.Meta) {
	defer flush()

	display(columns{"IP", "NETWORK", "NETWORK SIZE", "TYPE"})

	if len(info) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range info {
		display(columns{
			info[i].IP,
			info[i].Network,
			info[i].NetworkSize,
			info[i].Type,
		})
	}

	Meta(meta)
}

func BareMetalVNCUrl(vnc *govultr.VNCUrl) {
	defer flush()

	display(columns{"VNC URL"})
	display(columns{vnc.URL})
}

// BareMetalVPC2List Generate a printer display of all VPC 2.0 networks attached to a given server
func BareMetalVPC2List(vpc2s []govultr.VPC2Info) {
	defer flush()

	display(columns{"ID", "MAC ADDRESS", "IP ADDRESS"})

	if len(vpc2s) == 0 {
		display(columns{"---", "---", "---"})
		return
	}

	for i := range vpc2s {
		display(columns{
			vpc2s[i].ID,
			vpc2s[i].MacAddress,
			vpc2s[i].IPAddress,
		})
	}
}
