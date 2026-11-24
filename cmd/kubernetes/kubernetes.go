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
	long = `Get all available commands for Kubernetes`

	createLong    = `Create kubernetes cluster on your Vultr account`
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
	--node-pools="quantity:1,plan:vc2-4c-8gb,label:main-node-pool/quantity:5,plan:vc2-2c-4gb,label:worker-pool, \
		auto-scaler:true,min-nodes:5,max-nodes:10"

	Using node labels 
	--node-pools="quantity:5,plan:vc2-2c-4gb,label:worker-pool,auto-scaler:true,min-nodes:5,max-nodes:10, \
		node-labels:application=identity-service|worker-size=small"
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
	vultr-cli kubernetes node-pool update ffd31f18-5f77-454c-9064-212f942c3c34 abd31f18-3f77-454c-9064-212f942c3c34 \
		--quantity=4 --node-labels="application=id-service,environment=development"

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

	appsAvailableLong = `Retrieves a list of all available kubernetes cluster
applications. Further information may be retrieved with the 'details' command.`
	appsAvailableExample = `
	#Full example
	vultr-cli kubernetes apps list
`

	appsDetailsLong    = `Retrieves additional details related to the application.`
	appsDetailsExample = `
	#Full example
	vultr-cli kubernetes apps details c6a00ea5-9a54-481c-a79c-e4a1a74e8f44
`

	appsDeploymentListLong    = `Retrieves a list of deployments on a kubernetes cluster.`
	appsDeploymentListExample = `
	#Full example
	vultr-cli kubernetes apps deployments list 8cf93f5f-875e-4aea-8a84-ec2dd592b5ae
`

	appsDeploymentGetLong    = `Retrieves details of a deployment on a kubernetes cluster.`
	appsDeploymentGetExample = `
	#Full example
	vultr-cli kubernetes apps deployments get 8cf93f5f-875e-4aea-8a84-ec2dd592b5ae 4130764b-5276-4552-546f-32513239732b
`

	appsDeploymentCreateLong    = `Intiate the deployment of the app onto the kubernetes cluster.`
	appsDeploymentCreateExample = `
	#Full example
	vultr-cli kubernetes apps deployments create 8cf93f5f-875e-4aea-8a84-ec2dd592b5ae \
		--app-id 4130764b-5276-4552-546f-32513239732b
`

	appsDeploymentUpgradeLong = `Intiate the upgrade process for an app deployment on the kubernetes cluster. 
This will automatically select the latest version of the app for the upgrade.`
	appsDeploymentUpgradeExample = `
	#Full example
	vultr-cli kubernetes apps deployments upgrade 8cf93f5f-875e-4aea-8a84-ec2dd592b5ae 4130764b-5276-4552-546f-32513239732b
`

	appsDeploymentDeleteLong    = `Intiate the delete process for an app deployment.`
	appsDeploymentDeleteExample = `
	#Full example
	vultr-cli kubernetes apps deployments delete 8cf93f5f-875e-4aea-8a84-ec2dd592b5ae 4130764b-5276-4552-546f-32513239732b
`

	npLabelListExample = `
	# Full example
	vultr-cli kubernetes node-pool label list 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa
`

	npLabelGetExample = `
	# Full example
	vultr-cli kubernetes node-pool label get 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa \
		162b7e698c7fe33235649d8557c37fcc
`

	npLabelCreateExample = `
	# Full example
	vultr-cli kubernetes node-pool label create 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa \
		--key "my-key" --value "my-value"
`

	npLabelDeleteExample = `
	# Full example
	vultr-cli kubernetes node-pool label delete 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa \
		162b7e698c7fe33235649d8557c37fcc
`

	npTaintListExample = `
	# Full example
	vultr-cli kubernetes node-pool taint list 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa
`

	npTaintGetExample = `
	# Full example
	vultr-cli kubernetes node-pool taint get 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa \
		162b7e698c7fe33235649d8557c37fcc
`

	npTaintCreateExample = `
	# Full example
	vultr-cli kubernetes node-pool taint create 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa \
		--key "my-key" --value "my-value" --effect "NoSchedule"
`

	npTaintDeleteExample = `
	# Full example
	vultr-cli kubernetes node-pool taint delete 3a350f61-55d9-46a2-ae7b-7d2afc808995 84b1454f-c53a-4f5f-8f4b-fefa1a8f73fa \
		162b7e698c7fe33235649d8557c37fcc
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth() {
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
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
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
			if errFw != nil {
				return fmt.Errorf("error parsing flag 'enable-firewall' for kubernetes cluster create : %v", errFw)
			}

			vpcID, err := cmd.Flags().GetString("vpc-id")
			if err != nil {
				return fmt.Errorf("error parsing flag 'vpc-id' for kubernetes create : %v", err)
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
				VPCID:           vpcID,
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

	create.Flags().String(
		"vpc-id",
		"",
		`(optional) the VPC ID to use when initializing node private networking. A new VPC is created by default.`,
	)

	create.Flags().StringArrayP(
		"node-pools",
		"n",
		[]string{},
		`a comma-separated, key-value pair list of node pools. At least one node pool is required. 
At least one node is required in node pool. Use / between each new node pool.  
E.g: 'plan:vhf-8c-32gb,label:mynodepool,tag:my-tag,quantity:3/plan:vhf-8c-32gb,label:mynodepool2,quantity:3`,
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
				return errors.New("please provide a cluster ID")
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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
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

			summarize, errSu := cmd.Flags().GetBool("summarize")
			if errSu != nil {
				return fmt.Errorf("error parsing flag 'summarize' for kubernetes node pool list : %v", errSu)
			}

			nps, meta, err := o.nodePools()
			if err != nil {
				return fmt.Errorf("error getting node pool list : %v", err)
			}

			var data printer.ResourceOutput
			if summarize {
				data = &NodePoolsSummaryPrinter{NodePools: nps, Meta: meta}
			} else {
				data = &NodePoolsPrinter{NodePools: nps, Meta: meta}
			}

			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	npList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	npList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)
	npList.Flags().BoolP("summarize", "", false, "(optional) Summarize the list output. One line per node pool.")

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
				return fmt.Errorf(
					"error parsing flag 'auto-scaler' for kubernetes cluster node pool create : %v",
					errAu,
				)
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
				return fmt.Errorf(
					"error parsing flag 'node-labels' for kubernetes cluster node pool create : %v",
					errNl,
				)
			}

			vpcOnly, err := cmd.Flags().GetBool("vpc-only")
			if err != nil {
				return fmt.Errorf("error parsing flag 'vpc-only' for kubernetes node pool create : %v", err)
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

			if vpcOnly {
				o.npCreateReq.VPCOnly = govultr.BoolToBoolPtr(true)
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

	npCreate.Flags().IntP(
		"quantity",
		"q",
		1,
		"Number of nodes in your node pool. Note that at least one node is required for a node pool.",
	)
	if err := npCreate.MarkFlagRequired("quantity"); err != nil {
		fmt.Printf("error marking kubernetes node-pool create 'quantity' flag required: %v\n", err)
		os.Exit(1)
	}

	npCreate.Flags().BoolP("auto-scaler", "", false, "Enable the auto scaler with your cluster")
	npCreate.Flags().IntP("min-nodes", "", 1, "Minimum nodes for auto scaler")
	npCreate.Flags().IntP("max-nodes", "", 1, "Maximum nodes for auto scaler")
	npCreate.Flags().StringToString(
		"node-labels",
		nil,
		"A key=value comma separated string of labels to apply to the nodes in this node pool",
	)

	npCreate.Flags().Bool("vpc-only", false,
		`Enable VPC only network access on this node pool. This requires the VPC used by the cluster to 
have an active NAT gateway and, when enabled, will mean that the nodes will not have a public address 
or network interface. All traffic from the nodes will be routed, instead, through the NAT gateway.`,
	)

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
				return fmt.Errorf(
					"error parsing flag 'auto-scaler' for kubernetes cluster node pool update : %v",
					errAu,
				)
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
				return fmt.Errorf(
					"error parsing flag 'node-labels' for kubernetes cluster node pool update : %v",
					errNl,
				)
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
				o.npUpdateReq.Labels = npLabels //nolint
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
	npUpdate.Flags().StringToString(
		"node-labels",
		nil,
		"A key=value comma separated string of labels to apply to the nodes in this node pool",
	)

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

	// Node pool label
	npLabel := &cobra.Command{
		Use:     "label",
		Short:   "Commands to manage node pool labels",
		Aliases: []string{"labels"},
	}

	// Node pool label list
	npLabelList := &cobra.Command{
		Use:     "list <Cluster ID> <Node Pool ID>",
		Short:   "List all labels on the node pool",
		Example: npLabelListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and a node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			labels, err := o.nodePoolLabelList()
			if err != nil {
				return fmt.Errorf("error listing kubernetes node pool labels : %v", err)
			}

			data := &LabelsPrinter{Labels: labels}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Node pool label get
	npLabelGet := &cobra.Command{
		Use:     "get <Cluster ID> <Node Pool ID> <Label ID>",
		Short:   "Get a label on the node pool",
		Example: npLabelGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("please provide a cluster ID, a node pool ID and a label ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, err := o.nodePoolLabelGet()
			if err != nil {
				return fmt.Errorf("error getting kubernetes node pool label : %v", err)
			}

			data := &LabelPrinter{Label: label}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Node pool label create
	npLabelCreate := &cobra.Command{
		Use:     "create <Cluster ID> <Node Pool ID>",
		Short:   "Create a label in a node pool",
		Example: npLabelCreateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and a node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := cmd.Flags().GetString("key")
			if err != nil {
				return fmt.Errorf("error parsing flag 'key' for kubernetes node pool label create : %v", err)
			}

			val, err := cmd.Flags().GetString("value")
			if err != nil {
				return fmt.Errorf("error parsing flag 'value' for kubernetes node pool label create : %v", err)
			}

			o.npLabelReq = &govultr.NodePoolLabelReq{
				Key:   key,
				Value: val,
			}

			label, err := o.nodePoolLabelCreate()
			if err != nil {
				return fmt.Errorf("error creating kubernetes node pool label : %v", err)
			}

			data := &LabelPrinter{Label: label}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	npLabelCreate.Flags().String("key", "", "The key for the label")
	if err := npLabelCreate.MarkFlagRequired("key"); err != nil {
		fmt.Printf("error marking kubernetes node pool label create 'key' flag required: %v", err)
		os.Exit(1)
	}

	npLabelCreate.Flags().String("value", "", "The value for the label")
	if err := npLabelCreate.MarkFlagRequired("value"); err != nil {
		fmt.Printf("error marking kubernetes node pool label create 'value' flag required: %v", err)
		os.Exit(1)
	}

	// Node pool label delete
	npLabelDelete := &cobra.Command{
		Use:     "delete <Cluster ID> <Node Pool ID> <Label ID>",
		Short:   "Delete a node pool label",
		Example: npLabelDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("please provide a cluster ID, a node pool ID and a label ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.nodePoolLabelDelete(); err != nil {
				return fmt.Errorf("error deleting kubernetes node pool label : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes node pool label has been deleted"), nil)

			return nil
		},
	}

	npLabel.AddCommand(
		npLabelList,
		npLabelGet,
		npLabelCreate,
		npLabelDelete,
	)

	// Node pool taint
	npTaint := &cobra.Command{
		Use:     "taint",
		Short:   "Commands to manage node pool taints",
		Aliases: []string{"taints"},
	}

	// Node pool taint list
	npTaintList := &cobra.Command{
		Use:     "list <Cluster ID> <Node Pool ID>",
		Short:   "List all taints on the node pool",
		Example: npTaintListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and a node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			taints, err := o.nodePoolTaintList()
			if err != nil {
				return fmt.Errorf("error listing kubernetes node pool taints : %v", err)
			}

			data := &TaintsPrinter{Taints: taints}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Node pool taint get
	npTaintGet := &cobra.Command{
		Use:     "get <Cluster ID> <Node Pool ID> <Taint ID>",
		Short:   "Get a taint on the node pool",
		Example: npTaintGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("please provide a cluster ID, a node pool ID and a taint ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			taint, err := o.nodePoolTaintGet()
			if err != nil {
				return fmt.Errorf("error getting kubernetes node pool taint : %v", err)
			}

			data := &TaintPrinter{Taint: taint}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Node pool taint create
	npTaintCreate := &cobra.Command{
		Use:     "create <Cluster ID> <Node Pool ID>",
		Short:   "Create a taint in a node pool",
		Example: npTaintCreateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and a node pool ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := cmd.Flags().GetString("key")
			if err != nil {
				return fmt.Errorf("error parsing flag 'key' for kubernetes node pool taint create : %v", err)
			}

			val, err := cmd.Flags().GetString("value")
			if err != nil {
				return fmt.Errorf("error parsing flag 'value' for kubernetes node pool taint create : %v", err)
			}

			effect, err := cmd.Flags().GetString("effect")
			if err != nil {
				return fmt.Errorf("error parsing flag 'effect' for kubernetes node pool taint create : %v", err)
			}

			o.npTaintReq = &govultr.NodePoolTaintReq{
				Key:    key,
				Value:  val,
				Effect: effect,
			}

			taint, err := o.nodePoolTaintCreate()
			if err != nil {
				return fmt.Errorf("error creating kubernetes node pool taint : %v", err)
			}

			data := &TaintPrinter{Taint: taint}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	npTaintCreate.Flags().String("key", "", "The key for the taint")
	if err := npTaintCreate.MarkFlagRequired("key"); err != nil {
		fmt.Printf("error marking kubernetes node pool taint create 'key' flag required: %v", err)
		os.Exit(1)
	}

	npTaintCreate.Flags().String("value", "", "The value for the taint")
	if err := npTaintCreate.MarkFlagRequired("value"); err != nil {
		fmt.Printf("error marking kubernetes node pool taint create 'value' flag required: %v", err)
		os.Exit(1)
	}

	npTaintCreate.Flags().String(
		"effect",
		"",
		"The effect of the taint. One of 'NoSchedule', 'PreferNoSchedule' and 'NoExecute'",
	)
	if err := npTaintCreate.MarkFlagRequired("effect"); err != nil {
		fmt.Printf("error marking kubernetes node pool taint create 'effect' flag required: %v", err)
		os.Exit(1)
	}

	// Node pool taint delete
	npTaintDelete := &cobra.Command{
		Use:     "delete <Cluster ID> <Node Pool ID> <Taint ID>",
		Short:   "Delete a node pool taint",
		Example: npTaintDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("please provide a cluster ID, a node pool ID and a taint ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.nodePoolTaintDelete(); err != nil {
				return fmt.Errorf("error deleting kubernetes node pool taint : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Kubernetes node pool taint has been deleted"), nil)

			return nil
		},
	}

	npTaint.AddCommand(
		npTaintList,
		npTaintGet,
		npTaintCreate,
		npTaintDelete,
	)

	nodepool.AddCommand(
		npList,
		npGet,
		npCreate,
		npUpdate,
		npDelete,
		node,
		npLabel,
		npTaint,
	)

	// Apps
	apps := &cobra.Command{
		Use:     "apps",
		Short:   "Commands to manage apps in a kubernetes cluster",
		Aliases: []string{"app"},
	}

	// Apps list available
	appsList := &cobra.Command{
		Use:     "list",
		Short:   "List catalog of available apps for deployment",
		Aliases: []string{"available"},
		Long:    appsAvailableLong,
		Example: appsAvailableExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			apps, err := o.appsListAvailable()
			if err != nil {
				return fmt.Errorf("error retrieving available kubernetes applications list : %v", err)
			}

			data := &AppsPrinter{Apps: apps}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Apps show available details
	appsDetails := &cobra.Command{
		Use:     "details <APP ID>",
		Short:   "Show detailed information for an app",
		Long:    appsDetailsLong,
		Example: appsDetailsExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide an app ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := o.appsGetAvailable()
			if err != nil {
				return fmt.Errorf(
					"error retrieving available kubernetes application details : %v",
					err,
				)
			}

			data := &AppPrinter{App: app}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	appsDeployment := &cobra.Command{
		Use:     "deployment",
		Short:   "Commands for managing kubernetes app deployments",
		Aliases: []string{"deployments"},
	}

	// Apps deployments list
	appsDeploymentList := &cobra.Command{
		Use:     "list <Cluster ID>",
		Short:   "List all app deployments on the cluster",
		Long:    appsDeploymentListLong,
		Example: appsDeploymentListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := o.appsListDeployments()
			if err != nil {
				return fmt.Errorf("error retrieving list of application deployments on kubernetes cluster : %v", err)
			}

			data := &AppDeploymentsPrinter{Deployments: deps}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Apps get deployment
	appsDeploymentGet := &cobra.Command{
		Use:     "get <Cluster ID> <Deployment ID>",
		Short:   "Get the details of an app deployment on the cluster",
		Long:    appsDeploymentGetLong,
		Example: appsDeploymentGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and deployment ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dep, err := o.appsGetDeployment()
			if err != nil {
				return fmt.Errorf("error retrieving application deployment on kubernetes cluster : %v", err)
			}

			data := &AppDeploymentPrinter{Deployment: dep}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Apps create deployment
	appsDeploymentCreate := &cobra.Command{
		Use:     "create <Cluster ID>",
		Short:   "Create a deployment in the cluster",
		Long:    appsDeploymentCreateLong,
		Example: appsDeploymentCreateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a cluster ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appID, err := cmd.Flags().GetString("app-id")
			if err != nil {
				return fmt.Errorf("error parsing flag 'app-id' for kubernetes cluster app deployment create : %v", err)
			}

			config, err := cmd.Flags().GetString("config")
			if err != nil {
				return fmt.Errorf("error parsing flag 'config' for kubernetes cluster app deployment create : %v", err)
			}

			timeout, err := cmd.Flags().GetString("timeout")
			if err != nil {
				return fmt.Errorf("error parsing flag 'timeout' for kubernetes cluster app deployment create : %v", err)
			}

			o.appDeploymentReq = &govultr.ClusterApplicationDeploymentReq{
				ID:      appID,
				Config:  config,
				Timeout: timeout,
			}

			dep, err := o.appsCreateDeployment()
			if err != nil {
				return fmt.Errorf("error creating application deployment on kubernetes cluster : %v", err)
			}

			data := &AppDeploymentPrinter{Deployment: dep}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	appsDeploymentCreate.Flags().String("app-id", "", "The app ID to be deployed onto the kubernetes cluster")
	if err := appsDeploymentCreate.MarkFlagRequired("app-id"); err != nil {
		fmt.Printf("error marking app-id flag required for kubernetes app deployment create : %v\n", err)
		os.Exit(1)
	}

	appsDeploymentCreate.Flags().String("timeout", "", "The timeout to be used when creating the deployment")
	appsDeploymentCreate.Flags().String("config", "", "Custom Helm values to send with the deployment")

	// Apps upgrade deployment
	appsDeploymentUpgrade := &cobra.Command{
		Use:     "upgrade <Cluster ID> <Deployment ID>",
		Short:   "Upgrade a deployment in the cluster to the latest available version",
		Long:    appsDeploymentUpgradeLong,
		Example: appsDeploymentUpgradeExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and a deployment ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.appsUpgradeDeployment(); err != nil {
				return fmt.Errorf("error upgrading application deployment on kubernetes cluster : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Application deployment upgrade is underway"), nil)

			return nil
		},
	}

	// Apps delete deployment
	appsDeploymentDelete := &cobra.Command{
		Use:     "delete <Cluster ID> <Deployment ID>",
		Short:   "Delete a deployment in the cluster",
		Long:    appsDeploymentDeleteLong,
		Example: appsDeploymentDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a cluster ID and a deployment ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.appsDeleteDeployment(); err != nil {
				return fmt.Errorf("error deleting the application deployment on kubernetes cluster : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Application deployment deletion is underway"), nil)

			return nil
		},
	}

	appsDeployment.AddCommand(
		appsDeploymentList,
		appsDeploymentGet,
		appsDeploymentCreate,
		appsDeploymentUpgrade,
		appsDeploymentDelete,
	)

	apps.AddCommand(
		appsList,
		appsDetails,
		appsDeployment,
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
		apps,
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

		switch field {
		case "plan":
			nodeData.Plan = val
		case "quantity":
			port, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool quantity: %v", err)
			}
			nodeData.NodeQuantity = port
		case "label":
			nodeData.Label = val
		case "tag":
			nodeData.Tag = val
		case "node-labels":
			nodeData.Labels = formatNodeLabels(val)
		case "auto-scaler":
			v, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool auto-scaler: %v", err)
			}
			nodeData.AutoScaler = govultr.BoolToBoolPtr(v)
		case "min-nodes":
			v, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool min-nodes: %v", err)
			}
			nodeData.MinNodes = v
		case "max-nodes":
			v, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for max-nodes: %v", err)
			}
			nodeData.MaxNodes = v
		case "vpc-only":
			v, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for node pool vpc-only : %v", err)
			}
			nodeData.VPCOnly = govultr.BoolToBoolPtr(v)
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
	Base             *cli.Base
	CreateReq        *govultr.ClusterReq
	UpdateReq        *govultr.ClusterReqUpdate
	UpgradeReq       *govultr.ClusterUpgradeReq
	npCreateReq      *govultr.NodePoolReq
	npUpdateReq      *govultr.NodePoolReqUpdate
	npLabelReq       *govultr.NodePoolLabelReq
	npTaintReq       *govultr.NodePoolTaintReq
	appDeploymentReq *govultr.ClusterApplicationDeploymentReq
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
	np, _, err := o.Base.Client.Kubernetes.UpdateNodePool(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.npUpdateReq) //nolint:lll
	return np, err
}

func (o *options) nodePoolDelete() error {
	return o.Base.Client.Kubernetes.DeleteNodePool(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) nodePoolNodeDelete() error {
	return o.Base.Client.Kubernetes.DeleteNodePoolInstance(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2]) //nolint:lll
}

func (o *options) nodePoolNodeRecycle() error {
	return o.Base.Client.Kubernetes.RecycleNodePoolInstance(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2]) //nolint:lll
}

func (o *options) nodePoolLabelList() ([]govultr.NodePoolLabel, error) {
	labels, _, err := o.Base.Client.Kubernetes.ListNodePoolLabels(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return labels, err
}

func (o *options) nodePoolLabelGet() (*govultr.NodePoolLabel, error) {
	label, _, err := o.Base.Client.Kubernetes.GetNodePoolLabel(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2]) //nolint:lll
	return label, err
}

func (o *options) nodePoolLabelCreate() (*govultr.NodePoolLabel, error) {
	label, _, err := o.Base.Client.Kubernetes.CreateNodePoolLabel(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.npLabelReq) //nolint:lll
	return label, err
}

func (o *options) nodePoolLabelDelete() error {
	return o.Base.Client.Kubernetes.DeleteNodePoolLabel(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2]) //nolint:lll
}

func (o *options) nodePoolTaintList() ([]govultr.NodePoolTaint, error) {
	taints, _, err := o.Base.Client.Kubernetes.ListNodePoolTaints(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return taints, err
}

func (o *options) nodePoolTaintGet() (*govultr.NodePoolTaint, error) {
	taint, _, err := o.Base.Client.Kubernetes.GetNodePoolTaint(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2]) //nolint:lll
	return taint, err
}

func (o *options) nodePoolTaintCreate() (*govultr.NodePoolTaint, error) {
	taint, _, err := o.Base.Client.Kubernetes.CreateNodePoolTaint(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.npTaintReq) //nolint:lll
	return taint, err
}

func (o *options) nodePoolTaintDelete() error {
	return o.Base.Client.Kubernetes.DeleteNodePoolTaint(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.Base.Args[2]) //nolint:lll
}

func (o *options) appsListAvailable() ([]govultr.ClusterApplication, error) {
	apps, _, err := o.Base.Client.Kubernetes.ListClusterApplicationsAvailable(o.Base.Context)
	return apps, err
}

func (o *options) appsListDeployments() ([]govultr.ClusterApplicationDeployment, error) {
	deps, _, err := o.Base.Client.Kubernetes.ListClusterApplicationDeployments(o.Base.Context, o.Base.Args[0])
	return deps, err
}

func (o *options) appsGetDeployment() (*govultr.ClusterApplicationDeployment, error) {
	dep, _, err := o.Base.Client.Kubernetes.GetClusterApplicationDeployment(o.Base.Context, o.Base.Args[0], o.Base.Args[1]) //nolint:lll
	return dep, err
}

func (o *options) appsCreateDeployment() (*govultr.ClusterApplicationDeployment, error) {
	dep, _, err := o.Base.Client.Kubernetes.CreateClusterApplicationDeployment(o.Base.Context, o.Base.Args[0], o.appDeploymentReq) //nolint:lll
	return dep, err
}

func (o *options) appsUpgradeDeployment() error {
	return o.Base.Client.Kubernetes.UpgradeClusterApplicationDeployment(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) appsDeleteDeployment() error {
	return o.Base.Client.Kubernetes.DeleteClusterApplicationDeployment(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) appsGetAvailable() (*govultr.ClusterApplication, error) {
	apps, _, err := o.Base.Client.Kubernetes.ListClusterApplicationsAvailable(o.Base.Context)
	if err != nil {
		return nil, err
	}

	for i := range apps {
		if apps[i].ID == o.Base.Args[0] {
			return &apps[i], nil
		}
	}

	return nil, fmt.Errorf("app not found")
}
