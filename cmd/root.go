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
	"fmt"
	"os"

	"github.com/vultr/vultr-cli/cmd/operatingSystems"
	"github.com/vultr/vultr-cli/cmd/sshkeys"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/cmd/applications"
	"github.com/vultr/vultr-cli/cmd/plans"
	"github.com/vultr/vultr-cli/cmd/regions"
	"github.com/vultr/vultr-cli/cmd/users"
	"github.com/vultr/vultr-cli/cmd/version"
	"github.com/vultr/vultr-cli/pkg/cli"
)

var (
	cfgFile string
	output  string
	base    *cli.Base
	client  *govultr.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vultr-cli",
	Short:   "vultr-cli is a command line interface for the Vultr API",
	Long:    ``,
	Aliases: []string{"vultrctl"},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	initConfig()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configHome(), "config file (default is $HOME/.vultr-cli.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringVar(&output, "output", "text", "out of data json | yaml | text. text is default")
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(applications.NewCmdApplications(base))
	rootCmd.AddCommand(Backups())
	rootCmd.AddCommand(BareMetal())
	rootCmd.AddCommand(BlockStorageCmd())
	rootCmd.AddCommand(DNS())
	rootCmd.AddCommand(Firewall())
	rootCmd.AddCommand(ISO())
	rootCmd.AddCommand(LoadBalancer())
	rootCmd.AddCommand(Network())
	rootCmd.AddCommand(operatingSystems.NewCmdOS(base))
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
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var token string
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

	token = viper.GetString("api-key")
	if token == "" {
		token = os.Getenv("VULTR_API_KEY")
	}

	if token == "" {
		fmt.Println("Please export your VULTR API key as an environment variable or add `api-key` to your config file, eg:")
		fmt.Println("export VULTR_API_KEY='<api_key_from_vultr_account>'")
		os.Exit(1)
	}

	base = cli.NewCLIBase(token, viper.GetString("output"))
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
	configHome, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	configHome = fmt.Sprintf("%s/.vultr-cli.yaml", configHome)
	if _, err := os.Stat(configHome); err != nil {
		f, err := os.Create(configHome)
		if err != nil {
			os.Exit(1)
		}
		defer f.Close()

	}

	return configHome
}
