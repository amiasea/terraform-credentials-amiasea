// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewInitCmd creates the init command.
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize tfcred storage and choose default domain",
		Run: func(cmd *cobra.Command, _ []string) {
			domain, _ := cmd.Flags().GetString("domain")
			if domain == "" {
				domain = promptDefaultDomain()
			}
			store.Init(domain)
			fmt.Printf("[tfcred] initialized with default domain %s\n", domain)
		},
	}

	cmd.Flags().String("domain", "", "default Terraform domain")
	return cmd
}
