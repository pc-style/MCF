#!/bin/bash
"""
Enhanced status line showing model, project, git status, and MCP connections.
"""

# Read JSON input
input=$(cat)

# Helper functions
get_field() { echo "$input" | jq -r ".$1 // \"\""; }
get_nested_field() { echo "$input" | jq -r ".$1.$2 // \"\""; }

# Extract basic info
MODEL=$(get_nested_field "model" "display_name")
CURRENT_DIR=$(get_nested_field "workspace" "current_dir")
PROJECT_NAME=${CURRENT_DIR##*/}

# Git information
GIT_INFO=""
if [ -d ".git" ]; then
    BRANCH=$(git branch --show-current 2>/dev/null)
    if [ -n "$BRANCH" ]; then
        # Check for uncommitted changes
        if ! git diff-index --quiet HEAD 2>/dev/null; then
            GIT_STATUS="*"
        else
            GIT_STATUS=""
        fi
        GIT_INFO=" | 🌿 $BRANCH$GIT_STATUS"
    fi
fi

# MCP status indicator
MCP_INFO=""
if command -v claude >/dev/null 2>&1; then
    MCP_COUNT=$(claude mcp list 2>/dev/null | grep -c "✓ Connected" || echo "0")
    if [ "$MCP_COUNT" -gt "0" ]; then
        MCP_INFO=" | 🔌 $MCP_COUNT MCP"
    fi
fi

# Cost information (if available)
COST_INFO=""
TOTAL_COST=$(get_nested_field "cost" "total_cost_usd")
if [ "$TOTAL_COST" != "null" ] && [ "$TOTAL_COST" != "" ]; then
    COST_FORMATTED=$(printf "%.3f" "$TOTAL_COST")
    COST_INFO=" | 💰 \$${COST_FORMATTED}"
fi

# Combine all parts
echo "[$MODEL] 📁 $PROJECT_NAME$GIT_INFO$MCP_INFO$COST_INFO"
