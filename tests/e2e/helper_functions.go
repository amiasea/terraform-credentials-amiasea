package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	dummyToken = "FAKETOKENXXXXX.atlasv1.abcdefghijklmnopqrstuvwxyzABCDE"

	tfPluginDir = filepath.Join(
		os.Getenv("APPDATA"),
		"terraform.d",
		"plugins",
	)

	tfCredBinary = filepath.Join(
		tfPluginDir,
		"terraform-credentials-tfcred.exe",
	)

	tfConfigFile = filepath.Join(
		repoRoot(),
		"tests",
		"e2e",
		"terraform.tfrc.json",
	)

	tfCredContextDir = filepath.Join(
		repoRoot(),
		"tests",
		"e2e",
		"tfcred_context",
	)

	workspaceDir = filepath.Join(
		"workspace",
	)
)

func repoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(
		dir,
		"..",
		"..",
	)
}

func copyTfcred() error {
	source := filepath.Join(
		repoRoot(),
		"dist",
		"tfcred_windows_amd64_v1",
		"terraform-credentials-tfcred.exe",
	)

	return copyFile(
		source,
		tfCredBinary,
	)
}

func copyFile(
	source string,
	target string,
) error {
	input, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	return os.WriteFile(
		target,
		input,
		0o755,
	)
}

func runTfcred(
	t *testing.T,
	args ...string,
) string {
	t.Helper()

	cmd := exec.Command(
		tfCredBinary,
		args...,
	)

	cmd.Dir = workspaceDir

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

	terraform := "terraform"

	if cliPath := os.Getenv("TERRAFORM_CLI_PATH"); cliPath != "" {
		terraform = filepath.Join(cliPath, "terraform.exe")
	}

	cmd := exec.Command(
		terraform,
		tfArgs...,
	)

	output, err := cmd.CombinedOutput()

	t.Logf(
		"terraform output:\n%s",
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

func purgeTfcred(
	t *testing.T,
) {
	t.Helper()

	_, _ = exec.Command(
		tfCredBinary,
		"purge",
		"--force",
	).CombinedOutput()
}
