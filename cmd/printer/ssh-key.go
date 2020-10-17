package printer

import (
	"github.com/vultr/govultr/v2"
)

func SSHKeys(ssh []govultr.SSHKey, meta *govultr.Meta) {
	col := columns{"ID", "DATE CREATED", "NAME", "KEY"}
	display(col)
	for _, s := range ssh {
		display(columns{s.ID, s.DateCreated, s.Name, s.SSHKey})
	}

	Meta(meta)
	flush()
}

func SSHKey(ssh *govultr.SSHKey) {
	display(columns{"ID", "DATE CREATED", "NAME", "KEY"})
	display(columns{ssh.ID, ssh.DateCreated, ssh.Name, ssh.SSHKey})

	flush()
}
