package printer

import (
	"github.com/vultr/govultr/v3"
)

func BillingHistory(bh []govultr.History, meta *govultr.Meta) {
	col := columns{"ID", "DATE", "TYPE", "DESCRIPTION", "AMOUNT", "BALANCE"}
	display(col)
	for _, b := range bh {
		display(columns{b.ID, b.Date, b.Type, b.Description, b.Amount, b.Balance})
	}

	Meta(meta)
	flush()
}

func Invoices(inv []govultr.Invoice, meta *govultr.Meta) {
	col := columns{"ID", "DATE", "DESCRIPTION", "AMOUNT", "BALANCE"}
	display(col)
	for _, i := range inv {
		display(columns{i.ID, i.Date, i.Description, i.Amount, i.Balance})
	}

	Meta(meta)
	flush()
}

func Invoice(i *govultr.Invoice) {
	display(columns{"ID", "DATE", "DESCRIPTION", "AMOUNT", "BALANCE"})
	display(columns{i.ID, i.Date, i.Description, i.Amount, i.Balance})

	flush()
}

func InvoiceItems(inv []govultr.InvoiceItem, meta *govultr.Meta) {
	col := columns{"DESCRIPTION", "PRODUCT", "START DATE", "END DATE", "UNITS", "UNIT TYPE", "UNIT PRICE", "TOTAL"}
	display(col)
	for _, i := range inv {
		display(columns{i.Description, i.Product, i.StartDate, i.EndDate, i.Units, i.UnitType, i.UnitPrice, i.Total})
	}

	Meta(meta)
	flush()
}
