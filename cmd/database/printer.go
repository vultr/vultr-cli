package database

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// DBsPrinter ...
type DBsPrinter struct {
	DBs  []govultr.Database `json:"databases"`
	Meta *govultr.Meta      `json:"meta"`
}

// JSON ...
func (d *DBsPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DBsPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DBsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (d *DBsPrinter) Data() [][]string { //nolint:funlen,gocyclo
	if len(d.DBs) == 0 {
		return [][]string{0: {"No databases"}}
	}

	var data [][]string
	for i := range d.DBs {
		data = append(data,
			[]string{"ID", d.DBs[i].ID},
			[]string{"DATE CREATED", d.DBs[i].DateCreated},
			[]string{"PLAN", d.DBs[i].Plan},
			[]string{"PLAN DISK", strconv.Itoa(d.DBs[i].PlanDisk)},
			[]string{"PLAN RAM", strconv.Itoa(d.DBs[i].PlanRAM)},
			[]string{"PLAN VCPUS", strconv.Itoa(d.DBs[i].PlanVCPUs)},
		)

		if d.DBs[i].DatabaseEngine == "kafka" {
			data = append(data, []string{"PLAN BROKERS", strconv.Itoa(d.DBs[i].PlanBrokers)})
		} else {
			data = append(data, []string{"PLAN REPLICAS", strconv.Itoa(*d.DBs[i].PlanReplicas)})
		}

		data = append(data,
			[]string{"REGION", d.DBs[i].Region},
			[]string{"DATABASE ENGINE", d.DBs[i].DatabaseEngine},
			[]string{"DATABASE ENGINE VERSION", d.DBs[i].DatabaseEngineVersion},
			[]string{"VPC ID", d.DBs[i].VPCID},
			[]string{"STATUS", d.DBs[i].Status},
			[]string{"LABEL", d.DBs[i].Label},
			[]string{"TAG", d.DBs[i].Tag},
			[]string{"DB NAME", d.DBs[i].DBName},
		)

		if d.DBs[i].DatabaseEngine == "ferretpg" {
			data = append(data,
				[]string{" "},
				[]string{"FERRETDB CREDENTIALS"},
				[]string{"HOST", d.DBs[i].FerretDBCredentials.Host},
				[]string{"PORT", strconv.Itoa(d.DBs[i].FerretDBCredentials.Port)},
				[]string{"USER", d.DBs[i].FerretDBCredentials.User},
				[]string{"PASSWORD", d.DBs[i].FerretDBCredentials.Password},
				[]string{"PUBLIC IP", d.DBs[i].FerretDBCredentials.PublicIP},
			)

			if d.DBs[i].FerretDBCredentials.PrivateIP != "" {
				data = append(data,
					[]string{"PRIVATE IP", d.DBs[i].FerretDBCredentials.PrivateIP},
				)
			}

			data = append(data, []string{" "})
		}

		data = append(data, []string{"HOST", d.DBs[i].Host})

		if d.DBs[i].PublicHost != "" {
			data = append(data, []string{"PUBLIC HOST", d.DBs[i].PublicHost})
		}

		data = append(data, []string{"PORT", d.DBs[i].Port})

		if d.DBs[i].DatabaseEngine == "kafka" {
			data = append(data, []string{"SASL PORT", d.DBs[i].SASLPort})
		}

		data = append(data,
			[]string{"USER", d.DBs[i].User},
			[]string{"PASSWORD", d.DBs[i].Password},
		)

		if d.DBs[i].DatabaseEngine == "kafka" {
			data = append(data,
				[]string{"ACCESS KEY", d.DBs[i].AccessKey},
				[]string{"ACCESS CERT", d.DBs[i].AccessCert},
			)

			if d.DBs[i].EnableKafkaREST != nil {
				data = append(data, []string{"ENABLE KAFKA REST", strconv.FormatBool(*d.DBs[i].EnableKafkaREST)})

				if d.DBs[i].KafkaRESTURI != "" {
					data = append(data, []string{"KAFKA REST URI", d.DBs[i].KafkaRESTURI})
				}
			}

			if d.DBs[i].EnableSchemaRegistry != nil {
				data = append(data, []string{"ENABLE SCHEMA REGISTRY", strconv.FormatBool(*d.DBs[i].EnableSchemaRegistry)})

				if d.DBs[i].SchemaRegistryURI != "" {
					data = append(data, []string{"SCHEMA REGISTRY URI", d.DBs[i].SchemaRegistryURI})
				}
			}

			if d.DBs[i].EnableKafkaConnect != nil {
				data = append(data, []string{"ENABLE KAFKA CONNECT", strconv.FormatBool(*d.DBs[i].EnableKafkaConnect)})
			}
		}

		data = append(data,
			[]string{"MAINTENANCE DOW", d.DBs[i].MaintenanceDOW},
			[]string{"MAINTENANCE TIME", d.DBs[i].MaintenanceTime},
		)

		if d.DBs[i].BackupHour != nil {
			data = append(data, []string{"BACKUP HOUR", *d.DBs[i].BackupHour})
		}

		if d.DBs[i].BackupMinute != nil {
			data = append(data, []string{"BACKUP MINUTE", *d.DBs[i].BackupMinute})
		}

		data = append(data,
			[]string{"LATEST BACKUP", d.DBs[i].LatestBackup},
			[]string{"TRUSTED IPS", printer.ArrayOfStringsToString(d.DBs[i].TrustedIPs)},
			[]string{"CA CERTIFICATE", d.DBs[i].CACertificate},
		)

		if d.DBs[i].DatabaseEngine == "mysql" {
			data = append(data,
				[]string{"MYSQL SQL MODES", printer.ArrayOfStringsToString(d.DBs[i].MySQLSQLModes)},
				[]string{"MYSQL REQUIRE PRIMARY KEY", strconv.FormatBool(*d.DBs[i].MySQLRequirePrimaryKey)},
				[]string{"MYSQL SLOW QUERY LOG", strconv.FormatBool(*d.DBs[i].MySQLSlowQueryLog)},
			)

			if *d.DBs[i].MySQLSlowQueryLog {
				data = append(data,
					[]string{"MYSQL LONG QUERY TIME", strconv.Itoa(d.DBs[i].MySQLLongQueryTime)},
				)
			}
		}

		if d.DBs[i].DatabaseEngine == "pg" && len(d.DBs[i].PGAvailableExtensions) > 0 {
			data = append(data,
				[]string{" "},
				[]string{"PG AVAILABLE EXTENSIONS"},
				[]string{"NAME", "VERSIONS"},
			)

			for j := range d.DBs[i].PGAvailableExtensions {
				if len(d.DBs[i].PGAvailableExtensions[j].Versions) > 0 {
					data = append(data, []string{
						d.DBs[i].PGAvailableExtensions[j].Name,
						printer.ArrayOfStringsToString(d.DBs[i].PGAvailableExtensions[j].Versions)})
				} else {
					data = append(data, []string{d.DBs[i].PGAvailableExtensions[j].Name, ""})
				}
			}
			data = append(data, []string{" "})
		}

		if d.DBs[i].DatabaseEngine == "valkey" {
			data = append(data, []string{"EVICTION POLICY", d.DBs[i].EvictionPolicy})
		}

		data = append(data, []string{"CLUSTER TIME ZONE", d.DBs[i].ClusterTimeZone})

		if len(d.DBs[i].ReadReplicas) > 0 {
			data = append(data,
				[]string{" "},
				[]string{"READ REPLICAS"},
			)

			for j := range d.DBs[i].ReadReplicas {
				data = append(data,
					[]string{"ID", d.DBs[i].ReadReplicas[j].ID},
					[]string{"DATE CREATED", d.DBs[i].ReadReplicas[j].DateCreated},
					[]string{"PLAN", d.DBs[i].ReadReplicas[j].Plan},
					[]string{"PLAN DISK", strconv.Itoa(d.DBs[i].ReadReplicas[j].PlanDisk)},
					[]string{"PLAN RAM", strconv.Itoa(d.DBs[i].ReadReplicas[j].PlanRAM)},
					[]string{"PLAN VCPUS", strconv.Itoa(d.DBs[i].ReadReplicas[j].PlanVCPUs)},
					[]string{"PLAN REPLICAS", strconv.Itoa(*d.DBs[i].ReadReplicas[j].PlanReplicas)},
					[]string{"REGION", d.DBs[i].ReadReplicas[j].Region},
					[]string{"DATABASE ENGINE", d.DBs[i].ReadReplicas[j].DatabaseEngine},
					[]string{"DATABASE ENGINE VERSION", d.DBs[i].ReadReplicas[j].DatabaseEngineVersion},
					[]string{"VPC ID", d.DBs[i].ReadReplicas[j].VPCID},
					[]string{"STATUS", d.DBs[i].ReadReplicas[j].Status},
					[]string{"LABEL", d.DBs[i].ReadReplicas[j].Label},
					[]string{"TAG", d.DBs[i].ReadReplicas[j].Tag},
					[]string{"DB NAME", d.DBs[i].ReadReplicas[j].DBName},
				)

				if d.DBs[i].ReadReplicas[j].DatabaseEngine == "ferretpg" {
					data = append(data,
						[]string{" "},

						[]string{"FERRETDB CREDENTIALS"},
						[]string{"HOST", d.DBs[i].ReadReplicas[j].FerretDBCredentials.Host},
						[]string{"PORT", strconv.Itoa(d.DBs[i].ReadReplicas[j].FerretDBCredentials.Port)},
						[]string{"USER", d.DBs[i].ReadReplicas[j].FerretDBCredentials.User},
						[]string{"PASSWORD", d.DBs[i].ReadReplicas[j].FerretDBCredentials.Password},
						[]string{"PUBLIC IP", d.DBs[i].ReadReplicas[j].FerretDBCredentials.PublicIP},
					)

					if d.DBs[i].ReadReplicas[j].FerretDBCredentials.PrivateIP != "" {
						data = append(data,
							[]string{"PRIVATE IP", d.DBs[i].ReadReplicas[j].FerretDBCredentials.PrivateIP},
						)
					}

					data = append(data, []string{" "})
				}

				data = append(data, []string{"HOST", d.DBs[i].ReadReplicas[j].Host})

				if d.DBs[i].ReadReplicas[j].PublicHost != "" {
					data = append(data, []string{"PUBLIC HOST", d.DBs[i].ReadReplicas[j].PublicHost})
				}

				data = append(data,
					[]string{"PORT", d.DBs[i].ReadReplicas[j].Port},
					[]string{"USER", d.DBs[i].ReadReplicas[j].User},
					[]string{"PASSWORD", d.DBs[i].ReadReplicas[j].Password},
					[]string{"MAINTENANCE DOW", d.DBs[i].ReadReplicas[j].MaintenanceDOW},
					[]string{"MAINTENANCE TIME", d.DBs[i].ReadReplicas[j].MaintenanceTime},
				)

				if d.DBs[i].ReadReplicas[j].BackupHour != nil {
					data = append(data, []string{"BACKUP HOUR", *d.DBs[i].ReadReplicas[j].BackupHour})
				}

				if d.DBs[i].ReadReplicas[j].BackupMinute != nil {
					data = append(data, []string{"BACKUP MINUTE", *d.DBs[i].ReadReplicas[j].BackupMinute})
				}

				data = append(data,
					[]string{"LATEST BACKUP", d.DBs[i].ReadReplicas[j].LatestBackup},
					[]string{"TRUSTED IPS", printer.ArrayOfStringsToString(d.DBs[i].ReadReplicas[j].TrustedIPs)},
					[]string{"CA CERTIFICATE", d.DBs[i].ReadReplicas[j].CACertificate},
				)

				if d.DBs[i].ReadReplicas[j].DatabaseEngine == "mysql" {
					data = append(data,
						[]string{"MYSQL SQL MODES", printer.ArrayOfStringsToString(d.DBs[i].ReadReplicas[j].MySQLSQLModes)},
						[]string{"MYSQL REQUIRE PRIMARY KEY", strconv.FormatBool(*d.DBs[i].ReadReplicas[j].MySQLRequirePrimaryKey)},
						[]string{"MYSQL SLOW QUERY LOG", strconv.FormatBool(*d.DBs[i].ReadReplicas[j].MySQLSlowQueryLog)},
					)

					if *d.DBs[i].ReadReplicas[j].MySQLSlowQueryLog {
						data = append(data, []string{"MYSQL LONG QUERY TIME", strconv.Itoa(d.DBs[i].ReadReplicas[j].MySQLLongQueryTime)})
					}
				}

				if d.DBs[i].ReadReplicas[j].DatabaseEngine == "pg" && len(d.DBs[i].ReadReplicas[j].PGAvailableExtensions) > 0 {
					data = append(data,
						[]string{" "},
						[]string{"PG AVAILABLE EXTENSIONS"},
						[]string{"NAME", "VERSIONS"},
					)

					for k := range d.DBs[i].ReadReplicas[j].PGAvailableExtensions {
						if len(d.DBs[i].ReadReplicas[j].PGAvailableExtensions[k].Versions) > 0 {
							data = append(data, []string{
								d.DBs[i].ReadReplicas[j].PGAvailableExtensions[k].Name,
								printer.ArrayOfStringsToString(d.DBs[i].ReadReplicas[j].PGAvailableExtensions[k].Versions),
							})
						} else {
							data = append(data, []string{d.DBs[i].ReadReplicas[j].PGAvailableExtensions[k].Name, ""})
						}
					}

					data = append(data, []string{" "})
				}

				if d.DBs[i].ReadReplicas[j].DatabaseEngine == "valkey" {
					data = append(data, []string{"EVICTION POLICY", d.DBs[i].ReadReplicas[j].EvictionPolicy})
				}

				data = append(data, []string{"CLUSTER TIME ZONE", d.DBs[i].ReadReplicas[j].ClusterTimeZone})
			}
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (d *DBsPrinter) Paging() [][]string {
	paging := &printer.Total{Total: d.Meta.Total}
	return paging.Compose()
}

// ======================================

// DBPrinter ...
type DBPrinter struct {
	DB *govultr.Database `json:"database"`
}

// JSON ...
func (d *DBPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DBPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DBPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (d *DBPrinter) Data() [][]string { //nolint:funlen,gocyclo
	var data [][]string
	data = append(data,
		[]string{"ID", d.DB.ID},
		[]string{"DATE CREATED", d.DB.DateCreated},
		[]string{"PLAN", d.DB.Plan},
		[]string{"PLAN DISK", strconv.Itoa(d.DB.PlanDisk)},
		[]string{"PLAN RAM", strconv.Itoa(d.DB.PlanRAM)},
		[]string{"PLAN VCPUS", strconv.Itoa(d.DB.PlanVCPUs)},
	)

	if d.DB.DatabaseEngine == "kafka" {
		data = append(data, []string{"PLAN BROKERS", strconv.Itoa(d.DB.PlanBrokers)})
	} else {
		data = append(data, []string{"PLAN REPLICAS", strconv.Itoa(*d.DB.PlanReplicas)})
	}

	data = append(data,
		[]string{"REGION", d.DB.Region},
		[]string{"DATABASE ENGINE", d.DB.DatabaseEngine},
		[]string{"DATABASE ENGINE VERSION", d.DB.DatabaseEngineVersion},
		[]string{"VPC ID", d.DB.VPCID},
		[]string{"STATUS", d.DB.Status},
		[]string{"LABEL", d.DB.Label},
		[]string{"TAG", d.DB.Tag},
		[]string{"DB NAME", d.DB.DBName},
	)

	if d.DB.DatabaseEngine == "ferretpg" {
		data = append(data,
			[]string{" "},
			[]string{"FERRETDB CREDENTIALS"},
			[]string{"HOST", d.DB.FerretDBCredentials.Host},
			[]string{"PORT", strconv.Itoa(d.DB.FerretDBCredentials.Port)},
			[]string{"USER", d.DB.FerretDBCredentials.User},
			[]string{"PASSWORD", d.DB.FerretDBCredentials.Password},
			[]string{"PUBLIC IP", d.DB.FerretDBCredentials.PublicIP},
		)

		if d.DB.FerretDBCredentials.PrivateIP != "" {
			data = append(data,
				[]string{"PRIVATE IP", d.DB.FerretDBCredentials.PrivateIP},
			)
		}

		data = append(data, []string{" "})
	}

	data = append(data, []string{"HOST", d.DB.Host})

	if d.DB.PublicHost != "" {
		data = append(data, []string{"PUBLIC HOST", d.DB.PublicHost})
	}

	data = append(data, []string{"PORT", d.DB.Port})

	if d.DB.DatabaseEngine == "kafka" {
		data = append(data, []string{"SASL PORT", d.DB.SASLPort})
	}

	data = append(data,
		[]string{"USER", d.DB.User},
		[]string{"PASSWORD", d.DB.Password},
	)

	if d.DB.DatabaseEngine == "kafka" {
		data = append(data,
			[]string{"ACCESS KEY", d.DB.AccessKey},
			[]string{"ACCESS CERT", d.DB.AccessCert},
		)

		if d.DB.DatabaseEngine == "kafka" {
			data = append(data, []string{"SASL PORT", d.DB.SASLPort})

			if d.DB.EnableKafkaREST != nil {
				data = append(data, []string{"ENABLE KAFKA REST", strconv.FormatBool(*d.DB.EnableKafkaREST)})

				if d.DB.KafkaRESTURI != "" {
					data = append(data, []string{"KAFKA REST URI", d.DB.KafkaRESTURI})
				}
			}

			if d.DB.EnableSchemaRegistry != nil {
				data = append(data, []string{"ENABLE SCHEMA REGISTRY", strconv.FormatBool(*d.DB.EnableSchemaRegistry)})

				if d.DB.SchemaRegistryURI != "" {
					data = append(data, []string{"SCHEMA REGISTRY URI", d.DB.SchemaRegistryURI})
				}
			}

			if d.DB.EnableKafkaConnect != nil {
				data = append(data, []string{"ENABLE KAFKA CONNECT", strconv.FormatBool(*d.DB.EnableKafkaConnect)})
			}
		}
	}

	data = append(data,
		[]string{"MAINTENANCE DOW", d.DB.MaintenanceDOW},
		[]string{"MAINTENANCE TIME", d.DB.MaintenanceTime},
	)

	if d.DB.BackupHour != nil {
		data = append(data, []string{"BACKUP HOUR", *d.DB.BackupHour})
	}

	if d.DB.BackupMinute != nil {
		data = append(data, []string{"BACKUP MINUTE", *d.DB.BackupMinute})
	}

	data = append(data,
		[]string{"LATEST BACKUP", d.DB.LatestBackup},
		[]string{"TRUSTED IPS", printer.ArrayOfStringsToString(d.DB.TrustedIPs)},
		[]string{"CA CERTIFICATE", d.DB.CACertificate},
	)

	if d.DB.DatabaseEngine == "mysql" {
		data = append(data,
			[]string{"MYSQL SQL MODES", printer.ArrayOfStringsToString(d.DB.MySQLSQLModes)},
			[]string{"MYSQL REQUIRE PRIMARY KEY", strconv.FormatBool(*d.DB.MySQLRequirePrimaryKey)},
			[]string{"MYSQL SLOW QUERY LOG", strconv.FormatBool(*d.DB.MySQLSlowQueryLog)},
		)

		if *d.DB.MySQLSlowQueryLog {
			data = append(data,
				[]string{"MYSQL LONG QUERY TIME", strconv.Itoa(d.DB.MySQLLongQueryTime)},
			)
		}
	}

	if d.DB.DatabaseEngine == "pg" && len(d.DB.PGAvailableExtensions) > 0 {
		data = append(data,
			[]string{" "},
			[]string{"PG AVAILABLE EXTENSIONS"},
			[]string{"NAME", "VERSIONS"},
		)

		for i := range d.DB.PGAvailableExtensions {
			if len(d.DB.PGAvailableExtensions[i].Versions) > 0 {
				data = append(data, []string{
					d.DB.PGAvailableExtensions[i].Name,
					printer.ArrayOfStringsToString(d.DB.PGAvailableExtensions[i].Versions)})
			} else {
				data = append(data, []string{d.DB.PGAvailableExtensions[i].Name, ""})
			}
		}
		data = append(data, []string{" "})
	}

	if d.DB.DatabaseEngine == "valkey" {
		data = append(data, []string{"EVICTION POLICY", d.DB.EvictionPolicy})
	}

	data = append(data, []string{"CLUSTER TIME ZONE", d.DB.ClusterTimeZone})

	if len(d.DB.ReadReplicas) > 0 {
		data = append(data,
			[]string{" "},
			[]string{"READ REPLICAS"},
		)

		for i := range d.DB.ReadReplicas {
			data = append(data,
				[]string{"ID", d.DB.ReadReplicas[i].ID},
				[]string{"DATE CREATED", d.DB.ReadReplicas[i].DateCreated},
				[]string{"PLAN", d.DB.ReadReplicas[i].Plan},
				[]string{"PLAN DISK", strconv.Itoa(d.DB.ReadReplicas[i].PlanDisk)},
				[]string{"PLAN RAM", strconv.Itoa(d.DB.ReadReplicas[i].PlanRAM)},
				[]string{"PLAN VCPUS", strconv.Itoa(d.DB.ReadReplicas[i].PlanVCPUs)},
				[]string{"PLAN REPLICAS", strconv.Itoa(*d.DB.ReadReplicas[i].PlanReplicas)},
				[]string{"REGION", d.DB.ReadReplicas[i].Region},
				[]string{"DATABASE ENGINE", d.DB.ReadReplicas[i].DatabaseEngine},
				[]string{"DATABASE ENGINE VERSION", d.DB.ReadReplicas[i].DatabaseEngineVersion},
				[]string{"VPC ID", d.DB.ReadReplicas[i].VPCID},
				[]string{"STATUS", d.DB.ReadReplicas[i].Status},
				[]string{"LABEL", d.DB.ReadReplicas[i].Label},
				[]string{"TAG", d.DB.ReadReplicas[i].Tag},
				[]string{"DB NAME", d.DB.ReadReplicas[i].DBName},
			)

			if d.DB.ReadReplicas[i].DatabaseEngine == "ferretpg" {
				data = append(data,
					[]string{" "},

					[]string{"FERRETDB CREDENTIALS"},
					[]string{"HOST", d.DB.ReadReplicas[i].FerretDBCredentials.Host},
					[]string{"PORT", strconv.Itoa(d.DB.ReadReplicas[i].FerretDBCredentials.Port)},
					[]string{"USER", d.DB.ReadReplicas[i].FerretDBCredentials.User},
					[]string{"PASSWORD", d.DB.ReadReplicas[i].FerretDBCredentials.Password},
					[]string{"PUBLIC IP", d.DB.ReadReplicas[i].FerretDBCredentials.PublicIP},
				)

				if d.DB.ReadReplicas[i].FerretDBCredentials.PrivateIP != "" {
					data = append(data,
						[]string{"PRIVATE IP", d.DB.ReadReplicas[i].FerretDBCredentials.PrivateIP},
					)
				}

				data = append(data, []string{" "})
			}

			data = append(data, []string{"HOST", d.DB.ReadReplicas[i].Host})

			if d.DB.ReadReplicas[i].PublicHost != "" {
				data = append(data, []string{"PUBLIC HOST", d.DB.ReadReplicas[i].PublicHost})
			}

			data = append(data,
				[]string{"USER", d.DB.ReadReplicas[i].User},
				[]string{"PASSWORD", d.DB.ReadReplicas[i].Password},
				[]string{"PORT", d.DB.ReadReplicas[i].Port},
				[]string{"MAINTENANCE DOW", d.DB.ReadReplicas[i].MaintenanceDOW},
				[]string{"MAINTENANCE TIME", d.DB.ReadReplicas[i].MaintenanceTime},
			)

			if d.DB.ReadReplicas[i].BackupHour != nil {
				data = append(data, []string{"BACKUP HOUR", *d.DB.ReadReplicas[i].BackupHour})
			}

			if d.DB.ReadReplicas[i].BackupMinute != nil {
				data = append(data, []string{"BACKUP MINUTE", *d.DB.ReadReplicas[i].BackupMinute})
			}

			data = append(data,
				[]string{"LATEST BACKUP", d.DB.ReadReplicas[i].LatestBackup},
				[]string{"TRUSTED IPS", printer.ArrayOfStringsToString(d.DB.ReadReplicas[i].TrustedIPs)},
				[]string{"CA CERTIFICATE", d.DB.ReadReplicas[i].CACertificate},
			)

			if d.DB.ReadReplicas[i].DatabaseEngine == "mysql" {
				data = append(data,
					[]string{"MYSQL SQL MODES", printer.ArrayOfStringsToString(d.DB.ReadReplicas[i].MySQLSQLModes)},
					[]string{"MYSQL REQUIRE PRIMARY KEY", strconv.FormatBool(*d.DB.ReadReplicas[i].MySQLRequirePrimaryKey)},
					[]string{"MYSQL SLOW QUERY LOG", strconv.FormatBool(*d.DB.ReadReplicas[i].MySQLSlowQueryLog)},
				)

				if *d.DB.ReadReplicas[i].MySQLSlowQueryLog {
					data = append(data, []string{"MYSQL LONG QUERY TIME", strconv.Itoa(d.DB.ReadReplicas[i].MySQLLongQueryTime)})
				}
			}

			if d.DB.ReadReplicas[i].DatabaseEngine == "pg" && len(d.DB.ReadReplicas[i].PGAvailableExtensions) > 0 {
				data = append(data,
					[]string{" "},
					[]string{"PG AVAILABLE EXTENSIONS"},
					[]string{"NAME", "VERSIONS"},
				)

				for j := range d.DB.ReadReplicas[i].PGAvailableExtensions {
					if len(d.DB.ReadReplicas[i].PGAvailableExtensions[j].Versions) > 0 {
						data = append(data, []string{
							d.DB.ReadReplicas[i].PGAvailableExtensions[j].Name,
							printer.ArrayOfStringsToString(d.DB.ReadReplicas[i].PGAvailableExtensions[j].Versions),
						})
					} else {
						data = append(data, []string{d.DB.ReadReplicas[i].PGAvailableExtensions[j].Name, ""})
					}
				}

				data = append(data, []string{" "})
			}

			if d.DB.ReadReplicas[i].DatabaseEngine == "valkey" {
				data = append(data, []string{"EVICTION POLICY", d.DB.ReadReplicas[i].EvictionPolicy})
			}

			data = append(data, []string{"CLUSTER TIME ZONE", d.DB.ReadReplicas[i].ClusterTimeZone})
		}
	}

	return data
}

// Paging ...
func (d *DBPrinter) Paging() [][]string {
	return nil
}

// ======================================

// DBsSummaryPrinter ...
type DBsSummaryPrinter struct {
	DBs  []govultr.Database `json:"databases"`
	Meta *govultr.Meta      `json:"meta"`
}

// JSON ...
func (d *DBsSummaryPrinter) JSON() []byte {
	return printer.MarshalObject(d, "json")
}

// YAML ...
func (d *DBsSummaryPrinter) YAML() []byte {
	return printer.MarshalObject(d, "yaml")
}

// Columns ...
func (d *DBsSummaryPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"LABEL",
		"STATUS",
		"ENGINE",
		"VERSION",
	}}
}

// Data ...
func (d *DBsSummaryPrinter) Data() [][]string {
	if len(d.DBs) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range d.DBs {
		data = append(data, []string{

			d.DBs[i].ID,
			d.DBs[i].Region,
			d.DBs[i].Label,
			d.DBs[i].Status,
			d.DBs[i].DatabaseEngine,
			d.DBs[i].DatabaseEngineVersion,
		})
	}

	return data
}

// Paging ...
func (d *DBsSummaryPrinter) Paging() [][]string {
	paging := &printer.Total{Total: d.Meta.Total}
	return paging.Compose()
}

// ======================================

// PlansPrinter ...
type PlansPrinter struct {
	Plans []govultr.DatabasePlan `json:"plans"`
	Meta  *govultr.Meta          `json:"meta"`
}

// JSON ...
func (p *PlansPrinter) JSON() []byte {
	return printer.MarshalObject(p, "json")
}

// YAML ...
func (p *PlansPrinter) YAML() []byte {
	return printer.MarshalObject(p, "yaml")
}

// Columns ...
func (p *PlansPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (p *PlansPrinter) Data() [][]string {
	if len(p.Plans) == 0 {
		return [][]string{0: {"No database plans available"}}
	}

	var data [][]string
	for i := range p.Plans {
		data = append(data,
			[]string{"ID", p.Plans[i].ID},
			[]string{"NUMBER OF NODES", strconv.Itoa(p.Plans[i].NumberOfNodes)},
			[]string{"TYPE", p.Plans[i].Type},
		)

		if !*p.Plans[i].SupportedEngines.Kafka {
			data = append(data,
				[]string{"VCPU COUNT", strconv.Itoa(p.Plans[i].VCPUCount)},
				[]string{"RAM", strconv.Itoa(p.Plans[i].RAM)},
			)
		}

		if !*p.Plans[i].SupportedEngines.Valkey {
			data = append(data,
				[]string{"DISK", strconv.Itoa(p.Plans[i].Disk)},
			)
		}

		data = append(data,
			[]string{"MONTHLY COST", strconv.Itoa(p.Plans[i].MonthlyCost)},
			[]string{" "},
			[]string{"SUPPORTED ENGINES"},
			[]string{"MYSQL", strconv.FormatBool(*p.Plans[i].SupportedEngines.MySQL)},
			[]string{"PG", strconv.FormatBool(*p.Plans[i].SupportedEngines.PG)},
			[]string{"VALKEY", strconv.FormatBool(*p.Plans[i].SupportedEngines.Valkey)},
			[]string{"KAFKA", strconv.FormatBool(*p.Plans[i].SupportedEngines.Kafka)},
			[]string{" "},
		)

		if *p.Plans[i].SupportedEngines.MySQL || *p.Plans[i].SupportedEngines.PG {
			data = append(data,
				[]string{"MAX CONNECTIONS"},
				[]string{"MYSQL", strconv.Itoa(p.Plans[i].MaxConnections.MySQL)},
				[]string{"PG", strconv.Itoa(p.Plans[i].MaxConnections.PG)},
				[]string{" "},
			)
		}

		data = append(data,
			[]string{"LOCATIONS", printer.ArrayOfStringsToString(p.Plans[i].Locations)},
			[]string{"---------------------------"},
		)
	}

	return data
}

// Paging ...
func (p *PlansPrinter) Paging() [][]string {
	paging := &printer.Total{Total: p.Meta.Total}
	return paging.Compose()
}

// ======================================

// UsagePrinter ...
type UsagePrinter struct {
	Usage *govultr.DatabaseUsage `json:"usage"`
}

// JSON ...
func (u *UsagePrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UsagePrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UsagePrinter) Columns() [][]string {
	return nil
}

// Data ...
func (u *UsagePrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"DISK USAGE"},
		[]string{"CURRENT (GB)", strconv.FormatFloat(
			float64(u.Usage.Disk.CurrentGB),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
		[]string{"MAXIMUM (GB)", strconv.FormatFloat(
			float64(u.Usage.Disk.MaxGB),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
		[]string{"PERCENTAGE", strconv.FormatFloat(
			float64(u.Usage.Disk.Percentage),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
		[]string{" "},
		[]string{"MEMORY USAGE"},
		[]string{"CURRENT (MB)", strconv.FormatFloat(
			float64(u.Usage.Memory.CurrentMB),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
		[]string{"MAXIMUM (MB)", strconv.FormatFloat(
			float64(u.Usage.Memory.MaxMB),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
		[]string{"PERCENTAGE", strconv.FormatFloat(
			float64(u.Usage.Memory.Percentage),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
		[]string{" "},
		[]string{"CPU USAGE"},
		[]string{"PERCENTAGE", strconv.FormatFloat(
			float64(u.Usage.CPU.Percentage),
			'f',
			utils.FloatPrecision,
			utils.FloatBitDepth,
		)},
	)

	return data
}

// Paging ...
func (u *UsagePrinter) Paging() [][]string {
	return nil
}

// ======================================

// UsersPrinter ...
type UsersPrinter struct {
	Users []govultr.DatabaseUser `json:"users"`
	Meta  *govultr.Meta          `json:"meta"`
}

// JSON ...
func (u *UsersPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UsersPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UsersPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (u *UsersPrinter) Data() [][]string {
	if len(u.Users) == 0 {
		return [][]string{0: {"No database users"}}
	}

	var data [][]string
	for i := range u.Users {
		data = append(data,
			[]string{"USERNAME", u.Users[i].Username},
			[]string{"PASSWORD", u.Users[i].Password},
		)

		if u.Users[i].Encryption != "" {
			data = append(data, []string{"ENCRYPTION", u.Users[i].Encryption})
		}

		if u.Users[i].AccessControl != nil {
			data = append(data,
				[]string{"ACCESS CONTROL"},
				[]string{"ACL CATEGORIES", printer.ArrayOfStringsToString(u.Users[i].AccessControl.ACLCategories)},
				[]string{"ACL CHANNELS", printer.ArrayOfStringsToString(u.Users[i].AccessControl.ACLChannels)},
				[]string{"ACL COMMANDS", printer.ArrayOfStringsToString(u.Users[i].AccessControl.ACLCommands)},
				[]string{"ACL KEYS", printer.ArrayOfStringsToString(u.Users[i].AccessControl.ACLKeys)},
			)
		}

		if u.Users[i].Permission != "" {
			data = append(data, []string{"PERMISSION", u.Users[i].Permission})
		}

		if u.Users[i].AccessKey != "" {
			data = append(data, []string{"ACCESS KEY", u.Users[i].AccessKey})
		}

		if u.Users[i].AccessCert != "" {
			data = append(data, []string{"ACCESS CERT", u.Users[i].AccessCert})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (u *UsersPrinter) Paging() [][]string {
	paging := &printer.Total{Total: u.Meta.Total}
	return paging.Compose()
}

// ======================================

// UserPrinter ...
type UserPrinter struct {
	User *govultr.DatabaseUser `json:"user"`
}

// JSON ...
func (u *UserPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UserPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UserPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (u *UserPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"USERNAME", u.User.Username},
		[]string{"PASSWORD", u.User.Password},
	)

	if u.User.Encryption != "" {
		data = append(data, []string{"ENCRYPTION", u.User.Encryption})
	}

	if u.User.AccessControl != nil {
		data = append(data,
			[]string{"ACCESS CONTROL"},
			[]string{"ACL CATEGORIES", printer.ArrayOfStringsToString(u.User.AccessControl.ACLCategories)},
			[]string{"ACL CHANNELS", printer.ArrayOfStringsToString(u.User.AccessControl.ACLChannels)},
			[]string{"ACL COMMANDS", printer.ArrayOfStringsToString(u.User.AccessControl.ACLCommands)},
			[]string{"ACL KEYS", printer.ArrayOfStringsToString(u.User.AccessControl.ACLKeys)},
		)
	}

	if u.User.Permission != "" {
		data = append(data, []string{"PERMISSION", u.User.Permission})
	}

	if u.User.AccessKey != "" {
		data = append(data, []string{"ACCESS KEY", u.User.AccessKey})
	}

	if u.User.AccessCert != "" {
		data = append(data, []string{"ACCESS CERT", u.User.AccessCert})
	}

	return data
}

// Paging ...
func (u *UserPrinter) Paging() [][]string {
	return nil
}

// ======================================

// LogicalDBsPrinter ...
type LogicalDBsPrinter struct {
	DBs  []govultr.DatabaseDB `json:"dbs"`
	Meta *govultr.Meta        `json:"meta"`
}

// JSON ...
func (l *LogicalDBsPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LogicalDBsPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LogicalDBsPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
	}}
}

// Data ...
func (l *LogicalDBsPrinter) Data() [][]string {
	if len(l.DBs) == 0 {
		return [][]string{0: {"---"}}
	}

	var data [][]string
	for i := range l.DBs {
		data = append(data, []string{l.DBs[i].Name})
	}

	return data
}

// Paging ...
func (l *LogicalDBsPrinter) Paging() [][]string {
	paging := &printer.Total{Total: l.Meta.Total}
	return paging.Compose()
}

// ======================================

// LogicalDBPrinter ...
type LogicalDBPrinter struct {
	DB *govultr.DatabaseDB `json:"db"`
}

// JSON ...
func (l *LogicalDBPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LogicalDBPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LogicalDBPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
	}}
}

// Data ...
func (l *LogicalDBPrinter) Data() [][]string {
	return [][]string{0: {l.DB.Name}}
}

// Paging ...
func (l *LogicalDBPrinter) Paging() [][]string {
	return nil
}

// ======================================

// TopicsPrinter ...
type TopicsPrinter struct {
	Topics []govultr.DatabaseTopic `json:"topics"`
	Meta   *govultr.Meta           `json:"meta"`
}

// JSON ...
func (t *TopicsPrinter) JSON() []byte {
	return printer.MarshalObject(t, "json")
}

// YAML ...
func (t *TopicsPrinter) YAML() []byte {
	return printer.MarshalObject(t, "yaml")
}

// Columns ...
func (t *TopicsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (t *TopicsPrinter) Data() [][]string {
	if len(t.Topics) == 0 {
		return [][]string{0: {"No database topics"}}
	}

	var data [][]string
	for i := range t.Topics {
		data = append(data,
			[]string{"NAME", t.Topics[i].Name},
			[]string{"PARTITIONS", strconv.Itoa(t.Topics[i].Partitions)},
			[]string{"REPLICATION", strconv.Itoa(t.Topics[i].Replication)},
			[]string{"RETENTION HOURS", strconv.Itoa(t.Topics[i].RetentionHours)},
			[]string{"RETENTION BYTES", strconv.Itoa(t.Topics[i].RetentionBytes)},
		)

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (t *TopicsPrinter) Paging() [][]string {
	paging := &printer.Total{Total: t.Meta.Total}
	return paging.Compose()
}

// ======================================

// TopicPrinter ...
type TopicPrinter struct {
	Topic *govultr.DatabaseTopic `json:"topic"`
}

// JSON ...
func (t *TopicPrinter) JSON() []byte {
	return printer.MarshalObject(t, "json")
}

// YAML ...
func (t *TopicPrinter) YAML() []byte {
	return printer.MarshalObject(t, "yaml")
}

// Columns ...
func (t *TopicPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (t *TopicPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"NAME", t.Topic.Name},
		[]string{"PARTITIONS", strconv.Itoa(t.Topic.Partitions)},
		[]string{"REPLICATION", strconv.Itoa(t.Topic.Replication)},
		[]string{"RETENTION HOURS", strconv.Itoa(t.Topic.RetentionHours)},
		[]string{"RETENTION BYTES", strconv.Itoa(t.Topic.RetentionBytes)},
	)

	return data
}

// Paging ...
func (t *TopicPrinter) Paging() [][]string {
	return nil
}

// ======================================

// QuotasPrinter ...
type QuotasPrinter struct {
	Quotas []govultr.DatabaseQuota `json:"quotas"`
	Meta   *govultr.Meta           `json:"meta"`
}

// JSON ...
func (q *QuotasPrinter) JSON() []byte {
	return printer.MarshalObject(q, "json")
}

// YAML ...
func (q *QuotasPrinter) YAML() []byte {
	return printer.MarshalObject(q, "yaml")
}

// Columns ...
func (q *QuotasPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (q *QuotasPrinter) Data() [][]string {
	if len(q.Quotas) == 0 {
		return [][]string{0: {"No database quotas"}}
	}

	var data [][]string
	for i := range q.Quotas {
		data = append(data,
			[]string{"CLIENT ID", q.Quotas[i].ClientID},
			[]string{"CONSUMER BYTE RATE", strconv.Itoa(q.Quotas[i].ConsumerByteRate)},
			[]string{"PRODUCER BYTE RATE", strconv.Itoa(q.Quotas[i].ProducerByteRate)},
			[]string{"REQUEST PERCENTAGE", strconv.Itoa(q.Quotas[i].RequestPercentage)},
			[]string{"USER", q.Quotas[i].User},
		)

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (q *QuotasPrinter) Paging() [][]string {
	paging := &printer.Total{Total: q.Meta.Total}
	return paging.Compose()
}

// ======================================

// QuotaPrinter ...
type QuotaPrinter struct {
	Quota *govultr.DatabaseQuota `json:"quota"`
}

// JSON ...
func (q *QuotaPrinter) JSON() []byte {
	return printer.MarshalObject(q, "json")
}

// YAML ...
func (q *QuotaPrinter) YAML() []byte {
	return printer.MarshalObject(q, "yaml")
}

// Columns ...
func (q *QuotaPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (q *QuotaPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"CLIENT ID", q.Quota.ClientID},
		[]string{"CONSUMER BYTE RATE", strconv.Itoa(q.Quota.ConsumerByteRate)},
		[]string{"PRODUCER BYTE RATE", strconv.Itoa(q.Quota.ProducerByteRate)},
		[]string{"REQUEST PERCENTAGE", strconv.Itoa(q.Quota.RequestPercentage)},
		[]string{"USER", q.Quota.User},
	)

	return data
}

// Paging ...
func (q *QuotaPrinter) Paging() [][]string {
	return nil
}

// ======================================

// UpdatesPrinter ...
type UpdatesPrinter struct {
	Updates []string `json:"available_updates"`
}

// JSON ...
func (u *UpdatesPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UpdatesPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UpdatesPrinter) Columns() [][]string {
	return [][]string{0: {"AVAILABLE UPDATES"}}
}

// Data ...
func (u *UpdatesPrinter) Data() [][]string {
	var data [][]string

	for i := range u.Updates {
		data = append(data, []string{u.Updates[i]})
	}

	return data
}

// Paging ...
func (u *UpdatesPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AlertsPrinter ...
type AlertsPrinter struct {
	Alerts []govultr.DatabaseAlert `json:"alerts"`
}

// JSON ...
func (a *AlertsPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AlertsPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AlertsPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
	}}
}

// Data ...
func (a *AlertsPrinter) Data() [][]string {
	if len(a.Alerts) == 0 {
		return [][]string{0: {"No active database alerts"}}
	}

	var data [][]string
	for i := range a.Alerts {
		data = append(data,
			[]string{"TIMESTAMP", a.Alerts[i].Timestamp},
			[]string{"MESSAGE TYPE", a.Alerts[i].MessageType},
			[]string{"DESCRIPTION", a.Alerts[i].Description},
		)

		if a.Alerts[i].Recommendation != "" {
			data = append(data, []string{"RECOMMENDATION", a.Alerts[i].Recommendation})
		}

		if a.Alerts[i].MaintenanceScheduled != "" {
			data = append(data, []string{"MAINTENANCE SCHEDULED", a.Alerts[i].MaintenanceScheduled})
		}

		if a.Alerts[i].ResourceType != "" {
			data = append(data, []string{"RESOURCE TYPE", a.Alerts[i].ResourceType})
		}

		if a.Alerts[i].TableCount != 0 {
			data = append(data, []string{"TABLE COUNT", strconv.Itoa(a.Alerts[i].TableCount)})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (a *AlertsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// MigrationPrinter ...
type MigrationPrinter struct {
	Migration *govultr.DatabaseMigration `json:"migration"`
}

// JSON ...
func (m *MigrationPrinter) JSON() []byte {
	return printer.MarshalObject(m, "json")
}

// YAML ...
func (m *MigrationPrinter) YAML() []byte {
	return printer.MarshalObject(m, "yaml")
}

// Columns ...
func (m *MigrationPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (m *MigrationPrinter) Data() [][]string {
	var data [][]string
	data = append(data, []string{"STATUS", m.Migration.Status})

	if m.Migration.Method != "" {
		data = append(data, []string{"METHOD", m.Migration.Method})
	}

	if m.Migration.Error != "" {
		data = append(data, []string{"ERROR", m.Migration.Error})
	}

	data = append(data,
		[]string{" "},
		[]string{"CREDENTIALS"},
		[]string{"HOST", m.Migration.Credentials.Host},
		[]string{"PORT", strconv.Itoa(m.Migration.Credentials.Port)},
		[]string{"USERNAME", m.Migration.Credentials.Username},
		[]string{"PASSWORD", m.Migration.Credentials.Password},
	)

	if m.Migration.Credentials.Database != "" {
		data = append(data, []string{"DATABASE", m.Migration.Credentials.Database})
	}

	if m.Migration.Credentials.IgnoredDatabases != "" {
		data = append(data, []string{"IGNORED DATABASES", m.Migration.Credentials.IgnoredDatabases})
	}

	data = append(data, []string{"SSL", strconv.FormatBool(*m.Migration.Credentials.SSL)})

	return data
}

// Paging ...
func (m *MigrationPrinter) Paging() [][]string {
	return nil
}

// ======================================

// BackupPrinter ...
type BackupPrinter struct {
	Backup *govultr.DatabaseBackups `json:"backups"`
}

// JSON ...
func (b *BackupPrinter) JSON() []byte {
	return printer.MarshalObject(b, "json")
}

// YAML ...
func (b *BackupPrinter) YAML() []byte {
	return printer.MarshalObject(b, "yaml")
}

// Columns ...
func (b *BackupPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (b *BackupPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"LATEST BACKUP"},
		[]string{"DATE", b.Backup.LatestBackup.Date},
		[]string{"TIME", b.Backup.LatestBackup.Time},
		[]string{" "},
		[]string{"OLDEST BACKUP"},
		[]string{"DATE", b.Backup.OldestBackup.Date},
		[]string{"TIME", b.Backup.OldestBackup.Time},
	)

	return data
}

// Paging ...
func (b *BackupPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ConnectionsPrinter ...
type ConnectionsPrinter struct {
	Connections     *govultr.DatabaseConnections     `json:"connections"`
	ConnectionPools []govultr.DatabaseConnectionPool `json:"connection_pools"`
	Meta            *govultr.Meta                    `json:"meta"`
}

// JSON ...
func (c *ConnectionsPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConnectionsPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConnectionsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ConnectionsPrinter) Data() [][]string {
	var data [][]string

	data = append(data,
		[]string{"CONNECTIONS"},
		[]string{"USED", strconv.Itoa(c.Connections.Used)},
		[]string{"AVAILABLE", strconv.Itoa(c.Connections.Available)},
		[]string{"MAX", strconv.Itoa(c.Connections.Max)},

		[]string{" "},
		[]string{"CONNECTION POOLS"},
	)

	for i := range c.ConnectionPools {
		data = append(data,
			[]string{"NAME", c.ConnectionPools[i].Name},
			[]string{"DATABASE", c.ConnectionPools[i].Database},
			[]string{"USERNAME", c.ConnectionPools[i].Username},
			[]string{"MODE", c.ConnectionPools[i].Mode},
			[]string{"SIZE", strconv.Itoa(c.ConnectionPools[i].Size)},
			[]string{"---------------------------"},
		)
	}

	return data
}

// Paging ...
func (c *ConnectionsPrinter) Paging() [][]string {
	paging := &printer.Total{Total: c.Meta.Total}
	return paging.Compose()
}

// ======================================

// ConnectionPoolPrinter ...
type ConnectionPoolPrinter struct {
	ConnectionPool *govultr.DatabaseConnectionPool `json:"connection_pool"`
}

// JSON ...
func (c *ConnectionPoolPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConnectionPoolPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConnectionPoolPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
		"DATABASE",
		"USERNAME",
		"MODE",
		"SIZE",
	}}
}

// Data ...
func (c *ConnectionPoolPrinter) Data() [][]string {
	return [][]string{0: {
		c.ConnectionPool.Name,
		c.ConnectionPool.Database,
		c.ConnectionPool.Username,
		c.ConnectionPool.Mode,
		strconv.Itoa(c.ConnectionPool.Size),
	}}
}

// Paging ...
func (c *ConnectionPoolPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AdvancedOptionsPrinter ...
type AdvancedOptionsPrinter struct {
	Configured *govultr.DatabaseAdvancedOptions `json:"configured_options"`
	Available  []govultr.AvailableOption        `json:"available_options"`
}

// JSON ...
func (a *AdvancedOptionsPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AdvancedOptionsPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AdvancedOptionsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (a *AdvancedOptionsPrinter) Data() [][]string { //nolint:dupl
	var data [][]string

	if a.Configured == (&govultr.DatabaseAdvancedOptions{}) {
		data = append(data, []string{"CONFIGURED OPTIONS", "None"})
	} else {
		data = append(data, []string{"CONFIGURED OPTIONS"})
		v := reflect.ValueOf(*a.Configured)
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				switch v.Field(i).Kind() {
				case reflect.Pointer:
					data = append(data, []string{v.Type().Field(i).Name, strconv.FormatBool(v.Field(i).Elem().Interface().(bool))})
				case reflect.Int:
					data = append(data, []string{v.Type().Field(i).Name, strconv.Itoa(v.Field(i).Interface().(int))})
				case reflect.Float64:
					data = append(data,
						[]string{
							v.Type().Field(i).Name,
							strconv.FormatFloat(float64(v.Field(i).Interface().(float64)), 'f', utils.FloatPrecision, 32),
						},
					)
				default:
					data = append(data, []string{v.Type().Field(i).Name, v.Field(i).Interface().(string)})
				}
			}
		}
	}

	data = append(data,
		[]string{" "},
		[]string{"AVAILABLE OPTIONS"},
	)

	for i := range a.Available {
		data = append(data,
			[]string{"NAME", a.Available[i].Name},
			[]string{"TYPE", a.Available[i].Type},
		)

		if a.Available[i].Type == "enum" {
			data = append(data,
				[]string{"ENUMERALS", printer.ArrayOfStringsToString(a.Available[i].Enumerals)},
			)
		}

		if a.Available[i].Type == "int" || a.Available[i].Type == "float" {
			data = append(data,
				[]string{"MIN VALUE", strconv.FormatFloat(float64(*a.Available[i].MinValue), 'f', utils.FloatPrecision, 32)},
				[]string{"MAX VALUE", strconv.FormatFloat(float64(*a.Available[i].MaxValue), 'f', utils.FloatPrecision, 32)},
			)
		}

		if len(a.Available[i].AltValues) > 0 {
			data = append(data, []string{"ALT VALUES", printer.ArrayOfIntsToString(a.Available[i].AltValues)})
		}

		if a.Available[i].Units != "" {
			data = append(data, []string{"UNITS", a.Available[i].Units})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (a *AdvancedOptionsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AdvancedOptionsKafkaRESTPrinter ...
type AdvancedOptionsKafkaRESTPrinter struct {
	Configured *govultr.DatabaseKafkaRESTAdvancedOptions `json:"configured_options"`
	Available  []govultr.AvailableOption                 `json:"available_options"`
}

// JSON ...
func (a *AdvancedOptionsKafkaRESTPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AdvancedOptionsKafkaRESTPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AdvancedOptionsKafkaRESTPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (a *AdvancedOptionsKafkaRESTPrinter) Data() [][]string { //nolint:dupl
	var data [][]string

	if a.Configured == (&govultr.DatabaseKafkaRESTAdvancedOptions{}) {
		data = append(data, []string{"CONFIGURED OPTIONS", "None"})
	} else {
		data = append(data, []string{"CONFIGURED OPTIONS"})
		v := reflect.ValueOf(*a.Configured)
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				switch v.Field(i).Kind() {
				case reflect.Pointer:
					data = append(data, []string{v.Type().Field(i).Name, strconv.FormatBool(v.Field(i).Elem().Interface().(bool))})
				case reflect.Int:
					data = append(data, []string{v.Type().Field(i).Name, strconv.Itoa(v.Field(i).Interface().(int))})
				case reflect.Float64:
					data = append(data,
						[]string{
							v.Type().Field(i).Name,
							strconv.FormatFloat(float64(v.Field(i).Interface().(float64)), 'f', utils.FloatPrecision, 32),
						},
					)
				default:
					data = append(data, []string{v.Type().Field(i).Name, v.Field(i).Interface().(string)})
				}
			}
		}
	}

	data = append(data,
		[]string{" "},
		[]string{"AVAILABLE OPTIONS"},
	)

	for i := range a.Available {
		data = append(data,
			[]string{"NAME", a.Available[i].Name},
			[]string{"TYPE", a.Available[i].Type},
		)

		if a.Available[i].Type == "enum" {
			data = append(data,
				[]string{"ENUMERALS", printer.ArrayOfStringsToString(a.Available[i].Enumerals)},
			)
		}

		if a.Available[i].Type == "int" || a.Available[i].Type == "float" {
			data = append(data,
				[]string{"MIN VALUE", strconv.FormatFloat(float64(*a.Available[i].MinValue), 'f', utils.FloatPrecision, 32)},
				[]string{"MAX VALUE", strconv.FormatFloat(float64(*a.Available[i].MaxValue), 'f', utils.FloatPrecision, 32)},
			)
		}

		if len(a.Available[i].AltValues) > 0 {
			data = append(data, []string{"ALT VALUES", printer.ArrayOfIntsToString(a.Available[i].AltValues)})
		}

		if a.Available[i].Units != "" {
			data = append(data, []string{"UNITS", a.Available[i].Units})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (a *AdvancedOptionsKafkaRESTPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AdvancedOptionsSchemaRegistryPrinter ...
type AdvancedOptionsSchemaRegistryPrinter struct {
	Configured *govultr.DatabaseSchemaRegistryAdvancedOptions `json:"configured_options"`
	Available  []govultr.AvailableOption                      `json:"available_options"`
}

// JSON ...
func (a *AdvancedOptionsSchemaRegistryPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AdvancedOptionsSchemaRegistryPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AdvancedOptionsSchemaRegistryPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (a *AdvancedOptionsSchemaRegistryPrinter) Data() [][]string { //nolint:dupl
	var data [][]string

	if a.Configured == (&govultr.DatabaseSchemaRegistryAdvancedOptions{}) {
		data = append(data, []string{"CONFIGURED OPTIONS", "None"})
	} else {
		data = append(data, []string{"CONFIGURED OPTIONS"})
		v := reflect.ValueOf(*a.Configured)
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				switch v.Field(i).Kind() {
				case reflect.Pointer:
					data = append(data, []string{v.Type().Field(i).Name, strconv.FormatBool(v.Field(i).Elem().Interface().(bool))})
				case reflect.Int:
					data = append(data, []string{v.Type().Field(i).Name, strconv.Itoa(v.Field(i).Interface().(int))})
				case reflect.Float64:
					data = append(data,
						[]string{
							v.Type().Field(i).Name,
							strconv.FormatFloat(float64(v.Field(i).Interface().(float64)), 'f', utils.FloatPrecision, 32),
						},
					)
				default:
					data = append(data, []string{v.Type().Field(i).Name, v.Field(i).Interface().(string)})
				}
			}
		}
	}

	data = append(data,
		[]string{" "},
		[]string{"AVAILABLE OPTIONS"},
	)

	for i := range a.Available {
		data = append(data,
			[]string{"NAME", a.Available[i].Name},
			[]string{"TYPE", a.Available[i].Type},
		)

		if a.Available[i].Type == "enum" {
			data = append(data,
				[]string{"ENUMERALS", printer.ArrayOfStringsToString(a.Available[i].Enumerals)},
			)
		}

		if a.Available[i].Type == "int" || a.Available[i].Type == "float" {
			data = append(data,
				[]string{"MIN VALUE", strconv.FormatFloat(float64(*a.Available[i].MinValue), 'f', utils.FloatPrecision, 32)},
				[]string{"MAX VALUE", strconv.FormatFloat(float64(*a.Available[i].MaxValue), 'f', utils.FloatPrecision, 32)},
			)
		}

		if len(a.Available[i].AltValues) > 0 {
			data = append(data, []string{"ALT VALUES", printer.ArrayOfIntsToString(a.Available[i].AltValues)})
		}

		if a.Available[i].Units != "" {
			data = append(data, []string{"UNITS", a.Available[i].Units})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (a *AdvancedOptionsSchemaRegistryPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AdvancedOptionsKafkaConnectPrinter ...
type AdvancedOptionsKafkaConnectPrinter struct {
	Configured *govultr.DatabaseKafkaConnectAdvancedOptions `json:"configured_options"`
	Available  []govultr.AvailableOption                    `json:"available_options"`
}

// JSON ...
func (a *AdvancedOptionsKafkaConnectPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AdvancedOptionsKafkaConnectPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AdvancedOptionsKafkaConnectPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (a *AdvancedOptionsKafkaConnectPrinter) Data() [][]string { //nolint:dupl
	var data [][]string

	if a.Configured == (&govultr.DatabaseKafkaConnectAdvancedOptions{}) {
		data = append(data, []string{"CONFIGURED OPTIONS", "None"})
	} else {
		data = append(data, []string{"CONFIGURED OPTIONS"})
		v := reflect.ValueOf(*a.Configured)
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				switch v.Field(i).Kind() {
				case reflect.Pointer:
					data = append(data, []string{v.Type().Field(i).Name, strconv.FormatBool(v.Field(i).Elem().Interface().(bool))})
				case reflect.Int:
					data = append(data, []string{v.Type().Field(i).Name, strconv.Itoa(v.Field(i).Interface().(int))})
				case reflect.Float64:
					data = append(data,
						[]string{
							v.Type().Field(i).Name,
							strconv.FormatFloat(float64(v.Field(i).Interface().(float64)), 'f', utils.FloatPrecision, 32),
						},
					)
				default:
					data = append(data, []string{v.Type().Field(i).Name, v.Field(i).Interface().(string)})
				}
			}
		}
	}

	data = append(data,
		[]string{" "},
		[]string{"AVAILABLE OPTIONS"},
	)

	for i := range a.Available {
		data = append(data,
			[]string{"NAME", a.Available[i].Name},
			[]string{"TYPE", a.Available[i].Type},
		)

		if a.Available[i].Type == "enum" {
			data = append(data,
				[]string{"ENUMERALS", printer.ArrayOfStringsToString(a.Available[i].Enumerals)},
			)
		}

		if a.Available[i].Type == "int" || a.Available[i].Type == "float" {
			data = append(data,
				[]string{"MIN VALUE", strconv.FormatFloat(float64(*a.Available[i].MinValue), 'f', utils.FloatPrecision, 32)},
				[]string{"MAX VALUE", strconv.FormatFloat(float64(*a.Available[i].MaxValue), 'f', utils.FloatPrecision, 32)},
			)
		}

		if len(a.Available[i].AltValues) > 0 {
			data = append(data, []string{"ALT VALUES", printer.ArrayOfIntsToString(a.Available[i].AltValues)})
		}

		if a.Available[i].Units != "" {
			data = append(data, []string{"UNITS", a.Available[i].Units})
		}

		data = append(data, []string{"---------------------------"})
	}

	return data
}

// Paging ...
func (a *AdvancedOptionsKafkaConnectPrinter) Paging() [][]string {
	return nil
}

// ======================================

// VersionsPrinter ...
type VersionsPrinter struct {
	Versions []string `json:"available_versions"`
}

// JSON ...
func (v *VersionsPrinter) JSON() []byte {
	return printer.MarshalObject(v, "json")
}

// YAML ...
func (v *VersionsPrinter) YAML() []byte {
	return printer.MarshalObject(v, "yaml")
}

// Columns ...
func (v *VersionsPrinter) Columns() [][]string {
	return [][]string{0: {
		"AVAILABLE VERSIONS",
	}}
}

// Data ...
func (v *VersionsPrinter) Data() [][]string {
	var data [][]string
	for i := range v.Versions {
		data = append(data, []string{v.Versions[i]})
	}

	return data
}

// Paging ...
func (v *VersionsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AvailableConnectorsPrinter ...
type AvailableConnectorsPrinter struct {
	AvailableConnectors []govultr.DatabaseAvailableConnector `json:"available_connectors"`
}

// JSON ...
func (c *AvailableConnectorsPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *AvailableConnectorsPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *AvailableConnectorsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *AvailableConnectorsPrinter) Data() [][]string {
	if len(c.AvailableConnectors) == 0 {
		return [][]string{0: {"No database connectors available"}}
	}

	var data [][]string
	for i := range c.AvailableConnectors {
		data = append(data,
			[]string{"CLASS", c.AvailableConnectors[i].Class},
			[]string{"TITLE", c.AvailableConnectors[i].Title},
			[]string{"VERSION", c.AvailableConnectors[i].Version},
			[]string{"TYPE", c.AvailableConnectors[i].Type},
			[]string{"DOC URL", c.AvailableConnectors[i].DocURL},
			[]string{"---------------------------"},
		)
	}

	return data
}

// Paging ...
func (c *AvailableConnectorsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ConnectorConfigSchemaPrinter ...
type ConnectorConfigSchemaPrinter struct {
	ConfigurationSchema []govultr.DatabaseConnectorConfigurationOption `json:"configuration_schema"`
}

// JSON ...
func (c *ConnectorConfigSchemaPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConnectorConfigSchemaPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConnectorConfigSchemaPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ConnectorConfigSchemaPrinter) Data() [][]string {
	if len(c.ConfigurationSchema) == 0 {
		return [][]string{0: {"No database connector configuration schema available"}}
	}

	var data [][]string
	for i := range c.ConfigurationSchema {
		data = append(data,
			[]string{"NAME", c.ConfigurationSchema[i].Name},
			[]string{"TYPE", c.ConfigurationSchema[i].Type},
			[]string{"REQUIRED", strconv.FormatBool(c.ConfigurationSchema[i].Required)},
			[]string{"DEFAULT VALUE", c.ConfigurationSchema[i].DefaultValue},
			[]string{"DESCRIPTION", c.ConfigurationSchema[i].Description},
			[]string{"---------------------------"},
		)
	}

	return data
}

// Paging ...
func (c *ConnectorConfigSchemaPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ConnectorsPrinter ...
type ConnectorsPrinter struct {
	Connectors []govultr.DatabaseConnector `json:"connectors"`
	Meta       *govultr.Meta               `json:"meta"`
}

// JSON ...
func (c *ConnectorsPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConnectorsPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConnectorsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ConnectorsPrinter) Data() [][]string {
	if len(c.Connectors) == 0 {
		return [][]string{0: {"No database connectors available"}}
	}

	var data [][]string
	for i := range c.Connectors {
		data = append(data,
			[]string{"NAME", c.Connectors[i].Name},
			[]string{"CLASS", c.Connectors[i].Class},
			[]string{"TOPICS", c.Connectors[i].Topics},
		)

		if c.Connectors[i].Config != nil {
			data = append(data,
				[]string{" "},
				[]string{"CONFIG"},
			)

			for key, value := range c.Connectors[i].Config {
				data = append(data,
					[]string{strings.ToUpper(key), value.(string)},
				)
			}
		} else {
			data = append(data,
				[]string{"CONFIG", "NONE"},
			)
		}

		data = append(data,
			[]string{"---------------------------"},
		)
	}

	return data
}

// Paging ...
func (d *ConnectorsPrinter) Paging() [][]string {
	paging := &printer.Total{Total: d.Meta.Total}
	return paging.Compose()
}

// ======================================

// ConnectorPrinter ...
type ConnectorPrinter struct {
	Connector *govultr.DatabaseConnector `json:"connector"`
}

// JSON ...
func (c *ConnectorPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConnectorPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConnectorPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ConnectorPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"NAME", c.Connector.Name},
		[]string{"CLASS", c.Connector.Class},
		[]string{"TOPICS", c.Connector.Topics},
	)

	if c.Connector.Config != nil {
		data = append(data,
			[]string{" "},
			[]string{"CONFIG"},
		)

		for key, value := range c.Connector.Config {
			data = append(data,
				[]string{strings.ToUpper(key), value.(string)},
			)
		}
	} else {
		data = append(data,
			[]string{"CONFIG", "NONE"},
		)
	}

	return data
}

// Paging ...
func (c *ConnectorPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ConnectorStatusPrinter ...
type ConnectorStatusPrinter struct {
	ConnectorStatus *govultr.DatabaseConnectorStatus `json:"connector_status"`
}

// JSON ...
func (c *ConnectorStatusPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConnectorStatusPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConnectorStatusPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ConnectorStatusPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"STATE", c.ConnectorStatus.State},
	)

	if c.ConnectorStatus.Tasks != nil {
		data = append(data,
			[]string{"TASKS"},
		)

		for i := range c.ConnectorStatus.Tasks {
			data = append(data,
				[]string{"ID", strconv.Itoa(c.ConnectorStatus.Tasks[i].ID)},
				[]string{"STATE", c.ConnectorStatus.Tasks[i].State},
				[]string{"TRACE", c.ConnectorStatus.Tasks[i].Trace},
			)
		}
	}

	return data
}

// Paging ...
func (c *ConnectorStatusPrinter) Paging() [][]string {
	return nil
}
