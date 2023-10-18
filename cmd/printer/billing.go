package printer

import (
	"github.com/vultr/govultr/v3"
)

func BillingHistory(bh []govultr.History, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE", "TYPE", "DESCRIPTION", "AMOUNT", "BALANCE"})

	if len(bh) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range bh {
		display(columns{
			bh[i].ID,
			bh[i].Date,
			bh[i].Type,
			bh[i].Description,
			bh[i].Amount,
			bh[i].Balance,
		})
	}

	Meta(meta)
}

func Invoices(inv []govultr.Invoice, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE", "DESCRIPTION", "AMOUNT", "BALANCE"})

	if len(inv) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range inv {
		display(columns{
			inv[i].ID,
			inv[i].Date,
			inv[i].Description,
			inv[i].Amount,
			inv[i].Balance,
		})
	}

	Meta(meta)
}

func Invoice(i *govultr.Invoice) {
	defer flush()

	display(columns{"ID", "DATE", "DESCRIPTION", "AMOUNT", "BALANCE"})
	display(columns{i.ID, i.Date, i.Description, i.Amount, i.Balance})
}

func InvoiceItems(inv []govultr.InvoiceItem, meta *govultr.Meta) {
	defer flush()

	display(columns{"DESCRIPTION", "PRODUCT", "START DATE", "END DATE", "UNITS", "UNIT TYPE", "UNIT PRICE", "TOTAL"})

	if len(inv) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range inv {
		display(columns{
			inv[i].Description,
			inv[i].Product,
			inv[i].StartDate,
			inv[i].EndDate,
			inv[i].Units,
			inv[i].UnitType,
			inv[i].UnitPrice,
			inv[i].Total,
		})
	}

	Meta(meta)
}
