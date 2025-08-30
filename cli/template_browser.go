package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// Template browser states
type TemplateBrowserState int

const (
	StateTemplateList TemplateBrowserState = iota
	StateTemplatePreview
	StateParameterForm
	StateInstallProgress
	StateInstallComplete
	StateTemplateError
)

// Template browser model
type TemplateBrowserModel struct {
	state            TemplateBrowserState
	width            int
	height           int
	templates        []TemplateInfo
	list             list.Model
	viewport         viewport.Model
	textInputs       []textinput.Model
	progress         progress.Model
	currentField     int
	selectedTemplate *TemplateInfo
	parameterValues  map[string]string
	installationLog  []string
	errorMsg         string
	ready            bool
}

// Template information structure
type TemplateInfo struct {
	Name         string                 `yaml:"name" json:"name"`
	Type         string                 `yaml:"type" json:"type"`
	Version      string                 `yaml:"version" json:"version"`
	Description  string                 `yaml:"description" json:"description"`
	Author       string                 `yaml:"author" json:"author"`
	License      string                 `yaml:"license" json:"license"`
	Tags         []string               `yaml:"tags" json:"tags"`
	Files        []TemplateFile         `yaml:"files" json:"files"`
	Parameters   []TemplateParameter    `yaml:"parameters" json:"parameters"`
	Dependencies []string               `yaml:"dependencies" json:"dependencies"`
	Scripts      map[string]string      `yaml:"scripts" json:"scripts"`
	Metadata     map[string]interface{} `yaml:"metadata" json:"metadata"`
	FilePath     string                 `json:"file_path"` // Internal use
}

type TemplateFile struct {
	Path        string `yaml:"path" json:"path"`
	Type        string `yaml:"type" json:"type"` // file, directory, symlink
	Content     string `yaml:"content" json:"content"`
	Template    bool   `yaml:"template" json:"template"`
	Permissions string `yaml:"permissions" json:"permissions"`
	Executable  bool   `yaml:"executable" json:"executable"`
}

type TemplateParameter struct {
	Key         string      `yaml:"key" json:"key"`
	Label       string      `yaml:"label" json:"label"`
	Type        string      `yaml:"type" json:"type"`
	Description string      `yaml:"description" json:"description"`
	Required    bool        `yaml:"required" json:"required"`
	Default     interface{} `yaml:"default" json:"default"`
	Options     []string    `yaml:"options" json:"options"`
	Placeholder string      `yaml:"placeholder" json:"placeholder"`
	Validation  string      `yaml:"validation" json:"validation"`
	HelpText    string      `yaml:"help_text" json:"help_text"`
}

// List item for template
type templateItem struct {
	template TemplateInfo
}

func (i templateItem) Title() string       { return i.template.Name }
func (i templateItem) Description() string { return i.template.Description }
func (i templateItem) FilterValue() string { return i.template.Name + " " + i.template.Description }

// Create new template browser
func NewTemplateBrowserModel() TemplateBrowserModel {
	// Initialize list
	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "MCF Templates"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true).
		Margin(0, 0, 1, 2)

	// Initialize viewport for preview
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(1, 2)

	// Initialize progress bar
	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = 50

	return TemplateBrowserModel{
		state:           StateTemplateList,
		list:            l,
		viewport:        vp,
		progress:        prog,
		parameterValues: make(map[string]string),
		installationLog: []string{},
		ready:           false,
	}
}

func (m TemplateBrowserModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadTemplates(),
		textinput.Blink,
	)
}

func (m TemplateBrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listHeight := m.height - 8
		if listHeight < 10 {
			listHeight = 10
		}
		m.list.SetSize(msg.Width-4, listHeight)

		viewportWidth := msg.Width - 6
		viewportHeight := m.height - 10
		if viewportHeight < 5 {
			viewportHeight = 5
		}
		m.viewport.Width = viewportWidth
		m.viewport.Height = viewportHeight

		if !m.ready {
			m.ready = true
		}

	case TemplatesLoadedMsg:
		m.templates = msg.Templates
		items := make([]list.Item, len(m.templates))
		for i, template := range m.templates {
			items[i] = templateItem{template: template}
		}
		m.list.SetItems(items)

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case InstallProgressMsg:
		m.installationLog = append(m.installationLog, msg.Message)
		m.progress.SetPercent(msg.Progress)

	case InstallCompleteMsg:
		if msg.Success {
			m.state = StateInstallComplete
		} else {
			m.state = StateTemplateError
			m.errorMsg = msg.Error
		}

	default:
		switch m.state {
		case StateTemplateList:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case StateTemplatePreview:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		case StateParameterForm:
			if len(m.textInputs) > 0 && m.currentField < len(m.textInputs) {
				var cmd tea.Cmd
				m.textInputs[m.currentField], cmd = m.textInputs[m.currentField].Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m TemplateBrowserModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateTemplateList:
		switch msg.String() {
		case "enter":
			if selected := m.list.SelectedItem(); selected != nil {
				item := selected.(templateItem)
				m.selectedTemplate = &item.template
				m.state = StateTemplatePreview
				m.viewport.SetContent(m.renderTemplatePreview())
				return m, nil
			}
		case "esc", "q":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

	case StateTemplatePreview:
		switch msg.String() {
		case "enter", "i":
			// Start installation process
			if m.selectedTemplate != nil && len(m.selectedTemplate.Parameters) > 0 {
				m.state = StateParameterForm
				m.initializeParameterForm()
			} else {
				m.state = StateInstallProgress
				return m, m.startInstallation()
			}
		case "esc", "backspace":
			m.state = StateTemplateList
		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case StateParameterForm:
		switch msg.String() {
		case "enter":
			// Move to next field or start installation
			if m.currentField < len(m.textInputs)-1 {
				m.textInputs[m.currentField].Blur()
				m.currentField++
				m.textInputs[m.currentField].Focus()
			} else {
				// Collect parameter values
				m.collectParameterValues()
				m.state = StateInstallProgress
				return m, m.startInstallation()
			}
		case "shift+tab", "up":
			if m.currentField > 0 {
				m.textInputs[m.currentField].Blur()
				m.currentField--
				m.textInputs[m.currentField].Focus()
			}
		case "tab", "down":
			if m.currentField < len(m.textInputs)-1 {
				m.textInputs[m.currentField].Blur()
				m.currentField++
				m.textInputs[m.currentField].Focus()
			}
		case "esc":
			m.state = StateTemplatePreview
		default:
			var cmd tea.Cmd
			if m.currentField < len(m.textInputs) {
				m.textInputs[m.currentField], cmd = m.textInputs[m.currentField].Update(msg)
				return m, cmd
			}
		}

	case StateInstallProgress:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	case StateInstallComplete, StateTemplateError:
		switch msg.String() {
		case "enter", "esc", "q":
			m.state = StateTemplateList
			m.resetInstallation()
		case "r":
			if m.state == StateTemplateError {
				m.state = StateInstallProgress
				m.resetInstallation()
				return m, m.startInstallation()
			}
		}
	}

	return m, nil
}

func (m TemplateBrowserModel) View() string {
	if !m.ready {
		return "Loading templates..."
	}

	switch m.state {
	case StateTemplateList:
		return m.renderTemplateList()
	case StateTemplatePreview:
		return m.renderTemplatePreviewView()
	case StateParameterForm:
		return m.renderParameterForm()
	case StateInstallProgress:
		return m.renderInstallProgress()
	case StateInstallComplete:
		return m.renderInstallComplete()
	case StateTemplateError:
		return m.renderInstallError()
	default:
		return "Unknown state"
	}
}

func (m TemplateBrowserModel) renderTemplateList() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Align(lipgloss.Center)

	header := headerStyle.Render("🧩 MCF Template Browser")

	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Found %d templates • Enter to preview • q to quit", len(m.templates)))

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		header,
		"",
		m.list.View(),
		"",
		statusBar,
	)

	return content
}

func (m TemplateBrowserModel) renderTemplatePreviewView() string {
	if m.selectedTemplate == nil {
		return "No template selected"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	title := titleStyle.Render(fmt.Sprintf("📋 Template Preview: %s", m.selectedTemplate.Name))

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Enter/i to install • Backspace to go back • ↑↓ to scroll")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		m.viewport.View(),
		"",
		instructions,
	)

	return content
}

func (m TemplateBrowserModel) renderTemplatePreview() string {
	if m.selectedTemplate == nil {
		return "No template selected"
	}

	t := m.selectedTemplate

	var content strings.Builder

	// Basic information
	content.WriteString(lipgloss.NewStyle().Bold(true).Render("📄 Template Information"))
	content.WriteString("\n\n")

	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))

	content.WriteString(labelStyle.Render("Name: ") + t.Name + "\n")
	content.WriteString(labelStyle.Render("Type: ") + t.Type + "\n")
	if t.Version != "" {
		content.WriteString(labelStyle.Render("Version: ") + t.Version + "\n")
	}
	if t.Author != "" {
		content.WriteString(labelStyle.Render("Author: ") + t.Author + "\n")
	}
	content.WriteString(labelStyle.Render("Description: ") + t.Description + "\n")

	// Tags
	if len(t.Tags) > 0 {
		content.WriteString("\n" + labelStyle.Render("Tags: "))
		tagStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

		for i, tag := range t.Tags {
			if i > 0 {
				content.WriteString(" ")
			}
			content.WriteString(tagStyle.Render(tag))
		}
		content.WriteString("\n")
	}

	// Parameters
	if len(t.Parameters) > 0 {
		content.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("⚙️ Parameters"))
		content.WriteString("\n\n")

		for _, param := range t.Parameters {
			paramName := param.Key
			if param.Required {
				paramName += " *"
			}

			content.WriteString(labelStyle.Render(paramName + ": "))
			content.WriteString(param.Description)

			if param.Default != nil {
				content.WriteString(infoStyle.Render(fmt.Sprintf(" (default: %v)", param.Default)))
			}
			content.WriteString("\n")
		}
	}

	// Dependencies
	if len(t.Dependencies) > 0 {
		content.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("📦 Dependencies"))
		content.WriteString("\n\n")
		for _, dep := range t.Dependencies {
			content.WriteString("• " + dep + "\n")
		}
	}

	// Files structure
	if len(t.Files) > 0 {
		content.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("📁 File Structure"))
		content.WriteString("\n\n")

		fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
		dirStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

		for _, file := range t.Files {
			var prefix string
			var style lipgloss.Style

			switch file.Type {
			case "directory":
				prefix = "📁 "
				style = dirStyle
			case "symlink":
				prefix = "🔗 "
				style = fileStyle
			default:
				prefix = "📄 "
				style = fileStyle
			}

			content.WriteString(prefix + style.Render(file.Path))
			if file.Template {
				content.WriteString(infoStyle.Render(" (template)"))
			}
			content.WriteString("\n")
		}
	}

	// Scripts
	if len(t.Scripts) > 0 {
		content.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("🔧 Scripts"))
		content.WriteString("\n\n")
		for name, script := range t.Scripts {
			content.WriteString(labelStyle.Render(name + ": "))
			// Truncate long scripts for preview
			if len(script) > 50 {
				script = script[:47] + "..."
			}
			content.WriteString(infoStyle.Render(script) + "\n")
		}
	}

	return content.String()
}

func (m TemplateBrowserModel) renderParameterForm() string {
	if m.selectedTemplate == nil {
		return "No template selected"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	title := titleStyle.Render(fmt.Sprintf("📝 Configure Template: %s", m.selectedTemplate.Name))

	var formItems []string

	for i, param := range m.selectedTemplate.Parameters {
		fieldStyle := lipgloss.NewStyle().Margin(1, 0)

		labelText := param.Label
		if param.Required {
			labelText += " *"
		}

		label := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Render(labelText)

		description := lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true).
			Render(param.Description)

		var input string
		if i < len(m.textInputs) {
			input = m.textInputs[i].View()
		}

		helpText := ""
		if param.HelpText != "" {
			helpText = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("  " + param.HelpText)
		}

		fieldContent := lipgloss.JoinVertical(lipgloss.Left,
			label,
			description,
			input,
			helpText,
		)

		formItems = append(formItems, fieldStyle.Render(fieldContent))
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Tab/↓↑ to navigate • Enter to install • Esc to go back")

	progress := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Render(fmt.Sprintf("Step %d of %d", m.currentField+1, len(m.selectedTemplate.Parameters)))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		progress,
		"",
		lipgloss.JoinVertical(lipgloss.Left, formItems...),
		"",
		instructions,
	)

	return content
}

func (m TemplateBrowserModel) renderInstallProgress() string {
	if m.selectedTemplate == nil {
		return "No template selected"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	title := titleStyle.Render(fmt.Sprintf("⚡ Installing Template: %s", m.selectedTemplate.Name))

	progressBar := m.progress.View()
	progressText := fmt.Sprintf("%.0f%% complete", m.progress.Percent()*100)

	// Installation log
	logStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252")).
		Padding(1).
		Width(m.width - 6).
		Height(m.height - 15)

	var logContent strings.Builder
	logLines := m.installationLog
	if len(logLines) > m.height-20 {
		logLines = logLines[len(logLines)-(m.height-20):]
	}

	for _, line := range logLines {
		logContent.WriteString(line + "\n")
	}

	log := logStyle.Render(logContent.String())

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Ctrl+C to cancel installation")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		progressBar+" "+progressText,
		"",
		log,
		"",
		instructions,
	)

	return content
}

func (m TemplateBrowserModel) renderInstallComplete() string {
	if m.selectedTemplate == nil {
		return "No template selected"
	}

	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("34"))

	title := successStyle.Render(fmt.Sprintf("✅ Template Installed Successfully: %s", m.selectedTemplate.Name))

	summaryStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252")).
		Padding(1).
		Margin(1, 0)

	var summary strings.Builder
	summary.WriteString("Installation Summary:\n\n")
	summary.WriteString(fmt.Sprintf("Template: %s\n", m.selectedTemplate.Name))
	summary.WriteString(fmt.Sprintf("Type: %s\n", m.selectedTemplate.Type))
	summary.WriteString(fmt.Sprintf("Files created: %d\n", len(m.selectedTemplate.Files)))

	if len(m.parameterValues) > 0 {
		summary.WriteString("\nParameters used:\n")
		for key, value := range m.parameterValues {
			summary.WriteString(fmt.Sprintf("• %s: %s\n", key, value))
		}
	}

	nextSteps := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render(`
Next steps:
• Review the created files and directories
• Run any post-installation scripts if needed
• Start developing with your new template!`)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Enter to return to template list • q to quit")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		summaryStyle.Render(summary.String()),
		nextSteps,
		"",
		instructions,
	)

	return content
}

func (m TemplateBrowserModel) renderInstallError() string {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	title := errorStyle.Render("❌ Template Installation Failed")

	errorDetails := lipgloss.NewStyle().
		Background(lipgloss.Color("52")).
		Foreground(lipgloss.Color("255")).
		Padding(1).
		Margin(1, 0).
		Render(m.errorMsg)

	troubleshooting := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render(`
Troubleshooting:
• Check file permissions in the target directory
• Ensure all dependencies are available
• Verify template configuration is valid
• Check available disk space`)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("r to retry • Enter to return to list • q to quit")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		errorDetails,
		troubleshooting,
		"",
		instructions,
	)

	return content
}

// Helper methods

func (m *TemplateBrowserModel) initializeParameterForm() {
	if m.selectedTemplate == nil {
		return
	}

	m.textInputs = make([]textinput.Model, len(m.selectedTemplate.Parameters))
	m.currentField = 0

	for i, param := range m.selectedTemplate.Parameters {
		ti := textinput.New()
		ti.Placeholder = param.Placeholder
		if param.Placeholder == "" && param.Default != nil {
			ti.Placeholder = fmt.Sprintf("%v", param.Default)
		}
		ti.CharLimit = 200
		ti.Width = 50

		if param.Type == "password" {
			ti.EchoMode = textinput.EchoPassword
		}

		if i == 0 {
			ti.Focus()
		}

		m.textInputs[i] = ti
	}
}

func (m *TemplateBrowserModel) collectParameterValues() {
	if m.selectedTemplate == nil {
		return
	}

	m.parameterValues = make(map[string]string)

	for i, param := range m.selectedTemplate.Parameters {
		if i < len(m.textInputs) {
			value := strings.TrimSpace(m.textInputs[i].Value())
			if value == "" && param.Default != nil {
				value = fmt.Sprintf("%v", param.Default)
			}
			m.parameterValues[param.Key] = value
		}
	}
}

func (m *TemplateBrowserModel) resetInstallation() {
	m.installationLog = []string{}
	m.progress.SetPercent(0)
	m.errorMsg = ""
	m.parameterValues = make(map[string]string)
	m.textInputs = []textinput.Model{}
	m.currentField = 0
}

// Commands

func (m TemplateBrowserModel) loadTemplates() tea.Cmd {
	return func() tea.Msg {
		templates, err := loadTemplatesFromDirectory()
		if err != nil {
			return TemplatesLoadedMsg{Templates: []TemplateInfo{}, Error: err.Error()}
		}
		return TemplatesLoadedMsg{Templates: templates}
	}
}

func (m TemplateBrowserModel) startInstallation() tea.Cmd {
	return func() tea.Msg {
		if m.selectedTemplate == nil {
			return InstallCompleteMsg{
				Success: false,
				Error:   "No template selected",
			}
		}

		// Simulate installation process
		go func() {
			// This would normally be replaced with actual template installation logic
			time.Sleep(100 * time.Millisecond)
			// Send progress updates...
			// For demonstration, we'll just complete successfully
		}()

		return InstallProgressMsg{
			Message:  "Starting template installation...",
			Progress: 0.0,
		}
	}
}

// Message types

type TemplatesLoadedMsg struct {
	Templates []TemplateInfo
	Error     string
}

type InstallProgressMsg struct {
	Message  string
	Progress float64
}

// Template loading functions

func loadTemplatesFromDirectory() ([]TemplateInfo, error) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	templatesDir := filepath.Join(wd, ".claude", "templates")

	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// Return built-in templates if directory doesn't exist
		return getBuiltInTemplates(), nil
	}

	var templates []TemplateInfo

	err = filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Look for template files (.yaml, .yml, .json)
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		template, err := loadTemplate(path)
		if err != nil {
			// Log error but continue loading other templates
			fmt.Printf("Warning: failed to load template %s: %v\n", path, err)
			return nil
		}

		template.FilePath = path
		templates = append(templates, template)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk templates directory: %w", err)
	}

	// Add built-in templates
	templates = append(templates, getBuiltInTemplates()...)

	// Sort templates by name
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

func loadTemplate(filePath string) (TemplateInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TemplateInfo{}, fmt.Errorf("failed to read template file: %w", err)
	}

	var template TemplateInfo
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		err = json.Unmarshal(data, &template)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &template)
	default:
		return TemplateInfo{}, fmt.Errorf("unsupported template format: %s", ext)
	}

	if err != nil {
		return TemplateInfo{}, fmt.Errorf("failed to parse template: %w", err)
	}

	// Validate required fields
	if template.Name == "" {
		template.Name = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}
	if template.Type == "" {
		template.Type = "general"
	}
	if template.Description == "" {
		template.Description = "No description available"
	}

	return template, nil
}

func getBuiltInTemplates() []TemplateInfo {
	return []TemplateInfo{
		{
			Name:        "MCF Basic Project",
			Type:        "mcf-project",
			Version:     "1.0.0",
			Description: "Basic MCF project with Claude Code integration",
			Author:      "MCF Team",
			Tags:        []string{"basic", "starter", "mcf"},
			Parameters: []TemplateParameter{
				{
					Key:         "project_name",
					Label:       "Project Name",
					Type:        "text",
					Description: "Name of your new project",
					Required:    true,
					Placeholder: "my-mcf-project",
				},
				{
					Key:         "description",
					Label:       "Description",
					Type:        "text",
					Description: "Brief description of your project",
					Required:    false,
					Placeholder: "An amazing MCF project",
				},
			},
			Files: []TemplateFile{
				{Path: ".claude", Type: "directory"},
				{Path: ".claude/agents", Type: "directory"},
				{Path: ".claude/commands", Type: "directory"},
				{Path: ".claude/hooks", Type: "directory"},
				{Path: "README.md", Type: "file", Template: true},
				{Path: ".gitignore", Type: "file"},
			},
		},
		{
			Name:        "Next.js + MCF",
			Type:        "web-nextjs",
			Version:     "1.0.0",
			Description: "Next.js project with MCF integration and Claude Code workflows",
			Author:      "MCF Team",
			Tags:        []string{"nextjs", "react", "web", "typescript"},
			Parameters: []TemplateParameter{
				{
					Key:         "project_name",
					Label:       "Project Name",
					Type:        "text",
					Description: "Name of your Next.js project",
					Required:    true,
					Placeholder: "my-next-app",
				},
				{
					Key:         "typescript",
					Label:       "Use TypeScript",
					Type:        "boolean",
					Description: "Enable TypeScript support",
					Default:     true,
				},
				{
					Key:         "tailwind",
					Label:       "Include Tailwind CSS",
					Type:        "boolean",
					Description: "Add Tailwind CSS for styling",
					Default:     true,
				},
			},
			Files: []TemplateFile{
				{Path: "package.json", Type: "file", Template: true},
				{Path: "next.config.js", Type: "file", Template: true},
				{Path: "app", Type: "directory"},
				{Path: "app/page.tsx", Type: "file", Template: true},
				{Path: ".claude", Type: "directory"},
				{Path: ".claude/agents", Type: "directory"},
			},
		},
		{
			Name:        "Go CLI + MCF",
			Type:        "cli-go",
			Version:     "1.0.0",
			Description: "Go CLI application with MCF development workflows",
			Author:      "MCF Team",
			Tags:        []string{"go", "cli", "cobra"},
			Parameters: []TemplateParameter{
				{
					Key:         "project_name",
					Label:       "Project Name",
					Type:        "text",
					Description: "Name of your Go CLI project",
					Required:    true,
					Placeholder: "my-go-cli",
				},
				{
					Key:         "module_name",
					Label:       "Go Module Name",
					Type:        "text",
					Description: "Go module path (e.g., github.com/user/project)",
					Required:    true,
					Placeholder: "github.com/user/my-go-cli",
				},
			},
			Files: []TemplateFile{
				{Path: "go.mod", Type: "file", Template: true},
				{Path: "main.go", Type: "file", Template: true},
				{Path: "cmd", Type: "directory"},
				{Path: "cmd/root.go", Type: "file", Template: true},
				{Path: ".claude", Type: "directory"},
				{Path: ".claude/agents", Type: "directory"},
			},
		},
		{
			Name:        "Python API + MCF",
			Type:        "api-python",
			Version:     "1.0.0",
			Description: "Python FastAPI project with MCF integration",
			Author:      "MCF Team",
			Tags:        []string{"python", "api", "fastapi"},
			Parameters: []TemplateParameter{
				{
					Key:         "project_name",
					Label:       "Project Name",
					Type:        "text",
					Description: "Name of your Python API project",
					Required:    true,
					Placeholder: "my-python-api",
				},
				{
					Key:         "python_version",
					Label:       "Python Version",
					Type:        "select",
					Description: "Target Python version",
					Options:     []string{"3.8", "3.9", "3.10", "3.11", "3.12"},
					Default:     "3.11",
				},
			},
			Files: []TemplateFile{
				{Path: "requirements.txt", Type: "file", Template: true},
				{Path: "main.py", Type: "file", Template: true},
				{Path: "app", Type: "directory"},
				{Path: "app/__init__.py", Type: "file"},
				{Path: "app/routers", Type: "directory"},
				{Path: ".claude", Type: "directory"},
			},
		},
	}
}
