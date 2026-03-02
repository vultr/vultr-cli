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
	return nil
}

// Data ...
func (b *BlockStoragesPrinter) Data() [][]string {
	if len(b.BlockStorages) == 0 {
		return [][]string{0: {"No block storages"}}
	}

	var data [][]string
	for i := range b.BlockStorages {
		data = append(data,
			[]string{"ID", b.BlockStorages[i].ID},
			[]string{"REGION ID", b.BlockStorages[i].Region},
			[]string{"INSTANCE ID", b.BlockStorages[i].AttachedToInstance},
			[]string{"INSTANCE IP", b.BlockStorages[i].AttachedToInstanceIP},
			[]string{"INSTANCE LABEL", b.BlockStorages[i].AttachedToInstanceLabel},
			[]string{"SIZE GB", strconv.Itoa(b.BlockStorages[i].SizeGB)},
			[]string{"STATUS", b.BlockStorages[i].Status},
			[]string{"LABEL", b.BlockStorages[i].Label},
			[]string{"BLOCK TYPE", b.BlockStorages[i].BlockType},
			[]string{"DATE CREATED", b.BlockStorages[i].DateCreated},
			[]string{"MONTHLY COST", fmt.Sprintf("$%v", b.BlockStorages[i].Cost)},
			[]string{"PENDING CHARGES", fmt.Sprintf("$%v", b.BlockStorages[i].PendingCharges)},
			[]string{"MOUNT ID", b.BlockStorages[i].MountID},
			[]string{"OS ID", b.BlockStorages[i].MountID},
			[]string{"SNAPSHOT ID", b.BlockStorages[i].MountID},
			[]string{"BOOTABLE", strconv.FormatBool(b.BlockStorages[i].Bootable)},
		)

		data = append(data, []string{"---------------------------"})
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
	return nil
}

// Data ...
func (b *BlockStoragePrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", b.BlockStorage.ID},
		[]string{"REGION ID", b.BlockStorage.Region},
		[]string{"INSTANCE ID", b.BlockStorage.AttachedToInstance},
		[]string{"INSTANCE IP", b.BlockStorage.AttachedToInstanceIP},
		[]string{"INSTANCE LABEL", b.BlockStorage.AttachedToInstanceLabel},
		[]string{"SIZE GB", strconv.Itoa(b.BlockStorage.SizeGB)},
		[]string{"STATUS", b.BlockStorage.Status},
		[]string{"LABEL", b.BlockStorage.Label},
		[]string{"BLOCK TYPE", b.BlockStorage.BlockType},
		[]string{"DATE CREATED", b.BlockStorage.DateCreated},
		[]string{"MONTHLY COST", fmt.Sprintf("$%v", b.BlockStorage.Cost)},
		[]string{"PENDING CHARGES", fmt.Sprintf("$%v", b.BlockStorage.PendingCharges)},
		[]string{"MOUNT ID", b.BlockStorage.MountID},
		[]string{"OS ID", b.BlockStorage.MountID},
		[]string{"SNAPSHOT ID", b.BlockStorage.MountID},
		[]string{"BOOTABLE", strconv.FormatBool(b.BlockStorage.Bootable)},
	)

	return data
}

// Paging ...
func (b *BlockStoragePrinter) Paging() [][]string {
	return nil
}
