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

	"github.com/spf13/cobra"
)

// Firewall represents the firewall command
func Firewall() *cobra.Command {
	firewallCmd := &cobra.Command{
		Use:     "firewall",
		Short:   "firewall is used to access firewall commands",
		Long:    ``,
		Aliases: []string{"fw"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if auth := cmd.Context().Value("authenticated"); auth != true {
				return fmt.Errorf(apiKeyError)
			}
			return nil
		},
	}

	firewallCmd.AddCommand(FirewallGroup(), FirewallRule())

	return firewallCmd
}
