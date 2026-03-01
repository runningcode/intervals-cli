package errors

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ExitOK        = 0
	ExitError     = 1
	ExitUsage     = 2
	ExitDataError = 3
)

type ToolError struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func Exit(code int, errCode string, msg string) {
	e := ToolError{Error: msg, Code: errCode}
	b, _ := json.Marshal(e)
	fmt.Fprintln(os.Stderr, string(b))
	os.Exit(code)
}

func ExitWithDetails(code int, errCode string, msg string, details interface{}) {
	e := ToolError{Error: msg, Code: errCode, Details: details}
	b, _ := json.Marshal(e)
	fmt.Fprintln(os.Stderr, string(b))
	os.Exit(code)
}
