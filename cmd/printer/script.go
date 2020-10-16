package printer

import (
	"github.com/vultr/govultr"
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
	display(columns{"ID", "DATE CREATED", "DATE MODIFIED", "TYPE", "NAME"})
	display(columns{script.ID, script.DateCreated, script.DateModified, script.Type, script.Name})
	flush()
}
