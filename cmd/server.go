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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr"
	"github.com/vultr/vultr-cli/cmd/printer"
)

// Server represents the server command
func Server() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "commands to interact with servers on vultr",
		Long:  ``,
	}

	serverCmd.AddCommand(serverStart, serverStop, serverRestart, serverReinstall, serverTag, serverDelete, serverLabel, serverBandwidth, serverList, serverInfo, updateFwgGroup, serverRestore, serverCreate)

	serverTag.Flags().StringP("tag", "t", "", "tag you want to set for a given instance")
	serverTag.MarkFlagRequired("tag")

	serverLabel.Flags().StringP("label", "l", "", "label you want to set for a given instance")
	serverLabel.MarkFlagRequired("label")

	serverIPV4List.Flags().StringP("public", "p", "", "include information about the public network adapter : True or False")

	updateFwgGroup.Flags().StringP("instance-id", "i", "", "instance id of the server you want to use")
	updateFwgGroup.Flags().StringP("firewall-group-id", "f", "", "firewall group id that you want to assign. 0 Value will unset the firewall-group")
	updateFwgGroup.MarkFlagRequired("instance-id")
	updateFwgGroup.MarkFlagRequired("firewall-group-id")

	serverRestore.Flags().StringP("backup", "b", "", "id of backup you wish to restore the instance with")
	serverRestore.Flags().StringP("snapshot", "s", "", "id of snapshot you wish to restore the instance with")

	serverCreate.Flags().IntP("region", "r", 0, "region id you wish to have the instance created in")
	serverCreate.Flags().IntP("plan", "p", 0, "plan id you wish to the instance to have")
	serverCreate.Flags().IntP("os", "o", 0, "os id you wish the instance to have")
	serverCreate.MarkFlagRequired("region")
	serverCreate.MarkFlagRequired("plan")
	// Optional Params
	serverCreate.Flags().StringP("ipxe", "", "", "ff you've selected the 'custom' operating system, this can be set to chainload the specified URL on bootup")
	serverCreate.Flags().IntP("iso", "", 0, "iso ID you want to create the instance with")
	serverCreate.Flags().StringP("snapshot", "", "", "snapshot ID you want to create the instance with")
	serverCreate.Flags().StringP("script-id", "", "", "script id of the startup script")
	serverCreate.Flags().StringP("ipv6", "", "", "enable ipv6 | true or false")
	serverCreate.Flags().StringP("private-network", "", "", "enable private network | true or false")
	serverCreate.Flags().StringArrayP("network", "", []string{}, "network IDs you want to assign to the instance")
	serverCreate.Flags().StringP("label", "l", "", "label you want to give this instance")
	serverCreate.Flags().StringArrayP("ssh-keys", "s", []string{}, "ssh keys you want to assign to the instance")
	serverCreate.Flags().StringP("auto-backup", "b", "", "enable auto backups | true or false")
	serverCreate.Flags().StringP("app", "a", "", "application ID you want this instance to have")
	serverCreate.Flags().StringP("userdata", "u", "", "userdata you want to give this instance")
	serverCreate.Flags().StringP("notify", "n", "", "notify when server has been created | true or false")
	serverCreate.Flags().StringP("ddos", "d", "", "enable ddos protected | true or false")
	serverCreate.Flags().StringP("reserved-ipv4", "", "", "ip address of the floating IP to use as the main IP for this instance")
	serverCreate.Flags().StringP("host", "", "", "The hostname to assign to this instance")
	serverCreate.Flags().StringP("tag", "t", "", "The tag to assign to this instance")
	serverCreate.Flags().StringP("firewall-group", "", "", "The firewall group to assign to this instance")

	// Sub commands for OS
	osCmd := &cobra.Command{
		Use:   "os",
		Short: "list and update operating system for an instance",
		Long:  ``,
	}

	osCmd.AddCommand(osList, osUpdate)
	osUpdate.Flags().StringP("os", "o", "", "operating system ID you wish to use")
	osUpdate.MarkFlagRequired("os")
	serverCmd.AddCommand(osCmd)

	// Sub commands for App
	appCMD := &cobra.Command{
		Use:   "app",
		Short: "list and update application for an instance",
		Long:  ``,
	}
	appCMD.AddCommand(appList, appUpdate, appInfo)
	appUpdate.Flags().StringP("app", "a", "", "appplication ID you wish to use")
	appUpdate.MarkFlagRequired("app")
	serverCmd.AddCommand(appCMD)

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
	serverCmd.AddCommand(backupCMD)

	// IPV4 Subcommands
	isoCmd := &cobra.Command{
		Use:   "iso",
		Short: "attach/detach ISOs to a given instance",
		Long:  ``,
	}
	isoCmd.AddCommand(isoStatus, isoAttach, isoDetach)
	isoAttach.Flags().StringP("iso-id", "i", "", "id of the ISO you wish to attach")
	isoAttach.MarkFlagRequired("iso-id")
	serverCmd.AddCommand(isoCmd)

	ipv4Cmd := &cobra.Command{
		Use:   "ipv4",
		Short: "list/create/delete ipv4 on instance",
		Long:  ``,
	}
	ipv4Cmd.AddCommand(serverIPV4List, createIpv4, deleteIpv4)
	createIpv4.Flags().StringP("reboot", "r", "no", "whether to reboot server after adding ipv4 address")
	deleteIpv4.Flags().StringP("ipv4", "i", "", "ipv4 address you wish to delete")
	deleteIpv4.MarkFlagRequired("ipv4")
	serverCmd.AddCommand(ipv4Cmd)

	// IPV6 Subcommands
	ipv6Cmd := &cobra.Command{
		Use:   "ipv6",
		Short: "commands for ipv6 on instance",
		Long:  ``,
	}
	ipv6Cmd.AddCommand(serverIPV6List)
	serverCmd.AddCommand(ipv6Cmd)

	// Plans SubCommands
	plansCmd := &cobra.Command{
		Use:   "plans",
		Short: "update/list plans for an instance",
		Long:  ``,
	}
	plansCmd.AddCommand(listPlans, upgradePlan)
	upgradePlan.Flags().StringP("plan", "p", "", "plan id that you wish to upgrade to")
	upgradePlan.MarkFlagRequired("plan")
	serverCmd.AddCommand(plansCmd)

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
	serverCmd.AddCommand(reverseCmd)

	userdataCmd := &cobra.Command{
		Use:   "userdata",
		Short: "commands to handle userdata on an instance",
		Long:  ``,
	}
	userdataCmd.AddCommand(setUserData, getUserData)
	setUserData.Flags().StringP("userdata", "d", "/dev/stdin", "file to read userdata from")
	serverCmd.AddCommand(userdataCmd)

	return serverCmd
}

var serverStart = &cobra.Command{
	Use:   "start <instanceID>",
	Short: "starts a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Start(context.TODO(), id)

		if err != nil {
			fmt.Printf("error starting server : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Started up server")
	},
}

var serverStop = &cobra.Command{
	Use:   "stop <instanceID>",
	Short: "stops a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Halt(context.TODO(), id)

		if err != nil {
			fmt.Printf("error stopping server : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Stopped the server")
	},
}

var serverRestart = &cobra.Command{
	Use:   "restart <instanceID>",
	Short: "restart a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Reboot(context.TODO(), id)

		if err != nil {
			fmt.Printf("error rebooting server : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Rebooted server")
	},
}

var serverReinstall = &cobra.Command{
	Use:   "reinstall <instanceID>",
	Short: "reinstall a server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		err := client.Server.Reinstall(context.TODO(), id)

		if err != nil {
			fmt.Printf("error reinstalling server : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Reinstalled server")
	},
}

var serverTag = &cobra.Command{
	Use:   "tag <instanceID>",
	Short: "add/modify tag on server",
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
		err := client.Server.SetTag(context.TODO(), id, tag)

		if err != nil {
			fmt.Printf("error adding tag to server : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Tagged server with : %s\n", tag)
	},
}

var serverDelete = &cobra.Command{
	Use:     "delete <instanceID>",
	Short:   "delete/destroy a server",
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
		err := client.Server.Delete(context.TODO(), id)

		if err != nil {
			fmt.Printf("error deleting server : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted server")
	},
}

var serverLabel = &cobra.Command{
	Use:   "label <instanceID>",
	Short: "label a server",
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
		err := client.Server.SetLabel(context.TODO(), id, label)

		if err != nil {
			fmt.Printf("error labeling server : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Labeled server with : %s\n", label)
	},
}

var serverBandwidth = &cobra.Command{
	Use:   "bandwidth <instanceID>",
	Short: "bandwidth for server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		bw, err := client.Server.Bandwidth(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting bandwidth for server : %v\n", err)
			os.Exit(1)
		}

		printer.ServerBandwidth(bw)
	},
}

var serverIPV4List = &cobra.Command{
	Use:     "list <instanceID>",
	Aliases: []string{"v4"},
	Short:   "list ipv4 for a server",
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		public, _ := cmd.Flags().GetString("public")

		pub := false
		if strings.ToLower(public) == "true" {
			pub = true
		}

		v4, err := client.Server.IPV4Info(context.TODO(), id, pub)

		if err != nil {
			fmt.Printf("error getting ipv4 info : %v\n", err)
			os.Exit(1)
		}

		printer.ServerIPV4(v4)
	},
}

var serverIPV6List = &cobra.Command{
	Use:     "list <instanceID>",
	Aliases: []string{"v6"},
	Short:   "list ipv6 for a server",
	Long:    ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		v6, err := client.Server.IPV6Info(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting ipv6 info : %v\n", err)
			os.Exit(1)
		}

		printer.ServerIPV6(v6)
	},
}

var serverList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all available servers",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := client.Server.List(context.TODO())

		if err != nil {
			fmt.Printf("error getting list of servers : %v\n", err)
			os.Exit(1)
		}

		printer.ServerList(s)
	},
}

var serverInfo = &cobra.Command{
	Use:   "info <instanceID>",
	Short: "info about a specific server",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		s, err := client.Server.GetServer(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting server : %v\n", err)
			os.Exit(1)
		}

		printer.ServerInfo(s)
	},
}

var updateFwgGroup = &cobra.Command{
	Use:   "update-firewall-group",
	Short: "assign a firewall group to server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("instance-id")
		fwgID, _ := cmd.Flags().GetString("firewall-group-id")

		err := client.Server.SetFirewallGroup(context.TODO(), id, fwgID)

		if err != nil {
			fmt.Printf("error setting firewall group : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Added firewall-group")
	},
}

var osList = &cobra.Command{
	Use:   "list <instanceID>",
	Short: "list available operating systems this instance can be changed to",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		o, err := client.Server.ListOS(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting os list : %v\n", err)
			os.Exit(1)
		}

		printer.OsList(o)
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
		osID, _ := cmd.Flags().GetString("os")

		err := client.Server.ChangeOS(context.TODO(), id, osID)

		if err != nil {
			fmt.Printf("error updating os : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated OS")
	},
}

var appList = &cobra.Command{
	Use:   "list <instanceID>",
	Short: "list available applications this instance can be changed to",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		a, err := client.Server.ListApps(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting os list : %v\n", err)
			os.Exit(1)
		}

		printer.AppList(a)
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
		appID, _ := cmd.Flags().GetString("app")

		err := client.Server.ChangeApp(context.TODO(), id, appID)

		if err != nil {
			fmt.Printf("error updating application : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated Application")
	},
}

var appInfo = &cobra.Command{
	Use:   "info <instanceID>",
	Short: "gets information about the application on the instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		info, err := client.Server.AppInfo(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting application info : %v\n", err)
			os.Exit(1)
		}

		printer.ServerAppInfo(info)
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

		info, err := client.Server.GetBackupSchedule(context.TODO(), id)

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
			errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		crontType, _ := cmd.Flags().GetString("type")
		hour, _ := cmd.Flags().GetInt("hour")
		dow, _ := cmd.Flags().GetInt("dow")
		dom, _ := cmd.Flags().GetInt("dom")

		backup := &govultr.BackupSchedule{
			CronType: crontType,
			Hour:     hour,
			Dow:      dow,
			Dom:      dom,
		}

		err := client.Server.SetBackupSchedule(context.TODO(), id, backup)

		if err != nil {
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

		info, err := client.Server.IsoStatus(context.TODO(), id)

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

		err := client.Server.IsoAttach(context.TODO(), id, iso)

		if err != nil {
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
		err := client.Server.IsoDetach(context.TODO(), id)

		if err != nil {
			fmt.Printf("error detaching iso : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("ISO has been detached")
	},
}

var serverRestore = &cobra.Command{
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

		if backup == "" && snapshot == "" {
			fmt.Println("at least one flag must be provided (snapshot or backup)")
			os.Exit(1)
		} else if backup != "" && snapshot != "" {
			fmt.Println("one flag must be provided not both (snapshot or backup)")
			os.Exit(1)
		}

		var err error

		if snapshot != "" {
			err = client.Server.RestoreSnapshot(context.TODO(), id, snapshot)
		} else {
			err = client.Server.RestoreBackup(context.TODO(), id, backup)
		}

		if err != nil {
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
		reboot, _ := cmd.Flags().GetString("reboot")
		_, err := client.Server.AddIPV4(context.TODO(), id, reboot)

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
		err := client.Server.DestroyIPV4(context.TODO(), id, ip)

		if err != nil {
			fmt.Printf("error deleting ipv4 : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("IPV4 has been deleted")
	},
}

var listPlans = &cobra.Command{
	Use:   "list <instanceID>",
	Short: "list available plans for instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		plans, err := client.Server.ListUpgradePlan(context.TODO(), id)

		if err != nil {
			fmt.Printf("error getting plans : %v\n", err)
			os.Exit(1)
		}

		printer.PlansList(plans)
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
		err := client.Server.UpgradePlan(context.TODO(), id, plan)

		if err != nil {
			fmt.Printf("error upgrading plans : %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Upgraded plan")
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
		err := client.Server.SetDefaultReverseIPV4(context.TODO(), id, ip)

		if err != nil {
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
		rip, err := client.Server.ListReverseIPV6(context.TODO(), id)

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
		err := client.Server.DeleteReverseIPV6(context.TODO(), id, ip)

		if err != nil {
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
		err := client.Server.SetReverseIPV4(context.TODO(), id, ip, entry)

		if err != nil {
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
		err := client.Server.SetReverseIPV6(context.TODO(), id, ip, entry)

		if err != nil {
			fmt.Printf("error setting reverse dns ipv6 entry : %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Set reverse DNS entry for ipv6 address")
	},
}

var serverCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a server instance",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			osAppID  = 186
			osIsoID  = 159
			osSnapID = 164
		)

		region, _ := cmd.Flags().GetInt("region")
		plan, _ := cmd.Flags().GetInt("plan")
		osID, _ := cmd.Flags().GetInt("os")

		// Optional
		ipxe, _ := cmd.Flags().GetString("ipxe")
		iso, _ := cmd.Flags().GetInt("iso")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		script, _ := cmd.Flags().GetString("script-id")
		ipv6, _ := cmd.Flags().GetString("ipv6")
		privateNetwork, _ := cmd.Flags().GetString("private-network")
		networks, _ := cmd.Flags().GetStringArray("network")
		label, _ := cmd.Flags().GetString("label")
		ssh, _ := cmd.Flags().GetStringArray("ssh-keys")
		backup, _ := cmd.Flags().GetString("auto-backup")
		app, _ := cmd.Flags().GetString("app")
		userData, _ := cmd.Flags().GetString("userdata")
		notify, _ := cmd.Flags().GetString("notify")
		ddos, _ := cmd.Flags().GetString("ddos")
		ipv4, _ := cmd.Flags().GetString("reserved-ipv4")
		host, _ := cmd.Flags().GetString("host")
		tag, _ := cmd.Flags().GetString("tag")
		fwg, _ := cmd.Flags().GetString("firewall-group")

		osOptions := map[string]string{"app_id": app, "snapshot_id": snapshot}

		if iso != 0 {
			osOptions["iso_id"] = strconv.Itoa(iso)
		}

		osOption, err := optionCheck(osOptions)

		if err != nil {
			fmt.Printf("error creating instance : %v", err)
			os.Exit(1)
		}

		opt := &govultr.ServerOptions{
			IPXEChain:            ipxe,
			IsoID:                iso,
			SnapshotID:           snapshot,
			ScriptID:             script,
			NetworkID:            networks,
			Label:                label,
			SSHKeyIDs:            ssh,
			AppID:                app,
			UserData:             userData,
			ReservedIPV4:         ipv4,
			Hostname:             host,
			Tag:                  tag,
			FirewallGroupID:      fwg,
			EnableIPV6:           false,
			DDOSProtection:       false,
			NotifyActivate:       false,
			AutoBackups:          false,
			EnablePrivateNetwork: false,
		}

		// If no osOptions were selected and osID has a real value then set the osOptions to os_id
		if osOption == "" && osID != 0 {
			osOption = "os_id"
		} else if osOption == "" && osID == 0 {
			fmt.Printf("error creating instance: an os ID must be provided")
			os.Exit(1)
		}

		var osOpt int

		switch osOption {
		case "os_id":
			osOpt = osID

		case "app_id":
			osOpt = osAppID

		case "iso_id":
			osOpt = osIsoID

		case "snapshot_id":
			osOpt = osSnapID

		default:
			fmt.Println("Error occurred while getting your intended os type")
			os.Exit(1)
		}

		if strings.ToLower(ipv6) == "true" {
			opt.EnableIPV6 = true
		}
		if strings.ToLower(ddos) == "true" {
			opt.DDOSProtection = true
		}
		if strings.ToLower(notify) == "true" {
			opt.NotifyActivate = true
		}
		if strings.ToLower(ipv6) == "true" {
			opt.EnableIPV6 = true
		}
		if strings.ToLower(backup) == "true" {
			opt.AutoBackups = true
		}
		if strings.ToLower(privateNetwork) == "true" {
			opt.EnablePrivateNetwork = true
		}

		server, err := client.Server.Create(context.TODO(), region, plan, osOpt, opt)

		if err != nil {
			fmt.Printf("error creating instance : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance created - ID : %s\n", server.InstanceID)
	},
}

var setUserData = &cobra.Command{
	Use:   "set <instanceID>",
	Short: "Set the user-data of a server",
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

		err = client.Server.SetUserData(context.TODO(), args[0], base64.StdEncoding.EncodeToString(rawData))

		if err != nil {
			fmt.Printf("error setting user-data : %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Set user-data for server")
	},
}

var getUserData = &cobra.Command{
	Use:   "get <instanceID>",
	Short: "Get the user-data of a server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide an instanceID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		userData, err := client.Server.GetUserData(context.TODO(), args[0])

		if err != nil {
			fmt.Printf("error getting user-data : %v\n", err)
			os.Exit(1)
		}

		printer.UserData(userData)
	},
}

func optionCheck(options map[string]string) (string, error) {
	result := []string{}
	for k, v := range options {
		if v != "" {
			result = append(result, k)
		}
	}

	if len(result) > 1 {
		return "", fmt.Errorf("Too many options have been selected : %v : please select one", result)
	}

	// Return back an empty slice so we can possibly add in osID
	if len(result) == 0 {
		return "", nil
	}

	return result[0], nil
}
