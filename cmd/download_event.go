package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	ierrors "github.com/runningcode/intervals-cli/internal/errors"
	"github.com/runningcode/intervals-cli/internal/output"
)

var downloadEventCmd = &cobra.Command{
	Use:   "download-event",
	Short: "Download a workout file (zwo, mrc, erg, fit)",
	Long: `Download a workout file from a calendar event and save it to disk.

Note: use --dl-format (not --format) to specify the file type, since --format
is the global flag controlling JSON vs text output.

Required flags:
  --id string          event ID
  --dl-format string   file format to download: zwo, mrc, erg, or fit
  --output string      path to write the downloaded file

Examples:
  intervals-cli download-event --id 98765 --dl-format zwo --output workout.zwo
  intervals-cli download-event --id 98765 --dl-format fit --output workout.fit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := requireClient()

		id, _ := cmd.Flags().GetString("id")
		dlFormat, _ := cmd.Flags().GetString("dl-format")
		outFile, _ := cmd.Flags().GetString("output")

		if id == "" || dlFormat == "" || outFile == "" {
			ierrors.Exit(ierrors.ExitUsage, "ERR_USAGE", "--id, --dl-format, and --output are required")
		}

		path := c.AthletePath("/events/" + id + "/download/" + dlFormat)
		data, err := c.GetRaw(path, nil)
		if err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_API", fmt.Sprintf("download event: %v", err))
		}

		if err := os.WriteFile(outFile, data, 0644); err != nil {
			ierrors.Exit(ierrors.ExitError, "ERR_IO", fmt.Sprintf("write file: %v", err))
		}

		if format == "text" {
			fmt.Printf("Downloaded to %s (%d bytes)\n", outFile, len(data))
			return nil
		}

		if err := output.WriteJSON(map[string]interface{}{
			"path":  outFile,
			"bytes": len(data),
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
	downloadEventCmd.Flags().String("id", "", "event ID")
	downloadEventCmd.Flags().String("dl-format", "", "file format to download (zwo, mrc, erg, fit)")
	downloadEventCmd.Flags().String("output", "", "path to write the downloaded file")
	rootCmd.AddCommand(downloadEventCmd)
}
