# MCF CLI API Reference

## Table of Contents
- [Command Reference](#command-reference)
- [Service APIs](#service-apis)
- [Configuration APIs](#configuration-apis)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Command Reference

### `mcf config` - Profile Management

#### `mcf config list`
List all available profiles.

**Output:**
```
MCF Configuration Profiles:

  agentwise
  mcf
  proxy

Total: 3 profile(s)
```

#### `mcf config show <profile-name>`
Display detailed information about a specific profile.

**Parameters:**
- `profile-name`: Name of the profile to display

**Output:**
```json
Profile: mcf
Configuration:
  CLAUDE_CONFIG_DIR: /Users/username/mcf/.claude
  ANTHROPIC_BASE_URL: https://api.anthropic.com
  ANTHROPIC_AUTH_TOKEN: ********
  ANTHROPIC_MODEL: claude-3-5-sonnet-20241022
```

#### `mcf config create <profile-name>`
Create a new profile with interactive configuration.

**Parameters:**
- `profile-name`: Name for the new profile

**Interactive Prompts:**
- CLAUDE_CONFIG_DIR path
- Anthropic API settings
- Model preferences

#### `mcf config delete <profile-name>`
Delete an existing profile.

**Parameters:**
- `profile-name`: Name of the profile to delete

#### `mcf config clone <source-profile> <new-profile>`
Clone an existing profile to create a new one.

**Parameters:**
- `source-profile`: Profile to clone from
- `new-profile`: Name for the new profile

#### `mcf config set-default <profile-name>`
Set the default profile for `mcf run` command.

**Parameters:**
- `profile-name`: Profile to set as default

#### `mcf config edit <profile-name>`
Interactively edit an existing profile.

**Parameters:**
- `profile-name`: Profile to edit

#### `mcf config validate <profile-name>`
Validate profile configuration and report issues.

**Parameters:**
- `profile-name`: Profile to validate

### `mcf run` - Claude Code Execution

#### Basic Usage
```bash
mcf run --config <profile-name>
```

#### Advanced Options
```bash
mcf run --config <profile-name> --dangerous-skip
mcf run --config <profile-name> --no-interactive
mcf run --config <profile-name> --working-dir /path/to/project
```

#### Pass-through Arguments
```bash
mcf run --config <profile-name> -- --help
mcf run --config <profile-name> -- --version
mcf run --config <profile-name> -- /serena:analyze
```

**Parameters:**
- `--config <profile>`: Profile to use for Claude execution
- `--dangerous-skip`: Skip Claude permission checks
- `--no-interactive`: Run in non-interactive mode
- `--working-dir <path>`: Set working directory for Claude
- `--`: Separator for Claude-specific arguments

### `mcf project` - Project Management

#### `mcf project list`
List all managed projects.

#### `mcf project show <project-name>`
Show detailed project information.

#### `mcf project create <project-name>`
Create a new project workspace.

#### `mcf project delete <project-name>`
Delete a project workspace.

#### `mcf project switch <project-name>`
Switch to a different project workspace.

#### `mcf project current`
Show the current active project.

#### `mcf project discover`
Discover projects in the current directory.

#### `mcf project stats`
Show project statistics and metrics.

### `mcf mcp` - MCP Server Management

#### `mcf mcp list`
List all configured MCP servers.

#### `mcf mcp show <server-name>`
Show detailed MCP server information.

#### `mcf mcp start <server-name>`
Start an MCP server.

#### `mcf mcp stop <server-name>`
Stop an MCP server.

#### `mcf mcp restart <server-name>`
Restart an MCP server.

#### `mcf mcp status <server-name>`
Check MCP server status.

#### `mcf mcp health <server-name>`
Check MCP server health.

#### `mcf mcp logs <server-name>`
View MCP server logs.

#### `mcf mcp install <server-name>`
Install an MCP server.

#### `mcf mcp remove <server-name>`
Remove an MCP server.

#### `mcf mcp config <server-name>`
Configure an MCP server.

### `mcf install` - Self-Installation

Install MCF CLI to `~/.local/bin/mcf`.

**Options:**
- `--force`: Force reinstallation even if already installed

### `mcf status` - System Status

Check MCF installation and system status.

**Checks:**
- Installation completeness
- Configuration validity
- Profile availability
- MCP server status
- Dependencies

## Service APIs

### ConfigurationService

#### `saveProfile(name: string, profile: MCFProfile): Promise<void>`
Save a profile to the filesystem.

**Parameters:**
- `name`: Profile name
- `profile`: Profile configuration object

#### `loadProfile(name: string): Promise<MCFProfile>`
Load a profile from the filesystem.

**Parameters:**
- `name`: Profile name

**Returns:** Profile configuration object

#### `listProfiles(): Promise<string[]>`
List all available profile names.

**Returns:** Array of profile names

#### `deleteProfile(name: string): Promise<void>`
Delete a profile from the filesystem.

**Parameters:**
- `name`: Profile name to delete

#### `validateProfile(profile: MCFProfile): Promise<ValidationResult>`
Validate profile configuration.

**Parameters:**
- `profile`: Profile to validate

**Returns:** Validation result with errors and warnings

### ClaudeService

#### `runClaude(options: ClaudeRunOptions): Promise<ClaudeRunResult>`
Execute Claude Code with specified options.

**Parameters:**
- `options`: Claude execution options

**Options:**
```typescript
interface ClaudeRunOptions {
  profile?: string;
  dangerousSkip?: boolean;
  interactive?: boolean;
  workingDirectory?: string;
  additionalArgs?: string[];
}
```

**Returns:**
```typescript
interface ClaudeRunResult {
  exitCode: number;
  executionTime: number;
  output?: string;
}
```

### ProjectService

#### `createProject(name: string, path?: string): Promise<Project>`
Create a new project workspace.

#### `loadProject(name: string): Promise<Project>`
Load project configuration.

#### `listProjects(): Promise<string[]>`
List all project names.

#### `deleteProject(name: string): Promise<void>`
Delete a project workspace.

#### `switchProject(name: string): Promise<void>`
Switch to a different project.

#### `getCurrentProject(): Promise<Project | null>`
Get the currently active project.

### MCPService

#### `installServer(name: string): Promise<void>`
Install an MCP server.

#### `startServer(name: string): Promise<void>`
Start an MCP server.

#### `stopServer(name: string): Promise<void>`
Stop an MCP server.

#### `getServerStatus(name: string): Promise<ServerStatus>`
Get MCP server status.

#### `getServerHealth(name: string): Promise<HealthStatus>`
Check MCP server health.

#### `getServerLogs(name: string): Promise<string[]>`
Retrieve MCP server logs.

## Configuration APIs

### Profile Configuration

```typescript
interface MCFProfile {
  name: string;
  config: {
    claude?: {
      configDirectory?: string;
      model?: string;
      smallFastModel?: string;
      dangerousSkip?: boolean;
      environment?: Record<string, string>;
      flags?: string[];
    };
    mcp?: {
      servers?: string[];
    };
  };
}
```

### Project Configuration

```typescript
interface Project {
  name: string;
  path: string;
  config: {
    description?: string;
    tags?: string[];
    settings?: Record<string, any>;
  };
}
```

### MCP Server Configuration

```typescript
interface MCPServer {
  name: string;
  command: string;
  args?: string[];
  env?: Record<string, string>;
  config?: Record<string, any>;
}
```

## Error Handling

### Error Types

#### `CommandError`
General command execution errors.

```typescript
class CommandError extends Error {
  constructor(
    message: string,
    public readonly code?: string,
    public readonly details?: Record<string, any>
  ) {
    super(message);
    this.name = 'CommandError';
  }
}
```

#### `ConfigurationError`
Configuration-related errors.

#### `ClaudeServiceError`
Claude Code execution errors.

#### `ValidationError`
Input validation errors.

#### `MCPServiceError`
MCP server errors.

### Error Codes

- `PROFILE_NOT_FOUND`: Profile doesn't exist
- `PROFILE_EXISTS`: Profile already exists
- `INVALID_PROFILE`: Profile configuration is invalid
- `CLAUDE_EXECUTION_FAILED`: Claude Code execution failed
- `MCP_SERVER_ERROR`: MCP server operation failed
- `PERMISSION_DENIED`: Insufficient permissions
- `VALIDATION_FAILED`: Input validation failed

## Examples

### Profile Management
```bash
# Create and configure a new profile
mcf config create work-profile
mcf config edit work-profile

# Use the profile
mcf run --config work-profile -- --help
```

### Project Workflow
```bash
# Create a new project
mcf project create my-app
mcf project switch my-app

# Work with Claude in project context
mcf run --config work-profile --working-dir $(mcf project current)
```

### MCP Server Management
```bash
# Install and start Serena MCP server
mcf mcp install serena
mcf mcp start serena
mcf mcp health serena

# Use with Claude
mcf run --config work-profile -- /serena:analyze
```

### Advanced Usage
```bash
# Run Claude with custom arguments
mcf run --config work-profile -- --debug --verbose /help

# Batch operations
for profile in dev staging prod; do
  mcf config create $profile
  mcf config edit $profile
done

# Validate all profiles
mcf config list | xargs -I {} mcf config validate {}
```

## Response Codes

### Success Codes
- `0`: Command executed successfully
- `1`: General error
- `2`: Command not found
- `3`: Invalid arguments
- `4`: Configuration error
- `5`: Permission denied
- `6`: Network error
- `7`: Dependency missing

### Claude Exit Codes
- `0`: Claude executed successfully
- `1`: Claude execution failed
- `130`: Claude interrupted (Ctrl+C)
- `137`: Claude killed (SIGKILL)
