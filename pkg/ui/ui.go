package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	timestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))              // Grey
	eventStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))    // Blue/Cyan
	fileStyle      = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("245")) // Light Grey
	successStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))    // Green
)

func PrintFileChange(path string) {
	ts := timestampStyle.Render(time.Now().Format("15:04:05"))
	icon := eventStyle.Render("⚡ Reload")
	file := fileStyle.Render(path)

	fmt.Printf("%s  %s  %s\n", ts, icon, file)
}

func PrintSuccess(msg string) {
	ts := timestampStyle.Render(time.Now().Format("15:04:05"))
	icon := successStyle.Render("✔")
	fmt.Printf("%s  %s  %s\n", ts, icon, msg)
}
