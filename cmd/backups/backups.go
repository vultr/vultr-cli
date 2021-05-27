// Copyright © 2021 The Vultr-cli Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backups

import (
	"context"

	"github.com/vultr/vultr-cli/pkg/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
	"github.com/vultr/vultr-cli/cmd/utils"
)

var (
	backupLong    = ``
	backupExample = ``

	getLong    = `Get a backup from your Vultr account based on it's ID.`
	getExample = `
		# Full example
		vultr-cli backup get 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
		
		# Shortened with alias commands
		vultr-cli b g 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
	`
	listLong    = `List backups from your Vultr account.`
	listExample = `
		# Full example
		vultr-cli backup list
		
		# Shortened with alias commands
		vultr-cli b l
	`
)

// BackupOptionsInterface ...
type BackupOptionsInterface interface {
	validate(cmd *cobra.Command, args []string)
	List() ([]govultr.Backup, *govultr.Meta, error)
	Get() (*govultr.Backup, error)
}

// Options for backups
type BackupOptions struct {
	Base *cli.Base
}

// NewBackupOptions ...
func NewBackupOptions(base *cli.Base) *BackupOptions {
	return &BackupOptions{Base: base}
}

// NewCmdBackup ...
func NewCmdBackup(base *cli.Base) *cobra.Command {
	b := NewBackupOptions(base)

	cmd := &cobra.Command{
		Use:     "backup",
		Aliases: []string{"backups", "backup", "b"},
		Short:   "backup commands",
		Long:    backupLong,
		Example: backupExample,
	}

	// Get Command
	get := &cobra.Command{
		Use:     "get {backupID}",
		Short:   "get a backup",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Run: func(cmd *cobra.Command, args []string) {
			b.validate(cmd, args)
			backup, err := b.Get()
			if err != nil {
				printer.Error(err)
			}
			b.Base.Printer.Display(&BackupPrinter{Backup: backup}, err)
		},
	}

	//list
	list := &cobra.Command{
		Use:     "list",
		Short:   "list backups",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			b.validate(cmd, args)
			b.Base.Options = utils.GetPaging(cmd)
			backup, meta, err := b.List()
			if err != nil {
				printer.Error(err)
			}
			b.Base.Printer.Display(&BackupsPrinter{Backups: backup, Meta: meta}, err)
		},
	}
	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	cmd.AddCommand(get, list)
	return cmd
}

func (b *BackupOptions) validate(cmd *cobra.Command, args []string) {
	b.Base.Args = args
	b.Base.Printer.Output = viper.GetString("output")

	if cmd.Use == "list" {
		b.Base.Options = utils.GetPaging(cmd)
	}
}

// Get a single backup based on ID
func (b *BackupOptions) Get() (*govultr.Backup, error) {
	backup, err := b.Base.Client.Backup.Get(context.Background(), b.Base.Args[0])
	if err != nil {
		return nil, err
	}
	return backup, nil
}

// Get list of backups
func (b *BackupOptions) List() ([]govultr.Backup, *govultr.Meta, error) {
	backup, meta, err := b.Base.Client.Backup.List(context.Background(), b.Base.Options)
	if err != nil {
		return nil, nil, err
	}
	return backup, meta, nil
}
