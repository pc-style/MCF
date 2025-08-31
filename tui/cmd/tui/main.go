package main

import (
	"flag"
	"log"
	"os"

	"mcf-dev/tui/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	debugFlag := flag.Bool("debug", false, "Enable debug logging to stdout")
	logDirFlag := flag.String("log-dir", "", "Directory for log files (default: <mcf-root>/logs)")
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Set debug mode via environment variable if flag is set
	if *debugFlag {
		os.Setenv("MCF_TUI_DEBUG", "true")
	}
	if *logDirFlag != "" {
		os.Setenv("MCF_TUI_LOG_DIR", *logDirFlag)
	}

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
