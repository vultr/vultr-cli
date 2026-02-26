// Package instance provides the command for the CLI to control instances
package instance

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/ip"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/userdata"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `Get commands available to instance`
	example = `
	# Full example
	vultr-cli instance
	`
	listLong      = ``
	listExample   = ``
	getLong       = ``
	getExample    = ``
	createLong    = `Create a new instance with specified plan, region and os (from image, snapshot, app or ISO)`
	createExample = `
	# Full example
	vultr-cli instance create --region="ewr" --plan="vc2-2c-4gb" --os=1743

	You must pass one of these in addition to the required --region and --plan flags:
		--os
		--snapshot
		--iso
		--app
		--image

	# Shortened example with aliases
	vultr-cli instance c -r="ewr" -p="vc2-2c-4gb" -o=1743

	# Full example with attached VPCs
	vultr-cli instance create --region="ewr" --plan="vc2-2c-4gb" --os=1743 \
		--vpc-ids="08422775-5be0-4371-afba-64b03f9ad22d,13a45caa-9c06-4b5d-8f76-f5281ab172b7"

	# Full example with assigned ssh keys
	vultr-cli instance create --region="ewr" --plan="vc2-2c-4gb" --os=1743 \
		--ssh-keys="a14b6539-5583-41e8-a035-c07a76897f2b,be624232-56c7-4d5c-bf87-9bdaae7a1fbd"

	# Block devices options
	The --block-devices option allows you to pass in options for any number of
	block storage devices when creating an instance with a VX1 plan. The options
	are passed in a delimited string.  Different block devices are delimited by a
	slash (/). The options for each device are delimited by a comma (,) and each
	option is defined by colon (:).

	For example:

	Local Only
	--block-devices="block-id:local,bootable:true"

	Local Boot + New Block
	--block-devices="block-id:local,bootable:true/disk-size:50,label:new-block-label"

	New Bootable Block
	--block-devices="disk-size:50,label:new-bootable-block,bootable:true"

	Existing Bootable Block
	--block-devices="block-id:BLOCK_DEVICE_ID,bootable:true"

	Existing Bootable Block + Local NVMe
	--block-devices="block-id:local/block-id:BLOCK_DEVICE_ID,bootable:true"
	`
	deleteLong    = ``
	deleteExample = ``
	tagsLong      = `Modify the tags of the specified instance`
	tagsExample   = `
	# Full example
	vultr-cli instance tags <instanceID> --tags="example-tag-1,example-tag-2"

	# Shortened example with aliases
	vultr-cli instance tags <instanceID> -t="example-tag-1,example-tag-2"
	`

	userDataSetLong    = ``
	userDataSetExample = ``
	userDataGetLong    = ``
	userDataGetExample = ``
	vpcAttachLong      = `Attaches an existing VPC to the specified instance`
	vpcAttachExample   = `
	# Full example
	vultr-cli instance vpc attach <instanceID> --vpc-id="2126b7d9-5e2a-491e-8840-838aa6b5f294"
	`
	vpcDetachLong    = `Detaches an existing VPC from the specified instance`
	vpcDetachExample = `
	# Full example
	vultr-cli instance vpc detach <instanceID> --vpc-id="2126b7d9-5e2a-491e-8840-838aa6b5f294"
	`

	vpc2AttachLong    = `Attaches an existing VPC 2.0 network to the specified instance`
	vpc2AttachExample = `
	# Full example
	vultr-cli instance vpc2 attach <instanceID> --vpc-id="2126b7d9-5e2a-491e-8840-838aa6b5f294"
	`
	vpc2DetachLong    = `Detaches an existing VPC 2.0 network from the specified instance`
	vpc2DetachExample = `
	# Full example
	vultr-cli instance vpc2 detach <instanceID> --vpc-id="2126b7d9-5e2a-491e-8840-838aa6b5f294"
	`
)

// NewCmdInstance ...
func NewCmdInstance(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "instance",
		Short:   "Commands to interact with instances",
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
		Short:   "List all instances",
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			instances, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error getting instance list : %v", err)
			}

			data := &InstancesPrinter{Instances: instances, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
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
		Use:     "get <Instance ID>",
		Short:   "Get info on an instance",
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			instance, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting instance : %v", err)
			}

			data := &InstancePrinter{Instance: instance}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create an instance",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for instance create : %v", errRe)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing flag 'plan' for instance create : %v", errPl)
			}

			osID, errOs := cmd.Flags().GetInt("os")
			if errOs != nil {
				return fmt.Errorf("error parsing flag 'os' for instance create : %v", errOs)
			}

			ipxe, errIP := cmd.Flags().GetString("ipxe")
			if errIP != nil {
				return fmt.Errorf("error parsing flag 'ipxe' for instance create : %v", errIP)
			}

			iso, errIs := cmd.Flags().GetString("iso")
			if errIs != nil {
				return fmt.Errorf("error parsing flag 'iso' for instance create : %v", errIs)
			}

			snapshot, errSn := cmd.Flags().GetString("snapshot")
			if errSn != nil {
				return fmt.Errorf("error parsing flag 'snapshot' for instance create : %v", errSn)
			}

			script, errSc := cmd.Flags().GetString("script-id")
			if errSc != nil {
				return fmt.Errorf("error parsing flag 'script' for instance create : %v", errSc)
			}

			ipv6, errIv := cmd.Flags().GetBool("ipv6")
			if errIv != nil {
				return fmt.Errorf("error parsing flag 'ipv6' for instance create : %v", errIv)
			}

			vpcEnable, errVp := cmd.Flags().GetBool("vpc-enable")
			if errVp != nil {
				return fmt.Errorf("error parsing flag 'vpc-enable' for instance create : %v", errVp)
			}

			vpcAttach, errVp := cmd.Flags().GetStringSlice("vpc-ids")
			if errVp != nil {
				return fmt.Errorf("error parsing flag 'vpc-ids' for instance create : %v", errVp)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for instance create : %v", errLa)
			}

			ssh, errSs := cmd.Flags().GetStringSlice("ssh-keys")
			if errSs != nil {
				return fmt.Errorf("error parsing flag 'ssh-keys' for instance create : %v", errSs)
			}

			backup, errBa := cmd.Flags().GetBool("auto-backup")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'auto-backup' for instance create : %v", errBa)
			}

			app, errAp := cmd.Flags().GetInt("app")
			if errAp != nil {
				return fmt.Errorf("error parsing flag 'app' for instance create : %v", errAp)
			}

			image, errIm := cmd.Flags().GetString("image")
			if errIm != nil {
				return fmt.Errorf("error parsing flag 'image' for instance create : %v", errIm)
			}

			userData, errUs := cmd.Flags().GetString("userdata")
			if errUs != nil {
				return fmt.Errorf("error parsing flag 'userData' for instance create : %v", errUs)
			}

			notify, errNo := cmd.Flags().GetBool("notify")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'notify' for instance create : %v", errNo)
			}

			ddos, errDd := cmd.Flags().GetBool("ddos")
			if errDd != nil {
				return fmt.Errorf("error parsing flag 'ddos' for instance create : %v", errDd)
			}

			ipv4, errIi := cmd.Flags().GetString("reserved-ipv4")
			if errIi != nil {
				return fmt.Errorf("error parsing flag 'reserved-ipv4' for instance create : %v", errIi)
			}

			host, errHo := cmd.Flags().GetString("host")
			if errHo != nil {
				return fmt.Errorf("error parsing flag 'host' for instance create : %v", errHo)
			}

			tags, errTa := cmd.Flags().GetStringSlice("tags")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tags' for instance create : %v", errTa)
			}

			fwg, errFw := cmd.Flags().GetString("firewall-group")
			if errFw != nil {
				return fmt.Errorf("error parsing flag 'firewall-group' for instance create : %v", errFw)
			}

			blockDevices, errBD := cmd.Flags().GetStringArray("block-devices")
			if errBD != nil {
				return fmt.Errorf("error parsing flag 'block-devices' for kubernetes cluster create : %v", errBD)
			}

			o.CreateReq = &govultr.InstanceCreateReq{
				Plan:            plan,
				Region:          region,
				OsID:            osID,
				ISOID:           iso,
				SnapshotID:      snapshot,
				AppID:           app,
				ImageID:         image,
				IPXEChainURL:    ipxe,
				ScriptID:        script,
				Label:           label,
				SSHKeys:         ssh,
				UserData:        userData,
				ReservedIPv4:    ipv4,
				Hostname:        host,
				Tags:            tags,
				FirewallGroupID: fwg,
				EnableIPv6:      govultr.BoolToBoolPtr(ipv6),
				DDOSProtection:  govultr.BoolToBoolPtr(ddos),
				ActivationEmail: govultr.BoolToBoolPtr(notify),
				Backups:         "disabled",
				EnableVPC:       govultr.BoolToBoolPtr(vpcEnable),
				AttachVPC:       vpcAttach,
			}

			if backup {
				o.CreateReq.Backups = "enabled"
			}

			if userData != "" {
				o.CreateReq.UserData = base64.StdEncoding.EncodeToString([]byte(userData))
			}

			if len(blockDevices) > 0 {
				bds, errFb := formatBlockDevices(blockDevices)
				if errFb != nil {
					return fmt.Errorf("error in block devices formating : %v", errFb)
				}

				o.CreateReq.BlockDevices = bds
			}

			instance, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating instance : %v", err)
			}

			data := &InstancePrinter{Instance: instance}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("region", "r", "", "The ID of the region in which to create the instance")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking instance create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("plan", "p", "", "The plan ID with which to create the instance")
	if err := create.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking instance create 'plan' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().IntP("os", "", 0, "os id you wish the instance to have")
	create.Flags().StringP("iso", "", "", "iso ID you want to create the instance with")
	create.Flags().StringP("snapshot", "", "", "snapshot ID you want to create the instance with")
	create.Flags().IntP("app", "a", 0, "application ID you want this instance to have")
	create.Flags().StringP("image", "", "", "image ID of the application that will be installed on the server.")
	create.MarkFlagsMutuallyExclusive("os", "iso", "snapshot", "app", "image")
	create.MarkFlagsOneRequired("os", "iso", "snapshot", "app", "image")

	create.Flags().StringP(
		"ipxe",
		"",
		"",
		"if you've selected the 'custom' operating system, this can be set to chainload the specified URL on bootup",
	)
	create.Flags().StringP("script-id", "", "", "script id of the startup script")
	create.Flags().BoolP("ipv6", "", false, "enable ipv6 | true or false")
	create.Flags().BoolP("vpc-enable", "", false, "enable VPC | true or false")
	create.Flags().StringSliceP("vpc-ids", "", []string{}, "VPC IDs you want to assign to the instance")
	create.Flags().StringP("label", "l", "", "label you want to give this instance")
	create.Flags().StringSliceP("ssh-keys", "s", []string{}, "ssh keys you want to assign to the instance")
	create.Flags().BoolP("auto-backup", "b", false, "enable auto backups | true or false")
	create.Flags().StringP(
		"userdata",
		"u",
		"",
		"plain text userdata you want to give this instance",
	)
	create.Flags().BoolP("notify", "n", false, "notify when instance has been created | true or false")
	create.Flags().BoolP("ddos", "d", false, "enable ddos protection | true or false")
	create.Flags().StringP("reserved-ipv4", "", "", "ID of the floating IP to use as the main IP for this instance")
	create.Flags().StringP("host", "", "", "The hostname to assign to this instance")
	create.Flags().StringSliceP("tags", "", []string{}, "A comma-separated list of tags to assign to this instance")
	create.Flags().StringP("firewall-group", "", "", "The firewall group to assign to this instance")

	create.Flags().StringArray(
		"block-devices",
		[]string{},
		`a comma-separated, key-value pair list of block devices. At least one block is required for VX1 plans.`,
	)

	// Update
	// update := &cobra.Command{}

	// Delete
	del := &cobra.Command{
		Use:     "delete <Instance ID>",
		Short:   "Delete an instance",
		Aliases: []string{"destroy"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance has been deleted"), nil)

			return nil
		},
	}

	// Label
	label := &cobra.Command{
		Use:   "label <Instance ID>",
		Short: "Label an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for instance update : %v", errLa)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				Label: label,
			}

			instance, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating instance label : %v", err)
			}

			data := &InstancePrinter{Instance: instance}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	label.Flags().StringP("label", "l", "", "The label you want to set on an instance")
	if err := label.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking instance label 'label' flag required: %v", err)
		os.Exit(1)
	}

	// Tags
	tags := &cobra.Command{
		Use:     "tags <Instance ID>",
		Short:   "Update tags on an instance",
		Long:    tagsLong,
		Example: tagsExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			tags, errTa := cmd.Flags().GetStringSlice("tags")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tags' for instance update : %v", errTa)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				Tags: tags,
			}

			instance, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating instance tags : %v", err)
			}

			data := &InstancePrinter{Instance: instance}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	tags.Flags().StringSliceP("tags", "t", []string{}, "A comma separated list of tags to apply to the instance")
	if err := tags.MarkFlagRequired("tags"); err != nil {
		fmt.Printf("error marking instance tags 'tags' flag required: %v", err)
		os.Exit(1)
	}

	// User Data
	userData := &cobra.Command{
		Use:   "user-data",
		Short: "Commands to manage user data on an instance",
	}

	// User Data Get
	userDataGet := &cobra.Command{
		Use:     "get <Instance ID>",
		Short:   "Get the user data on an instance",
		Long:    userDataGetLong,
		Example: userDataGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ud, err := o.userData()
			if err != nil {
				return fmt.Errorf("error getting instance user data : %v", err)
			}

			data := &userdata.UserDataPrinter{UserData: *ud}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// User Data Set
	userDataSet := &cobra.Command{
		Use:     "set <Instance ID>",
		Short:   "Update user data on an instance",
		Long:    userDataSetLong,
		Example: userDataSetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			userDataPath, errPa := cmd.Flags().GetString("userdata")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'userdata' for instance update : %v", errPa)
			}

			userDataPath = filepath.Clean(userDataPath)

			rawData, errRe := os.ReadFile(userDataPath)
			if errRe != nil {
				return fmt.Errorf("error reading user-data : %v", errRe)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				UserData: base64.StdEncoding.EncodeToString(rawData),
			}

			_, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating instance user data : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance user data has been updated"), nil)

			return nil
		},
	}

	userDataSet.Flags().StringP("userdata", "d", "/dev/stdin", "The file to read userdata from")

	userData.AddCommand(
		userDataGet,
		userDataSet,
	)

	// Start
	start := &cobra.Command{
		Use:   "start <Instance ID>",
		Short: "Start an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.start(); err != nil {
				return fmt.Errorf("error starting instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance started"), nil)

			return nil
		},
	}

	// Stop
	stop := &cobra.Command{
		Use:   "stop <Instance ID>",
		Short: "Stop an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.stop(); err != nil {
				return fmt.Errorf("error stopping instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance stopped"), nil)

			return nil
		},
	}

	// Restart
	restart := &cobra.Command{
		Use:   "restart <Instance ID>",
		Short: "Restart an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.restart(); err != nil {
				return fmt.Errorf("error restarting instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance restarted"), nil)

			return nil
		},
	}

	// ISO
	iso := &cobra.Command{
		Use:   "iso",
		Short: "Manage ISOs on an instance",
	}

	// ISO Status
	isoStatus := &cobra.Command{
		Use:   "status <Instance ID>",
		Short: "Get ISO status",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			iso, err := o.isoStatus()
			if err != nil {
				return fmt.Errorf("error getting instance iso status : %v", err)
			}

			data := &ISOPrinter{ISO: *iso}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// ISO Attach
	isoAttach := &cobra.Command{
		Use:   "attach <Instance ID>",
		Short: "Attach ISO to an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			iso, errIs := cmd.Flags().GetString("iso-id")
			if errIs != nil {
				return fmt.Errorf("error parsing flag 'iso' for instance iso attach: %v", errIs)
			}

			o.ISOAttachID = iso

			if err := o.isoAttach(); err != nil {
				return fmt.Errorf("error attaching iso to instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("ISO attached to instance"), nil)

			return nil
		},
	}

	isoAttach.Flags().StringP("iso-id", "i", "", "id of the ISO you wish to attach")
	if err := isoAttach.MarkFlagRequired("iso-id"); err != nil {
		fmt.Printf("error marking instance iso attach 'iso-id' flag required: %v", err)
		os.Exit(1)
	}

	iso.AddCommand(
		isoStatus,
		isoAttach,
	)

	// Backup
	backup := &cobra.Command{
		Use:   "backup",
		Short: "List and create backup schedules for an instance",
	}

	// Backup Get
	backupGet := &cobra.Command{
		Use:   "get <Instance ID>",
		Short: "Get the backup schedule for an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bk, err := o.backups()
			if err != nil {
				return fmt.Errorf("error getting instance backups : %v", err)
			}

			data := &BackupPrinter{Backup: *bk}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Backup Create
	backupCreate := &cobra.Command{
		Use:   "create <Instance ID>",
		Short: "Create a backup schedule for an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			crontType, errCr := cmd.Flags().GetString("type")
			if errCr != nil {
				return fmt.Errorf("error parsing flag 'crontType' for instance backup create : %v", errCr)
			}

			hour, errHo := cmd.Flags().GetInt("hour")
			if errHo != nil {
				return fmt.Errorf("error parsing flag 'hour' for instance backup create : %v", errHo)
			}

			dow, errDo := cmd.Flags().GetInt("dow")
			if errDo != nil {
				return fmt.Errorf("error parsing flag 'dow' for instance backup create : %v", errDo)
			}

			dom, errDo := cmd.Flags().GetInt("dom")
			if errDo != nil {
				return fmt.Errorf("error parsing flag 'dom' for instance backup create : %v", errDo)
			}

			o.BackupCreateReq = &govultr.BackupScheduleReq{
				Type: crontType,
				Hour: govultr.IntToIntPtr(hour),
				Dow:  govultr.IntToIntPtr(dow),
				Dom:  dom,
			}

			if err := o.backupCreate(); err != nil {
				return fmt.Errorf("error getting instance backups : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance backup created"), nil)

			return nil
		},
	}

	backupCreate.Flags().StringP(
		"type",
		"t",
		"",
		`type string Backup cron type. Can be one of 'daily', 'weekly', 'monthly', 
'daily_alt_even', or 'daily_alt_odd'.`,
	)
	if err := backupCreate.MarkFlagRequired("type"); err != nil {
		fmt.Printf("error marking instance backup create 'type' flag required: %v", err)
		os.Exit(1)
	}
	backupCreate.Flags().IntP(
		"hour",
		"",
		0,
		"Hour value (0-23). Applicable to crons: 'daily', 'weekly', 'monthly', 'daily_alt_even', 'daily_alt_odd'",
	)
	backupCreate.Flags().IntP("dow", "w", 0, "Day-of-week value (0-6). Applicable to crons: 'weekly'")
	backupCreate.Flags().IntP("dom", "m", 0, "Day-of-month value (1-28). Applicable to crons: 'monthly'")

	backup.AddCommand(
		backupGet,
		backupCreate,
	)

	// Restore
	restore := &cobra.Command{
		Use:   "restore <Instance ID>",
		Short: "Restore instance from backup or snapshot",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			backup, errBa := cmd.Flags().GetString("backup")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'backup' for instance restore : %v", errBa)
			}

			snapshot, errSn := cmd.Flags().GetString("snapshot")
			if errSn != nil {
				return fmt.Errorf("error parsing flag 'snapshot' for instance restore : %v", errSn)
			}

			o.RestoreReq = &govultr.RestoreReq{
				SnapshotID: snapshot,
				BackupID:   backup,
			}

			if err := o.restore(); err != nil {
				return fmt.Errorf("error restoring instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance restored"), nil)

			return nil
		},
	}

	restore.Flags().StringP("backup", "b", "", "id of backup you wish to restore the instance with")
	restore.Flags().StringP("snapshot", "s", "", "id of snapshot you wish to restore the instance with")
	restore.MarkFlagsOneRequired("backup", "snapshot")
	restore.MarkFlagsMutuallyExclusive("backup", "snapshot")

	// Reinstall
	reinstall := &cobra.Command{
		Use:   "reinstall <Instance ID>",
		Short: "Reinstall an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname, errHo := cmd.Flags().GetString("host")
			if errHo != nil {
				return fmt.Errorf("error parsing flag 'host' for instance reinstall : %v", errHo)
			}

			o.ReinstallReq = &govultr.ReinstallReq{}
			if cmd.Flags().Changed("host") {
				o.ReinstallReq.Hostname = hostname
			}

			if err := o.reinstall(); err != nil {
				return fmt.Errorf("error reinstalling instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Instance reinstalled"), nil)

			return nil
		},
	}

	reinstall.Flags().StringP("host", "", "", "The hostname to assign to this instance")

	// Operating System
	operatingSystem := &cobra.Command{
		Use:   "os",
		Short: "Operating system commands for an instance",
	}

	// OS List
	osList := &cobra.Command{
		Use:   "list <Instance ID>",
		Short: "List available operating systems",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrades, err := o.upgrades()
			if err != nil {
				return fmt.Errorf("error getting instance os list : %v", err)
			}

			data := &OSsPrinter{OSs: upgrades.OS}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// OS Change
	osChange := &cobra.Command{
		Use:   "change <Instance ID>",
		Short: "Change operating system",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			osID, errOs := cmd.Flags().GetInt("os")
			if errOs != nil {
				return fmt.Errorf("error parsing flag 'osID' for instance os change : %v", errOs)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				OsID: osID,
			}

			_, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating instance os : %v", err)
			}

			o.Base.Printer.Display(printer.Info("OS change complete"), nil)

			return nil
		},
	}

	osChange.Flags().IntP("os", "", 0, "operating system ID you wish to use")
	if err := osChange.MarkFlagRequired("os"); err != nil {
		fmt.Printf("error marking instance os update 'os' flag required: %v", err)
		os.Exit(1)
	}

	operatingSystem.AddCommand(
		osList,
		osChange,
	)

	// Application
	app := &cobra.Command{
		Use:     "app",
		Aliases: []string{"image"},
		Short:   "Application commands for an instance",
	}

	// App List
	appList := &cobra.Command{
		Use:   "list <Instance ID>",
		Short: "List available applications",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrades, err := o.upgrades()
			if err != nil {
				return fmt.Errorf("error getting instance applications list : %v", err)
			}

			data := &AppsPrinter{Apps: upgrades.Applications}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// App Change
	appChange := &cobra.Command{
		Use:   "change <Instance ID>",
		Short: "Change application",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appID, errAp := cmd.Flags().GetInt("app")
			if errAp != nil {
				return fmt.Errorf("error parsing flag 'app' for instance application change : %v", errAp)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				AppID: appID,
			}

			_, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating instance application : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Application change complete"), nil)

			return nil
		},
	}

	appChange.Flags().IntP("app", "", 0, "Application ID you wish to use")
	if err := appChange.MarkFlagRequired("app"); err != nil {
		fmt.Printf("error marking instance app update 'app' flag required: %v", err)
		os.Exit(1)
	}

	app.AddCommand(
		appList,
		appChange,
	)

	// Plan
	plan := &cobra.Command{
		Use:   "plan",
		Short: "Plan commands for an instance",
	}

	// Plan List
	planList := &cobra.Command{
		Use:   "list <Instance ID>",
		Short: "List available plans",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrades, err := o.upgrades()
			if err != nil {
				return fmt.Errorf("error getting instance applications list : %v", err)
			}

			data := &PlansPrinter{Plans: upgrades.Plans}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Plan Upgrade
	planUpgrade := &cobra.Command{
		Use:   "upgrade <Instance ID>",
		Short: "Upgrade instance plan",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing flag 'plan' for instance plan upgrade : %v", errPl)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				Plan: plan,
			}

			_, err := o.update()
			if err != nil {
				return fmt.Errorf("error upgrading plan on instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Plan upgrade complete"), nil)

			return nil
		},
	}

	planUpgrade.Flags().String("plan", "", "The plan ID you wish to use")
	if err := planUpgrade.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking instance plan upgrade 'plan' flag required: %v", err)
		os.Exit(1)
	}

	plan.AddCommand(
		planList,
		planUpgrade,
	)

	// IPv4
	ipv4 := &cobra.Command{
		Use:   "ipv4",
		Short: "IPv4 instance commands",
	}

	// IPv4 List
	ipv4List := &cobra.Command{
		Use:     "list <Instance ID>",
		Aliases: []string{"v4"},
		Short:   "List IPv4 for an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			v4s, meta, err := o.ipv4s()
			if err != nil {
				return fmt.Errorf("error getting ipv4 list for instance : %v", err)
			}

			data := &ip.IPv4sPrinter{IPv4s: v4s, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	ipv4List.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	ipv4List.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// IPv4 Create
	ipv4Create := &cobra.Command{
		Use:   "create <Instance ID>",
		Short: "Create IPv4 for instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			reboot, errRe := cmd.Flags().GetBool("reboot")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'reboot' for instance ipv4 create: %v", errRe)
			}

			if cmd.Flags().Changed("reboot") {
				o.Reboot = &reboot
			}

			if err := o.ipv4Create(); err != nil {
				return fmt.Errorf("error creating instance ipv4 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("IPv4 has been created"), nil)

			return nil
		},
	}

	ipv4Create.Flags().Bool("reboot", false, "whether to reboot instance after adding ipv4 address")

	// IPv4 Delete
	ipv4Delete := &cobra.Command{
		Use:   "delete <Instance ID> <IPv4 Address>",
		Short: "Delete IPv4 on an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and the IP address")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.ipv4Delete(); err != nil {
				return fmt.Errorf("error deleting ipv4 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("IPv4 has been deleted"), nil)

			return nil
		},
	}

	ipv4.AddCommand(
		ipv4List,
		ipv4Create,
		ipv4Delete,
	)

	// IPv6
	ipv6 := &cobra.Command{
		Use:   "ipv6",
		Short: "Display instance IPv6 info",
	}

	// IPv6 List
	ipv6List := &cobra.Command{
		Use:     "list <Instance ID>",
		Aliases: []string{"v6"},
		Short:   "List IPv6 for an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			v6s, meta, err := o.ipv6s()
			if err != nil {
				return fmt.Errorf("error getting ipv6 list for instance : %v", err)
			}

			data := &ip.IPv6sPrinter{IPv6s: v6s, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	ipv6List.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	ipv6List.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	ipv6.AddCommand(
		ipv6List,
	)

	// Reverse DNS
	reverseDNS := &cobra.Command{
		Use:   "reverse-dns",
		Short: "Commands to handle reverse DNS on an instance",
	}

	rDNSIPv4Default := &cobra.Command{
		Use:   "default-ipv4 <Instance ID> <IPv4 Address>",
		Short: "Set a reverse DNS entry for an IPv4 address of an instance to the original setting",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and IP address")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.reverseDNSDefault(); err != nil {
				return fmt.Errorf("error setting default reverse dns ipv4 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Reverse DNS defaulse IPv4 has been set"), nil)

			return nil
		},
	}

	// Reverse DNS IPv4 Set
	rDNSIPv4Set := &cobra.Command{
		Use:   "set-ipv4 <Instance ID> <IPv4 Address>",
		Short: "Set a reverse DNS IPv4 address entry for an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and an IPv4 address")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, errEn := cmd.Flags().GetString("entry")
			if errEn != nil {
				return fmt.Errorf("error parsing flag 'entry' for instance reverse dns ipv4 set : %v", errEn)
			}

			o.ReverseDNSReq = &govultr.ReverseIP{
				IP:      o.Base.Args[1],
				Reverse: entry,
			}

			if err := o.reverseDNSIPv4Create(); err != nil {
				return fmt.Errorf("error creating reverse dns ipv4 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Reverse DNS IPv4 has been set"), nil)

			return nil
		},
	}

	rDNSIPv4Set.Flags().StringP("entry", "e", "", "reverse dns entry")
	if err := rDNSIPv4Set.MarkFlagRequired("entry"); err != nil {
		fmt.Printf("error marking instance reverse-dns set-ipv4 'entry' flag required: %v", err)
		os.Exit(1)
	}

	// Reverse DNS IPv6 List
	rDNSIPv6List := &cobra.Command{
		Use:   "list-ipv6 <Instance ID>",
		Short: "List the IPv6 reverse DNS entries for an instance",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ips, err := o.reverseDNSIPv6List()
			if err != nil {
				return fmt.Errorf("error retrieving list of reverse dns ipv6 : %v", err)
			}

			data := &ReverseIPsPrinter{ReverseIPs: ips}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Reverse DNS IPv6 Set
	rDNSIPv6Set := &cobra.Command{
		Use:   "set-ipv6 <Instance ID> <IPv6 Address>",
		Short: "Set a reverse DNS IPv6 address entry for an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and an IPv6 address")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, errEn := cmd.Flags().GetString("entry")
			if errEn != nil {
				return fmt.Errorf("error parsing flag 'entry' for instance reverse dns ipv6 set : %v", errEn)
			}

			o.ReverseDNSReq = &govultr.ReverseIP{
				IP:      o.Base.Args[1],
				Reverse: entry,
			}

			if err := o.reverseDNSIPv6Create(); err != nil {
				return fmt.Errorf("error creating reverse dns ipv6 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Reverse DNS IPv6 has been set"), nil)

			return nil
		},
	}

	rDNSIPv6Set.Flags().StringP("entry", "e", "", "reverse dns entry")
	if err := rDNSIPv6Set.MarkFlagRequired("entry"); err != nil {
		fmt.Printf("error marking instance reverse-dns set-ipv6 'entry' flag required: %v", err)
		os.Exit(1)
	}

	// Reverse DNS IPv6 Delete
	rDNSIPv6Delete := &cobra.Command{
		Use:   "delete-ipv6 <Instance ID>, <IPv6 Address>",
		Short: "Remove a reverse DNS IPv6 address entry for an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID and an IPv6 address")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.reverseDNSIPv6Delete(); err != nil {
				return fmt.Errorf("error deleting reverse dns ipv6 : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Reverse DNS IPv6 has been deleted"), nil)

			return nil
		},
	}

	reverseDNS.AddCommand(
		rDNSIPv4Default,
		rDNSIPv4Set,
		rDNSIPv6List,
		rDNSIPv6Delete,
		rDNSIPv6Set,
	)

	// Firewall Group
	firewallGroup := &cobra.Command{
		Use:   "update-firewall-group",
		Short: "Assign a firewall group to instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			fwgID, errID := cmd.Flags().GetString("firewall-group-id")
			if errID != nil {
				return fmt.Errorf(
					"error parsing flag 'firewall-group-id' for instance firewall group assignment : %v",
					errID,
				)
			}

			o.UpdateReq = &govultr.InstanceUpdateReq{
				FirewallGroupID: fwgID,
			}

			if _, err := o.update(); err != nil {
				return fmt.Errorf("error updating fire wall group on instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Firewall group assigned to instance"), nil)

			return nil
		},
	}

	firewallGroup.Flags().StringP(
		"firewall-group-id",
		"f",
		"",
		"firewall group id that you want to assign. 0 Value will unset the firewall-group",
	)
	if err := firewallGroup.MarkFlagRequired("firewall-group-id"); err != nil {
		fmt.Printf("error marking instance firewall group 'firewall-group-id' flag required: %v", err)
		os.Exit(1)
	}

	// VPC
	vpc := &cobra.Command{
		Use:   "vpc",
		Short: "Commands to handle vpcs on an instance",
	}

	// VPC Attach
	vpcAttach := &cobra.Command{
		Use:     "attach <Instance ID>",
		Short:   "Attach a VPC to an instance",
		Long:    vpcAttachLong,
		Example: vpcAttachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.vpcAttach(); err != nil {
				return fmt.Errorf("error attaching vpc to instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("VPC attached to instance"), nil)

			return nil
		},
	}

	// VPC Detach
	vpcDetach := &cobra.Command{
		Use:     "detach <Instance ID>",
		Short:   "Detach a VPC from an instance",
		Long:    vpcDetachLong,
		Example: vpcDetachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.vpcDetach(); err != nil {
				return fmt.Errorf("error detaching vpc from instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("VPC detached from instance"), nil)

			return nil
		},
	}

	vpc.AddCommand(
		vpcAttach,
		vpcDetach,
	)

	// VPC2
	vpc2 := &cobra.Command{
		Use:        "vpc2",
		Short:      "Commands to handle vpc2s on an instance",
		Deprecated: "all vpc2 commands should be migrated to vpc.",
	}

	// VPC List
	vpc2List := &cobra.Command{
		Use:     "list <Instance ID>",
		Aliases: []string{"l"},
		Short:   "List all VPC2 networks attached to an instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			vpc2s, meta, err := o.vpc2s()
			if err != nil {
				return fmt.Errorf("error getting vpc2 list for instance : %v", err)
			}

			data := &VPC2sPrinter{VPC2s: vpc2s, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
		Deprecated: "all vpc2 commands should be migrated to vpc.",
	}

	// VPC2 Attach
	vpc2Attach := &cobra.Command{
		Use:     "attach <Instance ID>, <VPC2 ID>",
		Short:   "Attach a VPC2 to an instance",
		Long:    vpc2AttachLong,
		Example: vpc2AttachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and a VPC ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, errIP := cmd.Flags().GetString("ip-address")
			if errIP != nil {
				return fmt.Errorf("error parsing flag 'ip-address' for vpc2 instance attach : %v", errIP)
			}

			o.VPC2Req = &govultr.AttachVPC2Req{ //nolint:staticcheck
				VPCID:     o.Base.Args[1],
				IPAddress: &ip,
			}

			if err := o.vpc2Attach(); err != nil {
				return fmt.Errorf("error attaching vpc2 to instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("VPC2 attached to instance"), nil)

			return nil
		},
		Deprecated: "all vpc2 commands should be migrated to vpc.",
	}

	vpc2Attach.Flags().StringP(
		"ip-address",
		"i",
		"",
		"the IP address to use for this instance on the attached VPC 2.0 network",
	)

	// VPC2 Detach
	vpc2Detach := &cobra.Command{
		Use:     "detach <Instance ID>",
		Short:   "Detach a VPC2 from an instance",
		Long:    vpc2DetachLong,
		Example: vpc2DetachExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide an instance ID and a VPC2 ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.vpc2Detach(); err != nil {
				return fmt.Errorf("error detaching vpc2 from instance : %v", err)
			}

			o.Base.Printer.Display(printer.Info("VPC2 detached from instance"), nil)

			return nil
		},
		Deprecated: "all vpc2 commands should be migrated to vpc.",
	}

	vpc2.AddCommand(
		vpc2List,
		vpc2Attach,
		vpc2Detach,
	)

	// Bandwidth
	bandwidth := &cobra.Command{
		Use:   "bandwidth <Instance ID>",
		Short: "Get bandwidth usage ",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an instance ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bw, err := o.bandwidth()
			if err != nil {
				return fmt.Errorf("error getting bandwidth details : %v", err)
			}

			data := &BandwidthPrinter{Bandwidth: *bw}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	cmd.AddCommand(
		list,
		get,
		create,
		del,
		label,
		tags,
		userData,
		start,
		stop,
		restart,
		iso,
		backup,
		restore,
		reinstall,
		operatingSystem,
		app,
		plan,
		ipv4,
		ipv6,
		reverseDNS,
		firewallGroup,
		vpc,
		vpc2,
		bandwidth,
	)

	return cmd
}

// formatBlockDevices parses block devices into proper format
func formatBlockDevices(blockDevices []string) ([]govultr.InstanceBlockDevice, error) {
	var formattedList []govultr.InstanceBlockDevice
	bdList := strings.Split(blockDevices[0], "/")

	for _, r := range bdList {
		bdData := strings.Split(r, ",")

		if len(bdData) < 1 || len(bdData) > 4 {
			return nil, fmt.Errorf(
				`unable to format block devices. valid fields are block-id, bootable, disk-size, and label`,
			)
		}

		formattedBDData, errFormat := formatBlockDeviceData(bdData)
		if errFormat != nil {
			return nil, errFormat
		}

		formattedList = append(formattedList, *formattedBDData)
	}

	return formattedList, nil
}

// formatBlockDeviceData loops over the parse strings for a block device and returns the formatted struct
func formatBlockDeviceData(bd []string) (*govultr.InstanceBlockDevice, error) {
	bdData := &govultr.InstanceBlockDevice{}
	for _, f := range bd {
		bdDataKeyVal := strings.Split(f, ":")

		if len(bdDataKeyVal) != 2 {
			return nil, fmt.Errorf("invalid block device format")
		}

		field := bdDataKeyVal[0]
		val := bdDataKeyVal[1]

		switch field {
		case "block-id":
			bdData.BlockID = val
		case "bootable":
			bootable, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for block device bootable: %v", err)
			}
			bdData.Bootable = bootable
		case "disk-size":
			diskSize, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for block device disk size: %v", err)
			}
			bdData.DiskSize = diskSize
		case "label":
			bdData.Label = val
		}
	}

	return bdData, nil
}

type options struct {
	Base            *cli.Base
	CreateReq       *govultr.InstanceCreateReq
	UpdateReq       *govultr.InstanceUpdateReq
	BackupCreateReq *govultr.BackupScheduleReq
	RestoreReq      *govultr.RestoreReq
	ReinstallReq    *govultr.ReinstallReq
	ISOAttachID     string
	Reboot          *bool
	ReverseDNSReq   *govultr.ReverseIP
	VPC2Req         *govultr.AttachVPC2Req //nolint:staticcheck
}

func (o *options) list() ([]govultr.Instance, *govultr.Meta, error) {
	insts, meta, _, err := o.Base.Client.Instance.List(o.Base.Context, o.Base.Options)
	return insts, meta, err
}

func (o *options) get() (*govultr.Instance, error) {
	inst, _, err := o.Base.Client.Instance.Get(o.Base.Context, o.Base.Args[0])
	return inst, err
}

func (o *options) create() (*govultr.Instance, error) {
	inst, _, err := o.Base.Client.Instance.Create(o.Base.Context, o.CreateReq)
	return inst, err
}

func (o *options) update() (*govultr.Instance, error) {
	inst, _, err := o.Base.Client.Instance.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
	return inst, err
}

func (o *options) del() error {
	return o.Base.Client.Instance.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) userData() (*govultr.UserData, error) {
	ud, _, err := o.Base.Client.Instance.GetUserData(o.Base.Context, o.Base.Args[0])
	return ud, err
}

func (o *options) start() error {
	return o.Base.Client.Instance.Start(o.Base.Context, o.Base.Args[0])
}

func (o *options) stop() error {
	return o.Base.Client.Instance.Halt(o.Base.Context, o.Base.Args[0])
}

func (o *options) restart() error {
	return o.Base.Client.Instance.Reboot(o.Base.Context, o.Base.Args[0])
}

func (o *options) backups() (*govultr.BackupSchedule, error) {
	bk, _, err := o.Base.Client.Instance.GetBackupSchedule(o.Base.Context, o.Base.Args[0])
	return bk, err
}

func (o *options) backupCreate() error {
	_, err := o.Base.Client.Instance.SetBackupSchedule(o.Base.Context, o.Base.Args[0], o.BackupCreateReq)
	return err
}

func (o *options) restore() error {
	_, err := o.Base.Client.Instance.Restore(o.Base.Context, o.Base.Args[0], o.RestoreReq)
	return err
}

func (o *options) reinstall() error {
	_, _, err := o.Base.Client.Instance.Reinstall(o.Base.Context, o.Base.Args[0], o.ReinstallReq)
	return err
}

func (o *options) isoStatus() (*govultr.Iso, error) {
	iso, _, err := o.Base.Client.Instance.ISOStatus(o.Base.Context, o.Base.Args[0])
	return iso, err
}

func (o *options) isoAttach() error {
	_, err := o.Base.Client.Instance.AttachISO(o.Base.Context, o.Base.Args[0], o.ISOAttachID)
	return err
}

func (o *options) upgrades() (*govultr.Upgrades, error) {
	oss, _, err := o.Base.Client.Instance.GetUpgrades(o.Base.Context, o.Base.Args[0])
	return oss, err
}

func (o *options) ipv4s() ([]govultr.IPv4, *govultr.Meta, error) {
	ips, meta, _, err := o.Base.Client.Instance.ListIPv4(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return ips, meta, err
}

func (o *options) ipv4Create() error {
	_, _, err := o.Base.Client.Instance.CreateIPv4(o.Base.Context, o.Base.Args[0], o.Reboot)
	return err
}

func (o *options) ipv4Delete() error {
	return o.Base.Client.Instance.DeleteIPv4(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) ipv6s() ([]govultr.IPv6, *govultr.Meta, error) {
	ips, meta, _, err := o.Base.Client.Instance.ListIPv6(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return ips, meta, err
}

func (o *options) reverseDNSDefault() error {
	return o.Base.Client.Instance.DefaultReverseIPv4(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) reverseDNSIPv4Create() error {
	return o.Base.Client.Instance.CreateReverseIPv4(o.Base.Context, o.Base.Args[0], o.ReverseDNSReq)
}

func (o *options) reverseDNSIPv6List() ([]govultr.ReverseIP, error) {
	ips, _, err := o.Base.Client.Instance.ListReverseIPv6(o.Base.Context, o.Base.Args[0])
	return ips, err
}

func (o *options) reverseDNSIPv6Create() error {
	return o.Base.Client.Instance.CreateReverseIPv6(o.Base.Context, o.Base.Args[0], o.ReverseDNSReq)
}

func (o *options) reverseDNSIPv6Delete() error {
	return o.Base.Client.Instance.DeleteReverseIPv6(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) vpcAttach() error {
	return o.Base.Client.Instance.AttachVPC(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) vpcDetach() error {
	return o.Base.Client.Instance.DetachVPC(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) vpc2s() ([]govultr.VPC2Info, *govultr.Meta, error) { //nolint:staticcheck
	vpc2s, meta, _, err := o.Base.Client.Instance.ListVPC2Info(o.Base.Context, o.Base.Args[0], o.Base.Options) //nolint:staticcheck,lll
	return vpc2s, meta, err
}

func (o *options) vpc2Attach() error {
	return o.Base.Client.Instance.AttachVPC2(o.Base.Context, o.Base.Args[0], o.VPC2Req) //nolint:staticcheck
}

func (o *options) vpc2Detach() error {
	return o.Base.Client.Instance.DetachVPC2(o.Base.Context, o.Base.Args[0], o.Base.Args[1]) //nolint:staticcheck
}

func (o *options) bandwidth() (*govultr.Bandwidth, error) {
	bw, _, err := o.Base.Client.Instance.GetBandwidth(o.Base.Context, o.Base.Args[0])
	return bw, err
}
