// Package database is used by the CLI to control databases
package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `Get commands available to database`
	example = `
	# Full example
	vultr-cli database
	`
	listLong    = `Get all databases on your Vultr account`
	listExample = `
	# Full example
	vultr-cli database list

	# Summarized view
	vultr-cli database list --summarize
	`
	createLong    = `Create a new Managed Database with specified plan, region, and database engine/version`
	createExample = `
	# Full example
	vultr-cli database create --database-engine="mysql" --database-engine-version="8" --region="ewr" \
	    --plan="vultr-dbaas-startup-cc-1-55-2" --label="example-db"

	# Full example with custom MySQL settings
	vultr-cli database create --database-engine="mysql" --database-engine-version="8" --region="ewr" \
	    --plan="vultr-dbaas-startup-cc-1-55-2" --label="example-db" --mysql-slow-query-log="true" \
		--mysql-long-query-time="2"
	`
	updateLong    = `Updates a Managed Database with the supplied information`
	updateExample = `
	# Full example
	vultr-cli database update --region="sea" --plan="vultr-dbaas-startup-cc-2-80-4"

	# Full example with custom MySQL settings
	vultr-cli database update --mysql-slow-query-log="true" --mysql-long-query-time="2"
	`
)

// NewCmdDatabase provides the CLI command for database functions
func NewCmdDatabase(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "database",
		Short:   "Commands to manage databases",
		Long:    long,
		Example: example,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Short:   "List databases",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			summarize, errSu := cmd.Flags().GetBool("summarize")
			if errSu != nil {
				return fmt.Errorf("error parsing flag 'summarize' for database list : %v", errSu)
			}

			dbs, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving database list : %v", err)
			}

			var data printer.ResourceOutput
			if summarize {
				data = &DBsSummaryPrinter{DBs: dbs, Meta: meta}
			} else {
				data = &DBsPrinter{DBs: dbs, Meta: meta}
			}

			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().BoolP("summarize", "", false, "(optional) Summarize the list output. One line per database")

	// Get
	get := &cobra.Command{
		Use:     "get <Database ID>",
		Short:   "Retrieve a database",
		Aliases: []string{"g"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving database : %v", err)
			}

			data := &DBPrinter{DB: db}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create database",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			engine, errEn := cmd.Flags().GetString("database-engine")
			if errEn != nil {
				return fmt.Errorf("error parsing flag 'database-engine' for database create : %v", errEn)
			}

			engineVersion, errEg := cmd.Flags().GetString("database-engine-version")
			if errEg != nil {
				return fmt.Errorf("error parsing flag 'database-engine-version' for database create : %v", errEg)
			}

			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for database create : %v", errRe)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing flag 'plan' for database create : %v", errPl)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for database create : %v", errLa)
			}

			// Optional
			tag, errTa := cmd.Flags().GetString("tag")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tag' for database create : %v", errTa)
			}

			vpc, errVp := cmd.Flags().GetString("vpc-id")
			if errVp != nil {
				return fmt.Errorf("error parsing flag 'vpc-id' for database create : %v", errVp)
			}

			maintenanceDOW, errMa := cmd.Flags().GetString("maintenance-dow")
			if errMa != nil {
				return fmt.Errorf("error parsing flag 'maintenance-dow' for database create : %v", errMa)
			}

			maintenanceTime, errMt := cmd.Flags().GetString("maintenance-time")
			if errMt != nil {
				return fmt.Errorf("error parsing flag 'maintenance-time' for database create : %v", errMt)
			}

			trustedIPs, errTr := cmd.Flags().GetStringSlice("trusted-ips")
			if errTr != nil {
				return fmt.Errorf("error parsing flag 'trusted-ips' for database create : %v", errTr)
			}

			mysqlSQLModes, errMy := cmd.Flags().GetStringSlice("mysql-sql-modes")
			if errMy != nil {
				return fmt.Errorf("error parsing flag 'mysql-sql-modes' for database create : %v", errMy)
			}

			mysqlRequirePrimaryKey, errMq := cmd.Flags().GetBool("mysql-require-primary-key")
			if errMq != nil {
				return fmt.Errorf("error parsing flag 'mysql-require-primary-key' for database create : %v", errMq)
			}

			mySQLSlowQueryLog, errMl := cmd.Flags().GetBool("mysql-slow-query-log")
			if errMl != nil {
				return fmt.Errorf("error parsing flag 'mysql-slow-query-log' for database create : %v", errMl)
			}

			mySQLLongQueryTime, errMt := cmd.Flags().GetInt("mysql-long-query-time")
			if errMt != nil {
				return fmt.Errorf("error parsing flag 'mysql-long-query-time' for database create : %v", errMt)
			}

			redisEvictionPolicy, errEe := cmd.Flags().GetString("redis-eviction-policy")
			if errEe != nil {
				return fmt.Errorf("error parsing flag 'redis-eviction-policy' for database create : %v", errEe)
			}

			o.CreateReq = &govultr.DatabaseCreateReq{
				DatabaseEngine:         engine,
				DatabaseEngineVersion:  engineVersion,
				Region:                 region,
				Plan:                   plan,
				Label:                  label,
				Tag:                    tag,
				VPCID:                  vpc,
				MaintenanceDOW:         maintenanceDOW,
				MaintenanceTime:        maintenanceTime,
				TrustedIPs:             trustedIPs,
				MySQLSQLModes:          mysqlSQLModes,
				MySQLRequirePrimaryKey: &mysqlRequirePrimaryKey,
				MySQLSlowQueryLog:      &mySQLSlowQueryLog,
				MySQLLongQueryTime:     mySQLLongQueryTime,
				RedisEvictionPolicy:    redisEvictionPolicy,
			}

			db, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating database : %v", err)
			}

			data := &DBPrinter{DB: db}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("database-engine", "e", "", "database engine for the new manaaged database")
	if err := create.MarkFlagRequired("database-engine"); err != nil {
		fmt.Printf("error marking database create 'database-engine' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("database-engine-version", "v", "", "database engine version for the new manaaged database")
	if err := create.MarkFlagRequired("database-engine-version"); err != nil {
		fmt.Printf("error marking database create 'database-engine-version' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("region", "r", "", "region id for the new managed database")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking database create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("plan", "p", "", "plan id for the new managed database")
	if err := create.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking database create 'plan' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("label", "l", "", "label for the new managed database")
	if err := create.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking database create 'label' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().String("tag", "t", "tag for the new managed database")
	create.Flags().String("vpc-id", "", "vpc id for the new managed database")
	create.Flags().String("maintenance-dow", "", "maintenance day of week for the new managed database")
	create.Flags().String("maintenance-time", "", "maintenance time for the new managed database")
	create.Flags().StringSlice(
		"trusted-ips",
		[]string{},
		"comma-separated list of trusted ip addresses for the new managed database",
	)
	create.Flags().StringSlice(
		"mysql-sql-modes",
		[]string{},
		"comma-separated list of sql modes for the new managed database",
	)
	create.Flags().Bool(
		"mysql-require-primary-key",
		true,
		"enable requiring primary keys for the new mysql managed database",
	)
	create.Flags().Bool(
		"mysql-slow-query-log",
		false,
		"enable slow query logging for the new mysql managed database",
	)
	create.Flags().Int(
		"mysql-long-query-time",
		0,
		"long query time for the new mysql managed database when slow query logging is enabled",
	)
	create.Flags().String("redis-eviction-policy", "", "eviction policy for the new redis managed database")

	// Update
	update := &cobra.Command{
		Use:     "update <Database ID>",
		Short:   "Update a database",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for database update : %v", errRe)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing flag 'plan' for database update : %v", errPl)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for database update : %v", errLa)
			}

			tag, errTa := cmd.Flags().GetString("tag")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tag' for database update : %v", errTa)
			}

			maintenanceDOW, errMa := cmd.Flags().GetString("maintenance-dow")
			if errMa != nil {
				return fmt.Errorf("error parsing flag 'maintenance-dow' for database update : %v", errMa)
			}

			maintenanceTime, errMt := cmd.Flags().GetString("maintenance-time")
			if errMt != nil {
				return fmt.Errorf("error parsing flag 'maintenance-time' for database update : %v", errMt)
			}

			clusterTimeZone, errTz := cmd.Flags().GetString("cluster-time-zone")
			if errTz != nil {
				return fmt.Errorf("error parsing flag 'cluster-time-zone' for database update : %v", errTz)
			}

			trustedIPs, errTr := cmd.Flags().GetStringSlice("trusted-ips")
			if errTr != nil {
				return fmt.Errorf("error parsing flag 'trusted-ips' for database update : %v", errTr)
			}

			mysqlSQLModes, errMy := cmd.Flags().GetStringSlice("mysql-sql-modes")
			if errMy != nil {
				return fmt.Errorf("error parsing flag 'mysql-sql-modes' for database update : %v", errMy)
			}

			mySQLLongQueryTime, errMt := cmd.Flags().GetInt("mysql-long-query-time")
			if errMt != nil {
				return fmt.Errorf("error parsing flag 'mysql-long-query-time' for database update : %v", errMt)
			}

			redisEvictionPolicy, errEe := cmd.Flags().GetString("redis-eviction-policy")
			if errEe != nil {
				return fmt.Errorf("error parsing flag 'redis-eviction-policy' for database update : %v", errEe)
			}

			o.UpdateReq = &govultr.DatabaseUpdateReq{}

			if cmd.Flags().Changed("region") {
				o.UpdateReq.Region = region
			}

			if cmd.Flags().Changed("plan") {
				o.UpdateReq.Plan = plan
			}

			if cmd.Flags().Changed("label") {
				o.UpdateReq.Label = label
			}

			if cmd.Flags().Changed("tag") {
				o.UpdateReq.Tag = tag
			}

			if cmd.Flags().Changed("maintenance-dow") {
				o.UpdateReq.MaintenanceDOW = maintenanceDOW
			}

			if cmd.Flags().Changed("maintenance-time") {
				o.UpdateReq.MaintenanceTime = maintenanceTime
			}

			if cmd.Flags().Changed("cluster-time-zone") {
				o.UpdateReq.ClusterTimeZone = clusterTimeZone
			}

			if cmd.Flags().Changed("trusted-ips") {
				o.UpdateReq.TrustedIPs = trustedIPs
			}

			if cmd.Flags().Changed("mysql-sql-modes") {
				o.UpdateReq.MySQLSQLModes = mysqlSQLModes
			}

			if cmd.Flags().Changed("mysql-long-query-time") {
				o.UpdateReq.MySQLLongQueryTime = mySQLLongQueryTime
			}

			if cmd.Flags().Changed("redis-eviction-policy") {
				o.UpdateReq.RedisEvictionPolicy = redisEvictionPolicy
			}

			db, err := o.update()
			if err != nil {
				return fmt.Errorf("error updating database : %v", err)
			}

			data := &DBPrinter{DB: db}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	update.Flags().StringP("region", "r", "", "region id for the managed database")
	update.Flags().StringP("plan", "p", "", "plan id for the managed database")
	update.Flags().StringP("label", "l", "", "label for the managed database")
	update.Flags().StringP("tag", "t", "", "tag for the managed database")
	update.Flags().String("vpc-id", "", "vpc id for the managed database")
	update.Flags().String("maintenance-dow", "", "maintenance day of week for the managed database")
	update.Flags().String("maintenance-time", "", "maintenance time for the managed database")
	update.Flags().String("cluster-time-zone", "", "configured time zone for the managed database")
	update.Flags().StringSlice(
		"trusted-ips",
		[]string{},
		"comma-separated list of trusted ip addresses for the managed database",
	)
	update.Flags().StringSlice(
		"mysql-sql-modes",
		[]string{},
		"comma-separated list of sql modes for the managed database",
	)
	update.Flags().Bool(
		"mysql-require-primary-key",
		true,
		"enable requiring primary keys for the mysql managed database",
	)
	update.Flags().Bool("mysql-slow-query-log", false, "enable slow query logging for the mysql managed database")
	update.Flags().Int(
		"mysql-long-query-time",
		0,
		"long query time for the mysql managed database when slow query logging is enabled",
	)
	update.Flags().String("redis-eviction-policy", "", "eviction policy for the redis managed database")

	// Delete
	del := &cobra.Command{
		Use:     "delete <Database ID>",
		Short:   "Delete a database",
		Aliases: []string{"destroy", "d"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting database : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Database has been deleted"), nil)

			return nil
		},
	}

	// Plan
	plan := &cobra.Command{
		Use:   "plan",
		Short: "Commands to access database plans",
	}

	// Plan List
	planList := &cobra.Command{
		Use:     "list",
		Short:   "List database plans",
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			plans, meta, err := o.listPlans()
			if err != nil {
				return fmt.Errorf("error retrieving database plans : %v", err)
			}

			data := &PlansPrinter{Plans: plans, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	plan.AddCommand(
		planList,
	)

	// User
	user := &cobra.Command{
		Use:   "user",
		Short: "Commands to handle database users",
	}

	// User List
	userList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List database users",
		RunE: func(cmd *cobra.Command, args []string) error {
			us, meta, err := o.listUsers()
			if err != nil {
				return fmt.Errorf("error retrieving database users : %v", err)
			}

			data := &UsersPrinter{Users: us, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// User Get
	userGet := &cobra.Command{
		Use:   "get <Database ID> <User Name>",
		Short: "Get a database user",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a user name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			us, err := o.getUser()
			if err != nil {
				return fmt.Errorf("error retrieving database user : %v", err)
			}

			data := &UserPrinter{User: us}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// User Create
	userCreate := &cobra.Command{
		Use:   "create <Database ID>",
		Short: "Create a database user",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			username, errUs := cmd.Flags().GetString("username")
			if errUs != nil {
				return fmt.Errorf("error parsing flag 'username' for database user create : %v", errUs)
			}

			password, errPa := cmd.Flags().GetString("password")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'password' for database user create : %v", errPa)
			}

			encryption, errEn := cmd.Flags().GetString("encryption")
			if errEn != nil {
				return fmt.Errorf("error parsing flag 'encryption' for database user create : %v", errEn)
			}

			o.UserCreateReq = &govultr.DatabaseUserCreateReq{
				Username:   username,
				Password:   password,
				Encryption: encryption,
			}

			us, err := o.createUser()
			if err != nil {
				return fmt.Errorf("error creating database user : %v", err)
			}

			data := &UserPrinter{User: us}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}
	userCreate.Flags().StringP("username", "u", "", "username for the new manaaged database user")
	userCreate.Flags().StringP(
		"password",
		"p",
		"",
		"password for the new manaaged database user (omit or leave empty to generate a random secure password)",
	)
	userCreate.Flags().StringP("encryption", "e", "", "encryption type for the new managed database user (MySQL only)")

	// User Update
	userUpdate := &cobra.Command{
		Use:   "update <Database ID> <User Name>",
		Short: "Update a database user",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a user name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			password, errPa := cmd.Flags().GetString("password")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'password' for database user create : %v", errPa)
			}

			o.UserUpdateReq = &govultr.DatabaseUserUpdateReq{
				Password: password,
			}

			us, err := o.updateUser()
			if err != nil {
				return fmt.Errorf("error updating database user : %v", err)
			}

			data := &UserPrinter{User: us}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	userUpdate.Flags().StringP(
		"password",
		"p",
		"",
		"password for the new manaaged database user (omit or leave empty to generate a random secure password)",
	)

	// User Delete
	userDelete := &cobra.Command{
		Use:   "delete <Database ID> <User Name>",
		Short: "Delete a database user",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a user name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.delUser(); err != nil {
				return fmt.Errorf("error deleting database user : %v", err)
			}

			o.Base.Printer.Display(printer.Info("User deleted"), nil)

			return nil
		},
	}

	// User ACL
	userACL := &cobra.Command{
		Use:   "acl",
		Short: "Commands to manage database user access control (Redis only)",
	}

	// User ACL Update
	userACLUpdate := &cobra.Command{
		Use:   "update <Database ID> <User Name>",
		Short: "Update a database user ACL",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a user name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			categories, errCa := cmd.Flags().GetStringSlice("redis-acl-categories")
			if errCa != nil {
				return fmt.Errorf("error parsing flag 'redis-acl-categories' for database user create : %v", errCa)
			}

			channels, errCh := cmd.Flags().GetStringSlice("redis-acl-channels")
			if errCh != nil {
				return fmt.Errorf("error parsing flag 'redis-acl-channels' for database user create : %v", errCh)
			}

			commands, errCo := cmd.Flags().GetStringSlice("redis-acl-commands")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'redis-acl-commands' for database user create : %v", errCo)
			}

			keys, errKe := cmd.Flags().GetStringSlice("redis-acl-keys")
			if errKe != nil {
				return fmt.Errorf("error parsing flag 'redis-acl-keys' for database user create : %v", errKe)
			}

			o.UserUpdateACLReq = &govultr.DatabaseUserACLReq{}

			if cmd.Flags().Changed("redis-acl-categories") {
				o.UserUpdateACLReq.RedisACLCategories = &categories
			}

			if cmd.Flags().Changed("redis-acl-channels") {
				o.UserUpdateACLReq.RedisACLChannels = &channels
			}

			if cmd.Flags().Changed("redis-acl-commands") {
				o.UserUpdateACLReq.RedisACLCommands = &commands
			}

			if cmd.Flags().Changed("redis-acl-keys") {
				o.UserUpdateACLReq.RedisACLKeys = &keys
			}

			us, err := o.updateUserACL()
			if err != nil {
				return fmt.Errorf("error updating database user acl : %v", err)
			}

			data := &UserPrinter{User: us}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	userACLUpdate.Flags().StringSlice(
		"redis-acl-categories",
		[]string{},
		"list of rules for command categories",
	)
	userACLUpdate.Flags().StringSlice(
		"redis-acl-channels",
		[]string{},
		"list of publish/subscribe channel patterns",
	)
	userACLUpdate.Flags().StringSlice(
		"redis-acl-commands",
		[]string{},
		"list of rules for individual commands",
	)
	userACLUpdate.Flags().StringSlice(
		"redis-acl-keys",
		[]string{},
		"list of key access rules",
	)

	userACLUpdate.MarkFlagsOneRequired(
		"redis-acl-categories",
		"redis-acl-channels",
		"redis-acl-commands",
		"redis-acl-keys",
	)

	userACL.AddCommand(
		userACLUpdate,
	)

	user.AddCommand(
		userList,
		userGet,
		userCreate,
		userUpdate,
		userDelete,
		userACL,
	)

	// Logical Database
	db := &cobra.Command{
		Use:   "db",
		Short: "Commands to handle database logical DBs",
	}

	// Logical DB List
	dbList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List logical databases",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dbs, meta, err := o.listDBs()
			if err != nil {
				return fmt.Errorf("error retrieving logical databases: %v", err)
			}

			data := &LogicalDBsPrinter{DBs: dbs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Logical DB Create
	dbCreate := &cobra.Command{
		Use:   "create <Database ID>",
		Short: "Create a logical database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for logical database create : %v", errNa)
			}

			o.DBCreateReq = &govultr.DatabaseDBCreateReq{
				Name: name,
			}

			db, err := o.createDB()
			if err != nil {
				return fmt.Errorf("error creating a logical database : %v", err)
			}

			data := &LogicalDBPrinter{DB: db}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	dbCreate.Flags().StringP("name", "n", "", "name of the new logical database within the manaaged database")
	if err := dbCreate.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking logical database create 'name' flag required: %v", err)
		os.Exit(1)
	}

	// Logical DB Delete
	dbDel := &cobra.Command{
		Use:   "delete <Database ID> <DB Name>",
		Short: "Delete a logical database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a DB name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.delDB(); err != nil {
				return fmt.Errorf("error deleting logical database : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Logical DB deleted"), nil)

			return nil
		},
	}

	db.AddCommand(
		dbList,
		dbCreate,
		dbDel,
	)

	// Usage
	usage := &cobra.Command{
		Use:   "usage",
		Short: "Commands to display database usage information",
	}

	// Usage Get
	usageGet := &cobra.Command{
		Use:   "get <Database ID>",
		Short: "Get database usage",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			us, err := o.getUsage()
			if err != nil {
				return fmt.Errorf("error retrieving database usage  : %v", err)
			}

			data := &UsagePrinter{Usage: us}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	usage.AddCommand(
		usageGet,
	)

	// Maintenance
	maintenance := &cobra.Command{
		Use:   "maintenance",
		Short: "Commands to handle database maintenance updates",
	}

	// Maintenance List
	maintenanceList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List maintenance updates for a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			upds, err := o.listMaintUpdates()
			if err != nil {
				return fmt.Errorf("error retrieving database maintenance updates : %v", err)
			}

			data := &UpdatesPrinter{Updates: upds}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Maintenance Start
	maintenanceStart := &cobra.Command{
		Use:   "start <Database ID>",
		Short: "Start database maintenance update",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			message, err := o.startMaintUpdate()
			if err != nil {
				return fmt.Errorf("error starting database maintenance update: %v", err)
			}

			o.Base.Printer.Display(printer.Info(message), nil)

			return nil
		},
	}

	maintenance.AddCommand(
		maintenanceList,
		maintenanceStart,
	)

	// Alert
	alert := &cobra.Command{
		Use:   "alert",
		Short: "Commands to handle database alerts",
	}

	// Alert List
	alertList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List database alerts",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			period, errPe := cmd.Flags().GetString("period")
			if errPe != nil {
				return fmt.Errorf("error parsing flag 'period' for alert list : %v", errPe)
			}

			o.AlertsReq = &govultr.DatabaseListAlertsReq{
				Period: period,
			}

			als, err := o.listAlerts()
			if err != nil {
				return fmt.Errorf("error retrieving database alerts : %v", err)
			}

			data := &AlertsPrinter{Alerts: als}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	alertList.Flags().StringP(
		"period",
		"p",
		"",
		"period (day, week, month, year) for viewing service alerts for a manaaged database",
	)
	if err := alertList.MarkFlagRequired("period"); err != nil {
		fmt.Printf("error marking alert list 'period' flag required: %v", err)
		os.Exit(1)
	}

	alert.AddCommand(
		alertList,
	)

	// Migration
	migration := &cobra.Command{
		Use:   "migration",
		Short: "Commands to handle database migrations",
	}

	// Migration Get
	migrationGet := &cobra.Command{
		Use:   "get <Database ID>",
		Short: "Get migration status of a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mig, err := o.getMigrationStatus()
			if err != nil {
				return fmt.Errorf("error retrieving database migration status : %v", err)
			}

			data := &MigrationPrinter{Migration: mig}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Migration Start
	migrationStart := &cobra.Command{
		Use:   "start <Database ID>",
		Short: "Get migration status of a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, errHo := cmd.Flags().GetString("host")
			if errHo != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errHo)
			}

			port, errPo := cmd.Flags().GetInt("port")
			if errPo != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errPo)
			}

			username, errUs := cmd.Flags().GetString("username")
			if errUs != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errUs)
			}

			password, errPa := cmd.Flags().GetString("password")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errPa)
			}

			database, errDa := cmd.Flags().GetString("database")
			if errDa != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errDa)
			}

			ignored, errIg := cmd.Flags().GetString("ignored-dbs")
			if errIg != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errIg)
			}

			ssl, errSs := cmd.Flags().GetBool("ssl")
			if errSs != nil {
				return fmt.Errorf("error parsing flag 'encryption' for migration start : %v", errSs)
			}

			o.MigrationReq = &govultr.DatabaseMigrationStartReq{
				Host:             host,
				Port:             port,
				Username:         username,
				Password:         password,
				Database:         database,
				IgnoredDatabases: ignored,
				SSL:              &ssl,
			}

			mig, err := o.startMigration()
			if err != nil {
				return fmt.Errorf("error retrieving database migration status : %v", err)
			}

			data := &MigrationPrinter{Migration: mig}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	migrationStart.Flags().String("host", "", "source host for the manaaged database migration")
	migrationStart.Flags().Int("port", 0, "source port for the manaaged database migration")
	migrationStart.Flags().String(
		"username",
		"",
		"source username for the manaaged database migration (uses `default` for Redis if omitted)",
	)
	migrationStart.Flags().String("password", "", "source password for the manaaged database migration")
	migrationStart.Flags().String(
		"database",
		"",
		"source database for the manaaged database migration (MySQL/PostgreSQL only)",
	)
	migrationStart.Flags().String(
		"ignored-dbs",
		"",
		"comma-separated list of ignored databases for the manaaged database migration (MySQL/PostgreSQL only)",
	)
	migrationStart.Flags().Bool("ssl", true, "source ssl requirement for the manaaged database migration")

	// Migration Detach
	migrationDetach := &cobra.Command{
		Use:   "detach <Database ID>",
		Short: "Detach a migration from a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.detachMigration(); err != nil {
				return fmt.Errorf("error detaching migration from database : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Migration detached"), nil)

			return nil
		},
	}

	migration.AddCommand(
		migrationGet,
		migrationStart,
		migrationDetach,
	)

	// Read Replica
	readReplica := &cobra.Command{
		Use:   "read-replica",
		Short: "Commands to handle database read replicas",
	}

	// Read Replica Add
	readReplicaCreate := &cobra.Command{
		Use:   "create <Database ID>",
		Short: "Create a read-only replica of a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for read-replica create : %v", errRe)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for read-replica create : %v", errLa)
			}

			o.ReadReplicaCreateReq = &govultr.DatabaseAddReplicaReq{
				Region: region,
				Label:  label,
			}

			rr, err := o.createReadReplica()
			if err != nil {
				return fmt.Errorf("error creating database read replica: %v", err)
			}

			data := &DBPrinter{DB: rr}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	readReplicaCreate.Flags().StringP("region", "r", "", "region id for the new managed database read replica")
	if err := readReplicaCreate.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking read replica create 'region' flag required: %v", err)
		os.Exit(1)
	}

	readReplicaCreate.Flags().StringP("label", "l", "", "label for the new managed database read replica")
	if err := readReplicaCreate.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking read replica create 'label' flag required: %v", err)
		os.Exit(1)
	}

	// Read Replica Promote
	readReplicaPromote := &cobra.Command{
		Use:   "promote <Database ID>",
		Short: "Promote a read-only replica of a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.promoteReadReplica(); err != nil {
				return fmt.Errorf("error promoting database read replica: %v", err)
			}

			o.Base.Printer.Display(printer.Info("Read replica has been promoted"), nil)

			return nil
		},
	}

	readReplica.AddCommand(
		readReplicaCreate,
		readReplicaPromote,
	)

	// Backup
	backup := &cobra.Command{
		Use:   "backup",
		Short: "Commands to handle database backups, restores and forks",
	}

	// Backup Get
	backupGet := &cobra.Command{
		Use:   "get <Database ID>",
		Short: "Get get latest and oldest database backup",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bk, err := o.getBackup()
			if err != nil {
				return fmt.Errorf("error retrieving database backups : %v", err)
			}

			data := &BackupPrinter{Backup: bk}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Backup Restore
	backupRestore := &cobra.Command{
		Use:   "restore <Database ID>",
		Short: "Restore a database backup",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for backup restore : %v", errLa)
			}

			rtype, errRt := cmd.Flags().GetString("type")
			if errRt != nil {
				return fmt.Errorf("error parsing flag 'type' for backup restore : %v", errRt)
			}

			date, errDa := cmd.Flags().GetString("date")
			if errDa != nil {
				return fmt.Errorf("error parsing flag 'date' for backup restore : %v", errDa)
			}

			time, errTi := cmd.Flags().GetString("time")
			if errTi != nil {
				return fmt.Errorf("error parsing flag 'time' for backup restore : %v", errTi)
			}

			o.BackupReq = &govultr.DatabaseBackupRestoreReq{
				Label: label,
				Type:  rtype,
				Date:  date,
				Time:  time,
			}

			bk, err := o.restoreBackup()
			if err != nil {
				return fmt.Errorf("error restoring database from backup : %v", err)
			}

			data := &DBPrinter{DB: bk}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	backupRestore.Flags().String("label", "", "label for the new managed database restored from backup")
	if err := backupRestore.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking backup restore 'label' flag required: %v", err)
		os.Exit(1)
	}

	backupRestore.Flags().String(
		"type",
		"",
		"restoration type: `pitr` for point-in-time recovery or `basebackup` for latest backup (default)",
	)
	backupRestore.Flags().String("date", "", "backup date to use for point-in-time recovery")
	backupRestore.Flags().String("time", "", "backup time to use for point-in-time recovery")

	// Backup Fork
	backupFork := &cobra.Command{
		Use:   "fork <Database ID>",
		Short: "Fork a database from backup",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errDa := cmd.Flags().GetString("region")
			if errDa != nil {
				return fmt.Errorf("error parsing flag 'region' for backup fork: %v", errDa)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing flag 'time' for backup fork: %v", errPl)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for backup fork : %v", errLa)
			}

			rtype, errRt := cmd.Flags().GetString("type")
			if errRt != nil {
				return fmt.Errorf("error parsing flag 'type' for backup fork: %v", errRt)
			}

			date, errDa := cmd.Flags().GetString("date")
			if errDa != nil {
				return fmt.Errorf("error parsing flag 'date' for backup fork: %v", errDa)
			}

			time, errTi := cmd.Flags().GetString("time")
			if errTi != nil {
				return fmt.Errorf("error parsing flag 'time' for backup fork: %v", errTi)
			}

			o.ForkReq = &govultr.DatabaseForkReq{
				Label:  label,
				Region: region,
				Plan:   plan,
				Type:   rtype,
				Date:   date,
				Time:   time,
			}

			db, err := o.fork()
			if err != nil {
				return fmt.Errorf("error forking database from backup : %v", err)
			}

			data := &DBPrinter{DB: db}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	backupFork.Flags().String("label", "", "label for the new managed database forked from the backup")
	if err := backupFork.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking backup fork 'label' flag required: %v", err)
		os.Exit(1)
	}

	backupFork.Flags().String("region", "", "region id for the new managed database forked from the backup")
	if err := backupFork.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking backup fork 'region' flag required: %v", err)
		os.Exit(1)
	}

	backupFork.Flags().String("plan", "", "plan id for the new managed database forked from the backup")
	if err := backupFork.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking backup fork 'label' flag required: %v", err)
		os.Exit(1)
	}

	backupFork.Flags().String(
		"type",
		"",
		"restoration type: `pitr` for point-in-time recovery or `basebackup` for latest backup (default)",
	)
	backupFork.Flags().String("date", "", "backup date to use for point-in-time recovery")
	backupFork.Flags().String("time", "", "backup time to use for point-in-time recovery")

	backup.AddCommand(
		backupGet,
		backupRestore,
		backupFork,
	)

	// Connection Pool
	connectionPool := &cobra.Command{
		Use:   "connection-pool",
		Short: "Commands to handle PostgreSQL database connection pools",
	}

	// Connection Pool List
	connectionPoolList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List connection pools within a PostgreSQL managed database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cns, pools, meta, err := o.listConnectionPools()
			if err != nil {
				return fmt.Errorf("error retrieving connection pool data : %v", err)
			}

			data := &ConnectionsPrinter{Connections: cns, ConnectionPools: pools, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Connection Pool Get
	connectionPoolGet := &cobra.Command{
		Use:   "get <Database ID> <Pool Name>",
		Short: "Get a database connection pool",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a pool name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cnp, err := o.getConnectionPool()
			if err != nil {
				return fmt.Errorf("error retrieving connection pool: %v", err)
			}

			data := &ConnectionPoolPrinter{ConnectionPool: cnp}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Connection Pool Create
	connectionPoolCreate := &cobra.Command{
		Use:   "create <Database ID>",
		Short: "Create a database connection pool",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing flag 'name' for connection pool create : %v", errNa)
			}

			database, errDa := cmd.Flags().GetString("database")
			if errDa != nil {
				return fmt.Errorf("error parsing flag 'database' for connection pool create : %v", errDa)
			}

			username, errUs := cmd.Flags().GetString("username")
			if errUs != nil {
				return fmt.Errorf("error parsing flag 'username' for connection pool create : %v", errUs)
			}

			mode, errMo := cmd.Flags().GetString("mode")
			if errMo != nil {
				return fmt.Errorf("error parsing flag 'mode' for connection pool create : %v", errMo)
			}

			size, errSi := cmd.Flags().GetInt("size")
			if errSi != nil {
				return fmt.Errorf("error parsing flag 'size' for connection pool create : %v", errSi)
			}

			o.ConnectionPoolCreateReq = &govultr.DatabaseConnectionPoolCreateReq{
				Name:     name,
				Database: database,
				Username: username,
				Mode:     mode,
				Size:     size,
			}

			cnp, err := o.createConnectionPool()
			if err != nil {
				return fmt.Errorf("error creating connection pool: %v", err)
			}

			data := &ConnectionPoolPrinter{ConnectionPool: cnp}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	connectionPoolCreate.Flags().StringP("name", "n", "", "name for the new managed database connection pool")
	if err := connectionPoolCreate.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking connection pool create 'name' flag required: %v", err)
		os.Exit(1)
	}

	connectionPoolCreate.Flags().StringP("database", "d", "", "database for the new managed database connection pool")
	if err := connectionPoolCreate.MarkFlagRequired("database"); err != nil {
		fmt.Printf("error marking connection pool create 'database' flag required: %v", err)
		os.Exit(1)
	}

	connectionPoolCreate.Flags().StringP("username", "u", "", "username for the new managed database connection pool")
	if err := connectionPoolCreate.MarkFlagRequired("username"); err != nil {
		fmt.Printf("error marking connection pool create 'username' flag required: %v", err)
		os.Exit(1)
	}

	connectionPoolCreate.Flags().StringP("mode", "m", "", "mode for the new managed database connection pool")
	if err := connectionPoolCreate.MarkFlagRequired("mode"); err != nil {
		fmt.Printf("error marking connection pool create 'mode' flag required: %v", err)
		os.Exit(1)
	}

	connectionPoolCreate.Flags().IntP("size", "s", 0, "size for the new managed database connection pool")
	if err := connectionPoolCreate.MarkFlagRequired("size"); err != nil {
		fmt.Printf("error marking connection pool create 'size' flag required: %v", err)
		os.Exit(1)
	}

	// Connection Pool Update
	connectionPoolUpdate := &cobra.Command{
		Use:   "update <Database ID> <Pool Name>",
		Short: "Update a database connection pool",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID pool name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			database, errDa := cmd.Flags().GetString("database")
			if errDa != nil {
				return fmt.Errorf("error parsing flag 'database' for connection pool update : %v", errDa)
			}

			username, errUs := cmd.Flags().GetString("username")
			if errUs != nil {
				return fmt.Errorf("error parsing flag 'username' for connection pool update : %v", errUs)
			}

			mode, errMo := cmd.Flags().GetString("mode")
			if errMo != nil {
				return fmt.Errorf("error parsing flag 'mode' for connection pool update : %v", errMo)
			}

			size, errSi := cmd.Flags().GetInt("size")
			if errSi != nil {
				return fmt.Errorf("error parsing flag 'size' for connection pool update : %v", errSi)
			}

			o.ConnectionPoolCreateReq = &govultr.DatabaseConnectionPoolCreateReq{}

			if cmd.Flags().Changed("database") {
				o.ConnectionPoolCreateReq.Database = database
			}

			if cmd.Flags().Changed("username") {
				o.ConnectionPoolCreateReq.Username = username
			}

			if cmd.Flags().Changed("mode") {
				o.ConnectionPoolCreateReq.Mode = mode
			}

			if cmd.Flags().Changed("size") {
				o.ConnectionPoolCreateReq.Size = size
			}

			cnp, err := o.updateConnectionPool()
			if err != nil {
				return fmt.Errorf("error updating connection pool : %v", err)
			}

			data := &ConnectionPoolPrinter{ConnectionPool: cnp}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	connectionPoolUpdate.Flags().StringP("database", "d", "", "database for the managed database connection pool")
	connectionPoolUpdate.Flags().StringP("username", "u", "", "username for the managed database connection pool")
	connectionPoolUpdate.Flags().StringP("mode", "m", "", "mode for the managed database connection pool")
	connectionPoolUpdate.Flags().IntP("size", "s", 0, "size for the managed database connection pool")

	connectionPoolUpdate.MarkFlagsOneRequired(
		"database",
		"username",
		"mode",
		"size",
	)

	// Connection Pool Delete
	connectionPoolDelete := &cobra.Command{
		Use:   "delete <Database ID> <Pool Name>",
		Short: "Delete a database connection pool",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("please provide a database ID and a pool name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.delConnectionPool(); err != nil {
				return fmt.Errorf("error deleting connection pool : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Connection pool has been deleted"), nil)

			return nil
		},
	}

	connectionPool.AddCommand(
		connectionPoolList,
		connectionPoolGet,
		connectionPoolCreate,
		connectionPoolUpdate,
		connectionPoolDelete,
	)

	// Advanced Option
	advancedOption := &cobra.Command{
		Use:   "advanced-option",
		Short: "Commands to handle managed database advanced options",
	}

	// Advanced Option List
	advancedOptionList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List all available and configured advanced options for a managed database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cur, avail, err := o.listAdvancedOptions()
			if err != nil {
				return fmt.Errorf("error retrieving database options : %v", err)
			}

			data := &AdvancedOptionsPrinter{Configured: cur, Available: avail}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Advanced Option Update
	advancedOptionUpdate := &cobra.Command{
		Use:   "update <Database ID>",
		Short: "Update advanced options for a managed database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			autovacuumAnalyzeScaleFactor, errAu := cmd.Flags().GetFloat32("autovacuum-analyze-scale-factor")
			if errAu != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-analyze-scale-factor' for advanced options update : %v",
					errAu,
				)
			}

			autovacuumAnalyzeThreshold, errAt := cmd.Flags().GetInt("autovacuum-analyze-threshold")
			if errAt != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-analyze-threshold' for advanced options update : %v",
					errAt,
				)
			}

			autovacuumFreezeMaxAge, errAo := cmd.Flags().GetInt("autovacuum-freeze-max-age")
			if errAo != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-freeze-max-age' for advanced options update : %v",
					errAo,
				)
			}

			autovacuumMaxWorkers, errAv := cmd.Flags().GetInt("autovacuum-max-workers")
			if errAv != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-max-workers' for advanced options update : %v",
					errAv,
				)
			}

			autovacuumNaptime, errAa := cmd.Flags().GetInt("autovacuum-naptime")
			if errAa != nil {
				return fmt.Errorf("error parsing flag 'autovacuum-naptime' for advanced options update : %v", errAa)
			}

			autovacuumVacuumCostDelay, errAc := cmd.Flags().GetInt("autovacuum-vacuum-cost-delay")
			if errAc != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-vacuum-cost-delay' for advanced options update : %v",
					errAc,
				)
			}

			autovacuumVacuumCostLimit, errAm := cmd.Flags().GetInt("autovacuum-vacuum-cost-limit")
			if errAm != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-vacuum-cost-limit' for advanced options update : %v",
					errAm,
				)
			}

			autovacuumVacuumScaleFactor, errAb := cmd.Flags().GetFloat32("autovacuum-vacuum-scale-factor")
			if errAb != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-vacuum-scale-factor' for advanced options update : %v",
					errAb,
				)
			}

			autovacuumVacuumThreshold, errAz := cmd.Flags().GetInt("autovacuum-vacuum-threshold")
			if errAz != nil {
				return fmt.Errorf(
					"error parsing flag 'autovacuum-vacuum-threshold' for advanced options update : %v",
					errAz,
				)
			}

			bgwriterDelay, errBg := cmd.Flags().GetInt("bgwriter-delay")
			if errBg != nil {
				return fmt.Errorf("error parsing flag 'bgwriter-delay' for advanced options update : %v", errBg)
			}

			bgwriterFlushAfter, errBw := cmd.Flags().GetInt("bgwriter-flush-after")
			if errBw != nil {
				return fmt.Errorf("error parsing flag 'bgwriter-flush-after' for advanced options update : %v", errBw)
			}

			bgwriterLruMaxpages, errBr := cmd.Flags().GetInt("bgwriter-lru-maxpages")
			if errBr != nil {
				return fmt.Errorf("error parsing flag 'bgwriter-lru-maxpages' for advanced options update : %v", errBr)
			}

			bgwriterLruMultiplier, errBi := cmd.Flags().GetFloat32("bgwriter-lru-multiplier")
			if errBi != nil {
				return fmt.Errorf(
					"error parsing flag 'bgwriter-lru-multiplier' for advanced options update : %v",
					errBi,
				)
			}

			connectTimeout, errCn := cmd.Flags().GetInt("connect-timeout")
			if errCn != nil {
				return fmt.Errorf("error parsing flag 'connect-timeout' for advanced options update : %v", errCn)
			}

			deadlockTimeout, errDe := cmd.Flags().GetInt("deadlock-timeout")
			if errDe != nil {
				return fmt.Errorf("error parsing flag 'deadlock-timeout' for advanced options update : %v", errDe)
			}

			defaultToastCompression, errDf := cmd.Flags().GetString("default-toast-compression")
			if errDf != nil {
				return fmt.Errorf(
					"error parsing flag 'default-toast-compression' for advanced options update : %v",
					errDf,
				)
			}

			groupConcatMaxLen, errGr := cmd.Flags().GetInt("group-concat-max-len")
			if errGr != nil {
				return fmt.Errorf("error parsing flag 'group-concat-max-len' for advanced options update : %v", errGr)
			}

			idleInTransactionSessionTimeout, errIl := cmd.Flags().GetInt("idle-in-transaction-session-timeout")
			if errIl != nil {
				return fmt.Errorf(
					"error parsing flag 'idle-in-transaction-session-timeout' for advanced options update : %v",
					errIl,
				)
			}

			innoDBChangeBufferMaxSize, errIcb := cmd.Flags().GetInt("innodb-change-buffer-max-size")
			if errIcb != nil {
				return fmt.Errorf("error parsing flag 'innodb-change-buffer-max-size' for advanced options update : %v", errIcb)
			}

			innoDBFlushNeighbors, errIfn := cmd.Flags().GetInt("innodb-flush-neighbors")
			if errIfn != nil {
				return fmt.Errorf("error parsing flag 'innodb-flush-neighbors' for advanced options update : %v", errIfn)
			}

			innoDBFTMinTokenSize, errImts := cmd.Flags().GetInt("innodb-ft-min-token-size")
			if errImts != nil {
				return fmt.Errorf("error parsing flag 'innodb-ft-min-token-size' for advanced options update : %v", errImts)
			}

			innoDBFTServerStopwordTable, errIsst := cmd.Flags().GetString("innodb-ft-server-stopword-table")
			if errIsst != nil {
				return fmt.Errorf("error parsing flag 'innodb-ft-server-stopword-table' for advanced options update : %v", errIsst)
			}

			innoDBLockWaitTimeout, errIlwt := cmd.Flags().GetInt("innodb-lock-wait-timeout")
			if errIlwt != nil {
				return fmt.Errorf("error parsing flag 'innodb-lock-wait-timeout' for advanced options update : %v", errIlwt)
			}

			innoDBLogBufferSize, errIlbs := cmd.Flags().GetInt("innodb-log-buffer-size")
			if errIlbs != nil {
				return fmt.Errorf("error parsing flag 'innodb-log-buffer-size' for advanced options update : %v", errIlbs)
			}

			innoDBOnlineAlterLogMaxSize, errIoa := cmd.Flags().GetInt("innodb-online-alter-log-max-size")
			if errIoa != nil {
				return fmt.Errorf("error parsing flag 'innodb-online-alter-log-max-size' for advanced options update : %v", errIoa)
			}

			innoDBPrintAllDeadlocks, errIpad := cmd.Flags().GetBool("innodb-print-all-deadlocks")
			if errIpad != nil {
				return fmt.Errorf("error parsing flag 'innodb-print-all-deadlocks' for advanced options update : %v", errIpad)
			}

			innoDBReadIOThreads, errIrio := cmd.Flags().GetInt("innodb-read-io-threads")
			if errIrio != nil {
				return fmt.Errorf("error parsing flag 'innodb-read-io-threads' for advanced options update : %v", errIrio)
			}

			innoDBRollbackOnTimeout, errIrot := cmd.Flags().GetBool("innodb-rollback-on-timeout")
			if errIrot != nil {
				return fmt.Errorf("error parsing flag 'innodb-rollback-on-timeout' for advanced options update : %v", errIrot)
			}

			innoDBThreadConcurrency, errItc := cmd.Flags().GetInt("innodb-thread-concurrency")
			if errItc != nil {
				return fmt.Errorf("error parsing flag 'innodb-thread-concurrency' for advanced options update : %v", errItc)
			}

			innoDBWriteIOThreads, errIwio := cmd.Flags().GetInt("innodb-write-io-threads")
			if errIwio != nil {
				return fmt.Errorf("error parsing flag 'innodb-write-io-threads' for advanced options update : %v", errIwio)
			}

			interactiveTimeout, errIit := cmd.Flags().GetInt("interactive-timeout")
			if errIit != nil {
				return fmt.Errorf("error parsing flag 'interactive-timeout' for advanced options update : %v", errIit)
			}

			internalTmpMemStorageEngine, errItm := cmd.Flags().GetString("internal-tmp-mem-storage-engine")
			if errItm != nil {
				return fmt.Errorf("error parsing flag 'internal-tmp-mem-storage-engine' for advanced options update : %v", errItm)
			}

			jit, errJi := cmd.Flags().GetBool("jit")
			if errJi != nil {
				return fmt.Errorf("error parsing flag 'jit' for advanced options update : %v", errJi)
			}

			logAutovacuumMinDuration, errLo := cmd.Flags().GetInt("log-autovacuum-min-duration")
			if errLo != nil {
				return fmt.Errorf(
					"error parsing flag 'log-autovacuum-min-duration' for advanced options update : %v",
					errLo,
				)
			}

			logErrorVerbosity, errLg := cmd.Flags().GetString("log-error-verbosity")
			if errLg != nil {
				return fmt.Errorf("error parsing flag 'log-error-verbosity' for advanced options update : %v", errLg)
			}

			logLinePrefix, errLl := cmd.Flags().GetString("log-line-prefix")
			if errLl != nil {
				return fmt.Errorf("error parsing flag 'log-line-prefix' for advanced options update : %v", errLl)
			}

			logMinDurationStatement, errLm := cmd.Flags().GetInt("log-min-duration-statement")
			if errLm != nil {
				return fmt.Errorf(
					"error parsing flag 'log-min-duration-statement' for advanced options update : %v",
					errLm,
				)
			}

			maxAllowedPacket, errMap := cmd.Flags().GetInt("max-allowed-packet")
			if errMap != nil {
				return fmt.Errorf("error parsing flag 'max-allowed-packet' for advanced options update : %v", errMap)
			}

			maxFilesPerProcess, errMa := cmd.Flags().GetInt("max-files-per-process")
			if errMa != nil {
				return fmt.Errorf("error parsing flag 'max-files-per-process' for advanced options update : %v", errMa)
			}

			maxHeapTableSize, errMht := cmd.Flags().GetInt("max-heap-table-size")
			if errMht != nil {
				return fmt.Errorf("error parsing flag 'max-heap-table-size' for advanced options update : %v", errMht)
			}

			maxLocksPerTransaction, errMx := cmd.Flags().GetInt("max-locks-per-transaction")
			if errMx != nil {
				return fmt.Errorf(
					"error parsing flag 'max-locks-per-transaction' for advanced options update : %v",
					errMx,
				)
			}

			maxLogicalReplicationWorkers, errMl := cmd.Flags().GetInt("max-logical-replication-workers")
			if errMl != nil {
				return fmt.Errorf(
					"error parsing flag 'max-logical-replication-workers' for advanced options update : %v",
					errMl,
				)
			}

			maxParallelWorkers, errMo := cmd.Flags().GetInt("max-parallel-workers")
			if errMo != nil {
				return fmt.Errorf("error parsing flag 'max-parallel-workers' for advanced options update : %v", errMo)
			}

			maxParallelWorkersPerGather, errMp := cmd.Flags().GetInt("max-parallel-workers-per-gather")
			if errMp != nil {
				return fmt.Errorf(
					"error parsing flag 'max-parallel-workers-per-gather' for advanced options update : %v",
					errMp,
				)
			}

			maxPredLocksPerTransaction, errMr := cmd.Flags().GetInt("max-pred-locks-per-transaction")
			if errMr != nil {
				return fmt.Errorf(
					"error parsing flag 'max-pred-locks-per-transaction' for advanced options update : %v",
					errMr,
				)
			}

			maxPreparedTransactions, errMe := cmd.Flags().GetInt("max-prepared-transactions")
			if errMe != nil {
				return fmt.Errorf(
					"error parsing flag 'max-prepared-transactions' for advanced options update : %v",
					errMe,
				)
			}

			maxReplicationSlots, errMi := cmd.Flags().GetInt("max-replication-slots")
			if errMi != nil {
				return fmt.Errorf("error parsing flag 'max-replication-slots' for advanced options update : %v", errMi)
			}

			maxStackDepth, errMs := cmd.Flags().GetInt("max-stack-depth")
			if errMs != nil {
				return fmt.Errorf("error parsing flag 'max-stack-depth' for advanced options update : %v", errMs)
			}

			maxStandbyArchiveDelay, errMv := cmd.Flags().GetInt("max-standby-archive-delay")
			if errMv != nil {
				return fmt.Errorf(
					"error parsing flag 'max-standby-archive-delay' for advanced options update : %v",
					errMv,
				)
			}

			maxStandbyStreamingDelay, errMy := cmd.Flags().GetInt("max-standby-streaming-delay")
			if errMy != nil {
				return fmt.Errorf(
					"error parsing flag 'max-standby-streaming-delay' for advanced options update : %v",
					errMy,
				)
			}

			maxWalSenders, errMd := cmd.Flags().GetInt("max-wal-senders")
			if errMd != nil {
				return fmt.Errorf("error parsing flag 'max-wal-senders' for advanced options update : %v", errMd)
			}

			maxWorkerProcesses, errMs := cmd.Flags().GetInt("max-worker-processes")
			if errMs != nil {
				return fmt.Errorf("error parsing flag 'max-worker-processes' for advanced options update : %v", errMs)
			}

			netBufferLength, errNbl := cmd.Flags().GetInt("net-buffer-length")
			if errNbl != nil {
				return fmt.Errorf("error parsing flag 'net-buffer-length' for advanced options update : %v", errNbl)
			}

			netReadTimeout, errNrt := cmd.Flags().GetInt("net-read-timeout")
			if errNrt != nil {
				return fmt.Errorf("error parsing flag 'net-read-timeout' for advanced options update : %v", errNrt)
			}

			netWriteTimeout, errNwt := cmd.Flags().GetInt("net-write-timeout")
			if errNwt != nil {
				return fmt.Errorf("error parsing flag 'net-write-timeout' for advanced options update : %v", errNwt)
			}

			pgPartmanBGWInterval, errPg := cmd.Flags().GetInt("pg-partman-bgw-interval")
			if errPg != nil {
				return fmt.Errorf(
					"error parsing flag 'pg-partman-bgw-interval' for advanced options update : %v",
					errPg,
				)
			}

			pgPartmanBGWRole, errPp := cmd.Flags().GetString("pg-partman-bgw-role")
			if errPp != nil {
				return fmt.Errorf("error parsing flag 'pg-partman-bgw-role' for advanced options update : %v", errPp)
			}

			pgStatStatementsTrack, errPs := cmd.Flags().GetString("pg-stat-statements-track")
			if errPs != nil {
				return fmt.Errorf(
					"error parsing flag 'pg-stat-statements-track' for advanced options update : %v",
					errPs,
				)
			}

			sortBufferSize, errSbs := cmd.Flags().GetInt("sort-buffer-size")
			if errSbs != nil {
				return fmt.Errorf("error parsing flag 'sort-buffer-size' for advanced options update : %v", errSbs)
			}

			tempFileLimit, errTe := cmd.Flags().GetInt("temp-file-limit")
			if errTe != nil {
				return fmt.Errorf("error parsing flag 'temp-file-limit' for advanced options update : %v", errTe)
			}

			tmpTableSize, errTts := cmd.Flags().GetInt("tmp-table-size")
			if errTts != nil {
				return fmt.Errorf("error parsing flag 'tmp-table-size' for advanced options update : %v", errTts)
			}

			trackActivityQuerySize, errTr := cmd.Flags().GetInt("track-activity-query-size")
			if errTr != nil {
				return fmt.Errorf(
					"error parsing flag 'track-activity-query-size' for advanced options update : %v",
					errTr,
				)
			}

			trackCommitTimestamp, errTa := cmd.Flags().GetString("track-commit-timestamp")
			if errTa != nil {
				return fmt.Errorf(
					"error parsing flag 'track-commit-timestamp' for advanced options update : %v",
					errTa,
				)
			}

			trackFunctions, errTc := cmd.Flags().GetString("track-functions")
			if errTc != nil {
				return fmt.Errorf("error parsing flag 'track-functions' for advanced options update : %v", errTc)
			}

			trackIOTiming, errTi := cmd.Flags().GetString("track-io-timing")
			if errTi != nil {
				return fmt.Errorf("error parsing flag 'track-io-timing' for advanced options update : %v", errTi)
			}

			waitTimeout, errWt := cmd.Flags().GetInt("wait-timeout")
			if errWt != nil {
				return fmt.Errorf("error parsing flag 'wait-timeout' for advanced options update : %v", errWt)
			}

			walSenderTimeout, errWa := cmd.Flags().GetInt("wal-sender-timeout")
			if errWa != nil {
				return fmt.Errorf("error parsing flag 'wal-sender-timeout' for advanced options update : %v", errWa)
			}

			walWriterDelay, errWl := cmd.Flags().GetInt("wal-writer-delay")
			if errWl != nil {
				return fmt.Errorf("error parsing flag 'wal-writer-delay' for advanced options update : %v", errWl)
			}

			o.AdvancedOptionsReq = &govultr.DatabaseAdvancedOptions{}

			if cmd.Flags().Changed("autovacuum-analyze-scale-factor") {
				o.AdvancedOptionsReq.AutovacuumAnalyzeScaleFactor = autovacuumAnalyzeScaleFactor
			}

			if cmd.Flags().Changed("autovacuum-analyze-threshold") {
				o.AdvancedOptionsReq.AutovacuumAnalyzeThreshold = autovacuumAnalyzeThreshold
			}

			if cmd.Flags().Changed("autovacuum-freeze-max-age") {
				o.AdvancedOptionsReq.AutovacuumFreezeMaxAge = autovacuumFreezeMaxAge
			}

			if cmd.Flags().Changed("autovacuum-max-workers") {
				o.AdvancedOptionsReq.AutovacuumMaxWorkers = autovacuumMaxWorkers
			}

			if cmd.Flags().Changed("autovacuum-naptime") {
				o.AdvancedOptionsReq.AutovacuumNaptime = autovacuumNaptime
			}

			if cmd.Flags().Changed("autovacuum-vacuum-cost-delay") {
				o.AdvancedOptionsReq.AutovacuumVacuumCostDelay = autovacuumVacuumCostDelay
			}

			if cmd.Flags().Changed("autovacuum-vacuum-cost-limit") {
				o.AdvancedOptionsReq.AutovacuumVacuumCostLimit = autovacuumVacuumCostLimit
			}

			if cmd.Flags().Changed("autovacuum-vacuum-scale-factor") {
				o.AdvancedOptionsReq.AutovacuumVacuumScaleFactor = autovacuumVacuumScaleFactor
			}

			if cmd.Flags().Changed("autovacuum-vacuum-threshold") {
				o.AdvancedOptionsReq.AutovacuumVacuumThreshold = autovacuumVacuumThreshold
			}

			if cmd.Flags().Changed("bgwriter-delay") {
				o.AdvancedOptionsReq.BGWRITERDelay = bgwriterDelay
			}

			if cmd.Flags().Changed("bgwriter-flush-after") {
				o.AdvancedOptionsReq.BGWRITERFlushAFter = bgwriterFlushAfter
			}

			if cmd.Flags().Changed("bgwriter-lru-maxpages") {
				o.AdvancedOptionsReq.BGWRITERLRUMaxPages = bgwriterLruMaxpages
			}

			if cmd.Flags().Changed("bgwriter-lru-multiplier") {
				o.AdvancedOptionsReq.BGWRITERLRUMultiplier = bgwriterLruMultiplier
			}

			if cmd.Flags().Changed("connect-timeout") {
				o.AdvancedOptionsReq.ConnectTimeout = connectTimeout
			}

			if cmd.Flags().Changed("deadlock-timeout") {
				o.AdvancedOptionsReq.DeadlockTimeout = deadlockTimeout
			}

			if cmd.Flags().Changed("default-toast-compression") {
				o.AdvancedOptionsReq.DefaultToastCompression = defaultToastCompression
			}

			if cmd.Flags().Changed("group-concat-max-len") {
				o.AdvancedOptionsReq.GroupConcatMaxLen = groupConcatMaxLen
			}

			if cmd.Flags().Changed("idle-in-transaction-session-timeout") {
				o.AdvancedOptionsReq.IdleInTransactionSessionTimeout = idleInTransactionSessionTimeout
			}

			if cmd.Flags().Changed("innodb-change-buffer-max-size") {
				o.AdvancedOptionsReq.InnoDBChangeBufferMaxSize = innoDBChangeBufferMaxSize
			}

			if cmd.Flags().Changed("innodb-flush-neighbors") {
				o.AdvancedOptionsReq.InnoDBFlushNeighbors = innoDBFlushNeighbors
			}

			if cmd.Flags().Changed("innodb-ft-min-token-size") {
				o.AdvancedOptionsReq.InnoDBFTMinTokenSize = innoDBFTMinTokenSize
			}

			if cmd.Flags().Changed("innodb-ft-server-stopword-table") {
				o.AdvancedOptionsReq.InnoDBFTServerStopwordTable = innoDBFTServerStopwordTable
			}

			if cmd.Flags().Changed("innodb-lock-wait-timeout") {
				o.AdvancedOptionsReq.InnoDBLockWaitTimeout = innoDBLockWaitTimeout
			}

			iif cmd.Flags().Changed("innodb-log-buffer-size") {
				o.AdvancedOptionsReq.InnoDBLogBufferSize = innoDBLogBufferSize
			}

			if cmd.Flags().Changed("innodb-online-alter-log-max-size") {
				o.AdvancedOptionsReq.InnoDBOnlineAlterLogMaxSize = innoDBOnlineAlterLogMaxSize
			}

			if cmd.Flags().Changed("innodb-print-all-deadlocks") {
				o.AdvancedOptionsReq.InnoDBPrintAllDeadlocks = innoDBPrintAllDeadlocks
			}

			if cmd.Flags().Changed("innodb-read-io-threads") {
				o.AdvancedOptionsReq.InnoDBReadIOThreads = innoDBReadIOThreads
			}

			if cmd.Flags().Changed("innodb-rollback-on-timeout") {
				o.AdvancedOptionsReq.InnoDBRollbackOnTimeout = innoDBRollbackOnTimeout
			}

			if cmd.Flags().Changed("innodb-thread-concurrency") {
				o.AdvancedOptionsReq.InnoDBThreadConcurrency = innoDBThreadConcurrency
			}

			if cmd.Flags().Changed("innodb-write-io-threads") {
				o.AdvancedOptionsReq.InnoDBWriteIOThreads = innoDBWriteIOThreads
			}

			if cmd.Flags().Changed("interactive-timeout") {
				o.AdvancedOptionsReq.InteractiveTimeout = interactiveTimeout
			}

			if cmd.Flags().Changed("internal-tmp-mem-storage-engine") {
				o.AdvancedOptionsReq.InternalTmpMemStorageEngine = internalTmpMemStorageEngine
			}

			if cmd.Flags().Changed("jit") {
				o.AdvancedOptionsReq.Jit = nil
			}

			if cmd.Flags().Changed("log-autovacuum-min-duration") {
				o.AdvancedOptionsReq.LogAutovacuumMinDuration = logAutovacuumMinDuration
			}

			if cmd.Flags().Changed("log-error-verbosity") {
				o.AdvancedOptionsReq.LogErrorVerbosity = logErrorVerbosity
			}

			if cmd.Flags().Changed("log-line-prefix") {
				o.AdvancedOptionsReq.LogLinePrefix = logLinePrefix
			}

			if cmd.Flags().Changed("log-min-duration-statement") {
				o.AdvancedOptionsReq.LogMinDurationStatement = logMinDurationStatement
			}

			if cmd.Flags().Changed("max-allowed-packet") {
				o.AdvancedOptionsReq.MaxAllowedPacket = maxAllowedPacket
			}

			if cmd.Flags().Changed("max-files-per-process") {
				o.AdvancedOptionsReq.MaxFilesPerProcess = maxFilesPerProcess
			}

			if cmd.Flags().Changed("max-heap-table-size") {
				o.AdvancedOptionsReq.MaxHeapTableSize = maxHeapTableSize
			}

			if cmd.Flags().Changed("max-locks-per-transaction") {
				o.AdvancedOptionsReq.MaxLocksPerTransaction = maxLocksPerTransaction
			}

			if cmd.Flags().Changed("max-logical-replication-workers") {
				o.AdvancedOptionsReq.MaxLogicalReplicationWorkers = maxLogicalReplicationWorkers
			}

			if cmd.Flags().Changed("max-parallel-workers") {
				o.AdvancedOptionsReq.MaxParallelWorkers = maxParallelWorkers
			}

			if cmd.Flags().Changed("max-parallel-workers-per-gather") {
				o.AdvancedOptionsReq.MaxParallelWorkersPerGather = maxParallelWorkersPerGather
			}

			if cmd.Flags().Changed("max-pred-locks-per-transaction") {
				o.AdvancedOptionsReq.MaxPredLocksPerTransaction = maxPredLocksPerTransaction
			}

			if cmd.Flags().Changed("max-prepared-transactions") {
				o.AdvancedOptionsReq.MaxPreparedTransactions = maxPreparedTransactions
			}

			if cmd.Flags().Changed("max-replication-slots") {
				o.AdvancedOptionsReq.MaxReplicationSlots = maxReplicationSlots
			}

			if cmd.Flags().Changed("max-stack-depth") {
				o.AdvancedOptionsReq.MaxStackDepth = maxStackDepth
			}

			if cmd.Flags().Changed("max-standby-archive-delay") {
				o.AdvancedOptionsReq.MaxStandbyArchiveDelay = maxStandbyArchiveDelay
			}

			if cmd.Flags().Changed("max-standby-streaming-delay") {
				o.AdvancedOptionsReq.MaxStandbyStreamingDelay = maxStandbyStreamingDelay
			}

			if cmd.Flags().Changed("max-wal-senders") {
				o.AdvancedOptionsReq.MaxWalSenders = maxWalSenders
			}

			if cmd.Flags().Changed("max-worker-processes") {
				o.AdvancedOptionsReq.MaxWorkerProcesses = maxWorkerProcesses
			}

			if cmd.Flags().Changed("net-buffer-length") {
				o.AdvancedOptionsReq.NetBufferLength = netBufferLength
			}

			if cmd.Flags().Changed("net-read-timeout") {
				o.AdvancedOptionsReq.NetReadTimeout = netReadTimeout
			}

			if cmd.Flags().Changed("net-write-timeout") {
				o.AdvancedOptionsReq.NetWriteTimeout = netWriteTimeout
			}

			if cmd.Flags().Changed("pg-partman-bgw-interval") {
				o.AdvancedOptionsReq.PGPartmanBGWInterval = pgPartmanBGWInterval
			}

			if cmd.Flags().Changed("pg-partman-bgw-role") {
				o.AdvancedOptionsReq.PGPartmanBGWRole = pgPartmanBGWRole
			}

			if cmd.Flags().Changed("pg-stat-statements-track") {
				o.AdvancedOptionsReq.PGStateStatementsTrack = pgStatStatementsTrack
			}

			if cmd.Flags().Changed("sort-buffer-size") {
				o.AdvancedOptionsReq.SortBufferSize = sortBufferSize
			}

			if cmd.Flags().Changed("temp-file-limit") {
				o.AdvancedOptionsReq.TempFileLimit = tempFileLimit
			}

			if cmd.Flags().Changed("tmp-table-size") {
				o.AdvancedOptionsReq.TmpTableSize = tmpTableSize
			}

			if cmd.Flags().Changed("track-activity-query-size") {
				o.AdvancedOptionsReq.TrackActivityQuerySize = trackActivityQuerySize
			}

			if cmd.Flags().Changed("track-commit-timestamp") {
				o.AdvancedOptionsReq.TrackCommitTimestamp = trackCommitTimestamp
			}

			if cmd.Flags().Changed("track-functions") {
				o.AdvancedOptionsReq.TrackFunctions = trackFunctions
			}

			if cmd.Flags().Changed("track-io-timing") {
				o.AdvancedOptionsReq.TrackIOTiming = trackIOTiming
			}

			if cmd.Flags().Changed("wait-timeout") {
				o.AdvancedOptionsReq.WaitTimeout = waitTimeout
			}

			if cmd.Flags().Changed("wal-sender-timeout") {
				o.AdvancedOptionsReq.WALSenderTImeout = walSenderTimeout
			}

			if cmd.Flags().Changed("wal-writer-delay") {
				o.AdvancedOptionsReq.WALWriterDelay = walWriterDelay
			}

			if cmd.Flags().Changed("jit") {
				o.AdvancedOptionsReq.Jit = &jit
			}

			cur, avail, err := o.updateAdvancedOptions()
			if err != nil {
				return fmt.Errorf("error updating database advanced options : %v", err)
			}

			data := &AdvancedOptionsPrinter{Configured: cur, Available: avail}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	advancedOptionUpdate.Flags().Float32(
		"autovacuum-analyze-scale-factor",
		0,
		"set the managed postgresql configuration value for autovacuum_analyze_scale_factor",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-analyze-threshold",
		0,
		"set the managed postgresql configuration value for autovacuum_analyze_threshold",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-freeze-max-age",
		0,
		"set the managed postgresql configuration value for autovacuum_freeze_max_age",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-max-workers",
		0,
		"set the managed postgresql configuration value for autovacuum_max_workers",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-naptime",
		0,
		"set the managed postgresql configuration value for autovacuum_naptime",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-vacuum-cost-delay",
		0,
		"set the managed postgresql configuration value for autovacuum_vacuum_cost_delay",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-vacuum-cost-limit",
		0,
		"set the managed postgresql configuration value for autovacuum_vacuum_cost_limit",
	)
	advancedOptionUpdate.Flags().Float32(
		"autovacuum-vacuum-scale-factor",
		0,
		"set the managed postgresql configuration value for autovacuum_vacuum_scale_factor",
	)
	advancedOptionUpdate.Flags().Int(
		"autovacuum-vacuum-threshold",
		0,
		"set the managed postgresql configuration value for autovacuum_vacuum_threshold",
	)
	advancedOptionUpdate.Flags().Int(
		"bgwriter-delay",
		0,
		"set the managed postgresql configuration value for bgwriter_delay",
	)
	advancedOptionUpdate.Flags().Int(
		"bgwriter-flush-after",
		0,
		"set the managed postgresql configuration value for bgwriter_flush_after",
	)
	advancedOptionUpdate.Flags().Int(
		"bgwriter-lru-maxpages",
		0,
		"set the managed postgresql configuration value for bgwriter_lru_maxpages",
	)
	advancedOptionUpdate.Flags().Float32(
		"bgwriter-lru-multiplier",
		0,
		"set the managed postgresql configuration value for bgwriter_lru_multiplier",
	)
	advancedOptionUpdate.Flags().Int(
		"connect-timeout",
		0,
		"set the managed mysql configuration value for connect_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"deadlock-timeout",
		0,
		"set the managed postgresql configuration value for deadlock_timeout",
	)
	advancedOptionUpdate.Flags().String(
		"default-toast-compression",
		"",
		"set the managed postgresql configuration value for default_toast_compression",
	)
	advancedOptionUpdate.Flags().Int(
		"group-concat-max-len",
		0,
		"set the managed mysql configuration value for group_concat_max_len",
	)
	advancedOptionUpdate.Flags().Int(
		"idle-in-transaction-session-timeout",
		0,
		"set the managed postgresql configuration value for idle_in_transaction_session_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-change-buffer-max-size",
		0,
		"set the managed mysql configuration value for innodb_change_buffer_max_size",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-flush-neighbors",
		0,
		"set the managed mysql configuration value for innodb_flush_neighbors",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-ft-min-token-size",
		0,
		"set the managed mysql configuration value for innodb_ft_min_token_size",
	)
	advancedOptionUpdate.Flags().String(
		"innodb-ft-server-stopword-table",
		"",
		"set the managed mysql configuration value for innodb_ft_server_stopword_table",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-lock-wait-timeout",
		0,
		"set the managed mysql configuration value for innodb_lock_wait_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-log-buffer-size",
		0,
		"set the managed mysql configuration value for innodb_log_buffer_size",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-online-alter-log-max-size",
		0,
		"set the managed mysql configuration value for innodb_online_alter_log_max_size",
	)
	advancedOptionUpdate.Flags().Bool(
		"innodb-print-all-deadlocks",
		false,
		"set the managed mysql configuration value for innodb_print_all_deadlocks",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-read-io-threads",
		0,
		"set the managed mysql configuration value for innodb_read_io_threads",
	)
	advancedOptionUpdate.Flags().Bool(
		"innodb-rollback-on-timeout",
		false,
		"set the managed mysql configuration value for innodb_rollback_on_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-thread-concurrency",
		0,
		"set the managed mysql configuration value for innodb_thread_concurrency",
	)
	advancedOptionUpdate.Flags().Int(
		"innodb-write-io-threads",
		0,
		"set the managed mysql configuration value for innodb_write_io_threads",
	)
	advancedOptionUpdate.Flags().Int(
		"interactive-timeout",
		0,
		"set the managed mysql configuration value for interactive_timeout",
	)
	advancedOptionUpdate.Flags().String(
		"internal-tmp-mem-storage-engine",
		"",
		"set the managed mysql configuration value for internal_tmp_mem_storage_engine",
	)
	advancedOptionUpdate.Flags().Bool(
		"jit",
		false,
		"set the managed postgresql configuration value for jit",
	)
	advancedOptionUpdate.Flags().Int(
		"log-autovacuum-min-duration",
		0,
		"set the managed postgresql configuration value for log_autovacuum_min_duration",
	)
	advancedOptionUpdate.Flags().String(
		"log-error-verbosity",
		"",
		"set the managed postgresql configuration value for log_error_verbosity",
	)
	advancedOptionUpdate.Flags().String(
		"log-line-prefix",
		"",
		"set the managed postgresql configuration value for log_line_prefix",
	)
	advancedOptionUpdate.Flags().Int(
		"log-min-duration-statement",
		0,
		"set the managed postgresql configuration value for log_min_duration_statement",
	)
	advancedOptionUpdate.Flags().Int(
		"max-allowed-packet",
		0,
		"set the managed mysql configuration value for max_allowed_packet",
	)
	advancedOptionUpdate.Flags().Int(
		"max-files-per-process",
		0,
		"set the managed postgresql configuration value for max_files_per_process",
	)
	advancedOptionUpdate.Flags().Int(
		"max-heap-table-size",
		0,
		"set the managed mysql configuration value for max_heap_table_size",
	)
	advancedOptionUpdate.Flags().Int(
		"max-locks-per-transaction",
		0,
		"set the managed postgresql configuration value for max_locks_per_transaction",
	)
	advancedOptionUpdate.Flags().Int(
		"max-logical-replication-workers",
		0,
		"set the managed postgresql configuration value for max_logical_replication_workers",
	)
	advancedOptionUpdate.Flags().Int(
		"max-parallel-workers",
		0,
		"set the managed postgresql configuration value for max_parallel_workers",
	)
	advancedOptionUpdate.Flags().Int(
		"max-parallel-workers-per-gather",
		0,
		"set the managed postgresql configuration value for max_parallel_workers_per_gather",
	)
	advancedOptionUpdate.Flags().Int(
		"max-pred-locks-per-transaction",
		0,
		"set the managed postgresql configuration value for max_pred_locks_per_transaction",
	)
	advancedOptionUpdate.Flags().Int(
		"max-prepared-transactions",
		0,
		"set the managed postgresql configuration value for max_prepared_transactions",
	)
	advancedOptionUpdate.Flags().Int(
		"max-replication-slots",
		0,
		"set the managed postgresql configuration value for max_replication_slots",
	)
	advancedOptionUpdate.Flags().Int(
		"max-stack-depth",
		0,
		"set the managed postgresql configuration value for max_stack_depth",
	)
	advancedOptionUpdate.Flags().Int(
		"max-standby-archive-delay",
		0,
		"set the managed postgresql configuration value for max_standby_archive_delay",
	)
	advancedOptionUpdate.Flags().Int(
		"max-standby-streaming-delay",
		0,
		"set the managed postgresql configuration value for max_standby_streaming_delay",
	)
	advancedOptionUpdate.Flags().Int(
		"max-wal-senders",
		0,
		"set the managed postgresql configuration value for max_wal_senders",
	)
	advancedOptionUpdate.Flags().Int(
		"max-worker-processes",
		0,
		"set the managed postgresql configuration value for max_worker_processes",
	)
	advancedOptionUpdate.Flags().Int(
		"net-buffer-length",
		0,
		"set the managed mysql configuration value for net_buffer_length",
	)
	advancedOptionUpdate.Flags().Int(
		"net-read-timeout",
		0,
		"set the managed mysql configuration value for net_read_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"net-write-timeout",
		0,
		"set the managed mysql configuration value for net_write_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"pg-partman-bgw-interval",
		0,
		"set the managed postgresql configuration value for pg_partman_bgw.interval",
	)
	advancedOptionUpdate.Flags().String(
		"pg-partman-bgw-role",
		"",
		"set the managed postgresql configuration value for pg_partman_bgw.role",
	)
	advancedOptionUpdate.Flags().String(
		"pg-stat-statements-track",
		"",
		"set the managed postgresql configuration value for pg_stat_statements.track",
	)
	advancedOptionUpdate.Flags().Int(
		"sort-buffer-size",
		0,
		"set the managed mysql configuration value for sort_buffer_size",
	)
	advancedOptionUpdate.Flags().Int(
		"temp-file-limit",
		0,
		"set the managed postgresql configuration value for temp_file_limit",
	)
	advancedOptionUpdate.Flags().Int(
		"tmp-table-size",
		0,
		"set the managed mysql configuration value for tmp_table_size",
	)
	advancedOptionUpdate.Flags().Int(
		"track-activity-query-size",
		0,
		"set the managed postgresql configuration value for track_activity_query_size",
	)
	advancedOptionUpdate.Flags().String(
		"track-commit-timestamp",
		"",
		"set the managed postgresql configuration value for track_commit_timestamp",
	)
	advancedOptionUpdate.Flags().String(
		"track-functions",
		"",
		"set the managed postgresql configuration value for track_functions",
	)
	advancedOptionUpdate.Flags().String(
		"track-io-timing",
		"",
		"set the managed postgresql configuration value for track_io_timing",
	)
	advancedOptionUpdate.Flags().Int(
		"wait-timeout",
		0,
		"set the managed mysql configuration value for wait_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"wal-sender-timeout",
		0,
		"set the managed postgresql configuration value for wal_sender_timeout",
	)
	advancedOptionUpdate.Flags().Int(
		"wal-writer-delay",
		0,
		"set the managed postgresql configuration value for wal_writer_delay",
	)

	advancedOption.AddCommand(
		advancedOptionList,
		advancedOptionUpdate,
	)

	// Version
	version := &cobra.Command{
		Use:   "version",
		Short: "Commands to handle database version upgrades",
	}

	// Version List
	versionList := &cobra.Command{
		Use:   "list <Database ID>",
		Short: "List all version upgrades for a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vs, err := o.listVersions()
			if err != nil {
				return fmt.Errorf("error retrieving database versions : %v", err)
			}

			data := &VersionsPrinter{Versions: vs}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Version Upgrade
	versionUpgrade := &cobra.Command{
		Use:   "upgrade <Database ID>",
		Short: "Start a version upgrade on a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a database ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			version, errVe := cmd.Flags().GetString("version")
			if errVe != nil {
				return fmt.Errorf("error parsing flag 'version' for database upgrade : %v", errVe)
			}

			o.UpgradeReq = &govultr.DatabaseVersionUpgradeReq{
				Version: version,
			}

			msg, err := o.upgradeVersion()
			if err != nil {
				return fmt.Errorf("error starting database version upgrade : %v", err)
			}

			o.Base.Printer.Display(printer.Info(msg), nil)

			return nil
		},
	}

	versionUpgrade.Flags().StringP("version", "v", "", "version of the manaaged database to upgrade to")
	if err := versionUpgrade.MarkFlagRequired("version"); err != nil {
		fmt.Printf("error marking version upgrade 'version' flag required: %v", err)
		os.Exit(1)
	}

	version.AddCommand(
		versionList,
		versionUpgrade,
	)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		user,
		db,
		usage,
		maintenance,
		plan,
		alert,
		migration,
		readReplica,
		backup,
		connectionPool,
		advancedOption,
		version,
	)

	return cmd
}

type options struct {
	Base                    *cli.Base
	CreateReq               *govultr.DatabaseCreateReq
	UpdateReq               *govultr.DatabaseUpdateReq
	UserCreateReq           *govultr.DatabaseUserCreateReq
	UserUpdateReq           *govultr.DatabaseUserUpdateReq
	UserUpdateACLReq        *govultr.DatabaseUserACLReq
	DBCreateReq             *govultr.DatabaseDBCreateReq
	AlertsReq               *govultr.DatabaseListAlertsReq
	MigrationReq            *govultr.DatabaseMigrationStartReq
	ReadReplicaCreateReq    *govultr.DatabaseAddReplicaReq
	BackupReq               *govultr.DatabaseBackupRestoreReq
	ForkReq                 *govultr.DatabaseForkReq
	ConnectionPoolCreateReq *govultr.DatabaseConnectionPoolCreateReq
	ConnectionPoolUpdateReq *govultr.DatabaseConnectionPoolUpdateReq
	AdvancedOptionsReq      *govultr.DatabaseAdvancedOptions
	UpgradeReq              *govultr.DatabaseVersionUpgradeReq
}

func (o *options) list() ([]govultr.Database, *govultr.Meta, error) {
	dbs, meta, _, err := o.Base.Client.Database.List(o.Base.Context, nil)
	return dbs, meta, err
}

func (o *options) get() (*govultr.Database, error) {
	db, _, err := o.Base.Client.Database.Get(o.Base.Context, o.Base.Args[0])
	return db, err
}

func (o *options) create() (*govultr.Database, error) {
	db, _, err := o.Base.Client.Database.Create(o.Base.Context, o.CreateReq)
	return db, err
}

func (o *options) update() (*govultr.Database, error) {
	db, _, err := o.Base.Client.Database.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
	return db, err
}

func (o *options) del() error {
	return o.Base.Client.Database.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) listPlans() ([]govultr.DatabasePlan, *govultr.Meta, error) {
	plans, meta, _, err := o.Base.Client.Database.ListPlans(o.Base.Context, nil)
	return plans, meta, err
}

func (o *options) listUsers() ([]govultr.DatabaseUser, *govultr.Meta, error) {
	users, meta, _, err := o.Base.Client.Database.ListUsers(o.Base.Context, o.Base.Args[0])
	return users, meta, err
}

func (o *options) getUser() (*govultr.DatabaseUser, error) {
	user, _, err := o.Base.Client.Database.GetUser(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return user, err
}

func (o *options) createUser() (*govultr.DatabaseUser, error) {
	user, _, err := o.Base.Client.Database.CreateUser(o.Base.Context, o.Base.Args[0], o.UserCreateReq)
	return user, err
}

func (o *options) updateUser() (*govultr.DatabaseUser, error) {
	user, _, err := o.Base.Client.Database.UpdateUser(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.UserUpdateReq)
	return user, err
}

func (o *options) delUser() error {
	return o.Base.Client.Database.DeleteUser(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) updateUserACL() (*govultr.DatabaseUser, error) {
	user, _, err := o.Base.Client.Database.UpdateUserACL(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.UserUpdateACLReq) //nolint:lll
	return user, err
}

func (o *options) listDBs() ([]govultr.DatabaseDB, *govultr.Meta, error) {
	dbs, meta, _, err := o.Base.Client.Database.ListDBs(o.Base.Context, o.Base.Args[0])
	return dbs, meta, err
}

func (o *options) createDB() (*govultr.DatabaseDB, error) {
	db, _, err := o.Base.Client.Database.CreateDB(o.Base.Context, o.Base.Args[0], o.DBCreateReq)
	return db, err
}

func (o *options) delDB() error {
	return o.Base.Client.Database.DeleteDB(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) getUsage() (*govultr.DatabaseUsage, error) {
	usage, _, err := o.Base.Client.Database.GetUsage(o.Base.Context, o.Base.Args[0])
	return usage, err
}

func (o *options) listMaintUpdates() ([]string, error) {
	updates, _, err := o.Base.Client.Database.ListMaintenanceUpdates(o.Base.Context, o.Base.Args[0])
	return updates, err
}

func (o *options) startMaintUpdate() (string, error) {
	updates, _, err := o.Base.Client.Database.StartMaintenance(o.Base.Context, o.Base.Args[0])
	return updates, err
}

func (o *options) listAlerts() ([]govultr.DatabaseAlert, error) {
	alerts, _, err := o.Base.Client.Database.ListServiceAlerts(o.Base.Context, o.Base.Args[0], o.AlertsReq)
	return alerts, err
}

func (o *options) getMigrationStatus() (*govultr.DatabaseMigration, error) {
	status, _, err := o.Base.Client.Database.GetMigrationStatus(o.Base.Context, o.Base.Args[0])
	return status, err
}

func (o *options) startMigration() (*govultr.DatabaseMigration, error) {
	status, _, err := o.Base.Client.Database.StartMigration(o.Base.Context, o.Base.Args[0], o.MigrationReq)
	return status, err
}

func (o *options) detachMigration() error {
	return o.Base.Client.Database.DetachMigration(o.Base.Context, o.Base.Args[0])
}

func (o *options) createReadReplica() (*govultr.Database, error) {
	db, _, err := o.Base.Client.Database.AddReadOnlyReplica(o.Base.Context, o.Base.Args[0], o.ReadReplicaCreateReq)
	return db, err
}

func (o *options) promoteReadReplica() error {
	return o.Base.Client.Database.PromoteReadReplica(o.Base.Context, o.Base.Args[0])
}

func (o *options) getBackup() (*govultr.DatabaseBackups, error) {
	backup, _, err := o.Base.Client.Database.GetBackupInformation(o.Base.Context, o.Base.Args[0])
	return backup, err
}

func (o *options) restoreBackup() (*govultr.Database, error) {
	db, _, err := o.Base.Client.Database.RestoreFromBackup(o.Base.Context, o.Base.Args[0], o.BackupReq)
	return db, err
}

func (o *options) fork() (*govultr.Database, error) {
	db, _, err := o.Base.Client.Database.Fork(o.Base.Context, o.Base.Args[0], o.ForkReq)
	return db, err
}

func (o *options) listConnectionPools() (*govultr.DatabaseConnections, []govultr.DatabaseConnectionPool, *govultr.Meta, error) { //nolint:lll
	cons, pool, meta, _, err := o.Base.Client.Database.ListConnectionPools(o.Base.Context, o.Base.Args[0])
	return cons, pool, meta, err
}

func (o *options) getConnectionPool() (*govultr.DatabaseConnectionPool, error) {
	pool, _, err := o.Base.Client.Database.GetConnectionPool(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return pool, err
}

func (o *options) createConnectionPool() (*govultr.DatabaseConnectionPool, error) {
	pool, _, err := o.Base.Client.Database.CreateConnectionPool(o.Base.Context, o.Base.Args[0], o.ConnectionPoolCreateReq) //nolint:lll
	return pool, err
}

func (o *options) updateConnectionPool() (*govultr.DatabaseConnectionPool, error) {
	pool, _, err := o.Base.Client.Database.UpdateConnectionPool(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.ConnectionPoolUpdateReq) //nolint:lll
	return pool, err
}

func (o *options) delConnectionPool() error {
	return o.Base.Client.Database.DeleteConnectionPool(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) listAdvancedOptions() (*govultr.DatabaseAdvancedOptions, []govultr.AvailableOption, error) {
	cur, avail, _, err := o.Base.Client.Database.ListAdvancedOptions(o.Base.Context, o.Base.Args[0])
	return cur, avail, err
}

func (o *options) updateAdvancedOptions() (*govultr.DatabaseAdvancedOptions, []govultr.AvailableOption, error) {
	cur, avail, _, err := o.Base.Client.Database.UpdateAdvancedOptions(o.Base.Context, o.Base.Args[0], o.AdvancedOptionsReq) //nolint:lll
	return cur, avail, err
}

func (o *options) listVersions() ([]string, error) {
	vers, _, err := o.Base.Client.Database.ListAvailableVersions(o.Base.Context, o.Base.Args[0])
	return vers, err
}

func (o *options) upgradeVersion() (string, error) {
	up, _, err := o.Base.Client.Database.StartVersionUpgrade(o.Base.Context, o.Base.Args[0], o.UpgradeReq)
	return up, err
}
