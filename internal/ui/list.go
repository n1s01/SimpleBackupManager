package ui

import (
	"fmt"
	"sort"
	"time"

	"backup-tool/internal/backup"
	"backup-tool/internal/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	backups  []*config.BackupMetadata
	cursor   int
	selected map[int]struct{}
	choice   string
	quitting bool
}

type action int

const (
	actionLoad action = iota
	actionRename
	actionDelete
	actionCancel
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD")).
			Bold(true).
			PaddingLeft(2)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#EE6FF8")).
				Bold(true)

	quitTextStyle = lipgloss.NewStyle().
			Margin(1, 0, 2, 4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

func initialModel(backups []*config.BackupMetadata) model {
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return model{
		backups:  backups,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			m.choice = "quit"
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.backups)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.backups) > 0 {
				m.choice = fmt.Sprintf("load:%d", m.cursor)
				return m, tea.Quit
			}

		case "r":
			if len(m.backups) > 0 {
				m.choice = fmt.Sprintf("rename:%d", m.cursor)
				return m, tea.Quit
			}

		case "d":
			if len(m.backups) > 0 {
				m.choice = fmt.Sprintf("delete:%d", m.cursor)
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Exiting...")
	}

	if len(m.backups) == 0 {
		return titleStyle.Render("Project Backups") + "\n\n" +
			"No backups found. Create first backup with: backup create\n"
	}

	s := titleStyle.Render("Project Backups") + "\n\n"

	for i, backup := range m.backups {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		displayName := backup.Name
		if displayName == "" {
			displayName = backup.CreatedAt.Format("2006-01-02 15:04:05")
		}

		size := formatSize(backup.Size)
		age := formatAge(backup.CreatedAt)

		line := fmt.Sprintf("%s %-30s | %-8s | %s",
			cursor, displayName, size, age)

		if m.cursor == i {
			s += selectedItemStyle.Render(line) + "\n"
		} else {
			s += itemStyle.Render(line) + "\n"
		}
	}

	s += "\n"
	s += helpStyle.Render("↑/↓: navigate • Enter: load • r: rename • d: delete • q: quit")

	return s
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatAge(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d min ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d h ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		} else if days < 7 {
			return fmt.Sprintf("%d days ago", days)
		} else {
			return t.Format("2006-01-02")
		}
	}
}

func RunListUI(projectPath string) (string, error) {
	backups, err := backup.LoadBackupMetadata(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to load backup list: %v", err)
	}

	p := tea.NewProgram(initialModel(backups))
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("UI error: %v", err)
	}

	if finalModel, ok := m.(model); ok {
		return finalModel.choice, nil
	}

	return "", nil
}
