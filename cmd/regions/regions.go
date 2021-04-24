//todo examples short long
// move viper get into the validation function to clean it up
// run revive on this

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

package regions

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
	"github.com/vultr/vultr-cli/cmd/utils"
)

type Interface interface {
	Availability() (*govultr.PlanAvailability, error)
	List() ([]govultr.Region, *govultr.Meta, error)
	validate(cmd *cobra.Command, args []string)
}

type Options struct {
	Args     []string
	Client   *govultr.Client
	Options  *govultr.ListOptions
	Printer  *printer.Output
	PlanType string
}

func NewRegionOptions(client *govultr.Client) *Options {
	return &Options{Client: client, Printer: &printer.Output{}}
}

func NewCmdRegion(client *govultr.Client) *cobra.Command {
	o := NewRegionOptions(client)

	cmd := &cobra.Command{
		Use:     "regions",
		Short:   "get regions",
		Aliases: []string{"r", "region"},
		Long:    ``,
		Example: ``,
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "list regions",
		Aliases: []string{"l"},
		Long:    ``,
		Example: ``,
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			o.Options = utils.GetPaging(cmd)
			o.Printer.Output = viper.GetString("output")
			regions, meta, err := o.List()
			data := &printer.Regions{Regions: regions, Meta: meta}
			o.Printer.Display(data, err)
		},
	}
	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	availability := &cobra.Command{
		Use:     "availability <regionID>",
		Short:   "list available plans in region",
		Aliases: []string{"a"},
		Long:    ``,
		Example: ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a regionID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			o.Printer.Output = viper.GetString("output")
			avail, err := o.Availability()
			data := &printer.RegionsAvailability{AvailablePlans: avail}
			o.Printer.Display(data, err)
		},
	}
	availability.Flags().StringP("type", "t", "", "type of plans for which to include availability. Possible values: 'vc2', 'vdc, 'vhf', 'vbm'. Defaults to all Instances plans.")

	cmd.AddCommand(list, availability)
	return cmd
}

func (o *Options) validate(cmd *cobra.Command, args []string) {
	o.PlanType, _ = cmd.Flags().GetString("type")
	o.Args = args
}

func (o *Options) List() ([]govultr.Region, *govultr.Meta, error) {
	return o.Client.Region.List(context.Background(), o.Options)
}

func (o *Options) Availability() (*govultr.PlanAvailability, error) {
	return o.Client.Region.Availability(context.Background(), o.Args[0], o.PlanType)
}
