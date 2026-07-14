// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewListCmd creates the list command.
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured contexts",
		Run: func(_ *cobra.Command, _ []string) {
			store.List()
		},
	}
}
