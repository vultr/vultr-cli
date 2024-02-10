// Package users provides the functionality for the CLI to access users
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

var (
	userLong    = ``
	userExample = ``

	createLong    = `Create a sub user on your Vultr account.`
	createExample = `
		# Full Example
		vultr-cli users create --email="vultrcli@vultr.com" --name="Vultr-cli" \ 
			--password="Password123" --api-enabled="true" --acl="manage_users,billing"

		# Shortened with alias commands
		vultr-cli users create -e="vultrcli@vultr.com" -n="Vultr-cli" \
			-p="Password123" --api-enabled="true" --acl="manage_users,billing"
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
		vultr-cli users update --email="vultrcli@vultr.com" --name="Vultr-cli" --password="Password123" \
			--api-enabled="false" --acl="manage_users,billing"

		# Shortened with alias commands
		vultr-cli u u -e="vultrcli@vultr.com" -n="Vultr-cli" -p="Password123" \
			--api-enabled="false" --acl="manage_users,billing"
	`

	deleteLong    = `Delete a sub user from your vultr account based on it's ID.'`
	deleteExample = `
		# Full example
		vultr-cli users delete 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
		
		# Shortened with alias commands
		vultr-cli u d 821fae4d-2a0f-4b0e-8ffd-2fe59d67d4b2
	`
)

// NewCmdUser ...
func NewCmdUser(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "user",
		Short:   "User commands",
		Aliases: []string{"users", "u"},
		Long:    userLong,
		Example: userExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Short:   "List all users",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			user, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving user list : %v", err)
			}

			data := &UsersPrinter{Users: user, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}
	list.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)

	// Get
	get := &cobra.Command{
		Use:     "get <User ID>",
		Short:   "Get a user",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a user ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			user, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving user : %v", err)
			}

			data := &UserPrinter{User: *user}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a user",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			email, errEm := cmd.Flags().GetString("email")
			if errEm != nil {
				return fmt.Errorf("error parsing flag 'email' for user create : %v", errEm)
			}

			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for user create : %v", errNa)
			}

			pass, errPa := cmd.Flags().GetString("password")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'password' for user create : %v", errPa)
			}

			api, errAp := cmd.Flags().GetBool("api-enabled")
			if errAp != nil {
				return fmt.Errorf("error parsing flag 'api-enabled' for user create : %v", errAp)
			}

			acl, errAc := cmd.Flags().GetStringSlice("acl")
			if errAc != nil {
				return fmt.Errorf("error parsing flag 'acl' for user create : %v", errAc)
			}

			o.CreateReq = &govultr.UserReq{
				Email:    email,
				Name:     name,
				Password: pass,
				ACL:      acl,
			}

			if cmd.Flags().Changed("api-enabled") {
				o.CreateReq.APIEnabled = &api
			}

			user, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating user : %v", err)
			}

			data := &UserPrinter{User: *user}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("email", "e", "", "(required) User email")
	if err := create.MarkFlagRequired("email"); err != nil {
		fmt.Printf("error marking user create 'email' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("name", "n", "", "(required) User name")
	if err := create.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking user create 'name' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("password", "p", "", "(required) User password")
	if err := create.MarkFlagRequired("password"); err != nil {
		fmt.Printf("error marking user create 'password' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().Bool("api-enabled", false, "User API access enabled")
	create.Flags().StringSliceP(
		"acl",
		"l",
		nil,
		`User access control list in a comma separated list. Possible values:
manage_users subscriptions_view subscriptions billing support provisioning dns abuse upgrade firewall alerts objstore loadbalancer`,
	)

	// Update
	update := &cobra.Command{
		Use:     "update <User ID>",
		Short:   "Update a user",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a user ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			email, errEm := cmd.Flags().GetString("email")
			if errEm != nil {
				return fmt.Errorf("error parsing flag 'email' for user create : %v", errEm)
			}

			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for user create : %v", errNa)
			}

			pass, errPa := cmd.Flags().GetString("password")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'password' for user create : %v", errPa)
			}

			api, errAp := cmd.Flags().GetBool("api-enabled")
			if errAp != nil {
				return fmt.Errorf("error parsing flag 'api-enabled' for user create : %v", errAp)
			}

			acl, errAc := cmd.Flags().GetStringSlice("acl")
			if errAc != nil {
				return fmt.Errorf("error parsing flag 'acl' for user create : %v", errAc)
			}

			o.UpdateReq = &govultr.UserReq{}

			if cmd.Flags().Changed("email") {
				o.UpdateReq.Email = email
			}

			if cmd.Flags().Changed("name") {
				o.UpdateReq.Name = name
			}

			if cmd.Flags().Changed("api-enabled") {
				o.UpdateReq.APIEnabled = &api
			}

			if cmd.Flags().Changed("acl") {
				o.UpdateReq.ACL = acl
			}

			if cmd.Flags().Changed("password") {
				o.UpdateReq.Password = pass
			}

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating user : %v", err)
			}

			o.Base.Printer.Display(printer.Info("User updated"), nil)

			return nil
		},
	}
	update.Flags().StringP("email", "e", "", "User email")
	update.Flags().StringP("name", "n", "", "User name")
	update.Flags().StringP("password", "p", "", "User password")
	update.Flags().Bool("api-enabled", false, "API access enabled")
	update.Flags().StringSliceP(
		"acl",
		"l",
		nil,
		`User access control list in a comma separated list. Possible values:
manage_users subscriptions_view subscriptions billing support provisioning dns abuse upgrade firewall alerts objstore loadbalancer`,
	)

	update.MarkFlagsOneRequired(
		"email",
		"name",
		"password",
		"api-enabled",
		"acl",
	)

	// Delete
	del := &cobra.Command{
		Use:     "delete <user ID>",
		Short:   "Delete a user",
		Aliases: []string{"d"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a user ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting user : %v", err)
			}

			o.Base.Printer.Display(printer.Info("User deleted"), nil)

			return nil
		},
	}

	cmd.AddCommand(
		create,
		get,
		list,
		update,
		del,
	)

	return cmd
}

type options struct {
	Base      *cli.Base
	CreateReq *govultr.UserReq
	UpdateReq *govultr.UserReq
}

func (o *options) list() ([]govultr.User, *govultr.Meta, error) {
	users, meta, _, err := o.Base.Client.User.List(context.Background(), o.Base.Options)
	return users, meta, err
}

func (o *options) get() (*govultr.User, error) {
	user, _, err := o.Base.Client.User.Get(context.Background(), o.Base.Args[0])
	return user, err
}

func (o *options) create() (*govultr.User, error) {
	user, _, err := o.Base.Client.User.Create(o.Base.Context, o.CreateReq)
	return user, err
}

func (o *options) update() error {
	return o.Base.Client.User.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
}

func (o *options) del() error {
	return o.Base.Client.User.Delete(o.Base.Context, o.Base.Args[0])
}
