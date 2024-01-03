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
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// BareMetalOS represents the baremetal operating system commands
func BareMetalOS() *cobra.Command {
	bareMetalOSCmd := &cobra.Command{
		Use:     "operatingSystems",
		Short:   "operatingSystems is used to access bare metal server operating system commands",
		Aliases: []string{"o"},
	}

	bareMetalOSCmd.AddCommand(bareMetalOSChange, bareMetalOSChangeList)

	return bareMetalOSCmd
}

var bareMetalOSChange = &cobra.Command{
	Use:     "change <bareMetalID> <osID>",
	Short:   "Change a bare metal server's operating system",
	Aliases: []string{"c"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a bareMetalID and osID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		osid, _ := strconv.Atoi(args[1])
		options := &govultr.BareMetalUpdate{
			OsID: osid,
		}

		if _, _, err := client.BareMetalServer.Update(context.TODO(), args[0], options); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server's operating system changed")
	},
}

var bareMetalOSChangeList = &cobra.Command{
	Use:   "list <bareMetalID>",
	Short: "available operating systems a bare metal server can change to.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an bareMetalID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		list, _, err := client.BareMetalServer.GetUpgrades(context.TODO(), id)

		if err != nil {
			fmt.Printf("error listing available operatingSystems : %v\n", err)
			os.Exit(1)
		}

		printer.OsList(list.OS)
	},
}
