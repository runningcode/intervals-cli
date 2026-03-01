package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var activityStreamsCmd = &cobra.Command{
	Use:   "activity-streams",
	Short: "Get activity streams (time-series data)",
	Long: `Get time-series streams for an activity (power, heart rate, cadence, etc.).

The response structure varies by activity type and what data was recorded.
Always outputs JSON regardless of --format, as the nested stream structure
does not have a useful tabular representation.

Required flags:
  --id string   activity ID

Examples:
  intervals-cli activity-streams --id 12345678`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--id is required")
		}

		data, err := c.Get(c.AthletePath("/activities/"+id+"/streams"), nil)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("get activity streams: %v", err))
		}

		var streams interface{}
		if err := json.Unmarshal(data, &streams); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if err := output.WriteJSON(streams, output.Metadata{
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
	activityStreamsCmd.Flags().String("id", "", "activity ID")
	rootCmd.AddCommand(activityStreamsCmd)
}
