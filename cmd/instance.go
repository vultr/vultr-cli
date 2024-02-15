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
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// Instance represents the instance command
func Instance() *cobra.Command {
	instanceCmd := &cobra.Command{
		Use:   "instance",
		Short: "commands to interact with instances on vultr",
		Long:  ``,
	}

	instanceCmd.AddCommand(instanceStart, instanceStop, instanceRestart, instanceReinstall, instanceTag, instanceDelete, instanceLabel, instanceBandwidth, instanceList, instanceInfo, updateFwgGroup, instanceRestore, instanceCreate)

	instanceTag.Flags().StringP("tag", "t", "", "tag you want to set for a given instance")
	instanceTag.MarkFlagRequired("tag")

	instanceLabel.Flags().StringP("label", "l", "", "label you want to set for a given instance")
	instanceLabel.MarkFlagRequired("label")

	updateFwgGroup.Flags().StringP("instance-id", "i", "", "instance id of the instance you want to use")
	updateFwgGroup.Flags().StringP("firewall-group-id", "f", "", "firewall group id that you want to assign. 0 Value will unset the firewall-group")
	updateFwgGroup.MarkFlagRequired("instance-id")
	updateFwgGroup.MarkFlagRequired("firewall-group-id")

	instanceRestore.Flags().StringP("backup", "b", "", "id of backup you wish to restore the instance with")
	instanceRestore.Flags().StringP("snapshot", "s", "", "id of snapshot you wish to restore the instance with")

	instanceCreate.Flags().StringP("region", "r", "", "region id you wish to have the instance created in")
	instanceCreate.Flags().StringP("plan", "p", "", "plan id you wish the instance to have")
	instanceCreate.Flags().IntP("operatingSystems", "o", 0, "operatingSystems id you wish the instance to have")
	instanceCreate.MarkFlagRequired("region")
	instanceCreate.MarkFlagRequired("plan")

	// Optional Params
	instanceCreate.Flags().StringP("ipxe", "", "", "if you've selected the 'custom' operating system, this can be set to chainload the specified URL on bootup")
	instanceCreate.Flags().StringP("iso", "", "", "iso ID you want to create the instance with")
	instanceCreate.Flags().StringP("snapshot", "", "", "snapshot ID you want to create the instance with")
	instanceCreate.Flags().StringP("script-id", "", "", "script id of the startup script")
	instanceCreate.Flags().BoolP("ipv6", "", false, "enable ipv6 | true or false")
	instanceCreate.Flags().BoolP("private-network", "", false, "enable private network | true or false")
	instanceCreate.Flags().StringArrayP("network", "", []string{}, "network IDs you want to assign to the instance")
	instanceCreate.Flags().StringP("label", "l", "", "label you want to give this instance")
	instanceCreate.Flags().StringArrayP("ssh-keys", "s", []string{}, "ssh keys you want to assign to the instance")
	instanceCreate.Flags().BoolP("auto-backup", "b", false, "enable auto backups | true or false")
	instanceCreate.Flags().IntP("app", "a", 0, "application ID you want this instance to have")
	instanceCreate.Flags().StringP("userdata", "u", "", "plain text userdata you want to give this instance which the CLI will base64 encode")
	instanceCreate.Flags().BoolP("notify", "n", true, "notify when instance has been created | true or false")
	instanceCreate.Flags().BoolP("ddos", "d", false, "enable ddos protection | true or false")
	instanceCreate.Flags().StringP("reserved-ipv4", "", "", "ip address of the floating IP to use as the main IP for this instance")
	instanceCreate.Flags().StringP("host", "", "", "The hostname to assign to this instance")
	instanceCreate.Flags().StringP("tag", "t", "", "The tag to assign to this instance")
	instanceCreate.Flags().StringP("firewall-group", "", "", "The firewall group to assign to this instance")

	instanceList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	instanceList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	instanceIPV4List.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	instanceIPV4List.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	instanceIPV6List.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	instanceIPV6List.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Sub commands for OS
	osCmd := &cobra.Command{
		Use:   "operatingSystems",
		Short: "update operating system for an instance",
		Long:  ``,
	}

	osCmd.AddCommand(osUpdate, osUpdateList)
	osUpdate.Flags().IntP("operatingSystems", "o", 0, "operating system ID you wish to use")
	osUpdate.MarkFlagRequired("operatingSystems")
	instanceCmd.AddCommand(osCmd)

	// Sub commands for App
	appCMD := &cobra.Command{
		Use:   "app",
		Short: "update application for an instance",
		Long:  ``,
	}
	appCMD.AddCommand(appUpdate, appUpdateList)
	appUpdate.Flags().IntP("app", "a", 0, "application ID you wish to use")
	appUpdate.MarkFlagRequired("app")
	instanceCmd.AddCommand(appCMD)

	// Sub commands for Backup
	backupCMD := &cobra.Command{
		Use:   "backup",
		Short: "list and create backup schedules for an instance",
		Long:  ``,
	}
	backupCMD.AddCommand(backupGet, backupCreate)
	backupCreate.Flags().StringP("type", "t", "", "type string Backup cron type. Can be one of 'daily', 'weekly', 'monthly', 'daily_alt_even', or 'daily_alt_odd'.")
	backupCreate.MarkFlagRequired("type")
	backupCreate.Flags().IntP("hour", "o", 0, "Hour value (0-23). Applicable to crons: 'daily', 'weekly', 'monthly', 'daily_alt_even', 'daily_alt_odd'")
	backupCreate.Flags().IntP("dow", "w", 0, "Day-of-week value (0-6). Applicable to crons: 'weekly'")
	backupCreate.Flags().IntP("dom", "m", 0, "Day-of-month value (1-28). Applicable to crons: 'monthly'")
	instanceCmd.AddCommand(backupCMD)

	// IPV4 Subcommands
	isoCmd := &cobra.Command{
		Use:   "iso",
		Short: "attach/detach ISOs to a given instance",
		Long:  ``,
	}
	isoCmd.AddCommand(isoStatus, isoAttach, isoDetach)
	isoAttach.Flags().StringP("iso-id", "i", "", "id of the ISO you wish to attach")
	isoAttach.MarkFlagRequired("iso-id")
	instanceCmd.AddCommand(isoCmd)

	ipv4Cmd := &cobra.Command{
		Use:   "ipv4",
		Short: "list/create/delete ipv4 on instance",
		Long:  ``,
	}
	ipv4Cmd.AddCommand(instanceIPV4List, createIpv4, deleteIpv4)
	createIpv4.Flags().Bool("reboot", false, "whether to reboot instance after adding ipv4 address")
	deleteIpv4.Flags().StringP("ipv4", "i", "", "ipv4 address you wish to delete")
	deleteIpv4.MarkFlagRequired("ipv4")
	instanceCmd.AddCommand(ipv4Cmd)

	// IPV6 Subcommands
	ipv6Cmd := &cobra.Command{
		Use:   "ipv6",
		Short: "commands for ipv6 on instance",
		Long:  ``,
	}
	ipv6Cmd.AddCommand(instanceIPV6List)
	instanceCmd.AddCommand(ipv6Cmd)

	// Plans SubCommands
	plansCmd := &cobra.Command{
		Use:   "plan",
		Short: "update/list plans for an instance",
		Long:  ``,
	}
	plansCmd.AddCommand(upgradePlan, upgradePlanList)
	upgradePlan.Flags().StringP("plan", "p", "", "plan id that you wish to upgrade to")
	upgradePlan.MarkFlagRequired("plan")
	instanceCmd.AddCommand(plansCmd)

	// ReverseDNS SubCommands
	reverseCmd := &cobra.Command{
		Use:   "reverse-dns",
		Short: "commands to handle reverse-dns on an instance",
		Long:  ``,
	}
	reverseCmd.AddCommand(defaultIpv4, listIpv6, deleteIpv6, setIpv4, setIpv6)
	defaultIpv4.Flags().StringP("ip", "i", "", "iPv4 address used in the reverse DNS update")
	defaultIpv4.MarkFlagRequired("ip")
	deleteIpv6.Flags().StringP("ip", "i", "", "ipv6 address you wish to delete")

	defaultIpv4.MarkFlagRequired("ip")
	setIpv4.Flags().StringP("ip", "i", "", "ip address you wish to set a reverse DNS on")
	setIpv4.Flags().StringP("entry", "e", "", "reverse dns entry")
	setIpv4.MarkFlagRequired("ip")
	setIpv4.MarkFlagRequired("entry")

	setIpv6.Flags().StringP("ip", "i", "", "ip address you wish to set a reverse DNS on")
	setIpv6.Flags().StringP("entry", "e", "", "reverse dns entry")
	setIpv6.MarkFlagRequired("ip")
	setIpv6.MarkFlagRequired("entry")
	instanceCmd.AddCommand(reverseCmd)

	userdataCmd := &cobra.Command{
		Use:   "user-data",
		Short: "commands to handle userdata on an instance",
		Long:  ``,
	}
	userdataCmd.AddCommand(setUserData, getUserData)
	setUserData.Flags().StringP("userdata", "d", "/dev/stdin", "file to read userdata from")
	instanceCmd.AddCommand(userdataCmd)

	return instanceCmd
}

var instanceStart = &cobra.Command{
	Use:   "start <instanceID>",
	Short: "starts an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Instance.Start(context.Background(), id); err != nil {
			fmt.Printf("error starting instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Started up instance")
	},
}

var instanceStop = &cobra.Command{
	Use:   "stop <instanceID>",
	Short: "stops an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Instance.Halt(context.Background(), id); err != nil {
			fmt.Printf("error stopping instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Stopped the instance")
	},
}

var instanceRestart = &cobra.Command{
	Use:   "restart <instanceID>",
	Short: "restart an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Instance.Reboot(context.Background(), id); err != nil {
			fmt.Printf("error rebooting instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Rebooted instance")
	},
}

var instanceReinstall = &cobra.Command{
	Use:   "reinstall <instanceID>",
	Short: "reinstall an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Instance.Reinstall(context.Background(), id); err != nil {
			fmt.Printf("error reinstalling instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Reinstalled instance")
	},
}

var instanceTag = &cobra.Command{
	Use:   "tag <instanceID>",
	Short: "add/modify tag on instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		tag, _ := cmd.Flags().GetString("tag")
		options := &govultr.InstanceUpdateReq{
			Tag: tag,
		}

		if err := client.Instance.Update(context.Background(), id, options); err != nil {
			fmt.Printf("error adding tag to instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Tagged instance with : %s\n", tag)
	},
}

var instanceDelete = &cobra.Command{
	Use:     "delete <instanceID>",
	Short:   "delete/destroy an instance",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Instance.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted instance")
	},
}

var instanceLabel = &cobra.Command{
	Use:   "label <instanceID>",
	Short: "label an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		label, _ := cmd.Flags().GetString("label")
		options := &govultr.InstanceUpdateReq{
			Label: label,
		}

		if err := client.Instance.Update(context.Background(), id, options); err != nil {
			fmt.Printf("error labeling instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Labeled instance with : %s\n", label)
	},
}

var instanceBandwidth = &cobra.Command{
	Use:   "bandwidth <instanceID>",
	Short: "bandwidth for instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		bw, err := client.Instance.GetBandwidth(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting bandwidth for instance : %v\n", err)
			os.Exit(1)
		}

		printer.InstanceBandwidth(bw)
	},
}

var instanceIPV4List = &cobra.Command{
	Use:     "list <instanceID>",
	Aliases: []string{"v4"},
	Short:   "list ipv4 for an instance",
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := getPaging(cmd)
		v4, meta, err := client.Instance.ListIPv4(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error getting ipv4 info : %v\n", err)
			os.Exit(1)
		}

		printer.InstanceIPV4(v4, meta)
	},
}

var instanceIPV6List = &cobra.Command{
	Use:     "list <instanceID>",
	Aliases: []string{"v6"},
	Short:   "list ipv6 for an instance",
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := getPaging(cmd)
		v6, meta, err := client.Instance.ListIPv6(context.TODO(), id, options)
		if err != nil {
			fmt.Printf("error getting ipv6 info : %v\n", err)
			os.Exit(1)
		}

		printer.InstanceIPV6(v6, meta)
	},
}

var instanceList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all available instances",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		s, meta, err := client.Instance.List(context.TODO(), options)
		if err != nil {
			fmt.Printf("error getting list of instances : %v\n", err)
			os.Exit(1)
		}

		printer.InstanceList(s, meta)
	},
}

var instanceInfo = &cobra.Command{
	Use:   "get <instanceID>",
	Short: "get info about a specific instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		s, err := client.Instance.Get(context.TODO(), id)
		if err != nil {
			fmt.Printf("error getting instance : %v\n", err)
			os.Exit(1)
		}

		printer.Instance(s)
	},
}

var updateFwgGroup = &cobra.Command{
	Use:   "update-firewall-group",
	Short: "assign a firewall group to instance",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("instance-id")
		fwgID, _ := cmd.Flags().GetString("firewall-group-id")

		options := &govultr.InstanceUpdateReq{
			FirewallGroupID: fwgID,
		}

		if err := client.Instance.Update(context.TODO(), id, options); err != nil {
			fmt.Printf("error setting firewall group : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated firewall group")
	},
}

var osUpdate = &cobra.Command{
	Use:   "change <instanceID>",
	Short: "changes operating system",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		osID, _ := cmd.Flags().GetInt("operatingSystems")

		options := &govultr.InstanceUpdateReq{
			OsID: osID,
		}

		if err := client.Instance.Update(context.TODO(), id, options); err != nil {
			fmt.Printf("error updating operatingSystems : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated OS")
	},
}

var osUpdateList = &cobra.Command{
	Use:   "list <instanceID>",
	Short: "available operating systems an instance can change to.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		list, err := client.Instance.GetUpgrades(context.TODO(), id)

		if err != nil {
			fmt.Printf("error listing available operatingSystems : %v\n", err)
			os.Exit(1)
		}

		printer.OsList(list.OS)
	},
}

var appUpdate = &cobra.Command{
	Use:   "change <instanceID>",
	Short: "changes application",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		appID, _ := cmd.Flags().GetInt("app")

		options := &govultr.InstanceUpdateReq{
			AppID: appID,
		}

		if err := client.Instance.Update(context.TODO(), id, options); err != nil {
			fmt.Printf("error updating application : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated Application")
	},
}

var appUpdateList = &cobra.Command{
	Use:   "list <instanceID>",
	Short: "available apps an instance can change to.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		list, err := client.Instance.GetUpgrades(context.TODO(), id)

		if err != nil {
			fmt.Printf("error listing available applications : %v\n", err)
			os.Exit(1)
		}

		printer.AppList(list.Applications)
	},
}

var backupGet = &cobra.Command{
	Use:   "get <instanceID>",
	Short: "get backup schedules on a given instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		info, err := client.Instance.GetBackupSchedule(context.TODO(), id)
		if err != nil {
			fmt.Printf("error getting application info : %v\n", err)
			os.Exit(1)
		}

		printer.BackupsGet(info)
	},
}

var backupCreate = &cobra.Command{
	Use:   "create <instanceID>",
	Short: "create backup schedule on a given instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		crontType, _ := cmd.Flags().GetString("type")
		hour, _ := cmd.Flags().GetInt("hour")
		dow, _ := cmd.Flags().GetInt("dow")
		dom, _ := cmd.Flags().GetInt("dom")

		backup := &govultr.BackupScheduleReq{
			Type: crontType,
			Hour: hour,
			Dow:  dow,
			Dom:  dom,
		}

		if err := client.Instance.SetBackupSchedule(context.TODO(), id, backup); err != nil {
			fmt.Printf("error creating backup schedule : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Created backup schedule")
	},
}

var isoStatus = &cobra.Command{
	Use:   "status <instanceID>",
	Short: "current ISO state",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		info, err := client.Instance.ISOStatus(context.TODO(), id)
		if err != nil {
			fmt.Printf("error getting iso state info : %v\n", err)
			os.Exit(1)
		}

		printer.IsoStatus(info)
	},
}

var isoAttach = &cobra.Command{
	Use:   "attach <instanceID>",
	Short: "attach ISO to instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		iso, _ := cmd.Flags().GetString("iso-id")
		if err := client.Instance.AttachISO(context.TODO(), id, iso); err != nil {
			fmt.Printf("error attaching iso : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("ISO has been attached")
	},
}

var isoDetach = &cobra.Command{
	Use:   "detach <instanceID>",
	Short: "detach ISO from instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Instance.DetachISO(context.TODO(), id); err != nil {
			fmt.Printf("error detaching iso : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("ISO has been detached")
	},
}

var instanceRestore = &cobra.Command{
	Use:   "restore <instanceID>",
	Short: "restore instance from backup/snapshot",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		backup, _ := cmd.Flags().GetString("backup")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		options := &govultr.RestoreReq{}

		if backup == "" && snapshot == "" {
			fmt.Println("at least one flag must be provided (snapshot or backup)")
			os.Exit(1)
		} else if backup != "" && snapshot != "" {
			fmt.Println("one flag must be provided not both (snapshot or backup)")
			os.Exit(1)
		}

		if snapshot != "" {
			options.SnapshotID = snapshot
		} else {
			options.BackupID = backup
		}

		if err := client.Instance.Restore(context.TODO(), id, options); err != nil {
			fmt.Printf("error restoring instance : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Instance has been restored")
	},
}

var createIpv4 = &cobra.Command{
	Use:   "create <instanceID>",
	Short: "create ipv4 for instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		reboot, _ := cmd.Flags().GetBool("reboot")

		_, err := client.Instance.CreateIPv4(context.TODO(), id, govultr.BoolToBoolPtr(reboot))
		if err != nil {
			fmt.Printf("error creating ipv4 : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("IPV4 has been created")
	},
}

var deleteIpv4 = &cobra.Command{
	Use:     "delete <instanceID>",
	Short:   "delete ipv4 for instance",
	Aliases: []string{"destroy"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ip, _ := cmd.Flags().GetString("ipv4")

		if err := client.Instance.DeleteIPv4(context.TODO(), id, ip); err != nil {
			fmt.Printf("error deleting ipv4 : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("IPV4 has been deleted")
	},
}

var upgradePlan = &cobra.Command{
	Use:   "upgrade <instanceID>",
	Short: "upgrade plan for instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		plan, _ := cmd.Flags().GetString("plan")

		options := &govultr.InstanceUpdateReq{
			Plan: plan,
		}

		if err := client.Instance.Update(context.TODO(), id, options); err != nil {
			fmt.Printf("error upgrading plans : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Upgraded plan")
	},
}

var upgradePlanList = &cobra.Command{
	Use:   "list <instanceID>",
	Short: "available plans an instance can upgrade to.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		list, err := client.Instance.GetUpgrades(context.TODO(), id)

		if err != nil {
			fmt.Printf("error listing available plans : %v\n", err)
			os.Exit(1)
		}

		printer.PlansList(list.Plans)
	},
}

var defaultIpv4 = &cobra.Command{
	Use:   "default-ipv4 <instanceID>",
	Short: "Set a reverse DNS entry for an IPv4 address of an instance to the original setting",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ip, _ := cmd.Flags().GetString("ip")

		if err := client.Instance.DefaultReverseIPv4(context.TODO(), id, ip); err != nil {
			fmt.Printf("error setting default reverse dns : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Set default reserve dns")
	},
}

var listIpv6 = &cobra.Command{
	Use:   "list-ipv6 <instanceID>",
	Short: "List the IPv6 reverse DNS entries for an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		rip, err := client.Instance.ListReverseIPv6(context.TODO(), id)
		if err != nil {
			fmt.Printf("error getting the reverse ipv6 list: %v\n", err)
			os.Exit(1)
		}
		printer.ReverseIpv6(rip)
	},
}

var deleteIpv6 = &cobra.Command{
	Use:     "delete-ipv6 <instanceID>",
	Short:   "Remove a reverse DNS entry for an IPv6 address for an instance",
	Aliases: []string{"destroy-ipv6"},
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ip, _ := cmd.Flags().GetString("ip")
		if err := client.Instance.DeleteReverseIPv6(context.TODO(), id, ip); err != nil {
			fmt.Printf("error deleting reverse ipv6 entry : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted reverse DNS IPV6 entry")
	},
}

var setIpv4 = &cobra.Command{
	Use:   "set-ipv4 <instanceID>",
	Short: "Set a reverse DNS entry for an IPv4 address for an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ip, _ := cmd.Flags().GetString("ip")
		entry, _ := cmd.Flags().GetString("entry")

		options := &govultr.ReverseIP{
			IP:      ip,
			Reverse: entry,
		}

		if err := client.Instance.CreateReverseIPv4(context.TODO(), id, options); err != nil {
			fmt.Printf("error setting reverse dns ipv4 entry : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Set reverse DNS entry for ipv4 address")
	},
}

var setIpv6 = &cobra.Command{
	Use:   "set-ipv6 <instanceID>",
	Short: "Set a reverse DNS entry for an IPv6 address for an instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ip, _ := cmd.Flags().GetString("ip")
		entry, _ := cmd.Flags().GetString("entry")

		options := &govultr.ReverseIP{
			IP:      ip,
			Reverse: entry,
		}

		if err := client.Instance.CreateReverseIPv6(context.TODO(), id, options); err != nil {
			fmt.Printf("error setting reverse dns ipv6 entry : %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Set reverse DNS entry for ipv6 address")
	},
}

var instanceCreate = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		plan, _ := cmd.Flags().GetString("plan")
		osID, _ := cmd.Flags().GetInt("operatingSystems")

		// Optional
		ipxe, _ := cmd.Flags().GetString("ipxe")
		iso, _ := cmd.Flags().GetString("iso")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		script, _ := cmd.Flags().GetString("script-id")
		ipv6, _ := cmd.Flags().GetBool("ipv6")
		privateNetwork, _ := cmd.Flags().GetBool("private-network")
		networks, _ := cmd.Flags().GetStringArray("network")
		label, _ := cmd.Flags().GetString("label")
		ssh, _ := cmd.Flags().GetStringArray("ssh-keys")
		backup, _ := cmd.Flags().GetBool("auto-backup")
		app, _ := cmd.Flags().GetInt("app")
		userData, _ := cmd.Flags().GetString("userdata")
		notify, _ := cmd.Flags().GetBool("notify")
		ddos, _ := cmd.Flags().GetBool("ddos")
		ipv4, _ := cmd.Flags().GetString("reserved-ipv4")
		host, _ := cmd.Flags().GetString("host")
		tag, _ := cmd.Flags().GetString("tag")
		fwg, _ := cmd.Flags().GetString("firewall-group")

		osOptions := map[string]interface{}{"iso_id": iso, "os_id": osID, "app_id": app, "snapshot_id": snapshot}

		if iso != "" {
			osOptions["iso_id"] = iso
		}

		osOption, err := optionCheck(osOptions)
		if err != nil {
			fmt.Printf("error creating instance : %v\n", err)
			os.Exit(1)
		}

		opt := &govultr.InstanceCreateReq{
			Plan:                 plan,
			Region:               region,
			IPXEChainURL:         ipxe,
			ISOID:                iso,
			SnapshotID:           snapshot,
			ScriptID:             script,
			AttachPrivateNetwork: networks,
			Label:                label,
			SSHKeys:              ssh,
			AppID:                app,
			UserData:             userData,
			ReservedIPv4:         ipv4,
			Hostname:             host,
			Tag:                  tag,
			FirewallGroupID:      fwg,
			EnableIPv6:           govultr.BoolToBoolPtr(false),
			DDOSProtection:       govultr.BoolToBoolPtr(false),
			ActivationEmail:      govultr.BoolToBoolPtr(false),
			Backups:              "disabled",
			EnablePrivateNetwork: govultr.BoolToBoolPtr(false),
		}

		// If no osOptions were selected and osID has a real value then set the osOptions to os_id
		if osOption == "os_id" && osID != 0 {
			opt.OsID = osID
		} else if osOption == "" && osID == 0 {
			fmt.Printf("error creating instance: an os_id, snapshot_id, iso_id, or app_id must be provided\n")
			os.Exit(1)
		}

		if ipv6 {
			opt.EnableIPv6 = govultr.BoolToBoolPtr(true)
		}
		if ddos {
			opt.DDOSProtection = govultr.BoolToBoolPtr(true)
		}
		if notify {
			opt.ActivationEmail = govultr.BoolToBoolPtr(true)
		}
		if backup {
			opt.Backups = "enabled"
		}
		if privateNetwork {
			opt.EnablePrivateNetwork = govultr.BoolToBoolPtr(true)
		}

		if userData != "" {
			opt.UserData = base64.StdEncoding.EncodeToString([]byte(userData))
		}

		//region, plan, osOpt, opt
		instance, err := client.Instance.Create(context.TODO(), opt)
		if err != nil {
			fmt.Printf("error creating instance : %v\n", err)
			os.Exit(1)
		}

		printer.Instance(instance)
	},
}

var setUserData = &cobra.Command{
	Use:   "set <instanceID>",
	Short: "Set the plain text user-data of an instance",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		userData, _ := cmd.Flags().GetString("userdata")

		rawData, err := ioutil.ReadFile(userData)
		if err != nil {
			fmt.Printf("error reading user-data : %v\n", err)
			os.Exit(1)
		}

		options := &govultr.InstanceUpdateReq{
			UserData: base64.StdEncoding.EncodeToString(rawData),
		}

		if err = client.Instance.Update(context.TODO(), args[0], options); err != nil {
			fmt.Printf("error setting user-data : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Set user-data for instance")
	},
}

var getUserData = &cobra.Command{
	Use:   "get <instanceID>",
	Short: "Get the user-data of an instance",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		userData, err := client.Instance.GetUserData(context.TODO(), args[0])
		if err != nil {
			fmt.Printf("error getting user-data : %v\n", err)
			os.Exit(1)
		}

		printer.UserData(userData)
	},
}

func optionCheck(options map[string]interface{}) (string, error) {
	var result []string
	for k, v := range options {
		switch v.(type) {
		case int:
			if v != 0 {
				result = append(result, k)
			}
		case string:
			if v != "" {
				result = append(result, k)
			}
		}
	}

	if len(result) > 1 {
		return "", fmt.Errorf("too many options have been selected : %v : please select one", result)
	}

	// Return back an empty slice so we can possibly add in osID
	if len(result) == 0 {
		return "", nil
	}

	return result[0], nil
}
