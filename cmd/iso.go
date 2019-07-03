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

	"github.com/vultr/vultr-cli/cmd/printer"

	"github.com/spf13/cobra"
)

// isoCmd represents the iso command
var isoCmd = &cobra.Command{
	Use:   "iso",
	Short: "iso is used to access iso commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iso called")
	},
}

func Iso() *cobra.Command {
	isoCmd = &cobra.Command{
		Use:   "iso",
		Short: "iso is used to access iso commands",
		Long:  ``,
	}

	isoCmd.AddCommand(isoPrivate)
	isoCmd.AddCommand(isoPublic)
	isoCmd.AddCommand(isoCreate)
	isoCreate.Flags().StringP("url", "u", "", "url from where the ISO will be downloaded")
	isoCreate.MarkFlagRequired("url")

	isoCmd.AddCommand(isoDelete)

	return isoCmd
}

var isoPrivate = &cobra.Command{
	Use:   "private",
	Short: "list all private ISOs available",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		isos, err := client.ISO.List(context.TODO())

		if err != nil {
			fmt.Printf("error getting private ISOs : %v", err)
			os.Exit(1)
		}

		printer.IsoPrivate(isos)
	},
}

var isoPublic = &cobra.Command{
	Use:   "public",
	Short: "list all public ISOs available",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		isos, err := client.ISO.GetPublicList(context.TODO())

		if err != nil {
			fmt.Printf("error getting public ISOs : %v", err)
			os.Exit(1)
		}

		printer.IsoPublic(isos)
	},
}

var isoCreate = &cobra.Command{
	Use:   "create",
	Short: "create iso from url",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		_, err := client.ISO.CreateFromURL(context.TODO(), url)

		if err != nil {
			fmt.Printf("error creating ISOs : %v", err)
			os.Exit(1)
		}

		fmt.Println("ISO is in the process of being created")
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

		i, _ := strconv.Atoi(id)
		err := client.ISO.Delete(context.TODO(), i)

		if err != nil {
			fmt.Printf("error deleting ISOs : %v", err)
			os.Exit(1)
		}

		fmt.Println("ISO has been deleted")
	},
}
