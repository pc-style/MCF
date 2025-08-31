# MCF CLI Tool

A comprehensive command-line interface for MCF (My Claude Flow) with Claude Code integration, profile management, and project automation.

## Features

- 🔧 **Profile Management**: Manage multiple Claude Code configurations
- 🚀 **Seamless Claude Integration**: Run Claude with profile-specific settings
- 📁 **Project Management**: Create and manage development workspaces
- 🔌 **MCP Server Support**: Manage Model Context Protocol servers
- 📦 **Zero Dependencies**: Single-file executable with no external dependencies
- 🛠️ **Self-Installation**: Install globally with `mcf install`

## Installation

### Global Installation (Recommended)

```bash
# Install globally via npm
npm install -g @pc-style/mcf-cli

# Or run directly with npx (no installation needed)
npx @pc-style/mcf-cli install
```

### Local Development

```bash
git clone https://github.com/pc-style/MCF.git
cd MCF/cli
npm install
npm link  # Creates global symlink for development
```

## Quick Start

```bash
# Install MCF globally
mcf install

# List available profiles
mcf config list

# Create a new profile
mcf config create my-workspace

# Run Claude with specific profile
mcf run --config my-workspace

# Check status
mcf status
```

## Usage

### Profile Management

```bash
# List all profiles
mcf config list

# Show profile details
mcf config show <profile-name>

# Create new profile
mcf config create <profile-name>

# Delete profile
mcf config delete <profile-name>

# Clone profile
mcf config clone <source-profile> <new-profile>

# Set default profile
mcf config set-default <profile-name>

# Edit profile interactively
mcf config edit <profile-name>

# Validate profile configuration
mcf config validate <profile-name>
```

### Claude Code Integration

```bash
# Run Claude with default profile
mcf run

# Run with specific profile
mcf run --config <profile-name>

# Run with dangerous skip (no permissions)
mcf run --config <profile-name> --dangerous-skip

# Pass arguments to Claude
mcf run --config <profile-name> -- --help
mcf run --config <profile-name> -- --version
mcf run --config <profile-name> -- /serena:status

# Non-interactive mode
mcf run --config <profile-name> --no-interactive

# Run in specific working directory
mcf run --config <profile-name> --working-dir /path/to/project
```

### Project Management

```bash
# List all projects
mcf project list

# Show project details
mcf project show <project-name>

# Create new project
mcf project create <project-name>

# Delete project
mcf project delete <project-name>

# Switch to project
mcf project switch <project-name>

# Show current project
mcf project current

# Discover projects in directory
mcf project discover

# Show project statistics
mcf project stats
```

### MCP Server Management

```bash
# List MCP servers
mcf mcp list

# Show server details
mcf mcp show <server-name>

# Start MCP server
mcf mcp start <server-name>

# Stop MCP server
mcf mcp stop <server-name>

# Restart MCP server
mcf mcp restart <server-name>

# Check server status
mcf mcp status <server-name>

# Show server health
mcf mcp health <server-name>

# View server logs
mcf mcp logs <server-name>

# Install MCP server
mcf mcp install <server-name>

# Remove MCP server
mcf mcp remove <server-name>

# Configure server
mcf mcp config <server-name>
```

### Legacy Commands

```bash
# MCF installation (legacy)
mcf install

# MCF setup (legacy)
mcf setup

# MCF status check
mcf status

# Template management (legacy)
mcf templates
mcf templates list
mcf templates info <template-name>
mcf templates init <template-name>
```

## Configuration Profiles

Profiles allow you to manage different Claude Code configurations:

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

## Environment Variables

MCF CLI automatically sets these environment variables when running Claude:

- `CLAUDE_CONFIG_DIR` - Profile-specific Claude configuration directory
- `ANTHROPIC_BASE_URL` - Anthropic API base URL
- `ANTHROPIC_AUTH_TOKEN` - Anthropic API token
- `ANTHROPIC_MODEL` - Claude model to use
- `ANTHROPIC_SMALL_FAST_MODEL` - Small fast model for quick tasks

## NPX Usage

You can run MCF CLI without installation using npx:

```bash
# Install MCF
npx @pc-style/mcf-cli install

# List profiles
npx @pc-style/mcf-cli config list

# Run Claude
npx @pc-style/mcf-cli run --config default

# Any other command
npx @pc-style/mcf-cli <command> [args...]
```

## Development Workflow

1. **Install MCF**: `mcf install`
2. **Create Profile**: `mcf config create my-workspace`
3. **Configure Profile**: `mcf config edit my-workspace`
4. **Run Claude**: `mcf run --config my-workspace`
5. **Manage Projects**: `mcf project create my-project`
6. **Check Status**: `mcf status` (anytime)

## Integration with MCF Ecosystem

This CLI tool is designed to work seamlessly with the existing MCF ecosystem:

- **install.sh** - Wrapped by `mcf install`
- **claude-mcf.sh** - Wrapped by `mcf run`
- **template-engine.py** - Wrapped by `mcf templates`
- **Serena MCP Server** - Integrated via `mcf mcp` commands

## Error Handling

The CLI provides comprehensive error handling:

- **Profile Issues** - Suggests configuration fixes
- **Missing Dependencies** - Provides installation commands
- **Permission Problems** - Offers troubleshooting steps
- **Claude Errors** - Shows detailed error messages
- **MCP Server Issues** - Provides health check suggestions

## Requirements

- **Node.js** >= 14.0.0
- **bash** (for running installation scripts)
- **Python 3** (for template engine)
- **git** (for cloning repositories)
- **Claude Code** (for `mcf run` commands)

## Advanced Usage

### Custom Profile Configuration

```bash
# Create profile with custom Claude config directory
mcf config create custom-profile

# Edit the profile to set CLAUDE_CONFIG_DIR
mcf config edit custom-profile

# The profile will automatically set:
# CLAUDE_CONFIG_DIR=/Users/username/.mcf/profiles/custom-profile/.claude
```

### MCP Server Integration

```bash
# Install Serena MCP server
mcf mcp install serena

# Start the server
mcf mcp start serena

# Check server health
mcf mcp health serena

# Use in Claude session
mcf run --config my-profile -- /serena:analyze
```

### Batch Operations

```bash
# Create multiple profiles
for profile in dev staging prod; do
  mcf config create $profile
done

# Validate all profiles
mcf config list | xargs -I {} mcf config validate {}
```

## Troubleshooting

### Common Issues

**"Command not found: mcf"**
```bash
# Install globally
npm install -g @pc-style/mcf-cli

# Or use npx
npx @pc-style/mcf-cli <command>
```

**"Profile not found"**
```bash
# List available profiles
mcf config list

# Create new profile
mcf config create <profile-name>
```

**"Claude execution failed"**
```bash
# Check Claude installation
claude --version

# Validate profile
mcf config validate <profile-name>
```

## Contributing

Contributions welcome! Please see the main MCF repository for guidelines.

## License

ISC - Same as MCF framework