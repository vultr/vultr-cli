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

package plans

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
	"github.com/vultr/vultr-cli/cmd/utils"
)

var (
	planLong    = `Plans will retrieve available plans for instances or bare metal.`
	planExample = `
	# Example
	vultr-cli plans
	`

	listLong    = `List all available instances plans on Vultr.`
	listExample = `
		# Full example
		vultr-cli plans list
		
		# Full example with paging
		vultr-cli plans list --type="vhf" --per-page=1 --cursor="bmV4dF9fdmhmLTFjLTFnYg=="

		# Shortened with alias commands
		vultr-cli p l
	`

	metalLong    = `List all available bare metal plans on Vultr.`
	metalExample = `
	# Full example
	vultr-cli plans metal
		
	# Full example with paging
	vultr-cli plans metal --per-page=1 --cursor="bmV4dF9fdmJtLTRjLTMyZ2I="

	# Shortened with alias commands
	vultr-cli p m
	`
)

// PlanOptionsInterface interface
type PlanOptionsInterface interface {
	validate(cmd *cobra.Command, args []string)
	List() ([]govultr.Plan, *govultr.Meta, error)
	MetalList() ([]govultr.BareMetalPlan, *govultr.Meta, error)
}

// PlanOptions struct specific for plans
type PlanOptions struct {
	Args     []string
	Client   *govultr.Client
	Options  *govultr.ListOptions
	Printer  *printer.Output
	PlanType string
}

// NewPlanOptions returns a PlanOptions struct
func NewPlanOptions(client *govultr.Client) *PlanOptions {
	return &PlanOptions{Client: client, Printer: &printer.Output{}}
}

// NewCmdPlan returns the cobra command for Plans
func NewCmdPlan(client *govultr.Client) *cobra.Command {
	p := NewPlanOptions(client)

	cmd := &cobra.Command{
		Use:     "plans",
		Short:   "get information about Vultr plans",
		Aliases: []string{"p", "plan"},
		Long:    planLong,
		Example: planExample,
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "list all instance plans",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			p.Printer.Output = viper.GetString("output")
			p.validate(cmd, args)
			plans, meta, err := p.List()
			data := &printer.Plans{Plan: plans, Meta: meta}
			p.Printer.Display(data, err)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")
	list.Flags().StringP("type", "t", "", "(optional) The type of plans to return. Possible values: 'vc2', 'vdc', 'vhf', 'dedicated'. Defaults to all Instances plans.")

	metal := &cobra.Command{
		Use:     "metal",
		Short:   "list all bare metal plans",
		Aliases: []string{"m"},
		Long:    metalLong,
		Example: metalExample,
		Run: func(cmd *cobra.Command, args []string) {
			p.validate(cmd, args)
			m, meta, err := p.MetalList()
			data := &printer.BaremetalPlans{
				Plan: m,
				Meta: meta,
			}
			p.Printer.Display(data, err)
		},
	}
	metal.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	metal.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	cmd.AddCommand(list, metal)
	return cmd
}

func (p *PlanOptions) validate(cmd *cobra.Command, args []string) {
	p.PlanType, _ = cmd.Flags().GetString("type")
	p.Options = utils.GetPaging(cmd)
	p.Args = args
}

// List retrieves all available instance plans
func (p *PlanOptions) List() ([]govultr.Plan, *govultr.Meta, error) {
	plans, meta, err := p.Client.Plan.List(context.Background(), p.PlanType, p.Options)
	if err != nil {
		return nil, nil, err
	}

	return plans, meta, nil
}

// MetalList retrieves all available bare metal plans
func (p *PlanOptions) MetalList() ([]govultr.BareMetalPlan, *govultr.Meta, error) {
	plans, meta, err := p.Client.Plan.ListBareMetal(context.Background(), p.Options)
	if err != nil {
		return nil, nil, err
	}

	return plans, meta, nil
}
