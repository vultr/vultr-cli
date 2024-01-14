// Package printer provides the console printing functionality for the CLI
package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/vultr/govultr/v3"
	"gopkg.in/yaml.v3"
)

const (
	twMinWidth       int    = 0
	twTabWidth       int    = 8
	twPadding        int    = 2
	twPadChar        byte   = '\t'
	twFlags          uint   = 0
	emptyPlaceholder string = "---"
	JSONIndent       string = "    "
)

type ResourceOutput interface {
	JSON() []byte
	YAML() []byte
	Columns() [][]string
	Data() [][]string
	Paging() [][]string
}

type Printer interface {
	display(values columns, lengths []int)
	flush()
}

type Output struct {
	Type     string
	Resource ResourceOutput
	Output   string
}

type Paging struct {
	Total      int
	CursorNext string
	CursorPrev string
}

type columns []interface{}

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

// Display confirms the output format then displays the ResourceOutput data to
// the CLI.  If there is an error, that is displayed instead via Error
func (o *Output) Display(r ResourceOutput, err error) {
	defer o.flush()

	if err != nil {
		//todo move this so it can follow the flow of the other printers and support json/yaml
		Error(err)
	}

	if strings.ToLower(o.Output) == "json" {
		o.displayNonText(r.JSON())
		os.Exit(0)
	} else if strings.ToLower(o.Output) == "yaml" {
		o.displayNonText(r.YAML())
		os.Exit(0)
	}

	o.display(r.Columns())
	o.display(r.Data())
	if r.Paging() != nil {
		o.display(r.Paging())
	}
}

func (o *Output) display(d [][]string) {
	for n := range d {
		for i := range d[n] {
			format := "\t%s"
			if i == 0 {
				format = "%s"
			}
			fmt.Fprintf(tw, format, fmt.Sprintf("%v", d[n][i]))
		}
		fmt.Fprintf(tw, "\n")
	}
}

func (o *Output) flush() {
	tw.Flush()
}

func (o *Output) displayNonText(data []byte) {
	fmt.Printf("%s\n", string(data))
}

// NewPaging validates and intializes the paging data
func NewPaging(total int, next *string, prev *string) *Paging {
	p := new(Paging)
	p.Total = total

	if next != nil {
		p.CursorNext = *next
	} else {
		p.CursorNext = emptyPlaceholder
	}

	if prev != nil {
		p.CursorPrev = *prev
	} else {
		p.CursorPrev = emptyPlaceholder
	}

	return p
}

// Compose returns the paging data for output
func (p *Paging) Compose() [][]string {
	cNext := p.CursorNext
	if cNext == "" {
		cNext = "---"
	}

	cPrev := p.CursorPrev
	if cPrev == "" {
		cPrev = "---"
	}

	return [][]string{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {strconv.Itoa(p.Total), cNext, cPrev},
	}
}

func MarshalObject(input interface{}, format string) []byte {
	var output []byte
	if format == "json" {
		j, errJ := json.MarshalIndent(input, "", JSONIndent)
		if errJ != nil {
			panic(fmt.Errorf("error marshalling JSON : %v", errJ))
		}
		output = j

	} else if format == "yaml" {
		y, errY := yaml.Marshal(input)
		if errY != nil {
			panic(fmt.Errorf("error marshalling YAML : %v", errY))
		}
		output = y
	}

	return output
}

// ArrayOfStringsToString will build a delimited string from an array for
// display in the printer functions.  It defaults to comma-delimited and
// enclosed in square brackets to maintain consistency with array Fprintf
func ArrayOfStringsToString(a []string) string {
	delimiter := ", "
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(strings.Join(a, delimiter))
	sb.WriteString("]")

	return sb.String()
}

// OLD funcs to be re-written //////////////////////////////////////////////////////////////
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

func flush() {
	if err := tw.Flush(); err != nil {
		panic("could not flush buffer")
	}
	tw.Flush()
}

// Meta prints out the pagination details TODO: old
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
