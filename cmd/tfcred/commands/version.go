// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of tfcred",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Println(cmd.Root().Version)
		},
	}
}
