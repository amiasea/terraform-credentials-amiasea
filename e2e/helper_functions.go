package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

var (
	dummyToken = "FAKETOKENXXXXX.atlasv1.abcdefghijklmnopqrstuvwxyzABCDE"

	localAppData = os.Getenv("LOCALAPPDATA")

	tfCredInstallDir = filepath.Join(
		localAppData,
		"Programs",
		"amiasea",
		"tfcred",
	)

	tfCredBinary = filepath.Join(
		tfCredInstallDir,
		"tfcred.exe",
	)

	tfCredDataDir = filepath.Join(
		localAppData,
		"amiasea",
		"tfcred",
	)

	tfCredContextsFile = filepath.Join(
		tfCredDataDir,
		"tfcred_contexts.json",
	)

	terraformConfigFile = filepath.Join(
		tfCredDataDir,
		"terraform.tfrc",
	)

	tfHelperCredBinary = filepath.Join(
		os.Getenv("APPDATA"),
		"terraform.d",
		"plugins",
		"terraform-credentials-tfcred.exe",
	)

	workspaceDir = "workspace"
)

func runTfcred(
	t *testing.T,
	args ...string,
) string {
	t.Helper()

	cmd := exec.Command(
		tfCredBinary,
		args...,
	)

	output, err := cmd.CombinedOutput()

	t.Logf(
		"tfcred %s output:\n%s",
		strings.Join(args, " "),
		output,
	)

	if err != nil {
		t.Fatalf(
			"tfcred %v failed:\n%s\n%v",
			args,
			output,
			err,
		)
	}

	return string(output)
}

func runTfcredAt(
	t *testing.T,
	workspace string,
	args ...string,
) string {
	t.Helper()

	cmd := exec.Command(
		tfCredBinary,
		args...,
	)

	if workspace != "" {
		cmd.Dir = workspace
	}

	output, err := cmd.CombinedOutput()

	t.Logf(
		"tfcred %s output:\n%s",
		strings.Join(args, " "),
		output,
	)

	if err != nil {
		t.Fatalf(
			"tfcred %v failed:\n%s\n%v",
			args,
			output,
			err,
		)
	}

	return string(output)
}

func runTerraform(
	t *testing.T,
	args ...string,
) {
	t.Helper()

	tfArgs := []string{
		fmt.Sprintf(
			"-chdir=%s",
			workspaceDir,
		),
	}

	tfArgs = append(
		tfArgs,
		args...,
	)

	cmd := exec.Command(
		"terraform",
		tfArgs...,
	)

	output, err := cmd.CombinedOutput()

	t.Logf(
		"terraform %s output:\n%s",
		strings.Join(args, " "),
		output,
	)

	if err != nil {
		t.Fatalf(
			"terraform %v failed:\n%s\n%v",
			args,
			output,
			err,
		)
	}
}

func assertContains(
	t *testing.T,
	output string,
	expected string,
) {
	t.Helper()

	if !strings.Contains(output, expected) {
		t.Fatalf(
			"expected output to contain:\n%s\n\nactual:\n%s",
			expected,
			output,
		)
	}
}

func assertNotContains(
	t *testing.T,
	output string,
	unexpected string,
) {
	t.Helper()

	if strings.Contains(output, unexpected) {
		t.Fatalf(
			"expected output NOT to contain:\n%s\n\nactual:\n%s",
			unexpected,
			output,
		)
	}
}

func purgeTfcredAfterTest(
	t *testing.T,
) {
	t.Helper()

	if t.Failed() {
		dumpTfcredEnvironment(t)
		dumpVaultState(t)
	}

	purgeTfcred(t)
}

func dumpTfcredEnvironment(
	t *testing.T,
) {
	t.Helper()

	t.Log("=== tfcred E2E Failure Environment Dump ===")

	t.Logf(
		"tfcred install directory: %s",
		tfCredInstallDir,
	)

	t.Logf(
		"tfcred binary: %s",
		tfCredBinary,
	)

	t.Logf(
		"tfcred data directory: %s",
		tfCredDataDir,
	)

	t.Logf(
		"context file: %s",
		tfCredContextsFile,
	)

	t.Logf(
		"terraform config: %s",
		terraformConfigFile,
	)

	t.Logf(
		"terraform helper: %s",
		tfHelperCredBinary,
	)

	t.Logf(
		"process TF_CLI_CONFIG_FILE: %s",
		os.Getenv("TF_CLI_CONFIG_FILE"),
	)

	t.Logf(
		"process TF_CONTEXT: %s",
		os.Getenv("TF_CONTEXT"),
	)

	if contexts, err := os.ReadFile(tfCredContextsFile); err != nil {
		t.Logf(
			"context file read failed: %v",
			err,
		)
	} else {
		t.Logf(
			"context file contents:\n%s",
			contexts,
		)
	}

	output := runTfcred(
		t,
		"current",
	)

	t.Logf(
		"tfcred current:\n%s",
		output,
	)

	output = runTfcred(
		t,
		"context",
	)

	t.Logf(
		"tfcred context:\n%s",
		output,
	)

	if config, err := os.ReadFile(terraformConfigFile); err != nil {
		t.Logf(
			"terraform config read failed: %v",
			err,
		)
	} else {
		t.Logf(
			"terraform config contents:\n%s",
			config,
		)
	}

	t.Log("==========================================")
}

func purgeTfcred(
	t *testing.T,
) {
	t.Helper()

	cmd := exec.Command(
		tfCredBinary,
		"purge",
		"--force",
	)

	output, err := cmd.CombinedOutput()

	t.Logf(
		"tfcred purge output:\n%s",
		output,
	)

	if err != nil {
		t.Fatalf(
			"tfcred purge failed:\n%s\n%v",
			output,
			err,
		)
	}
}

func dumpVaultState(
	t *testing.T,
) {
	t.Helper()

	t.Log("=== tfcred Vault State Dump ===")

	data, err := os.ReadFile(tfCredContextsFile)
	if err != nil {
		t.Logf(
			"context file unavailable: %v",
			err,
		)
		return
	}

	var contexts struct {
		Contexts map[string]struct {
			Org       string `json:"org"`
			TokenType string `json:"tokenType"`
			Domain    string `json:"domain"`
		} `json:"contexts"`
	}

	if err := json.Unmarshal(data, &contexts); err != nil {
		t.Logf(
			"failed parsing context file: %v",
			err,
		)
		return
	}

	for name, context := range contexts.Contexts {
		vaultKey := fmt.Sprintf(
			"tfcred:domain:%s:%s:%s",
			strings.ReplaceAll(
				context.Domain,
				".",
				"_",
			),
			context.TokenType,
			context.Org,
		)

		token, err := keyring.Get(
			"tfcred",
			vaultKey,
		)

		switch {
		case err != nil:
			t.Logf(
				"context=%s vault_key=%s status=missing error=%v",
				name,
				vaultKey,
				err,
			)

		default:
			t.Logf(
				"context=%s vault_key=%s status=present token_length=%d",
				name,
				vaultKey,
				len(token),
			)
		}
	}

	t.Log("==============================")
}
