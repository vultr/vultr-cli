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

package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

const (
	Version string = "v2.21.0"
)

// Interface for version
type Interface interface {
	Get() string
}

// Options for version
type Options struct {
	Version string
	Printer *printer.Output
}

var (
	long = `Displays current version of the Vultr-CLI`

	example = `
	# example
	vultr-cli version
	
	# Shortened with alias commands
	vultr-cli v
	`
)

// NewVersionOptions returns a VersionOptions struct
func NewVersionOptions() *Options {
	return &Options{Printer: &printer.Output{}}
}

// NewCmdVersion returns cobra command for version
func NewCmdVersion() *cobra.Command {
	v := NewVersionOptions()
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "vultr-cli version",
		Long:    long,
		Example: example,
		Run: func(cmd *cobra.Command, args []string) {
			v.Printer.Output = viper.GetString("output")
			v.Printer.Display(&VersionPrinter{Version: v.Get()}, nil)
		},
	}

	return cmd
}

// Get the version for vultr-cli
func (v *Options) Get() string {
	return fmt.Sprintf("Vultr-CLI %s", Version)
}
