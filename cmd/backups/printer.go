package backups

import (
	"encoding/json"

	"github.com/go-yaml/yaml"
	"github.com/vultr/govultr/v2"
)

// BackupsPrinter ...
type BackupsPrinter struct {
	Backups []govultr.Backup `json:"backups"`
	Meta    *govultr.Meta
}

// JSON ...
func (b *BackupsPrinter) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

// Yaml ...
func (b *BackupsPrinter) Yaml() []byte {
	yam, err := yaml.Marshal(b)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (b *BackupsPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "DATE CREATED", "NAME", "KEY"}}
}

// Data ...
func (b *BackupsPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, v := range b.Backups {
		data[k] = []interface{}{v.ID, v.DateCreated, v.Description, v.Size, v.Status}
	}
	return data
}

// Paging ...
func (b *BackupsPrinter) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {b.Meta.Total, b.Meta.Links.Next, b.Meta.Links.Prev},
	}
}

// BackupPrinter ...
type BackupPrinter struct {
	Backup *govultr.Backup `json:"backup"`
}

// JSON ...
func (b BackupPrinter) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

// Yaml ...
func (b BackupPrinter) Yaml() []byte {
	yam, err := yaml.Marshal(b)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (b BackupPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "DATE CREATED", "DESCRIPTION", "SIZE", "STATUS"}}
}

// Data ...
func (b BackupPrinter) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {b.Backup.ID, b.Backup.DateCreated, b.Backup.Description, b.Backup.Size, b.Backup.Status}}
}

// Paging ...
func (b BackupPrinter) Paging() map[int][]interface{} {
	return nil
}
