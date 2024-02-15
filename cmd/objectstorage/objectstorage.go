// Package objectstorage provides the object storage commands for the CLI
package objectstorage

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

// NewCmdObjectStorage provides the CLI command for object storage functions
func NewCmdObjectStorage(base *cli.Base) *cobra.Command { //nolint:gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:   "object-storage",
		Short: "object storage commands",
		Long:  `object-storage is used to interact with object storages`,
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
		Use:   "list",
		Short: "retrieves a list of active object storages",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			oss, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving object storage list : %v", err)
			}

			data := &ObjectStoragesPrinter{ObjectStorages: oss, Meta: meta}
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
		Use:   "get <Object Storage ID>",
		Short: "retrieves a given object storage",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an object storage ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			os, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting object storage info : %v", err)
			}

			data := &ObjectStoragePrinter{ObjectStorage: os}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:   "create",
		Short: "create a new object storage",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterID, errCl := cmd.Flags().GetInt("cluster-id")
			if errCl != nil {
				return fmt.Errorf("error parsing flag 'cluster-id' for object storage create : %v", errCl)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for object storage create : %v", errLa)
			}

			o.ClusterID = clusterID
			o.Label = label

			os, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating object storage : %v", err)
			}

			data := &ObjectStoragePrinter{ObjectStorage: os}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("label", "l", "", "label you want your object storage to have")
	create.Flags().IntP("cluster-id", "i", 0, "ID of the cluster in which to create the object storage")
	if err := create.MarkFlagRequired("cluster-id"); err != nil {
		printer.Error(fmt.Errorf("error marking object storage create 'cluster-id' flag required : %v", err))
		os.Exit(1)
	}

	// Label
	label := &cobra.Command{
		Use:   "label <Object Storage ID>",
		Short: "change the label for object storage subscription",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an object storage ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for object storage label : %v", errLa)
			}

			o.Label = label
			if err := o.update(); err != nil {
				return fmt.Errorf("error updating object storage label : %v", err)
			}

			o.Base.Printer.Display(printer.Info("object storage label has been set"), nil)
			return nil
		},
	}

	label.Flags().StringP("label", "l", "", "label you want your object storage to have")
	if err := label.MarkFlagRequired("label"); err != nil {
		printer.Error(fmt.Errorf("error marking object storage update 'label' flag required: %v", err))
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <Object Storage ID>",
		Short:   "delete specified object storage",
		Aliases: []string{"destroy"},
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an object storage ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("unable to delete object storage : %v", err)
			}

			o.Base.Printer.Display(printer.Info("object storage has been deleted"), nil)
			return nil
		},
	}

	// Regenerate Keys
	regenerateKeys := &cobra.Command{
		Use:   "regenerate-keys <Object Storage ID>",
		Short: "regenerate the S3 API keys for object storage",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an object storage ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := o.regenerateKeys()
			if err != nil {
				return fmt.Errorf("unable to regenerate keys for object storage : %v", err)
			}

			data := &ObjectStorageKeysPrinter{Keys: key}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// List Clusters
	listClusters := &cobra.Command{
		Use:   "list-clusters",
		Short: "retrieve a list of all available object storage clusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			clusters, meta, err := o.listClusters()
			if err != nil {
				return fmt.Errorf("error retrieving object storage cluster list : %v", err)
			}

			data := &ObjectStorageClustersPrinter{Clusters: clusters, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	cmd.AddCommand(
		list,
		get,
		create,
		label,
		del,
		regenerateKeys,
		listClusters,
	)

	return cmd
}

type options struct {
	Base      *cli.Base
	ClusterID int
	Label     string
}

func (o *options) list() ([]govultr.ObjectStorage, *govultr.Meta, error) {
	oss, meta, _, err := o.Base.Client.ObjectStorage.List(o.Base.Context, o.Base.Options)
	return oss, meta, err
}

func (o *options) get() (*govultr.ObjectStorage, error) {
	os, _, err := o.Base.Client.ObjectStorage.Get(o.Base.Context, o.Base.Args[0])
	return os, err
}

func (o *options) create() (*govultr.ObjectStorage, error) {
	os, _, err := o.Base.Client.ObjectStorage.Create(o.Base.Context, o.ClusterID, o.Label)
	return os, err
}

func (o *options) update() error {
	return o.Base.Client.ObjectStorage.Update(o.Base.Context, o.Base.Args[0], o.Label)
}

func (o *options) del() error {
	return o.Base.Client.ObjectStorage.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) listClusters() ([]govultr.ObjectStorageCluster, *govultr.Meta, error) {
	clusters, meta, _, err := o.Base.Client.ObjectStorage.ListCluster(o.Base.Context, o.Base.Options)
	return clusters, meta, err
}

func (o *options) regenerateKeys() (*govultr.S3Keys, error) {
	keys, _, err := o.Base.Client.ObjectStorage.RegenerateKeys(o.Base.Context, o.Base.Args[0])
	return keys, err
}
