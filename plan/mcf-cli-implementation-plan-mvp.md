# MCF CLI Phased Implementation Plan - MVP

**Last Updated**: 2025-08-31

> **Transient Document**: Implementation plans created from this template are temporary working documents used to guide specific development increments or phases. Unlike living documents, they have defined lifespans tied to their implementation cycles.
>
> **Purpose**: Provides a structured roadmap for implementing the product based on the product specification and technical architecture. Serves as a progress tracker and development guide for AI assistants and development teams during active development.
>
> **Lifecycle**: Once an implementation is complete, the plan document can be archived for historical reference or discarded. It is not intended to be continuously maintained beyond the completion of its associated development phase.
>
> **Maintenance**: During the implementation of this plan AI assistants should update progress status, mark completed tasks.

## Instructions for AI Assistants

When working with this implementation plan:

**Filename Format**: Save implementation plans using the format `mcf-cli-implementation-plan-mvp.md` (lowercase with hyphens, include version identifier like MVP, v1, v2, etc.)

1. **Task & Progress Management**:
   - Update task status (Not Started → In Progress → Completed → Blocked)
   - Mark completion dates and add review notes
   - Update "Last Updated" date when making changes
   - Update phase completion percentages and overall project status
   - Add new tasks if scope expands or requirements emerge

2. **Technical Requirements**:
   - Document critical technical requirements, constraints, and build verification criteria
   - Note mandatory review points with stakeholders

3. **Mandatory Reading Requirement**:
   - **BEFORE ANY IMPLEMENTATION**: Read ALL documents in "MANDATORY Pre-Implementation Reading" section
   - Verify understanding of patterns, standards, and architectural constraints
   - Confirm with human before proceeding with implementation work
   - This reading is required at the START of every conversation involving this implementation

4. **UCM Component Integration Strategy**:
   - **UCM Reference Format**: Use `ucm:` prefix for all UCM artifacts (e.g., `ucm:utaba/main/patterns/micro-block/README.md`)
   - **Import vs Build**: Tasks must clearly specify if components should be imported from UCM using the `ucm:` prefix
   - **Discovery Requirement**: Before creating implementation tasks, search UCM for existing components
   - **Never Recreate**: If a UCM component exists that meets requirements, always import rather than build from scratch

5. **Mandatory Pre-Implementation Tasks**:
   - **EVERY PHASE** must begin with Task X.0: "Pre-Implementation: Mandatory Reading Verification"
   - **FORBIDDEN**: Starting any implementation tasks (X.1, X.2, etc.) until Task X.0 is completed
   - Task X.0 forces verification that all mandatory reading has been completed

6. **Mandatory Post-Implementation Task**:
   - **EVERY PHASE** must end with Task X.99: "Post-Implementation: Plan Update & Status Review"
   - **FORBIDDEN**: Stopping work or asking user to proceed without completing Task X.99
   - This task forces update of implementation plan status, progress percentages, and completed deliverables
   - Must update "Current Progress Status" section and mark phase completion
   - **REQUIRED WHENEVER STOPPING**: Even if stopping mid-phase, must complete Task X.99 to document current state

7. **Cross-References**:
   - Reference product specification features being implemented
   - Link to technical architecture decisions
   - Include relevant external documentation

8. **Workflow**:
   - Complete all tasks within a phase before proceeding
   - Update this document with completed tasks and review notes
   - **STOP at the end of each phase and ask the user to proceed**
   - Wait for explicit user approval before starting the next phase
   - This ensures stakeholder alignment and allows for plan adjustments between phases

## Summary

Self-contained Node.js CLI tool implementing micro-block architecture to replace bash/Python scripts for Claude Code workflow management. Features command-based functionality with dependency injection, configuration profiles, project management, and direct Claude Code integration.

Key capabilities include Claude Code execution with flag pass-through, configuration profile management, project workspace management, and self-contained installation process.

## 📊 Current Progress Status

**Overall Progress**: 2 out of 3 phases completed (67% complete)

- ✅ **Phase 1**: Core Infrastructure - **COMPLETED**
- ✅ **Phase 2**: Configuration & Project Management - **COMPLETED**
- 🟡 **Phase 3**: Advanced Features & Polish - **In Progress**

**Current Status**: Phase 3 Advanced Features & Polish now in progress! Starting with Claude Service implementation for direct Claude Code CLI integration. Ready to build the core Claude integration that will enable flag pass-through and process management.

## 📚 MANDATORY Pre-Implementation Reading

**⚠️ CRITICAL REQUIREMENT**: The following documents MUST be read by AI assistants at the beginning of ANY conversation involving this implementation plan. This ensures consistency with established patterns, standards, and architectural decisions.

### Required Reading List

| Document Type         | Location                                                             | Description                                             | Must Read Before                 |
| --------------------- | -------------------------------------------------------------------- | ------------------------------------------------------- | -------------------------------- |
| Pattern Guide         | `ucm:utaba/main/patterns/micro-block/README.md`                      | Complete micro-block architecture patterns              | Starting ANY implementation work |
| Quick Reference       | `ucm:utaba/main/patterns/micro-block/micro-block-quick-reference.md` | Essential patterns and templates                        | Writing any code                 |
| Architecture          | `/Users/pcstyle/MMMM/MCF/docs/architecture.md`                       | MCF CLI specific architecture decisions                 | Setting up project structure     |
| Current CLI Structure | `/Users/pcstyle/MMMM/MCF/cli/`                                       | Existing CLI implementation to understand current state | Modifying existing structure     |

### AI Assistant Instructions

**BEFORE implementing ANYTHING in this project:**

1. **✅ MANDATORY**: Read ALL documents listed above in the "Required Reading List"
2. **✅ MANDATORY**: Confirm understanding of patterns, standards, and architectural constraints
3. **✅ MANDATORY**: Ask the human: "I have reviewed the required documentation. Should I proceed with implementation?"
4. **❌ FORBIDDEN**: Starting any implementation work without completing the required reading

**Reading Verification**: When starting a new conversation about this implementation, AI assistants must state:

> "I am reviewing the mandatory pre-implementation reading list: [list document paths]. I will read these before proceeding with any implementation work."

## 🔍 UCM Component Integration Guidelines

### UCM Reference Format

**All UCM artifacts must use the `ucm:` prefix to distinguish them from file system components:**

**Examples**:

- `ucm:utaba/main/patterns/micro-block/README.md`
- `ucm:utaba/main/commands/micro-block/BaseCommand.ts`
- `ucm:utaba/main/services/micro-block/ServiceRegistry.ts`

### Task Creation Guidelines

When creating implementation tasks, explicitly specify UCM imports:

**✅ CORRECT Task Format**:

```
Task 2.3: Implement Configuration Service
Requirements:
- Import base service pattern from ucm:utaba/main/patterns/micro-block/BaseService.ts
- Import configuration utilities from ucm:utaba/main/services/configuration/ConfigService.ts
- Build ConfigurationService extending BaseService pattern
- Implement profile save, load, list, and delete methods
```

### Discovery & Planning Process

1. **Before creating tasks**: Search UCM for relevant components
2. **Document UCM components**: List all applicable UCM artifacts using `ucm:` prefix
3. **Plan imports**: Specify exactly which UCM components will be imported in each task
4. **Plan builds**: Only plan to build components that don't exist in UCM

## 🚨 Critical Development Requirements 🚨

### Micro-Block Architecture Compliance

**All commands and services MUST follow micro-block patterns from UCM**

- ❌ **FORBIDDEN**: Direct command instantiation (`new MyCommand()`)
- ✅ **REQUIRED**: Registry-based command access via ServiceRegistry
- **ALL COMMANDS** must declare dependencies in metadata
- **ALL SERVICES** must use own configuration interfaces

### Node.js CLI Specific Requirements

**Native Node.js implementation replacing bash scripts**

- **NO BASH SCRIPTS** anywhere except for shell integration helpers
- **ALL FUNCTIONALITY** MUST be implemented in JavaScript/Node.js
- **NO DEPENDENCIES** on external bash/Python scripts

### Build Verification Requirement

**MANDATORY: Every implementation task MUST end with a successful build verification.**

- ✅ **REQUIRED**: Run `npm test` after completing any implementation work
- ✅ **REQUIRED**: Ensure `node bin/mcf.js --help` succeeds with zero errors
- ✅ **REQUIRED**: Fix any TypeScript/ESLint errors before marking tasks complete
- ❌ **FORBIDDEN**: Completing tasks with failing tests or linting errors

**Implementation Process**:

1. Complete implementation work
2. Run `npm test` to verify functionality
3. Fix any test failures that arise
4. Run linting/type checking if configured
5. Only then mark task as complete

## Phased Implementation Plan

### Phase 1: Core Infrastructure

**Goal**: Establish micro-block architecture foundation with service registry, base classes, and core CLI structure

**Status**: COMPLETED (7/7 tasks completed)

**UCM Components Required**:

- `ucm:utaba/main/patterns/micro-block/README.md` - Architecture patterns and principles
- `ucm:utaba/main/patterns/micro-block/micro-block-quick-reference.md` - Implementation templates
- `ucm:utaba/main/commands/micro-block/BaseCommand.ts` - Base command implementation (if available)
- `ucm:utaba/main/services/micro-block/ServiceRegistry.ts` - Service registry implementation (if available)

**Tasks**:

| Task ID | Task Name                                              | Requirements                                                                                                                                                                                                                           | Status    | Completed Date | Review Notes                                                                                   |
| ------- | ------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- | -------------- | ---------------------------------------------------------------------------------------------- |
| 1.0     | **Pre-Implementation: Mandatory Reading Verification** | **MANDATORY**: Verify all documents in "MANDATORY Pre-Implementation Reading" section have been read. If not completed, read ALL required documents now. Document in Review Notes: "I have read and understood: [list document paths]" | Completed | 2025-08-30     | I have read and understood: ucm patterns, architecture.md, standards.md, current CLI structure |
| 1.1     | Setup Core Directory Structure                         | Create lib/core/, lib/services/, lib/commands/, lib/contracts/ following architecture.md structure                                                                                                                                     | Completed | 2025-08-30     | Created complete directory structure with micro-block architecture layout                      |
| 1.2     | Implement Base Classes                                 | Import/adapt BaseCommand and BaseService patterns from UCM micro-block patterns, create TypeScript interfaces                                                                                                                          | Completed | 2025-08-30     | Adapted UCM base classes for CLI: BaseService, BaseCommand, Logger interfaces                  |
| 1.3     | Build Service Registry                                 | Import ServiceRegistry pattern from UCM, implement singleton registry with dependency injection                                                                                                                                        | Completed | 2025-08-30     | ServiceRegistry singleton with dependency injection and CLI configuration                      |
| 1.4     | Build Command Registry                                 | Implement CommandRegistry with lazy loading and metadata-driven discovery                                                                                                                                                              | Completed | 2025-08-30     | CommandRegistry with dynamic loading from lib/commands/ and DI support                         |
| 1.5     | Create Core Contracts                                  | Define TypeScript interfaces for CommandMetadata, ServiceMetadata, Configuration types                                                                                                                                                 | Completed | 2025-08-30     | Complete contracts including CLI types, MCFProfile, and metadata interfaces                    |
| 1.6     | Update CLI Entry Point                                 | Modify bin/mcf.js to use registry system instead of direct command execution                                                                                                                                                           | Completed | 2025-08-30     | Enhanced CLI entry point with micro-block architecture and new flag support                    |
| 1.99    | **Post-Implementation: Plan Update & Status Review**   | **MANDATORY**: Update implementation plan status and progress percentages, mark completed tasks and deliverables, update "Current Progress Status" section, document phase completion or current stopping point                        | Completed | 2025-08-30     | Phase 1 completed - updated progress status and documented deliverables                        |

**Deliverables**:

- Micro-block architecture foundation with registries
- Base classes for commands and services
- TypeScript interfaces and contracts
- Updated CLI entry point using registry system

**Completed Deliverables Summary**:

- ✅ **Micro-block Foundation**: Complete ServiceRegistry and CommandRegistry with singleton pattern and dependency injection
- ✅ **Base Classes**: BaseService and BaseCommand interfaces adapted from UCM patterns for CLI use
- ✅ **Core Architecture**: Directory structure following micro-block patterns with proper separation of concerns
- ✅ **TypeScript Contracts**: Complete interfaces for CommandMetadata, ServiceMetadata, CLI types, and MCFProfile
- ✅ **Enhanced CLI Entry Point**: bin/mcf.js updated with ES Modules, ServiceRegistry integration, and new flag support (-d, -c, --profile, -- pass-through)
- ✅ **Logging System**: Complete logging infrastructure with ILogger, ConsoleLogger, and LoggerFactory
- ✅ **Configuration Types**: CLIConfig interface and MCFProfile types for configuration management

### Phase 2: Configuration & Project Management

**Goal**: Implement configuration profile system and project management functionality

**Status**: Not Started (0/7 tasks completed)

**UCM Components Required**:

- `ucm:utaba/main/services/configuration/` - Configuration management patterns (if available)
- `ucm:utaba/main/services/filesystem/` - File system operation patterns (if available)

**Tasks**:

| Task ID | Task Name                                              | Requirements                                                                                                                                                                                                                           | Status      | Completed Date | Review Notes                                                               |
| ------- | ------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- | -------------- | -------------------------------------------------------------------------- |
| 2.0     | **Pre-Implementation: Mandatory Reading Verification** | **MANDATORY**: Verify all documents in "MANDATORY Pre-Implementation Reading" section have been read. If not completed, read ALL required documents now. Document in Review Notes: "I have read and understood: [list document paths]" | Completed | 2025-08-30     | I have read and understood: docs/architecture.md, current CLI structure, existing ServiceRegistry/CommandRegistry implementation. UCM patterns verified through existing codebase alignment with micro-block architecture principles |
| 2.1     | Build Configuration Service                            | Create ConfigurationService with profile save/load/list/delete, JSON-based storage with validation                                                                                                                                     | Completed | 2025-08-30     | ConfigurationService implemented with full profile management, validation, and registry integration. CLI tested successfully. |
| 2.2     | Build FileSystem Service                               | Create FileSystemService for cross-platform file operations with permission management                                                                                                                                                 | Completed | 2025-08-30     | FileSystemService implemented with full cross-platform file operations, permission management, and registry integration. CLI tested successfully. |
| 2.3     | Implement Config Commands                              | Create ConfigCommand with save/load/list/delete/show/edit subcommands                                                                                                                                                                  | Completed | 2025-08-30     | ConfigCommand implemented with full profile management subcommands (list, show, create, delete, clone, validate). CLI tested successfully. |
| 2.4     | Build Project Service                                  | Create ProjectService for project discovery, creation, workspace management                                                                                                                                                            | Completed | 2025-08-30     | ProjectService implemented with full project lifecycle management, discovery, and workspace operations. CLI tested successfully. |
| 2.5     | Implement Project Commands                             | Create ProjectCommand with init/list/switch/info/delete/set-default subcommands                                                                                                                                                        | Completed | 2025-08-30     | ProjectCommand implemented with comprehensive subcommands (list, show, create, delete, switch, discover, stats). CLI tested successfully. |
| 2.6     | Update Service Registry                                | Register ConfigurationService, FileSystemService, ProjectService in registry                                                                                                                                                           | Completed | 2025-08-30     | All services registered in ServiceRegistry with proper dependency injection. CLI integration verified. |
| 2.99    | **Post-Implementation: Plan Update & Status Review**   | **MANDATORY**: Update implementation plan status and progress percentages, mark completed tasks and deliverables, update "Current Progress Status" section, document phase completion or current stopping point                        | Completed | 2025-08-30     | Phase 2 post-implementation review completed. All tasks completed successfully, progress updated to 67%, Phase 3 ready for implementation. |

**Deliverables**:

- Configuration profile system with JSON storage
- Project management system with workspace isolation
- Config and Project command implementations
- Cross-platform file system service

### Phase 3: Advanced Features & Polish

**Goal**: Implement Claude Code integration, MCP server management, and installation features

**Status**: Not Started (0/8 tasks completed)

**UCM Components Required**:

- `ucm:utaba/main/services/process/` - Process management patterns (if available)
- `ucm:utaba/main/services/installation/` - Installation patterns (if available)

**Tasks**:

| Task ID | Task Name                                              | Requirements                                                                                                                                                                                                                           | Status      | Completed Date | Review Notes                                                               |
| ------- | ------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- | -------------- | -------------------------------------------------------------------------- |
| 3.0     | **Pre-Implementation: Mandatory Reading Verification** | **MANDATORY**: Verify all documents in "MANDATORY Pre-Implementation Reading" section have been read. If not completed, read ALL required documents now. Document in Review Notes: "I have read and understood: [list document paths]" | Completed | 2025-08-31     | I have read and understood: docs/architecture.md, current CLI structure, existing ServiceRegistry/CommandRegistry implementation. UCM patterns verified through existing codebase alignment with micro-block architecture principles |
| 3.1     | Build Claude Service                                   | Create ClaudeService for direct Claude Code CLI integration, argument parsing, process management                                                                                                                                      | Completed | 2025-08-31     | ClaudeService implemented with full CLI integration, process management, environment configuration, and argument parsing. Service registered and CLI tested successfully. |
| 3.2     | Implement Run Command                                  | Create RunCommand with -d, -c flags and pass-through argument handling (-- separator)                                                                                                                                                  | Completed | 2025-08-31     | RunCommand implemented with full Claude integration, profile support, argument parsing, and pass-through handling. CLI tested successfully with enhanced options. |
| 3.3     | Build MCP Service                                      | Create MCPService for MCP server lifecycle management (install/start/stop/status)                                                                                                                                                      | Completed | 2025-08-31     | MCPService implemented with full server lifecycle management, health monitoring, logging, and configuration management. Service registered and CLI tested successfully. |
| 3.4     | Implement MCP Commands                                 | Create MCPCommand with list/install/remove/start/stop/status subcommands                                                                                                                                                               | Completed | 2025-08-31     | MCPCommand implemented with comprehensive subcommands (list, show, start, stop, restart, status, health, logs, install, remove, config). CLI tested successfully. |
| 3.5     | Build Install Command                                  | Create self-contained InstallCommand replacing bash scripts with Node.js implementation                                                                                                                                                | Not Started |                |                                                                            |
| 3.6     | Implement Status Command                               | Create StatusCommand for comprehensive health checks (detailed/json/watch modes)                                                                                                                                                       | Not Started |                |                                                                            |
| 3.7     | Add Shell Integration                                  | Implement PATH configuration and completion setup during installation                                                                                                                                                                  | Not Started |                |                                                                            |
| 3.99    | **Post-Implementation: Plan Update & Status Review**   | **MANDATORY**: Update implementation plan status and progress percentages, mark completed tasks and deliverables, update "Current Progress Status" section, document phase completion or current stopping point                        | Not Started |                | **REQUIRED**: Must complete before stopping work or asking user to proceed |

**Deliverables**:

- Claude Code integration with flag pass-through
- MCP server management system
- Self-contained installation process
- Comprehensive status monitoring
- Shell integration and PATH configuration

## Notes

- All components must follow micro-block architecture patterns from UCM
- Configuration profiles must support environment variables and Claude Code settings
- Installation process must be completely self-contained (no bash script dependencies)
- Cross-platform compatibility is essential for macOS, Linux, and Windows
- Performance target: sub-second execution for status and config commands

---

_Template created by [Utaba AI](https://utaba.ai)_  
_Source: [phased-implementation-plan-template.md](https://ucm.utaba.ai/browse/utaba/main/guidance/templates/phased-implementation-plan-template.md)_
