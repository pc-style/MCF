package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Centralized message types for component communication
type MessageType string

const (
	// Application lifecycle messages
	MsgAppInit        MessageType = "app_init"
	MsgAppShutdown    MessageType = "app_shutdown"
	MsgModeTransition MessageType = "mode_transition"

	// Configuration messages
	MsgConfigLoaded  MessageType = "config_loaded"
	MsgConfigSaved   MessageType = "config_saved"
	MsgConfigChanged MessageType = "config_changed"
	MsgConfigError   MessageType = "config_error"

	// Installation messages
	MsgInstallStart    MessageType = "install_start"
	MsgInstallProgress MessageType = "install_progress"
	MsgInstallComplete MessageType = "install_complete"
	MsgInstallError    MessageType = "install_error"

	// UI interaction messages
	MsgUIError   MessageType = "ui_error"
	MsgUISuccess MessageType = "ui_success"
	MsgUIWarning MessageType = "ui_warning"
	MsgUIRefresh MessageType = "ui_refresh"
)

// Base message interface for type safety
type MCFMessage interface {
	tea.Msg
	Type() MessageType
	Timestamp() time.Time
}

// Base message struct
type BaseMessage struct {
	MsgType   MessageType `json:"type"`
	CreatedAt time.Time   `json:"created_at"`
}

func (m BaseMessage) Type() MessageType {
	return m.MsgType
}

func (m BaseMessage) Timestamp() time.Time {
	return m.CreatedAt
}

// Application Messages
type AppInitMessage struct {
	BaseMessage
	Version     string                 `json:"version"`
	Environment map[string]interface{} `json:"environment"`
}

func NewAppInitMessage(version string, env map[string]interface{}) AppInitMessage {
	return AppInitMessage{
		BaseMessage: BaseMessage{MsgType: MsgAppInit, CreatedAt: time.Now()},
		Version:     version,
		Environment: env,
	}
}

type ModeTransitionMessage struct {
	BaseMessage
	FromMode ApplicationMode `json:"from_mode"`
	ToMode   ApplicationMode `json:"to_mode"`
	Context  interface{}     `json:"context,omitempty"`
}

func NewModeTransitionMessage(from, to ApplicationMode, context interface{}) ModeTransitionMessage {
	return ModeTransitionMessage{
		BaseMessage: BaseMessage{MsgType: MsgModeTransition, CreatedAt: time.Now()},
		FromMode:    from,
		ToMode:      to,
		Context:     context,
	}
}

// Configuration Messages
type ConfigurationMessage struct {
	BaseMessage
	Success bool                   `json:"success"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Section string                 `json:"section,omitempty"`
}

func NewConfigLoadedMessage(success bool, data map[string]interface{}, err string) ConfigurationMessage {
	return ConfigurationMessage{
		BaseMessage: BaseMessage{MsgType: MsgConfigLoaded, CreatedAt: time.Now()},
		Success:     success,
		Error:       err,
		Data:        data,
	}
}

func NewConfigSavedMessage(success bool, err string, section string) ConfigurationMessage {
	return ConfigurationMessage{
		BaseMessage: BaseMessage{MsgType: MsgConfigSaved, CreatedAt: time.Now()},
		Success:     success,
		Error:       err,
		Section:     section,
	}
}

// Installation Messages
type InstallationMessage struct {
	BaseMessage
	Success     bool        `json:"success"`
	Error       string      `json:"error,omitempty"`
	Progress    float64     `json:"progress"`
	CurrentStep string      `json:"current_step,omitempty"`
	StepIndex   int         `json:"step_index"`
	Context     interface{} `json:"context,omitempty"`
}

func NewInstallProgressMessage(progress float64, step string, stepIndex int) InstallationMessage {
	return InstallationMessage{
		BaseMessage: BaseMessage{MsgType: MsgInstallProgress, CreatedAt: time.Now()},
		Progress:    progress,
		CurrentStep: step,
		StepIndex:   stepIndex,
	}
}

func NewInstallCompleteMessage(success bool, err string) InstallationMessage {
	return InstallationMessage{
		BaseMessage: BaseMessage{MsgType: MsgInstallComplete, CreatedAt: time.Now()},
		Success:     success,
		Error:       err,
		Progress:    1.0,
	}
}

// UI Messages
type UIMessage struct {
	BaseMessage
	Level   string      `json:"level"` // error, success, warning, info
	Title   string      `json:"title"`
	Message string      `json:"message"`
	Context interface{} `json:"context,omitempty"`
}

func NewUIErrorMessage(title, message string, context interface{}) UIMessage {
	return UIMessage{
		BaseMessage: BaseMessage{MsgType: MsgUIError, CreatedAt: time.Now()},
		Level:       "error",
		Title:       title,
		Message:     message,
		Context:     context,
	}
}

func NewUISuccessMessage(title, message string) UIMessage {
	return UIMessage{
		BaseMessage: BaseMessage{MsgType: MsgUISuccess, CreatedAt: time.Now()},
		Level:       "success",
		Title:       title,
		Message:     message,
	}
}

func NewUIWarningMessage(title, message string) UIMessage {
	return UIMessage{
		BaseMessage: BaseMessage{MsgType: MsgUIWarning, CreatedAt: time.Now()},
		Level:       "warning",
		Title:       title,
		Message:     message,
	}
}

// Message Bus for cross-component communication
type MessageBus struct {
	subscribers map[MessageType][]chan MCFMessage
}

func NewMessageBus() *MessageBus {
	return &MessageBus{
		subscribers: make(map[MessageType][]chan MCFMessage),
	}
}

func (mb *MessageBus) Subscribe(msgType MessageType, ch chan MCFMessage) {
	mb.subscribers[msgType] = append(mb.subscribers[msgType], ch)
}

func (mb *MessageBus) Publish(msg MCFMessage) {
	if subscribers, exists := mb.subscribers[msg.Type()]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- msg:
			default:
				// Non-blocking send
			}
		}
	}
}

// Global message bus instance
var GlobalMessageBus = NewMessageBus()
