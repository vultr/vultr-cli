package printer

import (
	"github.com/vultr/govultr"
)

func Application(apps []govultr.Application, meta *govultr.Meta) {
	col := columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME"}
	display(col)
	for _, a := range apps {
		display(columns{a.ID, a.Name, a.ShortName, a.DeployName})
	}

	Meta(meta)
	flush()
}
