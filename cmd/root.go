// Package cmd implements the command line commands relevant to the vultr-cli
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/vultr-cli/v3/cmd/account"
	"github.com/vultr/vultr-cli/v3/cmd/applications"
	"github.com/vultr/vultr-cli/v3/cmd/backups"
	"github.com/vultr/vultr-cli/v3/cmd/baremetal"
	"github.com/vultr/vultr-cli/v3/cmd/billing"
	"github.com/vultr/vultr-cli/v3/cmd/blockstorage"
	"github.com/vultr/vultr-cli/v3/cmd/containerregistry"
	"github.com/vultr/vultr-cli/v3/cmd/database"
	"github.com/vultr/vultr-cli/v3/cmd/dns"
	"github.com/vultr/vultr-cli/v3/cmd/firewall"
	"github.com/vultr/vultr-cli/v3/cmd/inference"
	"github.com/vultr/vultr-cli/v3/cmd/instance"
	"github.com/vultr/vultr-cli/v3/cmd/iso"
	"github.com/vultr/vultr-cli/v3/cmd/kubernetes"
	"github.com/vultr/vultr-cli/v3/cmd/loadbalancer"
	"github.com/vultr/vultr-cli/v3/cmd/marketplace"
	"github.com/vultr/vultr-cli/v3/cmd/objectstorage"
	"github.com/vultr/vultr-cli/v3/cmd/operatingsystems"
	"github.com/vultr/vultr-cli/v3/cmd/plans"
	"github.com/vultr/vultr-cli/v3/cmd/regions"
	"github.com/vultr/vultr-cli/v3/cmd/reservedip"
	"github.com/vultr/vultr-cli/v3/cmd/script"
	"github.com/vultr/vultr-cli/v3/cmd/snapshot"
	"github.com/vultr/vultr-cli/v3/cmd/sshkeys"
	"github.com/vultr/vultr-cli/v3/cmd/users"
	"github.com/vultr/vultr-cli/v3/cmd/version"
	"github.com/vultr/vultr-cli/v3/cmd/vpc"
	"github.com/vultr/vultr-cli/v3/cmd/vpc2"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

const (
	userAgent          = "vultr-cli/" + version.Version
	perPageDefault int = 100
)

var (
	cfgFile string
	output  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "vultr-cli",
	Short:        "vultr-cli is a command line interface for the Vultr API",
	Long:         ``,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	configPath := configHome()

	// init the config file with viper
	initConfig()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configPath, "config file (default is $HOME/.vultr-cli.yaml)")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Printf("error binding root pflag 'config': %v\n", err)
	}

	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "output format [ text | json | yaml ]")
	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")); err != nil {
		fmt.Printf("error binding root pflag 'output': %v\n", err)
	}

	base := cli.NewCLIBase(
		os.Getenv("VULTR_API_KEY"),
		userAgent,
		output,
	)

	rootCmd.AddCommand(
		account.NewCmdAccount(base),
		applications.NewCmdApplications(base),
		backups.NewCmdBackups(base),
		baremetal.NewCmdBareMetal(base),
		billing.NewCmdBilling(base),
		blockstorage.NewCmdBlockStorage(base),
		containerregistry.NewCmdContainerRegistry(base),
		database.NewCmdDatabase(base),
		dns.NewCmdDNS(base),
		firewall.NewCmdFirewall(base),
		inference.NewCmdInference(base),
		iso.NewCmdISO(base),
		kubernetes.NewCmdKubernetes(base),
		loadbalancer.NewCmdLoadBalancer(base),
		marketplace.NewCmdMarketplace(base),
		operatingsystems.NewCmdOS(base),
		objectstorage.NewCmdObjectStorage(base),
		plans.NewCmdPlan(base),
		regions.NewCmdRegion(base),
		reservedip.NewCmdReservedIP(base),
		script.NewCmdScript(base),
		instance.NewCmdInstance(base),
		snapshot.NewCmdSnapshot(base),
		sshkeys.NewCmdSSHKey(base),
		users.NewCmdUser(base),
		version.NewCmdVersion(base),
		vpc.NewCmdVPC(base),
		vpc2.NewCmdVPC2(base),
	)
}

// initConfig reads in config file to viper if it exists
func initConfig() {
	configPath := viper.GetString("config")

	if configPath == "" {
		cfgDir, err := os.UserHomeDir()
		if err != nil {
			os.Exit(1)
		}
		configPath = fmt.Sprintf("%s/.vultr-cli.yaml", cfgDir)
	}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error Reading in file:", viper.ConfigFileUsed())
	}
}

func configHome() string {
	// check for a config file in the user config directory
	configFolder, errConfig := os.UserConfigDir()
	if errConfig != nil {
		fmt.Printf("Unable to determine default user config directory : %v", errConfig)
		os.Exit(1)
	}

	configFile := fmt.Sprintf("%s/vultr-cli.yaml", configFolder)
	if _, err := os.Stat(configFile); err == nil {
		// if one exists, return the path
		return configFile
	}

	// check for a config file at ~/.vultr-cli.yaml
	configFolder, errHome := os.UserHomeDir()
	if errHome != nil {
		fmt.Printf("Unable to check user config in home directory: %v", errHome)
		os.Exit(1)
	}

	configFile = fmt.Sprintf("%s/.vultr-cli.yaml", configFolder)
	if _, err := os.Stat(configFile); err != nil {
		// if it doesn't exist, create one
		f, err := os.Create(filepath.Clean(configFile))
		if err != nil {
			fmt.Printf("Unable to create default config file : %v", err)
			os.Exit(1)
		}

		defer func() {
			if errCls := f.Close(); errCls != nil {
				fmt.Printf("failed to close config file.. error: %v", errCls)
			}
		}()
	}

	return configFile
}
