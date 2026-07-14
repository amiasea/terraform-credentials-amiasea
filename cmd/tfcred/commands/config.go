// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewConfigCmd creates the config command.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show or set configuration",
		Run: func(cmd *cobra.Command, _ []string) {
			defaultDomain, _ := cmd.Flags().GetString("default-domain")
			show, _ := cmd.Flags().GetBool("show")

			if show || defaultDomain == "" {
				f := store.Load()
				fmt.Printf("default_domain=%s\n", f.DefaultDomain)
				return
			}

			if !isSupportedDomain(defaultDomain) {
				fmt.Printf("[tfcred][error] unsupported domain: %s\n", defaultDomain)
				// Cobra will handle exit code via root
				return
			}

			store.SetDefaultDomain(defaultDomain)
			fmt.Printf("[tfcred] default domain set to %s\n", defaultDomain)
		},
	}

	cmd.Flags().String("default-domain", "", "set the default Terraform domain")
	cmd.Flags().Bool("show", false, "show current configuration")
	return cmd
}
