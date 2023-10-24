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
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// SSHKey represents the ssh-key command
func SSHKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh-key",
		Aliases: []string{"ssh"},
		Short:   "ssh-key commands",
		Long:    `ssh-key is used to access SSH key commands`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Context().Value(ctxAuthKey{}).(bool) == false {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	cmd.AddCommand(sshCreate, sshDelete, sshGet, sshList, sshUpdate)

	sshCreate.Flags().StringP("name", "n", "", "Name of the SSH key")
	sshCreate.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")

	if err := sshCreate.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking ssh-key create 'name' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := sshCreate.MarkFlagRequired("key"); err != nil {
		fmt.Printf("error marking ssh-key create 'key' flag required: %v\n", err)
		os.Exit(1)
	}

	sshUpdate.Flags().StringP("name", "n", "", "Name of the SSH key")
	sshUpdate.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")

	sshList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	sshList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return cmd
}

// Create SSH key command
var sshCreate = &cobra.Command{
	Use:   "create",
	Short: "Create an SSH key",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")
		options := &govultr.SSHKeyReq{
			Name:   name,
			SSHKey: key,
		}

		id, _, err := client.SSHKey.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.SSHKey(id)
	},
}

// Delete SSH key command
var sshDelete = &cobra.Command{
	Use:     "delete <sshKeyID>",
	Short:   "Delete an SSH key",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an sshKeyID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.SSHKey.Delete(context.Background(), args[0]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("SSH key has been deleted")
	},
}

// Get SSH key command
var sshGet = &cobra.Command{
	Use:   "get <sshKeyID>",
	Short: "Get an SSH key",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an sshKeyID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ssh, _, err := client.SSHKey.Get(context.Background(), args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.SSHKey(ssh)
	},
}

// List all SSH keys command
var sshList = &cobra.Command{
	Use:   "list",
	Short: "List all SSH keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, _, err := client.SSHKey.List(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.SSHKeys(list, meta)
	},
}

// Update SSH key command
var sshUpdate = &cobra.Command{
	Use:   "update <sshKeyID>",
	Short: "Update SSH key",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an sshKeyID")
		}
		return nil
	},
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")

		s := &govultr.SSHKeyReq{
			Name: name,
		}

		if key != "" {
			s.SSHKey = key
		}

		if err := client.SSHKey.Update(context.Background(), id, s); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("SSH key has been updated")
	},
}
