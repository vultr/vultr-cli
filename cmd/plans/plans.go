// Package plans provides the functionality for the CLI to access plans
package plans

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
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

	metalListLong    = `Get plans for bare-metal servers`
	metalListExample = `
	# Full example
	vultr-cli plans metal
		
	# Full example with paging
	vultr-cli plans metal --per-page=1 --cursor="bmV4dF9fdmJtLTRjLTMyZ2I="

	# Shortened with alias commands
	vultr-cli p m
	`
)

// NewCmdPlan returns the cobra command for Plans
func NewCmdPlan(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "plans",
		Short:   "get information about Vultr plans",
		Aliases: []string{"p", "plan"},
		Long:    planLong,
		Example: planExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			return nil
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all instance plans",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			planType, errTy := cmd.Flags().GetString("type")
			if errTy != nil {
				return fmt.Errorf("error parsing flag 'type' for plan list: %v", errTy)
			}

			o.PlanType = planType

			plans, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error getting plans : %v", err)
			}

			data := &PlansPrinter{Plans: plans, Meta: meta}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)
	list.Flags().StringP(
		"type",
		"t",
		"",
		"(optional) The type of plans to return. Possible values: 'vc2', 'vdc', 'vhf', 'dedicated'. Defaults to all Instances plans.",
	)

	metal := &cobra.Command{
		Use:     "metal",
		Short:   "List all bare metal plans",
		Aliases: []string{"m"},
		Long:    metalListLong,
		Example: metalListExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			m, meta, err := o.metalList()
			if err != nil {
				return fmt.Errorf("error getting bare metal plans : %v", err)
			}

			data := &MetalPlansPrinter{Plans: m, Meta: meta}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	metal.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	metal.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)

	cmd.AddCommand(list, metal)
	return cmd
}

type options struct {
	Base     *cli.Base
	PlanType string
}

func (o *options) list() ([]govultr.Plan, *govultr.Meta, error) {
	plans, meta, _, err := o.Base.Client.Plan.List(context.Background(), o.PlanType, o.Base.Options)
	return plans, meta, err
}

func (o *options) metalList() ([]govultr.BareMetalPlan, *govultr.Meta, error) {
	plans, meta, _, err := o.Base.Client.Plan.ListBareMetal(context.Background(), o.Base.Options)
	return plans, meta, err
}
