package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ApplicationMode int

const (
	ModeMainMenu ApplicationMode = iota
	ModeInstaller
	ModeConfigurator
	ModeTemplateBrowser
	ModeRunner
)

type MainModel struct {
	mode        ApplicationMode
	cursor      int
	choice      string
	subModel    tea.Model
	configMgr   *ConfigManager
	projectPath string
	initialized bool
}

func NewMainModel() MainModel {
	wd, _ := os.Getwd()
	return MainModel{
		mode:        ModeMainMenu,
		cursor:      0,
		projectPath: wd,
		configMgr:   NewConfigManager(wd),
	}
}

func (m MainModel) Init() tea.Cmd {
	// Check if MCF is already initialized
	if m.isInitialized() {
		m.initialized = true
	}
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global messages first
	if globalMsg, ok := msg.(MCFMessage); ok {
		cmd := m.handleGlobalMessage(globalMsg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Handle window size changes
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		GlobalState.SetWindowSize(sizeMsg.Width, sizeMsg.Height)
	}

	// Handle mode-specific logic
	var modeCmd tea.Cmd
	var updatedModel tea.Model
	switch m.mode {
	case ModeMainMenu:
		updatedModel, modeCmd = m.handleMainMenu(msg)
		if mainModel, ok := updatedModel.(MainModel); ok {
			m = mainModel
		}
	case ModeInstaller:
		updatedModel, modeCmd = m.handleInstaller(msg)
		if mainModel, ok := updatedModel.(MainModel); ok {
			m = mainModel
		}
	case ModeConfigurator:
		updatedModel, modeCmd = m.handleConfigurator(msg)
		if mainModel, ok := updatedModel.(MainModel); ok {
			m = mainModel
		}
	case ModeTemplateBrowser:
		updatedModel, modeCmd = m.handleTemplateBrowser(msg)
		if mainModel, ok := updatedModel.(MainModel); ok {
			m = mainModel
		}
	case ModeRunner:
		updatedModel, modeCmd = m.handleRunner(msg)
		if mainModel, ok := updatedModel.(MainModel); ok {
			m = mainModel
		}
	default:
		return m, tea.Batch(cmds...)
	}

	if modeCmd != nil {
		cmds = append(cmds, modeCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var baseView string

	switch m.mode {
	case ModeMainMenu:
		baseView = m.mainMenuView()
	case ModeInstaller:
		if m.subModel != nil {
			baseView = m.subModel.View()
		} else {
			baseView = "Starting installer..."
		}
	case ModeConfigurator:
		if m.subModel != nil {
			baseView = m.subModel.View()
		} else {
			baseView = "Starting configurator..."
		}
	case ModeTemplateBrowser:
		if m.subModel != nil {
			baseView = m.subModel.View()
		} else {
			baseView = "Starting template browser..."
		}
	case ModeRunner:
		if m.subModel != nil {
			baseView = m.subModel.View()
		} else {
			baseView = "Starting MCF runner..."
		}
	default:
		baseView = "Unknown mode"
	}

	// Overlay notifications on top of base view
	return m.renderWithNotifications(baseView)
}

func (m MainModel) handleMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			maxOptions := 4
			if !m.initialized {
				maxOptions = 3 // Hide "Run" option if not initialized
			}
			if m.cursor < maxOptions-1 {
				m.cursor++
			}
		case "enter", " ":
			choice := m.getAvailableChoices()[m.cursor]
			switch choice {
			case "🚀 Run Claude MCF":
				m.mode = ModeRunner
				return m.startRunner()
			case "🧩 Template Browser":
				m.mode = ModeTemplateBrowser
				return m.startTemplateBrowser()
			case "⚙️  Configure":
				m.mode = ModeConfigurator
				return m.startConfigurator()
			case "📦 Install/Setup":
				m.mode = ModeInstaller
				return m.startInstaller()
			}
		}
	}
	return m, nil
}

func (m MainModel) handleInstaller(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.subModel == nil {
		return m, nil
	}

	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Check for installer completion
	if installMsg, ok := msg.(InstallCompleteMsg); ok {
		if installMsg.Success {
			m.initialized = true
			m.mode = ModeMainMenu
			m.subModel = nil
		}
	}

	return m, cmd
}

func (m MainModel) handleConfigurator(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Initialize configurator sub-model if not present
	if m.subModel == nil {
		configurator := NewConfiguratorModel(m.configMgr)
		m.subModel = configurator
		return m, configurator.Init()
	}

	// Handle escape key to return to main menu
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		// Check if configurator is in its base state, otherwise let it handle esc first
		if configurator, ok := m.subModel.(ConfiguratorModel); ok {
			if configurator.state == ConfigStateSectionList {
				m.mode = ModeMainMenu
				m.subModel = nil
				return m, nil
			}
		}
	}

	// Delegate to configurator model
	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Check for configuration-specific completion messages
	if configuratorModel, ok := m.subModel.(ConfiguratorModel); ok {
		// Handle configurator state changes
		switch configuratorModel.state {
		case ConfigStateError:
			// Add error notification
			if len(configuratorModel.errors) > 0 {
				for field, errMsg := range configuratorModel.errors {
					notification := Notification{
						Type:      "error",
						Title:     "Configuration Error",
						Message:   fmt.Sprintf("%s: %s", field, errMsg),
						Timestamp: time.Now(),
					}
					GlobalState.AddNotification(notification)
				}
			}
		}

		// Update global state with configuration changes
		if configuratorModel.unsavedChanges {
			GlobalState.SetConfigurationDirty(true)
		}
	}

	return m, cmd
}

func (m MainModel) handleTemplateBrowser(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Initialize template browser sub-model if not present
	if m.subModel == nil {
		templateBrowser := NewTemplateBrowserModel()
		m.subModel = templateBrowser
		return m, templateBrowser.Init()
	}

	// Handle escape key to return to main menu from template list state
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		if templateBrowser, ok := m.subModel.(TemplateBrowserModel); ok {
			if templateBrowser.state == StateTemplateList {
				m.mode = ModeMainMenu
				m.subModel = nil
				return m, nil
			}
		}
	}

	// Delegate to template browser model
	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Handle template browser state transitions and notifications
	if templateBrowser, ok := m.subModel.(TemplateBrowserModel); ok {
		switch templateBrowser.state {
		case StateInstallComplete:
			// Show success notification
			if templateBrowser.selectedTemplate != nil {
				notification := Notification{
					Type:      "success",
					Title:     "Template Installed",
					Message:   fmt.Sprintf("Successfully installed template: %s", templateBrowser.selectedTemplate.Name),
					Timestamp: time.Now(),
				}
				GlobalState.AddNotification(notification)

				// Publish template installation success to message bus
				successMsg := NewUISuccessMessage(
					"Template Installation Complete",
					fmt.Sprintf("Template '%s' has been installed successfully", templateBrowser.selectedTemplate.Name),
				)
				GlobalMessageBus.Publish(successMsg)
			}

		case StateTemplateError:
			// Show error notification
			if templateBrowser.errorMsg != "" {
				notification := Notification{
					Type:      "error",
					Title:     "Template Installation Failed",
					Message:   templateBrowser.errorMsg,
					Timestamp: time.Now(),
				}
				GlobalState.AddNotification(notification)

				// Publish template installation error to message bus
				errorMsg := NewUIErrorMessage(
					"Template Installation Error",
					templateBrowser.errorMsg,
					map[string]interface{}{
						"template": templateBrowser.selectedTemplate,
						"context":  "template_installation",
					},
				)
				GlobalMessageBus.Publish(errorMsg)
			}
		}
	}

	return m, cmd
}

func (m MainModel) handleRunner(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Initialize MCF runner sub-model if not present
	if m.subModel == nil {
		runner := NewMCFRunnerModel(m.projectPath)
		m.subModel = runner
		return m, runner.Init()
	}

	// Handle escape key to return to main menu from operation select state
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		if runner, ok := m.subModel.(MCFRunnerModel); ok {
			if runner.state == MCFStateOperationSelect {
				m.mode = ModeMainMenu
				m.subModel = nil
				return m, nil
			}
		}
	}

	// Handle global quit key combination
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
		// Allow graceful shutdown from any runner state
		m.mode = ModeMainMenu
		m.subModel = nil
		return m, nil
	}

	// Delegate to MCF runner model
	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Handle MCF runner state transitions and notifications
	if runner, ok := m.subModel.(MCFRunnerModel); ok {
		switch runner.state {
		case MCFStateExecuting:
			// Operation is executing - could add a notification here if needed
			if runner.selected != nil {
				// Use notification system instead of non-existent method
				notification := Notification{
					Type:      "info",
					Title:     "Operation Running",
					Message:   fmt.Sprintf("Executing: %s", runner.selected.name),
					Timestamp: time.Now(),
				}
				GlobalState.AddNotification(notification)
			}

		case MCFStateResults:
			if runner.result != nil {
				if runner.result.Success {
					// Show success notification
					notification := Notification{
						Type:      "success",
						Title:     "Operation Completed",
						Message:   fmt.Sprintf("Successfully executed: %s", runner.selected.name),
						Timestamp: time.Now(),
					}
					GlobalState.AddNotification(notification)

					// Publish operation success to message bus
					successMsg := NewUISuccessMessage(
						"MCF Operation Complete",
						fmt.Sprintf("Operation '%s' completed successfully in %v",
							runner.selected.name,
							runner.result.Duration.Round(time.Millisecond)),
					)
					GlobalMessageBus.Publish(successMsg)
				} else {
					// Show error notification
					notification := Notification{
						Type:      "error",
						Title:     "Operation Failed",
						Message:   fmt.Sprintf("Failed to execute: %s", runner.selected.name),
						Timestamp: time.Now(),
					}
					GlobalState.AddNotification(notification)

					// Publish operation error to message bus
					errorMsg := NewUIErrorMessage(
						"MCF Operation Error",
						runner.result.Error,
						map[string]interface{}{
							"operation": runner.selected.name,
							"duration":  runner.result.Duration,
							"context":   "mcf_execution",
						},
					)
					GlobalMessageBus.Publish(errorMsg)
				}
			}

		case MCFStateError:
			// Handle system errors
			if runner.error != nil {
				notification := Notification{
					Type:      "error",
					Title:     "System Error",
					Message:   runner.error.Error(),
					Timestamp: time.Now(),
				}
				GlobalState.AddNotification(notification)

				// Publish system error to message bus
				errorMsg := NewUIErrorMessage(
					"MCF System Error",
					runner.error.Error(),
					map[string]interface{}{
						"context": "mcf_runner_system",
					},
				)
				GlobalMessageBus.Publish(errorMsg)
			}

		case MCFStateOperationSelect:
			// Operation selection state - ready for user interaction
		}

		// Handle MCF-specific messages from the runner
		if mcfMsg, ok := msg.(MCFExecutionCompleteMsg); ok {
			// Process execution completion
			result := mcfMsg.Result

			// Update project state if operation was successful and made changes
			if result.Success && runner.selected != nil {
				// Add success notification for different operation types
				var successMsg string
				switch runner.selected.opType {
				case MCFOpTypeAgent:
					// Agent interactions might update knowledge base
					successMsg = "Agent interaction completed successfully"
				case MCFOpTypeTemplate:
					// Template operations might modify project structure
					successMsg = "Template operation completed successfully"
				case MCFOpTypeCommand:
					// Command operations might modify configuration or files
					successMsg = "Command executed successfully"
				}

				if successMsg != "" {
					notification := Notification{
						Type:      "success",
						Title:     "Operation Completed",
						Message:   successMsg,
						Timestamp: time.Now(),
					}
					GlobalState.AddNotification(notification)
				}
			}
		}
	}

	return m, cmd
}

func (m MainModel) mainMenuView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Align(lipgloss.Center)

	title := titleStyle.Render("🚀 MCF - Multi Component Framework")
	subtitle := subtitleStyle.Render("Development Automation Platform with Claude Code Integration")

	var status string
	if m.initialized {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Render("✅ MCF is initialized and ready")
	} else {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Render("⚠️  MCF needs to be set up first")
	}

	menu := m.renderMenu()

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Use arrow keys to navigate, Enter to select, q to quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		subtitle,
		"",
		status,
		"",
		menu,
		"",
		instructions,
		"",
	)
}

func (m MainModel) renderMenu() string {
	choices := m.getAvailableChoices()
	var items []string

	for i, choice := range choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "▶ "
		}

		style := lipgloss.NewStyle()
		if m.cursor == i {
			style = style.Foreground(lipgloss.Color("99")).Bold(true)
		} else {
			style = style.Foreground(lipgloss.Color("243"))
		}

		items = append(items, cursor+style.Render(choice))
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (m MainModel) getAvailableChoices() []string {
	if m.initialized {
		return []string{
			"🚀 Run Claude MCF",
			"🧩 Template Browser",
			"⚙️  Configure",
			"📦 Install/Setup",
		}
	}
	return []string{
		"🧩 Template Browser",
		"⚙️  Configure",
		"📦 Install/Setup",
	}
}

func (m MainModel) startInstaller() (MainModel, tea.Cmd) {
	installer := NewInstallerModel()
	m.subModel = installer
	return m, installer.Init()
}

func (m MainModel) startTemplateBrowser() (MainModel, tea.Cmd) {
	templateBrowser := NewTemplateBrowserModel()
	m.subModel = templateBrowser
	return m, templateBrowser.Init()
}

func (m MainModel) startConfigurator() (MainModel, tea.Cmd) {
	configurator := NewConfiguratorModel(m.configMgr)
	m.subModel = configurator
	return m, configurator.Init()
}

func (m MainModel) startRunner() (MainModel, tea.Cmd) {
	runner := NewMCFRunnerModel(m.projectPath)
	m.subModel = runner
	return m, runner.Init()
}

func (m MainModel) isInitialized() bool {
	// Check if .claude directory and essential files exist
	claudeDir := filepath.Join(m.projectPath, ".claude")
	settingsFile := filepath.Join(claudeDir, "settings.json")

	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		return false
	}

	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func main() {
	// Check command line arguments for direct modes
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install", "setup":
			// Direct to installer
			installer := NewInstallerModel()
			p := tea.NewProgram(installer)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Installer error: %v\n", err)
				os.Exit(1)
			}
			return
		case "config", "configure":
			// Direct to configuration
			schemaPath := "config-schema.yaml"
			if err := RunConfigurationEditor(schemaPath); err != nil {
				fmt.Printf("Configuration error: %v\n", err)
				os.Exit(1)
			}
			return
		case "run":
			// Direct to runner
			fmt.Println("Runner mode not implemented yet")
			return
		case "test", "--test":
			// Run component tests
			RunComponentTests()
			return
		case "test-comprehensive", "--test-comprehensive":
			// Run comprehensive tests
			RunComprehensiveTests()
			return
		case "help", "--help", "-h":
			printHelp()
			return
		}
	}

	// Start main menu
	m := NewMainModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	help := `MCF - Multi Component Framework

A sophisticated development automation platform built around Claude Code integration.

USAGE:
    mcf [COMMAND]

COMMANDS:
    install, setup     Run the interactive installer
    config, configure  Configure MCF settings
    run               Run MCF (requires setup)
    help              Show this help message

INTERACTIVE MODE:
    Run 'mcf' without arguments to start the interactive menu.

FEATURES:
    🧠 AI Agent System        - 9 specialized agents for development tasks
    ⚡ Custom Command System  - 50+ workflow automation commands  
    🔧 Intelligent Hooks     - Event-driven automation system
    📚 Semantic Analysis     - Serena integration for code understanding
    🛡️ Security Features     - Input validation and audit logging

For more information, visit: https://github.com/your-repo/mcf
`
	fmt.Print(help)
}

// handleGlobalMessage processes MCF messages for cross-component communication
func (m MainModel) handleGlobalMessage(msg MCFMessage) tea.Cmd {
	switch msg.Type() {
	case MsgModeTransition:
		if transMsg, ok := msg.(ModeTransitionMessage); ok {
			m.mode = transMsg.ToMode
			GlobalState.SetMode(transMsg.ToMode)
		}
	case MsgInstallComplete:
		if installMsg, ok := msg.(InstallationMessage); ok {
			if installMsg.Success {
				m.initialized = true
				// Transition back to main menu
				m.mode = ModeMainMenu
				m.subModel = nil
			}
		}
	case MsgUIError:
		if uiMsg, ok := msg.(UIMessage); ok {
			// Add error notification
			notification := Notification{
				Type:      "error",
				Title:     uiMsg.Title,
				Message:   uiMsg.Message,
				Timestamp: uiMsg.Timestamp(),
			}
			GlobalState.AddNotification(notification)
		}
	case MsgUISuccess:
		if uiMsg, ok := msg.(UIMessage); ok {
			// Add success notification
			notification := Notification{
				Type:      "success",
				Title:     uiMsg.Title,
				Message:   uiMsg.Message,
				Timestamp: uiMsg.Timestamp(),
			}
			GlobalState.AddNotification(notification)
		}
	case MsgUIWarning:
		if uiMsg, ok := msg.(UIMessage); ok {
			// Add warning notification
			notification := Notification{
				Type:      "warning",
				Title:     uiMsg.Title,
				Message:   uiMsg.Message,
				Timestamp: uiMsg.Timestamp(),
			}
			GlobalState.AddNotification(notification)
		}
	}
	return nil
}

// renderWithNotifications overlays notifications on the base view
func (m MainModel) renderWithNotifications(baseView string) string {
	notifications := GlobalState.GetNotifications()
	if len(notifications) == 0 {
		return baseView
	}

	// Get window size for positioning
	width, height := GlobalState.GetWindowSize()
	if width == 0 || height == 0 {
		// Fallback dimensions
		width = 80
		height = 24
	}

	// Render notifications as overlay
	notificationView := m.renderNotifications(notifications, width)

	// Position notifications at the top-right of the screen
	lines := strings.Split(baseView, "\n")
	notificationLines := strings.Split(notificationView, "\n")

	// Ensure we have enough lines in the base view
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Overlay notifications on the first few lines
	for i, notifLine := range notificationLines {
		if i < len(lines) && i < 5 { // Show max 5 notification lines
			// Pad the base line and append notification
			baseLine := lines[i]
			if len(baseLine) < width-40 {
				baseLine += strings.Repeat(" ", width-40-len(baseLine))
			}
			lines[i] = baseLine + notifLine
		}
	}

	return strings.Join(lines, "\n")
}

// renderNotifications creates a notification overlay
func (m MainModel) renderNotifications(notifications []Notification, width int) string {
	if len(notifications) == 0 {
		return ""
	}

	var notifLines []string

	for i, notif := range notifications {
		if i >= 3 { // Show max 3 notifications
			break
		}

		var style lipgloss.Style
		var icon string

		switch notif.Type {
		case "error":
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Background(lipgloss.Color("52")).
				Padding(0, 1)
			icon = "❌"
		case "success":
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("34")).
				Background(lipgloss.Color("22")).
				Padding(0, 1)
			icon = "✅"
		case "warning":
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("220")).
				Background(lipgloss.Color("58")).
				Padding(0, 1)
			icon = "⚠️"
		default:
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Background(lipgloss.Color("235")).
				Padding(0, 1)
			icon = "ℹ️"
		}

		// Format notification text
		text := fmt.Sprintf("%s %s", icon, notif.Title)
		if len(text) > 35 {
			text = text[:32] + "..."
		}

		notifLines = append(notifLines, style.Render(text))
	}

	return strings.Join(notifLines, "\n")
}
