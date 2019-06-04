package printer

import (
	"github.com/vultr/govultr"
)

func Script(script []govultr.StartupScript) {
	col := columns{"SCRIPTID", "TYPE", "NAME", "SCRIPT"}
	display(col)
	for _, s := range script {
		display(columns{s.ScriptID, s.Type, s.Name, s.Script})
	}
	flush()
}
