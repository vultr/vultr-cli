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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

const (
	userAgent          = "vultr-cli/" + version
	perPageDefault int = 100
	//nolint: gosec
	apiKeyError string = `
Please export your VULTR API key as an environment variable or add 'api-key' to your config file, eg:
export VULTR_API_KEY='<api_key_from_vultr_account>'
	`
)

// ctxAuthKey is the context key for the authorized token check
type ctxAuthKey struct{}

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
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configHome(), "config file (default is $HOME/.vultr-cli.yaml)")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Printf("error binding root pflag 'config': %v\n", err)
	}
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(Applications())
	rootCmd.AddCommand(Backups())
	rootCmd.AddCommand(BareMetal())
	rootCmd.AddCommand(Billing())
	rootCmd.AddCommand(BlockStorageCmd())
	rootCmd.AddCommand(ContainerRegistry())
	rootCmd.AddCommand(Database())
	rootCmd.AddCommand(DNS())
	rootCmd.AddCommand(Firewall())
	rootCmd.AddCommand(ISO())
	rootCmd.AddCommand(Kubernetes())
	rootCmd.AddCommand(LoadBalancer())
	rootCmd.AddCommand(Marketplace())
	rootCmd.AddCommand(Network())
	rootCmd.AddCommand(Os())
	rootCmd.AddCommand(ObjectStorageCmd())
	rootCmd.AddCommand(Plans())
	rootCmd.AddCommand(Regions())
	rootCmd.AddCommand(ReservedIP())
	rootCmd.AddCommand(Script())
	rootCmd.AddCommand(Instance())
	rootCmd.AddCommand(Snapshot())
	rootCmd.AddCommand(SSHKey())
	rootCmd.AddCommand(User())
	rootCmd.AddCommand(VPC())
	rootCmd.AddCommand(VPC2())

	ctx := initConfig()
	rootCmd.SetContext(ctx)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() context.Context {
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

	ctx := context.Background()

	if token == "" {
		client = govultr.NewClient(nil)
		ctx = context.WithValue(ctx, ctxAuthKey{}, false)
	} else {
		config := &oauth2.Config{}
		ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: token})
		client = govultr.NewClient(oauth2.NewClient(ctx, ts))
		ctx = context.WithValue(ctx, ctxAuthKey{}, true)
	}

	client.SetRateLimit(1 * time.Second)
	client.SetUserAgent(userAgent)

	return ctx
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
