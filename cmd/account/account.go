// Package account provides the account functionality for the CLI
package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	accountLong    = `Retrieve information about your account.`
	accountExample = `
	# Full example
	vultr-cli account
	`
)

// NewCmdAccount creates a cobra command for Account
func NewCmdAccount(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "account",
		Short:   "Display account information",
		Long:    accountLong,
		Example: accountExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			account, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving account information : %v", err)
			}

			o.Base.Printer.Display(&AccountPrinter{Account: account}, nil)

			return nil
		},
	}

	return cmd
}

type options struct {
	Base *cli.Base
}

func (o *options) get() (*govultr.Account, error) {
	account, _, err := o.Base.Client.Account.Get(context.Background())
	return account, err
}
