// Package log provides application logging utilities.
package log

import (
	"fmt"
	"os"
	"time"
)

var (
	green = "\033[32m"
	reset = "\033[0m"
)

// Info logs an informational message to standard error with a specific format.
func Info(msg string) {
	fmt.Fprintln(os.Stderr, green+"[tfcred] "+msg+reset)
}

// Err logs an error message to standard error with a specific format.
func Err(msg string) {
	fmt.Fprintln(os.Stderr, "[tfcred][error] "+msg)
}

// AppendFile appends a timestamped error message to the specified file.
func AppendFile(path string, msg string) {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)

	fmt.Fprintf(
		file,
		"%s %s\n",
		timestamp,
		msg,
	)
}
