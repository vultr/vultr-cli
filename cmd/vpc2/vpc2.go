// Package vpc2 provides functionality for the CLI to control VPC2s
package vpc2

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
	long    = `Get commands available to vpc2`
	example = `
	# Full example
	vultr-cli vpc2
	`
	getLong       = ``
	getExample    = ``
	createLong    = `Create a new VPC 2.0 network with specified region, description, and network settings`
	createExample = `
	# Full example
	vultr-cli vpc2 create --region="ewr" --description="example-vpc" --ip-type="v4" --ip-block="10.99.0.0" --prefix-length="24"
	`
	updateLong    = `Updates a VPC 2.0 network with the supplied information`
	updateExample = `
	# Full example
	vultr-cli vpc2 update 84fee086-6691-417a-b2db-e2a71061fa17 --description="example-vpc"
	`
	deleteLong         = ``
	deleteExample      = ``
	nodesAttachLong    = `Attaches multiple nodes to a VPC 2.0 network`
	nodesAttachExample = `
	# Full example
	vultr-cli vpc2 nodes attach 84fee086-6691-417a-b2db-e2a71061fa17 \
		--nodes="35dbcffe-58bf-46fe-bd68-964d95488dd8,1f5d784a-1011-430c-a2e2-39ba045abe3c"
	`
	nodesDetachLong    = `Detaches multiple nodes from a VPC 2.0 network`
	nodesDetachExample = `
	# Full example
	vultr-cli vpc2 nodes detach 84fee086-6691-417a-b2db-e2a71061fa17 \
		--nodes="35dbcffe-58bf-46fe-bd68-964d95488dd8,1f5d784a-1011-430c-a2e2-39ba045abe3c"
	`
)

// NewCmdVPC2 provides the CLI command for VPC2 functions
func NewCmdVPC2(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "vpc2",
		Short:   "Commands to interact with VPC2 on vultr",
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
		Aliases: []string{"l"},
		Short:   "List all available VPC2 networks",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			vpc2s, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving vpc2 list : %v", err)
			}

			data := &VPC2sPrinter{VPC2s: vpc2s, Meta: meta}
			o.Base.Printer.Display(data, nil)

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

	// Get
	get := &cobra.Command{
		Use:     "get <VPC2 ID>",
		Aliases: []string{"g"},
		Short:   "Get info on a VPC2",
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vpc2, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving vpc2 : %v", err)
			}

			data := &VPC2Printer{VPC2: vpc2}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a VPC2 network",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for vpc2 create : %v", errRe)
			}

			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for vpc2 create : %v", errDe)
			}

			ipType, errTy := cmd.Flags().GetString("ip-type")
			if errTy != nil {
				return fmt.Errorf("error parsing flag 'ip-type' for vpc2 create : %v", errTy)
			}

			ipBlock, errBl := cmd.Flags().GetString("ip-block")
			if errBl != nil {
				return fmt.Errorf("error parsing flag 'ip-block' for vpc2 create : %v", errBl)
			}

			prefixLen, errPr := cmd.Flags().GetInt("prefix-length")
			if errPr != nil {
				return fmt.Errorf("error parsing flag 'prefix-length' for vpc2 create : %v", errPr)
			}

			o.CreateReq = &govultr.VPC2Req{
				Region:       region,
				Description:  description,
				IPType:       ipType,
				IPBlock:      ipBlock,
				PrefixLength: prefixLen,
			}

			vpc2, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating vpc2 : %v", err)
			}

			data := &VPC2Printer{VPC2: vpc2}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("region", "r", "", "The ID of the region in which to create the VPC2")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking vpc create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("description", "d", "", "description for the new VPC2 network")
	create.Flags().StringP("ip-type", "", "", "IP type for the new VPC2 network")
	create.Flags().StringP("ip-block", "", "", "subnet IP address for the new VPC2 network")
	create.Flags().IntP("prefix-length", "", 0, "number of bits for the netmask in CIDR notation for the new VPC2 network")

	// Update
	update := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Update a VPC2",
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provid a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for vpc2 update : %v", errDe)
			}

			o.Description = description

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating vpc2 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("VPC2 has been updated"), nil)

			return nil
		},
	}

	update.Flags().StringP("description", "d", "", "The description of the VPC2")
	if err := update.MarkFlagRequired("description"); err != nil {
		fmt.Printf("error marking vpc2 update 'description' flag required: %v", err)
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <VPC2 ID>",
		Aliases: []string{"destroy", "d"},
		Short:   "delete a VPC2",
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting vpc2 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("vpc2 has been deleted"), nil)

			return nil
		},
	}

	// Nodes
	nodes := &cobra.Command{
		Use:   "nodes",
		Short: "Commands to handle nodes attached to a VPC2 network",
	}

	// Nodes List
	nodesList := &cobra.Command{
		Use:     "list <VPC2 ID>",
		Aliases: []string{"l"},
		Short:   "List all nodes attached to a VPC2 network",
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			nodes, meta, err := o.listNodes()
			if err != nil {
				return fmt.Errorf("error retrieving vpc2 nodes list : %v", err)
			}

			data := &VPC2NodesPrinter{Nodes: nodes, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	nodesList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	nodesList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	// Nodes Attach
	nodesAttach := &cobra.Command{
		Use:     "attach <VPC2 ID>",
		Short:   "Attach nodes to a VPC2 network",
		Aliases: []string{"a"},
		Long:    nodesAttachLong,
		Example: nodesAttachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, errNo := cmd.Flags().GetStringSlice("nodes")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'nodes' for VPC2 nodes detach : %v", errNo)
			}

			o.AttachDetachReq = &govultr.VPC2AttachDetachReq{
				Nodes: nodes,
			}

			if err := o.attachNodes(); err != nil {
				return fmt.Errorf("error attaching nodes to VPC2 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Nodes have been attached"), nil)

			return nil

		},
	}

	nodesAttach.Flags().StringSliceP("nodes", "n", []string{}, "the instance IDs you wish to attach to the VPC2 network")
	if err := nodesAttach.MarkFlagRequired("nodes"); err != nil {
		fmt.Printf("error marking vpc2 nodes attach 'nodes' flag required: %v", err)
		os.Exit(1)
	}

	// Nodes Detach
	nodesDetach := &cobra.Command{
		Use:     "detach <VPC2 ID>",
		Short:   "Detach nodes from a VPC2 network",
		Aliases: []string{"d"},
		Long:    nodesDetachLong,
		Example: nodesDetachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, errNo := cmd.Flags().GetStringSlice("nodes")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'nodes' for VPC2 nodes detach : %v", errNo)
			}

			o.AttachDetachReq = &govultr.VPC2AttachDetachReq{
				Nodes: nodes,
			}

			if err := o.detachNodes(); err != nil {
				return fmt.Errorf("error detaching nodes from VPC2 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Nodes have been detached"), nil)

			return nil

		},
	}

	nodesDetach.Flags().StringSliceP("nodes", "n", []string{}, "the instance IDs you wish to attach to the VPC2 network")
	if err := nodesDetach.MarkFlagRequired("nodes"); err != nil {
		fmt.Printf("error marking vpc2 nodes detach 'nodes' flag required: %v", err)
		os.Exit(1)
	}

	nodes.AddCommand(
		nodesList,
		nodesAttach,
		nodesDetach,
	)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		nodes,
	)

	return cmd
}

type options struct {
	Base            *cli.Base
	CreateReq       *govultr.VPC2Req
	AttachDetachReq *govultr.VPC2AttachDetachReq
	Description     string
}

func (o *options) list() ([]govultr.VPC2, *govultr.Meta, error) {
	vpc2s, meta, _, err := o.Base.Client.VPC2.List(o.Base.Context, o.Base.Options)
	return vpc2s, meta, err
}

func (o *options) get() (*govultr.VPC2, error) {
	vpc2, _, err := o.Base.Client.VPC2.Get(o.Base.Context, o.Base.Args[0])
	return vpc2, err
}

func (o *options) create() (*govultr.VPC2, error) {
	vpc2, _, err := o.Base.Client.VPC2.Create(o.Base.Context, o.CreateReq)
	return vpc2, err
}

func (o *options) update() error {
	return o.Base.Client.VPC2.Update(o.Base.Context, o.Base.Args[0], o.Description)
}

func (o *options) listNodes() ([]govultr.VPC2Node, *govultr.Meta, error) {
	nodes, meta, _, err := o.Base.Client.VPC2.ListNodes(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return nodes, meta, err
}

func (o *options) attachNodes() error {
	return o.Base.Client.VPC2.Attach(o.Base.Context, o.Base.Args[0], o.AttachDetachReq)
}

func (o *options) detachNodes() error {
	return o.Base.Client.VPC2.Detach(o.Base.Context, o.Base.Args[0], o.AttachDetachReq)
}

func (o *options) del() error {
	return o.Base.Client.VPC.Delete(o.Base.Context, o.Base.Args[0])
}
