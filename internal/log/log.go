// Package log provides application logging utilities.
package log

import (
	"fmt"
	"os"
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
