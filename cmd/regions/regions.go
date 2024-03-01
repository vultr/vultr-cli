// Package regions provides the functionality for the CLI to access regions
package regions

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	regionLong    = `Get all available regions for Vultr.`
	regionExample = `
	# Full example
	vultr-cli regions
	`

	listLong    = `List all regions that Vultr has available.`
	listExample = `
	# Full example
	vultr-cli regions list
	
	# Full example with paging
	vultr-cli regions list --per-page=1 --cursor="bmV4dF9fQU1T"

	# Shortened with alias commands
	vultr-cli r l
	`

	availLong    = `Get all available plans in a given region.`
	availExample = `
	# Full example
	vultr-cli regions availability ewr
	
	# Full example with paging
	vultr-cli regions availability ewr 

	# Shortened with alias commands
	vultr-cli r a ewr
	`
)

// NewCmdRegion creates a cobra command for Regions
func NewCmdRegion(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "regions",
		Short:   "Display regions information",
		Aliases: []string{"r", "region"},
		Long:    regionLong,
		Example: regionExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			return nil
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List regions",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			regions, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving region list : %v", err)
			}

			data := &RegionsPrinter{Regions: regions, Meta: meta}
			o.Base.Printer.Display(data, err)

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

	availability := &cobra.Command{
		Use:     "availability <Region ID>",
		Short:   "List available plans by region",
		Aliases: []string{"a"},
		Long:    availLong,
		Example: availExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a region ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			avail, err := o.availability()
			if err != nil {
				return fmt.Errorf("error retrieving region availability : %v", err)
			}

			data := &RegionsAvailabilityPrinter{Plans: avail}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}
	availability.Flags().StringP(
		"type",
		"t",
		"",
		`type of plans for which to include availability. Possible values: 
'vc2', 'vdc, 'vhf', 'vbm'. Defaults to all Instances plans.`,
	)

	cmd.AddCommand(
		list,
		availability,
	)

	return cmd
}

type options struct {
	Base     *cli.Base
	PlanType string
}

func (o *options) list() ([]govultr.Region, *govultr.Meta, error) {
	list, meta, _, err := o.Base.Client.Region.List(context.Background(), o.Base.Options)
	return list, meta, err
}

func (o *options) availability() (*govultr.PlanAvailability, error) {
	avail, _, err := o.Base.Client.Region.Availability(context.Background(), o.Base.Args[0], o.PlanType)
	return avail, err
}
