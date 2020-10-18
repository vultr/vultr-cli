package printer

import (
	"github.com/vultr/govultr/v2"
)

func Snapshot(snapshot *govultr.Snapshot) {
	display(columns{"ID", "DATE CREATED", "SIZE", "STATUS", "OSID", "APPID", "DESCRIPTION"})
	display(columns{snapshot.ID, snapshot.DateCreated, snapshot.Size, snapshot.Status, snapshot.OsID, snapshot.AppID, snapshot.Description})

	flush()
}

func Snapshots(snapshot []govultr.Snapshot, meta *govultr.Meta) {
	col := columns{"ID", "DATE CREATED", "SIZE", "STATUS", "OSID", "APPID", "DESCRIPTION"}
	display(col)
	for _, s := range snapshot {
		display(columns{s.ID, s.DateCreated, s.Size, s.Status, s.OsID, s.AppID, s.Description})
	}

	Meta(meta)
	flush()
}
