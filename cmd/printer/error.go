package printer

import (
	"fmt"
	"os"
)

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
