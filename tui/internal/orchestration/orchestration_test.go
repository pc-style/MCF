package orchestration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutils "mcf-dev/tui/internal/testing"
)

// OrchestrationTestSuite provides integration testing for agent orchestration
type OrchestrationTestSuite struct {
	suite.Suite
	orchestrator *testutils.MockAgentOrchestrator
	configMgr    *testutils.MockConfigManager
	logger       *testutils.TestLogger
	ctx          context.Context
}

func (suite *OrchestrationTestSuite) SetupTest() {
	suite.orchestrator = testutils.NewMockAgentOrchestrator()
	suite.configMgr = testutils.NewMockConfigManager()
	suite.logger = testutils.NewTestLogger(suite.T())
	suite.ctx = context.Background()
}

func (suite *OrchestrationTestSuite) TestAgentRegistrationFlow() {
	// Test complete agent registration workflow

	agent := &testutils.MockAgent{
		ID:           "test-agent-1",
		Name:         "Test Agent",
		Capabilities: []string{"test", "analysis"},
		Status:       "ready",
		Metadata:     map[string]string{"version": "1.0.0"},
		LastSeen:     time.Now(),
	}

	// Register agent
	suite.orchestrator.RegisterAgent(agent)

	// Verify registration
	agents := suite.orchestrator.GetAgents()
	suite.Require().Contains(agents, "test-agent-1")

	registeredAgent := agents["test-agent-1"]
	suite.Equal("Test Agent", registeredAgent.Name)
	suite.Contains(registeredAgent.Capabilities, "test")
	suite.Contains(registeredAgent.Capabilities, "analysis")
}

func (suite *OrchestrationTestSuite) TestTaskSubmissionAndExecution() {
	// Test task submission and execution workflow

	// Register test agent
	agent := &testutils.MockAgent{
		ID:           "worker-agent",
		Name:         "Worker Agent",
		Capabilities: []string{"processing", "analysis"},
		Status:       "active",
	}
	suite.orchestrator.RegisterAgent(agent)

	// Create task
	task := testutils.Task{
		ID:      "task-1",
		AgentID: "worker-agent",
		Type:    "processing",
		Payload: map[string]interface{}{
			"data":   "test-data",
			"action": "process",
		},
		Status: "pending",
	}

	// Mock task submission
	suite.orchestrator.On("SubmitTask", suite.ctx, mock.MatchedBy(func(t testutils.Task) bool {
		return t.ID == "task-1" && t.AgentID == "worker-agent"
	})).Return(nil)

	// Submit task
	err := suite.orchestrator.SubmitTask(suite.ctx, task)
	suite.NoError(err)

	// Verify task was submitted
	tasks := suite.orchestrator.GetTasks()
	suite.Len(tasks, 1)
	suite.Equal("task-1", tasks[0].ID)
	suite.Equal("worker-agent", tasks[0].AgentID)
	suite.Equal("pending", tasks[0].Status)

	suite.orchestrator.AssertExpectations(suite.T())
}

func (suite *OrchestrationTestSuite) TestMultipleAgentCoordination() {
	// Test coordination between multiple agents

	agents := []*testutils.MockAgent{
		{
			ID:           "analyzer-1",
			Name:         "Data Analyzer 1",
			Capabilities: []string{"analysis", "data-processing"},
			Status:       "active",
		},
		{
			ID:           "analyzer-2",
			Name:         "Data Analyzer 2",
			Capabilities: []string{"analysis", "visualization"},
			Status:       "active",
		},
		{
			ID:           "coordinator",
			Name:         "Task Coordinator",
			Capabilities: []string{"coordination", "management"},
			Status:       "active",
		},
	}

	// Register all agents
	for _, agent := range agents {
		suite.orchestrator.RegisterAgent(agent)
	}

	// Create coordinated tasks
	tasks := []testutils.Task{
		{
			ID:      "analysis-task-1",
			AgentID: "analyzer-1",
			Type:    "analysis",
			Payload: map[string]interface{}{"dataset": "data1"},
		},
		{
			ID:      "analysis-task-2",
			AgentID: "analyzer-2",
			Type:    "analysis",
			Payload: map[string]interface{}{"dataset": "data2"},
		},
		{
			ID:      "coordination-task",
			AgentID: "coordinator",
			Type:    "coordination",
			Payload: map[string]interface{}{
				"subtasks": []string{"analysis-task-1", "analysis-task-2"},
			},
		},
	}

	// Mock task submissions
	for _, task := range tasks {
		suite.orchestrator.On("SubmitTask", suite.ctx, mock.MatchedBy(func(t testutils.Task) bool {
			return t.ID == task.ID
		})).Return(nil)
	}

	// Submit all tasks
	for _, task := range tasks {
		err := suite.orchestrator.SubmitTask(suite.ctx, task)
		suite.NoError(err)
	}

	// Verify all agents and tasks
	registeredAgents := suite.orchestrator.GetAgents()
	suite.Len(registeredAgents, 3)

	submittedTasks := suite.orchestrator.GetTasks()
	suite.Len(submittedTasks, 3)

	suite.orchestrator.AssertExpectations(suite.T())
}

func (suite *OrchestrationTestSuite) TestAgentFailureHandling() {
	// Test handling of agent failures and recovery

	agent := &testutils.MockAgent{
		ID:           "failing-agent",
		Name:         "Failing Agent",
		Capabilities: []string{"processing"},
		Status:       "active",
	}
	suite.orchestrator.RegisterAgent(agent)

	// Create task that will fail
	failingTask := testutils.Task{
		ID:      "failing-task",
		AgentID: "failing-agent",
		Type:    "processing",
		Payload: map[string]interface{}{"will_fail": true},
	}

	// Mock task submission with failure
	suite.orchestrator.On("SubmitTask", suite.ctx, mock.MatchedBy(func(t testutils.Task) bool {
		return t.ID == "failing-task"
	})).Return(assert.AnError)

	// Submit failing task
	err := suite.orchestrator.SubmitTask(suite.ctx, failingTask)
	suite.Error(err)

	// Agent should still be registered despite task failure
	agents := suite.orchestrator.GetAgents()
	suite.Contains(agents, "failing-agent")

	suite.orchestrator.AssertExpectations(suite.T())
}

func (suite *OrchestrationTestSuite) TestConfigurationIntegration() {
	// Test integration with configuration management

	// Setup configuration mocks
	suite.configMgr.On("Get", "orchestration.max_agents").Return(10, true)
	suite.configMgr.On("Get", "orchestration.task_timeout").Return(30, true)
	suite.configMgr.On("Get", "orchestration.retry_count").Return(3, true)

	// Simulate configuration-dependent behavior
	maxAgents, exists := suite.configMgr.Get("orchestration.max_agents")
	suite.True(exists)
	suite.Equal(10, maxAgents)

	timeout, exists := suite.configMgr.Get("orchestration.task_timeout")
	suite.True(exists)
	suite.Equal(30, timeout)

	retries, exists := suite.configMgr.Get("orchestration.retry_count")
	suite.True(exists)
	suite.Equal(3, retries)

	suite.configMgr.AssertExpectations(suite.T())
}

func TestOrchestrationSuite(t *testing.T) {
	suite.Run(t, new(OrchestrationTestSuite))
}

// Performance tests for orchestration
func TestOrchestrationPerformance(t *testing.T) {
	t.Run("should handle rapid agent registrations", func(t *testing.T) {
		orchestrator := testutils.NewMockAgentOrchestrator()

		benchmark := testutils.NewPerformanceBenchmark("agent_registration", func() error {
			agent := &testutils.MockAgent{
				ID:           "perf-agent",
				Name:         "Performance Test Agent",
				Capabilities: []string{"test"},
				Status:       "active",
			}
			orchestrator.RegisterAgent(agent)
			return nil
		}).WithIterations(100)

		benchmark.Run(t)
	})

	t.Run("should handle concurrent task submissions", func(t *testing.T) {
		orchestrator := testutils.NewMockAgentOrchestrator()
		ctx := context.Background()

		// Register test agent
		agent := &testutils.MockAgent{
			ID:           "concurrent-agent",
			Name:         "Concurrent Agent",
			Capabilities: []string{"processing"},
			Status:       "active",
		}
		orchestrator.RegisterAgent(agent)

		// Mock multiple task submissions
		orchestrator.On("SubmitTask", ctx, mock.Anything).Return(nil).Times(50)

		benchmark := testutils.NewPerformanceBenchmark("concurrent_tasks", func() error {
			task := testutils.Task{
				ID:      "perf-task",
				AgentID: "concurrent-agent",
				Type:    "processing",
				Payload: map[string]interface{}{"data": "test"},
			}
			return orchestrator.SubmitTask(ctx, task)
		}).WithIterations(50)

		benchmark.Run(t)
		orchestrator.AssertExpectations(t)
	})
}

// Integration tests
func TestOrchestrationIntegration(t *testing.T) {
	t.Run("should handle complete orchestration workflow", func(t *testing.T) {
		orchestrator := testutils.NewMockAgentOrchestrator()
		ctx := context.Background()

		// Step 1: Register agents
		agents := []*testutils.MockAgent{
			{ID: "agent-1", Name: "Agent 1", Capabilities: []string{"task1"}},
			{ID: "agent-2", Name: "Agent 2", Capabilities: []string{"task2"}},
		}

		for _, agent := range agents {
			orchestrator.RegisterAgent(agent)
		}

		// Step 2: Verify agent registration
		registeredAgents := orchestrator.GetAgents()
		require.Len(t, registeredAgents, 2)

		// Step 3: Submit tasks
		tasks := []testutils.Task{
			{ID: "task-1", AgentID: "agent-1", Type: "task1"},
			{ID: "task-2", AgentID: "agent-2", Type: "task2"},
		}

		orchestrator.On("SubmitTask", ctx, mock.Anything).Return(nil).Times(len(tasks))

		for _, task := range tasks {
			err := orchestrator.SubmitTask(ctx, task)
			assert.NoError(t, err)
		}

		// Step 4: Verify task submission
		submittedTasks := orchestrator.GetTasks()
		assert.Len(t, submittedTasks, 2)

		orchestrator.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkAgentRegistration(b *testing.B) {
	orchestrator := testutils.NewMockAgentOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agent := &testutils.MockAgent{
			ID:           "bench-agent",
			Name:         "Benchmark Agent",
			Capabilities: []string{"benchmark"},
			Status:       "active",
		}
		orchestrator.RegisterAgent(agent)
	}
}

func BenchmarkTaskSubmission(b *testing.B) {
	orchestrator := testutils.NewMockAgentOrchestrator()
	ctx := context.Background()

	// Setup
	agent := &testutils.MockAgent{ID: "bench-agent", Status: "active"}
	orchestrator.RegisterAgent(agent)
	orchestrator.On("SubmitTask", ctx, mock.Anything).Return(nil).Times(b.N)

	task := testutils.Task{
		ID:      "bench-task",
		AgentID: "bench-agent",
		Type:    "benchmark",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.SubmitTask(ctx, task)
	}
}
