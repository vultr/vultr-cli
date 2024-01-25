// Package reservedip provides the reserved-ip commands to the CLI
package reservedip

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
	long    = `Get all available commands for reserved IPs`
	example = `
	# Full example
	vultr-cli reserved-ip

	# Shortened with aliased commands
	vultr-cli rip
	`

	createLong    = `Create a reserved IP on your Vultr account`
	createExample = `
	# Full Example
	vultr-cli reserved-ip create --region="yto" --type="v4" --label="new IP"

	# Shortened with alias commands
	vultr-cli rip c -r="yto" -t="v4" -l="new IP"
	`

	getLong    = `Get info for a reserved IP on your Vultr account`
	getExample = `
	# Full example
	vultr-cli reserved-ip get 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5

	# Shortened with alias commands
	vultr-cli rip g 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5
	`

	listLong    = `List all reserved IPs on your Vultr account`
	listExample = `
	# Full example
	vultr-cli reserved-ip list

	# Shortened with alias commands
	vultr-cli rip l
	`

	attachLong    = `Attach a reserved IP to an instance on your Vultr account`
	attachExample = `
	# Full example
	vultr-cli reserved-ip attach 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 --instance-id="2b9bf5fb-1644-4e0a-b706-1116ab64d783"

	# Shortened with alias commands
	vultr-cli rip a 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 -i="2b9bf5fb-1644-4e0a-b706-1116ab64d783"
	`

	detachLong    = `Detach a reserved IP from an instance on your Vultr account`
	detachExample = `
	# Full example
	vultr-cli reserved-ip detach 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5

	# Shortened with alias commands
	vultr-cli rip d 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5
	`

	convertLong    = `Convert an instance IP to a reserved IP on your Vultr account`
	convertExample = `
	# Full example
	vultr-cli reserved-ip convert --ip="192.0.2.123" --label="new label converted"

	# Shortened with alias commands
	vultr-cli rip v -i="192.0.2.123" -l="new label converted"
	`

	updateLong    = `Update a reserved IP on your Vultr account`
	updateExample = `
	# Full example
	vultr-cli reserved-ip update 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 --label="new label"

	# Shortened with alias commands
	vultr-cli rip u 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5 -l="new label"
	`

	deleteLong    = `Delete a reserved IP from your Vultr account`
	deleteExample = `
	# Full example
	vultr-cli reserved-ip delete 6a31648d-ebfa-4d43-9a00-9c9f0e5048f5
	`
)

// NewCmdReservedIP provides the CLI command for reserved IP functions
func NewCmdReservedIP(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "reserved-ip",
		Aliases: []string{"rip"},
		Short:   "reserved-ip lets you interact with reserved-ip ",
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
		Short:   "list all reserved IPs",
		Long:    listLong,
		Example: listExample,
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			ips, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving reserved IP list : %v", err)
			}

			data := &ReservedIPsPrinter{IPs: ips, Meta: meta}
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
		Use:     "get <Reserved IP ID>",
		Short:   "get a reserved IP",
		Long:    getLong,
		Example: getExample,
		Aliases: []string{"g"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a reserved IP ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting reserved IP : %v", err)
			}

			data := &ReservedIPPrinter{IP: ip}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create ",
		Short:   "create reserved IP",
		Long:    createLong,
		Example: createExample,
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for reserved-ip create : %v", errRe)
			}

			ipType, errIp := cmd.Flags().GetString("type")
			if errIp != nil {
				return fmt.Errorf("error parsing flag 'type' for reserved-ip create : %v", errIp)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for reserved-ip create : %v", errLa)
			}

			o.CreateReq = &govultr.ReservedIPReq{
				Region: region,
				IPType: ipType,
				Label:  label,
			}

			ip, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating reserved IP : %v", err)
			}

			data := &ReservedIPPrinter{IP: ip}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("region", "r", "", "id of region")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking reserved-ip create 'region' flag required: %v", err)
		os.Exit(1)
	}
	create.Flags().StringP("type", "t", "", "type of IP : v4 or v6")
	if err := create.MarkFlagRequired("type"); err != nil {
		fmt.Printf("error marking reserved-ip create 'type' flag required: %v", err)
		os.Exit(1)
	}
	create.Flags().StringP("label", "l", "", "label")

	// Update
	update := &cobra.Command{
		Use:     "update <Reserved IP ID>",
		Short:   "update reserved IP",
		Long:    updateLong,
		Example: updateExample,
		Aliases: []string{"u"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a reserved IP ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for reserved-ip update : %v", errLa)
			}

			o.UpdateReq = &govultr.ReservedIPUpdateReq{
				Label: govultr.StringToStringPtr(label),
			}

			ip, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating reserved IP : %v", err)
			}

			data := &ReservedIPPrinter{IP: ip}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	update.Flags().StringP("label", "l", "", "label")
	if err := update.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking reserved-ip update 'label' flag required: %v", err)
		os.Exit(1)
	}

	// Attach
	attach := &cobra.Command{
		Use:     "attach <Reserved IP ID>",
		Short:   "attach a reserved IP to an instance",
		Long:    attachLong,
		Example: attachExample,
		Aliases: []string{"a"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a reserved IP ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceID, errIn := cmd.Flags().GetString("instance-id")
			if errIn != nil {
				return fmt.Errorf("error parsing flag 'instance-id' for reserved-ip attach : %v", errIn)
			}

			o.InstanceID = instanceID

			if err := o.attach(); err != nil {
				return fmt.Errorf("error attaching reserved IP : %v", err)
			}

			o.Base.Printer.Display(printer.Info("reserved IP has been attached to instance"), nil)

			return nil
		},
	}

	attach.Flags().StringP("instance-id", "i", "", "id of instance you want to attach")
	if err := attach.MarkFlagRequired("instance-id"); err != nil {
		fmt.Printf("error marking reserved-ip attach 'instance-id' flag required: %v", err)
		os.Exit(1)
	}

	// Detach
	detach := &cobra.Command{
		Use:     "detach <Reserved IP ID>",
		Short:   "detach a reserved IP from an instance",
		Long:    detachLong,
		Example: detachExample,
		Aliases: []string{"d"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a reserved IP ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.detach(); err != nil {
				return fmt.Errorf("error detaching reserved IP : %v", err)
			}

			o.Base.Printer.Display(printer.Info("reserved IP has been detached"), nil)

			return nil
		},
	}

	// Convert
	convert := &cobra.Command{
		Use:     "convert ",
		Short:   "convert IP address to reserved IP",
		Long:    convertLong,
		Example: convertExample,
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, errIp := cmd.Flags().GetString("ip")
			if errIp != nil {
				return fmt.Errorf("error parsing flag 'ip' for reserved-ip convert : %v", errIp)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for reserved-ip convert : %v", errLa)
			}

			o.ConvertReq = &govultr.ReservedIPConvertReq{
				IPAddress: ip,
				Label:     label,
			}

			newIP, err := o.convert()
			if err != nil {
				return fmt.Errorf("error converting reserved IP : %v", err)
			}

			data := &ReservedIPPrinter{IP: newIP}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	convert.Flags().StringP("ip", "i", "", "ip you wish to convert")
	if err := convert.MarkFlagRequired("ip"); err != nil {
		fmt.Printf("error marking reserved-ip convert 'ip' flag required: %v", err)
		os.Exit(1)
	}
	convert.Flags().StringP("label", "l", "", "label")

	// Delete
	del := &cobra.Command{
		Use:     "delete <Reserved IP ID>",
		Short:   "delete a reserved ip",
		Long:    deleteLong,
		Example: deleteExample,
		Aliases: []string{"destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a reserved IP ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.detach(); err != nil {
				return fmt.Errorf("error detaching reserved IP : %v", err)
			}

			o.Base.Printer.Display(printer.Info("reserved IP has been detached"), nil)

			return nil
		},
	}

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		attach,
		detach,
		convert,
		del,
	)

	return cmd
}

type options struct {
	Base       *cli.Base
	CreateReq  *govultr.ReservedIPReq
	UpdateReq  *govultr.ReservedIPUpdateReq
	ConvertReq *govultr.ReservedIPConvertReq
	InstanceID string
}

func (o *options) list() ([]govultr.ReservedIP, *govultr.Meta, error) {
	rips, meta, _, err := o.Base.Client.ReservedIP.List(o.Base.Context, o.Base.Options)
	return rips, meta, err
}

func (o *options) get() (*govultr.ReservedIP, error) {
	rip, _, err := o.Base.Client.ReservedIP.Get(o.Base.Context, o.Base.Args[0])
	return rip, err
}

func (o *options) create() (*govultr.ReservedIP, error) {
	rip, _, err := o.Base.Client.ReservedIP.Create(o.Base.Context, o.CreateReq)
	return rip, err
}

func (o *options) update() (*govultr.ReservedIP, error) {
	rip, _, err := o.Base.Client.ReservedIP.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
	return rip, err
}

func (o *options) attach() error {
	return o.Base.Client.ReservedIP.Attach(o.Base.Context, o.Base.Args[0], o.InstanceID)
}

func (o *options) detach() error {
	return o.Base.Client.ReservedIP.Detach(o.Base.Context, o.Base.Args[0])
}

func (o *options) convert() (*govultr.ReservedIP, error) {
	ip, _, err := o.Base.Client.ReservedIP.Convert(o.Base.Context, o.ConvertReq)
	return ip, err
}

func (o *options) del() error {
	return o.Base.Client.ReservedIP.Delete(o.Base.Context, o.Base.Args[0])
}
