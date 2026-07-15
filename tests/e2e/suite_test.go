package e2e

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := prepareE2EEnvironment(); err != nil {
		panic(
			"failed to prepare e2e environment: " + err.Error(),
		)
	}

	if err := os.Setenv(
		"TF_CLI_CONFIG_FILE",
		tfConfigFile,
	); err != nil {
		panic(
			"failed to set TF_CLI_CONFIG_FILE: " + err.Error(),
		)
	}

	if err := os.Setenv(
		"TF_CRED_CONTEXT_DIR",
		tfCredContextDir,
	); err != nil {
		panic(
			"failed to set TF_CRED_CONTEXT_DIR: " + err.Error(),
		)
	}

	code := m.Run()

	cleanupPaths()

	os.Exit(code)
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
