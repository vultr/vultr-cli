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

// DnsDomain represents the domain sub command
func DnsDomain() *cobra.Command {
	dnsDomainCmd := &cobra.Command{
		Use:   "domain",
		Short: "dns domain",
		Long:  ``,
	}

	dnsDomainCmd.AddCommand(domainCreate, domainDelete, secEnable, secInfo, domainList, soaInfo, soaUpdate)

	// Create
	domainCreate.Flags().StringP("domain", "d", "", "name of the domain")
	domainCreate.MarkFlagRequired("domain")
	domainCreate.Flags().StringP("ip", "i", "", "instance ip you want to assign this domain to")
	domainCreate.MarkFlagRequired("ip")

	//Dns Sec
	secEnable.Flags().StringP("enabled", "e", "", "set whether dns sec is enabled or not. true or false")
	secEnable.MarkFlagRequired("enabled")

	// Soa Update
	soaUpdate.Flags().StringP("ns-primary", "n", "", "primary nameserver to store in the SOA record")
	soaUpdate.MarkFlagRequired("ns-primary")
	soaUpdate.Flags().StringP("email", "e", "", "administrative email to store in the SOA record")

	return dnsDomainCmd
}

var domainCreate = &cobra.Command{
	Use:   "create",
	Short: "create a domain",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		instance, _ := cmd.Flags().GetString("ip")

		if err := client.DNSDomain.Create(context.TODO(), domain, instance); err != nil {
			fmt.Printf("error creating dns domain : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("created dns domain ")
	},
}

var domainDelete = &cobra.Command{
	Use:     "delete <domainName>",
	Short:   "delete a domain",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a domain name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		err := client.DNSDomain.Delete(context.TODO(), domain)

		if err != nil {
			fmt.Printf("error delete dns domain : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("deleted dns domain ")
	},
}

var secEnable = &cobra.Command{
	Use:   "dnssec <domainName>",
	Short: "enable/disable dnssec",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a domain name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		enabled, _ := cmd.Flags().GetString("enabled")

		enable := false
		if enabled == "true" {
			enable = true
		}

		err := client.DNSDomain.ToggleDNSSec(context.TODO(), domain, enable)

		if err != nil {
			fmt.Printf("error toggling dnssec : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("toggled dns sec")
	},
}

var secInfo = &cobra.Command{
	Use:   "dnssec-info <domainName>",
	Short: "get dns sec info",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a domain name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		info, err := client.DNSDomain.DNSSecInfo(context.TODO(), domain)

		if err != nil {
			fmt.Printf("error toggling dnssec : %v\n", err)
			os.Exit(1)
		}
		printer.SecInfo(info)
	},
}

var domainList = &cobra.Command{
	Use:   "list <domainName>",
	Short: "get dns sec info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := client.DNSDomain.List(context.TODO())

		if err != nil {
			fmt.Printf("error getting dns domains : %v\n", err)
			os.Exit(1)
		}
		printer.DomainList(list)
	},
}

var soaInfo = &cobra.Command{
	Use:   "soa-info <domainName>",
	Short: "get dns soa info",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a domain name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		info, err := client.DNSDomain.GetSoa(context.TODO(), domain)

		if err != nil {
			fmt.Printf("error toggling dnssec : %v\n", err)
			os.Exit(1)
		}
		printer.SoaInfo(info)
	},
}

var soaUpdate = &cobra.Command{
	Use:   "soa-update <domainName>",
	Short: "update soa for a domain",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a domain name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		nsPrimary, _ := cmd.Flags().GetString("ns-primary")
		email, _ := cmd.Flags().GetString("email")

		err := client.DNSDomain.UpdateSoa(context.TODO(), domain, nsPrimary, email)

		if err != nil {
			fmt.Printf("error toggling dnssec : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("updated SOA")
	},
}
