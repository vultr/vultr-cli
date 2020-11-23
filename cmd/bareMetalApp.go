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
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// BareMetalApp represents the baremetal app commands
func BareMetalApp() *cobra.Command {
	bareMetalAppCmd := &cobra.Command{
		Use:     "app",
		Short:   "app is used to access bare metal server application commands",
		Aliases: []string{"a"},
	}

	bareMetalAppCmd.AddCommand(bareMetalAppChange, bareMetalAppChangeList)

	return bareMetalAppCmd
}

var bareMetalAppChange = &cobra.Command{
	Use:     "change <bareMetalID> <appID>",
	Short:   "Change a bare metal server's application",
	Aliases: []string{"c"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a bareMetalID and appID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := strconv.Atoi(args[1])
		options := &govultr.BareMetalUpdate{
			AppID: appID,
		}

		_, err := client.BareMetalServer.Update(context.TODO(), args[0], options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server's application changed")
	},
}

var bareMetalAppChangeList = &cobra.Command{
	Use:   "list <bareMetalID>",
	Short: "available apps a bare metal server can change to.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		list, err := client.BareMetalServer.GetUpgrades(context.TODO(), id)

		if err != nil {
			fmt.Printf("error listing available applications : %v\n", err)
			os.Exit(1)
		}

		printer.AppList(list.Applications)
	},
}
