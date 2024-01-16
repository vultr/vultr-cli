// Package backups provides access to the backups for the CLI
package backups

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
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
	updateLong     = ``
	updateExample  = ``
)

type BackupsOptionsInterface interface {
	setOptions(cmd *cobra.Command, args []string)
	List() []govultr.Backup
	Get() *govultr.Backup
}

// BackupOptions ...
type BackupsOptions struct {
	Base *cli.Base
}

// NewBackupOptions ...
func NewBackupsOptions(base *cli.Base) *BackupsOptions {
	return &BackupsOptions{Base: base}
}

// NewCmdBackup ...
func NewCmdBackups(base *cli.Base) *cobra.Command {
	o := NewBackupsOptions(base)

	cmd := &cobra.Command{
		Use:     "backups",
		Aliases: []string{"backup", "b"},
		Short:   "user commands",
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
		Short:   "list all backups",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)
			backups, meta, err := o.List()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving backups list : %v", err))
				os.Exit(1)
			}
			data := &BackupsPrinter{Backups: backups, Meta: meta}
			o.Base.Printer.Display(data, err)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Get
	get := &cobra.Command{
		Use:   "get",
		Short: "get a backup",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a backup ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			backup, err := o.Get()
			if err != nil {
				panic(fmt.Errorf("error retrieving backup : %v", err))
			}

			data := &BackupPrinter{Backup: backup}
			o.Base.Printer.Display(data, err)
		},
	}

	cmd.AddCommand(get, list)
	return cmd
}

func (b *BackupsOptions) List() ([]govultr.Backup, *govultr.Meta, error) {
	backups, meta, _, err := b.Base.Client.Backup.List(b.Base.Context, b.Base.Options)
	if err != nil {
		fmt.Printf("Error with backups list request : %v\n", err)
		os.Exit(1)
	}

	return backups, meta, nil
}

func (b *BackupsOptions) Get() (*govultr.Backup, error) {
	backup, _, err := b.Base.Client.Backup.Get(b.Base.Context, b.Base.Args[0])
	if err != nil {
		fmt.Printf("Error with backups get request : %v\n", err)
		os.Exit(1)
	}

	return backup, nil
}
