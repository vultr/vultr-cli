// Package kubernetes provides functionality for the CLI to control VKE clusters
package kubernetes

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `Get all available commands for Kubernetes`
	example = `
	# Full example
	vultr-cli kubernetes
	`

	createLong = `Create kubernetes cluster on your Vultr account`
	//nolint:lll
	createExample = `
	# Full example
	vultr-cli kubernetes create --label="my-cluster" --region="ewr" --version="v1.29.2+1" \
		--node-pools="quantity:3,plan:vc2-2c-4gb,label:my-nodepool,tag:my-tag"

	# Shortened with alias commands
	vultr-cli k c -l="my-cluster" -r="ewr" -v="v1.29.2+1" -n="quantity:3,plan:vc2-2c-4gb,label:my-nodepool,tag:my-tag"

	# Node pool options
	The --node-pools option allows you to pass in options for any number of
	node pools when creating a cluster. The options are passed in a delimited
	string.  Different node pools are delimited by a slash (/). The options for
	each node pool are delimited by a comma (,) and each option is defined by
	colon (:).  If provided, the node pool options can also parse out the
	node-labels params which are delimited by a pipe (|).

	Available options are documented in the 'kubernetes node-pool create --help' 

	For example:

	Multiple node pools
	--node-pools="quantity:1,plan:vc2-4c-8gb,label:main-node-pool/quantity:5,plan:vc2-2c-4gb,label:worker-pool,auto-scaler:true,min-nodes:5,max-nodes:10"

	Using node labels 
	--node-pools="quantity:5,plan:vc2-2c-4gb,label:worker-pool,auto-scaler:true,min-nodes:5,max-nodes:10,node-labels:application=identity-service|worker-size=small"
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

	# Delete a specific kubernetes cluster and all linked load balancers and block storages off your Vultr Account
	vultr-cli kubernetes delete-with-resources ffd31f18-5f77-454c-9065-212f942c3c35
	`
	getConfigLong    = `Returns a base64 encoded config of a specified kubernetes cluster on your Vultr Account`
	getConfigExample = `
	
	# Full example
	vultr-cli kubernetes config ffd31f18-5f77-454c-9065-212f942c3c35
	vultr-cli kubernetes config ffd31f18-5f77-454c-9065-212f942c3c35 --output-file /your/path/

	# Shortened with alias commands
	vultr-cli k config ffd31f18-5f77-454c-9065-212f942c3c35
	vultr-cli k config  ffd31f18-5f77-454c-9065-212f942c3c35 -o /your/path/
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
	vultr-cli kubernetes upgrades start d4908765-b82a-4e7d-83d9-c0bc4c6a36d0 --version="v1.29.2+1"

	# Shortened with alias commands
	vultr-cli k e s d4908765-b82a-4e7d-83d9-c0bc4c6a36d0 -v="v1.29.2+1"
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
	vultr-cli kubernetes node-pool create ffd31f18-5f77-454c-9064-212f942c3c34 --label="nodepool" --quantity=3  \
		--plan="vc2-1c-2gb" --node-labels="application=id-service,environment=development"

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
	vultr-cli kubernetes node-pool update ffd31f18-5f77-454c-9064-212f942c3c34 abd31f18-3f77-454c-9064-212f942c3c34 --quantity=4 \
		--node-labels="application=id-service,environment=development"

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

	nodeLong    = `Get all available commands for Kubernetes node pool nodes`
	nodeExample = `
	# Full example
	vultr-cli kubernetes node-pool node

	# Shortened with alias commands
	vultr-cli k n node
	`

	nodeDeleteLong    = `Delete a specific node pool node in a kubernetes cluster`
	nodeDeleteExample = `
	# Full example
	vultr-cli kubernetes node-pool node delete ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli k n node d ffd31f18-5f77-454c-9065-212f942c3c35
	`

	nodeRecycleLong    = `Recycles a specific node pool node in a kubernetes cluster`
	nodeRecycleExample = `
	# Full example
	vultr-cli kubernetes node-pool node recycle ffd31f18-5f77-454c-9065-212f942c3c35

	# Shortened with alias commands
	vultr-cli k n node r ffd31f18-5f77-454c-9065-212f942c3c35
	`
)

const (
	kubeconfigFilePermission = 0600
	kubeconfigDirPermission  = 0755
)

// NewCmdKubernetes provides the CLI command for VKE functions
func NewCmdKubernetes(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "kubernetes",
		Aliases: []string{"k"},
		Short:   "Commands to manage kubernetes clusters",
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
		Short:   "List kubernetes clusters",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			summarize, errSu := cmd.Flags().GetBool("summarize")
			if errSu != nil {
				return fmt.Errorf("error parsing flag 'summarize' for kubernetes list : %v", errSu)
			}

			k8s, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving kubernetes clusters list : %v", err)
			}

			var data printer.ResourceOutput
			if summarize {
				data = &ClustersSummaryPrinter{Clusters: k8s, Meta: meta}
			} else {
				data = &ClustersPrinter{Clusters: k8s, Meta: meta}
			}

			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)
	list.Flags().BoolP("summarize", "", false, "(optional) Summarize the list output. One line per cluster.")

	// Get
	get := &cobra.Command{
		Use:     "get <Cluster ID>",
		Short:   "Retrieve a kubernetes cluster",
		Long:    getLong,
		Example: getExample,
		Aliases: []string{"g"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			k8, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving kubernetes cluster : %v", err)
			}

			data := &ClusterPrinter{Cluster: k8}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create kubernetes cluster",
		Long:    createLong,
		Example: createExample,
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for kubernetes cluster create : %v", errLa)
			}

			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'region' for kubernetes cluster create : %v", errRe)
			}

			nodepools, errNP := cmd.Flags().GetStringArray("node-pools")
			if errNP != nil {
				return fmt.Errorf("error parsing flag 'node-pools' for kubernetes cluster create : %v", errNP)
			}

			version, errVe := cmd.Flags().GetString("version")
			if errVe != nil {
				return fmt.Errorf("error parsing flag 'version' for kubernetes cluster create : %v", errVe)
			}

			ha, errHi := cmd.Flags().GetBool("high-avail")
			if errHi != nil {
				return fmt.Errorf("error parsing flag 'high-avail' for kubernetes cluster create : %v", errHi)
			}

			fw, errFw := cmd.Flags().GetBool("enable-firewall")
			if errHi != nil {
				return fmt.Errorf("error parsing flag 'enable-firewall' for kubernetes cluster create : %v", errFw)
			}

			nps, errFm := formatNodePools(nodepools)
			if errFm != nil {
				return fmt.Errorf("error in node pool formating : %v", errFm)
			}

			o.CreateReq = &govultr.ClusterReq{
				Label:           label,
				Region:          region,
				NodePools:       nps,
				Version:         version,
				HAControlPlanes: ha,
				EnableFirewall:  fw,
			}

			k8, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating kubernetes cluster : %v", err)
			}

			data := &ClusterPrinter{Cluster: k8}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("label", "l", "", "label for your kubernetes cluster")
	if err := create.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking kubernetes create 'label' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("region", "r", "", "region you want your kubernetes cluster to be located in")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking kubernetes create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("version", "v", "", "the kubernetes version you want for your cluster")
	if err := create.MarkFlagRequired("version"); err != nil {
		fmt.Printf("error marking kubernetes create 'version' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().Bool(
		"high-avail",
		false,
		`(optional, default false) whether or not the cluster should be deployed with multiple, 
highly available, control planes`,
	)

	create.Flags().BoolP(
		"enable-firewall",
		"f",
		false,
		`(optional, default false) whether a firewall group should be created for the cluster`,
	)

	create.Flags().StringArrayP(
		"node-pools",
		"n",
		[]string{},
		`a comma-separated, key-value pair list of node pools. At least one node pool is required. At least one node is
required in node pool. Use / between each new node pool.  E.g: 
'plan:vhf-8c-32gb,label:mynodepool,tag:my-tag,quantity:3/plan:vhf-8c-32gb,label:mynodepool2,quantity:3`,
	)
	if err := create.MarkFlagRequired("node-pools"); err != nil {
		fmt.Printf("error marking kubernetes create 'ns-primary' flag required: %v", err)
		os.Exit(1)
	}

	// Update
	update := &cobra.Command{
		Use:     "update <Cluster ID>",
		Short:   "Update a kubernetes cluster",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for kubernetes cluster update : %v", errLa)
			}

			o.UpdateReq = &govultr.ClusterReqUpdate{
				Label: label,
			}

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating kubernetes cluster : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes cluster has been updated"), nil)

			return nil
		},
	}

	update.Flags().StringP("label", "l", "", "label for your kubernetes cluster")
	if err := update.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking kubernetes update 'label' flag required: %v", err)
		os.Exit(1)
	}

	// Delete
	del := &cobra.Command{
		Use:     "delete <Cluster ID>",
		Short:   "Delete a kubernetes cluster",
		Aliases: []string{"destroy", "d"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			withRes, errRe := cmd.Flags().GetBool("delete-resources")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'delete-resource' for kubernetes cluster delete: %v", errRe)
			}

			if withRes {
				if err := o.delWithRes(); err != nil {
					return fmt.Errorf("error deleting kubernetes cluster and resources : %v", err)
				}
			} else {
				if err := o.del(); err != nil {
					return fmt.Errorf("error deleting kubernetes cluster : %v", err)
				}
			}

			o.Base.Printer.Display(printer.Info("Kubernetes cluster has been deleted"), nil)

			return nil
		},
	}

	del.Flags().BoolP("delete-resources", "r", false, "delete a kubernetes cluster and related resources")

	// Config
	config := &cobra.Command{
		Use:     "config <Cluster ID>",
		Short:   "Get a kubernetes cluster config",
		Long:    getConfigLong,
		Example: getConfigExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a clusterID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path, errPa := cmd.Flags().GetString("output-file")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'output-file' for kubernetes cluster config : %v", errPa)
			}

			config, err := o.config()
			if err != nil {
				return fmt.Errorf("error retrieving kubernetes cluster config : %v", err)
			}

			if path != "" {
				dir := filepath.Dir(path)
				if errDi := os.MkdirAll(dir, kubeconfigDirPermission); errDi != nil {
					return fmt.Errorf("error creating directory for kubeconfig : %v", errDi)
				}

				kubeConfigData, errDe := base64.StdEncoding.DecodeString(config.KubeConfig)
				if errDe != nil {
					return fmt.Errorf("error decoding kubeconfig : %v", errDe)
				}

				if errWr := os.WriteFile(path, kubeConfigData, kubeconfigFilePermission); errWr != nil {
					return fmt.Errorf("error writing kubeconfig to %s : %v", path, errWr)
				}
			} else {
				data := &ConfigPrinter{Config: config}
				o.Base.Printer.Display(data, nil)
			}

			return nil
		},
	}

	config.Flags().StringP("output-file", "", "", "(optional) the file path to write kubeconfig to")

	// Versions
	versions := &cobra.Command{
		Use:     "versions",
		Short:   "List supported kubernetes versions",
		Long:    getVersionsLong,
		Example: getVersionsExample,
		Aliases: []string{"v"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// override parent pre-run auth check
			utils.SetOptions(o.Base, cmd, args)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := o.versions()
			if err != nil {
				return fmt.Errorf("error retrieving list of kubernetes versions : %v", err)
			}

			data := &VersionsPrinter{Versions: versions.Versions}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Upgrades
	upgrades := &cobra.Command{
		Use:     "upgrades",
		Aliases: []string{"upgrade", "e"},
		Short:   `Commands for kubernetes version upgrades`,
		Long:    upgradesLong,
		Example: upgradesExample,
	}

	// Upgrades List
	upgradesList := &cobra.Command{
		Use:     "list <Cluster ID>",
		Short:   "Get available upgrades for a cluster",
		Long:    getUpgradesLong,
		Example: getUpgradesExample,
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrades, err := o.upgrades()
			if err != nil {
				return fmt.Errorf("error retrieving the available kubernetes upgrades : %v", err)
			}

			data := &UpgradesPrinter{Upgrades: upgrades}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Upgrade start
	upgradeStart := &cobra.Command{
		Use:     "start <clusterID>",
		Short:   "Perform an upgrade on a cluster",
		Long:    upgradeLong,
		Example: upgradeExample,
		Aliases: []string{"s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			version, errVe := cmd.Flags().GetString("version")
			if errVe != nil {
				return fmt.Errorf("error parsing flag 'version' for kubernetes upgrade start : %v", errVe)
			}

			o.UpgradeReq = &govultr.ClusterUpgradeReq{
				UpgradeVersion: version,
			}

			if err := o.upgrade(); err != nil {
				return fmt.Errorf("error starting the kubernetes upgrade : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes upgrade has been intiated"), nil)

			return nil
		},
	}

	upgradeStart.Flags().StringP("version", "v", "", "the version to upgrade the cluster to")
	if err := upgradeStart.MarkFlagRequired("version"); err != nil {
		fmt.Printf("error marking kubernetes upgrade 'version' flag required: %v", err)
		os.Exit(1)
	}

	upgrades.AddCommand(
		upgradesList,
		upgradeStart,
	)

	// Node Pools
	nodepool := &cobra.Command{
		Use:     "node-pool",
		Aliases: []string{"n"},
		Short:   "Commands for kubernetes cluster node pools",
		Long:    nodepoolLong,
		Example: nodepoolExample,
	}

	// Node Pools List
	npList := &cobra.Command{
		Use:     "list <Cluster ID>",
		Short:   "List node pools",
		Aliases: []string{"l"},
		Long:    listNPLong,
		Example: listNPExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			nps, meta, err := o.nodePools()
			if err != nil {
				return fmt.Errorf("error getting node pool list : %v", err)
			}

			data := &NodePoolsPrinter{NodePools: nps, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	npList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	npList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)

	// Node Pool Get
	npGet := &cobra.Command{
		Use:     "get <Cluster ID> <Node Pool ID>",
		Short:   "Get a node pool in kubernetes cluster",
		Aliases: []string{"g"},
		Long:    getNPLong,
		Example: getNPExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			np, err := o.nodePool()
			if err != nil {
				return fmt.Errorf("error getting node pool : %v", err)
			}

			data := &NodePoolPrinter{NodePool: np}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Node Pool Create
	npCreate := &cobra.Command{
		Use:     "create <Cluster ID>",
		Short:   "Create a node pool in a kubernetes cluster",
		Aliases: []string{"c"},
		Long:    createNPLong,
		Example: createNPExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			quantity, errQu := cmd.Flags().GetInt("quantity")
			if errQu != nil {
				return fmt.Errorf("error parsing flag 'quantity' for kubernetes cluster node pool create : %v", errQu)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for kubernetes cluster node pool create : %v", errLa)
			}

			tag, errTa := cmd.Flags().GetString("tag")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tag' for kubernetes cluster node pool create : %v", errTa)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing flag 'plan' for kubernetes cluster node pool create : %v", errPl)
			}

			autoscaler, errAu := cmd.Flags().GetBool("auto-scaler")
			if errAu != nil {
				return fmt.Errorf("error parsing flag 'auto-scaler' for kubernetes cluster node pool create : %v", errAu)
			}

			minNodes, errMi := cmd.Flags().GetInt("min-nodes")
			if errMi != nil {
				return fmt.Errorf("error parsing flag 'min-nodes' for kubernetes cluster node pool create : %v", errMi)
			}

			maxNodes, errMa := cmd.Flags().GetInt("max-nodes")
			if errMa != nil {
				return fmt.Errorf("error parsing flag 'max-nodes' for kubernetes cluster node pool create : %v", errMa)
			}

			npLabels, errNl := cmd.Flags().GetStringToString("node-labels")
			if errNl != nil {
				return fmt.Errorf("error parsing flag 'node-labels' for kubernetes cluster node pool create : %v", errNl)
			}

			o.npCreateReq = &govultr.NodePoolReq{
				NodeQuantity: quantity,
				Label:        label,
				Plan:         plan,
				Tag:          tag,
				AutoScaler:   govultr.BoolToBoolPtr(false),
				MinNodes:     minNodes,
				MaxNodes:     maxNodes,
				Labels:       npLabels,
			}

			if autoscaler {
				o.npCreateReq.AutoScaler = govultr.BoolToBoolPtr(true)
			}

			np, err := o.nodePoolCreate()
			if err != nil {
				return fmt.Errorf("error creating kubernetes cluster node pool : %v", err)
			}

			data := &NodePoolPrinter{NodePool: np}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	npCreate.Flags().StringP("label", "l", "", "label you want for your node pool.")
	if err := npCreate.MarkFlagRequired("label"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'label' flag required: %v\n", err)
		os.Exit(1)
	}

	npCreate.Flags().StringP("tag", "t", "", "tag you want for your node pool.")

	npCreate.Flags().StringP("plan", "p", "", "the plan you want for your node pool.")
	if err := npCreate.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'plan' flag required: %v\n", err)
		os.Exit(1)
	}

	npCreate.Flags().IntP("quantity", "q", 1, "Number of nodes in your node pool. Note that at least one node is required for a node pool.")
	if err := npCreate.MarkFlagRequired("quantity"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'quantity' flag required: %v\n", err)
		os.Exit(1)
	}

	npCreate.Flags().BoolP("auto-scaler", "", false, "Enable the auto scaler with your cluster")
	npCreate.Flags().IntP("min-nodes", "", 1, "Minimum nodes for auto scaler")
	npCreate.Flags().IntP("max-nodes", "", 1, "Maximum nodes for auto scaler")
	npCreate.Flags().StringToString("node-labels", nil, "A key=value comma separated string of labels to apply to the nodes in this node pool")

	// Node Pool Update
	npUpdate := &cobra.Command{
		Use:     "update <Cluster ID> <Node Pool ID>",
		Short:   "Update a cluster node pool",
		Aliases: []string{"u"},
		Long:    updateNPLong,
		Example: updateNPExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			quantity, errQu := cmd.Flags().GetInt("quantity")
			if errQu != nil {
				return fmt.Errorf("error parsing flag 'quantity' for kubernetes cluster node pool update : %v", errQu)
			}

			tag, errTa := cmd.Flags().GetString("tag")
			if errTa != nil {
				return fmt.Errorf("error parsing flag 'tag' for kubernetes cluster node pool update : %v", errTa)
			}

			autoscaler, errAu := cmd.Flags().GetBool("auto-scaler")
			if errAu != nil {
				return fmt.Errorf("error parsing flag 'auto-scaler' for kubernetes cluster node pool update : %v", errAu)
			}

			minNodes, errMi := cmd.Flags().GetInt("min-nodes")
			if errMi != nil {
				return fmt.Errorf("error parsing flag 'min-nodes' for kubernetes cluster node pool update : %v", errMi)
			}

			maxNodes, errMa := cmd.Flags().GetInt("max-nodes")
			if errMa != nil {
				return fmt.Errorf("error parsing flag 'max-nodes' for kubernetes cluster node pool update : %v", errMa)
			}

			npLabels, errNl := cmd.Flags().GetStringToString("node-labels")
			if errNl != nil {
				return fmt.Errorf("error parsing flag 'node-labels' for kubernetes cluster node pool update : %v", errNl)
			}

			o.npUpdateReq = &govultr.NodePoolReqUpdate{}

			if cmd.Flags().Changed("quantity") {
				o.npUpdateReq.NodeQuantity = quantity
			}

			if cmd.Flags().Changed("auto-scaler") {
				o.npUpdateReq.AutoScaler = govultr.BoolToBoolPtr(autoscaler)
			}

			if cmd.Flags().Changed("tag") {
				o.npUpdateReq.Tag = govultr.StringToStringPtr(tag)
			}

			if cmd.Flags().Changed("min-nodes") {
				o.npUpdateReq.MinNodes = minNodes
			}

			if cmd.Flags().Changed("max-nodes") {
				o.npUpdateReq.MaxNodes = maxNodes
			}

			if cmd.Flags().Changed("node-labels") {
				o.npUpdateReq.Labels = npLabels
			}

			np, err := o.nodePoolUpdate()
			if err != nil {
				return fmt.Errorf("error updating kubernetes cluster node pool : %v", err)
			}

			data := &NodePoolPrinter{NodePool: np}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	npUpdate.Flags().IntP(
		"quantity",
		"q",
		1,
		"Number of nodes in your node pool. Note that at least one node is required for a node pool.",
	)
	npUpdate.Flags().StringP("tag", "t", "", "The tag the node pool")
	npUpdate.Flags().BoolP("auto-scaler", "", false, "Enable the auto scaler with your cluster")
	npUpdate.Flags().IntP("min-nodes", "", 1, "Minimum nodes for auto scaler")
	npUpdate.Flags().IntP("max-nodes", "", 1, "Maximum nodes for auto scaler")
	npUpdate.Flags().StringToString("node-labels", nil, "A key=value comma separated string of labels to apply to the nodes in this node pool")

	npUpdate.MarkFlagsOneRequired("quantity", "tag", "auto-scaler", "min-nodes", "max-nodes", "node-labels")

	// Node Pool Delete
	npDelete := &cobra.Command{
		Use:     "delete <Cluster ID> <Node Pool ID>",
		Short:   "Delete a cluster node pool",
		Aliases: []string{"destroy", "d"},
		Long:    deleteNPLong,
		Example: deleteNPExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.nodePoolDelete(); err != nil {
				return fmt.Errorf("error deleting kubernetes cluster node pool : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes node pool has been deleted"), nil)

			return nil
		},
	}

	// Node
	node := &cobra.Command{
		Use:     "node",
		Short:   "Commands to manage node pool nodes",
		Long:    nodeLong,
		Example: nodeExample,
	}

	// Node Pool Node Delete
	nodeDelete := &cobra.Command{
		Use:     "delete <Cluster ID> <Node Pool ID> <Node ID>",
		Short:   "Delete a node in a cluster node pool",
		Aliases: []string{"destroy", "d"},
		Long:    nodeDeleteLong,
		Example: nodeDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("please provide a cluster ID, a node pool ID and a node ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.nodePoolNodeDelete(); err != nil {
				return fmt.Errorf("error deleting kubernetes cluster node pool node : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes node pool node has been deleted"), nil)

			return nil
		},
	}

	// Node Pool Node Recycle
	nodeRecycle := &cobra.Command{
		Use:     "recycle <clusterID> <nodePoolID> <nodeID>",
		Short:   "Recycle a node in a cluster node pool",
		Aliases: []string{"r"},
		Long:    nodeRecycleLong,
		Example: nodeRecycleExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("please provide a cluster ID, a node pool ID and a node ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.nodePoolNodeRecycle(); err != nil {
				return fmt.Errorf("error recycling kubernetes cluster node pool node : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes node pool node has been recycled"), nil)

			return nil
		},
	}

	node.AddCommand(
		nodeDelete,
		nodeRecycle,
	)

	nodepool.AddCommand(
		npList,
		npGet,
		npCreate,
		npUpdate,
		npDelete,
		node,
	)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		config,
		nodepool,
		versions,
		upgrades,
	)

	return cmd
}

// formatNodePools parses node pools into proper format
func formatNodePools(nodePools []string) ([]govultr.NodePoolReq, error) {
	var formattedList []govultr.NodePoolReq
	npList := strings.Split(nodePools[0], "/")

	for _, r := range npList {
		nodeData := strings.Split(r, ",")

		if len(nodeData) < 3 || len(nodeData) > 8 {
			return nil, fmt.Errorf(
				`unable to format node pool. each node pool must include label, quantity, and plan.
Optionally you can include tag, node-labels, auto-scaler, min-nodes and max-nodes`,
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
func formatNodeData(node []string) (*govultr.NodePoolReq, error) { //nolint:gocyclo
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
		case field == "node-labels":
			nodeData.Labels = formatNodeLabels(val)
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

// formatNodeLabels parses the node-labels option from the cluster create nodepool formatted string
func formatNodeLabels(nl string) map[string]string {
	data := make(map[string]string)
	labels := strings.Split(nl, "|")

	for i := range labels {
		label := strings.Split(labels[i], "=")
		data[label[0]] = label[1]
	}

	return data
}

type options struct {
	Base        *cli.Base
	CreateReq   *govultr.ClusterReq
	UpdateReq   *govultr.ClusterReqUpdate
	UpgradeReq  *govultr.ClusterUpgradeReq
	npCreateReq *govultr.NodePoolReq
	npUpdateReq *govultr.NodePoolReqUpdate
}

func (o *options) list() ([]govultr.Cluster, *govultr.Meta, error) {
	k8s, meta, _, err := o.Base.Client.Kubernetes.ListClusters(o.Base.Context, o.Base.Options)
	return k8s, meta, err
}

func (o *options) get() (*govultr.Cluster, error) {
	k8, _, err := o.Base.Client.Kubernetes.GetCluster(o.Base.Context, o.Base.Args[0])
	return k8, err
}

func (o *options) create() (*govultr.Cluster, error) {
	k8, _, err := o.Base.Client.Kubernetes.CreateCluster(o.Base.Context, o.CreateReq)
	return k8, err
}

func (o *options) update() error {
	return o.Base.Client.Kubernetes.UpdateCluster(o.Base.Context, o.Base.Args[0], o.UpdateReq)
}

func (o *options) del() error {
	return o.Base.Client.Kubernetes.DeleteCluster(o.Base.Context, o.Base.Args[0])
}

func (o *options) delWithRes() error {
	return o.Base.Client.Kubernetes.DeleteClusterWithResources(o.Base.Context, o.Base.Args[0])
}

func (o *options) config() (*govultr.KubeConfig, error) {
	kc, _, err := o.Base.Client.Kubernetes.GetKubeConfig(o.Base.Context, o.Base.Args[0])
	return kc, err
}

func (o *options) versions() (*govultr.Versions, error) {
	versions, _, err := o.Base.Client.Kubernetes.GetVersions(o.Base.Context)
	return versions, err
}

func (o *options) upgrades() ([]string, error) {
	ups, _, err := o.Base.Client.Kubernetes.GetUpgrades(o.Base.Context, o.Base.Args[0])
	return ups, err
}

func (o *options) upgrade() error {
	return o.Base.Client.Kubernetes.Upgrade(o.Base.Context, o.Base.Args[0], o.UpgradeReq)
}

func (o *options) nodePools() ([]govultr.NodePool, *govultr.Meta, error) {
	nps, meta, _, err := o.Base.Client.Kubernetes.ListNodePools(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return nps, meta, err
}

func (o *options) nodePool() (*govultr.NodePool, error) {
	np, _, err := o.Base.Client.Kubernetes.GetNodePool(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return np, err
}

func (o *options) nodePoolCreate() (*govultr.NodePool, error) {
	np, _, err := o.Base.Client.Kubernetes.CreateNodePool(o.Base.Context, o.Base.Args[0], o.npCreateReq)
	return np, err
}

func (o *options) nodePoolUpdate() (*govultr.NodePool, error) {
	np, _, err := o.Base.Client.Kubernetes.UpdateNodePool(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.npUpdateReq)
	return np, err
}

func (o *options) nodePoolDelete() error {
	return o.Base.Client.Kubernetes.DeleteNodePool(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) nodePoolNodeDelete() error {
	return o.Base.Client.Kubernetes.DeleteNodePoolInstance(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2])
}

func (o *options) nodePoolNodeRecycle() error {
	return o.Base.Client.Kubernetes.RecycleNodePoolInstance(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2])
}
