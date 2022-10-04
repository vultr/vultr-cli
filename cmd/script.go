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
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// Script represents the script command
func Script() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "script",
		Aliases: []string{"ss"},
		Short:   "startup script commands",
		Long:    `script is used to access startup script commands`,
	}

	cmd.AddCommand(scriptCreate, scriptGet, scriptDelete, scriptList, scriptUpdate)

	scriptCreate.Flags().StringP("name", "n", "", "Name of the newly created startup script.")
	scriptCreate.Flags().StringP("script", "s", "", "Startup script contents.")
	scriptCreate.Flags().StringP("type", "t", "", "(Optional) Type of startup script. Possible values: 'boot', 'pxe'. Default is 'boot'.")

	if err := scriptCreate.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking script create 'name' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := scriptCreate.MarkFlagRequired("script"); err != nil {
		fmt.Printf("error marking script create 'script' flag required: %v\n", err)
		os.Exit(1)
	}

	scriptUpdate.Flags().StringP("name", "n", "", "Name of the startup script.")
	scriptUpdate.Flags().StringP("script", "s", "", "Startup script contents.")
	scriptUpdate.Flags().StringP("type", "t", "", "Type of startup script. Possible values: 'boot', 'pxe'. Default is 'boot'.")

	scriptList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	scriptList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

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

		options := &govultr.StartupScriptReq{
			Name:   name,
			Script: script,
			Type:   scriptType,
		}

		startup, err := client.StartupScript.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Script(startup)
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
		if err := client.StartupScript.Delete(context.Background(), id); err != nil {
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
		options := getPaging(cmd)
		list, meta, err := client.StartupScript.List(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.ScriptList(list, meta)
	},
}

// Displays the contents of a specified script
var scriptGet = &cobra.Command{
	Use:   "get <scriptID>",
	Short: "Displays the contents of specified script",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a scriptID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		script, err := client.StartupScript.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Script(script)
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
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		script, _ := cmd.Flags().GetString("script")
		scriptType, _ := cmd.Flags().GetString("type")

		s := &govultr.StartupScriptReq{}

		if name != "" {
			s.Name = name
		}

		if script != "" {
			s.Script = script
		}

		if scriptType != "" {
			s.Type = scriptType
		}

		if err := client.StartupScript.Update(context.Background(), id, s); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Startup script has been updated")
	},
}
