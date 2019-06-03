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

// User represents the user command
func User() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"u"},
		Short:   "user commands",
		Long:    `user is used to access user commands`,
	}

	cmd.AddCommand(userCreate)
	cmd.AddCommand(userDelete)
	cmd.AddCommand(userList)
	cmd.AddCommand(userUpdate)

	userCreate.Flags().StringP("email", "e", "", "User email")
	userCreate.Flags().StringP("name", "n", "", "User name")
	userCreate.Flags().StringP("password", "p", "", "User password")
	userCreate.Flags().StringP("api-enabled", "a", "", "Toggle User API Access")
	userCreate.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users,subscriptions,billing,support,provisioning,dns,abuse,upgrade,firewall,alerts]")

	userCreate.MarkFlagRequired("email")
	userCreate.MarkFlagRequired("name")
	userCreate.MarkFlagRequired("password")

	userUpdate.Flags().StringP("email", "e", "", "User email")
	userUpdate.Flags().StringP("name", "n", "", "User name")
	userUpdate.Flags().StringP("password", "p", "", "User password")
	userUpdate.Flags().StringP("api-enabled", "a", "", "Toggle User API Access")
	userUpdate.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users,subscriptions,billing,support,provisioning,dns,abuse,upgrade,firewall,alerts]")
	return cmd
}

// Create user command
var userCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a user",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		api, _ := cmd.Flags().GetString("api-enabled")
		acl, _ := cmd.Flags().GetStringSlice("acl")

		id, err := client.User.Create(context.TODO(), email, name, password, api, acl)

		if err != nil {
			fmt.Println("error creating user")
			os.Exit(1)
		}

		fmt.Printf("User has been created : %s", id.UserID)
	},
}

// Delete User command
var userDelete = &cobra.Command{
	Use:   "delete <userID>",
	Short: "Delete a user",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a userID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.User.Delete(context.TODO(), id)

		if err != nil {
			fmt.Println("error deleting user")
			os.Exit(1)
		}

		fmt.Println("User has been deleted")
	},
}

// List all Users command
var userList = &cobra.Command{
	Use:   "list",
	Short: "List all available users",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := client.User.List(context.TODO())

		if err != nil {
			fmt.Errorf("error while grabbing users %v", err)
			os.Exit(1)
		}

		printer.User(list)
	},
}

// Update User command
var userUpdate = &cobra.Command{
	Use:   "update <userID>",
	Short: "Update User",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a userID")
		}
		return nil
	},
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		api, _ := cmd.Flags().GetString("api-enabled")
		acl, _ := cmd.Flags().GetStringSlice("acl")

		user := new(govultr.User)
		id := args[0]

		user.UserID = id

		if name != "" {
			user.Name = name
		}

		if email != "" {
			user.Email = email
		}

		if password != "" {
			user.Password = password
		}

		if api == "true" {
			user.APIEnabled = "yes"
		} else if api == "false" {
			user.APIEnabled = "no"
		}

		if acl != nil {
			user.ACL = acl
		}

		err := client.User.Update(context.TODO(), user)

		if err != nil {
			fmt.Println("error updating user")
			os.Exit(1)
		}

		fmt.Println("User has been updated")
	},
}
