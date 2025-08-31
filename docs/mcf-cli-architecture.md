# MCF CLI Architecture and Features

## Overview
The MCF CLI (My Claude Flow Command Line Interface) is a comprehensive tool for managing Claude Code configurations, profiles, and development workflows.

## Key Features

### 🔧 Profile Management System
- **Multi-profile support**: Manage different Claude Code configurations
- **Profile-specific CLAUDE_CONFIG_DIR**: Each profile can have its own Claude configuration directory
- **Environment variable injection**: Automatic setting of ANTHROPIC_BASE_URL, ANTHROPIC_AUTH_TOKEN, ANTHROPIC_MODEL
- **JSON-based configuration**: Profiles stored in `~/.mcf/profiles/` directory

### 🚀 Seamless Claude Integration
- **Pass-through arguments**: Forward arguments to Claude Code after `--` separator
- **Process management**: Proper child process spawning with signal handling
- **Environment configuration**: Automatic environment variable setup based on profile
- **Error handling**: Comprehensive error reporting for Claude execution failures

### 📁 Project Management
- **Project discovery**: Automatic detection of project directories
- **Project metadata**: JSON-based project configuration files
- **Workspace switching**: Easy switching between different development workspaces
- **Project statistics**: Analysis of project structure and file counts

### 🔌 MCP Server Integration
- **Serena MCP support**: Built-in integration with Serena semantic analysis server
- **Server lifecycle management**: Start, stop, restart MCP servers
- **Health monitoring**: Server status and health checks
- **Log management**: Access to server logs and debugging information

### 📦 Zero Dependencies
- **Single-file executable**: All functionality bundled in `mcf-standalone-pure.js`
- **No external dependencies**: Works without npm install
- **Self-installation**: `mcf install` copies executable to `~/.local/bin/mcf`
- **Cross-platform compatibility**: Works on macOS, Linux, and Windows

## Architecture

### Micro-Block Architecture
- **Command micro-blocks**: Each command is a self-contained module
- **Service registry**: Dependency injection through singleton service registry
- **Interface contracts**: Clear input/output contracts for all components
- **Error contracts**: Specific error classes with detailed context

### Service Registry Pattern
```javascript
// Service registration
ServiceRegistry.getInstance().register('IConfigurationService', configService);

// Service access
const configService = ServiceRegistry.getInstance().get('IConfigurationService');
```

### Command Pattern
```javascript
export class ExampleCommand {
  static readonly metadata = {
    name: 'ExampleCommand',
    dependencies: {
      services: ['IConfigurationService'],
      commands: [],
      external: []
    }
  };

  async execute() {
    // Command implementation
  }
}
```

## Configuration System

### Profile Structure
```json
{
  "name": "work",
  "config": {
    "claude": {
      "configDirectory": "/Users/username/work/.claude",
      "model": "claude-3-5-sonnet-20241022",
      "smallFastModel": "claude-3-5-haiku-20241022",
      "flags": ["--dangerous-skip"],
      "environment": {
        "ANTHROPIC_BASE_URL": "https://api.anthropic.com",
        "ANTHROPIC_AUTH_TOKEN": "your-token-here"
      }
    },
    "mcp": {
      "servers": ["serena", "filesystem"]
    }
  }
}
```

### Environment Variables
- `CLAUDE_CONFIG_DIR`: Profile-specific Claude configuration directory
- `ANTHROPIC_BASE_URL`: Anthropic API base URL
- `ANTHROPIC_AUTH_TOKEN`: Anthropic API token
- `ANTHROPIC_MODEL`: Claude model to use
- `ANTHROPIC_SMALL_FAST_MODEL`: Small fast model for quick tasks

## Command Structure

### Core Commands
- `mcf config`: Profile management
- `mcf run`: Execute Claude with profile
- `mcf project`: Project management
- `mcf mcp`: MCP server management
- `mcf install`: Self-installation
- `mcf status`: System status check

### Command Options
- `--config <profile>`: Specify profile to use
- `--dangerous-skip`: Skip Claude permissions
- `--no-interactive`: Run in non-interactive mode
- `--working-dir <path>`: Specify working directory

## Error Handling

### Error Classes
- `CommandError`: General command execution errors
- `ConfigurationError`: Configuration-related errors
- `ClaudeServiceError`: Claude execution errors
- `MCPServiceError`: MCP server errors
- `ValidationError`: Input validation errors

### Error Context
```javascript
throw new CommandError(
  'Profile not found',
  'PROFILE_NOT_FOUND',
  { profileName, availableProfiles }
);
```

## Development Workflow

1. **Installation**: `npm install -g pc-style-mcf-cli` or `npx pc-style-mcf-cli install`
2. **Configuration**: Create profiles with `mcf config create <name>`
3. **Setup**: Configure environment variables and settings
4. **Execution**: Run Claude with `mcf run --config <profile>`
5. **Management**: Use project and MCP commands for advanced features

## Integration Points

### Claude Code Integration
- Environment variable injection
- Profile-based configuration
- Argument forwarding
- Process lifecycle management

### MCP Server Integration
- Server discovery and management
- Health monitoring
- Log aggregation
- Semantic analysis capabilities

### Project Ecosystem
- Template system integration
- Status line enhancement
- Hook system support
- Development workflow automation

## Performance Characteristics

- **Startup time**: <100ms for basic commands
- **Memory usage**: ~50MB for typical usage
- **Disk usage**: ~1MB for installation
- **Network**: Minimal, only for MCP server communication

## Security Considerations

- Environment variable validation
- Path traversal protection
- Command injection prevention
- Secure credential handling
- Permission escalation prevention

## Future Enhancements

- Plugin system for custom commands
- GUI interface options
- Cloud synchronization
- Team collaboration features
- Advanced profiling and analytics
