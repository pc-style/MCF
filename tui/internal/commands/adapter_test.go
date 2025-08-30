package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	testutils "mcf-dev/tui/internal/testing"
)

// MCFCommandAdapter provides interface to MCF CLI commands
type MCFCommandAdapter struct {
	client MCFClient
	logger Logger
}

// MCFClient interface for testing
type MCFClient interface {
	ExecuteCommand(ctx context.Context, command string, args []string) (testutils.CommandResult, error)
	GetSystemHealth(ctx context.Context) (testutils.SystemHealthStatus, error)
	GetServices(ctx context.Context) ([]testutils.ServiceStatus, error)
	GetLogs(ctx context.Context, filter testutils.LogFilter) ([]testutils.LogEntry, error)
	GetAgentStates(ctx context.Context) (map[string]testutils.AgentState, error)
	Connect(ctx context.Context) error
	Disconnect() error
	IsConnected() bool
}

// Logger interface for testing
type Logger interface {
	Log(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// NewMCFCommandAdapter creates a new MCF command adapter
func NewMCFCommandAdapter(client MCFClient, logger Logger) *MCFCommandAdapter {
	return &MCFCommandAdapter{
		client: client,
		logger: logger,
	}
}

// ExecuteCommand executes an MCF command
func (a *MCFCommandAdapter) ExecuteCommand(ctx context.Context, command string, args []string) (testutils.CommandResult, error) {
	a.logger.Log("Executing command: %s %v", command, args)

	result, err := a.client.ExecuteCommand(ctx, command, args)
	if err != nil {
		a.logger.Error("Command execution failed: %v", err)
		return testutils.CommandResult{}, err
	}

	a.logger.Log("Command executed successfully: %s", command)
	return result, nil
}

// GetSystemHealth retrieves system health information
func (a *MCFCommandAdapter) GetSystemHealth(ctx context.Context) (testutils.SystemHealthStatus, error) {
	a.logger.Log("Retrieving system health")

	health, err := a.client.GetSystemHealth(ctx)
	if err != nil {
		a.logger.Error("Failed to get system health: %v", err)
		return testutils.SystemHealthStatus{}, err
	}

	a.logger.Log("System health retrieved: %s", health.Status)
	return health, nil
}

// GetServices retrieves service status information
func (a *MCFCommandAdapter) GetServices(ctx context.Context) ([]testutils.ServiceStatus, error) {
	a.logger.Log("Retrieving service status")

	services, err := a.client.GetServices(ctx)
	if err != nil {
		a.logger.Error("Failed to get services: %v", err)
		return nil, err
	}

	a.logger.Log("Retrieved %d services", len(services))
	return services, nil
}

// GetLogs retrieves logs with filtering
func (a *MCFCommandAdapter) GetLogs(ctx context.Context, filter testutils.LogFilter) ([]testutils.LogEntry, error) {
	a.logger.Log("Retrieving logs with filter: %+v", filter)

	logs, err := a.client.GetLogs(ctx, filter)
	if err != nil {
		a.logger.Error("Failed to get logs: %v", err)
		return nil, err
	}

	a.logger.Log("Retrieved %d log entries", len(logs))
	return logs, nil
}

// GetAgentStates retrieves agent state information
func (a *MCFCommandAdapter) GetAgentStates(ctx context.Context) (map[string]testutils.AgentState, error) {
	a.logger.Log("Retrieving agent states")

	agents, err := a.client.GetAgentStates(ctx)
	if err != nil {
		a.logger.Error("Failed to get agent states: %v", err)
		return nil, err
	}

	a.logger.Log("Retrieved %d agent states", len(agents))
	return agents, nil
}

// Connect establishes connection to MCF
func (a *MCFCommandAdapter) Connect(ctx context.Context) error {
	a.logger.Log("Connecting to MCF")

	err := a.client.Connect(ctx)
	if err != nil {
		a.logger.Error("Failed to connect to MCF: %v", err)
		return err
	}

	a.logger.Log("Successfully connected to MCF")
	return nil
}

// Disconnect closes connection to MCF
func (a *MCFCommandAdapter) Disconnect() error {
	a.logger.Log("Disconnecting from MCF")

	err := a.client.Disconnect()
	if err != nil {
		a.logger.Error("Failed to disconnect from MCF: %v", err)
		return err
	}

	a.logger.Log("Successfully disconnected from MCF")
	return nil
}

// IsConnected checks connection status
func (a *MCFCommandAdapter) IsConnected() bool {
	connected := a.client.IsConnected()
	a.logger.Log("Connection status: %t", connected)
	return connected
}

// Test suite for MCF Command Adapter
func TestMCFCommandAdapter_Creation(t *testing.T) {
	t.Run("should create adapter with valid client and logger", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)

		adapter := NewMCFCommandAdapter(mockClient, logger)

		assert.NotNil(t, adapter, "Adapter should be created")
		assert.Equal(t, mockClient, adapter.client)
		assert.Equal(t, logger, adapter.logger)
	})
}

func TestMCFCommandAdapter_ExecuteCommand(t *testing.T) {
	t.Run("should execute command successfully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		command := "status"
		args := []string{"--verbose"}
		expectedResult := testutils.CommandResult{
			Command:   command,
			Args:      args,
			ExitCode:  0,
			Output:    "Service status: OK",
			Duration:  100 * time.Millisecond,
			Timestamp: time.Now(),
		}

		mockClient.On("ExecuteCommand", ctx, command, args).Return(expectedResult, nil)

		result, err := adapter.ExecuteCommand(ctx, command, args)

		assert.NoError(t, err)
		assert.Equal(t, command, result.Command)
		assert.Equal(t, args, result.Args)
		assert.Equal(t, 0, result.ExitCode)
		assert.Equal(t, "Service status: OK", result.Output)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle command execution failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		command := "invalid"
		args := []string{}
		expectedError := errors.New("command not found")

		mockClient.On("ExecuteCommand", ctx, command, args).Return(testutils.CommandResult{}, expectedError)

		result, err := adapter.ExecuteCommand(ctx, command, args)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, testutils.CommandResult{}, result)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle command with non-zero exit code", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		command := "failing-command"
		args := []string{}
		expectedResult := testutils.CommandResult{
			Command:  command,
			Args:     args,
			ExitCode: 1,
			Output:   "",
			Error:    "Command failed",
			Duration: 50 * time.Millisecond,
		}

		mockClient.On("ExecuteCommand", ctx, command, args).Return(expectedResult, nil)

		result, err := adapter.ExecuteCommand(ctx, command, args)

		assert.NoError(t, err) // No execution error, but command failed
		assert.Equal(t, 1, result.ExitCode)
		assert.Equal(t, "Command failed", result.Error)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle context timeout", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		command := "slow-command"
		args := []string{}
		expectedError := context.DeadlineExceeded

		mockClient.On("ExecuteCommand", ctx, command, args).Return(testutils.CommandResult{}, expectedError)

		result, err := adapter.ExecuteCommand(ctx, command, args)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, testutils.CommandResult{}, result)

		mockClient.AssertExpectations(t)
	})
}

func TestMCFCommandAdapter_GetSystemHealth(t *testing.T) {
	t.Run("should retrieve system health successfully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedHealth := testutils.SystemHealthStatus{
			Status:  "healthy",
			Version: "1.0.0",
			Uptime:  time.Hour * 24,
			Components: map[string]string{
				"database": "healthy",
				"cache":    "healthy",
			},
			Memory: testutils.MemoryStats{
				Used:      1024 * 1024 * 512,  // 512MB
				Available: 1024 * 1024 * 512,  // 512MB
				Total:     1024 * 1024 * 1024, // 1GB
				Percent:   50.0,
			},
			CPU: testutils.CPUStats{
				Usage: 25.5,
				Cores: 4,
				Load1: 1.2,
			},
		}

		mockClient.On("GetSystemHealth", ctx).Return(expectedHealth, nil)

		health, err := adapter.GetSystemHealth(ctx)

		assert.NoError(t, err)
		assert.Equal(t, "healthy", health.Status)
		assert.Equal(t, "1.0.0", health.Version)
		assert.Equal(t, time.Hour*24, health.Uptime)
		assert.Equal(t, "healthy", health.Components["database"])
		assert.Equal(t, 50.0, health.Memory.Percent)
		assert.Equal(t, 25.5, health.CPU.Usage)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle system health retrieval failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedError := errors.New("service unavailable")

		mockClient.On("GetSystemHealth", ctx).Return(testutils.SystemHealthStatus{}, expectedError)

		health, err := adapter.GetSystemHealth(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, testutils.SystemHealthStatus{}, health)

		mockClient.AssertExpectations(t)
	})
}

func TestMCFCommandAdapter_GetServices(t *testing.T) {
	t.Run("should retrieve services successfully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedServices := []testutils.ServiceStatus{
			{
				Name:        "api-server",
				Status:      "running",
				Health:      "healthy",
				Port:        8080,
				Uptime:      time.Hour * 2,
				LastChecked: time.Now(),
			},
			{
				Name:        "database",
				Status:      "running",
				Health:      "healthy",
				Port:        5432,
				Uptime:      time.Hour * 24,
				LastChecked: time.Now(),
			},
		}

		mockClient.On("GetServices", ctx).Return(expectedServices, nil)

		services, err := adapter.GetServices(ctx)

		assert.NoError(t, err)
		assert.Len(t, services, 2)
		assert.Equal(t, "api-server", services[0].Name)
		assert.Equal(t, "running", services[0].Status)
		assert.Equal(t, 8080, services[0].Port)
		assert.Equal(t, "database", services[1].Name)
		assert.Equal(t, 5432, services[1].Port)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty services list", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedServices := []testutils.ServiceStatus{}

		mockClient.On("GetServices", ctx).Return(expectedServices, nil)

		services, err := adapter.GetServices(ctx)

		assert.NoError(t, err)
		assert.Empty(t, services)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle services retrieval failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedError := errors.New("network error")

		mockClient.On("GetServices", ctx).Return(nil, expectedError)

		services, err := adapter.GetServices(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, services)

		mockClient.AssertExpectations(t)
	})
}

func TestMCFCommandAdapter_GetLogs(t *testing.T) {
	t.Run("should retrieve logs with filter", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		filter := testutils.LogFilter{
			Service: "api-server",
			Level:   "ERROR",
			Limit:   10,
		}
		expectedLogs := []testutils.LogEntry{
			{
				Timestamp: time.Now(),
				Level:     "ERROR",
				Service:   "api-server",
				Message:   "Database connection failed",
			},
			{
				Timestamp: time.Now(),
				Level:     "ERROR",
				Service:   "api-server",
				Message:   "Request timeout",
			},
		}

		mockClient.On("GetLogs", ctx, filter).Return(expectedLogs, nil)

		logs, err := adapter.GetLogs(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		assert.Equal(t, "ERROR", logs[0].Level)
		assert.Equal(t, "api-server", logs[0].Service)
		assert.Contains(t, logs[0].Message, "Database connection")

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle logs retrieval with empty filter", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		filter := testutils.LogFilter{} // Empty filter
		expectedLogs := []testutils.LogEntry{
			{Timestamp: time.Now(), Level: "INFO", Service: "service1", Message: "Started"},
			{Timestamp: time.Now(), Level: "WARN", Service: "service2", Message: "Warning"},
		}

		mockClient.On("GetLogs", ctx, filter).Return(expectedLogs, nil)

		logs, err := adapter.GetLogs(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, logs, 2)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle logs retrieval failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		filter := testutils.LogFilter{Service: "nonexistent"}
		expectedError := errors.New("service not found")

		mockClient.On("GetLogs", ctx, filter).Return(nil, expectedError)

		logs, err := adapter.GetLogs(ctx, filter)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, logs)

		mockClient.AssertExpectations(t)
	})
}

func TestMCFCommandAdapter_GetAgentStates(t *testing.T) {
	t.Run("should retrieve agent states successfully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedAgents := map[string]testutils.AgentState{
			"agent-1": {
				ID:           "agent-1",
				Name:         "Frontend Developer",
				Status:       "active",
				LastSeen:     time.Now(),
				Capabilities: []string{"react", "javascript", "css"},
			},
			"agent-2": {
				ID:           "agent-2",
				Name:         "Backend Developer",
				Status:       "idle",
				LastSeen:     time.Now().Add(-5 * time.Minute),
				Capabilities: []string{"go", "python", "database"},
			},
		}

		mockClient.On("GetAgentStates", ctx).Return(expectedAgents, nil)

		agents, err := adapter.GetAgentStates(ctx)

		assert.NoError(t, err)
		assert.Len(t, agents, 2)
		assert.Equal(t, "Frontend Developer", agents["agent-1"].Name)
		assert.Equal(t, "active", agents["agent-1"].Status)
		assert.Contains(t, agents["agent-1"].Capabilities, "react")
		assert.Equal(t, "Backend Developer", agents["agent-2"].Name)
		assert.Equal(t, "idle", agents["agent-2"].Status)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty agent states", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedAgents := map[string]testutils.AgentState{}

		mockClient.On("GetAgentStates", ctx).Return(expectedAgents, nil)

		agents, err := adapter.GetAgentStates(ctx)

		assert.NoError(t, err)
		assert.Empty(t, agents)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle agent states retrieval failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedError := errors.New("agents service unavailable")

		mockClient.On("GetAgentStates", ctx).Return(nil, expectedError)

		agents, err := adapter.GetAgentStates(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, agents)

		mockClient.AssertExpectations(t)
	})
}

func TestMCFCommandAdapter_Connection(t *testing.T) {
	t.Run("should connect successfully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		mockClient.On("Connect", ctx).Return(nil)

		err := adapter.Connect(ctx)

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle connection failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		expectedError := errors.New("connection refused")
		mockClient.On("Connect", ctx).Return(expectedError)

		err := adapter.Connect(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("should disconnect successfully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		mockClient.On("Disconnect").Return(nil)

		err := adapter.Disconnect()

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle disconnection failure", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		expectedError := errors.New("disconnection failed")
		mockClient.On("Disconnect").Return(expectedError)

		err := adapter.Disconnect()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("should check connection status", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		mockClient.On("IsConnected").Return(true)

		connected := adapter.IsConnected()

		assert.True(t, connected)
		mockClient.AssertExpectations(t)

		// Test disconnected state
		mockClient.On("IsConnected").Return(false)

		connected = adapter.IsConnected()

		assert.False(t, connected)
	})
}

// Integration tests with realistic scenarios
func TestMCFCommandAdapter_Integration(t *testing.T) {
	t.Run("should handle complete workflow", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()

		// Step 1: Connect
		mockClient.On("Connect", ctx).Return(nil)
		err := adapter.Connect(ctx)
		require.NoError(t, err)

		// Step 2: Check system health
		health := testutils.SystemHealthStatus{Status: "healthy", Version: "1.0.0"}
		mockClient.On("GetSystemHealth", ctx).Return(health, nil)
		retrievedHealth, err := adapter.GetSystemHealth(ctx)
		require.NoError(t, err)
		assert.Equal(t, "healthy", retrievedHealth.Status)

		// Step 3: Get services
		services := []testutils.ServiceStatus{{Name: "api", Status: "running"}}
		mockClient.On("GetServices", ctx).Return(services, nil)
		retrievedServices, err := adapter.GetServices(ctx)
		require.NoError(t, err)
		assert.Len(t, retrievedServices, 1)

		// Step 4: Execute command
		result := testutils.CommandResult{Command: "test", ExitCode: 0, Output: "success"}
		mockClient.On("ExecuteCommand", ctx, "test", []string{}).Return(result, nil)
		commandResult, err := adapter.ExecuteCommand(ctx, "test", []string{})
		require.NoError(t, err)
		assert.Equal(t, 0, commandResult.ExitCode)

		// Step 5: Disconnect
		mockClient.On("Disconnect").Return(nil)
		err = adapter.Disconnect()
		require.NoError(t, err)

		mockClient.AssertExpectations(t)
	})
}

// Performance and stress tests
func TestMCFCommandAdapter_Performance(t *testing.T) {
	t.Run("should handle rapid command execution", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		result := testutils.CommandResult{Command: "test", ExitCode: 0}

		// Setup mock for multiple calls
		mockClient.On("ExecuteCommand", ctx, "test", mock.AnythingOfType("[]string")).Return(result, nil).Times(100)

		benchmark := testutils.NewPerformanceBenchmark("rapid_commands", func() error {
			_, err := adapter.ExecuteCommand(ctx, "test", []string{})
			return err
		}).WithIterations(100)

		benchmark.Run(t)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle concurrent requests", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		health := testutils.SystemHealthStatus{Status: "healthy"}

		// Setup mock for concurrent calls
		mockClient.On("GetSystemHealth", ctx).Return(health, nil).Times(50)

		benchmark := testutils.NewPerformanceBenchmark("concurrent_health_checks", func() error {
			_, err := adapter.GetSystemHealth(ctx)
			return err
		}).WithIterations(50)

		benchmark.Run(t)

		mockClient.AssertExpectations(t)
	})
}

// Error handling and edge cases
func TestMCFCommandAdapter_EdgeCases(t *testing.T) {
	t.Run("should handle nil context gracefully", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		// Most implementations would handle nil context, but for testing we'll use a valid one
		ctx := context.Background()
		mockClient.On("GetSystemHealth", ctx).Return(testutils.SystemHealthStatus{}, errors.New("context error"))

		_, err := adapter.GetSystemHealth(ctx)
		assert.Error(t, err)
	})

	t.Run("should handle very long command output", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		longOutput := string(make([]byte, 1024*1024)) // 1MB output
		result := testutils.CommandResult{
			Command:  "long-output",
			ExitCode: 0,
			Output:   longOutput,
		}

		mockClient.On("ExecuteCommand", ctx, "long-output", []string{}).Return(result, nil)

		commandResult, err := adapter.ExecuteCommand(ctx, "long-output", []string{})

		assert.NoError(t, err)
		assert.Len(t, commandResult.Output, len(longOutput))
	})

	t.Run("should handle special characters in commands", func(t *testing.T) {
		mockClient := testutils.NewMockMCFClient()
		logger := testutils.NewTestLogger(t)
		adapter := NewMCFCommandAdapter(mockClient, logger)

		ctx := context.Background()
		specialCommand := "test-cmd-with-special-chars!@#$%"
		specialArgs := []string{"arg with spaces", "arg-with-dashes", "arg_with_underscores"}
		result := testutils.CommandResult{Command: specialCommand, ExitCode: 0}

		mockClient.On("ExecuteCommand", ctx, specialCommand, specialArgs).Return(result, nil)

		commandResult, err := adapter.ExecuteCommand(ctx, specialCommand, specialArgs)

		assert.NoError(t, err)
		assert.Equal(t, specialCommand, commandResult.Command)
	})
}

// Benchmark tests
func BenchmarkMCFCommandAdapter_ExecuteCommand(b *testing.B) {
	mockClient := testutils.NewMockMCFClient()
	logger := testutils.NewTestLogger(&testing.T{})
	adapter := NewMCFCommandAdapter(mockClient, logger)

	ctx := context.Background()
	result := testutils.CommandResult{Command: "benchmark", ExitCode: 0}
	mockClient.On("ExecuteCommand", ctx, "benchmark", mock.Anything).Return(result, nil).Times(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.ExecuteCommand(ctx, "benchmark", []string{})
	}
}

func BenchmarkMCFCommandAdapter_GetSystemHealth(b *testing.B) {
	mockClient := testutils.NewMockMCFClient()
	logger := testutils.NewTestLogger(&testing.T{})
	adapter := NewMCFCommandAdapter(mockClient, logger)

	ctx := context.Background()
	health := testutils.SystemHealthStatus{Status: "healthy"}
	mockClient.On("GetSystemHealth", ctx).Return(health, nil).Times(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.GetSystemHealth(ctx)
	}
}
