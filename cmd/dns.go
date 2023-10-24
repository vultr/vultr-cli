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
	"errors"

	"github.com/spf13/cobra"
)

// DNS represents the dns command
func DNS() *cobra.Command {
	dnsCmd := &cobra.Command{
		Use:   "dns",
		Short: "dns is used to access dns commands",
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Context().Value(ctxAuthKey{}).(bool) {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	dnsCmd.AddCommand(DNSDomain())
	dnsCmd.AddCommand(DNSRecord())
	return dnsCmd
}
