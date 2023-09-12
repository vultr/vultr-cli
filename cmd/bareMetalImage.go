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
)

// BareMetalImage represents the baremetal image commands
func BareMetalImage() *cobra.Command {
	bareMetalImageCmd := &cobra.Command{
		Use:     "image",
		Short:   "image is used to access bare metal server image commands",
		Aliases: []string{"i"},
	}

	bareMetalImageCmd.AddCommand(bareMetalImageChange, bareMetalAppChangeList)

	return bareMetalImageCmd
}

var bareMetalImageChange = &cobra.Command{
	Use:     "change <bareMetalID> <imageID>",
	Short:   "Change a bare metal server's application",
	Aliases: []string{"c"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a bareMetalID and imageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		imageID := args[1]
		options := &govultr.BareMetalUpdate{
			ImageID: imageID,
		}

		if _, _, err := client.BareMetalServer.Update(context.TODO(), args[0], options); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("bare metal server's application changed")
	},
}
