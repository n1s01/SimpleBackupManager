package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type ProjectConfig struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	BackupPath string    `json:"backup_path"`
	Excludes   []string  `json:"excludes"`
}

type BackupMetadata struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	FilePath  string    `json:"file_path"`
}

type GlobalConfig struct {
	DefaultExcludes []string `json:"default_excludes"`
}

const (
	ConfigFileName = ".backup-config.json"
	AppDataDir     = "ProjectBackup"
)

func GetAppDataPath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("переменная окружения APPDATA не найдена")
	}

	backupDir := filepath.Join(appData, AppDataDir)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию бэкапов: %v", err)
	}

	return backupDir, nil
}

func GetProjectBackupPath(projectID string) (string, error) {
	appDataPath, err := GetAppDataPath()
	if err != nil {
		return "", err
	}

	projectPath := filepath.Join(appDataPath, projectID)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию проекта: %v", err)
	}

	return projectPath, nil
}

func LoadProjectConfig(dir string) (*ProjectConfig, error) {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфигурацию: %v", err)
	}

	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию: %v", err)
	}

	return &config, nil
}

func (c *ProjectConfig) Save(dir string) error {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать конфигурацию: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("не удалось сохранить конфигурацию: %v", err)
	}

	return nil
}

func NewProjectConfig(name string) *ProjectConfig {
	return &ProjectConfig{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		Excludes: []string{
			"node_modules/",
			".git/",
			"*.tmp",
			"*.log",
			"build/",
			"dist/",
			".env*",
			"*.exe",
			"*.dll",
		},
	}
}

func GetDefaultExcludes() []string {
	return []string{
		"node_modules/",
		".git/",
		"*.tmp",
		"*.log",
		"build/",
		"dist/",
		".env*",
		"*.exe",
		"*.dll",
		"target/",
		"bin/",
		"obj/",
		".vs/",
		".vscode/",
		"__pycache__/",
		"*.pyc",
	}
}
