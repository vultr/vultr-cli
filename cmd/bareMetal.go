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
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
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
		bareMetalHalt,
		bareMetalStart,
		bareMetalGet,
		bareMetalGetVNCUrl,
		bareMetalListIPV4,
		bareMetalListIPV6,
		bareMetalList,
		BareMetalOS(),
		bareMetalReboot,
		bareMetalReinstall,
		BareMetalUserData(),
	)

	// create server
	bareMetalCreate.Flags().StringP("region", "r", "", "ID of the region where the server will be created.")
	bareMetalCreate.MarkFlagRequired("region")
	bareMetalCreate.Flags().StringP("plan", "p", "", "ID of the plan that the server will subscribe to.")
	bareMetalCreate.MarkFlagRequired("plan")
	bareMetalCreate.Flags().IntP("os", "o", 0, "ID of the operating system that will be installed on the server.")
	bareMetalCreate.Flags().StringP("script", "s", "", "(optional) ID of the startup script that will run after the server is created.")
	bareMetalCreate.Flags().StringP("snapshot", "", "", "(optional) ID of the snapshot that the server will be restored from.")
	bareMetalCreate.Flags().StringP("ipv6", "i", "", "(optional) Whether IPv6 is enabled on the server. Possible values: 'yes', 'no'. Defaults to 'no'.")
	bareMetalCreate.Flags().StringP("label", "l", "", "(optional) The label to assign to the server.")
	bareMetalCreate.Flags().StringSliceP("ssh", "k", []string{}, "(optional) Comma separated list of SSH key IDs that will be added to the server.")
	bareMetalCreate.Flags().IntP("app", "a", 0, "(optional) ID of the application that will be installed on the server.")
	bareMetalCreate.Flags().StringP("userdata", "u", "", "(optional) A generic data store, which some provisioning tools and cloud operating systems use as a configuration file.")
	bareMetalCreate.Flags().StringP("notify", "n", "", "(optional) Whether an activation email will be sent when the server is ready. Possible values: 'yes', 'no'. Defaults to 'yes'.")
	bareMetalCreate.Flags().StringP("hostname", "m", "", "(optional) The hostname to assign to the server.")
	bareMetalCreate.Flags().StringP("tag", "t", "", "(optional) The tag to assign to the server.")
	bareMetalCreate.Flags().StringP("ripv4", "v", "", "(optional) IP address of the floating IP to use as the main IP of this server.")

	bareMetalList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	bareMetalList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	bareMetalListIPV4.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	bareMetalListIPV4.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	bareMetalListIPV6.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	bareMetalListIPV6.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return bareMetalCmd
}

var bareMetalCreate = &cobra.Command{
	Use:     "create",
	Short:   "create a bare metal server",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		osID, _ := cmd.Flags().GetInt("os")
		script, _ := cmd.Flags().GetString("script")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		ipv6, _ := cmd.Flags().GetString("ipv6")
		label, _ := cmd.Flags().GetString("label")
		sshKeys, _ := cmd.Flags().GetStringSlice("ssh")
		app, _ := cmd.Flags().GetInt("app")
		userdata, _ := cmd.Flags().GetString("userdata")
		notify, _ := cmd.Flags().GetString("notify")
		hostname, _ := cmd.Flags().GetString("hostname")
		tag, _ := cmd.Flags().GetString("tag")
		ripv4, _ := cmd.Flags().GetString("ripv4")

		options := &govultr.BareMetalCreate{
			StartupScriptID: script,
			Plan:            plan,
			SnapshotID:      snapshot,
			Label:           label,
			SSHKeyIDs:       sshKeys,
			Hostname:        hostname,
			Tag:             tag,
			ReservedIPv4:    ripv4,
			OsID:            osID,
			Region:          region,
			AppID:           app,
		}

		if userdata != "" {
			options.UserData = base64.StdEncoding.EncodeToString([]byte(userdata))
		}

		if notify == "yes" {
			options.ActivationEmail = govultr.BoolToBoolPtr(true)
		}

		if ipv6 == "yes" {
			options.EnableIPv6 = govultr.BoolToBoolPtr(true)
		}

		osOptions := map[string]interface{}{"app_id": app, "snapshot_id": snapshot, "os_id": osID}
		osOption, err := optionCheckBM(osOptions)

		if err != nil {
			fmt.Printf("error creating bare metal server : %v\n", err)
			os.Exit(1)
		}

		// If no osOptions were selected and osID has a real value then set the osOptions to os_id
		if osOption == "" && osID == 0 {
			fmt.Printf("error creating bare metal server : an app, snapshot, or os ID must be provided\n")
			os.Exit(1)
		}

		bm, err := client.BareMetalServer.Create(context.TODO(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetal(bm)
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
		if err := client.BareMetalServer.Delete(context.TODO(), args[0]); err != nil {
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
		options := getPaging(cmd)
		list, meta, err := client.BareMetalServer.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalList(list, meta)
	},
}

var bareMetalGet = &cobra.Command{
	Use:   "get <bareMetalID>",
	Short: "Get a bare metal server by <bareMetalID>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := client.BareMetalServer.Get(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetal(srv)
	},
}

var bareMetalGetVNCUrl = &cobra.Command{
	Use:   "vnc <bareMetalID>",
	Short: "Get a bare metal server's VNC url by <bareMetalID>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		vnc, err := client.BareMetalServer.GetVNCUrl(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(vnc.URL)
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
		bw, err := client.BareMetalServer.GetBandwidth(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalBandwidth(bw)
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
		if err := client.BareMetalServer.Halt(context.TODO(), args[0]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server halted.")
	},
}

var bareMetalStart = &cobra.Command{
	Use:     "start <bareMetalID>",
	Short:   "Start a bare metal server.",
	Long:    `Start a bare metal server.`,
	Aliases: []string{"h"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.BareMetalServer.Start(context.TODO(), args[0]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server started.")
	},
}

var bareMetalListIPV4 = &cobra.Command{
	Use:   "ipv4 <bareMetalID>",
	Short: "List the IPv4 information of a bare metal server.",
	Long:  `List the IPv4 information of a bare metal server. IP information is only available for bare metal servers in the "active" state.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		info, meta, err := client.BareMetalServer.ListIPv4s(context.TODO(), args[0], options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalIPV4Info(info, meta)
	},
}

var bareMetalListIPV6 = &cobra.Command{
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
		options := getPaging(cmd)
		info, meta, err := client.BareMetalServer.ListIPv6s(context.TODO(), args[0], options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.BareMetalIPV6Info(info, meta)
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
		if err := client.BareMetalServer.Reboot(context.TODO(), args[0]); err != nil {
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
		if _, err := client.BareMetalServer.Reinstall(context.TODO(), args[0]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server reinstalled.")
	},
}

func optionCheckBM(options map[string]interface{}) (string, error) {
	result := []string{}

	for k, v := range options {
		switch v.(type) {
		case int:
			if v != 0 {
				result = append(result, k)
			}
		case string:
			if v != "" {
				result = append(result, k)
			}
		}
	}

	if len(result) > 1 {
		return "", fmt.Errorf("Too many options have been selected : %v : please select one", result)
	}

	// Return back an empty slice so we can possibly add in osID
	if len(result) == 0 {
		return "", nil
	}

	return result[0], nil
}
