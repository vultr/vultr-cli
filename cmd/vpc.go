// Copyright © 2022 The Vultr-cli Authors
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
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	vpcLong    = `Access information about VPCs on the account and perform CRUD operations`
	vpcExample = `
	# Full example
	vultr-cli vpc
	`
	vpcGetLong    = `Display information for a specific VPC`
	vpcGetExample = `
	# Full example
	vultr-cli vpc get 9fd4dcf5-7108-4641-9969-b2b9a8f77990

	# Shortened example with aliases
	vultr-cli vpc g 9fd4dcf5-7108-4641-9969-b2b9a8f77990
	`
	vpcCreateLong    = `Create a new VPC with desired options`
	vpcCreateExample = `
	# Full example
	vultr-cli vpc create --region="ewr" --description="Example VPC" --subnet="10.200.0.0" --size=24

	--region is required.  Everything else is optional

	# Shortened example with aliases
	vultr-cli vpc c -r="ewr" -d="Example VPC" -s="10.200.0.0" -z=24
	`
	vpcUpdateLong    = `Update an existing VPC with the supplied information`
	vpcUpdateExample = `
	# Full example
	vultr-cli vpc update fe8cfe1d-b25c-4c3c-8dfe-e5784bade8d9 --description="Example Updated VPC"

	# Shortned example with aliases
	vultr-cli vpc u fe8cfe1d-b25c-4c3c-8dfe-e5784bade8d9 -d="Example Updated VPC"
	`
	vpcDeleteLong    = `Delete an existing VPC`
	vpcDeleteExample = `
	#Full example
	vultr-cli vpc delete 6b8d8af9-e74a-4829-850d-647f75a056ca

	#Shortened example with aliases
	vultr-cli vpc d 6b8d8af9-e74a-4829-850d-647f75a056ca
	`
	vpcListLong    = `List all available VPC information on the account`
	vpcListExample = `
	# Full example
	vultr-cli vpc list

	# Shortened example with aliases
	vultr-cli vpc l
	`
)

// VPC represents the vpc command
func VPC() *cobra.Command {
	vpcCmd := &cobra.Command{
		Use:     "vpc",
		Short:   "Interact with VPCs",
		Long:    vpcLong,
		Example: vpcExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Context().Value(ctxAuthKey{}).(bool) == false {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	vpcCmd.AddCommand(vpcGet, vpcList, vpcDelete, vpcCreate, vpcUpdate)
	vpcCreate.Flags().StringP("region", "r", "", "The ID of the region in which to create the VPC")
	vpcCreate.Flags().StringP("description", "d", "", "The description of the VPC")
	vpcCreate.Flags().StringP("subnet", "s", "", "The IPv4 VPC in CIDR notation.")
	vpcCreate.Flags().IntP("size", "z", 0, "The number of bits for the netmask in CIDR notation.")
	if err := vpcCreate.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking vpc create 'region' flag required: %v\n", err)
		os.Exit(1)
	}

	vpcList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	vpcList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	vpcUpdate.Flags().StringP("description", "d", "", "The description of the VPC")
	if err := vpcUpdate.MarkFlagRequired("description"); err != nil {
		fmt.Printf("error marking vpc update 'description' flag required: %v\n", err)
		os.Exit(1)
	}

	return vpcCmd
}

var vpcGet = &cobra.Command{
	Use:     "get <VPCID>",
	Aliases: []string{"g"},
	Short:   "get a VPC",
	Long:    vpcGetLong,
	Example: vpcGetExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a VPCID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		vpc, _, err := client.VPC.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting VPC : %v\n", err)
			os.Exit(1)
		}

		printer.VPC(vpc)
	},
}

var vpcList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all VPCs",
	Long:    vpcListLong,
	Example: vpcListExample,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		vpc, meta, _, err := client.VPC.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting VPC list : %v\n", err)
			os.Exit(1)
		}

		printer.VPCList(vpc, meta)
	},
}

var vpcDelete = &cobra.Command{
	Use:     "delete <VPCID>",
	Aliases: []string{"destroy", "d"},
	Short:   "delete a VPC",
	Long:    vpcDeleteLong,
	Example: vpcDeleteExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a VPCID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.VPC.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting VPC: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted VPC")
	},
}

var vpcCreate = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "create a VPC",
	Long:    vpcCreateLong,
	Example: vpcCreateExample,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		description, _ := cmd.Flags().GetString("description")
		subnet, _ := cmd.Flags().GetString("subnet")
		size, _ := cmd.Flags().GetInt("size")

		options := &govultr.VPCReq{
			Region:       region,
			Description:  description,
			V4Subnet:     subnet,
			V4SubnetMask: size,
		}

		vpc, _, err := client.VPC.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating VPC: %v\n", err)
			os.Exit(1)
		}

		printer.VPC(vpc)
	},
}

var vpcUpdate = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "update a VPC",
	Long:    vpcUpdateLong,
	Example: vpcUpdateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provid a VPC ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		description, _ := cmd.Flags().GetString("description")

		err := client.VPC.Update(context.Background(), args[0], description)
		if err != nil {
			fmt.Printf("error updating VPC: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("VPC updated")
	},
}
