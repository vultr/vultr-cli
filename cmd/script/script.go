// Package script provides the script commands to the CLI
package script

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

// NewCmdScript provides the CLI command for startup script functions
func NewCmdScript(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "script",
		Aliases: []string{"ss", "startup-script"},
		Short:   "startup script commands",
		Long:    `script is used to access startup script commands`,
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
		Short: "list all startup scripts",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			scripts, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving startup script list : %v", err)
			}

			data := &ScriptsPrinter{Scripts: scripts, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	// Get
	get := &cobra.Command{
		Use:   "get <Script ID>",
		Short: "display the contents of specified script",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a script ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			script, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting startup script : %v", err)
			}

			data := &ScriptPrinter{Script: script}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:   "create",
		Short: "create a startup script",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for script create : %v", errNa)
			}

			script, errSc := cmd.Flags().GetString("script")
			if errSc != nil {
				return fmt.Errorf("error parsing flag 'script' for script create : %v", errSc)
			}

			sType, errST := cmd.Flags().GetString("type")
			if errST != nil {
				return fmt.Errorf("error parsing flag 'type' for script create : %v", errST)
			}

			o.ScriptReq = &govultr.StartupScriptReq{
				Name:   name,
				Script: script,
				Type:   sType,
			}

			sNew, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating startup script : %v", err)
			}

			data := &ScriptPrinter{Script: sNew}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("name", "n", "", "Name of the newly created startup script.")
	create.Flags().StringP("script", "s", "", "Startup script contents.")
	create.Flags().StringP("type", "t", "", "(Optional) Type of startup script. Possible values: 'boot', 'pxe'. Default is 'boot'.")

	if err := create.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking script create 'name' flag required: %v", err)
		os.Exit(1)
	}

	if err := create.MarkFlagRequired("script"); err != nil {
		fmt.Printf("error marking script create 'script' flag required: %v", err)
		os.Exit(1)
	}

	// Update
	update := &cobra.Command{
		Use:   "update <Script ID>",
		Short: "update startup script",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a script ID")
			}
			return nil
		},
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for script update : %v", errNa)
			}

			script, errSc := cmd.Flags().GetString("script")
			if errSc != nil {
				return fmt.Errorf("error parsing flag 'script' for script update : %v", errSc)
			}

			sType, errST := cmd.Flags().GetString("type")
			if errST != nil {
				return fmt.Errorf("error parsing flag 'type' for script update : %v", errST)
			}

			o.ScriptReq = &govultr.StartupScriptReq{
				Name:   name,
				Script: script,
				Type:   sType,
			}

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating startup script : %v", err)
			}

			o.Base.Printer.Display(printer.Info("startup script has been updated"), nil)

			return nil
		},
	}

	update.Flags().StringP("name", "n", "", "Name of the startup script.")
	update.Flags().StringP("script", "s", "", "Startup script contents.")
	update.Flags().StringP("type", "t", "", "Type of startup script. Possible values: 'boot', 'pxe'. Default is 'boot'.")

	if err := update.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking script update 'name' flag required: %v", err)
		os.Exit(1)
	}

	if err := update.MarkFlagRequired("script"); err != nil {
		fmt.Printf("error marking script update 'script' flag required: %v", err)
		os.Exit(1)
	}

	if err := update.MarkFlagRequired("type"); err != nil {
		fmt.Printf("error marking script update 'type' flag required: %v", err)
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <Script ID>",
		Short:   "delete a startup script",
		Aliases: []string{"destroy"},
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a script ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting startup script : %v", err)
			}

			o.Base.Printer.Display(printer.Info("startup script has been deleted"), nil)

			return nil
		},
	}

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
	)

	return cmd
}

type options struct {
	Base      *cli.Base
	ScriptReq *govultr.StartupScriptReq
}

func (o *options) list() ([]govultr.StartupScript, *govultr.Meta, error) {
	scripts, meta, _, err := o.Base.Client.StartupScript.List(o.Base.Context, o.Base.Options)
	return scripts, meta, err
}

func (o *options) get() (*govultr.StartupScript, error) {
	script, _, err := o.Base.Client.StartupScript.Get(o.Base.Context, o.Base.Args[0])
	return script, err
}

func (o *options) create() (*govultr.StartupScript, error) {
	script, _, err := o.Base.Client.StartupScript.Create(o.Base.Context, o.ScriptReq)
	return script, err
}

func (o *options) update() error {
	return o.Base.Client.StartupScript.Update(o.Base.Context, o.Base.Args[0], o.ScriptReq)
}

func (o *options) del() error {
	return o.Base.Client.StartupScript.Delete(o.Base.Context, o.Base.Args[0])
}
