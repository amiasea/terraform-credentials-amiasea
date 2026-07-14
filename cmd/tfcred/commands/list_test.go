package commands

import (
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestListCommand_Empty(t *testing.T) {
	setupCommandTest(t)

	cmd := NewListCmd()

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
		"no contexts configured",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}

func TestListCommand_WithContexts(t *testing.T) {
	setupCommandTest(t)

	store.Add(
		"production",
		"amiasea",
		"org",
		"app.terraform.io",
		"",
	)

	cmd := NewListCmd()

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
		"production",
	) {
		t.Fatalf(
			"expected context in output:\n%s",
			output,
		)
	}
}
