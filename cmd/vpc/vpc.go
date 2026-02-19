// Package vpc provides functionality for the CLI to control VPCs
package vpc

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
	long    = `Access information about VPCs on the account and perform CRUD operations`
	example = `
	# Full example
	vultr-cli vpc
	`
	listLong    = `List all available VPC information on the account`
	listExample = `
	# Full example
	vultr-cli vpc list

	# Shortened example with aliases
	vultr-cli vpc l
	`
	getLong    = `Display information for a specific VPC`
	getExample = `
	# Full example
	vultr-cli vpc get 9fd4dcf5-7108-4641-9969-b2b9a8f77990

	# Shortened example with aliases
	vultr-cli vpc g 9fd4dcf5-7108-4641-9969-b2b9a8f77990
	`
	createLong    = `Create a new VPC with desired options`
	createExample = `
	# Full example
	vultr-cli vpc create --region="ewr" --description="Example VPC" --subnet="10.200.0.0" --size=24

	--region is required.  Everything else is optional

	# Shortened example with aliases
	vultr-cli vpc c -r="ewr" -d="Example VPC" -s="10.200.0.0" -z=24
	`
	updateLong    = `Update an existing VPC with the supplied information`
	updateExample = `
	# Full example
	vultr-cli vpc update fe8cfe1d-b25c-4c3c-8dfe-e5784bade8d9 --description="Example Updated VPC"

	# Shortned example with aliases
	vultr-cli vpc u fe8cfe1d-b25c-4c3c-8dfe-e5784bade8d9 -d="Example Updated VPC"
	`
	deleteLong    = `Delete an existing VPC`
	deleteExample = `
	#Full example
	vultr-cli vpc delete 6b8d8af9-e74a-4829-850d-647f75a056ca

	#Shortened example with aliases
	vultr-cli vpc d 6b8d8af9-e74a-4829-850d-647f75a056ca
	`
)

// NewCmdVPC provides the CLI command for VPC functions
func NewCmdVPC(base *cli.Base) *cobra.Command { //nolint:gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "vpc",
		Short:   "Commands to manage VPCs",
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
		Short:   "List all VPCs",
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			vpcs, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving vpc list : %v", err)
			}

			data := &VPCsPrinter{VPCs: vpcs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Get
	get := &cobra.Command{
		Use:     "get <VPC ID>",
		Aliases: []string{"g"},
		Short:   "Get a VPC",
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vpc, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving vpc : %v", err)
			}

			data := &VPCPrinter{VPC: vpc}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a VPC",
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for vpc create : %v", errRe)
			}

			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for vpc create : %v", errDe)
			}

			subnet, errSu := cmd.Flags().GetString("subnet")
			if errSu != nil {
				return fmt.Errorf("error parsing flag 'subnet' for vpc create : %v", errSu)
			}

			size, errSi := cmd.Flags().GetInt("size")
			if errSi != nil {
				return fmt.Errorf("error parsing flag 'size' for vpc create : %v", errSi)
			}

			o.CreateReq = &govultr.VPCReq{
				Region:       region,
				Description:  description,
				V4Subnet:     subnet,
				V4SubnetMask: size,
			}

			vpc, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating vpc : %v", err)
			}

			data := &VPCPrinter{VPC: vpc}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("region", "r", "", "The ID of the region in which to create the VPC")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking vpc create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("description", "d", "", "The description of the VPC")
	create.Flags().StringP("subnet", "s", "", "The IPv4 VPC in CIDR notation.")
	create.Flags().IntP("size", "z", 0, "The number of bits for the netmask in CIDR notation.")

	// Update
	update := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Update a VPC",
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provid a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for vpc update : %v", errDe)
			}

			o.Description = description

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating vpc : %v", err)
			}

			o.Base.Printer.Display(printer.Info("vpc has been updated"), nil)

			return nil
		},
	}

	update.Flags().StringP("description", "d", "", "The description of the VPC")
	if err := update.MarkFlagRequired("description"); err != nil {
		fmt.Printf("error marking vpc update 'description' flag required: %v", err)
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <VPC ID>",
		Aliases: []string{"destroy", "d"},
		Short:   "Delete a VPC",
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting vpc : %v", err)
			}

			o.Base.Printer.Display(printer.Info("vpc has been deleted"), nil)

			return nil
		},
	}

	// NAT Gateway
	natGateway := &cobra.Command{
		Use:   "nat-gateway",
		Short: "Commands to handle NAT Gateways",
	}

	// NAT Gateway List
	natGatewayList := &cobra.Command{
		Use:   "list <VPC ID>",
		Short: "List NAT Gateways",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ng, meta, err := o.listNATGateways()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateways : %v", err)
			}

			data := &NATGatewaysPrinter{NATGateways: ng, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Get
	natGatewayGet := &cobra.Command{
		Use:   "get <VPC ID> <NAT Gateway ID>",
		Short: "Get a NAT Gateway",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ng, err := o.getNATGateway()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateway : %v", err)
			}

			data := &NATGatewayPrinter{NATGateway: ng}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Create
	natGatewayCreate := &cobra.Command{
		Use:   "create <VPC ID>",
		Short: "Create a NAT Gateway",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for NAT Gateway create : %v", errLa)
			}

			tag, errTa := cmd.Flags().GetString("tag")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tag' for NAT Gateway create : %v", errTa)
			}

			o.NATGatewayReq = &govultr.NATGatewayReq{
				Label: label,
				Tag:   tag,
			}

			ng, err := o.createNATGateway()
			if err != nil {
				return fmt.Errorf("error creating NAT Gateway : %v", err)
			}

			data := &NATGatewayPrinter{NATGateway: ng}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	natGatewayCreate.Flags().StringP("label", "l", "", "label for the new NAT Gateway subscription")
	natGatewayCreate.Flags().StringP("tag", "t", "", "tag for the new NAT Gateway subscription")

	// NAT Gateway Update
	natGatewayUpdate := &cobra.Command{
		Use:   "update <VPC ID> <NAT Gateway ID>",
		Short: "Update a NAT Gateway",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for NAT Gateway update : %v", errLa)
			}

			tag, errTa := cmd.Flags().GetString("tag")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tag' for NAT Gateway update : %v", errTa)
			}

			o.NATGatewayReq = &govultr.NATGatewayReq{}

			if cmd.Flags().Changed("label") {
				o.NATGatewayReq.Label = label
			}

			if cmd.Flags().Changed("tag") {
				o.NATGatewayReq.Tag = tag
			}

			ng, err := o.updateNATGateway()
			if err != nil {
				return fmt.Errorf("error updating NAT Gateway : %v", err)
			}

			data := &NATGatewayPrinter{NATGateway: ng}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	natGatewayUpdate.Flags().StringP("label", "l", "", "label for the NAT Gateway subnscription")
	natGatewayUpdate.Flags().StringP("tag", "t", "", "tag for the NAT Gateway subnscription")

	// NAT Gateway Delete
	natGatewayDelete := &cobra.Command{
		Use:   "delete <VPC ID> <NAT Gateway ID>",
		Short: "Delete a NAT Gateway",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.delNATGateway(); err != nil {
				return fmt.Errorf("error deleting NAT Gateway : %v", err)
			}

			o.Base.Printer.Display(printer.Info("NAT Gateway deleted"), nil)

			return nil
		},
	}

	natGateway.AddCommand(
		natGatewayList,
		natGatewayGet,
		natGatewayCreate,
		natGatewayUpdate,
		natGatewayDelete,
	)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		natGateway,
	)

	return cmd
}

type options struct {
	Base                  *cli.Base
	CreateReq             *govultr.VPCReq
	Description           string
	NATGatewayReq         *govultr.NATGatewayReq
	FirewallRuleCreateReq *govultr.NATGatewayFirewallRuleCreateReq
	FirewallRuleUpdateReq *govultr.NATGatewayFirewallRuleUpdateReq
	PortForwardingRuleReq *govultr.NATGatewayPortForwardingRuleReq
}

func (o *options) list() ([]govultr.VPC, *govultr.Meta, error) {
	vpcs, meta, _, err := o.Base.Client.VPC.List(o.Base.Context, o.Base.Options)
	return vpcs, meta, err
}

func (o *options) get() (*govultr.VPC, error) {
	vpc, _, err := o.Base.Client.VPC.Get(o.Base.Context, o.Base.Args[0])
	return vpc, err
}

func (o *options) create() (*govultr.VPC, error) {
	vpc, _, err := o.Base.Client.VPC.Create(o.Base.Context, o.CreateReq)
	return vpc, err
}

func (o *options) update() error {
	return o.Base.Client.VPC.Update(o.Base.Context, o.Base.Args[0], o.Description)
}

func (o *options) del() error {
	return o.Base.Client.VPC.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) listNATGateways() ([]govultr.NATGateway, *govultr.Meta, error) {
	natGateways, meta, _, err := o.Base.Client.VPC.ListNATGateways(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return natGateways, meta, err
}

func (o *options) getNATGateway() (*govultr.NATGateway, error) {
	natGateway, _, err := o.Base.Client.VPC.GetNATGateway(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return natGateway, err
}

func (o *options) createNATGateway() (*govultr.NATGateway, error) {
	natGateway, _, err := o.Base.Client.VPC.CreateNATGateway(o.Base.Context, o.Base.Args[0], o.NATGatewayReq)
	return natGateway, err
}

func (o *options) updateNATGateway() (*govultr.NATGateway, error) {
	natGateway, _, err := o.Base.Client.VPC.UpdateNATGateway(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.NATGatewayReq)
	return natGateway, err
}

func (o *options) delNATGateway() error {
	return o.Base.Client.VPC.DeleteNATGateway(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}
