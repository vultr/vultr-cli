// Package baremetal provides functionality to perform operations on
// bare metal servers through the CLI
package baremetal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/applications"
	"github.com/vultr/vultr-cli/v3/cmd/ip"
	"github.com/vultr/vultr-cli/v3/cmd/operatingsystems"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/userdata"
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

	startLong     = ``
	startExample  = ``
	rebootLong    = `This is a hard reboot, which means that the server is powered off, then back on.`
	rebootExample = ``

	reinstallLong = `Reinstall the operating system on a bare metal server.
All data will be permanently lost, but the IP address will remain the same.
There is no going back from this call.`
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

	applicationLong    = ``
	applicationExample = ``

	applicationChangeLong    = ``
	applicationChangeExample = ``

	applicationListLong    = ``
	applicationListExample = ``

	imageLong          = ``
	imageExample       = ``
	imageChangeLong    = ``
	imageChangeExample = ``

	operatingSystemLong          = ``
	operatingSystemExample       = ``
	operatingSystemChangeLong    = ``
	operatingSystemChangeExample = ``
	operatingSystemListLong      = ``
	operatingSystemListExample   = ``

	userDataLong       = ``
	userDataExample    = ``
	userDataGetLong    = ``
	userDataGetExample = ``
	userDataSetLong    = ``
	userDataSetExample = ``

	vncLong    = ``
	vncExample = ``

	bandwidthLong    = ``
	bandwidthExample = ``

	ipv4Long    = `IP information is only available for bare metal servers in the "active" state.`
	ipv4Example = ``

	ipv6Long = `
List the IPv6 information of a bare metal server. IP information is only available for bare metal servers in the "active" state.
`
	ipv6Example = ``

	vpc2Long        = ``
	vpc2ListLong    = ``
	vpc2ListExample = ``
)

// NewCmdBareMetal ...
func NewCmdBareMetal(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := options{Base: base}

	cmd := &cobra.Command{
		Use:     "bare-metal",
		Short:   "Commands to manage bare metal servers",
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
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			list, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal list : %v", err)
			}
			data := &BareMetalsPrinter{BareMetals: list, Meta: meta}
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

	// Get
	get := &cobra.Command{
		Use:     "get <Bare Metal ID>",
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
		RunE: func(cmd *cobra.Command, args []string) error {
			bm, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal : %v", err)
			}
			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a bare metal server",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			req, errParse := parseCreateFlags(cmd)
			if errParse != nil {
				return fmt.Errorf("error parsing flags for bare metal create : %v", errParse)
			}

			o.CreateReq = req

			bm, err := o.create()
			if err != nil {
				return fmt.Errorf("error with bare metal create : %v", err)
			}

			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)

			return nil
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
		fmt.Printf("error marking bare metal create 'region' flag required: %v", err)
		os.Exit(1)
	}

	if err := create.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking bare metal create 'plan' flag required: %v", err)
		os.Exit(1)
	}

	installFlags := []string{"app", "snapshot", "os", "image"}
	create.MarkFlagsMutuallyExclusive(installFlags...)
	create.MarkFlagsOneRequired(installFlags...)

	// Delete
	del := &cobra.Command{
		Use:     "delete <Bare Metal ID>",
		Short:   "Delete a bare metal server",
		Aliases: []string{"destroy"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting bare metal : %v", err)
			}
			o.Base.Printer.Display(printer.Info("bare metal server has been deleted"), nil)

			return nil
		},
	}

	// Halt
	halt := &cobra.Command{
		Use:     "halt <Bare Metal ID>",
		Short:   "Halt a bare metal server.",
		Aliases: []string{"h"},
		Long:    haltLong,
		Example: haltExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.halt(); err != nil {
				return fmt.Errorf("error halting bare metal : %v", err)
			}
			o.Base.Printer.Display(printer.Info("bare metal server has been halted"), nil)

			return nil
		},
	}

	// Start
	start := &cobra.Command{
		Use:     "start <Bare Metal ID>",
		Short:   "Start a bare metal server.",
		Long:    startLong,
		Example: startExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.start(); err != nil {
				return fmt.Errorf("error starting bare metal : %v", err)
			}

			o.Base.Printer.Display(printer.Info("bare metal server has been started"), nil)

			return nil
		},
	}

	// Reboot
	reboot := &cobra.Command{
		Use:     "reboot <Bare Metal ID>",
		Short:   "Reboot a bare metal server.",
		Aliases: []string{"r"},
		Long:    rebootLong,
		Example: rebootExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.reboot(); err != nil {
				return fmt.Errorf("error rebooting bare metal : %v", err)
			}

			o.Base.Printer.Display(printer.Info("bare metal server has been rebooted"), nil)

			return nil
		},
	}

	// Reinstall
	reinstall := &cobra.Command{
		Use:     "reinstall <bareMetalID>",
		Short:   "Reinstall the operating system on a bare metal server.",
		Long:    reinstallLong,
		Example: reinstallExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.reinstall(); err != nil {
				return fmt.Errorf("error reinstalling bare metal : %v", err)
			}

			o.Base.Printer.Display(printer.Info("bare metal server has initiated reinstallation"), nil)

			return nil
		},
	}

	// Application
	application := &cobra.Command{
		Use:     "app",
		Short:   "Commands to modify applications on a bare metal server",
		Aliases: []string{"a", "application"},
		Long:    applicationLong,
		Example: applicationExample,
	}

	// Application Change
	applicationChange := &cobra.Command{
		Use:     "change <Bare Metal ID>",
		Short:   "Change a bare metal server application",
		Aliases: []string{"c"},
		Long:    applicationChangeLong,
		Example: applicationChangeExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appID, err := cmd.Flags().GetInt("app")
			if err != nil {
				return fmt.Errorf("error parsing app flag for bare metal app ID change : %v", err)
			}

			o.UpdateReq = &govultr.BareMetalUpdate{
				AppID: appID,
			}

			bm, err := o.update()
			if err != nil {
				return fmt.Errorf("error with bare metal update : %v", err)
			}

			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	applicationChange.Flags().IntP(
		"app",
		"a",
		0,
		"ID of the application that will be installed on the server",
	)

	if err := applicationChange.MarkFlagRequired("app"); err != nil {
		fmt.Printf("error marking bare metal 'app' flag required : %v", err)
		os.Exit(1)
	}

	// Application List
	applicationList := &cobra.Command{
		Use:     "list <Bare Metal ID>",
		Short:   "List available apps for a bare metal server",
		Long:    applicationListLong,
		Example: applicationListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrades, err := o.getUpgrades()
			if err != nil {
				return fmt.Errorf("error with bare metal get upgrades : %v", err)
			}
			data := &applications.ApplicationsPrinter{Applications: upgrades.Applications}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	application.AddCommand(applicationChange, applicationList)

	// Image
	image := &cobra.Command{
		Use:     "image",
		Short:   "Commands to manage the image on a bare metal server",
		Aliases: []string{"i"},
		Long:    imageLong,
		Example: imageExample,
	}

	// Image Change
	imageChange := &cobra.Command{
		Use:     "change <Bare Metal ID>",
		Short:   "Change a bare metal server's image",
		Aliases: []string{"c"},
		Long:    imageChangeLong,
		Example: imageChangeExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			imageID, err := cmd.Flags().GetString("image")
			if err != nil {
				return fmt.Errorf("error parsing image flag for bare metal image change : %v", err)
			}

			o.UpdateReq = &govultr.BareMetalUpdate{
				ImageID: imageID,
			}

			bm, err := o.update()
			if err != nil {
				return fmt.Errorf("error with bare metal image update : %v", err)
			}

			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	imageChange.Flags().StringP(
		"image",
		"i",
		"",
		"ID of the image that will be installed on the server",
	)

	if err := imageChange.MarkFlagRequired("image"); err != nil {
		fmt.Printf("error marking bare metal 'image' flag required : %v", err)
		os.Exit(1)
	}

	image.AddCommand(imageChange)

	// OS
	operatingSystem := &cobra.Command{
		Use:     "os",
		Short:   "Commands to manage the operating system on a bare metal server",
		Aliases: []string{"o"},
		Long:    operatingSystemLong,
		Example: operatingSystemExample,
	}

	// OS Change
	operatingSystemChange := &cobra.Command{
		Use:     "change <Bare Metal ID>",
		Short:   "Change a bare metal server's operating system",
		Aliases: []string{"c"},
		Long:    operatingSystemChangeLong,
		Example: operatingSystemChangeExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			operatingSystemID, err := cmd.Flags().GetInt("os")
			if err != nil {
				return fmt.Errorf("error parsing os flag for bare metal os change : %v", err)
			}

			o.UpdateReq = &govultr.BareMetalUpdate{
				OsID: operatingSystemID,
			}

			bm, err := o.update()
			if err != nil {
				return fmt.Errorf("error with bare metal os update : %v", err)
			}

			data := &BareMetalPrinter{BareMetal: *bm}
			o.Base.Printer.Display(data, err)

			return nil
		},
	}

	operatingSystemChange.Flags().IntP(
		"os",
		"o",
		0,
		"ID of the operating system that will be installed on the server",
	)
	if err := operatingSystemChange.MarkFlagRequired("os"); err != nil {
		fmt.Printf("error marking bare metal 'os' flag required : %v", err)
		os.Exit(1)
	}

	// Operating System List
	operatingSystemList := &cobra.Command{
		Use:     "list <Bare Metal ID>",
		Short:   "List the available operating systems for a bare metal server",
		Long:    operatingSystemListLong,
		Example: operatingSystemListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrades, err := o.getUpgrades()
			if err != nil {
				return fmt.Errorf("error with bare metal get upgrades : %v", err)
			}

			data := &operatingsystems.OSPrinter{OperatingSystems: upgrades.OS}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	operatingSystem.AddCommand(operatingSystemChange, operatingSystemList)

	// User Data
	userData := &cobra.Command{
		Use:     "user-data",
		Short:   "Commands for bare metal server user data",
		Aliases: []string{"u"},
		Long:    userDataLong,
		Example: userDataExample,
	}

	// User Data Get
	userDataGet := &cobra.Command{
		Use:     "get <Bare Metal ID>",
		Short:   "Get the user data of a bare metal server",
		Aliases: []string{"g"},
		Long:    userDataGetLong,
		Example: userDataGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ud, err := o.getUserData()
			if err != nil {
				return fmt.Errorf("error with bare metal get user data : %v", err)
			}

			data := &userdata.UserDataPrinter{UserData: *ud}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// User Data Set
	userDataSet := &cobra.Command{
		Use:     "set <Bare Metal ID>",
		Short:   "Set the plain text user-data of a bare metal server",
		Long:    userDataSetLong,
		Example: userDataSetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString("user-data")
			if err != nil {
				return fmt.Errorf("error parsing user-data flag for bare metal user data set : %v", err)
			}

			rawData, err := os.ReadFile(filepath.Clean(path))
			if err != nil {
				return fmt.Errorf("error reading user-data : %v", err)
			}

			o.UpdateReq = &govultr.BareMetalUpdate{
				UserData: base64.StdEncoding.EncodeToString(rawData),
			}

			_, errUpdate := o.update()
			if err != nil {
				return fmt.Errorf("error updating bare metal user-data : %v", errUpdate)
			}

			o.Base.Printer.Display(printer.Info("bare metal server user data has been set"), nil)

			return nil
		},
	}

	userDataSet.Flags().StringP("user-data", "d", "/dev/stdin", "file to read userdata from")
	if err := userDataSet.MarkFlagRequired("user-data"); err != nil {
		fmt.Printf("error marking bare metal 'user-data' flag required: %v", err)
		os.Exit(1)
	}

	userData.AddCommand(userDataGet, userDataSet)

	// VNC URL
	vnc := &cobra.Command{
		Use:     "vnc <Bare Metal ID>",
		Short:   "Get a bare metal server's VNC url",
		Long:    vncLong,
		Example: vncExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vnc, err := o.getVNCURL()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal VNC URL : %v", err)
			}

			data := &BareMetalVNCPrinter{VNC: *vnc}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Bandwidth
	bandwidth := &cobra.Command{
		Use:     "bandwidth <Bare Metal ID>",
		Short:   "Get a bare metal server's bandwidth usage",
		Aliases: []string{"b"},
		Long:    bandwidthLong,
		Example: bandwidthExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bw, err := o.getBandwidth()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal bandwidth usage : %v", err)
			}

			data := &BareMetalBandwidthPrinter{Bandwidth: *bw}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Tags
	tags := &cobra.Command{
		Use:     "tags <Bare Metal ID>",
		Short:   "Add or modify tags on the bare metal server",
		Long:    tagsLong,
		Example: tagsExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			tags, _ := cmd.Flags().GetStringSlice("tags")
			o.UpdateReq = &govultr.BareMetalUpdate{
				Tags: tags,
			}

			_, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating bare metal tags : %v", err)
			}

			o.Base.Printer.Display(printer.Info("bare metal server tags have been updated"), nil)

			return nil
		},
	}

	tags.Flags().StringSliceP("tags", "t", []string{}, "A comma separated list of tags to apply to the server")
	if err := tags.MarkFlagRequired("tags"); err != nil {
		fmt.Printf("error marking bare metal 'tags' flag required: %v", err)
		os.Exit(1)
	}

	// IPv4 Addresses
	ipv4 := &cobra.Command{
		Use:     "ipv4 <Bare Metal ID>",
		Short:   "List the IPv4 information of a bare metal server",
		Long:    ipv4Long,
		Example: ipv4Example,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			ipv4, meta, err := o.getIPv4Addresses()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal IPv4 information : %v", err)
			}
			data := &ip.IPv4sPrinter{IPv4s: ipv4, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// IPv6 Addresses
	ipv6 := &cobra.Command{
		Use:     "ipv6 <Bare Metal ID>",
		Short:   "Get the IPv6 information of a bare metal server",
		Long:    ipv6Long,
		Example: ipv6Example,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			ipv6, meta, err := o.getIPv6Addresses()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal IPv6 information : %v", err)
			}
			data := &ip.IPv6sPrinter{IPv6s: ipv6, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// VPC2
	vpc2 := &cobra.Command{
		Use:   "vpc2",
		Short: "Commands to manage VPC2s on bare metal servers",
		Long:  vpc2Long,
	}

	// VPC2 List
	vpc2List := &cobra.Command{
		Use:     "list <Bare Metal ID>",
		Short:   "List all VPC2 networks attached to a server",
		Aliases: []string{"l"},
		Long:    vpc2ListLong,
		Example: vpc2ListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vpc2s, err := o.vpc2NetworksList()
			if err != nil {
				return fmt.Errorf("error retrieving bare metal vpc2 information : %v", err)
			}
			data := &BareMetalVPC2sPrinter{VPC2s: vpc2s}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// VPC2 Attach
	vpc2Attach := &cobra.Command{
		Use:     "attach <Bare Metal ID>",
		Short:   "Attach a VPC2 network to a server",
		Long:    vpc2AttachLong,
		Example: vpc2AttachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vpcID, errID := cmd.Flags().GetString("vpc-id")
			if errID != nil {
				return fmt.Errorf("error parsing vpc-id flag for bare metal VPC2 attach : %v", errID)
			}

			IPAddress, errIP := cmd.Flags().GetString("ip-address")
			if errIP != nil {
				return fmt.Errorf("error parsing ip-address flag for bare metal VPC2 attach : %v", errIP)
			}

			o.VPC2Req = &govultr.AttachVPC2Req{
				VPCID:     vpcID,
				IPAddress: &IPAddress,
			}

			if err := o.vpc2NetworksAttach(); err != nil {
				return fmt.Errorf("error attaching bare metal to VPC2 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("bare metal server has been attached to VPC2 network"), nil)

			return nil
		},
	}

	vpc2Attach.Flags().StringP("vpc-id", "v", "", "the ID of the VPC 2.0 network you wish to attach")
	vpc2Attach.Flags().StringP("ip-address", "i", "", "the IP address to use for this server on the attached VPC 2.0 network")
	if errVPC := vpc2Attach.MarkFlagRequired("vpc-id"); errVPC != nil {
		fmt.Printf("error marking bare metal 'vpc-id' flag required for attach : %v", errVPC)
		os.Exit(1)
	}

	// VPC2 Detach
	vpc2Detach := &cobra.Command{
		Use:     "detach <Bare Metal ID>",
		Short:   "Detach a VPC2 network from a server",
		Long:    vpc2DetachLong,
		Example: vpc2DetachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a bare metal ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vpcID, errID := cmd.Flags().GetString("vpc-id")
			if errID != nil {
				return fmt.Errorf("error parsing vpc-id flag for bare metal VPC2 detach : %v", errID)
			}

			o.VPC2ID = vpcID
			if err := o.vpc2NetworksDetach(); err != nil {
				return fmt.Errorf("error detaching bare metal VPC2 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("bare metal server has been detached from VPC2 network"), nil)

			return nil
		},
	}

	vpc2Detach.Flags().StringP("vpc-id", "v", "", "the ID of the VPC 2.0 network you wish to detach")
	if errVPC2 := vpc2Detach.MarkFlagRequired("vpc-id"); errVPC2 != nil {
		fmt.Printf("error marking bare metal 'vpc-id' flag required for detach : %v", errVPC2)
		os.Exit(1)
	}

	vpc2.AddCommand(vpc2List, vpc2Attach, vpc2Detach)

	cmd.AddCommand(
		get,
		list,
		create,
		del,
		halt,
		start,
		reboot,
		reinstall,
		application,
		image,
		operatingSystem,
		userData,
		vnc,
		bandwidth,
		tags,
		ipv4,
		ipv6,
		vpc2,
	)

	return cmd
}

type options struct {
	Base      *cli.Base
	CreateReq *govultr.BareMetalCreate
	UpdateReq *govultr.BareMetalUpdate
	VPC2Req   *govultr.AttachVPC2Req
	VPC2ID    string
}

func (b *options) list() ([]govultr.BareMetalServer, *govultr.Meta, error) {
	bms, meta, _, err := b.Base.Client.BareMetalServer.List(b.Base.Context, b.Base.Options)
	return bms, meta, err
}

func (b *options) get() (*govultr.BareMetalServer, error) {
	bm, _, err := b.Base.Client.BareMetalServer.Get(b.Base.Context, b.Base.Args[0])
	return bm, err
}

func (b *options) create() (*govultr.BareMetalServer, error) {
	bm, _, err := b.Base.Client.BareMetalServer.Create(b.Base.Context, b.CreateReq)
	return bm, err
}

func (b *options) update() (*govultr.BareMetalServer, error) {
	bm, _, err := b.Base.Client.BareMetalServer.Update(b.Base.Context, b.Base.Args[0], b.UpdateReq)
	return bm, err
}

func (b *options) del() error {
	return b.Base.Client.BareMetalServer.Delete(b.Base.Context, b.Base.Args[0])
}

func (b *options) halt() error {
	return b.Base.Client.BareMetalServer.Halt(b.Base.Context, b.Base.Args[0])
}

func (b *options) start() error {
	return b.Base.Client.BareMetalServer.Start(b.Base.Context, b.Base.Args[0])
}

func (b *options) reboot() error {
	return b.Base.Client.BareMetalServer.Reboot(b.Base.Context, b.Base.Args[0])
}

func (b *options) reinstall() error {
	_, _, err := b.Base.Client.BareMetalServer.Reinstall(b.Base.Context, b.Base.Args[0])
	return err
}

func (b *options) getUpgrades() (*govultr.Upgrades, error) {
	list, _, err := b.Base.Client.BareMetalServer.GetUpgrades(b.Base.Context, b.Base.Args[0])
	return list, err
}

func (b *options) getUserData() (*govultr.UserData, error) {
	ud, _, err := b.Base.Client.BareMetalServer.GetUserData(b.Base.Context, b.Base.Args[0])
	return ud, err
}

func (b *options) getVNCURL() (*govultr.VNCUrl, error) {
	url, _, err := b.Base.Client.BareMetalServer.GetVNCUrl(b.Base.Context, b.Base.Args[0])
	return url, err
}

func (b *options) getBandwidth() (*govultr.Bandwidth, error) {
	bw, _, err := b.Base.Client.BareMetalServer.GetBandwidth(b.Base.Context, b.Base.Args[0])
	return bw, err
}

func (b *options) getIPv4Addresses() ([]govultr.IPv4, *govultr.Meta, error) {
	ips, meta, _, err := b.Base.Client.BareMetalServer.ListIPv4s(b.Base.Context, b.Base.Args[0], b.Base.Options)
	return ips, meta, err
}

func (b *options) getIPv6Addresses() ([]govultr.IPv6, *govultr.Meta, error) {
	ips, meta, _, err := b.Base.Client.BareMetalServer.ListIPv6s(b.Base.Context, b.Base.Args[0], b.Base.Options)
	return ips, meta, err
}

func (b *options) vpc2NetworksList() ([]govultr.VPC2Info, error) {
	vpc2s, _, err := b.Base.Client.BareMetalServer.ListVPC2Info(b.Base.Context, b.Base.Args[0])
	return vpc2s, err
}

func (b *options) vpc2NetworksAttach() error {
	return b.Base.Client.BareMetalServer.AttachVPC2(b.Base.Context, b.Base.Args[0], b.VPC2Req)
}

func (b *options) vpc2NetworksDetach() error {
	return b.Base.Client.BareMetalServer.DetachVPC2(b.Base.Context, b.Base.Args[0], b.VPC2ID)
}

// ============================

func parseCreateFlags(cmd *cobra.Command) (*govultr.BareMetalCreate, error) { //nolint:funlen,gocyclo
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return nil, fmt.Errorf("error parsing region flag for bare metal create : %v", err)
	}

	plan, err := cmd.Flags().GetString("plan")
	if err != nil {
		return nil, fmt.Errorf("error parsing plan flag for bare metal create : %v", err)
	}

	osID, err := cmd.Flags().GetInt("os")
	if err != nil {
		return nil, fmt.Errorf("error parsing os flag for bare metal create : %v", err)
	}

	script, err := cmd.Flags().GetString("script")
	if err != nil {
		return nil, fmt.Errorf("error parsing script flag bare metal create : %v", err)
	}

	snapshot, err := cmd.Flags().GetString("snapshot")
	if err != nil {
		return nil, fmt.Errorf("error parsing snapshot flag for bare metal create : %v", err)
	}

	ipv6, err := cmd.Flags().GetString("ipv6")
	if err != nil {
		return nil, fmt.Errorf("error parsing ipv6 flag for bare metal create : %v", err)
	}

	label, err := cmd.Flags().GetString("label")
	if err != nil {
		return nil, fmt.Errorf("error parsing label flag for bare metal create : %v", err)
	}

	sshKeys, err := cmd.Flags().GetStringSlice("ssh")
	if err != nil {
		return nil, fmt.Errorf("error parsing ssh flag for bare metal create : %v", err)
	}

	app, err := cmd.Flags().GetInt("app")
	if err != nil {
		return nil, fmt.Errorf("error parsing app flag for bare metal create : %v", err)
	}

	userdata, err := cmd.Flags().GetString("userdata")
	if err != nil {
		return nil, fmt.Errorf("error parsing userdata flag for bare metal create : %v", err)
	}

	notify, err := cmd.Flags().GetString("notify")
	if err != nil {
		return nil, fmt.Errorf("error parsing notify flag for bare metal create : %v", err)
	}

	hostname, err := cmd.Flags().GetString("hostname")
	if err != nil {
		return nil, fmt.Errorf("error parsing hostname flag for bare metal create : %v", err)
	}

	tags, err := cmd.Flags().GetStringSlice("tags")
	if err != nil {
		return nil, fmt.Errorf("error parsing tags flag for bare metal create : %v", err)
	}

	ripv4, err := cmd.Flags().GetString("ripv4")
	if err != nil {
		return nil, fmt.Errorf("error parsing ripv4 flag for bare metal create : %v", err)
	}

	pxe, err := cmd.Flags().GetBool("persistent_pxe")
	if err != nil {
		return nil, fmt.Errorf("error parsing persistent_pxe flag for bare metal create : %v", err)
	}

	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return nil, fmt.Errorf("error parsing image flag for bare metal create : %v", err)
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

	return options, nil
}

func parseUpdateFlags(cmd *cobra.Command) (*govultr.BareMetalUpdate, error) { //nolint:unused
	osID, err := cmd.Flags().GetInt("os")
	if err != nil {
		return nil, fmt.Errorf("error parsing os flag for bare metal update : %v", err)
	}

	label, err := cmd.Flags().GetString("label")
	if err != nil {
		return nil, fmt.Errorf("error parsing label flag for bare metal update : %v", err)
	}

	app, err := cmd.Flags().GetInt("app")
	if err != nil {
		return nil, fmt.Errorf("error parsing app flag for bare metal update : %v", err)
	}

	userdata, err := cmd.Flags().GetString("userdata")
	if err != nil {
		return nil, fmt.Errorf("error parsing userdata flag for bare metal update : %v", err)
	}

	tags, err := cmd.Flags().GetStringSlice("tags")
	if err != nil {
		return nil, fmt.Errorf("error parsing tags flag for bare metal update : %v", err)
	}

	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return nil, fmt.Errorf("error parsing image flag for bare metal update : %v", err)
	}

	options := &govultr.BareMetalUpdate{
		AppID:   app,
		OsID:    osID,
		ImageID: image,
		Label:   label,
		Tags:    tags,
	}
	if userdata != "" {
		options.UserData = base64.StdEncoding.EncodeToString([]byte(userdata))
	}

	return options, nil
}
