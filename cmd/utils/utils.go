// Package utils provides some common utility functions
package utils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

const (
	PerPageDefault   int = 100
	DecimalPrecision int = 4
)

// SetOptions initializes values used in all CLI commands
func SetOptions(b *cli.Base, cmd *cobra.Command, args []string) {
	b.Args = args
	b.Printer.Output = viper.GetString("output")
}

// GetFirewallSource parses the source and if empty, returns 'anywhere'
func GetFirewallSource(source string) string {
	if source == "" {
		return "anywhere"
	}
	return source
}

// FormatFirewalNetwork returns the subnet and size of a network in CIDR
// notation
func FormatFirewallNetwork(subnet string, size int) string {
	return fmt.Sprintf("%s/%d", subnet, size)
}
