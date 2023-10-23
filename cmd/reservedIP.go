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
	reservedIPLong    = `Get all available commands for reserved IPs`
	reservedIPExample = `
	# Full example
	vultr-cli reserved-ip

	# Shortened with aliased commands
	vultr-cli rip
	`

	reservedIPCreateLong    = `Create a reserved IP on your Vultr account`
	reservedIPCreateExample = `
	# Full Example
	vultr-cli reserved-ip create --region="yto" --type="v4" --label="new IP"

	# Shortened with alias commands
	vultr-cli rip c -r="yto" -t="v4" -l="new IP"
	`

	reservedIPGetLong    = `Get info for a reserved IP on your Vultr account`
	reservedIPGetExample = `
	# Full example
	vultr-cli reserved-ip get 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5

	# Shortened with alias commands
	vultr-cli rip g 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5
	`

	reservedIPListLong    = `List all reserved IPs on your Vultr account`
	reservedIPListExample = `
	# Full example
	vultr-cli reserved-ip list

	# Shortened with alias commands
	vultr-cli rip l
	`

	reservedIPAttachLong    = `Attach a reserved IP to an instance on your Vultr account`
	reservedIPAttachExample = `
	# Full example
	vultr-cli reserved-ip attach 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 --instance-id="2b9bf5fb-1644-4e0a-b706-1116ab64d783"

	# Shortened with alias commands
	vultr-cli rip a 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 -i="2b9bf5fb-1644-4e0a-b706-1116ab64d783"
	`

	reservedIPDetachLong    = `Detach a reserved IP from an instance on your Vultr account`
	reservedIPDetachExample = `
	# Full example
	vultr-cli reserved-ip detach 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5

	# Shortened with alias commands
	vultr-cli rip d 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5
	`

	reservedIPConvertLong    = `Convert an instance IP to a reserved IP on your Vultr account`
	reservedIPConvertExample = `
	# Full example
	vultr-cli reserved-ip convert --ip="192.0.2.123" --label="new label converted"

	# Shortened with alias commands
	vultr-cli rip v -i="192.0.2.123" -l="new label converted"
	`

	reservedIPUpdateLong    = `Update a reserved IP on your Vultr account`
	reservedIPUpdateExample = `
	# Full example
	vultr-cli reserved-ip update 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 --label="new label"

	# Shortened with alias commands
	vultr-cli rip u 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 -l="new label"
	`

	reservedIPDeleteLong    = `Delete a reserved IP from your Vultr account`
	reservedIPDeleteExample = `
	# Full example
	vultr-cli reserved-ip delete 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5
	`
)

// ReservedIP represents the reservedip command
func ReservedIP() *cobra.Command {
	reservedIPCmd := &cobra.Command{
		Use:     "reserved-ip",
		Aliases: []string{"rip"},
		Short:   "reserved-ip lets you interact with reserved-ip ",
		Long:    reservedIPLong,
		Example: reservedIPExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if auth := cmd.Context().Value("authenticated"); auth != true {
				return fmt.Errorf(apiKeyError)
			}
			return nil
		},
	}

	reservedIPCmd.AddCommand(
		reservedIPGet,
		reservedIPList,
		reservedIPDelete,
		reservedIPAttach,
		reservedIPDetach,
		reservedIPConvert,
		reservedIPCreate,
		reservedIPUpdate,
	)

	// List
	reservedIPList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	reservedIPList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	// Attach
	reservedIPAttach.Flags().StringP("instance-id", "i", "", "id of instance you want to attach")
	if err := reservedIPAttach.MarkFlagRequired("instance-id"); err != nil {
		fmt.Printf("error marking reserved-ip attach 'instance-id' flag required: %v\n", err)
		os.Exit(1)
	}

	// Convert
	reservedIPConvert.Flags().StringP("ip", "i", "", "ip you wish to convert")
	if err := reservedIPConvert.MarkFlagRequired("ip"); err != nil {
		fmt.Printf("error marking reserved-ip convert 'ip' flag required: %v\n", err)
		os.Exit(1)
	}
	reservedIPConvert.Flags().StringP("label", "l", "", "label")

	// Create
	reservedIPCreate.Flags().StringP("region", "r", "", "id of region")
	if err := reservedIPCreate.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking reserved-ip create 'region' flag required: %v\n", err)
		os.Exit(1)
	}
	reservedIPCreate.Flags().StringP("type", "t", "", "type of IP : v4 or v6")
	if err := reservedIPCreate.MarkFlagRequired("type"); err != nil {
		fmt.Printf("error marking reserved-ip create 'type' flag required: %v\n", err)
		os.Exit(1)
	}
	reservedIPCreate.Flags().StringP("label", "l", "", "label")

	// Update
	reservedIPUpdate.Flags().StringP("label", "l", "", "label")
	if err := reservedIPUpdate.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking reserved-ip update 'label' flag required: %v\n", err)
		os.Exit(1)
	}

	return reservedIPCmd
}

var reservedIPGet = &cobra.Command{
	Use:     "get <reservedIPID>",
	Short:   "get a reserved IP",
	Long:    reservedIPGetLong,
	Example: reservedIPGetExample,
	Aliases: []string{"g"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reservedIP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		rip, _, err := client.ReservedIP.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting reserved IP : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(rip)
	},
}

var reservedIPList = &cobra.Command{
	Use:     "list",
	Short:   "list all reserved IPs",
	Long:    reservedIPListLong,
	Example: reservedIPListExample,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		rip, meta, _, err := client.ReservedIP.List(context.Background(), options)
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
	Long:    reservedIPDeleteLong,
	Example: reservedIPDeleteExample,
	Aliases: []string{"destroy"},
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
	Use:     "attach <reservedIPID>",
	Short:   "attach a reservedIP to an instance",
	Long:    reservedIPAttachLong,
	Example: reservedIPAttachExample,
	Aliases: []string{"a"},
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
	Use:     "detach <reservedIPID>",
	Short:   "detach a reservedIP to an instance",
	Long:    reservedIPDetachLong,
	Example: reservedIPDetachExample,
	Aliases: []string{"d"},
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
	Use:     "convert ",
	Short:   "convert IP address to reservedIP",
	Long:    reservedIPConvertLong,
	Example: reservedIPConvertExample,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		label, _ := cmd.Flags().GetString("label")
		options := &govultr.ReservedIPConvertReq{
			IPAddress: ip,
			Label:     label,
		}

		r, _, err := client.ReservedIP.Convert(context.Background(), options)
		if err != nil {
			fmt.Printf("error converting IP to reserved IPs : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(r)
	},
}

var reservedIPCreate = &cobra.Command{
	Use:     "create ",
	Short:   "create reservedIP",
	Long:    reservedIPCreateLong,
	Example: reservedIPCreateExample,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		ipType, _ := cmd.Flags().GetString("type")
		label, _ := cmd.Flags().GetString("label")

		options := &govultr.ReservedIPReq{
			Region: region,
			IPType: ipType,
			Label:  label,
		}

		r, _, err := client.ReservedIP.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating reserved IPs : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(r)
	},
}

var reservedIPUpdate = &cobra.Command{
	Use:     "update <reservedIPID>",
	Short:   "update reservedIP",
	Long:    reservedIPUpdateLong,
	Example: reservedIPUpdateExample,
	Aliases: []string{"u"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reserved IP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")
		ip := args[0]

		options := &govultr.ReservedIPUpdateReq{
			Label: govultr.StringToStringPtr(label),
		}

		r, _, err := client.ReservedIP.Update(context.Background(), ip, options)
		if err != nil {
			fmt.Printf("error updating reserved IPs : %v\n", err)
			os.Exit(1)
		}

		printer.ReservedIP(r)
	},
}
