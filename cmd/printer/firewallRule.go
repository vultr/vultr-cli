package printer

import (
	"fmt"

	"github.com/vultr/govultr/v3"
)

func FirewallRules(fwr []govultr.FirewallRule, meta *govultr.Meta) {
	defer flush()

	display(columns{"RULE NUMBER", "ACTION", "TYPE", "PROTOCOL", "PORT", "NETWORK", "SOURCE", "NOTES"})

	if len(fwr) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range fwr {
		display(columns{
			fwr[i].ID,
			fwr[i].Action,
			fwr[i].IPType,
			fwr[i].Protocol,
			fwr[i].Port,
			getFirewallNetwork(fwr[i].Subnet, fwr[i].SubnetSize),
			getFirewallSource(fwr[i].Source),
			fwr[i].Notes,
		})
	}

	Meta(meta)
}

func FirewallRule(fwr *govultr.FirewallRule) {
	defer flush()

	display(columns{"RULE NUMBER", "ACTION", "TYPE", "PROTOCOL", "PORT", "NETWORK", "SOURCE", "NOTES"})
	display(columns{
		fwr.ID,
		fwr.Action,
		fwr.IPType,
		fwr.Protocol,
		fwr.Port,
		getFirewallNetwork(fwr.Subnet, fwr.SubnetSize),
		getFirewallSource(fwr.Source),
		fwr.Notes,
	})
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
