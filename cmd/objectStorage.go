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
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

// ObjectStorageCmd represents the objStorageCmd command
func ObjectStorageCmd() *cobra.Command {
	objStorageCmd := &cobra.Command{
		Use:     "object-storage",
		Aliases: []string{"objStorage"},
		Short:   "object storage commands",
		Long:    `object-storage is used to interact with the object-storage api`,
	}

	objStorageCmd.AddCommand(objStorageCreate, objStorageLabelSet, objStorageList, objStorageGet, objStorageClusterList, objStorageS3KeyRegenerate, objStorageDestroy)

	// Create
	objStorageCreate.Flags().StringP("label", "l", "", "label you want your object storage to have")
	objStorageCreate.Flags().IntP("obj-store-clusterid", "o", 0, "obj-store-clusterid you want to create the object storage in")
	if err := objStorageCreate.MarkFlagRequired("obj-store-clusterid"); err != nil {
		fmt.Printf("error marking object-storage create 'obj-store-clusterid' flag required: %v\n", err)
		os.Exit(1)
	}

	// Label
	objStorageLabelSet.Flags().StringP("label", "l", "", "label you want your object storage to have")
	if err := objStorageLabelSet.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking object-storage create 'label' flag required: %v\n", err)
		os.Exit(1)
	}

	// List
	objStorageList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	objStorageList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Regenerate
	objStorageS3KeyRegenerate.Flags().StringP("s3-access-key", "s", "", "access key for a given object storage subscription")

	// Cluster List
	objStorageClusterList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	objStorageClusterList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return objStorageCmd
}

var objStorageCreate = &cobra.Command{
	Use:   "create",
	Short: "create a new object storage subscription",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		objectStoreClusterID, _ := cmd.Flags().GetInt("obj-store-clusterid")
		label, _ := cmd.Flags().GetString("label")

		objStorage, _, err := client.ObjectStorage.Create(context.TODO(), objectStoreClusterID, label)
		if err != nil {
			fmt.Printf("error creating object storage : %v\n", err)
			os.Exit(1)
		}

		printer.SingleObjectStorage(objStorage)
	},
}

var objStorageLabelSet = &cobra.Command{
	Use:   "label <objectStorageID>",
	Short: "change the label for object storage subscription",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an objectStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		label, _ := cmd.Flags().GetString("label")

		err := client.ObjectStorage.Update(context.TODO(), id, label)
		if err != nil {
			fmt.Printf("error setting label : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("set label on object storage : %v\n", id)
	},
}

var objStorageList = &cobra.Command{
	Use:   "list",
	Short: "retrieves a list of active object storage subscriptions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		objStorage, meta, _, err := client.ObjectStorage.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting object storage : %v\n", err)
			os.Exit(1)
		}

		printer.ObjectStorages(objStorage, meta)
	},
}

var objStorageGet = &cobra.Command{
	Use:   "get <objectStorageID>",
	Short: "retrieves a given object storage subscription",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an objectStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		objStorage, _, err := client.ObjectStorage.Get(context.TODO(), id)
		if err != nil {
			fmt.Printf("error getting object storage : %v\n", err)
			os.Exit(1)
		}

		printer.SingleObjectStorage(objStorage)
	},
}

var objStorageClusterList = &cobra.Command{
	Use:   "list-cluster",
	Short: "retrieves a list of object storage clusters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		cluster, meta, _, err := client.ObjectStorage.ListCluster(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting object storage clusters : %v\n", err)
			os.Exit(1)
		}

		printer.ObjectStorageClusterList(cluster, meta)
	},
}

var objStorageS3KeyRegenerate = &cobra.Command{
	Use:   "s3key-regenerate <objectStorageID>",
	Short: "regenerate the S3 API keys of an object storage subscription",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an objectStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		s3Keys, _, err := client.ObjectStorage.RegenerateKeys(context.TODO(), id)
		if err != nil {
			fmt.Printf("error regenerating object storage keys : %v\n", err)
			os.Exit(1)
		}

		printer.ObjStorageS3KeyRegenerate(s3Keys)
	},
}

var objStorageDestroy = &cobra.Command{
	Use:   "delete <objectStorageID>",
	Short: "deletes an object storage subscription",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an objectStorageID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.ObjectStorage.Delete(context.TODO(), id); err != nil {
			fmt.Printf("error destroying object storage subscription : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("destroyed object storage subscription")
	},
}
