//go:build mage

package main

import (
	"fmt"
	"os/exec"
)

// BuildDev builds the tfcred binary for local development.
func BuildDev() error {
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags",
		"-X main.version=dev",
		"-o",
		"terraform-credentials-tfcred.exe",
		"./cmd/tfcred",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	fmt.Println("[mage] Building terraform-credentials-tfcred.exe...")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	fmt.Println("[mage] Build complete.")
	return nil
}

// InstallDev installs the locally built tfcred binary into the Terraform CLI environment.
func InstallDev() error {
	fmt.Println("[mage] Installing local tfcred binary...")

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-ExecutionPolicy",
		"Bypass",
		"-File",
		".\\scripts\\install.ps1",
		"-BinaryPath",
		".\\terraform-credentials-tfcred.exe",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install.ps1 failed: %w", err)
	}

	fmt.Println("[mage] Installation complete.")
	return nil
}

// Dev performs a local build and installation.
func Dev() error {
	if err := BuildDev(); err != nil {
		return err
	}

	return InstallDev()
}

// pluginDir := e2e.DefaultPluginDir()

// _, err := e2e.InstallTfcred(
// 	"terraform-credentials-tfcred.exe",
// 	pluginDir,
// )
