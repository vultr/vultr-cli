// Package billing provides the account billing operations and
// functionality for the CLI
package billing

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// BillingHistoryPrinter ...
type BillingHistoryPrinter struct {
	Billing []govultr.History `json:"billing_history"`
	Meta    *govultr.Meta     `json:"meta"`
}

// JSON ...
func (b *BillingHistoryPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BillingHistoryPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BillingHistoryPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE",
		"TYPE",
		"DESCRIPTION",
		"AMOUNT",
		"BALANCE",
	}}
}

// Data ...
func (b *BillingHistoryPrinter) Data() [][]string {
	if len(b.Billing) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range b.Billing {
		data = append(data, []string{
			strconv.Itoa(b.Billing[i].ID),
			b.Billing[i].Date,
			b.Billing[i].Type,
			b.Billing[i].Description,
			strconv.FormatFloat(float64(b.Billing[i].Amount), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(b.Billing[i].Balance), 'f', utils.FloatPrecision, 32),
		})
	}
	return data
}

// Paging ...
func (b *BillingHistoryPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(b.Meta).Compose()
}

// ======================================

// BillingInvoicesPrinter ...
type BillingInvoicesPrinter struct {
	Invoices []govultr.Invoice `json:"billing_invoices"`
	Meta     *govultr.Meta
}

// JSON ...
func (b *BillingInvoicesPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BillingInvoicesPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BillingInvoicesPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE",
		"DESCRIPTION",
		"AMOUNT",
		"BALANCE",
	}}
}

// Data ...
func (b *BillingInvoicesPrinter) Data() [][]string {
	if len(b.Invoices) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range b.Invoices {
		data = append(data, []string{
			strconv.Itoa(b.Invoices[i].ID),
			b.Invoices[i].Date,
			b.Invoices[i].Description,
			strconv.FormatFloat(float64(b.Invoices[i].Amount), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(b.Invoices[i].Balance), 'f', utils.FloatPrecision, 32),
		})
	}
	return data
}

// Paging ...
func (b *BillingInvoicesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(b.Meta).Compose()
}

// ======================================

// BillingInvoicePrinter ...
type BillingInvoicePrinter struct {
	Invoice govultr.Invoice `json:"billing_invoice"`
}

// JSON ...
func (b *BillingInvoicePrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BillingInvoicePrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BillingInvoicePrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE",
		"DESCRIPTION",
		"AMOUNT",
		"BALANCE",
	}}
}

// Data ...
func (b *BillingInvoicePrinter) Data() [][]string {
	return [][]string{0: {
		strconv.Itoa(b.Invoice.ID),
		b.Invoice.Date,
		b.Invoice.Description,
		strconv.FormatFloat(float64(b.Invoice.Amount), 'f', utils.FloatPrecision, 32),
		strconv.FormatFloat(float64(b.Invoice.Balance), 'f', utils.FloatPrecision, 32),
	}}
}

// Paging ...
func (b *BillingInvoicePrinter) Paging() [][]string {
	return nil
}

// ======================================

// BillingInvoiceItemsPrinter ...
type BillingInvoiceItemsPrinter struct {
	InvoiceItems []govultr.InvoiceItem `json:"invoice_items"`
	Meta         *govultr.Meta         `json:"meta"`
}

// JSON ...
func (b *BillingInvoiceItemsPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BillingInvoiceItemsPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BillingInvoiceItemsPrinter) Columns() [][]string {
	return [][]string{0: {
		"DESCRIPTION",
		"PRODUCT",
		"START DATE",
		"END DATE",
		"UNITS",
		"UNIT TYPE",
		"UNIT PRICE",
		"TOTAL",
	}}
}

// Data ...
func (b *BillingInvoiceItemsPrinter) Data() [][]string {
	if len(b.InvoiceItems) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range b.InvoiceItems {
		data = append(data, []string{
			b.InvoiceItems[i].Description,
			b.InvoiceItems[i].Product,
			b.InvoiceItems[i].StartDate,
			b.InvoiceItems[i].EndDate,
			strconv.Itoa(b.InvoiceItems[i].Units),
			b.InvoiceItems[i].UnitType,
			strconv.FormatFloat(float64(b.InvoiceItems[i].UnitPrice), 'f', utils.FloatPrecision, 32),
			strconv.FormatFloat(float64(b.InvoiceItems[i].Total), 'f', utils.FloatPrecision, 32),
		})
	}

	return data
}

// Paging ...
func (b *BillingInvoiceItemsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(b.Meta).Compose()
}
