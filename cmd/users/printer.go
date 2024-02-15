package users

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// UsersPrinter ...
type UsersPrinter struct {
	Users []govultr.User `json:"users"`
	Meta  *govultr.Meta  `json:"meta"`
}

// JSON ...
func (u *UsersPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UsersPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
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

// ======================================

// UserPrinter ...
type UserPrinter struct {
	User govultr.User `json:"user"`
}

// JSON ...
func (u *UserPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UserPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UserPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"EMAIL",
		"API",
		"ACL",
	}}
}

// Data ...
func (u *UserPrinter) Data() [][]string {
	return [][]string{0: {
		u.User.ID,
		u.User.Name,
		u.User.Email,
		strconv.FormatBool(*u.User.APIEnabled),
		printer.ArrayOfStringsToString(u.User.ACL),
	}}
}

// Paging ...
func (u *UserPrinter) Paging() [][]string {
	return nil
}
