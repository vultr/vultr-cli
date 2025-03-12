// Package objectstorage provides the object storage commands for the CLI
package objectstorage

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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
		Short: "Commands to manage object storage",
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
		Short: "Retrieve all active object storages",
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
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Get
	get := &cobra.Command{
		Use:   "get <Object Storage ID>",
		Short: "Retrieve a given object storage",
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
		Short: "Create a new object storage",
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

			tierID, errTID := cmd.Flags().GetInt("tier-id")
			if errTID != nil {
				return fmt.Errorf("error parsing flag 'tier-id' for object storage create : %v", errTID)
			}

			o.ClusterID = clusterID
			o.Label = label
			o.TierID = tierID

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
	create.Flags().IntP("tier-id", "t", 1, "Tier ID used to create the object storage tiers")
	if err := create.MarkFlagRequired("cluster-id"); err != nil {
		printer.Error(fmt.Errorf("error marking object storage create 'cluster-id' flag required : %v", err))
		os.Exit(1)
	}

	// Label
	label := &cobra.Command{
		Use:   "label <Object Storage ID>",
		Short: "Change the label for object storage",
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
		Short:   "Delete an object storage",
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
		Short: "Regenerate the S3 API keys for an object storage",
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
		Short: "Retrieve a list of all available object storage clusters",
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

	// List Cluster Tiers
	listClusterTiers := &cobra.Command{
		Use:   "list-cluster-tiers [clusterID]",
		Short: "Retrieve a list of all available object storage tiers on a specific cluster",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a Cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid clusterID: %v", err)
			}

			o.Base.Options = utils.GetPaging(cmd)

			clusterTiers, err := o.listClusterTiers(clusterID)
			if err != nil {
				return fmt.Errorf("error retrieving object storage cluster tier list: %v", err)
			}

			data := &ObjectStorageClusterTiersPrinter{ClusterTiers: clusterTiers}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}
	// List Tiers
	listTiers := &cobra.Command{
		Use:   "list-tiers",
		Short: "Retrieve a list of all available object storage tiers",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			tiers, err := o.listTiers()
			if err != nil {
				return fmt.Errorf("error retrieving object storage tier list : %v", err)
			}

			data := &ObjectStorageTiersPrinter{Tiers: tiers}
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
		listClusterTiers,
		listTiers,
	)

	return cmd
}

type options struct {
	Base      *cli.Base
	ClusterID int
	TierID    int
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
	OSreq := &govultr.ObjectStorageReq{
		ClusterID: o.ClusterID,
		TierID:    o.TierID,
		Label:     o.Label,
	}
	os, _, err := o.Base.Client.ObjectStorage.Create(o.Base.Context, OSreq)
	return os, err
}

func (o *options) update() error {
	OSreq := &govultr.ObjectStorageReq{
		Label: o.Label,
	}

	return o.Base.Client.ObjectStorage.Update(o.Base.Context, o.Base.Args[0], OSreq)
}

func (o *options) del() error {
	return o.Base.Client.ObjectStorage.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) listClusters() ([]govultr.ObjectStorageCluster, *govultr.Meta, error) {
	clusters, meta, _, err := o.Base.Client.ObjectStorage.ListCluster(o.Base.Context, o.Base.Options)
	return clusters, meta, err
}

func (o *options) listTiers() ([]govultr.ObjectStorageTier, error) {
	tiers, _, err := o.Base.Client.ObjectStorage.ListTiers(o.Base.Context)
	return tiers, err
}

func (o *options) listClusterTiers(clusterID int) ([]govultr.ObjectStorageTier, error) {
	clusterTiers, _, err := o.Base.Client.ObjectStorage.ListClusterTiers(o.Base.Context, clusterID)
	return clusterTiers, err
}

func (o *options) regenerateKeys() (*govultr.S3Keys, error) {
	keys, _, err := o.Base.Client.ObjectStorage.RegenerateKeys(o.Base.Context, o.Base.Args[0])
	return keys, err
}
