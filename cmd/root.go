package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backup",
	Short: "CLI tool for creating and managing project backups",
	Long: `backup - CLI tool for creating, managing and restoring project backups.

Supported commands:
  init   - initialize directory for backups
  create - create new project backup
  list   - display list of all backups
  load   - load backup into current directory`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
