package printer

import "github.com/vultr/govultr"

func ApiKey(api *govultr.API) {
	col := columns{"NAME", "EMAIL", "ACLS"}
	display(col)
	display(columns{api.Name, api.Email, api.ACL})
	flush()
}
