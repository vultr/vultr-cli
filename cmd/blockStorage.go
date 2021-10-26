// Copyright © 2019 The Vultr-cli Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// BlockStorageCmd represents the blockStorage command
func BlockStorageCmd() *cobra.Command {

	bsCmd := &cobra.Command{
		Use:     "block-storage",
		Aliases: []string{"bs"},
		Short:   "block storage commands",
		Long:    `block-storage is used to interact with the block-storage api`,
	}

	bsCmd.AddCommand(bsAttach, bsCreate, bsDelete, bsDetach, bsLabelSet, bsList, bsGet, bsResize)

	// List
	bsList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	bsList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Attach
	bsAttach.Flags().StringP("instance", "i", "", "instance id you want to attach to")
	bsAttach.Flags().Bool("live", false, "attach Block Storage without restarting the Instance.")
	bsAttach.MarkFlagRequired("instance")

	// Detach
	bsDetach.Flags().Bool("live", false, "detach block storage from instance without a restart")

	// Create
	bsCreate.Flags().StringP("region", "r", "", "regionID you want to create the block storage in")
	bsCreate.MarkFlagRequired("region")

	bsCreate.Flags().IntP("size", "s", 0, "size of the block storage you want to create")
	bsCreate.MarkFlagRequired("size")

	bsCreate.Flags().StringP("label", "l", "", "label you want to give the block storage")

	// Label
	bsLabelSet.Flags().StringP("label", "l", "", "label you want your block storage to have")
	bsLabelSet.MarkFlagRequired("label")

	// Resize
	bsResize.Flags().IntP("size", "s", 0, "size you want your block storage to be")
	bsResize.MarkFlagRequired("size")

	return bsCmd
}

var bsAttach = &cobra.Command{
	Use:   "attach <blockStorageID>",
	Short: "attaches a block storage to an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		instance, _ := cmd.Flags().GetString("instance")
		live, _ := cmd.Flags().GetBool("live")

		bsAttach := &govultr.BlockStorageAttach{
			InstanceID: instance,
			Live:       govultr.BoolToBoolPtr(live),
		}

		if err := client.BlockStorage.Attach(context.Background(), id, bsAttach); err != nil {
			fmt.Printf("error attaching block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("attached block storage")
	},
}

var bsCreate = &cobra.Command{
	Use:   "create",
	Short: "create a new block storage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		size, _ := cmd.Flags().GetInt("size")
		label, _ := cmd.Flags().GetString("label")

		bsCreate := &govultr.BlockStorageCreate{
			Region: region,
			SizeGB: size,
			Label:  label,
		}

		bs, err := client.BlockStorage.Create(context.Background(), bsCreate)
		if err != nil {
			fmt.Printf("error creating block storage : %v\n", err)
			os.Exit(1)
		}

		printer.SingleBlockStorage(bs)

	},
}

var bsDelete = &cobra.Command{
	Use:     "delete <blockStorageID>",
	Short:   "delete a block storage",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.BlockStorage.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("deleted block storage")
	},
}

var bsDetach = &cobra.Command{
	Use:   "detach <blockStorageID>",
	Short: "detaches a block storage from an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		live, _ := cmd.Flags().GetBool("live")

		bsDetach := &govultr.BlockStorageDetach{
			Live: govultr.BoolToBoolPtr(live),
		}

		if err := client.BlockStorage.Detach(context.Background(), id, bsDetach); err != nil {
			fmt.Printf("error detaching block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("detached block storage")
	},
}

var bsLabelSet = &cobra.Command{
	Use:   "label <blockStorageID>",
	Short: "sets a label for a block storage",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		label, _ := cmd.Flags().GetString("label")

		options := &govultr.BlockStorageUpdate{
			Label: label,
		}

		if err := client.BlockStorage.Update(context.Background(), id, options); err != nil {
			fmt.Printf("error setting label : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("set label on block storage : %s\n", id)
	},
}

// List all of individual block storage
var bsList = &cobra.Command{
	Use:   "list",
	Short: "retrieves a list of active block storage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		bs, meta, err := client.BlockStorage.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting block storage : %v\n", err)
			os.Exit(1)
		}

		printer.BlockStorage(bs, meta)
	},
}

// Get a block storage
var bsGet = &cobra.Command{
	Use:   "get <blockStorageID>",
	Short: "retrieves a block storage",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		bs, err := client.BlockStorage.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting block storage : %v\n", err)
			os.Exit(1)
		}

		printer.SingleBlockStorage(bs)
	},
}

var bsResize = &cobra.Command{
	Use:   "resize <blockStorageID>",
	Short: "resize a block storage",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		size, _ := cmd.Flags().GetInt("size")

		options := &govultr.BlockStorageUpdate{
			SizeGB: size,
		}

		if err := client.BlockStorage.Update(context.Background(), id, options); err != nil {
			fmt.Printf("error resizing block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("resized block storage")
	},
}
