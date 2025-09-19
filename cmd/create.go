package cmd

import (
	"fmt"
	"os"

	"backup-tool/internal/backup"
	"backup-tool/internal/config"
	"backup-tool/internal/ui"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new project backup",
	Long: `The create command archives current project into ZIP file.
Excludes standard folders (node_modules, .git, build, dist etc.)
and saves archive to backup directory.`,
	Run: runCreate,
}

var backupName string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&backupName, "name", "n", "", "Backup name (optional)")
}

func runCreate(cmd *cobra.Command, args []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to get current directory: %v", err)))
		return
	}

	if _, configErr := config.LoadProjectConfig(currentDir); configErr != nil {
		fmt.Println(ui.Error("Project not initialized. Run 'backup init' first."))
		return
	}

	fmt.Println(ui.Info("Preparing to create backup..."))

	var bar *progressbar.ProgressBar

	metadata, err := backup.CreateBackup(currentDir, backupName, func(progress backup.ArchiveProgress) {
		if bar == nil {
			bar = progressbar.NewOptions(progress.Total,
				progressbar.OptionSetDescription("Archiving"),
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
			bar.Describe(fmt.Sprintf("Archiving: %s", progress.File))
		}
	})

	if err != nil {
		fmt.Printf("\n%s\n", ui.Error(fmt.Sprintf("Failed to create backup: %v", err)))
		return
	}

	if bar != nil {
		bar.Finish()
	}

	fmt.Printf("\n%s\n", ui.Success("Backup successfully created!"))
	fmt.Println()
	fmt.Println(ui.Label("Name", getDisplayName(metadata)))
	fmt.Println(ui.Label("Size", fmt.Sprintf("%.2f MB", float64(metadata.Size)/(1024*1024))))
	fmt.Println(ui.Label("Created", metadata.CreatedAt.Format("2006-01-02 15:04:05")))
	fmt.Println(ui.Label("Path", metadata.FilePath))
}

func getDisplayName(metadata *config.BackupMetadata) string {
	if metadata.Name != "" {
		return metadata.Name
	}
	return metadata.CreatedAt.Format("2006-01-02 15:04:05")
}
