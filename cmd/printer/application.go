package printer

import (
	"github.com/vultr/govultr/v3"
)

func Application(apps []govultr.Application, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME", "TYPE", "VENDOR", "IMAGE ID"})

	if len(apps) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range apps {
		display(columns{
			apps[i].ID,
			apps[i].Name,
			apps[i].ShortName,
			apps[i].DeployName,
			apps[i].Type,
			apps[i].Vendor,
			apps[i].ImageID,
		})
	}

	Meta(meta)
}
