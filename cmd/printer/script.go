package printer

import (
	"github.com/vultr/govultr/v3"
)

func ScriptList(script []govultr.StartupScript, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "DATE MODIFIED", "TYPE", "NAME"})

	if len(script) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range script {
		display(columns{
			script[i].ID,
			script[i].DateCreated,
			script[i].DateModified,
			script[i].Type,
			script[i].Name,
		})
	}

	Meta(meta)
}

func Script(script *govultr.StartupScript) {
	defer flush()

	display(columns{"ID", script.ID})
	display(columns{"DATE CREATED", script.DateCreated})
	display(columns{"DATE MODIFIED", script.DateModified})
	display(columns{"TYPE", script.Type})
	display(columns{"NAME", script.Name})
	display(columns{"SCRIPT", script.Script})
}
