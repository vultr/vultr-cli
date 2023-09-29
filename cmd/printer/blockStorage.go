package printer

import (
	"fmt"

	"github.com/vultr/govultr/v3"
)

func BlockStorage(bs []govultr.BlockStorage, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "BLOCK TYPE", "DATE CREATED", "MONTHLY COST", "MOUNT ID"})

	if len(bs) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range bs {
		cost := fmt.Sprintf("$%v", bs[i].Cost)
		display(columns{
			bs[i].ID,
			bs[i].Region,
			bs[i].AttachedToInstance,
			bs[i].SizeGB,
			bs[i].Status,
			bs[i].Label,
			bs[i].BlockType,
			bs[i].DateCreated,
			cost,
			bs[i].MountID,
		})
	}

	Meta(meta)
}

func SingleBlockStorage(b *govultr.BlockStorage) {
	defer flush()

	display(columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "BLOCK TYPE", "DATE CREATED", "MONTHLY COST", "MOUNT ID"})
	cost := fmt.Sprintf("$%v", b.Cost)
	display(columns{b.ID, b.Region, b.AttachedToInstance, b.SizeGB, b.Status, b.Label, b.BlockType, b.DateCreated, cost, b.MountID})
}
