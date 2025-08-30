package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"mcf-dev/tui/internal/ui"
)

func (m MCFModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Performance tracking
	now := time.Now().UnixMilli()
	defer func() {
		m.lastInteractionTime = now
	}()

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

		// Add simulated log entry
		if time.Now().Unix()%10 == 0 { // Every 10 seconds
			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "system",
				Message:   "Periodic system check completed",
			})
		}

		cmds = append(cmds, tickCmd())
	}

	return m, tea.Batch(cmds...)
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
		// Execute selected quick action
		selected := m.dashboard.GetSelectedQuickAction()
		if selected >= 0 && selected < len(ui.QuickActions) {
			action := ui.QuickActions[selected]
			m.commandInput.AddToHistory(action.Command)
			// In real implementation, execute the command

			// Add to log
			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "dashboard",
				Message:   "Executed quick action: " + action.Label,
			})
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
		// Start/stop agent
		selectedAgent := m.agentsList.GetSelectedItem()
		if selectedAgent != nil {
			newStatus := "active"
			if selectedAgent.Status == "active" {
				newStatus = "idle"
			}

			// Update agent status (in real implementation, call actual agent control)
			items := []ui.ListItem{}
			for _, item := range []ui.ListItem{
				{Title: "orchestrator", Status: "active", Description: "Main coordination agent - managing task distribution"},
				{Title: "frontend-developer", Status: "active", Description: "Building React components and UI interfaces"},
				{Title: "backend-developer", Status: "idle", Description: "API development and database management"},
				{Title: "test-engineer", Status: "active", Description: "Running automated tests and quality checks"},
				{Title: "system-architect", Status: "idle", Description: "System design and architecture planning"},
				{Title: "go-tui-expert", Status: "active", Description: "Terminal UI development and optimization"},
			} {
				if item.Title == selectedAgent.Title {
					item.Status = newStatus
				}
				items = append(items, item)
			}
			m.agentsList.SetItems(items)

			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "agents",
				Message:   selectedAgent.Title + " status changed to " + newStatus,
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
		// Re-execute command
		selectedCommand := m.commandsList.GetSelectedItem()
		if selectedCommand != nil {
			m.commandInput.AddToHistory(selectedCommand.Title)

			m.logViewer.AddLog(ui.LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "commands",
				Message:   "Re-executed: " + selectedCommand.Title,
			})
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
