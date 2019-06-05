package printer

import (
	"github.com/vultr/govultr"
)

func FirewallGroup(fwg []govultr.FirewallGroup) {
	col := columns{"FIREWALLGROUPID", "DATE CREATED", "DATE MODIFIED", "INSTANCE COUNT", "RULE COUNT", "MAX RULE COUNT", "DESCRIPTION"}
	display(col)
	for _, f := range fwg {
		display(columns{f.FirewallGroupID, f.DateCreated, f.DateModified, f.InstanceCount, f.RuleCount, f.MaxRuleCount, f.Description})
	}
	flush()
}
