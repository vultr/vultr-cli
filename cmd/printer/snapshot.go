package printer

import (
	"github.com/vultr/govultr/v3"
)

func Snapshots(snapshot []govultr.Snapshot, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "SIZE", "COMPRESSED SIZE", "STATUS", "OSID", "APPID", "DESCRIPTION"})

	if len(snapshot) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range snapshot {
		display(columns{
			snapshot[i].ID,
			snapshot[i].DateCreated,
			snapshot[i].Size,
			snapshot[i].CompressedSize,
			snapshot[i].Status,
			snapshot[i].OsID,
			snapshot[i].AppID,
			snapshot[i].Description,
		})
	}

	Meta(meta)
}

func Snapshot(snapshot *govultr.Snapshot) {
	defer flush()
	display(columns{"ID", "DATE CREATED", "SIZE", "COMPRESSED SIZE", "STATUS", "OSID", "APPID", "DESCRIPTION"})
	display(columns{
		snapshot.ID,
		snapshot.DateCreated,
		snapshot.Size,
		snapshot.CompressedSize,
		snapshot.Status,
		snapshot.OsID,
		snapshot.AppID,
		snapshot.Description,
	})
}
