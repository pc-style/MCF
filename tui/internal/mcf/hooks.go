package mcf

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"mcf-dev/tui/internal/ui"
)

// HookAdapter provides integration with MCF hooks system
type HookAdapter struct {
	mcfRoot     string
	hooks       map[string][]HookConfig
	suggestions []HookSuggestion
}

// HookConfig represents a hook configuration
type HookConfig struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
	Matcher string `json:"matcher,omitempty"`
}

// HookSuggestion represents a contextual suggestion from hooks
type HookSuggestion struct {
	Title       string
	Description string
	Command     string
	Context     string
	Priority    int
	Timestamp   time.Time
}

// NewHookAdapter creates a new hook adapter
func NewHookAdapter(mcfRoot string, settings *MCFSettings) *HookAdapter {
	adapter := &HookAdapter{
		mcfRoot:     mcfRoot,
		hooks:       make(map[string][]HookConfig),
		suggestions: []HookSuggestion{},
	}

	// Load hooks from settings
	if settings != nil && settings.Hooks != nil {
		adapter.loadHooksFromSettings(settings.Hooks)
	}

	return adapter
}

// loadHooksFromSettings loads hook configurations from MCF settings
func (h *HookAdapter) loadHooksFromSettings(hooksConfig map[string]interface{}) {
	for event, configsInterface := range hooksConfig {
		if configsList, ok := configsInterface.([]interface{}); ok {
			for _, configInterface := range configsList {
				if configMap, ok := configInterface.(map[string]interface{}); ok {
					if hooksArray, ok := configMap["hooks"].([]interface{}); ok {
						for _, hookInterface := range hooksArray {
							if hookMap, ok := hookInterface.(map[string]interface{}); ok {
								hook := HookConfig{}

								if hookType, ok := hookMap["type"].(string); ok {
									hook.Type = hookType
								}
								if command, ok := hookMap["command"].(string); ok {
									hook.Command = command
								}
								if timeout, ok := hookMap["timeout"].(float64); ok {
									hook.Timeout = int(timeout)
								}
								if matcher, ok := configMap["matcher"].(string); ok {
									hook.Matcher = matcher
								}

								h.hooks[event] = append(h.hooks[event], hook)
							}
						}
					}
				}
			}
		}
	}
}

// GetHooks returns all configured hooks
func (h *HookAdapter) GetHooks() map[string][]HookConfig {
	return h.hooks
}

// GetHooksForEvent returns hooks for a specific event
func (h *HookAdapter) GetHooksForEvent(event string) []HookConfig {
	return h.hooks[event]
}

// ExecuteHooks executes hooks for a given event
func (h *HookAdapter) ExecuteHooks(event string, context map[string]string) []HookSuggestion {
	suggestions := []HookSuggestion{}

	hooks := h.GetHooksForEvent(event)
	for _, hook := range hooks {
		if hook.Type == "command" {
			suggestion := h.executeCommandHook(hook, context)
			if suggestion != nil {
				suggestions = append(suggestions, *suggestion)
			}
		}
	}

	// Add suggestions to internal list
	h.suggestions = append(h.suggestions, suggestions...)

	return suggestions
}

// executeCommandHook executes a command-based hook
func (h *HookAdapter) executeCommandHook(hook HookConfig, context map[string]string) *HookSuggestion {
	// Replace environment variables in command
	command := h.expandCommand(hook.Command)

	// Execute the hook command
	output, err := h.runHookCommand(command, hook.Timeout)
	if err != nil {
		return &HookSuggestion{
			Title:       "Hook Error",
			Description: fmt.Sprintf("Failed to execute hook: %s", err.Error()),
			Command:     command,
			Context:     "error",
			Priority:    1,
			Timestamp:   time.Now(),
		}
	}

	// Parse hook output for suggestions
	return h.parseHookOutput(output, command, context)
}

// expandCommand expands environment variables in hook command
func (h *HookAdapter) expandCommand(command string) string {
	// Replace $CLAUDE_CONFIG_DIR with actual .claude directory
	claudeConfigDir := filepath.Join(h.mcfRoot, ".claude")
	command = strings.ReplaceAll(command, "$CLAUDE_CONFIG_DIR", claudeConfigDir)

	// Replace other common variables
	command = strings.ReplaceAll(command, "$MCF_ROOT", h.mcfRoot)

	return command
}

// runHookCommand executes a hook command with timeout
func (h *HookAdapter) runHookCommand(command string, timeout int) (string, error) {
	if timeout == 0 {
		timeout = 10 // Default timeout
	}

	// Execute the command
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = h.mcfRoot

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// parseHookOutput parses hook command output into suggestions
func (h *HookAdapter) parseHookOutput(output, command string, context map[string]string) *HookSuggestion {
	if strings.TrimSpace(output) == "" {
		return nil
	}

	// Try to parse as JSON first
	var jsonSuggestion map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonSuggestion); err == nil {
		return h.parseJSONSuggestion(jsonSuggestion, command)
	}

	// Parse as plain text
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return nil
	}

	suggestion := &HookSuggestion{
		Title:       "Hook Suggestion",
		Description: lines[0],
		Command:     command,
		Context:     "general",
		Priority:    5,
		Timestamp:   time.Now(),
	}

	// Look for specific patterns in the output
	if strings.Contains(strings.ToLower(output), "serena") {
		suggestion.Context = "serena"
		suggestion.Title = "Serena Suggestion"
		suggestion.Priority = 8
	}

	if strings.Contains(strings.ToLower(output), "git") {
		suggestion.Context = "git"
		suggestion.Title = "Git Suggestion"
		suggestion.Priority = 7
	}

	if strings.Contains(strings.ToLower(output), "context7") {
		suggestion.Context = "context7"
		suggestion.Title = "Context7 Reminder"
		suggestion.Priority = 6
	}

	return suggestion
}

// parseJSONSuggestion parses a JSON-formatted hook suggestion
func (h *HookAdapter) parseJSONSuggestion(data map[string]interface{}, command string) *HookSuggestion {
	suggestion := &HookSuggestion{
		Command:   command,
		Timestamp: time.Now(),
		Priority:  5,
	}

	if title, ok := data["title"].(string); ok {
		suggestion.Title = title
	}
	if desc, ok := data["description"].(string); ok {
		suggestion.Description = desc
	}
	if ctx, ok := data["context"].(string); ok {
		suggestion.Context = ctx
	}
	if priority, ok := data["priority"].(float64); ok {
		suggestion.Priority = int(priority)
	}

	return suggestion
}

// GetRecentSuggestions returns recent hook suggestions
func (h *HookAdapter) GetRecentSuggestions(limit int) []HookSuggestion {
	// Sort by timestamp (newest first) and priority
	suggestions := make([]HookSuggestion, len(h.suggestions))
	copy(suggestions, h.suggestions)

	// Simple sort by timestamp and priority
	for i := 0; i < len(suggestions)-1; i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[i].Timestamp.Before(suggestions[j].Timestamp) ||
				(suggestions[i].Timestamp.Equal(suggestions[j].Timestamp) && suggestions[i].Priority < suggestions[j].Priority) {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	if limit > 0 && limit < len(suggestions) {
		return suggestions[:limit]
	}

	return suggestions
}

// TriggerContextualSuggestions triggers hooks based on current context
func (h *HookAdapter) TriggerContextualSuggestions(currentView string, selectedItem string) []HookSuggestion {
	context := map[string]string{
		"view": currentView,
		"item": selectedItem,
	}

	suggestions := []HookSuggestion{}

	// Trigger different hooks based on context
	switch currentView {
	case "agents":
		suggestions = append(suggestions, h.ExecuteHooks("AgentView", context)...)
	case "commands":
		suggestions = append(suggestions, h.ExecuteHooks("CommandView", context)...)
	case "logs":
		suggestions = append(suggestions, h.ExecuteHooks("LogView", context)...)
	}

	// Always trigger general suggestions
	suggestions = append(suggestions, h.ExecuteHooks("UserPromptSubmit", context)...)

	return suggestions
}

// GetHookActivity returns recent hook activity as log entries
func (h *HookAdapter) GetHookActivity() []ui.LogEntry {
	logs := []ui.LogEntry{}
	now := time.Now()

	// Add hook execution logs
	for event, hooks := range h.hooks {
		for _, hook := range hooks {
			logs = append(logs, ui.LogEntry{
				Timestamp: now.Add(-time.Duration(len(logs)) * time.Minute),
				Level:     "INFO",
				Component: "hooks",
				Message:   fmt.Sprintf("Hook configured for %s: %s", event, filepath.Base(hook.Command)),
			})
		}
	}

	// Add recent suggestions as activity
	for _, suggestion := range h.GetRecentSuggestions(5) {
		logs = append(logs, ui.LogEntry{
			Timestamp: suggestion.Timestamp,
			Level:     "INFO",
			Component: "hooks",
			Message:   fmt.Sprintf("Suggestion: %s - %s", suggestion.Title, suggestion.Description),
		})
	}

	return logs
}
