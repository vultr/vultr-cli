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
	"github.com/vultr/govultr"
)

// BareMetalOS represents the baremetal operating system commands
func BareMetalOS() *cobra.Command {
	bareMetalOSCmd := &cobra.Command{
		Use:     "os",
		Short:   "os is used to access bare metal server operating system commands",
		Aliases: []string{"o"},
	}

	bareMetalOSCmd.AddCommand(bareMetalOSChange)

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

		if err := client.BareMetalServer.Update(context.TODO(), args[0], options); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server's operating system changed")
	},
}
