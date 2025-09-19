package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"backup-tool/internal/config"
	"backup-tool/internal/ui"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display interactive list of all backups",
	Long: `The list command shows all created backups in interactive mode.
Allows to select backup and perform actions:
- Enter: load backup
- r: rename backup  
- d: delete backup
- q: quit`,
	Run: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to get current directory: %v", err)))
		return
	}

	if _, err := config.LoadProjectConfig(currentDir); err != nil {
		fmt.Println(ui.Error("Project not initialized. Run 'backup init' first."))
		return
	}

	choice, err := ui.RunListUI(currentDir)
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Error: %v", err)))
		return
	}

	if choice == "" || choice == "quit" {
		return
	}

	parts := strings.Split(choice, ":")
	if len(parts) != 2 {
		return
	}

	action := parts[0]
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Invalid index: %v", err)))
		return
	}

	switch action {
	case "load":
		fmt.Println(ui.Progress(fmt.Sprintf("Loading backup (index: %d)...", index)))
		fmt.Println(ui.Warning("Load function will be implemented in 'backup load' command"))
	case "rename":
		fmt.Println(ui.Progress(fmt.Sprintf("Renaming backup (index: %d)...", index)))
		fmt.Println(ui.Warning("Rename function will be added later"))
	case "delete":
		fmt.Println(ui.Progress(fmt.Sprintf("Deleting backup (index: %d)...", index)))
		fmt.Println(ui.Warning("Delete function will be added later"))
	}
}
