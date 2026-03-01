package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List calendar events for a date range",
	Long: `List calendar events for a date range.

Required flags:
  --oldest string   start date, inclusive (YYYY-MM-DD)
  --newest string   end date, inclusive (YYYY-MM-DD)

Optional flags:
  --category string   filter by event category (WORKOUT, NOTE, RACE, etc.)

Examples:
  intervals-cli events --oldest 2024-02-01 --newest 2024-02-29
  intervals-cli events --oldest 2024-02-01 --newest 2024-02-29 --category WORKOUT
  intervals-cli events --oldest 2024-02-01 --newest 2024-02-29 --format text`,
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
		if v, _ := cmd.Flags().GetString("category"); v != "" {
			params.Set("category", v)
		}

		data, err := c.Get(c.AthletePath("/events"), params)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("list events: %v", err))
		}

		var events []map[string]interface{}
		if err := json.Unmarshal(data, &events); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			headers := []string{"ID", "NAME", "CATEGORY", "START"}
			var rows [][]string
			for _, e := range events {
				rows = append(rows, []string{
					str(e["id"]),
					str(e["name"]),
					str(e["category"]),
					str(e["start_date_local"]),
				})
			}
			output.WriteText(headers, rows)
			return nil
		}

		if err := output.WriteJSON(events, output.Metadata{
			Count:   len(events),
			Tool:    "intervals-cli",
			Version: version,
		}); err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_OUTPUT", err.Error())
		}
		return nil
	},
}

func init() {
	eventsCmd.Flags().String("oldest", "", "start date, inclusive (YYYY-MM-DD)")
	eventsCmd.Flags().String("newest", "", "end date, inclusive (YYYY-MM-DD)")
	eventsCmd.Flags().String("category", "", "filter by category (WORKOUT, NOTE, RACE, etc.)")
	rootCmd.AddCommand(eventsCmd)
}
