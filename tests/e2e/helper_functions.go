package e2e

import (
	"archive/zip"
	"fmt"
	"io"
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

	tfHelperCredBinary = filepath.Join(
		tfPluginDir,
		"terraform-credentials-tfcred.exe",
	)

	tfInstallDir = filepath.Join(
		repoRoot(),
		"tests",
		"e2e",
		"install",
	)

	tfCredBinary = filepath.Join(
		tfInstallDir,
		"tfcred.exe",
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

	workspaceDir = "workspace"
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

func prepareE2EEnvironment() error {
	cleanupPaths()

	return extractTfcredPackage()
}

func cleanupPaths() {
	_ = os.RemoveAll(tfInstallDir)
	_ = os.RemoveAll(tfCredContextDir)
	_ = os.Remove(tfHelperCredBinary)
}

func extractTfcredPackage() error {
	archive := filepath.Join(
		repoRoot(),
		"dist",
		"terraform-credentials-tfcred_windows_amd64.zip",
	)

	if err := os.MkdirAll(tfInstallDir, 0o755); err != nil {
		return err
	}

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		target := filepath.Join(
			tfInstallDir,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		source, err := file.Open()
		if err != nil {
			return err
		}

		destination, err := os.OpenFile(
			target,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			0o755,
		)
		if err != nil {
			source.Close()
			return err
		}

		_, err = io.Copy(
			destination,
			source,
		)

		source.Close()
		destination.Close()

		if err != nil {
			return err
		}

		entries, err := os.ReadDir(tfInstallDir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fmt.Println("[e2e install]", entry.Name())
		}
	}

	return nil
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

	cmd := exec.Command(
		"terraform",
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

	if _, err := os.Stat(tfCredBinary); err != nil {
		return
	}

	_, _ = exec.Command(
		tfCredBinary,
		"purge",
		"--force",
	).CombinedOutput()
}
