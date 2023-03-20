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

	instanceCreate.Flags().StringP("region", "r", "", "region id you wish to have the instance created in")
	instanceCreate.Flags().StringP("plan", "p", "", "plan id you wish the instance to have")
	instanceCreate.Flags().IntP("os", "o", 0, "os id you wish the instance to have")
	if err := instanceCreate.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking instance create 'region' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := instanceCreate.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking instance create 'plan' flag required: %v\n", err)
		os.Exit(1)
	}

	// Optional Params
	instanceCreate.Flags().StringP("ipxe", "", "", "if you've selected the 'custom' operating system, this can be set to chainload the specified URL on bootup")
	instanceCreate.Flags().StringP("iso", "", "", "iso ID you want to create the instance with")
	instanceCreate.Flags().StringP("snapshot", "", "", "snapshot ID you want to create the instance with")
	instanceCreate.Flags().StringP("script-id", "", "", "script id of the startup script")
	instanceCreate.Flags().BoolP("ipv6", "", false, "enable ipv6 | true or false")
	instanceCreate.Flags().BoolP("private-network", "", false, "Deprecated: use vpc-enable instead. enable private network | true or false")
	instanceCreate.Flags().StringSliceP("network", "", []string{}, "Deprecated: use vpc-ids instead. network IDs you want to assign to the instance")
	instanceCreate.Flags().BoolP("vpc-enable", "", false, "enable VPC | true or false")
	instanceCreate.Flags().StringSliceP("vpc-ids", "", []string{}, "VPC IDs you want to assign to the instance")
	instanceCreate.Flags().StringP("label", "l", "", "label you want to give this instance")
	instanceCreate.Flags().StringSliceP("ssh-keys", "s", []string{}, "ssh keys you want to assign to the instance")
	instanceCreate.Flags().BoolP("auto-backup", "b", false, "enable auto backups | true or false")
	instanceCreate.Flags().IntP("app", "a", 0, "application ID you want this instance to have")
	instanceCreate.Flags().StringP("image", "", "", "(optional) image ID of the application that will be installed on the server.")
	instanceCreate.Flags().StringP("userdata", "u", "", "plain text userdata you want to give this instance which the CLI will base64 encode")
	instanceCreate.Flags().BoolP("notify", "n", false, "notify when instance has been created | true or false")
	instanceCreate.Flags().BoolP("ddos", "d", false, "enable ddos protection | true or false")
	instanceCreate.Flags().StringP("reserved-ipv4", "", "", "ip address of the floating IP to use as the main IP for this instance")
	instanceCreate.Flags().StringP("host", "", "", "The hostname to assign to this instance")
	instanceCreate.Flags().StringP("tag", "t", "", "Deprecated: use tags instead. The tag to assign to this instance")
	instanceCreate.Flags().StringSliceP("tags", "", []string{}, "A comma-separated list of tags to assign to this instance")
	instanceCreate.Flags().StringP("firewall-group", "", "", "The firewall group to assign to this instance")

	instanceList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	instanceList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	instanceIPV4List.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	instanceIPV4List.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	instanceIPV6List.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	instanceIPV6List.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Sub commands for OS
	osCmd := &cobra.Command{
		Use:   "os",
		Short: "update operating system for an instance",
		Long:  ``,
	}

	osCmd.AddCommand(osUpdate, osUpdateList)
	osUpdate.Flags().IntP("os", "o", 0, "operating system ID you wish to use")
	if err := osUpdate.MarkFlagRequired("os"); err != nil {
		fmt.Printf("error marking instance os update 'os' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(osCmd)

	// Sub commands for App
	appCMD := &cobra.Command{
		Use:   "app",
		Short: "update application for an instance",
		Long:  ``,
	}
	appCMD.AddCommand(appUpdate, appUpdateList)
	appUpdate.Flags().IntP("app", "a", 0, "application ID you wish to use")
	if err := appUpdate.MarkFlagRequired("app"); err != nil {
		fmt.Printf("error marking instance app update 'app' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(appCMD)

	// Sub commands for Image
	imageCMD := &cobra.Command{
		Use:   "image",
		Short: "update image for an instance",
		Long:  ``,
	}
	imageCMD.AddCommand(imageUpdate, appUpdateList)
	imageUpdate.Flags().StringP("image", "", "", "application image ID you wish to use")
	if err := imageUpdate.MarkFlagRequired("image"); err != nil {
		fmt.Printf("error marking instance image update 'image' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(imageCMD)

	// Sub commands for Backup
	backupCMD := &cobra.Command{
		Use:   "backup",
		Short: "list and create backup schedules for an instance",
		Long:  ``,
	}
	backupCMD.AddCommand(backupGet, backupCreate)
	backupCreate.Flags().StringP("type", "t", "", "type string Backup cron type. Can be one of 'daily', 'weekly', 'monthly', 'daily_alt_even', or 'daily_alt_odd'.")
	if err := backupCreate.MarkFlagRequired("type"); err != nil {
		fmt.Printf("error marking instance backup create 'type' flag required: %v\n", err)
		os.Exit(1)
	}
	backupCreate.Flags().IntP("hour", "o", 0, "Hour value (0-23). Applicable to crons: 'daily', 'weekly', 'monthly', 'daily_alt_even', 'daily_alt_odd'")
	backupCreate.Flags().IntP("dow", "w", 0, "Day-of-week value (0-6). Applicable to crons: 'weekly'")
	backupCreate.Flags().IntP("dom", "m", 0, "Day-of-month value (1-28). Applicable to crons: 'monthly'")
	databaseCmd.AddCommand(backupCMD)

	// IPV4 Subcommands
	isoCmd := &cobra.Command{
		Use:   "iso",
		Short: "attach/detach ISOs to a given instance",
		Long:  ``,
	}
	isoCmd.AddCommand(isoStatus, isoAttach, isoDetach)
	isoAttach.Flags().StringP("iso-id", "i", "", "id of the ISO you wish to attach")
	if err := isoAttach.MarkFlagRequired("iso-id"); err != nil {
		fmt.Printf("error marking instance iso attach 'iso-id' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(isoCmd)

	ipv4Cmd := &cobra.Command{
		Use:   "ipv4",
		Short: "list/create/delete ipv4 on instance",
		Long:  ``,
	}
	ipv4Cmd.AddCommand(instanceIPV4List, createIpv4, deleteIpv4)
	createIpv4.Flags().Bool("reboot", false, "whether to reboot instance after adding ipv4 address")
	deleteIpv4.Flags().StringP("ipv4", "i", "", "ipv4 address you wish to delete")
	if err := deleteIpv4.MarkFlagRequired("ipv4"); err != nil {
		fmt.Printf("error marking instance delete IPv4 'ipv4' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(ipv4Cmd)

	// IPV6 Subcommands
	ipv6Cmd := &cobra.Command{
		Use:   "ipv6",
		Short: "commands for ipv6 on instance",
		Long:  ``,
	}
	ipv6Cmd.AddCommand(instanceIPV6List)
	databaseCmd.AddCommand(ipv6Cmd)

	// Plans SubCommands
	plansCmd := &cobra.Command{
		Use:   "plan",
		Short: "update/list plans for an instance",
		Long:  ``,
	}
	plansCmd.AddCommand(upgradePlan, upgradePlanList)
	upgradePlan.Flags().StringP("plan", "p", "", "plan id that you wish to upgrade to")
	if err := upgradePlan.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking instance plan upgrace 'plan' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(plansCmd)

	// ReverseDNS SubCommands
	reverseCmd := &cobra.Command{
		Use:   "reverse-dns",
		Short: "commands to handle reverse-dns on an instance",
		Long:  ``,
	}
	reverseCmd.AddCommand(defaultIpv4, listIpv6, deleteIpv6, setIpv4, setIpv6)
	defaultIpv4.Flags().StringP("ip", "i", "", "iPv4 address used in the reverse DNS update")
	if err := defaultIpv4.MarkFlagRequired("ip"); err != nil {
		fmt.Printf("error marking instance reverse-dns ipv4 'ip' flag required: %v\n", err)
		os.Exit(1)
	}
	deleteIpv6.Flags().StringP("ip", "i", "", "ipv6 address you wish to delete")

	if err := defaultIpv4.MarkFlagRequired("ip"); err != nil {
		fmt.Printf("error marking instance reverse-dns default-ipv4 'ip' flag required: %v\n", err)
		os.Exit(1)
	}
	setIpv4.Flags().StringP("ip", "i", "", "ip address you wish to set a reverse DNS on")
	setIpv4.Flags().StringP("entry", "e", "", "reverse dns entry")
	if err := setIpv4.MarkFlagRequired("ip"); err != nil {
		fmt.Printf("error marking instance reverse-dns set-ipv4 'ip' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := setIpv4.MarkFlagRequired("entry"); err != nil {
		fmt.Printf("error marking instance reverse-dns set-ipv4 'entry' flag required: %v\n", err)
		os.Exit(1)
	}

	setIpv6.Flags().StringP("ip", "i", "", "ip address you wish to set a reverse DNS on")
	setIpv6.Flags().StringP("entry", "e", "", "reverse dns entry")
	if err := setIpv6.MarkFlagRequired("ip"); err != nil {
		fmt.Printf("error marking instance reverse-dns set-ipv6 'ip' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := setIpv6.MarkFlagRequired("entry"); err != nil {
		fmt.Printf("error marking instance reverse-dns set-ipv6 'entry' flag required: %v\n", err)
		os.Exit(1)
	}
	databaseCmd.AddCommand(reverseCmd)

	userdataCmd := &cobra.Command{
		Use:   "user-data",
		Short: "commands to handle userdata on an instance",
		Long:  ``,
	}
	userdataCmd.AddCommand(setUserData, getUserData)
	setUserData.Flags().StringP("userdata", "d", "/dev/stdin", "file to read userdata from")
	databaseCmd.AddCommand(userdataCmd)

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
