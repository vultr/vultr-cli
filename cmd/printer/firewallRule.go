package printer

import (
	"fmt"

	"github.com/vultr/govultr/v2"
)

func FirewallRules(fwr []govultr.FirewallRule, meta *govultr.Meta) {
	col := columns{"RULE NUMBER", "ACTION", "TYPE", "PROTOCOL", "PORT", "NETWORK", "SOURCE", "NOTES"}
	display(col)
	for _, f := range fwr {
		display(columns{f.ID, f.Action, f.IPType, f.Protocol, f.Port, getFirewallNetwork(f.Subnet, f.SubnetSize), getFirewallSource(f.Source), f.Notes})
	}

	Meta(meta)
	flush()
}

func FirewallRule(fwr *govultr.FirewallRule) {
	col := columns{"RULE NUMBER", "ACTION", "TYPE", "PROTOCOL", "PORT", "NETWORK", "SOURCE", "NOTES"}
	display(col)
	display(columns{fwr.ID, fwr.Action, fwr.IPType, fwr.Protocol, fwr.Port, getFirewallNetwork(fwr.Subnet, fwr.SubnetSize), getFirewallSource(fwr.Source), fwr.Notes})
	flush()
}

func getFirewallSource(source string) string {
	if source == "" {
		return "anywhere"
	}
	return source
}

func getFirewallNetwork(subnet string, size int) string {
	return fmt.Sprintf("%s/%d", subnet, size)
}
