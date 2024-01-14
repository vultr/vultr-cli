package baremetal

import (
	"encoding/base64"
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
	long    = `Show all commands available to bare-metal`
	example = `
	# Full example
	vultr-cli bare-metal
	`

	listLong      = ``
	listExample   = ``
	getLong       = ``
	getExample    = ``
	createLong    = ``
	createExample = ``
	deleteLong    = ``
	deleteExample = ``

	haltLong = `
	Halt a bare metal server. This is a hard power off, meaning that the power
	to the machine is severed.  The data on the machine will not be modified,
	and you will still be billed for the machine.
	`

	haltExample = ``

	startLong        = ``
	startExample     = ``
	rebootLong       = ``
	rebootExample    = ``
	reinstallLong    = ``
	reinstallExample = ``

	tagsLong    = `Update the tags on a bare metal server`
	tagsExample = `
	# Full example
	vultr-cli bare-metal tags <bareMetalID> tags="tag-1,tag-2"

	# Shortened example with aliases
	vultr-cli bm tags <bareMetalID> -t="tag-1,tag-2"
	`

	vpc2AttachLong    = `Attaches an existing VPC 2.0 network to the specified bare metal server`
	vpc2AttachExample = `
	# Full example
	vultr-cli bare-metal vpc2 attach <bareMetalID> --vpc-id="2126b7d9-5e2a-491e-8840-838aa6b5f294"
	`
	vpc2DetachLong    = `Detaches an existing VPC 2.0 network from the specified bare metal server`
	vpc2DetachExample = `
	# Full example
	vultr-cli bare-metal vpc2 detach <bareMetalID> --vpc-id="2126b7d9-5e2a-491e-8840-838aa6b5f294"
	`
)

// app         app is used to access bare metal server application commands
// bandwidth   Get a bare metal server's bandwidth usage
// create      create a bare metal server
// delete      Delete a bare metal server
// get         Get a bare metal server by <bareMetalID>
// halt        Halt a bare metal server.
// image       image is used to access bare metal server image commands
// ipv4        List the IPv4 information of a bare metal server.
// ipv6        List the IPv6 information of a bare metal server.
// list        List all bare metal servers.
// os          os is used to access bare metal server operating system commands
// reboot      Reboot a bare metal server. This is a hard reboot, which means that the server is powered off, then back on.
// reinstall   Reinstall the operating system on a bare metal server.
// start       Start a bare metal server.
// tags        Add or modify tags on the bare metal server.
// user-data   user-data is used to access bare metal server user-data commands
// vnc         Get a bare metal server's VNC url by <bareMetalID>
// vpc2        commands to handle vpc 2.0 on a server

type BareMetalOptionsInterface interface {
	setOptions(cmd *cobra.Command, args []string)
	List() []govultr.BareMetalServer
	Get() *govultr.BareMetalServer
	Create() *govultr.BareMetalServer
	// Update() *govultr.BareMetalServer
	Delete() bool
	Halt() bool
	Start() bool
	Reboot() bool
	Reinstall() bool
	Application() string
	OS() string
	Image() string
	IPv4() string
	IPv6() string
	Tags() []string
	UserData() string
	VNC() string
	VPC2() string
	Bandwidth() string
}

// BareMetalOptions ...
type BareMetalOptions struct {
	Base      *cli.Base
	CreateReq *govultr.BareMetalCreate
	UpdateReq *govultr.BareMetalUpdate
}

// NewBareMetalOptions ...
func NewBareMetalOptions(base *cli.Base) *BareMetalOptions {
	return &BareMetalOptions{Base: base}
}

// NewCmdBareMetal ...
func NewCmdBareMetal(base *cli.Base) *cobra.Command {
	o := NewBareMetalOptions(base)

	cmd := &cobra.Command{
		Use:     "bare-metal",
		Short:   "bare-metal is used to access bare metal server commands",
		Aliases: []string{"bm"},
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
		Short:   "List all bare metal servers.",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)
			list, meta, err := o.List()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving bare metal list : %v", err))
				os.Exit(1)
			}
			data := &BareMetalsPrinter{BareMetals: list, Meta: meta}
			o.Base.Printer.Display(data, err)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Get
	get := &cobra.Command{
		Use:     "get <BARE_METAL_ID>",
		Short:   "Get a bare metal server by ID",
		Aliases: []string{"l"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			bm, err := o.Get()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving bare metal : %v", err))
				os.Exit(1)
			}
			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "create a bare metal server",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			o.CreateReq = parseCreateFlags(cmd)
			bm, err := o.Create()
			if err != nil {
				printer.Error(fmt.Errorf("error with bare metal create : %v", err))
				os.Exit(1)
			}

			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)
		},
	}

	create.Flags().StringP("region", "r", "", "ID of the region where the server will be created.")
	create.Flags().StringP("plan", "p", "", "ID of the plan that the server will subscribe to.")
	create.Flags().Int("os", 0, "ID of the operating system that will be installed on the server.")
	create.Flags().StringP(
		"script",
		"s",
		"",
		"(optional) ID of the startup script that will run after the server is created.",
	)
	create.Flags().StringP(
		"snapshot",
		"",
		"",
		"(optional) ID of the snapshot that the server will be restored from.",
	)
	create.Flags().StringP(
		"ipv6",
		"i",
		"",
		"(optional) Whether IPv6 is enabled on the server. Possible values: 'yes', 'no'. Defaults to 'no'.",
	)
	create.Flags().StringP("label", "l", "", "(optional) The label to assign to the server.")
	create.Flags().StringSliceP(
		"ssh",
		"k",
		[]string{},
		"(optional) Comma separated list of SSH key IDs that will be added to the server.",
	)
	create.Flags().IntP(
		"app",
		"a",
		0,
		"(optional) ID of the application that will be installed on the server.",
	)
	create.Flags().StringP("image", "", "", "(optional) Image ID of the application that will be installed on the server.")
	create.Flags().StringP(
		"userdata",
		"u",
		"",
		"(optional) A generic data store, which some provisioning tools and cloud operating systems use as a configuration file.",
	)
	create.Flags().StringP(
		"notify",
		"n",
		"",
		"(optional) Whether an activation email will be sent when the server is ready. Possible values: 'yes', 'no'. Defaults to 'yes'.",
	)
	create.Flags().StringP("hostname", "m", "", "(optional) The hostname to assign to the server.")
	create.Flags().StringP("tag", "t", "", "Deprecated: use `tags` instead. (optional) The tag to assign to the server.")
	create.Flags().StringSliceP("tags", "", []string{}, "(optional) A comma separated list of tags to assign to the server.")
	create.Flags().StringP("ripv4", "v", "", "(optional) IP address of the floating IP to use as the main IP of this server.")
	create.Flags().BoolP("persistent_pxe", "x", false, "enable persistent_pxe | true or false")

	if err := create.MarkFlagRequired("region"); err != nil {
		printer.Error(fmt.Errorf("error marking bare metal create 'region' flag required: %v", err))
		os.Exit(1)
	}

	if err := create.MarkFlagRequired("plan"); err != nil {
		printer.Error(fmt.Errorf("error marking bare metal create 'plan' flag required: %v", err))
		os.Exit(1)
	}

	installFlags := []string{"app", "snapshot", "os", "image"}
	create.MarkFlagsMutuallyExclusive(installFlags...)
	create.MarkFlagsOneRequired(installFlags...)

	// Delete
	del := &cobra.Command{
		Use:     "delete <BARE_METAL_ID>",
		Short:   "Delete a bare metal server",
		Aliases: []string{"destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Delete(); err != nil {
				printer.Error(fmt.Errorf("error deleting bare metal : %v", err))
				os.Exit(1)
			}
			o.Base.Printer.Display(printer.Info("bare metal server has been deleted"), nil)
		},
	}

	// Halt
	halt := &cobra.Command{
		Use:     "halt <BARE_METAL_ID>",
		Short:   "Halt a bare metal server.",
		Aliases: []string{"h"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Halt(); err != nil {
				printer.Error(fmt.Errorf("error halting bare metal : %v", err))
				os.Exit(1)
			}
			o.Base.Printer.Display(printer.Info("bare metal server has been halted"), nil)
		},
	}

	cmd.AddCommand(get, list, create, del, halt)
	return cmd
}

// List ...
func (b *BareMetalOptions) List() ([]govultr.BareMetalServer, *govultr.Meta, error) {
	bms, meta, _, err := b.Base.Client.BareMetalServer.List(b.Base.Context, b.Base.Options)
	return bms, meta, err
}

// Get ...
func (b *BareMetalOptions) Get() (*govultr.BareMetalServer, error) {
	bm, _, err := b.Base.Client.BareMetalServer.Get(b.Base.Context, b.Base.Args[0])
	return bm, err
}

// Create ...
func (b *BareMetalOptions) Create() (*govultr.BareMetalServer, error) {
	bm, _, err := b.Base.Client.BareMetalServer.Create(b.Base.Context, b.CreateReq)
	return bm, err
}

// Delete ...
func (b *BareMetalOptions) Delete() error {
	return b.Base.Client.BareMetalServer.Delete(b.Base.Context, b.Base.Args[0])
}

// Halt ...
func (b *BareMetalOptions) Halt() error {
	return b.Base.Client.BareMetalServer.Halt(b.Base.Context, b.Base.Args[0])
}

// Start ...
func (b *BareMetalOptions) Start() error {
	return b.Base.Client.BareMetalServer.Start(b.Base.Context, b.Base.Args[0])
}

// Reboot ...
func (b *BareMetalOptions) Reboot() error {
	return b.Base.Client.BareMetalServer.Reboot(b.Base.Context, b.Base.Args[0])
}

// Reinstall...
func (b *BareMetalOptions) Reinstall() (*govultr.BareMetalServer, error) {
	bm, _, err := b.Base.Client.BareMetalServer.Reinstall(b.Base.Context, b.Base.Args[0])
	return bm, err
}

// ============================

func parseCreateFlags(cmd *cobra.Command) *govultr.BareMetalCreate {
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing region flag for bare metal create : %v", err))
		os.Exit(1)
	}

	plan, err := cmd.Flags().GetString("plan")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing plan flag for bare metal create : %v", err))
		os.Exit(1)
	}

	osID, err := cmd.Flags().GetInt("os")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing os flag for bare metal create : %v", err))
		os.Exit(1)
	}

	script, err := cmd.Flags().GetString("script")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing script flag bare metal create : %v", err))
		os.Exit(1)
	}

	snapshot, err := cmd.Flags().GetString("snapshot")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing snapshot flag for bare metal create : %v", err))
		os.Exit(1)
	}

	ipv6, err := cmd.Flags().GetString("ipv6")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing ipv6 flag for bare metal create : %v", err))
		os.Exit(1)
	}

	label, err := cmd.Flags().GetString("label")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing label flag for bare metal create : %v", err))
		os.Exit(1)
	}

	sshKeys, err := cmd.Flags().GetStringSlice("ssh")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing ssh flag for bare metal create : %v", err))
		os.Exit(1)
	}

	app, err := cmd.Flags().GetInt("app")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing app flag for bare metal create : %v", err))
		os.Exit(1)
	}

	userdata, err := cmd.Flags().GetString("userdata")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing userdata flag for bare metal create : %v", err))
		os.Exit(1)
	}

	notify, err := cmd.Flags().GetString("notify")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing notify flag for bare metal create : %v", err))
		os.Exit(1)
	}

	hostname, err := cmd.Flags().GetString("hostname")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing hostname flag for bare metal create : %v", err))
		os.Exit(1)
	}

	tags, err := cmd.Flags().GetStringSlice("tags")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing tags flag for bare metal create : %v", err))
		os.Exit(1)
	}

	ripv4, err := cmd.Flags().GetString("ripv4")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing ripv4 flag for bare metal create : %v", err))
		os.Exit(1)
	}

	pxe, err := cmd.Flags().GetBool("persistenterrpxe")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing persistenterrpxe flag for bare metal create : %v", err))
		os.Exit(1)
	}

	image, err := cmd.Flags().GetString("image")
	if err != nil {
		printer.Error(fmt.Errorf("error parsing image flag for bare metal create : %v", err))
		os.Exit(1)
	}

	options := &govultr.BareMetalCreate{
		StartupScriptID: script,
		Plan:            plan,
		SnapshotID:      snapshot,
		AppID:           app,
		OsID:            osID,
		ImageID:         image,
		Label:           label,
		SSHKeyIDs:       sshKeys,
		Hostname:        hostname,
		Tags:            tags,
		ReservedIPv4:    ripv4,
		Region:          region,
		PersistentPxe:   govultr.BoolToBoolPtr(pxe),
	}
	if userdata != "" {
		options.UserData = base64.StdEncoding.EncodeToString([]byte(userdata))
	}

	if notify == "yes" {
		options.ActivationEmail = govultr.BoolToBoolPtr(true)
	}

	if ipv6 == "yes" {
		options.EnableIPv6 = govultr.BoolToBoolPtr(true)
	}

	return options
}
