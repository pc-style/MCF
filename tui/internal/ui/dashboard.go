package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Dashboard data structures
type AgentStatus struct {
	Name        string
	Status      string
	LastSeen    time.Time
	TasksActive int
	TasksTotal  int
}

type SystemHealth struct {
	MCFVersion   string
	Uptime       time.Duration
	MemoryUsage  float64
	CPUUsage     float64
	DiskUsage    float64
	ActiveAgents int
	TotalAgents  int
	SerenaStatus string
	ClaudeStatus string
}

type RecentActivity struct {
	Timestamp time.Time
	Type      string // "command", "agent", "error", "info"
	Message   string
	Details   string
}

type Dashboard struct {
	theme               *Theme
	systemHealth        SystemHealth
	agentStatuses       []AgentStatus
	recentActivity      []RecentActivity
	commandHistory      []string
	selectedQuickAction int
	showHelp            bool
	lastRefresh         time.Time
}

func NewDashboard(theme *Theme) *Dashboard {
	return &Dashboard{
		theme:               theme,
		systemHealth:        SystemHealth{},
		agentStatuses:       []AgentStatus{},
		recentActivity:      []RecentActivity{},
		commandHistory:      []string{},
		selectedQuickAction: 0,
		showHelp:            false,
		lastRefresh:         time.Now(),
	}
}

func (d *Dashboard) Update() {
	d.lastRefresh = time.Now()
	// Real data will be updated via UpdateWithMCFData
}

// SetSystemHealth updates system health with real data
func (d *Dashboard) SetSystemHealth(version string, serenaStatus string, activeAgents, totalAgents int) {
	d.systemHealth = SystemHealth{
		MCFVersion:   version,
		Uptime:       time.Since(time.Now().Add(-2 * time.Hour)), // TODO: Get real uptime
		MemoryUsage:  0,                                          // TODO: Get real system metrics
		CPUUsage:     0,
		DiskUsage:    0,
		ActiveAgents: activeAgents,
		TotalAgents:  totalAgents,
		SerenaStatus: serenaStatus,
		ClaudeStatus: "active", // TODO: Get real Claude status
	}
}

// SetAgentStatuses updates agent statuses with real data
func (d *Dashboard) SetAgentStatuses(agentStatuses []AgentStatus) {
	d.agentStatuses = agentStatuses
}

// AddRecentActivity adds a new activity entry
func (d *Dashboard) AddRecentActivity(activityType, message, details string) {
	activity := RecentActivity{
		Timestamp: time.Now(),
		Type:      activityType,
		Message:   message,
		Details:   details,
	}

	// Add to beginning and limit to 10 entries
	d.recentActivity = append([]RecentActivity{activity}, d.recentActivity...)
	if len(d.recentActivity) > 10 {
		d.recentActivity = d.recentActivity[:10]
	}
}

// SetCommandHistory updates command history with real commands
func (d *Dashboard) SetCommandHistory(commands []string) {
	d.commandHistory = commands
	// Limit to last 5 commands
	if len(d.commandHistory) > 5 {
		d.commandHistory = d.commandHistory[:5]
	}
}

func (d *Dashboard) SetSelectedQuickAction(idx int) {
	if idx >= 0 && idx < len(QuickActions) {
		d.selectedQuickAction = idx
	}
}

func (d *Dashboard) GetSelectedQuickAction() int {
	return d.selectedQuickAction
}

func (d *Dashboard) ToggleHelp() {
	d.showHelp = !d.showHelp
}

// Render dashboard components
func (d *Dashboard) Render(width, height int) string {
	if d.showHelp {
		return d.renderHelp(width, height)
	}

	// Calculate layout dimensions
	leftPanelWidth := width * 2 / 3
	rightPanelWidth := width - leftPanelWidth - 2
	topPanelHeight := height / 2
	bottomPanelHeight := height - topPanelHeight - 4

	// Render main content areas
	systemPanel := d.renderSystemHealth(leftPanelWidth/2-1, topPanelHeight-2)
	agentsPanel := d.renderAgentStatus(leftPanelWidth/2-1, topPanelHeight-2)
	activityPanel := d.renderRecentActivity(leftPanelWidth-2, bottomPanelHeight-2)
	quickActionsPanel := d.renderQuickActions(rightPanelWidth-2, height-4)

	// Layout panels
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, systemPanel, agentsPanel)
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, topRow, activityPanel)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, quickActionsPanel)

	return mainView
}

func (d *Dashboard) renderSystemHealth(width, height int) string {
	content := ""

	// MCF Version and uptime
	content += d.theme.Body.Render(fmt.Sprintf("MCF %s", d.systemHealth.MCFVersion)) + "\n"
	content += d.theme.Muted.Render(fmt.Sprintf("Uptime: %v", d.systemHealth.Uptime.Truncate(time.Second))) + "\n\n"

	// System metrics with progress bars
	content += d.theme.Subtitle.Render("System Metrics") + "\n"

	memBar := RenderProgressBar(d.systemHealth.MemoryUsage/100, width-15, d.theme)
	content += fmt.Sprintf("Memory: %s %.1f%%\n", memBar, d.systemHealth.MemoryUsage)

	cpuBar := RenderProgressBar(d.systemHealth.CPUUsage/100, width-15, d.theme)
	content += fmt.Sprintf("CPU:    %s %.1f%%\n", cpuBar, d.systemHealth.CPUUsage)

	diskBar := RenderProgressBar(d.systemHealth.DiskUsage/100, width-15, d.theme)
	content += fmt.Sprintf("Disk:   %s %.1f%%\n", diskBar, d.systemHealth.DiskUsage)

	content += "\n"

	// Integration status
	content += d.theme.Subtitle.Render("Integrations") + "\n"
	content += "Serena:  " + RenderStatusIndicator(d.systemHealth.SerenaStatus, d.theme) + "\n"
	content += "Claude:  " + RenderStatusIndicator(d.systemHealth.ClaudeStatus, d.theme) + "\n"

	return RenderBox(content, "System Health", width, height, d.theme)
}

func (d *Dashboard) renderAgentStatus(width, height int) string {
	content := ""

	content += d.theme.Subtitle.Render(fmt.Sprintf("Agents (%d/%d active)",
		d.systemHealth.ActiveAgents, d.systemHealth.TotalAgents)) + "\n\n"

	for _, agent := range d.agentStatuses {
		// Agent name and status
		statusIndicator := RenderStatusIndicator(agent.Status, d.theme)
		content += fmt.Sprintf("%-18s %s\n", agent.Name, statusIndicator)

		// Task progress if active
		if agent.TasksActive > 0 {
			progress := float64(agent.TasksActive) / float64(agent.TasksTotal)
			progressBar := RenderProgressBar(progress, width-25, d.theme)
			content += fmt.Sprintf("  Tasks: %s %d/%d\n", progressBar, agent.TasksActive, agent.TasksTotal)
		}

		// Last seen
		if agent.Status != "active" {
			lastSeen := time.Since(agent.LastSeen).Truncate(time.Second)
			content += d.theme.Muted.Render(fmt.Sprintf("  Last seen: %v ago\n", lastSeen))
		}

		content += "\n"
	}

	return RenderBox(content, "Agent Status", width, height, d.theme)
}

func (d *Dashboard) renderRecentActivity(width, height int) string {
	content := ""

	content += d.theme.Subtitle.Render("Recent Activity") + "\n\n"

	maxItems := (height - 6) // Account for title and padding
	if maxItems < 1 {
		maxItems = 1 // Ensure at least 1 item can be shown
	}

	items := d.recentActivity
	if len(items) > maxItems {
		items = items[:maxItems]
	}

	for _, activity := range items {
		timestamp := activity.Timestamp.Format("15:04:05")

		var typeStyle lipgloss.Style
		var icon string
		switch activity.Type {
		case "command":
			typeStyle = d.theme.Info
			icon = "⚡"
		case "agent":
			typeStyle = d.theme.Success
			icon = "🤖"
		case "error":
			typeStyle = d.theme.Error
			icon = "❌"
		default:
			typeStyle = d.theme.Muted
			icon = "ℹ"
		}

		content += d.theme.Muted.Render(fmt.Sprintf("[%s] ", timestamp))
		content += typeStyle.Render(icon + " " + activity.Message)

		if activity.Details != "" && width > 60 {
			content += "\n" + d.theme.Muted.Render("    "+activity.Details)
		}
		content += "\n"
	}

	return RenderBox(content, "Recent Activity", width, height, d.theme)
}

func (d *Dashboard) renderQuickActions(width, height int) string {
	content := d.theme.Subtitle.Render("Quick Actions") + "\n\n"

	for i, action := range QuickActions {
		var style lipgloss.Style
		if i == d.selectedQuickAction {
			style = d.theme.ListItemActive
		} else {
			style = d.theme.ListItem
		}

		line := fmt.Sprintf("%s  %s", action.Key, action.Label)

		if i == d.selectedQuickAction {
			content += style.Render("► "+line) + "\n"
			content += d.theme.Muted.Render("  "+action.Description) + "\n"
			content += d.theme.Info.Render("  "+action.Command) + "\n"
		} else {
			content += style.Render("  "+line) + "\n"
		}

		if i < len(QuickActions)-1 {
			content += "\n"
		}
	}

	// Command history section
	content += "\n" + d.theme.Subtitle.Render("Recent Commands") + "\n\n"

	historyCount := min(5, len(d.commandHistory))
	for i := 0; i < historyCount; i++ {
		cmd := d.commandHistory[i]
		if len(cmd) > width-6 {
			cmd = cmd[:width-9] + "..."
		}
		content += d.theme.Muted.Render(fmt.Sprintf("  %s\n", cmd))
	}

	// Footer with refresh time
	content += "\n" + d.theme.Muted.Render(
		fmt.Sprintf("Last refresh: %s", d.lastRefresh.Format("15:04:05")))

	return RenderBox(content, "Actions & History", width, height, d.theme)
}

func (d *Dashboard) renderHelp(width, height int) string {
	help := `
MCF TUI Dashboard Help

NAVIGATION:
  Tab / Shift+Tab  Switch between views
  Esc              Return to dashboard
  :                Open command bar
  ?                Toggle this help
  q / Ctrl+C       Quit application

DASHBOARD:
  j/k or ↑/↓       Navigate quick actions
  Enter            Execute selected action
  r                Refresh system status
  1-6              Quick action shortcuts

QUICK ACTIONS:
  1  Agent Status   View all agent statuses
  2  Start Serena   Start Serena integration
  3  Run Tests      Execute test suite  
  4  Deploy         Deploy application
  5  View Logs      View recent logs
  6  Health Check   System health check

The dashboard shows:
• System health metrics (CPU, memory, disk usage)
• Agent status and active tasks
• Recent activity and command history
• Integration status (Serena, Claude)
• Quick action shortcuts

Press ? again to return to dashboard.
`

	return d.theme.Panel.
		Width(width - 4).
		Height(height - 2).
		Render(strings.TrimSpace(help))
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
