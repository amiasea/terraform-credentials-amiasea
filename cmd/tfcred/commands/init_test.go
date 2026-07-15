package commands

import (
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestInitCommand_WithDomain(t *testing.T) {
	// Setup minimal environment
	t.Setenv("APPDATA", t.TempDir())
	t.Setenv("LOCALAPPDATA", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	// Use fake vault if your store needs it
	store.SetVault(newFakeVault())

	cmd := NewInitCmd()

	cmd.SetArgs([]string{"--domain", "app.terraform.io"})

	output := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "Initialization completed successfully") {
		t.Fatalf("unexpected output:\n%s", output)
	}

	// Basic check that store was initialized
	config := store.Load()
	if config.DefaultDomain != "app.terraform.io" {
		t.Fatalf("expected DefaultDomain = app.terraform.io, got %s", config.DefaultDomain)
	}
}

func TestInitCommand_NoDomainFlag(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	t.Setenv("LOCALAPPDATA", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	store.SetVault(newFakeVault())

	cmd := NewInitCmd()
	cmd.SetArgs([]string{}) // no --domain flag

	output := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "Initialization completed successfully") {
		t.Fatalf("expected success message, got:\n%s", output)
	}
}

func TestInitCommand_Structure(t *testing.T) {
	cmd := NewInitCmd()

	if cmd.Use != "init" {
		t.Errorf("expected Use='init', got %s", cmd.Use)
	}
	if cmd.Flags().Lookup("domain") == nil {
		t.Error("missing --domain flag")
	}
}
