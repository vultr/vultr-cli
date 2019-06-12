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

// SSHKey represents the ssh-key command
func SSHKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh-key",
		Aliases: []string{"ssh"},
		Short:   "ssh-key commands",
		Long:    `ssh-key is used to access SSH key commands`,
	}

	cmd.AddCommand(sshCreate)
	cmd.AddCommand(sshDelete)
	cmd.AddCommand(sshList)
	cmd.AddCommand(sshUpdate)

	sshCreate.Flags().StringP("name", "n", "", "Name of the SSH key")
	sshCreate.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")

	sshCreate.MarkFlagRequired("name")
	sshCreate.MarkFlagRequired("key")

	sshUpdate.Flags().StringP("name", "n", "", "Name of the SSH key")
	sshUpdate.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")

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

		id, err := client.SSHKey.Create(context.TODO(), name, key)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("SSH key has been created : %s", id.SSHKeyID)
	},
}

// Delete SSH key command
var sshDelete = &cobra.Command{
	Use:   "delete <sshKeyID>",
	Short: "Delete an SSH key",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an sshKeyID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.SSHKey.Delete(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("SSH key has been deleted")
	},
}

// List all SSH keys command
var sshList = &cobra.Command{
	Use:   "list",
	Short: "List all SSH keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := client.SSHKey.List(context.TODO())

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.SSHKey(list)
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
		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")

		s := new(govultr.SSHKey)
		s.SSHKeyID = args[0]

		if name != "" {
			s.Name = name
		}

		if key != "" {
			s.Key = key
		}

		err := client.SSHKey.Update(context.TODO(), s)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("SSH key has been updated")
	},
}
