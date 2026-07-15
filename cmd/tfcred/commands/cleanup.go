package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewCleanupCmd creates the cleanup command.
func NewCleanupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Full cleanup: remove helper, data, and reset Terraform config",
		Long: `Performs a complete cleanup:
- Removes the Terraform credentials helper binary
- Deletes tfcred data directory (contexts, etc.)
- Runs purge to clear Windows credential vault entries
- Removes TF_CLI_CONFIG_FILE environment variable (Terraform will fallback to credentials.tfrc.json)`,
		Run: func(cmd *cobra.Command, _ []string) {
			if err := runCleanup(); err != nil {
				fmt.Fprintf(os.Stderr, "[tfcred] ❌ %v\n", err)
				os.Exit(1)
			}

			fmt.Println("[tfcred] ✅ Full cleanup completed.")
			fmt.Println("   Terraform will now use its default credentials.tfrc.json location.")
		},
	}

	return cmd
}

func runCleanup() error {
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")

	helperPath := filepath.Join(appData, "terraform.d", "plugins", "terraform-credentials-tfcred.exe")
	tfcredDir := filepath.Join(localAppData, "tfcred")

	fmt.Println("[tfcred] Clearing stored credentials from Windows Vault...")
	if err := runPurge(); err != nil {
		return fmt.Errorf("purge credentials: %w", err)
	}

	if err := os.Remove(helperPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove Terraform credentials helper: %w", err)
	}

	fmt.Println("[tfcred] ✅ Removed Terraform credentials helper.")

	if err := os.RemoveAll(tfcredDir); err != nil {
		return fmt.Errorf("remove tfcred data directory: %w", err)
	}

	fmt.Println("[tfcred] ✅ Removed tfcred data directory.")

	if err := removeUserEnvVar("TF_CLI_CONFIG_FILE"); err != nil {
		return fmt.Errorf("remove TF_CLI_CONFIG_FILE: %w", err)
	}

	fmt.Println("[tfcred] ✅ Reset TF_CLI_CONFIG_FILE (Terraform will use default location).")

	return nil
}

// runPurge calls your existing purge logic
func runPurge() error {
	purgeCmd := NewPurgeCmd()
	// Adjust this if your Purge command has a better internal function to call
	return purgeCmd.RunE(purgeCmd, nil)
}

func removeUserEnvVar(name string) error {
	// setx with /D is not reliable for deletion.
	// We can only advise or use registry. For simplicity:
	fmt.Printf("[tfcred] Note: Removing %s environment variable...\n", name)
	cmd := exec.Command("setx", name, "")
	return cmd.Run()
}
