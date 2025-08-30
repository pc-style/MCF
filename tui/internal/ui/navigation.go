package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Navigation views
type View int

const (
	DashboardView View = iota
	AgentsView
	CommandsView
	LogsView
	ConfigView
	CommandBarView
)

var ViewNames = map[View]string{
	DashboardView:  "Dashboard",
	AgentsView:     "Agents",
	CommandsView:   "Commands",
	LogsView:       "Logs",
	ConfigView:     "Config",
	CommandBarView: "Command Bar",
}

// Navigation state and history
type Navigation struct {
	currentView View
	viewHistory []View
	breadcrumb  []string
	theme       *Theme
}

func NewNavigation(theme *Theme) *Navigation {
	return &Navigation{
		currentView: DashboardView,
		viewHistory: []View{DashboardView},
		breadcrumb:  []string{"Dashboard"},
		theme:       theme,
	}
}

// Navigation methods
func (n *Navigation) SetView(view View) {
	if view != n.currentView {
		n.viewHistory = append(n.viewHistory, view)
		n.currentView = view
		n.updateBreadcrumb()
	}
}

func (n *Navigation) GoBack() View {
	if len(n.viewHistory) > 1 {
		n.viewHistory = n.viewHistory[:len(n.viewHistory)-1]
		n.currentView = n.viewHistory[len(n.viewHistory)-1]
		n.updateBreadcrumb()
	}
	return n.currentView
}

func (n *Navigation) GetCurrentView() View {
	return n.currentView
}

func (n *Navigation) GetViewName(view View) string {
	if name, ok := ViewNames[view]; ok {
		return name
	}
	return "Unknown"
}

func (n *Navigation) updateBreadcrumb() {
	n.breadcrumb = []string{}
	for _, view := range n.viewHistory {
		n.breadcrumb = append(n.breadcrumb, n.GetViewName(view))
	}
}

// Render navigation components
func (n *Navigation) RenderTabBar(width int) string {
	tabs := []string{}
	tabWidth := (width - 10) / 5 // 5 main views, leave margin

	views := []View{DashboardView, AgentsView, CommandsView, LogsView, ConfigView}

	for _, view := range views {
		name := n.GetViewName(view)
		if len(name) > tabWidth-4 {
			name = name[:tabWidth-4] + "..."
		}

		var style lipgloss.Style
		if view == n.currentView {
			style = n.theme.TabActive
		} else {
			style = n.theme.TabInactive
		}

		tabs = append(tabs, style.Width(tabWidth).Align(lipgloss.Center).Render(name))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (n *Navigation) RenderBreadcrumb() string {
	if len(n.breadcrumb) <= 1 {
		return ""
	}

	parts := make([]string, len(n.breadcrumb))
	for i, part := range n.breadcrumb {
		if i == len(n.breadcrumb)-1 {
			// Current view - highlighted
			parts[i] = n.theme.Info.Render(part)
		} else {
			// Previous views - muted
			parts[i] = n.theme.Muted.Render(part)
		}
	}

	breadcrumbText := strings.Join(parts, n.theme.Muted.Render(" › "))
	return n.theme.Breadcrumb.Render("📍 " + breadcrumbText)
}

func (n *Navigation) RenderShortcuts() string {
	shortcuts := []string{
		"Tab: Next View",
		"Shift+Tab: Previous View",
		"Esc: Back/Dashboard",
		":: Command Bar",
		"?: Help",
		"q: Quit",
	}

	shortcutText := ""
	for i, shortcut := range shortcuts {
		if i > 0 {
			shortcutText += " │ "
		}
		parts := strings.SplitN(shortcut, ": ", 2)
		if len(parts) == 2 {
			shortcutText += n.theme.Info.Render(parts[0]) + ": " + n.theme.Muted.Render(parts[1])
		} else {
			shortcutText += n.theme.Muted.Render(shortcut)
		}
	}

	return n.theme.Panel.
		Width(lipgloss.Width(shortcutText) + 4).
		Render(shortcutText)
}

// Quick action menu
type QuickAction struct {
	Key         string
	Label       string
	Description string
	Command     string
}

var QuickActions = []QuickAction{
	{"1", "Agents Status", "View all agent statuses", "mcf agents status"},
	{"2", "Start Serena", "Start Serena integration", "mcf serena start"},
	{"3", "Run Tests", "Execute test suite", "mcf test"},
	{"4", "Deploy", "Deploy application", "mcf deploy"},
	{"5", "Logs", "View recent logs", "mcf logs tail"},
	{"6", "Health Check", "System health check", "mcf health"},
}

func (n *Navigation) RenderQuickActions(selectedIdx int) string {
	content := n.theme.Subtitle.Render("Quick Actions") + "\n\n"

	for i, action := range QuickActions {
		var style lipgloss.Style
		if i == selectedIdx {
			style = n.theme.ListItemActive
		} else {
			style = n.theme.ListItem
		}

		line := fmt.Sprintf("%s  %s", action.Key, action.Label)
		if i == selectedIdx {
			line += "\n    " + n.theme.Muted.Render(action.Description)
		}

		content += style.Render(line) + "\n"
	}

	return n.theme.Panel.Render(content)
}

// Help system
func (n *Navigation) RenderHelp(width int) string {
	sections := map[string][]string{
		"Navigation": {
			"Tab / Shift+Tab - Switch between views",
			"Esc - Go back or return to dashboard",
			": - Open command bar",
			"? - Toggle this help",
			"q / Ctrl+C - Quit application",
		},
		"Dashboard": {
			"j/k or ↑/↓ - Navigate quick actions",
			"Enter - Execute selected action",
			"r - Refresh system status",
			"1-6 - Quick action shortcuts",
		},
		"Agents View": {
			"j/k or ↑/↓ - Navigate agent list",
			"Enter - View agent details",
			"s - Start/stop selected agent",
			"r - Refresh agent status",
			"l - View agent logs",
		},
		"Commands View": {
			"j/k or ↑/↓ - Navigate command history",
			"Enter - Re-execute command",
			"d - Delete command from history",
			"c - Clear command history",
			"/ - Search commands",
		},
		"Logs View": {
			"j/k or ↑/↓ - Scroll logs",
			"g/G - Go to top/bottom",
			"/ - Search logs",
			"f - Follow/unfollow logs",
			"c - Clear log view",
		},
	}

	content := n.theme.Title.Render("MCF TUI Help") + "\n\n"

	for section, items := range sections {
		content += n.theme.Subtitle.Render(section) + "\n"
		for _, item := range items {
			parts := strings.SplitN(item, " - ", 2)
			if len(parts) == 2 {
				content += "  " + n.theme.Info.Render(parts[0]) + " - " + n.theme.Body.Render(parts[1]) + "\n"
			} else {
				content += "  " + n.theme.Body.Render(item) + "\n"
			}
		}
		content += "\n"
	}

	return n.theme.Panel.
		Width(width - 4).
		Height(25).
		Render(content)
}
