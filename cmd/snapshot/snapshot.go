// Package snapshot provides functionality for the CLI to control snapshots
package snapshot

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

// NewCmdSnapshot provides the CLI command for snapshot functions
func NewCmdSnapshot(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "snapshot",
		Short:   "Commands to interact with snapshots",
		Aliases: []string{"sn"},
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
		Short: "List all snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			snaps, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving snapshot list : %v", err)
			}

			data := &SnapshotsPrinter{Snapshots: snaps, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Get
	get := &cobra.Command{
		Use:   "get <Snapshot ID>",
		Short: "Get a snapshot",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a snapshot ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshot, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving snapshot : %v", err)
			}

			data := &SnapshotPrinter{Snapshot: snapshot}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:   "create",
		Short: "Create a snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, errID := cmd.Flags().GetString("id")
			if errID != nil {
				return fmt.Errorf("error parsing flag 'id' for create : %v", errID)
			}

			desc, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'description' for create : %v", errDe)
			}

			o.Req = &govultr.SnapshotReq{
				InstanceID:  id,
				Description: desc,
			}

			snapshot, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating snapshot : %v", err)
			}

			data := &SnapshotPrinter{Snapshot: snapshot}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("id", "i", "", "ID of the virtual machine to create a snapshot from.")
	if err := create.MarkFlagRequired("id"); err != nil {
		fmt.Printf("error marking snapshot create 'id' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("description", "d", "", "(optional) Description of snapshot contents")

	// Create URL
	createURL := &cobra.Command{
		Use:   "create-url",
		Short: "Create a snapshot from a URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			url, errUR := cmd.Flags().GetString("url")
			if errUR != nil {
				return fmt.Errorf("error parsing flag 'url' for createURL : %v", errUR)
			}

			o.URLReq = &govultr.SnapshotURLReq{
				URL: url,
			}

			snapshot, err := o.createURL()
			if err != nil {
				return fmt.Errorf("error creating snapshot from URL : %v", err)
			}

			data := &SnapshotPrinter{Snapshot: snapshot}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	createURL.Flags().StringP("url", "u", "", "Remote URL from where the snapshot will be downloaded.")
	if err := createURL.MarkFlagRequired("url"); err != nil {
		fmt.Printf("error marking snapshot create 'url' flag required: %v", err)
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <Snapshot ID>",
		Short:   "Delete a snapshot",
		Aliases: []string{"destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a snapshot ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting snapshot : %v", err)
			}

			o.Base.Printer.Display(printer.Info("snapshot has been deleted"), nil)

			return nil
		},
	}

	cmd.AddCommand(
		list,
		get,
		create,
		createURL,
		del,
	)

	return cmd
}

type options struct {
	Base   *cli.Base
	Req    *govultr.SnapshotReq
	URLReq *govultr.SnapshotURLReq
}

func (o *options) list() ([]govultr.Snapshot, *govultr.Meta, error) {
	snapshots, meta, _, err := o.Base.Client.Snapshot.List(o.Base.Context, o.Base.Options)
	return snapshots, meta, err
}

func (o *options) get() (*govultr.Snapshot, error) {
	snapshot, _, err := o.Base.Client.Snapshot.Get(o.Base.Context, o.Base.Args[0])
	return snapshot, err
}

func (o *options) create() (*govultr.Snapshot, error) {
	snapshot, _, err := o.Base.Client.Snapshot.Create(o.Base.Context, o.Req)
	return snapshot, err
}

func (o *options) createURL() (*govultr.Snapshot, error) {
	snapshot, _, err := o.Base.Client.Snapshot.CreateFromURL(o.Base.Context, o.URLReq)
	return snapshot, err
}

func (o *options) del() error {
	return o.Base.Client.Snapshot.Delete(o.Base.Context, o.Base.Args[0])
}
