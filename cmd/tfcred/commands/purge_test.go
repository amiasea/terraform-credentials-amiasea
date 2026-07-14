package commands

import (
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestPurgeCommand_Force(t *testing.T) {
	vault := setupCommandTest(t)

	token := "FAKETOKENXXXXX.atlasv1.abcdefghijklmnopqrstuvwxyzABCDE"

	store.Add(
		"production",
		"amiasea",
		"org",
		"app.terraform.io",
		token,
	)

	key := store.TokenVaultKey(
		"app.terraform.io",
		"org",
		"amiasea",
	)

	if _, exists := vault.tokens[key]; !exists {
		t.Fatal("expected token before purge")
	}

	cmd := NewPurgeCmd()

	cmd.SetArgs([]string{
		"--force",
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
		"SUCCESS: System entirely purged",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}

	config := store.Load()

	if len(config.Contexts) != 0 {
		t.Fatalf(
			"expected contexts to be empty",
		)
	}

	if len(config.Directories) != 0 {
		t.Fatalf(
			"expected directories to be empty",
		)
	}

	if _, exists := vault.tokens[key]; exists {
		t.Fatal("expected token to be removed")
	}
}
