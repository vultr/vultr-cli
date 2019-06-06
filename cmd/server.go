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

	"github.com/spf13/cobra"
)

// Server represents the server command
func Server() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "A brief description of your command",
		Long:  ``,
	}

	serverCmd.AddCommand(serverStart, serverStop, serverRestart, serverReinstall, serverTag, serverDelete, serverLabel)

	serverTag.Flags().StringP("tag", "t", "", "tag you want to set for a given instance")
	serverTag.MarkFlagRequired("tag")

	serverLabel.Flags().StringP("label", "l", "", "label you want to set for a given instance")
	serverLabel.MarkFlagRequired("label")
	return serverCmd
}

/*
todo list
todo show grab single
todo list-ipv4
todo list ipv6
todo bandwidth

*/

var serverStart = &cobra.Command{
	Use:   "start <instanceID>",
	Short: "starts a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Start(context.TODO(), id)

		if err != nil {
			fmt.Printf("error starting server : %v", err)
		}

		fmt.Println("Started up server")
	},
}

var serverStop = &cobra.Command{
	Use:   "stop <instanceID>",
	Short: "stops a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Halt(context.TODO(), id)

		if err != nil {
			fmt.Printf("error stopping server : %v", err)
		}

		fmt.Println("Stopped the server")
	},
}

var serverRestart = &cobra.Command{
	Use:   "restart <instanceID>",
	Short: "restart a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Reboot(context.TODO(), id)

		if err != nil {
			fmt.Printf("error rebooting server : %v", err)
		}

		fmt.Println("Rebooted server")
	},
}

var serverReinstall = &cobra.Command{
	Use:   "reinstall <instanceID>",
	Short: "reinstall a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Reinstall(context.TODO(), id)

		if err != nil {
			fmt.Printf("error reinstalling server : %v", err)
		}

		fmt.Println("Reinstalled server")
	},
}

var serverTag = &cobra.Command{
	Use:   "tag <instanceID>",
	Short: "add/modify tag on server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		tag, _ := cmd.Flags().GetString("tag")
		err := client.Server.SetTag(context.TODO(), id, tag)

		if err != nil {
			fmt.Printf("error adding tag to server : %v", err)
		}

		fmt.Printf("Tagged server with : %s", tag)
	},
}

var serverDelete = &cobra.Command{
	Use:   "delete <instanceID>",
	Short: "delete a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		err := client.Server.Delete(context.TODO(), id)

		if err != nil {
			fmt.Printf("error deleting server : %v", err)
		}

		fmt.Println("Deleted server")
	},
}

var serverLabel = &cobra.Command{
	Use:   "label <instanceID>",
	Short: "label a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		label, _ := cmd.Flags().GetString("label")
		err := client.Server.SetLabel(context.TODO(), id, label)

		if err != nil {
			fmt.Printf("error labeling server : %v", err)
		}

		fmt.Printf("Labeled server with : %s", label)
	},
}

//backup                 get and set backup schedules
//create                 create a new virtual machine
//os                     show and change OS on a virtual machine
//app                    show and change application on a virtual machine
//iso                    attach/detach ISO of a virtual machine
//restore                restore from backup/snapshot
//create-ipv4            add a new IPv4 address to a virtual machine
//delete-ipv4            remove IPv4 address from a virtual machine
//reverse-dns            modify reverse DNS entries
//set-firewall-group     set firewall group of a virtual machine
//unset-firewall-group   remove virtual machine from firewall group
//upgrade-plan           upgrade plan of a virtual machine
