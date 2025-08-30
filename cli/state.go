package main

import (
	"encoding/json"
	"sync"
	"time"
)

// ApplicationState manages shared state across TUI components
type ApplicationState struct {
	mutex              sync.RWMutex
	currentMode        ApplicationMode
	previousMode       ApplicationMode
	initializationTime time.Time
	lastActivity       time.Time

	// Configuration state
	configurationLoaded bool
	configurationDirty  bool
	configurationError  error
	globalConfig        *GlobalConfig
	projectConfig       *ProjectConfig
	localConfig         *LocalConfig

	// Installation state
	installationActive   bool
	installationProgress float64
	installationStep     string
	installationError    error

	// UI state
	windowWidth   int
	windowHeight  int
	notifications []Notification

	// Component states
	mainModel         *MainModel
	installerModel    *InstallerModel
	configuratorModel *ConfiguratorModel

	// Session information
	sessionID       string
	startTime       time.Time
	projectPath     string
	workingDir      string
	userPreferences map[string]interface{}
}

// Notification represents a UI notification
type Notification struct {
	ID        string        `json:"id"`
	Type      string        `json:"type"` // success, error, warning, info
	Title     string        `json:"title"`
	Message   string        `json:"message"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
	Dismissed bool          `json:"dismissed"`
}

// StateSnapshot represents a point-in-time state for debugging/recovery
type StateSnapshot struct {
	Timestamp            time.Time              `json:"timestamp"`
	Mode                 ApplicationMode        `json:"mode"`
	ConfigurationLoaded  bool                   `json:"configuration_loaded"`
	InstallationActive   bool                   `json:"installation_active"`
	InstallationProgress float64                `json:"installation_progress"`
	ProjectPath          string                 `json:"project_path"`
	Notifications        []Notification         `json:"notifications"`
	UserPreferences      map[string]interface{} `json:"user_preferences"`
}

// Global state instance
var GlobalState = NewApplicationState()

func NewApplicationState() *ApplicationState {
	return &ApplicationState{
		currentMode:        ModeMainMenu,
		initializationTime: time.Now(),
		lastActivity:       time.Now(),
		sessionID:          generateSessionID(),
		startTime:          time.Now(),
		notifications:      make([]Notification, 0),
		userPreferences:    make(map[string]interface{}),
	}
}

func generateSessionID() string {
	// Simple session ID generation
	return time.Now().Format("20060102-150405")
}

// Mode management
func (s *ApplicationState) GetCurrentMode() ApplicationMode {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.currentMode
}

func (s *ApplicationState) SetMode(mode ApplicationMode) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.previousMode = s.currentMode
	s.currentMode = mode
	s.lastActivity = time.Now()

	// Publish mode transition message
	msg := NewModeTransitionMessage(s.previousMode, s.currentMode, nil)
	GlobalMessageBus.Publish(msg)
}

func (s *ApplicationState) GetPreviousMode() ApplicationMode {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.previousMode
}

// Configuration management
func (s *ApplicationState) IsConfigurationLoaded() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.configurationLoaded
}

func (s *ApplicationState) SetConfigurationLoaded(loaded bool, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.configurationLoaded = loaded
	s.configurationError = err
	s.lastActivity = time.Now()

	// Publish configuration loaded message
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	msg := NewConfigLoadedMessage(loaded, nil, errStr)
	GlobalMessageBus.Publish(msg)
}

func (s *ApplicationState) IsConfigurationDirty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.configurationDirty
}

func (s *ApplicationState) SetConfigurationDirty(dirty bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.configurationDirty = dirty
	s.lastActivity = time.Now()
}

func (s *ApplicationState) GetConfiguration() (*GlobalConfig, *ProjectConfig, *LocalConfig) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.globalConfig, s.projectConfig, s.localConfig
}

func (s *ApplicationState) SetConfiguration(global *GlobalConfig, project *ProjectConfig, local *LocalConfig) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.globalConfig = global
	s.projectConfig = project
	s.localConfig = local
	s.lastActivity = time.Now()
}

// Installation management
func (s *ApplicationState) IsInstallationActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.installationActive
}

func (s *ApplicationState) SetInstallationActive(active bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.installationActive = active
	s.lastActivity = time.Now()
}

func (s *ApplicationState) GetInstallationProgress() (float64, string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.installationProgress, s.installationStep
}

func (s *ApplicationState) SetInstallationProgress(progress float64, step string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.installationProgress = progress
	s.installationStep = step
	s.lastActivity = time.Now()

	// Publish progress message
	msg := NewInstallProgressMessage(progress, step, 0)
	GlobalMessageBus.Publish(msg)
}

func (s *ApplicationState) SetInstallationError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.installationError = err
	s.installationActive = false
	s.lastActivity = time.Now()

	if err != nil {
		// Publish error message
		msg := NewInstallCompleteMessage(false, err.Error())
		GlobalMessageBus.Publish(msg)
	}
}

// UI management
func (s *ApplicationState) GetWindowSize() (int, int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.windowWidth, s.windowHeight
}

func (s *ApplicationState) SetWindowSize(width, height int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.windowWidth = width
	s.windowHeight = height
}

// Notification management
func (s *ApplicationState) AddNotification(notification Notification) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if notification.ID == "" {
		notification.ID = generateSessionID() + "-" + time.Now().Format("150405")
	}
	if notification.Timestamp.IsZero() {
		notification.Timestamp = time.Now()
	}

	s.notifications = append(s.notifications, notification)
	s.lastActivity = time.Now()

	// Publish UI message
	switch notification.Type {
	case "error":
		msg := NewUIErrorMessage(notification.Title, notification.Message, nil)
		GlobalMessageBus.Publish(msg)
	case "success":
		msg := NewUISuccessMessage(notification.Title, notification.Message)
		GlobalMessageBus.Publish(msg)
	case "warning":
		msg := NewUIWarningMessage(notification.Title, notification.Message)
		GlobalMessageBus.Publish(msg)
	}
}

func (s *ApplicationState) GetNotifications() []Notification {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Filter out dismissed notifications
	active := make([]Notification, 0, len(s.notifications))
	for _, notification := range s.notifications {
		if !notification.Dismissed {
			active = append(active, notification)
		}
	}
	return active
}

func (s *ApplicationState) DismissNotification(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.notifications {
		if s.notifications[i].ID == id {
			s.notifications[i].Dismissed = true
			break
		}
	}
}

// Component model management
func (s *ApplicationState) SetMainModel(model *MainModel) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.mainModel = model
}

func (s *ApplicationState) GetMainModel() *MainModel {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.mainModel
}

func (s *ApplicationState) SetInstallerModel(model *InstallerModel) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.installerModel = model
}

func (s *ApplicationState) GetInstallerModel() *InstallerModel {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.installerModel
}

func (s *ApplicationState) SetConfiguratorModel(model *ConfiguratorModel) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.configuratorModel = model
}

func (s *ApplicationState) GetConfiguratorModel() *ConfiguratorModel {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.configuratorModel
}

// Session management
func (s *ApplicationState) GetSessionInfo() (string, time.Time, time.Duration) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.sessionID, s.startTime, time.Since(s.startTime)
}

func (s *ApplicationState) GetProjectInfo() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.projectPath
}

func (s *ApplicationState) SetProjectPath(path string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.projectPath = path
	s.lastActivity = time.Now()
}

// User preferences
func (s *ApplicationState) GetUserPreference(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.userPreferences[key]
	return value, exists
}

func (s *ApplicationState) SetUserPreference(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.userPreferences[key] = value
	s.lastActivity = time.Now()
}

// State snapshot for debugging/recovery
func (s *ApplicationState) CreateSnapshot() StateSnapshot {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return StateSnapshot{
		Timestamp:            time.Now(),
		Mode:                 s.currentMode,
		ConfigurationLoaded:  s.configurationLoaded,
		InstallationActive:   s.installationActive,
		InstallationProgress: s.installationProgress,
		ProjectPath:          s.projectPath,
		Notifications:        s.notifications,
		UserPreferences:      s.userPreferences,
	}
}

func (s *ApplicationState) ToJSON() ([]byte, error) {
	snapshot := s.CreateSnapshot()
	return json.MarshalIndent(snapshot, "", "  ")
}

// Health check
func (s *ApplicationState) IsHealthy() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Basic health checks
	if time.Since(s.lastActivity) > 10*time.Minute {
		return false // No activity for too long
	}

	if s.configurationError != nil {
		return false // Configuration error
	}

	if s.installationError != nil {
		return false // Installation error
	}

	return true
}

func (s *ApplicationState) GetHealthStatus() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"healthy":              s.IsHealthy(),
		"uptime":               time.Since(s.startTime).String(),
		"last_activity":        s.lastActivity,
		"current_mode":         s.currentMode,
		"configuration_loaded": s.configurationLoaded,
		"installation_active":  s.installationActive,
		"notification_count":   len(s.notifications),
		"session_id":           s.sessionID,
	}
}
