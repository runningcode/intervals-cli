package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Get a single activity by ID",
	Long: `Get a single activity by ID.

Required flags:
  --id string   activity ID

Examples:
  intervals-cli activity --id 12345678
  intervals-cli activity --id 12345678 --format text`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--id is required")
		}

		data, err := c.Get(c.AthletePath("/activities/"+id), nil)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("get activity: %v", err))
		}

		var activity map[string]interface{}
		if err := json.Unmarshal(data, &activity); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			headers := []string{"FIELD", "VALUE"}
			rows := [][]string{
				{"id", str(activity["id"])},
				{"name", str(activity["name"])},
				{"type", str(activity["type"])},
				{"date", str(activity["start_date_local"])},
				{"distance", str(activity["distance"])},
				{"moving_time", str(activity["moving_time"])},
				{"icu_training_load", str(activity["icu_training_load"])},
			}
			output.WriteText(headers, rows)
			return nil
		}

		if err := output.WriteJSON(activity, output.Metadata{
			Count:   1,
			Tool:    "intervals-cli",
			Version: version,
		}); err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_OUTPUT", err.Error())
		}
		return nil
	},
}

func init() {
	activityCmd.Flags().String("id", "", "activity ID")
	rootCmd.AddCommand(activityCmd)
}
