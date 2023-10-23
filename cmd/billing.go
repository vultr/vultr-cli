// Copyright © 2021 The Vultr-cli Authors
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

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	billingLong    = `Get all available commands for billing`
	billingExample = `
	# Full example
	vultr-cli billing
	`

	historyLong    = `Get all available commands for billing history`
	historyExample = `
	# Full example
	vultr-cli billing history

	# Shortened with alias commands
	vultr-cli billing h
	`

	historyListLong    = `Retrieve a list of all billing history on your account`
	historyListExample = `
	# Full example
	vultr-cli billing history list

	# Shortened with alias commands
	vultr-cli billing h l
	`

	invoicesLong    = `Get all available commands for billing invoices`
	invoicesExample = `
	# Full example
	vultr-cli billing invoice

	# Shortened with alias commands
	vultr-cli billing i
	`

	invoiceListLong    = `Retrieve a list of all invoices on your account`
	invoiceListExample = `
	# Full example
	vultr-cli billing invoice list

	# Shortened with alias commands
	vultr-cli billing i l
	`

	invoiceGetLong    = `Get a specific invoice on your account`
	invoiceGetExample = `
	# Full example
	vultr-cli billing invoice get 123456

	# Shortened with alias commands
	vultr-cli billing i g 123456
	`

	invoiceItemsListLong    = `Retrieve a list of invoice items from a specific invoice on your account`
	invoiceItemsListExample = `
	# Full example
	vultr-cli billing invoice items 123456

	# Shortened with alias commands
	vultr-cli billing i i 123456
	`
)

// Billing represents the billing command
func Billing() *cobra.Command {
	billingCmd := &cobra.Command{
		Use:     "billing",
		Short:   "Display billing information",
		Long:    billingLong,
		Example: billingExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if auth := cmd.Context().Value("authenticated"); auth != true {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	historyCmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "Display billing history information",
		Long:    historyLong,
		Example: historyExample,
	}

	historyCmd.AddCommand(billingHistoryList)

	invoiceCmd := &cobra.Command{
		Use:     "invoice",
		Aliases: []string{"i"},
		Short:   "Display invoice information",
		Long:    invoicesLong,
		Example: invoicesExample,
	}

	invoicesList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	invoicesList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	invoiceItemsList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	invoiceItemsList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	invoiceCmd.AddCommand(invoicesList, invoiceGet, invoiceItemsList)

	billingHistoryList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	billingHistoryList.Flags().IntP(
		"per-page",
		"p",
		perPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	billingCmd.AddCommand(historyCmd, invoiceCmd)

	return billingCmd
}

var billingHistoryList = &cobra.Command{
	Use:     "list",
	Short:   "list billing history",
	Aliases: []string{"l"},
	Long:    historyListLong,
	Example: historyListExample,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		history, meta, _, err := client.Billing.ListHistory(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting billing history : %v\n", err)
			os.Exit(1)
		}

		printer.BillingHistory(history, meta)
	},
}

var invoicesList = &cobra.Command{
	Use:     "list",
	Short:   "list billing invoices",
	Aliases: []string{"l"},
	Long:    invoiceListLong,
	Example: invoiceListExample,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		history, meta, _, err := client.Billing.ListInvoices(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting invoices : %v\n", err)
			os.Exit(1)
		}

		printer.Invoices(history, meta)
	},
}

var invoiceItemsList = &cobra.Command{
	Use:     "items <invoiceID>",
	Short:   "list invoice items",
	Aliases: []string{"i"},
	Long:    invoiceItemsListLong,
	Example: invoiceItemsListExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an invoiceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := strconv.Atoi(args[0])
		options := getPaging(cmd)
		items, meta, _, err := client.Billing.ListInvoiceItems(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error getting invoice items : %v\n", err)
			os.Exit(1)
		}

		printer.InvoiceItems(items, meta)
	},
}

var invoiceGet = &cobra.Command{
	Use:     "get",
	Short:   "get invoice",
	Aliases: []string{"g"},
	Long:    invoiceGetLong,
	Example: invoiceGetExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an invoiceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		invoice, _, err := client.Billing.GetInvoice(context.Background(), args[0])
		if err != nil {
			fmt.Printf("error getting invoice : %v\n", err)
			os.Exit(1)
		}

		printer.Invoice(invoice)
	},
}
