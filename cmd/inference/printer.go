package inference

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// InferenceSubsPrinter ...
type InferenceSubsPrinter struct {
	InferenceSubs []govultr.Inference `json:"subscriptions"`
}

// JSON ...
func (inf *InferenceSubsPrinter) JSON() []byte {
	return printer.MarshalObject(inf, "json")
}

// YAML ...
func (inf *InferenceSubsPrinter) YAML() []byte {
	return printer.MarshalObject(inf, "yaml")
}

// Columns ...
func (inf *InferenceSubsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"LABEL",
		"API KEY",
	}}
}

// Data ...
func (inf *InferenceSubsPrinter) Data() [][]string {
	if len(inf.InferenceSubs) == 0 {
		return [][]string{0: {"---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range inf.InferenceSubs {
		data = append(data, []string{
			inf.InferenceSubs[i].ID,
			inf.InferenceSubs[i].DateCreated,
			inf.InferenceSubs[i].Label,
			inf.InferenceSubs[i].APIKey,
		})
	}

	return data
}

// Paging ...
func (inf *InferenceSubsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// InferenceSubPrinter ...
type InferenceSubPrinter struct {
	InferenceSub *govultr.Inference `json:"subscription"`
}

// JSON ...
func (inf *InferenceSubPrinter) JSON() []byte {
	return printer.MarshalObject(inf, "json")
}

// YAML ...
func (inf *InferenceSubPrinter) YAML() []byte {
	return printer.MarshalObject(inf, "yaml")
}

// Columns ...
func (inf *InferenceSubPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"LABEL",
		"API KEY",
	}}
}

// Data ...
func (inf *InferenceSubPrinter) Data() [][]string {
	var data [][]string
	data = append(data, []string{
		inf.InferenceSub.ID,
		inf.InferenceSub.DateCreated,
		inf.InferenceSub.Label,
		inf.InferenceSub.APIKey,
	})

	return data
}

// Paging ...
func (inf *InferenceSubPrinter) Paging() [][]string {
	return nil
}

// ======================================

// UsagePrinter ...
type UsagePrinter struct {
	Usage *govultr.InferenceUsage `json:"usage"`
}

// JSON ...
func (u *UsagePrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UsagePrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UsagePrinter) Columns() [][]string {
	return nil
}

// Data ...
func (u *UsagePrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"CHAT USAGE"},
		[]string{"CURRENT TOKENS", strconv.FormatInt(int64(u.Usage.Chat.CurrentTokens), 10)},
		[]string{"MONTHLY ALLOTMENT", strconv.FormatInt(int64(u.Usage.Chat.MonthlyAllotment), 10)},
		[]string{"OVERAGE", strconv.FormatInt(int64(u.Usage.Chat.Overage), 10)},
		[]string{" "},
		[]string{"AUDIO USAGE"},
		[]string{"TTS CHARACTERS", strconv.FormatInt(int64(u.Usage.Audio.TTSCharacters), 10)},
		[]string{"TTS (SM) CHARACTERS", strconv.FormatInt(int64(u.Usage.Audio.TTSSMCharacters), 10)},
	)

	return data
}

// Paging ...
func (u *UsagePrinter) Paging() [][]string {
	return nil
}
