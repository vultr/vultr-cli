// Copyright © 2019 The Vultr-cli Authors
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

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// Backups  represents the application command
func Backups() *cobra.Command {
	backupsCmd := &cobra.Command{
		Use:     "backups",
		Aliases: []string{"b"},
		Short:   "Display backups",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Context().Value(ctxAuthKey{}).(bool) {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	backupsCmd.AddCommand(backupsList, backupsGet)

	backupsList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	backupsList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return backupsCmd
}

var backupsList = &cobra.Command{
	Use:     "list",
	Short:   "list backups",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		backups, meta, _, err := client.Backup.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting backups : %v\n", err)
			os.Exit(1)
		}

		printer.Backups(backups, meta)
	},
}

var backupsGet = &cobra.Command{
	Use:   "get",
	Short: "get backup",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a backupID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		backup, _, err := client.Backup.Get(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting backup : %v\n", err)
			os.Exit(1)
		}

		printer.Backup(backup)
	},
}
