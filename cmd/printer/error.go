package printer

import (
	"encoding/json"
	"fmt"
	"os"
)

type T struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func Error(err error) {
	// TODO make errors uniform
	// t := errorToStruct(err)

	col := columns{"ERROR MESSAGE", "STATUS CODE"}
	display(col)
	// display(columns{t.Error, t.Status})
	fmt.Printf("%v", err)
	flush()

	os.Exit(1)
}

func errorToStruct(err error) *T { //nolint:unused
	t := &T{}
	if errMar := json.Unmarshal([]byte(err.Error()), t); errMar != nil {
		panic(errMar)
	}
	return t
}
