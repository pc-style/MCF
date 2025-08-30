# MCF CLI Examples

This document provides practical examples of using the MCF CLI tool.

## Getting Started

### 1. Quick Start with NPX (No Installation)

```bash
# Check if MCF is already installed
npx @pc-style/mcf-cli status

# Install MCF
npx @pc-style/mcf-cli install

# Configure MCF
npx @pc-style/mcf-cli setup

# Start using MCF
npx @pc-style/mcf-cli run
```

### 2. Global Installation

```bash
# Install CLI globally
npm install -g @pc-style/mcf-cli

# Now use 'mcf' directly
mcf status
mcf install
mcf setup
mcf run
```

## Complete Setup Flow

### First Time Setup

```bash
# 1. Install MCF framework
mcf install

# Output:
# 🚀 MCF Framework Installer
# ✅ MCF installation completed successfully!
# 
# Next steps:
#   • Run mcf setup to configure MCF
#   • Run mcf run to start a MCF session
#   • Run mcf status to check installation status

# 2. Configure MCF
mcf setup

# Interactive prompts:
# ? Enable MCF intelligent hooks system? (Y/n)
# ? Enable enhanced status line? (Y/n) 
# ? Choose output style: Explanatory (recommended)
# ? Add MCF to your shell PATH? (Y/n)

# 3. Verify installation
mcf status

# Output shows all green checkmarks:
# ✅ MCF is fully operational!

# 4. Start your first MCF session
mcf run
```

### Non-Interactive Installation

```bash
# For automation or CI/CD
mcf install --yes

# This skips all prompts and uses defaults
```

## Template Management Examples

### List Available Templates

```bash
mcf templates

# Or explicitly:
mcf templates list

# Output:
# 📁 Available Templates:
# 
#   react-app       - Modern React application with TypeScript
#   node-api        - Express.js REST API with middleware
#   python-cli      - Python CLI application with Click
#   
# Total: 3 templates
```

### Get Template Information

```bash
mcf templates info react-app

# Output:
# 📋 Template: react-app
# 📝 Description: Modern React application with TypeScript
# 🏷️  Category: Frontend
# ⚙️  Prerequisites: Node.js, npm
# 🔧 Variables:
#   • project_name: Enter project name
#   • author: Enter author name
#   • description: Enter project description
# 📋 Steps: 5
```

### Initialize from Template

```bash
mcf templates init react-app

# This will prompt for variables and create the project
# ✅ Template 'react-app' initialized successfully!
```

## Status Checking Examples

### Healthy Installation

```bash
mcf status

# Output:
# 📊 MCF Status Check
# 
# 🔍 Installation Status
#   ✅ MCF directory found (/home/user/mcf)
# 🔍 Core Files
#   ✅ Main runner script found
#   ✅ Settings configuration found
#   ✅ Template engine found
# 🔍 Directory Structure
#   ✅ Templates directory found (3 items)
#   ✅ Hooks directory found (12 items)
#   ✅ Scripts directory found (1 items)
# 🔍 Configuration
#   ✅ Hooks system configured (6 hook types)
#   ✅ Status line enabled
#   ✅ Output style: explanatory
# 🔍 Templates
#   ✅ 3 templates available
# 🔍 Shell Integration
#   ✅ Shell integration configured
# 
# ✅ MCF is fully operational!
```

### Problematic Installation

```bash
mcf status

# Output:
# 📊 MCF Status Check
# 
# 🔍 Installation Status
#   ❌ MCF directory not found
#      Run mcf install to install MCF
# 
# ❌ MCF has some issues that need attention
# 
# Recommended actions:
#   • Run mcf install to install MCF
#   • Run mcf setup to configure MCF
```

## Troubleshooting Examples

### Installation Issues

```bash
# If installation fails
mcf install

# Common issues and solutions:
# - Permission denied: Don't run as root
# - Network issues: Check internet connection
# - Disk space: Ensure at least 100MB free
# - Missing tools: Install git, curl, python3
```

### Configuration Issues

```bash
# If setup fails
mcf setup

# Will show specific error messages and suggestions
# - Corrupted settings: Backup and recreate
# - Missing directories: Reinstall MCF
```

### Runtime Issues

```bash
# If MCF session fails to start
mcf run

# Common messages:
# ❌ MCF is not installed.
# Run mcf install first.
#
# ⚠️  MCF is not configured.
# Run mcf setup to configure MCF first.
```

## Integration with Existing Workflows

### CI/CD Pipeline

```bash
#!/bin/bash
# In your CI/CD script

# Install MCF non-interactively
npx @pc-style/mcf-cli install --yes

# Verify installation
npx @pc-style/mcf-cli status

# Use templates in automated builds
npx @pc-style/mcf-cli templates init $TEMPLATE_NAME
```

### Development Workflow

```bash
# Morning routine
mcf status          # Check everything is OK
mcf run            # Start MCF session

# Template work
mcf templates      # See available templates
mcf templates init project-type

# Check system health
mcf status         # Periodic health checks
```

### Maintenance

```bash
# Update/repair installation
mcf install        # Reinstalls/updates MCF

# Reconfigure
mcf setup          # Update configuration

# Health check
mcf status         # Verify everything works
```

## Advanced Usage

### Custom Install Location

The MCF CLI uses the standard `~/mcf` location, but you can work with different project directories:

```bash
# MCF installs to ~/mcf but works in any directory
cd /path/to/my/project
mcf run  # Runs Claude with MCF config in this directory
```

### Template Development

```bash
# After creating a custom template in ~/mcf/templates/
mcf templates info my-custom-template
mcf templates init my-custom-template
```

### Multiple Projects

```bash
# MCF works with multiple projects
cd ~/project1
mcf run

# Later...
cd ~/project2  
mcf run  # Same MCF config, different project
```