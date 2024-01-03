package sshkeys

import (
	"encoding/json"

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

// Yaml ...
func (s *SSHKeysPrinter) Yaml() []byte {
	yam, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (s *SSHKeysPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "DATE CREATED", "NAME", "KEY"}}
}

// Data ...
func (s *SSHKeysPrinter) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, v := range s.SSHKeys {
		data[k] = []interface{}{v.ID, v.DateCreated, v.Name, v.SSHKey}
	}
	return data
}

// Paging ...
func (s *SSHKeysPrinter) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {s.Meta.Total, s.Meta.Links.Next, s.Meta.Links.Prev},
	}
}

// SSHKeyPrinter ...
type SSHKeyPrinter struct {
	SSHKey *govultr.SSHKey `json:"ssh_key"`
}

// JSON ...
func (s SSHKeyPrinter) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

// Yaml ...
func (s SSHKeyPrinter) Yaml() []byte {
	yam, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

// Columns ...
func (s SSHKeyPrinter) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "DATE CREATED", "NAME", "KEY"}}
}

// Data ...
func (s SSHKeyPrinter) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {s.SSHKey.ID, s.SSHKey.DateCreated, s.SSHKey.Name, s.SSHKey.SSHKey}}
}

// Paging ...
func (s SSHKeyPrinter) Paging() map[int][]interface{} {
	return nil
}
