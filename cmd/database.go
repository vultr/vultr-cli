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
	databaseUpdateLong    = `Create a new Managed Database with specified plan, region, and database engine/version`
	databaseUpdateExample = `
	# Full example
	vultr-cli database update --region="sea" --plan="vultr-dbaas-startup-cc-2-80-4"

	# Full example with custom MySQL settings
	vultr-cli database update --mysql-slow-query-log="true" --mysql-long-query-time="2"
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

	databaseCmd.AddCommand(
		databasePlanList,
		databaseList, databaseCreate, databaseInfo, databaseUpdate, databaseDelete,
		databaseUserList, databaseUserCreate, databaseUserInfo, databaseUserUpdate, databaseUserDelete)

	// Plan list flags
	databasePlanList.Flags().StringP("engine", "e", "", "(optional) Filter by database engine type.")
	databasePlanList.Flags().IntP("nodes", "n", 0, "(optional) Filter by number of nodes.")
	databasePlanList.Flags().StringP("region", "r", "", "(optional) Filter by region.")

	// Database list flags
	databaseList.Flags().StringP("label", "l", "", "(optional) Filter by label.")
	databaseList.Flags().StringP("tag", "t", "", "(optional) Filter by tag.")
	databaseList.Flags().StringP("region", "r", "", "(optional) Filter by region.")

	// Database create flags
	databaseCreate.Flags().StringP("database-engine", "e", "", "database engine for the new manaaged database")
	databaseCreate.Flags().StringP("database-engine-version", "v", "", "database engine version for the new manaaged database")
	databaseCreate.Flags().StringP("region", "r", "", "region id for the new managed database")
	databaseCreate.Flags().StringP("plan", "p", "", "plan id for the new managed database")
	databaseCreate.Flags().StringP("label", "l", "", "label for the new managed database")
	databaseCreate.Flags().StringP("tag", "t", "", "tag for the new managed database")
	databaseCreate.Flags().StringP("maintenance-dow", "", "", "maintenance day of week for the new managed database")
	databaseCreate.Flags().StringP("maintenance-time", "", "", "maintenance time for the new managed database")
	databaseCreate.Flags().StringSliceP("trusted-ips", "", []string{}, "comma-separated list of trusted ip addresses for the new managed database")
	databaseCreate.Flags().StringSliceP("mysql-sql-modes", "", []string{}, "comma-separated list of sql modes for the new managed database")
	databaseCreate.Flags().BoolP("mysql-require-primary-key", "", true, "enable requiring primary keys for the new mysql managed database")
	databaseCreate.Flags().BoolP("mysql-slow-query-log", "", false, "enable slow query logging for the new mysql managed database")
	databaseCreate.Flags().StringP("mysql-long-query-time", "", "", "long query time for the new mysql managed database when slow query logging is enabled")
	databaseCreate.Flags().StringP("redis-eviction-policy", "", "", "eviction policy for the new redis managed database")

	// Database update flags
	databaseUpdate.Flags().StringP("database-engine", "e", "", "database engine for the manaaged database")
	databaseUpdate.Flags().StringP("database-engine-version", "v", "", "database engine version for the manaaged database")
	databaseUpdate.Flags().StringP("region", "r", "", "region id for the managed database")
	databaseUpdate.Flags().StringP("plan", "p", "", "plan id for the managed database")
	databaseUpdate.Flags().StringP("label", "l", "", "label for the managed database")
	databaseUpdate.Flags().StringP("tag", "t", "", "tag for the managed database")
	databaseUpdate.Flags().StringP("maintenance-dow", "", "", "maintenance day of week for the managed database")
	databaseUpdate.Flags().StringP("maintenance-time", "", "", "maintenance time for the managed database")
	databaseUpdate.Flags().StringP("cluster-time-zone", "", "", "configured time zone for the managed database")
	databaseUpdate.Flags().StringSliceP("trusted-ips", "", []string{}, "comma-separated list of trusted ip addresses for the managed database")
	databaseUpdate.Flags().StringSliceP("mysql-sql-modes", "", []string{}, "comma-separated list of sql modes for the managed database")
	databaseUpdate.Flags().BoolP("mysql-require-primary-key", "", true, "enable requiring primary keys for the mysql managed database")
	databaseUpdate.Flags().BoolP("mysql-slow-query-log", "", false, "enable slow query logging for the mysql managed database")
	databaseUpdate.Flags().StringP("mysql-long-query-time", "", "", "long query time for the mysql managed database when slow query logging is enabled")
	databaseUpdate.Flags().StringP("redis-eviction-policy", "", "", "eviction policy for the redis managed database")

	// Database user create flags
	databaseUserCreate.Flags().StringP("username", "u", "", "username for the new manaaged database user")
	databaseUserCreate.Flags().StringP("password", "p", "", "password for the new manaaged database user (omit or leave empty to generate a random secure password)")
	databaseUserCreate.Flags().StringP("encryption", "e", "", "encryption type for the new managed database user (MySQL only)")

	// Database user update flags
	databaseUserUpdate.Flags().StringP("password", "p", "", "password for the new manaaged database user (leave empty to generate a random secure password)")

	return databaseCmd
}

var databasePlanList = &cobra.Command{
	Use:     "list-plans",
	Aliases: []string{"lp"},
	Short:   "list all available managed database plans",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		engine, _ := cmd.Flags().GetString("engine")
		nodes, _ := cmd.Flags().GetInt("nodes")
		region, _ := cmd.Flags().GetString("region")
		options := &govultr.DBPlanListOptions{
			Engine: engine,
			Nodes:  nodes,
			Region: region,
		}
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
		label, _ := cmd.Flags().GetString("label")
		tag, _ := cmd.Flags().GetString("tag")
		region, _ := cmd.Flags().GetString("region")
		options := &govultr.DBListOptions{
			Label:  label,
			Tag:    tag,
			Region: region,
		}
		s, meta, _, err := client.Database.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting list of databases : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseList(s, meta)
	},
}

var databaseCreate = &cobra.Command{
	Use:     "create",
	Short:   "Create a managed database",
	Aliases: []string{"c"},
	Long:    databaseCreateLong,
	Example: databaseCreateExample,
	Run: func(cmd *cobra.Command, args []string) {
		databaseEngine, _ := cmd.Flags().GetString("database-engine")
		databaseEngineVersion, _ := cmd.Flags().GetString("database-engine-version")
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		label, _ := cmd.Flags().GetString("label")

		// Optional
		tag, _ := cmd.Flags().GetString("tag")
		maintenanceDOW, _ := cmd.Flags().GetString("maintenance-dow")
		maintenanceTime, _ := cmd.Flags().GetString("maintenance-time")
		trustedIPs, _ := cmd.Flags().GetStringSlice("trusted-ips")
		mysqlSQLModes, _ := cmd.Flags().GetStringSlice("mysql-sql-modes")
		mysqlRequirePrimaryKey, _ := cmd.Flags().GetBool("mysql-require-primary-key")
		mySQLSlowQueryLog, _ := cmd.Flags().GetBool("mysql-slow-query-log")
		mySQLLongQueryTime, _ := cmd.Flags().GetInt("mysql-long-query-time")
		redisEvictionPolicy, _ := cmd.Flags().GetString("redis-eviction-policy")

		opt := &govultr.DatabaseCreateReq{
			DatabaseEngine:         databaseEngine,
			DatabaseEngineVersion:  databaseEngineVersion,
			Plan:                   plan,
			Region:                 region,
			Label:                  label,
			Tag:                    tag,
			MaintenanceDOW:         maintenanceDOW,
			MaintenanceTime:        maintenanceTime,
			TrustedIPs:             trustedIPs,
			MySQLSQLModes:          mysqlSQLModes,
			MySQLRequirePrimaryKey: govultr.BoolToBoolPtr(true),
			MySQLSlowQueryLog:      govultr.BoolToBoolPtr(false),
			MySQLLongQueryTime:     mySQLLongQueryTime,
			RedisEvictionPolicy:    redisEvictionPolicy,
		}

		if !mysqlRequirePrimaryKey {
			opt.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(false)
		}

		if mySQLSlowQueryLog {
			opt.MySQLSlowQueryLog = govultr.BoolToBoolPtr(true)
		}

		// Make the request
		database, _, err := client.Database.Create(context.TODO(), opt)
		if err != nil {
			fmt.Printf("error creating managed database : %v\n", err)
			os.Exit(1)
		}

		printer.Database(database)
	},
}

var databaseInfo = &cobra.Command{
	Use:   "get <databaseID>",
	Short: "get info about a specific managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Database.Get(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting managed database : %v\n", err)
			os.Exit(1)
		}

		printer.Database(s)
	},
}

var databaseUpdate = &cobra.Command{
	Use:     "update <databaseID>",
	Short:   "Update a managed database",
	Aliases: []string{"u"},
	Long:    databaseUpdateLong,
	Example: databaseUpdateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		databaseEngine, _ := cmd.Flags().GetString("database-engine")
		databaseEngineVersion, _ := cmd.Flags().GetString("database-engine-version")
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		label, _ := cmd.Flags().GetString("label")
		tag, _ := cmd.Flags().GetString("tag")
		maintenanceDOW, _ := cmd.Flags().GetString("maintenance-dow")
		maintenanceTime, _ := cmd.Flags().GetString("maintenance-time")
		clusterTimeZone, _ := cmd.Flags().GetString("cluster-time-zone")
		trustedIPs, _ := cmd.Flags().GetStringSlice("trusted-ips")
		mysqlSQLModes, _ := cmd.Flags().GetStringSlice("mysql-sql-modes")
		mysqlRequirePrimaryKey, _ := cmd.Flags().GetBool("mysql-require-primary-key")
		mySQLSlowQueryLog, _ := cmd.Flags().GetBool("mysql-slow-query-log")
		mySQLLongQueryTime, _ := cmd.Flags().GetInt("mysql-long-query-time")
		redisEvictionPolicy, _ := cmd.Flags().GetString("redis-eviction-policy")

		opt := &govultr.DatabaseUpdateReq{
			DatabaseEngine:         databaseEngine,
			DatabaseEngineVersion:  databaseEngineVersion,
			Plan:                   plan,
			Region:                 region,
			Label:                  label,
			Tag:                    tag,
			MaintenanceDOW:         maintenanceDOW,
			MaintenanceTime:        maintenanceTime,
			ClusterTimeZone:        clusterTimeZone,
			TrustedIPs:             trustedIPs,
			MySQLSQLModes:          mysqlSQLModes,
			MySQLRequirePrimaryKey: nil,
			MySQLSlowQueryLog:      nil,
			MySQLLongQueryTime:     mySQLLongQueryTime,
			RedisEvictionPolicy:    redisEvictionPolicy,
		}

		if mysqlRequirePrimaryKey {
			opt.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(true)
		} else if !mysqlRequirePrimaryKey {
			opt.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(false)
		}

		if mySQLSlowQueryLog {
			opt.MySQLSlowQueryLog = govultr.BoolToBoolPtr(true)
		} else if !mySQLSlowQueryLog {
			opt.MySQLSlowQueryLog = govultr.BoolToBoolPtr(false)
		}

		// Make the request
		database, _, err := client.Database.Update(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error updating managed database : %v\n", err)
			os.Exit(1)
		}

		printer.Database(database)
	},
}

var databaseDelete = &cobra.Command{
	Use:     "delete <databaseID>",
	Short:   "delete/destroy a managed database",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Database.Delete(context.Background(), args[0]); err != nil {
			fmt.Printf("error deleting managed database : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted managed database")
	},
}

var databaseUserList = &cobra.Command{
	Use:   "list-users <databaseID>",
	Short: "list all users within the managed databases",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, meta, _, err := client.Database.ListUsers(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of databases : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUserList(s, meta)
	},
}

var databaseUserCreate = &cobra.Command{
	Use:   "create-user <databaseID>",
	Short: "Create a user within a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")

		// Optional
		password, _ := cmd.Flags().GetString("password")
		encryption, _ := cmd.Flags().GetString("encryption")

		opt := &govultr.DatabaseUserCreateReq{
			Username:   username,
			Password:   password,
			Encryption: encryption,
		}

		// Make the request
		databaseUser, _, err := client.Database.CreateUser(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error creating managed database users : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUser(*databaseUser)
	},
}

var databaseUserInfo = &cobra.Command{
	Use:   "get-user <databaseID> <username>",
	Short: "get info about a specific user within a managed databases",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and username")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Database.GetUser(context.TODO(), args[0], args[1])
		if err != nil {
			fmt.Printf("error getting list of databases : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUser(*s)
	},
}

var databaseUserUpdate = &cobra.Command{
	Use:   "update-user <databaseID> <username>",
	Short: "Update a user within a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and username")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")

		opt := &govultr.DatabaseUserUpdateReq{
			Password: password,
		}

		// Make the request
		databaseUser, _, err := client.Database.UpdateUser(context.TODO(), args[0], args[1], opt)
		if err != nil {
			fmt.Printf("error updating managed database user : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUser(*databaseUser)
	},
}

var databaseUserDelete = &cobra.Command{
	Use:   "delete-user <databaseID> <username>",
	Short: "Delete a user within a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and username")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Database.DeleteUser(context.Background(), args[0], args[1]); err != nil {
			fmt.Printf("error deleting managed database user : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted managed database user")
	},
}
