// Package applications provides the application functionality for the CLI
package applications

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	appLong    = `Display all commands for applications`
	appExample = `
	# Full example
	vultr-cli applications
	`

	listLong    = `Display all available applications.`
	listExample = `
	# Full example
	vultr-cli applications list
		
	# Full example with paging
	vultr-cli applications list --per-page=1 --cursor="bmV4dF9fMg=="

	# Shortened with alias commands
	vultr-cli a l
	`
)

// NewCmdApplications creates cobra command for applications
func NewCmdApplications(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "apps",
		Aliases: []string{"a", "application", "applications", "app"},
		Short:   "display applications",
		Long:    appLong,
		Example: appExample,
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list applications",
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			apps, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving application list : %v", err)
			}

			data := &ApplicationsPrinter{Applications: apps, Meta: meta}
			o.Printer.Display(data, err)

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
	Base    *cli.Base
	Printer *printer.Output
}

func (o *options) list() ([]govultr.Application, *govultr.Meta, error) {
	list, meta, _, err := o.Base.Client.Application.List(context.Background(), o.Base.Options)
	return list, meta, err
}
