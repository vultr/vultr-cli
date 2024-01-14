package sshkeys

import (
	"encoding/json"
	"strconv"

	"github.com/vultr/govultr/v3"
	"gopkg.in/yaml.v3"
)

// SSHKeysPrinter ...
type SSHKeysPrinter struct {
	SSHKeys []govultr.SSHKey `json:"ssh_keys"`
	Meta    *govultr.Meta
}

// JSON ...
func (s *SSHKeysPrinter) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

// YAML ...
func (s *SSHKeysPrinter) YAML() []byte {
	yam, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yam
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
	return [][]string{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {strconv.Itoa(s.Meta.Total), s.Meta.Links.Next, s.Meta.Links.Prev},
	}
}

// SSHKeyPrinter ...
type SSHKeyPrinter struct {
	SSHKey *govultr.SSHKey `json:"ssh_key"`
}

// JSON ...
func (s SSHKeyPrinter) JSON() []byte {
	js, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return js
}

// YAML ...
func (s SSHKeyPrinter) YAML() []byte {
	yam, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yam
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
