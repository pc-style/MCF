package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutils "mcf-dev/tui/internal/testing"
)

// ConfigManager handles configuration management for the TUI application
type ConfigManager struct {
	configPath string
	config     map[string]interface{}
	logger     *testutils.TestLogger
}

// DefaultConfig represents the default configuration structure
var DefaultConfig = map[string]interface{}{
	"mcf": map[string]interface{}{
		"host":           "localhost",
		"port":           8080,
		"timeout":        30,
		"retry_attempts": 3,
		"retry_delay":    1000,
		"tls_enabled":    false,
		"api_version":    "v1",
	},
	"tui": map[string]interface{}{
		"theme":           "dark",
		"refresh_rate":    1000,
		"max_log_lines":   1000,
		"auto_scroll":     true,
		"show_timestamps": true,
		"default_view":    "dashboard",
	},
	"logging": map[string]interface{}{
		"level":     "info",
		"file_path": "mcf-tui.log",
		"max_size":  10,
		"max_age":   7,
		"max_files": 3,
		"compress":  true,
	},
	"performance": map[string]interface{}{
		"max_goroutines":     50,
		"cache_size":         100,
		"gc_percent":         100,
		"memory_limit_mb":    512,
		"cpu_limit_percent":  80.0,
		"disk_limit_percent": 90.0,
	},
}

// NewConfigManager creates a new configuration manager instance
func NewConfigManager(configPath string, logger *testutils.TestLogger) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
		config:     make(map[string]interface{}),
		logger:     logger,
	}
}

// Load loads configuration from file
func (c *ConfigManager) Load() error {
	if c.logger != nil {
		c.logger.Log("Loading configuration from %s", c.configPath)
	}

	// Check if config file exists
	if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
		if c.logger != nil {
			c.logger.Log("Config file does not exist, using defaults")
		}
		c.config = DefaultConfig
		return c.Save() // Create default config file
	}

	data, err := os.ReadFile(c.configPath)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("Failed to read config file: %v", err)
		}
		return err
	}

	err = json.Unmarshal(data, &c.config)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("Failed to unmarshal config: %v", err)
		}
		return err
	}

	if c.logger != nil {
		c.logger.Log("Configuration loaded successfully")
	}
	return nil
}

// Save saves configuration to file
func (c *ConfigManager) Save() error {
	if c.logger != nil {
		c.logger.Log("Saving configuration to %s", c.configPath)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if c.logger != nil {
			c.logger.Error("Failed to create config directory: %v", err)
		}
		return err
	}

	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		if c.logger != nil {
			c.logger.Error("Failed to marshal config: %v", err)
		}
		return err
	}

	err = os.WriteFile(c.configPath, data, 0644)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("Failed to write config file: %v", err)
		}
		return err
	}

	if c.logger != nil {
		c.logger.Log("Configuration saved successfully")
	}
	return nil
}

// Get retrieves a configuration value using dot notation (e.g., "mcf.host")
func (c *ConfigManager) Get(key string) (interface{}, bool) {
	return c.getNestedValue(c.config, key)
}

// Set sets a configuration value using dot notation
func (c *ConfigManager) Set(key string, value interface{}) error {
	if c.logger != nil {
		c.logger.Log("Setting config key %s to %v", key, value)
	}

	err := c.setNestedValue(c.config, key, value)
	if err == nil {
		return c.Save()
	}
	return err
}

// GetString retrieves a string configuration value
func (c *ConfigManager) GetString(key string) (string, error) {
	value, exists := c.Get(key)
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}

	return str, nil
}

// GetInt retrieves an integer configuration value
func (c *ConfigManager) GetInt(key string) (int, error) {
	value, exists := c.Get(key)
	if !exists {
		return 0, fmt.Errorf("key %s not found", key)
	}

	// Handle different numeric types from JSON unmarshaling
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("key %s is not a number", key)
	}
}

// GetBool retrieves a boolean configuration value
func (c *ConfigManager) GetBool(key string) (bool, error) {
	value, exists := c.Get(key)
	if !exists {
		return false, fmt.Errorf("key %s not found", key)
	}

	boolean, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("key %s is not a boolean", key)
	}

	return boolean, nil
}

// GetFloat retrieves a float configuration value
func (c *ConfigManager) GetFloat(key string) (float64, error) {
	value, exists := c.Get(key)
	if !exists {
		return 0, fmt.Errorf("key %s not found", key)
	}

	// Handle different numeric types
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("key %s is not a number", key)
	}
}

// Reset resets configuration to defaults
func (c *ConfigManager) Reset() error {
	if c.logger != nil {
		c.logger.Log("Resetting configuration to defaults")
	}

	c.config = make(map[string]interface{})
	for k, v := range DefaultConfig {
		c.config[k] = v
	}

	return c.Save()
}

// Backup creates a backup of the current configuration
func (c *ConfigManager) Backup(backupPath string) error {
	if c.logger != nil {
		c.logger.Log("Creating backup at %s", backupPath)
	}

	// Create backup directory if needed
	dir := filepath.Dir(backupPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, data, 0644)
}

// Validate validates the current configuration
func (c *ConfigManager) Validate() []error {
	var errors []error

	// Required string fields
	requiredStrings := map[string]string{
		"mcf.host":        "MCF host",
		"mcf.api_version": "API version",
		"tui.theme":       "TUI theme",
		"logging.level":   "Log level",
	}

	for key, description := range requiredStrings {
		if _, err := c.GetString(key); err != nil {
			errors = append(errors, fmt.Errorf("%s (%s) is required", description, key))
		}
	}

	// Validate numeric ranges
	port, err := c.GetInt("mcf.port")
	if err == nil {
		if port < 1 || port > 65535 {
			errors = append(errors, fmt.Errorf("mcf.port must be between 1 and 65535"))
		}
	}

	timeout, err := c.GetInt("mcf.timeout")
	if err == nil {
		if timeout < 1 || timeout > 300 {
			errors = append(errors, fmt.Errorf("mcf.timeout must be between 1 and 300 seconds"))
		}
	}

	return errors
}

// Helper methods for nested value operations
func (c *ConfigManager) getNestedValue(config map[string]interface{}, key string) (interface{}, bool) {
	keys := c.splitKey(key)
	current := config

	for i, k := range keys {
		if i == len(keys)-1 {
			value, exists := current[k]
			return value, exists
		}

		next, exists := current[k]
		if !exists {
			return nil, false
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, false
		}

		current = nextMap
	}

	return nil, false
}

func (c *ConfigManager) setNestedValue(config map[string]interface{}, key string, value interface{}) error {
	keys := c.splitKey(key)
	current := config

	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
			return nil
		}

		next, exists := current[k]
		if !exists {
			next = make(map[string]interface{})
			current[k] = next
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot set nested value: %s is not a map", k)
		}

		current = nextMap
	}

	return nil
}

func (c *ConfigManager) splitKey(key string) []string {
	// Simple dot notation split
	result := []string{}
	current := ""

	for _, char := range key {
		if char == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// Test Suite
type ConfigTestSuite struct {
	suite.Suite
	tempDir    string
	configPath string
	manager    *ConfigManager
	logger     *testutils.TestLogger
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.tempDir = suite.T().TempDir()
	suite.configPath = filepath.Join(suite.tempDir, "test-config.json")
	suite.logger = testutils.NewTestLogger(suite.T())
	suite.manager = NewConfigManager(suite.configPath, suite.logger)
}

func TestConfigManagerSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// Unit tests
func TestNewConfigManager(t *testing.T) {
	t.Run("should create new config manager", func(t *testing.T) {
		logger := testutils.NewTestLogger(t)
		manager := NewConfigManager("/tmp/test-config.json", logger)

		assert.NotNil(t, manager)
		assert.Equal(t, "/tmp/test-config.json", manager.configPath)
		assert.NotNil(t, manager.config)
		assert.Equal(t, logger, manager.logger)
	})
}

func TestConfigManager_DefaultConfig(t *testing.T) {
	t.Run("should have valid default configuration", func(t *testing.T) {
		assert.NotNil(t, DefaultConfig)
		assert.Contains(t, DefaultConfig, "mcf")
		assert.Contains(t, DefaultConfig, "tui")
		assert.Contains(t, DefaultConfig, "logging")
		assert.Contains(t, DefaultConfig, "performance")

		// Validate MCF section
		mcf := DefaultConfig["mcf"].(map[string]interface{})
		assert.Equal(t, "localhost", mcf["host"])
		assert.Equal(t, 8080, mcf["port"])
		assert.Equal(t, "v1", mcf["api_version"])

		// Validate TUI section
		tui := DefaultConfig["tui"].(map[string]interface{})
		assert.Equal(t, "dark", tui["theme"])
		assert.Equal(t, "dashboard", tui["default_view"])

		// Validate performance section with float values
		perf := DefaultConfig["performance"].(map[string]interface{})
		assert.Equal(t, 80.0, perf["cpu_limit_percent"])
		assert.Equal(t, 90.0, perf["disk_limit_percent"])
	})
}

func (suite *ConfigTestSuite) TestLoad() {
	suite.Run("should load default config when file doesn't exist", func() {
		err := suite.manager.Load()
		suite.NoError(err)

		// Should create the file
		_, err = os.Stat(suite.configPath)
		suite.NoError(err)

		// Should have default values
		host, exists := suite.manager.Get("mcf.host")
		suite.True(exists)
		suite.Equal("localhost", host)
	})

	suite.Run("should load existing config file", func() {
		// Create a test config file
		testConfig := map[string]interface{}{
			"mcf": map[string]interface{}{
				"host": "test-host",
				"port": 9000,
			},
		}

		data, err := json.MarshalIndent(testConfig, "", "  ")
		suite.NoError(err)

		err = os.WriteFile(suite.configPath, data, 0644)
		suite.NoError(err)

		// Load the config
		err = suite.manager.Load()
		suite.NoError(err)

		// Verify values
		host, exists := suite.manager.Get("mcf.host")
		suite.True(exists)
		suite.Equal("test-host", host)

		port, exists := suite.manager.Get("mcf.port")
		suite.True(exists)
		suite.Equal(float64(9000), port) // JSON unmarshaling gives float64
	})
}

func (suite *ConfigTestSuite) TestSave() {
	suite.Run("should save config to file", func() {
		suite.manager.config = DefaultConfig

		err := suite.manager.Save()
		suite.NoError(err)

		// Verify file was created
		_, err = os.Stat(suite.configPath)
		suite.NoError(err)

		// Verify content
		data, err := os.ReadFile(suite.configPath)
		suite.NoError(err)

		var loadedConfig map[string]interface{}
		err = json.Unmarshal(data, &loadedConfig)
		suite.NoError(err)

		suite.Contains(loadedConfig, "mcf")
		suite.Contains(loadedConfig, "tui")
	})
}

func (suite *ConfigTestSuite) TestGetSet() {
	suite.Run("should get and set values using dot notation", func() {
		testCases := []struct {
			key   string
			value interface{}
		}{
			{"mcf.host", "test-host"},
			{"mcf.port", 8080},
			{"mcf.tls_enabled", true},
			{"performance.cpu_limit_percent", 85.5},
			{"new.nested.value", "created"},
		}

		for _, tc := range testCases {
			// Set value
			err := suite.manager.Set(tc.key, tc.value)
			suite.NoError(err, "Failed to set %s", tc.key)

			// Get value
			retrievedValue, exists := suite.manager.Get(tc.key)
			suite.True(exists, "Key %s should exist", tc.key)
			if tc.key == "performance.cpu_limit_percent" {
				// Handle float comparison
				suite.InDelta(tc.value, retrievedValue, 0.01, "Value should match for %s", tc.key)
			} else {
				suite.Equal(tc.value, retrievedValue, "Value should match for %s", tc.key)
			}
		}
	})
}

func (suite *ConfigTestSuite) TestTypedGetters() {
	suite.Run("should retrieve typed values correctly", func() {
		// Set test values
		suite.manager.Set("test.string", "hello")
		suite.manager.Set("test.int", 42)
		suite.manager.Set("test.float", 3.14)
		suite.manager.Set("test.bool", true)

		// Test string getter
		str, err := suite.manager.GetString("test.string")
		suite.NoError(err)
		suite.Equal("hello", str)

		// Test int getter
		intVal, err := suite.manager.GetInt("test.int")
		suite.NoError(err)
		suite.Equal(42, intVal)

		// Test float getter
		floatVal, err := suite.manager.GetFloat("test.float")
		suite.NoError(err)
		suite.InDelta(3.14, floatVal, 0.01)

		// Test bool getter
		boolVal, err := suite.manager.GetBool("test.bool")
		suite.NoError(err)
		suite.True(boolVal)
	})
}

func (suite *ConfigTestSuite) TestValidation() {
	suite.Run("should validate configuration", func() {
		// Start with default config
		suite.manager.config = make(map[string]interface{})
		for k, v := range DefaultConfig {
			suite.manager.config[k] = v
		}

		errors := suite.manager.Validate()
		suite.Empty(errors, "Default config should be valid")

		// Test invalid port
		suite.manager.Set("mcf.port", 99999)
		errors = suite.manager.Validate()
		suite.NotEmpty(errors)

		// Test missing required field
		delete(suite.manager.config["mcf"].(map[string]interface{}), "host")
		errors = suite.manager.Validate()
		suite.NotEmpty(errors)
	})
}

func (suite *ConfigTestSuite) TestReset() {
	suite.Run("should reset to default configuration", func() {
		// Set some custom values
		suite.manager.Set("mcf.host", "custom-host")
		suite.manager.Set("custom.setting", "custom-value")

		// Reset
		err := suite.manager.Reset()
		suite.NoError(err)

		// Verify reset to defaults
		host, exists := suite.manager.Get("mcf.host")
		suite.True(exists)
		suite.Equal("localhost", host)

		// Custom setting should be gone
		_, exists = suite.manager.Get("custom.setting")
		suite.False(exists)
	})
}

func (suite *ConfigTestSuite) TestBackup() {
	suite.Run("should create backup of configuration", func() {
		// Set some config
		suite.manager.config = DefaultConfig

		backupPath := filepath.Join(suite.tempDir, "backup.json")
		err := suite.manager.Backup(backupPath)
		suite.NoError(err)

		// Verify backup file exists
		_, err = os.Stat(backupPath)
		suite.NoError(err)

		// Verify backup content
		data, err := os.ReadFile(backupPath)
		suite.NoError(err)

		var backupConfig map[string]interface{}
		err = json.Unmarshal(data, &backupConfig)
		suite.NoError(err)

		suite.Contains(backupConfig, "mcf")
		suite.Contains(backupConfig, "tui")
	})
}

// Integration test
func TestConfigManager_Integration(t *testing.T) {
	t.Run("should handle complete configuration lifecycle", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "integration-config.json")
		logger := testutils.NewTestLogger(t)

		// Step 1: Create and load manager
		manager := NewConfigManager(configPath, logger)
		err := manager.Load()
		assert.NoError(t, err)

		// Step 2: Modify configuration
		err = manager.Set("mcf.host", "integration-host")
		assert.NoError(t, err)

		err = manager.Set("mcf.port", 8888)
		assert.NoError(t, err)

		err = manager.Set("custom.setting", "custom-value")
		assert.NoError(t, err)

		// Step 3: Create new manager and load same config
		manager2 := NewConfigManager(configPath, logger)
		err = manager2.Load()
		assert.NoError(t, err)

		// Step 4: Verify configuration was preserved
		host, exists := manager2.Get("mcf.host")
		assert.True(t, exists)
		assert.Equal(t, "integration-host", host)

		port, exists := manager2.Get("mcf.port")
		assert.True(t, exists)
		assert.Equal(t, float64(8888), port) // JSON unmarshaling gives float64

		custom, exists := manager2.Get("custom.setting")
		assert.True(t, exists)
		assert.Equal(t, "custom-value", custom)

		// Step 5: Test validation
		errors := manager2.Validate()
		assert.Empty(t, errors)

		// Step 6: Test backup
		backupPath := filepath.Join(tempDir, "integration-backup.json")
		err = manager2.Backup(backupPath)
		assert.NoError(t, err)

		_, err = os.Stat(backupPath)
		assert.NoError(t, err)

		// Step 7: Reset and verify
		err = manager2.Reset()
		assert.NoError(t, err)

		host, exists = manager2.Get("mcf.host")
		assert.True(t, exists)
		assert.Equal(t, "localhost", host) // Back to default
	})
}

// Performance tests
func TestConfigManager_Performance(t *testing.T) {
	t.Run("should handle rapid configuration operations", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "perf-config.json")
		logger := testutils.NewTestLogger(t)
		manager := NewConfigManager(configPath, logger)

		err := manager.Load()
		require.NoError(t, err)

		benchmark := testutils.NewPerformanceBenchmark("config_operations", func() error {
			// Rapid set/get operations
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("test.key%d", i)
				value := fmt.Sprintf("value%d", i)

				err := manager.Set(key, value)
				if err != nil {
					return err
				}

				_, exists := manager.Get(key)
				if !exists {
					return fmt.Errorf("key %s should exist", key)
				}
			}
			return nil
		}).WithIterations(100)

		benchmark.Run(t)
	})

	t.Run("should handle large configurations efficiently", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "large-config.json")
		logger := testutils.NewTestLogger(t)
		manager := NewConfigManager(configPath, logger)

		// Create large configuration
		for i := 0; i < 1000; i++ {
			section := fmt.Sprintf("section%d", i)
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("%s.key%d", section, j)
				value := fmt.Sprintf("value-%d-%d", i, j)
				manager.Set(key, value)
			}
		}

		benchmark := testutils.NewPerformanceBenchmark("large_config_access", func() error {
			// Random access to configuration values
			for i := 0; i < 50; i++ {
				key := fmt.Sprintf("section%d.key%d", i%100, i%10)
				_, exists := manager.Get(key)
				if !exists {
					return fmt.Errorf("key %s should exist", key)
				}
			}
			return nil
		}).WithIterations(100)

		benchmark.Run(t)
	})
}

// Benchmark tests
func BenchmarkConfigManager_Set(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "bench-config.json")
	logger := testutils.NewTestLogger(&testing.T{})
	manager := NewConfigManager(configPath, logger)

	manager.Load()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("benchmark.key%d", i%100)
		value := fmt.Sprintf("value%d", i)
		manager.Set(key, value)
	}
}

func BenchmarkConfigManager_Get(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "bench-config.json")
	logger := testutils.NewTestLogger(&testing.T{})
	manager := NewConfigManager(configPath, logger)

	manager.Load()

	// Pre-populate with data
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("benchmark.key%d", i)
		value := fmt.Sprintf("value%d", i)
		manager.Set(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("benchmark.key%d", i%100)
		manager.Get(key)
	}
}
