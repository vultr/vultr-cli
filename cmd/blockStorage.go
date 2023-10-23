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
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	attachBlockStorageLong    = `Attaches a block storage resource to an specified instance`
	attachBlockStorageExample = `
	#Full example
	vultr-cli block-storage attach 67181686-5455-4ebb-81eb-7299f3506e2c --instance=a7898453-dd9e-4b47-bdab-9dd7a3448f1f

	#Shortened with aliased commands
	vultr-cli bs a 67181686-5455-4ebb-81eb-7299f3506e2c -i=a7898453-dd9e-4b47-bdab-9dd7a3448f1f
	`

	createBlockStorageLong    = `Create a new block storage resource in a specified region`
	createBlockStorageExample = `
	#Full example
	vultr-cli block-storage create --region='lax' --size=10

	#Full example with block-type
	vultr-cli block-storage create --region='lax' --size=10 --block-type='high_perf'

	#Shortened with aliased commands
	vultr-cli bs c -r='lax' -s=10

	#Shortened with aliased commands and block-type
	vultr-cli bs c -r='lax' -s=10 -b='high_perf'
	`

	deleteBlockStorageLong    = `Delete a block storage resource`
	deleteBlockStorageExample = `
	#Full example
	vultr-cli block-storage delete 67181686-5455-4ebb-81eb-7299f3506e2c

	#Shortened with aliased commands
	vultr-cli bs d 67181686-5455-4ebb-81eb-7299f3506e2c
	`

	detachBlockStorageLong    = `Detach a block storage resource from an instance`
	detachBlockStorageExample = `
	#Full example
	vultr-cli block-storage detach 67181686-5455-4ebb-81eb-7299f3506e2c

	#Shortened with aliased commands
	vultr-cli bs detach 67181686-5455-4ebb-81eb-7299f3506e2c
	`

	labelBlockStorageLong    = `Set a label for a block storage resource`
	labelBlockStorageExample = `
	#Full example
	vultr-cli block-storage label 67181686-5455-4ebb-81eb-7299f3506e2c --label="Example Label"

	#Shortened with aliased commands
	vultr-cli bs label 67181686-5455-4ebb-81eb-7299f3506e2c -l="Example Label"
	`

	listBlockStorageLong    = `Retrieves a list of active block storage resources`
	listBlockStorageExample = `
	#Full example
	vultr-cli block-storage list

	#Shortened with aliased commands
	vultr-cli bs l
	`

	getBlockStorageLong    = `Retrieves a specified block storage resource`
	getBlockStorageExample = `
	#Full example
	vultr-cli block-storage get 67181686-5455-4ebb-81eb-7299f3506e2c

	#Shortened with aliased commands
	vultr-cli bs g 67181686-5455-4ebb-81eb-7299f3506e2c
	`

	resizeBlockStorageLong    = `Resizes a specified block storage resource`
	resizeBlockStorageExample = `
	#Full example
	vultr-cli block-storage resize 67181686-5455-4ebb-81eb-7299f3506e2c --size=20

	#Shortened with aliased commands
	vultr-cli bs r 67181686-5455-4ebb-81eb-7299f3506e2c -s=20
	`
)

// BlockStorageCmd represents the blockStorage command
func BlockStorageCmd() *cobra.Command {
	bsCmd := &cobra.Command{
		Use:     "block-storage",
		Aliases: []string{"bs"},
		Short:   "block storage commands",
		Long:    `block-storage is used to interact with the block-storage api`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if auth := cmd.Context().Value("authenticated"); auth != true {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	bsCmd.AddCommand(bsAttach, bsCreate, bsDelete, bsDetach, bsLabelSet, bsList, bsGet, bsResize)

	// List
	bsList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	bsList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Attach
	bsAttach.Flags().StringP("instance", "i", "", "instance id you want to attach to")
	bsAttach.Flags().Bool("live", false, "attach Block Storage without restarting the Instance.")
	if err := bsAttach.MarkFlagRequired("instance"); err != nil {
		fmt.Printf("error marking block storage attach 'live' flag required: %v\n", err)
		os.Exit(1)
	}
	// Detach
	bsDetach.Flags().Bool("live", false, "detach block storage from instance without a restart")

	// Create
	bsCreate.Flags().StringP("region", "r", "", "regionID you want to create the block storage in")
	if err := bsCreate.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking block storage create 'region' flag required: %v\n", err)
		os.Exit(1)
	}

	bsCreate.Flags().IntP("size", "s", 0, "size of the block storage you want to create")
	if err := bsCreate.MarkFlagRequired("size"); err != nil {
		fmt.Printf("error marking block storage create 'size' flag required: %v\n", err)
		os.Exit(1)
	}

	bsCreate.Flags().StringP("label", "l", "", "label you want to give the block storage")

	bsCreate.Flags().StringP(
		"block-type",
		"b",
		"",
		`(optional) Block type you want to give the block storage.
		Possible values: 'high_perf', 'storage_opt'. Currently defaults to 'high_perf'.`,
	)

	// Label
	bsLabelSet.Flags().StringP("label", "l", "", "label you want your block storage to have")
	if err := bsLabelSet.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking block storage label set 'label' flag required: %v\n", err)
		os.Exit(1)
	}

	// Resize
	bsResize.Flags().IntP("size", "s", 0, "size you want your block storage to be")
	if err := bsResize.MarkFlagRequired("size"); err != nil {
		fmt.Printf("error marking block storage resize 'size' flag required: %v\n", err)
		os.Exit(1)
	}

	return bsCmd
}

var bsAttach = &cobra.Command{
	Use:     "attach <blockStorageID>",
	Short:   "attaches a block storage to an instance",
	Aliases: []string{"a"},
	Long:    attachBlockStorageLong,
	Example: attachBlockStorageExample,
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
	Use:     "create",
	Short:   "create a new block storage",
	Aliases: []string{"c"},
	Long:    createBlockStorageLong,
	Example: createBlockStorageExample,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		size, _ := cmd.Flags().GetInt("size")
		label, _ := cmd.Flags().GetString("label")
		blockType, _ := cmd.Flags().GetString("block-type")

		bsCreate := &govultr.BlockStorageCreate{
			Region:    region,
			SizeGB:    size,
			Label:     label,
			BlockType: blockType,
		}

		bs, _, err := client.BlockStorage.Create(context.Background(), bsCreate)
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
	Aliases: []string{"d", "destroy"},
	Long:    deleteBlockStorageLong,
	Example: deleteBlockStorageExample,
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
	Use:     "detach <blockStorageID>",
	Short:   "detaches a block storage from an instance",
	Long:    detachBlockStorageLong,
	Example: detachBlockStorageExample,
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
	Use:     "label <blockStorageID>",
	Short:   "sets a label for a block storage",
	Long:    labelBlockStorageLong,
	Example: labelBlockStorageExample,
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
	Use:     "list",
	Short:   "retrieves a list of active block storage",
	Aliases: []string{"l"},
	Long:    listBlockStorageLong,
	Example: listBlockStorageExample,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		bs, meta, _, err := client.BlockStorage.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting block storage : %v\n", err)
			os.Exit(1)
		}

		printer.BlockStorage(bs, meta)
	},
}

// Get a block storage
var bsGet = &cobra.Command{
	Use:     "get <blockStorageID>",
	Short:   "retrieves a block storage",
	Aliases: []string{"g"},
	Long:    getBlockStorageLong,
	Example: getBlockStorageExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		bs, _, err := client.BlockStorage.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting block storage : %v\n", err)
			os.Exit(1)
		}

		printer.SingleBlockStorage(bs)
	},
}

var bsResize = &cobra.Command{
	Use:     "resize <blockStorageID>",
	Short:   "resize a block storage",
	Aliases: []string{"r"},
	Long:    resizeBlockStorageLong,
	Example: resizeBlockStorageExample,
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
