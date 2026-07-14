// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewAddCmd creates the add command.
func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new context",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx, _ := cmd.Flags().GetString("context")
			org, _ := cmd.Flags().GetString("org")
			tokenType, _ := cmd.Flags().GetString("token-type")
			domain, _ := cmd.Flags().GetString("domain")
			token, _ := cmd.Flags().GetString("token")
			shouldSwitch, _ := cmd.Flags().GetBool("switch")
			force, _ := cmd.Flags().GetBool("force")

			if ctx == "" {
				fmt.Println("[tfcred][error] --context is required")
				os.Exit(1)
			}

			// User tokens should not have org
			if tokenType == "user" && org != "" {
				fmt.Println("[tfcred][error] --org should not be specified for 'user' token types")
				os.Exit(1)
			}

			if tokenType != "user" && tokenType != "default" && org == "" {
				fmt.Println("[tfcred][error] --org is required for non-user contexts")
				os.Exit(1)
			}

			if domain != "" && !isSupportedDomain(domain) {
				fmt.Printf("[tfcred][error] unsupported domain: %s\n", domain)
				os.Exit(1)
			}

			config := store.Load()
			if domain == "" {
				if config.DefaultDomain == "" {
					fmt.Println("[tfcred][error] no default domain configured; pass --domain")
					os.Exit(1)
				}
				domain = config.DefaultDomain
			}

			if token != "" && !isValidTokenFormat(token) {
				fmt.Println("[tfcred][error] invalid token format")
				os.Exit(1)
			}

			// Duplicate / uniqueness checks
			for name, existing := range config.Contexts {
				if name == ctx {
					continue
				}

				if tokenType == "org" && existing.TokenType == "org" &&
					existing.Org == org && existing.Domain == domain {
					fmt.Printf("[tfcred][error] Duplicate org token for %s on %s\n", org, domain)
					os.Exit(1)
				}

				if tokenType != "team" && existing.Org == org &&
					existing.TokenType == tokenType && existing.Domain == domain {
					fmt.Printf("[tfcred][error] Duplicate mapping already exists under context: %s\n", name)
					os.Exit(1)
				}
			}

			if !force {
				// Overwrite confirmation
				if _, exists := config.Contexts[ctx]; exists {
					fmt.Printf("[tfcred][warning] context '%s' already exists. Overwrite? [y/N]: ", ctx)
					var confirm string
					_, _ = fmt.Scanln(&confirm)
					if !strings.EqualFold(confirm, "y") && !strings.EqualFold(confirm, "yes") {
						fmt.Println("[tfcred] Aborted.")
						return
					}
				}
			}

			store.Add(ctx, org, tokenType, domain, token)

			fmt.Printf("[tfcred] Context '%s' configured successfully.\n", ctx)

			if shouldSwitch {
				cwd, err := os.Getwd()
				if err == nil {
					_ = store.BindDirectory(cwd, ctx)
					fmt.Printf("[tfcred] Current directory bound to context '%s'\n", ctx)
				}
			}
		},
	}

	cmd.Flags().String("context", "", "context name")
	cmd.Flags().String("org", "", "organization")
	cmd.Flags().String("token-type", "user", "user|team|org")
	cmd.Flags().String("domain", "", "Terraform domain")
	cmd.Flags().String("token", "", "optional token")
	cmd.Flags().Bool("switch", false, "switch to this context after adding")
	cmd.Flags().Bool("force", false, "overwrite existing, bypass interaction")

	_ = cmd.MarkFlagRequired("context")

	return cmd
}
