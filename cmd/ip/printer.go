// Package ip provides printers for server network addresses
package ip

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// IPv4Printer ...
type IPv4sPrinter struct {
	IPv4s []govultr.IPv4 `json:"ipv4s"`
	Meta  *govultr.Meta  `json:"meta"`
}

// JSON ...
func (i *IPv4sPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *IPv4sPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *IPv4sPrinter) Columns() [][]string {
	return [][]string{0: {
		"IP",
		"NETMASK",
		"GATEWAY",
		"TYPE",
	}}
}

// Data ...
func (i *IPv4sPrinter) Data() [][]string {
	var data [][]string
	for j := range i.IPv4s {
		data = append(data, []string{
			i.IPv4s[j].IP,
			i.IPv4s[j].Netmask,
			i.IPv4s[j].Gateway,
			i.IPv4s[j].Type,
		})
	}

	return data
}

// Paging ...
func (i *IPv4sPrinter) Paging() [][]string {
	return printer.NewPaging(i.Meta.Total, &i.Meta.Links.Next, &i.Meta.Links.Prev).Compose()
}

// ======================================

// IPv6sPrinter ...
type IPv6sPrinter struct {
	IPv6s []govultr.IPv6 `json:"ipv6s"`
	Meta  *govultr.Meta  `json:"meta"`
}

// JSON ...
func (i *IPv6sPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *IPv6sPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *IPv6sPrinter) Columns() [][]string {
	return [][]string{0: {
		"IP",
		"NETWORK",
		"NETWORK SIZE",
		"TYPE",
	}}
}

// Data ...
func (i *IPv6sPrinter) Data() [][]string {
	var data [][]string
	for j := range i.IPv6s {
		data = append(data, []string{
			i.IPv6s[j].IP,
			i.IPv6s[j].Network,
			strconv.Itoa(i.IPv6s[j].NetworkSize),
			i.IPv6s[j].Type,
		})
	}

	return data
}

// Paging ...
func (i *IPv6sPrinter) Paging() [][]string {
	return printer.NewPaging(i.Meta.Total, &i.Meta.Links.Next, &i.Meta.Links.Prev).Compose()
}
