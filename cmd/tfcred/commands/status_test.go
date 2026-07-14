package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestStatusCommand_NoActiveContext(t *testing.T) {
	setupCommandTest(t)

	cmd := NewStatusCmd()

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
		"No active context",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}

func TestStatusCommand_ActiveContext(t *testing.T) {
	setupCommandTest(t)

	store.Add(
		"production",
		"amiasea",
		"org",
		"app.terraform.io",
		"",
	)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed getting cwd: %v",
			err,
		)
	}

	if err := store.BindDirectory(
		cwd,
		"production",
	); err != nil {
		t.Fatalf(
			"failed binding: %v",
			err,
		)
	}

	cmd := NewStatusCmd()

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
		"context=production",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}
