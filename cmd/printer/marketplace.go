package printer

import (
	"github.com/vultr/govultr/v3"
)

// MarketplaceAppVariableList will generate a printer display of user-supplied variables for a Vultr Marketplace app
func MarketplaceAppVariableList(appVariables []govultr.MarketplaceAppVariable) {
	defer flush()

	if len(appVariables) == 0 {
		displayString("This app contains no user-supplied variables")
		return
	}

	for p := range appVariables {
		display(columns{"NAME", appVariables[p].Name})
		display(columns{"DESCRIPTION", appVariables[p].Description})
		display(columns{"TYPE", *appVariables[p].Required})
		display(columns{"---------------------------"})
	}
}
