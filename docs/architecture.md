# MCF CLI Architecture Specification

**Last Updated**: 2025-08-30

> **Living Document**: This architecture specification is actively maintained throughout the project lifecycle. It serves as the authoritative source for architectural decisions and should be updated when making significant architectural changes, technology stack modifications, or design pattern updates.
>
> **Purpose**: Provides high-level architectural guidance and technology stack decisions for AI assistants, developers, and stakeholders. Should be referenced in the main README.md as essential reading for understanding system design and architectural patterns.
>
> **Maintenance**: AI assistants should recommend updates when discovering architectural inconsistencies or when implementing features that require architectural changes to keep this document aligned with the actual system.

## Instructions for AI Assistants

When working with this architecture specification:

1. **Technology Stack Decisions**: Document technology choices with clear rationale and version constraints
2. **Architecture Pattern Application**: Reference UCM patterns or industry standards and note any adaptations
3. **Component Organization**: Structure components by business capability with clear interface boundaries

## Overview

MCF CLI is a self-contained Node.js command-line tool that replaces bash/Python scripts for managing Claude Code workflows. It implements a micro-block architecture adapted for CLI applications, providing command-based functionality with dependency injection, configuration management, project management, and MCP server orchestration. The architecture prioritizes modularity, testability, and AI-collaborative development patterns.

## Technology Stack

### Core Technologies

- **Runtime**: Node.js 14.0.0+
- **CLI Framework**: Commander.js 11.0.0+
- **Architecture Pattern**: Micro-Block Pattern (adapted from UCM `utaba/main/patterns/micro-block`)

### Key Dependencies

- **User Interface**: Chalk 4.1.2+ for terminal styling, Inquirer.js 8.2.6+ for interactive prompts
- **Loading Indicators**: Ora 5.4.1+ for progress spinners
- **Process Management**: Child processes for Claude Code CLI and MCP server integration
- **Configuration**: JSON-based profile system with validation

## Architecture Patterns

### Micro-Block Pattern (CLI Adaptation)

The micro-block architecture is adapted for CLI applications where each command represents a discrete, contract-driven component. Unlike web applications, CLI commands have direct user interaction and process lifecycle considerations.

#### Core Principles

- **Command-as-Component**: Each CLI command is a self-contained micro-block with explicit contracts
- **Service-Based Infrastructure**: Reusable services for configuration, project management, Claude integration
- **Registry-Driven Composition**: Commands discovered and composed via registry with dependency injection
- **Contract-First Design**: All interfaces defined before implementation with rich metadata

#### Component Structure

```
cli/
├── bin/                    # Entry point
├── lib/                    # Core implementation
│   ├── commands/          # Command micro-blocks
│   │   ├── install/       # Installation commands
│   │   ├── run/           # Runtime commands
│   │   ├── config/        # Configuration commands
│   │   └── project/       # Project management commands
│   ├── services/          # Infrastructure services
│   │   ├── configuration/ # Profile and settings management
│   │   ├── claude/        # Claude Code integration
│   │   ├── mcp/           # MCP server management
│   │   └── project/       # Project lifecycle
│   └── core/              # Registry and base classes
│       ├── registry/      # Service and command registries
│       ├── base/          # Base command and service classes
│       └── contracts/     # Interfaces and type definitions
```

## System Architecture

### Request Flow

```
CLI Entry (bin/mcf.js)
    ↓
Command Parser (Commander.js)
    ↓
Service Registry Initialization
    ↓
Command Registry Resolution
    ↓
Command Instantiation with Dependencies
    ↓
Command Execution
    ↓
Result Output / Process Management
```

### Component Design

#### Commands (Business Logic)

- **InstallCommand**: Self-contained MCF installation replacing bash scripts
- **RunCommand**: Claude Code execution with flag parsing and pass-through arguments
- **ConfigCommand**: Configuration profile management (save/load/list/delete)
- **ProjectCommand**: Project lifecycle management (init/switch/list/info)
- **MCPCommand**: MCP server lifecycle (install/start/stop/status)
- **TemplateCommand**: Template discovery and initialization

#### Services (Infrastructure)

- **ConfigurationService**: Profile storage, loading, and validation with environment variable management
- **ClaudeService**: Direct integration with Claude Code CLI, argument parsing, and process management
- **ProjectService**: Project discovery, creation, and workspace management
- **MCPService**: MCP server installation, lifecycle management, and status monitoring
- **FileSystemService**: Cross-platform file operations with permission management

#### Core Registry System

- **ServiceRegistry**: Singleton registry managing service lifecycle and dependency injection
- **CommandRegistry**: Lazy-loading command registry with metadata-driven discovery
- **ConfigurationRegistry**: Configuration profile registry with validation and inheritance

## Performance & Scalability

**Sub-second Command Execution**: Commands like `status` and `config list` execute in <500ms through caching and lazy loading.

**Parallel Operations**: Installation and setup tasks use parallel processing for network operations and file system tasks.

**Resource Management**: Services are instantiated on-demand and cached for the duration of command execution.

## Security

**Configuration Isolation**: Each profile maintains isolated environment variables and Claude configurations.

**Process Sandboxing**: Child processes (Claude Code, MCP servers) run with appropriate permissions and resource limits.

**Input Validation**: All user input validated before processing, with sanitization for shell command construction.

**File System Permissions**: Configuration files protected with appropriate Unix permissions (600 for sensitive data).

## Deployment

**NPM Distribution**: Published as `@pc-style/mcf-cli` with global installation support.

**NPX Compatibility**: Full functionality available via `npx @pc-style/mcf-cli` without installation.

**Cross-Platform Support**: Native Node.js implementation ensures compatibility across macOS, Linux, and Windows.

**Shell Integration**: Automatic PATH configuration and completion setup during installation.

---

_Template created by [Utaba AI](https://utaba.ai)_  
_Source: [architecture-template.md](https://ucm.utaba.ai/browse/utaba/main/guidance/templates/architecture-template.md)_
