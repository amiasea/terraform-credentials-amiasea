package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestCurrentCommand_NoBinding(t *testing.T) {
	setupCommandTest(t)

	cmd := NewCurrentCmd()

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
		"Current directory is not bound",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}

func TestCurrentCommand_BoundDirectory(t *testing.T) {
	setupCommandTest(t)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed getting cwd: %v",
			err,
		)
	}

	store.Add(
		"development",
		"",
		"user",
		"app.terraform.io",
		"",
	)

	if err := store.BindDirectory(
		cwd,
		"development",
	); err != nil {
		t.Fatalf(
			"failed binding directory: %v",
			err,
		)
	}

	cmd := NewCurrentCmd()

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
		"active_context=development",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}
