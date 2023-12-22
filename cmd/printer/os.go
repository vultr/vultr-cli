package printer

import (
	"github.com/vultr/govultr/v3"
)

func Os(vultrOS []govultr.OS, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "NAME", "ARCH", "FAMILY"})

	if len(vultrOS) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range vultrOS {
		display(columns{
			vultrOS[i].ID,
			vultrOS[i].Name,
			vultrOS[i].Arch,
			vultrOS[i].Family,
		})
	}

	Meta(meta)
}
