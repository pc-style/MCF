package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Interactive List Component
type InteractiveList struct {
	theme    *Theme
	items    []ListItem
	selected int
	focused  bool
	title    string
	height   int
}

type ListItem struct {
	Title       string
	Description string
	Status      string
	Metadata    map[string]interface{}
}

func NewInteractiveList(theme *Theme, title string, height int) *InteractiveList {
	return &InteractiveList{
		theme:    theme,
		items:    []ListItem{},
		selected: 0,
		focused:  false,
		title:    title,
		height:   height,
	}
}

func (l *InteractiveList) SetItems(items []ListItem) {
	l.items = items
	if l.selected >= len(items) {
		l.selected = len(items) - 1
	}
	if l.selected < 0 {
		l.selected = 0
	}
}

func (l *InteractiveList) SetFocus(focused bool) {
	l.focused = focused
}

func (l *InteractiveList) GetSelected() int {
	return l.selected
}

func (l *InteractiveList) GetSelectedItem() *ListItem {
	if l.selected >= 0 && l.selected < len(l.items) {
		return &l.items[l.selected]
	}
	return nil
}

func (l *InteractiveList) Update(msg tea.Msg) (*InteractiveList, tea.Cmd) {
	if !l.focused {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if l.selected > 0 {
				l.selected--
			}
		case "down", "j":
			if l.selected < len(l.items)-1 {
				l.selected++
			}
		case "home", "g":
			l.selected = 0
		case "end", "G":
			if len(l.items) > 0 {
				l.selected = len(l.items) - 1
			}
		}
	}

	return l, nil
}

func (l *InteractiveList) Render(width int) string {
	if len(l.items) == 0 {
		empty := l.theme.Muted.Render("No items available")
		return RenderBox(empty, l.title, width, l.height, l.theme)
	}

	content := ""
	visibleItems := l.height - 4 // Account for title and padding

	// Calculate scroll offset
	scrollOffset := 0
	if l.selected >= visibleItems {
		scrollOffset = l.selected - visibleItems + 1
	}

	// Render visible items
	for i := 0; i < visibleItems && i+scrollOffset < len(l.items); i++ {
		idx := i + scrollOffset
		item := l.items[idx]

		var style lipgloss.Style
		cursor := "  "

		if idx == l.selected && l.focused {
			style = l.theme.ListItemActive
			cursor = "► "
		} else {
			style = l.theme.ListItem
		}

		// Item line with cursor, title, and status
		line := cursor + item.Title
		if item.Status != "" {
			statusIndicator := RenderStatusIndicator(item.Status, l.theme)
			line = fmt.Sprintf("%-*s %s", width-20, line, statusIndicator)
		}

		content += style.Render(line) + "\n"

		// Show description for selected item
		if idx == l.selected && item.Description != "" {
			desc := "    " + item.Description
			if len(desc) > width-6 {
				desc = desc[:width-9] + "..."
			}
			content += l.theme.Muted.Render(desc) + "\n"
		}
	}

	// Scroll indicator
	if len(l.items) > visibleItems {
		scrollInfo := fmt.Sprintf(" (%d-%d of %d)",
			scrollOffset+1,
			min(scrollOffset+visibleItems, len(l.items)),
			len(l.items))
		content += "\n" + l.theme.Muted.Render(scrollInfo)
	}

	return RenderBox(content, l.title, width, l.height, l.theme)
}

// Command Input Component
type CommandInput struct {
	theme       *Theme
	input       textinput.Model
	suggestions []string
	selected    int
	focused     bool
	history     []string
	historyIdx  int
}

func NewCommandInput(theme *Theme) *CommandInput {
	input := textinput.New()
	input.Placeholder = "Enter MCF command..."
	input.Focus()

	return &CommandInput{
		theme:       theme,
		input:       input,
		suggestions: []string{},
		selected:    0,
		focused:     true,
		history:     []string{},
		historyIdx:  -1,
	}
}

func (c *CommandInput) SetFocus(focused bool) {
	c.focused = focused
	if focused {
		c.input.Focus()
	} else {
		c.input.Blur()
	}
}

func (c *CommandInput) SetSuggestions(suggestions []string) {
	c.suggestions = suggestions
	c.selected = 0
}

func (c *CommandInput) AddToHistory(command string) {
	if command != "" && (len(c.history) == 0 || c.history[0] != command) {
		c.history = append([]string{command}, c.history...)
		if len(c.history) > 50 { // Limit history size
			c.history = c.history[:50]
		}
	}
	c.historyIdx = -1
}

func (c *CommandInput) GetValue() string {
	return c.input.Value()
}

func (c *CommandInput) Clear() {
	c.input.SetValue("")
	c.suggestions = []string{}
	c.selected = 0
}

func (c *CommandInput) Update(msg tea.Msg) (*CommandInput, tea.Cmd) {
	if !c.focused {
		return c, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if len(c.suggestions) > 0 && c.selected > 0 {
				c.selected--
			} else if c.historyIdx < len(c.history)-1 {
				c.historyIdx++
				if c.historyIdx < len(c.history) {
					c.input.SetValue(c.history[c.historyIdx])
				}
			}
		case "down":
			if len(c.suggestions) > 0 && c.selected < len(c.suggestions)-1 {
				c.selected++
			} else if c.historyIdx > -1 {
				c.historyIdx--
				if c.historyIdx == -1 {
					c.input.SetValue("")
				} else {
					c.input.SetValue(c.history[c.historyIdx])
				}
			}
		case "tab":
			if len(c.suggestions) > 0 {
				c.input.SetValue(c.suggestions[c.selected])
				c.suggestions = []string{}
			}
		default:
			c.input, cmd = c.input.Update(msg)
			// Update suggestions based on input
			c.updateSuggestions()
		}
	default:
		c.input, cmd = c.input.Update(msg)
	}

	return c, cmd
}

func (c *CommandInput) updateSuggestions() {
	value := c.input.Value()
	if len(value) < 2 {
		c.suggestions = []string{}
		return
	}

	// Common MCF commands for suggestions
	commands := []string{
		"mcf agents status",
		"mcf agents start",
		"mcf agents stop",
		"mcf serena start",
		"mcf serena stop",
		"mcf serena status",
		"mcf deploy --stage dev",
		"mcf deploy --stage prod",
		"mcf test --coverage",
		"mcf test --unit",
		"mcf logs tail",
		"mcf logs search",
		"mcf health check",
		"mcf config show",
		"mcf config set",
	}

	suggestions := []string{}
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd), strings.ToLower(value)) {
			suggestions = append(suggestions, cmd)
		}
	}

	c.suggestions = suggestions
	c.selected = 0
}

func (c *CommandInput) Render(width int) string {
	var style lipgloss.Style
	if c.focused {
		style = c.theme.InputFocused
	} else {
		style = c.theme.Input
	}

	inputBox := style.Width(width - 4).Render(c.input.View())
	content := inputBox

	// Show suggestions
	if len(c.suggestions) > 0 && c.focused {
		content += "\n\n" + c.theme.Subtitle.Render("Suggestions:") + "\n"

		maxSuggestions := min(5, len(c.suggestions))
		for i := 0; i < maxSuggestions; i++ {
			suggestion := c.suggestions[i]
			if len(suggestion) > width-8 {
				suggestion = suggestion[:width-11] + "..."
			}

			var suggestionStyle lipgloss.Style
			if i == c.selected {
				suggestionStyle = c.theme.ListItemActive
			} else {
				suggestionStyle = c.theme.ListItem
			}

			content += suggestionStyle.Render(fmt.Sprintf("  %s", suggestion)) + "\n"
		}
	}

	// Show recent commands
	if len(c.history) > 0 && c.input.Value() == "" && c.focused {
		content += "\n" + c.theme.Subtitle.Render("Recent Commands:") + "\n"

		maxHistory := min(3, len(c.history))
		for i := 0; i < maxHistory; i++ {
			cmd := c.history[i]
			if len(cmd) > width-8 {
				cmd = cmd[:width-11] + "..."
			}
			content += c.theme.Muted.Render(fmt.Sprintf("  %s", cmd)) + "\n"
		}
	}

	return content
}

// Log Viewer Component
type LogViewer struct {
	theme       *Theme
	logs        []LogEntry
	filter      string
	following   bool
	scrollPos   int
	height      int
	width       int
	searchMode  bool
	searchInput textinput.Model
}

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Component string
	Message   string
}

func NewLogViewer(theme *Theme, height int) *LogViewer {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search logs..."

	return &LogViewer{
		theme:       theme,
		logs:        []LogEntry{},
		filter:      "",
		following:   false,
		scrollPos:   0,
		height:      height,
		searchMode:  false,
		searchInput: searchInput,
	}
}

func (lv *LogViewer) AddLog(entry LogEntry) {
	lv.logs = append(lv.logs, entry)
	if lv.following {
		lv.scrollToBottom()
	}

	// Limit log history
	if len(lv.logs) > 1000 {
		lv.logs = lv.logs[100:] // Keep last 900 entries
	}
}

func (lv *LogViewer) SetFollowing(following bool) {
	lv.following = following
	if following {
		lv.scrollToBottom()
	}
}

func (lv *LogViewer) SetSearchMode(enabled bool) {
	lv.searchMode = enabled
	if enabled {
		lv.searchInput.Focus()
	} else {
		lv.searchInput.Blur()
		lv.filter = ""
	}
}

func (lv *LogViewer) Clear() {
	lv.logs = []LogEntry{}
	lv.scrollPos = 0
}

func (lv *LogViewer) scrollToBottom() {
	if len(lv.logs) > lv.height-4 {
		lv.scrollPos = len(lv.logs) - (lv.height - 4)
	} else {
		lv.scrollPos = 0
	}
}

func (lv *LogViewer) Update(msg tea.Msg) (*LogViewer, tea.Cmd) {
	var cmd tea.Cmd

	if lv.searchMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				lv.filter = lv.searchInput.Value()
				lv.searchMode = false
				lv.searchInput.Blur()
			case "esc":
				lv.searchMode = false
				lv.searchInput.Blur()
				lv.filter = ""
			default:
				lv.searchInput, cmd = lv.searchInput.Update(msg)
			}
		default:
			lv.searchInput, cmd = lv.searchInput.Update(msg)
		}
		return lv, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if lv.scrollPos > 0 {
				lv.scrollPos--
				lv.following = false
			}
		case "down", "j":
			maxScroll := len(lv.logs) - (lv.height - 4)
			if lv.scrollPos < maxScroll {
				lv.scrollPos++
			}
		case "home", "g":
			lv.scrollPos = 0
			lv.following = false
		case "end", "G":
			lv.scrollToBottom()
		case "f":
			lv.following = !lv.following
			if lv.following {
				lv.scrollToBottom()
			}
		case "/":
			lv.searchMode = true
			lv.searchInput.Focus()
		case "c":
			lv.Clear()
		}
	}

	return lv, nil
}

func (lv *LogViewer) Render(width int) string {
	lv.width = width

	if lv.searchMode {
		return lv.renderSearchMode()
	}

	content := ""

	// Filter logs
	filteredLogs := lv.logs
	if lv.filter != "" {
		filteredLogs = []LogEntry{}
		for _, entry := range lv.logs {
			if strings.Contains(strings.ToLower(entry.Message), strings.ToLower(lv.filter)) ||
				strings.Contains(strings.ToLower(entry.Component), strings.ToLower(lv.filter)) {
				filteredLogs = append(filteredLogs, entry)
			}
		}
	}

	// Render visible logs
	visibleLines := lv.height - 4
	startIdx := lv.scrollPos
	endIdx := min(startIdx+visibleLines, len(filteredLogs))

	for i := startIdx; i < endIdx; i++ {
		entry := filteredLogs[i]
		content += lv.renderLogEntry(entry) + "\n"
	}

	// Status line
	statusLine := ""
	if lv.following {
		statusLine += lv.theme.Success.Render("● FOLLOWING")
	} else {
		statusLine += lv.theme.Muted.Render("● PAUSED")
	}

	if lv.filter != "" {
		statusLine += " " + lv.theme.Info.Render(fmt.Sprintf("Filter: '%s'", lv.filter))
	}

	statusLine += " " + lv.theme.Muted.Render(fmt.Sprintf("(%d/%d)",
		min(endIdx, len(filteredLogs)), len(filteredLogs)))

	title := "Logs " + statusLine
	return RenderBox(content, title, width, lv.height, lv.theme)
}

func (lv *LogViewer) renderSearchMode() string {
	content := lv.theme.Subtitle.Render("Search Logs") + "\n\n"
	content += lv.theme.InputFocused.Width(lv.width - 6).Render(lv.searchInput.View())
	content += "\n\n" + lv.theme.Muted.Render("Press Enter to search, Esc to cancel")

	return RenderBox(content, "Search", lv.width, 8, lv.theme)
}

func (lv *LogViewer) renderLogEntry(entry LogEntry) string {
	timestamp := entry.Timestamp.Format("15:04:05")

	var levelStyle lipgloss.Style
	switch entry.Level {
	case "ERROR":
		levelStyle = lv.theme.Error
	case "WARN":
		levelStyle = lv.theme.Warning
	case "INFO":
		levelStyle = lv.theme.Info
	default:
		levelStyle = lv.theme.Muted
	}

	line := fmt.Sprintf("[%s] %s %s: %s",
		lv.theme.Muted.Render(timestamp),
		levelStyle.Render(fmt.Sprintf("%-5s", entry.Level)),
		lv.theme.Info.Render(entry.Component),
		entry.Message,
	)

	if len(line) > lv.width-4 {
		line = line[:lv.width-7] + "..."
	}

	return line
}
