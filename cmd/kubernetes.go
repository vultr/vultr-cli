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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	kubernetesLong    = `Get all available commands for Kubernetes`
	kubernetesExample = `
	# Full example
	vultr-cli kubernetes
	`

	createLong    = `Create kubernetes cluster on your Vultr account`
	createExample = `
	# Full Example
	vultr-cli kubernetes create --label="my-cluster" --region="ewr" --version="v1.20.0+1" \
		--node-pools="quantity:3,plan:vc2-1c-2gb,label:my-nodepool,tag:my-tag"

	# Shortened with alias commands
	vultr-cli k c -l="my-cluster" -r="ewr" -v="v1.20.0+1" -n="quantity:3,plan:vc2-1c-2gb,label:my-nodepool,tag:my-tag"
	`

	getLong    = `Get a single kubernetes cluster from your account`
	getExample = `
	# Full example
	vultr-cli kubernetes get ffd31f18-5f77-454c-9064-212f942c3c34

	# Shortened with alias commands
	vultr-cli k g ffd31f18-5f77-454c-9064-212f942c3c34
	`

	listLong    = `Get all kubernetes clusters available on your Vultr account`
	listExample = `
	# Full example
	vultr-cli kubernetes list

	# Full example with paging
	vultr-cli kubernetes list --per-page=1 --cursor="bmV4dF9fQU1T"

	# Shortened with alias commands
	vultr-cli k l

	# Summarized view
	vultr-cli kubernetes list --summarize
	`

	updateLong    = `Update a specific kubernetes cluster on your Vultr Account`
	updateExample = `
	# Full example
	vultr-cli kubernetes update ffd31f18-5f77-454c-9065-212f942c3c35 --label="updated-label"

	# Shortened with alias commands
	vultr-cli k u ffd31f18-5f77-454c-9065-212f942c3c35 -l="updated-label"
	`

	deleteLong    = `Delete a specific kubernetes cluster off your Vultr Account`
	deleteExample = `
	# Full example
	vultr-cli kubernetes delete ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli k d ffd31f18-5f77-454c-9065-212f942c3c35'
	`

	deleteWithResourcesLong    = `Delete a specific kubernetes cluster and all linked load balancers and block storages off your Vultr Account`
	deleteWithResourcesExample = `
	# Full example
	vultr-cli kubernetes delete-with-resources ffd31f18-5f77-454c-9065-212f942c3c35
	`

	getConfigLong    = `Returns a base64 encoded config of a specified kubernetes cluster on your Vultr Account`
	getConfigExample = `
	
		# Full example
		vultr-cli kubernetes config ffd31f18-5f77-454c-9065-212f942c3c35
		vultr-cli kubernetes config --output-file /your/path/ ffd31f18-5f77-454c-9065-212f942c3c35
	
		# Use the default config location (~/.kube/config)
		vultr-cli kubernetes config --output-file ~/.kube/config ffd31f18-5f77-454c-9065-212f942c3c35
		vultr-cli kubernetes config ffd31f18-5f77-454c-9065-212f942c3c35
	
	
		# Shortened with alias commands
		vultr-cli k config ffd31f18-5f77-454c-9065-212f942c3c35
		vultr-cli k config -o /your/path/ ffd31f18-5f77-454c-9065-212f942c3c35
		`

	getVersionsLong    = `Returns a list of supported kubernetes versions you can deploy`
	getVersionsExample = `
	# Full example
	vultr-cli kubernetes versions

	# Shortened with alias commands
	vultr-cli k v
	`

	upgradesLong    = `Display available kubernetes upgrade commands`
	upgradesExample = `
	# Full example
	vultr-cli kubernetes upgrades

	# Shortened example with aliases
	vultr-cli k e
	`

	getUpgradesLong    = `Returns a list of available kubernetes version the cluster can be upgraded to`
	getUpgradesExample = `
	# Full example
	vultr-cli kubernetes upgrades list d4908765-b82a-4e7d-83d9-c0bc4c6a36d0

	# Shortened with alias commands
	vultr-cli k e l d4908765-b82a-4e7d-83d9-c0bc4c6a36d0
	`

	upgradeLong    = `Initiate an upgrade of the kubernetes version on a given cluster`
	upgradeExample = `
	# Full example
	vultr-cli kubernetes upgrades start d4908765-b82a-4e7d-83d9-c0bc4c6a36d0 --version="v1.23.5+3"

	# Shortened with alias commands
	vultr-cli k e s d4908765-b82a-4e7d-83d9-c0bc4c6a36d0 -v="v1.23.5+3"
	`

	nodepoolLong    = `Get all available commands for Kubernetes node pools`
	nodepoolExample = `
	# Full example
	vultr-cli kubernetes node-pool

	# Shortened with alias commands
	vultr-cli k n
	`

	createNPLong    = `Create node pool for your kubernetes cluster on your Vultr account`
	createNPExample = `
	# Full Example
	vultr-cli kubernetes node-pool create ffd31f18-5f77-454c-9064-212f942c3c34 --label="nodepool" --quantity=3  --plan="vc2-1c-2gb"

	# Shortened with alias commands
	vultr-cli k n c ffd31f18-5f77-454c-9064-212f942c3c34 -l="nodepool" -q=3  -p="vc2-1c-2gb"
	`

	getNPLong    = `Get a node pool in a single kubernetes cluster from your account`
	getNPExample = `
	# Full example
	vultr-cli kubernetes node-pool get ffd31f18-5f77-454c-9064-212f942c3c34 abd31f18-3f77-454c-9064-212f942c3c34
	# Shortened with alias commands
	vultr-cli k n g ffd31f18-5f77-454c-9064-212f942c3c34 abd31f18-3f77-454c-9064-212f942c3c34
	`

	listNPLong    = `Get all nodepools from a kubernetes cluster on your Vultr account`
	listNPExample = `
	# Full example
	vultr-cli kubernetes node-pool list ffd31f18-5f77-454c-9064-212f942c3c34

	# Full example with paging
	vultr-cli kubernetes node-pool list ffd31f18-5f77-454c-9064-212f942c3c34 --per-page=1 --cursor="bmV4dF9fQU1T"

	# Shortened with alias commands
	vultr-cli k n l ffd31f18-5f77-454c-9064-212f942c3c34
	`

	updateNPLong    = `Update a specific node pool in a kubernetes cluster on your Vultr Account`
	updateNPExample = `
	# Full example
	vultr-cli kubernetes node-pool update ffd31f18-5f77-454c-9064-212f942c3c34 abd31f18-3f77-454c-9064-212f942c3c34 --quantity=4

	# Shortened with alias commands
	vultr-cli k n u ffd31f18-5f77-454c-9065-212f942c3c35 abd31f18-3f77-454c-9064-212f942c3c34 --q=4
	`

	deleteNPLong    = `Delete a specific node pool in a kubernetes cluster off your Vultr Account`
	deleteNPExample = `
	# Full example
	vultr-cli kubernetes node-pool delete ffd31f18-5f77-454c-9065-212f942c3c35 abd31f18-3f77-454c-9064-212f942c3c34

	# Shortened with alias commands
	vultr-cli k n d ffd31f18-5f77-454c-9065-212f942c3c35 abd31f18-3f77-454c-9064-212f942c3c34'
	`

	nodepoolInstanceLong    = `Get all available commands for Kubernetes node pool instances`
	nodepoolInstanceExample = `
	# Full example
	vultr-cli kubernetes node-pool node

	# Shortened with alias commands
	vultr-cli k n node
	`

	deleteNPInstanceLong    = `Delete a specific node pool instance in a kubernetes cluster from your Vultr Account`
	deleteNPInstanceExample = `
	# Full example
	vultr-cli kubernetes node-pool node delete ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli k n node d ffd31f18-5f77-454c-9065-212f942c3c35
	`

	deleteNPInstanceRecycleLong    = `Recycles a specific node pool instance in a kubernetes cluster from your Vultr Account`
	deleteNPInstanceRecycleExample = `
	# Full example
	vultr-cli kubernetes node-pool node recycle ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli k n node r ffd31f18-5f77-454c-9065-212f942c3c35
	`
)

// Kubernetes represents the kubernetes command
func Kubernetes() *cobra.Command { //nolint: funlen
	kubernetesCmd := &cobra.Command{
		Use:     "kubernetes",
		Aliases: []string{"k"},
		Short:   "kubernetes is used to access kubernetes commands",
		Long:    kubernetesLong,
		Example: kubernetesExample,
	}

	kubernetesCmd.AddCommand(k8Create, k8Get, k8List, k8GetConfig, k8Update, k8Delete, k8DeleteWithResources, k8GetVersions)
	k8Create.Flags().StringP("label", "l", "", "label for your kubernetes cluster")
	k8Create.Flags().StringP("region", "r", "", "region you want your kubernetes cluster to be located in")
	k8Create.Flags().StringP("version", "v", "", "the kubernetes version you want for your cluster")
	k8GetConfig.Flags().StringVarP(&kubeconfigFilePath, "output-file", "o", "", "Optional file path to write kubeconfig to")
	k8Create.Flags().StringArrayP(
		"node-pools",
		"n",
		[]string{},
		`a comma-separated, key-value pair list of node pools. At least one node pool is required. At least one node is
		required in node pool. Use / between each new node pool.
		E.g: 'plan:vhf-8c-32gb,label:mynodepool,tag:my-tag,quantity:3/plan:vhf-8c-32gb,label:mynodepool2,quantity:3`,
	)

	if err := k8Create.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking kubernetes create 'label' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := k8Create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking kubernetes create 'region' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := k8Create.MarkFlagRequired("version"); err != nil {
		fmt.Printf("error marking kubernetes create 'version' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := k8Create.MarkFlagRequired("node-pools"); err != nil {
		fmt.Printf("error marking kubernetes create 'ns-primary' flag required: %v\n", err)
		os.Exit(1)
	}

	k8List.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	k8List.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")
	k8List.Flags().BoolP("summarize", "", false, "(optional) Summarize the list output. One line per cluster.")

	k8Update.Flags().StringP("label", "l", "", "label for your kubernetes cluster")
	if err := k8Update.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking kubernetes update 'label' flag required: %v\n", err)
		os.Exit(1)
	}

	// Sub command for upgrade functions
	k8UpgradeCmd := &cobra.Command{
		Use:     "upgrades",
		Aliases: []string{"upgrade", "e"},
		Short:   `upgrade commands for kubernetes version upgrades`,
		Long:    upgradesLong,
		Example: upgradesExample,
	}

	k8UpgradeCmd.AddCommand(k8Upgrade, k8GetUpgrades)
	k8Upgrade.Flags().StringP("version", "v", "", "the version to upgrade the cluster to")
	if err := k8Upgrade.MarkFlagRequired("version"); err != nil {
		fmt.Printf("error marking kubernetes upgrade 'version' flag required: %v\n", err)
		os.Exit(1)
	}
	kubernetesCmd.AddCommand(k8UpgradeCmd)

	// Node Pools SubCommands
	nodepoolsCmd := &cobra.Command{
		Use:     "node-pool",
		Aliases: []string{"n"},
		Short:   "node pools commands for a kubernetes cluster",
		Long:    nodepoolLong,
		Example: nodepoolExample,
	}

	// Node Pools
	npCreate.Flags().StringP("label", "l", "", "label you want for your node pool.")
	npCreate.Flags().StringP("tag", "t", "", "tag you want for your node pool.")
	npCreate.Flags().StringP("plan", "p", "", "the plan you want for your node pool.")
	npCreate.Flags().IntP("quantity", "q", 1, "Number of nodes in your node pool. Note that at least one node is required for a node pool.")
	npCreate.Flags().BoolP("auto-scaler", "", false, "Enable the auto scaler with your cluster")
	npCreate.Flags().IntP("min-nodes", "", 1, "Minimum nodes for auto scaler")
	npCreate.Flags().IntP("max-nodes", "", 1, "Maximum nodes for auto scaler")

	if err := npCreate.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'label' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := npCreate.MarkFlagRequired("quantity"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'quantity' flag required: %v\n", err)
		os.Exit(1)
	}
	if err := npCreate.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'plan' flag required: %v\n", err)
		os.Exit(1)
	}

	npList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	npList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	npUpdate.Flags().IntP("quantity", "q", 1, "Number of nodes in your node pool. Note that at least one node is required for a node pool.")
	npUpdate.Flags().StringP("tag", "t", "", "tag you want for your node pool.")
	npUpdate.Flags().BoolP("auto-scaler", "", false, "Enable the auto scaler with your cluster")
	npUpdate.Flags().IntP("min-nodes", "", 1, "Minimum nodes for auto scaler")
	npUpdate.Flags().IntP("max-nodes", "", 1, "Maximum nodes for auto scaler")

	// Node Instance SubCommands
	nodeCmd := &cobra.Command{
		Use:     "node",
		Short:   "delete/recycle instances in a cluster's node pool",
		Long:    nodepoolInstanceLong,
		Example: nodepoolInstanceExample,
	}

	nodeCmd.AddCommand(npInstanceDelete, npInstanceRecycle)
	nodepoolsCmd.AddCommand(nodeCmd, npCreate, npGet, npList, npDelete, npUpdate)
	kubernetesCmd.AddCommand(nodepoolsCmd)

	return kubernetesCmd
}

var k8Create = &cobra.Command{
	Use:     "create",
	Short:   "create kubernetes cluster",
	Long:    createLong,
	Example: createExample,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")
		region, _ := cmd.Flags().GetString("region")
		nodepools, _ := cmd.Flags().GetStringArray("node-pools")
		version, _ := cmd.Flags().GetString("version")

		nps, err := formatNodePools(nodepools)
		if err != nil {
			fmt.Printf("error creating kubernetes cluster: %v\n", err)
			os.Exit(1)
		}

		options := &govultr.ClusterReq{
			Label:     label,
			Region:    region,
			NodePools: nps,
			Version:   version,
		}

		kubernetes, _, err := client.Kubernetes.CreateCluster(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating kubernetes cluster : %v\n", err)
			os.Exit(1)
		}

		printer.Cluster(kubernetes)
	},
}

var k8List = &cobra.Command{
	Use:     "list <clusterID>",
	Short:   "list kubernetes clusters",
	Aliases: []string{"l"},
	Long:    listLong,
	Example: listExample,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		summarize, _ := cmd.Flags().GetBool("summarize")

		k8s, meta, _, err := client.Kubernetes.ListClusters(context.Background(), options)
		if err != nil {
			fmt.Printf("error listing kubernetes clusters : %v\n", err)
			os.Exit(1)
		}

		if summarize {
			printer.ClustersSummary(k8s, meta)
		} else {
			printer.Clusters(k8s, meta)
		}
	},
}

var k8Get = &cobra.Command{
	Use:     "get <clusterID>",
	Short:   "retrieves a kubernetes cluster",
	Long:    getLong,
	Example: getExample,
	Aliases: []string{"g"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		lb, _, err := client.Kubernetes.GetCluster(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting cluster : %v\n", err)
			os.Exit(1)
		}

		printer.Cluster(lb)
	},
}

var k8Update = &cobra.Command{
	Use:     "update <clusterID>",
	Short:   "updates a kubernetes cluster",
	Aliases: []string{"u"},
	Long:    updateLong,
	Example: updateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		label, _ := cmd.Flags().GetString("label")

		options := &govultr.ClusterReqUpdate{
			Label: label,
		}

		if err := client.Kubernetes.UpdateCluster(context.Background(), id, options); err != nil {
			fmt.Printf("error updating kubernetes cluster : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated kubernetes cluster")
	},
}

var k8Delete = &cobra.Command{
	Use:     "delete <clusterID>",
	Short:   "delete a kubernetes cluster",
	Aliases: []string{"destroy", "d"},
	Long:    deleteLong,
	Example: deleteExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Kubernetes.DeleteCluster(context.Background(), id); err != nil {
			fmt.Printf("error deleting kubernetes cluster : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("kubernetes cluster has been deleted")
	},
}

var k8DeleteWithResources = &cobra.Command{
	Use:     "delete-with-resources <clusterID>",
	Short:   "delete a kubernetes cluster and related resources",
	Long:    deleteWithResourcesLong,
	Example: deleteWithResourcesExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.Kubernetes.DeleteClusterWithResources(context.Background(), id); err != nil {
			fmt.Printf("error deleting kubernetes cluster : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("kubernetes cluster and related resources have been deleted")
	},
}

var kubeconfigFilePath string

const kubeconfigFilePermission = 0600

var k8GetConfig = &cobra.Command{
	Use:     "config <clusterID>",
	Short:   "gets a Kubernetes cluster's config",
	Long:    getConfigLong,
	Example: getConfigExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		config, _, err := client.Kubernetes.GetKubeConfig(context.Background(), id)
		if err != nil {
			fmt.Printf("error retrieving kube config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(config.KubeConfig)
		fmt.Println()

		if kubeconfigFilePath != "" {
			fmt.Printf("Writing kubeconfig to: %s\n", kubeconfigFilePath)

			// Write the kubeconfig to the specified file path
			err := os.WriteFile(kubeconfigFilePath, []byte(config.KubeConfig), kubeconfigFilePermission)
			if err != nil {
				fmt.Printf("\nError writing kubeconfig to %s: %v\n", kubeconfigFilePath, err)
				os.Exit(1)
			} else {
				fmt.Printf("Kubeconfig successfully written to %s\n", kubeconfigFilePath)
			}
		} else {
			home, _ := os.UserHomeDir()
			defaultKubeconfigPath := filepath.Join(home, ".kube", "config")
			fmt.Printf("Writing kubeconfig to the default path: %s\n", defaultKubeconfigPath)

			err := os.WriteFile(defaultKubeconfigPath, []byte(config.KubeConfig), kubeconfigFilePermission)
			if err != nil {
				fmt.Printf("\nError writing kubeconfig to %s: %v\n", defaultKubeconfigPath, err)
				os.Exit(1)
			} else {
				fmt.Printf("Kubeconfig successfully written to %s\n", defaultKubeconfigPath)
			}
		}
	},
}

var k8GetVersions = &cobra.Command{
	Use:     "versions",
	Short:   "gets supported kubernetes versions",
	Long:    getVersionsLong,
	Example: getVersionsExample,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		versions, _, err := client.Kubernetes.GetVersions(context.Background())
		if err != nil {
			fmt.Printf("error retrieving supported versions : %v\n", err)
			os.Exit(1)
		}

		printer.K8Versions(versions)
	},
}

var k8GetUpgrades = &cobra.Command{
	Use:     "list <clusterID>",
	Short:   "gets available upgrades for a cluster",
	Long:    getUpgradesLong,
	Example: getUpgradesExample,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		upgrades, _, err := client.Kubernetes.GetUpgrades(context.Background(), id)
		if err != nil {
			fmt.Printf("error retrieving available upgrades : %v\n", err)
			os.Exit(1)
		}

		printer.K8Upgrades(upgrades)
	},
}

var k8Upgrade = &cobra.Command{
	Use:     "start <clusterID>",
	Short:   "perform upgrade on a cluster",
	Long:    upgradeLong,
	Example: upgradeExample,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		version, _ := cmd.Flags().GetString("version")

		options := &govultr.ClusterUpgradeReq{
			UpgradeVersion: version,
		}

		if err := client.Kubernetes.Upgrade(context.Background(), id, options); err != nil {
			fmt.Printf("error performing cluster upgrade : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("kubernetes cluster upgrade has been initiated")
	},
}

var npCreate = &cobra.Command{
	Use:     "create <clusterID>",
	Short:   "creates a node pool in a kubernetes cluster",
	Aliases: []string{"c"},
	Long:    createNPLong,
	Example: createNPExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		quantity, _ := cmd.Flags().GetInt("quantity")
		label, _ := cmd.Flags().GetString("label")
		tag, _ := cmd.Flags().GetString("tag")
		plan, _ := cmd.Flags().GetString("plan")
		autoscaler, _ := cmd.Flags().GetBool("auto-scaler")
		min, _ := cmd.Flags().GetInt("min-nodes")
		max, _ := cmd.Flags().GetInt("max-nodes")

		options := &govultr.NodePoolReq{
			NodeQuantity: quantity,
			Label:        label,
			Plan:         plan,
			Tag:          tag,
			AutoScaler:   govultr.BoolToBoolPtr(false),
			MinNodes:     min,
			MaxNodes:     max,
		}

		if autoscaler {
			options.AutoScaler = govultr.BoolToBoolPtr(true)
		}

		np, _, err := client.Kubernetes.CreateNodePool(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error creating cluster node pool : %v\n", err)
			os.Exit(1)
		}

		printer.NodePool(np)
	},
}

var npUpdate = &cobra.Command{
	Use:     "update <clusterID> <nodePoolID>",
	Short:   "updates a cluster's node pool quantity",
	Aliases: []string{"u"},
	Long:    updateNPLong,
	Example: updateNPExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a clusterID and nodePoolID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		nodeID := args[1]
		quantity, _ := cmd.Flags().GetInt("quantity")
		tag, _ := cmd.Flags().GetString("tag")
		autoscaler, _ := cmd.Flags().GetBool("auto-scaler")
		min, _ := cmd.Flags().GetInt("min-nodes")
		max, _ := cmd.Flags().GetInt("max-nodes")

		options := &govultr.NodePoolReqUpdate{
			NodeQuantity: quantity,
			Tag:          govultr.StringToStringPtr(tag),
			AutoScaler:   govultr.BoolToBoolPtr(false),
			MinNodes:     min,
			MaxNodes:     max,
		}

		if autoscaler {
			options.AutoScaler = govultr.BoolToBoolPtr(true)
		}

		np, _, err := client.Kubernetes.UpdateNodePool(context.Background(), id, nodeID, options)
		if err != nil {
			fmt.Printf("error updating cluster node pool : %v\n", err)
			os.Exit(1)
		}

		printer.NodePool(np)
	},
}

var npDelete = &cobra.Command{
	Use:     "delete <clusterID> <nodeID>",
	Short:   "delete a cluster node pool",
	Aliases: []string{"destroy", "d"},
	Long:    deleteNPLong,
	Example: deleteNPExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a clusterID and nodeID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		nodePoolID := args[1]

		if err := client.Kubernetes.DeleteNodePool(context.Background(), id, nodePoolID); err != nil {
			fmt.Printf("error deleting cluster nodepool : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("kubernetes cluster has been deleted")
	},
}

var npInstanceDelete = &cobra.Command{
	Use:     "delete <clusterID> <nodePoolID> <nodeID>",
	Short:   "deletes a node in a cluster's node pool",
	Aliases: []string{"destroy", "d"},
	Long:    deleteNPInstanceLong,
	Example: deleteNPInstanceExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("please provide a clusterID, nodePoolID, and nodeID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		nodePoolID := args[1]
		nodeID := args[2]

		if err := client.Kubernetes.DeleteNodePoolInstance(context.Background(), id, nodePoolID, nodeID); err != nil {
			fmt.Printf("error deleting node pool node : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("node pool node deleted")
	},
}

var npInstanceRecycle = &cobra.Command{
	Use:     "recycle <clusterID> <nodePoolID> <nodeID>",
	Short:   "recycles a node in a cluster's node pool",
	Aliases: []string{"r"},
	Long:    deleteNPInstanceRecycleLong,
	Example: deleteNPInstanceRecycleExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("please provide a clusterID, nodePoolID, and nodeID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		nodePoolID := args[1]
		nodeID := args[2]

		if err := client.Kubernetes.RecycleNodePoolInstance(context.Background(), id, nodePoolID, nodeID); err != nil {
			fmt.Printf("error recycling node pool node : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("node pool node recycled")
	},
}

var npList = &cobra.Command{
	Use:     "list <clusterID>",
	Short:   "list nodepools",
	Aliases: []string{"l"},
	Long:    listNPLong,
	Example: listNPExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a clusterID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := getPaging(cmd)
		nps, meta, _, err := client.Kubernetes.ListNodePools(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error listing cluster node pools : %v\n", err)
			os.Exit(1)
		}

		printer.NodePools(nps, meta)
	},
}

var npGet = &cobra.Command{
	Use:     "get <clusterID> <nodePoolID>",
	Short:   "get nodepool in kubernetes cluster",
	Aliases: []string{"g"},
	Long:    getNPLong,
	Example: getNPExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a clusterID and nodePoolID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		nodeID := args[1]
		np, _, err := client.Kubernetes.GetNodePool(context.Background(), id, nodeID)
		if err != nil {
			fmt.Printf("error getting cluster node pool : %v\n", err)
			os.Exit(1)
		}

		printer.NodePool(np)
	},
}

// formatNodePools parses node pools into proper format
func formatNodePools(nodePools []string) ([]govultr.NodePoolReq, error) {
	var formattedList []govultr.NodePoolReq
	npList := strings.Split(nodePools[0], "/")

	for _, r := range npList {
		nodeData := strings.Split(r, ",")

		if len(nodeData) < 3 || len(nodeData) > 7 {
			return nil, fmt.Errorf(
				`unable to format node pool. each node pool must include label, quantity, and plan.
				Optionally you can include tag, auto-scaler, min-nodes and max-nodes`,
			)
		}

		formattedNodeData, errFormat := formatNodeData(nodeData)
		if errFormat != nil {
			return nil, errFormat
		}

		formattedList = append(formattedList, *formattedNodeData)
	}

	return formattedList, nil
}

// formatNodeData loops over the parse strings for a node and returns the formatted struct
func formatNodeData(node []string) (*govultr.NodePoolReq, error) {
	nodeData := &govultr.NodePoolReq{}
	for _, f := range node {
		nodeDataKeyVal := strings.Split(f, ":")

		if len(nodeDataKeyVal) != 2 && len(nodeDataKeyVal) != 3 {
			return nil, fmt.Errorf("invalid node pool format")
		}

		field := nodeDataKeyVal[0]
		val := nodeDataKeyVal[1]

		switch {
		case field == "plan":
			nodeData.Plan = val
		case field == "quantity":
			port, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool quantity: %v", err)
			}
			nodeData.NodeQuantity = port
		case field == "label":
			nodeData.Label = val
		case field == "tag":
			nodeData.Tag = val
		case field == "auto-scaler":
			v, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool auto-scaler: %v", err)
			}
			nodeData.AutoScaler = govultr.BoolToBoolPtr(v)
		case field == "min-nodes":
			v, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool min-nodes: %v", err)
			}
			nodeData.MinNodes = v
		case field == "max-nodes":
			v, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for max-nodes: %v", err)
			}
			nodeData.MaxNodes = v
		}
	}

	return nodeData, nil
}
