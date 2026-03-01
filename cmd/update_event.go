package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var updateEventCmd = &cobra.Command{
	Use:   "update-event",
	Short: "Update an existing calendar event",
	Long: `Update fields on an existing calendar event. Only flags that are
explicitly provided are sent to the API; omitted flags leave those fields unchanged.

Required flags:
  --id string   event ID

Optional flags (at least one required):
  --name string         new event name
  --description string  new event description
  --category string     new event category (WORKOUT, NOTE, RACE, etc.)
  --color string        new event color

Examples:
  intervals-cli update-event --id 98765 --name "Rescheduled VO2Max"
  intervals-cli update-event --id 98765 --description "Bumped FTP targets to 310w"`,
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
		if v, _ := cmd.Flags().GetString("category"); cmd.Flags().Changed("category") {
			body["category"] = v
		}
		if v, _ := cmd.Flags().GetString("color"); cmd.Flags().Changed("color") {
			body["color"] = v
		}

		if len(body) == 0 {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "at least one field to update is required")
		}

		data, err := c.Put(c.AthletePath("/events/"+id), body)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("update event: %v", err))
		}

		var event map[string]interface{}
		if err := json.Unmarshal(data, &event); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			fmt.Printf("Updated event %s\n", id)
			return nil
		}

		if err := output.WriteJSON(event, output.Metadata{
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
	updateEventCmd.Flags().String("id", "", "event ID")
	updateEventCmd.Flags().String("name", "", "new event name")
	updateEventCmd.Flags().String("description", "", "new event description")
	updateEventCmd.Flags().String("category", "", "new event category")
	updateEventCmd.Flags().String("color", "", "new event color")
	rootCmd.AddCommand(updateEventCmd)
}
