package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"mcf-dev/tui/internal/app"
	testutils "mcf-dev/tui/internal/testing"
)

// E2ETestSuite provides end-to-end testing for the MCF TUI application
type E2ETestSuite struct {
	suite.Suite
	tempDir       string
	configPath    string
	mockClient    *testutils.MockMCFClient
	orchestrator  *testutils.MockAgentOrchestrator
	configManager *testutils.MockConfigManager
	logger        *testutils.TestLogger
	ctx           context.Context
}

func (suite *E2ETestSuite) SetupSuite() {
	// Create temporary directory for test artifacts
	tempDir, err := os.MkdirTemp("", "mcf-tui-e2e-*")
	suite.Require().NoError(err)
	suite.tempDir = tempDir
	suite.configPath = filepath.Join(tempDir, "config.json")
	suite.ctx = context.Background()
}

func (suite *E2ETestSuite) SetupTest() {
	// Initialize mock dependencies
	suite.mockClient = testutils.NewMockMCFClient()
	suite.orchestrator = testutils.NewMockAgentOrchestrator()
	suite.configManager = testutils.NewMockConfigManager()
	suite.logger = testutils.NewTestLogger(suite.T())

	// Setup default mock data
	suite.setupDefaultMockData()
}

func (suite *E2ETestSuite) TearDownSuite() {
	// Clean up temporary directory
	os.RemoveAll(suite.tempDir)
}

func (suite *E2ETestSuite) setupDefaultMockData() {
	// System health data
	healthStatus := testutils.SystemHealthStatus{
		Status:  "healthy",
		Version: "1.0.0-test",
		Uptime:  time.Hour * 24,
		Components: map[string]string{
			"api":      "healthy",
			"database": "healthy",
			"cache":    "healthy",
		},
		Memory: testutils.MemoryStats{
			Used:      1024 * 1024 * 512,  // 512MB
			Available: 1024 * 1024 * 1536, // 1.5GB
			Total:     1024 * 1024 * 2048, // 2GB
			Percent:   25.0,
		},
		CPU: testutils.CPUStats{
			Usage:  15.5,
			Cores:  4,
			Load1:  0.8,
			Load5:  1.2,
			Load15: 1.0,
		},
	}
	suite.mockClient.SetSystemHealth(healthStatus)

	// Service data
	services := []testutils.ServiceStatus{
		{
			Name:        "api-server",
			Status:      "running",
			Uptime:      time.Hour * 24,
			Port:        8080,
			Health:      "healthy",
			LastChecked: time.Now(),
			Metadata:    map[string]string{"version": "1.0.0"},
		},
		{
			Name:        "database",
			Status:      "running",
			Uptime:      time.Hour * 72,
			Port:        5432,
			Health:      "healthy",
			LastChecked: time.Now(),
			Metadata:    map[string]string{"version": "13.0"},
		},
		{
			Name:        "cache",
			Status:      "running",
			Uptime:      time.Hour * 48,
			Port:        6379,
			Health:      "healthy",
			LastChecked: time.Now(),
			Metadata:    map[string]string{"version": "6.0"},
		},
	}
	for _, service := range services {
		suite.mockClient.AddService(service)
	}

	// Log data
	logs := []testutils.LogEntry{
		{
			Timestamp: time.Now().Add(-time.Hour),
			Level:     "info",
			Service:   "api-server",
			Message:   "Server started successfully",
			Fields:    map[string]string{"component": "main"},
		},
		{
			Timestamp: time.Now().Add(-30 * time.Minute),
			Level:     "warn",
			Service:   "database",
			Message:   "Connection pool near capacity",
			Fields:    map[string]string{"pool_size": "95"},
		},
		{
			Timestamp: time.Now().Add(-15 * time.Minute),
			Level:     "error",
			Service:   "api-server",
			Message:   "Request timeout occurred",
			Fields:    map[string]string{"endpoint": "/api/data", "duration": "30s"},
		},
	}
	for _, log := range logs {
		suite.mockClient.AddLog(log)
	}

	// Agent data
	agents := []*testutils.MockAgent{
		{
			ID:           "analyzer-agent",
			Name:         "Data Analyzer",
			Capabilities: []string{"analysis", "data-processing", "reporting"},
			Status:       "active",
			Metadata:     map[string]string{"version": "2.0.0", "type": "analyzer"},
			LastSeen:     time.Now(),
		},
		{
			ID:           "monitor-agent",
			Name:         "System Monitor",
			Capabilities: []string{"monitoring", "alerting", "health-check"},
			Status:       "active",
			Metadata:     map[string]string{"version": "1.5.0", "type": "monitor"},
			LastSeen:     time.Now(),
		},
	}
	for _, agent := range agents {
		suite.orchestrator.RegisterAgent(agent)
	}
}

func (suite *E2ETestSuite) TestCompleteApplicationWorkflow() {
	// This test simulates a complete user workflow from startup to shutdown

	// Step 1: Application Startup
	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Wait for initialization
	time.Sleep(100 * time.Millisecond)

	// Step 2: Window Size Setup
	// Simulate terminal resize
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Verify dashboard is displayed
	runner.WaitForOutput(suite.T(), "MCF TUI Development", 2*time.Second)

	// Step 3: Navigation Testing
	// Switch to command bar
	runner.SendKeypress(":")
	time.Sleep(50 * time.Millisecond)

	// Type help command
	runner.SendKeys("h", "e", "l", "p")
	time.Sleep(50 * time.Millisecond)

	// Press enter to execute
	runner.SendKeypress("\r")
	time.Sleep(100 * time.Millisecond)

	// Step 4: Return to dashboard
	runner.SendKeypress("\x1b") // Escape
	time.Sleep(50 * time.Millisecond)

	// Step 5: Test rapid navigation
	views := []string{":", "l", "o", "g", "s", "\r"} // :logs command
	for _, key := range views {
		runner.SendKeypress(key)
		time.Sleep(10 * time.Millisecond)
	}

	runner.SendKeypress("\x1b") // Return to dashboard
	time.Sleep(50 * time.Millisecond)

	// Step 6: Quit application
	runner.SendKeypress("q")
	time.Sleep(100 * time.Millisecond)

	// Test should complete without hanging
}

func (suite *E2ETestSuite) TestSystemHealthMonitoringWorkflow() {
	// Test complete health monitoring workflow

	// Setup mock expectations
	suite.mockClient.On("GetSystemHealth", suite.ctx).Return(
		suite.mockClient.GetSystemHealthStatus(), nil).Times(5)

	// Simulate health monitoring workflow
	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize window
	runner.SendKeypress(tea.WindowSizeMsg{Width: 100, Height: 30}.String())
	time.Sleep(50 * time.Millisecond)

	// Navigate through health information
	runner.SendKeys(":", "s", "t", "a", "t", "u", "s", "\r")
	time.Sleep(100 * time.Millisecond)

	// Check if health information is displayed
	runner.WaitForOutput(suite.T(), "healthy", 2*time.Second)

	// Return to dashboard and quit
	runner.SendKeys("\x1b", "q")
	time.Sleep(100 * time.Millisecond)
}

func (suite *E2ETestSuite) TestServiceManagementWorkflow() {
	// Test complete service management workflow

	// Setup mock expectations
	suite.mockClient.On("GetServices", suite.ctx).Return(
		suite.mockClient.GetServicesData(), nil).Times(3)

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize and navigate to services
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Access service management
	runner.SendKeys(":", "s", "e", "r", "v", "i", "c", "e", "s", "\r")
	time.Sleep(100 * time.Millisecond)

	// Navigate through services (simulate arrow key navigation)
	runner.SendKeys("j", "j", "k") // Down, down, up
	time.Sleep(50 * time.Millisecond)

	// Return and quit
	runner.SendKeys("\x1b", "q")
	time.Sleep(100 * time.Millisecond)
}

func (suite *E2ETestSuite) TestLogViewingWorkflow() {
	// Test complete log viewing workflow

	// Setup mock expectations
	filter := testutils.LogFilter{Limit: 100}
	suite.mockClient.On("GetLogs", suite.ctx, filter).Return(
		suite.mockClient.GetLogsData(), nil).Times(2)

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize and navigate to logs
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Access logs
	runner.SendKeys(":", "l", "o", "g", "s", "\r")
	time.Sleep(100 * time.Millisecond)

	// Test log filtering
	runner.SendKeys("/", "e", "r", "r", "o", "r") // Search for errors
	time.Sleep(50 * time.Millisecond)

	// Clear filter
	runner.SendKeys("\x1b", "/", "\r") // Escape, then clear search
	time.Sleep(50 * time.Millisecond)

	// Return and quit
	runner.SendKeys("\x1b", "q")
	time.Sleep(100 * time.Millisecond)
}

func (suite *E2ETestSuite) TestAgentOrchestrationWorkflow() {
	// Test complete agent orchestration workflow

	// Setup mock expectations for agent operations
	suite.orchestrator.On("GetAgents").Return(suite.orchestrator.GetAgents()).Times(2)

	task := testutils.Task{
		ID:      "test-task-1",
		Type:    "analysis",
		Payload: map[string]interface{}{"data": "test-dataset"},
	}
	suite.orchestrator.On("SubmitTask", suite.ctx, task).Return(nil).Once()

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize and navigate to agents
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Access agent management
	runner.SendKeys(":", "a", "g", "e", "n", "t", "s", "\r")
	time.Sleep(100 * time.Millisecond)

	// Simulate agent interaction
	runner.SendKeys("j", "k") // Navigate through agents
	time.Sleep(50 * time.Millisecond)

	// Submit a task
	runner.SendKeys("t") // Task submission key (hypothetical)
	time.Sleep(50 * time.Millisecond)

	// Return and quit
	runner.SendKeys("\x1b", "q")
	time.Sleep(100 * time.Millisecond)
}

func (suite *E2ETestSuite) TestErrorHandlingWorkflow() {
	// Test error handling throughout the application

	// Setup error scenarios
	errorClient := testutils.NewMockMCFClient()
	errorClient.SetConnected(false)
	errorClient.On("GetSystemHealth", suite.ctx).Return(
		testutils.SystemHealthStatus{}, fmt.Errorf("connection failed"))

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize window
	runner.SendKeypress(tea.WindowSizeMsg{Width: 100, Height: 30}.String())
	time.Sleep(50 * time.Millisecond)

	// Try to access system status (should handle error gracefully)
	runner.SendKeys(":", "s", "t", "a", "t", "u", "s", "\r")
	time.Sleep(100 * time.Millisecond)

	// Application should not crash
	runner.SendKeys("\x1b", "q")
	time.Sleep(100 * time.Millisecond)
}

func (suite *E2ETestSuite) TestPerformanceUnderLoad() {
	// Test application performance under simulated load

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize window
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Perform rapid navigation
	benchmark := testutils.NewPerformanceBenchmark("e2e_navigation", func() error {
		// Rapid view switching
		runner.SendKeys(":", "d", "a", "s", "h", "\r") // Dashboard
		runner.SendKeys(":", "l", "o", "g", "s", "\r") // Logs
		runner.SendKeys(":", "s", "e", "r", "v", "\r") // Services
		runner.SendKeypress("\x1b")                    // Return to dashboard
		return nil
	}).WithIterations(10).WithTimeout(30 * time.Second)

	benchmark.Run(suite.T())

	// Cleanup
	runner.SendKeypress("q")
	time.Sleep(100 * time.Millisecond)
}

func (suite *E2ETestSuite) TestConfigurationIntegration() {
	// Test configuration integration throughout the application

	// Setup configuration
	suite.configManager.On("Get", "tui.theme").Return("dark", true)
	suite.configManager.On("Get", "tui.refresh_rate").Return(1000, true)
	suite.configManager.On("Get", "mcf.host").Return("localhost", true)
	suite.configManager.On("Get", "mcf.port").Return(8080, true)

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Test configuration-dependent behavior
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Navigate to settings (hypothetical)
	runner.SendKeys(":", "c", "o", "n", "f", "i", "g", "\r")
	time.Sleep(100 * time.Millisecond)

	// Test theme switching
	runner.SendKeys("t") // Toggle theme (hypothetical)
	time.Sleep(50 * time.Millisecond)

	// Return and quit
	runner.SendKeys("\x1b", "q")
	time.Sleep(100 * time.Millisecond)

	// Verify configuration methods were called
	suite.configManager.AssertExpectations(suite.T())
}

func (suite *E2ETestSuite) TestLongRunningSession() {
	// Test application stability over extended period

	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(suite.T())
	defer runner.Stop()

	// Initialize
	runner.SendKeypress(tea.WindowSizeMsg{Width: 120, Height: 40}.String())
	time.Sleep(50 * time.Millisecond)

	// Simulate periodic activity over time
	for i := 0; i < 20; i++ {
		// Cycle through different views
		runner.SendKeys(":", "d", "a", "s", "h", "\r")
		time.Sleep(100 * time.Millisecond)

		runner.SendKeys(":", "l", "o", "g", "s", "\r")
		time.Sleep(100 * time.Millisecond)

		runner.SendKeypress("\x1b")
		time.Sleep(50 * time.Millisecond)

		// Simulate periodic refresh
		if i%5 == 0 {
			runner.SendKeys("r") // Refresh key (hypothetical)
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Application should still be responsive
	runner.SendKeypress("q")
	time.Sleep(100 * time.Millisecond)
}

// Run the test suite
func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

// Helper methods for mock data access
func (client *testutils.MockMCFClient) GetSystemHealthStatus() testutils.SystemHealthStatus {
	// This would typically be a method on the mock client
	return testutils.SystemHealthStatus{
		Status:  "healthy",
		Version: "1.0.0-test",
		Uptime:  time.Hour * 24,
	}
}

func (client *testutils.MockMCFClient) GetServicesData() []testutils.ServiceStatus {
	// This would typically be a method on the mock client
	return []testutils.ServiceStatus{
		{Name: "api-server", Status: "running", Health: "healthy"},
		{Name: "database", Status: "running", Health: "healthy"},
	}
}

func (client *testutils.MockMCFClient) GetLogsData() []testutils.LogEntry {
	// This would typically be a method on the mock client
	return []testutils.LogEntry{
		{Level: "info", Service: "api", Message: "Server started"},
		{Level: "error", Service: "api", Message: "Request failed"},
	}
}

// Individual integration tests
func TestApplicationStartup(t *testing.T) {
	t.Run("should start application successfully", func(t *testing.T) {
		model := app.InitialModel()
		runner := testutils.NewTestProgramRunner(model)

		runner.Start(t)
		time.Sleep(100 * time.Millisecond)

		// Should not crash or hang
		runner.Stop()
	})
}

func TestKeyboardNavigation(t *testing.T) {
	t.Run("should handle keyboard navigation correctly", func(t *testing.T) {
		model := app.InitialModel()
		runner := testutils.NewTestProgramRunner(model)

		runner.Start(t)
		defer runner.Stop()

		// Test various key combinations
		keys := []string{
			":",                // Command mode
			"h", "e", "l", "p", // Help command
			"\r",     // Enter
			"\x1b",   // Escape
			"j", "k", // Navigation
			"q", // Quit
		}

		for _, key := range keys {
			runner.SendKeypress(key)
			time.Sleep(10 * time.Millisecond)
		}
	})
}

func TestWindowResizing(t *testing.T) {
	t.Run("should handle window resize correctly", func(t *testing.T) {
		model := app.InitialModel()
		runner := testutils.NewTestProgramRunner(model)

		runner.Start(t)
		defer runner.Stop()

		// Test different window sizes
		sizes := []struct{ width, height int }{
			{80, 24},
			{120, 40},
			{200, 60},
			{40, 10}, // Very small
		}

		for _, size := range sizes {
			runner.SendKeypress(tea.WindowSizeMsg{
				Width:  size.width,
				Height: size.height,
			}.String())
			time.Sleep(50 * time.Millisecond)
		}

		runner.SendKeypress("q")
	})
}

func TestRapidInputHandling(t *testing.T) {
	t.Run("should handle rapid input without issues", func(t *testing.T) {
		model := app.InitialModel()
		runner := testutils.NewTestProgramRunner(model)

		runner.Start(t)
		defer runner.Stop()

		// Initialize window
		runner.SendKeypress(tea.WindowSizeMsg{Width: 100, Height: 30}.String())
		time.Sleep(50 * time.Millisecond)

		// Send rapid keypresses
		for i := 0; i < 100; i++ {
			runner.SendKeypress("j")
			if i%10 == 0 {
				time.Sleep(time.Millisecond) // Brief pause
			}
		}

		// Should still be responsive
		runner.SendKeypress("q")
		time.Sleep(100 * time.Millisecond)
	})
}

// Benchmark tests
func BenchmarkApplicationStartup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := app.InitialModel()
		runner := testutils.NewTestProgramRunner(model)

		runner.Start(&testing.T{})
		time.Sleep(10 * time.Millisecond)
		runner.Stop()
	}
}

func BenchmarkKeyboardInput(b *testing.B) {
	model := app.InitialModel()
	runner := testutils.NewTestProgramRunner(model)

	runner.Start(&testing.T{})
	defer runner.Stop()

	runner.SendKeypress(tea.WindowSizeMsg{Width: 100, Height: 30}.String())
	time.Sleep(50 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner.SendKeypress("j")
	}
}
