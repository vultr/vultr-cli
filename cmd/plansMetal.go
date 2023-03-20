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

var (
	plansMetalLong    = `Get commands available to metal`
	plansMetalExample = `
	#Full example
	vultr-cli plans metal

	#Shortened with aliased commands
	vultr-cli p m
	`

	plansMetalListLong    = `Get all Vultr bare metal plans`
	plansMetalListExample = `
	#Full example
	vultr-cli plans metal list

	#Full example with paging
	vultr-cli plans metal list --per-page=10 --cursor="bmV4dF9fQU1T"

	#Shortened with aliased commands
	vultr-cli p m l
	`
)

// PlansMetal represents the metal sub command
func PlansMetal() *cobra.Command {
	planMetalCmd := &cobra.Command{
		Use:     "metal",
		Aliases: []string{"m", "metals"},
		Short:   "metal is used to access bare metal commands",
		Long:    plansMetalLong,
		Example: plansMetalExample,
	}

	planMetalCmd.AddCommand(metalList)

	metalList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	metalList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return planMetalCmd
}

var metalList = &cobra.Command{
	Use:     "list",
	Short:   "list bare-metal plans",
	Long:    plansMetalListLong,
	Example: plansMetalListExample,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, _, err := client.Plan.ListBareMetal(context.TODO(), options)

		if err != nil {
			fmt.Printf("error getting bare metal plan list : %v\n", err)
			os.Exit(1)
		}

		printer.PlanBareMetal(list, meta)
	},
}
