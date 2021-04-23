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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/vultr-cli/cmd/printer"
)

type Interface interface {
	Get() string
}

type Options struct {
	Version string
	Printer *printer.Output
}

func newVersionOptions() *Options {
	return &Options{Printer: &printer.Output{}}
}

func NewCmdVersion() *cobra.Command {
	v := newVersionOptions()
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display current version of Vultr-cli",
		Long:  ``,
		Example: ``,
		Run: func(cmd *cobra.Command, args []string) {
			v.Printer.Output = viper.GetString("output")
			v.Printer.Display(&printer.Version{Version: v.Get()}, nil)
		},
	}

	return cmd
}

func (v *Options) Get() string {
	return "Vultr-cli v2.4.1"
}
