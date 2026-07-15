// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the init command.
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize tfcred storage and Terraform integration",
		Long: `Initializes tfcred, installs the Terraform credentials helper, 
and configures Terraform to use it.`,
		Run: func(cmd *cobra.Command, _ []string) {
			domain, _ := cmd.Flags().GetString("domain")
			if domain == "" {
				domain = promptDefaultDomain()
			}

			if err := runInit(domain); err != nil {
				fmt.Fprintf(os.Stderr, "[tfcred] ❌ %v\n", err)
				os.Exit(1)
			}

			fmt.Println("[tfcred] ✅ Initialization completed successfully.")
		},
	}

	cmd.Flags().String("domain", "", "default Terraform domain")
	return cmd
}

func runInit(defaultDomain string) error {
	// 1. Initialize internal store
	store.Init(defaultDomain)

	// 2. Define paths (matching install.ps1)
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")
	userProfile := os.Getenv("USERPROFILE")

	terraformPluginDir := filepath.Join(appData, "terraform.d", "plugins")
	tfcredProgramDir := filepath.Join(localAppData, "Programs", "tfcred")
	contextsPath := filepath.Join(localAppData, "tfcred", "contexts.json")
	tfrcPath := filepath.Join(userProfile, "terraform.tfrc.json")

	// 3. Create required directories
	dirs := []string{terraformPluginDir, tfcredProgramDir, filepath.Dir(contextsPath)}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 4. Copy the helper binary to Terraform plugin directory
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}

	helperPath := filepath.Join(terraformPluginDir, "terraform-credentials-tfcred.exe")
	if err := copyFile(exe, helperPath); err != nil {
		return fmt.Errorf("failed to install Terraform credentials helper: %w", err)
	}

	// 5. Upsert terraform.tfrc.json with credentials_helper
	if err := upsertTerraformConfig(tfrcPath); err != nil {
		return fmt.Errorf("failed to configure terraform.tfrc.json: %w", err)
	}

	// 6. Set user-level TF_CLI_CONFIG_FILE
	if err := setUserEnvironmentVariable("TF_CLI_CONFIG_FILE", tfrcPath); err != nil {
		return fmt.Errorf("failed to set TF_CLI_CONFIG_FILE: %w", err)
	}

	fmt.Printf("[tfcred] ✅ Helper installed: %s\n", helperPath)
	fmt.Printf("[tfcred] ✅ Terraform config: %s\n", tfrcPath)
	fmt.Printf("[tfcred] ✅ Contexts stored at: %s\n", contextsPath)

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o755)
}

func upsertTerraformConfig(path string) error {
	var config map[string]interface{}

	if data, err := os.ReadFile(path); err == nil {
		// File exists → parse it
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("corrupted terraform.tfrc.json: %w", err)
		}
	} else {
		// Fresh file
		config = make(map[string]interface{})
	}

	// Add/upsert credentials_helper
	config["credentials_helper"] = map[string]interface{}{
		"tfcred": map[string]interface{}{},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func setUserEnvironmentVariable(name, value string) error {
	// Use setx for persistent user-level environment variable on Windows
	cmd := exec.Command("setx", name, value)
	return cmd.Run()
}
