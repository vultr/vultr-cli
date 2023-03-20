package printer

import "github.com/vultr/govultr/v3"

func DatabasePlanList(databasePlans []govultr.DatabasePlan, meta *govultr.Meta) {
	for _, p := range databasePlans {
		display(columns{"ID", p.ID})
		display(columns{"NUMBER OF NODES", p.NumberOfNodes})
		display(columns{"TYPE", p.Type})
		display(columns{"VCPU COUNT", p.VCPUCount})
		display(columns{"RAM", p.RAM})
		display(columns{"DISK", p.Disk})
		display(columns{"MONTHLY COST", p.MonthlyCost})

		display(columns{" "})

		display(columns{"SUPPORTED ENGINES"})
		display(columns{"MYSQL", p.SupportedEngines.MySQL})
		display(columns{"PG", p.SupportedEngines.PG})
		display(columns{"REDIS", p.SupportedEngines.Redis})

		display(columns{" "})

		display(columns{"MAX CONNECTIONS"})
		display(columns{"MYSQL", p.MaxConnections.MySQL})
		display(columns{"PG", p.MaxConnections.PG})

		display(columns{" "})

		display(columns{"LOCATIONS", p.Locations})
		display(columns{"---------------------------"})
	}

	Meta(meta)
	flush()
}

func DatabaseList(databases []govultr.Database, meta *govultr.Meta) {
	for _, d := range databases {
		display(columns{"ID", d.ID})
		display(columns{"DATE CREATED", d.DateCreated})
		display(columns{"PLAN", d.Plan})
		display(columns{"PLAN DISK", d.PlanDisk})
		display(columns{"PLAN RAM", d.PlanRAM})
		display(columns{"PLAN VCPUS", d.PlanVCPUs})
		display(columns{"PLAN REPLICAS", d.PlanReplicas})
		display(columns{"REGION", d.Region})
		display(columns{"STATUS", d.Status})
		display(columns{"LABEL", d.Label})
		display(columns{"TAG", d.Tag})
		display(columns{"DATABASE ENGINE", d.DatabaseEngine})
		display(columns{"DATABASE ENGINE VERSION", d.DatabaseEngineVersion})
		display(columns{"DB NAME", d.DBName})
		display(columns{"HOST", d.Host})
		display(columns{"USER", d.User})
		display(columns{"PASSWORD", d.Password})
		display(columns{"PORT", d.Port})
		display(columns{"MAINTENANCE DOW", d.MaintenanceDOW})
		display(columns{"MAINTENANCE TIME", d.MaintenanceTime})
		display(columns{"LATEST BACKUP", d.LatestBackup})
		display(columns{"TRUSTED IPS", d.TrustedIPs})

		if d.DatabaseEngine == "mysql" {
			display(columns{"MYSQL SQL MODES", d.MySQLSQLModes})
			display(columns{"MYSQL REQUIRE PRIMARY KEY", d.MySQLRequirePrimaryKey})
			display(columns{"MYSQL SLOW QUERY LOG", d.MySQLSlowQueryLog})
			display(columns{"MYSQL LONG QUERY TIME", d.MySQLLongQueryTime})
		}

		if d.DatabaseEngine == "pg" {
			display(columns{" "})
			display(columns{"PG AVAILABLE EXTENSIONS"})
			for _, ext := range d.PGAvailableExtensions {
				display(columns{"NAME", "VERSIONS"})
				display(columns{ext.Name, ext.Versions})
			}
			display(columns{" "})
		}

		if d.DatabaseEngine == "redis" {
			display(columns{"REDIS EVICTION POLICY", d.RedisEvictionPolicy})
		}

		display(columns{"CLUSTER TIME ZONE", d.ClusterTimeZone})

		if len(d.ReadReplicas) > 0 {
			display(columns{" "})
			display(columns{"READ REPLICAS"})
			for _, r := range d.ReadReplicas {
				display(columns{"ID", r.ID})
				display(columns{"DATE CREATED", r.DateCreated})
				display(columns{"PLAN", r.Plan})
				display(columns{"PLAN DISK", r.PlanDisk})
				display(columns{"PLAN RAM", r.PlanRAM})
				display(columns{"PLAN VCPUS", r.PlanVCPUs})
				display(columns{"PLAN REPLICAS", r.PlanReplicas})
				display(columns{"REGION", r.Region})
				display(columns{"STATUS", r.Status})
				display(columns{"LABEL", r.Label})
				display(columns{"TAG", r.Tag})
				display(columns{"DATABASE ENGINE", r.DatabaseEngine})
				display(columns{"DATABASE ENGINE VERSION", r.DatabaseEngineVersion})
				display(columns{"DB NAME", r.DBName})
				display(columns{"HOST", r.Host})
				display(columns{"USER", r.User})
				display(columns{"PASSWORD", r.Password})
				display(columns{"PORT", r.Port})
				display(columns{"MAINTENANCE DOW", r.MaintenanceDOW})
				display(columns{"MAINTENANCE TIME", r.MaintenanceTime})
				display(columns{"LATEST BACKUP", r.LatestBackup})
				display(columns{"TRUSTED IPS", r.TrustedIPs})

				if r.DatabaseEngine == "mysql" {
					display(columns{"MYSQL SQL MODES", r.MySQLSQLModes})
					display(columns{"MYSQL REQUIRE PRIMARY KEY", r.MySQLRequirePrimaryKey})
					display(columns{"MYSQL SLOW QUERY LOG", r.MySQLSlowQueryLog})
					display(columns{"MYSQL LONG QUERY TIME", r.MySQLLongQueryTime})
				}

				if r.DatabaseEngine == "pg" {
					display(columns{"PG AVAILABLE EXTENSIONS"})
					for _, rext := range r.PGAvailableExtensions {
						display(columns{"NAME", "VERSIONS"})
						display(columns{rext.Name, rext.Versions})
					}
				}

				if r.DatabaseEngine == "redis" {
					display(columns{"REDIS EVICTION POLICY", r.RedisEvictionPolicy})
				}

				display(columns{"CLUSTER TIME ZONE", r.ClusterTimeZone})
			}
		}

		display(columns{"---------------------------"})
	}

	Meta(meta)
	flush()
}
