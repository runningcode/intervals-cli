package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var updateWellnessCmd = &cobra.Command{
	Use:   "update-wellness",
	Short: "Update a wellness record for a specific date",
	Long: `Update a wellness record for a specific date. Only flags that are
explicitly provided are sent to the API; omitted flags leave those fields unchanged.

Required flags:
  --date string   date to update (YYYY-MM-DD)

Optional flags (at least one required):
  --sleep-secs int    sleep duration in seconds (e.g. 28800 = 8 hours)
  --resting-hr int    resting heart rate in bpm
  --weight float      body weight in kg
  --fatigue int       fatigue level (1-5)
  --soreness int      soreness level (1-5)
  --stress int        stress level (1-5)
  --mood int          mood level (1-5)

Examples:
  intervals-cli update-wellness --date 2024-01-15 --sleep-secs 28800 --resting-hr 52
  intervals-cli update-wellness --date 2024-01-15 --fatigue 2 --mood 4 --stress 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--date is required")
		}

		body := map[string]interface{}{
			"id": date,
		}

		if v, _ := cmd.Flags().GetInt("sleep-secs"); cmd.Flags().Changed("sleep-secs") {
			body["sleepSecs"] = v
		}
		if v, _ := cmd.Flags().GetInt("resting-hr"); cmd.Flags().Changed("resting-hr") {
			body["restingHR"] = v
		}
		if v, _ := cmd.Flags().GetFloat64("weight"); cmd.Flags().Changed("weight") {
			body["weight"] = v
		}
		if v, _ := cmd.Flags().GetInt("fatigue"); cmd.Flags().Changed("fatigue") {
			body["fatigue"] = v
		}
		if v, _ := cmd.Flags().GetInt("soreness"); cmd.Flags().Changed("soreness") {
			body["soreness"] = v
		}
		if v, _ := cmd.Flags().GetInt("stress"); cmd.Flags().Changed("stress") {
			body["stress"] = v
		}
		if v, _ := cmd.Flags().GetInt("mood"); cmd.Flags().Changed("mood") {
			body["mood"] = v
		}

		data, err := c.Put(c.AthletePath("/wellness/"+date), body)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("update wellness: %v", err))
		}

		var record map[string]interface{}
		if err := json.Unmarshal(data, &record); err != nil {
			ierrors.Exit(ierrors.ExitDataError, "ERR_PARSE", fmt.Sprintf("parse response: %v", err))
		}

		if format == "text" {
			fmt.Printf("Updated wellness for %s\n", date)
			return nil
		}

		if err := output.WriteJSON(record, output.Metadata{
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
	updateWellnessCmd.Flags().String("date", "", "date to update (YYYY-MM-DD)")
	updateWellnessCmd.Flags().Int("sleep-secs", 0, "sleep duration in seconds")
	updateWellnessCmd.Flags().Int("resting-hr", 0, "resting heart rate in bpm")
	updateWellnessCmd.Flags().Float64("weight", 0, "body weight in kg")
	updateWellnessCmd.Flags().Int("fatigue", 0, "fatigue level (1-5)")
	updateWellnessCmd.Flags().Int("soreness", 0, "soreness level (1-5)")
	updateWellnessCmd.Flags().Int("stress", 0, "stress level (1-5)")
	updateWellnessCmd.Flags().Int("mood", 0, "mood level (1-5)")
	rootCmd.AddCommand(updateWellnessCmd)
}
