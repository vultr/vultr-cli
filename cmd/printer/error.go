package printer

import (
	"encoding/json"
	"os"
)

type T struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func Error(err error) {
	t := errorToStruct(err)

	col := columns{"ERROR MESSAGE", "STATUS CODE"}
	display(col)
	display(columns{t.Error, t.Status})
	flush()

	os.Exit(1)
}

func errorToStruct(err error) *T {
	t := &T{}
	if err := json.Unmarshal([]byte(err.Error()), t); err != nil {
		panic(err)
	}
	return t
}
