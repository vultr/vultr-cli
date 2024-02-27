package backups

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// BackupsPrinter ...
type BackupsPrinter struct {
	Backups []govultr.Backup `json:"backups"`
	Meta    *govultr.Meta    `json:"meta"`
}

// JSON ...
func (b *BackupsPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BackupsPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BackupsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"DESCRIPTION",
		"SIZE",
		"STATUS",
	}}
}

// Data ...
func (b *BackupsPrinter) Data() [][]string {
	data := [][]string{}
	for i := range b.Backups {
		data = append(data, []string{
			b.Backups[i].ID,
			b.Backups[i].DateCreated,
			b.Backups[i].Description,
			strconv.Itoa(b.Backups[i].Size),
			b.Backups[i].Status,
		})
	}
	return data
}

// Paging ...
func (b *BackupsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(b.Meta).Compose()
}

// ======================================

// BackupPrinter ...
type BackupPrinter struct {
	Backup *govultr.Backup `json:"backup"`
}

// JSON ...
func (b *BackupPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BackupPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BackupPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"DESCRIPTION",
		"SIZE",
		"STATUS",
	}}
}

// Data ...
func (b *BackupPrinter) Data() [][]string {
	return [][]string{0: {
		b.Backup.ID,
		b.Backup.DateCreated,
		b.Backup.Description,
		strconv.Itoa(b.Backup.Size),
		b.Backup.Status,
	}}
}

// Paging ...
func (b *BackupPrinter) Paging() [][]string {
	return nil
}
