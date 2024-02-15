package utils

import (
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
)

func GetPaging(cmd *cobra.Command) *govultr.ListOptions {
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
