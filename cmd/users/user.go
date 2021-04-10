// Copyright © 2021 The Vultr-cli Authors
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

package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
	"os"
)

var (
	userLong    = ``
	userExample = ``

	createLong    = ``
	createExample = ``

	getLong    = `Get a sub user from your Vultr account based on it's ID.`
	getExample = `
		# Full example
		vultr-cli users get 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
		
		# Shortened with alias commands
		vultr-cli u g 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
	`
)

type UserOptionsInterface interface {
	validate(cmd *cobra.Command, args []string)
	Get() *govultr.User
}

type UserOptions struct {
	Args   []string
	Client *govultr.Client
	User   struct {
		Email      string
		Name       string
		Password   string
		APIEnabled string
		ACLs       []string
	}
}

func NewUserOptions(client *govultr.Client) *UserOptions {
	return &UserOptions{Client: client}
}

func NewCmdUser(client *govultr.Client) *cobra.Command {
	u := NewUserOptions(client)

	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"users", "u"},
		Short:   "user commands",
		Long:    userLong,
		Example: userExample,
	}

	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a user",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			u.validate(cmd, args)
			//u.Create()
		},
	}
	create.Flags().StringP("email", "e", "", "User email")
	create.Flags().StringP("name", "n", "", "User name")
	create.Flags().StringP("password", "p", "", "User password")
	create.Flags().StringP("api-enabled", "a", "yes", "Toggle User API Access")
	create.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users,subscriptions,billing,support,provisioning,dns,abuse,upgrade,firewall,alerts")

	get := &cobra.Command{
		Use:     "get {userID}",
		Short:   "get a user",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a userID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			u.validate(cmd, args)
			printer.User(u.Get())
		},
	}

	//update
	//list
	//delete
	cmd.AddCommand(create, get)
	return cmd
}

func (u *UserOptions) validate(cmd *cobra.Command, args []string) {
	// do validation of flags into the struct
	u.Args = args
}

//func (u *UserOptions) Create() {
//
//}

// Get a single user based on ID
func (u *UserOptions) Get() *govultr.User {
	user, err := u.Client.User.Get(context.Background(), u.Args[0])
	if err != nil {
		fmt.Printf("error getting user : %v\n", err)
		os.Exit(1)
	}
	return user
}
