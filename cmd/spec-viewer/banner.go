package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintBanner(port, folder string) {
	primary := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))    // Indigo
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("250")) // Light Grey
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))    // White
	bold := lipgloss.NewStyle().Bold(true)

	asciiArt := `
  ___ ___  ___  ___   __   ___ _____      _____ ___ 
 / __| _ \/ _ \/ __|  \ \ / / | __\ \    / / __| _ \
 \__ \  _/  __/ (__    \ V /| | _| \ \/\/ /| _||   /
 |___/_|  \___|\___|    \_/ |_|___| \_/\_/ |___|_|_\
`
	fmt.Println(primary.Render(asciiArt))

	statusKey := secondary.Render("●  Server Running")
	urlKey := bold.Render("➜  Local:")
	urlValue := primary.Render(fmt.Sprintf("http://localhost:%s", port))
	folderKey := bold.Render("➜  Folder:")
	folderValue := secondary.Render(folder)
	quitMsg := secondary.Render("Press Ctrl+C to stop")

	content := fmt.Sprintf(`
 %s

 %s   %s
 %s   %s

 %s
`,
		accent.Render(statusKey),
		urlKey, urlValue,
		folderKey, folderValue,
		quitMsg,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Indigo border
		Padding(0, 1).
		Width(60).
		Render(content)

	fmt.Println(box)
	fmt.Println("")
}
