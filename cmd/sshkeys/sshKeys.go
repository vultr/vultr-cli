// Package sshkeys provides the commands for the CLI to access ssh keys
package sshkeys

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	sshLong    = `Get all available commands for SSH Keys`
	sshExample = `
	# Full example
	vultr-cli ssh-keys
	`

	createLong    = `Create a SSH key on your Vultr account`
	createExample = `
	# Full Example
	vultr-cli ssh create --name="ssh key name" --key="ssh-rsa AAAAB3NzaC1yc...."
	
	# Shortened with alias commands
	vultr-cli s c -n="ssh key name" -k="ssh-rsa AAAAB3NzaC1yc...."
	`

	getLong    = `Get a single SSH Key from your account`
	getExample = `
	# Full example
	vultr-cli ssh-key get ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli ssh g ffd31f18-5f77-454c-9065-212f942c3c35
	`

	listLong    = `Get all ssh keys available on your Vultr account`
	listExample = `
	# Full example
	vultr-cli ssh-key list 
	
	# Full example with paging
	vultr-cli ssh-key list --per-page=1 --cursor="bmV4dF9fQU1T"

	# Shortened with alias commands
	vultr-cli ssh l 	
	`

	updateLong    = `Update a specific SSH Key on your Vultr Account`
	updateExample = `
	# Full example
	vultr-cli ssh-key update ffd31f18-5f77-454c-9065-212f942c3c35 --name="updated name" --key="ssh-rsa AAAAB3NzaC1yc...."

	# Shortened with alias commands
	vultr-cli ssh u ffd31f18-5f77-454c-9065-212f942c3c35 --name="updated name" --key="ssh-rsa AAAAB3NzaC1yc...."
	`

	deleteLong    = `Delete a specific SSH Key off your Vultr Account`
	deleteExample = `
	# Full example
	vultr-cli ssh-key delete ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli ssh d ffd31f18-5f77-454c-9065-212f942c3c35
	`
)

// NewCmdSSHKey creates a cobra command for Regions
func NewCmdSSHKey(base *cli.Base) *cobra.Command { //nolint:gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "ssh-key",
		Short:   "Commands to access SSH key functions",
		Aliases: []string{"ssh", "ssh-keys", "sshkeys"},
		Long:    sshLong,
		Example: sshExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Short:   "List all SSH keys",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			list, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving ssh key list : %v", err)
			}

			data := &SSHKeysPrinter{SSHKeys: list, Meta: meta}
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
		Use:     "get <SSH Key ID>",
		Short:   "Get an SSH key",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an SSH Key ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting ssh key : %v", err)
			}

			o.Base.Printer.Display(&SSHKeyPrinter{SSHKey: key}, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create an SSH key",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for ssh key create : %v", errNa)
			}

			key, errKe := cmd.Flags().GetString("key")
			if errKe != nil {
				return fmt.Errorf("error parsing flag 'key' for ssh key create : %v", errKe)
			}

			o.SSHKeyReq = &govultr.SSHKeyReq{
				Name:   name,
				SSHKey: key,
			}

			k, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating ssh key : %v", err)
			}

			o.Base.Printer.Display(&SSHKeyPrinter{SSHKey: k}, err)

			return nil
		},
	}
	create.Flags().StringP("name", "n", "", "Name of the SSH key")
	create.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")
	create.MarkFlagsRequiredTogether("name", "key")

	// Update
	update := &cobra.Command{
		Use:     "update <SSH Key ID>",
		Short:   "Update SSH key",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an SSH Key ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for ssh key update : %v", errNa)
			}

			key, errKe := cmd.Flags().GetString("key")
			if errKe != nil {
				return fmt.Errorf("error parsing flag 'key' for ssh key update : %v", errKe)
			}

			o.SSHKeyReq = &govultr.SSHKeyReq{}

			if cmd.Flags().Changed("name") {
				o.SSHKeyReq.Name = name
			}

			if cmd.Flags().Changed("key") {
				o.SSHKeyReq.SSHKey = key
			}

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating ssh key : %v", err)
			}

			o.Base.Printer.Display(printer.Info("SSH Key has been updated"), nil)

			return nil
		},
	}
	update.Flags().StringP("name", "n", "", "Name of the SSH key")
	update.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")

	// Delete
	del := &cobra.Command{
		Use:     "delete <sshKeyID>",
		Short:   "Delete an SSH key",
		Aliases: []string{"destroy", "d"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an SSH Key ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting ssh key : %v", err)
			}

			o.Base.Printer.Display(printer.Info("SSH Key has been deleted"), nil)

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
	SSHKeyReq *govultr.SSHKeyReq
}

func (o *options) create() (*govultr.SSHKey, error) {
	key, _, err := o.Base.Client.SSHKey.Create(context.Background(), o.SSHKeyReq)
	return key, err
}

func (o *options) get() (*govultr.SSHKey, error) {
	key, _, err := o.Base.Client.SSHKey.Get(context.Background(), o.Base.Args[0])
	return key, err
}

func (o *options) list() ([]govultr.SSHKey, *govultr.Meta, error) {
	keys, meta, _, err := o.Base.Client.SSHKey.List(context.Background(), o.Base.Options)
	return keys, meta, err
}

func (o *options) update() error {
	return o.Base.Client.SSHKey.Update(context.Background(), o.Base.Args[0], o.SSHKeyReq)
}

func (o *options) del() error {
	return o.Base.Client.SSHKey.Delete(context.Background(), o.Base.Args[0])
}
