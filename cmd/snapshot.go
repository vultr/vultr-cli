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
	"github.com/vultr/vultr-cli/cmd/printer"
)

// Snapshot represents the snapshot command
func Snapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot",
		Aliases: []string{"sn"},
		Short:   "snapshot commands",
		Long:    `snapshot is used to access snapshot commands`,
	}

	cmd.AddCommand(snapshotCreate, snapshotCreateFromURL, snapshotGet, snapshotDelete, snapshotList)

	snapshotCreate.Flags().StringP("id", "i", "", "ID of the virtual machine to create a snapshot from.")
	snapshotCreate.Flags().StringP("description", "d", "", "(optional) Description of snapshot contents")
	snapshotCreate.MarkFlagRequired("id")

	snapshotCreateFromURL.Flags().StringP("url", "u", "", "Remote URL from where the snapshot will be downloaded.")
	snapshotCreateFromURL.MarkFlagRequired("url")

	snapshotList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	snapshotList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	return cmd
}

// Create snapshot command
var snapshotCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a snapshot",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		desc, _ := cmd.Flags().GetString("description")

		options := &govultr.SnapshotReq{
			InstanceID:  id,
			Description: desc,
		}

		s, err := client.Snapshot.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Snapshot(s)
	},
}

// Create snapshot from URL command
var snapshotCreateFromURL = &cobra.Command{
	Use:   "create-url",
	Short: "Create a snapshot from a URL",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		options := &govultr.SnapshotURLReq{
			URL: url,
		}

		s, err := client.Snapshot.CreateFromURL(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Snapshot(s)
	},
}

// Delete snapshot command
var snapshotDelete = &cobra.Command{
	Use:     "delete <snapshotID>",
	Short:   "Delete a snapshot",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a snapshotID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Snapshot.Delete(context.Background(), id); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Snapshot has been deleted")
	},
}

// snapshotGet snapshot command
var snapshotGet = &cobra.Command{
	Use:   "get <snapshotID>",
	Short: "Get a snapshot",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a snapshotID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		snapshot, err := client.Snapshot.Get(context.Background(), args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Snapshot(snapshot)
	},
}

// List all snapshots command
var snapshotList = &cobra.Command{
	Use:   "list",
	Short: "List all snapshots",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, err := client.Snapshot.List(context.Background(), options)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		printer.Snapshots(list, meta)
	},
}
