package sshkeys

import (
	"encoding/json"

	"github.com/go-yaml/yaml"
	"github.com/vultr/govultr/v2"
)

type SSHKeys struct {
	SSHKeys []govultr.SSHKey `json:"ssh_keys"`
	Meta    *govultr.Meta
}

func (s *SSHKeys) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (s *SSHKeys) Yaml() []byte {
	yam, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (s *SSHKeys) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "DATE CREATED", "NAME", "KEY"}}
}

func (s *SSHKeys) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, v := range s.SSHKeys {
		data[k] = []interface{}{v.ID, v.DateCreated, v.Name, v.SSHKey}
	}
	return data
}

func (s *SSHKeys) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {s.Meta.Total, s.Meta.Links.Next, s.Meta.Links.Prev},
	}
}

type SSHKey struct {
	SSHKey *govultr.SSHKey `json:"ssh_key"`
}

func (s *SSHKey) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (s *SSHKey) Yaml() []byte {
	yam, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (s *SSHKey) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "DATE CREATED", "NAME", "KEY"}}
}

func (s *SSHKey) Data() map[int][]interface{} {
	return map[int][]interface{}{0: {s.SSHKey.ID, s.SSHKey.DateCreated, s.SSHKey.Name, s.SSHKey.SSHKey}}
}

func (s *SSHKey) Paging() map[int][]interface{} {
	return nil
}
