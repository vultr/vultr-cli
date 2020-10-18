package printer

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vultr/govultr/v2"
)

type Printer interface {
	display(values columns, lengths []int)
	flush()
}

var tw = new(tabwriter.Writer)

func init() {
	tw.Init(os.Stdout, 0, 8, 2, '\t', 0)
}

type columns []interface{}

func display(values columns) {

	for i, value := range values {
		format := "\t%s"
		if i == 0 {
			format = "%s"
		}
		fmt.Fprintf(tw, format, fmt.Sprintf("%v", value))
	}
	fmt.Fprintf(tw, "\n")
}

func flush() {
	tw.Flush()
}

func Meta(meta *govultr.Meta) {
	display(columns{"======================================"})
	col := columns{"TOTAL", "NEXT PAGE", "PREV PAGE"}
	display(col)

	display(columns{meta.Total, meta.Links.Next, meta.Links.Prev})
}
