package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MCF Runner states
type MCFRunnerState int

const (
	MCFStateOperationSelect MCFRunnerState = iota
	MCFStateParameterInput
	MCFStateExecuting
	MCFStateResults
	MCFStateError
)

// MCF Operation types
type MCFOperationType int

const (
	MCFOpTypeAgent MCFOperationType = iota
	MCFOpTypeCommand
	MCFOpTypeTemplate
)

// MCF Operation item for list display
type MCFOperationItem struct {
	name        string
	description string
	opType      MCFOperationType
	path        string
	params      []MCFParameterDef
}

func (i MCFOperationItem) FilterValue() string { return i.name }
func (i MCFOperationItem) Title() string       { return i.name }
func (i MCFOperationItem) Description() string { return i.description }

// MCF Parameter definition
type MCFParameterDef struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Description string   `json:"description"`
	Default     string   `json:"default"`
	Options     []string `json:"options"`
}

// MCF Execution result
type MCFExecutionResult struct {
	Success   bool          `json:"success"`
	Output    string        `json:"output"`
	Error     string        `json:"error"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// MCF Messages
type MCFExecutionCompleteMsg struct {
	Result MCFExecutionResult
}

type MCFLoadOperationsMsg struct {
	Operations []MCFOperationItem
}

// MCF Runner model
type MCFRunnerModel struct {
	state        MCFRunnerState
	list         list.Model
	operations   []MCFOperationItem
	selected     *MCFOperationItem
	params       map[string]string
	inputs       []textinput.Model
	currentInput int
	spinner      spinner.Model
	result       *MCFExecutionResult
	error        error
	width        int
	height       int
	projectPath  string
}

func NewMCFRunnerModel(projectPath string) MCFRunnerModel {
	// Setup list
	items := []list.Item{}
	l := list.New(items, mcfItemDelegate{}, 0, 0)
	l.Title = "MCF Operations"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = mcfTitleStyle
	l.Styles.PaginationStyle = mcfPaginationStyle
	l.Styles.HelpStyle = mcfHelpStyle

	// Setup spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return MCFRunnerModel{
		state:       MCFStateOperationSelect,
		list:        l,
		operations:  []MCFOperationItem{},
		params:      make(map[string]string),
		inputs:      []textinput.Model{},
		spinner:     s,
		projectPath: projectPath,
	}
}

func (m MCFRunnerModel) Init() tea.Cmd {
	return tea.Batch(
		loadMCFOperations(m.projectPath),
		m.spinner.Tick,
	)
}

func (m MCFRunnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)

	case MCFLoadOperationsMsg:
		items := make([]list.Item, len(msg.Operations))
		for i, op := range msg.Operations {
			items[i] = op
		}
		m.operations = msg.Operations
		m.list.SetItems(items)

	case MCFExecutionCompleteMsg:
		m.result = &msg.Result
		m.state = MCFStateResults
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case MCFStateOperationSelect:
			return m.handleMCFOperationSelect(msg)
		case MCFStateParameterInput:
			return m.handleMCFParameterInput(msg)
		case MCFStateResults, MCFStateError:
			return m.handleMCFResultsView(msg)
		}
	}

	switch m.state {
	case MCFStateOperationSelect:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

	case MCFStateParameterInput:
		for i := range m.inputs {
			var cmd tea.Cmd
			m.inputs[i], cmd = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}

	case MCFStateExecuting:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MCFRunnerModel) handleMCFOperationSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		return m, tea.Quit
	case "enter":
		if item, ok := m.list.SelectedItem().(MCFOperationItem); ok {
			m.selected = &item
			if len(item.params) > 0 {
				return m.setupMCFParameterInput()
			} else {
				return m.executeMCFOperation()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MCFRunnerModel) handleMCFParameterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = MCFStateOperationSelect
		return m, nil
	case "tab", "shift+tab", "up", "down":
		return m.navigateMCFInputs(msg.String())
	case "enter":
		if m.currentInput == len(m.inputs)-1 {
			// Last input, execute operation
			m.collectMCFParameterValues()
			return m.executeMCFOperation()
		} else {
			// Move to next input
			m.currentInput++
			m.updateMCFInputFocus()
		}
	}

	var cmd tea.Cmd
	if m.currentInput < len(m.inputs) {
		m.inputs[m.currentInput], cmd = m.inputs[m.currentInput].Update(msg)
	}
	return m, cmd
}

func (m MCFRunnerModel) handleMCFResultsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.state = MCFStateOperationSelect
		m.result = nil
		m.error = nil
		return m, nil
	case "r":
		// Re-run the operation
		if m.selected != nil {
			if len(m.selected.params) > 0 {
				return m.setupMCFParameterInput()
			} else {
				return m.executeMCFOperation()
			}
		}
	}
	return m, nil
}

func (m MCFRunnerModel) navigateMCFInputs(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "tab", "down":
		m.currentInput++
		if m.currentInput >= len(m.inputs) {
			m.currentInput = 0
		}
	case "shift+tab", "up":
		m.currentInput--
		if m.currentInput < 0 {
			m.currentInput = len(m.inputs) - 1
		}
	}
	m.updateMCFInputFocus()
	return m, nil
}

func (m MCFRunnerModel) setupMCFParameterInput() (MCFRunnerModel, tea.Cmd) {
	m.state = MCFStateParameterInput
	m.inputs = make([]textinput.Model, len(m.selected.params))

	for i, param := range m.selected.params {
		input := textinput.New()
		input.Placeholder = param.Description
		input.CharLimit = 256
		input.Width = 50

		if param.Default != "" {
			input.SetValue(param.Default)
		}

		m.inputs[i] = input
	}

	m.currentInput = 0
	m.updateMCFInputFocus()

	return m, m.inputs[0].Focus()
}

func (m MCFRunnerModel) updateMCFInputFocus() {
	for i := range m.inputs {
		if i == m.currentInput {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m MCFRunnerModel) collectMCFParameterValues() {
	m.params = make(map[string]string)
	for i, param := range m.selected.params {
		if i < len(m.inputs) {
			m.params[param.Name] = m.inputs[i].Value()
		}
	}
}

func (m MCFRunnerModel) executeMCFOperation() (MCFRunnerModel, tea.Cmd) {
	m.state = MCFStateExecuting
	return m, executeMCFOperation(*m.selected, m.params)
}

func (m MCFRunnerModel) View() string {
	switch m.state {
	case MCFStateOperationSelect:
		return m.mcfOperationSelectView()
	case MCFStateParameterInput:
		return m.mcfParameterInputView()
	case MCFStateExecuting:
		return m.mcfExecutingView()
	case MCFStateResults:
		return m.mcfResultsView()
	case MCFStateError:
		return m.mcfErrorView()
	default:
		return "Unknown state"
	}
}

func (m MCFRunnerModel) mcfOperationSelectView() string {
	if len(m.operations) == 0 {
		return mcfLoadingStyle.Render("Loading operations...")
	}

	header := mcfHeaderStyle.Render("🚀 MCF Runner - Select Operation")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.list.View(),
		"",
		mcfHelpStyle.Render("↑/↓: navigate • enter: select • q: quit • /: filter"),
	)
}

func (m MCFRunnerModel) mcfParameterInputView() string {
	if m.selected == nil {
		return "No operation selected"
	}

	header := mcfHeaderStyle.Render(fmt.Sprintf("⚙️ Configure: %s", m.selected.name))

	var inputs []string
	for i, param := range m.selected.params {
		label := param.Name
		if param.Required {
			label += " *"
		}

		labelStyle := mcfInputLabelStyle
		if i == m.currentInput {
			labelStyle = mcfFocusedInputLabelStyle
		}

		inputs = append(inputs, labelStyle.Render(label))
		inputs = append(inputs, m.inputs[i].View())

		if param.Description != "" {
			inputs = append(inputs, mcfHelpStyle.Render(param.Description))
		}
		inputs = append(inputs, "")
	}

	footer := mcfHelpStyle.Render("tab/↑/↓: navigate • enter: next/execute • esc: back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		strings.Join(inputs, "\n"),
		footer,
	)
}

func (m MCFRunnerModel) mcfExecutingView() string {
	if m.selected == nil {
		return "No operation selected"
	}

	header := mcfHeaderStyle.Render(fmt.Sprintf("🔄 Executing: %s", m.selected.name))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		"",
		header,
		"",
		lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " Please wait..."),
		"",
		"",
	)
}

func (m MCFRunnerModel) mcfResultsView() string {
	if m.result == nil {
		return "No results available"
	}

	var header string
	if m.result.Success {
		header = mcfSuccessHeaderStyle.Render("✅ Execution Successful")
	} else {
		header = mcfErrorHeaderStyle.Render("❌ Execution Failed")
	}

	info := mcfInfoStyle.Render(fmt.Sprintf(
		"Operation: %s | Duration: %v | Time: %s",
		m.selected.name,
		m.result.Duration.Round(time.Millisecond),
		m.result.Timestamp.Format("15:04:05"),
	))

	var content string
	if m.result.Success {
		if m.result.Output != "" {
			content = mcfOutputBoxStyle.Render(m.result.Output)
		} else {
			content = mcfSuccessStyle.Render("Operation completed successfully (no output)")
		}
	} else {
		content = mcfErrorBoxStyle.Render(m.result.Error)
	}

	footer := mcfHelpStyle.Render("r: run again • esc/q: back to operations")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		info,
		"",
		content,
		"",
		footer,
	)
}

func (m MCFRunnerModel) mcfErrorView() string {
	header := mcfErrorHeaderStyle.Render("💥 System Error")

	var errorMsg string
	if m.error != nil {
		errorMsg = m.error.Error()
	} else {
		errorMsg = "Unknown error occurred"
	}

	content := mcfErrorBoxStyle.Render(errorMsg)
	footer := mcfHelpStyle.Render("esc/q: back to operations")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		content,
		"",
		footer,
	)
}

// MCF Commands
func loadMCFOperations(projectPath string) tea.Cmd {
	return func() tea.Msg {
		var operations []MCFOperationItem

		claudeDir := filepath.Join(projectPath, ".claude")

		// Load agents
		agentsDir := filepath.Join(claudeDir, "agents")
		if agents, err := loadMCFAgents(agentsDir); err == nil {
			operations = append(operations, agents...)
		}

		// Load commands
		commandsDir := filepath.Join(claudeDir, "commands")
		if commands, err := loadMCFCommands(commandsDir); err == nil {
			operations = append(operations, commands...)
		}

		// Load templates
		templatesDir := filepath.Join(claudeDir, "commands", "templates")
		if templates, err := loadMCFTemplates(templatesDir); err == nil {
			operations = append(operations, templates...)
		}

		return MCFLoadOperationsMsg{Operations: operations}
	}
}

func loadMCFAgents(agentsDir string) ([]MCFOperationItem, error) {
	var agents []MCFOperationItem

	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		path := filepath.Join(agentsDir, entry.Name())

		// Read first few lines to get description
		description := getMCFFileDescription(path)

		agents = append(agents, MCFOperationItem{
			name:        name,
			description: description,
			opType:      MCFOpTypeAgent,
			path:        path,
			params: []MCFParameterDef{
				{Name: "prompt", Type: "text", Required: true, Description: "Your request or question for the agent"},
			},
		})
	}

	return agents, nil
}

func loadMCFCommands(commandsDir string) ([]MCFOperationItem, error) {
	var commands []MCFOperationItem

	// Walk through command directories
	err := filepath.Walk(commandsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		// Skip templates directory (handled separately)
		if strings.Contains(path, "templates") {
			return nil
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Get relative path from commands directory
		relPath, _ := filepath.Rel(commandsDir, path)
		name := strings.TrimSuffix(relPath, ".md")
		name = strings.ReplaceAll(name, "/", " > ")

		description := getMCFFileDescription(path)

		// Try to parse parameters from command file
		params := parseMCFCommandParams(path)

		commands = append(commands, MCFOperationItem{
			name:        name,
			description: description,
			opType:      MCFOpTypeCommand,
			path:        path,
			params:      params,
		})

		return nil
	})

	return commands, err
}

func loadMCFTemplates(templatesDir string) ([]MCFOperationItem, error) {
	var templates []MCFOperationItem

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		path := filepath.Join(templatesDir, entry.Name())
		description := getMCFFileDescription(path)

		// Templates typically have configurable parameters
		params := parseMCFTemplateParams(path)

		templates = append(templates, MCFOperationItem{
			name:        "Template: " + name,
			description: description,
			opType:      MCFOpTypeTemplate,
			path:        path,
			params:      params,
		})
	}

	return templates, nil
}

func executeMCFOperation(item MCFOperationItem, params map[string]string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		var result MCFExecutionResult
		result.Timestamp = start

		switch item.opType {
		case MCFOpTypeAgent:
			result = executeMCFAgent(item, params)
		case MCFOpTypeCommand:
			result = executeMCFCommand(item, params)
		case MCFOpTypeTemplate:
			result = executeMCFTemplate(item, params)
		default:
			result = MCFExecutionResult{
				Success: false,
				Error:   "Unknown operation type",
			}
		}

		result.Duration = time.Since(start)
		return MCFExecutionCompleteMsg{Result: result}
	}
}

func executeMCFAgent(item MCFOperationItem, params map[string]string) MCFExecutionResult {
	prompt, exists := params["prompt"]
	if !exists || prompt == "" {
		return MCFExecutionResult{
			Success: false,
			Error:   "Prompt is required for agent execution",
		}
	}

	// Execute using Claude Code CLI with the agent
	cmd := exec.Command("claude", "agent", filepath.Base(strings.TrimSuffix(item.path, ".md")), prompt)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return MCFExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Agent execution failed: %v\nOutput: %s", err, string(output)),
		}
	}

	return MCFExecutionResult{
		Success: true,
		Output:  string(output),
	}
}

func executeMCFCommand(item MCFOperationItem, params map[string]string) MCFExecutionResult {
	// Build command arguments
	args := []string{"command"}

	// Get command name from path
	relPath, _ := filepath.Rel(filepath.Join(filepath.Dir(item.path), ".."), item.path)
	commandName := strings.TrimSuffix(relPath, ".md")
	commandName = strings.ReplaceAll(commandName, "/", ":")

	args = append(args, commandName)

	// Add parameters as arguments
	for name, value := range params {
		if value != "" {
			args = append(args, fmt.Sprintf("--%s=%s", name, value))
		}
	}

	cmd := exec.Command("claude", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return MCFExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Command execution failed: %v\nOutput: %s", err, string(output)),
		}
	}

	return MCFExecutionResult{
		Success: true,
		Output:  string(output),
	}
}

func executeMCFTemplate(item MCFOperationItem, params map[string]string) MCFExecutionResult {
	// Templates are usually applied via command system
	cmd := exec.Command("claude", "template", "apply", filepath.Base(strings.TrimSuffix(item.path, ".md")))

	// Set environment variables for template parameters
	env := os.Environ()
	for name, value := range params {
		env = append(env, fmt.Sprintf("TEMPLATE_%s=%s", strings.ToUpper(name), value))
	}
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		return MCFExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Template execution failed: %v\nOutput: %s", err, string(output)),
		}
	}

	return MCFExecutionResult{
		Success: true,
		Output:  string(output),
	}
}

// MCF Utility functions
func getMCFFileDescription(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "No description available"
	}
	defer file.Close()

	// Read first few lines to extract description
	content := make([]byte, 512)
	n, _ := file.Read(content)

	lines := strings.Split(string(content[:n]), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "```") {
			if len(line) > 80 {
				return line[:77] + "..."
			}
			return line
		}
	}

	return "No description available"
}

func parseMCFCommandParams(path string) []MCFParameterDef {
	// Basic parameter parsing - look for parameter definitions in markdown
	// This is a simplified implementation
	return []MCFParameterDef{
		{Name: "input", Type: "text", Required: false, Description: "Input parameter for the command"},
	}
}

func parseMCFTemplateParams(path string) []MCFParameterDef {
	// Parse template parameters from file content
	return []MCFParameterDef{
		{Name: "name", Type: "text", Required: true, Description: "Template name parameter"},
		{Name: "type", Type: "select", Required: false, Description: "Template type", Options: []string{"component", "service", "utility"}},
	}
}

// Custom MCF list item delegate
type mcfItemDelegate struct{}

func (d mcfItemDelegate) Height() int                             { return 2 }
func (d mcfItemDelegate) Spacing() int                            { return 1 }
func (d mcfItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d mcfItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(MCFOperationItem)
	if !ok {
		return
	}

	var icon string
	switch i.opType {
	case MCFOpTypeAgent:
		icon = "🤖"
	case MCFOpTypeCommand:
		icon = "⚡"
	case MCFOpTypeTemplate:
		icon = "📄"
	}

	title := fmt.Sprintf("%s %s", icon, i.name)
	desc := i.description

	var rendered string
	if index == m.Index() {
		rendered = mcfSelectedItemStyle.Render(fmt.Sprintf("%s\n%s", title, desc))
	} else {
		rendered = mcfNormalItemStyle.Render(fmt.Sprintf("%s\n%s", title, desc))
	}

	fmt.Fprint(w, rendered)
}

// MCF Styles
var (
	mcfTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Padding(0, 1)

	mcfHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)

	mcfSuccessHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("34")).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("34")).
				Padding(0, 1)

	mcfErrorHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("196")).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("196")).
				Padding(0, 1)

	mcfInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	mcfHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	mcfPaginationStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	mcfSelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99")).
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("99")).
				Padding(0, 1)

	mcfNormalItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Padding(0, 1)

	mcfInputLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Bold(true)

	mcfFocusedInputLabelStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("99")).
					Bold(true)

	mcfLoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center)

	mcfSuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34"))

	mcfOutputBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238")).
				Padding(1).
				MaxWidth(80)

	mcfErrorBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("196")).
				Padding(1).
				MaxWidth(80).
				Foreground(lipgloss.Color("196"))
)
