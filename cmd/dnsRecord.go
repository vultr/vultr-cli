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
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// dnsRecordCmd represents the dnsRecord command
func DnsRecord() *cobra.Command {
	dnsRecordCmd := &cobra.Command{
		Use:   "record",
		Short: "dns record",
		Long:  ``,
	}

	dnsRecordCmd.AddCommand(recordCreate, recordlist, recordDelete, recordUpdate)

	// Create
	recordCreate.Flags().StringP("domain", "m", "", "name of domain you want to create this record for")
	recordCreate.Flags().StringP("type", "t", "", "record type you want to create : Possible values A, AAAA, CNAME, NS, MX, SRV, TXT CAA, SSHFP")
	recordCreate.Flags().StringP("name", "n", "", "name of record")
	recordCreate.Flags().StringP("data", "d", "", "data for the record")
	recordCreate.MarkFlagRequired("domain")
	recordCreate.MarkFlagRequired("type")
	recordCreate.MarkFlagRequired("name")
	recordCreate.MarkFlagRequired("data")
	// Create Optional
	recordCreate.Flags().IntP("ttl", "", 0, "time to live for the record")
	recordCreate.Flags().IntP("priority", "p", 0, "only required for MX and SRV")

	// Update
	recordUpdate.Flags().StringP("name", "n", "", "name of record")
	recordUpdate.Flags().StringP("data", "d", "", "data for the record")
	recordUpdate.Flags().IntP("ttl", "", 0, "time to live for the record")
	recordUpdate.Flags().IntP("priority", "p", 0, "only required for MX and SRV")

	return dnsRecordCmd
}

// Temporary solution to determine if the record type is TXT, in order to
// add quotes around the value. The API does not accept TXT records without
//	quotes.
var regRecordTxt = regexp.MustCompile("([A-Z]|=|_)")

var recordCreate = &cobra.Command{
	Use:   "create",
	Short: "create a dns record",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		rType, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		data, _ := cmd.Flags().GetString("data")
		// Record data for TXT must be enclosed in quotes
		if data[0] != '"' && data[len(data)-1] != '"' && regRecordTxt.Match([]byte(data)) {
			data = fmt.Sprintf("\"%s\"", data)
		}
		ttl, _ := cmd.Flags().GetInt("ttl")
		priority, _ := cmd.Flags().GetInt("priority")

		err := client.DNSRecord.Create(context.TODO(), domain, rType, name, data, ttl, priority)

		if err != nil {
			fmt.Printf("error while creating dns record : %v", err)
			os.Exit(1)
		}
		fmt.Println("created dns record")
	},
}

var recordlist = &cobra.Command{
	Use:   "list <domainName>",
	Short: "list all dns records",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a domain name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		records, err := client.DNSRecord.List(context.TODO(), domain)
		if err != nil {
			fmt.Printf("error while getting dns records : %v", err)
			os.Exit(1)
		}

		printer.DnsRecordsList(records)
	},
}

var recordDelete = &cobra.Command{
	Use:     "delete <domainName> <recordID>",
	Short:   "delete dns record",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a domainName & recordID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		id := args[1]
		err := client.DNSRecord.Delete(context.TODO(), domain, id)

		if err != nil {
			fmt.Printf("error while deleting dns record : %v", err)
			os.Exit(1)
		}

		fmt.Println("deleted dns record")
	},
}

var recordUpdate = &cobra.Command{
	Use:   "update <domainName> <recordID>",
	Short: "update dns record",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a domainName & recordID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		id := args[1]
		name, _ := cmd.Flags().GetString("name")
		data, _ := cmd.Flags().GetString("data")
		ttl, _ := cmd.Flags().GetInt("ttl")
		priority, _ := cmd.Flags().GetInt("priority")

		updates := &govultr.DNSRecord{}
		i, _ := strconv.Atoi(id)
		updates.RecordID = i

		if name != "" {
			updates.Name = name
		}

		if data != "" {
			// Record data for TXT must be enclosed in quotes
			if data[0] != '"' && data[len(data)-1] != '"' && regRecordTxt.Match([]byte(data)) {
				data = fmt.Sprintf("\"%s\"", data)
			}
			updates.Data = data
		}

		if ttl != 0 {
			updates.TTL = ttl
		}

		if priority != 0 {
			updates.Priority = priority
		}

		err := client.DNSRecord.Update(context.TODO(), domain, updates)

		if err != nil {
			fmt.Printf("error updating dns record : %v", err)
			os.Exit(1)
		}

		fmt.Println("updated dns record")
	},
}
