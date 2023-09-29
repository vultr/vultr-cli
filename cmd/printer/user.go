package printer

import (
	"github.com/vultr/govultr/v3"
)

func Users(user []govultr.User, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "NAME", "EMAIL", "API", "ACL"})

	if len(user) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range user {
		display(columns{
			user[i].ID,
			user[i].Name,
			user[i].Email,
			*user[i].APIEnabled,
			user[i].ACL,
		})
	}

	Meta(meta)
}

func User(user *govultr.User) {
	defer flush()

	display(columns{"ID", "NAME", "EMAIL", "API", "ACL"})
	display(columns{user.ID, user.Name, user.Email, *user.APIEnabled, user.ACL})
}
