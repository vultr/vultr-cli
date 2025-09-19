// Package logs provides the functionality for the CLI to access logs from the API
package logs

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	logsLong    = `All commands related to the logs functionality`
	logsExample = `
	vultr-cli logs
	`

	listLong    = `Retrieve all logs between the start and end timestamps, filtered by UUID, type or log level`
	listExample = `
	# Full example
	vultr-cli logs list --start '2025-08-26T00:00:00Z' --end '2025-09-13T00:30:00Z' \
		--uuid '8b903420-b2e3-4e4f-9f88-19efb30e1237' --type 'instances'

	# Shortened with aliased commands
	vultr-cli logs list -s '2025-08-26T00:00:00Z' -e '2025-09-13T00:30:00Z' \
		-u '8b903420-b2e3-4e4f-9f88-19efb30e1237' -t 'instances'
	`
)

// NewCmdLogs provides the logs command to the CLI
func NewCmdLogs(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "logs",
		Short:   "Commands to view logs",
		Aliases: []string{"log"},
		Long:    logsLong,
		Example: logsExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}

			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Short:   "Retrieve logs",
		Aliases: []string{"l", "get"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			start, err := cmd.Flags().GetString("start")
			if err != nil {
				return fmt.Errorf("error parsing flag 'start' for logs list : %v", err)
			}

			end, err := cmd.Flags().GetString("end")
			if err != nil {
				return fmt.Errorf("error parsing flag 'end' for logs list : %v", err)
			}

			level, err := cmd.Flags().GetString("level")
			if err != nil {
				return fmt.Errorf("error parsing flag 'level' for logs list : %v", err)
			}

			resType, err := cmd.Flags().GetString("type")
			if err != nil {
				return fmt.Errorf("error parsing flag 'type' for logs list : %v", err)
			}

			uuid, err := cmd.Flags().GetString("uuid")
			if err != nil {
				return fmt.Errorf("error parsing flag 'uuid' for logs list : %v", err)
			}

			o.LogsOptions = govultr.LogsOptions{
				StartTime:    start,
				EndTime:      end,
				LogLevel:     level,
				ResourceType: resType,
				ResourceID:   uuid,
			}

			logs, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving logs list : %v", err)
			}

			data := &LogsPrinter{Logs: logs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}
	list.Flags().StringP("start", "s", "", `
a UTC timestamp for the start of the time period from which to return logs (ex. 2025-06-26T00:00:00Z)
`)
	if err := list.MarkFlagRequired("start"); err != nil {
		fmt.Printf("error marking logs list 'start' flag required: %v", err)
		os.Exit(1)
	}

	list.Flags().StringP("end", "e", "", `
a UTC timestamp for the end of the time period from which to return logs (ex. 2025-06-28T00:00:00Z). 
Must be later than the start timestamp and may not exceed 30 days and 1 hour difference from start.
`)
	if err := list.MarkFlagRequired("end"); err != nil {
		fmt.Printf("error marking logs list 'end' flag required: %v", err)
		os.Exit(1)
	}

	list.Flags().StringP("level", "l", "", "filter logs by a level (info, debug, warning, error, critical)")
	list.Flags().StringP("type", "t", "", "filter logs by a resource type")
	list.Flags().StringP("uuid", "u", "", "filter logs by a resource UUID")

	cmd.AddCommand(
		list,
	)

	return cmd
}

type options struct {
	Base        *cli.Base
	LogsOptions govultr.LogsOptions
}

func (o *options) list() ([]govultr.Log, *govultr.LogsMeta, error) {
	logs, meta, _, err := o.Base.Client.Logs.List(context.Background(), o.LogsOptions)
	return logs, meta, err
}
