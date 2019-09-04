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
	"github.com/vultr/vultr-cli/cmd/printer"
)

// Plans represents the plans command
func Plans() *cobra.Command {
	planCmd := &cobra.Command{
		Use:     "plans",
		Short:   "get information about Vultr plans",
		Aliases: []string{"p"},
	}

	planCmd.AddCommand(planList)

	planList.Flags().StringP("type", "t", "", "(optional) The type of plans to return. Possible values: 'bare-metal', 'vc2', 'vdc2', 'ssd', 'dedicated'. Defaults to all VPS plans.")

	return planCmd
}

var planList = &cobra.Command{
	Use:     "list",
	Short:   "list plans",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		planType, _ := cmd.Flags().GetString("type")

		if planType == "bare-metal" {
			list, err := client.Plan.GetBareMetalList(context.TODO())

			if err != nil {
				fmt.Printf("error getting bare metal plan list : %v", err)
				os.Exit(1)
			}

			printer.PlanBareMetal(list)
		} else {
			list, err := client.Plan.List(context.TODO(), planType)

			if err != nil {
				fmt.Printf("error getting plan list : %v", err)
				os.Exit(1)
			}

			printer.Plan(list)
		}
	},
}
