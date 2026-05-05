package inference

import (
	"fmt"
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
		[]string{"---------------------------"},
		[]string{"CHAT USAGE (CURRENT MONTH)"},
		[]string{"---------------------------"},
	)

	if len(u.Usage.ChatByModel.CurrentMonth) == 0 {
		data = append(data,
			[]string{"None"},
			[]string{" "},
		)
	} else {
		for i := range u.Usage.ChatByModel.CurrentMonth {
			data = append(data,
				[]string{"MODEL", u.Usage.ChatByModel.CurrentMonth[i].Model},
				[]string{"TOKENS", strconv.Itoa(u.Usage.ChatByModel.CurrentMonth[i].Tokens)},
				[]string{"INPUT TOKENS", strconv.Itoa(u.Usage.ChatByModel.CurrentMonth[i].InputTokens)},
				[]string{"OUTPUT PRICE", strconv.Itoa(u.Usage.ChatByModel.CurrentMonth[i].OutputPrice)},
				[]string{"INPUT PRICE", strconv.Itoa(u.Usage.ChatByModel.CurrentMonth[i].InputPrice)},
				[]string{" "},
			)
		}
	}

	data = append(data,
		[]string{"---------------------------"},
		[]string{"CHAT USAGE (PREVIOUS MONTH)"},
		[]string{"---------------------------"},
	)

	if len(u.Usage.ChatByModel.PreviousMonth) == 0 {
		data = append(data,
			[]string{"None"},
			[]string{" "},
		)
	} else {
		for i := range u.Usage.ChatByModel.PreviousMonth {
			data = append(data,
				[]string{"MODEL", u.Usage.ChatByModel.PreviousMonth[i].Model},
				[]string{"TOKENS", strconv.Itoa(u.Usage.ChatByModel.PreviousMonth[i].Tokens)},
				[]string{"INPUT TOKENS", strconv.Itoa(u.Usage.ChatByModel.PreviousMonth[i].InputTokens)},
				[]string{"OUTPUT PRICE", strconv.Itoa(u.Usage.ChatByModel.PreviousMonth[i].OutputPrice)},
				[]string{"INPUT PRICE", strconv.Itoa(u.Usage.ChatByModel.PreviousMonth[i].InputPrice)},
				[]string{" "},
			)
		}
	}

	data = append(data,
		[]string{"---------------------------"},
		[]string{"AUDIO USAGE"},
		[]string{"---------------------------"},
		[]string{"TTS CHARACTERS", strconv.Itoa(u.Usage.Audio.TTSCharacters)},
		[]string{"TTS (SM) CHARACTERS", strconv.Itoa(u.Usage.Audio.TTSSMCharacters)},
		[]string{" "},
		[]string{"---------------------------"},
		[]string{"IMAGE USAGE"},
		[]string{"---------------------------"},
		[]string{"MEGAPIXELS", fmt.Sprintf("%g", u.Usage.Image.Megapixels)},
		[]string{"MEGAPIXELS (SM)", fmt.Sprintf("%g", u.Usage.Image.SMMegapixels)},
	)

	return data
}

// Paging ...
func (u *UsagePrinter) Paging() [][]string {
	return nil
}
