package account

import (
	"encoding/json"
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"gopkg.in/yaml.v3"
)

// AccountPrinter ...
type AccountPrinter struct {
	Account *govultr.Account `json:"account"`
}

// JSON ...
func (s *AccountPrinter) JSON() []byte {
	js, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return js
}

// YAML ...
func (s *AccountPrinter) YAML() []byte {
	yml, err := yaml.Marshal(s)
	if err != nil {
		panic("move this into byte")
	}
	return yml
}

// Columns ...
func (a *AccountPrinter) Columns() [][]string {
	return [][]string{0: {
		"BALANCE",
		"PENDING CHARGES",
		"LAST PAYMENT DATE",
		"LAST PAYMENT AMOUNT",
		"NAME",
		"EMAIL",
		"ACLS",
	}}
}

// Data ...
func (a *AccountPrinter) Data() [][]string {
	return [][]string{0: {
		strconv.FormatFloat(float64(a.Account.Balance), 'f', 2, 32),
		strconv.FormatFloat(float64(a.Account.PendingCharges), 'f', 2, 32),
		a.Account.LastPaymentDate,
		strconv.FormatFloat(float64(a.Account.LastPaymentAmount), 'f', 2, 32),
		a.Account.Name,
		a.Account.Email,
		printer.ArrayOfStringsToString(a.Account.ACL),
	}}
}

// Paging ...
func (a *AccountPrinter) Paging() [][]string {
	return nil
}
