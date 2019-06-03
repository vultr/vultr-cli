package printer

import (
	"github.com/vultr/govultr"
)

func Os(vultrOS []govultr.OS) {
	col := columns{"OSID", "NAME", "ARCH", "FAMILY", "WINDOWS"}
	display(col)
	for _, os := range vultrOS {
		display(columns{os.OsID, os.Name, os.Arch, os.Family, os.Windows})
	}
	flush()
}
