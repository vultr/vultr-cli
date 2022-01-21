package printer

import (
	"github.com/vultr/govultr/v2"
)

func ScriptList(script []govultr.StartupScript, meta *govultr.Meta) {
	col := columns{"ID", "DATE CREATED", "DATE MODIFIED", "TYPE", "NAME"}
	display(col)
	for _, s := range script {
		display(columns{s.ID, s.DateCreated, s.DateModified, s.Type, s.Name})
	}

	Meta(meta)
	flush()
}

func Script(script *govultr.StartupScript) {
	display(columns{"ID", script.ID})
	display(columns{"DATE CREATED", script.DateCreated})
	display(columns{"DATE MODIFIED", script.DateModified})
	display(columns{"TYPE", script.Type})
	display(columns{"NAME", script.Name})
	display(columns{"SCRIPT", script.Script})

	flush()
}
