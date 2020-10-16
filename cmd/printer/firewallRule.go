package printer

import (
	"github.com/vultr/govultr"
)

func FirewallRules(fwr []govultr.FirewallRule, meta *govultr.Meta) {
	col := columns{"RULE NUMBER", "ACTION", "PROTOCOL", "PORT", "NETWORK", "NOTES"}
	display(col)
	for _, f := range fwr {
		display(columns{f.ID, f.Action, f.Protocol, f.Port, f.Subnet, f.Notes})
	}

	Meta(meta)
	flush()
}

func FirewallRule(fwr *govultr.FirewallRule) {
	col := columns{"RULE NUMBER", "ACTION", "PROTOCOL", "PORT", "NETWORK", "NOTES"}
	display(col)

	display(columns{fwr.ID, fwr.Action, fwr.Protocol, fwr.Port, fwr.Subnet, fwr.Notes})
	flush()
}
