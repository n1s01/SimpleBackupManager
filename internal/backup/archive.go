package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backup-tool/internal/config"
)

type ArchiveProgress struct {
	Current int
	Total   int
	File    string
}

func ShouldExclude(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		if strings.HasSuffix(pattern, "/") {
			if strings.Contains(path, pattern) {
				return true
			}
		} else if strings.Contains(pattern, "*") {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if matched {
				return true
			}
		} else if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func CountFiles(rootPath string, excludePatterns []string) (int, error) {
	count := 0
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(rootPath, path)
		if ShouldExclude(relPath, excludePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

func CreateBackup(projectPath string, backupName string, progressCallback func(ArchiveProgress)) (*config.BackupMetadata, error) {
	projectConfig, err := config.LoadProjectConfig(projectPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфигурацию проекта: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("backup_%s", timestamp)
	if backupName != "" {
		fileName = fmt.Sprintf("backup_%s_%s", timestamp, backupName)
	}
	fileName += ".zip"

	backupPath := filepath.Join(projectConfig.BackupPath, fileName)

	totalFiles, err := CountFiles(projectPath, projectConfig.Excludes)
	if err != nil {
		return nil, fmt.Errorf("не удалось подсчитать файлы: %v", err)
	}

	zipFile, err := os.Create(backupPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать архив: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	processedFiles := 0

	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		if ShouldExclude(relPath, projectConfig.Excludes) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if progressCallback != nil {
			progressCallback(ArchiveProgress{
				Current: processedFiles,
				Total:   totalFiles,
				File:    relPath,
			})
		}

		fileInArchive, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileOnDisk, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileOnDisk.Close()

		_, err = io.Copy(fileInArchive, fileOnDisk)
		if err != nil {
			return err
		}

		processedFiles++
		return nil
	})

	if err != nil {
		os.Remove(backupPath)
		return nil, fmt.Errorf("ошибка при создании архива: %v", err)
	}

	zipWriter.Close()
	zipFile.Close()

	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о файле: %v", err)
	}

	metadata := &config.BackupMetadata{
		ID:        fmt.Sprintf("%d", time.Now().Unix()),
		Name:      backupName,
		Size:      fileInfo.Size(),
		CreatedAt: time.Now(),
		FilePath:  backupPath,
	}

	return metadata, nil
}

func LoadBackupMetadata(projectPath string) ([]*config.BackupMetadata, error) {
	projectConfig, err := config.LoadProjectConfig(projectPath)
	if err != nil {
		return nil, err
	}

	var backups []*config.BackupMetadata

	err = filepath.Walk(projectConfig.BackupPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".zip" && strings.HasPrefix(info.Name(), "backup_") {
			backup := &config.BackupMetadata{
				ID:        fmt.Sprintf("%d", info.ModTime().Unix()),
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
				FilePath:  path,
			}

			fileName := strings.TrimSuffix(info.Name(), ".zip")
			parts := strings.Split(fileName, "_")
			if len(parts) > 2 {
				backup.Name = strings.Join(parts[2:], "_")
			}

			backups = append(backups, backup)
		}
		return nil
	})

	return backups, err
}

func RestoreBackup(backupPath string, targetPath string, progressCallback func(ArchiveProgress)) error {
	reader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть архив: %v", err)
	}
	defer reader.Close()

	totalFiles := len(reader.File)
	processedFiles := 0

	for _, file := range reader.File {
		if progressCallback != nil {
			progressCallback(ArchiveProgress{
				Current: processedFiles,
				Total:   totalFiles,
				File:    file.Name,
			})
		}

		path := filepath.Join(targetPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			fileReader.Close()
			return err
		}

		_, err = io.Copy(targetFile, fileReader)
		fileReader.Close()
		targetFile.Close()

		if err != nil {
			return err
		}

		processedFiles++
	}

	return nil
}
