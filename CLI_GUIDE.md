# MCF CLI Tool

MCF now includes a modern, user-friendly CLI tool that provides easy installation, configuration, and management via NPX!

## Quick Start

### Option 1: Run with NPX (Recommended)

No installation required - run directly:

```bash
# Install MCF
npx @pc-style/mcf-cli install

# Configure MCF  
npx @pc-style/mcf-cli setup

# Start MCF session
npx @pc-style/mcf-cli run

# Check status anytime
npx @pc-style/mcf-cli status

# Manage templates
npx @pc-style/mcf-cli templates
```

### Option 2: Global Installation

Install once, use anywhere:

```bash
# Install globally
npm install -g @pc-style/mcf-cli

# Now use 'mcf' directly
mcf install
mcf setup  
mcf run
```

## Available Commands

- **`mcf install`** - Install MCF framework (wraps install.sh)
- **`mcf setup`** - Interactive configuration wizard
- **`mcf run`** - Start MCF session (wraps claude-mcf.sh)  
- **`mcf templates`** - Manage project templates (wraps template-engine.py)
- **`mcf status`** - Comprehensive health check

## Benefits

✅ **Modern Interface** - Clean, colorful output with progress indicators  
✅ **NPX Support** - No installation required, always latest version  
✅ **Error Handling** - Helpful error messages and troubleshooting guidance  
✅ **Cross-Platform** - Works on Windows, macOS, and Linux  
✅ **Preserves Functionality** - Wraps existing scripts without changing them

## Documentation

See the [CLI directory](cli/) for complete documentation:
- [README.md](cli/README.md) - Full CLI reference
- [EXAMPLES.md](cli/EXAMPLES.md) - Usage examples and workflows

## Legacy Installation

The original bash installation method still works:

```bash
curl -fsSL https://raw.githubusercontent.com/pc-style/MCF/main/install.sh | bash
```

But we recommend using the CLI tool for the best experience!