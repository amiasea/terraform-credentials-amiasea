package commands

import (
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersionCmd()

	root := cmd.Root()
	root.Version = "test-version"

	output := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	})

	if !strings.Contains(
		output,
		"test-version",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}
