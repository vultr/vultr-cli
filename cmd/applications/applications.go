package applications

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

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
	"github.com/vultr/vultr-cli/cmd/utils"
)

var (
	appLong    = `Display all commands for applications`
	appExample = `
	# Full example
	vultr-cli applications
	`

	listLong    = `Display all available applications.`
	listExample = `
	# Full example
	vultr-cli applications list
		
	# Full example with paging
	vultr-cli applications list --per-page=1 --cursor="bmV4dF9fMg=="

	# Shortened with alias commands
	vultr-cli a l
	`
)

// Interface for regions
type Interface interface {
	List() ([]govultr.Application, *govultr.Meta, error)
	validate(cmd *cobra.Command, args []string)
}

// Options for regions
type Options struct {
	Args    []string
	Client  *govultr.Client
	Options *govultr.ListOptions
	Printer *printer.Output
}

func NewApplicationOptions(client *govultr.Client) *Options {
	return &Options{Client: client, Printer: &printer.Output{}}
}

// NewCmdApplications creates cobra command for applications
func NewCmdApplications(client *govultr.Client) *cobra.Command {
	o := NewApplicationOptions(client)
	cmd := &cobra.Command{
		Use:     "apps",
		Aliases: []string{"a", "application", "applications", "app"},
		Short:   "display applications",
		Long:    appLong,
		Example: appExample,
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list applications",
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			apps, meta, err := o.List()
			data := &printer.Applications{Applications: apps, Meta: meta}
			o.Printer.Display(data, err)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	cmd.AddCommand(list)
	return cmd
}

func (o *Options) validate(cmd *cobra.Command, args []string) {
	o.Printer.Output = viper.GetString("output")
	o.Options = utils.GetPaging(cmd)
	o.Args = args
}

// List all applications
func (o *Options) List() ([]govultr.Application, *govultr.Meta, error) {
	return o.Client.Application.List(context.Background(), o.Options)
}
