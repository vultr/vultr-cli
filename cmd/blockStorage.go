// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/vultr/vultr-cli/cmd/printer"
)

// blockStorageCmd represents the blockStorage command
func BlockStorageCmd() *cobra.Command {

	bsCmd := &cobra.Command{
		Use:     "block-storage",
		Aliases: []string{"bs"},
		Short:   "block storage commands",
		Long:    `block-storage is used to interact with the block-storage api`,
	}

	bsCmd.AddCommand(bsAttach, bsCreate, bsDelete, bsDetach, bsLabelSet, bsList, bsResize)

	bsList.Flags().StringP("instance", "i", "", "get the block storage that is attached to a given instance id")

	// Attach
	bsAttach.Flags().StringP("instance", "i", "", "instance id you want to attach to")
	bsAttach.MarkFlagRequired("instance")

	bsAttach.Flags().StringP("live", "l", "", "attach block storage to the instance without a restart (yes or no)")

	// Detach
	bsDetach.Flags().StringP("live", "l", "", "detach block storage from instance without a restart")

	// Create
	bsCreate.Flags().IntP("region", "r", 0, "regionID you want to create the block storage in")
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
	Short: "",
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
		live, _ := cmd.Flags().GetString("live")

		err := client.BlockStorage.Attach(context.TODO(), id, instance, live)
		if err != nil {
			fmt.Printf("error attaching block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("attached block storage")
	},
}

var bsCreate = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetInt("region")
		size, _ := cmd.Flags().GetInt("size")
		label, _ := cmd.Flags().GetString("label")

		bs, err := client.BlockStorage.Create(context.TODO(), region, size, label)
		if err != nil {
			fmt.Printf("error creating block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("created block storage - ID : %v\n", bs.BlockStorageID)
	},
}

var bsDelete = &cobra.Command{
	Use:     "delete <blockStorageID>",
	Short:   "",
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
		if err := client.BlockStorage.Delete(context.TODO(), id); err != nil {
			fmt.Printf("error deleting block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("deleted block storage")
	},
}

var bsDetach = &cobra.Command{
	Use:   "detach <blockStorageID>",
	Short: "",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a blockStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		live, _ := cmd.Flags().GetString("live")

		err := client.BlockStorage.Detach(context.TODO(), id, live)
		if err != nil {
			fmt.Printf("error detaching block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("detached block storage")
	},
}

var bsLabelSet = &cobra.Command{
	Use:   "label <blockStorageID>",
	Short: "",
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

		err := client.BlockStorage.SetLabel(context.TODO(), id, label)
		if err != nil {
			fmt.Printf("error setting label : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("set label on block storage : %s\n", id)
	},
}

// List all of individual block storage devices
var bsList = &cobra.Command{
	Use:   "list",
	Short: "retrieves a list of active block storage devices",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		instance, _ := cmd.Flags().GetString("instance")

		if instance == "" {
			bs, err := client.BlockStorage.List(context.TODO())
			if err != nil {
				fmt.Printf("error getting block storage : %v\n", err)
				os.Exit(1)
			}

			printer.BlockStorage(bs)

		} else {
			bs, err := client.BlockStorage.Get(context.TODO(), instance)
			if err != nil {
				fmt.Printf("error getting block storage : %v\n", err)
				os.Exit(1)
			}

			printer.SingleBlockStorage(bs)
		}

	},
}

var bsResize = &cobra.Command{
	Use:   "resize <blockStorageID>",
	Short: "resize a block storage device",
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

		if err := client.BlockStorage.Resize(context.TODO(), id, size); err != nil {
			fmt.Printf("error resizing block storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("resized block storage device")
	},
}
