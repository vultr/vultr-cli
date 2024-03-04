// Package iso provides the ISO related commands to the CLI
package iso

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

// NewCmdISO provides the CLI command for ISO functions
func NewCmdISO(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:   "iso",
		Short: "Commands to manage ISOs",
		Long:  ``,
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
		Use:   "list",
		Short: "List all private ISOs available",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			isos, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving private ISO list : %v", err)
			}

			data := &ISOsPrinter{ISOs: isos, Meta: meta}
			o.Base.Printer.Display(data, nil)

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

	// Get
	get := &cobra.Command{
		Use:   "get <ISO ID>",
		Short: "Get a private ISO by ID",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an ISO ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			iso, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting ISO : %v", err)
			}

			data := &ISOPrinter{ISO: *iso}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:   "create",
		Short: "Create an ISO from url",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			url, errUR := cmd.Flags().GetString("url")
			if errUR != nil {
				return fmt.Errorf("error parsing flag 'url' for ISO create : %v", errUR)
			}

			o.CreateReq = &govultr.ISOReq{URL: url}

			iso, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating ISO : %v", err)
			}

			data := &ISOPrinter{ISO: *iso}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("url", "u", "", "url from where the ISO will be downloaded")
	if err := create.MarkFlagRequired("url"); err != nil {
		printer.Error(fmt.Errorf("error marking iso create 'url' flag required : %v", err))
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <ISO ID>",
		Short:   "Delete a private ISO",
		Aliases: []string{"destroy"},
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an ISO ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting ISO : %v", err)
			}

			o.Base.Printer.Display(printer.Info("ISO has been deleted"), nil)
			return nil
		},
	}

	// Public ISOs
	public := &cobra.Command{
		Use:   "public",
		Short: "List all public ISOs",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			isos, meta, err := o.listPublic()
			if err != nil {
				return fmt.Errorf("error retrieving public ISO list : %v", err)
			}

			data := &PublicISOsPrinter{ISOs: isos, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	cmd.AddCommand(list, get, create, del, public)

	return cmd
}

type options struct {
	Base      *cli.Base
	CreateReq *govultr.ISOReq
}

func (o *options) list() ([]govultr.ISO, *govultr.Meta, error) {
	isos, meta, _, err := o.Base.Client.ISO.List(o.Base.Context, o.Base.Options)
	return isos, meta, err
}

func (o *options) get() (*govultr.ISO, error) {
	iso, _, err := o.Base.Client.ISO.Get(o.Base.Context, o.Base.Args[0])
	return iso, err
}

func (o *options) create() (*govultr.ISO, error) {
	iso, _, err := o.Base.Client.ISO.Create(o.Base.Context, o.CreateReq)
	return iso, err
}

func (o *options) del() error {
	return o.Base.Client.ISO.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) listPublic() ([]govultr.PublicISO, *govultr.Meta, error) {
	isos, meta, _, err := o.Base.Client.ISO.ListPublic(o.Base.Context, o.Base.Options)
	return isos, meta, err
}
