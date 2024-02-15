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

// DNSDomain represents the domain sub command
func DNSDomain() *cobra.Command {
	dnsDomainCmd := &cobra.Command{
		Use:   "domain",
		Short: "dns domain",
		Long:  ``,
	}

	dnsDomainCmd.AddCommand(domainCreate, domainGet, domainDelete, secEnable, secInfo, domainList, soaInfo, soaUpdate)

	// Create
	domainCreate.Flags().StringP("domain", "d", "", "name of the domain")
	domainCreate.MarkFlagRequired("domain")
	domainCreate.Flags().StringP("ip", "i", "", "instance ip you want to assign this domain to")

	// Dns Sec
	secEnable.Flags().StringP("enabled", "e", "", "set whether dns sec is enabled or not. true or false")
	secEnable.MarkFlagRequired("enabled")

	// Soa Update
	soaUpdate.Flags().StringP("ns-primary", "n", "", "primary nameserver to store in the SOA record")
	soaUpdate.MarkFlagRequired("ns-primary")
	soaUpdate.Flags().StringP("email", "e", "", "administrative email to store in the SOA record")

	// List
	domainList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	domainList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return dnsDomainCmd
}

var domainCreate = &cobra.Command{
	Use:   "create",
	Short: "create a domain",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		ip, _ := cmd.Flags().GetString("ip")

		options := &govultr.DomainReq{
			Domain: domain,
			IP:     ip,
		}

		dns, err := client.Domain.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating dns domain : %v\n", err)
			os.Exit(1)
		}

		printer.Domain(dns)
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
		if err := client.Domain.Delete(context.Background(), domain); err != nil {
			fmt.Printf("error delete dns domain : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("deleted dns domain")
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
		if err := client.Domain.Update(context.Background(), domain, enabled); err != nil {
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
		info, err := client.Domain.GetDNSSec(context.Background(), domain)
		if err != nil {
			fmt.Printf("error getting dnssec info : %v\n", err)
			os.Exit(1)
		}

		printer.SecInfo(info)
	},
}

var domainGet = &cobra.Command{
	Use:   "get <domainName>",
	Short: "get a domain",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		domain, err := client.Domain.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting domain : %v\n", err)
			os.Exit(1)
		}

		printer.Domain(domain)
	},
}

var domainList = &cobra.Command{
	Use:   "list",
	Short: "get list of domains",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, err := client.Domain.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting domains : %v\n", err)
			os.Exit(1)
		}

		printer.DomainList(list, meta)
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
		info, err := client.Domain.GetSoa(context.Background(), domain)
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

		soaUpdate := &govultr.Soa{
			NSPrimary: nsPrimary,
			Email:     email,
		}

		if err := client.Domain.UpdateSoa(context.Background(), domain, soaUpdate); err != nil {
			fmt.Printf("error toggling dnssec : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("updated SOA")
	},
}
