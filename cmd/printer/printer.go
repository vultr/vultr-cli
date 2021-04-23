package printer

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vultr/govultr/v2"
)

type ResourceOutput interface {
	Json()
	Yaml()
	Columns() map[int][]interface{}
	Data() map[int][]interface{}
	Paging() map[int][]interface{}
}

type Printer interface {
	display(values columns, lengths []int)
	flush()
}

type Output struct {
	Type     string
	Resource ResourceOutput
}

type columns2 map[int][]interface{}
type columns []interface{}
var tw = new(tabwriter.Writer)

func init() {
	tw.Init(os.Stdout, 0, 8, 2, '\t', 0)
}

func (o *Output) Display(r ResourceOutput, err error) {
	if err != nil {
		Error(err)
	}

	o.display(r.Columns())
	o.display(r.Data())
	if r.Paging() != nil {
		o.display(r.Paging())
	}
	o.flush()
}

func (o *Output) display(d columns2) {
	for _, values := range d {
		for i, value := range values {
			format := "\t%s"
			if i == 0 {
				format = "%s"
			}
			fmt.Fprintf(tw, format, fmt.Sprintf("%v", value))
		}
		fmt.Fprintf(tw, "\n")
	}
}

func (o *Output) flush(){
	tw.Flush()
}




////////////////////////////////////////////////////////////////
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
