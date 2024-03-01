// Package operatingsystems provides the operating systems functionality for
// the CLI
package operatingsystems

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `OS will retrieve available operating systems that can be deployed`
	example = `
	# Example
	vultr-cli os
	`

	listLong    = `List all operating systems available to be deployed on Vultr`
	listExample = `
	# Full example
	vultr-cli os list
		
	# Full example with paging
	vultr-cli os list --per-page=1 --cursor="bmV4dF9fMTI0" 

	# Shortened with alias commands
	vultr-cli o l
	`
)

// NewCmdOS provides the command for operating systems to the CLI
func NewCmdOS(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "os",
		Short:   "Display available operating systems",
		Aliases: []string{"o"},
		Long:    long,
		Example: example,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Short:   "List all operating systems",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			os, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error getting operating systems : %v", err)
			}

			data := &OSPrinter{OperatingSystems: os, Meta: meta}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)

	cmd.AddCommand(list)
	return cmd
}

type options struct {
	Base *cli.Base
}

func (o *options) list() ([]govultr.OS, *govultr.Meta, error) {
	list, meta, _, err := o.Base.Client.OS.List(context.Background(), o.Base.Options)
	return list, meta, err
}
