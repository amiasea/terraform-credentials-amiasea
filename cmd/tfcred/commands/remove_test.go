package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestRemoveCommand_RemovesContext(t *testing.T) {
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
		t.Fatalf("expected token to be vaulted")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed getting cwd: %v", err)
	}

	if err := store.BindDirectory(
		cwd,
		"production",
	); err != nil {
		t.Fatalf(
			"failed binding directory: %v",
			err,
		)
	}

	cmd := NewRemoveCmd()

	cmd.SetArgs([]string{
		"production",
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
		"removed successfully",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}

	config := store.Load()

	if _, exists := config.Contexts["production"]; exists {
		t.Fatal("expected context to be removed")
	}

	if _, exists := config.Directories[cwd]; exists {
		t.Fatal("expected directory binding to be removed")
	}

	if _, exists := vault.tokens[key]; exists {
		t.Fatal("expected token to be deleted")
	}
}
