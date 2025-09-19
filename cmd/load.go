package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"backup-tool/internal/backup"
	"backup-tool/internal/config"
	"backup-tool/internal/ui"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load backup into current directory",
	Long: `The load command restores project from selected backup.
Without parameters opens interactive backup list.
With --name parameter loads specific backup by name.

WARNING: All files in current directory will be deleted!`,
	Run: runLoad,
}

var loadBackupName string

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringVarP(&loadBackupName, "name", "n", "", "Backup name to load")
}

func runLoad(cmd *cobra.Command, args []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to get current directory: %v", err)))
		return
	}

	projectConfig, err := config.LoadProjectConfig(currentDir)
	if err != nil {
		fmt.Println(ui.Error("Project not initialized. Run 'backup init' first."))
		return
	}

	backups, err := backup.LoadBackupMetadata(currentDir)
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to load backup list: %v", err)))
		return
	}

	if len(backups) == 0 {
		fmt.Println(ui.Info("No backups found. Create first backup with: backup create"))
		return
	}

	var selectedBackup *config.BackupMetadata

	if loadBackupName != "" {
		for _, b := range backups {
			if b.Name == loadBackupName {
				selectedBackup = b
				break
			}
		}
		if selectedBackup == nil {
			fmt.Println(ui.Error(fmt.Sprintf("Backup with name '%s' not found", loadBackupName)))
			return
		}
	} else {
		choice, err := ui.RunListUI(currentDir)
		if err != nil {
			fmt.Println(ui.Error(fmt.Sprintf("Error: %v", err)))
			return
		}

		if choice == "" || choice == "quit" {
			return
		}

		parts := strings.Split(choice, ":")
		if len(parts) != 2 || parts[0] != "load" {
			fmt.Println(ui.Error("Invalid choice"))
			return
		}

		index, err := strconv.Atoi(parts[1])
		if err != nil || index < 0 || index >= len(backups) {
			fmt.Println(ui.Error("Invalid backup index"))
			return
		}

		selectedBackup = backups[index]
	}

	displayName := selectedBackup.Name
	if displayName == "" {
		displayName = selectedBackup.CreatedAt.Format("2006-01-02 15:04:05")
	}

	fmt.Println(ui.Warning("WARNING: All files in current directory will be deleted!"))
	fmt.Println()
	fmt.Println(ui.Label("Backup to load", displayName))
	fmt.Println(ui.Label("Created", selectedBackup.CreatedAt.Format("2006-01-02 15:04:05")))
	fmt.Println(ui.Label("Size", fmt.Sprintf("%.2f MB", float64(selectedBackup.Size)/(1024*1024))))
	fmt.Println()

	fmt.Print(ui.ValueStyle.Render("Continue? (y/N): "))
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println(ui.Warning("Operation cancelled"))
		return
	}

	fmt.Println(ui.Info("Clearing current directory..."))
	if err := clearDirectory(currentDir, projectConfig.ID); err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to clear directory: %v", err)))
		return
	}

	fmt.Println(ui.Info("Restoring from backup..."))

	var bar *progressbar.ProgressBar

	err = backup.RestoreBackup(selectedBackup.FilePath, currentDir, func(progress backup.ArchiveProgress) {
		if bar == nil {
			bar = progressbar.NewOptions(progress.Total,
				progressbar.OptionSetDescription("Restoring"),
				progressbar.OptionSetWidth(50),
				progressbar.OptionShowCount(),
				progressbar.OptionShowIts(),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "█",
					SaucerPadding: "░",
					BarStart:      "▐",
					BarEnd:        "▌",
				}))
		}
		bar.Set(progress.Current)
		if progress.Current < progress.Total {
			bar.Describe(fmt.Sprintf("Restoring: %s", progress.File))
		}
	})

	if err != nil {
		fmt.Printf("\n%s\n", ui.Error(fmt.Sprintf("Restore failed: %v", err)))
		return
	}

	if bar != nil {
		bar.Finish()
	}

	fmt.Printf("\n%s\n", ui.Success("Backup successfully loaded!"))
	fmt.Println()
	fmt.Println(ui.Label("Restored backup", displayName))
	fmt.Println(ui.Label("Directory", currentDir))
}

func clearDirectory(dir string, projectID string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.Name() == config.ConfigFileName {
			continue
		}

		if entry.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}
