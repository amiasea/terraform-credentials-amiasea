package commands

import (
	"testing"
)

func TestInitCommand_Structure(t *testing.T) {
	cmd := NewInitCmd()

	if cmd.Use != "init" {
		t.Errorf("expected Use='init', got %s", cmd.Use)
	}
	if cmd.Flags().Lookup("domain") == nil {
		t.Error("missing --domain flag")
	}
}
