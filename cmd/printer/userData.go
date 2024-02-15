package printer

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/vultr/govultr/v2"
)

func UserData(u *govultr.UserData) {
	display(columns{"USERDATA"})
	data, err := base64.StdEncoding.DecodeString(u.Data)
	if err != nil {
		fmt.Printf("Error decoding user-data: %v\n", err)
		os.Exit(1)
	}
	display(columns{string(data)})
	flush()
}
