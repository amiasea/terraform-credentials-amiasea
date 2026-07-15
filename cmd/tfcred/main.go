// Package main
package main

import (
	"fmt"
	"os"

	"github.com/amiasea/terraform-credentials-tfcred/cmd/tfcred/commands"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	// ==============================================================================
	// Interactive Human/Developer Entry Points (Cobra Framework Execution)
	// ==============================================================================
	rootCmd := &cobra.Command{
		Use:     "tfcred",
		Short:   "Terraform credential helper and context manager",
		Long:    `tfcred makes it easy to switch between multiple Terraform Cloud / Enterprise tokens using named contexts.`,
		Version: version,
	}

	rootCmd.AddCommand(
		commands.NewVersionCmd(),
		commands.NewInitCmd(),
		commands.NewConfigCmd(),
		commands.NewAddCmd(),
		commands.NewRemoveCmd(),
		commands.NewPurgeCmd(),
		commands.NewListCmd(),
		commands.NewSwitchCmd(),
		commands.NewStatusCmd(),
		commands.NewCurrentCmd(),
		commands.NewContextCmd(),
    commands.NewCleanupCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
