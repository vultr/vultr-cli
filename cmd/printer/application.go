package printer

import (
	"github.com/vultr/govultr/v2"
)

func Application(apps []govultr.Application, meta *govultr.Meta) {
	col := columns{"ID", "NAME", "SHORT NAME", "DEPLOY NAME", "TYPE", "VENDOR", "IMAGE ID"}
	display(col)
	for _, a := range apps {
		display(columns{a.ID, a.Name, a.ShortName, a.DeployName, a.Type, a.Vendor, a.ImageID})
	}

	Meta(meta)
	flush()
}
