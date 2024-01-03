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

package sshkeys

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	vultr-cli ssh create --name="ssh key name" --key "ssh-rsa AAAAB3NzaC1yc...."
	
	# Shortened with alias commands
	vultr-cli s c -n="ssh key name" -k "ssh-rsa AAAAB3NzaC1yc...."
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

// Interface for ssh-keys
type Interface interface {
	validate(cmd *cobra.Command, args []string)
	Create() (*govultr.SSHKey, error)
	Get() (*govultr.SSHKey, error)
	List() ([]govultr.SSHKey, *govultr.Meta, error)
	Update() error
	Delete() error
}

// Options for ssh-keys
type Options struct {
	Base   *cli.Base
	SSHKey *govultr.SSHKeyReq
}

// NewSSHKeyOptions returns Options struct
func NewSSHKeyOptions(base *cli.Base) *Options {
	return &Options{Base: base}
}

// NewCmdSSHKey creates a cobra command for Regions
func NewCmdSSHKey(base *cli.Base) *cobra.Command {
	o := NewSSHKeyOptions(base)

	cmd := &cobra.Command{
		Use:     "ssh-key",
		Aliases: []string{"ssh", "ssh-keys", "sshkeys"},
		Short:   "ssh-key commands",
		Long:    sshLong,
		Example: sshExample,
	}

	create := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create an SSH key",
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			key, err := o.Create()
			o.Base.Printer.Display(&SSHKeyPrinter{SSHKey: key}, err)
		},
	}
	create.Flags().StringP("name", "n", "", "Name of the SSH key")
	create.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")
	_ = create.MarkFlagRequired("name")
	_ = create.MarkFlagRequired("key")

	get := &cobra.Command{
		Use:     "get <sshKeyID>",
		Short:   "Get an SSH key",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an sshKeyID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			key, err := o.Get()
			o.Base.Printer.Display(&SSHKeyPrinter{SSHKey: key}, err)
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all SSH keys",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			list, meta, err := o.List()
			data := &SSHKeysPrinter{SSHKeys: list, Meta: meta}
			o.Base.Printer.Display(data, err)
		},
	}
	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	update := &cobra.Command{
		Use:     "update <sshKeyID>",
		Short:   "Update SSH key",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an sshKeyID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			o.Base.Printer.Display(&printer.Generic{Message: "SSH Key has been updated"}, o.Update())
		},
	}
	update.Flags().StringP("name", "n", "", "Name of the SSH key")
	update.Flags().StringP("key", "k", "", "SSH public key (in authorized_keys format)")

	deleteCmd := &cobra.Command{
		Use:     "delete <sshKeyID>",
		Short:   "Delete an SSH key",
		Aliases: []string{"destroy", "d"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an sshKeyID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			o.Base.Printer.Display(&printer.Generic{Message: "SSH Key has been deleted"}, o.Delete())
		},
	}

	cmd.AddCommand(create, get, list, update, update, deleteCmd)
	return cmd
}

func (o *Options) validate(cmd *cobra.Command, args []string) {
	fields := map[string]bool{"create": true, "update <sshKeyID>": true}
	if fields[cmd.Use] {
		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")
		o.SSHKey = &govultr.SSHKeyReq{
			Name: name,
		}

		// On updates we do not want this to be empty
		if key != "" {
			o.SSHKey.SSHKey = key
		}
	}

	if cmd.Use == "list" {
		o.Base.Options = utils.GetPaging(cmd)
	}

	o.Base.Args = args
	o.Base.Printer.Output = viper.GetString("output")
}

// Create a ssh key
func (o *Options) Create() (*govultr.SSHKey, error) {
	return o.Base.Client.SSHKey.Create(context.Background(), o.SSHKey)
}

// Get a specific ssh key on your account
func (o *Options) Get() (*govultr.SSHKey, error) {
	return o.Base.Client.SSHKey.Get(context.Background(), o.Base.Args[0])
}

// List all ssh keys on your account.
func (o *Options) List() ([]govultr.SSHKey, *govultr.Meta, error) {
	return o.Base.Client.SSHKey.List(context.Background(), o.Base.Options)
}

// Update a specific ssh key on your account
func (o *Options) Update() error {
	fmt.Println(o.SSHKey)
	return o.Base.Client.SSHKey.Update(context.Background(), o.Base.Args[0], o.SSHKey)
}

// Delete a specific ssh key on your account
func (o *Options) Delete() error {
	return o.Base.Client.SSHKey.Delete(context.Background(), o.Base.Args[0])
}
