package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var activitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "List activities for a date range",
	Long: `List activities for a date range.

Required flags:
  --oldest string   start date, inclusive (YYYY-MM-DD)
  --newest string   end date, inclusive (YYYY-MM-DD)

Examples:
  intervals-cli activities --oldest 2024-01-01 --newest 2024-01-31
  intervals-cli activities --oldest 2024-01-01 --newest 2024-01-31 --format text`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		oldest, _ := cmd.Flags().GetString("oldest")
		newest, _ := cmd.Flags().GetString("newest")

		if oldest == "" || newest == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--oldest and --newest are required")
		}

		params := url.Values{
			"oldest": {oldest},
			"newest": {newest},
		}

		data, err := c.Get(c.AthletePath("/activities"), params)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("list activities: %v", err))
		}

		var activities []map[string]interface{}
		if err := json.Unmarshal(data, &activities); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			headers := []string{"ID", "NAME", "TYPE", "DATE"}
			var rows [][]string
			for _, a := range activities {
				rows = append(rows, []string{
					str(a["id"]),
					str(a["name"]),
					str(a["type"]),
					str(a["start_date_local"]),
				})
			}
			output.WriteText(headers, rows)
			return nil
		}

		if err := output.WriteJSON(activities, output.Metadata{
			Count:   len(activities),
			Tool:    "intervals-cli",
			Version: version,
		}); err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_OUTPUT", err.Error())
		}
		return nil
	},
}

func init() {
	activitiesCmd.Flags().String("oldest", "", "start date, inclusive (YYYY-MM-DD)")
	activitiesCmd.Flags().String("newest", "", "end date, inclusive (YYYY-MM-DD)")
	rootCmd.AddCommand(activitiesCmd)
}

func str(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
