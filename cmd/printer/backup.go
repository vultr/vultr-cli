package printer

import (
	"github.com/vultr/govultr/v2"
)

func Backups(bs []govultr.Backup, meta *govultr.Meta) {
	col := columns{"ID", "DATE CREATED", "DESCRIPTION", "SIZE", "STATUS"}
	display(col)
	for _, b := range bs {
		display(columns{b.ID, b.DateCreated, b.Description, b.Size, b.Status})
	}

	Meta(meta)
	flush()
}

func Backup(bs *govultr.Backup) {
	col := columns{"ID", "DATE CREATED", "DESCRIPTION", "SIZE", "STATUS"}
	display(col)

	flush()
}
