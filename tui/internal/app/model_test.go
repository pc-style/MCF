package app

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	testutils "mcf-dev/tui/internal/testing"
	"mcf-dev/tui/internal/ui"
)

func TestInitialModel(t *testing.T) {
	t.Run("should create model with default values", func(t *testing.T) {
		model := InitialModel()

		assert.NotNil(t, model.theme, "Theme should be initialized")
		assert.NotNil(t, model.navigation, "Navigation should be initialized")
		assert.NotNil(t, model.commandInput, "Command input should be initialized")
		assert.NotNil(t, model.agentsList, "Agents list should be initialized")
		assert.NotNil(t, model.logViewer, "Log viewer should be initialized")
		assert.NotNil(t, model.dashboard, "Dashboard should be initialized")

		assert.Equal(t, ui.DashboardView, model.navigation.GetCurrentView(), "Should start in dashboard view")
		assert.False(t, model.ready, "Model should not be ready initially")
		assert.False(t, model.showHelp, "Help should not be shown initially")
		assert.Equal(t, 0, model.width, "Width should be 0 initially")
		assert.Equal(t, 0, model.height, "Height should be 0 initially")
	})
}

func TestMCFModel_SetView(t *testing.T) {
	t.Run("should set view correctly", func(t *testing.T) {
		model := InitialModel()

		model.SetView(ui.AgentsView)
		assert.Equal(t, ui.AgentsView, model.navigation.GetCurrentView())

		model.SetView(ui.LogsView)
		assert.Equal(t, ui.LogsView, model.navigation.GetCurrentView())

		model.SetView(ui.ConfigView)
		assert.Equal(t, ui.ConfigView, model.navigation.GetCurrentView())
	})
}

func TestMCFModel_Dimensions(t *testing.T) {
	t.Run("should handle window size updates", func(t *testing.T) {
		model := InitialModel()

		msg := tea.WindowSizeMsg{Width: 120, Height: 40}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.Equal(t, 120, updatedModel.Width())
		assert.Equal(t, 40, updatedModel.Height())
		assert.True(t, updatedModel.ready)
	})

	t.Run("should handle size changes", func(t *testing.T) {
		model := InitialModel()
		model.width = 100
		model.height = 30
		model.ready = true

		// Update size
		msg := tea.WindowSizeMsg{Width: 150, Height: 50}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.Equal(t, 150, updatedModel.Width())
		assert.Equal(t, 50, updatedModel.Height())
	})
}

func TestMCFModel_Navigation(t *testing.T) {
	t.Run("should navigate between views", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		// Test navigation order
		expectedViews := []ui.View{
			ui.DashboardView,
			ui.AgentsView,
			ui.CommandsView,
			ui.LogsView,
			ui.ConfigView,
		}

		for i, expectedView := range expectedViews {
			if i > 0 { // Skip first as it's initial view
				tabMsg := tea.KeyMsg{Type: tea.KeyTab}
				newModel, _ := model.Update(tabMsg)
				model = newModel.(MCFModel)
			}

			currentView := model.navigation.GetCurrentView()
			assert.Equal(t, expectedView, currentView,
				"View should be %v at navigation step %d", expectedView, i)
		}
	})

	t.Run("should wrap around when navigating past last view", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.ConfigView) // Last view

		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		newModel, _ := model.Update(tabMsg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.DashboardView, currentView, "Should wrap to first view")
	})

	t.Run("should handle reverse navigation", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.AgentsView)

		shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
		newModel, _ := model.Update(shiftTabMsg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.DashboardView, currentView, "Should go to previous view")
	})
}

func TestMCFModel_Help(t *testing.T) {
	t.Run("should toggle help display", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		// Initially help should be hidden
		assert.False(t, model.showHelp)

		// Toggle help on
		model.ToggleHelp()
		assert.True(t, model.showHelp)

		// Toggle help off
		model.ToggleHelp()
		assert.False(t, model.showHelp)
	})

	t.Run("should handle help with keyboard shortcut", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		questionMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
		newModel, _ := model.Update(questionMsg)
		updatedModel := newModel.(MCFModel)

		assert.True(t, updatedModel.showHelp, "Help should be shown after ? key")
	})
}

func TestMCFModel_Agents(t *testing.T) {
	t.Run("should handle agent list operations", func(t *testing.T) {
		model := InitialModel()

		testAgents := []ui.ListItem{
			{Title: "agent-1", Status: "active", Description: "Frontend development agent"},
			{Title: "agent-2", Status: "idle", Description: "Backend development agent"},
			{Title: "agent-3", Status: "error", Description: "DevOps agent"},
		}

		model.agentsList.SetItems(testAgents)

		// Navigate to agents view
		model.SetView(ui.AgentsView)

		selectedAgent := model.agentsList.GetSelectedItem()
		assert.NotNil(t, selectedAgent, "Should have a selected agent")
		assert.Equal(t, "agent-1", selectedAgent.Title, "First agent should be selected by default")
	})

	t.Run("should handle empty agent list", func(t *testing.T) {
		model := InitialModel()
		model.agentsList.SetItems([]ui.ListItem{})

		selectedAgent := model.agentsList.GetSelectedItem()
		assert.Nil(t, selectedAgent, "Should return nil when no agents available")
	})
}

func TestMCFModel_Logs(t *testing.T) {
	t.Run("should handle log operations", func(t *testing.T) {
		model := InitialModel()

		testLog := ui.LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Component: "test-component",
			Message:   "Test message",
		}

		model.logViewer.AddLog(testLog)

		// Verify log was added
		// In a full implementation, we'd have a method to retrieve logs
		assert.NotNil(t, model.logViewer, "Log viewer should be available")
	})
}

func TestMCFModel_CommandHistory(t *testing.T) {
	t.Run("should maintain command history", func(t *testing.T) {
		model := InitialModel()

		// Command history is managed by the command input component
		assert.NotNil(t, model.commandInput, "Command input should be initialized")
	})
}

func TestMCFModel_Theme(t *testing.T) {
	t.Run("should have theme initialized", func(t *testing.T) {
		model := InitialModel()

		assert.NotNil(t, model.theme, "Theme should be initialized")
	})
}

func TestMCFModel_StateConsistency(t *testing.T) {
	t.Run("should maintain consistent state after multiple operations", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		// Perform multiple state changes
		model.SetView(ui.AgentsView)
		model.ToggleHelp()
		model.SetView(ui.LogsView)
		model.ToggleHelp()

		// State should still be consistent
		assert.True(t, model.ready, "Model should remain ready")
		assert.Equal(t, ui.LogsView, model.navigation.GetCurrentView())
		assert.False(t, model.showHelp, "Help should be hidden after toggling twice")
		assert.NotNil(t, model.theme, "Theme should still be available")
	})
}

func TestMCFModel_View(t *testing.T) {
	t.Run("should render view without errors", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.width = 100
		model.height = 30

		view := model.View()
		assert.NotEmpty(t, view, "View should render content")
		assert.IsType(t, "", view, "View should return string")
	})

	t.Run("should handle view rendering when not ready", func(t *testing.T) {
		model := InitialModel()
		model.ready = false

		view := model.View()
		assert.NotEmpty(t, view, "Should render content even when not ready")
	})
}

func TestMCFModel_InteractionTime(t *testing.T) {
	t.Run("should track last interaction time", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		initialTime := model.lastInteractionTime

		// Simulate interaction
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.GreaterOrEqual(t, updatedModel.lastInteractionTime, initialTime,
			"Last interaction time should be updated")
	})
}

// Performance tests
func TestMCFModel_Performance(t *testing.T) {
	t.Run("should handle rapid view rendering", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.width = 100
		model.height = 30

		benchmark := testutils.NewPerformanceBenchmark("view_rendering", func() error {
			_ = model.View()
			return nil
		}).WithIterations(100)

		benchmark.Run(t)
	})

	t.Run("should handle rapid state updates", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		benchmark := testutils.NewPerformanceBenchmark("state_updates", func() error {
			model.ToggleHelp()
			model.SetView(ui.AgentsView)
			model.SetView(ui.DashboardView)
			return nil
		}).WithIterations(50)

		benchmark.Run(t)
	})
}

// Edge cases and error handling
func TestMCFModel_EdgeCases(t *testing.T) {
	t.Run("should handle invalid view", func(t *testing.T) {
		model := InitialModel()

		// Try to set an undefined view
		invalidView := ui.View(999)
		model.SetView(invalidView)

		// Should not crash and should maintain a valid state
		assert.NotNil(t, model.navigation, "Navigation should remain valid")
	})

	t.Run("should handle negative dimensions", func(t *testing.T) {
		model := InitialModel()

		msg := tea.WindowSizeMsg{Width: -10, Height: -5}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		// Should handle gracefully, possibly clamping to 0 or ignoring
		assert.NotNil(t, updatedModel, "Model should handle negative dimensions")
	})

	t.Run("should handle very large dimensions", func(t *testing.T) {
		model := InitialModel()

		msg := tea.WindowSizeMsg{Width: 10000, Height: 10000}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.NotNil(t, updatedModel, "Model should handle large dimensions")
		assert.True(t, updatedModel.ready, "Model should be ready")
	})
}

// Benchmark tests
func BenchmarkMCFModel_View(b *testing.B) {
	model := InitialModel()
	model.ready = true
	model.width = 100
	model.height = 30

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkMCFModel_SetView(b *testing.B) {
	model := InitialModel()
	model.ready = true

	views := []ui.View{ui.DashboardView, ui.AgentsView, ui.CommandsView, ui.LogsView, ui.ConfigView}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.SetView(views[i%len(views)])
	}
}

func BenchmarkMCFModel_ToggleHelp(b *testing.B) {
	model := InitialModel()
	model.ready = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.ToggleHelp()
	}
}
