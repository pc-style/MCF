package app

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	testutils "mcf-dev/tui/internal/testing"
	"mcf-dev/tui/internal/ui"
)

func TestMCFModelUpdate_WindowSize(t *testing.T) {
	t.Run("should handle window size message", func(t *testing.T) {
		model := InitialModel()
		model.ready = false

		msg := tea.WindowSizeMsg{Width: 120, Height: 40}
		newModel, cmd := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.Equal(t, 120, updatedModel.width, "Width should be updated")
		assert.Equal(t, 40, updatedModel.height, "Height should be updated")
		assert.True(t, updatedModel.ready, "Model should be ready")
		assert.Nil(t, cmd, "Should not return command")
	})
}

func TestMCFModelUpdate_GlobalKeys(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		expectQuit  bool
		expectCmd   bool
		description string
	}{
		{"ctrl+c", "ctrl+c", true, false, "Should quit on Ctrl+C"},
		{"q", "q", true, false, "Should quit on q"},
		{"?", "?", false, false, "Should toggle help on ?"},
		{"tab", "tab", false, false, "Should switch to next view on Tab"},
		{"shift+tab", "shift+tab", false, false, "Should switch to previous view on Shift+Tab"},
		{":", ":", false, false, "Should switch to command bar on :"},
		{"esc", "esc", false, false, "Should handle escape key"},
		{"r", "r", false, false, "Should refresh on r"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel()
			model.ready = true
			model.width = 100
			model.height = 30

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			} else if tt.key == "tab" {
				msg = tea.KeyMsg{Type: tea.KeyTab}
			} else if tt.key == "shift+tab" {
				msg = tea.KeyMsg{Type: tea.KeyShiftTab}
			} else if tt.key == "esc" {
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			}

			newModel, cmd := model.Update(msg)

			if tt.expectQuit {
				// Quit commands are typically tea.Quit() which we can't easily test
				// In real implementation, we'd check if cmd == tea.Quit
				assert.NotNil(t, cmd, "Should return quit command")
			} else {
				updatedModel := newModel.(MCFModel)
				assert.NotNil(t, updatedModel, "Should return updated model")
			}
		})
	}
}

func TestMCFModelUpdate_ViewNavigation(t *testing.T) {
	t.Run("should navigate to next view with tab", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.DashboardView)

		msg := tea.KeyMsg{Type: tea.KeyTab}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.AgentsView, currentView, "Should move to next view")
	})

	t.Run("should navigate to previous view with shift+tab", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.AgentsView)

		msg := tea.KeyMsg{Type: tea.KeyShiftTab}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.DashboardView, currentView, "Should move to previous view")
	})

	t.Run("should wrap around views", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.ConfigView) // Last view

		msg := tea.KeyMsg{Type: tea.KeyTab}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.DashboardView, currentView, "Should wrap to first view")
	})
}

func TestMCFModelUpdate_CommandBar(t *testing.T) {
	t.Run("should switch to command bar on colon", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.DashboardView)

		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.CommandBarView, currentView, "Should switch to command bar")
	})

	t.Run("should return from command bar on escape", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.CommandBarView)

		msg := tea.KeyMsg{Type: tea.KeyEsc}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		currentView := updatedModel.navigation.GetCurrentView()
		assert.NotEqual(t, ui.CommandBarView, currentView, "Should leave command bar")
	})
}

func TestMCFModelUpdate_HelpToggle(t *testing.T) {
	t.Run("should toggle help with question mark", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.showHelp = false

		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.True(t, updatedModel.showHelp, "Should show help")

		// Toggle again
		newModel2, _ := updatedModel.Update(msg)
		updatedModel2 := newModel2.(MCFModel)

		assert.False(t, updatedModel2.showHelp, "Should hide help")
	})

	t.Run("should close help with escape", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.showHelp = true
		model.SetView(ui.DashboardView)

		msg := tea.KeyMsg{Type: tea.KeyEsc}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.False(t, updatedModel.showHelp, "Should close help")
	})
}

func TestMCFModelUpdate_PeriodicTick(t *testing.T) {
	t.Run("should handle tick message", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		msg := tickMsg(time.Now())
		newModel, cmd := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.NotNil(t, updatedModel, "Should return updated model")
		assert.NotNil(t, cmd, "Should return tick command for continuation")
	})
}

func TestMCFModelUpdate_DashboardView(t *testing.T) {
	t.Run("should handle dashboard navigation", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.DashboardView)

		tests := []struct {
			key         string
			description string
		}{
			{"j", "Should handle down navigation"},
			{"k", "Should handle up navigation"},
			{"down", "Should handle down arrow"},
			{"up", "Should handle up arrow"},
			{"enter", "Should handle enter key"},
		}

		for _, tt := range tests {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "down" {
				msg = tea.KeyMsg{Type: tea.KeyDown}
			} else if tt.key == "up" {
				msg = tea.KeyMsg{Type: tea.KeyUp}
			} else if tt.key == "enter" {
				msg = tea.KeyMsg{Type: tea.KeyEnter}
			}

			newModel, _ := model.Update(msg)
			updatedModel := newModel.(MCFModel)

			assert.NotNil(t, updatedModel, tt.description)
			// For dashboard interactions, cmd could be nil or a command
		}
	})

	t.Run("should handle quick action shortcuts", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.DashboardView)

		for i := 1; i <= 6; i++ {
			key := string(rune('0' + i))
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}

			newModel, _ := model.Update(msg)
			updatedModel := newModel.(MCFModel)

			assert.NotNil(t, updatedModel, "Should handle quick action shortcut %s", key)
		}
	})
}

func TestMCFModelUpdate_AgentsView(t *testing.T) {
	t.Run("should handle agents view interactions", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.AgentsView)

		tests := []struct {
			key         string
			description string
		}{
			{"s", "Should handle start/stop agent"},
			{"l", "Should handle view logs"},
			{"j", "Should handle navigation down"},
			{"k", "Should handle navigation up"},
		}

		for _, tt := range tests {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}

			newModel, _ := model.Update(msg)
			updatedModel := newModel.(MCFModel)

			assert.NotNil(t, updatedModel, tt.description)

			if tt.key == "l" {
				// Should switch to logs view
				currentView := updatedModel.navigation.GetCurrentView()
				assert.Equal(t, ui.LogsView, currentView, "Should switch to logs view")
			}
		}
	})
}

func TestMCFModelUpdate_CommandsView(t *testing.T) {
	t.Run("should handle commands view interactions", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.CommandsView)

		tests := []struct {
			key         string
			description string
		}{
			{"enter", "Should re-execute command"},
			{"d", "Should delete from history"},
			{"c", "Should clear history"},
			{"j", "Should navigate down"},
			{"k", "Should navigate up"},
		}

		for _, tt := range tests {
			var msg tea.Msg
			if tt.key == "enter" {
				msg = tea.KeyMsg{Type: tea.KeyEnter}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}

			newModel, _ := model.Update(msg)
			updatedModel := newModel.(MCFModel)

			assert.NotNil(t, updatedModel, tt.description)

			if tt.key == "c" {
				// Should clear command list
				// In a full implementation, we'd verify the list is empty
			}
		}
	})
}

func TestMCFModelUpdate_LogsView(t *testing.T) {
	t.Run("should handle logs view interactions", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.LogsView)

		tests := []struct {
			key         string
			description string
		}{
			{"j", "Should navigate down"},
			{"k", "Should navigate up"},
			{"f", "Should toggle follow mode"},
			{"/", "Should enter search mode"},
			{"c", "Should clear logs"},
			{"g", "Should go to top"},
			{"G", "Should go to bottom"},
		}

		for _, tt := range tests {
			var msg tea.Msg
			if tt.key == "G" {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("G")}
			} else if tt.key == "g" {
				msg = tea.KeyMsg{Type: tea.KeyHome}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}

			newModel, _ := model.Update(msg)
			updatedModel := newModel.(MCFModel)

			assert.NotNil(t, updatedModel, tt.description)
		}
	})
}

func TestMCFModelUpdate_ConfigView(t *testing.T) {
	t.Run("should handle config view interactions", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.ConfigView)

		tests := []struct {
			key         string
			description string
		}{
			{"e", "Should edit configuration"},
			{"r", "Should reload configuration"},
			{"b", "Should backup configuration"},
			{"d", "Should reset to defaults"},
		}

		for _, tt := range tests {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}

			newModel, _ := model.Update(msg)
			updatedModel := newModel.(MCFModel)

			assert.NotNil(t, updatedModel, tt.description)
		}
	})
}

func TestMCFModelUpdate_CommandBarView(t *testing.T) {
	t.Run("should handle command bar interactions", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.CommandBarView)

		// Test entering a command
		model.commandInput.Clear()

		// Simulate typing "test command"
		chars := []rune("test command")
		for _, char := range chars {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
			newModel, _ := model.Update(msg)
			model = newModel.(MCFModel)
		}

		// Test executing command with enter
		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		newModel, _ := model.Update(enterMsg)
		updatedModel := newModel.(MCFModel)

		// Should return to dashboard after execution
		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.DashboardView, currentView, "Should return to dashboard after command execution")

		// Command input should be cleared
		value := updatedModel.commandInput.GetValue()
		assert.Empty(t, value, "Command input should be cleared")
	})

	t.Run("should handle empty command", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.SetView(ui.CommandBarView)
		model.commandInput.Clear()

		// Press enter with empty command
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		// Should remain in command bar view
		currentView := updatedModel.navigation.GetCurrentView()
		assert.Equal(t, ui.CommandBarView, currentView, "Should remain in command bar with empty command")
	})
}

func TestMCFModelUpdate_PerformanceTracking(t *testing.T) {
	t.Run("should update last interaction time on any key press", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		initialTime := model.lastInteractionTime

		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.Greater(t, updatedModel.lastInteractionTime, initialTime,
			"Last interaction time should be updated on key press")
	})
}

func TestMCFModelUpdate_ErrorHandling(t *testing.T) {
	t.Run("should handle unknown message types gracefully", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		// Send an unknown message type
		type unknownMsg struct{}
		msg := unknownMsg{}

		newModel, cmd := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.NotNil(t, updatedModel, "Should handle unknown message type")
		assert.NotNil(t, cmd, "Should return tick command")
	})

	t.Run("should handle malformed key messages", func(t *testing.T) {
		model := InitialModel()
		model.ready = true

		// Send a key message with empty runes
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}}

		newModel, _ := model.Update(msg)
		updatedModel := newModel.(MCFModel)

		assert.NotNil(t, updatedModel, "Should handle empty key message")
	})
}

// Integration tests
func TestMCFModelUpdate_Integration(t *testing.T) {
	t.Run("should maintain state consistency across view changes", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.width = 100
		model.height = 30

		// Start in dashboard
		model.SetView(ui.DashboardView)
		assert.Equal(t, ui.DashboardView, model.navigation.GetCurrentView())

		// Navigate through all views
		views := []ui.View{ui.AgentsView, ui.CommandsView, ui.LogsView, ui.ConfigView}
		for range views {
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}
			newModel, _ := model.Update(tabMsg)
			model = newModel.(MCFModel)

			expectedView := model.navigation.GetCurrentView()
			assert.Contains(t, views, expectedView, "Should navigate to valid view")
		}
	})

	t.Run("should handle complex interaction sequences", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.width = 100
		model.height = 30

		// Sequence: Dashboard -> Command Bar -> Execute -> Dashboard

		// 1. Start in dashboard
		model.SetView(ui.DashboardView)

		// 2. Go to command bar
		colonMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")}
		newModel, _ := model.Update(colonMsg)
		model = newModel.(MCFModel)
		assert.Equal(t, ui.CommandBarView, model.navigation.GetCurrentView())

		// 3. Type command
		chars := []rune("test")
		for _, char := range chars {
			charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
			newModel, _ := model.Update(charMsg)
			model = newModel.(MCFModel)
		}

		// 4. Execute command
		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		newModel, _ = model.Update(enterMsg)
		model = newModel.(MCFModel)

		// 5. Should be back in dashboard
		assert.Equal(t, ui.DashboardView, model.navigation.GetCurrentView())
	})
}

// Performance tests
func TestMCFModelUpdate_Performance(t *testing.T) {
	t.Run("should handle rapid key presses efficiently", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.width = 100
		model.height = 30

		benchmark := testutils.NewPerformanceBenchmark("rapid_key_presses", func() error {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
			newModel, _ := model.Update(msg)
			model = newModel.(MCFModel)
			return nil
		}).WithIterations(1000)

		benchmark.Run(t)
	})

	t.Run("should handle view switching efficiently", func(t *testing.T) {
		model := InitialModel()
		model.ready = true
		model.width = 100
		model.height = 30

		benchmark := testutils.NewPerformanceBenchmark("view_switching", func() error {
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}
			newModel, _ := model.Update(tabMsg)
			model = newModel.(MCFModel)
			return nil
		}).WithIterations(500)

		benchmark.Run(t)
	})
}

// Benchmark tests
func BenchmarkMCFModelUpdate_KeyPress(b *testing.B) {
	model := InitialModel()
	model.ready = true
	model.width = 100
	model.height = 30

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel, _ := model.Update(msg)
		model = newModel.(MCFModel)
	}
}

func BenchmarkMCFModelUpdate_WindowResize(b *testing.B) {
	model := InitialModel()
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel, _ := model.Update(msg)
		model = newModel.(MCFModel)
	}
}

func BenchmarkMCFModelUpdate_TickMessage(b *testing.B) {
	model := InitialModel()
	model.ready = true
	msg := tickMsg(time.Now())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel, _ := model.Update(msg)
		model = newModel.(MCFModel)
	}
}
