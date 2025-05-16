// Package account provides the account functionality for the CLI
package account

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	accountLong    = `Account related commands`
	accountExample = `
	# Full example
	vultr-cli account
	`
	accountInfoLong    = `Retrieve information about your account.`
	accountInfoExample = `
	# Full example
	vultr-cli account info
	`
	accountBandwidthLong    = `Retrieve information about accout bandwidth`
	accountBandwidthExample = `
	# Full example
	vultr-cli account bandwidth
	`
)

// NewCmdAccount creates a cobra command for Account
func NewCmdAccount(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "account",
		Short:   "Commands related to account information",
		Long:    accountLong,
		Example: accountExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	info := &cobra.Command{
		Use:     "info",
		Short:   "Display account information",
		Long:    accountInfoLong,
		Example: accountInfoExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			account, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving account information : %v", err)
			}

			o.Base.Printer.Display(&AccountPrinter{Account: account}, nil)

			return nil
		},
	}

	bandwidth := &cobra.Command{
		Use:     "bandwidth",
		Short:   "Display account bandwidth",
		Long:    accountBandwidthLong,
		Example: accountBandwidthExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			bandwidth, err := o.getBandwidth()
			if err != nil {
				return fmt.Errorf("error retrieving account bandwidth : %v", err)
			}

			o.Base.Printer.Display(&AccountBandwidthPrinter{Bandwidth: bandwidth}, nil)

			return nil
		},
	}

	cmd.AddCommand(
		info,
		bandwidth,
	)

	return cmd
}

type options struct {
	Base *cli.Base
}

func (o *options) get() (*govultr.Account, error) {
	account, _, err := o.Base.Client.Account.Get(o.Base.Context)
	return account, err
}

func (o *options) getBandwidth() (*govultr.AccountBandwidth, error) {
	bw, _, err := o.Base.Client.Account.GetBandwidth(o.Base.Context)
	return bw, err
}
