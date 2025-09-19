package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF88")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD93D")).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#74C0FC")).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ADB5BD"))

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	HintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#868E96")).
			Italic(true)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	ProgressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D9FF"))

	SecondaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7293"))
)

func Success(text string) string {
	return SuccessStyle.Render("[SUCCESS] " + text)
}

func Error(text string) string {
	return ErrorStyle.Render("[ERROR] " + text)
}

func Warning(text string) string {
	return WarningStyle.Render("[WARNING] " + text)
}

func Info(text string) string {
	return InfoStyle.Render("[INFO] " + text)
}

func Label(label, value string) string {
	return LabelStyle.Render(label+": ") + ValueStyle.Render(value)
}

func Hint(text string) string {
	return HintStyle.Render("TIP: " + text)
}

func Progress(text string) string {
	return ProgressStyle.Render("[PROGRESS] " + text)
}
