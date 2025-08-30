package testing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockMCFClient provides a mock implementation of the MCF client
type MockMCFClient struct {
	mock.Mock
	mu            sync.RWMutex
	systemHealth  SystemHealthStatus
	services      []ServiceStatus
	logs          []LogEntry
	agentStates   map[string]AgentState
	configuration map[string]interface{}
	connected     bool
	responseDelay time.Duration
}

// SystemHealthStatus represents the overall system health
type SystemHealthStatus struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	Uptime     time.Duration     `json:"uptime"`
	Components map[string]string `json:"components"`
	Memory     MemoryStats       `json:"memory"`
	CPU        CPUStats          `json:"cpu"`
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Uptime      time.Duration     `json:"uptime"`
	Port        int               `json:"port"`
	Health      string            `json:"health"`
	Metadata    map[string]string `json:"metadata"`
	LastChecked time.Time         `json:"last_checked"`
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields"`
}

// AgentState represents the state of an agent
type AgentState struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Status       string            `json:"status"`
	LastSeen     time.Time         `json:"last_seen"`
	Capabilities []string          `json:"capabilities"`
	Metadata     map[string]string `json:"metadata"`
}

// MemoryStats represents memory statistics
type MemoryStats struct {
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Total     uint64  `json:"total"`
	Percent   float64 `json:"percent"`
}

// CPUStats represents CPU statistics
type CPUStats struct {
	Usage  float64 `json:"usage"`
	Cores  int     `json:"cores"`
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// NewMockMCFClient creates a new mock MCF client
func NewMockMCFClient() *MockMCFClient {
	return &MockMCFClient{
		systemHealth: SystemHealthStatus{
			Status:  "healthy",
			Version: "1.0.0-test",
			Uptime:  time.Hour * 24,
			Components: map[string]string{
				"database": "healthy",
				"cache":    "healthy",
				"api":      "healthy",
			},
			Memory: MemoryStats{
				Used:      1024 * 1024 * 500,  // 500MB
				Available: 1024 * 1024 * 1500, // 1.5GB
				Total:     1024 * 1024 * 2000, // 2GB
				Percent:   25.0,
			},
			CPU: CPUStats{
				Usage:  15.5,
				Cores:  4,
				Load1:  0.8,
				Load5:  1.2,
				Load15: 1.0,
			},
		},
		services:      []ServiceStatus{},
		logs:          []LogEntry{},
		agentStates:   make(map[string]AgentState),
		configuration: make(map[string]interface{}),
		connected:     true,
		responseDelay: 0,
	}
}

// Connect simulates connecting to MCF
func (m *MockMCFClient) Connect(ctx context.Context) error {
	args := m.Called(ctx)

	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	m.mu.Lock()
	m.connected = true
	m.mu.Unlock()

	return args.Error(0)
}

// Disconnect simulates disconnecting from MCF
func (m *MockMCFClient) Disconnect() error {
	args := m.Called()

	m.mu.Lock()
	m.connected = false
	m.mu.Unlock()

	return args.Error(0)
}

// IsConnected returns the connection status
func (m *MockMCFClient) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

// GetSystemHealth returns mock system health data
func (m *MockMCFClient) GetSystemHealth(ctx context.Context) (SystemHealthStatus, error) {
	args := m.Called(ctx)

	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := args.Error(1); err != nil {
		return SystemHealthStatus{}, err
	}

	return m.systemHealth, nil
}

// GetServices returns mock service data
func (m *MockMCFClient) GetServices(ctx context.Context) ([]ServiceStatus, error) {
	args := m.Called(ctx)

	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return m.services, nil
}

// GetLogs returns mock log data
func (m *MockMCFClient) GetLogs(ctx context.Context, filter LogFilter) ([]LogEntry, error) {
	args := m.Called(ctx, filter)

	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := args.Error(1); err != nil {
		return nil, err
	}

	// Apply filter to logs
	filteredLogs := m.filterLogs(m.logs, filter)
	return filteredLogs, nil
}

// GetAgentStates returns mock agent state data
func (m *MockMCFClient) GetAgentStates(ctx context.Context) (map[string]AgentState, error) {
	args := m.Called(ctx)

	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return m.agentStates, nil
}

// ExecuteCommand simulates command execution
func (m *MockMCFClient) ExecuteCommand(ctx context.Context, command string, args []string) (CommandResult, error) {
	mockArgs := m.Called(ctx, command, args)

	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	if err := mockArgs.Error(1); err != nil {
		return CommandResult{}, err
	}

	// Simulate command execution
	result := CommandResult{
		Command:   command,
		Args:      args,
		ExitCode:  0,
		Output:    fmt.Sprintf("Mock output for command: %s", command),
		Error:     "",
		Duration:  time.Millisecond * 100,
		Timestamp: time.Now(),
	}

	return result, nil
}

// LogFilter represents log filtering criteria
type LogFilter struct {
	Service   string
	Level     string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Command   string        `json:"command"`
	Args      []string      `json:"args"`
	ExitCode  int           `json:"exit_code"`
	Output    string        `json:"output"`
	Error     string        `json:"error"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// Test data setup methods
func (m *MockMCFClient) SetSystemHealth(health SystemHealthStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.systemHealth = health
}

func (m *MockMCFClient) AddService(service ServiceStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.services = append(m.services, service)
}

func (m *MockMCFClient) AddLog(entry LogEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, entry)
}

func (m *MockMCFClient) SetAgentState(id string, state AgentState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.agentStates[id] = state
}

func (m *MockMCFClient) SetConnected(connected bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = connected
}

func (m *MockMCFClient) SetResponseDelay(delay time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responseDelay = delay
}

// Helper methods
func (m *MockMCFClient) filterLogs(logs []LogEntry, filter LogFilter) []LogEntry {
	var filtered []LogEntry

	for _, log := range logs {
		if filter.Service != "" && log.Service != filter.Service {
			continue
		}
		if filter.Level != "" && log.Level != filter.Level {
			continue
		}
		if !filter.StartTime.IsZero() && log.Timestamp.Before(filter.StartTime) {
			continue
		}
		if !filter.EndTime.IsZero() && log.Timestamp.After(filter.EndTime) {
			continue
		}

		filtered = append(filtered, log)

		if filter.Limit > 0 && len(filtered) >= filter.Limit {
			break
		}
	}

	return filtered
}

// MockAgentOrchestrator provides mock agent orchestration functionality
type MockAgentOrchestrator struct {
	mock.Mock
	mu        sync.RWMutex
	agents    map[string]*MockAgent
	tasks     []Task
	connected bool
}

// Task represents an orchestration task
type Task struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	Status      string                 `json:"status"`
	Result      interface{}            `json:"result"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// MockAgent represents a mock agent
type MockAgent struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Capabilities []string          `json:"capabilities"`
	Status       string            `json:"status"`
	Metadata     map[string]string `json:"metadata"`
	LastSeen     time.Time         `json:"last_seen"`
}

// NewMockAgentOrchestrator creates a new mock agent orchestrator
func NewMockAgentOrchestrator() *MockAgentOrchestrator {
	return &MockAgentOrchestrator{
		agents:    make(map[string]*MockAgent),
		tasks:     []Task{},
		connected: true,
	}
}

// RegisterAgent registers a mock agent
func (o *MockAgentOrchestrator) RegisterAgent(agent *MockAgent) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.agents[agent.ID] = agent
}

// GetAgents returns all registered agents
func (o *MockAgentOrchestrator) GetAgents() map[string]*MockAgent {
	o.mu.RLock()
	defer o.mu.RUnlock()

	agents := make(map[string]*MockAgent)
	for k, v := range o.agents {
		agents[k] = v
	}
	return agents
}

// SubmitTask submits a task for execution
func (o *MockAgentOrchestrator) SubmitTask(ctx context.Context, task Task) error {
	args := o.Called(ctx, task)

	o.mu.Lock()
	defer o.mu.Unlock()

	task.CreatedAt = time.Now()
	task.Status = "pending"
	o.tasks = append(o.tasks, task)

	return args.Error(0)
}

// GetTasks returns all tasks
func (o *MockAgentOrchestrator) GetTasks() []Task {
	o.mu.RLock()
	defer o.mu.RUnlock()

	tasks := make([]Task, len(o.tasks))
	copy(tasks, o.tasks)
	return tasks
}

// MockConfigManager provides mock configuration management
type MockConfigManager struct {
	mock.Mock
	mu     sync.RWMutex
	config map[string]interface{}
}

// NewMockConfigManager creates a new mock configuration manager
func NewMockConfigManager() *MockConfigManager {
	return &MockConfigManager{
		config: make(map[string]interface{}),
	}
}

// Get retrieves a configuration value
func (c *MockConfigManager) Get(key string) (interface{}, bool) {
	args := c.Called(key)

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.config[key]
	return value, exists && args.Bool(1)
}

// Set sets a configuration value
func (c *MockConfigManager) Set(key string, value interface{}) error {
	args := c.Called(key, value)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.config[key] = value
	return args.Error(0)
}

// GetAll returns all configuration values
func (c *MockConfigManager) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	config := make(map[string]interface{})
	for k, v := range c.config {
		config[k] = v
	}
	return config
}

// Save persists the configuration
func (c *MockConfigManager) Save() error {
	args := c.Called()
	return args.Error(0)
}

// Load loads the configuration
func (c *MockConfigManager) Load() error {
	args := c.Called()
	return args.Error(0)
}
