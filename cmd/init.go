package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"backup-tool/internal/config"
	"backup-tool/internal/ui"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize current directory for backup management",
	Long: `The init command creates project configuration for the backup system.
Creates .backup-config.json file with unique project ID
and sets up backup directory in %APPDATA%/ProjectBackup.`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to get current directory: %v", err)))
		return
	}

	configPath := filepath.Join(currentDir, config.ConfigFileName)
	if _, statErr := os.Stat(configPath); statErr == nil {
		fmt.Println(ui.Warning("Project already initialized in this directory"))
		fmt.Println()

		existingConfig, loadErr := config.LoadProjectConfig(currentDir)
		if loadErr != nil {
			fmt.Println(ui.Error(fmt.Sprintf("Failed to read existing configuration: %v", loadErr)))
			return
		}

		fmt.Println(ui.Label("Project Name", existingConfig.Name))
		fmt.Println(ui.Label("Project ID", existingConfig.ID))
		fmt.Println(ui.Label("Created", existingConfig.CreatedAt.Format("2006-01-02 15:04:05")))
		return
	}

	projectName := filepath.Base(currentDir)
	projectConfig := config.NewProjectConfig(projectName)

	backupPath, err := config.GetProjectBackupPath(projectConfig.ID)
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to create backup directory: %v", err)))
		return
	}

	projectConfig.BackupPath = backupPath

	if err := projectConfig.Save(currentDir); err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to save configuration: %v", err)))
		return
	}

	fmt.Println(ui.Success("Project successfully initialized!"))
	fmt.Println()
	fmt.Println(ui.Label("Project Name", projectConfig.Name))
	fmt.Println(ui.Label("Project ID", projectConfig.ID))
	fmt.Println(ui.Label("Backup Path", backupPath))
	fmt.Println()
	fmt.Println(ui.Hint("Now you can create your first backup with: backup create"))
}
