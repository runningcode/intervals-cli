package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var updateActivityCmd = &cobra.Command{
	Use:   "update-activity",
	Short: "Update an existing activity",
	Long: `Update fields on an existing activity. Only flags that are explicitly
provided are sent to the API; omitted flags leave those fields unchanged.

Required flags:
  --id string   activity ID

Optional flags (at least one required):
  --name string         new activity name
  --description string  new activity description
  --type string         new activity type (Run, Ride, Swim, etc.)

Examples:
  intervals-cli update-activity --id 12345678 --name "Easy Run"
  intervals-cli update-activity --id 12345678 --name "Recovery Ride" --description "Zone 2 only"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--id is required")
		}

		body := map[string]interface{}{}

		if v, _ := cmd.Flags().GetString("name"); cmd.Flags().Changed("name") {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); cmd.Flags().Changed("description") {
			body["description"] = v
		}
		if v, _ := cmd.Flags().GetString("type"); cmd.Flags().Changed("type") {
			body["type"] = v
		}

		if len(body) == 0 {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "at least one field to update is required")
		}

		data, err := c.Put(c.AthletePath("/activities/"+id), body)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("update activity: %v", err))
		}

		var activity map[string]interface{}
		if err := json.Unmarshal(data, &activity); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			fmt.Printf("Updated activity %s\n", id)
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
	updateActivityCmd.Flags().String("id", "", "activity ID")
	updateActivityCmd.Flags().String("name", "", "new activity name")
	updateActivityCmd.Flags().String("description", "", "new activity description")
	updateActivityCmd.Flags().String("type", "", "new activity type")
	rootCmd.AddCommand(updateActivityCmd)
}
