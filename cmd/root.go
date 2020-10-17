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
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

var cfgFile string
var client *govultr.Client

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
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initClient)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vultr-cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(Applications())
	rootCmd.AddCommand(Backups())
	rootCmd.AddCommand(BareMetal())
	rootCmd.AddCommand(BlockStorageCmd())
	rootCmd.AddCommand(Dns())
	rootCmd.AddCommand(Firewall())
	rootCmd.AddCommand(Iso())
	rootCmd.AddCommand(LoadBalancer())
	rootCmd.AddCommand(Network())
	rootCmd.AddCommand(Os())
	rootCmd.AddCommand(ObjectStorageCmd())
	rootCmd.AddCommand(Plans())
	rootCmd.AddCommand(Regions())
	rootCmd.AddCommand(ReservedIP())
	rootCmd.AddCommand(Script())
	rootCmd.AddCommand(Server())
	rootCmd.AddCommand(Snapshot())
	rootCmd.AddCommand(SSHKey())
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(User())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vultr-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vultr-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initClient() {
	apiKey := os.Getenv("VULTR_API_KEY")
	if apiKey == "" {
		fmt.Println("Please export your VULTR API key as an environment variable, eg:")
		fmt.Println("export VULTR_API_KEY='<api_key_from_vultr_account>'")
		os.Exit(1)
	}

	config := &oauth2.Config{}
	ts := config.TokenSource(context.Background(), &oauth2.Token{AccessToken: apiKey})
	client = govultr.NewClient(oauth2.NewClient(context.Background(), ts))

	client.SetRateLimit(1 * time.Second)
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
