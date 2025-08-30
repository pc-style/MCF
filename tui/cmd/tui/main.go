package main

import (
	"log"

	"mcf-dev/tui/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the TUI application
	model := app.InitialModel()

	// Create the program with alt screen and mouse support
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start the program (this is the correct way - Run() calls Start() internally)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
