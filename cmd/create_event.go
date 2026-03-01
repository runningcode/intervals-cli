package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var createEventCmd = &cobra.Command{
	Use:   "create-event",
	Short: "Create a calendar event",
	Long: `Create a calendar event (workout, note, race, etc.).

Required flags:
  --name string       event name
  --category string   event category (WORKOUT, NOTE, RACE, etc.)
  --start string      start date/time in ISO 8601 format (e.g. 2024-02-15T18:00:00)

Optional flags:
  --description string   event description or workout notes
  --color string         event color (hex or named color)

Examples:
  intervals-cli create-event --name "VO2Max Intervals" --category WORKOUT --start 2024-02-15T18:00:00
  intervals-cli create-event --name "A Race" --category RACE --start 2024-03-10T09:00:00 --description "Target: sub-3hr"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		name, _ := cmd.Flags().GetString("name")
		category, _ := cmd.Flags().GetString("category")
		start, _ := cmd.Flags().GetString("start")

		if name == "" || category == "" || start == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--name, --category, and --start are required")
		}

		body := map[string]interface{}{
			"name":             name,
			"category":         category,
			"start_date_local": start,
		}

		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		if v, _ := cmd.Flags().GetString("color"); v != "" {
			body["color"] = v
		}

		data, err := c.Post(c.AthletePath("/events"), body)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("create event: %v", err))
		}

		var event map[string]interface{}
		if err := json.Unmarshal(data, &event); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			fmt.Printf("Created event %s\n", str(event["id"]))
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
	createEventCmd.Flags().String("name", "", "event name")
	createEventCmd.Flags().String("category", "", "event category (WORKOUT, NOTE, RACE, etc.)")
	createEventCmd.Flags().String("start", "", "start date/time in ISO 8601 format")
	createEventCmd.Flags().String("description", "", "event description or workout notes")
	createEventCmd.Flags().String("color", "", "event color")
	rootCmd.AddCommand(createEventCmd)
}
