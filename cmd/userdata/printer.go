// Package userdata provides a printer for server user data
package userdata

import (
	"encoding/base64"
	"fmt"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// UserDataPrinter ...
type UserDataPrinter struct {
	UserData govultr.UserData `json:"user_data"`
}

// JSON ...
func (u *UserDataPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UserDataPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UserDataPrinter) Columns() [][]string {
	return [][]string{0: {
		"USERDATA",
	}}
}

// Data ...
func (u *UserDataPrinter) Data() [][]string {
	ud, err := base64.StdEncoding.DecodeString(u.UserData.Data)
	if err != nil {
		printer.Error(fmt.Errorf("error decoding base64 user data : %v", err))
	}

	return [][]string{0: {
		string(ud),
	}}
}

// Paging ...
func (u *UserDataPrinter) Paging() [][]string {
	return nil
}
