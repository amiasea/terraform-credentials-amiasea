// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewPurgeCmd creates the purge command.
func NewPurgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge all contexts, bindings, and tokens",
		Run: func(cmd *cobra.Command, _ []string) {
			force, _ := cmd.Flags().GetBool("force")

			if !force {
				fmt.Println("[tfcred][warning] This will wipe ALL contexts, directory bindings, and vaulted tokens!")
				fmt.Print("Are you absolutely sure? [y/N]: ")
				var confirm string
				_, _ = fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" && confirm != "yes" {
					fmt.Println("[tfcred] Purge aborted.")
					return
				}
			}

			fmt.Println("[tfcred] Wiping all data...")
			store.PurgeAll()
			fmt.Println("[tfcred] SUCCESS: System entirely purged.")
		},
	}

	cmd.Flags().Bool("force", false, "skip confirmation prompt")
	return cmd
}
