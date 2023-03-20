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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	databaseLong    = `Get commands available to database`
	databaseExample = `
	# Full example
	vultr-cli database
	`
	databaseCreateLong    = `Create a new Managed Database with specified plan, region, and database engine/version`
	databaseCreateExample = `
	# Full example
	vultr-cli database create --database-engine="mysql" --database-engine-version="8" --region="ewr" --plan="vultr-dbaas-startup-cc-1-55-2" --label="example-db"

	# Full example with custom MySQL settings
	vultr-cli database create --database-engine="mysql" --database-engine-version="8" --region="ewr" --plan="vultr-dbaas-startup-cc-1-55-2" --label="example-db" --mysql-slow-query-log="true" --mysql-long-query-time="2"
	`
)

// Instance represents the instance command
func Database() *cobra.Command {
	databaseCmd := &cobra.Command{
		Use:     "database",
		Short:   "commands to interact with managed databases on vultr",
		Long:    databaseLong,
		Example: databaseExample,
	}

	databaseCmd.AddCommand(databaseList, databasePlanList)

	databasePlanList.Flags().StringP("engine", "e", "", "(optional) Filter by database engine type.")
	databasePlanList.Flags().StringP("nodes", "n", "", "(optional) Filter by number of nodes.")
	databasePlanList.Flags().StringP("region", "r", "", "(optional) Filter by region.")

	databaseList.Flags().StringP("label", "l", "", "(optional) Filter by label.")
	databaseList.Flags().StringP("tag", "t", "", "(optional) Filter by tag.")
	databaseList.Flags().StringP("region", "r", "", "(optional) Filter by region.")

	return databaseCmd
}

var databasePlanList = &cobra.Command{
	Use:     "list-plans",
	Aliases: []string{"l"},
	Short:   "list all available managed database plans",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := &govultr.DBPlanListOptions{}
		s, meta, _, err := client.Database.ListPlans(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting list of database plans : %v\n", err)
			os.Exit(1)
		}

		printer.DatabasePlanList(s, meta)
	},
}

var databaseList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all available managed databases",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := &govultr.DBListOptions{}
		s, meta, _, err := client.Database.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting list of databases : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseList(s, meta)
	},
}
