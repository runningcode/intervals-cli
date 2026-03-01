package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/runningcode/intervals-cli/internal/client"
	"github.com/runningcode/intervals-cli/internal/config"
	ierrors "github.com/runningcode/intervals-cli/internal/errors"
)

var (
	version = "dev"
	format  string
	quiet   bool

	apiClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "intervals-cli",
	Short: "CLI for the Intervals.icu API",
	Long: `intervals-cli wraps the Intervals.icu REST API for managing
athlete training data: activities, events, wellness, and sport settings.

All output is wrapped in a JSON envelope with "data" and "metadata" fields.
Use --format text for human-readable tables.

Exit codes:
  0  success
  1  general error
  2  invalid usage
  3  data error`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&format, "format", "json", "output format: json or text")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-essential output")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf(`{"version": "%s"}`, version) + "\n")
}

// requireClient loads config, applies env var overrides, and returns a
// ready-to-use API client. Commands that talk to the API call this in RunE.
func requireClient() *client.Client {
	if apiClient != nil {
		return apiClient
	}

	cfg, err := config.Load()
	if err != nil {
		ierrors.Exit(ierrors.ExitError, "ERR_CONFIG", fmt.Sprintf("load config: %v", err))
	}

	if v := os.Getenv("INTERVALS_ATHLETE_ID"); v != "" {
		cfg.AthleteID = v
	}
	if v := os.Getenv("INTERVALS_API_KEY"); v != "" {
		cfg.APIKey = v
	}

	if cfg.AthleteID == "" || cfg.APIKey == "" {
		ierrors.Exit(ierrors.ExitUsage, "ERR_AUTH",
			"athlete ID and API key required; run 'intervals-cli config' or set INTERVALS_ATHLETE_ID / INTERVALS_API_KEY")
	}

	apiClient = client.New(cfg.AthleteID, cfg.APIKey)
	return apiClient
}
