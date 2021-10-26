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

// Regions represents the region command
func Regions() *cobra.Command {
	regionCmd := &cobra.Command{
		Use:   "regions",
		Short: "get regions",
		Long:  `regions lets you get information on all data centers`,
	}

	regionCmd.AddCommand(regionList)
	regionCmd.AddCommand(regionAvailability)

	regionAvailability.Flags().StringP("type", "t", "", "type of plans for which to include availability. Possible values: vc2, vdc2, bare-metal")

	regionList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	regionList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return regionCmd
}

var regionList = &cobra.Command{
	Use:   "list",
	Short: "list regions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, err := client.Region.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting region list : %v\n", err)
			os.Exit(1)
		}

		printer.Regions(list, meta)
	},
}

var regionAvailability = &cobra.Command{
	Use:   "availability <regionID>",
	Short: "list availability in region",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a regionID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		regionID := args[0]
		availType, _ := cmd.Flags().GetString("type")
		availability, err := client.Region.Availability(context.Background(), regionID, availType)
		if err != nil {
			fmt.Printf("error getting availability : %v\n", err)
			os.Exit(1)
		}

		printer.RegionAvailability(availability)
	},
}
