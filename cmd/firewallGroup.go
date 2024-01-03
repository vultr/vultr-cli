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
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// FirewallGroup represents the firewall group commands
func FirewallGroup() *cobra.Command {
	firewallGroupCmd := &cobra.Command{
		Use:     "group",
		Short:   "group is used to access firewall group commands",
		Long:    ``,
		Aliases: []string{"g"},
	}

	firewallGroupCmd.AddCommand(firewallGroupCreate, firewallGroupDelete, firewallGroupGet, firewallGroupUpdate, firewallGroupList)

	firewallGroupCreate.Flags().StringP("description", "d", "", "(optional) Description of firewall group.")

	firewallGroupList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	firewallGroupList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	return firewallGroupCmd
}

var firewallGroupCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a firewall group",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		description, _ := cmd.Flags().GetString("description")
		options := &govultr.FirewallGroupReq{
			Description: description,
		}

		fwg, _, err := client.FirewallGroup.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallGroup(fwg)
	},
}

var firewallGroupDelete = &cobra.Command{
	Use:     "delete <firewallGroupID>",
	Short:   "Delete a firewall group",
	Aliases: []string{"d", "destroy"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a firewallGroupID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.FirewallGroup.Delete(context.Background(), args[0]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Firewall group has been deleted")
	},
}

var firewallGroupUpdate = &cobra.Command{
	Use:     "update <firewallGroupID> <description>",
	Short:   "Update firewall group description",
	Aliases: []string{"u"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a firewallGroupID and description")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		description := args[1]
		options := &govultr.FirewallGroupReq{
			Description: description,
		}

		if err := client.FirewallGroup.Update(context.Background(), args[0], options); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Firewall group has been updated")
	},
}

var firewallGroupGet = &cobra.Command{
	Use:   "get <firewallGroupID>",
	Short: "Get firewall group",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a firewallGroupID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, _, err := client.FirewallGroup.List(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallGroups(list, meta)
	},
}

var firewallGroupList = &cobra.Command{
	Use:     "list",
	Short:   "List all firewall groups",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, _, err := client.FirewallGroup.List(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallGroups(list, meta)
	},
}
