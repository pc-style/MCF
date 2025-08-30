#!/bin/bash

PROJECT_DIR=$(pwd)
INSTALL_DIR="$HOME/mcf"
INSTALL_CMD="curl -fsSL https://raw.githubusercontent.com/pc-style/MCF/main/install.sh | bash"

if [ ! -f "$INSTALL_DIR/.claude/settings.json" ]; then
    echo "❌ Claude MCF not found"
    echo "Please run '$INSTALL_CMD'"
    exit 1
fi
mkdir -p "$INSTALL_DIR/.claude/bookmarks"

#if first run, we need to authenticate
if [ ! -f "$INSTALL_DIR/.claude/bookmarks/.first-run.txt" ]; then
    echo "🚀 First run detected"
    echo ""
    echo "!! YOU WILL HAVE TO LOG IN AGAIN !!"
    sleep 2
    echo "🚀 Running Claude MCF..."
fi

RUN_CMD="cd $PROJECT_DIR && CLAUDE_CONFIG_DIR=$INSTALL_DIR/.claude claude"
eval "$RUN_CMD"
touch "$INSTALL_DIR/.claude/bookmarks/.first-run.txt"

echo "✅ Claude MCF ready"