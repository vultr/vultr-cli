package users

import (
	"encoding/json"
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// SSHKeysPrinter ...
type UsersPrinter struct {
	Users []govultr.User `json:"users"`
	Meta  *govultr.Meta
}

// JSON ...
func (u *UsersPrinter) JSON() []byte {
	js, err := json.MarshalIndent(u, "", printer.JSONIndent)
	if err != nil {
		panic("move this into byte")
	}

	return js
}

// YAML ...
func (u *UsersPrinter) YAML() []byte {
	yam, err := yaml.Marshal(u)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (u *UsersPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"EMAIL",
		"API",
		"ACL",
	}}
}

// Data ...
func (u *UsersPrinter) Data() [][]string {
	data := [][]string{}
	for i := range u.Users {
		data = append(data, []string{
			u.Users[i].ID,
			u.Users[i].Name,
			u.Users[i].Email,
			strconv.FormatBool(*u.Users[i].APIEnabled),
			printer.ArrayOfStringsToString(u.Users[i].ACL),
		})
	}
	return data
}

// Paging ...
func (u *UsersPrinter) Paging() [][]string {
	return printer.NewPaging(u.Meta.Total, &u.Meta.Links.Next, &u.Meta.Links.Prev).Compose()
}
