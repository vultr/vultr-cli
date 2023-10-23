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
	databaseListLong    = `Get all databases on your Vultr account`
	databaseListExample = `
	# Full example
	vultr-cli database list

	# Summarized view
	vultr-cli database list --summarize
	`
	databaseCreateLong    = `Create a new Managed Database with specified plan, region, and database engine/version`
	databaseCreateExample = `
	# Full example
	vultr-cli database create --database-engine="mysql" --database-engine-version="8" --region="ewr" \
	    --plan="vultr-dbaas-startup-cc-1-55-2" --label="example-db"

	# Full example with custom MySQL settings
	vultr-cli database create --database-engine="mysql" --database-engine-version="8" --region="ewr" \
	    --plan="vultr-dbaas-startup-cc-1-55-2" --label="example-db" --mysql-slow-query-log="true" --mysql-long-query-time="2"
	`
	databaseUpdateLong    = `Updates a Managed Database with the supplied information`
	databaseUpdateExample = `
	# Full example
	vultr-cli database update --region="sea" --plan="vultr-dbaas-startup-cc-2-80-4"

	# Full example with custom MySQL settings
	vultr-cli database update --mysql-slow-query-log="true" --mysql-long-query-time="2"
	`
)

// Database represents the database command
func Database() *cobra.Command { //nolint:funlen
	databaseCmd := &cobra.Command{
		Use:     "database",
		Short:   "commands to interact with managed databases on vultr",
		Long:    databaseLong,
		Example: databaseExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if auth := cmd.Context().Value("authenticated"); auth != true {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	databaseCmd.AddCommand(databaseList, databaseCreate, databaseInfo, databaseUpdate, databaseDelete)

	// Plan list flags
	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "commands to handle managed database plans",
		Long:  ``,
	}
	planCmd.AddCommand(databasePlanList)
	databasePlanList.Flags().StringP("engine", "e", "", "(optional) Filter by database engine type.")
	databasePlanList.Flags().IntP("nodes", "n", 0, "(optional) Filter by number of nodes.")
	databasePlanList.Flags().StringP("region", "r", "", "(optional) Filter by region.")
	databaseCmd.AddCommand(planCmd)

	// Database list flags
	databaseList.Flags().StringP("label", "l", "", "(optional) Filter by label.")
	databaseList.Flags().StringP("tag", "t", "", "(optional) Filter by tag.")
	databaseList.Flags().StringP("region", "r", "", "(optional) Filter by region.")
	databaseList.Flags().BoolP("summarize", "", false, "(optional) Summarize the list output. One line per database.")

	// Database create flags
	databaseCreate.Flags().StringP("database-engine", "e", "", "database engine for the new manaaged database")
	databaseCreate.Flags().StringP("database-engine-version", "v", "", "database engine version for the new manaaged database")
	databaseCreate.Flags().StringP("region", "r", "", "region id for the new managed database")
	databaseCreate.Flags().StringP("plan", "p", "", "plan id for the new managed database")
	databaseCreate.Flags().StringP("label", "l", "", "label for the new managed database")
	databaseCreate.Flags().StringP("tag", "t", "", "tag for the new managed database")
	databaseCreate.Flags().StringP("vpc-id", "", "", "vpc id for the new managed database")
	databaseCreate.Flags().StringP("maintenance-dow", "", "", "maintenance day of week for the new managed database")
	databaseCreate.Flags().StringP("maintenance-time", "", "", "maintenance time for the new managed database")
	databaseCreate.Flags().StringSliceP("trusted-ips", "", []string{},
		"comma-separated list of trusted ip addresses for the new managed database")
	databaseCreate.Flags().StringSliceP("mysql-sql-modes", "", []string{}, "comma-separated list of sql modes for the new managed database")
	databaseCreate.Flags().BoolP("mysql-require-primary-key", "", true, "enable requiring primary keys for the new mysql managed database")
	databaseCreate.Flags().BoolP("mysql-slow-query-log", "", false, "enable slow query logging for the new mysql managed database")
	databaseCreate.Flags().StringP("mysql-long-query-time", "", "",
		"long query time for the new mysql managed database when slow query logging is enabled")
	databaseCreate.Flags().StringP("redis-eviction-policy", "", "", "eviction policy for the new redis managed database")

	// Database update flags
	databaseUpdate.Flags().StringP("region", "r", "", "region id for the managed database")
	databaseUpdate.Flags().StringP("plan", "p", "", "plan id for the managed database")
	databaseUpdate.Flags().StringP("label", "l", "", "label for the managed database")
	databaseUpdate.Flags().StringP("tag", "t", "", "tag for the managed database")
	databaseUpdate.Flags().StringP("vpc-id", "", "", "vpc id for the managed database")
	databaseUpdate.Flags().StringP("maintenance-dow", "", "", "maintenance day of week for the managed database")
	databaseUpdate.Flags().StringP("maintenance-time", "", "", "maintenance time for the managed database")
	databaseUpdate.Flags().StringP("cluster-time-zone", "", "", "configured time zone for the managed database")
	databaseUpdate.Flags().StringSliceP("trusted-ips", "", []string{}, "comma-separated list of trusted ip addresses for the managed database")
	databaseUpdate.Flags().StringSliceP("mysql-sql-modes", "", []string{}, "comma-separated list of sql modes for the managed database")
	databaseUpdate.Flags().BoolP("mysql-require-primary-key", "", true, "enable requiring primary keys for the mysql managed database")
	databaseUpdate.Flags().BoolP("mysql-slow-query-log", "", false, "enable slow query logging for the mysql managed database")
	databaseUpdate.Flags().StringP("mysql-long-query-time", "", "",
		"long query time for the mysql managed database when slow query logging is enabled")
	databaseUpdate.Flags().StringP("redis-eviction-policy", "", "", "eviction policy for the redis managed database")

	// Database user flags
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "commands to handle managed database users",
		Long:  ``,
	}
	userCmd.AddCommand(databaseUserList, databaseUserCreate, databaseUserInfo, databaseUserUpdate, databaseUserDelete)
	databaseUserCreate.Flags().StringP("username", "u", "", "username for the new manaaged database user")
	databaseUserCreate.Flags().StringP("password", "p", "",
		"password for the new manaaged database user (omit or leave empty to generate a random secure password)")
	databaseUserCreate.Flags().StringP("encryption", "e", "", "encryption type for the new managed database user (MySQL only)")
	databaseUserUpdate.Flags().StringP("password", "p", "",
		"new password for the manaaged database user (leave empty to generate a random secure password)")
	databaseCmd.AddCommand(userCmd)

	// Database logical db flags
	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "commands to handle managed database logical dbs",
		Long:  ``,
	}
	dbCmd.AddCommand(databaseDBList, databaseDBCreate, databaseDBInfo, databaseDBDelete)
	databaseDBCreate.Flags().StringP("name", "n", "", "name of the new logical database within the manaaged database")
	databaseCmd.AddCommand(dbCmd)

	// Database maintenance update commands
	maintenanceCmd := &cobra.Command{
		Use:   "maintenance",
		Short: "commands to handle managed database maintenance updates",
		Long:  ``,
	}
	maintenanceCmd.AddCommand(databaseMaintenanceUpdatesList, databaseStartMaintenance)
	databaseCmd.AddCommand(maintenanceCmd)

	// Database alerts flags
	alertCmd := &cobra.Command{
		Use:   "alert",
		Short: "commands to handle managed database alerts",
		Long:  ``,
	}
	alertCmd.AddCommand(databaseAlertsList)
	databaseAlertsList.Flags().StringP("period", "p", "", "period (day, week, month, year) for viewing service alerts for a manaaged database")
	databaseCmd.AddCommand(alertCmd)

	// Database migration flags
	migrationsCmd := &cobra.Command{
		Use:   "migration",
		Short: "commands to handle managed database migrations",
		Long:  ``,
	}
	migrationsCmd.AddCommand(databaseMigrationStatus, databaseStartMigration, databaseDetachMigration)
	databaseStartMigration.Flags().StringP("host", "", "", "source host for the manaaged database migration")
	databaseStartMigration.Flags().IntP("port", "", 0, "source port for the manaaged database migration")
	databaseStartMigration.Flags().StringP("username", "", "",
		"source username for the manaaged database migration (uses `default` for Redis if omitted)")
	databaseStartMigration.Flags().StringP("password", "", "", "source password for the manaaged database migration")
	databaseStartMigration.Flags().StringP("database", "", "", "source database for the manaaged database migration (MySQL/PostgreSQL only)")
	databaseStartMigration.Flags().StringP("ignored-dbs", "", "",
		"comma-separated list of ignored databases for the manaaged database migration (MySQL/PostgreSQL only)")
	databaseStartMigration.Flags().BoolP("ssl", "", true, "source ssl requirement for the manaaged database migration")
	databaseCmd.AddCommand(migrationsCmd)

	// Database read replica flags
	readReplicaCmd := &cobra.Command{
		Use:   "read-replica",
		Short: "commands to handle managed database read replicas",
		Long:  ``,
	}
	readReplicaCmd.AddCommand(databaseAddReadReplica)
	databaseAddReadReplica.Flags().StringP("region", "r", "", "region id for the new managed database read replica")
	databaseAddReadReplica.Flags().StringP("label", "l", "", "label for the new managed database read replica")
	databaseCmd.AddCommand(readReplicaCmd)

	// Database backup and restore flags
	backupsCmd := &cobra.Command{
		Use:   "backup",
		Short: "commands to handle managed database backups, restorations, and forks",
		Long:  ``,
	}
	backupsCmd.AddCommand(databaseGetBackupInfo, databaseRestoreFromBackup, databaseFork)
	databaseRestoreFromBackup.Flags().StringP("label", "", "", "label for the new managed database restored from backup")
	databaseRestoreFromBackup.Flags().StringP("type", "", "",
		"restoration type: `pitr` for point-in-time recovery or `basebackup` for latest backup (default)")
	databaseRestoreFromBackup.Flags().StringP("date", "", "", "backup date to use for point-in-time recovery")
	databaseRestoreFromBackup.Flags().StringP("time", "", "", "backup time to use for point-in-time recovery")
	databaseFork.Flags().StringP("label", "", "", "label for the new managed database forked from the backup")
	databaseFork.Flags().StringP("region", "", "", "region id for the new managed database forked from the backup")
	databaseFork.Flags().StringP("plan", "", "", "plan id for the new managed database forked from the backup")
	databaseFork.Flags().StringP("type", "", "",
		"restoration type: `pitr` for point-in-time recovery or `basebackup` for latest backup (default)")
	databaseFork.Flags().StringP("date", "", "", "backup date to use for point-in-time recovery")
	databaseFork.Flags().StringP("time", "", "", "backup time to use for point-in-time recovery")
	databaseCmd.AddCommand(backupsCmd)

	// Database connection pools flags
	connectionPoolsCmd := &cobra.Command{
		Use:   "connection-pool",
		Short: "commands to handle PostgreSQL managed database connection pools",
		Long:  ``,
	}

	connectionPoolsCmd.AddCommand(
		databaseConnectionPoolList,
		databaseConnectionPoolCreate,
		databaseConnectionPoolInfo,
		databaseConnectionPoolUpdate,
		databaseConnectionPoolDelete,
	)

	databaseConnectionPoolCreate.Flags().StringP("name", "n", "", "name for the new managed database connection pool")
	databaseConnectionPoolCreate.Flags().StringP("database", "d", "", "database for the new managed database connection pool")
	databaseConnectionPoolCreate.Flags().StringP("username", "u", "", "username for the new managed database connection pool")
	databaseConnectionPoolCreate.Flags().StringP("mode", "m", "", "mode for the new managed database connection pool")
	databaseConnectionPoolCreate.Flags().IntP("size", "s", 0, "size for the new managed database connection pool")
	databaseConnectionPoolUpdate.Flags().StringP("database", "d", "", "database for the managed database connection pool")
	databaseConnectionPoolUpdate.Flags().StringP("username", "u", "", "username for the managed database connection pool")
	databaseConnectionPoolUpdate.Flags().StringP("mode", "m", "", "mode for the managed database connection pool")
	databaseConnectionPoolUpdate.Flags().IntP("size", "s", 0, "size for the managed database connection pool")
	databaseCmd.AddCommand(connectionPoolsCmd)

	// Database PostgreSQL advanced option flags
	advancedOptionsCmd := &cobra.Command{
		Use:   "advanced-option",
		Short: "commands to handle PostgreSQL managed database advanced options",
		Long:  ``,
	}
	advancedOptionsCmd.AddCommand(databaseAdvancedOptionsList, databaseAdvancedOptionsUpdate)
	databaseAdvancedOptionsUpdate.Flags().Float32P("autovacuum-analyze-scale-factor", "", 0,
		"set the managed postgresql configuration value for autovacuum_analyze_scale_factor")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-analyze-threshold", "", 0,
		"set the managed postgresql configuration value for autovacuum_analyze_threshold")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-freeze-max-age", "", 0,
		"set the managed postgresql configuration value for autovacuum_freeze_max_age")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-max-workers", "", 0,
		"set the managed postgresql configuration value for autovacuum_max_workers")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-naptime", "", 0,
		"set the managed postgresql configuration value for autovacuum_naptime")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-vacuum-cost-delay", "", 0,
		"set the managed postgresql configuration value for autovacuum_vacuum_cost_delay")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-vacuum-cost-limit", "", 0,
		"set the managed postgresql configuration value for autovacuum_vacuum_cost_limit")
	databaseAdvancedOptionsUpdate.Flags().Float32P("autovacuum-vacuum-scale-factor", "", 0,
		"set the managed postgresql configuration value for autovacuum_vacuum_scale_factor")
	databaseAdvancedOptionsUpdate.Flags().IntP("autovacuum-vacuum-threshold", "", 0,
		"set the managed postgresql configuration value for autovacuum_vacuum_threshold")
	databaseAdvancedOptionsUpdate.Flags().IntP("bgwriter-delay", "", 0,
		"set the managed postgresql configuration value for bgwriter_delay")
	databaseAdvancedOptionsUpdate.Flags().IntP("bgwriter-flush-after", "", 0,
		"set the managed postgresql configuration value for bgwriter_flush_after")
	databaseAdvancedOptionsUpdate.Flags().IntP("bgwriter-lru-maxpages", "", 0,
		"set the managed postgresql configuration value for bgwriter_lru_maxpages")
	databaseAdvancedOptionsUpdate.Flags().Float32P("bgwriter-lru-multiplier", "", 0,
		"set the managed postgresql configuration value for bgwriter_lru_multiplier")
	databaseAdvancedOptionsUpdate.Flags().IntP("deadlock-timeout", "", 0,
		"set the managed postgresql configuration value for deadlock_timeout")
	databaseAdvancedOptionsUpdate.Flags().StringP("default-toast-compression", "", "",
		"set the managed postgresql configuration value for default_toast_compression")
	databaseAdvancedOptionsUpdate.Flags().IntP("idle-in-transaction-session-timeout", "", 0,
		"set the managed postgresql configuration value for idle_in_transaction_session_timeout")
	databaseAdvancedOptionsUpdate.Flags().BoolP("jit", "", false,
		"set the managed postgresql configuration value for jit")
	databaseAdvancedOptionsUpdate.Flags().IntP("log-autovacuum-min-duration", "", 0,
		"set the managed postgresql configuration value for log_autovacuum_min_duration")
	databaseAdvancedOptionsUpdate.Flags().StringP("log-error-verbosity", "", "",
		"set the managed postgresql configuration value for log_error_verbosity")
	databaseAdvancedOptionsUpdate.Flags().StringP("log-line-prefix", "", "",
		"set the managed postgresql configuration value for log_line_prefix")
	databaseAdvancedOptionsUpdate.Flags().IntP("log-min-duration-statement", "", 0,
		"set the managed postgresql configuration value for log_min_duration_statement")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-files-per-process", "", 0,
		"set the managed postgresql configuration value for max_files_per_process")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-locks-per-transaction", "", 0,
		"set the managed postgresql configuration value for max_locks_per_transaction")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-logical-replication-workers", "", 0,
		"set the managed postgresql configuration value for max_logical_replication_workers")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-parallel-workers", "", 0,
		"set the managed postgresql configuration value for max_parallel_workers")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-parallel-workers-per-gather", "", 0,
		"set the managed postgresql configuration value for max_parallel_workers_per_gather")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-pred-locks-per-transaction", "", 0,
		"set the managed postgresql configuration value for max_pred_locks_per_transaction")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-prepared-transactions", "", 0,
		"set the managed postgresql configuration value for max_prepared_transactions")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-replication-slots", "", 0,
		"set the managed postgresql configuration value for max_replication_slots")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-stack-depth", "", 0,
		"set the managed postgresql configuration value for max_stack_depth")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-standby-archive-delay", "", 0,
		"set the managed postgresql configuration value for max_standby_archive_delay")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-standby-streaming-delay", "", 0,
		"set the managed postgresql configuration value for max_standby_streaming_delay")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-wal-senders", "", 0,
		"set the managed postgresql configuration value for max_wal_senders")
	databaseAdvancedOptionsUpdate.Flags().IntP("max-worker-processes", "", 0,
		"set the managed postgresql configuration value for max_worker_processes")
	databaseAdvancedOptionsUpdate.Flags().IntP("pg-partman-bgw-interval", "", 0,
		"set the managed postgresql configuration value for pg_partman_bgw.interval")
	databaseAdvancedOptionsUpdate.Flags().StringP("pg-partman-bgw-role", "", "",
		"set the managed postgresql configuration value for pg_partman_bgw.role")
	databaseAdvancedOptionsUpdate.Flags().StringP("pg-stat-statements-track", "", "",
		"set the managed postgresql configuration value for pg_stat_statements.track")
	databaseAdvancedOptionsUpdate.Flags().IntP("temp-file-limit", "", 0,
		"set the managed postgresql configuration value for temp_file_limit")
	databaseAdvancedOptionsUpdate.Flags().IntP("track-activity-query-size", "", 0,
		"set the managed postgresql configuration value for track_activity_query_size")
	databaseAdvancedOptionsUpdate.Flags().StringP("track-commit-timestamp", "", "",
		"set the managed postgresql configuration value for track_commit_timestamp")
	databaseAdvancedOptionsUpdate.Flags().StringP("track-functions", "", "",
		"set the managed postgresql configuration value for track_functions")
	databaseAdvancedOptionsUpdate.Flags().StringP("track-io-timing", "", "",
		"set the managed postgresql configuration value for track_io_timing")
	databaseAdvancedOptionsUpdate.Flags().IntP("wal-sender-timeout", "", 0,
		"set the managed postgresql configuration value for wal_sender_timeout")
	databaseAdvancedOptionsUpdate.Flags().IntP("wal-writer-delay", "", 0,
		"set the managed postgresql configuration value for wal_writer_delay")
	databaseCmd.AddCommand(advancedOptionsCmd)

	// Database version upgrade flags
	versionUpgradeCmd := &cobra.Command{
		Use:   "version",
		Short: "commands to handle managed database version upgrades",
		Long:  ``,
	}
	versionUpgradeCmd.AddCommand(databaseAvailableVersionsList, databaseStartVersionUpgrade)
	databaseStartVersionUpgrade.Flags().StringP("version", "v", "", "version of the manaaged database to upgrade to")
	databaseCmd.AddCommand(versionUpgradeCmd)

	return databaseCmd
}

var databasePlanList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
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
	Long:    databaseListLong,
	Example: databaseListExample,
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")
		tag, _ := cmd.Flags().GetString("tag")
		region, _ := cmd.Flags().GetString("region")
		options := &govultr.DBListOptions{
			Label:  label,
			Tag:    tag,
			Region: region,
		}
		summarize, _ := cmd.Flags().GetBool("summarize")

		s, meta, _, err := client.Database.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting list of databases : %v\n", err)
			os.Exit(1)
		}

		if summarize {
			printer.DatabaseListSummary(s, meta)
		} else {
			printer.DatabaseList(s, meta)
		}
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
		vpc, _ := cmd.Flags().GetString("vpc-id")
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
			Region:                 region,
			Plan:                   plan,
			Label:                  label,
			Tag:                    tag,
			VPCID:                  vpc,
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
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		label, _ := cmd.Flags().GetString("label")
		tag, _ := cmd.Flags().GetString("tag")
		vpc, _ := cmd.Flags().GetString("vpc-id")
		vpcSet := cmd.Flags().Lookup("vpc-id").Changed
		maintenanceDOW, _ := cmd.Flags().GetString("maintenance-dow")
		maintenanceTime, _ := cmd.Flags().GetString("maintenance-time")
		clusterTimeZone, _ := cmd.Flags().GetString("cluster-time-zone")
		trustedIPs, _ := cmd.Flags().GetStringSlice("trusted-ips")
		mysqlSQLModes, _ := cmd.Flags().GetStringSlice("mysql-sql-modes")
		mysqlRequirePrimaryKey, _ := cmd.Flags().GetBool("mysql-require-primary-key")
		mysqlRequirePrimaryKeySet := cmd.Flags().Lookup("mysql-require-primary-key").Changed
		mySQLSlowQueryLog, _ := cmd.Flags().GetBool("mysql-slow-query-log")
		mySQLSlowQueryLogSet := cmd.Flags().Lookup("mysql-slow-query-log").Changed
		mySQLLongQueryTime, _ := cmd.Flags().GetInt("mysql-long-query-time")
		redisEvictionPolicy, _ := cmd.Flags().GetString("redis-eviction-policy")

		opt := &govultr.DatabaseUpdateReq{
			Region:              region,
			Plan:                plan,
			Label:               label,
			Tag:                 tag,
			MaintenanceDOW:      maintenanceDOW,
			MaintenanceTime:     maintenanceTime,
			ClusterTimeZone:     clusterTimeZone,
			TrustedIPs:          trustedIPs,
			MySQLSQLModes:       mysqlSQLModes,
			MySQLLongQueryTime:  mySQLLongQueryTime,
			RedisEvictionPolicy: redisEvictionPolicy,
		}

		if vpcSet {
			opt.VPCID = govultr.StringToStringPtr(vpc)
		}

		if mysqlRequirePrimaryKeySet && mysqlRequirePrimaryKey {
			opt.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(true)
		} else if mysqlRequirePrimaryKeySet && !mysqlRequirePrimaryKey {
			opt.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(false)
		}

		if mySQLSlowQueryLogSet && mySQLSlowQueryLog {
			opt.MySQLSlowQueryLog = govultr.BoolToBoolPtr(true)
		} else if mySQLSlowQueryLogSet && !mySQLSlowQueryLog {
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
	Use:   "list <databaseID>",
	Short: "list all users within a managed database",
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
			fmt.Printf("error getting list of database users : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUserList(s, meta)
	},
}

var databaseUserCreate = &cobra.Command{
	Use:   "create <databaseID>",
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
			fmt.Printf("error creating managed database user : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUser(*databaseUser)
	},
}

var databaseUserInfo = &cobra.Command{
	Use:   "get <databaseID> <username>",
	Short: "get info about a specific user within a managed database",
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
			fmt.Printf("error getting managed database user : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUser(*s)
	},
}

var databaseUserUpdate = &cobra.Command{
	Use:   "update <databaseID> <username>",
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
	Use:   "delete <databaseID> <username>",
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

var databaseDBList = &cobra.Command{
	Use:   "list <databaseID>",
	Short: "list all logical databases within a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, meta, _, err := client.Database.ListDBs(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of logical databases : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseDBList(s, meta)
	},
}

var databaseDBCreate = &cobra.Command{
	Use:   "create <databaseID>",
	Short: "Create a logical database within a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		opt := &govultr.DatabaseDBCreateReq{
			Name: name,
		}

		// Make the request
		databaseDB, _, err := client.Database.CreateDB(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error creating logical database : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseDB(*databaseDB)
	},
}

var databaseDBInfo = &cobra.Command{
	Use:   "get <databaseID> <dbname>",
	Short: "get info about a specific logicla database within a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and dbname")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Database.GetDB(context.TODO(), args[0], args[1])
		if err != nil {
			fmt.Printf("error getting logical database : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseDB(*s)
	},
}

var databaseDBDelete = &cobra.Command{
	Use:   "delete <databaseID> <dbname>",
	Short: "Delete a logical database within a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and dbname")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Database.DeleteDB(context.Background(), args[0], args[1]); err != nil {
			fmt.Printf("error deleting logical database : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted logical database")
	},
}

var databaseMaintenanceUpdatesList = &cobra.Command{
	Use:   "list <databaseID>",
	Short: "list all available maintenance updates for a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Database.ListMaintenanceUpdates(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of available maintenance updates : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseUpdates(s)
	},
}

var databaseStartMaintenance = &cobra.Command{
	Use:   "start <databaseID>",
	Short: "Initialize maintenance updates for a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		message, _, err := client.Database.StartMaintenance(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error starting maintenance updates : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseMessage(message)
	},
}

var databaseAlertsList = &cobra.Command{
	Use:   "list <databaseID>",
	Short: "List service alerts for a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		period, _ := cmd.Flags().GetString("period")

		opt := &govultr.DatabaseListAlertsReq{
			Period: period,
		}

		// Make the request
		databaseAlerts, _, err := client.Database.ListServiceAlerts(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error listing managed database alerts : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseAlertsList(databaseAlerts)
	},
}

var databaseMigrationStatus = &cobra.Command{
	Use:   "get <databaseID>",
	Short: "Get the current migration status of a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		databaseMigration, _, err := client.Database.GetMigrationStatus(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error retrieving managed database migration status : %v\n", err)
			os.Exit(1)
		}

		if databaseMigration == nil {
			printer.DatabaseMessage("There is currently no active migration configured for this Managed Database.")
		} else {
			printer.DatabaseMigrationStatus(databaseMigration)
		}
	},
}

var databaseStartMigration = &cobra.Command{
	Use:   "start <databaseID>",
	Short: "Start a migration for a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		IgnoredDatabases, _ := cmd.Flags().GetString("ignored-dbs")
		ssl, _ := cmd.Flags().GetBool("ssl")

		opt := &govultr.DatabaseMigrationStartReq{
			Host:             host,
			Port:             port,
			Username:         username,
			Password:         password,
			Database:         database,
			IgnoredDatabases: IgnoredDatabases,
			SSL:              nil,
		}

		if ssl {
			opt.SSL = govultr.BoolToBoolPtr(true)
		} else if !ssl {
			opt.SSL = govultr.BoolToBoolPtr(false)
		}

		// Make the request
		databaseMigration, _, err := client.Database.StartMigration(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error starting migration : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseMigrationStatus(databaseMigration)
	},
}

var databaseDetachMigration = &cobra.Command{
	Use:   "detach <databaseID>",
	Short: "Detach a migration from a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Database.DetachMigration(context.Background(), args[0]); err != nil {
			fmt.Printf("error deleting managed database user : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Detached migration")
	},
}

var databaseAddReadReplica = &cobra.Command{
	Use:   "create <databaseID>",
	Short: "Add a read-only replica to a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		label, _ := cmd.Flags().GetString("label")

		opt := &govultr.DatabaseAddReplicaReq{
			Region: region,
			Label:  label,
		}

		// Make the request
		database, _, err := client.Database.AddReadOnlyReplica(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error creating read-only replica : %v\n", err)
			os.Exit(1)
		}

		printer.Database(database)
	},
}

var databaseGetBackupInfo = &cobra.Command{
	Use:   "get <databaseID>",
	Short: "Get the latest and oldest available backups for a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		databaseBackups, _, err := client.Database.GetBackupInformation(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error retrieving backup information : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseBackupInfo(databaseBackups)
	},
}

var databaseRestoreFromBackup = &cobra.Command{
	Use:   "restore <databaseID>",
	Short: "Create a new managed database from a backup",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")
		rtype, _ := cmd.Flags().GetString("type")
		date, _ := cmd.Flags().GetString("date")
		time, _ := cmd.Flags().GetString("time")

		opt := &govultr.DatabaseBackupRestoreReq{
			Label: label,
			Type:  rtype,
			Date:  date,
			Time:  time,
		}

		// Make the request
		database, _, err := client.Database.RestoreFromBackup(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error creating managed database from backup : %v\n", err)
			os.Exit(1)
		}

		printer.Database(database)
	},
}

var databaseFork = &cobra.Command{
	Use:   "fork <databaseID>",
	Short: "Fork a Managed Database to a new subscription from a backup.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		rtype, _ := cmd.Flags().GetString("type")
		date, _ := cmd.Flags().GetString("date")
		time, _ := cmd.Flags().GetString("time")

		opt := &govultr.DatabaseForkReq{
			Label:  label,
			Region: region,
			Plan:   plan,
			Type:   rtype,
			Date:   date,
			Time:   time,
		}

		// Make the request
		database, _, err := client.Database.Fork(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error forking managed database : %v\n", err)
			os.Exit(1)
		}

		printer.Database(database)
	},
}

var databaseConnectionPoolList = &cobra.Command{
	Use:   "list <databaseID>",
	Short: "list all connection pools within a PostgreSQL managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c, s, meta, _, err := client.Database.ListConnectionPools(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of connection pools : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseConnectionPoolList(c, s, meta)
	},
}

var databaseConnectionPoolCreate = &cobra.Command{
	Use:   "create <databaseID>",
	Short: "Create a connection pool within a PostgreSQL managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		database, _ := cmd.Flags().GetString("database")
		username, _ := cmd.Flags().GetString("username")
		mode, _ := cmd.Flags().GetString("mode")
		size, _ := cmd.Flags().GetInt("size")

		opt := &govultr.DatabaseConnectionPoolCreateReq{
			Name:     name,
			Database: database,
			Username: username,
			Mode:     mode,
			Size:     size,
		}

		// Make the request
		databaseConnectionPool, _, err := client.Database.CreateConnectionPool(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error creating managed database connection pool : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseConnectionPool(*databaseConnectionPool)
	},
}

var databaseConnectionPoolInfo = &cobra.Command{
	Use:   "get <databaseID> <poolName>",
	Short: "get info about a specific connection pool within a PostgreSQL managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and poolName")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Database.GetConnectionPool(context.TODO(), args[0], args[1])
		if err != nil {
			fmt.Printf("error getting managed database connection pool : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseConnectionPool(*s)
	},
}

var databaseConnectionPoolUpdate = &cobra.Command{
	Use:   "update <databaseID> <poolName>",
	Short: "Update a connection pool within a PostgreSQL managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and poolName")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		database, _ := cmd.Flags().GetString("database")
		username, _ := cmd.Flags().GetString("username")
		mode, _ := cmd.Flags().GetString("mode")
		size, _ := cmd.Flags().GetInt("size")

		opt := &govultr.DatabaseConnectionPoolUpdateReq{
			Database: database,
			Username: username,
			Mode:     mode,
			Size:     size,
		}

		// Make the request
		databaseConnectionPool, _, err := client.Database.UpdateConnectionPool(context.TODO(), args[0], args[1], opt)
		if err != nil {
			fmt.Printf("error updating managed database connection pool : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseConnectionPool(*databaseConnectionPool)
	},
}

var databaseConnectionPoolDelete = &cobra.Command{
	Use:   "delete <databaseID> <poolName>",
	Short: "Delete a connection pool within a PostgreSQL managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a databaseID and poolName")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Database.DeleteConnectionPool(context.Background(), args[0], args[1]); err != nil {
			fmt.Printf("error deleting managed database connection pool : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted managed database connection pool")
	},
}

var databaseAdvancedOptionsList = &cobra.Command{
	Use:   "list <databaseID>",
	Short: "list all available and configured advanced options for a PostgreSQL managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c, a, _, err := client.Database.ListAdvancedOptions(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of advanced options : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseAdvancedOptions(c, a)
	},
}

var databaseAdvancedOptionsUpdate = &cobra.Command{
	Use:   "update <databaseID>",
	Short: "Configure advanced options for a PostgreSQL managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		autovacuumAnalyzeScaleFactor, _ := cmd.Flags().GetFloat32("autovacuum-analyze-scale-factor")
		autovacuumAnalyzeThreshold, _ := cmd.Flags().GetInt("autovacuum-analyze-threshold")
		autovacuumFreezeMaxAge, _ := cmd.Flags().GetInt("autovacuum-freeze-max-age")
		autovacuumMaxWorkers, _ := cmd.Flags().GetInt("autovacuum-max-workers")
		autovacuumNaptime, _ := cmd.Flags().GetInt("autovacuum-naptime")
		autovacuumVacuumCostDelay, _ := cmd.Flags().GetInt("autovacuum-vacuum-cost-delay")
		autovacuumVacuumCostLimit, _ := cmd.Flags().GetInt("autovacuum-vacuum-cost-limit")
		autovacuumVacuumScaleFactor, _ := cmd.Flags().GetFloat32("autovacuum-vacuum-scale-factor")
		autovacuumVacuumThreshold, _ := cmd.Flags().GetInt("autovacuum-vacuum-threshold")
		bgwriterDelay, _ := cmd.Flags().GetInt("bgwriter-delay")
		bgwriterFlushAfter, _ := cmd.Flags().GetInt("bgwriter-flush-after")
		bgwriterLruMaxpages, _ := cmd.Flags().GetInt("bgwriter-lru-maxpages")
		bgwriterLruMultiplier, _ := cmd.Flags().GetFloat32("bgwriter-lru-multiplier")
		deadlockTimeout, _ := cmd.Flags().GetInt("deadlock-timeout")
		defaultToastCompression, _ := cmd.Flags().GetString("default-toast-compression")
		idleInTransactionSessionTimeout, _ := cmd.Flags().GetInt("idle-in-transaction-session-timeout")
		jit, _ := cmd.Flags().GetBool("jit")
		logAutovacuumMinDuration, _ := cmd.Flags().GetInt("log-autovacuum-min-duration")
		logErrorVerbosity, _ := cmd.Flags().GetString("log-error-verbosity")
		logLinePrefix, _ := cmd.Flags().GetString("log-line-prefix")
		logMinDurationStatement, _ := cmd.Flags().GetInt("log-min-duration-statement")
		maxFilesPerProcess, _ := cmd.Flags().GetInt("max-files-per-process")
		maxLocksPerTransaction, _ := cmd.Flags().GetInt("max-locks-per-transaction")
		maxLogicalReplicationWorkers, _ := cmd.Flags().GetInt("max-logical-replication-workers")
		maxParallelWorkers, _ := cmd.Flags().GetInt("max-parallel-workers")
		maxParallelWorkersPerGather, _ := cmd.Flags().GetInt("max-parallel-workers-per-gather")
		maxPredLocksPerTransaction, _ := cmd.Flags().GetInt("max-pred-locks-per-transaction")
		maxPreparedTransactions, _ := cmd.Flags().GetInt("max-prepared-transactions")
		maxReplicationSlots, _ := cmd.Flags().GetInt("max-replication-slots")
		maxStackDepth, _ := cmd.Flags().GetInt("max-stack-depth")
		maxStandbyArchiveDelay, _ := cmd.Flags().GetInt("max-standby-archive-delay")
		maxStandbyStreamingDelay, _ := cmd.Flags().GetInt("max-standby-streaming-delay")
		maxWalSenders, _ := cmd.Flags().GetInt("max-wal-senders")
		maxWorkerProcesses, _ := cmd.Flags().GetInt("max-worker-processes")
		pgPartmanBGWInterval, _ := cmd.Flags().GetInt("pg-partman-bgw-interval")
		pgPartmanBGWRole, _ := cmd.Flags().GetString("pg-partman-bgw-role")
		pgStatStatementsTrack, _ := cmd.Flags().GetString("pg-stat-statements-track")
		tempFileLimit, _ := cmd.Flags().GetInt("temp-file-limit")
		trackActivityQuerySize, _ := cmd.Flags().GetInt("track-activity-query-size")
		trackCommitTimestamp, _ := cmd.Flags().GetString("track-commit-timestamp")
		trackFunctions, _ := cmd.Flags().GetString("track-functions")
		trackIOTiming, _ := cmd.Flags().GetString("track-io-timing")
		walSenderTimeout, _ := cmd.Flags().GetInt("wal-sender-timeout")
		walWriterDelay, _ := cmd.Flags().GetInt("wal-writer-delay")

		opt := &govultr.DatabaseAdvancedOptions{
			AutovacuumAnalyzeScaleFactor:    autovacuumAnalyzeScaleFactor,
			AutovacuumAnalyzeThreshold:      autovacuumAnalyzeThreshold,
			AutovacuumFreezeMaxAge:          autovacuumFreezeMaxAge,
			AutovacuumMaxWorkers:            autovacuumMaxWorkers,
			AutovacuumNaptime:               autovacuumNaptime,
			AutovacuumVacuumCostDelay:       autovacuumVacuumCostDelay,
			AutovacuumVacuumCostLimit:       autovacuumVacuumCostLimit,
			AutovacuumVacuumScaleFactor:     autovacuumVacuumScaleFactor,
			AutovacuumVacuumThreshold:       autovacuumVacuumThreshold,
			BGWRITERDelay:                   bgwriterDelay,
			BGWRITERFlushAFter:              bgwriterFlushAfter,
			BGWRITERLRUMaxPages:             bgwriterLruMaxpages,
			BGWRITERLRUMultiplier:           bgwriterLruMultiplier,
			DeadlockTimeout:                 deadlockTimeout,
			DefaultToastCompression:         defaultToastCompression,
			IdleInTransactionSessionTimeout: idleInTransactionSessionTimeout,
			Jit:                             nil,
			LogAutovacuumMinDuration:        logAutovacuumMinDuration,
			LogErrorVerbosity:               logErrorVerbosity,
			LogLinePrefix:                   logLinePrefix,
			LogMinDurationStatement:         logMinDurationStatement,
			MaxFilesPerProcess:              maxFilesPerProcess,
			MaxLocksPerTransaction:          maxLocksPerTransaction,
			MaxLogicalReplicationWorkers:    maxLogicalReplicationWorkers,
			MaxParallelWorkers:              maxParallelWorkers,
			MaxParallelWorkersPerGather:     maxParallelWorkersPerGather,
			MaxPredLocksPerTransaction:      maxPredLocksPerTransaction,
			MaxPreparedTransactions:         maxPreparedTransactions,
			MaxReplicationSlots:             maxReplicationSlots,
			MaxStackDepth:                   maxStackDepth,
			MaxStandbyArchiveDelay:          maxStandbyArchiveDelay,
			MaxStandbyStreamingDelay:        maxStandbyStreamingDelay,
			MaxWalSenders:                   maxWalSenders,
			MaxWorkerProcesses:              maxWorkerProcesses,
			PGPartmanBGWInterval:            pgPartmanBGWInterval,
			PGPartmanBGWRole:                pgPartmanBGWRole,
			PGStateStatementsTrack:          pgStatStatementsTrack,
			TempFileLimit:                   tempFileLimit,
			TrackActivityQuerySize:          trackActivityQuerySize,
			TrackCommitTimestamp:            trackCommitTimestamp,
			TrackFunctions:                  trackFunctions,
			TrackIOTiming:                   trackIOTiming,
			WALSenderTImeout:                walSenderTimeout,
			WALWriterDelay:                  walWriterDelay,
		}

		if jit {
			opt.Jit = govultr.BoolToBoolPtr(true)
		} else if !jit {
			opt.Jit = govultr.BoolToBoolPtr(false)
		}

		// Make the request
		c, a, _, err := client.Database.UpdateAdvancedOptions(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error starting migration : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseAdvancedOptions(c, a)
	},
}

var databaseAvailableVersionsList = &cobra.Command{
	Use:   "list <databaseID>",
	Short: "list all available version upgrades for a managed database",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, _, err := client.Database.ListAvailableVersions(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting list of available version upgrades : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseAvailableVersions(s)
	},
}

var databaseStartVersionUpgrade = &cobra.Command{
	Use:   "upgrade <databaseID>",
	Short: "Initialize a version upgrade for a managed database",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a databaseID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetString("version")

		opt := &govultr.DatabaseVersionUpgradeReq{
			Version: version,
		}

		message, _, err := client.Database.StartVersionUpgrade(context.TODO(), args[0], opt)
		if err != nil {
			fmt.Printf("error starting version upgrade : %v\n", err)
			os.Exit(1)
		}

		printer.DatabaseMessage(message)
	},
}
