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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// Applications  represents the application command
func Applications() *cobra.Command {
	appsCmd := &cobra.Command{
		Use:     "apps",
		Aliases: []string{"a", "application", "applications"},
		Short:   "Display all available applications",
	}

	appsCmd.AddCommand(appsList)

	appsList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	appsList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return appsCmd
}

var appsList = &cobra.Command{
	Use:     "list",
	Short:   "list applications",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		apps, meta, _, err := client.Application.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting available applications : %v\n", err)
			os.Exit(1)
		}

		printer.Application(apps, meta)
	},
}
