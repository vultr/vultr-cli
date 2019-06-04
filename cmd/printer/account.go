package printer

import (
	"github.com/vultr/govultr"
)

func Account(account *govultr.Account) {
	col := columns{"BALANCE", "PENDING CHARGES", "LAST PAYMENT DATE", "LAST PAYMENT AMOUNT"}
	display(col)
	display(columns{account.Balance, account.PendingCharges, account.LastPaymentDate, account.LastPaymentAmount})
	flush()
}
