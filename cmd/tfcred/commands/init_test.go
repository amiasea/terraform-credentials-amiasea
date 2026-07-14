package commands

import (
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestInitCommand_WithDomain(t *testing.T) {
	t.Setenv(
		"TF_CRED_CONTEXT_DIR",
		t.TempDir(),
	)

	store.SetVault(newFakeVault())

	cmd := NewInitCmd()

	cmd.SetArgs([]string{
		"--domain",
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
		"initialized with default domain app.eu.terraform.io",
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
