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

// ReservedIP represents the reservedip command
func ReservedIP() *cobra.Command {
	reservedIPCmd := &cobra.Command{
		Use:     "reserved-ip",
		Aliases: []string{"rip"},
		Short:   "reserved-ip lets you interact with reserved-ip ",
		Long:    ``,
	}

	reservedIPCmd.AddCommand(reservedIPGet, reservedIPList, reservedIPDelete, reservedIPAttach, reservedIPDetach, reservedIPConvert, reservedIPCreate)

	// List
	reservedIPList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	reservedIPList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Attach
	reservedIPAttach.Flags().StringP("instance-id", "i", "", "id of instance you want to attach")
	reservedIPAttach.MarkFlagRequired("instance-id")

	// Convert
	reservedIPConvert.Flags().StringP("ip", "i", "", "ip you wish to convert")
	reservedIPConvert.MarkFlagRequired("ip")
	reservedIPConvert.Flags().StringP("label", "l", "", "label")

	// Create
	reservedIPCreate.Flags().StringP("region", "r", "", "id of region")
	reservedIPCreate.MarkFlagRequired("region")
	reservedIPCreate.Flags().StringP("type", "t", "", "type of IP : v4 or v6")
	reservedIPCreate.MarkFlagRequired("type")
	reservedIPCreate.Flags().StringP("label", "l", "", "label")

	return reservedIPCmd
}

var reservedIPGet = &cobra.Command{
	Use:   "get <reservedIPID",
	Short: "get a reserved IP",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reservedIP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		rip, err := client.ReservedIP.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting reserved IP : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(rip)
	},
}

var reservedIPList = &cobra.Command{
	Use:   "list",
	Short: "list all reserved IPs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		rip, meta, err := client.ReservedIP.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting reserved IPs : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIPList(rip, meta)
	},
}

var reservedIPDelete = &cobra.Command{
	Use:     "delete <reservedIPID>",
	Short:   "delete a reserved ip",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reservedIP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		if err := client.ReservedIP.Delete(context.Background(), ip); err != nil {
			fmt.Printf("error getting reserved IPs : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted reserved ip")
	},
}

var reservedIPAttach = &cobra.Command{
	Use:   "attach <reservedIPID>",
	Short: "attach a reservedIP to an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reservedIP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		instance, _ := cmd.Flags().GetString("instance-id")
		if err := client.ReservedIP.Attach(context.Background(), ip, instance); err != nil {
			fmt.Printf("error attaching reserved IPs : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Attached reserved ip")
	},
}

var reservedIPDetach = &cobra.Command{
	Use:   "detach <reservedIPID>",
	Short: "detach a reservedIP to an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reservedIP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		if err := client.ReservedIP.Detach(context.Background(), ip); err != nil {
			fmt.Printf("error detaching reserved IPs : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Detached reserved ip")
	},
}

var reservedIPConvert = &cobra.Command{
	Use:   "convert ",
	Short: "convert IP address to reservedIP",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		label, _ := cmd.Flags().GetString("label")
		options := &govultr.ReservedIPConvertReq{
			IPAddress: ip,
			Label:     label,
		}

		r, err := client.ReservedIP.Convert(context.Background(), options)
		if err != nil {
			fmt.Printf("error converting IP to reserved IPs : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(r)
	},
}

var reservedIPCreate = &cobra.Command{
	Use:   "create ",
	Short: "create reservedIP",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		ipType, _ := cmd.Flags().GetString("type")
		label, _ := cmd.Flags().GetString("label")

		options := &govultr.ReservedIPReq{
			Region: region,
			IPType: ipType,
			Label:  label,
		}

		r, err := client.ReservedIP.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating reserved IPs : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(r)
	},
}
