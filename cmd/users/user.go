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
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

//todo move the lens checks into a function
//todo move error checking intos own package?
//todo add cli verbiage to vars
//todo tests
//todo break subcommands into their own vars

var (
	userLong    = ``
	userExample = ``

	createLong    = `Create a sub user on your Vultr account.`
	createExample = `
		# Full Example
		vultr-cli users create --email="vultrcli@vultr.com" --name="Vultr-cli" --password="Password123" --api-enabled="yes" --acl="manage_users,billing"

		# Shortened with alias commands
		vultr-cli users create -e="vultrcli@vultr.com" -n="Vultr-cli" -p="Password123" -a-enabled="yes" -l="manage_users,billing"
	`

	getLong    = `Get a sub user from your Vultr account based on it's ID.`
	getExample = `
		# Full example
		vultr-cli users get 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
		
		# Shortened with alias commands
		vultr-cli u g 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
	`
	listLong    = ``
	listExample = `
		
	`

	updateLong    = `Update a sub user from your Vultr account based on it's ID.`
	updateExample = `
		# Full Example
		vultr-cli users update --email="vultrcli@vultr.com" --name="Vultr-cli" --password="Password123" --api-enabled="yes" --acl="manage_users,billing"

		# Shortened with alias commands
		vultr-cli u u -e="vultrcli@vultr.com" -n="Vultr-cli" -p="Password123" -a-enabled="yes" -l="manage_users,billing"
	`

	deleteLong    = `Delete a sub user from your vultr account based on it's ID.'`
	deleteExample = `
		# Full example
		vultr-cli users delete 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
		
		# Shortened with alias commands
		vultr-cli u d 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
	`
)

// UserOptionsInterface ...
type UserOptionsInterface interface {
	validate(cmd *cobra.Command, args []string)
	Create() *govultr.User
	Get() *govultr.User
}

// UserOptions ...
type UserOptions struct {
	Base *cli.Base
	User *govultr.UserReq
}

// NewUserOptions ...
func NewUserOptions(base *cli.Base) *UserOptions {
	return &UserOptions{Base: base}
}

// NewCmdUser ...
func NewCmdUser(base *cli.Base) *cobra.Command {
	u := NewUserOptions(base)

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
			user, err := u.Create()
			if err != nil {
				printer.Error(err)
			}
			printer.User(user)
		},
	}
	create.Flags().StringP("email", "e", "", "(required) User email")
	_ = create.MarkFlagRequired("email")

	create.Flags().StringP("name", "n", "", "(required) User name")
	_ = create.MarkFlagRequired("name")

	create.Flags().StringP("password", "p", "", "(required) User password")
	_ = create.MarkFlagRequired("password")

	create.Flags().StringP("api-enabled", "a", "yes", "Toggle User API Access")
	create.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users subscriptions_view subscriptions billing support provisioning dns abuse upgrade firewall alerts objstore loadbalancer")

	// Get Command
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
			user, err := u.Create()
			if err != nil {
				printer.Error(err)
			}
			printer.User(user)
		},
	}

	//list
	list := &cobra.Command{
		Use:     "list",
		Short:   "list users",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			u.validate(cmd, args)
			u.Base.Options = utils.GetPaging(cmd)
			user, meta, err := u.List()
			if err != nil {
				printer.Error(err)
			}
			printer.Users(user, meta)
		},
	}
	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	//update
	update := &cobra.Command{
		Use:     "update {userID}",
		Short:   "update a user",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a userID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			u.validate(cmd, args)
			if err := u.Update(); err != nil {
				printer.Error(err)
			}
			fmt.Println("updated user")
		},
	}
	update.Flags().StringP("email", "e", "", "User email")
	update.Flags().StringP("name", "n", "", "User name")
	update.Flags().StringP("password", "p", "", "User password")
	update.Flags().StringP("api-enabled", "a", "yes", "Toggle User API Access")
	update.Flags().StringSliceP("acl", "l", []string{}, "User access control list in a comma separated list. Possible values manage_users subscriptions_view subscriptions billing support provisioning dns abuse upgrade firewall alerts objstore loadbalancer")

	// Delete Command
	del := &cobra.Command{
		Use:     "delete {userID}",
		Short:   "delete a user",
		Aliases: []string{"d"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a userID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			u.validate(cmd, args)
			if err := u.Delete(); err != nil {
				printer.Error(err)
				os.Exit(1)
			}
			fmt.Println("User has been deleted")
		},
	}

	cmd.AddCommand(create, get, list, update, del)
	return cmd
}

func (u *UserOptions) validate(cmd *cobra.Command, args []string) {
	u.User.Name, _ = cmd.Flags().GetString("name")
	u.User.Email, _ = cmd.Flags().GetString("email")
	u.User.Password, _ = cmd.Flags().GetString("password")
	u.User.ACL, _ = cmd.Flags().GetStringSlice("acl")

	api, _ := cmd.Flags().GetString("api-enabled")
	if api == "yes" {
		u.User.APIEnabled = govultr.BoolToBoolPtr(true)
	} else {
		u.User.APIEnabled = govultr.BoolToBoolPtr(false)
	}

	u.Base.Args = args
}

// Create ...
func (u *UserOptions) Create() (*govultr.User, error) {
	user, _, err := u.Base.Client.User.Create(context.Background(), u.User)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Get a single user based on ID
func (u *UserOptions) Get() (*govultr.User, error) {
	user, _, err := u.Base.Client.User.Get(context.Background(), u.Base.Args[0])
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserOptions) List() ([]govultr.User, *govultr.Meta, error) {
	user, meta, _, err := u.Base.Client.User.List(context.Background(), u.Base.Options)
	if err != nil {
		return nil, nil, err
	}
	return user, meta, nil
}

// Update ...
func (u *UserOptions) Update() error {
	if err := u.Base.Client.User.Update(context.Background(), u.Base.Args[0], u.User); err != nil {
		return err
	}
	return nil
}

// Delete ...
func (u *UserOptions) Delete() error {
	if err := u.Base.Client.User.Delete(context.Background(), u.Base.Args[0]); err != nil {
		return err
	}
	return nil
}
