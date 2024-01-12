package backups

import (
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// BackupsPrinter ...
type BackupsPrinter struct {
	Backups []govultr.Backup `json:"backups"`
	Meta    *govultr.Meta
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
func (b *BackupsPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"ID",
		"DATE CREATED",
		"DESCRIPTION",
		"SIZE",
		"STATUS",
	}}
}

// Data ...
func (b *BackupsPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, v := range b.Backups {
		data[k] = []interface{}{
			v.ID,
			v.DateCreated,
			v.Description,
			v.Size,
			v.Status,
		}
	}
	return data
}

// Paging ...
func (b *BackupsPrinter) Paging() map[int][]interface{} {
	return printer.NewPaging(b.Meta.Total, &b.Meta.Links.Next, &b.Meta.Links.Prev).Compose()
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
func (b *BackupPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {
		"ID",
		"DATE CREATED",
		"DESCRIPTION",
		"SIZE",
		"STATUS",
	}}
}

// Data ...
func (b *BackupPrinter) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {
		b.Backup.ID,
		b.Backup.DateCreated,
		b.Backup.Description,
		b.Backup.Size,
		b.Backup.Status,
	}}
}

// Paging ...
func (b *BackupPrinter) Paging() map[int][]interface{} {
	return nil
}
