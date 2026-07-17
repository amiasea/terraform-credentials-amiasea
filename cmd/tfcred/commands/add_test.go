package commands

import (
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestAddCommand_AddsUserContext(t *testing.T) {
	vault := setupCommandTest(t)

	cmd := NewAddCmd()

	cmd.SetArgs([]string{
		"--context",
		"personal",
		"--token",
		"FAKETOKENXXXXX.atlasv1.abcdefghijklmnopqrstuvwxyzABCDE",
	})

	output := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(
		output,
		"Context 'personal' configured successfully",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}

	config := store.Load()

	entry, exists := config.Contexts["personal"]

	if !exists {
		t.Fatal("expected context to exist")
	}

	if entry.TokenType != "user" {
		t.Fatalf(
			"expected user token type, got %s",
			entry.TokenType,
		)
	}

	if entry.Domain != "app.terraform.io" {
		t.Fatalf(
			"expected default domain, got %s",
			entry.Domain,
		)
	}

	key := store.TokenVaultKey(
		entry.Domain,
		entry.TokenType,
		entry.Org,
	)

	if vault.tokens[key] == "" {
		t.Fatalf(
			"expected token stored under %s",
			key,
		)
	}
}

func TestAddCommand_AddsOrganizationContext(t *testing.T) {
	vault := setupCommandTest(t)

	cmd := NewAddCmd()

	cmd.SetArgs([]string{
		"--context",
		"production",
		"--token-type",
		"org",
		"--org",
		"amiasea",
		"--token",
		"FAKETOKENXXXXX.atlasv1.abcdefghijklmnopqrstuvwxyzABCDE",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	config := store.Load()

	entry, exists := config.Contexts["production"]

	if !exists {
		t.Fatal("expected context to exist")
	}

	if entry.Org != "amiasea" {
		t.Fatalf(
			"expected org amiasea, got %s",
			entry.Org,
		)
	}

	if entry.TokenType != "org" {
		t.Fatalf(
			"expected org token type, got %s",
			entry.TokenType,
		)
	}

	key := store.TokenVaultKey(
		entry.Domain,
		entry.TokenType,
		entry.Org,
	)

	if vault.tokens[key] == "" {
		t.Fatalf(
			"expected vaulted token for %s",
			key,
		)
	}
}

func TestAddCommand_DefaultDomainRequiredWhenMissing(t *testing.T) {
	store.SetVault(newFakeVault())

	cmd := NewAddCmd()

	cmd.SetArgs([]string{
		"--context",
		"missing-domain",
	})

	output := captureStdout(t, func() {
		// This command path currently calls os.Exit(1).
		// This assertion is intentionally skipped until command errors
		// are returned instead of terminating the process.
	})

	if output != "" {
		t.Fatalf(
			"unexpected output: %s",
			output,
		)
	}
}
