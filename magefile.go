//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	e2eOutputDir = "dist-dev\\e2e"
)

func repoRoot() string {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return root
}

func e2eInstallerPath() string {
	return filepath.Join(
		repoRoot(),
		e2eOutputDir,
		"tfcred-bootstrap.exe",
	)
}

// BuildE2E builds the local tfcred-bootstrap installer package used for end-to-end testing.
func BuildE2E() error {
	fmt.Println("[mage] Building E2E installer package...")

	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-ExecutionPolicy",
		"Bypass",
		"-File",
		".\\scripts\\build-e2e.ps1",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"build-e2e.ps1 failed: %w",
			err,
		)
	}

	fmt.Println("[mage] E2E build complete.")

	return nil
}

// InstallE2E installs the locally built tfcred package using tfcred-bootstrap.exe.
func InstallE2E() error {
	installer := e2eInstallerPath()

	if _, err := os.Stat(installer); err != nil {
		return fmt.Errorf(
			"E2E installer not found: %w (run mage BuildE2E first)",
			err,
		)
	}

	fmt.Println("[mage] Installing E2E package...")

	cmd := exec.Command(
		installer,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"tfcred-bootstrap.exe failed: %w",
			err,
		)
	}

	fmt.Println("[mage] E2E installation complete.")

	return nil
}

// UninstallE2E uninstalls the locally installed tfcred package using tfcred-bootstrap.exe.
func UninstallE2E() error {
	installer := e2eInstallerPath()

	if _, err := os.Stat(installer); err != nil {
		return fmt.Errorf(
			"E2E installer not found: %w (run mage BuildE2E first)",
			err,
		)
	}

	fmt.Println("[mage] Uninstalling E2E package...")

	cmd := exec.Command(
		installer,
		"uninstall",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"tfcred-bootstrap.exe uninstall failed: %w",
			err,
		)
	}

	fmt.Println("[mage] E2E uninstall complete.")

	return nil
}

// E2E prepares and installs the local package for E2E testing.
func E2E() error {
	if err := BuildE2E(); err != nil {
		return err
	}

	return InstallE2E()
}
