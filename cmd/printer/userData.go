package printer

import (
	"encoding/base64"
	"fmt"

	"github.com/vultr/govultr/v3"
)

func UserData(u *govultr.UserData) {
	defer flush()

	display(columns{"USERDATA"})
	data, err := base64.StdEncoding.DecodeString(u.Data)
	if err != nil {
		displayString(fmt.Sprintf("Error decoding user-data: %v\n", err))
		return
	}
	display(columns{string(data)})
}
