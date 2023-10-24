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
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	vpc2Long    = `Get commands available to vpc2`
	vpc2Example = `
	# Full example
	vultr-cli vpc2
	`
	vpc2CreateLong    = `Create a new VPC 2.0 network with specified region, description, and network settings`
	vpc2CreateExample = `
	# Full example
	vultr-cli vpc2 create --region="ewr" --description="example-vpc" --ip-type="v4" --ip-block="10.99.0.0" --prefix-length="24"
	`
	vpc2UpdateLong    = `Updates a VPC 2.0 network with the supplied information`
	vpc2UpdateExample = `
	# Full example
	vultr-cli vpc2 update 84fee086-6691-417a-b2db-e2a71061fa17 --description="example-vpc"
	`
	vpc2NodesAttachLong    = `Attaches multiple nodes to a VPC 2.0 network`
	vpc2NodesAttachExample = `
	# Full example
	vultr-cli vpc2 nodes attach 84fee086-6691-417a-b2db-e2a71061fa17 \
		--nodes="35dbcffe-58bf-46fe-bd68-964d95488dd8,1f5d784a-1011-430c-a2e2-39ba045abe3c"
	`
	vpc2NodesDetachLong    = `Detaches multiple nodes from a VPC 2.0 network`
	vpc2NodesDetachExample = `
	# Full example
	vultr-cli vpc2 nodes detach 84fee086-6691-417a-b2db-e2a71061fa17 \
		--nodes="35dbcffe-58bf-46fe-bd68-964d95488dd8,1f5d784a-1011-430c-a2e2-39ba045abe3c"
	`
)

// VPC2 represents the VPC2 command
func VPC2() *cobra.Command {
	vpc2Cmd := &cobra.Command{
		Use:     "vpc2",
		Short:   "commands to interact with vpc2 on vultr",
		Long:    vpc2Long,
		Example: vpc2Example,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Context().Value(ctxAuthKey{}).(bool) {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	vpc2Cmd.AddCommand(vpc2List, vpc2Create, vpc2Info, vpc2Update, vpc2Delete)

	// VPC2 list flags
	vpc2List.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	vpc2List.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	// VPC2 create flags
	vpc2Create.Flags().StringP("region", "r", "", "region id for the new vpc network")
	vpc2Create.Flags().StringP("description", "d", "", "description for the new vpc network")
	vpc2Create.Flags().StringP("ip-type", "", "", "IP type for the new vpc network")
	vpc2Create.Flags().StringP("ip-block", "", "", "subnet IP address for the new vpc network")
	vpc2Create.Flags().IntP("prefix-length", "", 0, "number of bits for the netmask in CIDR notation for the new vpc network")

	// VPC2 update flags
	vpc2Update.Flags().StringP("description", "d", "", "description for the vpc network")

	// VPC2 node commands
	vpc2NodeCmd := &cobra.Command{
		Use:   "nodes",
		Short: "commands to handle nodes attached to a vpc 2.0 network",
		Long:  ``,
	}
	vpc2NodeCmd.AddCommand(vpc2NodesList, vpc2NodesAttach, vpc2NodesDetach)
	vpc2NodesList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	vpc2NodesList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)
	vpc2NodesAttach.Flags().StringSliceP("nodes", "n", []string{}, "instance ids you wish to attach to the vpc network")
	vpc2NodesDetach.Flags().StringSliceP("nodes", "n", []string{}, "instance ids you wish to detach from the vpc network")
	vpc2Cmd.AddCommand(vpc2NodeCmd)

	return vpc2Cmd
}

var vpc2List = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all available VPC 2.0 networks",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		s, meta, _, err := client.VPC2.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting list of vpc 2.0 networks : %v\n", err)
			os.Exit(1)
		}

		printer.VPC2List(s, meta)
	},
}

var vpc2Create = &cobra.Command{
	Use:     "create",
	Short:   "Create a VPC 2.0 network",
	Aliases: []string{"c"},
	Long:    vpc2CreateLong,
	Example: vpc2CreateExample,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")

		// Optional
		description, _ := cmd.Flags().GetString("description")
		ipType, _ := cmd.Flags().GetString("ip-type")
		ipBlock, _ := cmd.Flags().GetString("ip-block")
		prefixLength, _ := cmd.Flags().GetInt("prefix-length")

		opt := &govultr.VPC2Req{
			Region:       region,
			Description:  description,
			IPType:       ipType,
			IPBlock:      ipBlock,
			PrefixLength: prefixLength,
		}

		// Make the request
		vpc2, _, err := client.VPC2.Create(context.TODO(), opt)
		if err != nil {
			fmt.Printf("error creating VPC 2.0 network : %v\n", err)
			os.Exit(1)
		}

		printer.VPC2(vpc2)
	},
}

var vpc2Info = &cobra.Command{
	Use:   "get <vpc2ID>",
	Short: "get info about a specific VPC 2.0 network",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a vpc2ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.VPC2.Get(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting VPC 2.0 network : %v\n", err)
			os.Exit(1)
		}

		printer.VPC2(s)
	},
}

var vpc2Update = &cobra.Command{
	Use:     "update <vpc2ID>",
	Short:   "Update a VPC 2.0 network",
	Aliases: []string{"u"},
	Long:    vpc2UpdateLong,
	Example: vpc2UpdateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a vpc2ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		description, _ := cmd.Flags().GetString("description")

		// Make the request
		err := client.VPC2.Update(context.TODO(), args[0], description)
		if err != nil {
			fmt.Printf("error updating VPC 2.0 network : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("VPC 2.0 network updated")
	},
}

var vpc2Delete = &cobra.Command{
	Use:     "delete <vpc2ID>",
	Short:   "delete/destroy a VPC 2.0 network",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a vpc2ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.VPC2.Delete(context.Background(), args[0]); err != nil {
			fmt.Printf("error deleting VPC 2.0 network : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted VPC 2.0 network")
	},
}

var vpc2NodesList = &cobra.Command{
	Use:     "list <vpc2ID>",
	Aliases: []string{"l"},
	Short:   "list all nodes attached to a VPC 2.0 network",
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a vpc2ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		s, meta, _, err := client.VPC2.ListNodes(context.TODO(), args[0], options)
		if err != nil {
			fmt.Printf("error getting list of vpc 2.0 network nodes : %v\n", err)
			os.Exit(1)
		}

		printer.VPC2ListNodes(s, meta)
	},
}

var vpc2NodesAttach = &cobra.Command{
	Use:     "attach <vpc2ID>",
	Short:   "Attach nodes to a VPC 2.0 network",
	Aliases: []string{"a"},
	Long:    vpc2NodesAttachLong,
	Example: vpc2NodesAttachExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a vpc2ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		nodes, _ := cmd.Flags().GetStringSlice("nodes")

		opt := &govultr.VPC2AttachDetachReq{
			Nodes: nodes,
		}

		// Make the request
		err := client.VPC2.Attach(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error attaching nodes to VPC 2.0 network : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Nodes attached to VPC 2.0 network")
	},
}

var vpc2NodesDetach = &cobra.Command{
	Use:     "detach <vpc2ID>",
	Short:   "Detach nodes from a VPC 2.0 network",
	Aliases: []string{"d"},
	Long:    vpc2NodesDetachLong,
	Example: vpc2NodesDetachExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a vpc2ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		nodes, _ := cmd.Flags().GetStringSlice("nodes")

		opt := &govultr.VPC2AttachDetachReq{
			Nodes: nodes,
		}

		// Make the request
		err := client.VPC2.Detach(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error detaching nodes from VPC 2.0 network : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Nodes detached from VPC 2.0 network")
	},
}
