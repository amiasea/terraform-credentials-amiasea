package e2e

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv(
		"TF_CLI_CONFIG_FILE",
		tfConfigFile,
	)

	os.Setenv(
		"TF_CRED_CONTEXT_DIR",
		tfCredContextDir,
	)

	if err := buildTfcred(); err != nil {
		panic(err)
	}

	os.Exit(
		m.Run(),
	)
}

func TestE2E(t *testing.T) {
	t.Run(
		"credentials_helper",
		func(t *testing.T) {
			runCredentialsHelperBasic(t)
		},
	)

	t.Run(
		"configuration",
		func(t *testing.T) {
			runConfigurationTests(t)
		},
	)

	t.Run(
		"context_lifecycle",
		func(t *testing.T) {
			runContextLifecycleTests(t)
		},
	)

	t.Run(
		"cleanup",
		func(t *testing.T) {
			runCleanupTests(t)
		},
	)
}
