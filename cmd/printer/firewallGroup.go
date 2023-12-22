package printer

import (
	"github.com/vultr/govultr/v3"
)

func FirewallGroups(fwg []govultr.FirewallGroup, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "DATE MODIFIED", "INSTANCE COUNT", "RULE COUNT", "MAX RULE COUNT", "DESCRIPTION"})

	if len(fwg) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range fwg {
		display(columns{
			fwg[i].ID,
			fwg[i].DateCreated,
			fwg[i].DateModified,
			fwg[i].InstanceCount,
			fwg[i].RuleCount,
			fwg[i].MaxRuleCount,
			fwg[i].Description,
		})
	}

	Meta(meta)
}

func FirewallGroup(fwg *govultr.FirewallGroup) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "DATE MODIFIED", "INSTANCE COUNT", "RULE COUNT", "MAX RULE COUNT", "DESCRIPTION"})
	display(columns{fwg.ID, fwg.DateCreated, fwg.DateModified, fwg.InstanceCount, fwg.RuleCount, fwg.MaxRuleCount, fwg.Description})
}
