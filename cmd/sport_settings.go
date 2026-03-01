package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var sportSettingsCmd = &cobra.Command{
	Use:   "sport-settings",
	Short: "Get sport settings",
	Long: `Get training settings for a specific sport, including threshold values
(FTP, LTHR), heart rate zones, and power zones.

Required flags:
  --sport string   sport type (Ride, Run, Swim, Walk, VirtualRide, etc.)

Examples:
  intervals-cli sport-settings --sport Ride
  intervals-cli sport-settings --sport Run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		sport, _ := cmd.Flags().GetString("sport")
		if sport == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--sport is required")
		}

		params := url.Values{
			"sport": {sport},
		}

		data, err := c.Get(c.AthletePath("/sport-settings"), params)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("get sport settings: %v", err))
		}

		var settings interface{}
		if err := json.Unmarshal(data, &settings); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
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
	sportSettingsCmd.Flags().String("sport", "", "sport type (Ride, Run, Swim, etc.)")
	rootCmd.AddCommand(sportSettingsCmd)
}
