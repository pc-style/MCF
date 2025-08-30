//go:build ignore

package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Testing configuration schema loading...")

	// Test schema loading
	schemaPath := "config-schema.yaml"
	schema, err := LoadConfigSchema(schemaPath)
	if err != nil {
		log.Fatalf("Failed to load schema: %v", err)
	}

	fmt.Printf("✅ Successfully loaded schema with %d sections\n", len(schema.Sections))

	for i, section := range schema.Sections {
		fmt.Printf("  Section %d: %s (%d fields)\n", i+1, section.Name, len(section.Fields))
	}

	// Test configuration manager
	configMgr := NewConfigManager(".")
	fmt.Println("✅ Successfully created configuration manager")

	// Test configurator model creation
	configurator := NewConfiguratorModel(configMgr)
	fmt.Println("✅ Successfully created configurator model")

	// Test field count
	fmt.Printf("✅ Configurator has %d sections\n", len(configurator.sections))

	fmt.Println("🎉 All configuration components initialized successfully!")
}
