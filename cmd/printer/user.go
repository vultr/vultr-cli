package printer

import (
	"github.com/vultr/govultr"
)

func User(user []govultr.User) {
	col := columns{"USERID", "NAME", "EMAIL", "API", "ACL"}
	display(col)
	for _, u := range user {
		display(columns{u.UserID, u.Name, u.Email, u.APIEnabled, u.ACL})
	}
	flush()
}
