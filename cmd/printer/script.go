package printer

import (
	"github.com/vultr/govultr"
)

func Script(script []govultr.StartupScript) {
	col := columns{"SCRIPTID", "DATE CREATED", "DATE MODIFIED", "TYPE", "NAME"}
	display(col)
	for _, s := range script {
		display(columns{s.ScriptID, s.DateCreated, s.DateModified, s.Type, s.Name})
	}
	flush()
}
