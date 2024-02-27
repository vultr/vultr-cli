package sshkeys

import (
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// SSHKeysPrinter ...
type SSHKeysPrinter struct {
	SSHKeys []govultr.SSHKey `json:"ssh_keys"`
	Meta    *govultr.Meta
}

// JSON ...
func (s *SSHKeysPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s *SSHKeysPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s *SSHKeysPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"NAME",
		"KEY",
	}}
}

// Data ...
func (s *SSHKeysPrinter) Data() [][]string {
	data := [][]string{}
	for i := range s.SSHKeys {
		data = append(data, []string{
			s.SSHKeys[i].ID,
			s.SSHKeys[i].DateCreated,
			s.SSHKeys[i].Name,
			s.SSHKeys[i].SSHKey,
		})
	}
	return data
}

// Paging ...
func (s *SSHKeysPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(s.Meta).Compose()
}

// SSHKeyPrinter ...
type SSHKeyPrinter struct {
	SSHKey *govultr.SSHKey `json:"ssh_key"`
}

// JSON ...
func (s SSHKeyPrinter) JSON() []byte {
	return printer.MarshalObject(s, "json")
}

// YAML ...
func (s SSHKeyPrinter) YAML() []byte {
	return printer.MarshalObject(s, "yaml")
}

// Columns ...
func (s SSHKeyPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"NAME",
		"KEY",
	}}
}

// Data ...
func (s SSHKeyPrinter) Data() [][]string {
	return [][]string{0: {
		s.SSHKey.ID,
		s.SSHKey.DateCreated,
		s.SSHKey.Name,
		s.SSHKey.SSHKey,
	}}
}

// Paging ...
func (s SSHKeyPrinter) Paging() [][]string {
	return nil
}
