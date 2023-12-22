// Package printer provides the console printing functionality for the CLI
package printer

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/vultr/govultr/v3"
)

const (
	twMinWidth int  = 0
	twTabWidth int  = 8
	twPadding  int  = 2
	twPadChar  byte = '\t'
	twFlags    uint = 0
)

type Printer interface {
	display(values columns, lengths []int)
	flush()
}

var tw = new(tabwriter.Writer)

func init() {
	tw.Init(
		os.Stdout,
		twMinWidth,
		twTabWidth,
		twPadding,
		twPadChar,
		twFlags,
	)
}

// columns is a type to contain the strings
type columns []interface{}

// display loops over the value `columns` and Fprintf the output to tabwriter
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

// displayString will `Fprintln` a string to the tabwriter
func displayString(message string) {
	fmt.Fprintln(tw, message)
}

// arrayOfStringsToString will build a delimited string from an array for
// display in the printer functions.  Defaulted to comma-delimited and enclosed
// in square brackets to maintain consistency with array Fprintf
func arrayOfStringsToString(a []string) string {
	delimiter := ", "
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(strings.Join(a, delimiter))
	sb.WriteString("]")

	return sb.String()
}

// flush calls the tabwriter `Flush()` to write output
func flush() {
	if err := tw.Flush(); err != nil {
		panic("could not flush buffer")
	}
}

// Meta prints out the pagination details
func Meta(meta *govultr.Meta) {
	var pageNext string
	var pagePrev string

	if meta.Links.Next == "" {
		pageNext = "---"
	} else {
		pageNext = meta.Links.Next
	}

	if meta.Links.Prev == "" {
		pagePrev = "---"
	} else {
		pagePrev = meta.Links.Prev
	}

	displayString("======================================")
	display(columns{"TOTAL", "NEXT PAGE", "PREV PAGE"})
	display(columns{meta.Total, pageNext, pagePrev})
}

// MetaDBaaS prints out the pagination details used by database commands
func MetaDBaaS(meta *govultr.Meta) {
	displayString("======================================")
	display(columns{"TOTAL"})

	display(columns{meta.Total})
}
