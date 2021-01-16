package printer

import (
	"github.com/vultr/govultr/v2"
)

func Account(account *govultr.Account) {
	col := columns{"BALANCE", "PENDING CHARGES", "LAST PAYMENT DATE", "LAST PAYMENT AMOUNT", "NAME", "EMAIL", "ACLS"}
	display(col)
	display(columns{account.Balance, account.PendingCharges, account.LastPaymentDate, account.LastPaymentAmount, account.Name, account.Email, account.ACL})
	flush()
}
