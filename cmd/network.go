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
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// networkCmd represents the network command
func Network() *cobra.Command {
	networkCmd := &cobra.Command{
		Use:   "network",
		Short: "network interacts with network actions",
		Long:  ``,
	}

	networkCmd.AddCommand(networkGet, networkList, networkDelete, networkCreate)
	networkCreate.Flags().StringP("region-id", "r", "", "id of the region you wish to create the network")
	networkCreate.Flags().StringP("description", "d", "", "description of the network")
	networkCreate.Flags().StringP("subnet", "s", "", "The IPv4 network in CIDR notation.")
	networkCreate.Flags().IntP("size", "z", 0, "The number of bits for the netmask in CIDR notation.")
	networkCreate.MarkFlagRequired("region-id")
	networkCreate.MarkFlagRequired("description")
	networkCreate.MarkFlagRequired("subnet")
	networkCreate.MarkFlagRequired("size")

	networkList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	networkList.Flags().IntP("per-page", "p", 25, "(optional) Number of items requested per page. Default and Max are 25.")

	return networkCmd
}

var networkGet = &cobra.Command{
	Use:   "get <networkID>",
	Short: "get a private network",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a networkID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		network, err := client.Network.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting network : %v\n", err)
			os.Exit(1)
		}

		printer.Network(network)
	},
}

var networkList = &cobra.Command{
	Use:   "list",
	Short: "list all private networks",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		network, meta, err := client.Network.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting network list : %v\n", err)
			os.Exit(1)
		}

		printer.NetworkList(network, meta)
	},
}

var networkDelete = &cobra.Command{
	Use:     "delete <networkID>",
	Short:   "delete a private network",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a networkID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Network.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting network : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted network")
	},
}

var networkCreate = &cobra.Command{
	Use:   "create",
	Short: "create a private network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region-id")
		description, _ := cmd.Flags().GetString("description")
		subnet, _ := cmd.Flags().GetString("subnet")
		size, _ := cmd.Flags().GetInt("size")

		options := &govultr.NetworkReq{
			Region:       region,
			Description:  description,
			V4Subnet:     subnet,
			V4SubnetMask: size,
		}

		network, err := client.Network.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating network : %v\n", err)
			os.Exit(1)
		}

		printer.Network(network)
	},
}
