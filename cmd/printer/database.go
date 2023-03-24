package printer

import "github.com/vultr/govultr/v3"

// DatabasePlanList will generate a printer display of available Managed Database plans
func DatabasePlanList(databasePlans []govultr.DatabasePlan, meta *govultr.Meta) {
	for p := range databasePlans {
		display(columns{"ID", databasePlans[p].ID})
		display(columns{"NUMBER OF NODES", databasePlans[p].NumberOfNodes})
		display(columns{"TYPE", databasePlans[p].Type})
		display(columns{"VCPU COUNT", databasePlans[p].VCPUCount})
		display(columns{"RAM", databasePlans[p].RAM})
		display(columns{"DISK", databasePlans[p].Disk})
		display(columns{"MONTHLY COST", databasePlans[p].MonthlyCost})

		display(columns{" "})

		display(columns{"SUPPORTED ENGINES"})
		display(columns{"MYSQL", *databasePlans[p].SupportedEngines.MySQL})
		display(columns{"PG", *databasePlans[p].SupportedEngines.PG})
		display(columns{"REDIS", *databasePlans[p].SupportedEngines.Redis})

		if !*databasePlans[p].SupportedEngines.Redis {
			display(columns{" "})
			display(columns{"MAX CONNECTIONS"})
			display(columns{"MYSQL", databasePlans[p].MaxConnections.MySQL})
			display(columns{"PG", databasePlans[p].MaxConnections.PG})
			display(columns{" "})
		}

		display(columns{"LOCATIONS", databasePlans[p].Locations})
		display(columns{"---------------------------"})
	}

	MetaDBaaS(meta)
	flush()
}

// DatabaseList will generate a printer display of all Managed Databases on the account
func DatabaseList(databases []govultr.Database, meta *govultr.Meta) {
	for d := range databases {
		display(columns{"ID", databases[d].ID})
		display(columns{"DATE CREATED", databases[d].DateCreated})
		display(columns{"PLAN", databases[d].Plan})
		display(columns{"PLAN DISK", databases[d].PlanDisk})
		display(columns{"PLAN RAM", databases[d].PlanRAM})
		display(columns{"PLAN VCPUS", databases[d].PlanVCPUs})
		display(columns{"PLAN REPLICAS", databases[d].PlanReplicas})
		display(columns{"REGION", databases[d].Region})
		display(columns{"STATUS", databases[d].Status})
		display(columns{"LABEL", databases[d].Label})
		display(columns{"TAG", databases[d].Tag})
		display(columns{"DATABASE ENGINE", databases[d].DatabaseEngine})
		display(columns{"DATABASE ENGINE VERSION", databases[d].DatabaseEngineVersion})
		display(columns{"DB NAME", databases[d].DBName})
		display(columns{"HOST", databases[d].Host})
		display(columns{"USER", databases[d].User})
		display(columns{"PASSWORD", databases[d].Password})
		display(columns{"PORT", databases[d].Port})
		display(columns{"MAINTENANCE DOW", databases[d].MaintenanceDOW})
		display(columns{"MAINTENANCE TIME", databases[d].MaintenanceTime})
		display(columns{"LATEST BACKUP", databases[d].LatestBackup})
		display(columns{"TRUSTED IPS", databases[d].TrustedIPs})

		if databases[d].DatabaseEngine == "mysql" {
			display(columns{"MYSQL SQL MODES", databases[d].MySQLSQLModes})
			display(columns{"MYSQL REQUIRE PRIMARY KEY", *databases[d].MySQLRequirePrimaryKey})
			display(columns{"MYSQL SLOW QUERY LOG", *databases[d].MySQLSlowQueryLog})
			if *databases[d].MySQLSlowQueryLog {
				display(columns{"MYSQL LONG QUERY TIME", databases[d].MySQLLongQueryTime})
			}
		}

		if databases[d].DatabaseEngine == "pg" && len(databases[d].PGAvailableExtensions) > 0 {
			display(columns{" "})
			display(columns{"PG AVAILABLE EXTENSIONS"})
			display(columns{"NAME", "VERSIONS"})
			for _, ext := range databases[d].PGAvailableExtensions {
				if len(ext.Versions) > 0 {
					display(columns{ext.Name, ext.Versions})
				} else {
					display(columns{ext.Name, ""})
				}
			}
			display(columns{" "})
		}

		if databases[d].DatabaseEngine == "redis" {
			display(columns{"REDIS EVICTION POLICY", databases[d].RedisEvictionPolicy})
		}

		display(columns{"CLUSTER TIME ZONE", databases[d].ClusterTimeZone})

		if len(databases[d].ReadReplicas) > 0 {
			display(columns{" "})
			display(columns{"READ REPLICAS"})
			for r := range databases[d].ReadReplicas {
				display(columns{"ID", databases[d].ReadReplicas[r].ID})
				display(columns{"DATE CREATED", databases[d].ReadReplicas[r].DateCreated})
				display(columns{"PLAN", databases[d].ReadReplicas[r].Plan})
				display(columns{"PLAN DISK", databases[d].ReadReplicas[r].PlanDisk})
				display(columns{"PLAN RAM", databases[d].ReadReplicas[r].PlanRAM})
				display(columns{"PLAN VCPUS", databases[d].ReadReplicas[r].PlanVCPUs})
				display(columns{"PLAN REPLICAS", databases[d].ReadReplicas[r].PlanReplicas})
				display(columns{"REGION", databases[d].ReadReplicas[r].Region})
				display(columns{"STATUS", databases[d].ReadReplicas[r].Status})
				display(columns{"LABEL", databases[d].ReadReplicas[r].Label})
				display(columns{"TAG", databases[d].ReadReplicas[r].Tag})
				display(columns{"DATABASE ENGINE", databases[d].ReadReplicas[r].DatabaseEngine})
				display(columns{"DATABASE ENGINE VERSION", databases[d].ReadReplicas[r].DatabaseEngineVersion})
				display(columns{"DB NAME", databases[d].ReadReplicas[r].DBName})
				display(columns{"HOST", databases[d].ReadReplicas[r].Host})
				display(columns{"USER", databases[d].ReadReplicas[r].User})
				display(columns{"PASSWORD", databases[d].ReadReplicas[r].Password})
				display(columns{"PORT", databases[d].ReadReplicas[r].Port})
				display(columns{"MAINTENANCE DOW", databases[d].ReadReplicas[r].MaintenanceDOW})
				display(columns{"MAINTENANCE TIME", databases[d].ReadReplicas[r].MaintenanceTime})
				display(columns{"LATEST BACKUP", databases[d].ReadReplicas[r].LatestBackup})
				display(columns{"TRUSTED IPS", databases[d].ReadReplicas[r].TrustedIPs})

				if databases[d].ReadReplicas[r].DatabaseEngine == "mysql" {
					display(columns{"MYSQL SQL MODES", databases[d].ReadReplicas[r].MySQLSQLModes})
					display(columns{"MYSQL REQUIRE PRIMARY KEY", *databases[d].ReadReplicas[r].MySQLRequirePrimaryKey})
					display(columns{"MYSQL SLOW QUERY LOG", *databases[d].ReadReplicas[r].MySQLSlowQueryLog})
					if *databases[d].ReadReplicas[r].MySQLSlowQueryLog {
						display(columns{"MYSQL LONG QUERY TIME", databases[d].ReadReplicas[r].MySQLLongQueryTime})
					}
				}

				if databases[d].ReadReplicas[r].DatabaseEngine == "pg" && len(databases[d].ReadReplicas[r].PGAvailableExtensions) > 0 {
					display(columns{" "})
					display(columns{"PG AVAILABLE EXTENSIONS"})
					display(columns{"NAME", "VERSIONS"})
					for _, ext := range databases[d].ReadReplicas[r].PGAvailableExtensions {
						if len(ext.Versions) > 0 {
							display(columns{ext.Name, ext.Versions})
						} else {
							display(columns{ext.Name, ""})
						}
					}
					display(columns{" "})
				}

				if databases[d].ReadReplicas[r].DatabaseEngine == "redis" {
					display(columns{"REDIS EVICTION POLICY", databases[d].ReadReplicas[r].RedisEvictionPolicy})
				}

				display(columns{"CLUSTER TIME ZONE", databases[d].ReadReplicas[r].ClusterTimeZone})
			}
		}

		display(columns{"---------------------------"})
	}

	MetaDBaaS(meta)
	flush()
}

// Database will generate a printer display of a given Managed Database
func Database(database *govultr.Database) {
	display(columns{"ID", database.ID})
	display(columns{"DATE CREATED", database.DateCreated})
	display(columns{"PLAN", database.Plan})
	display(columns{"PLAN DISK", database.PlanDisk})
	display(columns{"PLAN RAM", database.PlanRAM})
	display(columns{"PLAN VCPUS", database.PlanVCPUs})
	display(columns{"PLAN REPLICAS", database.PlanReplicas})
	display(columns{"REGION", database.Region})
	display(columns{"STATUS", database.Status})
	display(columns{"LABEL", database.Label})
	display(columns{"TAG", database.Tag})
	display(columns{"DATABASE ENGINE", database.DatabaseEngine})
	display(columns{"DATABASE ENGINE VERSION", database.DatabaseEngineVersion})
	display(columns{"DB NAME", database.DBName})
	display(columns{"HOST", database.Host})
	display(columns{"USER", database.User})
	display(columns{"PASSWORD", database.Password})
	display(columns{"PORT", database.Port})
	display(columns{"MAINTENANCE DOW", database.MaintenanceDOW})
	display(columns{"MAINTENANCE TIME", database.MaintenanceTime})
	display(columns{"LATEST BACKUP", database.LatestBackup})
	display(columns{"TRUSTED IPS", database.TrustedIPs})

	if database.DatabaseEngine == "mysql" {
		display(columns{"MYSQL SQL MODES", database.MySQLSQLModes})
		display(columns{"MYSQL REQUIRE PRIMARY KEY", *database.MySQLRequirePrimaryKey})
		display(columns{"MYSQL SLOW QUERY LOG", *database.MySQLSlowQueryLog})
		if *database.MySQLSlowQueryLog {
			display(columns{"MYSQL LONG QUERY TIME", database.MySQLLongQueryTime})
		}
	}

	if database.DatabaseEngine == "pg" && len(database.PGAvailableExtensions) > 0 {
		display(columns{" "})
		display(columns{"PG AVAILABLE EXTENSIONS"})
		display(columns{"NAME", "VERSIONS"})
		for _, ext := range database.PGAvailableExtensions {
			if len(ext.Versions) > 0 {
				display(columns{ext.Name, ext.Versions})
			} else {
				display(columns{ext.Name, ""})
			}
		}
		display(columns{" "})
	}

	if database.DatabaseEngine == "redis" {
		display(columns{"REDIS EVICTION POLICY", database.RedisEvictionPolicy})
	}

	display(columns{"CLUSTER TIME ZONE", database.ClusterTimeZone})

	if len(database.ReadReplicas) > 0 {
		display(columns{" "})
		display(columns{"READ REPLICAS"})
		for r := range database.ReadReplicas {
			display(columns{"ID", database.ReadReplicas[r].ID})
			display(columns{"DATE CREATED", database.ReadReplicas[r].DateCreated})
			display(columns{"PLAN", database.ReadReplicas[r].Plan})
			display(columns{"PLAN DISK", database.ReadReplicas[r].PlanDisk})
			display(columns{"PLAN RAM", database.ReadReplicas[r].PlanRAM})
			display(columns{"PLAN VCPUS", database.ReadReplicas[r].PlanVCPUs})
			display(columns{"PLAN REPLICAS", database.ReadReplicas[r].PlanReplicas})
			display(columns{"REGION", database.ReadReplicas[r].Region})
			display(columns{"STATUS", database.ReadReplicas[r].Status})
			display(columns{"LABEL", database.ReadReplicas[r].Label})
			display(columns{"TAG", database.ReadReplicas[r].Tag})
			display(columns{"DATABASE ENGINE", database.ReadReplicas[r].DatabaseEngine})
			display(columns{"DATABASE ENGINE VERSION", database.ReadReplicas[r].DatabaseEngineVersion})
			display(columns{"DB NAME", database.ReadReplicas[r].DBName})
			display(columns{"HOST", database.ReadReplicas[r].Host})
			display(columns{"USER", database.ReadReplicas[r].User})
			display(columns{"PASSWORD", database.ReadReplicas[r].Password})
			display(columns{"PORT", database.ReadReplicas[r].Port})
			display(columns{"MAINTENANCE DOW", database.ReadReplicas[r].MaintenanceDOW})
			display(columns{"MAINTENANCE TIME", database.ReadReplicas[r].MaintenanceTime})
			display(columns{"LATEST BACKUP", database.ReadReplicas[r].LatestBackup})
			display(columns{"TRUSTED IPS", database.ReadReplicas[r].TrustedIPs})

			if database.ReadReplicas[r].DatabaseEngine == "mysql" {
				display(columns{"MYSQL SQL MODES", database.ReadReplicas[r].MySQLSQLModes})
				display(columns{"MYSQL REQUIRE PRIMARY KEY", *database.ReadReplicas[r].MySQLRequirePrimaryKey})
				display(columns{"MYSQL SLOW QUERY LOG", *database.ReadReplicas[r].MySQLSlowQueryLog})
				if *database.ReadReplicas[r].MySQLSlowQueryLog {
					display(columns{"MYSQL LONG QUERY TIME", database.ReadReplicas[r].MySQLLongQueryTime})
				}
			}

			if database.ReadReplicas[r].DatabaseEngine == "pg" && len(database.ReadReplicas[r].PGAvailableExtensions) > 0 {
				display(columns{" "})
				display(columns{"PG AVAILABLE EXTENSIONS"})
				display(columns{"NAME", "VERSIONS"})
				for _, ext := range database.ReadReplicas[r].PGAvailableExtensions {
					if len(ext.Versions) > 0 {
						display(columns{ext.Name, ext.Versions})
					} else {
						display(columns{ext.Name, ""})
					}
				}
				display(columns{" "})
			}

			if database.ReadReplicas[r].DatabaseEngine == "redis" {
				display(columns{"REDIS EVICTION POLICY", database.ReadReplicas[r].RedisEvictionPolicy})
			}

			display(columns{"CLUSTER TIME ZONE", database.ReadReplicas[r].ClusterTimeZone})
		}
	}

	flush()
}

// DatabaseUserList will generate a printer display of users within a Managed Database
func DatabaseUserList(databaseUsers []govultr.DatabaseUser, meta *govultr.Meta) {
	for u := range databaseUsers {
		display(columns{"USERNAME", databaseUsers[u].Username})
		display(columns{"PASSWORD", databaseUsers[u].Password})
		if databaseUsers[u].Encryption != "" {
			display(columns{"ENCRYPTION", databaseUsers[u].Encryption})
		}
		display(columns{"---------------------------"})
	}

	MetaDBaaS(meta)
	flush()
}

// DatabaseUser will generate a printer display of a given user within a Managed Database
func DatabaseUser(databaseUser govultr.DatabaseUser) {
	display(columns{"USERNAME", databaseUser.Username})
	display(columns{"PASSWORD", databaseUser.Password})
	if databaseUser.Encryption != "" {
		display(columns{"ENCRYPTION", databaseUser.Encryption})
	}

	flush()
}

// DatabaseDBList will generate a printer display of logical databases within a Managed Database cluster
func DatabaseDBList(databaseDBs []govultr.DatabaseDB, meta *govultr.Meta) {
	for u := range databaseDBs {
		display(columns{"NAME", databaseDBs[u].Name})
		display(columns{"---------------------------"})
	}

	MetaDBaaS(meta)
	flush()
}

// DatabaseDB will generate a printer display of a given logical database within a Managed Database cluster
func DatabaseDB(databaseDB govultr.DatabaseDB) {
	display(columns{"NAME", databaseDB.Name})
	flush()
}

// DatabaseUpdates will generate a printer display of available updates for a Managed Database cluster
func DatabaseUpdates(databaseUpdates []string) {
	display(columns{"AVAILABLE UPDATES", databaseUpdates})
	flush()
}

// DatabaseMessage will generate a printer display of a generic information message for a Managed Database cluster
func DatabaseMessage(message string) {
	display(columns{"MESSAGE", message})
	flush()
}

// DatabaseAlertsList will generate a printer display of service alerts for a Managed Database
func DatabaseAlertsList(databaseAlerts []govultr.DatabaseAlert) {
	for a := range databaseAlerts {
		display(columns{"TIMESTAMP", databaseAlerts[a].Timestamp})
		display(columns{"MESSAGE TYPE", databaseAlerts[a].MessageType})
		display(columns{"DESCRIPTION", databaseAlerts[a].Description})

		if databaseAlerts[a].Recommendation != "" {
			display(columns{"RECOMMENDATION", databaseAlerts[a].Recommendation})
		}

		if databaseAlerts[a].MaintenanceScheduled != "" {
			display(columns{"MAINTENANCE SCHEDULED", databaseAlerts[a].MaintenanceScheduled})
		}

		if databaseAlerts[a].ResourceType != "" {
			display(columns{"RESOURCE TYPE", databaseAlerts[a].ResourceType})
		}

		if databaseAlerts[a].TableCount != 0 {
			display(columns{"TABLE COUNT", databaseAlerts[a].TableCount})
		}

		display(columns{"---------------------------"})
	}

	flush()
}

// DatabaseMigrationStatus will generate a printer display of the current migration status of a Managed Database cluster
func DatabaseMigrationStatus(databaseMigration *govultr.DatabaseMigration) {
	display(columns{"STATUS", databaseMigration.Status})

	if databaseMigration.Method != "" {
		display(columns{"METHOD", databaseMigration.Method})
	}

	if databaseMigration.Error != "" {
		display(columns{"ERROR", databaseMigration.Error})
	}

	display(columns{" "})

	display(columns{"CREDENTIALS"})
	display(columns{"HOST", databaseMigration.Credentials.Host})
	display(columns{"PORT", databaseMigration.Credentials.Port})
	display(columns{"USERNAME", databaseMigration.Credentials.Username})
	display(columns{"PASSWORD", databaseMigration.Credentials.Password})

	if databaseMigration.Credentials.Database != "" {
		display(columns{"DATABASE", databaseMigration.Credentials.Database})
	}

	if databaseMigration.Credentials.IgnoredDatabases != "" {
		display(columns{"IGNORED DATABASES", databaseMigration.Credentials.IgnoredDatabases})
	}

	display(columns{"SSL", *databaseMigration.Credentials.SSL})

	flush()
}

// DatabaseBackupInfo will generate a printer display of the latest and oldest backups for a Managed Database cluster
func DatabaseBackupInfo(databaseBackups *govultr.DatabaseBackups) {
	display(columns{"LATEST BACKUP"})
	display(columns{"DATE", databaseBackups.LatestBackup.Date})
	display(columns{"TIME", databaseBackups.LatestBackup.Time})

	display(columns{" "})

	display(columns{"OLDEST BACKUP"})
	display(columns{"DATE", databaseBackups.OldestBackup.Date})
	display(columns{"TIME", databaseBackups.OldestBackup.Time})

	flush()
}

// DatabaseConnectionPoolList will generate a printer display of connection pools within a PostgreSQL Managed Database
func DatabaseConnectionPoolList(databaseConnections *govultr.DatabaseConnections, databaseConnectionPools []govultr.DatabaseConnectionPool, meta *govultr.Meta) {
	display(columns{"CONNECTIONS"})
	display(columns{"USED", databaseConnections.Used})
	display(columns{"AVAILABLE", databaseConnections.Available})
	display(columns{"MAX", databaseConnections.Max})

	display(columns{" "})
	display(columns{"CONNECTION POOLS"})

	for u := range databaseConnectionPools {
		display(columns{"NAME", databaseConnectionPools[u].Name})
		display(columns{"DATABASE", databaseConnectionPools[u].Database})
		display(columns{"USERNAME", databaseConnectionPools[u].Username})
		display(columns{"MODE", databaseConnectionPools[u].Mode})
		display(columns{"SIZE", databaseConnectionPools[u].Size})
		display(columns{"---------------------------"})
	}

	MetaDBaaS(meta)
	flush()
}

// DatabaseConnectionPool will generate a printer display of a given connection pool within a PostgreSQL Managed Database
func DatabaseConnectionPool(databaseConnectionPool govultr.DatabaseConnectionPool) {
	display(columns{"NAME", databaseConnectionPool.Name})
	display(columns{"DATABASE", databaseConnectionPool.Database})
	display(columns{"USERNAME", databaseConnectionPool.Username})
	display(columns{"MODE", databaseConnectionPool.Mode})
	display(columns{"SIZE", databaseConnectionPool.Size})

	flush()
}
