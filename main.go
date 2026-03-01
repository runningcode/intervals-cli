package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/runningcode/intervals-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Commands call ierrors.Exit directly for expected failures, so an
		// error here is unexpected (e.g. cobra flag parsing). Emit error
		// JSON so the process never exits silently.
		b, _ := json.Marshal(map[string]string{"error": err.Error(), "code": "ERR_UNKNOWN"})
		fmt.Fprintln(os.Stderr, string(b))
		os.Exit(1)
	}
}
