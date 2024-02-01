package snapshot

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// SnapshotsPrinter ...
type SnapshotsPrinter struct {
	Snapshots []govultr.Snapshot `json:"snapshots"`
	Meta      *govultr.Meta      `json:"meta"`
}

// JSON ...
func (s *SnapshotsPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *SnapshotsPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *SnapshotsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"SIZE",
		"COMPRESSED SIZE",
		"STATUS",
		"OSID",
		"APPID",
		"DESCRIPTION",
	}}
}

// Data ...
func (s *SnapshotsPrinter) Data() [][]string {
	if len(s.Snapshots) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---"}}

	}

	var data [][]string
	for i := range s.Snapshots {
		data = append(data, []string{
			s.Snapshots[i].ID,
			s.Snapshots[i].DateCreated,
			strconv.Itoa(s.Snapshots[i].Size),
			strconv.Itoa(s.Snapshots[i].CompressedSize),
			s.Snapshots[i].Status,
			strconv.Itoa(s.Snapshots[i].OsID),
			strconv.Itoa(s.Snapshots[i].AppID),
			s.Snapshots[i].Description,
		})
	}

	return data
}

// Paging ...
func (s *SnapshotsPrinter) Paging() [][]string {
	return printer.NewPaging(s.Meta.Total, &s.Meta.Links.Next, &s.Meta.Links.Prev).Compose()
}

// ======================================

// SnapshotPrinter ...
type SnapshotPrinter struct {
	Snapshot *govultr.Snapshot `json:"snapshot"`
}

// JSON ...
func (s *SnapshotPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *SnapshotPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *SnapshotPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"SIZE",
		"COMPRESSED SIZE",
		"STATUS",
		"OSID",
		"APPID",
		"DESCRIPTION",
	}}
}

// Data ...
func (s *SnapshotPrinter) Data() [][]string {
	return [][]string{0: {
		s.Snapshot.ID,
		s.Snapshot.DateCreated,
		strconv.Itoa(s.Snapshot.Size),
		strconv.Itoa(s.Snapshot.CompressedSize),
		s.Snapshot.Status,
		strconv.Itoa(s.Snapshot.OsID),
		strconv.Itoa(s.Snapshot.AppID),
		s.Snapshot.Description,
	}}
}

// Paging ...
func (s *SnapshotPrinter) Paging() [][]string {
	return nil
}
