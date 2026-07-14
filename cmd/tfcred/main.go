// Package main
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amiasea/terraform-credentials-tfcred/cmd/tfcred/commands"
	"github.com/amiasea/terraform-credentials-tfcred/internal/log"
	"github.com/amiasea/terraform-credentials-tfcred/internal/resolve"
	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
	"github.com/amiasea/terraform-credentials-tfcred/internal/tfcontext"

	"github.com/spf13/cobra"
)

var version = "dev"

type creds struct {
	Token string `json:"token"`
}

func main() {
	// ==============================================================================
	// 1. Intercept Terraform CLI Native Protocol Verbs (Before Cobra Triggers)
	// ==============================================================================
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "get":
			if len(os.Args) < 3 {
				log.Err("get command requires a hostname/domain (e.g. app.terraform.io)")
				os.Exit(1)
			}
			requestedDomain := os.Args[2]

			// 1. Resolve context from current working directory
			cwd, err := os.Getwd()
			if err != nil {
				log.Err(fmt.Sprintf("failed to get current working directory: %v", err))
				os.Exit(1)
			}

			contextKey, found := store.ResolveContextByDir(cwd)
			if !found {
				// Spec-Compliant Fallback: Directory is unbound, return empty token safely
				_ = json.NewEncoder(os.Stdout).Encode(creds{Token: ""})
				os.Exit(0)
			}

			// 2. Load context metadata registry
			storeFile := store.Load()
			entry, ok := storeFile.Contexts[contextKey]
			if !ok {
				log.Err(fmt.Sprintf("unknown context key '%s'", contextKey))
				os.Exit(1)
			}

			// 3. Fallback resolution rules for domain targets
			contextDomain := entry.Domain
			if contextDomain == "" {
				contextDomain = storeFile.DefaultDomain
			}
			if contextDomain == "" {
				contextDomain = "app.terraform.io"
			}

			// 4. Verify domain match
			if contextDomain != requestedDomain {
				// Spec-Compliant Fallback: Token is for a completely different host
				_ = json.NewEncoder(os.Stdout).Encode(creds{Token: ""})
				os.Exit(0)
			}

			// 5. Fetch secure secret string from Windows Vault via internal package
			activeCtx := tfcontext.Context{Key: contextKey}
			token, err := resolve.Resolve(activeCtx)

			fmt.Fprintln(
				os.Stderr,
				"resolved token:",
				token,
			)
			if err != nil {
				log.Err(err.Error())
				os.Exit(1)
			}

			// 6. Output pure JSON payload onto stdout for Terraform to parse
			_ = json.NewEncoder(os.Stdout).Encode(creds{Token: token})
			os.Exit(0)

		case "store", "forget":
			// Gracefully swallow background payload blocks to keep them read-only safe
			os.Exit(0)
		}
	}

	// ==============================================================================
	// 2. Interactive Human/Developer Entry Points (Cobra Framework Execution)
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
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
