# MCF CLI Tool

A comprehensive command-line interface for MCF (Multi Component Framework) installation, configuration, and management.

## Installation

### Global Installation (Recommended)

```bash
# Install globally via npm
npm install -g @pc-style/mcf-cli

# Or run directly with npx (no installation needed)
npx @pc-style/mcf-cli --help
```

### Local Development

```bash
git clone https://github.com/pc-style/MCF.git
cd MCF/cli
npm install
npm link  # Creates global symlink for development
```

## Usage

### Basic Commands

```bash
# Show help
mcf --help

# Check MCF installation status
mcf status

# Install MCF framework
mcf install

# Configure MCF (run after installation)
mcf setup

# Start MCF session
mcf run

# Manage templates
mcf templates
mcf templates list
mcf templates info <template-name>
mcf templates init <template-name>
```

### Installation Options

```bash
# Interactive installation (default)
mcf install

# Non-interactive installation (skip prompts)
mcf install --yes
```

### Template Management

```bash
# List all available templates
mcf templates
mcf templates list

# Get information about a specific template
mcf templates info react-app

# Initialize project from template
mcf templates init react-app
```

## Commands Reference

### `mcf install [options]`

Installs the MCF framework to `~/mcf/` directory.

**Options:**
- `--yes, -y` - Skip interactive prompts and proceed automatically

**What it does:**
- Downloads and installs MCF components
- Sets up directory structure
- Configures permissions
- Installs dependencies (Serena MCP server)

### `mcf setup`

Interactive configuration wizard for MCF.

**What it does:**
- Configures hooks system
- Sets up status line
- Chooses output style
- Sets up shell integration

### `mcf run`

Starts an MCF session with Claude integration.

**What it does:**
- Validates installation
- Starts Claude with MCF configuration
- Handles first-run authentication
- Provides enhanced development environment

### `mcf templates [action] [name]`

Manages MCF project templates.

**Actions:**
- `list` - List available templates (default)
- `info <name>` - Show template information
- `init <name>` - Initialize project from template

### `mcf status`

Comprehensive status check of MCF installation.

**What it checks:**
- Installation completeness
- Core files presence
- Directory structure
- Configuration status
- Template availability
- Shell integration

## NPX Usage

You can run MCF CLI without installation using npx:

```bash
# Install MCF
npx @pc-style/mcf-cli install

# Check status
npx @pc-style/mcf-cli status

# Run setup
npx @pc-style/mcf-cli setup

# Any other command
npx @pc-style/mcf-cli <command>
```

## Development Workflow

1. **Install MCF**: `mcf install`
2. **Configure**: `mcf setup`
3. **Start Session**: `mcf run`
4. **Check Status**: `mcf status` (anytime)
5. **Manage Templates**: `mcf templates` (as needed)

## Integration with MCF Ecosystem

This CLI tool is designed to work seamlessly with the existing MCF bash and Python scripts:

- **install.sh** - Wrapped by `mcf install`
- **claude-mcf.sh** - Wrapped by `mcf run`
- **template-engine.py** - Wrapped by `mcf templates`

The CLI provides a modern, user-friendly interface while preserving all existing functionality.

## Error Handling

The CLI provides helpful error messages and suggestions:

- **Missing dependencies** - Suggests installation commands
- **Configuration issues** - Points to setup command
- **Permission problems** - Provides troubleshooting steps
- **Missing files** - Recommends reinstallation

## Requirements

- **Node.js** >= 14.0.0
- **bash** (for running installation scripts)
- **Python 3** (for template engine)
- **git** (for cloning repositories)

## License

ISC - Same as MCF framework

## Contributing

Contributions welcome! Please see the main MCF repository for guidelines.