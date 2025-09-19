package logs

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// LogsPrinter ...
type LogsPrinter struct {
	Logs []govultr.Log     `json:"logs"`
	Meta *govultr.LogsMeta `json:"meta"`
}

// JSON ...
func (l *LogsPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LogsPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LogsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (l *LogsPrinter) Data() [][]string {
	if len(l.Logs) == 0 {
		return [][]string{0: {"No logs returned"}}
	}

	var data [][]string
	for i := range l.Logs {
		data = append(data,
			[]string{"---------------------------"},
			[]string{"UUID", l.Logs[i].ResourceID},
			[]string{"TYPE", l.Logs[i].ResourceType},
			[]string{"LEVEL", l.Logs[i].Level},
			[]string{"MESSAGE", l.Logs[i].Message},
			[]string{"TIMESTAMP", l.Logs[i].Timestamp},
			[]string{" "},
			[]string{"USER ID", l.Logs[i].Metadata.UserID},
			[]string{"USER NAME", l.Logs[i].Metadata.UserName},
			[]string{"IP", l.Logs[i].Metadata.IPAddress},
			[]string{"HTTP CODE", strconv.Itoa(l.Logs[i].Metadata.HTTPStatusCode)},
			[]string{"HTTP METHOD", l.Logs[i].Metadata.Method},
			[]string{"REQUEST PATH", l.Logs[i].Metadata.RequestPath},
			[]string{"REQUEST BODY", l.Logs[i].Metadata.RequestBody},
			[]string{"PARAMETERS", l.Logs[i].Metadata.QueryParameters},
		)
	}

	return data
}

// Paging ...
func (l *LogsPrinter) Paging() [][]string {
	if l.Meta.TotalCount == 0 {
		return nil
	}

	var data [][]string
	data = append(data,
		[]string{"======================================"},
		[]string{"CONTINUE TIME", l.Meta.ContinueTime},
		[]string{"RETURNED COUNT", strconv.Itoa(l.Meta.ReturnedCount)},
		[]string{"UNRETURNED COUNT", strconv.Itoa(l.Meta.UnreturnedCount)},
		[]string{"TOTAL COUNT", strconv.Itoa(l.Meta.TotalCount)},
	)

	return data
}
