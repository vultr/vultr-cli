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

	natGatewayLong    = `Access information about NAT Gateways on the account's VPC network and perform CRUD operations`
	natGatewayExample = `
	# Full example
	vultr-cli vpc nat-gateway

	# Shortened example with alias
	vultr-cli vpc ng
	`
	natGatewayListLong    = `List all available NAT Gateway information on the account's VPC network`
	natGatewayListExample = `
	# Full example
	vultr-cli vpc nat-gateway list e3512e83-64e9-4d3e-a401-9b86f2e09b1d

	# Shortened example with aliases
	vultr-cli vpc ng l e3512e83-64e9-4d3e-a401-9b86f2e09b1d
	`
	natGatewayGetLong    = `Display information for a specific NAT Gateway`
	natGatewayGetExample = `
	# Full example
	vultr-cli vpc nat-gateway get e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0

	# Shortened example with aliases
	vultr-cli vpc ng g e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0
	`
	natGatewayCreateLong    = `Create a new NAT Gateway with desired options`
	natGatewayCreateExample = `
	# Full example
	vultr-cli vpc nat-gateway create e3512e83-64e9-4d3e-a401-9b86f2e09b1d --label="example-label" --tag="example tag"

	--label and --tag are optional

	# Shortened example with aliases
	vultr-cli vpc ng c e3512e83-64e9-4d3e-a401-9b86f2e09b1d -l="example-label" -t="example tag"
	`
	natGatewayUpdateLong    = `Update an existing NAT Gateway with the supplied information`
	natGatewayUpdateExample = `
	# Full example
	vultr-cli vpc nat-gateway update e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 --tag="example updated tag"

	# Shortned example with aliases
	vultr-cli vpc ng u e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 -t="example updated tag"
	`
	natGatewayDeleteLong    = `Delete an existing NAT Gateway`
	natGatewayDeleteExample = `
	#Full example
	vultr-cli vpc nat-gateway delete e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0

	#Shortened example with aliases
	vultr-cli vpc ng d e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0
	`

	pfrLong    = `Access information about Port Forwarding Rules on the account's NAT Gateway and perform CRUD operations`
	pfrExample = `
	# Full example
	vultr-cli vpc nat-gateway port-forwarding-rule e3512e83-64e9-4d3e-a401-9b86f2e09b1d

	# Shortened example with alias
	vultr-cli vpc ng pfr e3512e83-64e9-4d3e-a401-9b86f2e09b1d
	`
	pfrListLong    = `List all available Port Forwarding Rule information on the account's NAT Gateway`
	pfrListExample = `
	# Full example
	vultr-cli vpc nat-gateway port-forwarding-rule list e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0

	# Shortened example with aliases
	vultr-cli vpc ng pfr l e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0
	`
	pfrGetLong    = `Display information for a specific Port Forwarding Rule`
	pfrGetExample = `
	# Full example
	vultr-cli vpc nat-gateway port-forwarding-rule get e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 bf66d7cf-7587-4acd-a053-a2ea34e33f1d

	# Shortened example with aliases
	vultr-cli vpc ng pfr g e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 bf66d7cf-7587-4acd-a053-a2ea34e33f1d
	`
	pfrCreateLong    = `Create a new Port Forwarding Rule with desired options`
	pfrCreateExample = `
	# Full example
	vultr-cli vpc nat-gateway port-forwarding-rule create e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 --name="example-rule" --description="example desc" --internal-ip="10.1.2.3" --protocol="tcp" --external-port="123" --internal-port="555" --enabled="true"

	--description is optional. Everything else is required

	# Shortened example with aliases
	vultr-cli vpc ng pfr c e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 -n="example-rule" -d="example desc" --internal-ip="10.1.2.3" -p="tcp" --external-port="123" --internal-port="555" -e="true"
	`
	pfrUpdateLong    = `Update an existing Port Forwarding Rule with the supplied information`
	pfrUpdateExample = `
	# Full example
	vultr-cli vpc nat-gateway port-forwarding-rule update e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 bf66d7cf-7587-4acd-a053-a2ea34e33f1d --name="example-rule-updated"

	# Shortned example with aliases
	vultr-cli vpc ng pfr u e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 bf66d7cf-7587-4acd-a053-a2ea34e33f1d -n="example-rule-updated"
	`
	pfrDeleteLong    = `Delete an existing Port Forwarding Rule`
	pfrDeleteExample = `
	#Full example
	vultr-cli vpc nat-gateway port-forwarding-rule delete e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 bf66d7cf-7587-4acd-a053-a2ea34e33f1d

	#Shortened example with aliases
	vultr-cli vpc ng pfr d e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 bf66d7cf-7587-4acd-a053-a2ea34e33f1d
	`

	fwrLong    = `Access information about Firewall Rules on the account's NAT Gateway and perform CRUD operations`
	fwrExample = `
	# Full example
	vultr-cli vpc nat-gateway port-forwarding-rule e3512e83-64e9-4d3e-a401-9b86f2e09b1d

	# Shortened example with alias
	vultr-cli vpc ng fr e3512e83-64e9-4d3e-a401-9b86f2e09b1d
	`
	fwrListLong    = `List all available Firewall Rule information on the account's NAT Gateway`
	fwrListExample = `
	# Full example
	vultr-cli vpc nat-gateway firewall-rule list e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0

	# Shortened example with aliases
	vultr-cli vpc ng fr l e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0
	`
	fwrGetLong    = `Display information for a specific Firewall Rule`
	fwrGetExample = `
	# Full example
	vultr-cli vpc nat-gateway firewall-rule get e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 74fb2217-de13-4ca2-8065-0c16281a7849

	# Shortened example with aliases
	vultr-cli vpc ng fr g e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 74fb2217-de13-4ca2-8065-0c16281a7849
	`
	fwrCreateLong    = `Create a new Firewall Rule with desired options`
	fwrCreateExample = `
	# Full example
	vultr-cli vpc nat-gateway firewall-rule create e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 --protocol="tcp" --port="123" --subnet="1.2.3.4" --subnet-size="24" --notes="example rule"

	--notes is optional. Everything else is required

	# Shortened example with aliases
	vultr-cli vpc ng fr c e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 --protocol="tcp" -p="123" -s="1.2.3.4" --subnet-size="24" -n="example rule"
	`
	fwrUpdateLong    = `Update an existing Firewall Rule with the supplied information`
	fwrUpdateExample = `
	# Full example
	vultr-cli vpc nat-gateway firewall-rule update e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 74fb2217-de13-4ca2-8065-0c16281a7849 --name="example-rule-updated"

	# Shortned example with aliases
	vultr-cli vpc ng fr u e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 74fb2217-de13-4ca2-8065-0c16281a7849 -n="example-rule-updated"
	`
	fwrDeleteLong    = `Delete an existing Firewall Rule`
	fwrDeleteExample = `
	#Full example
	vultr-cli vpc nat-gateway firewall-rule delete e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 74fb2217-de13-4ca2-8065-0c16281a7849

	#Shortened example with aliases
	vultr-cli vpc ng fr d e3512e83-64e9-4d3e-a401-9b86f2e09b1d efc811e7-d07d-45b1-b50e-19a0c2936ce0 74fb2217-de13-4ca2-8065-0c16281a7849
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
		Use:     "nat-gateway",
		Aliases: []string{"ng"},
		Short:   "Commands to handle NAT Gateways",
		Long:    natGatewayLong,
		Example: natGatewayExample,
	}

	// NAT Gateway List
	natGatewayList := &cobra.Command{
		Use:     "list <VPC ID>",
		Aliases: []string{"l"},
		Short:   "List NAT Gateways",
		Long:    natGatewayListLong,
		Example: natGatewayListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ngs, meta, err := o.listNATGateways()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateways : %v", err)
			}

			data := &NATGatewaysPrinter{NATGateways: ngs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Get
	natGatewayGet := &cobra.Command{
		Use:     "get <VPC ID> <NAT Gateway ID>",
		Aliases: []string{"g"},
		Short:   "Get a NAT Gateway",
		Long:    natGatewayGetLong,
		Example: natGatewayGetExample,
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
		Use:     "create <VPC ID>",
		Aliases: []string{"c"},
		Short:   "Create a NAT Gateway",
		Long:    natGatewayCreateLong,
		Example: natGatewayCreateExample,
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
		Use:     "update <VPC ID> <NAT Gateway ID>",
		Aliases: []string{"u"},
		Short:   "Update a NAT Gateway",
		Long:    natGatewayUpdateLong,
		Example: natGatewayUpdateExample,
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
		Use:     "delete <VPC ID> <NAT Gateway ID>",
		Aliases: []string{"d"},
		Short:   "Delete a NAT Gateway",
		Long:    natGatewayDeleteLong,
		Example: natGatewayDeleteExample,
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

	// NAT Gateway Port Forwarding Rule
	portForwardingRule := &cobra.Command{
		Use:     "port-forwarding-rule",
		Aliases: []string{"pfr"},
		Short:   "Commands to handle NAT Gateway port forwarding rules",
		Long:    pfrLong,
		Example: pfrExample,
	}

	// NAT Gateway Port Forwarding Rule List
	portForwardingRuleList := &cobra.Command{
		Use:     "list <VPC ID>",
		Aliases: []string{"l"},
		Short:   "List NAT Gateway port forwarding rules",
		Long:    pfrListLong,
		Example: pfrListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pfrs, meta, err := o.listNATGatewayPortForwardingRules()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateway port forwarding rules : %v", err)
			}

			data := &NATGatewayPortForwardingRulesPrinter{PortForwardingRules: pfrs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Port Forwarding Rule Get
	portForwardingRuleGet := &cobra.Command{
		Use:     "get <VPC ID> <NAT Gateway ID> <Port Forwarding Rule ID>",
		Aliases: []string{"g"},
		Short:   "Get a NAT Gateway port forwarding rule",
		Long:    pfrGetLong,
		Example: pfrGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("please provide a VPC ID, a NAT Gateway ID, and a Port Forwarding Rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pfr, err := o.getNATGatewayPortForwardingRule()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateway port forwarding rule : %v", err)
			}

			data := &NATGatewayPortForwardingRulePrinter{PortForwardingRule: pfr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Port Forwarding Rule Create
	portForwardingRuleCreate := &cobra.Command{
		Use:     "create <VPC ID> <NAT Gateway ID>",
		Aliases: []string{"c"},
		Short:   "Create a NAT Gateway port forwarding rule",
		Long:    pfrCreateLong,
		Example: pfrCreateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for NAT Gateway port forwarding rule create : %v", errNa)
			}

			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for NAT Gateway port forwarding rule create : %v", errDe)
			}

			internalIP, errIIP := cmd.Flags().GetString("internal-ip")
			if errIIP != nil {
				return fmt.Errorf("error parsing flag 'internal-ip' for NAT Gateway port forwarding rule create : %v", errIIP)
			}

			protocol, errPr := cmd.Flags().GetString("protocol")
			if errPr != nil {
				return fmt.Errorf("error parsing flag 'protocol' for NAT Gateway port forwarding rule create : %v", errPr)
			}

			externalPort, errEP := cmd.Flags().GetInt("external-port")
			if errEP != nil {
				return fmt.Errorf("error parsing flag 'external-port' for NAT Gateway port forwarding rule create : %v", errEP)
			}

			internalPort, errIP := cmd.Flags().GetInt("internal-port")
			if errIP != nil {
				return fmt.Errorf("error parsing flag 'internal-port' for NAT Gateway port forwarding rule create : %v", errIP)
			}

			enabled, errEn := cmd.Flags().GetBool("enabled")
			if errEn != nil {
				return fmt.Errorf("error parsing flag 'enabled' for NAT Gateway port forwarding rule create : %v", errEn)
			}

			o.PortForwardingRuleReq = &govultr.NATGatewayPortForwardingRuleReq{
				Name:         name,
				Description:  description,
				InternalIP:   internalIP,
				Protocol:     protocol,
				ExternalPort: externalPort,
				InternalPort: internalPort,
				Enabled:      &enabled,
			}

			pfr, err := o.createNATGatewayPortForwardingRule()
			if err != nil {
				return fmt.Errorf("error creating NAT Gateway port forwarding rule : %v", err)
			}

			data := &NATGatewayPortForwardingRulePrinter{PortForwardingRule: pfr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	portForwardingRuleCreate.Flags().StringP("name", "n", "", "name for the new NAT Gateway port forwarding rule")
	portForwardingRuleCreate.Flags().StringP("description", "d", "", "description for the new NAT Gateway port forwarding rule")
	portForwardingRuleCreate.Flags().String("internal-ip", "", "internal IP for the new NAT Gateway port forwarding rule")
	portForwardingRuleCreate.Flags().StringP("protocol", "p", "", "protocol for the new NAT Gateway port forwarding rule")
	portForwardingRuleCreate.Flags().Int("external-port", 0, "external port for the new NAT Gateway port forwarding rule")
	portForwardingRuleCreate.Flags().Int("internal-port", 0, "internal port for the new NAT Gateway port forwarding rule")
	portForwardingRuleCreate.Flags().BoolP("enabled", "e", true, "name for the new NAT Gateway port forwarding rule")

	// NAT Gateway Port Forwarding Rule Update
	portForwardingRuleUpdate := &cobra.Command{
		Use:     "update <VPC ID> <NAT Gateway ID> <Port Forwarding Rule ID>",
		Aliases: []string{"u"},
		Short:   "Update a NAT Gateway port forwarding rule",
		Long:    pfrUpdateLong,
		Example: pfrUpdateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("please provide a VPC ID, a NAT Gateway ID, and a Port Forwarding Rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for NAT Gateway port forwarding rule update : %v", errNa)
			}

			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for NAT Gateway port forwarding rule update : %v", errDe)
			}

			internalIP, errIIP := cmd.Flags().GetString("internal-ip")
			if errIIP != nil {
				return fmt.Errorf("error parsing flag 'internal-ip' for NAT Gateway port forwarding rule update : %v", errIIP)
			}

			protocol, errPr := cmd.Flags().GetString("protocol")
			if errPr != nil {
				return fmt.Errorf("error parsing flag 'protocol' for NAT Gateway port forwarding rule update : %v", errPr)
			}

			externalPort, errEP := cmd.Flags().GetInt("external-port")
			if errEP != nil {
				return fmt.Errorf("error parsing flag 'external-port' for NAT Gateway port forwarding rule update : %v", errEP)
			}

			internalPort, errIP := cmd.Flags().GetInt("internal-port")
			if errIP != nil {
				return fmt.Errorf("error parsing flag 'internal-port' for NAT Gateway port forwarding rule update : %v", errIP)
			}

			enabled, errEn := cmd.Flags().GetBool("enabled")
			if errEn != nil {
				return fmt.Errorf("error parsing flag 'enabled' for NAT Gateway port forwarding rule update : %v", errEn)
			}

			o.PortForwardingRuleReq = &govultr.NATGatewayPortForwardingRuleReq{}

			if cmd.Flags().Changed("name") {
				o.PortForwardingRuleReq.Name = name
			}

			if cmd.Flags().Changed("description") {
				o.PortForwardingRuleReq.Description = description
			}

			if cmd.Flags().Changed("internal-ip") {
				o.PortForwardingRuleReq.InternalIP = internalIP
			}

			if cmd.Flags().Changed("protocol") {
				o.PortForwardingRuleReq.Protocol = protocol
			}

			if cmd.Flags().Changed("external-port") {
				o.PortForwardingRuleReq.ExternalPort = externalPort
			}

			if cmd.Flags().Changed("internal-port") {
				o.PortForwardingRuleReq.InternalPort = internalPort
			}

			if cmd.Flags().Changed("enabled") {
				o.PortForwardingRuleReq.Enabled = &enabled
			}

			pfr, err := o.updateNATGatewayPortForwardingRule()
			if err != nil {
				return fmt.Errorf("error updating NAT Gateway port forwarding rule : %v", err)
			}

			data := &NATGatewayPortForwardingRulePrinter{PortForwardingRule: pfr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	portForwardingRuleUpdate.Flags().StringP("name", "n", "", "name for the NAT Gateway port forwarding rule")
	portForwardingRuleUpdate.Flags().StringP("description", "d", "", "description for the NAT Gateway port forwarding rule")
	portForwardingRuleUpdate.Flags().String("internal-ip", "", "internal IP for the NAT Gateway port forwarding rule")
	portForwardingRuleUpdate.Flags().StringP("protocol", "p", "", "protocol for the NAT Gateway port forwarding rule")
	portForwardingRuleUpdate.Flags().Int("external-port", 0, "external port for the NAT Gateway port forwarding rule")
	portForwardingRuleUpdate.Flags().Int("internal-port", 0, "internal port for the NAT Gateway port forwarding rule")
	portForwardingRuleUpdate.Flags().BoolP("enabled", "e", true, "name for the NAT Gateway port forwarding rule")

	// NAT Gateway Port Forwarding Rule Delete
	portForwardingRuleDelete := &cobra.Command{
		Use:     "delete <VPC ID> <NAT Gateway ID> <Port Forwarding Rule ID>",
		Aliases: []string{"d"},
		Short:   "Delete a NAT Gateway port forwarding rule",
		Long:    pfrDeleteLong,
		Example: pfrDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("please provide a VPC ID, a NAT Gateway ID, and a Port Forwarding Rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.delNATGatewayPortForwardingRule(); err != nil {
				return fmt.Errorf("error deleting NAT Gateway port forwarding rule : %v", err)
			}

			o.Base.Printer.Display(printer.Info("NAT Gateway port forwarding rule deleted"), nil)

			return nil
		},
	}

	portForwardingRule.AddCommand(
		portForwardingRuleList,
		portForwardingRuleGet,
		portForwardingRuleCreate,
		portForwardingRuleUpdate,
		portForwardingRuleDelete,
	)

	// NAT Gateway Firewall Rule
	firewallRule := &cobra.Command{
		Use:     "firewall-rule",
		Aliases: []string{"fr"},
		Short:   "Commands to handle NAT Gateway firewall rules",
		Long:    fwrLong,
		Example: fwrExample,
	}

	// NAT Gateway Firewall Rule List
	firewallRuleList := &cobra.Command{
		Use:     "list <VPC ID>",
		Aliases: []string{"l"},
		Short:   "List NAT Gateway firewall rules",
		Long:    fwrListLong,
		Example: fwrListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fwrs, meta, err := o.listNATGatewayFirewallRules()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateway firewall rules : %v", err)
			}

			data := &NATGatewayFirewallRulesPrinter{FirewallRules: fwrs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Firewall Rule Get
	firewallRuleGet := &cobra.Command{
		Use:     "get <VPC ID> <NAT Gateway ID> <Firewall Rule ID>",
		Aliases: []string{"g"},
		Short:   "Get a NAT Gateway firewall rule",
		Long:    fwrGetLong,
		Example: fwrGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("please provide a VPC ID, a NAT Gateway ID, and a Firewall Rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fwr, err := o.getNATGatewayFirewallRule()
			if err != nil {
				return fmt.Errorf("error retrieving NAT Gateway firewall rule : %v", err)
			}

			data := &NATGatewayFirewallRulePrinter{FirewallRule: fwr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// NAT Gateway Firewall Rule Create
	firewallRuleCreate := &cobra.Command{
		Use:     "create <VPC ID> <NAT Gateway ID>",
		Aliases: []string{"c"},
		Short:   "Create a NAT Gateway firewall rule",
		Long:    fwrCreateLong,
		Example: fwrCreateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a VPC ID and a NAT Gateway ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			protocol, errPr := cmd.Flags().GetString("protocol")
			if errPr != nil {
				return fmt.Errorf("error parsing flag 'protocol' for NAT Gateway firewall rule create : %v", errPr)
			}

			subnet, errSu := cmd.Flags().GetString("subnet")
			if errSu != nil {
				return fmt.Errorf("error parsing flag 'subnet' for NAT Gateway firewall rule create : %v", errSu)
			}

			subnetSize, errSS := cmd.Flags().GetInt("subnet-size")
			if errSS != nil {
				return fmt.Errorf("error parsing flag 'protocol' for NAT Gateway firewall rule create : %v", errSS)
			}

			port, errPo := cmd.Flags().GetString("port")
			if errPo != nil {
				return fmt.Errorf("error parsing flag 'port' for NAT Gateway firewall rule create : %v", errPo)
			}

			notes, errNo := cmd.Flags().GetString("notes")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'notes' for NAT Gateway firewall rule create : %v", errNo)
			}

			o.FirewallRuleCreateReq = &govultr.NATGatewayFirewallRuleCreateReq{
				Protocol:   protocol,
				Subnet:     subnet,
				SubnetSize: subnetSize,
				Port:       port,
				Notes:      notes,
			}

			fwr, err := o.createNATGatewayFirewallRule()
			if err != nil {
				return fmt.Errorf("error creating NAT Gateway firewall rule : %v", err)
			}

			data := &NATGatewayFirewallRulePrinter{FirewallRule: fwr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	firewallRuleCreate.Flags().String("protocol", "", "protocol for the new NAT Gateway firewall rule")
	firewallRuleCreate.Flags().StringP("subnet", "s", "", "subnet for the new NAT Gateway firewall rule")
	firewallRuleCreate.Flags().Int("subnet-size", 0, "subnet size for the new NAT Gateway firewall rule")
	firewallRuleCreate.Flags().StringP("port", "p", "", "port or port range for the new NAT Gateway firewall rule")
	firewallRuleCreate.Flags().StringP("notes", "n", "", "notes for the new NAT Gateway firewall rule")

	// NAT Gateway Firewall Rule Update
	firewallRuleUpdate := &cobra.Command{
		Use:     "update <VPC ID> <NAT Gateway ID> <Firewall Rule ID>",
		Aliases: []string{"u"},
		Short:   "Update a NAT Gateway firewall rule",
		Long:    fwrUpdateLong,
		Example: fwrUpdateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("please provide a VPC ID, a NAT Gateway ID, and a Firewall Rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			notes, errNo := cmd.Flags().GetString("notes")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'notes' for NAT Gateway firewall rule update : %v", errNo)
			}

			o.FirewallRuleUpdateReq = &govultr.NATGatewayFirewallRuleUpdateReq{}

			if cmd.Flags().Changed("notes") {
				o.FirewallRuleUpdateReq.Notes = notes
			}

			fwr, err := o.updateNATGatewayFirewallRule()
			if err != nil {
				return fmt.Errorf("error updating NAT Gateway firewall rule : %v", err)
			}

			data := &NATGatewayFirewallRulePrinter{FirewallRule: fwr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	firewallRuleUpdate.Flags().StringP("notes", "n", "", "notes for the NAT Gateway firewall rule")

	// NAT Gateway Firewall Rule Delete
	firewallRuleDelete := &cobra.Command{
		Use:     "delete <VPC ID> <NAT Gateway ID> <Firewall Rule ID>",
		Aliases: []string{"d"},
		Short:   "Delete a NAT Gateway firewall rule",
		Long:    fwrDeleteLong,
		Example: fwrDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("please provide a VPC ID, a NAT Gateway ID, and a Firewall Rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.delNATGatewayFirewallRule(); err != nil {
				return fmt.Errorf("error deleting NAT Gateway firewall rule : %v", err)
			}

			o.Base.Printer.Display(printer.Info("NAT Gateway firewall rule deleted"), nil)

			return nil
		},
	}

	firewallRule.AddCommand(
		firewallRuleList,
		firewallRuleGet,
		firewallRuleCreate,
		firewallRuleUpdate,
		firewallRuleDelete,
	)

	natGateway.AddCommand(
		natGatewayList,
		natGatewayGet,
		natGatewayCreate,
		natGatewayUpdate,
		natGatewayDelete,
		portForwardingRule,
		firewallRule,
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
	PortForwardingRuleReq *govultr.NATGatewayPortForwardingRuleReq
	FirewallRuleCreateReq *govultr.NATGatewayFirewallRuleCreateReq
	FirewallRuleUpdateReq *govultr.NATGatewayFirewallRuleUpdateReq
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

func (o *options) listNATGatewayPortForwardingRules() ([]govultr.NATGatewayPortForwardingRule, *govultr.Meta, error) {
	portForwardingRules, meta, _, err := o.Base.Client.VPC.ListNATGatewayPortForwardingRules(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Options)
	return portForwardingRules, meta, err
}

func (o *options) getNATGatewayPortForwardingRule() (*govultr.NATGatewayPortForwardingRule, error) {
	portForwardingRule, _, err := o.Base.Client.VPC.GetNATGatewayPortForwardingRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2])
	return portForwardingRule, err
}

func (o *options) createNATGatewayPortForwardingRule() (*govultr.NATGatewayPortForwardingRule, error) {
	portForwardingRule, _, err := o.Base.Client.VPC.CreateNATGatewayPortForwardingRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.PortForwardingRuleReq)
	return portForwardingRule, err
}

func (o *options) updateNATGatewayPortForwardingRule() (*govultr.NATGatewayPortForwardingRule, error) {
	portForwardingRule, _, err := o.Base.Client.VPC.UpdateNATGatewayPortForwardingRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2], o.PortForwardingRuleReq)
	return portForwardingRule, err
}

func (o *options) delNATGatewayPortForwardingRule() error {
	return o.Base.Client.VPC.DeleteNATGatewayPortForwardingRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2])
}

func (o *options) listNATGatewayFirewallRules() ([]govultr.NATGatewayFirewallRule, *govultr.Meta, error) {
	firewallRules, meta, _, err := o.Base.Client.VPC.ListNATGatewayFirewallRules(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Options)
	return firewallRules, meta, err
}

func (o *options) getNATGatewayFirewallRule() (*govultr.NATGatewayFirewallRule, error) {
	firewallRule, _, err := o.Base.Client.VPC.GetNATGatewayFirewallRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2])
	return firewallRule, err
}

func (o *options) createNATGatewayFirewallRule() (*govultr.NATGatewayFirewallRule, error) {
	firewallRule, _, err := o.Base.Client.VPC.CreateNATGatewayFirewallRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.FirewallRuleCreateReq)
	return firewallRule, err
}

func (o *options) updateNATGatewayFirewallRule() (*govultr.NATGatewayFirewallRule, error) {
	firewallRule, _, err := o.Base.Client.VPC.UpdateNATGatewayFirewallRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2], o.FirewallRuleUpdateReq)
	return firewallRule, err
}

func (o *options) delNATGatewayFirewallRule() error {
	return o.Base.Client.VPC.DeleteNATGatewayFirewallRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2])
}
