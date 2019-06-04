package printer

import (
	"github.com/vultr/govultr"
)

func Application(apps []govultr.Application) {
	col := columns{"APPID", "NAME", "SHORT NAME", "DEPLOY NAME"}
	display(col)
	for _, a := range apps {
		display(columns{a.AppID, a.Name, a.ShortName, a.DeployName})
	}
	flush()
}
