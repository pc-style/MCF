# MCF Interactive Terminal Interface

A sophisticated BubbleTea-based terminal user interface for the Multi Component Framework (MCF), featuring complete project automation, configuration management, and Claude Code integration.

## Features

### **Core Interface Components**

- **Interactive Installation Wizard** - 11-step guided setup with progress tracking
- **Live Configuration Editor** - Schema-driven forms with real-time validation
- **MCF Operation Runner** - Execute agents, commands, and templates with visual feedback
- **Template Browser** - Browse, preview, and install project templates with parameter input
- **Multi-Modal Navigation** - Main menu system with escape-key navigation

### **Advanced Functionality**

- **Global State Management** - Cross-component communication and coordination
- **Message Bus Architecture** - Event-driven updates and notifications
- **Error Handling** - Comprehensive error recovery with user-friendly messages
- **Performance Optimization** - Efficient rendering with caching and differential updates
- **Notification System** - Toast-style notifications with visual indicators

## Getting Started

### **Build and Launch**

```bash
# Build the main TUI application
go build -o mcf

# Launch interactive interface
./mcf
```

### **Direct Mode Access**

```bash
# Direct to installation wizard
./mcf install

# Direct to configuration editor
./mcf config

# Direct to MCF runner (requires setup)
./mcf run

# Show help
./mcf help
```

## Interface Navigation

### **Main Menu Navigation**

- `↑/k` - Move cursor up
- `↓/j` - Move cursor down
- `Enter/Space` - Select option
- `q/Ctrl+C` - Quit application

### **Component-Specific Controls**

#### **Installation Wizard**

- `Enter` - Proceed to next step
- `Ctrl+C` - Cancel installation

#### **Configuration Editor**

- `Tab/↑/↓` - Navigate between fields
- `Enter` - Edit field value
- `Esc` - Return to field selection
- `Ctrl+S` - Save configuration
- `Ctrl+R` - Reset to defaults

#### **MCF Runner**

- `↑/↓` - Select operation (agent/command/template)
- `Enter` - Choose operation
- `Tab` - Navigate parameter fields
- `Ctrl+R` - Execute operation
- `Esc` - Cancel/return to selection

#### **Template Browser**

- `↑/↓` - Browse templates
- `Enter` - View template details
- `i` - Install template
- `p` - Preview template files
- `Tab` - Navigate parameters

## Architecture Overview

### **Component Structure**

```
cli/
├── main.go                   # Main application and routing
├── installer.go              # 11-step installation wizard
├── configurator.go           # Live configuration editor
├── mcf_runner.go             # MCF operation runner
├── template_browser.go       # Template browser and installer
├── config.go                 # Configuration management
├── global_state.go           # Global state management
├── message_bus.go            # Cross-component messaging
└── notifications.go          # Notification system
```

### **Key Design Patterns**

#### **Bubble Tea Model-View-Update**

```go
type Model struct {
    // Component state
}

func (m Model) Init() tea.Cmd { /* ... */ }
func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) { /* ... */ }
func (m Model) View() string { /* ... */ }
```

#### **Message Bus Communication**

```go
// Cross-component messaging
GlobalMessageBus.Publish(ConfigSavedMessage{...})
GlobalMessageBus.Subscribe(MsgConfigChanged, handler)
```

#### **State Management**

```go
// Centralized state
GlobalState.SetMode(ModeConfigurator)
notifications := GlobalState.GetNotifications()
```

## Dependencies

### **Core Framework**

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - Terminal app framework (v1.3.6)
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components library (v0.18.0)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions and layout (v1.1.0)

### **Additional Libraries**

- `gopkg.in/yaml.v3` - YAML configuration parsing
- `encoding/json` - JSON configuration handling
- `os/exec` - System command execution
- `path/filepath` - Cross-platform path handling

## Integration Points

### **Claude Code Integration**

- **Settings Management**: Reads/writes `.claude/settings.json`
- **Agent Discovery**: Auto-discovers agents from `.claude/agents/`
- **Command Integration**: Loads commands from `.claude/commands/`
- **Hook System**: Configures and manages `.claude/hooks/`

### **MCF Framework Integration**

- **Template System**: Manages `.claude/templates/`
- **Configuration Schema**: Uses `config-schema.yaml` for form generation
- **Project Structure**: Creates and manages MCF directory layout
- **Installation Automation**: 11-step setup process

### **External Tools**

- **Git Integration**: Repository initialization and configuration
- **Serena MCP**: Semantic code analysis setup
- **Shell Integration**: PATH and alias configuration
- **File Permissions**: Executable permissions for scripts and hooks

## Performance Characteristics

- **Memory Efficient**: Differential rendering and component caching
- **Responsive UI**: Async operations with progress indicators
- **Scalable Architecture**: Message bus prevents tight coupling
- **Error Resilient**: Comprehensive error handling and recovery
