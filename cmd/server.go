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

	serverCmd.AddCommand(serverStart, serverStop)

	return serverCmd
}

/*
todo reboot
todo reinstall
todo tag
todo delete
todo rename
todo list
todo show grab single
todo list-ipv4
todo list ipv6

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
			fmt.Printf("error starting instance : %v", err)
		}

		fmt.Println("Started up instance")
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
			fmt.Printf("error stopping instance : %v", err)
		}

		fmt.Println("Stopped up instance")
	},
}

//backup                 get and set backup schedules
//create                 create a new virtual machine
//rename                 rename a virtual machine
//os                     show and change OS on a virtual machine
//app                    show and change application on a virtual machine
//iso                    attach/detach ISO of a virtual machine
//restore                restore from backup/snapshot
//bandwidth              list bandwidth used by a virtual machine
//create-ipv4            add a new IPv4 address to a virtual machine
//delete-ipv4            remove IPv4 address from a virtual machine
//reverse-dns            modify reverse DNS entries
//set-firewall-group     set firewall group of a virtual machine
//unset-firewall-group   remove virtual machine from firewall group
//upgrade-plan           upgrade plan of a virtual machine
