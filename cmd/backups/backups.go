// Package backups provides access to the backups for the CLI
package backups

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	backupsLong    = ``
	backupsExample = ``
	listLong       = ``
	listExample    = ``
	getLong        = ``
	getExample     = ``
)

// NewCmdBackups provides the backup command for the CLI
func NewCmdBackups(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "backups",
		Aliases: []string{"backup", "b"},
		Short:   "Display backups",
		Long:    backupsLong,
		Example: backupsExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Short:   "List all backups",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			backups, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving backups list : %v", err)
			}
			data := &BackupsPrinter{Backups: backups, Meta: meta}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Get
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a backup",
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a backup ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			backup, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving backup : %v", err)
			}

			data := &BackupPrinter{Backup: backup}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	cmd.AddCommand(get, list)
	return cmd
}

type options struct {
	Base *cli.Base
}

func (b *options) list() ([]govultr.Backup, *govultr.Meta, error) {
	backups, meta, _, err := b.Base.Client.Backup.List(b.Base.Context, b.Base.Options)
	return backups, meta, err
}

func (b *options) get() (*govultr.Backup, error) {
	backup, _, err := b.Base.Client.Backup.Get(b.Base.Context, b.Base.Args[0])
	return backup, err
}
