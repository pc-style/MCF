package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mcf-dev/tui/internal/mcf"
	"mcf-dev/tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application model
type MCFModel struct {
	// Core application state
	ready  bool
	width  int
	height int

	// MCF integration
	mcfAdapter *mcf.MCFAdapter

	// UI components
	theme        *ui.Theme
	navigation   *ui.Navigation
	dashboard    *ui.Dashboard
	agentsList   *ui.InteractiveList
	commandsList *ui.InteractiveList
	logViewer    *ui.LogViewer
	commandInput *ui.CommandInput

	// View state
	showHelp bool

	// Performance tracking
	lastInteractionTime int64
}

func InitialModel() MCFModel {
	theme := ui.NewTheme()
	navigation := ui.NewNavigation(theme)
	dashboard := ui.NewDashboard(theme)

	// Initialize MCF adapter
	mcfRoot := findMCFRoot()
	mcfAdapter, err := mcf.NewMCFAdapter(mcfRoot)
	if err != nil {
		// Fallback to mock data if MCF adapter fails
		mcfAdapter = nil
	}

	// Initialize components
	agentsList := ui.NewInteractiveList(theme, "Agents", 20)
	commandsList := ui.NewInteractiveList(theme, "Command History", 20)
	logViewer := ui.NewLogViewer(theme, 20)
	commandInput := ui.NewCommandInput(theme)

	// Setup initial data (will use real MCF data if adapter is available)
	setupInitialData(agentsList, commandsList, logViewer, mcfAdapter)

	model := MCFModel{
		mcfAdapter:   mcfAdapter,
		theme:        theme,
		navigation:   navigation,
		dashboard:    dashboard,
		agentsList:   agentsList,
		commandsList: commandsList,
		logViewer:    logViewer,
		commandInput: commandInput,
		showHelp:     false,
	}

	// Initialize dashboard with real MCF data
	if mcfAdapter != nil {
		updateDashboardWithRealData(&model)
	}

	return model
}

func setupInitialData(agentsList *ui.InteractiveList, commandsList *ui.InteractiveList, logViewer *ui.LogViewer, mcfAdapter *mcf.MCFAdapter) {
	if mcfAdapter != nil {
		// Use real MCF data
		setupRealData(agentsList, commandsList, logViewer, mcfAdapter)
	} else {
		// Fallback to mock data
		setupMockData(agentsList, commandsList, logViewer)
	}
}

func setupRealData(agentsList *ui.InteractiveList, commandsList *ui.InteractiveList, logViewer *ui.LogViewer, mcfAdapter *mcf.MCFAdapter) {
	// Setup real agents data from MCF
	realAgents := mcfAdapter.GetAgents()
	agentItems := make([]ui.ListItem, len(realAgents))
	for i, agent := range realAgents {
		agentItems[i] = ui.ListItem{
			Title:       agent.Name,
			Status:      agent.Status,
			Description: agent.Description,
		}
	}
	agentsList.SetItems(agentItems)

	// Setup real commands data from MCF
	commandsByCategory := mcfAdapter.GetCommandsByCategory()
	commandItems := []ui.ListItem{}

	// Add commands from each category
	for _, cmds := range commandsByCategory {
		for _, cmd := range cmds {
			commandItems = append(commandItems, ui.ListItem{
				Title:       cmd.Name,
				Status:      "available",
				Description: cmd.Description,
			})
		}
	}

	// If no commands found, add some defaults
	if len(commandItems) == 0 {
		commandItems = []ui.ListItem{
			{Title: "serena:status", Status: "available", Description: "Check Serena integration status"},
			{Title: "orchestration:status", Status: "available", Description: "Check orchestration system status"},
			{Title: "agent:status", Status: "available", Description: "Check agent statuses"},
		}
	}
	commandsList.SetItems(commandItems)

	// Add real system logs from MCF
	logs := mcfAdapter.GetSystemLogs(10)
	for _, entry := range logs {
		logViewer.AddLog(entry)
	}
}

func setupMockData(agentsList *ui.InteractiveList, commandsList *ui.InteractiveList, logViewer *ui.LogViewer) {
	// Setup agents data (fallback mock data)
	agentItems := []ui.ListItem{
		{Title: "orchestrator", Status: "active", Description: "Main coordination agent - managing task distribution"},
		{Title: "frontend-developer", Status: "active", Description: "Building React components and UI interfaces"},
		{Title: "backend-developer", Status: "idle", Description: "API development and database management"},
		{Title: "test-engineer", Status: "active", Description: "Running automated tests and quality checks"},
		{Title: "system-architect", Status: "idle", Description: "System design and architecture planning"},
		{Title: "go-tui-expert", Status: "active", Description: "Terminal UI development and optimization"},
	}
	agentsList.SetItems(agentItems)

	// Setup command history (fallback mock data)
	commandItems := []ui.ListItem{
		{Title: "mcf agents status", Description: "Check status of all agents"},
		{Title: "mcf serena start", Description: "Start Serena integration service"},
		{Title: "mcf deploy --stage dev", Description: "Deploy to development environment"},
		{Title: "mcf test --coverage", Description: "Run tests with coverage report"},
		{Title: "mcf logs tail -f", Description: "Follow application logs in real-time"},
	}
	commandsList.SetItems(commandItems)

	// Add sample log entries (fallback mock data)
	now := time.Now()
	sampleLogs := []ui.LogEntry{
		{Timestamp: now.Add(-10 * time.Minute), Level: "INFO", Component: "orchestrator", Message: "MCF system initialized successfully"},
		{Timestamp: now.Add(-8 * time.Minute), Level: "INFO", Component: "frontend-dev", Message: "Started building React dashboard components"},
		{Timestamp: now.Add(-5 * time.Minute), Level: "WARN", Component: "test-engineer", Message: "Test coverage below threshold (85% < 90%)"},
		{Timestamp: now.Add(-3 * time.Minute), Level: "INFO", Component: "go-tui-expert", Message: "TUI dashboard layout optimized for responsive design"},
		{Timestamp: now.Add(-1 * time.Minute), Level: "INFO", Component: "system", Message: "Health check passed - all systems operational"},
		{Timestamp: now, Level: "INFO", Component: "orchestrator", Message: "Ready to accept new tasks"},
	}

	for _, entry := range sampleLogs {
		logViewer.AddLog(entry)
	}
}

// findMCFRoot attempts to find the MCF project root directory
func findMCFRoot() string {
	// Start from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "/Users/pcstyle/mcf-dev" // fallback
	}

	// Walk up the directory tree looking for .claude directory
	current := cwd
	for {
		claudeDir := filepath.Join(current, ".claude")
		if _, err := os.Stat(claudeDir); err == nil {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root directory
			break
		}
		current = parent
	}

	// Fallback to known path
	return "/Users/pcstyle/mcf-dev"
}

// updateDashboardWithRealData populates dashboard with real MCF data
func updateDashboardWithRealData(model *MCFModel) {
	if model.mcfAdapter == nil {
		return
	}

	// Get real system info
	version := model.mcfAdapter.GetVersion()
	if version == "" || version == "unknown" {
		version = "MCF Development"
	}

	serenaStatus := model.mcfAdapter.GetSerenaStatus()

	// Get real agent data
	agents := model.mcfAdapter.GetAgents()
	activeCount := 0
	agentStatuses := []ui.AgentStatus{}

	for _, agent := range agents {
		if agent.Status == "active" {
			activeCount++
		}

		// Convert MCF agent to dashboard agent status
		agentStatus := ui.AgentStatus{
			Name:        agent.Name,
			Status:      agent.Status,
			LastSeen:    agent.LastActive,
			TasksActive: 0, // TODO: Get real task data
			TasksTotal:  0,
		}
		agentStatuses = append(agentStatuses, agentStatus)
	}

	// Update dashboard with real data
	model.dashboard.SetSystemHealth(version, serenaStatus, activeCount, len(agents))
	model.dashboard.SetAgentStatuses(agentStatuses)

	// Get real command history from discovered Claude commands
	commands := model.mcfAdapter.GetCommands()
	commandNames := []string{}

	// Prioritize common Claude commands
	priorityCommands := []string{"agent:auto", "serena:init", "project:analyze", "orchestration:status", "project:deploy"}

	for _, priority := range priorityCommands {
		if _, exists := commands[priority]; exists {
			commandNames = append(commandNames, priority)
		}
		if len(commandNames) >= 5 {
			break
		}
	}

	// Fill remaining slots with other discovered commands
	if len(commandNames) < 5 {
		for name := range commands {
			// Skip if already added
			found := false
			for _, existing := range commandNames {
				if existing == name {
					found = true
					break
				}
			}
			if !found {
				commandNames = append(commandNames, name)
				if len(commandNames) >= 5 {
					break
				}
			}
		}
	}

	model.dashboard.SetCommandHistory(commandNames)

	// Add initial activity (only once during initialization)
	model.dashboard.AddRecentActivity("info", "MCF TUI started", fmt.Sprintf("Loaded %d agents, %d commands", len(agents), len(commands)))
}

func (m MCFModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// Periodic update command
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func (m *MCFModel) SetView(view ui.View) {
	m.navigation.SetView(view)

	// Update component focus based on current view
	currentView := m.navigation.GetCurrentView()
	m.agentsList.SetFocus(currentView == ui.AgentsView)
	m.commandsList.SetFocus(currentView == ui.CommandsView)
	m.commandInput.SetFocus(currentView == ui.CommandBarView)

	if currentView == ui.LogsView {
		// Logs view doesn't need explicit focus
	}
}

func (m *MCFModel) Width() int  { return m.width }
func (m *MCFModel) Height() int { return m.height }

func (m *MCFModel) ToggleHelp() {
	m.showHelp = !m.showHelp
	if m.navigation.GetCurrentView() == ui.DashboardView {
		m.dashboard.ToggleHelp()
	}
}

// View rendering methods
func (m MCFModel) View() string {
	if !m.ready {
		return "Initializing MCF TUI..."
	}

	// Global help overlay
	if m.showHelp && m.navigation.GetCurrentView() != ui.DashboardView {
		return m.navigation.RenderHelp(m.width)
	}

	// Main layout
	header := m.renderHeader()
	content := m.renderContent()
	footer := m.renderFooter()

	// Combine layout
	mainHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer) - 2
	content = lipgloss.NewStyle().Height(mainHeight).Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func (m MCFModel) renderHeader() string {
	// Title and navigation
	title := m.theme.Title.Render("MCF Development Framework")

	// Tab bar
	tabBar := m.navigation.RenderTabBar(m.width)

	// Breadcrumb (if applicable)
	breadcrumb := m.navigation.RenderBreadcrumb()

	header := title + "\n" + tabBar
	if breadcrumb != "" {
		header += "\n" + breadcrumb
	}

	return header
}

func (m MCFModel) renderContent() string {
	currentView := m.navigation.GetCurrentView()
	contentWidth := m.width
	contentHeight := m.height - 8 // Account for header and footer

	// Ensure minimum height to prevent crashes
	if contentHeight < 10 {
		contentHeight = 10
	}

	switch currentView {
	case ui.DashboardView:
		return m.dashboard.Render(contentWidth, contentHeight)

	case ui.AgentsView:
		return m.renderAgentsView(contentWidth, contentHeight)

	case ui.CommandsView:
		return m.renderCommandsView(contentWidth, contentHeight)

	case ui.LogsView:
		return m.logViewer.Render(contentWidth)

	case ui.ConfigView:
		return m.renderConfigView(contentWidth, contentHeight)

	case ui.CommandBarView:
		return m.renderCommandBar(contentWidth, contentHeight)

	default:
		return m.theme.Error.Render("Unknown view")
	}
}

func (m MCFModel) renderAgentsView(width, height int) string {
	// Main agents list
	agentsList := m.agentsList.Render(width * 2 / 3)

	// Agent details panel
	selectedAgent := m.agentsList.GetSelectedItem()
	detailsContent := ""

	if selectedAgent != nil {
		detailsContent += m.theme.Subtitle.Render("Agent Details") + "\n\n"
		detailsContent += m.theme.Body.Render("Name: "+selectedAgent.Title) + "\n"
		detailsContent += "Status: " + ui.RenderStatusIndicator(selectedAgent.Status, m.theme) + "\n\n"
		detailsContent += m.theme.Muted.Render(selectedAgent.Description) + "\n\n"

		// Agent actions
		detailsContent += m.theme.Subtitle.Render("Actions") + "\n"
		detailsContent += m.theme.ListItem.Render("s - Start/Stop Agent") + "\n"
		detailsContent += m.theme.ListItem.Render("r - Restart Agent") + "\n"
		detailsContent += m.theme.ListItem.Render("l - View Logs") + "\n"
		detailsContent += m.theme.ListItem.Render("c - Configure Agent") + "\n"
	}

	detailsPanel := ui.RenderBox(detailsContent, "Agent Details", width/3, height-2, m.theme)

	return lipgloss.JoinHorizontal(lipgloss.Top, agentsList, detailsPanel)
}

func (m MCFModel) renderCommandsView(width, height int) string {
	// Commands history list
	commandsList := m.commandsList.Render(width * 2 / 3)

	// Command details and actions
	selectedCommand := m.commandsList.GetSelectedItem()
	detailsContent := ""

	if selectedCommand != nil {
		detailsContent += m.theme.Subtitle.Render("Command Details") + "\n\n"
		detailsContent += m.theme.Info.Render(selectedCommand.Title) + "\n\n"
		detailsContent += m.theme.Muted.Render(selectedCommand.Description) + "\n\n"

		// Command actions
		detailsContent += m.theme.Subtitle.Render("Actions") + "\n"
		detailsContent += m.theme.ListItem.Render("Enter - Re-execute Command") + "\n"
		detailsContent += m.theme.ListItem.Render("e - Edit Command") + "\n"
		detailsContent += m.theme.ListItem.Render("d - Delete from History") + "\n"
		detailsContent += m.theme.ListItem.Render("c - Copy to Clipboard") + "\n"
	}

	detailsPanel := ui.RenderBox(detailsContent, "Command Actions", width/3, height-2, m.theme)

	return lipgloss.JoinHorizontal(lipgloss.Top, commandsList, detailsPanel)
}

func (m MCFModel) renderConfigView(width, height int) string {
	content := m.theme.Title.Render("MCF Configuration") + "\n\n"

	if m.mcfAdapter != nil {
		// Real configuration from MCF
		settings := m.mcfAdapter.GetSettings()

		// System Settings
		content += m.theme.Subtitle.Render("System Settings") + "\n"
		version := "unknown"
		if settings != nil {
			version = settings.Version
		}
		if version == "" {
			version = m.mcfAdapter.GetVersion()
		}
		content += m.theme.Body.Render("• MCF Version: "+version) + "\n"
		content += m.theme.Body.Render("• Working Directory: "+findMCFRoot()) + "\n"

		if settings != nil {
			content += m.theme.Body.Render("• Output Style: "+settings.OutputStyle) + "\n"
		}
		content += m.theme.Body.Render("• Auto-refresh: Enabled (5s)") + "\n\n"

		// Agent Configuration
		content += m.theme.Subtitle.Render("Agent Configuration") + "\n"
		agents := m.mcfAdapter.GetAgents()
		content += m.theme.Body.Render(fmt.Sprintf("• Total Agents: %d", len(agents))) + "\n"

		activeCount := 0
		for _, agent := range agents {
			if agent.Status == "active" {
				activeCount++
			}
		}
		content += m.theme.Body.Render(fmt.Sprintf("• Active Agents: %d", activeCount)) + "\n"
		content += m.theme.Body.Render(fmt.Sprintf("• Available Agents: %d", len(agents)-activeCount)) + "\n\n"

		// Serena Integration
		content += m.theme.Subtitle.Render("Serena Integration") + "\n"
		serenaStatus := m.mcfAdapter.GetSerenaStatus()
		content += m.theme.Body.Render("• Status: ") + ui.RenderStatusIndicator(serenaStatus, m.theme) + "\n"

		if serenaAdapter := m.mcfAdapter.GetSerenaAdapter(); serenaAdapter != nil {
			if serenaAdapter.IsEnabled() {
				content += m.theme.Body.Render("• Host: localhost:8080") + "\n"
				content += m.theme.Body.Render("• Service: Available") + "\n"
			} else {
				content += m.theme.Body.Render("• Service: Not available") + "\n"
			}
		}
		content += "\n"

		// Commands
		commands := m.mcfAdapter.GetCommands()
		commandsByCategory := m.mcfAdapter.GetCommandsByCategory()
		content += m.theme.Subtitle.Render("Commands") + "\n"
		content += m.theme.Body.Render(fmt.Sprintf("• Total Commands: %d", len(commands))) + "\n"
		content += m.theme.Body.Render(fmt.Sprintf("• Categories: %d", len(commandsByCategory))) + "\n\n"
	} else {
		// Fallback configuration display
		content += m.theme.Subtitle.Render("System Settings") + "\n"
		content += m.theme.Body.Render("• MCF Version: v1.0.0 (fallback)") + "\n"
		content += m.theme.Body.Render("• Working Directory: "+findMCFRoot()) + "\n"
		content += m.theme.Body.Render("• Log Level: INFO") + "\n"
		content += m.theme.Body.Render("• Auto-refresh: Enabled (5s)") + "\n\n"

		content += m.theme.Subtitle.Render("Status") + "\n"
		content += m.theme.Body.Render("• MCF Integration: ") + ui.RenderStatusIndicator("disconnected", m.theme) + "\n"
		content += m.theme.Body.Render("• Serena Integration: ") + ui.RenderStatusIndicator("disconnected", m.theme) + "\n\n"
	}

	content += m.theme.Subtitle.Render("Actions") + "\n"
	content += m.theme.ListItem.Render("r - Reload Configuration") + "\n"
	content += m.theme.ListItem.Render("s - Show Settings File") + "\n"
	content += m.theme.ListItem.Render("h - Show Hooks Configuration") + "\n"

	return ui.RenderBox(content, "", width, height, m.theme)
}

func (m MCFModel) renderCommandBar(width, height int) string {
	// Command input
	commandInput := m.commandInput.Render(width)

	// Instructions
	instructions := "\n" + m.theme.Muted.Render("Enter MCF command and press Enter to execute")
	instructions += "\n" + m.theme.Muted.Render("Use ↑/↓ to browse history and suggestions")
	instructions += "\n" + m.theme.Muted.Render("Press Esc to return to previous view")

	content := commandInput + instructions

	return ui.RenderBox(content, "Command Bar", width, height, m.theme)
}

func (m MCFModel) renderFooter() string {
	// Status and shortcuts
	currentView := m.navigation.GetCurrentView()

	status := ""
	switch currentView {
	case ui.DashboardView:
		status = "Dashboard - Press ? for help"
	case ui.AgentsView:
		status = "Agents - j/k to navigate, Enter for details"
	case ui.CommandsView:
		status = "Commands - j/k to navigate, Enter to execute"
	case ui.LogsView:
		status = "Logs - j/k to scroll, f to follow, / to search"
	case ui.ConfigView:
		status = "Config - System configuration and settings"
	case ui.CommandBarView:
		status = "Command Bar - Type MCF commands"
	}

	shortcuts := "Tab: Next View │ ?: Help │ q: Quit"

	footerContent := m.theme.Muted.Render(status + " │ " + shortcuts)
	return footerContent
}
