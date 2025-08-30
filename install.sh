#!/bin/bash
# MCF (My Claude Framework) - One-Command Install Script
# Usage: curl -fsSL https://raw.githubusercontent.com/pc-style/MCF/refs/heads/main/install.sh | bash
# Version: 1.0.0
# Security Level: High

IFS=$'\n\t'

# Configuration
readonly SCRIPT_VERSION="1.0.0"
readonly MCF_REPO="${MCF_REPO_URL:-https://github.com/pc-style/MCF.git}"
readonly MCF_DIR="${MCF_INSTALL_DIR:-$HOME/mcf}"
readonly REQUIRED_FILES=("claude-mcf.sh" "README.md" "scripts" "templates" ".claude")
readonly LOG_FILE="/tmp/mcf-install-$(date +%s).log"

# Color codes for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Logging functions
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >> "$LOG_FILE"
}

info() {
    echo -e "${BLUE}ℹ️  INFO:${NC} $*"
    log "INFO: $*"
}

warn() {
    echo -e "${YELLOW}⚠️  WARN:${NC} $*"
    log "WARN: $*"
}

error() {
    echo -e "${RED}❌ ERROR:${NC} $*"
    log "ERROR: $*"
}

success() {
    echo -e "${GREEN}✅ SUCCESS:${NC} $*"
    log "SUCCESS: $*"
}

# Security functions
verify_environment() {
    log "Starting environment verification"
    
    # Don't run as root for security
    if [[ $EUID -eq 0 ]]; then
        error "Do not run this script as root for security reasons"
        exit 1
    fi
    
    # Check required tools
    for cmd in git curl; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            error "Required tool not found: $cmd"
            exit 1
        fi
    done
    
    # Check disk space (need at least 100MB)
    local available_space
    available_space=$(df "$HOME" 2>/dev/null | tail -1 | awk '{print $4}' || echo 0)
    if [[ $available_space -lt 102400 ]]; then
        error "Insufficient disk space (need at least 100MB in $HOME)"
        exit 1
    fi
    
    log "Environment verification passed"
}

sanitize_environment() {
    log "Sanitizing environment"
    
    # Remove potentially dangerous environment variables
    unset BASH_ENV ENV CDPATH
    export PATH="/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
    
    # Set secure umask
    umask 077
}

create_secure_temp() {
    log "Creating secure temporary directory"
    
    local temp_dir
    if command -v mktemp >/dev/null 2>&1; then
        temp_dir=$(mktemp -d -t mcf-install.XXXXXXXXXX)
    else
        temp_dir="/tmp/mcf-install-$$-$(date +%s)"
        mkdir -m 700 "$temp_dir"
    fi
    
    if [[ ! -d "$temp_dir" ]]; then
        error "Failed to create temporary directory"
        exit 1
    fi
    
    echo "$temp_dir"
}

secure_cleanup() {
    local temp_dir="$1"
    if [[ -n "$temp_dir" && -d "$temp_dir" ]]; then
        log "Cleaning up temporary directory: $temp_dir"
        rm -rf "$temp_dir" 2>/dev/null || true
    fi
}

# Validation functions
validate_repo_url() {
    local url="$1"
    if [[ ! $url =~ ^https://github\.com/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+\.git$ ]]; then
        error "Invalid GitHub repository URL: $url"
        exit 1
    fi
}

validate_directory_path() {
    local path="$1"
    # Convert to absolute path and check it's under HOME
    local abs_path
    local abs_home

    # Convert to absolute path - handle both existing and non-existing paths
    if [[ "$path" = /* ]]; then
        abs_path="$path"
    else
        abs_path="$PWD/$path"
    fi

    abs_home=$(realpath "$HOME" 2>/dev/null || echo "$HOME")

    if [[ -z "$abs_path" || -z "$abs_home" ]]; then
        error "Invalid directory path: $path"
        exit 1
    fi

    # Check if path is under HOME directory
    case "$abs_path" in
        "$abs_home"/*|"$abs_home") ;;
        *) error "Installation directory must be under HOME: $path"; exit 1 ;;
    esac
}

# Installation functions
show_banner() {
    cat << 'EOF'
🚀 MCF (My Claude Framework) Installer
=====================================

This installer will:
• Clone the MCF repository to a temporary location
• Create the MCF directory in your home folder
• Install core MCF components:
  - claude-mcf.sh (main script)
  - README.md (documentation)
  - scripts/ (utility scripts)
  - templates/ (project templates)
  - .claude/ (configuration)
• Set up proper permissions
• Clean up temporary files

EOF
}

get_user_consent() {
    cat << EOF
📋 Installation Details:
• Repository: $MCF_REPO
• Install Directory: $MCF_DIR
• Log File: $LOG_FILE

EOF

    # Check if running interactively
    if [[ -t 0 && -t 1 ]]; then
        # Interactive mode - ask for confirmation
        read -p "Do you want to proceed with the installation? [y/N]: " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            info "Installation cancelled by user"
            exit 0
        fi
    else
        # Non-interactive mode (pipe) - auto-proceed with warning
        warn "Running in non-interactive mode - proceeding automatically"
        info "Auto-proceeding with installation (use --help for options)"
        sleep 2
    fi
}

clone_repository() {
    local temp_dir="$1"
    local repo_dir="$temp_dir/mcf-repo"

    # Check if we're already in the MCF repository
    if git remote -v 2>/dev/null | grep -q "pc-style/MCF.git" && [[ -f "claude-mcf.sh" && -f "README.md" && -d "scripts" && -d "templates" && -d ".claude" ]]; then
        info "Already in MCF repository - using current directory..." >&2
        log "Using current directory as source: $PWD"
        echo "$PWD"
        return
    fi

    info "Cloning MCF repository..." >&2
    log "Cloning from $MCF_REPO to $repo_dir"

    if ! git clone --depth 1 --quiet "$MCF_REPO" "$repo_dir" 2>>"$LOG_FILE"; then
        error "Failed to clone repository: $MCF_REPO"
        error "Check the log file: $LOG_FILE"
        exit 1
    fi

    echo "$repo_dir"
}

verify_required_files() {
    local repo_dir="$1"

    info "Verifying required files..."
    log "Checking required files in $repo_dir"

    local missing_files=()
    for file in "${REQUIRED_FILES[@]}"; do
        if [[ ! -e "$repo_dir/$file" ]]; then
            missing_files+=("$file")
        fi
    done

    if [[ ${#missing_files[@]} -gt 0 ]]; then
        error "Missing required files: ${missing_files[*]}"
        exit 1
    fi

    log "All required files verified"
}

backup_existing_installation() {
    if [[ -d "$MCF_DIR" ]]; then
        local backup_dir="${MCF_DIR}.backup.$(date +%s)"
        warn "Existing MCF installation found"
        info "Creating backup: $backup_dir"
        
        if ! mv "$MCF_DIR" "$backup_dir" 2>>"$LOG_FILE"; then
            error "Failed to create backup of existing installation"
            exit 1
        fi
        
        log "Backup created: $backup_dir"
    fi
}

install_mcf_files() {
    local repo_dir="$1"
    
    info "Installing MCF components..."
    log "Installing from $repo_dir to $MCF_DIR"
    
    # Create MCF directory
    if ! mkdir -p "$MCF_DIR" 2>>"$LOG_FILE"; then
        error "Failed to create MCF directory: $MCF_DIR"
        exit 1
    fi
    
    # Copy required files and directories
    for item in "${REQUIRED_FILES[@]}"; do
        info "Installing: $item"
        if ! cp -r "$repo_dir/$item" "$MCF_DIR/" 2>>"$LOG_FILE"; then
            error "Failed to install: $item"
            exit 1
        fi
        log "Installed: $item"
    done
}

install_serena_mcp_server() {
    info "Installing Serena MCP server..."
    log "Installing Serena MCP server in $MCF_DIR/.claude"
    CLAUDE_CONFIG_DIR="$MCF_DIR/.claude" claude mcp add --transport http context7 https://mcp.context7.com/mcp
    CLAUDE_CONFIG_DIR="$MCF_DIR/.claude" claude mcp add serena --scope user -- uvx --from git+https://github.com/oraios/serena serena start-mcp-server
}

setup_permissions() {
    info "Setting up permissions..."
    log "Setting up permissions in $MCF_DIR"
    
    # Make main script executable
    if [[ -f "$MCF_DIR/claude-mcf.sh" ]]; then
        chmod +x "$MCF_DIR/claude-mcf.sh" 2>>"$LOG_FILE"
        log "Made claude-mcf.sh executable"
    fi
    
    # Set permissions for scripts directory
    if [[ -d "$MCF_DIR/scripts" ]]; then
        find "$MCF_DIR/scripts" -name "*.sh" -exec chmod +x {} \; 2>>"$LOG_FILE"
        log "Made scripts executable"
    fi
    
    # Set permissions for hooks
    if [[ -d "$MCF_DIR/.claude/hooks" ]]; then
        find "$MCF_DIR/.claude/hooks" -name "*.py" -exec chmod +x {} \; 2>>"$LOG_FILE"
        find "$MCF_DIR/.claude/hooks" -name "*.sh" -exec chmod +x {} \; 2>>"$LOG_FILE"
        log "Made hooks executable"
    fi
    
    # Protect sensitive configuration files
    if [[ -f "$MCF_DIR/.claude/settings.json" ]]; then
        chmod 600 "$MCF_DIR/.claude/settings.json" 2>>"$LOG_FILE"
        log "Protected settings.json"
    fi
}

setup_shell_integration() {
    info "Setting up shell integration..."

    # Create ~/.local/bin directory if it doesn't exist
    local local_bin_dir="$HOME/.local/bin"
    if [[ ! -d "$local_bin_dir" ]]; then
        info "Creating ~/.local/bin directory..."
        if ! mkdir -p "$local_bin_dir" 2>>"$LOG_FILE"; then
            warn "Failed to create ~/.local/bin directory"
        else
            log "Created ~/.local/bin directory"
        fi
    fi

    # Copy claude-mcf.sh to ~/.local/bin/mcf
    if [[ -f "$MCF_DIR/claude-mcf.sh" ]]; then
        info "Installing mcf command to ~/.local/bin..."
        if cp "$MCF_DIR/claude-mcf.sh" "$local_bin_dir/mcf" 2>>"$LOG_FILE"; then
            chmod +x "$local_bin_dir/mcf" 2>>"$LOG_FILE"
            log "Installed mcf command to ~/.local/bin/mcf"

            # Check if ~/.local/bin is in PATH
            if ! echo "$PATH" | tr ':' '\n' | grep -q "^$local_bin_dir$"; then
                warn "~/.local/bin is not in your PATH"
                warn "You may need to add it to your shell profile:"
                warn "  export PATH=\"$local_bin_dir:\$PATH\""
            fi
        else
            warn "Failed to install mcf command to ~/.local/bin"
        fi
    fi

    # Create convenient mcf alias as fallback
    local alias_line="alias mcf='$MCF_DIR/claude-mcf.sh'"
    local shell_config=""

    # Detect shell and config file
    case "$SHELL" in
        */bash) shell_config="$HOME/.bashrc" ;;
        */zsh) shell_config="$HOME/.zshrc" ;;
        */fish) shell_config="$HOME/.config/fish/config.fish" ;;
        *) warn "Unknown shell: $SHELL. Please manually add alias: $alias_line" ;;
    esac

    if [[ -n "$shell_config" && -f "$shell_config" ]]; then
        if ! grep -q "alias mcf=" "$shell_config" 2>/dev/null; then
            echo "" >> "$shell_config"
            echo "# MCF alias (fallback)" >> "$shell_config"
            echo "$alias_line" >> "$shell_config"
            info "Added mcf alias to $shell_config (as fallback)"
            log "Added alias to $shell_config"
        else
            log "Alias already exists in $shell_config"
        fi
    fi
}

show_completion_summary() {
    success "MCF installation completed successfully!"

    cat << EOF

📁 Installation Summary:
• MCF Directory: $MCF_DIR
• Main Command: ~/.local/bin/mcf (in PATH)
• Fallback Script: $MCF_DIR/claude-mcf.sh
• Documentation: $MCF_DIR/README.md
• Configuration: $MCF_DIR/.claude/
• Log File: $LOG_FILE

🚀 Next Steps:
1. Restart your terminal or run: source ~/.bashrc (or ~/.zshrc)
2. Use 'mcf' command anywhere in your terminal
3. Read the documentation: $MCF_DIR/README.md

💡 Quick Start:
   mcf --help                 # Show help
   mcf init                   # Initialize a new project
   mcf status                 # Check system status

🔧 Troubleshooting:
• Check the log file: $LOG_FILE
• If 'mcf' command not found, ensure ~/.local/bin is in your PATH:
  export PATH="\$HOME/.local/bin:\$PATH"
• Visit: https://github.com/pc-style/MCF/issues

EOF
}

# Trap for cleanup
cleanup_on_exit() {
    local exit_code=$?
    if [[ -n "${TEMP_DIR:-}" ]]; then
        secure_cleanup "$TEMP_DIR"
    fi
    
    if [[ $exit_code -ne 0 ]]; then
        error "Installation failed with exit code: $exit_code"
        error "Check the log file for details: $LOG_FILE"
    fi
}

# Main installation function
main() {
    local temp_dir repo_dir
    local auto_yes=false

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --yes|-y)
                auto_yes=true
                shift
                ;;
            --help|-h)
                cat << EOF
MCF (My Claude Framework) Installer - Version $SCRIPT_VERSION

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --yes, -y    Skip interactive prompts and proceed automatically
    --help, -h   Show this help message

EXAMPLES:
    $0                    # Interactive installation
    $0 --yes             # Non-interactive installation
    curl -fsSL https://raw.githubusercontent.com/pc-style/MCF/main/install.sh | bash -s -- --yes

For more information, visit: https://github.com/pc-style/MCF
EOF
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                error "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    # Set up cleanup trap
    trap cleanup_on_exit EXIT INT TERM

    # Start installation
    log "MCF installation started (version $SCRIPT_VERSION)"

    # Security and validation
    verify_environment
    sanitize_environment
    validate_repo_url "$MCF_REPO"
    validate_directory_path "$MCF_DIR"

    # User interaction
    show_banner
    if [[ "$auto_yes" == "true" ]]; then
        info "Auto-yes mode enabled - skipping confirmation prompt"
    else
        get_user_consent
    fi
    
    # Create secure workspace
    temp_dir=$(create_secure_temp)
    TEMP_DIR="$temp_dir"  # For cleanup trap
    
    # Installation process
    repo_dir=$(clone_repository "$temp_dir")
    verify_required_files "$repo_dir"
    backup_existing_installation
    install_mcf_files "$repo_dir"
    setup_permissions
    install_serena_mcp_server
    setup_shell_integration
    
    # Complete installation
    log "MCF installation completed successfully"
    show_completion_summary
}

# Execute main function
main "$@"
