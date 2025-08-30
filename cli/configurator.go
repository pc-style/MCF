package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Configuration Editor UI Component
type ConfiguratorModel struct {
	state          ConfiguratorState
	sections       []ConfigSection
	currentSection int
	currentField   int
	textInputs     []textinput.Model
	listModel      list.Model
	help           help.Model
	configManager  *ConfigManager
	schema         *ConfigSchema
	fieldValues    map[string]interface{}
	errors         map[string]string
	unsavedChanges bool
	previewMode    bool
	searchQuery    string
	searchInput    textinput.Model
	width          int
	height         int
	keyMap         ConfiguratorKeyMap
}

type ConfiguratorState int

const (
	ConfigStateSectionList ConfiguratorState = iota
	ConfigStateFieldEdit
	ConfigStatePreview
	ConfigStateSearch
	ConfigStateHelp
	ConfigStateConfirmSave
	ConfigStateError
)

type ConfigSection struct {
	Name        string
	Description string
	Fields      []ConfigurationField
	Active      bool
	Errors      map[string]string
}

type ConfigurationField struct {
	Key         string
	Label       string
	Type        FieldType
	Value       interface{}
	Default     interface{}
	Required    bool
	Description string
	Placeholder string
	HelpText    string
	Options     []FieldOption
	Validator   *FieldValidator
	Sensitive   bool
	ReadOnly    bool
}

type ConfiguratorKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Escape   key.Binding
	Save     key.Binding
	Preview  key.Binding
	Search   key.Binding
	Help     key.Binding
	Reset    key.Binding
	Quit     key.Binding
}

func NewConfiguratorKeyMap() ConfiguratorKeyMap {
	return ConfiguratorKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/edit"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev field"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel/back"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Preview: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "preview"),
		),
		Search: key.NewBinding(
			key.WithKeys("ctrl+f", "/"),
			key.WithHelp("ctrl+f", "search"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Reset: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "reset"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k ConfiguratorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Save, k.Preview, k.Search, k.Quit}
}

func (k ConfiguratorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Tab, k.ShiftTab, k.Escape},
		{k.Save, k.Preview, k.Search, k.Help},
		{k.Reset, k.Quit},
	}
}

func NewConfiguratorModel(configManager *ConfigManager) ConfiguratorModel {
	keyMap := NewConfiguratorKeyMap()

	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search configuration fields..."
	searchInput.CharLimit = 50

	// Load configuration schema
	schema := createDefaultConfigSchema()

	// Create sections from schema
	sections := createConfigSections(schema, configManager)

	// Create list model for sections
	items := make([]list.Item, len(sections))
	for i, section := range sections {
		items[i] = ConfigSectionItem{
			name:        section.Name,
			description: section.Description,
			fieldCount:  len(section.Fields),
		}
	}

	delegate := ConfigSectionDelegate{}
	listModel := list.New(items, delegate, 0, 0)
	listModel.Title = "Configuration Sections"
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)

	model := ConfiguratorModel{
		state:          ConfigStateSectionList,
		sections:       sections,
		currentSection: 0,
		currentField:   0,
		listModel:      listModel,
		help:           help.New(),
		configManager:  configManager,
		schema:         schema,
		fieldValues:    make(map[string]interface{}),
		errors:         make(map[string]string),
		searchInput:    searchInput,
		keyMap:         keyMap,
	}

	// Register this model with the global state for cross-component coordination
	GlobalState.SetConfiguratorModel(&model)

	return model
}

func (m ConfiguratorModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.loadConfigurationValues(),
	)
}

func (m ConfiguratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listModel.SetWidth(msg.Width)
		m.listModel.SetHeight(msg.Height - 6) // Leave space for header/footer
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch m.state {
		case ConfigStateSectionList:
			return m.handleSectionListInput(msg)
		case ConfigStateFieldEdit:
			return m.handleFieldEditInput(msg)
		case ConfigStatePreview:
			return m.handlePreviewInput(msg)
		case ConfigStateSearch:
			return m.handleSearchInput(msg)
		case ConfigStateHelp:
			return m.handleHelpInput(msg)
		case ConfigStateConfirmSave:
			return m.handleConfirmSaveInput(msg)
		case ConfigStateError:
			return m.handleErrorInput(msg)
		}

	case ConfigLoadedMsg:
		m.fieldValues = msg.Values
		m.refreshSectionsFromValues()

	case ConfigSavedMsg:
		m.unsavedChanges = false
		if msg.Success {
			// Show success message briefly, then return to section list
			return m, tea.Sequence(
				tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return ConfigSaveCompleteMsg{}
				}),
			)
		} else {
			m.state = ConfigStateError
			m.errors["save"] = msg.Error
		}

	case ConfigSaveCompleteMsg:
		m.state = ConfigStateSectionList

	case ValidationErrorMsg:
		m.errors[msg.Field] = msg.Error
	}

	// Update components based on current state
	switch m.state {
	case ConfigStateSectionList:
		m.listModel, cmd = m.listModel.Update(msg)
		cmds = append(cmds, cmd)

	case ConfigStateFieldEdit:
		if m.currentField < len(m.textInputs) {
			m.textInputs[m.currentField], cmd = m.textInputs[m.currentField].Update(msg)
			cmds = append(cmds, cmd)
		}

	case ConfigStateSearch:
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleMCFMessage processes global MCF messages for cross-component coordination
func (m ConfiguratorModel) handleMCFMessage(msg MCFMessage) tea.Cmd {
	switch msg.Type() {
	case MsgModeTransition:
		if transMsg, ok := msg.(ModeTransitionMessage); ok {
			// React to mode transitions if needed
			if transMsg.FromMode == ModeConfigurator {
				// Save any unsaved changes when leaving configurator
				if m.unsavedChanges {
					// Could auto-save or warn user
					return m.saveConfiguration()
				}
			}
		}

	case MsgConfigLoaded:
		if configMsg, ok := msg.(ConfigurationMessage); ok {
			// Another component loaded configuration, refresh our view
			if configMsg.Success && configMsg.Data != nil {
				return m.loadConfigurationValues()
			}
		}

	case MsgUIError:
		if uiMsg, ok := msg.(UIMessage); ok {
			// Handle UI errors that might affect configuration
			if uiMsg.Context != nil {
				if errorContext, ok := uiMsg.Context.(map[string]interface{}); ok {
					if section, exists := errorContext["config_section"]; exists {
						// Set error state for specific section
						m.errors[section.(string)] = uiMsg.Message
						m.state = ConfigStateError
					}
				}
			}
		}

	case MsgAppShutdown:
		// Handle graceful shutdown
		if m.unsavedChanges {
			// Could trigger a save confirmation dialog
			return func() tea.Msg {
				return UIMessage{
					BaseMessage: BaseMessage{MsgType: MsgUIWarning, CreatedAt: time.Now()},
					Level:       "warning",
					Title:       "Unsaved Changes",
					Message:     "Configuration has unsaved changes",
				}
			}
		}
	}
	return nil
}

func (m ConfiguratorModel) View() string {
	switch m.state {
	case ConfigStateSectionList:
		return m.sectionListView()
	case ConfigStateFieldEdit:
		return m.fieldEditView()
	case ConfigStatePreview:
		return m.previewView()
	case ConfigStateSearch:
		return m.searchView()
	case ConfigStateHelp:
		return m.helpView()
	case ConfigStateConfirmSave:
		return m.confirmSaveView()
	case ConfigStateError:
		return m.errorView()
	default:
		return "Unknown state"
	}
}

func (m ConfiguratorModel) sectionListView() string {
	var sections []string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Align(lipgloss.Center)

	header := headerStyle.Render("⚙️ MCF Configuration")

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Align(lipgloss.Center)

	var statusMsg string
	if m.unsavedChanges {
		statusMsg = "⚠️ Unsaved changes"
	} else {
		statusMsg = "✅ All changes saved"
	}
	status := statusStyle.Render(statusMsg)

	// Main content - list of sections
	listContent := m.listModel.View()

	// Help footer
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpContent := helpStyle.Render(m.help.View(m.keyMap))

	sections = append(sections, header, "", status, "", listContent, "", helpContent)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m ConfiguratorModel) fieldEditView() string {
	if m.currentSection >= len(m.sections) {
		return "Invalid section"
	}

	section := m.sections[m.currentSection]

	// Section header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	header := headerStyle.Render(fmt.Sprintf("⚙️ %s", section.Name))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true)

	description := descStyle.Render(section.Description)

	// Progress indicator
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220"))

	progress := progressStyle.Render(fmt.Sprintf("Field %d of %d", m.currentField+1, len(section.Fields)))

	var content []string
	content = append(content, header, description, "", progress, "")

	// Current field
	if m.currentField < len(section.Fields) {
		field := section.Fields[m.currentField]

		// Field label
		labelStyle := lipgloss.NewStyle().Bold(true)
		if field.Required {
			labelStyle = labelStyle.Foreground(lipgloss.Color("196"))
		} else {
			labelStyle = labelStyle.Foreground(lipgloss.Color("34"))
		}

		label := labelStyle.Render(field.Label)
		if field.Required {
			label += " *"
		}

		content = append(content, label)

		// Field description
		if field.Description != "" {
			fieldDescStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Italic(true)
			content = append(content, fieldDescStyle.Render(field.Description))
		}

		// Field input
		content = append(content, "")
		content = append(content, m.renderFieldInput(field))

		// Field error
		if err, exists := m.errors[field.Key]; exists {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)
			content = append(content, "", errorStyle.Render(fmt.Sprintf("❌ %s", err)))
		}

		// Field help text
		if field.HelpText != "" {
			helpStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true)
			content = append(content, "", helpStyle.Render(fmt.Sprintf("💡 %s", field.HelpText)))
		}
	}

	// Navigation help
	content = append(content, "", "")
	navStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	content = append(content, navStyle.Render("Tab: Next field • Shift+Tab: Previous field • Enter: Edit • Esc: Back to sections"))

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m ConfiguratorModel) previewView() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Align(lipgloss.Center)

	header := headerStyle.Render("📋 Configuration Preview")

	var sections []string
	sections = append(sections, header, "")

	// Generate preview for each section
	for _, section := range m.sections {
		if len(section.Fields) == 0 {
			continue
		}

		sectionStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("220")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			PaddingTop(1)

		sections = append(sections, sectionStyle.Render(section.Name))

		for _, field := range section.Fields {
			value := m.getFieldValue(field.Key)
			displayValue := m.formatValueForDisplay(field, value)

			if field.Sensitive && displayValue != "" {
				displayValue = strings.Repeat("*", len(displayValue))
			}

			fieldStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				PaddingLeft(2)

			if displayValue == "" {
				displayValue = lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Italic(true).
					Render("<not set>")
			}

			fieldLine := fmt.Sprintf("%s: %s", field.Label, displayValue)
			sections = append(sections, fieldStyle.Render(fieldLine))
		}
		sections = append(sections, "")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	footer := footerStyle.Render("Ctrl+S: Save • Esc: Back to editing")
	sections = append(sections, footer)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m ConfiguratorModel) searchView() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	header := headerStyle.Render("🔍 Search Configuration")

	searchLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("Search query:")

	searchBox := m.searchInput.View()

	var results []string
	if m.searchQuery != "" {
		results = m.searchFields(m.searchQuery)
	}

	var content []string
	content = append(content, header, "", searchLabel, searchBox, "")

	if len(results) > 0 {
		resultStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Bold(true)
		content = append(content, resultStyle.Render("Search Results:"), "")

		for _, result := range results {
			content = append(content, lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				PaddingLeft(2).
				Render(result))
		}
	} else if m.searchQuery != "" {
		noResultsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Italic(true)
		content = append(content, noResultsStyle.Render("No matching fields found"))
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	content = append(content, "", helpStyle.Render("Enter: Search • Esc: Cancel"))

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m ConfiguratorModel) helpView() string {
	return m.help.View(m.keyMap)
}

func (m ConfiguratorModel) confirmSaveView() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("220")).
		Align(lipgloss.Center)

	header := headerStyle.Render("💾 Save Configuration")

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Bold(true)

	warning := warningStyle.Render("Are you sure you want to save these changes?")

	// Show summary of changes
	var changes []string
	for key, value := range m.fieldValues {
		changes = append(changes, fmt.Sprintf("• %s: %v", key, value))
	}

	var content []string
	content = append(content, header, "", warning, "")

	if len(changes) > 0 {
		changesStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
		content = append(content, changesStyle.Render("Changes to save:"))
		for _, change := range changes {
			content = append(content, lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				PaddingLeft(2).
				Render(change))
		}
		content = append(content, "")
	}

	instructionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	content = append(content, instructionStyle.Render("Y: Confirm save • N: Cancel • Esc: Back"))

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m ConfiguratorModel) errorView() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")).
		Align(lipgloss.Center)

	header := headerStyle.Render("❌ Configuration Error")

	var content []string
	content = append(content, header, "")

	for field, err := range m.errors {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
		content = append(content, errorStyle.Render(fmt.Sprintf("%s: %s", field, err)))
	}

	instructionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	content = append(content, "", instructionStyle.Render("Any key: Continue"))

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

// Input handlers
func (m ConfiguratorModel) handleSectionListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Enter):
		if selected, ok := m.listModel.SelectedItem().(ConfigSectionItem); ok {
			// Find section index
			for i, section := range m.sections {
				if section.Name == selected.name {
					m.currentSection = i
					m.currentField = 0
					m.state = ConfigStateFieldEdit
					m.initializeTextInputs()
					break
				}
			}
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Preview):
		m.state = ConfigStatePreview
		return m, nil

	case key.Matches(msg, m.keyMap.Search):
		m.state = ConfigStateSearch
		m.searchInput.Focus()
		return m, nil

	case key.Matches(msg, m.keyMap.Save):
		if m.unsavedChanges {
			m.state = ConfigStateConfirmSave
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Help):
		m.state = ConfigStateHelp
		return m, nil

	case key.Matches(msg, m.keyMap.Quit):
		if m.unsavedChanges {
			// TODO: Show unsaved changes warning
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(msg)
	return m, cmd
}

func (m ConfiguratorModel) handleFieldEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	section := m.sections[m.currentSection]

	switch {
	case key.Matches(msg, m.keyMap.Tab):
		m.nextField()
		return m, nil

	case key.Matches(msg, m.keyMap.ShiftTab):
		m.previousField()
		return m, nil

	case key.Matches(msg, m.keyMap.Enter):
		if m.currentField < len(section.Fields) {
			field := section.Fields[m.currentField]
			return m.handleFieldValueChange(field)
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Escape):
		m.state = ConfigStateSectionList
		return m, nil

	case key.Matches(msg, m.keyMap.Save):
		m.state = ConfigStateConfirmSave
		return m, nil

	case key.Matches(msg, m.keyMap.Preview):
		m.state = ConfigStatePreview
		return m, nil
	}

	return m, nil
}

func (m ConfiguratorModel) handlePreviewInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Escape):
		m.state = ConfigStateFieldEdit
		return m, nil

	case key.Matches(msg, m.keyMap.Save):
		m.state = ConfigStateConfirmSave
		return m, nil
	}

	return m, nil
}

func (m ConfiguratorModel) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Enter):
		m.searchQuery = m.searchInput.Value()
		return m, nil

	case key.Matches(msg, m.keyMap.Escape):
		m.state = ConfigStateSectionList
		m.searchInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m ConfiguratorModel) handleHelpInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.state = ConfigStateSectionList
	return m, nil
}

func (m ConfiguratorModel) handleConfirmSaveInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m, m.saveConfiguration()
	case "n", "N", "esc":
		m.state = ConfigStateFieldEdit
		return m, nil
	}
	return m, nil
}

func (m ConfiguratorModel) handleErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.state = ConfigStateFieldEdit
	return m, nil
}

// Helper functions
func (m *ConfiguratorModel) nextField() {
	section := m.sections[m.currentSection]
	if m.currentField < len(section.Fields)-1 {
		m.currentField++
	} else {
		// Move to next section
		if m.currentSection < len(m.sections)-1 {
			m.currentSection++
			m.currentField = 0
			m.initializeTextInputs()
		}
	}
}

func (m *ConfiguratorModel) previousField() {
	if m.currentField > 0 {
		m.currentField--
	} else {
		// Move to previous section
		if m.currentSection > 0 {
			m.currentSection--
			section := m.sections[m.currentSection]
			m.currentField = len(section.Fields) - 1
			m.initializeTextInputs()
		}
	}
}

func (m *ConfiguratorModel) initializeTextInputs() {
	section := m.sections[m.currentSection]
	m.textInputs = make([]textinput.Model, len(section.Fields))

	for i, field := range section.Fields {
		input := textinput.New()
		input.Placeholder = field.Placeholder
		input.CharLimit = 200

		if field.Type == FieldTypePassword {
			input.EchoMode = textinput.EchoPassword
		}

		// Set current value
		if value := m.getFieldValue(field.Key); value != nil {
			input.SetValue(fmt.Sprintf("%v", value))
		} else if field.Default != nil {
			input.SetValue(fmt.Sprintf("%v", field.Default))
		}

		if i == m.currentField {
			input.Focus()
		}

		m.textInputs[i] = input
	}
}

func (m ConfiguratorModel) renderFieldInput(field ConfigurationField) string {
	switch field.Type {
	case FieldTypeSelect:
		return m.renderSelectField(field)
	case FieldTypeMulti:
		return m.renderMultiSelectField(field)
	case FieldTypeBool:
		return m.renderBooleanField(field)
	default:
		return m.renderTextField(field)
	}
}

func (m ConfiguratorModel) renderTextField(field ConfigurationField) string {
	if m.currentField < len(m.textInputs) {
		return m.textInputs[m.currentField].View()
	}

	value := m.getFieldValue(field.Key)
	if value == nil {
		value = field.Default
	}

	displayValue := fmt.Sprintf("%v", value)
	if field.Sensitive {
		displayValue = strings.Repeat("*", len(displayValue))
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		Width(40)

	return style.Render(displayValue)
}

func (m ConfiguratorModel) renderSelectField(field ConfigurationField) string {
	value := m.getFieldValue(field.Key)

	var options []string
	for _, option := range field.Options {
		prefix := "  "
		if fmt.Sprintf("%v", value) == option.Value {
			prefix = "▶ "
		}
		options = append(options, fmt.Sprintf("%s%s", prefix, option.Label))
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1).
		Border(lipgloss.RoundedBorder())

	return style.Render(strings.Join(options, "\n"))
}

func (m ConfiguratorModel) renderMultiSelectField(field ConfigurationField) string {
	value := m.getFieldValue(field.Key)
	selectedValues := make(map[string]bool)

	if values, ok := value.([]string); ok {
		for _, v := range values {
			selectedValues[v] = true
		}
	}

	var options []string
	for _, option := range field.Options {
		prefix := "☐ "
		if selectedValues[option.Value] {
			prefix = "☑ "
		}
		options = append(options, fmt.Sprintf("%s%s", prefix, option.Label))
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1).
		Border(lipgloss.RoundedBorder())

	return style.Render(strings.Join(options, "\n"))
}

func (m ConfiguratorModel) renderBooleanField(field ConfigurationField) string {
	value := m.getFieldValue(field.Key)

	var boolValue bool
	if value != nil {
		if b, ok := value.(bool); ok {
			boolValue = b
		} else if str, ok := value.(string); ok {
			boolValue, _ = strconv.ParseBool(str)
		}
	}

	var options []string
	if boolValue {
		options = []string{"▶ Yes", "  No"}
	} else {
		options = []string{"  Yes", "▶ No"}
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1).
		Border(lipgloss.RoundedBorder())

	return style.Render(strings.Join(options, "\n"))
}

func (m ConfiguratorModel) getFieldValue(key string) interface{} {
	if value, exists := m.fieldValues[key]; exists {
		return value
	}
	return nil
}

func (m ConfiguratorModel) formatValueForDisplay(field ConfigurationField, value interface{}) string {
	if value == nil {
		return ""
	}

	switch field.Type {
	case FieldTypeBool:
		if b, ok := value.(bool); ok {
			if b {
				return "Yes"
			}
			return "No"
		}
	case FieldTypeMulti:
		if values, ok := value.([]string); ok {
			return strings.Join(values, ", ")
		}
	}

	return fmt.Sprintf("%v", value)
}

func (m ConfiguratorModel) handleFieldValueChange(field ConfigurationField) (ConfiguratorModel, tea.Cmd) {
	// This would typically open a dedicated input interface for the field
	// For now, we'll just mark changes
	m.unsavedChanges = true
	delete(m.errors, field.Key) // Clear any existing error

	// Update global dirty state
	GlobalState.SetConfigurationDirty(true)

	// Publish configuration changed message
	changeMsg := ConfigurationMessage{
		BaseMessage: BaseMessage{MsgType: MsgConfigChanged, CreatedAt: time.Now()},
		Success:     true,
		Section:     m.sections[m.currentSection].Name,
		Data: map[string]interface{}{
			"field": field.Key,
			"value": field.Value,
		},
	}
	GlobalMessageBus.Publish(changeMsg)

	return m, nil
}

func (m ConfiguratorModel) searchFields(query string) []string {
	var results []string
	query = strings.ToLower(query)

	for _, section := range m.sections {
		for _, field := range section.Fields {
			if strings.Contains(strings.ToLower(field.Label), query) ||
				strings.Contains(strings.ToLower(field.Description), query) ||
				strings.Contains(strings.ToLower(field.Key), query) {
				results = append(results, fmt.Sprintf("%s > %s", section.Name, field.Label))
			}
		}
	}

	return results
}

func (m *ConfiguratorModel) refreshSectionsFromValues() {
	// Update field values in sections
	for i, section := range m.sections {
		for j, field := range section.Fields {
			if value, exists := m.fieldValues[field.Key]; exists {
				m.sections[i].Fields[j].Value = value
			}
		}
	}
}

// Tea commands
func (m ConfiguratorModel) loadConfigurationValues() tea.Cmd {
	return func() tea.Msg {
		if err := m.configManager.Load(); err != nil {
			// Publish configuration load error to message bus
			errorMsg := NewConfigLoadedMessage(false, make(map[string]interface{}), err.Error())
			GlobalMessageBus.Publish(errorMsg)

			// Update global state
			GlobalState.SetConfigurationLoaded(false, err)

			// Add error notification
			notification := Notification{
				Type:      "error",
				Title:     "Configuration Load Failed",
				Message:   fmt.Sprintf("Failed to load configuration: %v", err),
				Timestamp: time.Now(),
			}
			GlobalState.AddNotification(notification)

			return ConfigLoadedMsg{
				Success: false,
				Error:   err.Error(),
				Values:  make(map[string]interface{}),
			}
		}

		// Extract values from layered config
		values := make(map[string]interface{})
		// Add logic to extract from m.configManager.config

		// Publish configuration load success to message bus
		successMsg := NewConfigLoadedMessage(true, values, "")
		GlobalMessageBus.Publish(successMsg)

		// Update global state
		GlobalState.SetConfigurationLoaded(true, nil)

		return ConfigLoadedMsg{
			Success: true,
			Values:  values,
		}
	}
}

func (m ConfiguratorModel) saveConfiguration() tea.Cmd {
	return func() tea.Msg {
		// Apply field values to config manager
		m.applyValuesToConfig()

		if err := m.configManager.Save(); err != nil {
			// Publish configuration save error to message bus
			errorMsg := NewConfigSavedMessage(false, err.Error(), "all")
			GlobalMessageBus.Publish(errorMsg)

			// Also update global state
			GlobalState.SetConfigurationDirty(false) // Reset dirty flag even on error

			return ConfigSavedMsg{
				Success: false,
				Error:   err.Error(),
			}
		}

		// Publish configuration save success to message bus
		successMsg := NewConfigSavedMessage(true, "", "all")
		GlobalMessageBus.Publish(successMsg)

		// Update global state
		GlobalState.SetConfigurationDirty(false)

		// Add success notification
		notification := Notification{
			Type:      "success",
			Title:     "Configuration Saved",
			Message:   "Configuration changes have been saved successfully",
			Timestamp: time.Now(),
		}
		GlobalState.AddNotification(notification)

		return ConfigSavedMsg{
			Success: true,
		}
	}
}

func (m *ConfiguratorModel) applyValuesToConfig() {
	// Apply field values to the configuration manager
	// This would involve mapping field values to the appropriate config sections
	for key, value := range m.fieldValues {
		// Map to appropriate config section based on key prefix or schema
		// For example:
		// - global.* -> m.configManager.config.Global
		// - project.* -> m.configManager.config.Project
		// - local.* -> m.configManager.config.Local
		_ = key
		_ = value
	}
}

// Message types
type ConfigLoadedMsg struct {
	Success bool
	Error   string
	Values  map[string]interface{}
}

type ConfigSavedMsg struct {
	Success bool
	Error   string
}

type ConfigSaveCompleteMsg struct{}

type ValidationErrorMsg struct {
	Field string
	Error string
}

// List item for configuration sections
type ConfigSectionItem struct {
	name        string
	description string
	fieldCount  int
}

func (i ConfigSectionItem) FilterValue() string { return i.name }

type ConfigSectionDelegate struct{}

func (d ConfigSectionDelegate) Height() int                             { return 3 }
func (d ConfigSectionDelegate) Spacing() int                            { return 1 }
func (d ConfigSectionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ConfigSectionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ConfigSectionItem)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(strs ...string) string {
			if len(strs) > 0 {
				return selectedItemStyle.Render("▶ " + strs[0])
			}
			return selectedItemStyle.Render("▶ ")
		}
	}

	title := i.name
	desc := fmt.Sprintf("%s (%d fields)", i.description, i.fieldCount)

	fmt.Fprint(w, fn(fmt.Sprintf("%s\n%s", title, desc)))
}

var (
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99")).
				Bold(true).
				PaddingLeft(1)
)

// Schema creation helpers
func createDefaultConfigSchema() *ConfigSchema {
	return &ConfigSchema{
		Sections: []SectionSchema{
			{
				Name:        "Global Settings",
				Description: "System-wide MCF configuration",
				Fields: []FieldSchema{
					{
						Key:         "global.theme",
						Label:       "Theme",
						Type:        FieldTypeSelect,
						Description: "UI theme for MCF",
						Default:     "default",
						Options: []FieldOption{
							{Value: "default", Label: "Default", Default: true},
							{Value: "dark", Label: "Dark"},
							{Value: "light", Label: "Light"},
						},
					},
					{
						Key:         "global.editor",
						Label:       "Default Editor",
						Type:        FieldTypeSelect,
						Description: "Preferred code editor",
						Default:     "auto-detect",
						Options: []FieldOption{
							{Value: "auto-detect", Label: "Auto-detect", Default: true},
							{Value: "vscode", Label: "VS Code"},
							{Value: "vim", Label: "Vim"},
							{Value: "emacs", Label: "Emacs"},
						},
					},
				},
			},
			{
				Name:        "Project Configuration",
				Description: "Current project settings",
				Fields: []FieldSchema{
					{
						Key:         "project.name",
						Label:       "Project Name",
						Type:        FieldTypeText,
						Required:    true,
						Description: "Name of the current MCF project",
						Validation: ValidationRules{
							Required: true,
							MinLen:   func(i int) *int { return &i }(1),
							MaxLen:   func(i int) *int { return &i }(50),
						},
					},
					{
						Key:         "project.features",
						Label:       "Enabled Features",
						Type:        FieldTypeMulti,
						Description: "Features to enable for this project",
						Options: []FieldOption{
							{Value: "ai-agents", Label: "AI Agents", Default: true},
							{Value: "git-integration", Label: "Git Integration", Default: true},
							{Value: "serena", Label: "Serena Analysis"},
							{Value: "security", Label: "Security Features", Default: true},
						},
					},
				},
			},
			{
				Name:        "Developer Settings",
				Description: "Personal developer preferences",
				Fields: []FieldSchema{
					{
						Key:         "local.name",
						Label:       "Developer Name",
						Type:        FieldTypeText,
						Description: "Your name for git commits and signatures",
						Placeholder: "John Doe",
					},
					{
						Key:         "local.email",
						Label:       "Email Address",
						Type:        FieldTypeEmail,
						Description: "Your email address",
						Placeholder: "john@example.com",
						Validation: ValidationRules{
							Pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
						},
					},
				},
			},
		},
	}
}

func createConfigSections(schema *ConfigSchema, configManager *ConfigManager) []ConfigSection {
	var sections []ConfigSection

	for _, sectionSchema := range schema.Sections {
		section := ConfigSection{
			Name:        sectionSchema.Name,
			Description: sectionSchema.Description,
			Fields:      make([]ConfigurationField, len(sectionSchema.Fields)),
			Errors:      make(map[string]string),
		}

		for i, fieldSchema := range sectionSchema.Fields {
			field := ConfigurationField{
				Key:         fieldSchema.Key,
				Label:       fieldSchema.Label,
				Type:        fieldSchema.Type,
				Default:     fieldSchema.Default,
				Required:    fieldSchema.Required,
				Description: fieldSchema.Description,
				Placeholder: fieldSchema.Placeholder,
				HelpText:    fieldSchema.HelpText,
				Options:     fieldSchema.Options,
				Sensitive:   fieldSchema.Sensitive,
				Validator:   CreateValidator(fieldSchema.Validation),
			}

			section.Fields[i] = field
		}

		sections = append(sections, section)
	}

	return sections
}
