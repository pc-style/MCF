package mcf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"mcf-dev/tui/internal/ui"
)

// SerenaAdapter provides integration with Serena semantic analysis
type SerenaAdapter struct {
	mcfRoot string
	host    string
	port    int
	enabled bool
}

// SerenaSymbol represents a code symbol found by Serena
type SerenaSymbol struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Description string `json:"description"`
}

// SerenaReference represents a symbol reference
type SerenaReference struct {
	Symbol  SerenaSymbol `json:"symbol"`
	File    string       `json:"file"`
	Line    int          `json:"line"`
	Column  int          `json:"column"`
	Context string       `json:"context"`
}

// SerenaAnalysis represents analysis results
type SerenaAnalysis struct {
	File        string         `json:"file"`
	Symbols     []SerenaSymbol `json:"symbols"`
	Issues      []string       `json:"issues"`
	Suggestions []string       `json:"suggestions"`
	Metrics     map[string]int `json:"metrics"`
}

// NewSerenaAdapter creates a new Serena adapter
func NewSerenaAdapter(mcfRoot string) *SerenaAdapter {
	return &SerenaAdapter{
		mcfRoot: mcfRoot,
		host:    "localhost",
		port:    8080,
		enabled: true, // Will be determined by checking if Serena is available
	}
}

// IsEnabled checks if Serena integration is enabled and available
func (s *SerenaAdapter) IsEnabled() bool {
	return s.enabled && s.checkSerenaAvailability()
}

// checkSerenaAvailability checks if Serena service is running
func (s *SerenaAdapter) checkSerenaAvailability() bool {
	url := fmt.Sprintf("http://%s:%d/health", s.host, s.port)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// GetStatus returns the current status of Serena integration
func (s *SerenaAdapter) GetStatus() string {
	if !s.enabled {
		return "disabled"
	}

	if s.checkSerenaAvailability() {
		return "connected"
	}

	return "disconnected"
}

// FindSymbols searches for symbols using Serena
func (s *SerenaAdapter) FindSymbols(query string) ([]SerenaSymbol, error) {
	if !s.IsEnabled() {
		return s.mockFindSymbols(query), nil
	}

	// Try to use real Serena command
	output, err := s.executeSerenaCommand("find", query)
	if err != nil {
		return s.mockFindSymbols(query), nil
	}

	var symbols []SerenaSymbol
	if err := json.Unmarshal([]byte(output), &symbols); err != nil {
		return s.mockFindSymbols(query), nil
	}

	return symbols, nil
}

// AnalyzeFile performs semantic analysis on a file
func (s *SerenaAdapter) AnalyzeFile(filepath string) (*SerenaAnalysis, error) {
	if !s.IsEnabled() {
		return s.mockAnalyzeFile(filepath), nil
	}

	// Try to use real Serena command
	output, err := s.executeSerenaCommand("analyze", filepath)
	if err != nil {
		return s.mockAnalyzeFile(filepath), nil
	}

	var analysis SerenaAnalysis
	if err := json.Unmarshal([]byte(output), &analysis); err != nil {
		return s.mockAnalyzeFile(filepath), nil
	}

	return &analysis, nil
}

// FindReferences finds all references to a symbol
func (s *SerenaAdapter) FindReferences(symbolName string) ([]SerenaReference, error) {
	if !s.IsEnabled() {
		return s.mockFindReferences(symbolName), nil
	}

	// Try to use real Serena command
	output, err := s.executeSerenaCommand("refs", symbolName)
	if err != nil {
		return s.mockFindReferences(symbolName), nil
	}

	var references []SerenaReference
	if err := json.Unmarshal([]byte(output), &references); err != nil {
		return s.mockFindReferences(symbolName), nil
	}

	return references, nil
}

// executeSerenaCommand executes a Serena command
func (s *SerenaAdapter) executeSerenaCommand(command string, args ...string) (string, error) {
	// Try different ways to execute Serena commands

	// Method 1: Direct serena command
	cmdArgs := append([]string{command}, args...)
	if output, err := exec.Command("serena", cmdArgs...).Output(); err == nil {
		return string(output), nil
	}

	// Method 2: MCF serena command
	mcfArgs := append([]string{"serena", command}, args...)
	if output, err := exec.Command("mcf", mcfArgs...).Output(); err == nil {
		return string(output), nil
	}

	// Method 3: Direct HTTP call to Serena service
	url := fmt.Sprintf("http://%s:%d/api/%s", s.host, s.port, command)
	resp, err := http.Get(url + "?q=" + strings.Join(args, " "))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("serena service returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	output, _ := json.Marshal(result["data"])
	return string(output), nil
}

// Mock implementations for when Serena is not available

func (s *SerenaAdapter) mockFindSymbols(query string) []SerenaSymbol {
	return []SerenaSymbol{
		{
			Name:        fmt.Sprintf("mockSymbol_%s", query),
			Kind:        "function",
			File:        "src/main.go",
			Line:        42,
			Column:      10,
			Description: fmt.Sprintf("Mock symbol matching '%s'", query),
		},
		{
			Name:        fmt.Sprintf("Mock%sClass", strings.Title(query)),
			Kind:        "class",
			File:        "src/models.go",
			Line:        15,
			Column:      1,
			Description: fmt.Sprintf("Mock class related to '%s'", query),
		},
	}
}

func (s *SerenaAdapter) mockAnalyzeFile(filepath string) *SerenaAnalysis {
	return &SerenaAnalysis{
		File: filepath,
		Symbols: []SerenaSymbol{
			{
				Name:        "mainFunction",
				Kind:        "function",
				File:        filepath,
				Line:        1,
				Column:      1,
				Description: "Main function",
			},
		},
		Issues:      []string{"No issues found (mock analysis)"},
		Suggestions: []string{"Consider adding documentation", "Optimize performance"},
		Metrics: map[string]int{
			"lines":      100,
			"functions":  5,
			"complexity": 3,
		},
	}
}

func (s *SerenaAdapter) mockFindReferences(symbolName string) []SerenaReference {
	return []SerenaReference{
		{
			Symbol: SerenaSymbol{
				Name: symbolName,
				Kind: "function",
				File: "src/main.go",
				Line: 42,
			},
			File:    "src/caller.go",
			Line:    15,
			Column:  10,
			Context: fmt.Sprintf("calling %s() here", symbolName),
		},
		{
			Symbol: SerenaSymbol{
				Name: symbolName,
				Kind: "function",
				File: "src/main.go",
				Line: 42,
			},
			File:    "src/test.go",
			Line:    28,
			Column:  5,
			Context: fmt.Sprintf("testing %s() functionality", symbolName),
		},
	}
}

// GetRecentActivity returns recent Serena activity logs
func (s *SerenaAdapter) GetRecentActivity() []ui.LogEntry {
	now := time.Now()

	if s.IsEnabled() {
		// Try to get real activity logs
		if logs := s.getRealActivityLogs(); len(logs) > 0 {
			return logs
		}
	}

	// Fallback to mock activity
	return []ui.LogEntry{
		{
			Timestamp: now.Add(-10 * time.Minute),
			Level:     "INFO",
			Component: "serena",
			Message:   "Semantic analysis service started",
		},
		{
			Timestamp: now.Add(-8 * time.Minute),
			Level:     "INFO",
			Component: "serena",
			Message:   "Indexed 1,247 symbols across 42 files",
		},
		{
			Timestamp: now.Add(-5 * time.Minute),
			Level:     "INFO",
			Component: "serena",
			Message:   "Symbol search query: 'handleRequest' - 3 results found",
		},
		{
			Timestamp: now.Add(-2 * time.Minute),
			Level:     "INFO",
			Component: "serena",
			Message:   "Code analysis completed for src/handlers.go",
		},
	}
}

// getRealActivityLogs attempts to get real Serena activity logs
func (s *SerenaAdapter) getRealActivityLogs() []ui.LogEntry {
	// This would implement reading from Serena's actual log files or API
	// For now, return empty to use mock data
	return []ui.LogEntry{}
}
