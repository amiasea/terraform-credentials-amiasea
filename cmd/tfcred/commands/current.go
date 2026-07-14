// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewCurrentCmd creates the current command.
func NewCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Print the active context for current directory",
		Run: func(_ *cobra.Command, _ []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("[tfcred][error] %v\n", err)
				os.Exit(1)
			}

			contextKey, found := store.ResolveContextByDir(cwd)
			if !found {
				fmt.Println("[tfcred] Current directory is not bound to any context.")
				return
			}

			fmt.Printf("Active Directory: %s\n", cwd)
			fmt.Printf("active_context=%s\n", contextKey)
		},
	}
}
