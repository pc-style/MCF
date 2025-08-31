#!/bin/bash

# MCF TUI with Local Claude Proxy
# This script starts the local Claude proxy and runs the TUI with proper environment

set -e

echo "🚀 Starting MCF TUI with Local Claude Proxy..."

# Check if proxy is already running
if curl -s http://localhost:4141/health >/dev/null 2>&1; then
    echo "✅ Claude proxy already running on localhost:4141"
else
    echo "🔄 Starting Claude proxy..."
    # Start the proxy in background
    nohup npx copilot-api@latest start --claude-code >/dev/null 2>&1 &
    PROXY_PID=$!
    
    # Wait for proxy to start
    echo "⏳ Waiting for proxy to start..."
    for i in {1..30}; do
        if curl -s http://localhost:4141/health >/dev/null 2>&1; then
            echo "✅ Claude proxy started successfully"
            break
        fi
        sleep 1
        if [ $i -eq 30 ]; then
            echo "❌ Failed to start Claude proxy"
            exit 1
        fi
    done
fi

# Set environment variables and run TUI
export ANTHROPIC_BASE_URL=http://localhost:4141
export ANTHROPIC_AUTH_TOKEN=dummy
export CLAUDE_CONFIG_DIR="$HOME/mcf-dev/.claude"
export ANTHROPIC_MODEL=claude-3.5-sonnet
export ANTHROPIC_SMALL_FAST_MODEL=grok-code-fast-1

echo "🎯 Starting MCF TUI..."
echo "📝 Logs will be saved to: $HOME/mcf-dev/logs/"
echo "🔧 Debug mode: Use --debug flag for console output"
echo ""

# Run the TUI with the configured environment
./mcf-tui "$@"

echo ""
echo "👋 MCF TUI session ended"
