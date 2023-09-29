package printer

import (
	"github.com/vultr/govultr/v3"
)

func Account(account *govultr.Account) {
	defer flush()
	display(columns{"BALANCE", "PENDING CHARGES", "LAST PAYMENT DATE", "LAST PAYMENT AMOUNT", "NAME", "EMAIL", "ACLS"})
	display(columns{
		account.Balance,
		account.PendingCharges,
		account.LastPaymentDate,
		account.LastPaymentAmount,
		account.Name,
		account.Email,
		account.ACL,
	})
}
