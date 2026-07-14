// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewRemoveCmd creates the remove command.
func NewRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <context>",
		Short: "Remove a context and its bindings",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctxName := args[0]

			f := store.Load()
			if _, exists := f.Contexts[ctxName]; !exists {
				fmt.Printf("[tfcred][error] unknown context: %s\n", ctxName)
				os.Exit(1)
			}

			store.Remove(ctxName)
			fmt.Printf("[tfcred] Context '%s' and associated bindings removed successfully.\n", ctxName)
		},
	}
}
