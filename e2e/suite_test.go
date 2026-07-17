package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	logEnvironment()

	if err := validateInstallation(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func logEnvironment() {
	fmt.Println("=== tfcred E2E Environment ===")
	fmt.Printf("tfcred install directory: %s\n", tfCredInstallDir)
	fmt.Printf("tfcred binary:            %s\n", tfCredBinary)
	fmt.Printf("tfcred data directory:    %s\n", tfCredDataDir)
	fmt.Printf("contexts file:            %s\n", tfCredContextsFile)
	fmt.Printf("terraform config:         %s\n", terraformConfigFile)
	fmt.Printf("terraform helper:         %s\n", tfHelperCredBinary)
	fmt.Printf(
		"process TF_CLI_CONFIG_FILE: %s\n",
		os.Getenv("TF_CLI_CONFIG_FILE"),
	)

	if terraformVersion, err := exec.Command(
		"terraform",
		"version",
	).Output(); err == nil {
		fmt.Printf(
			"terraform version:\n%s",
			terraformVersion,
		)
	}

	fmt.Println("============================")
}

func validateInstallation() error {
	required := map[string]string{
		"tfcred binary":         tfCredBinary,
		"terraform helper":      tfHelperCredBinary,
		"tfcred data directory": tfCredDataDir,
		"terraform config":      terraformConfigFile,
	}

	for name, path := range required {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf(
				"%s missing: %s",
				name,
				path,
			)
		}
	}

	return nil
}

func TestE2E(t *testing.T) {
	t.Run(
		"installation",
		func(t *testing.T) {
			runInstallationTests(t)
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
		"terraform_credentials_helper",
		func(t *testing.T) {
			runCredentialsHelperTests(t)
		},
	)
}
