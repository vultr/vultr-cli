package printer

import (
	"github.com/vultr/govultr/v3"
)

func SSHKeys(ssh []govultr.SSHKey, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "NAME", "KEY"})

	if len(ssh) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range ssh {
		display(columns{
			ssh[i].ID,
			ssh[i].DateCreated,
			ssh[i].Name,
			ssh[i].SSHKey,
		})
	}

	Meta(meta)
}

func SSHKey(ssh *govultr.SSHKey) {
	defer flush()

	display(columns{"ID", "DATE CREATED", "NAME", "KEY"})
	display(columns{ssh.ID, ssh.DateCreated, ssh.Name, ssh.SSHKey})
}
