#!/bin/bash
"""
Test script to demonstrate the enhanced statusline with sample data
"""

# Sample JSON data that would normally come from Claude Code
create_test_json() {
    cat <<EOF
{
    "model": {
        "display_name": "claude-3.5-sonnet"
    },
    "workspace": {
        "current_dir": "$(pwd)"
    },
    "cost": {
        "total_cost_usd": 0.0234
    },
    "usage": {
        "total_tokens": 15432
    },
    "session_start_time": "$(date -d '45 minutes ago' 2>/dev/null || date)"
}
EOF
}

echo "🎨 ENHANCED STATUSLINE DEMO"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test 1: Basic statusline
echo "📍 Test 1: Basic Status"
create_test_json | ./.claude/hooks/enhanced_statusline.sh
echo ""

# Test 2: Different model
echo "📍 Test 2: With Opus Model"
create_test_json | jq '.model.display_name = "claude-3-opus"' | ./.claude/hooks/enhanced_statusline.sh
echo ""

# Test 3: High cost warning
echo "📍 Test 3: High Cost Warning"
create_test_json | jq '.cost.total_cost_usd = 2.5678' | ./.claude/hooks/enhanced_statusline.sh
echo ""

# Test 4: Long session
echo "📍 Test 4: Marathon Coding Session"
create_test_json | jq '.session_start_time = "'$(date -d '3 hours ago' 2>/dev/null || date)'"' | ./.claude/hooks/enhanced_statusline.sh
echo ""

# Test 5: With many tokens
echo "📍 Test 5: Heavy Token Usage"
create_test_json | jq '.usage.total_tokens = 125000 | .cost.total_cost_usd = 8.9234' | ./.claude/hooks/enhanced_statusline.sh
echo ""

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✨ Features demonstrated:"
echo "  • Dynamic model detection with custom emojis"
echo "  • Git status with detailed change counts"
echo "  • MCP connection monitoring"
echo "  • Cost tracking with color-coded warnings"
echo "  • Session duration tracking"
echo "  • Token usage display"
echo "  • Time-based greetings"
echo "  • Project statistics"
echo "  • Responsive width adaptation"
echo ""
echo "🎯 The statusline adapts to terminal width:"
echo "  • 80+ chars:  Compact mode"
echo "  • 120+ chars: Add mood/greeting"
echo "  • 150+ chars: Add project stats"
echo "  • 180+ chars: Add system metrics"
