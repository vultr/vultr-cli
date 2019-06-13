package printer

import (
	"fmt"
	"github.com/vultr/govultr"
)

func BlockStorage(bs []govultr.BlockStorage) {
	col := columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "DATE CREATED", "MONTHLY COST"}
	display(col)
	for _, b := range bs {
		cost := fmt.Sprintf("$%s", b.CostPerMonth)
		display(columns{b.BlockStorageID, b.RegionID, b.InstanceID, b.SizeGB, b.Status, b.Label, b.DateCreated, cost})
	}
	flush()
}

func SingleBlockStorage(b *govultr.BlockStorage) {
	col := columns{"ID", "REGION ID", "INSTANCE ID", "SIZE GB", "STATUS", "LABEL", "DATE CREATED", "MONTHLY COST"}
	display(col)
	cost := fmt.Sprintf("$%s", b.CostPerMonth)
	display(columns{b.BlockStorageID, b.RegionID, b.InstanceID, b.SizeGB, b.Status, b.Label, b.DateCreated, cost})
	flush()
}
