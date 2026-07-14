package commands

import (
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestConfigCommand_SetDefaultDomain(t *testing.T) {
	setupCommandTest(t)

	cmd := NewConfigCmd()

	cmd.SetArgs([]string{
		"--default-domain",
		"app.eu.terraform.io",
	})

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
		"default domain set to app.eu.terraform.io",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}

	config := store.Load()

	if config.DefaultDomain != "app.eu.terraform.io" {
		t.Fatalf(
			"expected app.eu.terraform.io, got %s",
			config.DefaultDomain,
		)
	}
}

func TestConfigCommand_Show(t *testing.T) {
	setupCommandTest(t)

	cmd := NewConfigCmd()

	cmd.SetArgs([]string{
		"--show",
	})

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
		"default_domain=app.terraform.io",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}
