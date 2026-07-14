// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

// NewContextCmd creates the context diagnostic command.
func NewContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Show active context resolution details",
		Run:   runContextCmd,
	}

	cmd.Flags().Bool(
		"json",
		false,
		"output in JSON format",
	)

	cmd.Flags().Bool(
		"all",
		false,
		"include all configured contexts",
	)

	return cmd
}

func runContextCmd(
	cmd *cobra.Command,
	_ []string,
) {
	jsonOut, _ := cmd.Flags().GetBool("json")
	all, _ := cmd.Flags().GetBool("all")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("[tfcred][error] failed to determine current directory:", err)
		os.Exit(1)
	}

	f := store.Load()

	contextKey, found := store.ResolveContextByDir(cwd)

	report := map[string]any{
		"working_directory": cwd,
	}

	if found {
		entry, exists := f.Contexts[contextKey]

		if !exists {
			fmt.Printf(
				"[tfcred][error] directory is bound to unknown context '%s'\n",
				contextKey,
			)
			os.Exit(1)
		}

		report["active_context"] = contextKey
		report["organization"] = entry.Org
		report["token_type"] = entry.TokenType
		report["domain"] = entry.Domain
		report["vault_key"] = store.TokenVaultKey(
			entry.Domain,
			entry.TokenType,
			entry.Org,
		)
	} else {
		report["active_context"] = ""
	}

	if all {
		report["contexts"] = buildContextSummary(f.Contexts)
	}

	if jsonOut {
		output, _ := json.MarshalIndent(
			report,
			"",
			" ",
		)

		fmt.Println(string(output))
		return
	}

	printContextReport(
		report,
	)
}

func buildContextSummary(
	contexts map[string]store.Entry,
) []map[string]string {
	names := make(
		[]string,
		0,
		len(contexts),
	)

	for name := range contexts {
		names = append(
			names,
			name,
		)
	}

	sort.Strings(names)

	results := make(
		[]map[string]string,
		0,
		len(names),
	)

	for _, name := range names {
		entry := contexts[name]

		results = append(
			results,
			map[string]string{
				"context":      name,
				"organization": entry.Org,
				"token_type":   entry.TokenType,
				"domain":       entry.Domain,
				"vault_key": store.TokenVaultKey(
					entry.Domain,
					entry.TokenType,
					entry.Org,
				),
			},
		)
	}

	return results
}

func printContextReport(
	report map[string]any,
) {
	keys := []string{
		"working_directory",
		"active_context",
		"organization",
		"token_type",
		"domain",
		"vault_key",
	}

	for _, key := range keys {
		if value, exists := report[key]; exists {
			fmt.Printf(
				"%s=%v\n",
				key,
				value,
			)
		}
	}

	if contexts, exists := report["contexts"]; exists {
		fmt.Println()
		fmt.Println("contexts:")

		for _, context := range contexts.([]map[string]string) {
			fmt.Printf(
				"- %s (%s/%s) %s\n",
				context["context"],
				context["token_type"],
				context["organization"],
				context["domain"],
			)

			fmt.Printf(
				"  vault_key=%s\n",
				context["vault_key"],
			)
		}
	}
}
