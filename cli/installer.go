package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type InstallerModel struct {
	state          InstallerState
	currentStep    int
	totalSteps     int
	steps          []InstallStep
	progress       progress.Model
	spinner        spinner.Model
	textInput      textinput.Model
	list           list.Model
	config         MCFConfig
	errorMsg       string
	showError      bool
	installPath    string
	backupPath     string
	isUpgrade      bool
	nonInteractive bool
	currentField   int
	configFields   []ConfigField
	existingConfig *MCFConfig
}

type InstallerState int

const (
	StateWelcome InstallerState = iota
	StatePathSelection
	StateExistingDetection
	StateBackupConfirmation
	StateEnvironmentCheck
	StateConfiguration
	StateFeatureSelection
	StateInstallation
	StateShellIntegration
	StatePermissions
	StateCompletion
	StateError
)

type InstallStep struct {
	Name        string
	Description string
	Required    bool
	Status      StepStatus
	Action      func(context.Context, *InstallerModel) error
	Condition   func(*MCFConfig) bool // Optional condition for step execution
}

type StepStatus int

const (
	StatusPending StepStatus = iota
	StatusRunning
	StatusComplete
	StatusError
	StatusSkipped
)

type MCFConfig struct {
	ProjectName   string            `json:"project_name" yaml:"project_name"`
	ProjectPath   string            `json:"project_path" yaml:"project_path"`
	InstallPath   string            `json:"install_path" yaml:"install_path"`
	GitRepo       string            `json:"git_repo" yaml:"git_repo"`
	Features      []string          `json:"features" yaml:"features"`
	ClaudeAPIKey  string            `json:"claude_api_key,omitempty" yaml:"claude_api_key,omitempty"`
	SerenaEnabled bool              `json:"serena_enabled" yaml:"serena_enabled"`
	Preferences   UserPreferences   `json:"preferences" yaml:"preferences"`
	Shell         ShellConfig       `json:"shell" yaml:"shell"`
	Security      SecurityConfig    `json:"security" yaml:"security"`
	Environment   map[string]string `json:"environment" yaml:"environment"`
	Version       string            `json:"version" yaml:"version"`
	InstalledAt   time.Time         `json:"installed_at" yaml:"installed_at"`
	UpdatedAt     time.Time         `json:"updated_at" yaml:"updated_at"`
}

type UserPreferences struct {
	Theme          string `json:"theme" yaml:"theme"`
	Editor         string `json:"editor" yaml:"editor"`
	AutoUpdate     bool   `json:"auto_update" yaml:"auto_update"`
	TelemetryOpt   bool   `json:"telemetry_opt" yaml:"telemetry_opt"`
	ShowWelcome    bool   `json:"show_welcome" yaml:"show_welcome"`
	ConfirmDeletes bool   `json:"confirm_deletes" yaml:"confirm_deletes"`
	VerboseLogging bool   `json:"verbose_logging" yaml:"verbose_logging"`
}

type ShellConfig struct {
	Shell       string   `json:"shell" yaml:"shell"`
	AddToPath   bool     `json:"add_to_path" yaml:"add_to_path"`
	CreateAlias bool     `json:"create_alias" yaml:"create_alias"`
	Aliases     []string `json:"aliases" yaml:"aliases"`
	RCFile      string   `json:"rc_file" yaml:"rc_file"`
}

type SecurityConfig struct {
	ValidateInputs  bool     `json:"validate_inputs" yaml:"validate_inputs"`
	AuditLogging    bool     `json:"audit_logging" yaml:"audit_logging"`
	RequireConfirm  bool     `json:"require_confirm" yaml:"require_confirm"`
	MaxMemoryUsage  string   `json:"max_memory_usage" yaml:"max_memory_usage"`
	AllowedCommands []string `json:"allowed_commands" yaml:"allowed_commands"`
	RestrictedPaths []string `json:"restricted_paths" yaml:"restricted_paths"`
}

type ConfigField struct {
	Key         string
	Label       string
	Type        string
	Required    bool
	Default     string
	Validation  func(string) error
	Description string
	Sensitive   bool
}

type EnvironmentInfo struct {
	OS            string
	Shell         string
	HomeDir       string
	HasGit        bool
	HasCurl       bool
	HasClaude     bool
	ClaudeVersion string
	PythonVersion string
	NodeVersion   string
	GoVersion     string
}

func NewInstallerModel() InstallerModel {
	// Initialize progress bar
	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = 50

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.Focus()

	// Get user home directory for default install path
	usr, _ := user.Current()
	defaultInstallPath := filepath.Join(usr.HomeDir, "mcf")

	// Define comprehensive installation steps
	steps := []InstallStep{
		{
			Name:        "Environment Check",
			Description: "Verifying system requirements and dependencies",
			Required:    true,
			Status:      StatusPending,
			Action:      checkEnvironment,
		},
		{
			Name:        "Existing Installation Detection",
			Description: "Scanning for existing MCF installations",
			Required:    true,
			Status:      StatusPending,
			Action:      detectExistingInstallation,
		},
		{
			Name:        "Configuration Backup",
			Description: "Creating backup of existing configuration",
			Required:    false,
			Status:      StatusPending,
			Action:      backupExistingConfig,
			Condition:   func(cfg *MCFConfig) bool { return cfg.InstallPath != "" },
		},
		{
			Name:        "Directory Structure",
			Description: "Creating MCF directory structure",
			Required:    true,
			Status:      StatusPending,
			Action:      createDirectoryStructure,
		},
		{
			Name:        "Claude Code Integration",
			Description: "Setting up Claude Code CLI integration",
			Required:    true,
			Status:      StatusPending,
			Action:      setupClaudeCodeIntegration,
		},
		{
			Name:        "Serena MCP Server",
			Description: "Installing semantic analysis server",
			Required:    false,
			Status:      StatusPending,
			Action:      installSerenaMCP,
			Condition:   func(cfg *MCFConfig) bool { return cfg.SerenaEnabled },
		},
		{
			Name:        "Project Templates",
			Description: "Installing project templates and scaffolding",
			Required:    true,
			Status:      StatusPending,
			Action:      installProjectTemplates,
		},
		{
			Name:        "Command System",
			Description: "Setting up MCF command system",
			Required:    true,
			Status:      StatusPending,
			Action:      installCommandSystem,
		},
		{
			Name:        "Shell Integration",
			Description: "Configuring shell aliases and PATH",
			Required:    false,
			Status:      StatusPending,
			Action:      configureShellIntegration,
			Condition:   func(cfg *MCFConfig) bool { return cfg.Shell.AddToPath || cfg.Shell.CreateAlias },
		},
		{
			Name:        "File Permissions",
			Description: "Setting proper file permissions",
			Required:    true,
			Status:      StatusPending,
			Action:      setFilePermissions,
		},
		{
			Name:        "Configuration Validation",
			Description: "Validating installation and configuration",
			Required:    true,
			Status:      StatusPending,
			Action:      validateInstallation,
		},
	}

	// Configuration fields for the setup wizard
	configFields := []ConfigField{
		{
			Key:         "install_path",
			Label:       "Installation Directory",
			Type:        "path",
			Required:    true,
			Default:     defaultInstallPath,
			Description: "Directory where MCF will be installed",
			Validation:  validateInstallPath,
		},
		{
			Key:         "project_name",
			Label:       "Default Project Name",
			Type:        "text",
			Required:    true,
			Default:     "my-mcf-project",
			Description: "Default name for new MCF projects",
			Validation:  validateProjectName,
		},
		{
			Key:         "editor",
			Label:       "Preferred Editor",
			Type:        "select",
			Required:    false,
			Default:     "auto-detect",
			Description: "Your preferred code editor",
			Validation:  validateEditor,
		},
		{
			Key:         "shell",
			Label:       "Shell Type",
			Type:        "select",
			Required:    false,
			Default:     "auto-detect",
			Description: "Your shell for integration setup",
			Validation:  validateShell,
		},
		{
			Key:         "claude_api_key",
			Label:       "Claude API Key (Optional)",
			Type:        "password",
			Required:    false,
			Default:     "",
			Description: "Your Anthropic API key for Claude integration",
			Sensitive:   true,
			Validation:  validateClaudeAPIKey,
		},
	}

	return InstallerModel{
		state:        StateWelcome,
		currentStep:  0,
		totalSteps:   len(steps),
		steps:        steps,
		progress:     prog,
		spinner:      s,
		textInput:    ti,
		installPath:  defaultInstallPath,
		configFields: configFields,
		config: MCFConfig{
			Version:     "1.0.0",
			Environment: make(map[string]string),
			Preferences: UserPreferences{
				Theme:          "default",
				AutoUpdate:     true,
				ShowWelcome:    true,
				ConfirmDeletes: true,
				VerboseLogging: false,
			},
			Shell: ShellConfig{
				AddToPath:   true,
				CreateAlias: true,
				Aliases:     []string{"mcf", "claude-mcf"},
			},
			Security: SecurityConfig{
				ValidateInputs:  true,
				AuditLogging:    true,
				RequireConfirm:  true,
				MaxMemoryUsage:  "1GB",
				AllowedCommands: []string{"*"},
				RestrictedPaths: []string{"/etc", "/sys", "/proc"},
			},
		},
	}
}

func (m InstallerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
	)
}

func (m InstallerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateWelcome:
			return m.handleWelcomeInput(msg)
		case StatePathSelection:
			return m.handlePathSelectionInput(msg)
		case StateExistingDetection:
			return m.handleExistingDetectionInput(msg)
		case StateBackupConfirmation:
			return m.handleBackupConfirmationInput(msg)
		case StateConfiguration:
			return m.handleConfigurationInput(msg)
		case StateFeatureSelection:
			return m.handleFeatureSelectionInput(msg)
		case StateShellIntegration:
			return m.handleShellIntegrationInput(msg)
		case StateError:
			return m.handleErrorInput(msg)
		case StateCompletion:
			return m.handleCompletionInput(msg)
		default:
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case InstallCompleteMsg:
		if msg.Success {
			m.state = StateCompletion
		} else {
			m.state = StateError
			m.errorMsg = msg.Error
			m.showError = true
		}

	case StepCompleteMsg:
		if msg.Success {
			m.steps[msg.StepIndex].Status = StatusComplete
		} else {
			m.steps[msg.StepIndex].Status = StatusError
			m.state = StateError
			m.errorMsg = msg.Error
			m.showError = true
			return m, nil
		}

		// Move to next step or complete installation
		nextStep := m.findNextExecutableStep(msg.StepIndex + 1)
		if nextStep != -1 {
			m.currentStep = nextStep
			m.steps[m.currentStep].Status = StatusRunning
			cmds = append(cmds, m.runCurrentStep())
		} else {
			cmds = append(cmds, func() tea.Msg {
				return InstallCompleteMsg{Success: true}
			})
		}

	case ExistingInstallationMsg:
		if msg.Found {
			m.existingConfig = &msg.Config
			m.isUpgrade = true
			m.config.InstallPath = msg.Config.InstallPath
			m.state = StateBackupConfirmation
		} else {
			m.state = StateConfiguration
		}

	case EnvironmentCheckCompleteMsg:
		if msg.Success {
			m.state = StateExistingDetection
			cmds = append(cmds, m.checkForExistingInstallation())
		} else {
			m.state = StateError
			m.errorMsg = msg.Error
			m.showError = true
		}
	}

	// Update progress bar
	completedSteps := 0
	for _, step := range m.steps {
		if step.Status == StatusComplete {
			completedSteps++
		}
	}
	progressPercent := float64(completedSteps) / float64(m.totalSteps)
	m.progress.SetPercent(progressPercent)

	// Update other components
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m InstallerModel) View() string {
	switch m.state {
	case StateWelcome:
		return m.welcomeView()
	case StatePathSelection:
		return m.pathSelectionView()
	case StateExistingDetection:
		return m.existingDetectionView()
	case StateBackupConfirmation:
		return m.backupConfirmationView()
	case StateEnvironmentCheck:
		return m.environmentCheckView()
	case StateConfiguration:
		return m.configurationView()
	case StateFeatureSelection:
		return m.featureSelectionView()
	case StateInstallation:
		return m.installationView()
	case StateShellIntegration:
		return m.shellIntegrationView()
	case StatePermissions:
		return m.permissionsView()
	case StateCompletion:
		return m.completionView()
	case StateError:
		return m.errorView()
	default:
		return "Unknown state"
	}
}

func (m InstallerModel) welcomeView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Align(lipgloss.Center).
		Render("🚀 MCF - Multi Component Framework")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Align(lipgloss.Center).
		Render("Welcome to the Interactive Installer")

	description := `
MCF is a sophisticated development automation platform featuring:

• 🧠 AI Agent System with 9 specialized development agents
• ⚡ 50+ Custom Commands for comprehensive workflow automation  
• 🔧 Intelligent Hook System for event-driven automation
• 📚 Semantic Code Analysis with Serena MCP integration
• 🛡️ Advanced security features and input validation
• 🎯 Project templates and scaffolding tools

This installer will guide you through the complete setup process,
including customizable installation paths, feature selection, and
shell integration configuration.
`

	var statusMsg string
	if m.isUpgrade {
		statusMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Render("🔄 Existing installation detected - this will upgrade your MCF setup")
	} else {
		statusMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Render("🆕 Fresh installation - setting up MCF for the first time")
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press Enter to continue, q to quit")

	return lipgloss.JoinVertical(lipgloss.Center,
		"",
		title,
		"",
		subtitle,
		"",
		description,
		"",
		statusMsg,
		"",
		instructions,
		"",
	)
}

func (m InstallerModel) pathSelectionView() string {
	title := lipgloss.NewStyle().Bold(true).Render("📁 Installation Directory")

	current := fmt.Sprintf("Current path: %s", m.installPath)

	description := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("Choose where MCF will be installed. Default is ~/mcf")

	textInput := fmt.Sprintf("Path: %s", m.textInput.View())

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Enter custom path or press Enter to use default, Esc to go back")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		current,
		"",
		description,
		"",
		textInput,
		"",
		help,
	)
}

func (m InstallerModel) existingDetectionView() string {
	title := lipgloss.NewStyle().Bold(true).Render("🔍 Scanning for Existing Installation")

	spinner := m.spinner.View()
	status := fmt.Sprintf("%s Checking for existing MCF installations...", spinner)

	paths := []string{
		"• ~/mcf",
		"• ~/.mcf",
		"• Current directory",
		"• System-wide installations",
	}

	searchList := lipgloss.JoinVertical(lipgloss.Left, paths...)

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		status,
		"",
		"Scanning locations:",
		searchList,
	)
}

func (m InstallerModel) backupConfirmationView() string {
	title := lipgloss.NewStyle().Bold(true).Render("💾 Backup Existing Configuration")

	var existingInfo string
	if m.existingConfig != nil {
		existingInfo = fmt.Sprintf(`
Found existing MCF installation:
• Version: %s
• Location: %s
• Installed: %s
• Features: %s

A backup will be created before proceeding with the upgrade.
`,
			m.existingConfig.Version,
			m.existingConfig.InstallPath,
			m.existingConfig.InstalledAt.Format("2006-01-02 15:04:05"),
			strings.Join(m.existingConfig.Features, ", "),
		)
	}

	warning := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Render("⚠️  This will backup your current configuration and update MCF")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press Enter to create backup and continue, Esc to cancel")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		existingInfo,
		warning,
		"",
		instructions,
	)
}

func (m InstallerModel) environmentCheckView() string {
	title := lipgloss.NewStyle().Bold(true).Render("🔧 Environment Check")

	var items []string

	checks := []struct {
		name   string
		status string
	}{
		{"Operating System", runtime.GOOS},
		{"Architecture", runtime.GOARCH},
		{"Git availability", "checking..."},
		{"Curl availability", "checking..."},
		{"Claude Code CLI", "checking..."},
		{"Shell environment", "detecting..."},
		{"File permissions", "verifying..."},
	}

	for _, check := range checks {
		var icon string
		if check.status == "checking..." || check.status == "detecting..." || check.status == "verifying..." {
			icon = m.spinner.View()
		} else {
			icon = "✅"
		}
		items = append(items, fmt.Sprintf("%s %s: %s", icon, check.name, check.status))
	}

	checkList := lipgloss.JoinVertical(lipgloss.Left, items...)

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		checkList,
	)
}

func (m InstallerModel) configurationView() string {
	title := lipgloss.NewStyle().Bold(true).Render("⚙️ Configuration")

	if m.currentField >= len(m.configFields) {
		return "Configuration complete"
	}

	field := m.configFields[m.currentField]

	fieldTitle := fmt.Sprintf("Step %d of %d: %s", m.currentField+1, len(m.configFields), field.Label)

	description := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render(field.Description)

	var input string
	if field.Type == "password" {
		// Show asterisks for password fields
		input = strings.Repeat("*", len(m.textInput.Value()))
		input = fmt.Sprintf("%s: %s", field.Label, input)
	} else {
		input = fmt.Sprintf("%s: %s", field.Label, m.textInput.View())
	}

	var defaultInfo string
	if field.Default != "" {
		defaultInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true).
			Render(fmt.Sprintf("(default: %s)", field.Default))
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Enter value or press Tab to use default, Enter to continue")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		fieldTitle,
		"",
		description,
		"",
		input,
		defaultInfo,
		"",
		help,
	)
}

func (m InstallerModel) featureSelectionView() string {
	title := lipgloss.NewStyle().Bold(true).Render("🎯 Feature Selection")

	features := []struct {
		name        string
		description string
		enabled     bool
		required    bool
	}{
		{"🧠 AI Agent System", "Core development agents (Required)", true, true},
		{"🔧 Git Integration", "Automated workflows and hooks", true, false},
		{"📊 Serena MCP Server", "Advanced semantic code analysis", false, false},
		{"📚 Template System", "Project scaffolding and templates", true, false},
		{"🛡️ Security Features", "Input validation and audit logging", true, false},
		{"📈 Performance Monitoring", "Resource usage tracking", false, false},
		{"🌐 Web Dashboard", "Browser-based project management", false, false},
		{"🔄 Auto-Updates", "Automatic MCF updates", true, false},
	}

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("Select features to install (Space to toggle, Enter to continue)")

	var featureList []string
	for i, feature := range features {
		var prefix string
		var style lipgloss.Style

		if feature.required {
			prefix = "✓"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
		} else if feature.enabled {
			prefix = "☑"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
		} else {
			prefix = "☐"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
		}

		line := fmt.Sprintf("  %s %s", prefix, feature.name)
		if feature.required {
			line += " (Required)"
		}

		featureList = append(featureList, style.Render(line))

		// Add description
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)
		featureList = append(featureList, descStyle.Render(fmt.Sprintf("      %s", feature.description)))

		if i < len(features)-1 {
			featureList = append(featureList, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		subtitle,
		"",
		lipgloss.JoinVertical(lipgloss.Left, featureList...),
	)
}

func (m InstallerModel) installationView() string {
	title := lipgloss.NewStyle().Bold(true).Render("📦 Installing MCF")

	// Progress bar
	progressBar := m.progress.View()
	completedSteps := 0
	for _, step := range m.steps {
		if step.Status == StatusComplete {
			completedSteps++
		}
	}
	progressText := fmt.Sprintf(" %d/%d steps completed", completedSteps, m.totalSteps)

	// Current step info
	var currentStepInfo string
	if m.currentStep < len(m.steps) {
		currentStep := m.steps[m.currentStep]
		currentStepInfo = fmt.Sprintf("%s %s", m.spinner.View(), currentStep.Name)
		if currentStep.Description != "" {
			currentStepInfo += fmt.Sprintf("\n%s", lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Italic(true).
				Render(fmt.Sprintf("    %s", currentStep.Description)))
		}
	}

	// Step list
	var stepList []string
	for _, step := range m.steps {
		var status, color string

		switch step.Status {
		case StatusPending:
			status = "⏳"
			color = "243"
		case StatusRunning:
			status = "🔄"
			color = "99"
		case StatusComplete:
			status = "✅"
			color = "34"
		case StatusError:
			status = "❌"
			color = "196"
		case StatusSkipped:
			status = "⏭️"
			color = "220"
		}

		line := fmt.Sprintf("  %s %s", status, step.Name)
		if step.Status == StatusError {
			line += " (Failed)"
		} else if step.Status == StatusSkipped {
			line += " (Skipped)"
		}

		stepList = append(stepList, lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		progressBar+progressText,
		"",
		currentStepInfo,
		"",
		"Installation Steps:",
		lipgloss.JoinVertical(lipgloss.Left, stepList...),
	)
}

func (m InstallerModel) shellIntegrationView() string {
	title := lipgloss.NewStyle().Bold(true).Render("🐚 Shell Integration")

	description := `Configure shell integration for MCF:

• Add MCF to your PATH
• Create convenient aliases (mcf, claude-mcf)
• Set up completion scripts
• Configure environment variables`

	options := []string{
		"✓ Add to PATH",
		"✓ Create aliases",
		"☐ Install completions",
		"☐ Set global env vars",
	}

	optionsList := lipgloss.JoinVertical(lipgloss.Left, options...)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press Enter to apply settings, Space to toggle options")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		description,
		"",
		optionsList,
		"",
		instructions,
	)
}

func (m InstallerModel) permissionsView() string {
	title := lipgloss.NewStyle().Bold(true).Render("🔐 Setting File Permissions")

	items := []string{
		fmt.Sprintf("%s Setting executable permissions...", m.spinner.View()),
		"  • MCF binary files",
		"  • Shell scripts and hooks",
		"  • Template generators",
		"",
		fmt.Sprintf("%s Configuring directory permissions...", m.spinner.View()),
		"  • Configuration directories",
		"  • Cache and log directories",
		"  • Project template directories",
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}

func (m InstallerModel) completionView() string {
	success := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("34")).
		Render("🎉 MCF Installation Complete!")

	var installType string
	if m.isUpgrade {
		installType = "upgraded"
	} else {
		installType = "installed"
	}

	summary := fmt.Sprintf(`
MCF has been successfully %s and configured!

Installation Details:
• Location: %s
• Version: %s
• Features: %s
• Shell Integration: %s

Quick Start Commands:
  mcf --help                 Show all available commands
  mcf init                   Initialize MCF in current project
  mcf agent:list            List available AI agents
  mcf template:create       Create new project from template
  
Advanced Features:
  mcf serena:analyze        Run semantic code analysis
  mcf security:audit        Security audit and validation
  mcf performance:monitor   Resource usage monitoring

Configuration:
  Edit ~/.mcf/config.yaml to customize settings
  Visit %s/.claude/ for project-specific configuration

`,
		installType,
		m.config.InstallPath,
		m.config.Version,
		strings.Join(m.config.Features, ", "),
		func() string {
			if m.config.Shell.AddToPath {
				return "Enabled"
			}
			return "Manual setup required"
		}(),
		m.config.InstallPath,
	)

	var backupInfo string
	if m.isUpgrade && m.backupPath != "" {
		backupInfo = fmt.Sprintf("\n💾 Backup created: %s\n", m.backupPath)
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press any key to exit")

	return lipgloss.JoinVertical(lipgloss.Left,
		success,
		summary,
		backupInfo,
		instructions,
	)
}

func (m InstallerModel) errorView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")).
		Render("❌ Installation Error")

	errorDetails := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render(fmt.Sprintf("Error: %s", m.errorMsg))

	troubleshooting := `
Troubleshooting:
• Check your internet connection
• Ensure you have write permissions to the install directory
• Verify Claude Code CLI is properly installed
• Check system requirements (Git, Curl, etc.)

For help:
• Visit: https://github.com/your-repo/mcf/issues
• Run: mcf doctor (if partially installed)
• Check logs in ~/.mcf/logs/installer.log
`

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'r' to retry installation, 'q' to quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		errorDetails,
		troubleshooting,
		instructions,
	)
}

// Event messages
type InstallCompleteMsg struct {
	Success bool
	Error   string
}

type StepCompleteMsg struct {
	StepIndex int
	Success   bool
	Error     string
}

type ExistingInstallationMsg struct {
	Found  bool
	Config MCFConfig
}

type EnvironmentCheckCompleteMsg struct {
	Success bool
	Error   string
	Info    EnvironmentInfo
}

// Input handlers
func (m InstallerModel) handleWelcomeInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = StateEnvironmentCheck
		return m, m.startEnvironmentCheck()
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m InstallerModel) handlePathSelectionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Use current text or default
		path := strings.TrimSpace(m.textInput.Value())
		if path == "" {
			path = m.installPath
		}

		// Validate path
		if err := validateInstallPath(path); err != nil {
			m.errorMsg = err.Error()
			m.showError = true
			return m, nil
		}

		m.installPath = path
		m.config.InstallPath = path
		m.state = StateEnvironmentCheck
		return m, m.startEnvironmentCheck()
	case "esc":
		m.state = StateWelcome
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InstallerModel) handleExistingDetectionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// This state is automatic, just handle quit
	if msg.String() == "ctrl+c" || msg.String() == "q" {
		return m, tea.Quit
	}
	return m, nil
}

func (m InstallerModel) handleBackupConfirmationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = StateConfiguration
		return m, nil
	case "esc":
		m.state = StateWelcome
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m InstallerModel) handleConfigurationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Validate current field
		field := m.configFields[m.currentField]
		value := strings.TrimSpace(m.textInput.Value())

		// Use default if empty
		if value == "" {
			value = field.Default
		}

		// Validate
		if field.Validation != nil {
			if err := field.Validation(value); err != nil {
				m.errorMsg = err.Error()
				m.showError = true
				return m, nil
			}
		}

		// Apply value to config
		m.applyConfigField(field.Key, value)

		// Move to next field or complete
		m.currentField++
		if m.currentField >= len(m.configFields) {
			m.state = StateFeatureSelection
			return m, nil
		}

		// Reset input for next field
		m.textInput.SetValue("")
		if m.configFields[m.currentField].Default != "" {
			m.textInput.SetValue(m.configFields[m.currentField].Default)
		}
		return m, nil

	case "tab":
		// Use default value
		if m.currentField < len(m.configFields) {
			field := m.configFields[m.currentField]
			m.textInput.SetValue(field.Default)
		}
		return m, nil

	case "esc":
		if m.currentField > 0 {
			m.currentField--
			// Restore previous field value
			field := m.configFields[m.currentField]
			m.textInput.SetValue(field.Default)
		} else {
			m.state = StateWelcome
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InstallerModel) handleFeatureSelectionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = StateInstallation
		return m, m.startInstallation()
	case "space":
		// TODO: Toggle feature selection
		return m, nil
	case "esc":
		m.state = StateConfiguration
		m.currentField = len(m.configFields) - 1
		return m, nil
	}
	return m, nil
}

func (m InstallerModel) handleShellIntegrationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = StateInstallation
		return m, m.continueInstallation()
	case "space":
		// TODO: Toggle shell integration options
		return m, nil
	}
	return m, nil
}

func (m InstallerModel) handleErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		// Reset and retry
		m.state = StateWelcome
		m.currentStep = 0
		m.showError = false
		m.errorMsg = ""
		for i := range m.steps {
			m.steps[i].Status = StatusPending
		}
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m InstallerModel) handleCompletionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

// Helper functions
func (m *InstallerModel) applyConfigField(key, value string) {
	switch key {
	case "install_path":
		m.config.InstallPath = value
		m.installPath = value
	case "project_name":
		m.config.ProjectName = value
	case "editor":
		m.config.Preferences.Editor = value
	case "shell":
		m.config.Shell.Shell = value
	case "claude_api_key":
		m.config.ClaudeAPIKey = value
	}
}

func (m InstallerModel) findNextExecutableStep(startIndex int) int {
	for i := startIndex; i < len(m.steps); i++ {
		step := m.steps[i]
		// Check if step should be executed based on condition
		if step.Condition == nil || step.Condition(&m.config) {
			return i
		} else {
			// Mark as skipped
			m.steps[i].Status = StatusSkipped
		}
	}
	return -1
}

func (m InstallerModel) startEnvironmentCheck() tea.Cmd {
	return func() tea.Msg {
		// Perform environment checks
		info := EnvironmentInfo{
			OS:      runtime.GOOS,
			Shell:   os.Getenv("SHELL"),
			HomeDir: os.Getenv("HOME"),
		}

		// Check for required tools
		var err error
		_, err = exec.LookPath("git")
		info.HasGit = err == nil
		_, err = exec.LookPath("curl")
		info.HasCurl = err == nil

		// Check Claude CLI
		if path, err := exec.LookPath("claude"); err == nil {
			info.HasClaude = true
			// Try to get version
			if cmd := exec.Command(path, "--version"); cmd != nil {
				if output, err := cmd.Output(); err == nil {
					info.ClaudeVersion = strings.TrimSpace(string(output))
				}
			}
		}

		// Check other tools
		if cmd := exec.Command("python", "--version"); cmd != nil {
			if output, err := cmd.Output(); err == nil {
				info.PythonVersion = strings.TrimSpace(string(output))
			}
		}

		if cmd := exec.Command("node", "--version"); cmd != nil {
			if output, err := cmd.Output(); err == nil {
				info.NodeVersion = strings.TrimSpace(string(output))
			}
		}

		if cmd := exec.Command("go", "version"); cmd != nil {
			if output, err := cmd.Output(); err == nil {
				info.GoVersion = strings.TrimSpace(string(output))
			}
		}

		// Validate requirements
		var errors []string
		if !info.HasGit {
			errors = append(errors, "Git is required but not found")
		}
		if !info.HasCurl {
			errors = append(errors, "Curl is required but not found")
		}
		if !info.HasClaude {
			errors = append(errors, "Claude Code CLI is required but not found")
		}

		if len(errors) > 0 {
			return EnvironmentCheckCompleteMsg{
				Success: false,
				Error:   strings.Join(errors, ", "),
				Info:    info,
			}
		}

		return EnvironmentCheckCompleteMsg{
			Success: true,
			Info:    info,
		}
	}
}

func (m InstallerModel) checkForExistingInstallation() tea.Cmd {
	return func() tea.Msg {
		// Check multiple possible locations
		locations := []string{
			m.installPath,
			filepath.Join(os.Getenv("HOME"), ".mcf"),
			filepath.Join(os.Getenv("HOME"), "mcf"),
		}

		for _, location := range locations {
			configPath := filepath.Join(location, ".claude", "config.yaml")
			if _, err := os.Stat(configPath); err == nil {
				// Found existing installation
				var config MCFConfig
				if data, err := os.ReadFile(configPath); err == nil {
					if err := yaml.Unmarshal(data, &config); err == nil {
						return ExistingInstallationMsg{
							Found:  true,
							Config: config,
						}
					}
				}
			}
		}

		return ExistingInstallationMsg{Found: false}
	}
}

func (m InstallerModel) startInstallation() tea.Cmd {
	return m.runCurrentStep()
}

func (m InstallerModel) continueInstallation() tea.Cmd {
	return m.runCurrentStep()
}

func (m InstallerModel) runCurrentStep() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Find next executable step
		stepIndex := m.findNextExecutableStep(m.currentStep)
		if stepIndex == -1 {
			return InstallCompleteMsg{Success: true}
		}

		step := m.steps[stepIndex]
		err := step.Action(ctx, &m)

		return StepCompleteMsg{
			StepIndex: stepIndex,
			Success:   err == nil,
			Error: func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
		}
	}
}

// Step action implementations
func checkEnvironment(ctx context.Context, m *InstallerModel) error {
	// This is handled in startEnvironmentCheck, just return success
	return nil
}

func detectExistingInstallation(ctx context.Context, m *InstallerModel) error {
	// This is handled separately, just return success
	return nil
}

func backupExistingConfig(ctx context.Context, m *InstallerModel) error {
	if m.existingConfig == nil {
		return nil // Nothing to backup
	}

	// Create backup directory
	backupDir := filepath.Join(m.existingConfig.InstallPath, "backup")
	timestamp := time.Now().Format("20060102-150405")
	m.backupPath = filepath.Join(backupDir, fmt.Sprintf("mcf-backup-%s", timestamp))

	if err := os.MkdirAll(m.backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Backup .claude directory
	srcDir := filepath.Join(m.existingConfig.InstallPath, ".claude")
	dstDir := filepath.Join(m.backupPath, ".claude")

	return copyDir(srcDir, dstDir)
}

func createDirectoryStructure(ctx context.Context, m *InstallerModel) error {
	dirs := []string{
		m.config.InstallPath,
		filepath.Join(m.config.InstallPath, ".claude"),
		filepath.Join(m.config.InstallPath, ".claude", "agents"),
		filepath.Join(m.config.InstallPath, ".claude", "commands"),
		filepath.Join(m.config.InstallPath, ".claude", "commands", "context"),
		filepath.Join(m.config.InstallPath, ".claude", "commands", "gh"),
		filepath.Join(m.config.InstallPath, ".claude", "commands", "project"),
		filepath.Join(m.config.InstallPath, ".claude", "commands", "serena"),
		filepath.Join(m.config.InstallPath, ".claude", "commands", "templates"),
		filepath.Join(m.config.InstallPath, ".claude", "hooks"),
		filepath.Join(m.config.InstallPath, ".claude", "templates"),
		filepath.Join(m.config.InstallPath, ".claude", "cache"),
		filepath.Join(m.config.InstallPath, ".claude", "logs"),
		filepath.Join(m.config.InstallPath, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func setupClaudeCodeIntegration(ctx context.Context, m *InstallerModel) error {
	// Create Claude Code configuration
	configPath := filepath.Join(m.config.InstallPath, ".claude", "config.yaml")

	config := map[string]interface{}{
		"version":      m.config.Version,
		"install_path": m.config.InstallPath,
		"project_name": m.config.ProjectName,
		"features":     m.config.Features,
		"created_at":   time.Now().Format(time.RFC3339),
		"preferences":  m.config.Preferences,
		"shell":        m.config.Shell,
		"security":     m.config.Security,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	// Create settings.json for Claude Code
	settingsPath := filepath.Join(m.config.InstallPath, ".claude", "settings.json")
	settings := map[string]interface{}{
		"claude_api_key": m.config.ClaudeAPIKey,
		"model":          "claude-3-5-sonnet-20241022",
		"max_tokens":     8192,
		"temperature":    0.7,
		"timeout":        30,
		"retry_attempts": 3,
	}

	settingsData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, settingsData, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

func installSerenaMCP(ctx context.Context, m *InstallerModel) error {
	// Create Serena MCP configuration
	serenaDir := filepath.Join(m.config.InstallPath, ".serena")
	if err := os.MkdirAll(serenaDir, 0755); err != nil {
		return fmt.Errorf("failed to create .serena directory: %w", err)
	}

	// Create basic Serena configuration
	serenaConfig := map[string]interface{}{
		"server": map[string]interface{}{
			"name":    "serena-mcp",
			"command": "npx",
			"args":    []string{"@serena-ai/mcp-server"},
		},
		"capabilities": []string{
			"semantic_analysis",
			"code_understanding",
			"project_mapping",
		},
	}

	configData, err := yaml.Marshal(serenaConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal Serena config: %w", err)
	}

	configPath := filepath.Join(serenaDir, "config.yaml")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write Serena config: %w", err)
	}

	return nil
}

func installProjectTemplates(ctx context.Context, m *InstallerModel) error {
	templatesDir := filepath.Join(m.config.InstallPath, ".claude", "templates")

	// Create basic project templates
	templates := map[string]string{
		"basic/mcf-project.yaml": `
name: "{{.ProjectName}}"
type: "mcf-project"
description: "Basic MCF project template"
files:
  - path: ".claude/agents/"
    type: "directory"
  - path: ".claude/commands/"
    type: "directory"
  - path: "README.md"
    content: |
      # {{.ProjectName}}
      
      MCF-enabled project with AI agent system integration.
`,
		"web/next-mcf.yaml": `
name: "{{.ProjectName}}"
type: "next-mcf"
description: "Next.js project with MCF integration"
dependencies:
  - "next"
  - "react"
  - "react-dom"
files:
  - path: "package.json"
    content: |
      {
        "name": "{{.ProjectName}}",
        "scripts": {
          "dev": "next dev",
          "build": "next build",
          "mcf": "claude --project ."
        }
      }
`,
	}

	for templatePath, content := range templates {
		fullPath := filepath.Join(templatesDir, templatePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create template directory: %w", err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write template %s: %w", templatePath, err)
		}
	}

	return nil
}

func installCommandSystem(ctx context.Context, m *InstallerModel) error {
	commandsDir := filepath.Join(m.config.InstallPath, ".claude", "commands")

	// Create command scripts
	commands := map[string]string{
		"mcf-init.sh": `#!/bin/bash
# MCF project initialization script
set -e

PROJECT_NAME="${1:-$(basename $(pwd))}"
echo "🚀 Initializing MCF project: $PROJECT_NAME"

# Create .claude directory structure
mkdir -p .claude/{agents,commands,hooks,templates}

# Create basic configuration
cat > .claude/config.yaml << EOF
name: "$PROJECT_NAME"
type: "mcf-project"
features:
  - "ai-agents"
  - "command-system"
created_at: "$(date -Iseconds)"
EOF

echo "✅ MCF project initialized successfully!"
`,
		"mcf-doctor.sh": `#!/bin/bash
# MCF health check script
set -e

echo "🩺 MCF Health Check"
echo "==================="

# Check Claude CLI
if command -v claude >/dev/null 2>&1; then
    echo "✅ Claude CLI: $(claude --version)"
else
    echo "❌ Claude CLI: Not found"
fi

# Check configuration
if [[ -f .claude/config.yaml ]]; then
    echo "✅ MCF Configuration: Found"
else
    echo "❌ MCF Configuration: Missing"
fi

# Check directory structure
for dir in .claude/{agents,commands,hooks,templates}; do
    if [[ -d "$dir" ]]; then
        echo "✅ Directory $dir: Found"
    else
        echo "⚠️  Directory $dir: Missing"
    fi
done

echo "==================="
echo "🩺 Health check complete"
`,
	}

	for cmdName, content := range commands {
		cmdPath := filepath.Join(commandsDir, cmdName)
		if err := os.WriteFile(cmdPath, []byte(content), 0755); err != nil {
			return fmt.Errorf("failed to write command %s: %w", cmdName, err)
		}
	}

	return nil
}

func configureShellIntegration(ctx context.Context, m *InstallerModel) error {
	if !m.config.Shell.AddToPath && !m.config.Shell.CreateAlias {
		return nil // Nothing to configure
	}

	// Detect shell configuration file
	shell := m.config.Shell.Shell
	if shell == "auto-detect" {
		shell = filepath.Base(os.Getenv("SHELL"))
	}

	var rcFiles []string
	switch shell {
	case "bash":
		rcFiles = []string{".bashrc", ".bash_profile"}
	case "zsh":
		rcFiles = []string{".zshrc"}
	case "fish":
		rcFiles = []string{".config/fish/config.fish"}
	default:
		rcFiles = []string{".profile"}
	}

	homeDir := os.Getenv("HOME")
	var rcFile string

	// Find existing rc file
	for _, file := range rcFiles {
		path := filepath.Join(homeDir, file)
		if _, err := os.Stat(path); err == nil {
			rcFile = path
			break
		}
	}

	// Use first option if none found
	if rcFile == "" {
		rcFile = filepath.Join(homeDir, rcFiles[0])
	}

	// Prepare shell integration content
	var shellContent strings.Builder
	shellContent.WriteString("\n# MCF - Multi Component Framework")
	shellContent.WriteString("\n# Added by MCF installer")

	if m.config.Shell.AddToPath {
		binPath := filepath.Join(m.config.InstallPath, "bin")
		shellContent.WriteString(fmt.Sprintf("\nexport PATH=\"%s:$PATH\"", binPath))
	}

	if m.config.Shell.CreateAlias {
		mcfScript := filepath.Join(m.config.InstallPath, ".claude", "commands", "mcf-init.sh")
		for _, alias := range m.config.Shell.Aliases {
			shellContent.WriteString(fmt.Sprintf("\nalias %s='%s'", alias, mcfScript))
		}
	}

	shellContent.WriteString("\n# End MCF configuration\n")

	// Append to shell configuration file
	file, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open shell config file %s: %w", rcFile, err)
	}
	defer file.Close()

	if _, err := file.WriteString(shellContent.String()); err != nil {
		return fmt.Errorf("failed to write shell configuration: %w", err)
	}

	m.config.Shell.RCFile = rcFile
	return nil
}

func setFilePermissions(ctx context.Context, m *InstallerModel) error {
	// Set executable permissions on scripts
	scriptsDir := filepath.Join(m.config.InstallPath, ".claude", "commands")

	return filepath.Walk(scriptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".sh") || strings.HasSuffix(path, ".py")) {
			return os.Chmod(path, 0755)
		}
		return nil
	})
}

func validateInstallation(ctx context.Context, m *InstallerModel) error {
	// Verify all required directories exist
	requiredDirs := []string{
		filepath.Join(m.config.InstallPath, ".claude"),
		filepath.Join(m.config.InstallPath, ".claude", "commands"),
		filepath.Join(m.config.InstallPath, ".claude", "agents"),
	}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("required directory missing: %s", dir)
		}
	}

	// Verify configuration file
	configPath := filepath.Join(m.config.InstallPath, ".claude", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file missing: %s", configPath)
	}

	// Update final timestamps
	m.config.InstalledAt = time.Now()
	m.config.UpdatedAt = time.Now()

	return nil
}

// Validation functions
func validateInstallPath(path string) error {
	if path == "" {
		return fmt.Errorf("install path cannot be empty")
	}

	// Expand environment variables and ~
	expanded := os.ExpandEnv(path)
	if strings.HasPrefix(expanded, "~") {
		homeDir := os.Getenv("HOME")
		expanded = filepath.Join(homeDir, expanded[1:])
	}

	// Check if path is absolute
	if !filepath.IsAbs(expanded) {
		return fmt.Errorf("install path must be absolute")
	}

	// Check if parent directory exists and is writable
	parent := filepath.Dir(expanded)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		return fmt.Errorf("parent directory does not exist: %s", parent)
	}

	// Test write permissions by creating a temporary file
	testFile := filepath.Join(parent, ".mcf-install-test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("no write permission to directory: %s", parent)
	}
	os.Remove(testFile)

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if len(name) > 50 {
		return fmt.Errorf("project name cannot exceed 50 characters")
	}
	// Simple validation - letters, numbers, hyphens, underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
		}
	}
	return nil
}

func validateEditor(editor string) error {
	// Accept common editors or auto-detect
	validEditors := []string{"auto-detect", "vscode", "vim", "neovim", "emacs", "sublime", "atom"}
	for _, valid := range validEditors {
		if editor == valid {
			return nil
		}
	}
	return fmt.Errorf("unsupported editor: %s", editor)
}

func validateShell(shell string) error {
	// Accept common shells or auto-detect
	validShells := []string{"auto-detect", "bash", "zsh", "fish", "sh"}
	for _, valid := range validShells {
		if shell == valid {
			return nil
		}
	}
	return fmt.Errorf("unsupported shell: %s", shell)
}

func validateClaudeAPIKey(key string) error {
	// Optional field, so empty is allowed
	if key == "" {
		return nil
	}
	// Basic validation - should start with 'sk-ant-'
	if !strings.HasPrefix(key, "sk-ant-") {
		return fmt.Errorf("Claude API key should start with 'sk-ant-'")
	}
	if len(key) < 20 {
		return fmt.Errorf("Claude API key appears to be too short")
	}
	return nil
}

// Utility functions
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
