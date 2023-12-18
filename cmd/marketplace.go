// Copyright © 2023 The Vultr-cli Authors
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
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	marketplaceLong    = `Get commands available to marketplace`
	marketplaceExample = `
	# Full example
	vultr-cli marketplace
	`
	marketplaceAppLong    = `Get commands available to marketplace apps`
	marketplaceAppExample = `
	# Full example
	vultr-cli marketplace app
	`
	marketplaceListAppVariableLong    = `List all user-supplied variables for a given Vultr Marketplace app`
	marketplaceListAppVariableExample = `
	# Full example
	vultr-cli marketplace app list-variables exampleapp
	`
)

// Marketplace represents the marketplace command
func Marketplace() *cobra.Command { //nolint:funlen
	marketplaceCmd := &cobra.Command{
		Use:     "marketplace",
		Short:   "commands to interact with the vultr marketplace",
		Long:    marketplaceLong,
		Example: marketplaceExample,
	}

	// App flags
	marketplaceAppCmd := &cobra.Command{
		Use:     "app",
		Short:   "commands to interact with vultr marketplace apps",
		Long:    marketplaceAppLong,
		Example: marketplaceAppExample,
	}
	marketplaceAppCmd.AddCommand(marketplaceAppVariableList)
	marketplaceCmd.AddCommand(marketplaceAppCmd)

	return marketplaceCmd
}

var marketplaceAppVariableList = &cobra.Command{
	Use:     "list-variables",
	Aliases: []string{"l"},
	Short:   "list all user-supplied variables for a marketplace app",
	Long:    marketplaceListAppVariableLong,
	Example: marketplaceListAppVariableExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an imageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Marketplace.ListAppVariables(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of marketplace app variables : %v\n", err)
			os.Exit(1)
		}

		printer.MarketplaceAppVariableList(s)
	},
}
