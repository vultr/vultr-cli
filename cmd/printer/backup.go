package printer

import (
	"github.com/vultr/govultr/v3"
)

func Backups(bs []govultr.Backup, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "DESCRIPTION", "SIZE", "STATUS"})

	if len(bs) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range bs {
		display(columns{
			bs[i].ID,
			bs[i].DateCreated,
			bs[i].Description,
			bs[i].Size,
			bs[i].Status,
		})
	}

	Meta(meta)
}

func Backup(bs *govultr.Backup) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "DESCRIPTION", "SIZE", "STATUS"})
	display(columns{bs.ID, bs.DateCreated, bs.Description, bs.Size, bs.Status})
}
