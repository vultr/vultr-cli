package printer

import (
	"github.com/vultr/govultr"
)

func SSHKey(ssh []govultr.SSHKey) {
	col := columns{"SSHKEYID", "DATE CREATED", "NAME", "KEY"}
	display(col)
	for _, s := range ssh {
		display(columns{s.SSHKeyID, s.DateCreated, s.Name, s.Key})
	}
	flush()
}
