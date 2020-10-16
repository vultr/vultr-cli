package printer

import (
	"github.com/vultr/govultr"
)

func Users(user []govultr.User, meta *govultr.Meta) {
	col := columns{"ID", "NAME", "EMAIL", "API", "ACL"}
	display(col)
	for _, u := range user {
		display(columns{u.ID, u.Name, u.Email, u.APIEnabled, u.ACL})
	}

	Meta(meta)
	flush()
}

func User(user *govultr.User) {
	display(columns{"ID", "NAME", "EMAIL", "API", "ACL"})
	display(columns{user.ID, user.Name, user.Email, user.APIEnabled, user.ACL})

	flush()
}
