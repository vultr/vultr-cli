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

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"

	"github.com/spf13/cobra"
)

// ISO represents the iso command
func ISO() *cobra.Command {
	isoCmd := &cobra.Command{
		Use:   "iso",
		Short: "iso is used to access iso commands",
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if auth := cmd.Context().Value("authenticated"); auth != true {
				return fmt.Errorf(apiKeyError)
			}
			return nil
		},
	}

	isoCmd.AddCommand(isoCreate, isoDelete, isoPrivateGet, isoPrivateList, isoPublic)
	isoCreate.Flags().StringP("url", "u", "", "url from where the ISO will be downloaded")
	if err := isoCreate.MarkFlagRequired("url"); err != nil {
		fmt.Printf("error marking iso create 'url' flag required: %v\n", err)
		os.Exit(1)
	}

	isoPrivateList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	isoPrivateList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	isoPublic.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	isoPublic.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	return isoCmd
}

var isoPrivateGet = &cobra.Command{
	Use:   "get <isoID>",
	Short: "get private ISO <isoID>",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an ISO id")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		iso, _, err := client.ISO.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting ISO : %v\n", err)
			os.Exit(1)
		}

		printer.IsoPrivate(iso)
	},
}

var isoPrivateList = &cobra.Command{
	Use:   "list",
	Short: "list all private ISOs available",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		isos, meta, _, err := client.ISO.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting private ISOs : %v\n", err)
			os.Exit(1)
		}

		printer.IsoPrivates(isos, meta)
	},
}

var isoPublic = &cobra.Command{
	Use:   "public",
	Short: "list all public ISOs available",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		isos, meta, _, err := client.ISO.ListPublic(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting public ISOs : %v\n", err)
			os.Exit(1)
		}

		printer.IsoPublic(isos, meta)
	},
}

var isoCreate = &cobra.Command{
	Use:   "create",
	Short: "create iso from url",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		options := &govultr.ISOReq{
			URL: url,
		}

		iso, _, err := client.ISO.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating ISOs : %v\n", err)
			os.Exit(1)
		}

		printer.IsoPrivate(iso)
	},
}

var isoDelete = &cobra.Command{
	Use:     "delete <isoID>",
	Short:   "delete a private iso",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an isoID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.ISO.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting ISOs : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("ISO has been deleted")
	},
}
