// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewStatusCmd creates the status command.
func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current context status",
		Run: func(_ *cobra.Command, _ []string) {
			cwd, _ := os.Getwd()
			contextKey, found := store.ResolveContextByDir(cwd)
			if !found {
				fmt.Println("[tfcred] No active context bound to this directory.")
				return
			}

			f := store.Load()
			entry, ok := f.Contexts[contextKey]
			if !ok {
				fmt.Printf("[tfcred][error] unknown context key %s\n", contextKey)
				os.Exit(1)
			}

			fmt.Printf("[tfcred] context=%s type=%s org=%s domain=%s\n",
				contextKey, entry.TokenType, entry.Org, entry.Domain)
		},
	}
}
