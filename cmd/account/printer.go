package account

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// AccountPrinter ...
type AccountPrinter struct {
	Account *govultr.Account `json:"account"`
}

// JSON ...
func (a *AccountPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AccountPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
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
		strconv.FormatFloat(float64(a.Account.Balance), 'f', utils.FloatPrecision, 32),
		strconv.FormatFloat(float64(a.Account.PendingCharges), 'f', utils.FloatPrecision, 32),
		a.Account.LastPaymentDate,
		strconv.FormatFloat(float64(a.Account.LastPaymentAmount), 'f', utils.FloatPrecision, 32),
		a.Account.Name,
		a.Account.Email,
		printer.ArrayOfStringsToString(a.Account.ACL),
	}}
}

// Paging ...
func (a *AccountPrinter) Paging() [][]string {
	return nil
}
