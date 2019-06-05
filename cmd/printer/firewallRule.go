package printer

import (
	"github.com/vultr/govultr"
)

func FirewallRule(fwr []govultr.FirewallRule) {
	col := columns{"RULE NUMBER", "ACTION", "PROTOCOL", "PORT", "NETWORK", "NOTES"}
	display(col)
	for _, f := range fwr {
		display(columns{f.RuleNumber, f.Action, f.Protocol, f.Port, f.Network, f.Notes})
	}
	flush()
}
