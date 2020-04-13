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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr"
	"github.com/vultr/vultr-cli/cmd/printer"
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
	objStorageCreate.MarkFlagRequired("obj-store-clusterid")

	// Label
	objStorageLabelSet.Flags().StringP("label", "l", "", "label you want your object storage to have")
	objStorageLabelSet.MarkFlagRequired("label")

	// List
	objStorageList.Flags().StringP("include-s3", "i", "", "(optional) Whether to include s3 keys with each subscription entry. Possible values: 'yes', 'no'. Defaults to 'yes'.")

	// Get
	objStorageGet.Flags().StringP("include-s3", "i", "", "(optional) Whether to include s3 keys with subscription entry. Possible values: 'yes', 'no'. Defaults to 'yes'.")

	// Regenerate
	objStorageS3KeyRegenerate.Flags().StringP("s3-access-key", "s", "", "access key for a given object storage subscription")

	return objStorageCmd
}

var objStorageCreate = &cobra.Command{
	Use:   "create",
	Short: "create a new object storage subscription",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		objectStoreClusterID, _ := cmd.Flags().GetInt("obj-store-clusterid")
		label, _ := cmd.Flags().GetString("label")

		objStorage, err := client.ObjectStorage.Create(context.TODO(), objectStoreClusterID, label)

		if err != nil {
			fmt.Printf("error creating object storage : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("created object storage - ID : %v", objStorage.ID)
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
		id, _ := strconv.Atoi(args[0])
		label, _ := cmd.Flags().GetString("label")

		err := client.ObjectStorage.SetLabel(context.TODO(), id, label)
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
		includeS3, _ := cmd.Flags().GetString("include-s3")
		options := &govultr.ObjectListOptions{
			IncludeS3: true,
		}

		if strings.ToLower(includeS3) == "no" {
			options.IncludeS3 = false
		}

		objStorage, err := client.ObjectStorage.List(context.TODO(), options)

		if err != nil {
			fmt.Printf("error getting object storage : %v\n", err)
			os.Exit(1)
		}

		printer.ObjectStorage(objStorage, options)
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
		includeS3, _ := cmd.Flags().GetString("include-s3")
		options := &govultr.ObjectListOptions{
			IncludeS3: true,
		}

		if strings.ToLower(includeS3) == "no" {
			options.IncludeS3 = false
		}

		id, _ := strconv.Atoi(args[0])
		objStorage, err := client.ObjectStorage.Get(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting object storage : %v\n", err)
			os.Exit(1)
		}

		printer.SingleObjectStorage(objStorage, options)
	},
}

var objStorageClusterList = &cobra.Command{
	Use:   "list-cluster",
	Short: "retrieves a list of object storage clusters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, err := client.ObjectStorage.ListCluster(context.TODO())

		if err != nil {
			fmt.Printf("error getting object storage clusters : %v\n", err)
			os.Exit(1)
		}

		printer.ObjectStorageClusterList(cluster)
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
		id, _ := strconv.Atoi(args[0])
		s3AccessKey, _ := cmd.Flags().GetString("s3-access-key")
		s3Keys, err := client.ObjectStorage.RegenerateKeys(context.TODO(), id, s3AccessKey)

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
		id, _ := strconv.Atoi(args[0])
		err := client.ObjectStorage.Delete(context.TODO(), id)

		if err != nil {
			fmt.Printf("error destroying object storage subscription : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("destroyed object storage subscription")
	},
}
