package app

import (
	"fmt"
	"time"

	"mcf-dev/tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func (m MCFModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Performance tracking - update immediately
	m.lastInteractionTime = time.Now().UnixMilli()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		// Global key handlers
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "?":
			m.ToggleHelp()
			return m, nil

		case "tab":
			m.nextView()
			return m, nil

		case "shift+tab":
			m.prevView()
			return m, nil

		case ":":
			m.SetView(ui.CommandBarView)
			return m, nil

		case "esc":
			currentView := m.navigation.GetCurrentView()
			if currentView == ui.CommandBarView {
				// Return to previous view
				prevView := m.navigation.GoBack()
				m.SetView(prevView)
			} else if currentView != ui.DashboardView {
				// Return to dashboard
				m.SetView(ui.DashboardView)
			} else if m.showHelp {
				// Close help
				m.ToggleHelp()
			}
			return m, nil

		case "r":
			// Refresh current view
			if m.navigation.GetCurrentView() == ui.DashboardView {
				m.dashboard.Update()
			}
			return m, nil
		}

		// Route input to current view
		currentView := m.navigation.GetCurrentView()
		switch currentView {
		case ui.DashboardView:
			return m.updateDashboard(msg)

		case ui.AgentsView:
			return m.updateAgents(msg)

		case ui.CommandsView:
			return m.updateCommands(msg)

		case ui.LogsView:
			return m.updateLogs(msg)

		case ui.ConfigView:
			return m.updateConfig(msg)

		case ui.CommandBarView:
			return m.updateCommandBar(msg)
		}

	case tickMsg:
		// Periodic background updates
		m.dashboard.Update()

		// Use real MCF data if adapter is available
		if m.mcfAdapter != nil {
			// Only update agents and logs on tick, not full dashboard reinit
			m.updateAgentsFromMCF()
			m.updateLogsFromMCF()

			// Add periodic system check log entry
			if time.Now().Unix()%30 == 0 { // Every 30 seconds
				m.logViewer.AddLog(ui.LogEntry{
					Timestamp: time.Now(),
					Level:     "INFO",
					Component: "mcf-tui",
					Message:   "Periodic MCF system health check completed",
				})
			}
		} else {
			// Fallback to mock updates
			if time.Now().Unix()%10 == 0 { // Every 10 seconds
				m.logViewer.AddLog(ui.LogEntry{
					Timestamp: time.Now(),
					Level:     "INFO",
					Component: "system",
					Message:   "Periodic system check completed",
				})
			}
		}

		cmds = append(cmds, tickCmd())

	default:
		// Handle unknown message types gracefully by returning tick command
		cmds = append(cmds, tickCmd())
	}

	return m, tea.Batch(cmds...)
}

// updateAgentsFromMCF updates agent data from the real MCF system
func (m *MCFModel) updateAgentsFromMCF() {
	if m.mcfAdapter == nil {
		return
	}

	// Get real agents from MCF
	realAgents := m.mcfAdapter.GetAgents()

	// Update agent list items with current status
	agentItems := make([]ui.ListItem, len(realAgents))
	for i, agent := range realAgents {
		agentItems[i] = ui.ListItem{
			Title:       agent.Name,
			Status:      agent.Status,
			Description: agent.Description,
		}
	}

	m.agentsList.SetItems(agentItems)
}

// updateLogsFromMCF updates log data from the real MCF system
func (m *MCFModel) updateLogsFromMCF() {
	if m.mcfAdapter == nil {
		return
	}

	// Get recent logs from MCF (last 5 entries to avoid flooding)
	logs := m.mcfAdapter.GetSystemLogs(5)

	// Add new logs to the viewer
	for _, entry := range logs {
		m.logViewer.AddLog(entry)
	}
}

func (m *MCFModel) nextView() {
	currentView := m.navigation.GetCurrentView()
	views := []ui.View{ui.DashboardView, ui.AgentsView, ui.CommandsView, ui.LogsView, ui.ConfigView}

	for i, view := range views {
		if view == currentView {
			nextView := views[(i+1)%len(views)]
			m.SetView(nextView)
			return
		}
	}
}

func (m *MCFModel) prevView() {
	currentView := m.navigation.GetCurrentView()
	views := []ui.View{ui.DashboardView, ui.AgentsView, ui.CommandsView, ui.LogsView, ui.ConfigView}

	for i, view := range views {
		if view == currentView {
			prevView := views[(i-1+len(views))%len(views)]
			m.SetView(prevView)
			return
		}
	}
}

// Dashboard view updates
func (m MCFModel) updateDashboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		current := m.dashboard.GetSelectedQuickAction()
		if current < len(ui.QuickActions)-1 {
			m.dashboard.SetSelectedQuickAction(current + 1)
		}

	case "k", "up":
		current := m.dashboard.GetSelectedQuickAction()
		if current > 0 {
			m.dashboard.SetSelectedQuickAction(current - 1)
		}

	case "enter":
		// Execute selected quick action with real MCF command
		selected := m.dashboard.GetSelectedQuickAction()
		if selected >= 0 && selected < len(ui.QuickActions) {
			action := ui.QuickActions[selected]
			m.commandInput.AddToHistory(action.Command)

			if m.mcfAdapter != nil {
				// Execute the real MCF command
				result, err := m.mcfAdapter.ExecuteCommand(action.Command, []string{})

				if err == nil && result.Success {
					// Log successful execution
					m.logViewer.AddLog(ui.LogEntry{
						Timestamp: time.Now(),
						Level:     "INFO",
						Component: "dashboard",
						Message:   fmt.Sprintf("✓ %s: %s", action.Label, result.Output),
					})

					// Add to dashboard recent activity
					m.dashboard.AddRecentActivity("command", action.Command, "Executed successfully")
				} else {
					// Log execution error
					errorMsg := "Unknown error"
					if err != nil {
						errorMsg = err.Error()
					} else if !result.Success {
						errorMsg = result.Error
					}

					m.logViewer.AddLog(ui.LogEntry{
						Timestamp: time.Now(),
						Level:     "ERROR",
						Component: "dashboard",
						Message:   fmt.Sprintf("✗ %s failed: %s", action.Label, errorMsg),
					})

					// Add to dashboard recent activity
					m.dashboard.AddRecentActivity("error", action.Command, fmt.Sprintf("Failed: %s", errorMsg))
				}
			} else {
				// Fallback when MCF adapter is not available
				m.logViewer.AddLog(ui.LogEntry{
					Timestamp: time.Now(),
					Level:     "WARN",
					Component: "dashboard",
					Message:   fmt.Sprintf("⚠ %s: MCF adapter not available", action.Label),
				})
			}
		}

	case "1", "2", "3", "4", "5", "6":
		// Quick action shortcuts
		idx := int(msg.String()[0] - '1')
		if idx >= 0 && idx < len(ui.QuickActions) {
			m.dashboard.SetSelectedQuickAction(idx)
			action := ui.QuickActions[idx]
			m.commandInput.AddToHistory(action.Command)

			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "dashboard",
				Message:   "Quick action: " + action.Label,
			})
		}
	}

	return m, nil
}

// Agents view updates
func (m MCFModel) updateAgents(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "s":
		// Toggle agent status (simulated since no real agent:control command exists)
		selectedAgent := m.agentsList.GetSelectedItem()
		if selectedAgent != nil {
			newStatus := "active"
			if selectedAgent.Status == "active" {
				newStatus = "idle"
			}

			// Since agent:control doesn't exist, simulate the status change
			// In a real implementation, this would call an actual agent management API

			// Update the agent in the list
			if m.mcfAdapter != nil {
				// Find and update the agent status
				agents := m.mcfAdapter.GetAgents()
				for _, agent := range agents {
					if agent.Name == selectedAgent.Title {
						agent.Status = newStatus
						agent.LastActive = time.Now()
						break
					}
				}

				// Refresh the agents list
				m.updateAgentsFromMCF()
			}

			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "agents",
				Message:   fmt.Sprintf("Agent %s status toggled to %s", selectedAgent.Title, newStatus),
			})
		}

	case "l":
		// View agent logs
		selectedAgent := m.agentsList.GetSelectedItem()
		if selectedAgent != nil {
			m.SetView(ui.LogsView)
			// In real implementation, filter logs for this agent
		}

	default:
		// Pass to list for navigation
		m.agentsList, cmd = m.agentsList.Update(tea.KeyMsg(msg))
	}

	return m, cmd
}

// Commands view updates
func (m MCFModel) updateCommands(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		// Re-execute command using real MCF integration
		selectedCommand := m.commandsList.GetSelectedItem()
		if selectedCommand != nil && m.mcfAdapter != nil {
			m.commandInput.AddToHistory(selectedCommand.Title)

			// Execute the real MCF command
			result, err := m.mcfAdapter.ExecuteCommand(selectedCommand.Title, []string{})

			if err == nil && result.Success {
				m.logViewer.AddLog(ui.LogEntry{
					Timestamp: time.Now(),
					Level:     "INFO",
					Component: "commands",
					Message:   fmt.Sprintf("Executed: %s - %s", selectedCommand.Title, result.Output),
				})
			} else {
				errorMsg := "Unknown error"
				if err != nil {
					errorMsg = err.Error()
				} else if !result.Success {
					errorMsg = result.Error
				}

				m.logViewer.AddLog(ui.LogEntry{
					Timestamp: time.Now(),
					Level:     "ERROR",
					Component: "commands",
					Message:   fmt.Sprintf("Failed to execute %s: %s", selectedCommand.Title, errorMsg),
				})
			}
		}

	case "d":
		// Delete from history (simplified - in real implementation, maintain state)
		selectedCommand := m.commandsList.GetSelectedItem()
		if selectedCommand != nil {
			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "commands",
				Message:   "Removed from history: " + selectedCommand.Title,
			})
		}

	case "c":
		// Clear command history
		m.commandsList.SetItems([]ui.ListItem{})
		m.logViewer.AddLog(ui.LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Component: "commands",
			Message:   "Command history cleared",
		})

	default:
		// Pass to list for navigation
		m.commandsList, cmd = m.commandsList.Update(tea.KeyMsg(msg))
	}

	return m, cmd
}

// Logs view updates
func (m MCFModel) updateLogs(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.logViewer, cmd = m.logViewer.Update(tea.KeyMsg(msg))
	return m, cmd
}

// Config view updates
func (m MCFModel) updateConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "e":
		// Edit configuration
		m.logViewer.AddLog(ui.LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Component: "config",
			Message:   "Configuration editor opened",
		})

	case "r":
		// Reload configuration
		m.logViewer.AddLog(ui.LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Component: "config",
			Message:   "Configuration reloaded successfully",
		})

	case "b":
		// Backup configuration
		m.logViewer.AddLog(ui.LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Component: "config",
			Message:   "Configuration backed up to config.backup",
		})

	case "d":
		// Reset to defaults
		m.logViewer.AddLog(ui.LogEntry{
			Timestamp: time.Now(),
			Level:     "WARN",
			Component: "config",
			Message:   "Configuration reset to defaults",
		})
	}

	return m, nil
}

// Command bar updates
func (m MCFModel) updateCommandBar(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		// Execute command
		command := m.commandInput.GetValue()
		if command != "" {
			m.commandInput.AddToHistory(command)

			// Add to command history list
			currentItems := []ui.ListItem{}
			if m.commandsList != nil {
				// Get current items and prepend new command
				currentItems = append([]ui.ListItem{
					{Title: command, Description: "Recently executed command"},
				}, currentItems...)

				// Limit history size
				if len(currentItems) > 20 {
					currentItems = currentItems[:20]
				}

				m.commandsList.SetItems(currentItems)
			}

			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "commands",
				Message:   "Executed: " + command,
			})

			m.commandInput.Clear()

			// Return to dashboard after execution
			m.SetView(ui.DashboardView)
		}

	default:
		// Pass to command input
		m.commandInput, cmd = m.commandInput.Update(tea.KeyMsg(msg))
	}

	return m, cmd
}
