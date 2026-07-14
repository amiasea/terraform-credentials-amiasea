package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

func TestSwitchCommand_BindsDirectory(t *testing.T) {
	setupCommandTest(t)

	store.Add(
		"development",
		"",
		"user",
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

	cmd := NewSwitchCmd()

	cmd.SetArgs([]string{
		"development",
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
		"Context 'development' is now bound",
	) {
		t.Fatalf(
			"unexpected output:\n%s",
			output,
		)
	}

	contextKey, found := store.ResolveContextByDir(cwd)

	if !found {
		t.Fatal("expected directory binding")
	}

	if contextKey != "development" {
		t.Fatalf(
			"expected development, got %s",
			contextKey,
		)
	}
}
