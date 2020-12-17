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

	cmd.AddCommand(userCreate, userDelete, userGet, userList, userUpdate)

	userCreate.Flags().StringP("email", "e", "", "User email")
	userCreate.Flags().StringP("name", "n", "", "User name")
	userCreate.Flags().StringP("password", "p", "", "User password")
	userCreate.Flags().StringP("api-enabled", "a", "yes", "Toggle User API Access")
	userCreate.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users,subscriptions,billing,support,provisioning,dns,abuse,upgrade,firewall,alerts]")

	userCreate.MarkFlagRequired("email")
	userCreate.MarkFlagRequired("name")
	userCreate.MarkFlagRequired("password")

	userUpdate.Flags().StringP("email", "e", "", "User email")
	userUpdate.Flags().StringP("name", "n", "", "User name")
	userUpdate.Flags().StringP("password", "p", "", "User password")
	userUpdate.Flags().StringP("api-enabled", "a", "yes", "Toggle User API Access")
	userUpdate.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users,subscriptions,billing,support,provisioning,dns,abuse,upgrade,firewall,alerts]")

	userList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	userList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

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

		options := &govultr.UserReq{
			Name:     name,
			Email:    email,
			Password: password,
			ACL:      acl,
		}

		if api == "yes" {
			options.APIEnabled = govultr.BoolToBoolPtr(true)
		}

		user, err := client.User.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating user : %v\n", err)
			os.Exit(1)
		}

		printer.User(user)
	},
}

// Delete User command
var userDelete = &cobra.Command{
	Use:     "delete <userID>",
	Short:   "Delete a user",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a userID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.User.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting user : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("User has been deleted")
	},
}

// Get User command
var userGet = &cobra.Command{
	Use:   "get <userID>",
	Short: "Get a user",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a userID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		user, err := client.User.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting user : %v\n", err)
			os.Exit(1)
		}

		printer.User(user)
	},
}

// List all Users command
var userList = &cobra.Command{
	Use:   "list",
	Short: "List all available users",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, err := client.User.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error while grabbing users %v\n", err)
			os.Exit(1)
		}

		printer.Users(list, meta)
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
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		api, _ := cmd.Flags().GetString("api-enabled")
		acl, _ := cmd.Flags().GetStringSlice("acl")

		user := &govultr.UserReq{}

		if name != "" {
			user.Name = name
		}

		if email != "" {
			user.Email = email
		}

		if password != "" {
			user.Password = password
		}

		if api == "yes" {
			user.APIEnabled = govultr.BoolToBoolPtr(true)
		} else if api == "no" {
			user.APIEnabled = govultr.BoolToBoolPtr(false)
		}

		if acl != nil {
			user.ACL = acl
		}

		if err := client.User.Update(context.Background(), id, user); err != nil {
			fmt.Printf("error updating user : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("User has been updated")
	},
}
