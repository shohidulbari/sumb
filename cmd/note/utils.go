package note

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	ERROR_COLOR   = lipgloss.Color("#FF0000")
	WARN_COLOR    = lipgloss.Color("#FFFF00")
	SUCCESS_COLOR = lipgloss.Color("#00FF00")
	INFO_COLOR    = lipgloss.Color("#00FFFF")
)

var (
	Alert   = lipgloss.NewStyle().Foreground(ERROR_COLOR)
	Info    = lipgloss.NewStyle().Foreground(INFO_COLOR)
	Success = lipgloss.NewStyle().Foreground(SUCCESS_COLOR)
	Warn    = lipgloss.NewStyle().Foreground(WARN_COLOR)
)

func GetAlertWithUsageInfo(alertMsg string, usageInfo string) string {
	return fmt.Sprintf(
		"❌ %s ℹ️ %s",
		Alert.Render(alertMsg),
		Info.Render(usageInfo))
}
