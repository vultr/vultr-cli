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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/cmd/printer"
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
	return regionCmd
}

var regionList = &cobra.Command{
	Use:   "list",
	Short: "list regions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := client.Region.List(context.TODO())

		if err != nil {
			fmt.Printf("error getting region list : %v\n", err)
			os.Exit(1)
		}

		printer.Regions(list)
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

		var availability []int
		var err error

		region, _ := strconv.Atoi(regionID)

		switch availType {
		case "vc2":
			availability, err = client.Region.Vc2Availability(context.TODO(), region)

			if err != nil {
				fmt.Printf("error getting VC2 availability : %v\n", err)
				os.Exit(1)
			}
		case "v2c2":
			availability, err = client.Region.Vdc2Availability(context.TODO(), region)

			if err != nil {
				fmt.Printf("error getting VDC2 availability : %v\n", err)
				os.Exit(1)
			}

		case "bare-metal":
			availability, err = client.Region.BareMetalAvailability(context.TODO(), region)

			if err != nil {
				fmt.Printf("error getting bare-metal availability : %v\n", err)
				os.Exit(1)
			}

		default:
			availability, err = client.Region.Availability(context.TODO(), region, "")

			if err != nil {
				fmt.Printf("error getting availability : %v\n", err)
				os.Exit(1)
			}
		}

		printer.RegionAvailability(availability)
	},
}
