// Package cdn provides the CDN related commands to the CLI
package cdn

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
	pullLong    = ``
	pullExample = ``

	pushLong    = ``
	pushExample = ``
)

// NewCmdCDN provides the CLI command for CDN functions
func NewCmdCDN(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:   "cdn",
		Short: "Commands to manage your CDN zones",
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// Pull
	pull := &cobra.Command{
		Use:     "pull",
		Short:   "CDN pull zone commands",
		Long:    pullLong,
		Example: pullExample,
	}

	// Pull List
	pullList := &cobra.Command{
		Use:   "list",
		Short: "List all CDN pull zones",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			pullZones, meta, err := o.pullList()
			if err != nil {
				return fmt.Errorf("error retrieving cdn pull zone list : %v", err)
			}

			data := &PullZonesPrinter{PullZones: pullZones, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Pull Get
	pullGet := &cobra.Command{
		Use:   "get <ZONE ID>",
		Short: "Get a CDN pull zone by ID",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pullZone, err := o.pullGet()
			if err != nil {
				return fmt.Errorf("error getting cdn pull zone : %v", err)
			}

			data := &PullZonePrinter{PullZone: pullZone}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Pull Create
	pullCreate := &cobra.Command{
		Use:   "create",
		Short: "Create a CDN pull zone",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for cdn pull zone create : %v", errLa)
			}

			scheme, errSc := cmd.Flags().GetString("scheme")
			if errSc != nil {
				return fmt.Errorf("error parsing flag 'scheme' for cdn pull zone create : %v", errSc)
			}

			domain, errDo := cmd.Flags().GetString("domain")
			if errDo != nil {
				return fmt.Errorf("error parsing flag 'domain' for cdn pull zone create : %v", errDo)
			}

			cors, errCo := cmd.Flags().GetBool("cors")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'cors' for cdn pull zone create : %v", errCo)
			}

			gzip, errGz := cmd.Flags().GetBool("gzip")
			if errGz != nil {
				return fmt.Errorf("error parsing flag 'gzip' for cdn pull zone create : %v", errGz)
			}

			blai, errBa := cmd.Flags().GetBool("block-ai")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'block-ai' for cdn pull zone create : %v", errBa)
			}

			blbb, errBb := cmd.Flags().GetBool("block-bad-bots")
			if errBb != nil {
				return fmt.Errorf("error parsing flag 'block-bad-bots' for cdn pull zone create : %v", errBb)
			}

			o.ZoneReq = &govultr.CDNZoneReq{
				Label:        label,
				OriginScheme: scheme,
				OriginDomain: domain,
				CORS:         cors,
				GZIP:         gzip,
				BlockAI:      blai,
				BlockBadBots: blbb,
			}

			pullZone, err := o.pullCreate()
			if err != nil {
				return fmt.Errorf("error creating cdn pull zone : %v", err)
			}

			data := &PullZonePrinter{PullZone: pullZone}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	pullCreate.Flags().StringP("label", "l", "", "label to use for the pull zone")
	if err := pullCreate.MarkFlagRequired("label"); err != nil {
		printer.Error(fmt.Errorf("error marking cdn pull zone create 'label' flag required : %v", err))
		os.Exit(1)
	}

	pullCreate.Flags().StringP("scheme", "c", "", "the URI scheme of the origin domain. Must be 'http' or 'https'")
	if err := pullCreate.MarkFlagRequired("scheme"); err != nil {
		printer.Error(fmt.Errorf("error marking cdn pull zone create 'scheme' flag required : %v", err))
		os.Exit(1)
	}

	pullCreate.Flags().StringP("domain", "d", "", "the domain name from which the cdn content will be pulled")
	if err := pullCreate.MarkFlagRequired("domain"); err != nil {
		printer.Error(fmt.Errorf("error marking cdn pull zone create 'domain' flag required : %v", err))
		os.Exit(1)
	}

	pullCreate.Flags().BoolP("cors", "r", false, "enable cross-origin resource sharing")
	pullCreate.Flags().BoolP("gzip", "g", false, "enable gzip compression")
	pullCreate.Flags().BoolP("block-ai", "i", false, "block ai bots")
	pullCreate.Flags().BoolP("block-bad-bots", "t", false, "block potentially malicious bots")

	// Pull Update
	pullUpdate := &cobra.Command{
		Use:   "update",
		Short: "Update a CDN pull zone",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for cdn pull zone update : %v", errLa)
			}

			cors, errCo := cmd.Flags().GetBool("cors")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'cors' for cdn pull zone update : %v", errCo)
			}

			gzip, errGz := cmd.Flags().GetBool("gzip")
			if errGz != nil {
				return fmt.Errorf("error parsing flag 'gzip' for cdn pull zone update : %v", errGz)
			}

			blai, errBa := cmd.Flags().GetBool("block-ai")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'block-ai' for cdn pull zone update : %v", errBa)
			}

			blbb, errBb := cmd.Flags().GetBool("block-bad-bots")
			if errBb != nil {
				return fmt.Errorf("error parsing flag 'block-bad-bots' for cdn pull zone update : %v", errBb)
			}

			o.ZoneReq = &govultr.CDNZoneReq{}

			if cmd.Flags().Changed("label") {
				o.ZoneReq.Label = label
			}

			if cmd.Flags().Changed("cors") {
				o.ZoneReq.CORS = cors
			}

			if cmd.Flags().Changed("gzip") {
				o.ZoneReq.GZIP = gzip
			}

			if cmd.Flags().Changed("block-ai") {
				o.ZoneReq.BlockAI = blai
			}

			if cmd.Flags().Changed("block-bad-bots") {
				o.ZoneReq.BlockBadBots = blbb
			}

			pullZone, err := o.pullUpdate()
			if err != nil {
				return fmt.Errorf("error updating cdn pull zone : %v", err)
			}

			data := &PullZonePrinter{PullZone: pullZone}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	pullUpdate.Flags().StringP("label", "l", "", "label to use for the pull zone")
	pullUpdate.Flags().BoolP("cors", "r", false, "enable cross-origin resource sharing")
	pullUpdate.Flags().BoolP("gzip", "g", false, "enable gzip compression")
	pullUpdate.Flags().BoolP("block-ai", "i", false, "block ai bots")
	pullUpdate.Flags().BoolP("block-bad-bots", "t", false, "block potentially malicious bots")

	// Pull Purge
	pullPurge := &cobra.Command{
		Use:   "purge <ZONE ID>",
		Short: "Purge a CDN pull zone cache",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.pullPurge(); err != nil {
				return fmt.Errorf("error purging cdn pull zone : %v", err)
			}

			o.Base.Printer.Display(printer.Info("CDN pull zone has been purged"), nil)
			return nil
		},
	}

	// Pull Delete
	pullDel := &cobra.Command{
		Use:     "delete <ZONE ID>",
		Short:   "Delete a CDN pull zone",
		Aliases: []string{"destroy"},
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.pullDel(); err != nil {
				return fmt.Errorf("error deleting cdn pull zone : %v", err)
			}

			o.Base.Printer.Display(printer.Info("CDN pull zone has been deleted"), nil)
			return nil
		},
	}

	pull.AddCommand(
		pullList,
		pullGet,
		pullCreate,
		pullUpdate,
		pullPurge,
		pullDel,
	)

	// Push
	push := &cobra.Command{
		Use:     "push",
		Short:   "CDN push zone commands",
		Long:    pushLong,
		Example: pushExample,
	}

	// Push List
	pushList := &cobra.Command{
		Use:   "list",
		Short: "List all CDN push zones",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			pushZones, meta, err := o.pushList()
			if err != nil {
				return fmt.Errorf("error retrieving cdn push zone list : %v", err)
			}

			data := &PushZonesPrinter{PushZones: pushZones, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Push Get
	pushGet := &cobra.Command{
		Use:   "get <ZONE ID>",
		Short: "Get a CDN push zone by ID",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pushZone, err := o.pushGet()
			if err != nil {
				return fmt.Errorf("error getting cdn push zone : %v", err)
			}

			data := &PushZonePrinter{PushZone: pushZone}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Push Create
	pushCreate := &cobra.Command{
		Use:   "create",
		Short: "Create a CDN push zone",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for cdn push zone create : %v", errLa)
			}

			cors, errCo := cmd.Flags().GetBool("cors")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'cors' for cdn push zone create : %v", errCo)
			}

			gzip, errGz := cmd.Flags().GetBool("gzip")
			if errGz != nil {
				return fmt.Errorf("error parsing flag 'gzip' for cdn push zone create : %v", errGz)
			}

			blai, errBa := cmd.Flags().GetBool("block-ai")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'block-ai' for cdn push zone create : %v", errBa)
			}

			blbb, errBb := cmd.Flags().GetBool("block-bad-bots")
			if errBb != nil {
				return fmt.Errorf("error parsing flag 'block-bad-bots' for cdn push zone create : %v", errBb)
			}

			o.ZoneReq = &govultr.CDNZoneReq{
				Label:        label,
				CORS:         cors,
				GZIP:         gzip,
				BlockAI:      blai,
				BlockBadBots: blbb,
			}

			pushZone, err := o.pushCreate()
			if err != nil {
				return fmt.Errorf("error creating cdn push zone : %v", err)
			}

			data := &PushZonePrinter{PushZone: pushZone}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	pushCreate.Flags().StringP("label", "l", "", "label to use for the push zone")
	if err := pushCreate.MarkFlagRequired("label"); err != nil {
		printer.Error(fmt.Errorf("error marking cdn push zone create 'label' flag required : %v", err))
		os.Exit(1)
	}

	pushCreate.Flags().BoolP("cors", "r", false, "enable cross-origin resource sharing")
	pushCreate.Flags().BoolP("gzip", "g", false, "enable gzip compression")
	pushCreate.Flags().BoolP("block-ai", "i", false, "block ai bots")
	pushCreate.Flags().BoolP("block-bad-bots", "t", false, "block potentially malicious bots")

	// Push Update
	pushUpdate := &cobra.Command{
		Use:   "update",
		Short: "Update a CDN push zone",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for cdn push zone update : %v", errLa)
			}

			cors, errCo := cmd.Flags().GetBool("cors")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'cors' for cdn push zone update : %v", errCo)
			}

			gzip, errGz := cmd.Flags().GetBool("gzip")
			if errGz != nil {
				return fmt.Errorf("error parsing flag 'gzip' for cdn push zone update : %v", errGz)
			}

			blai, errBa := cmd.Flags().GetBool("block-ai")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'block-ai' for cdn push zone update : %v", errBa)
			}

			blbb, errBb := cmd.Flags().GetBool("block-bad-bots")
			if errBb != nil {
				return fmt.Errorf("error parsing flag 'block-bad-bots' for cdn push zone update : %v", errBb)
			}

			regions, errRg := cmd.Flags().GetStringSlice("regions")
			if errRg != nil {
				return fmt.Errorf("error parsing flag 'regions' for cdn push zone update : %v", errRg)
			}

			o.ZoneReq = &govultr.CDNZoneReq{}

			if cmd.Flags().Changed("label") {
				o.ZoneReq.Label = label
			}

			if cmd.Flags().Changed("cors") {
				o.ZoneReq.CORS = cors
			}

			if cmd.Flags().Changed("gzip") {
				o.ZoneReq.GZIP = gzip
			}

			if cmd.Flags().Changed("block-ai") {
				o.ZoneReq.BlockAI = blai
			}

			if cmd.Flags().Changed("block-bad-bots") {
				o.ZoneReq.BlockBadBots = blbb
			}

			if cmd.Flags().Changed("regions") {
				o.ZoneReq.Regions = regions
			}

			pushZone, err := o.pushUpdate()
			if err != nil {
				return fmt.Errorf("error updating cdn push zone : %v", err)
			}

			data := &PushZonePrinter{PushZone: pushZone}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	pushUpdate.Flags().StringP("label", "l", "", "label to use for the push zone")
	pushUpdate.Flags().BoolP("cors", "r", false, "enable cross-origin resource sharing")
	pushUpdate.Flags().BoolP("gzip", "g", false, "enable gzip compression")
	pushUpdate.Flags().BoolP("block-ai", "i", false, "block ai bots")
	pushUpdate.Flags().BoolP("block-bad-bots", "t", false, "block potentially malicious bots")
	pushUpdate.Flags().StringSliceP("regions", "n", nil, "a comma separated list of region IDs")

	// Push Delete
	pushDel := &cobra.Command{
		Use:     "delete <ZONE ID>",
		Short:   "Delete a CDN push zone",
		Aliases: []string{"destroy"},
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.pushDel(); err != nil {
				return fmt.Errorf("error deleting cdn push zone : %v", err)
			}

			o.Base.Printer.Display(printer.Info("CDN push zone has been deleted"), nil)
			return nil
		},
	}

	// Push Files List
	pushFilesList := &cobra.Command{
		Use:   "list-files <ZONE ID>",
		Short: "List all files in a CDN push zone",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fileData, err := o.pushFileList()
			if err != nil {
				return fmt.Errorf("error listing cdn push zone files : %v", err)
			}

			data := &PushZoneFilesPrinter{FileData: fileData}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Push Files Get
	pushFileGet := &cobra.Command{
		Use:   "get-file <ZONE ID> <FILE NAME>",
		Short: "Retrieve a file in a CDN push zone",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a zone ID and file name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := o.pushFileGet()
			if err != nil {
				return fmt.Errorf("error retrieving the cdn push zone file : %v", err)
			}

			data := &PushZoneFilePrinter{File: file}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Push File Delete
	pushFileDel := &cobra.Command{
		Use:   "delete-file <ZONE ID> <FILE NAME>",
		Short: "Delete a file in a CDN push zone",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a zone ID and a file name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.pushFileDelete(); err != nil {
				return fmt.Errorf("error deleting cdn push zone file : %v", err)
			}

			o.Base.Printer.Display(printer.Info("CDN push zone file has been deleted"), nil)
			return nil
		},
	}

	// Push Endpoint Create
	pushEndpointCreate := &cobra.Command{
		Use:   "create-endpoint <ZONE ID>",
		Short: "Create a file upload endpoint in a CDN push zone",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a zone ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for cdn pull zone endpoint create : %v", errNa)
			}

			size, errSi := cmd.Flags().GetInt("size")
			if errSi != nil {
				return fmt.Errorf("error parsing flag 'size' for cdn pull zone endpoint create : %v", errSi)
			}

			o.EndpointReq = &govultr.CDNZoneEndpointReq{
				Name: name,
				Size: size,
			}

			endpoint, err := o.pushFileEndpointCreate()
			if err != nil {
				return fmt.Errorf("error creating a cdn push zone file endpoint : %v", err)
			}

			data := &PushZoneEndpointPrinter{Endpoint: endpoint}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	pushEndpointCreate.Flags().StringP("name", "n", "", "the name of the file to be uploaded, including the extension")
	if err := pushEndpointCreate.MarkFlagRequired("name"); err != nil {
		printer.Error(fmt.Errorf("error marking cdn push zone endpoint create 'name' flag required : %v", err))
		os.Exit(1)
	}

	pushEndpointCreate.Flags().IntP("size", "s", 0, "the size of the file to be uploaded (must match the file)")
	if err := pushEndpointCreate.MarkFlagRequired("size"); err != nil {
		printer.Error(fmt.Errorf("error marking cdn push zone endpoint create 'size' flag required : %v", err))
		os.Exit(1)
	}

	push.AddCommand(
		pushList,
		pushGet,
		pushCreate,
		pushUpdate,
		pushDel,
		pushFilesList,
		pushFileGet,
		pushFileDel,
		pushEndpointCreate,
	)

	cmd.AddCommand(
		pull,
		push,
	)

	return cmd
}

type options struct {
	Base        *cli.Base
	ZoneReq     *govultr.CDNZoneReq
	EndpointReq *govultr.CDNZoneEndpointReq
}

func (o *options) pullList() ([]govultr.CDNZone, *govultr.Meta, error) {
	zones, meta, _, err := o.Base.Client.CDN.ListPullZones(o.Base.Context)
	return zones, meta, err
}

func (o *options) pullGet() (*govultr.CDNZone, error) {
	zone, _, err := o.Base.Client.CDN.GetPullZone(o.Base.Context, o.Base.Args[0])
	return zone, err
}

func (o *options) pullCreate() (*govultr.CDNZone, error) {
	zone, _, err := o.Base.Client.CDN.CreatePullZone(o.Base.Context, o.ZoneReq)
	return zone, err
}

func (o *options) pullUpdate() (*govultr.CDNZone, error) {
	zone, _, err := o.Base.Client.CDN.UpdatePullZone(o.Base.Context, o.Base.Args[0], o.ZoneReq)
	return zone, err
}

func (o *options) pullPurge() error {
	return o.Base.Client.CDN.PurgePullZone(o.Base.Context, o.Base.Args[0])
}

func (o *options) pullDel() error {
	return o.Base.Client.CDN.DeletePullZone(o.Base.Context, o.Base.Args[0])
}

func (o *options) pushList() ([]govultr.CDNZone, *govultr.Meta, error) {
	zones, meta, _, err := o.Base.Client.CDN.ListPushZones(o.Base.Context)
	return zones, meta, err
}

func (o *options) pushGet() (*govultr.CDNZone, error) {
	zone, _, err := o.Base.Client.CDN.GetPushZone(o.Base.Context, o.Base.Args[0])
	return zone, err
}

func (o *options) pushCreate() (*govultr.CDNZone, error) {
	zone, _, err := o.Base.Client.CDN.CreatePushZone(o.Base.Context, o.ZoneReq)
	return zone, err
}

func (o *options) pushUpdate() (*govultr.CDNZone, error) {
	zone, _, err := o.Base.Client.CDN.UpdatePushZone(o.Base.Context, o.Base.Args[0], o.ZoneReq)
	return zone, err
}

func (o *options) pushDel() error {
	return o.Base.Client.CDN.DeletePushZone(o.Base.Context, o.Base.Args[0])
}

func (o *options) pushFileList() (*govultr.CDNZoneFileData, error) {
	file, _, err := o.Base.Client.CDN.ListPushZoneFiles(o.Base.Context, o.Base.Args[0])
	return file, err
}

func (o *options) pushFileGet() (*govultr.CDNZoneFile, error) {
	file, _, err := o.Base.Client.CDN.GetPushZoneFile(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return file, err
}

func (o *options) pushFileDelete() error {
	return o.Base.Client.CDN.DeletePushZoneFile(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) pushFileEndpointCreate() (*govultr.CDNZoneEndpoint, error) {
	endpoint, _, err := o.Base.Client.CDN.CreatePushZoneFileEndpoint(o.Base.Context, o.Base.Args[0], o.EndpointReq)
	return endpoint, err
}
