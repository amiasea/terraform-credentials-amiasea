package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"[tfcred bootstrap] ❌ %v\n",
			err,
		)

		os.Exit(1)
	}
}

func run() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf(
			"failed to determine bootstrap location: %w",
			err,
		)
	}

	bootstrapDir := filepath.Dir(executable)

	scriptName := "install.ps1"

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "uninstall":
			scriptName = "uninstall.ps1"

		default:
			return fmt.Errorf(
				"unknown bootstrap argument: %s",
				os.Args[1],
			)
		}
	}

	scriptPath := filepath.Join(
		bootstrapDir,
		scriptName,
	)

	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf(
			"%s not found: %w",
			scriptName,
			err,
		)
	}

	fmt.Println("[tfcred bootstrap] Starting...")
	fmt.Printf(
		"[tfcred bootstrap] Version: %s\n",
		version,
	)
	fmt.Printf(
		"[tfcred bootstrap] Script: %s\n",
		scriptName,
	)

	args := []string{
		"-NoProfile",
		"-ExecutionPolicy",
		"Bypass",
		"-File",
		scriptPath,
	}

	if scriptName == "install.ps1" {
		args = append(
			args,
			"-Version",
			version,
		)
	}

	cmd := exec.Command(
		"powershell.exe",
		args...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"%s failed: %w",
			scriptName,
			err,
		)
	}

	fmt.Println(
		"[tfcred bootstrap] Completed successfully.",
	)

	return nil
}
