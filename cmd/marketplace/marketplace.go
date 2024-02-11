// Package marketplace provides the command for the CLI to access marketplace
// functionality
package marketplace

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `Get commands available to marketplace`
	example = `
	# Full example
	vultr-cli marketplace
	`
	appLong    = `Get commands available to marketplace apps`
	appExample = `
	# Full example
	vultr-cli marketplace app
	`
	listAppVariableLong    = `List all user-supplied variables for a given Vultr Marketplace app`
	listAppVariableExample = `
	# Full example
	vultr-cli marketplace app list-variables exampleapp
	`
)

// NewCmdMarketplace provides the CLI command for marketplace functions
func NewCmdMarketplace(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "marketplace",
		Short:   "Commands to interact with Marketplace functions",
		Long:    long,
		Example: example,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			return nil
		},
	}

	// App
	app := &cobra.Command{
		Use:     "app",
		Short:   "commands to interact with vultr marketplace apps",
		Long:    appLong,
		Example: appExample,
	}

	// List Variables
	listVariables := &cobra.Command{
		Use:     "list-variables",
		Short:   "List all user-supplied variables for a marketplace app",
		Aliases: []string{"l"},
		Long:    listAppVariableLong,
		Example: listAppVariableExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an image ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vars, err := o.listVariables()
			if err != nil {
				return fmt.Errorf("error getting list of marketplace app variables : %v", err)
			}

			o.Base.Printer.Display(&VariablesPrinter{Variables: vars}, nil)

			return nil
		},
	}

	app.AddCommand(
		listVariables,
	)

	cmd.AddCommand(
		app,
	)

	return cmd
}

type options struct {
	Base *cli.Base
}

func (o *options) listVariables() ([]govultr.MarketplaceAppVariable, error) {
	vars, _, err := o.Base.Client.Marketplace.ListAppVariables(o.Base.Context, o.Base.Args[0])
	return vars, err
}
