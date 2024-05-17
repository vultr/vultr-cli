// Package version provides functionality to display the CLI version
package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

const (
	Version string = "v3.1.0"
)

var (
	long = `Displays current version of the Vultr-CLI`

	example = `
	# example
	vultr-cli version
	
	# Shortened with alias commands
	vultr-cli v
	`
)

// NewCmdVersion returns cobra command for version
func NewCmdVersion(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Display the vultr-cli version",
		Aliases: []string{"v"},
		Long:    long,
		Example: example,
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Printer.Display(&VersionPrinter{Version: o.get()}, nil)
		},
	}

	return cmd
}

type options struct {
	Base    *cli.Base
	Version string
}

func (o *options) get() string {
	return fmt.Sprintf("Vultr-CLI %s", Version)
}
