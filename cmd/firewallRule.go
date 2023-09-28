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
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	firewallRuleLong    = `Show commands available for firewall rules`
	firewallRuleExample = `
	# Full example
	vultr-cli firewall rule

	# Shortened example with aliases
	vultr-cli fw r
	`
	firewallRuleCreateLong = `
	Create a new firewall rule in the provided firewall group

	If protocol is TCP or UDP, port must be provided.

	An ip-type of v4 or v6 must be supplied for all rules.
	`
	firewallRuleCreateExample = `
	# Full examples
	vultr-cli firewall rule create --id=f04ae5aa-ff6a-4078-900d-78cc17dca2d5 --ip-type=v4 --protocol=tcp --size=24 \
		--subnet=127.0.0.0 --port=30000

	vultr-cli firewall rule create --id=f04ae5aa-ff6a-4078-900d-78cc17dca2d5 --ip-type=v4 --protocol=icmp --size=24 --subnet=127.0.0.0

	# Shortened example with aliases
	vultr-cli fw r c -i=f04ae5aa-ff6a-4078-900d-78cc17dca2d5 -t=v4 -p=tcp -z=24 -s=127.0.0.0 -r=30000
	`
	firewallRuleDeleteLong    = `Delete a firewall rule in the provided firewall group`
	firewallRuleDeleteExample = `
	# Full example
	vultr-cli firewall rule delete 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3

	# Shortened example with aliases
	vultr-cli fw r d 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3
	`
	firewallRuleGetLong    = `Get a firewall rule in the provided firewall group`
	firewallRuleGetExample = `
	# Full example
	vultr-cli firewall rule get 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3

	# Shortened example with aliases
	vultr-cli fw r get 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3
	`
	firewallRuleListLong    = `List all firewall rules in the provided firewall group`
	firewallRuleListExample = `
	# Full example
	vultr-cli firewall rule list 704ac064-4ff2-49ca-a6e6-88262cca8f8a

	# Shortened example with aliases
	vultr-cli fw r l 704ac064-4ff2-49ca-a6e6-88262cca8f8a
	`
)

// FirewallRule represents the firewall rule commands
func FirewallRule() *cobra.Command {
	firewallRuleCmd := &cobra.Command{
		Use:     "rule",
		Short:   "rule is used to access firewall rule commands",
		Long:    firewallRuleLong,
		Example: firewallRuleExample,
		Aliases: []string{"r"},
	}

	firewallRuleCmd.AddCommand(firewallRuleCreate, firewallRuleDelete, firewallRuleGet, firewallRuleList)

	firewallRuleCreate.Flags().StringP("id", "i", "", "ID of the target firewall group.")
	firewallRuleCreate.Flags().StringP("protocol", "p", "", "Protocol type. Possible values: 'icmp', 'tcp', 'udp', 'gre'.")
	firewallRuleCreate.Flags().StringP("subnet", "s", "", "The IPv4 network in CIDR notation.")
	firewallRuleCreate.Flags().IntP("size", "z", 0, "The number of bits for the netmask in CIDR notation.")
	firewallRuleCreate.Flags().StringP(
		"source",
		"o",
		"",
		"(optional) When empty, uses value from subnet and size. If \"cloudflare\", allows all Cloudflare IP space through firewall.",
	)
	firewallRuleCreate.Flags().StringP("type", "", "", "Deprecated: use ip-type instead. The type of IP rule - v4 or v6.")
	firewallRuleCreate.Flags().StringP("ip-type", "t", "", "The type of IP rule - v4 or v6.")

	firewallRuleCreate.Flags().StringP(
		"port",
		"r",
		"",
		"(optional) TCP/UDP only. This field can be an integer value specifying a port or a colon separated port range.",
	)
	firewallRuleCreate.Flags().StringP("notes", "n", "", "(optional) This field supports notes up to 255 characters.")

	if err := firewallRuleCreate.MarkFlagRequired("id"); err != nil {
		fmt.Printf("error marking firewall rule create  'id' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := firewallRuleCreate.MarkFlagRequired("protocol"); err != nil {
		fmt.Printf("error marking firewall rule create 'protocol' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := firewallRuleCreate.MarkFlagRequired("subnet"); err != nil {
		fmt.Printf("error marking firewall rule create 'subnet' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := firewallRuleCreate.MarkFlagRequired("size"); err != nil {
		fmt.Printf("error marking firewall rule create 'size' flag required: %v\n", err)
		os.Exit(1)
	}

	firewallRuleList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	firewallRuleList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	return firewallRuleCmd
}

var firewallRuleCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a firewall rule",
	Long:    firewallRuleCreateLong,
	Example: firewallRuleCreateExample,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		protocol, _ := cmd.Flags().GetString("protocol")
		subnet, _ := cmd.Flags().GetString("subnet")
		ipTypeOld, _ := cmd.Flags().GetString("type")
		ipType, _ := cmd.Flags().GetString("ip-type")
		size, _ := cmd.Flags().GetInt("size")
		source, _ := cmd.Flags().GetString("source")
		port, _ := cmd.Flags().GetString("port")
		notes, _ := cmd.Flags().GetString("notes")

		options := &govultr.FirewallRuleReq{
			Protocol:   protocol,
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

		if ipTypeOld == "" && ipType == "" {
			fmt.Println("a firewall rule requires an IP type. Pass an --ip-type value of v4 or v6")
			os.Exit(1)
		}

		if ipTypeOld != "" && ipType != "" {
			fmt.Println("--type is deprecated. Instead, use only --ip-type")
			os.Exit(1)
		}

		if ipType != "" {
			options.IPType = ipType
		}

		if ipTypeOld != "" {
			options.IPType = ipTypeOld
		}

		fwr, _, err := client.FirewallRule.Create(context.Background(), id, options)
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
	Long:    firewallRuleDeleteLong,
	Example: firewallRuleDeleteExample,
	Aliases: []string{"d", "destroy"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a firewallGroupID and firewallRuleNumber")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		rule, _ := strconv.Atoi(args[1])
		if err := client.FirewallRule.Delete(context.Background(), args[0], rule); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Firewall rule has been deleted")
	},
}

var firewallRuleGet = &cobra.Command{
	Use:     "get <firewallGroupID> <firewallRuleNumber>",
	Short:   "Get firewall rule",
	Long:    firewallRuleGetLong,
	Example: firewallRuleGetExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a firewallGroupID and firewallRuleNumber")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ruleNumber, _ := strconv.Atoi(args[1])
		fwRule, _, err := client.FirewallRule.Get(context.Background(), args[0], ruleNumber)
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
	Long:    firewallRuleListLong,
	Example: firewallRuleListExample,
	Aliases: []string{"l"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a firewallGroupID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, _, err := client.FirewallRule.List(context.Background(), args[0], options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.FirewallRules(list, meta)
	},
}
