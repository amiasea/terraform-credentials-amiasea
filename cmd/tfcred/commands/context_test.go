package commands

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestContextCommand_JSONOutput(t *testing.T) {
	setupCommandTest(t)

	store.Add(
		"production",
		"amiasea",
		"org",
		"app.terraform.io",
		"",
	)

	cmd := NewContextCmd()

	cmd.SetArgs([]string{
		"--json",
	})

	output := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	})

	var report map[string]any

	if err := json.Unmarshal(
		[]byte(output),
		&report,
	); err != nil {
		t.Fatalf(
			"failed parsing json output: %v\n%s",
			err,
			output,
		)
	}

	if report["active_context"] != "" {
		t.Fatalf(
			"expected no active context, got %v",
			report["active_context"],
		)
	}
}

func TestContextCommand_AllContexts(t *testing.T) {
	setupCommandTest(t)

	store.Add(
		"alpha",
		"org-a",
		"org",
		"app.terraform.io",
		"",
	)

	store.Add(
		"beta",
		"org-b",
		"team",
		"app.terraform.io",
		"",
	)

	cmd := NewContextCmd()

	cmd.SetArgs([]string{
		"--json",
		"--all",
	})

	output := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	})

	var report map[string]any

	if err := json.Unmarshal(
		[]byte(output),
		&report,
	); err != nil {
		t.Fatalf(
			"invalid json: %v\n%s",
			err,
			output,
		)
	}

	contexts, ok := report["contexts"].([]any)

	if !ok {
		t.Fatalf(
			"expected contexts array, got %T",
			report["contexts"],
		)
	}

	if len(contexts) != 2 {
		t.Fatalf(
			"expected two contexts, got %d",
			len(contexts),
		)
	}
}

func TestContextCommand_TextOutput(t *testing.T) {
	setupCommandTest(t)

	cmd := NewContextCmd()

	output := captureStdout(t, func() {
		cmd.SetArgs(nil)

		if err := cmd.Execute(); err != nil {
			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	})

	if !strings.Contains(
		output,
		"active_context=",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}
}
