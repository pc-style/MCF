package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Test functions for the three main components

// Test Configuration Editor functionality
func TestConfigurationEditor() error {
	fmt.Println("Testing Configuration Editor...")

	// Test schema loading
	schemaPath := "config-schema.yaml"
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration schema file not found: %s", schemaPath)
	}

	// Test loading schema
	schema, err := loadEditorSchema(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration schema: %w", err)
	}

	if len(schema.Sections) == 0 {
		return fmt.Errorf("no sections found in schema")
	}

	// Test creating editor model
	_, err = NewEditorModel(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to create editor model: %w", err)
	}

	fmt.Printf("✓ Configuration Editor: Schema loaded with %d sections\n", len(schema.Sections))
	return nil
}

// Test Template Browser functionality
func TestTemplateBrowser() error {
	fmt.Println("Testing Template Browser...")

	// Test creating template browser model
	_ = NewTemplateBrowserModel()

	// Test loading templates (should return built-in templates if directory doesn't exist)
	templates, err := loadTemplatesFromDirectory()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	if len(templates) == 0 {
		return fmt.Errorf("no templates found (should at least have built-in templates)")
	}

	// Test template validation
	for _, template := range templates {
		if template.Name == "" {
			return fmt.Errorf("template found with empty name")
		}
		if template.Type == "" {
			return fmt.Errorf("template '%s' has empty type", template.Name)
		}
	}

	fmt.Printf("✓ Template Browser: Loaded %d templates\n", len(templates))

	// Test individual templates
	builtInTemplates := getBuiltInTemplates()
	if len(builtInTemplates) < 4 {
		return fmt.Errorf("expected at least 4 built-in templates, got %d", len(builtInTemplates))
	}

	fmt.Printf("✓ Template Browser: %d built-in templates available\n", len(builtInTemplates))
	return nil
}

// Test MCF Runner functionality
func TestMCFRunner() error {
	fmt.Println("Testing MCF Runner...")

	wd, _ := os.Getwd()

	// Test creating MCF runner model
	_ = NewMCFRunnerModel(wd)

	// Test operation loading (should handle missing .claude directory gracefully)
	claudeDir := filepath.Join(wd, ".claude")

	// Create a mock .claude directory structure for testing
	agentsDir := filepath.Join(claudeDir, "agents")
	commandsDir := filepath.Join(claudeDir, "commands")

	err := os.MkdirAll(agentsDir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create test agents directory: %w", err)
	}

	err = os.MkdirAll(commandsDir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create test commands directory: %w", err)
	}

	// Create a test agent file
	agentFile := filepath.Join(agentsDir, "test-agent.md")
	err = os.WriteFile(agentFile, []byte("# Test Agent\nThis is a test agent for MCF testing purposes.\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to create test agent file: %w", err)
	}

	// Test loading agents
	agents, err := loadMCFAgents(agentsDir)
	if err != nil {
		// Clean up and return error
		os.RemoveAll(claudeDir)
		return fmt.Errorf("failed to load MCF agents: %w", err)
	}

	if len(agents) == 0 {
		// Clean up and return error
		os.RemoveAll(claudeDir)
		return fmt.Errorf("no agents loaded from test directory")
	}

	// Test agent validation
	testAgent := agents[0]
	if testAgent.name != "test-agent" {
		// Clean up and return error
		os.RemoveAll(claudeDir)
		return fmt.Errorf("expected agent name 'test-agent', got '%s'", testAgent.name)
	}

	if testAgent.opType != MCFOpTypeAgent {
		// Clean up and return error
		os.RemoveAll(claudeDir)
		return fmt.Errorf("expected agent type %v, got %v", MCFOpTypeAgent, testAgent.opType)
	}

	// Clean up test directory
	os.RemoveAll(claudeDir)

	fmt.Printf("✓ MCF Runner: Loaded %d agents from test directory\n", len(agents))
	return nil
}

// Test basic navigation and state management
func TestBasicNavigation() error {
	fmt.Println("Testing Basic Navigation...")

	// Test main model creation
	mainModel := NewMainModel()

	if mainModel.mode != ModeMainMenu {
		return fmt.Errorf("expected main menu mode, got %v", mainModel.mode)
	}

	// Test mode transitions
	choices := mainModel.getAvailableChoices()
	if len(choices) == 0 {
		return fmt.Errorf("no menu choices available")
	}

	// Test initialization state
	initialized := mainModel.isInitialized()
	fmt.Printf("✓ Navigation: MCF initialization state: %v\n", initialized)
	fmt.Printf("✓ Navigation: Available menu choices: %d\n", len(choices))

	return nil
}

// Main test runner
func RunComponentTests() {
	fmt.Println("=== MCF Component Testing ===\n")

	tests := []struct {
		name string
		test func() error
	}{
		{"Configuration Editor", TestConfigurationEditor},
		{"Template Browser", TestTemplateBrowser},
		{"MCF Runner", TestMCFRunner},
		{"Basic Navigation", TestBasicNavigation},
	}

	passed := 0
	failed := 0

	for _, test := range tests {
		fmt.Printf("Running %s tests...\n", test.name)
		if err := test.test(); err != nil {
			fmt.Printf("❌ %s FAILED: %v\n\n", test.name, err)
			failed++
		} else {
			fmt.Printf("✅ %s PASSED\n\n", test.name)
			passed++
		}
	}

	fmt.Printf("=== Test Results ===\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total:  %d\n", passed+failed)

	if failed > 0 {
		fmt.Printf("\n⚠️  Some tests failed. Please review the issues above.\n")
		os.Exit(1)
	} else {
		fmt.Printf("\n🎉 All tests passed successfully!\n")
	}
}
