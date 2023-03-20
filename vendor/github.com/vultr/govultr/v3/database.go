package govultr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const databasePath = "/v2/databases"

// DatabaseService is the interface to interact with the Database endpoints on the Vultr API
// Link: https://www.vultr.com/api/#tag/managed-databases
type DatabaseService interface {
	ListPlans(ctx context.Context, options *DBPlanListOptions) ([]DatabasePlan, *Meta, *http.Response, error)

	List(ctx context.Context, options *DBListOptions) ([]Database, *Meta, *http.Response, error)
	Create(ctx context.Context, databaseReq *DatabaseCreateReq) (*Database, *http.Response, error)
	Get(ctx context.Context, databaseID string) (*Database, *http.Response, error)
	Update(ctx context.Context, databaseID string, databaseReq *DatabaseUpdateReq) (*Database, *http.Response, error)
	Delete(ctx context.Context, databaseID string) error

	ListUsers(ctx context.Context, databaseID string) ([]DatabaseUser, *Meta, *http.Response, error)
	CreateUser(ctx context.Context, databaseID string, databaseUserReq *DatabaseUserCreateReq) (*DatabaseUser, *http.Response, error)
	GetUser(ctx context.Context, databaseID string, username string) (*DatabaseUser, *http.Response, error)
	UpdateUser(ctx context.Context, databaseID string, username string, databaseUserReq *DatabaseUserUpdateReq) (*DatabaseUser, *http.Response, error)
	DeleteUser(ctx context.Context, databaseID string, username string) error

	ListDBs(ctx context.Context, databaseID string) ([]DatabaseDB, *Meta, *http.Response, error)
	CreateDB(ctx context.Context, databaseID string, databaseDBReq *DatabaseDBCreateReq) (*DatabaseDB, *http.Response, error)
	GetDB(ctx context.Context, databaseID string, dbname string) (*DatabaseDB, *http.Response, error)
	DeleteDB(ctx context.Context, databaseID string, dbname string) error

	ListMaintenanceUpdates(ctx context.Context, databaseID string) ([]string, *http.Response, error)
	StartMaintenance(ctx context.Context, databaseID string) (string, *http.Response, error)

	ListServiceAlerts(ctx context.Context, databaseID string, databaseAlertsReq *DatabaseListAlertsReq) ([]DatabaseAlert, *http.Response, error)

	GetMigrationStatus(ctx context.Context, databaseID string) (*DatabaseMigration, *http.Response, error)
	StartMigration(ctx context.Context, databaseID string, databaseMigrationReq *DatabaseMigrationStartReq) (*DatabaseMigration, *http.Response, error)
	DetachMigration(ctx context.Context, databaseID string) error

	AddReadOnlyReplica(ctx context.Context, databaseID string, databaseReplicaReq *DatabaseAddReplicaReq) (*Database, *http.Response, error)

	GetBackupInformation(ctx context.Context, databaseID string) (*DatabaseBackups, *http.Response, error)
	RestoreFromBackup(ctx context.Context, databaseID string, databaseRestoreReq *DatabaseBackupRestoreReq) (*Database, *http.Response, error)
	Fork(ctx context.Context, databaseID string, databaseForkReq *DatabaseForkReq) (*Database, *http.Response, error)

	ListConnectionPools(ctx context.Context, databaseID string) (*DatabaseConnections, []DatabaseConnectionPool, *Meta, *http.Response, error)
	CreateConnectionPool(ctx context.Context, databaseID string, databaseConnectionPoolReq *DatabaseConnectionPoolCreateReq) (*DatabaseConnectionPool, *http.Response, error)
	GetConnectionPool(ctx context.Context, databaseID string, poolName string) (*DatabaseConnectionPool, *http.Response, error)
	UpdateConnectionPool(ctx context.Context, databaseID string, poolName string, databaseConnectionPoolReq *DatabaseConnectionPoolUpdateReq) (*DatabaseConnectionPool, *http.Response, error)
	DeleteConnectionPool(ctx context.Context, databaseID string, poolName string) error

	ListAdvancedOptions(ctx context.Context, databaseID string) (*DatabaseAdvancedOptions, []AvailableOption, *http.Response, error)
	UpdateAdvancedOptions(ctx context.Context, databaseID string, databaseAdvancedOptionsReq *DatabaseAdvancedOptions) (*DatabaseAdvancedOptions, []AvailableOption, *http.Response, error)

	ListAvailableVersions(ctx context.Context, databaseID string) ([]string, *http.Response, error)
	StartVersionUpgrade(ctx context.Context, databaseID string, databaseVersionUpgradeReq *DatabaseVersionUpgradeReq) (string, *http.Response, error)
}

// DatabaseServiceHandler handles interaction with the server methods for the Vultr API
type DatabaseServiceHandler struct {
	client *Client
}

// DBPlanListOptions handles GET request parameters for the ListPlans endpoint
type DBPlanListOptions struct {
	Engine string `url:"engine,omitempty"`
	Nodes  int    `url:"nodes,omitempty"`
	Region string `url:"region,omitempty"`
}

// DatabasePlan represents a Managed Database plan
type DatabasePlan struct {
	ID               string           `json:"id"`
	NumberOfNodes    int              `json:"number_of_nodes"`
	Type             string           `json:"type"`
	VCPUCount        int              `json:"vcpu_count"`
	RAM              int              `json:"ram"`
	Disk             int              `json:"disk"`
	MonthlyCost      int              `json:"monthly_cost"`
	SupportedEngines SupportedEngines `json:"supported_engines"`
	MaxConnections   *MaxConnections  `json:"max_connections,omitempty"`
	Locations        []string         `json:"locations"`
}

// SupportedEngines represents an object containing supported database engine types for Managed Database plans
type SupportedEngines struct {
	MySQL *bool `json:"mysql"`
	PG    *bool `json:"pg"`
	Redis *bool `json:"redis"`
}

// MaxConnections represents an object containing the maximum number of connections by engine type for Managed Database plans
type MaxConnections struct {
	MySQL int `json:"mysql,omitempty"`
	PG    int `json:"pg,omitempty"`
}

// databasePlansBase holds the entire ListPlans API response
type databasePlansBase struct {
	DatabasePlans []DatabasePlan `json:"plans"`
	Meta          *Meta          `json:"meta"`
}

// DBListOptions handles GET request parameters for the List endpoint
type DBListOptions struct {
	Label  string `url:"label,omitempty"`
	Tag    string `url:"tag,omitempty"`
	Region string `url:"region,omitempty"`
}

// Database represents a Managed Database subscription
type Database struct {
	ID                     string        `json:"id"`
	DateCreated            string        `json:"date_created"`
	Plan                   string        `json:"plan"`
	PlanDisk               int           `json:"plan_disk"`
	PlanRAM                int           `json:"plan_ram"`
	PlanVCPUs              int           `json:"plan_vcpus"`
	PlanReplicas           int           `json:"plan_replicas"`
	Region                 string        `json:"region"`
	Status                 string        `json:"status"`
	Label                  string        `json:"label"`
	Tag                    string        `json:"tag"`
	DatabaseEngine         string        `json:"database_engine"`
	DatabaseEngineVersion  string        `json:"database_engine_version"`
	DBName                 string        `json:"dbname,omitempty"`
	Host                   string        `json:"host"`
	User                   string        `json:"user"`
	Password               string        `json:"password"`
	Port                   string        `json:"port"`
	MaintenanceDOW         string        `json:"maintenance_dow"`
	MaintenanceTime        string        `json:"maintenance_time"`
	LatestBackup           string        `json:"latest_backup"`
	TrustedIPs             []string      `json:"trusted_ips"`
	MySQLSQLModes          []string      `json:"mysql_sql_modes,omitempty"`
	MySQLRequirePrimaryKey *bool         `json:"mysql_require_primary_key,omitempty"`
	MySQLSlowQueryLog      *bool         `json:"mysql_slow_query_log,omitempty"`
	MySQLLongQueryTime     int           `json:"mysql_long_query_time,omitempty"`
	PGAvailableExtensions  []PGExtension `json:"pg_available_extensions,omitempty"`
	RedisEvictionPolicy    string        `json:"redis_eviction_policy,omitempty"`
	ClusterTimeZone        string        `json:"cluster_time_zone,omitempty"`
	ReadReplicas           []Database    `json:"read_replicas,omitempty"`
}

// PGExtension represents an object containing extension name and version information
type PGExtension struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions,omitempty"`
}

// databasesBase holds the entire List API response
type databasesBase struct {
	Databases []Database `json:"databases"`
	Meta      *Meta      `json:"meta"`
}

// databaseBase holds the entire Get API response
type databaseBase struct {
	Database *Database `json:"database"`
}

// DatabaseCreateReq struct used to create a database.
type DatabaseCreateReq struct {
	DatabaseEngine         string   `json:"database_engine,omitempty"`
	DatabaseEngineVersion  string   `json:"database_engine_version,omitempty"`
	Region                 string   `json:"region,omitempty"`
	Plan                   string   `json:"plan,omitempty"`
	Label                  string   `json:"label,omitempty"`
	Tag                    string   `json:"tag,omitempty"`
	MaintenanceDOW         string   `json:"maintenance_dow,omitempty"`
	MaintenanceTime        string   `json:"maintenance_time,omitempty"`
	TrustedIPs             []string `json:"trusted_ips,omitempty"`
	MySQLSQLModes          []string `json:"mysql_sql_modes,omitempty"`
	MySQLRequirePrimaryKey *bool    `json:"mysql_require_primary_key,omitempty"`
	MySQLSlowQueryLog      *bool    `json:"mysql_slow_query_log,omitempty"`
	MySQLLongQueryTime     int      `json:"mysql_long_query_time,omitempty"`
	RedisEvictionPolicy    string   `json:"redis_eviction_policy,omitempty"`
}

// DatabaseUpdateReq struct used to update a dataase.
type DatabaseUpdateReq struct {
	DatabaseEngine         string   `json:"database_engine,omitempty"`
	DatabaseEngineVersion  string   `json:"database_engine_version,omitempty"`
	Region                 string   `json:"region,omitempty"`
	Plan                   string   `json:"plan,omitempty"`
	Label                  string   `json:"label,omitempty"`
	Tag                    string   `json:"tag,omitempty"`
	MaintenanceDOW         string   `json:"maintenance_dow,omitempty"`
	MaintenanceTime        string   `json:"maintenance_time,omitempty"`
	ClusterTimeZone        string   `json:"cluster_time_zone,omitempty"`
	TrustedIPs             []string `json:"trusted_ips,omitempty"`
	MySQLSQLModes          []string `json:"mysql_sql_modes,omitempty"`
	MySQLRequirePrimaryKey *bool    `json:"mysql_require_primary_key,omitempty"`
	MySQLSlowQueryLog      *bool    `json:"mysql_slow_query_log,omitempty"`
	MySQLLongQueryTime     int      `json:"mysql_long_query_time,omitempty"`
	RedisEvictionPolicy    string   `json:"redis_eviction_policy,omitempty"`
}

// DatabaseUser represents a user within a Managed Database cluster
type DatabaseUser struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Encryption string `json:"encryption,omitempty"`
}

// databaseUserBase holds the API response for retrieving a single database user within a Managed Database
type databaseUserBase struct {
	DatabaseUser *DatabaseUser `json:"user"`
}

// databaseUsersBase holds the API response for retrieving a list of database users within a Managed Database
type databaseUsersBase struct {
	DatabaseUsers []DatabaseUser `json:"users"`
	Meta          *Meta          `json:"meta"`
}

// DatabaseUserCreateReq struct used to create a user within a Managed Database.
type DatabaseUserCreateReq struct {
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	Encryption string `json:"encryption,omitempty"`
}

// DatabaseUserUpdateReq struct used to update a user within a Managed Database.
type DatabaseUserUpdateReq struct {
	Password string `json:"password,omitempty"`
}

// DatabaseDB represents a logical database within a Managed Database cluster
type DatabaseDB struct {
	Name string `json:"name"`
}

// databaseDBBase holds the API response for retrieving a single logical database within a Managed Database
type databaseDBBase struct {
	DatabaseDB *DatabaseDB `json:"db"`
}

// databaseDBsBase holds the API response for retrieving a list of logical databases within a Managed Database
type databaseDBsBase struct {
	DatabaseDBs []DatabaseDB `json:"dbs"`
	Meta        *Meta        `json:"meta"`
}

// DatabaseDBCreateReq struct used to create a logical database within a Managed Database.
type DatabaseDBCreateReq struct {
	Name string `json:"name"`
}

// databaseDBsBase holds the API response for retrieving a list of available maintenance updates within a Managed Database
type databaseUpdatesBase struct {
	AvailableUpdates []string `json:"available_updates"`
}

// databaseMessage is a bsic object holding a return message for certain API endpoints
type databaseMessage struct {
	Message string `json:"message"`
}

// DatabaseAlert represents a service alert for a Managed Database cluster
type DatabaseAlert struct {
	Timestamp            string `json:"timestamp"`
	MessageType          string `json:"message_type"`
	Description          string `json:"description"`
	Recommendation       string `json:"recommendation,omitempty"`
	MaintenanceScheduled string `json:"maintenance_scheduled,omitempty"`
	ResourceType         string `json:"resource_type,omitempty"`
	TableCount           int    `json:"table_count,omitempty"`
}

// databaseDBsBase holds the API response for querying service alerts within a Managed Database
type databaseAlertsBase struct {
	DatabaseAlerts []DatabaseAlert `json:"alerts"`
}

// DatabaseListAlertsReq struct used to query service alerts for a Managed Database.
type DatabaseListAlertsReq struct {
	Period string `json:"period"`
}

// DatabaseMigration represents migration details for a Managed Database cluster
type DatabaseMigration struct {
	Status      string              `json:"status"`
	Method      string              `json:"method,omitempty"`
	Error       string              `json:"error,omitempty"`
	Credentials DatabaseCredentials `json:"credentials"`
}

// DatabaseCredentials represents migration credentials for migration within a Managed Database cluster
type DatabaseCredentials struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Database         string `json:"database,omitempty"`
	IgnoredDatabases string `json:"ignored_databases,omitempty"`
	SSL              *bool  `json:"ssl"`
}

// databaseMigrationBase represents a migration status object API response for a Managed Database
type databaseMigrationBase struct {
	Migration *DatabaseMigration `json:"migration"`
}

// DatabaseMigrationStartReq struct used to start a migration for a Managed Database.
type DatabaseMigrationStartReq struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Database         string `json:"database,omitempty"`
	IgnoredDatabases string `json:"ignored_databases,omitempty"`
	SSL              *bool  `json:"ssl"`
}

// DatabaseAddReplicaReq struct used to add a read-only replica to a Managed Database.
type DatabaseAddReplicaReq struct {
	Region string `json:"region,omitempty"`
	Label  string `json:"label,omitempty"`
}

// DatabaseBackups represents backup information for a Managed Database cluster
type DatabaseBackups struct {
	LatestBackup DatabaseBackup `json:"latest_backup,omitempty"`
	OldestBackup DatabaseBackup `json:"oldest_backup,omitempty"`
}

// DatabaseBackup represents individual backup details for a Managed Database cluster
type DatabaseBackup struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

// DatabaseBackupRestoreReq struct used to restore the backup of a Managed Database to a new subscription.
type DatabaseBackupRestoreReq struct {
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"`
	Date  string `json:"date,omitempty"`
	Time  string `json:"time,omitempty"`
}

// DatabaseForkReq struct used to fork a Managed Database to a new subscription from a backup.
type DatabaseForkReq struct {
	Label  string `json:"label,omitempty"`
	Region string `json:"region,omitempty"`
	Plan   string `json:"plan,omitempty"`
	Type   string `json:"type,omitempty"`
	Date   string `json:"date,omitempty"`
	Time   string `json:"time,omitempty"`
}

// DatabaseConnectionPool represents a PostgreSQL connection pool within a Managed Database cluster
type DatabaseConnectionPool struct {
	Name     string `json:"name"`
	Database string `json:"database"`
	Username string `json:"username"`
	Mode     string `json:"mode"`
	Size     int    `json:"size"`
}

// DatabaseConnections represents a an object containing used and available connections for a PostgreSQL Managed Database cluster
type DatabaseConnections struct {
	Used      int `json:"used"`
	Available int `json:"available"`
	Max       int `json:"max"`
}

// databaseConnectionPoolBase represents the API response for retrieving a single connection pool for a PostgreSQL Managed Database
type databaseConnectionPoolBase struct {
	ConnectionPool *DatabaseConnectionPool `json:"connection_pool"`
}

// databaseConnectionPoolBase represents the API response for retrieving all connection pool information for a PostgreSQL Managed Database
type databaseConnectionPoolsBase struct {
	Connections     *DatabaseConnections     `json:"connections"`
	ConnectionPools []DatabaseConnectionPool `json:"connection_pools"`
	Meta            *Meta                    `json:"meta"`
}

// DatabaseConnectionPoolCreateReq struct used to create a connection pool within a PostgreSQL Managed Database.
type DatabaseConnectionPoolCreateReq struct {
	Name     string `json:"name,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Size     int    `json:"size,omitempty"`
}

// DatabaseConnectionPoolUpdateReq struct used to update a connection pool within a PostgreSQL Managed Database.
type DatabaseConnectionPoolUpdateReq struct {
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Size     int    `json:"size,omitempty"`
}

// DatabaseAdvancedOptions represents user configurable advanced options within a PostgreSQL Managed Database cluster
type DatabaseAdvancedOptions struct {
	AutovacuumAnalyzeScaleFactor    float32 `json:"autovacuum_analyze_scale_factor,omitempty"`
	AutovacuumAnalyzeThreshold      int     `json:"autovacuum_analyze_threshold,omitempty"`
	AutovacuumFreezeMaxAge          int     `json:"autovacuum_freeze_max_age,omitempty"`
	AutovacuumMaxWorkers            int     `json:"autovacuum_max_workers,omitempty"`
	AutovacuumNaptime               int     `json:"autovacuum_naptime,omitempty"`
	AutovacuumVacuumCostDelay       int     `json:"autovacuum_vacuum_cost_delay,omitempty"`
	AutovacuumVacuumCostLimit       int     `json:"autovacuum_vacuum_cost_limit,omitempty"`
	AutovacuumVacuumScaleFactor     float32 `json:"autovacuum_vacuum_scale_factor,omitempty"`
	AutovacuumVacuumThreshold       int     `json:"autovacuum_vacuum_threshold,omitempty"`
	BGWRITERDelay                   int     `json:"bgwriter_delay,omitempty"`
	BGWRITERFlushAFter              int     `json:"bgwriter_flush_after,omitempty"`
	BGWRITERLRUMaxPages             int     `json:"bgwriter_lru_maxpages,omitempty"`
	BGWRITERLRUMultiplier           float32 `json:"bgwriter_lru_multiplier,omitempty"`
	DeadlockTimeout                 int     `json:"deadlock_timeout,omitempty"`
	DefaultToastCompression         string  `json:"default_toast_compression,omitempty"`
	IdleInTransactionSessionTimeout int     `json:"idle_in_transaction_session_timeout,omitempty"`
	Jit                             *bool   `json:"jit,omitempty"`
	LogAutovacuumMinDuration        int     `json:"log_autovacuum_min_duration,omitempty"`
	LogErrorVerbosity               string  `json:"log_error_verbosity,omitempty"`
	LogLinePrefix                   string  `json:"log_line_prefix,omitempty"`
	LogMinDurationStatement         int     `json:"log_min_duration_statement,omitempty"`
	MaxFilesPerProcess              int     `json:"max_files_per_process,omitempty"`
	MaxLocksPerTransaction          int     `json:"max_locks_per_transaction,omitempty"`
	MaxLogicalReplicationWorkers    int     `json:"max_logical_replication_workers,omitempty"`
	MaxParallelWorkers              int     `json:"max_parallel_workers,omitempty"`
	MaxParallelWorkersPerGather     int     `json:"max_parallel_workers_per_gather,omitempty"`
	MaxPredLocksPerTransaction      int     `json:"max_pred_locks_per_transaction,omitempty"`
	MaxPreparedTransactions         int     `json:"max_prepared_transactions,omitempty"`
	MaxReplicationSlots             int     `json:"max_replication_slots,omitempty"`
	MaxStackDepth                   int     `json:"max_stack_depth,omitempty"`
	MaxStandbyArchiveDelay          int     `json:"max_standby_archive_delay,omitempty"`
	MaxStandbyStreamingDelay        int     `json:"max_standby_streaming_delay,omitempty"`
	MaxWalSenders                   int     `json:"max_wal_senders,omitempty"`
	MaxWorkerProcesses              int     `json:"max_worker_processes,omitempty"`
	PGPartmanBGWInterval            int     `json:"pg_partman_bgw.interval,omitempty"`
	PGPartmanBGWRole                string  `json:"pg_partman_bgw.role,omitempty"`
	PGStateStatementsTrack          string  `json:"pg_stat_statements.track,omitempty"`
	TempFileLimit                   int     `json:"temp_file_limit,omitempty"`
	TrackActivityQuerySize          int     `json:"track_activity_query_size,omitempty"`
	TrackCommitTimestamp            string  `json:"track_commit_timestamp,omitempty"`
	TrackFunctions                  string  `json:"track_functions,omitempty"`
	TrackIOTiming                   string  `json:"track_io_timing,omitempty"`
	WALSenderTImeout                int     `json:"wal_sender_timeout,omitempty"`
	WALWriterDelay                  int     `json:"wal_writer_delay,omitempty"`
}

// AvailableOption represents an available advanced configuration option for a PostgreSQL Managed Database cluster
type AvailableOption struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Enumerals []string `json:"enumerals,omitempty"`
	MinValue  *int     `json:"min_value,omitempty"`
	MaxValue  *int     `json:"max_value,omitempty"`
	AltValues []int    `json:"alt_values,omitempty"`
	Units     string   `json:"units,omitempty"`
}

// databaseAdvancedOptionsBase represents the API response for PostgreSQL advanced configuration options for a Managed Database
type databaseAdvancedOptionsBase struct {
	ConfiguredOptions *DatabaseAdvancedOptions `json:"configured_options"`
	AvailableOptions  []AvailableOption        `json:"available_options"`
}

// DatabaseAvailableVersions represents available versions upgrades for a Managed Database cluster
type DatabaseAvailableVersions struct {
	AvailableVersions []string `json:"available_versions"`
}

// DatabaseVersionUpgradeReq struct used to initiate a version upgrade for a PostgreSQL Managed Database.
type DatabaseVersionUpgradeReq struct {
	Version string `json:"version,omitempty"`
}

// ListPlans retrieves all database plans
func (d *DatabaseServiceHandler) ListPlans(ctx context.Context, options *DBPlanListOptions) ([]DatabasePlan, *Meta, *http.Response, error) {
	uri := fmt.Sprintf("%s/plans", databasePath)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	databasePlans := new(databasePlansBase)
	resp, err := d.client.DoWithContext(ctx, req, databasePlans)
	if err != nil {
		return nil, nil, nil, err
	}

	return databasePlans.DatabasePlans, databasePlans.Meta, resp, nil
}

// List retrieves all databases on your account
func (d *DatabaseServiceHandler) List(ctx context.Context, options *DBListOptions) ([]Database, *Meta, *http.Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodGet, databasePath, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	databases := new(databasesBase)
	resp, err := d.client.DoWithContext(ctx, req, databases)
	if err != nil {
		return nil, nil, nil, err
	}

	return databases.Databases, databases.Meta, resp, nil
}

// Create will create the Managed Database with the given parameters
func (d *DatabaseServiceHandler) Create(ctx context.Context, databaseReq *DatabaseCreateReq) (*Database, *http.Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodPost, databasePath, databaseReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(databaseBase)
	resp, err := d.client.DoWithContext(ctx, req, database)
	if err != nil {
		return nil, nil, err
	}

	return database.Database, resp, nil
}

// Get will get the Managed Database with the given databaseID
func (d *DatabaseServiceHandler) Get(ctx context.Context, databaseID string) (*Database, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	database := new(databaseBase)
	resp, err := d.client.DoWithContext(ctx, req, database)
	if err != nil {
		return nil, nil, err
	}

	return database.Database, resp, nil
}

// Update will update the Managed Database with the given parameters
func (d *DatabaseServiceHandler) Update(ctx context.Context, databaseID string, databaseReq *DatabaseUpdateReq) (*Database, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPut, uri, databaseReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(databaseBase)
	resp, err := d.client.DoWithContext(ctx, req, database)
	if err != nil {
		return nil, nil, err
	}

	return database.Database, resp, nil
}

// Delete a Managed database. All data will be permanently lost.
func (d *DatabaseServiceHandler) Delete(ctx context.Context, databaseID string) error {
	uri := fmt.Sprintf("%s/%s", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	_, err = d.client.DoWithContext(ctx, req, nil)
	return err
}

// ListUsers retrieves all database users on your Managed Database.
func (d *DatabaseServiceHandler) ListUsers(ctx context.Context, databaseID string) ([]DatabaseUser, *Meta, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/users", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	databaseUsers := new(databaseUsersBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseUsers)
	if err != nil {
		return nil, nil, nil, err
	}

	return databaseUsers.DatabaseUsers, databaseUsers.Meta, resp, nil
}

// CreateUser will create a user within the Managed Database with the given parameters
func (d *DatabaseServiceHandler) CreateUser(ctx context.Context, databaseID string, databaseUserReq *DatabaseUserCreateReq) (*DatabaseUser, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/users", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseUserReq)
	if err != nil {
		return nil, nil, err
	}

	databaseUser := new(databaseUserBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseUser)
	if err != nil {
		return nil, nil, err
	}

	return databaseUser.DatabaseUser, resp, nil
}

// GetUser retrieves information on an individual user within a Managed Database based on a username and databaseID
func (d *DatabaseServiceHandler) GetUser(ctx context.Context, databaseID, username string) (*DatabaseUser, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/users/%s", databasePath, databaseID, username)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseUser := new(databaseUserBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseUser)
	if err != nil {
		return nil, nil, err
	}

	return databaseUser.DatabaseUser, resp, nil
}

// UpdateUser will update a user within the Managed Database with the given parameters
func (d *DatabaseServiceHandler) UpdateUser(ctx context.Context, databaseID, username string, databaseUserReq *DatabaseUserUpdateReq) (*DatabaseUser, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/users/%s", databasePath, databaseID, username)

	req, err := d.client.NewRequest(ctx, http.MethodPut, uri, databaseUserReq)
	if err != nil {
		return nil, nil, err
	}

	databaseUser := new(databaseUserBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseUser)
	if err != nil {
		return nil, nil, err
	}

	return databaseUser.DatabaseUser, resp, nil
}

// DeleteUser will delete a user within the Managed database. All data will be permanently lost.
func (d *DatabaseServiceHandler) DeleteUser(ctx context.Context, databaseID, username string) error {
	uri := fmt.Sprintf("%s/%s/users/%s", databasePath, databaseID, username)

	req, err := d.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	_, err = d.client.DoWithContext(ctx, req, nil)
	return err
}

// ListDBs retrieves all logical databases on your Managed Database.
func (d *DatabaseServiceHandler) ListDBs(ctx context.Context, databaseID string) ([]DatabaseDB, *Meta, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/dbs", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	databaseDBs := new(databaseDBsBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseDBs)
	if err != nil {
		return nil, nil, nil, err
	}

	return databaseDBs.DatabaseDBs, databaseDBs.Meta, resp, nil
}

// CreateDB will create a logical database within the Managed Database with the given parameters
func (d *DatabaseServiceHandler) CreateDB(ctx context.Context, databaseID string, databaseDBReq *DatabaseDBCreateReq) (*DatabaseDB, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/dbs", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseDBReq)
	if err != nil {
		return nil, nil, err
	}

	databaseDB := new(databaseDBBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseDB)
	if err != nil {
		return nil, nil, err
	}

	return databaseDB.DatabaseDB, resp, nil
}

// GetDB retrieves information on an individual logical database within a Managed Database based on a dbname and databaseID
func (d *DatabaseServiceHandler) GetDB(ctx context.Context, databaseID, dbname string) (*DatabaseDB, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/dbs/%s", databasePath, databaseID, dbname)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseDB := new(databaseDBBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseDB)
	if err != nil {
		return nil, nil, err
	}

	return databaseDB.DatabaseDB, resp, nil
}

// DeleteDB will delete a user within the Managed database
func (d *DatabaseServiceHandler) DeleteDB(ctx context.Context, databaseID, dbname string) error {
	uri := fmt.Sprintf("%s/%s/dbs/%s", databasePath, databaseID, dbname)

	req, err := d.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	_, err = d.client.DoWithContext(ctx, req, nil)
	return err
}

// ListMaintenanceUpdates retrieves all available maintenance updates for your Managed Database.
func (d *DatabaseServiceHandler) ListMaintenanceUpdates(ctx context.Context, databaseID string) ([]string, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/maintenance", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseUpdates := new(databaseUpdatesBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseUpdates)
	if err != nil {
		return nil, nil, err
	}

	return databaseUpdates.AvailableUpdates, resp, nil
}

// StartMaintenance will start the maintenance update process for your Managed Database
func (d *DatabaseServiceHandler) StartMaintenance(ctx context.Context, databaseID string) (string, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/maintenance", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return "", nil, err
	}

	databaseUpdates := new(databaseMessage)
	resp, err := d.client.DoWithContext(ctx, req, databaseUpdates)
	if err != nil {
		return "", nil, err
	}

	return databaseUpdates.Message, resp, nil
}

// ListServiceAlerts queries for service alerts for the Managed Database using the given parameters
func (d *DatabaseServiceHandler) ListServiceAlerts(ctx context.Context, databaseID string, databaseAlertsReq *DatabaseListAlertsReq) ([]DatabaseAlert, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/alerts", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseAlertsReq)
	if err != nil {
		return nil, nil, err
	}

	databaseAlerts := new(databaseAlertsBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseAlerts)
	if err != nil {
		return nil, nil, err
	}

	return databaseAlerts.DatabaseAlerts, resp, nil
}

// GetMigrationStatus retrieves the migration status for your Managed Database.
func (d *DatabaseServiceHandler) GetMigrationStatus(ctx context.Context, databaseID string) (*DatabaseMigration, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/migration", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseMigration := new(databaseMigrationBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseMigration)
	if err != nil {
		return nil, nil, err
	}

	return databaseMigration.Migration, resp, nil
}

// StartMigration will start a migration for the Managed Database using the given credentials.
func (d *DatabaseServiceHandler) StartMigration(ctx context.Context, databaseID string, databaseMigrationReq *DatabaseMigrationStartReq) (*DatabaseMigration, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/migration", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseMigrationReq)
	if err != nil {
		return nil, nil, err
	}

	databaseMigration := new(databaseMigrationBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseMigration)
	if err != nil {
		return nil, nil, err
	}

	return databaseMigration.Migration, resp, nil
}

// DetachMigration will detach a migration from the Managed database.
func (d *DatabaseServiceHandler) DetachMigration(ctx context.Context, databaseID string) error {
	uri := fmt.Sprintf("%s/%s/migration", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	_, err = d.client.DoWithContext(ctx, req, nil)
	return err
}

// AddReadOnlyReplica will add a read-only replica node to the Managed Database with the given parameters
func (d *DatabaseServiceHandler) AddReadOnlyReplica(ctx context.Context, databaseID string, databaseReplicaReq *DatabaseAddReplicaReq) (*Database, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/read-replica", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseReplicaReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(databaseBase)
	resp, err := d.client.DoWithContext(ctx, req, database)
	if err != nil {
		return nil, nil, err
	}

	return database.Database, resp, nil
}

// GetBackupInformation retrieves backup information for your Managed Database.
func (d *DatabaseServiceHandler) GetBackupInformation(ctx context.Context, databaseID string) (*DatabaseBackups, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/backups", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseBackups := new(DatabaseBackups)
	resp, err := d.client.DoWithContext(ctx, req, databaseBackups)
	if err != nil {
		return nil, nil, err
	}

	return databaseBackups, resp, nil
}

// RestoreFromBackup will create a new subscription of the same plan from a backup of the Managed Database using the given parameters
func (d *DatabaseServiceHandler) RestoreFromBackup(ctx context.Context, databaseID string, databaseRestoreReq *DatabaseBackupRestoreReq) (*Database, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/restore", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseRestoreReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(databaseBase)
	resp, err := d.client.DoWithContext(ctx, req, database)
	if err != nil {
		return nil, nil, err
	}

	return database.Database, resp, nil
}

// Fork will create a new subscription of any plan from a backup of the Managed Database using the given parameters
func (d *DatabaseServiceHandler) Fork(ctx context.Context, databaseID string, databaseForkReq *DatabaseForkReq) (*Database, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/fork", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseForkReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(databaseBase)
	resp, err := d.client.DoWithContext(ctx, req, database)
	if err != nil {
		return nil, nil, err
	}

	return database.Database, resp, nil
}

// ListConnectionPools retrieves all connection pools within your PostgreSQL Managed Database.
func (d *DatabaseServiceHandler) ListConnectionPools(ctx context.Context, databaseID string) (*DatabaseConnections, []DatabaseConnectionPool, *Meta, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/connection-pools", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	databaseConnectionPools := new(databaseConnectionPoolsBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseConnectionPools)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return databaseConnectionPools.Connections, databaseConnectionPools.ConnectionPools, databaseConnectionPools.Meta, resp, nil
}

// CreateConnectionPool will create a connection pool within the PostgreSQL Managed Database with the given parameters
func (d *DatabaseServiceHandler) CreateConnectionPool(ctx context.Context, databaseID string, databaseConnectionPoolReq *DatabaseConnectionPoolCreateReq) (*DatabaseConnectionPool, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/connection-pools", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseConnectionPoolReq)
	if err != nil {
		return nil, nil, err
	}

	databaseConnectionPool := new(databaseConnectionPoolBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseConnectionPool)
	if err != nil {
		return nil, nil, err
	}

	return databaseConnectionPool.ConnectionPool, resp, nil
}

// GetConnectionPool retrieves information on an individual connection pool within a PostgreSQL Managed Database based on a poolName and databaseID
func (d *DatabaseServiceHandler) GetConnectionPool(ctx context.Context, databaseID, poolName string) (*DatabaseConnectionPool, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/connection-pools/%s", databasePath, databaseID, poolName)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseConnectionPool := new(databaseConnectionPoolBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseConnectionPool)
	if err != nil {
		return nil, nil, err
	}

	return databaseConnectionPool.ConnectionPool, resp, nil
}

// UpdateConnectionPool will update a connection pool within the PostgreSQL Managed Database with the given parameters
func (d *DatabaseServiceHandler) UpdateConnectionPool(ctx context.Context, databaseID, poolName string, databaseConnectionPoolReq *DatabaseConnectionPoolUpdateReq) (*DatabaseConnectionPool, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/connection-pools/%s", databasePath, databaseID, poolName)

	req, err := d.client.NewRequest(ctx, http.MethodPut, uri, databaseConnectionPoolReq)
	if err != nil {
		return nil, nil, err
	}

	databaseConnectionPool := new(databaseConnectionPoolBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseConnectionPool)
	if err != nil {
		return nil, nil, err
	}

	return databaseConnectionPool.ConnectionPool, resp, nil
}

// DeleteConnectionPool will delete a user within the Managed database. All data will be permanently lost.
func (d *DatabaseServiceHandler) DeleteConnectionPool(ctx context.Context, databaseID, poolName string) error {
	uri := fmt.Sprintf("%s/%s/connection-pools/%s", databasePath, databaseID, poolName)

	req, err := d.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	_, err = d.client.DoWithContext(ctx, req, nil)
	return err
}

// ListAdvancedOptions retrieves all connection pools within your PostgreSQL Managed Database.
func (d *DatabaseServiceHandler) ListAdvancedOptions(ctx context.Context, databaseID string) (*DatabaseAdvancedOptions, []AvailableOption, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/advanced-options", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	databaseAdvancedOptions := new(databaseAdvancedOptionsBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseAdvancedOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	return databaseAdvancedOptions.ConfiguredOptions, databaseAdvancedOptions.AvailableOptions, resp, nil
}

// UpdateAdvancedOptions will update a connection pool within the PostgreSQL Managed Database with the given parameters
func (d *DatabaseServiceHandler) UpdateAdvancedOptions(ctx context.Context, databaseID string, databaseAdvancedOptionsReq *DatabaseAdvancedOptions) (*DatabaseAdvancedOptions, []AvailableOption, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/advanced-options", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPut, uri, databaseAdvancedOptionsReq)
	if err != nil {
		return nil, nil, nil, err
	}

	databaseAdvancedOptions := new(databaseAdvancedOptionsBase)
	resp, err := d.client.DoWithContext(ctx, req, databaseAdvancedOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	return databaseAdvancedOptions.ConfiguredOptions, databaseAdvancedOptions.AvailableOptions, resp, nil
}

// ListAvailableVersions retrieves all available version upgrades for your Managed Database.
func (d *DatabaseServiceHandler) ListAvailableVersions(ctx context.Context, databaseID string) ([]string, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/version-upgrade", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	databaseVersions := new(DatabaseAvailableVersions)
	resp, err := d.client.DoWithContext(ctx, req, databaseVersions)
	if err != nil {
		return nil, nil, err
	}

	return databaseVersions.AvailableVersions, resp, nil
}

// StartVersionUpgrade will start a migration for the Managed Database using the given credentials.
func (d *DatabaseServiceHandler) StartVersionUpgrade(ctx context.Context, databaseID string, databaseVersionUpgradeReq *DatabaseVersionUpgradeReq) (string, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s/version-upgrade", databasePath, databaseID)

	req, err := d.client.NewRequest(ctx, http.MethodPost, uri, databaseVersionUpgradeReq)
	if err != nil {
		return "", nil, err
	}

	databaseVersionUpgrade := new(databaseMessage)
	resp, err := d.client.DoWithContext(ctx, req, databaseVersionUpgrade)
	if err != nil {
		return "", nil, err
	}

	return databaseVersionUpgrade.Message, resp, nil
}
