// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

// ReservedIP represents the reservedip command
func ReservedIP() *cobra.Command {
	reservedIPCmd := &cobra.Command{
		Use:     "reserved-ip",
		Aliases: []string{"rip"},
		Short:   "",
		Long:    ``,
	}

	reservedIPCmd.AddCommand(reservedIPList, reservedIPDelete, reservedIPAttach, reservedIPDetach, reservedIPConvert, reservedIPCreate)

	// Attach
	reservedIPAttach.Flags().StringP("instance-id", "i", "", "id of instance you want to attach")
	reservedIPAttach.MarkFlagRequired("instance-id")

	// Detach
	reservedIPDetach.Flags().StringP("instance-id", "i", "", "id of instance you want to detach")
	reservedIPDetach.MarkFlagRequired("instance-id")

	// Convert
	reservedIPConvert.Flags().StringP("instance-id", "i", "", "id of instance")
	reservedIPConvert.MarkFlagRequired("instance-id")
	reservedIPConvert.Flags().StringP("ip", "", "", "ip you wish to convert")
	reservedIPConvert.MarkFlagRequired("ip")
	reservedIPConvert.Flags().StringP("label", "l", "", "label")

	// Create
	reservedIPCreate.Flags().IntP("region-id", "r", 0, "id of region")
	reservedIPCreate.MarkFlagRequired("region-id")
	reservedIPCreate.Flags().StringP("type", "t", "", "type of IP : v4 or v6")
	reservedIPCreate.MarkFlagRequired("type")
	reservedIPCreate.Flags().StringP("label", "l", "", "label")
	/*
	   convert
	   create
	*/

	return reservedIPCmd
}

var reservedIPList = &cobra.Command{
	Use:   "list",
	Short: "list all reserved IPs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rip, err := client.ReservedIP.List(context.TODO())

		if err != nil {
			fmt.Printf("error getting reserved IPs : %v", err)
			os.Exit(1)
		}

		printer.ReservedIPList(rip)
	},
}

var reservedIPDelete = &cobra.Command{
	Use:   "delete <reservedIPID>",
	Short: "delete a reserved ip",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a reservedIP ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		err := client.ReservedIP.Delete(context.TODO(), ip)

		if err != nil {
			fmt.Printf("error getting reserved IPs : %v", err)
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
		err := client.ReservedIP.Attach(context.TODO(), ip, instance)

		if err != nil {
			fmt.Printf("error attaching reserved IPs : %v", err)
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
		instance, _ := cmd.Flags().GetString("instance-id")
		err := client.ReservedIP.Detach(context.TODO(), ip, instance)

		if err != nil {
			fmt.Printf("error detaching reserved IPs : %v", err)
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
		instance, _ := cmd.Flags().GetString("instance-id")
		label, _ := cmd.Flags().GetString("label")

		r, err := client.ReservedIP.Convert(context.TODO(), ip, instance, label)

		if err != nil {
			fmt.Printf("error converting IP to reserved IPs : %v", err)
			os.Exit(1)
		}

		fmt.Println("Converted ip to reservedIP : %v", r.ReservedIPID)
	},
}

var reservedIPCreate = &cobra.Command{
	Use:   "create ",
	Short: "create reservedIP",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetInt("region-id")
		ipType, _ := cmd.Flags().GetString("type")
		label, _ := cmd.Flags().GetString("label")

		fmt.Println(region, ipType, label)

		r, err := client.ReservedIP.Create(context.TODO(), region, ipType, label)

		if err != nil {
			fmt.Printf("error creating reserved IPs : %v", err)
			os.Exit(1)
		}

		fmt.Println("Created ReservedIP : %v", r.ReservedIPID)
	},
}
