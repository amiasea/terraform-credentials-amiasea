// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewSwitchCmd creates the switch command.
func NewSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <context>",
		Short: "Switch active context for current directory",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctxName := args[0]

			f := store.Load()
			if _, ok := f.Contexts[ctxName]; !ok {
				fmt.Printf("[tfcred][error] unknown context: %s\n", ctxName)
				os.Exit(1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("[tfcred][error] failed to get current directory: %v\n", err)
				os.Exit(1)
			}

			if err := store.BindDirectory(cwd, ctxName); err != nil {
				fmt.Printf("[tfcred][error] failed to bind directory: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("[tfcred] Context '%s' is now bound to: %s\n", ctxName, cwd)
		},
	}
}
