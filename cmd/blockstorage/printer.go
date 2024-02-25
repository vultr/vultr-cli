package blockstorage

import (
	"fmt"
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// BlockStoragesPrinter ...
type BlockStoragesPrinter struct {
	BlockStorages []govultr.BlockStorage `json:"blocks"`
	Meta          *govultr.Meta          `json:"meta"`
}

// JSON ...
func (b *BlockStoragesPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BlockStoragesPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BlockStoragesPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION ID",
		"INSTANCE ID",
		"SIZE GB",
		"STATUS",
		"LABEL",
		"BLOCK TYPE",
		"DATE CREATED",
		"MONTHLY COST",
		"MOUNT ID",
	}}
}

// Data ...
func (b *BlockStoragesPrinter) Data() [][]string {
	if len(b.BlockStorages) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range b.BlockStorages {
		data = append(data, []string{
			b.BlockStorages[i].ID,
			b.BlockStorages[i].Region,
			b.BlockStorages[i].AttachedToInstance,
			strconv.Itoa(b.BlockStorages[i].SizeGB),
			b.BlockStorages[i].Status,
			b.BlockStorages[i].Label,
			b.BlockStorages[i].BlockType,
			b.BlockStorages[i].DateCreated,
			fmt.Sprintf("$%v", b.BlockStorages[i].Cost),
			b.BlockStorages[i].MountID,
		})
	}

	return data
}

// Paging ...
func (b *BlockStoragesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(b.Meta).Compose()
}

// ======================================

// BlockStoragePrinter ...
type BlockStoragePrinter struct {
	BlockStorage *govultr.BlockStorage `json:"block"`
}

// JSON ...
func (b *BlockStoragePrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BlockStoragePrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BlockStoragePrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION ID",
		"INSTANCE ID",
		"SIZE GB",
		"STATUS",
		"LABEL",
		"BLOCK TYPE",
		"DATE CREATED",
		"MONTHLY COST",
		"MOUNT ID",
	}}
}

// Data ...
func (b *BlockStoragePrinter) Data() [][]string {
	return [][]string{0: {
		b.BlockStorage.ID,
		b.BlockStorage.Region,
		b.BlockStorage.AttachedToInstance,
		strconv.Itoa(b.BlockStorage.SizeGB),
		b.BlockStorage.Status,
		b.BlockStorage.Label,
		b.BlockStorage.BlockType,
		b.BlockStorage.DateCreated,
		fmt.Sprintf("$%v", b.BlockStorage.Cost),
		b.BlockStorage.MountID,
	}}
}

// Paging ...
func (b *BlockStoragePrinter) Paging() [][]string {
	return nil
}
