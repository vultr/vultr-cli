package printer

import (
	"fmt"

	"github.com/vultr/govultr/v3"
)

func BlockStorage(bs []govultr.BlockStorage, meta *govultr.Meta) {
	col := columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "BLOCK TYPE", "DATE CREATED", "MONTHLY COST", "MOUNT ID"}
	display(col)
	for _, b := range bs {
		cost := fmt.Sprintf("$%v", b.Cost)
		display(columns{b.ID, b.Region, b.AttachedToInstance, b.SizeGB, b.Status, b.Label, b.BlockType, b.DateCreated, cost, b.MountID})
	}

	Meta(meta)
	flush()
}

func SingleBlockStorage(b *govultr.BlockStorage) {
	col := columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "BLOCK TYPE", "DATE CREATED", "MONTHLY COST", "MOUNT ID"}
	display(col)
	cost := fmt.Sprintf("$%v", b.Cost)
	display(columns{b.ID, b.Region, b.AttachedToInstance, b.SizeGB, b.Status, b.Label, b.BlockType, b.DateCreated, cost, b.MountID})
	flush()
}
