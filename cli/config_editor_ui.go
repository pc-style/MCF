package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// EditorModel represents the state of the configuration editor
type EditorModel struct {
	state          EditorState
	sections       []EditorSection
	currentSection int
	currentField   int
	textInput      textinput.Model
	listModel      list.Model
	config         map[string]interface{}
	errors         map[string]string
	showPreview    bool
	showHelp       bool
	width          int
	height         int
	schema         EditorSchema
}

// EditorState represents the current state of the editor
type EditorState int

const (
	EditorStateNavigation EditorState = iota
	EditorStateFieldEdit
	EditorStatePreview
	EditorStateHelp
	EditorStateValidation
)

// EditorSchema represents the loaded schema structure
type EditorSchema struct {
	Sections []EditorSection `yaml:"sections"`
}

// EditorSection represents a configuration section
type EditorSection struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Optional     bool              `yaml:"optional"`
	Fields       []EditorField     `yaml:"fields"`
	Dependencies []string          `yaml:"dependencies"`
	Conditionals []EditorCondition `yaml:"conditionals"`
	Visible      bool              `yaml:"-"`
	Valid        bool              `yaml:"-"`
}

// EditorField represents a configuration field
type EditorField struct {
	Key          string            `yaml:"key"`
	Label        string            `yaml:"label"`
	Type         string            `yaml:"type"`
	Required     bool              `yaml:"required"`
	Description  string            `yaml:"description"`
	Default      interface{}       `yaml:"default"`
	Placeholder  string            `yaml:"placeholder"`
	HelpText     string            `yaml:"help_text"`
	Options      []EditorOption    `yaml:"options"`
	Validation   EditorValidation  `yaml:"validation"`
	Conditionals []EditorCondition `yaml:"conditionals"`
	Visible      bool              `yaml:"-"`
	Valid        bool              `yaml:"-"`
	Value        interface{}       `yaml:"-"`
	Error        string            `yaml:"-"`
}

// EditorOption represents an option for select fields
type EditorOption struct {
	Value       string `yaml:"value"`
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Default     bool   `yaml:"default"`
}

// EditorValidation represents field validation rules
type EditorValidation struct {
	Pattern   string   `yaml:"pattern"`
	MinLength int      `yaml:"min_length"`
	MaxLength int      `yaml:"max_length"`
	Required  bool     `yaml:"required"`
	Enum      []string `yaml:"enum"`
}

// EditorCondition represents conditional field logic
type EditorCondition struct {
	Field     string      `yaml:"field"`
	Operation string      `yaml:"operation"`
	Value     interface{} `yaml:"value"`
	Action    string      `yaml:"action"`
}

// Color definitions for consistent styling
var (
	editorPrimaryColor   = lipgloss.Color("99")
	editorSecondaryColor = lipgloss.Color("243")
	editorErrorColor     = lipgloss.Color("196")
	editorSuccessColor   = lipgloss.Color("34")
	editorWarningColor   = lipgloss.Color("220")
	editorHighlightColor = lipgloss.Color("212")
)

// Style definitions
var (
	editorTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(editorPrimaryColor).
				Align(lipgloss.Center)

	editorSectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(editorPrimaryColor)

	editorFieldLabelStyle = lipgloss.NewStyle().
				Foreground(editorPrimaryColor).
				Bold(true)

	editorFieldDescStyle = lipgloss.NewStyle().
				Foreground(editorSecondaryColor).
				Italic(true)

	editorHelpStyle = lipgloss.NewStyle().
			Foreground(editorSecondaryColor)

	editorErrorStyle = lipgloss.NewStyle().
				Foreground(editorErrorColor).
				Bold(true)

	editorSuccessStyle = lipgloss.NewStyle().
				Foreground(editorSuccessColor)

	editorWarningStyle = lipgloss.NewStyle().
				Foreground(editorWarningColor)

	editorSelectedStyle = lipgloss.NewStyle().
				Background(editorHighlightColor).
				Foreground(lipgloss.Color("0"))

	editorPreviewStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(editorPrimaryColor).
				Padding(1)
)

// NewEditorModel creates a new configuration editor model
func NewEditorModel(schemaPath string) (*EditorModel, error) {
	// Load schema from file
	schema, err := loadEditorSchema(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.CharLimit = 200

	// Initialize configuration map
	config := make(map[string]interface{})

	// Initialize sections with default values
	for i := range schema.Sections {
		schema.Sections[i].Visible = true
		schema.Sections[i].Valid = true

		for j := range schema.Sections[i].Fields {
			field := &schema.Sections[i].Fields[j]
			field.Visible = true
			field.Valid = true

			// Set default values
			if field.Default != nil {
				field.Value = field.Default
				config[field.Key] = field.Default
			}
		}
	}

	// Create list items for navigation
	listItems := make([]list.Item, len(schema.Sections))
	for i, section := range schema.Sections {
		listItems[i] = EditorSectionListItem{
			title:       section.Name,
			description: section.Description,
			index:       i,
		}
	}

	listModel := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Configuration Sections"
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)

	return &EditorModel{
		state:          EditorStateNavigation,
		sections:       schema.Sections,
		currentSection: 0,
		currentField:   0,
		textInput:      ti,
		listModel:      listModel,
		config:         config,
		errors:         make(map[string]string),
		schema:         schema,
		width:          80,
		height:         24,
	}, nil
}

// EditorSectionListItem implements list.Item for section navigation
type EditorSectionListItem struct {
	title       string
	description string
	index       int
}

func (s EditorSectionListItem) FilterValue() string { return s.title }
func (s EditorSectionListItem) Title() string       { return s.title }
func (s EditorSectionListItem) Description() string { return s.description }

// loadEditorSchema loads the configuration schema from a YAML file
func loadEditorSchema(path string) (EditorSchema, error) {
	var schema EditorSchema

	data, err := os.ReadFile(path)
	if err != nil {
		return schema, err
	}

	err = yaml.Unmarshal(data, &schema)
	return schema, err
}

// Init initializes the model
func (m EditorModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}

// Update handles messages and updates the model
func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listModel.SetSize(msg.Width-4, msg.Height-8)

	case tea.KeyMsg:
		switch m.state {
		case EditorStateNavigation:
			return m.handleNavigationInput(msg)
		case EditorStateFieldEdit:
			return m.handleFieldEditInput(msg)
		case EditorStatePreview:
			return m.handlePreviewInput(msg)
		case EditorStateHelp:
			return m.handleHelpInput(msg)
		case EditorStateValidation:
			return m.handleValidationInput(msg)
		}

	case EditorValidationErrorMsg:
		// Handle validation error - could show a toast or update UI
		return m, nil

	case EditorSaveErrorMsg:
		// Handle save error - could show error message
		return m, nil

	case EditorSaveSuccessMsg:
		// Handle successful save - could quit or show success message
		return m, tea.Quit
	}

	// Update components
	switch m.state {
	case EditorStateNavigation:
		m.listModel, cmd = m.listModel.Update(msg)
		cmds = append(cmds, cmd)
	case EditorStateFieldEdit:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m EditorModel) View() string {
	switch m.state {
	case EditorStateNavigation:
		return m.navigationView()
	case EditorStateFieldEdit:
		return m.fieldEditView()
	case EditorStatePreview:
		return m.previewView()
	case EditorStateHelp:
		return m.helpView()
	case EditorStateValidation:
		return m.validationView()
	default:
		return "Unknown state"
	}
}

// navigationView renders the section navigation view
func (m EditorModel) navigationView() string {
	title := editorTitleStyle.Render("MCF Configuration Editor")

	var content strings.Builder
	content.WriteString(title + "\n\n")

	// Show current section details
	if m.currentSection < len(m.sections) {
		section := m.sections[m.currentSection]
		content.WriteString(editorSectionTitleStyle.Render(section.Name) + "\n")
		content.WriteString(editorFieldDescStyle.Render(section.Description) + "\n")

		// Section status
		statusIcon := "✓"
		statusColor := editorSuccessColor
		if !section.Valid {
			statusIcon = "❌"
			statusColor = editorErrorColor
		} else if section.Optional {
			statusIcon = "◯"
			statusColor = editorWarningColor
		}

		status := lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon + " ")
		if section.Optional {
			status += "Optional section"
		} else if section.Valid {
			status += "Complete"
		} else {
			status += "Validation errors"
		}
		content.WriteString(status + "\n\n")

		// Field summary
		completedFields := 0
		totalFields := 0
		for _, field := range section.Fields {
			if field.Visible {
				totalFields++
				if field.Value != nil && field.Valid {
					completedFields++
				}
			}
		}

		progress := fmt.Sprintf("Fields: %d/%d configured", completedFields, totalFields)
		content.WriteString(editorFieldDescStyle.Render(progress) + "\n\n")
	}

	// Navigation list
	content.WriteString(m.listModel.View() + "\n")

	// Help text
	help := editorHelpStyle.Render("Enter: Configure section • p: Preview config • h: Help • q: Quit")
	content.WriteString("\n" + help)

	return content.String()
}

// fieldEditView renders the field editing view
func (m EditorModel) fieldEditView() string {
	if m.currentSection >= len(m.sections) || m.currentField >= len(m.sections[m.currentSection].Fields) {
		return "Invalid field selection"
	}

	section := m.sections[m.currentSection]
	field := section.Fields[m.currentField]

	title := editorTitleStyle.Render(fmt.Sprintf("%s - %s", section.Name, field.Label))

	var content strings.Builder
	content.WriteString(title + "\n\n")

	// Field information
	content.WriteString(editorFieldLabelStyle.Render(field.Label))
	if field.Required {
		content.WriteString(editorErrorStyle.Render(" *"))
	}
	content.WriteString("\n")

	content.WriteString(editorFieldDescStyle.Render(field.Description) + "\n")

	if field.HelpText != "" {
		content.WriteString(editorHelpStyle.Render("💡 "+field.HelpText) + "\n")
	}

	content.WriteString("\n")

	// Current value display
	currentValue := "Not set"
	if field.Value != nil {
		currentValue = fmt.Sprintf("%v", field.Value)
	}
	content.WriteString(fmt.Sprintf("Current value: %s\n", currentValue))

	// Field input based on type
	switch field.Type {
	case "text", "email":
		content.WriteString("Enter new value:\n")
		content.WriteString(m.textInput.View() + "\n")

	case "boolean":
		content.WriteString("Value: ")
		if val, ok := field.Value.(bool); ok && val {
			content.WriteString(editorSuccessStyle.Render("true"))
		} else {
			content.WriteString(editorErrorStyle.Render("false"))
		}
		content.WriteString(" (Space to toggle)\n")

	case "select":
		content.WriteString("Available options:\n")
		for _, option := range field.Options {
			prefix := "  "
			style := editorHelpStyle
			if field.Value == option.Value {
				prefix = "→ "
				style = editorSelectedStyle
			}

			line := fmt.Sprintf("%s%s - %s", prefix, option.Label, option.Description)
			content.WriteString(style.Render(line) + "\n")
		}
		content.WriteString("Use ↑/↓ to navigate, Enter to select\n")

	case "multi_select":
		content.WriteString("Available options (Space to toggle):\n")
		selectedValues := make(map[string]bool)
		if vals, ok := field.Value.([]interface{}); ok {
			for _, v := range vals {
				selectedValues[fmt.Sprintf("%v", v)] = true
			}
		}

		for _, option := range field.Options {
			prefix := "☐ "
			style := editorHelpStyle
			if selectedValues[option.Value] {
				prefix = "☑ "
				style = editorSuccessStyle
			}

			line := fmt.Sprintf("%s%s - %s", prefix, option.Label, option.Description)
			content.WriteString(style.Render(line) + "\n")
		}
	}

	// Validation error
	if field.Error != "" {
		content.WriteString("\n")
		content.WriteString(editorErrorStyle.Render("❌ "+field.Error) + "\n")
	}

	// Field navigation
	fieldNav := fmt.Sprintf("Field %d of %d", m.currentField+1, len(section.Fields))
	content.WriteString("\n" + editorHelpStyle.Render(fieldNav) + "\n")

	// Help text
	help := "Enter: Save • Tab: Next field • Shift+Tab: Previous field • Esc: Back to sections"
	content.WriteString(editorHelpStyle.Render(help))

	return content.String()
}

// previewView renders the configuration preview
func (m EditorModel) previewView() string {
	title := editorTitleStyle.Render("Configuration Preview")

	var content strings.Builder
	content.WriteString(title + "\n\n")

	// Generate YAML preview
	yamlData, err := yaml.Marshal(m.config)
	if err != nil {
		content.WriteString(editorErrorStyle.Render("Error generating preview: " + err.Error()))
	} else {
		preview := editorPreviewStyle.Render(string(yamlData))
		content.WriteString(preview)
	}

	// Validation summary
	content.WriteString("\n\n")
	content.WriteString(editorSectionTitleStyle.Render("Validation Summary") + "\n")

	validSections := 0
	totalSections := len(m.sections)

	for _, section := range m.sections {
		statusIcon := "✓"
		statusColor := editorSuccessColor
		status := "Valid"

		if !section.Valid {
			statusIcon = "❌"
			statusColor = editorErrorColor
			status = "Has errors"
		} else {
			validSections++
		}

		line := fmt.Sprintf("%s %s: %s",
			lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon),
			section.Name,
			status,
		)
		content.WriteString(line + "\n")
	}

	summaryColor := editorSuccessColor
	if validSections < totalSections {
		summaryColor = editorErrorColor
	}

	summary := fmt.Sprintf("\n%d/%d sections valid", validSections, totalSections)
	content.WriteString(lipgloss.NewStyle().Foreground(summaryColor).Bold(true).Render(summary) + "\n")

	// Help text
	help := "s: Save config • v: Validate all • Esc: Back to sections"
	content.WriteString("\n" + editorHelpStyle.Render(help))

	return content.String()
}

// helpView renders the help view
func (m EditorModel) helpView() string {
	title := editorTitleStyle.Render("MCF Configuration Editor - Help")

	var content strings.Builder
	content.WriteString(title + "\n\n")

	helpSections := []struct {
		title string
		items []string
	}{
		{
			title: "Navigation",
			items: []string{
				"↑/↓ - Navigate sections or options",
				"Enter - Select section or confirm input",
				"Tab - Next field",
				"Shift+Tab - Previous field",
				"Esc - Go back / Cancel",
			},
		},
		{
			title: "Field Editing",
			items: []string{
				"Text fields - Type to enter value",
				"Boolean fields - Space to toggle",
				"Select fields - ↑/↓ to choose, Enter to select",
				"Multi-select - Space to toggle options",
			},
		},
		{
			title: "Global Commands",
			items: []string{
				"p - Preview current configuration",
				"h - Show this help",
				"v - Validate all fields",
				"s - Save configuration",
				"q - Quit (with confirmation if unsaved)",
			},
		},
		{
			title: "Validation",
			items: []string{
				"* Required fields must be filled",
				"Validation happens on field save",
				"Red indicators show errors",
				"Yellow indicators show warnings",
			},
		},
	}

	for _, section := range helpSections {
		content.WriteString(editorSectionTitleStyle.Render(section.title) + "\n")
		for _, item := range section.items {
			content.WriteString("  • " + item + "\n")
		}
		content.WriteString("\n")
	}

	help := "Press any key to return"
	content.WriteString(editorHelpStyle.Render(help))

	return content.String()
}

// validationView renders the validation results view
func (m EditorModel) validationView() string {
	title := editorTitleStyle.Render("Configuration Validation")

	var content strings.Builder
	content.WriteString(title + "\n\n")

	hasErrors := false

	for _, section := range m.sections {
		content.WriteString(editorSectionTitleStyle.Render(section.Name) + "\n")

		sectionHasErrors := false
		for _, field := range section.Fields {
			if !field.Visible {
				continue
			}

			statusIcon := "✓"
			statusColor := editorSuccessColor
			message := "Valid"

			if field.Error != "" {
				statusIcon = "❌"
				statusColor = editorErrorColor
				message = field.Error
				sectionHasErrors = true
				hasErrors = true
			} else if field.Required && field.Value == nil {
				statusIcon = "⚠"
				statusColor = editorWarningColor
				message = "Required field not set"
				sectionHasErrors = true
			}

			line := fmt.Sprintf("  %s %s: %s",
				lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon),
				field.Label,
				message,
			)
			content.WriteString(line + "\n")
		}

		if !sectionHasErrors {
			content.WriteString(editorSuccessStyle.Render("  All fields valid") + "\n")
		}
		content.WriteString("\n")
	}

	// Overall status
	if hasErrors {
		content.WriteString(editorErrorStyle.Render("❌ Configuration has validation errors") + "\n")
	} else {
		content.WriteString(editorSuccessStyle.Render("✅ Configuration is valid") + "\n")
	}

	help := "f: Fix first error • Esc: Back to sections"
	content.WriteString("\n" + editorHelpStyle.Render(help))

	return content.String()
}

// Input handlers
func (m EditorModel) handleNavigationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "enter":
		if selected, ok := m.listModel.SelectedItem().(EditorSectionListItem); ok {
			m.currentSection = selected.index
			m.currentField = 0
			m.state = EditorStateFieldEdit

			// Initialize text input for first field
			if len(m.sections[m.currentSection].Fields) > 0 {
				field := m.sections[m.currentSection].Fields[0]
				m.updateTextInputForField(field)
			}
		}

	case "p":
		m.state = EditorStatePreview

	case "h":
		m.state = EditorStateHelp

	case "v":
		m.validateAllFields()
		m.state = EditorStateValidation
	}

	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(msg)
	return m, cmd
}

func (m EditorModel) handleFieldEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.currentSection >= len(m.sections) || m.currentField >= len(m.sections[m.currentSection].Fields) {
		return m, nil
	}

	field := &m.sections[m.currentSection].Fields[m.currentField]

	switch msg.String() {
	case "esc":
		m.state = EditorStateNavigation
		return m, nil

	case "tab":
		m.nextField()

	case "shift+tab":
		m.prevField()

	case "enter":
		switch field.Type {
		case "text", "email":
			value := strings.TrimSpace(m.textInput.Value())
			if value != "" || !field.Required {
				field.Value = value
				m.config[field.Key] = value
				field.Error = m.validateField(*field)
				if field.Error == "" {
					field.Valid = true
					m.nextField()
				} else {
					field.Valid = false
				}
			}

		case "select":
			// Handle select field confirmation
			if field.Value != nil {
				field.Error = m.validateField(*field)
				if field.Error == "" {
					field.Valid = true
					m.nextField()
				} else {
					field.Valid = false
				}
			}
		}

	case " ":
		switch field.Type {
		case "boolean":
			currentVal := false
			if val, ok := field.Value.(bool); ok {
				currentVal = val
			}
			field.Value = !currentVal
			m.config[field.Key] = !currentVal
			field.Error = m.validateField(*field)
			field.Valid = field.Error == ""

		case "multi_select":
			// Handle multi-select toggle
			m.toggleMultiSelectOption(field)
		}

	case "up", "down":
		if field.Type == "select" {
			m.navigateSelectOptions(field, msg.String() == "down")
		} else if field.Type == "multi_select" {
			// Navigate through multi-select options
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m EditorModel) handlePreviewInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = EditorStateNavigation
	case "s":
		return m, m.saveConfiguration()
	case "v":
		m.validateAllFields()
		m.state = EditorStateValidation
	}
	return m, nil
}

func (m EditorModel) handleHelpInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.state = EditorStateNavigation
	return m, nil
}

func (m EditorModel) handleValidationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = EditorStateNavigation
	case "f":
		// Jump to first error
		m.jumpToFirstError()
		m.state = EditorStateFieldEdit
	}
	return m, nil
}

// Helper methods
func (m *EditorModel) nextField() {
	section := &m.sections[m.currentSection]
	for {
		m.currentField++
		if m.currentField >= len(section.Fields) {
			// Move to next section
			for i := m.currentSection + 1; i < len(m.sections); i++ {
				if m.sections[i].Visible && len(m.sections[i].Fields) > 0 {
					m.currentSection = i
					m.currentField = 0
					break
				}
			}
			// If no next section found, stay at current position
			if m.currentSection >= len(m.sections) || m.currentField >= len(m.sections[m.currentSection].Fields) {
				m.currentField = len(section.Fields) - 1
				return
			}
		}

		if section.Fields[m.currentField].Visible {
			break
		}
	}

	m.updateTextInputForField(m.sections[m.currentSection].Fields[m.currentField])
}

func (m *EditorModel) prevField() {
	for {
		m.currentField--
		if m.currentField < 0 {
			// Move to previous section
			for i := m.currentSection - 1; i >= 0; i-- {
				if m.sections[i].Visible && len(m.sections[i].Fields) > 0 {
					m.currentSection = i
					m.currentField = len(m.sections[i].Fields) - 1
					break
				}
			}
			// If no previous section found, stay at current position
			if m.currentSection < 0 || m.currentField < 0 {
				m.currentSection = 0
				m.currentField = 0
				return
			}
		}

		section := m.sections[m.currentSection]
		if m.currentField < len(section.Fields) && section.Fields[m.currentField].Visible {
			break
		}
	}

	m.updateTextInputForField(m.sections[m.currentSection].Fields[m.currentField])
}

func (m *EditorModel) updateTextInputForField(field EditorField) {
	m.textInput.Placeholder = field.Placeholder
	if field.Type == "email" {
		m.textInput.Placeholder = "user@example.com"
	}

	// Set current value
	if field.Value != nil {
		m.textInput.SetValue(fmt.Sprintf("%v", field.Value))
	} else {
		m.textInput.SetValue("")
	}

	// Set character limit based on validation
	if field.Validation.MaxLength > 0 {
		m.textInput.CharLimit = field.Validation.MaxLength
	} else {
		m.textInput.CharLimit = 200
	}

	m.textInput.Focus()
}

func (m *EditorModel) navigateSelectOptions(field *EditorField, down bool) {
	if len(field.Options) == 0 {
		return
	}

	currentIndex := -1
	for i, option := range field.Options {
		if field.Value == option.Value {
			currentIndex = i
			break
		}
	}

	if down {
		currentIndex++
		if currentIndex >= len(field.Options) {
			currentIndex = 0
		}
	} else {
		currentIndex--
		if currentIndex < 0 {
			currentIndex = len(field.Options) - 1
		}
	}

	field.Value = field.Options[currentIndex].Value
	m.config[field.Key] = field.Value
}

func (m *EditorModel) toggleMultiSelectOption(field *EditorField) {
	if len(field.Options) == 0 {
		return
	}

	// Get current selected values
	selectedValues := make(map[string]bool)
	if vals, ok := field.Value.([]interface{}); ok {
		for _, v := range vals {
			selectedValues[fmt.Sprintf("%v", v)] = true
		}
	} else if field.Value == nil {
		field.Value = []interface{}{}
	}

	// For simplicity, toggle the first option for demonstration
	// In a full implementation, you'd need to track which option is currently highlighted
	if len(field.Options) > 0 {
		firstOption := field.Options[0].Value
		if selectedValues[firstOption] {
			// Remove from selection
			delete(selectedValues, firstOption)
		} else {
			// Add to selection
			selectedValues[firstOption] = true
		}

		// Convert back to slice
		var newValues []interface{}
		for value := range selectedValues {
			newValues = append(newValues, value)
		}
		field.Value = newValues
		m.config[field.Key] = newValues
	}
}

func (m *EditorModel) validateField(field EditorField) string {
	if field.Value == nil {
		if field.Required {
			return "This field is required"
		}
		return ""
	}

	valueStr := fmt.Sprintf("%v", field.Value)

	// Required validation
	if field.Required && valueStr == "" {
		return "This field is required"
	}

	// Length validation
	if field.Validation.MinLength > 0 && len(valueStr) < field.Validation.MinLength {
		return fmt.Sprintf("Minimum length is %d characters", field.Validation.MinLength)
	}

	if field.Validation.MaxLength > 0 && len(valueStr) > field.Validation.MaxLength {
		return fmt.Sprintf("Maximum length is %d characters", field.Validation.MaxLength)
	}

	// Pattern validation
	if field.Validation.Pattern != "" {
		matched, err := regexp.MatchString(field.Validation.Pattern, valueStr)
		if err != nil {
			return "Invalid validation pattern"
		}
		if !matched {
			return "Value does not match required format"
		}
	}

	// Enum validation
	if len(field.Validation.Enum) > 0 {
		valid := false
		for _, enumVal := range field.Validation.Enum {
			if valueStr == enumVal {
				valid = true
				break
			}
		}
		if !valid {
			return "Value must be one of the allowed options"
		}
	}

	// Email validation
	if field.Type == "email" {
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailPattern, valueStr)
		if !matched {
			return "Please enter a valid email address"
		}
	}

	return ""
}

func (m *EditorModel) validateAllFields() {
	for i := range m.sections {
		section := &m.sections[i]
		section.Valid = true

		for j := range section.Fields {
			field := &section.Fields[j]
			if !field.Visible {
				continue
			}

			field.Error = m.validateField(*field)
			field.Valid = field.Error == ""

			if !field.Valid {
				section.Valid = false
			}
		}
	}
}

func (m *EditorModel) jumpToFirstError() {
	for i, section := range m.sections {
		for j, field := range section.Fields {
			if !field.Valid && field.Visible {
				m.currentSection = i
				m.currentField = j
				m.updateTextInputForField(field)
				return
			}
		}
	}
}

func (m EditorModel) saveConfiguration() tea.Cmd {
	return func() tea.Msg {
		// Validate all fields first
		m.validateAllFields()

		// Check if there are any validation errors
		hasErrors := false
		for _, section := range m.sections {
			if !section.Valid {
				hasErrors = true
				break
			}
		}

		if hasErrors {
			return EditorValidationErrorMsg{Error: "Configuration has validation errors"}
		}

		// Generate final configuration
		finalConfig := make(map[string]interface{})
		finalConfig["generated_at"] = time.Now().Format(time.RFC3339)
		finalConfig["schema_version"] = "1.0"

		for key, value := range m.config {
			if value != nil {
				finalConfig[key] = value
			}
		}

		// Save to file
		configPath := "mcf-config.yaml"
		yamlData, err := yaml.Marshal(finalConfig)
		if err != nil {
			return EditorSaveErrorMsg{Error: fmt.Sprintf("Failed to marshal configuration: %v", err)}
		}

		err = os.WriteFile(configPath, yamlData, 0644)
		if err != nil {
			return EditorSaveErrorMsg{Error: fmt.Sprintf("Failed to save configuration: %v", err)}
		}

		return EditorSaveSuccessMsg{Path: configPath}
	}
}

// Message types
type EditorValidationErrorMsg struct {
	Error string
}

type EditorSaveErrorMsg struct {
	Error string
}

type EditorSaveSuccessMsg struct {
	Path string
}

// RunConfigurationEditor starts the configuration editor with the given schema
func RunConfigurationEditor(schemaPath string) error {
	model, err := NewEditorModel(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to initialize config editor: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// Example usage function for demonstration
func ExampleConfigurationEditorUsage() {
	// This would be called from main.go or another entry point
	schemaPath := "./config-schema.yaml"

	if err := RunConfigurationEditor(schemaPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error running config editor: %v\n", err)
		os.Exit(1)
	}
}
