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

// Package cmd implements the command line commands relevant to the vultr-cli
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/account"
	"github.com/vultr/vultr-cli/v3/cmd/applications"
	"github.com/vultr/vultr-cli/v3/cmd/backups"
	baremetal "github.com/vultr/vultr-cli/v3/cmd/bareMetal"
	operatingsystems "github.com/vultr/vultr-cli/v3/cmd/operatingSystems"
	"github.com/vultr/vultr-cli/v3/cmd/plans"
	"github.com/vultr/vultr-cli/v3/cmd/regions"
	"github.com/vultr/vultr-cli/v3/cmd/sshkeys"
	"github.com/vultr/vultr-cli/v3/cmd/users"
	"github.com/vultr/vultr-cli/v3/cmd/version"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

const (
	userAgent          = "vultr-cli/" + version.Version
	perPageDefault int = 100
)

var (
	cfgFile string
	output  string
	client  *govultr.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vultr-cli",
	Short: "vultr-cli is a command line interface for the Vultr API",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// init the config file with viper
	initConfig()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configHome(), "config file (default is $HOME/.vultr-cli.yaml)")
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

	rootCmd.AddCommand(account.NewCmdAccount(base))
	rootCmd.AddCommand(applications.NewCmdApplications(base))
	rootCmd.AddCommand(backups.NewCmdBackups(base))
	rootCmd.AddCommand(baremetal.NewCmdBareMetal(base))
	rootCmd.AddCommand(Billing())
	rootCmd.AddCommand(BlockStorageCmd())
	rootCmd.AddCommand(ContainerRegistry())
	rootCmd.AddCommand(Database())
	rootCmd.AddCommand(DNS())
	rootCmd.AddCommand(Firewall())
	rootCmd.AddCommand(ISO())
	rootCmd.AddCommand(Kubernetes())
	rootCmd.AddCommand(LoadBalancer())
	rootCmd.AddCommand(Network())
	rootCmd.AddCommand(operatingsystems.NewCmdOS(base))
	rootCmd.AddCommand(ObjectStorageCmd())
	rootCmd.AddCommand(plans.NewCmdPlan(base))
	rootCmd.AddCommand(regions.NewCmdRegion(base))
	rootCmd.AddCommand(ReservedIP())
	rootCmd.AddCommand(Script())
	rootCmd.AddCommand(Instance())
	rootCmd.AddCommand(Snapshot())
	rootCmd.AddCommand(sshkeys.NewCmdSSHKey(base))
	rootCmd.AddCommand(users.NewCmdUser(base))
	rootCmd.AddCommand(version.NewCmdVersion())
	rootCmd.AddCommand(VPC())
	rootCmd.AddCommand(VPC2())
}

// initConfig reads in config file and ENV variables if set.
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

func getPaging(cmd *cobra.Command) *govultr.ListOptions {
	options := &govultr.ListOptions{}

	cursor, _ := cmd.Flags().GetString("cursor")
	perPage, _ := cmd.Flags().GetInt("per-page")

	if cursor != "" {
		options.Cursor = cursor
	}

	if perPage != 0 {
		options.PerPage = perPage
	}

	return options
}

func configHome() string {
	// check for a config file at ~/.config/vultr-cli.yaml
	configFolder, errConfig := os.UserConfigDir()
	if errConfig != nil {
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
		os.Exit(1)
	}

	configFile = fmt.Sprintf("%s/.vultr-cli.yaml", configFolder)
	if _, err := os.Stat(configFile); err != nil {
		// if it doesn't exist, create one
		f, err := os.Create(filepath.Clean(configFile))
		if err != nil {
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
