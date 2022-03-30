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
	plansLong = `Get commands available to plans`
	plansExample = `
	#Full example
	vultr-cli plans
	`

	plansListLong = `Get all Vultr plans`
	plansListExample = `
	#Full example
	vultr-cli plans list

	#Full example with paging
	vultr-cli plans list --type=vc2 --per-page=5 --cursor="bmV4dF9fdmMyLTJjLTRnYg=="
	
	#Shortened with aliased commands
	vultr-cli p l
	`
)

// Plans represents the plans command
func Plans() *cobra.Command {
	planCmd := &cobra.Command{
		Use:     "plans",
		Short:   "get information about Vultr plans",
		Aliases: []string{"p"},
		Long:	 plansLong,
		Example: plansExample,
	}

	planCmd.AddCommand(planList)
	planCmd.AddCommand(PlansMetal())

	planList.Flags().StringP("type", "t", "", "(optional) The type of plans to return. Possible values: 'bare-metal', 'vc2', 'vdc2', 'ssd', 'dedicated'. Defaults to all VPS plans.")

	planList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	planList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return planCmd
}

var planList = &cobra.Command{
	Use:     "list",
	Short:   "list plans",
	Aliases: []string{"l"},
	Long:	 plansListLong,
	Example: plansListExample,
	Run: func(cmd *cobra.Command, args []string) {
		planType, _ := cmd.Flags().GetString("type")
		options := getPaging(cmd)

		if planType == "bare-metal" {
			list, meta, err := client.Plan.ListBareMetal(context.TODO(), options)
			if err != nil {
				fmt.Printf("error getting bare metal plan list : %v\n", err)
				os.Exit(1)
			}

			printer.PlanBareMetal(list, meta)
		} else {
			list, meta, err := client.Plan.List(context.TODO(), planType, options)
			if err != nil {
				fmt.Printf("error getting plan list : %v\n", err)
				os.Exit(1)
			}

			printer.Plan(list, meta)
		}
	},
}
