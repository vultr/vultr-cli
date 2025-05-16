package account

import (
	"fmt"
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

// ======================================

// AccountPrinter ...
type AccountBandwidthPrinter struct {
	Bandwidth *govultr.AccountBandwidth `json:"account_bandwidth"`
}

// JSON ...
func (a *AccountBandwidthPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AccountBandwidthPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AccountBandwidthPrinter) Columns() [][]string {
	return [][]string{0: {
		"PERIOD",
		"GB IN",
		"GB OUT",
		"INSTANCE HOURS",
		"INSTANCE COUNT",
		"CREDITS INSTANCE",
		"CREDITS FREE",
		"CREDITS PURCHASED",
		"OVERAGE",
		"OVERAGE UNIT COST",
		"OVERAGE COST",
	}}
}

// Data ...
func (a *AccountBandwidthPrinter) Data() [][]string {
	previousPeriod := fmt.Sprintf(
		"%s - %s",
		utils.ParseAPITimestamp(a.Bandwidth.PreviousMonth.TimestampStart),
		utils.ParseAPITimestamp(a.Bandwidth.PreviousMonth.TimestampEnd),
	)

	currentPeriod := fmt.Sprintf(
		"%s - %s",
		utils.ParseAPITimestamp(a.Bandwidth.CurrentMonthToDate.TimestampStart),
		utils.ParseAPITimestamp(a.Bandwidth.CurrentMonthToDate.TimestampEnd),
	)

	projectedPeriod := fmt.Sprintf(
		"%s - %s",
		utils.ParseAPITimestamp(a.Bandwidth.CurrentMonthProjected.TimestampStart),
		utils.ParseAPITimestamp(a.Bandwidth.CurrentMonthProjected.TimestampEnd),
	)

	return [][]string{
		0: {
			previousPeriod,
			strconv.Itoa(a.Bandwidth.PreviousMonth.GBIn),
			strconv.Itoa(a.Bandwidth.PreviousMonth.GBOut),
			strconv.Itoa(a.Bandwidth.PreviousMonth.TotalInstanceHours),
			strconv.Itoa(a.Bandwidth.PreviousMonth.TotalInstanceCount),
			strconv.Itoa(a.Bandwidth.PreviousMonth.InstanceBandwidthCredits),
			strconv.Itoa(a.Bandwidth.PreviousMonth.FreeBandwidthCredits),
			strconv.Itoa(a.Bandwidth.PreviousMonth.PurchasedBandwidthCredits),
			strconv.FormatFloat(float64(a.Bandwidth.PreviousMonth.Overage), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(a.Bandwidth.PreviousMonth.OverageUnitCost), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(a.Bandwidth.PreviousMonth.OverageCost), 'f', utils.FloatPrecision, 32),
		},
		1: {
			currentPeriod,
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.GBIn),
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.GBOut),
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.TotalInstanceHours),
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.TotalInstanceCount),
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.InstanceBandwidthCredits),
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.FreeBandwidthCredits),
			strconv.Itoa(a.Bandwidth.CurrentMonthToDate.PurchasedBandwidthCredits),
			strconv.FormatFloat(float64(a.Bandwidth.CurrentMonthToDate.Overage), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(a.Bandwidth.CurrentMonthToDate.OverageUnitCost), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(a.Bandwidth.CurrentMonthToDate.OverageCost), 'f', utils.FloatPrecision, 32),
		},
		2: {
			projectedPeriod,
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.GBIn),
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.GBOut),
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.TotalInstanceHours),
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.TotalInstanceCount),
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.InstanceBandwidthCredits),
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.FreeBandwidthCredits),
			strconv.Itoa(a.Bandwidth.CurrentMonthProjected.PurchasedBandwidthCredits),
			strconv.FormatFloat(float64(a.Bandwidth.CurrentMonthProjected.Overage), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(a.Bandwidth.CurrentMonthProjected.OverageUnitCost), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(a.Bandwidth.CurrentMonthProjected.OverageCost), 'f', utils.FloatPrecision, 32),
		},
	}
}

// Paging ...
func (a *AccountBandwidthPrinter) Paging() [][]string {
	return nil
}
