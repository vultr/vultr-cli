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

// BareMetal represents the baremetal commands
func BareMetal() *cobra.Command {
	bareMetalCmd := &cobra.Command{
		Use:     "bare-metal",
		Short:   "bare-metal is used to access bare metal server commands",
		Aliases: []string{"bm"},
	}

	bareMetalCmd.AddCommand(
		BareMetalApp(),
		bareMetalBandwidth,
		bareMetalCreate,
		bareMetalDelete,
		bareMetalEnableIPv6,
		bareMetalHalt,
		bareMetalInfo,
		bareMetalIPV4Info,
		bareMetalIPV6Info,
		bareMetalList,
		BareMetalOS(),
		bareMetalReboot,
		bareMetalReinstall,
		bareMetalSetLabel,
		bareMetalSetTag,
		BareMetalUserData(),
	)

	// create server
	bareMetalCreate.Flags().StringP("region", "r", "", "ID of the region where the server will be created.")
	bareMetalCreate.MarkFlagRequired("region")
	bareMetalCreate.Flags().StringP("plan", "p", "", "ID of the plan that the server will subscribe to.")
	bareMetalCreate.MarkFlagRequired("plan")
	bareMetalCreate.Flags().StringP("os", "o", "", "ID of the operating system that will be installed on the server.")
	bareMetalCreate.MarkFlagRequired("os")
	bareMetalCreate.Flags().StringP("script", "s", "", "(optional) ID of the startup script that will run after the server is created.")
	bareMetalCreate.Flags().StringP("snapshot", "", "", "(optional) ID of the snapshot that the server will be restored from.")
	bareMetalCreate.Flags().StringP("ipv6", "i", "", "(optional) Whether IPv6 is enabled on the server. Possible values: 'yes', 'no'. Defaults to 'no'.")
	bareMetalCreate.Flags().StringP("label", "l", "", "(optional) The label to assign to the server.")
	bareMetalCreate.Flags().StringSliceP("ssh", "k", []string{}, "(optional) Comma separated list of SSH key IDs that will be added to the server.")
	bareMetalCreate.Flags().StringP("app", "a", "", "(optional) ID of the application that will be installed on the server.")
	bareMetalCreate.Flags().StringP("userdata", "u", "", "(optional) A generic data store, which some provisioning tools and cloud operating systems use as a configuration file.")
	bareMetalCreate.Flags().StringP("notify", "n", "", "(optional) Whether an activation email will be sent when the server is ready. Possible values: 'yes', 'no'. Defaults to 'yes'.")
	bareMetalCreate.Flags().StringP("hostname", "m", "", "(optional) The hostname to assign to the server.")
	bareMetalCreate.Flags().StringP("tag", "t", "", "(optional) The tag to assign to the server.")
	bareMetalCreate.Flags().StringP("ripv4", "v", "", "(optional) IP address of the floating IP to use as the main IP of this server.")

	return bareMetalCmd
}

var bareMetalCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a bare metal server",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		osID, _ := cmd.Flags().GetString("os")

		script, _ := cmd.Flags().GetString("script")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		ipv6, _ := cmd.Flags().GetString("ipv6")
		label, _ := cmd.Flags().GetString("label")
		sshKeys, _ := cmd.Flags().GetStringSlice("ssh")
		app, _ := cmd.Flags().GetString("app")
		userdata, _ := cmd.Flags().GetString("userdata")
		notify, _ := cmd.Flags().GetString("notify")
		hostname, _ := cmd.Flags().GetString("hostname")
		tag, _ := cmd.Flags().GetString("tag")
		ripv4, _ := cmd.Flags().GetString("ripv4")

		options := &govultr.BareMetalServerOptions{
			StartupScriptID: script,
			SnapshotID:      snapshot,
			EnableIPV6:      ipv6,
			Label:           label,
			SSHKeyIDs:       sshKeys,
			AppID:           app,
			UserData:        userdata,
			NotifyActivate:  notify,
			Hostname:        hostname,
			Tag:             tag,
			ReservedIPV4:    ripv4,
		}

		id, err := client.BareMetalServer.Create(context.TODO(), region, plan, osID, options)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("created bare metal server: %s\n", id.BareMetalServerID)
	},
}

var bareMetalDelete = &cobra.Command{
	Use:     "delete <bareMetalID>",
	Short:   "Delete a bare metal server",
	Aliases: []string{"destroy"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.Delete(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("deleted bare metal server")
	},
}

var bareMetalList = &cobra.Command{
	Use:     "list",
	Short:   "List all bare metal servers.",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		list, err := client.BareMetalServer.List(context.TODO())

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalList(list)
	},
}

var bareMetalInfo = &cobra.Command{
	Use:     "info <bareMetalID>",
	Short:   "Get a bare metal server's information",
	Aliases: []string{"i"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := client.BareMetalServer.GetServer(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetal(srv)
	},
}

var bareMetalBandwidth = &cobra.Command{
	Use:     "bandwidth <bareMetalID>",
	Short:   "Get a bare metal server's bandwidth usage",
	Aliases: []string{"b"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		bw, err := client.BareMetalServer.Bandwidth(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalBandwidth(bw)
	},
}

var bareMetalEnableIPv6 = &cobra.Command{
	Use:   "enable-ipv6 <bareMetalID>",
	Short: "Enables IPv6 networking on a bare metal server by assigning an IPv6 subnet to it.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.EnableIPV6(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("IPv6 enabled.")
	},
}

var bareMetalHalt = &cobra.Command{
	Use:   "halt <bareMetalID>",
	Short: "Halt a bare metal server.",
	Long: `Halt a bare metal server. This is a hard power off, meaning that the power to the machine is severed. 
	The data on the machine will not be modified, and you will still be billed for the machine.`,
	Aliases: []string{"h"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.Halt(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server halted.")
	},
}

var bareMetalIPV4Info = &cobra.Command{
	Use:     "ipv4 <bareMetalID>",
	Short:   "List the IPv4 information of a bare metal server.",
	Long:    `List the IPv4 information of a bare metal server. IP information is only available for bare metal servers in the "active" state.`,
	Aliases: []string{"h"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		info, err := client.BareMetalServer.IPV4Info(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalIPV4Info(info)
	},
}

var bareMetalIPV6Info = &cobra.Command{
	Use:   "ipv6 <bareMetalID>",
	Short: "List the IPv6 information of a bare metal server.",
	Long:  `List the IPv6 information of a bare metal server. IP information is only available for bare metal servers in the "active" state.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		info, err := client.BareMetalServer.IPV6Info(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalIPV6Info(info)
	},
}

var bareMetalReboot = &cobra.Command{
	Use:     "reboot <bareMetalID>",
	Short:   "Reboot a bare metal server. This is a hard reboot, which means that the server is powered off, then back on.",
	Aliases: []string{"r"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.Reboot(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server rebooted.")
	},
}

var bareMetalReinstall = &cobra.Command{
	Use:   "reinstall <bareMetalID>",
	Short: "Reinstall the operating system on a bare metal server.",
	Long: `Reinstall the operating system on a bare metal server. 
	All data will be permanently lost, but the IP address will remain the same. There is no going back from this call.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.Reinstall(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server reinstalled.")
	},
}

var bareMetalSetLabel = &cobra.Command{
	Use:   "set-label <bareMetalID> <label>",
	Short: "Set the label of a bare metal server.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a bareMetalID and label")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.SetLabel(context.TODO(), args[0], args[1])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server label set.")
	},
}

var bareMetalSetTag = &cobra.Command{
	Use:   "set-tag <bareMetalID> <tag>",
	Short: "Set the tag of a bare metal server.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a bareMetalID and tag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.BareMetalServer.SetTag(context.TODO(), args[0], args[1])

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server tag set.")
	},
}
