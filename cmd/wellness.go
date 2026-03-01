package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var wellnessCmd = &cobra.Command{
	Use:   "wellness",
	Short: "List wellness records for a date range",
	Long: `List wellness records for a date range. Each record corresponds to
one day and may include metrics such as CTL, ATL, resting HR, weight, and sleep.

Required flags:
  --oldest string   start date, inclusive (YYYY-MM-DD)
  --newest string   end date, inclusive (YYYY-MM-DD)

Examples:
  intervals-cli wellness --oldest 2024-01-01 --newest 2024-01-31
  intervals-cli wellness --oldest 2024-01-01 --newest 2024-01-31 --format text`,
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

		data, err := c.Get(c.AthletePath("/wellness"), params)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("list wellness: %v", err))
		}

		var records []map[string]interface{}
		if err := json.Unmarshal(data, &records); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			headers := []string{"DATE", "CTL", "ATL", "RESTING_HR", "WEIGHT", "SLEEP"}
			var rows [][]string
			for _, r := range records {
				rows = append(rows, []string{
					str(r["id"]),
					str(r["ctl"]),
					str(r["atl"]),
					str(r["restingHR"]),
					str(r["weight"]),
					str(r["sleepSecs"]),
				})
			}
			output.WriteText(headers, rows)
			return nil
		}

		if err := output.WriteJSON(records, output.Metadata{
			Count:   len(records),
			Tool:    "intervals-cli",
			Version: version,
		}); err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_OUTPUT", err.Error())
		}
		return nil
	},
}

func init() {
	wellnessCmd.Flags().String("oldest", "", "start date, inclusive (YYYY-MM-DD)")
	wellnessCmd.Flags().String("newest", "", "end date, inclusive (YYYY-MM-DD)")
	rootCmd.AddCommand(wellnessCmd)
}
