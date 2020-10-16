package printer

import (
	"fmt"

	"github.com/vultr/govultr"
)

func BlockStorage(bs []govultr.BlockStorage, meta *govultr.Meta) {
	col := columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "DATE CREATED", "MONTHLY COST"}
	display(col)
	for _, b := range bs {
		cost := fmt.Sprintf("$%v", b.Cost)
		display(columns{b.ID, b.Region, b.AttachedToInstance, b.SizeGB, b.Status, b.Label, b.DateCreated, cost})
	}

	Meta(meta)
	flush()
}

func SingleBlockStorage(b *govultr.BlockStorage) {
	col := columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "DATE CREATED", "MONTHLY COST"}
	display(col)
	cost := fmt.Sprintf("$%v", b.Cost)
	display(columns{b.ID, b.Region, b.AttachedToInstance, b.SizeGB, b.Status, b.Label, b.DateCreated, cost})
	flush()
}
