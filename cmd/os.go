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
	"log"

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// Os represents the iso command
func Os() *cobra.Command {
	osCmd := &cobra.Command{
		Use:   "os",
		Short: "os is used to access os commands",
		Long:  ``,
	}

	osCmd.AddCommand(osList)

	osList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	osList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return osCmd
}

// osList represents the list of available operating systems
var osList = &cobra.Command{
	Use:     "list",
	Short:   "list all available operating systems",
	Aliases: []string{"o"},
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		os, meta, _, err := client.OS.List(context.TODO(), options)
		if err != nil {
			log.Fatal(err)
		}

		printer.Os(os, meta)
	},
}
