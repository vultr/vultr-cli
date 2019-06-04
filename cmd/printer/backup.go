package printer

import (
	"github.com/vultr/govultr"
)

func Backup(bs []govultr.Backup) {
	col := columns{"ID", "DATE CREATED", "DESCRIPTION", "SIZE", "STATUS"}
	display(col)
	for _, b := range bs {
		display(columns{b.BackupID, b.DateCreated, b.Description, b.Size, b.Status})
	}
	flush()
}
