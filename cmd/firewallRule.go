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
	"github.com/vultr/govultr"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// FirewallRule represents the firewall rule commands
func FirewallRule() *cobra.Command {
	firewallRuleCmd := &cobra.Command{
		Use:     "rule",
		Short:   "rule is used to access firewall rule commands",
		Aliases: []string{"r"},
	}

	firewallRuleCmd.AddCommand(firewallRuleCreate, firewallRuleDelete, firewallRuleList)

	firewallRuleCreate.Flags().StringP("id", "i", "", "ID of the target firewall group.")
	firewallRuleCreate.Flags().StringP("protocol", "p", "", "Protocol type. Possible values: 'icmp', 'tcp', 'udp', 'gre'.")
	firewallRuleCreate.Flags().StringP("cidr", "c", "", "The IP subnet and mask you want to apply the rule to.")
	firewallRuleCreate.Flags().StringP("port", "o", "", "(optional) TCP/UDP only. This field can be an integer value specifying a port or a colon separated port range.")
	firewallRuleCreate.Flags().StringP("notes", "n", "", "(optional) This field supports notes up to 255 characters.")

	firewallRuleCreate.MarkFlagRequired("id")
	firewallRuleCreate.MarkFlagRequired("protocol")
	firewallRuleCreate.MarkFlagRequired("cidr")

	firewallRuleList.Flags().StringP("type", "t", "", "(optional) IP address type. Possible values: 'v4', 'v6'.")

	return firewallRuleCmd
}

var firewallRuleCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a firewall rule",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		protocol, _ := cmd.Flags().GetString("protocol")
		cidr, _ := cmd.Flags().GetString("cidr")
		port, _ := cmd.Flags().GetString("port")
		notes, _ := cmd.Flags().GetString("notes")

		fwr, err := client.FirewallRule.Create(context.TODO(), id, protocol, port, cidr, notes)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("created firewall rule number: %d\n", fwr.RuleNumber)
	},
}

var firewallRuleDelete = &cobra.Command{
	Use:     "delete <firewallGroupID> <firewallRuleNumber>",
	Short:   "Delete a firewall rule",
	Aliases: []string{"d", "destroy"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a firewallGroupID and firewallRuleNumber")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.FirewallRule.Delete(context.TODO(), args[0], args[1])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Firewall rule has been deleted")
	},
}

var firewallRuleList = &cobra.Command{
	Use:     "list <firewallGroupID>",
	Short:   "List all firewall rules",
	Aliases: []string{"l"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a firewallGroupID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ipType, _ := cmd.Flags().GetString("type")

		var err error
		var list []govultr.FirewallRule
		if ipType != "" {
			list, err = client.FirewallRule.ListByIPType(context.TODO(), args[0], ipType)
		} else {
			list, err = client.FirewallRule.List(context.TODO(), args[0])
		}

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallRule(list)
	},
}
