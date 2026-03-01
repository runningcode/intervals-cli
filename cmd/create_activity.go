package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var createActivityCmd = &cobra.Command{
	Use:   "create-activity",
	Short: "Create a manual activity",
	Long: `Create a manual activity entry.

Required flags:
  --name string    activity name
  --type string    activity type (Run, Ride, Swim, Walk, VirtualRide, etc.)
  --start string   start time in ISO 8601 format (e.g. 2024-01-15T06:00:00)

Optional flags:
  --distance float     distance in meters
  --moving-time int    moving time in seconds
  --description string activity description

Examples:
  intervals-cli create-activity --name "Morning Run" --type Run --start 2024-01-15T06:00:00
  intervals-cli create-activity --name "Long Ride" --type Ride --start 2024-01-20T08:00:00 --distance 80000 --moving-time 10800`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		name, _ := cmd.Flags().GetString("name")
		actType, _ := cmd.Flags().GetString("type")
		start, _ := cmd.Flags().GetString("start")

		if name == "" || actType == "" || start == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--name, --type, and --start are required")
		}

		body := map[string]interface{}{
			"name":             name,
			"type":             actType,
			"start_date_local": start,
		}

		if v, _ := cmd.Flags().GetFloat64("distance"); cmd.Flags().Changed("distance") {
			body["distance"] = v
		}
		if v, _ := cmd.Flags().GetInt("moving-time"); cmd.Flags().Changed("moving-time") {
			body["moving_time"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}

		data, err := c.Post(c.AthletePath("/activities"), body)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("create activity: %v", err))
		}

		var activity map[string]interface{}
		if err := json.Unmarshal(data, &activity); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			fmt.Printf("Created activity %s\n", str(activity["id"]))
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
	createActivityCmd.Flags().String("name", "", "activity name")
	createActivityCmd.Flags().String("type", "", "activity type (Run, Ride, Swim, etc.)")
	createActivityCmd.Flags().String("start", "", "start time in ISO 8601 format")
	createActivityCmd.Flags().Float64("distance", 0, "distance in meters")
	createActivityCmd.Flags().Int("moving-time", 0, "moving time in seconds")
	createActivityCmd.Flags().String("description", "", "activity description")
	rootCmd.AddCommand(createActivityCmd)
}
