package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var updateSportSettingsCmd = &cobra.Command{
	Use:   "update-sport-settings",
	Short: "Update sport settings",
	Long: `Update training settings for a specific sport. Only flags that are
explicitly provided are sent to the API; omitted flags leave those fields unchanged.

Required flags:
  --sport string   sport type (Ride, Run, Swim, Walk, VirtualRide, etc.)

Optional flags (at least one required):
  --ftp int          functional threshold power in watts (cycling)
  --lthr int         lactate threshold heart rate in bpm
  --max-hr int       maximum heart rate in bpm
  --resting-hr int   resting heart rate in bpm

Examples:
  intervals-cli update-sport-settings --sport Ride --ftp 310
  intervals-cli update-sport-settings --sport Run --lthr 168 --max-hr 185`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		sport, _ := cmd.Flags().GetString("sport")
		if sport == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--sport is required")
		}

		body := map[string]interface{}{}

		if v, _ := cmd.Flags().GetInt("ftp"); cmd.Flags().Changed("ftp") {
			body["ftp"] = v
		}
		if v, _ := cmd.Flags().GetInt("lthr"); cmd.Flags().Changed("lthr") {
			body["lthr"] = v
		}
		if v, _ := cmd.Flags().GetInt("max-hr"); cmd.Flags().Changed("max-hr") {
			body["max_hr"] = v
		}
		if v, _ := cmd.Flags().GetInt("resting-hr"); cmd.Flags().Changed("resting-hr") {
			body["resting_hr"] = v
		}

		if len(body) == 0 {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "at least one setting to update is required")
		}

		data, err := c.Put(c.AthletePath("/sport-settings/"+sport), body)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("update sport settings: %v", err))
		}

		var settings interface{}
		if err := json.Unmarshal(data, &settings); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			fmt.Printf("Updated %s settings\n", sport)
			return nil
		}

		if err := output.WriteJSON(settings, output.Metadata{
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
	updateSportSettingsCmd.Flags().String("sport", "", "sport type (Ride, Run, Swim, etc.)")
	updateSportSettingsCmd.Flags().Int("ftp", 0, "functional threshold power in watts")
	updateSportSettingsCmd.Flags().Int("lthr", 0, "lactate threshold heart rate in bpm")
	updateSportSettingsCmd.Flags().Int("max-hr", 0, "maximum heart rate in bpm")
	updateSportSettingsCmd.Flags().Int("resting-hr", 0, "resting heart rate in bpm")
	rootCmd.AddCommand(updateSportSettingsCmd)
}
