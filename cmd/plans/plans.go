// Package plans provides the functionality for the CLI to access plans
package plans

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
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

	metalLong    = `Get commands available to metal`
	metalExample = `
	#Full example
	vultr-cli plans metal

	#Shortened with aliased commands
	vultr-cli p m
	`

	metalListLong    = `List all available bare metal plans on Vultr.`
	metalListExample = `
	# Full example
	vultr-cli plans metal
		
	# Full example with paging
	vultr-cli plans metal --per-page=1 --cursor="bmV4dF9fdmJtLTRjLTMyZ2I="

	# Shortened with alias commands
	vultr-cli p m
	`
)

// PlanOptionsInterface implementes the command options for the plan command
type PlanOptionsInterface interface {
	validate(cmd *cobra.Command, args []string)
	List() ([]govultr.Plan, *govultr.Meta, error)
	MetalList() ([]govultr.BareMetalPlan, *govultr.Meta, error)
}

// PlanOptions represents the data used by the plan command
type PlanOptions struct {
	Base     *cli.Base
	PlanType string
}

// NewPlanOptions returns a PlanOptions struct
func NewPlanOptions(Base *cli.Base) *PlanOptions {
	return &PlanOptions{Base: Base}
}

// NewCmdPlan returns the cobra command for Plans
func NewCmdPlan(Base *cli.Base) *cobra.Command {
	p := NewPlanOptions(Base)

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
			p.validate(cmd, args)
			p.Base.Options = utils.GetPaging(cmd)
			plans, meta, err := p.List()
			data := &PlansPrinter{Plans: plans, Meta: meta}
			p.Base.Printer.Display(data, err)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")
	list.Flags().StringP("type", "t", "", "(optional) The type of plans to return. Possible values: 'vc2', 'vdc', 'vhf', 'dedicated'. Defaults to all Instances plans.")

	metal := &cobra.Command{
		Use:     "metal",
		Short:   "list all bare metal plans",
		Aliases: []string{"m"},
		Long:    metalListLong,
		Example: metalListExample,
		Run: func(cmd *cobra.Command, args []string) {
			p.validate(cmd, args)
			m, meta, err := p.MetalList()
			data := &MetalPlansPrinter{Plans: m, Meta: meta}
			p.Base.Printer.Display(data, err)
		},
	}
	metal.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	metal.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	cmd.AddCommand(list, metal)
	return cmd
}

func (p *PlanOptions) validate(cmd *cobra.Command, args []string) {
	p.Base.Args = args
	p.Base.Options = utils.GetPaging(cmd)
	p.PlanType, _ = cmd.Flags().GetString("type")
	p.Base.Printer.Output = viper.GetString("output")
}

// List retrieves all available instance plans
func (p *PlanOptions) List() ([]govultr.Plan, *govultr.Meta, error) {
	plans, meta, _, err := p.Base.Client.Plan.List(context.Background(), p.PlanType, p.Base.Options)
	if err != nil {
		return nil, nil, err
	}

	return plans, meta, nil
}

// MetalList retrieves all available bare metal plans
func (p *PlanOptions) MetalList() ([]govultr.BareMetalPlan, *govultr.Meta, error) {
	plans, meta, _, err := p.Base.Client.Plan.ListBareMetal(context.Background(), p.Base.Options)
	if err != nil {
		return nil, nil, err
	}

	return plans, meta, nil
}
