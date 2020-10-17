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
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// FirewallRule represents the firewall rule commands
func FirewallRule() *cobra.Command {
	firewallRuleCmd := &cobra.Command{
		Use:     "rule",
		Short:   "rule is used to access firewall rule commands",
		Aliases: []string{"r"},
	}

	firewallRuleCmd.AddCommand(firewallRuleCreate, firewallRuleDelete, firewallRuleGet, firewallRuleList)

	firewallRuleCreate.Flags().StringP("id", "i", "", "ID of the target firewall group.")
	firewallRuleCreate.Flags().StringP("protocol", "p", "", "Protocol type. Possible values: 'icmp', 'tcp', 'udp', 'gre'.")
	firewallRuleCreate.Flags().StringP("subnet", "s", "", "The IPv4 network in CIDR notation.")
	firewallRuleCreate.Flags().IntP("size", "z", 0, "The number of bits for the netmask in CIDR notation.")
	firewallRuleCreate.Flags().IntP("source", "o", 0, "(optional) When empty, uses value from subnet and size. If \"cloudflare\", allows all Cloudflare IP space through firewall.")
	firewallRuleCreate.Flags().StringP("type", "t", "", "The type of IP rule - v4 or v6.")

	firewallRuleCreate.Flags().StringP("port", "r", "", "(optional) TCP/UDP only. This field can be an integer value specifying a port or a colon separated port range.")
	firewallRuleCreate.Flags().StringP("notes", "n", "", "(optional) This field supports notes up to 255 characters.")

	firewallRuleCreate.MarkFlagRequired("id")
	firewallRuleCreate.MarkFlagRequired("protocol")
	firewallRuleCreate.MarkFlagRequired("subnet")
	firewallRuleCreate.MarkFlagRequired("size")
	firewallRuleCreate.MarkFlagRequired("type")

	firewallRuleList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	firewallRuleList.Flags().IntP("per-page", "p", 25, "(optional) Number of items requested per page. Default and Max are 25.")

	return firewallRuleCmd
}

var firewallRuleCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a firewall rule",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		protocol, _ := cmd.Flags().GetString("protocol")
		subnet, _ := cmd.Flags().GetString("subnet")
		ipType, _ := cmd.Flags().GetString("type")
		size, _ := cmd.Flags().GetInt("size")
		source, _ := cmd.Flags().GetString("source")
		port, _ := cmd.Flags().GetString("port")
		notes, _ := cmd.Flags().GetString("notes")

		options := &govultr.FirewallRuleReq{
			Protocol:   protocol,
			IPType:     ipType,
			Subnet:     subnet,
			SubnetSize: size,
			Notes:      notes,
		}

		if port != "" {
			options.Port = port
		}

		if source != "" {
			options.Source = source
		}

		fwr, err := client.FirewallRule.Create(context.TODO(), id, options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallRule(fwr)
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
		rule, _ := strconv.Atoi(args[1])
		if err := client.FirewallRule.Delete(context.TODO(), args[0], rule); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Firewall rule has been deleted")
	},
}

var firewallRuleGet = &cobra.Command{
	Use:   "get <firewallGroupID> <firewallRuleNumber>",
	Short: "Get firewall rule",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a firewallGroupID and firewallRuleNumber")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ruleNumber, _ := strconv.Atoi(args[1])
		fwRule, err := client.FirewallRule.Get(context.TODO(), args[0], ruleNumber)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallRule(fwRule)
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
		options := getPaging(cmd)
		list, meta, err := client.FirewallRule.List(context.TODO(), args[0], options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallRules(list, meta)
	},
}
