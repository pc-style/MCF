package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Configuration management with layered approach
type ConfigManager struct {
	globalPath  string
	projectPath string
	localPath   string
	config      *LayeredConfig
}

type LayeredConfig struct {
	Global  GlobalConfig  `json:"global" yaml:"global"`
	Project ProjectConfig `json:"project" yaml:"project"`
	Local   LocalConfig   `json:"local" yaml:"local"`
}

type GlobalConfig struct {
	Version     string                 `json:"version" yaml:"version"`
	Preferences UserPreferences        `json:"preferences" yaml:"preferences"`
	Defaults    map[string]interface{} `json:"defaults" yaml:"defaults"`
	UpdatedAt   time.Time              `json:"updated_at" yaml:"updated_at"`
}

type ProjectConfig struct {
	Name         string                 `json:"name" yaml:"name"`
	Type         string                 `json:"type" yaml:"type"`
	Features     []string               `json:"features" yaml:"features"`
	Integrations map[string]bool        `json:"integrations" yaml:"integrations"`
	Hooks        []string               `json:"hooks" yaml:"hooks"`
	Settings     map[string]interface{} `json:"settings" yaml:"settings"`
	CreatedAt    time.Time              `json:"created_at" yaml:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" yaml:"updated_at"`
}

type LocalConfig struct {
	Developer   DeveloperInfo          `json:"developer" yaml:"developer"`
	Environment map[string]string      `json:"environment" yaml:"environment"`
	Overrides   map[string]interface{} `json:"overrides" yaml:"overrides"`
	LastUsed    time.Time              `json:"last_used" yaml:"last_used"`
}

type DeveloperInfo struct {
	Name         string `json:"name" yaml:"name"`
	Email        string `json:"email" yaml:"email"`
	PreferredIDE string `json:"preferred_ide" yaml:"preferred_ide"`
	WorkingHours string `json:"working_hours" yaml:"working_hours"`
}

func NewConfigManager(projectPath string) *ConfigManager {
	homeDir, _ := os.UserHomeDir()

	return &ConfigManager{
		globalPath:  filepath.Join(homeDir, ".mcf", "config.yaml"),
		projectPath: filepath.Join(projectPath, ".claude", "config.yaml"),
		localPath:   filepath.Join(projectPath, ".claude", "local.yaml"),
		config:      &LayeredConfig{},
	}
}

func (cm *ConfigManager) Load() error {
	// Initialize with defaults
	cm.config = &LayeredConfig{
		Global:  cm.getGlobalDefaults(),
		Project: cm.getProjectDefaults(),
		Local:   cm.getLocalDefaults(),
	}

	// Load global config
	if err := cm.loadGlobal(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Load project config
	if err := cm.loadProject(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Load local config (developer-specific)
	if err := cm.loadLocal(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load local config: %w", err)
	}

	return nil
}

func (cm *ConfigManager) Save() error {
	// Save global config
	if err := cm.saveGlobal(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}

	// Save project config
	if err := cm.saveProject(); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	// Save local config
	if err := cm.saveLocal(); err != nil {
		return fmt.Errorf("failed to save local config: %w", err)
	}

	return nil
}

func (cm *ConfigManager) loadGlobal() error {
	data, err := os.ReadFile(cm.globalPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &cm.config.Global)
}

func (cm *ConfigManager) loadProject() error {
	data, err := os.ReadFile(cm.projectPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &cm.config.Project)
}

func (cm *ConfigManager) loadLocal() error {
	data, err := os.ReadFile(cm.localPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &cm.config.Local)
}

func (cm *ConfigManager) saveGlobal() error {
	cm.config.Global.UpdatedAt = time.Now()
	return cm.saveYAMLFile(cm.globalPath, cm.config.Global)
}

func (cm *ConfigManager) saveProject() error {
	cm.config.Project.UpdatedAt = time.Now()
	return cm.saveYAMLFile(cm.projectPath, cm.config.Project)
}

func (cm *ConfigManager) saveLocal() error {
	cm.config.Local.LastUsed = time.Now()
	return cm.saveYAMLFile(cm.localPath, cm.config.Local)
}

func (cm *ConfigManager) saveYAMLFile(path string, data interface{}) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, yamlData, 0644)
}

func (cm *ConfigManager) getGlobalDefaults() GlobalConfig {
	return GlobalConfig{
		Version: "1.0.0",
		Preferences: UserPreferences{
			Theme:        "default",
			Editor:       "auto-detect",
			AutoUpdate:   true,
			TelemetryOpt: false,
		},
		Defaults: map[string]interface{}{
			"timeout":        "30s",
			"retry_attempts": 3,
			"log_level":      "info",
		},
		UpdatedAt: time.Now(),
	}
}

func (cm *ConfigManager) getProjectDefaults() ProjectConfig {
	return ProjectConfig{
		Type:     "general",
		Features: []string{"core"},
		Integrations: map[string]bool{
			"git":      true,
			"serena":   false,
			"security": true,
		},
		Hooks:     []string{},
		Settings:  map[string]interface{}{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (cm *ConfigManager) getLocalDefaults() LocalConfig {
	return LocalConfig{
		Developer: DeveloperInfo{
			PreferredIDE: "auto-detect",
			WorkingHours: "9-17",
		},
		Environment: map[string]string{},
		Overrides:   map[string]interface{}{},
		LastUsed:    time.Now(),
	}
}

// Field validation framework
type FieldValidator struct {
	validators []ValidatorFunc
	sanitizers []SanitizerFunc
}

type ValidatorFunc func(string) error
type SanitizerFunc func(string) string

func NewFieldValidator() *FieldValidator {
	return &FieldValidator{
		validators: []ValidatorFunc{},
		sanitizers: []SanitizerFunc{TrimSpaceSanitizer()},
	}
}

func (fv *FieldValidator) AddValidator(validator ValidatorFunc) *FieldValidator {
	fv.validators = append(fv.validators, validator)
	return fv
}

func (fv *FieldValidator) AddSanitizer(sanitizer SanitizerFunc) *FieldValidator {
	fv.sanitizers = append(fv.sanitizers, sanitizer)
	return fv
}

func (fv *FieldValidator) Validate(value string) (string, error) {
	// Apply sanitizers first
	sanitized := value
	for _, sanitizer := range fv.sanitizers {
		sanitized = sanitizer(sanitized)
	}

	// Run validators
	for _, validator := range fv.validators {
		if err := validator(sanitized); err != nil {
			return sanitized, err
		}
	}

	return sanitized, nil
}

// Common validators
func RequiredValidator() ValidatorFunc {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("this field is required")
		}
		return nil
	}
}

func EmailValidator() ValidatorFunc {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return func(value string) error {
		if !emailRegex.MatchString(value) {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}
}

func ProjectNameValidator() ValidatorFunc {
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$`)
	return func(value string) error {
		if len(value) < 1 {
			return fmt.Errorf("project name cannot be empty")
		}
		if len(value) > 50 {
			return fmt.Errorf("project name cannot exceed 50 characters")
		}
		if !nameRegex.MatchString(value) {
			return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
		}
		return nil
	}
}

func PathValidator() ValidatorFunc {
	return func(value string) error {
		expanded := os.ExpandEnv(value)
		if !filepath.IsAbs(expanded) {
			if _, err := filepath.Abs(expanded); err != nil {
				return fmt.Errorf("invalid path: %w", err)
			}
		}
		return nil
	}
}

func DirectoryExistsValidator() ValidatorFunc {
	return func(value string) error {
		expanded := os.ExpandEnv(value)
		info, err := os.Stat(expanded)
		if err != nil {
			return fmt.Errorf("directory does not exist: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path is not a directory")
		}
		return nil
	}
}

func EnumValidator(validValues []string) ValidatorFunc {
	return func(value string) error {
		for _, valid := range validValues {
			if value == valid {
				return nil
			}
		}
		return fmt.Errorf("value must be one of: %s", strings.Join(validValues, ", "))
	}
}

// Common sanitizers
func TrimSpaceSanitizer() SanitizerFunc {
	return strings.TrimSpace
}

func LowercaseSanitizer() SanitizerFunc {
	return strings.ToLower
}

func ExpandPathSanitizer() SanitizerFunc {
	return func(value string) string {
		return os.ExpandEnv(value)
	}
}

// Configuration schema for dynamic forms
type ConfigSchema struct {
	Sections []SectionSchema `json:"sections" yaml:"sections"`
}

type SectionSchema struct {
	Name         string             `json:"name" yaml:"name"`
	Description  string             `json:"description" yaml:"description"`
	Fields       []FieldSchema      `json:"fields" yaml:"fields"`
	Dependencies []string           `json:"dependencies" yaml:"dependencies"`
	Optional     bool               `json:"optional" yaml:"optional"`
	Conditionals []ConditionalLogic `json:"conditionals" yaml:"conditionals"`
}

type FieldSchema struct {
	Key         string          `json:"key" yaml:"key"`
	Label       string          `json:"label" yaml:"label"`
	Type        FieldType       `json:"type" yaml:"type"`
	Required    bool            `json:"required" yaml:"required"`
	Description string          `json:"description" yaml:"description"`
	Default     interface{}     `json:"default" yaml:"default"`
	Validation  ValidationRules `json:"validation" yaml:"validation"`
	Options     []FieldOption   `json:"options" yaml:"options"`
	Placeholder string          `json:"placeholder" yaml:"placeholder"`
	HelpText    string          `json:"help_text" yaml:"help_text"`
	Sensitive   bool            `json:"sensitive" yaml:"sensitive"`
}

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeEmail    FieldType = "email"
	FieldTypePassword FieldType = "password"
	FieldTypePath     FieldType = "path"
	FieldTypeSelect   FieldType = "select"
	FieldTypeMulti    FieldType = "multi_select"
	FieldTypeBool     FieldType = "boolean"
	FieldTypeNumber   FieldType = "number"
)

type ValidationRules struct {
	Pattern  string   `json:"pattern" yaml:"pattern"`
	MinLen   *int     `json:"min_length" yaml:"min_length"`
	MaxLen   *int     `json:"max_length" yaml:"max_length"`
	Min      *float64 `json:"min" yaml:"min"`
	Max      *float64 `json:"max" yaml:"max"`
	Enum     []string `json:"enum" yaml:"enum"`
	Required bool     `json:"required" yaml:"required"`
}

type FieldOption struct {
	Value       string `json:"value" yaml:"value"`
	Label       string `json:"label" yaml:"label"`
	Description string `json:"description" yaml:"description"`
	Default     bool   `json:"default" yaml:"default"`
}

type ConditionalLogic struct {
	Field     string      `json:"field" yaml:"field"`
	Operation string      `json:"operation" yaml:"operation"`
	Value     interface{} `json:"value" yaml:"value"`
	Action    string      `json:"action" yaml:"action"`
}

// Load configuration schema from file
func LoadConfigSchema(schemaPath string) (*ConfigSchema, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var schema ConfigSchema
	if strings.HasSuffix(schemaPath, ".yaml") || strings.HasSuffix(schemaPath, ".yml") {
		err = yaml.Unmarshal(data, &schema)
	} else {
		err = json.Unmarshal(data, &schema)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return &schema, nil
}

// Create validator from validation rules
func CreateValidator(rules ValidationRules) *FieldValidator {
	validator := NewFieldValidator()

	if rules.Required {
		validator.AddValidator(RequiredValidator())
	}

	if rules.Pattern != "" {
		pattern := regexp.MustCompile(rules.Pattern)
		validator.AddValidator(func(value string) error {
			if !pattern.MatchString(value) {
				return fmt.Errorf("value does not match required pattern")
			}
			return nil
		})
	}

	if rules.MinLen != nil {
		minLen := *rules.MinLen
		validator.AddValidator(func(value string) error {
			if len(value) < minLen {
				return fmt.Errorf("value must be at least %d characters", minLen)
			}
			return nil
		})
	}

	if rules.MaxLen != nil {
		maxLen := *rules.MaxLen
		validator.AddValidator(func(value string) error {
			if len(value) > maxLen {
				return fmt.Errorf("value must be at most %d characters", maxLen)
			}
			return nil
		})
	}

	if len(rules.Enum) > 0 {
		validator.AddValidator(EnumValidator(rules.Enum))
	}

	return validator
}

// State persistence with migration support
type StatePersister struct {
	configPath string
	backupPath string
	migrations []StateMigration
}

type StateMigration struct {
	Version int
	Name    string
	Up      func(map[string]interface{}) error
	Down    func(map[string]interface{}) error
}

func NewStatePersister(configPath string) *StatePersister {
	return &StatePersister{
		configPath: configPath,
		backupPath: configPath + ".backup",
		migrations: []StateMigration{},
	}
}

func (sp *StatePersister) AddMigration(migration StateMigration) {
	sp.migrations = append(sp.migrations, migration)
}

func (sp *StatePersister) LoadWithMigration() (map[string]interface{}, error) {
	// Load current state
	state := make(map[string]interface{})

	if _, err := os.Stat(sp.configPath); os.IsNotExist(err) {
		// No config file exists, return empty state
		return state, nil
	}

	data, err := os.ReadFile(sp.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Run migrations
	if err := sp.runMigrations(state); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return state, nil
}

func (sp *StatePersister) runMigrations(state map[string]interface{}) error {
	currentVersion := sp.getCurrentVersion(state)

	for _, migration := range sp.migrations {
		if migration.Version > currentVersion {
			// Create backup before migration
			if err := sp.createBackup(); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}

			if err := migration.Up(state); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Name, err)
			}

			sp.setVersion(state, migration.Version)
		}
	}

	return nil
}

func (sp *StatePersister) getCurrentVersion(state map[string]interface{}) int {
	if version, exists := state["_version"]; exists {
		if v, ok := version.(float64); ok {
			return int(v)
		}
	}
	return 0
}

func (sp *StatePersister) setVersion(state map[string]interface{}, version int) {
	state["_version"] = version
}

func (sp *StatePersister) createBackup() error {
	if _, err := os.Stat(sp.configPath); os.IsNotExist(err) {
		return nil // No file to backup
	}

	data, err := os.ReadFile(sp.configPath)
	if err != nil {
		return err
	}

	return os.WriteFile(sp.backupPath, data, 0644)
}

func (sp *StatePersister) SaveState(state map[string]interface{}) error {
	// Add timestamp
	state["_updated_at"] = time.Now().Format(time.RFC3339)

	// Marshal to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write to temporary file first
	tempPath := sp.configPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, sp.configPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
