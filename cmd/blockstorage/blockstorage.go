// Package blockstorage provides the block storage functionality for
// the CLI
package blockstorage

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
	attachLong    = `Attaches a block storage resource to an specified instance`
	attachExample = `
	#Full example
	vultr-cli block-storage attach 67181686-5455-4ebb-81eb-7299f3506e2c --instance=a7898453-dd9e-4b47-bdab-9dd7a3448f1f

	#Shortened with aliased commands
	vultr-cli bs a 67181686-5455-4ebb-81eb-7299f3506e2c -i=a7898453-dd9e-4b47-bdab-9dd7a3448f1f
	`

	createLong    = `Create a new block storage resource in a specified region`
	createExample = `
	#Full example
	vultr-cli block-storage create --region='lax' --size=10 --label='your-label'

	#Full example with block-type
	vultr-cli block-storage create --region='lax' --size=10 --block-type='high_perf'

	#Shortened with aliased commands
	vultr-cli bs c -r='lax' -s=10 -l='your-label'

	#Shortened with aliased commands and block-type
	vultr-cli bs c -r='lax' -s=10 -b='high_perf'
	`

	deleteLong    = `Delete a block storage resource`
	deleteExample = `
	#Full example
	vultr-cli block-storage delete 67181686-5455-4ebb-81eb-7299f3506e2c

	#Shortened with aliased commands
	vultr-cli bs d 67181686-5455-4ebb-81eb-7299f3506e2c
	`

	detachLong    = `Detach a block storage resource from an instance`
	detachExample = `
	#Full example
	vultr-cli block-storage detach 67181686-5455-4ebb-81eb-7299f3506e2c

	#Shortened with aliased commands
	vultr-cli bs detach 67181686-5455-4ebb-81eb-7299f3506e2c
	`

	labelLong    = `Set a label for a block storage resource`
	labelExample = `
	#Full example
	vultr-cli block-storage label 67181686-5455-4ebb-81eb-7299f3506e2c --label="Example Label"

	#Shortened with aliased commands
	vultr-cli bs label 67181686-5455-4ebb-81eb-7299f3506e2c -l="Example Label"
	`

	listLong    = `Retrieves a list of active block storage resources`
	listExample = `
	#Full example
	vultr-cli block-storage list

	#Shortened with aliased commands
	vultr-cli bs l
	`

	getLong    = `Retrieves a specified block storage resource`
	getExample = `
	#Full example
	vultr-cli block-storage get 67181686-5455-4ebb-81eb-7299f3506e2c

	#Shortened with aliased commands
	vultr-cli bs g 67181686-5455-4ebb-81eb-7299f3506e2c
	`

	resizeLong    = `Resizes a specified block storage resource`
	resizeExample = `
	#Full example
	vultr-cli block-storage resize 67181686-5455-4ebb-81eb-7299f3506e2c --size=20

	#Shortened with aliased commands
	vultr-cli bs r 67181686-5455-4ebb-81eb-7299f3506e2c -s=20
	`
)

// NewCmdBlockStorage provides the command for block storage to the CLI
func NewCmdBlockStorage(base *cli.Base) *cobra.Command { //nolint:gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "block-storage",
		Aliases: []string{"bs"},
		Short:   "block storage commands",
		Long:    `block-storage is used to interact with the block-storage api`,
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
		Short:   "List block storage",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)
			bss, meta, err := o.list()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving block storage list : %v", err))
				os.Exit(1)
			}

			data := &BlockStoragesPrinter{BlockStorages: bss, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)

	// Get
	get := &cobra.Command{
		Use:     "get <Block Storage ID>",
		Short:   "Retrieve a block storage",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a block storage ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			bs, err := o.get()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving block storage : %v", err))
				os.Exit(1)
			}

			data := &BlockStoragePrinter{BlockStorage: bs}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a new block storage",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			reg, errRg := cmd.Flags().GetString("region")
			if errRg != nil {
				printer.Error(fmt.Errorf("error parsing 'region' flag for block storage create : %v", errRg))
				os.Exit(1)
			}

			size, errSz := cmd.Flags().GetInt("size")
			if errSz != nil {
				printer.Error(fmt.Errorf("error parsing 'size' flag for block storage create : %v", errSz))
				os.Exit(1)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				printer.Error(fmt.Errorf("error parsing 'label' flag for block storage create : %v", errLa))
				os.Exit(1)
			}

			blockType, errBt := cmd.Flags().GetString("block-type")
			if errBt != nil {
				printer.Error(fmt.Errorf("error parsing 'block-type' flag for block storage create : %v", errBt))
				os.Exit(1)
			}

			o.CreateReq = &govultr.BlockStorageCreate{
				Region:    reg,
				SizeGB:    size,
				Label:     label,
				BlockType: blockType,
			}

			bs, err := o.create()
			if err != nil {
				printer.Error(fmt.Errorf("error creating block storage : %v", err))
				os.Exit(1)
			}

			data := &BlockStoragePrinter{BlockStorage: bs}
			o.Base.Printer.Display(data, nil)
		},
	}

	create.Flags().StringP("region", "r", "", "ID of the region in which to create the block storage")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking block storage create 'region' flag required: %v\n", err)
		os.Exit(1)
	}

	create.Flags().IntP("size", "s", 0, "size of the block storage you want to create")
	if err := create.MarkFlagRequired("size"); err != nil {
		fmt.Printf("error marking block storage create 'size' flag required: %v\n", err)
		os.Exit(1)
	}

	create.Flags().StringP("label", "l", "", "label you want to give the block storage")

	create.Flags().StringP(
		"block-type",
		"b",
		"",
		`(optional) Block type you want to give the block storage.
		Possible values: 'high_perf', 'storage_opt'. Currently defaults to 'high_perf'.`,
	)

	// Delete
	del := &cobra.Command{
		Use:     "delete <Block Storage ID>",
		Short:   "Delete a block storage",
		Aliases: []string{"d", "destroy"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a block storage ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.del(); err != nil {
				printer.Error(fmt.Errorf("error deleting block storage : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("block storage has been deleted"), nil)
		},
	}

	// Attach
	attach := &cobra.Command{
		Use:     "attach <Block Storage ID>",
		Short:   "Attach a block storage to an instance",
		Aliases: []string{"a"},
		Long:    attachLong,
		Example: attachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a block storage ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			instance, errIe := cmd.Flags().GetString("instance")
			if errIe != nil {
				printer.Error(fmt.Errorf("error parsing 'instance' flag for block storage attach : %v", errIe))
				os.Exit(1)
			}

			live, errLe := cmd.Flags().GetBool("live")
			if errLe != nil {
				printer.Error(fmt.Errorf("error parsing 'live' flag for block storage attach : %v", errLe))
				os.Exit(1)
			}

			o.AttachReq = &govultr.BlockStorageAttach{
				InstanceID: instance,
				Live:       govultr.BoolToBoolPtr(live),
			}

			if err := o.attach(); err != nil {
				printer.Error(fmt.Errorf("error attaching block storage : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("block storage has been attached"), nil)
		},
	}

	attach.Flags().StringP("instance", "i", "", "instance ID to which to attach the block storage")
	if err := attach.MarkFlagRequired("instance"); err != nil {
		fmt.Printf("error marking block storage attach 'instance' flag required: %v\n", err)
		os.Exit(1)
	}

	attach.Flags().Bool("live", false, "attach block storage without restarting the instance")

	// Detach
	detach := &cobra.Command{
		Use:     "detach <Block Storage ID>",
		Short:   "Detach a block storage from an instance",
		Long:    detachLong,
		Example: detachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a block storage ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			live, errLe := cmd.Flags().GetBool("live")
			if errLe != nil {
				printer.Error(fmt.Errorf("error parsing 'live' flag for block storage detach : %v", errLe))
				os.Exit(1)
			}

			o.DetachReq = &govultr.BlockStorageDetach{
				Live: govultr.BoolToBoolPtr(live),
			}

			if err := o.detach(); err != nil {
				printer.Error(fmt.Errorf("error detaching block storage : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("block storage has been detached"), nil)
		},
	}

	detach.Flags().Bool("live", false, "detach block storage without a restarting instance")

	// Label
	label := &cobra.Command{
		Use:     "label <Block Storage ID>",
		Short:   "Label a block storage",
		Long:    labelLong,
		Example: labelExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a block storage ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			label, errLl := cmd.Flags().GetString("label")
			if errLl != nil {
				printer.Error(fmt.Errorf("error parsing 'label' flag for block storage : %v", errLl))
				os.Exit(1)
			}

			o.UpdateReq = &govultr.BlockStorageUpdate{
				Label: label,
			}

			if err := o.update(); err != nil {
				printer.Error(fmt.Errorf("error updating block storage label : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("block storage label has been updated"), nil)
		},
	}

	label.Flags().StringP("label", "l", "", "the label to apply to the block storage")
	if err := label.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking block storage label set 'label' flag required: %v\n", err)
		os.Exit(1)
	}

	// Resize
	resize := &cobra.Command{
		Use:     "resize <Block Storage ID>",
		Short:   "Resize a block storage",
		Aliases: []string{"r"},
		Long:    resizeLong,
		Example: resizeExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a block storage ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			size, errSz := cmd.Flags().GetInt("size")
			if errSz != nil {
				printer.Error(fmt.Errorf("error parsing 'size' flag for block storage resize : %v", errSz))
				os.Exit(1)
			}

			o.UpdateReq = &govultr.BlockStorageUpdate{
				SizeGB: size,
			}

			if err := o.update(); err != nil {
				printer.Error(fmt.Errorf("error resizing block storage : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("block storage has been resized"), nil)
		},
	}

	resize.Flags().IntP("size", "s", 0, "size you want your block storage to be")
	if err := resize.MarkFlagRequired("size"); err != nil {
		fmt.Printf("error marking block storage resize 'size' flag required: %v\n", err)
		os.Exit(1)
	}

	cmd.AddCommand(
		list,
		get,
		create,
		del,
		label,
		attach,
		detach,
		resize,
	)

	return cmd
}

type options struct {
	Base      *cli.Base
	CreateReq *govultr.BlockStorageCreate
	UpdateReq *govultr.BlockStorageUpdate
	AttachReq *govultr.BlockStorageAttach
	DetachReq *govultr.BlockStorageDetach
}

func (o *options) list() ([]govultr.BlockStorage, *govultr.Meta, error) {
	bs, meta, _, err := o.Base.Client.BlockStorage.List(o.Base.Context, o.Base.Options)
	return bs, meta, err
}

func (o *options) get() (*govultr.BlockStorage, error) {
	bs, _, err := o.Base.Client.BlockStorage.Get(o.Base.Context, o.Base.Args[0])
	return bs, err
}

func (o *options) create() (*govultr.BlockStorage, error) {
	bs, _, err := o.Base.Client.BlockStorage.Create(o.Base.Context, o.CreateReq)
	return bs, err
}

func (o *options) del() error {
	return o.Base.Client.BlockStorage.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) update() error {
	return o.Base.Client.BlockStorage.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
}

func (o *options) attach() error {
	return o.Base.Client.BlockStorage.Attach(o.Base.Context, o.Base.Args[0], o.AttachReq)
}

func (o *options) detach() error {
	return o.Base.Client.BlockStorage.Detach(o.Base.Context, o.Base.Args[0], o.DetachReq)
}
