package script

import (
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// ReservedIPsPrinter...
type ScriptsPrinter struct {
	Scripts []govultr.StartupScript `json:"startup_scripts"`
	Meta    *govultr.Meta           `json:"meta"`
}

// JSON ...
func (s *ScriptsPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *ScriptsPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *ScriptsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"DATE MODIFIED",
		"TYPE",
		"NAME",
	}}
}

// Data ...
func (s *ScriptsPrinter) Data() [][]string {
	if len(s.Scripts) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---"}}

	}

	var data [][]string
	for i := range s.Scripts {
		data = append(data, []string{
			s.Scripts[i].ID,
			s.Scripts[i].DateCreated,
			s.Scripts[i].DateModified,
			s.Scripts[i].Type,
			s.Scripts[i].Name,
		})
	}

	return data
}

// Paging ...
func (s *ScriptsPrinter) Paging() [][]string {
	return printer.NewPaging(s.Meta.Total, &s.Meta.Links.Next, &s.Meta.Links.Prev).Compose()
}

// ======================================

// ReservedIPPrinter...
type ScriptPrinter struct {
	Script *govultr.StartupScript `json:"startup_script"`
}

// JSON ...
func (s *ScriptPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *ScriptPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *ScriptPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (s *ScriptPrinter) Data() [][]string {
	return [][]string{
		0: {"ID", s.Script.ID},
		1: {"DATE CREATED", s.Script.DateCreated},
		2: {"DATE MODIFIED", s.Script.DateModified},
		3: {"TYPE", s.Script.Type},
		4: {"NAME", s.Script.Name},
		5: {"SCRIPT", s.Script.Script},
	}
}

// Paging ...
func (s *ScriptPrinter) Paging() [][]string {
	return nil
}
