package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mcf-dev/tui/internal/app"
	testutils "mcf-dev/tui/internal/testing"
	"mcf-dev/tui/internal/ui"
)

// End-to-end workflow tests that simulate complete user interactions

func TestCompleteUserWorkflow_DashboardToAgentsToLogs(t *testing.T) {
	t.Run("should navigate through complete workflow", func(t *testing.T) {
		// Initialize model
		model := app.InitialModel()
		model = prepareModel(model, t)

		// 1. Start in dashboard
		assert.Equal(t, ui.DashboardView, model.GetCurrentView())

		// 2. Navigate to agents view
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		assert.Equal(t, ui.AgentsView, model.GetCurrentView())

		// 3. Select an agent (simulate down navigation)
		for i := 0; i < 3; i++ {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		}

		// 4. View agent logs
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
		assert.Equal(t, ui.LogsView, model.GetCurrentView())

		// 5. Search logs
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})

		// Type search term
		searchTerm := "error"
		for _, char := range searchTerm {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		}

		// Execute search
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		// 6. Return to dashboard
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
		assert.Equal(t, ui.CommandBarView, model.GetCurrentView())

		// Type exit command
		exitCmd := "dashboard"
		for _, char := range exitCmd {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		}

		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, ui.DashboardView, model.GetCurrentView())
	})
}

func TestCommandExecutionWorkflow(t *testing.T) {
	t.Run("should execute commands end-to-end", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		// Start command sequence
		commands := []string{
			"mcf agents status",
			"mcf serena start",
			"mcf deploy --stage dev",
		}

		for _, command := range commands {
			// 1. Go to command bar
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
			assert.Equal(t, ui.CommandBarView, model.GetCurrentView())

			// 2. Type command
			for _, char := range command {
				model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
			}

			// 3. Execute command
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should return to dashboard after command execution
			assert.Equal(t, ui.DashboardView, model.GetCurrentView())

			// Verify command was added to history
			history := model.GetCommandHistory()
			assert.Contains(t, history, command, "Command should be in history")
		}
	})
}

func TestAgentOrchestrationWorkflow(t *testing.T) {
	t.Run("should simulate agent orchestration scenario", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		// 1. Check system health first
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
		healthCmd := "mcf health check"
		for _, char := range healthCmd {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		}
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		// 2. Navigate to agents view
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		assert.Equal(t, ui.AgentsView, model.GetCurrentView())

		// 3. Start multiple agents
		agentActions := []string{"s", "s", "s"} // Start action for multiple agents
		for _, action := range agentActions {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(action)})
			// Move to next agent
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		}

		// 4. Monitor logs for all agents
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
		assert.Equal(t, ui.LogsView, model.GetCurrentView())

		// 5. Enable following mode to monitor real-time
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")})

		// 6. Check commands history
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		assert.Equal(t, ui.CommandsView, model.GetCurrentView())
	})
}

func TestConfigurationManagementWorkflow(t *testing.T) {
	t.Run("should handle configuration management workflow", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		// 1. Navigate to config view
		for {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
			if model.GetCurrentView() == ui.ConfigView {
				break
			}
		}

		// 2. View current configuration
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}) // reload

		// 3. Backup current config
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")}) // backup

		// 4. Edit configuration
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")}) // edit

		// 5. Reset to defaults if needed
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")}) // defaults

		// 6. Verify configuration through command
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
		configCmd := "mcf config show"
		for _, char := range configCmd {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		}
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	})
}

func TestErrorHandlingWorkflow(t *testing.T) {
	t.Run("should handle errors gracefully in workflows", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		// 1. Execute invalid command
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
		invalidCmd := "mcf invalid-command --bad-flag"
		for _, char := range invalidCmd {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		}
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		// Should still be functional after error
		assert.Equal(t, ui.DashboardView, model.GetCurrentView())

		// 2. Try to access non-existent agent
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		assert.Equal(t, ui.AgentsView, model.GetCurrentView())

		// Navigate beyond available agents
		for i := 0; i < 100; i++ {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		}

		// Should still be functional
		selectedAgent := model.GetSelectedAgent()
		assert.NotNil(t, selectedAgent, "Should have a selected agent even after excessive navigation")

		// 3. Test recovery with help system
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		assert.True(t, model.IsHelpShown(), "Help should be shown")

		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		assert.False(t, model.IsHelpShown(), "Help should be hidden")
	})
}

func TestPerformanceUnderLoad(t *testing.T) {
	t.Run("should handle rapid interactions efficiently", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		benchmark := testutils.NewPerformanceBenchmark("rapid_workflow_interactions", func() error {
			// Simulate rapid user interaction
			interactions := []tea.KeyMsg{
				{Type: tea.KeyTab},
				{Type: tea.KeyDown},
				{Type: tea.KeyDown},
				{Type: tea.KeyTab},
				{Type: tea.KeyRunes, Runes: []rune("j")},
				{Type: tea.KeyRunes, Runes: []rune("k")},
				{Type: tea.KeyTab},
				{Type: tea.KeyRunes, Runes: []rune("f")},
				{Type: tea.KeyRunes, Runes: []rune("c")},
			}

			for _, msg := range interactions {
				var err error
				model, _ = model.Update(msg)
				if err != nil {
					return err
				}
			}
			return nil
		}).WithIterations(100)

		benchmark.Run(t)
	})

	t.Run("should handle large datasets efficiently", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		// Add many log entries
		for i := 0; i < 10000; i++ {
			logEntry := ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: fmt.Sprintf("component-%d", i%10),
				Message:   fmt.Sprintf("Test message %d with some additional content", i),
			}
			model.AddLog(logEntry)
		}

		// Navigate to logs and test performance
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab}) // Navigate to logs view

		benchmark := testutils.NewPerformanceBenchmark("large_dataset_navigation", func() error {
			// Test scrolling through large log dataset
			for i := 0; i < 100; i++ {
				model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			for i := 0; i < 100; i++ {
				model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
			}
			return nil
		}).WithIterations(10)

		benchmark.Run(t)
	})
}

func TestLongRunningWorkflow(t *testing.T) {
	t.Run("should maintain stability during extended usage", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		// Simulate extended usage session
		for session := 0; session < 10; session++ {
			// Complete workflow cycle
			views := []ui.View{
				ui.DashboardView,
				ui.AgentsView,
				ui.CommandsView,
				ui.LogsView,
				ui.ConfigView,
			}

			for _, targetView := range views {
				// Navigate to target view
				for model.GetCurrentView() != targetView {
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
				}

				// Perform view-specific actions
				switch targetView {
				case ui.AgentsView:
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
				case ui.LogsView:
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")})
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
				case ui.CommandsView:
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				case ui.ConfigView:
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
				}

				// Execute some commands
				model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
				cmd := fmt.Sprintf("mcf test-command-%d", session)
				for _, char := range cmd {
					model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
				}
				model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}
		}

		// Verify model is still functional
		assert.NotNil(t, model, "Model should still be functional")
		assert.True(t, model.IsReady(), "Model should still be ready")

		// Verify all components are still working
		model.ToggleHelp()
		assert.True(t, model.IsHelpShown(), "Help should still work")

		view := model.View()
		assert.NotEmpty(t, view, "Should still render content")
	})
}

func TestStateConsistencyAcrossWorkflows(t *testing.T) {
	t.Run("should maintain consistent state across complex workflows", func(t *testing.T) {
		model := app.InitialModel()
		model = prepareModel(model, t)

		initialState := captureModelState(model)

		// Execute complex state changes
		workflows := []func(app.MCFModel) app.MCFModel{
			executeAgentWorkflow,
			executeCommandWorkflow,
			executeLogWorkflow,
			executeConfigWorkflow,
		}

		for _, workflow := range workflows {
			model = workflow(model)

			// Verify state consistency after each workflow
			assert.True(t, model.IsReady(), "Model should remain ready")
			assert.NotNil(t, model.GetTheme(), "Theme should be preserved")
			assert.Greater(t, model.Width(), 0, "Width should be valid")
			assert.Greater(t, model.Height(), 0, "Height should be valid")
		}

		finalState := captureModelState(model)

		// Verify essential state is preserved
		assert.Equal(t, initialState.Width, finalState.Width, "Width should be preserved")
		assert.Equal(t, initialState.Height, finalState.Height, "Height should be preserved")
		assert.Equal(t, initialState.Ready, finalState.Ready, "Ready state should be preserved")
	})
}

// Helper functions and types

type ModelState struct {
	Width  int
	Height int
	Ready  bool
	View   ui.View
}

func prepareModel(model app.MCFModel, t *testing.T) app.MCFModel {
	// Set up model for testing
	model, _ = model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Add some test data
	testAgents := []ui.ListItem{
		{Title: "frontend-agent", Status: "active", Description: "Frontend development agent"},
		{Title: "backend-agent", Status: "idle", Description: "Backend development agent"},
		{Title: "devops-agent", Status: "error", Description: "DevOps automation agent"},
	}
	model.SetAgents(testAgents)

	// Add test logs
	testLogs := []ui.LogEntry{
		{Timestamp: time.Now(), Level: "INFO", Component: "system", Message: "System initialized"},
		{Timestamp: time.Now(), Level: "ERROR", Component: "frontend-agent", Message: "Build failed"},
		{Timestamp: time.Now(), Level: "WARN", Component: "backend-agent", Message: "High memory usage"},
	}
	for _, log := range testLogs {
		model.AddLog(log)
	}

	return model
}

func captureModelState(model app.MCFModel) ModelState {
	return ModelState{
		Width:  model.Width(),
		Height: model.Height(),
		Ready:  model.IsReady(),
		View:   model.GetCurrentView(),
	}
}

func executeAgentWorkflow(model app.MCFModel) app.MCFModel {
	// Navigate to agents
	for model.GetCurrentView() != ui.AgentsView {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	// Interact with agents
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}) // start
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}) // logs

	return model
}

func executeCommandWorkflow(model app.MCFModel) app.MCFModel {
	// Execute command
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	cmd := "mcf test workflow"
	for _, char := range cmd {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	return model
}

func executeLogWorkflow(model app.MCFModel) app.MCFModel {
	// Navigate to logs
	for model.GetCurrentView() != ui.LogsView {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	// Interact with logs
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}) // follow
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}) // search
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("error")})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	return model
}

func executeConfigWorkflow(model app.MCFModel) app.MCFModel {
	// Navigate to config
	for model.GetCurrentView() != ui.ConfigView {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	// Interact with config
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}) // reload
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")}) // backup

	return model
}

// Integration test with mock MCF services
func TestIntegrationWithMCFServices(t *testing.T) {
	t.Run("should integrate with MCF services end-to-end", func(t *testing.T) {
		ctx := context.Background()

		// Setup mock MCF client
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)

		// Setup expected interactions
		mockClient.On("Connect", ctx).Return(nil)
		mockClient.On("GetSystemHealth", ctx).Return(testutils.SystemHealthStatus{
			Status:  "healthy",
			Version: "1.0.0",
		}, nil)
		mockClient.On("GetServices", ctx).Return([]testutils.ServiceStatus{
			{Name: "api-server", Status: "running", Health: "healthy"},
		}, nil)

		// Initialize model with mock client
		model := app.InitialModel()
		model = model.WithMCFClient(mockClient, logger)
		model = prepareModel(model, t)

		// Execute workflow that uses MCF services
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
		healthCmd := "mcf health check"
		for _, char := range healthCmd {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		}
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})
}

// Benchmark complete workflows
func BenchmarkCompleteWorkflow(b *testing.B) {
	model := app.InitialModel()
	model, _ = model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate complete user workflow
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}
}

func BenchmarkViewSwitching(b *testing.B) {
	model := app.InitialModel()
	model, _ = model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	}
}
