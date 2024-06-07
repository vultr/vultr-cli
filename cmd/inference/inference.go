// Package inference is used by the CLI to control serverless inference subscriptions
package inference

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

var (
	long    = `Get commands available to inference`
	example = `
	# Full example
	vultr-cli inference
	`
	listLong    = `Get all serverless inference subscriptions on your Vultr account`
	listExample = `
	# Full example
	vultr-cli inference list
	`
	createLong    = `Create a new Serverless Inference subscription with specified label`
	createExample = `
	# Full example
	vultr-cli inference create --label="example-inference"
	`
	updateLong    = `Updates a Serverless Inference subscription with the supplied information`
	updateExample = `
	# Full example
	vultr-cli inference update --label="example-inference-updated"
	`
)

// NewCmdInference provides the CLI command for inference functions
func NewCmdInference(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "inference",
		Short:   "Commands to manage serverless inference",
		Long:    long,
		Example: example,
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
		Use:     "list",
		Short:   "List inference subscriptions",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			inferenceSubs, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving inference list : %v", err)
			}

			data := &InferenceSubsPrinter{InferenceSubs: inferenceSubs}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Get
	get := &cobra.Command{
		Use:     "get <Inference ID>",
		Short:   "Retrieve an inference subscription",
		Aliases: []string{"g"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an inference ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			inferenceSub, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving inference subscription : %v", err)
			}

			data := &InferenceSubPrinter{InferenceSub: inferenceSub}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create inference subscription",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for inference create : %v", errLa)
			}

			o.CreateUpdateReq = &govultr.InferenceCreateUpdateReq{
				Label: label,
			}

			inferenceSub, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating inference subscription : %v", err)
			}

			data := &InferenceSubPrinter{InferenceSub: inferenceSub}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("label", "l", "", "label for the new inference subscription")
	if err := create.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking inference create 'label' flag required: %v", err)
		os.Exit(1)
	}

	// Update
	update := &cobra.Command{
		Use:     "update <inference ID>",
		Short:   "Update an inference subscription",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an inference ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for inference update : %v", errLa)
			}

			o.CreateUpdateReq = &govultr.InferenceCreateUpdateReq{}

			if cmd.Flags().Changed("label") {
				o.CreateUpdateReq.Label = label
			}

			inferenceSub, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating inference : %v", err)
			}

			data := &InferenceSubPrinter{InferenceSub: inferenceSub}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	update.Flags().StringP("label", "l", "", "label for the inference subscription")

	// Delete
	del := &cobra.Command{
		Use:     "delete <inference ID>",
		Short:   "Delete an inference subscription",
		Aliases: []string{"destroy", "d"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an inference ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting inference subscription : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Inference subscription has been deleted"), nil)

			return nil
		},
	}

	// Usage
	usage := &cobra.Command{
		Use:   "usage",
		Short: "Commands to display inference subscription usage information",
	}

	// Usage Get
	usageGet := &cobra.Command{
		Use:   "get <inference ID>",
		Short: "Get inference subscription usage",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide an inference ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			us, err := o.getUsage()
			if err != nil {
				return fmt.Errorf("error retrieving inference subscription usage  : %v", err)
			}

			data := &UsagePrinter{Usage: us}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	usage.AddCommand(
		usageGet,
	)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		usage,
	)

	return cmd
}

type options struct {
	Base            *cli.Base
	CreateUpdateReq *govultr.InferenceCreateUpdateReq
}

func (o *options) list() ([]govultr.Inference, error) {
	inferenceSubs, _, err := o.Base.Client.Inference.List(o.Base.Context)
	return inferenceSubs, err
}

func (o *options) get() (*govultr.Inference, error) {
	inferenceSub, _, err := o.Base.Client.Inference.Get(o.Base.Context, o.Base.Args[0])
	return inferenceSub, err
}

func (o *options) create() (*govultr.Inference, error) {
	inferenceSub, _, err := o.Base.Client.Inference.Create(o.Base.Context, o.CreateUpdateReq)
	return inferenceSub, err
}

func (o *options) update() (*govultr.Inference, error) {
	inferenceSub, _, err := o.Base.Client.Inference.Update(o.Base.Context, o.Base.Args[0], o.CreateUpdateReq)
	return inferenceSub, err
}

func (o *options) del() error {
	return o.Base.Client.Inference.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) getUsage() (*govultr.InferenceUsage, error) {
	usage, _, err := o.Base.Client.Inference.GetUsage(o.Base.Context, o.Base.Args[0])
	return usage, err
}
