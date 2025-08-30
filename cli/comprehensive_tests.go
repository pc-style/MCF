package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Comprehensive tests for edge cases and error conditions
func RunComprehensiveTests() {
	fmt.Println("=== MCF Comprehensive Testing ===\n")

	tests := []struct {
		name string
		test func() error
	}{
		{"Configuration Schema Validation", TestConfigSchemaValidation},
		{"Template Parameter Validation", TestTemplateParameterValidation},
		{"MCF Operation Error Handling", TestMCFOperationErrorHandling},
		{"State Management", TestStateManagement},
		{"Message Bus Communication", TestMessageBusCommunication},
		{"File System Operations", TestFileSystemOperations},
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

	fmt.Printf("=== Comprehensive Test Results ===\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total:  %d\n", passed+failed)

	if failed > 0 {
		fmt.Printf("\n⚠️  Some comprehensive tests failed. Please review the issues above.\n")
		os.Exit(1)
	} else {
		fmt.Printf("\n🎉 All comprehensive tests passed successfully!\n")
	}
}

// Test configuration schema validation
func TestConfigSchemaValidation() error {
	fmt.Println("Testing configuration schema validation...")

	// Test valid schema loading
	schemaPath := "config-schema.yaml"
	schema, err := loadEditorSchema(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to load valid schema: %w", err)
	}

	// Validate schema structure
	if len(schema.Sections) == 0 {
		return fmt.Errorf("schema should have sections")
	}

	// Check that each section has required fields
	for _, section := range schema.Sections {
		if section.Name == "" {
			return fmt.Errorf("section found with empty name")
		}
		if section.Description == "" {
			return fmt.Errorf("section '%s' has empty description", section.Name)
		}
	}

	// Test invalid schema file
	invalidPath := "non-existent-schema.yaml"
	_, err = loadEditorSchema(invalidPath)
	if err == nil {
		return fmt.Errorf("should fail when loading non-existent schema file")
	}

	fmt.Printf("✓ Schema validation: Valid schema loaded with %d sections\n", len(schema.Sections))
	fmt.Printf("✓ Schema validation: Error handling works for invalid paths\n")

	return nil
}

// Test template parameter validation
func TestTemplateParameterValidation() error {
	fmt.Println("Testing template parameter validation...")

	templates := getBuiltInTemplates()
	if len(templates) == 0 {
		return fmt.Errorf("should have built-in templates")
	}

	// Test each template for valid structure
	for _, template := range templates {
		// Validate basic template fields
		if template.Name == "" {
			return fmt.Errorf("template found with empty name")
		}
		if template.Type == "" {
			return fmt.Errorf("template '%s' has empty type", template.Name)
		}
		if template.Description == "" {
			return fmt.Errorf("template '%s' has empty description", template.Name)
		}

		// Validate parameters
		for _, param := range template.Parameters {
			if param.Key == "" {
				return fmt.Errorf("template '%s' has parameter with empty key", template.Name)
			}
			if param.Label == "" {
				return fmt.Errorf("template '%s' has parameter '%s' with empty label", template.Name, param.Key)
			}
			if param.Type == "" {
				return fmt.Errorf("template '%s' has parameter '%s' with empty type", template.Name, param.Key)
			}

			// Validate parameter types
			validTypes := []string{"text", "boolean", "select", "password", "number"}
			validType := false
			for _, vt := range validTypes {
				if param.Type == vt {
					validType = true
					break
				}
			}
			if !validType {
				return fmt.Errorf("template '%s' has parameter '%s' with invalid type '%s'", template.Name, param.Key, param.Type)
			}

			// If type is select, should have options
			if param.Type == "select" && len(param.Options) == 0 {
				return fmt.Errorf("template '%s' has select parameter '%s' without options", template.Name, param.Key)
			}
		}

		// Validate files structure
		for _, file := range template.Files {
			if file.Path == "" {
				return fmt.Errorf("template '%s' has file with empty path", template.Name)
			}
			validFileTypes := []string{"file", "directory", "symlink"}
			validFileType := false
			for _, ft := range validFileTypes {
				if file.Type == ft {
					validFileType = true
					break
				}
			}
			if !validFileType {
				return fmt.Errorf("template '%s' has file '%s' with invalid type '%s'", template.Name, file.Path, file.Type)
			}
		}
	}

	// Test template loading from directory
	loadedTemplates, err := loadTemplatesFromDirectory()
	if err != nil {
		return fmt.Errorf("failed to load templates from directory: %w", err)
	}

	// Should include built-in templates
	if len(loadedTemplates) < len(templates) {
		return fmt.Errorf("loaded templates (%d) should include all built-in templates (%d)", len(loadedTemplates), len(templates))
	}

	fmt.Printf("✓ Template validation: All %d templates have valid structure\n", len(templates))
	fmt.Printf("✓ Template validation: Template loading works correctly\n")

	return nil
}

// Test MCF operation error handling
func TestMCFOperationErrorHandling() error {
	fmt.Println("Testing MCF operation error handling...")

	// Test loading operations from non-existent directory
	nonExistentPath := "/non/existent/path"
	agents, err := loadMCFAgents(nonExistentPath)
	if err == nil {
		return fmt.Errorf("should return error for non-existent agents directory")
	}
	if len(agents) != 0 {
		return fmt.Errorf("should return empty agents list for non-existent directory")
	}

	commands, err := loadMCFCommands(nonExistentPath)
	if err != nil && len(commands) != 0 {
		return fmt.Errorf("should handle non-existent commands directory gracefully")
	}

	// Test with valid but empty directory
	tempDir, err := os.MkdirTemp("", "mcf_test_")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	emptyAgents, err := loadMCFAgents(tempDir)
	if err != nil {
		return fmt.Errorf("should handle empty directory gracefully: %w", err)
	}
	if len(emptyAgents) != 0 {
		return fmt.Errorf("should return empty list for empty directory")
	}

	// Test MCF file description parsing
	testFilePath := filepath.Join(tempDir, "test.md")
	testContent := `# Test Agent
This is a test agent description.

Some more content here.`

	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}

	description := getMCFFileDescription(testFilePath)
	if description == "No description available" {
		return fmt.Errorf("should extract description from test file")
	}

	fmt.Printf("✓ MCF operations: Error handling works for missing directories\n")
	fmt.Printf("✓ MCF operations: File description parsing works\n")

	return nil
}

// Test state management
func TestStateManagement() error {
	fmt.Println("Testing state management...")

	// Test application state creation
	state := NewApplicationState()
	if state == nil {
		return fmt.Errorf("failed to create application state")
	}

	// Test mode management
	initialMode := state.GetCurrentMode()
	if initialMode != ModeMainMenu {
		return fmt.Errorf("expected initial mode to be ModeMainMenu, got %v", initialMode)
	}

	// Test mode transition
	state.SetMode(ModeConfigurator)
	currentMode := state.GetCurrentMode()
	if currentMode != ModeConfigurator {
		return fmt.Errorf("expected mode to be ModeConfigurator, got %v", currentMode)
	}

	previousMode := state.GetPreviousMode()
	if previousMode != initialMode {
		return fmt.Errorf("expected previous mode to be %v, got %v", initialMode, previousMode)
	}

	// Test configuration state
	if state.IsConfigurationLoaded() {
		return fmt.Errorf("configuration should not be loaded initially")
	}

	state.SetConfigurationLoaded(true, nil)
	if !state.IsConfigurationLoaded() {
		return fmt.Errorf("configuration should be loaded after setting")
	}

	// Test notifications
	notification := Notification{
		Type:    "success",
		Title:   "Test",
		Message: "Test notification",
	}

	state.AddNotification(notification)
	notifications := state.GetNotifications()
	if len(notifications) == 0 {
		return fmt.Errorf("should have notifications after adding")
	}

	// Test health status
	if !state.IsHealthy() {
		return fmt.Errorf("state should be healthy initially")
	}

	healthStatus := state.GetHealthStatus()
	if healthStatus["healthy"] != true {
		return fmt.Errorf("health status should indicate healthy state")
	}

	fmt.Printf("✓ State management: Mode transitions work correctly\n")
	fmt.Printf("✓ State management: Configuration state tracking works\n")
	fmt.Printf("✓ State management: Notifications system works\n")
	fmt.Printf("✓ State management: Health monitoring works\n")

	return nil
}

// Test message bus communication
func TestMessageBusCommunication() error {
	fmt.Println("Testing message bus communication...")

	// Test message bus creation
	bus := NewMessageBus()
	if bus == nil {
		return fmt.Errorf("failed to create message bus")
	}

	// Test message creation
	appMsg := NewAppInitMessage("1.0.0", map[string]interface{}{"test": true})
	if appMsg.Type() != MsgAppInit {
		return fmt.Errorf("expected message type %v, got %v", MsgAppInit, appMsg.Type())
	}

	configMsg := NewConfigLoadedMessage(true, nil, "")
	if configMsg.Type() != MsgConfigLoaded {
		return fmt.Errorf("expected message type %v, got %v", MsgConfigLoaded, configMsg.Type())
	}

	installMsg := NewInstallProgressMessage(0.5, "Testing", 1)
	if installMsg.Type() != MsgInstallProgress {
		return fmt.Errorf("expected message type %v, got %v", MsgInstallProgress, installMsg.Type())
	}

	uiMsg := NewUIErrorMessage("Test Error", "Test error message", nil)
	if uiMsg.Type() != MsgUIError {
		return fmt.Errorf("expected message type %v, got %v", MsgUIError, uiMsg.Type())
	}

	// Test message timestamps
	if appMsg.Timestamp().IsZero() {
		return fmt.Errorf("message should have timestamp")
	}

	// Test subscription and publishing
	msgChan := make(chan MCFMessage, 1)
	bus.Subscribe(MsgAppInit, msgChan)

	bus.Publish(appMsg)

	// Check if message was received (non-blocking)
	select {
	case receivedMsg := <-msgChan:
		if receivedMsg.Type() != MsgAppInit {
			return fmt.Errorf("received wrong message type")
		}
	default:
		return fmt.Errorf("message was not received through subscription")
	}

	fmt.Printf("✓ Message bus: Message creation works correctly\n")
	fmt.Printf("✓ Message bus: Subscription and publishing works\n")
	fmt.Printf("✓ Message bus: Message timestamps are set correctly\n")

	return nil
}

// Test file system operations
func TestFileSystemOperations() error {
	fmt.Println("Testing file system operations...")

	// Test with temp directory
	tempDir, err := os.MkdirTemp("", "mcf_fs_test_")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Test template loading from empty directory
	templatesDir := filepath.Join(tempDir, ".claude", "templates")
	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Create a test template file
	templateFile := filepath.Join(templatesDir, "test-template.yaml")
	templateContent := `name: "Test Template"
type: "test"
version: "1.0.0"
description: "A test template"
author: "Test Author"
parameters:
  - key: "test_param"
    label: "Test Parameter"
    type: "text"
    description: "A test parameter"
    required: true
files:
  - path: "test.txt"
    type: "file"
    content: "Test content"
`

	err = os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	// Test loading the template
	template, err := loadTemplate(templateFile)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	if template.Name != "Test Template" {
		return fmt.Errorf("expected template name 'Test Template', got '%s'", template.Name)
	}

	if len(template.Parameters) == 0 {
		return fmt.Errorf("template should have parameters")
	}

	if len(template.Files) == 0 {
		return fmt.Errorf("template should have files")
	}

	// Test invalid template file
	invalidTemplateFile := filepath.Join(templatesDir, "invalid.yaml")
	invalidContent := "invalid: yaml: content: ["
	err = os.WriteFile(invalidTemplateFile, []byte(invalidContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write invalid template file: %w", err)
	}

	_, err = loadTemplate(invalidTemplateFile)
	if err == nil {
		return fmt.Errorf("should fail when loading invalid template file")
	}

	// Test directory walking
	// Change to temp directory temporarily to test template loading
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	defer os.Chdir(originalWd)

	templates, err := loadTemplatesFromDirectory()
	if err != nil {
		return fmt.Errorf("failed to load templates from test directory: %w", err)
	}

	// Should have built-in templates plus our test template
	if len(templates) < 5 { // 4 built-in + 1 test template
		return fmt.Errorf("expected at least 5 templates, got %d", len(templates))
	}

	fmt.Printf("✓ File system: Template file loading works correctly\n")
	fmt.Printf("✓ File system: Invalid file handling works correctly\n")
	fmt.Printf("✓ File system: Directory walking works correctly\n")

	return nil
}
