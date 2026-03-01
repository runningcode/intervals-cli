package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/runningcode/intervals-cli/internal/config"
	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set athlete ID and API key",
	Long: `Save your Intervals.icu credentials to ~/.intervals-cli/config.json.

Credentials can also be provided via environment variables, which take
precedence over the config file:
  INTERVALS_ATHLETE_ID
  INTERVALS_API_KEY

Your athlete ID and API key are found under Settings → API in the
Intervals.icu web app.

Required flags:
  --athlete-id string   Intervals.icu athlete ID
  --api-key string      Intervals.icu API key

Examples:
  intervals-cli config --athlete-id i12345 --api-key abc123xyz`,
	RunE: func(cmd *cobra.Command, args []string) error {
		athleteID, _ := cmd.Flags().GetString("athlete-id")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if athleteID == "" || apiKey == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--athlete-id and --api-key are required")
		}

		cfg := &config.Config{
			AthleteID: athleteID,
			APIKey:    apiKey,
		}
		if err := config.Save(cfg); err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_CONFIG", fmt.Sprintf("save config: %v", err))
		}

		if format == "text" {
			fmt.Printf("Config saved to %s\n", config.Path())
			return nil
		}

		if err := output.WriteJSON(map[string]string{
			"path":   config.Path(),
			"status": "saved",
		}, output.Metadata{
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
	configCmd.Flags().String("athlete-id", "", "Intervals.icu athlete ID")
	configCmd.Flags().String("api-key", "", "Intervals.icu API key")
	rootCmd.AddCommand(configCmd)
}
