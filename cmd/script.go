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

// Script represents the script command
func Script() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "script",
		Aliases: []string{"ss"},
		Short:   "startup script commands",
		Long:    `script is used to access startup script commands`,
	}

	cmd.AddCommand(scriptCreate)
	cmd.AddCommand(scriptDelete)
	cmd.AddCommand(scriptList)
	cmd.AddCommand(scriptUpdate)
	cmd.AddCommand(scriptContents)

	scriptCreate.Flags().StringP("name", "n", "", "Name of the newly created startup script.")
	scriptCreate.Flags().StringP("script", "s", "", "Startup script contents.")
	scriptCreate.Flags().StringP("type", "t", "", "(Optional) Type of startup script. Possible values: 'boot', 'pxe'. Default is 'boot'.")

	scriptCreate.MarkFlagRequired("name")
	scriptCreate.MarkFlagRequired("script")

	scriptUpdate.Flags().StringP("name", "n", "", "Name of the startup script.")
	scriptUpdate.Flags().StringP("script", "s", "", "Startup script contents.")

	return cmd
}

// Create startup script command
var scriptCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a startup script",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		script, _ := cmd.Flags().GetString("script")
		scriptType, _ := cmd.Flags().GetString("type")

		id, err := client.StartupScript.Create(context.TODO(), name, script, scriptType)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Startup script has been created : %s\n", id.ScriptID)
	},
}

// Delete startup script command
var scriptDelete = &cobra.Command{
	Use:     "delete <scriptID>",
	Short:   "Delete a startup script",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a scriptID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.StartupScript.Delete(context.TODO(), id)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Startup script has been deleted")
	},
}

// List all startup scripts command
var scriptList = &cobra.Command{
	Use:   "list",
	Short: "List all startup scripts",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := client.StartupScript.List(context.TODO())

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Script(list)
	},
}

// Displays the contents of a specified script
var scriptContents = &cobra.Command{
	Use:   "contents <scriptID>",
	Short: "Displays the contents of specified script",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a scriptID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		list, err := client.StartupScript.List(context.TODO())

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		matchingID := false
		scriptContent := ""
		for _, key := range list {
			if args[0] == key.ScriptID {
				matchingID = true
				scriptContent = key.Script
				break
			}
		}

		if !matchingID {
			fmt.Println("Invalid scriptID")
			os.Exit(1)
		}

		fmt.Println(scriptContent)
	},
}

// Update startup script command
var scriptUpdate = &cobra.Command{
	Use:   "update <scriptID>",
	Short: "Update startup script",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a scriptID")
		}
		return nil
	},
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		script, _ := cmd.Flags().GetString("script")

		s := new(govultr.StartupScript)
		s.ScriptID = args[0]

		if name != "" {
			s.Name = name
		}

		if script != "" {
			s.Script = script
		}

		err := client.StartupScript.Update(context.TODO(), s)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Startup script has been updated")
	},
}
