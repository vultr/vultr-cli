package printer

import (
	"github.com/vultr/govultr/v3"
)

func FirewallGroups(fwg []govultr.FirewallGroup, meta *govultr.Meta) {
	col := columns{"ID", "DATE CREATED", "DATE MODIFIED", "INSTANCE COUNT", "RULE COUNT", "MAX RULE COUNT", "DESCRIPTION"}
	display(col)
	for _, f := range fwg {
		display(columns{f.ID, f.DateCreated, f.DateModified, f.InstanceCount, f.RuleCount, f.MaxRuleCount, f.Description})
	}

	Meta(meta)
	flush()
}

func FirewallGroup(fwg *govultr.FirewallGroup) {
	col := columns{"ID", "DATE CREATED", "DATE MODIFIED", "INSTANCE COUNT", "RULE COUNT", "MAX RULE COUNT", "DESCRIPTION"}
	display(col)

	display(columns{fwg.ID, fwg.DateCreated, fwg.DateModified, fwg.InstanceCount, fwg.RuleCount, fwg.MaxRuleCount, fwg.Description})
	flush()
}
