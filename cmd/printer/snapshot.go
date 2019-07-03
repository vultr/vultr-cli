package printer

import (
	"github.com/vultr/govultr"
)

func Snapshot(snapshot []govultr.Snapshot) {
	col := columns{"SNAPSHOTID", "DATE CREATED", "SIZE", "STATUS", "OSID", "APPID", "DESCRIPTION"}
	display(col)
	for _, s := range snapshot {
		display(columns{s.SnapshotID, s.DateCreated, s.Size, s.Status, s.OsID, s.AppID, s.Description})
	}
	flush()
}
