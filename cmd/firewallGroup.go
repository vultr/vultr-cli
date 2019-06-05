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
	"github.com/vultr/vultr-cli/cmd/printer"
)

// FirewallGroup represents the firewall group commands
func FirewallGroup() *cobra.Command {
	firewallGroupCmd := &cobra.Command{
		Use:     "group",
		Short:   "group is used to access firewall group commands",
		Long:    ``,
		Aliases: []string{"g"},
	}

	firewallGroupCmd.AddCommand(firewallGroupCreate, firewallGroupDelete, firewallGroupUpdate, firewallGroupList)

	firewallGroupCreate.Flags().StringP("description", "d", "", "(optional) Description of firewall group.")

	return firewallGroupCmd
}

var firewallGroupCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a firewall group",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		description, _ := cmd.Flags().GetString("description")

		fwg, err := client.FirewallGroup.Create(context.TODO(), description)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("created firewall group: %s\n", fwg.FirewallGroupID)
	},
}

var firewallGroupDelete = &cobra.Command{
	Use:     "delete <firewallGroupID>",
	Short:   "Delete a firewall group",
	Aliases: []string{"d"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a firewallGroupID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.FirewallGroup.Delete(context.TODO(), args[0])
		if err != nil {
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
		err := client.FirewallGroup.ChangeDescription(context.TODO(), args[0], args[1])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Firewall group has been updated")
	},
}

var firewallGroupList = &cobra.Command{
	Use:     "list",
	Short:   "List all firewall groups",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		list, err := client.FirewallGroup.List(context.TODO())

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallGroup(list)
	},
}
