package ui

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutils "mcf-dev/tui/internal/testing"
)

func TestInteractiveList_Creation(t *testing.T) {
	t.Run("should create list with default values", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		assert.NotNil(t, list, "List should be created")
		assert.Equal(t, "Test List", list.title)
		assert.Equal(t, 20, list.height)
		assert.Equal(t, 0, list.selected)
		assert.False(t, list.focused)
		assert.Empty(t, list.items)
	})
}

func TestInteractiveList_SetItems(t *testing.T) {
	t.Run("should set items correctly", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		items := []ListItem{
			{Title: "Item 1", Description: "First item", Status: "active"},
			{Title: "Item 2", Description: "Second item", Status: "idle"},
			{Title: "Item 3", Description: "Third item", Status: "error"},
		}

		list.SetItems(items)

		assert.Len(t, list.items, 3, "Should have 3 items")
		assert.Equal(t, "Item 1", list.items[0].Title)
		assert.Equal(t, "active", list.items[0].Status)
	})

	t.Run("should adjust selection when items change", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		// Set initial items and select last one
		items := []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"}, {Title: "Item 3"},
		}
		list.SetItems(items)
		list.selected = 2

		// Reduce items to 2
		items = []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"},
		}
		list.SetItems(items)

		assert.Equal(t, 1, list.selected, "Should adjust selection to valid index")
	})

	t.Run("should handle empty items", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		list.SetItems([]ListItem{})

		assert.Empty(t, list.items)
		assert.Equal(t, 0, list.selected)
	})
}

func TestInteractiveList_Selection(t *testing.T) {
	t.Run("should get selected index", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)
		list.selected = 5

		assert.Equal(t, 5, list.GetSelected())
	})

	t.Run("should get selected item", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		items := []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"}, {Title: "Item 3"},
		}
		list.SetItems(items)
		list.selected = 1

		selectedItem := list.GetSelectedItem()
		require.NotNil(t, selectedItem)
		assert.Equal(t, "Item 2", selectedItem.Title)
	})

	t.Run("should return nil for invalid selection", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)
		list.selected = 5 // No items, so this is invalid

		selectedItem := list.GetSelectedItem()
		assert.Nil(t, selectedItem)
	})
}

func TestInteractiveList_Focus(t *testing.T) {
	t.Run("should set focus correctly", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		list.SetFocus(true)
		assert.True(t, list.focused)

		list.SetFocus(false)
		assert.False(t, list.focused)
	})
}

func TestInteractiveList_Update(t *testing.T) {
	t.Run("should handle navigation when focused", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		items := []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"}, {Title: "Item 3"},
		}
		list.SetItems(items)
		list.SetFocus(true)

		// Test down navigation
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedList, cmd := list.Update(downMsg)
		assert.Equal(t, 1, updatedList.selected, "Should move down")
		assert.Nil(t, cmd)

		// Test up navigation
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		updatedList, cmd = updatedList.Update(upMsg)
		assert.Equal(t, 0, updatedList.selected, "Should move up")

		// Test 'j' key (vim-style down)
		jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
		updatedList, cmd = updatedList.Update(jMsg)
		assert.Equal(t, 1, updatedList.selected, "Should move down with j")

		// Test 'k' key (vim-style up)
		kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
		updatedList, cmd = updatedList.Update(kMsg)
		assert.Equal(t, 0, updatedList.selected, "Should move up with k")
	})

	t.Run("should handle boundary navigation", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		items := []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"}, {Title: "Item 3"},
		}
		list.SetItems(items)
		list.SetFocus(true)

		// Try to move up from first item
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		updatedList, _ := list.Update(upMsg)
		assert.Equal(t, 0, updatedList.selected, "Should stay at first item")

		// Move to last item and try to move down
		list.selected = 2
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedList, _ = list.Update(downMsg)
		assert.Equal(t, 2, updatedList.selected, "Should stay at last item")
	})

	t.Run("should handle home and end keys", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		items := []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"}, {Title: "Item 3"},
		}
		list.SetItems(items)
		list.SetFocus(true)
		list.selected = 1

		// Test home key
		homeMsg := tea.KeyMsg{Type: tea.KeyHome}
		updatedList, _ := list.Update(homeMsg)
		assert.Equal(t, 0, updatedList.selected, "Should go to first item")

		// Test end key
		endMsg := tea.KeyMsg{Type: tea.KeyEnd}
		updatedList, _ = updatedList.Update(endMsg)
		assert.Equal(t, 2, updatedList.selected, "Should go to last item")

		// Test 'g' key (vim-style home)
		list.selected = 1
		gMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("g")}
		updatedList, _ = list.Update(gMsg)
		assert.Equal(t, 0, updatedList.selected, "Should go to first item with g")

		// Test 'G' key (vim-style end)
		GMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("G")}
		updatedList, _ = updatedList.Update(GMsg)
		assert.Equal(t, 2, updatedList.selected, "Should go to last item with G")
	})

	t.Run("should not respond when not focused", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 20)

		items := []ListItem{
			{Title: "Item 1"}, {Title: "Item 2"}, {Title: "Item 3"},
		}
		list.SetItems(items)
		list.SetFocus(false) // Not focused

		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedList, cmd := list.Update(downMsg)
		assert.Equal(t, 0, updatedList.selected, "Should not change selection when not focused")
		assert.Nil(t, cmd)
	})
}

func TestInteractiveList_Render(t *testing.T) {
	t.Run("should render empty list", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Empty List", 10)

		rendered := list.Render(80)
		assert.Contains(t, rendered, "Empty List", "Should show title")
		assert.Contains(t, rendered, "No items available", "Should show empty message")
	})

	t.Run("should render list with items", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Test List", 10)

		items := []ListItem{
			{Title: "Item 1", Description: "First item", Status: "active"},
			{Title: "Item 2", Description: "Second item", Status: "idle"},
		}
		list.SetItems(items)
		list.SetFocus(true)

		rendered := list.Render(80)
		assert.Contains(t, rendered, "Test List", "Should show title")
		assert.Contains(t, rendered, "Item 1", "Should show first item")
		assert.Contains(t, rendered, "Item 2", "Should show second item")
		assert.Contains(t, rendered, "►", "Should show selection cursor")
	})

	t.Run("should handle scrolling for long lists", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Long List", 5) // Small height

		// Create many items
		items := []ListItem{}
		for i := 1; i <= 20; i++ {
			items = append(items, ListItem{Title: fmt.Sprintf("Item %d", i)})
		}
		list.SetItems(items)
		list.selected = 10 // Select item beyond visible area

		rendered := list.Render(80)
		assert.Contains(t, rendered, "of 20", "Should show scroll indicator")
	})

	t.Run("should show status indicators", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Status List", 10)

		items := []ListItem{
			{Title: "Active Item", Status: "active"},
			{Title: "Error Item", Status: "error"},
			{Title: "Idle Item", Status: "idle"},
		}
		list.SetItems(items)

		rendered := list.Render(80)
		// The exact status indicator symbols would depend on the RenderStatusIndicator function
		assert.Contains(t, rendered, "Active Item")
		assert.Contains(t, rendered, "Error Item")
		assert.Contains(t, rendered, "Idle Item")
	})
}

func TestCommandInput_Creation(t *testing.T) {
	t.Run("should create command input with defaults", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		assert.NotNil(t, commandInput, "Command input should be created")
		assert.Empty(t, commandInput.suggestions)
		assert.Equal(t, 0, commandInput.selected)
		assert.True(t, commandInput.focused) // Should be focused by default
		assert.Empty(t, commandInput.history)
		assert.Equal(t, -1, commandInput.historyIdx)
	})
}

func TestCommandInput_Focus(t *testing.T) {
	t.Run("should manage focus correctly", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		commandInput.SetFocus(true)
		assert.True(t, commandInput.focused)

		commandInput.SetFocus(false)
		assert.False(t, commandInput.focused)
	})
}

func TestCommandInput_Suggestions(t *testing.T) {
	t.Run("should set suggestions", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		suggestions := []string{"command 1", "command 2", "command 3"}
		commandInput.SetSuggestions(suggestions)

		assert.Equal(t, suggestions, commandInput.suggestions)
		assert.Equal(t, 0, commandInput.selected)
	})

	t.Run("should update suggestions based on input", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		// This tests the internal updateSuggestions method
		// We'll simulate typing to trigger suggestions
		commandInput.input.SetValue("mcf")
		commandInput.updateSuggestions()

		assert.NotEmpty(t, commandInput.suggestions, "Should have suggestions for 'mcf'")

		// Clear input should clear suggestions
		commandInput.input.SetValue("")
		commandInput.updateSuggestions()

		assert.Empty(t, commandInput.suggestions, "Should clear suggestions for short input")
	})
}

func TestCommandInput_History(t *testing.T) {
	t.Run("should add commands to history", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		commandInput.AddToHistory("command 1")
		commandInput.AddToHistory("command 2")

		assert.Len(t, commandInput.history, 2)
		assert.Equal(t, "command 2", commandInput.history[0]) // Most recent first
		assert.Equal(t, "command 1", commandInput.history[1])
	})

	t.Run("should not add duplicate consecutive commands", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		commandInput.AddToHistory("same command")
		commandInput.AddToHistory("same command")

		assert.Len(t, commandInput.history, 1, "Should not add duplicate consecutive commands")
	})

	t.Run("should not add empty commands", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		commandInput.AddToHistory("")

		assert.Empty(t, commandInput.history, "Should not add empty commands")
	})

	t.Run("should limit history size", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		// Add more than 50 commands (the limit)
		for i := 1; i <= 55; i++ {
			commandInput.AddToHistory(fmt.Sprintf("command %d", i))
		}

		assert.Len(t, commandInput.history, 50, "Should limit history to 50 items")
		assert.Equal(t, "command 55", commandInput.history[0], "Should keep most recent commands")
	})
}

func TestCommandInput_Value(t *testing.T) {
	t.Run("should get and clear value", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		commandInput.input.SetValue("test command")
		assert.Equal(t, "test command", commandInput.GetValue())

		commandInput.Clear()
		assert.Empty(t, commandInput.GetValue())
		assert.Empty(t, commandInput.suggestions)
		assert.Equal(t, 0, commandInput.selected)
	})
}

func TestCommandInput_Update(t *testing.T) {
	t.Run("should handle navigation in suggestions", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.SetFocus(true)

		suggestions := []string{"cmd1", "cmd2", "cmd3"}
		commandInput.SetSuggestions(suggestions)

		// Test down navigation
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedInput, _ := commandInput.Update(downMsg)
		assert.Equal(t, 1, updatedInput.selected)

		// Test up navigation
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		updatedInput, _ = updatedInput.Update(upMsg)
		assert.Equal(t, 0, updatedInput.selected)
	})

	t.Run("should handle history navigation", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.SetFocus(true)

		commandInput.AddToHistory("old command")
		commandInput.AddToHistory("new command")

		// Navigate up in history
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		updatedInput, _ := commandInput.Update(upMsg)
		assert.Equal(t, "new command", updatedInput.GetValue())

		// Navigate up again
		updatedInput, _ = updatedInput.Update(upMsg)
		assert.Equal(t, "old command", updatedInput.GetValue())

		// Navigate down
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedInput, _ = updatedInput.Update(downMsg)
		assert.Equal(t, "new command", updatedInput.GetValue())
	})

	t.Run("should handle tab completion", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.SetFocus(true)

		suggestions := []string{"mcf agents status", "mcf agents start"}
		commandInput.SetSuggestions(suggestions)
		commandInput.selected = 1

		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		updatedInput, _ := commandInput.Update(tabMsg)

		assert.Equal(t, "mcf agents start", updatedInput.GetValue())
		assert.Empty(t, updatedInput.suggestions, "Should clear suggestions after completion")
	})

	t.Run("should not respond when not focused", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.SetFocus(false)

		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedInput, cmd := commandInput.Update(downMsg)

		assert.Equal(t, commandInput, updatedInput, "Should not change when not focused")
		assert.Nil(t, cmd)
	})
}

func TestCommandInput_Render(t *testing.T) {
	t.Run("should render basic input", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.input.SetValue("test command")

		rendered := commandInput.Render(80)
		assert.NotEmpty(t, rendered, "Should render input")
	})

	t.Run("should render with suggestions", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.SetFocus(true)

		suggestions := []string{"mcf agents status", "mcf deploy"}
		commandInput.SetSuggestions(suggestions)

		rendered := commandInput.Render(80)
		assert.Contains(t, rendered, "Suggestions", "Should show suggestions header")
		assert.Contains(t, rendered, "mcf agents status", "Should show first suggestion")
		assert.Contains(t, rendered, "mcf deploy", "Should show second suggestion")
	})

	t.Run("should render with history when input is empty", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)
		commandInput.SetFocus(true)

		commandInput.AddToHistory("recent command 1")
		commandInput.AddToHistory("recent command 2")

		rendered := commandInput.Render(80)
		assert.Contains(t, rendered, "Recent Commands", "Should show recent commands header")
		assert.Contains(t, rendered, "recent command 2", "Should show most recent command")
	})

	t.Run("should handle long text truncation", func(t *testing.T) {
		theme := NewTheme()
		commandInput := NewCommandInput(theme)

		longCommand := "mcf deploy --environment production --region us-east-1 --confirm --verbose --dry-run --timeout 300"
		suggestions := []string{longCommand}
		commandInput.SetSuggestions(suggestions)

		rendered := commandInput.Render(40) // Narrow width
		assert.Contains(t, rendered, "...", "Should truncate long suggestions")
	})
}

func TestLogViewer_Creation(t *testing.T) {
	t.Run("should create log viewer with defaults", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		assert.NotNil(t, logViewer, "Log viewer should be created")
		assert.Equal(t, 20, logViewer.height)
		assert.Empty(t, logViewer.logs)
		assert.Empty(t, logViewer.filter)
		assert.False(t, logViewer.following)
		assert.Equal(t, 0, logViewer.scrollPos)
		assert.False(t, logViewer.searchMode)
	})
}

func TestLogViewer_AddLog(t *testing.T) {
	t.Run("should add log entries", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		entry := LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Component: "test",
			Message:   "Test message",
		}

		logViewer.AddLog(entry)

		assert.Len(t, logViewer.logs, 1)
		assert.Equal(t, "Test message", logViewer.logs[0].Message)
	})

	t.Run("should scroll to bottom when following", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 5) // Small height

		logViewer.SetFollowing(true)

		// Add many log entries
		for i := 1; i <= 10; i++ {
			entry := LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "test",
				Message:   fmt.Sprintf("Message %d", i),
			}
			logViewer.AddLog(entry)
		}

		// Should be scrolled to show latest entries
		assert.Greater(t, logViewer.scrollPos, 0, "Should scroll when following")
	})

	t.Run("should limit log history", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		// Add more than 1000 entries (the limit)
		for i := 1; i <= 1050; i++ {
			entry := LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "test",
				Message:   fmt.Sprintf("Message %d", i),
			}
			logViewer.AddLog(entry)
		}

		assert.LessOrEqual(t, len(logViewer.logs), 1000, "Should limit log history")
	})
}

func TestLogViewer_Following(t *testing.T) {
	t.Run("should set following mode", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		logViewer.SetFollowing(true)
		assert.True(t, logViewer.following)

		logViewer.SetFollowing(false)
		assert.False(t, logViewer.following)
	})
}

func TestLogViewer_SearchMode(t *testing.T) {
	t.Run("should set search mode", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		logViewer.SetSearchMode(true)
		assert.True(t, logViewer.searchMode)

		logViewer.SetSearchMode(false)
		assert.False(t, logViewer.searchMode)
		assert.Empty(t, logViewer.filter, "Should clear filter when exiting search")
	})
}

func TestLogViewer_Clear(t *testing.T) {
	t.Run("should clear logs", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		entry := LogEntry{Timestamp: time.Now(), Level: "INFO", Message: "Test"}
		logViewer.AddLog(entry)

		logViewer.Clear()

		assert.Empty(t, logViewer.logs)
		assert.Equal(t, 0, logViewer.scrollPos)
	})
}

func TestLogViewer_Update(t *testing.T) {
	t.Run("should handle search mode input", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)
		logViewer.SetSearchMode(true)

		// Type in search mode
		charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("error")}
		updatedViewer, _ := logViewer.Update(charMsg)

		// Press enter to confirm search
		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		updatedViewer, _ = updatedViewer.Update(enterMsg)

		assert.False(t, updatedViewer.searchMode, "Should exit search mode")
		assert.Equal(t, "error", updatedViewer.filter, "Should set filter")
	})

	t.Run("should handle navigation keys", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 5) // Small height

		// Add many entries to enable scrolling
		for i := 1; i <= 20; i++ {
			entry := LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Message:   fmt.Sprintf("Message %d", i),
			}
			logViewer.AddLog(entry)
		}

		// Test scrolling up
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		updatedViewer, _ := logViewer.Update(upMsg)
		assert.Equal(t, 0, updatedViewer.scrollPos, "Should not scroll up from top")
		assert.False(t, updatedViewer.following, "Should disable following when manually scrolling")

		// Test scrolling down
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedViewer, _ = updatedViewer.Update(downMsg)
		assert.Greater(t, updatedViewer.scrollPos, 0, "Should scroll down")
	})

	t.Run("should handle special keys", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		// Add test entries
		for i := 1; i <= 10; i++ {
			entry := LogEntry{Timestamp: time.Now(), Level: "INFO", Message: fmt.Sprintf("Message %d", i)}
			logViewer.AddLog(entry)
		}

		// Test follow toggle
		fMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}
		updatedViewer, _ := logViewer.Update(fMsg)
		assert.True(t, updatedViewer.following, "Should toggle following")

		// Test search
		slashMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
		updatedViewer, _ = updatedViewer.Update(slashMsg)
		assert.True(t, updatedViewer.searchMode, "Should enter search mode")

		// Test clear
		logViewer.SetSearchMode(false) // Exit search first
		cMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")}
		updatedViewer, _ = logViewer.Update(cMsg)
		assert.Empty(t, updatedViewer.logs, "Should clear logs")
	})
}

func TestLogViewer_Render(t *testing.T) {
	t.Run("should render empty log viewer", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		rendered := logViewer.Render(80)
		assert.Contains(t, rendered, "Logs", "Should show title")
		assert.Contains(t, rendered, "PAUSED", "Should show paused status")
	})

	t.Run("should render with log entries", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		entries := []LogEntry{
			{Timestamp: time.Now(), Level: "INFO", Component: "test", Message: "Info message"},
			{Timestamp: time.Now(), Level: "ERROR", Component: "test", Message: "Error message"},
			{Timestamp: time.Now(), Level: "WARN", Component: "test", Message: "Warning message"},
		}

		for _, entry := range entries {
			logViewer.AddLog(entry)
		}

		rendered := logViewer.Render(80)
		assert.Contains(t, rendered, "Info message", "Should show info message")
		assert.Contains(t, rendered, "Error message", "Should show error message")
		assert.Contains(t, rendered, "Warning message", "Should show warning message")
	})

	t.Run("should render search mode", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)
		logViewer.SetSearchMode(true)

		rendered := logViewer.Render(80)
		assert.Contains(t, rendered, "Search Logs", "Should show search title")
		assert.Contains(t, rendered, "Enter to search", "Should show search instructions")
	})

	t.Run("should show following status", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)
		logViewer.SetFollowing(true)

		rendered := logViewer.Render(80)
		assert.Contains(t, rendered, "FOLLOWING", "Should show following status")
	})

	t.Run("should show filter status", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)
		logViewer.filter = "error"

		rendered := logViewer.Render(80)
		assert.Contains(t, rendered, "Filter: 'error'", "Should show filter status")
	})
}

// Performance and benchmark tests
func TestUIComponents_Performance(t *testing.T) {
	t.Run("should handle rapid list updates", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Performance Test", 20)

		benchmark := testutils.NewPerformanceBenchmark("rapid_list_updates", func() error {
			items := []ListItem{{Title: "Test Item", Status: "active"}}
			list.SetItems(items)
			return nil
		}).WithIterations(1000)

		benchmark.Run(t)
	})

	t.Run("should render large lists efficiently", func(t *testing.T) {
		theme := NewTheme()
		list := NewInteractiveList(theme, "Large List", 20)

		// Create large item list
		items := []ListItem{}
		for i := 1; i <= 1000; i++ {
			items = append(items, ListItem{
				Title:       fmt.Sprintf("Item %d", i),
				Description: fmt.Sprintf("Description for item %d", i),
				Status:      "active",
			})
		}
		list.SetItems(items)

		benchmark := testutils.NewPerformanceBenchmark("large_list_render", func() error {
			_ = list.Render(80)
			return nil
		}).WithIterations(100)

		benchmark.Run(t)
	})

	t.Run("should handle rapid log additions", func(t *testing.T) {
		theme := NewTheme()
		logViewer := NewLogViewer(theme, 20)

		benchmark := testutils.NewPerformanceBenchmark("rapid_log_additions", func() error {
			entry := LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "test",
				Message:   "Performance test message",
			}
			logViewer.AddLog(entry)
			return nil
		}).WithIterations(1000)

		benchmark.Run(t)
	})
}

// Benchmark tests
func BenchmarkInteractiveListRender(b *testing.B) {
	theme := NewTheme()
	list := NewInteractiveList(theme, "Benchmark", 20)

	items := []ListItem{}
	for i := 1; i <= 100; i++ {
		items = append(items, ListItem{Title: fmt.Sprintf("Item %d", i), Status: "active"})
	}
	list.SetItems(items)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Render(80)
	}
}

func BenchmarkLogViewerAddLog(b *testing.B) {
	theme := NewTheme()
	logViewer := NewLogViewer(theme, 20)

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Component: "benchmark",
		Message:   "Benchmark message",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logViewer.AddLog(entry)
	}
}

func BenchmarkCommandInputUpdate(b *testing.B) {
	theme := NewTheme()
	commandInput := NewCommandInput(theme)
	commandInput.SetFocus(true)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		commandInput, _ = commandInput.Update(msg)
	}
}
