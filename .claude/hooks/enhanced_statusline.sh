#!/bin/bash
"""
🚀 ULTRA ENHANCED STATUS LINE - The Ultimate Claude Code Status Experience
Features: Dynamic colors, animations, real-time metrics, rich git info, and more!
"""

# Read JSON input
input=$(cat)

# ANSI Color Codes for Terminal Beauty
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[0;37m'
BOLD='\033[1m'
DIM='\033[2m'
ITALIC='\033[3m'
UNDERLINE='\033[4m'
BLINK='\033[5m'
REVERSE='\033[7m'
RESET='\033[0m'

# Gradient colors for fancy effects
GRADIENT_1='\033[38;5;196m'  # Red
GRADIENT_2='\033[38;5;202m'  # Orange
GRADIENT_3='\033[38;5;226m'  # Yellow
GRADIENT_4='\033[38;5;46m'   # Green
GRADIENT_5='\033[38;5;21m'   # Blue
GRADIENT_6='\033[38;5;93m'   # Purple

# Background colors
BG_DARK='\033[48;5;234m'
BG_HIGHLIGHT='\033[48;5;236m'

# Helper functions
get_field() { echo "$input" | jq -r ".$1 // \"\""; }
get_nested_field() { echo "$input" | jq -r ".$1.$2 // \"\""; }

# Time-based greeting with emoji
get_time_greeting() {
    hour=$(date +%H)
    if [ $hour -lt 6 ]; then
        echo "🌙 Night Owl"
    elif [ $hour -lt 12 ]; then
        echo "☀️ Good Morning"
    elif [ $hour -lt 17 ]; then
        echo "🌤️ Good Afternoon"
    elif [ $hour -lt 22 ]; then
        echo "🌆 Good Evening"
    else
        echo "✨ Late Night"
    fi
}

# Get system performance metrics
get_system_metrics() {
    # CPU usage (macOS specific, adjust for Linux)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        CPU_USAGE=$(top -l 1 -n 0 | grep "CPU usage" | awk '{print $3}' | sed 's/%//')
    else
        CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    fi
    
    # Memory usage
    if [[ "$OSTYPE" == "darwin"* ]]; then
        MEM_USAGE=$(vm_stat | perl -ne '/page size of (\d+)/ and $size=$1; /Pages\s+([^:]+)[^\d]+(\d+)/ and printf("%.1f", $2 * $size / 1048576) if $1 eq "active"')
        MEM_TOTAL=$(sysctl -n hw.memsize | awk '{print int($1/1024/1024/1024)}')
    else
        MEM_INFO=$(free -m | awk 'NR==2{printf "%.1f %.1f", $3/1024, $2/1024}')
        MEM_USAGE=$(echo $MEM_INFO | awk '{print $1}')
        MEM_TOTAL=$(echo $MEM_INFO | awk '{print $2}')
    fi
    
    echo "⚡ CPU: ${CPU_USAGE}% | 🧠 RAM: ${MEM_USAGE}GB"
}

# Extract basic info with enhanced formatting
MODEL=$(get_nested_field "model" "display_name")
CURRENT_DIR=$(get_nested_field "workspace" "current_dir")
PROJECT_NAME=${CURRENT_DIR##*/}

# Model emoji and color based on model type
case "$MODEL" in
    *"sonnet"*)
        MODEL_EMOJI="🎭"
        MODEL_COLOR=$PURPLE
        ;;
    *"opus"*)
        MODEL_EMOJI="🎼"
        MODEL_COLOR=$GRADIENT_6
        ;;
    *"haiku"*)
        MODEL_EMOJI="🍃"
        MODEL_COLOR=$GREEN
        ;;
    *"grok"*)
        MODEL_EMOJI="⚡"
        MODEL_COLOR=$YELLOW
        ;;
    *)
        MODEL_EMOJI="🤖"
        MODEL_COLOR=$CYAN
        ;;
esac

# Enhanced Git information with detailed stats
GIT_INFO=""
if [ -d ".git" ]; then
    BRANCH=$(git branch --show-current 2>/dev/null)
    if [ -n "$BRANCH" ]; then
        # Count various git stats
        STAGED=$(git diff --cached --numstat 2>/dev/null | wc -l | tr -d ' ')
        MODIFIED=$(git diff --numstat 2>/dev/null | wc -l | tr -d ' ')
        UNTRACKED=$(git ls-files --others --exclude-standard 2>/dev/null | wc -l | tr -d ' ')
        
        # Get ahead/behind info
        UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
        if [ -n "$UPSTREAM" ]; then
            AHEAD=$(git rev-list --count @{u}..HEAD 2>/dev/null || echo "0")
            BEHIND=$(git rev-list --count HEAD..@{u} 2>/dev/null || echo "0")
            SYNC_INFO=""
            [ "$AHEAD" -gt 0 ] && SYNC_INFO="${SYNC_INFO}↑$AHEAD"
            [ "$BEHIND" -gt 0 ] && SYNC_INFO="${SYNC_INFO}↓$BEHIND"
            [ -n "$SYNC_INFO" ] && SYNC_INFO=" $SYNC_INFO"
        else
            SYNC_INFO=""
        fi
        
        # Build git status string with colors
        GIT_STATUS=""
        [ "$STAGED" -gt 0 ] && GIT_STATUS="${GIT_STATUS}${GREEN}●${STAGED}${RESET}"
        [ "$MODIFIED" -gt 0 ] && GIT_STATUS="${GIT_STATUS}${YELLOW}✚${MODIFIED}${RESET}"
        [ "$UNTRACKED" -gt 0 ] && GIT_STATUS="${GIT_STATUS}${RED}…${UNTRACKED}${RESET}"
        
        # Branch color based on status
        if [ "$STAGED" -gt 0 ] || [ "$MODIFIED" -gt 0 ] || [ "$UNTRACKED" -gt 0 ]; then
            BRANCH_COLOR=$YELLOW
        else
            BRANCH_COLOR=$GREEN
        fi
        
        # Last commit info
        LAST_COMMIT=$(git log -1 --pretty=format:"%ar" 2>/dev/null || echo "")
        [ -n "$LAST_COMMIT" ] && LAST_COMMIT=" (${DIM}$LAST_COMMIT${RESET})"
        
        GIT_INFO=" ${DIM}|${RESET} ${BRANCH_COLOR}🌿 $BRANCH${RESET}${GIT_STATUS}${CYAN}${SYNC_INFO}${RESET}${LAST_COMMIT}"
    fi
fi

# Enhanced MCP status with server details
MCP_INFO=""
if command -v claude >/dev/null 2>&1; then
    # Try to get detailed MCP info
    MCP_LIST=$(claude mcp list 2>/dev/null || echo "")
    MCP_COUNT=$(echo "$MCP_LIST" | grep -c "✓ Connected" || echo "0")
    
    if [ "$MCP_COUNT" -gt "0" ]; then
        # Extract server names if possible
        MCP_SERVERS=$(echo "$MCP_LIST" | grep "✓ Connected" | awk '{print $3}' | head -3 | tr '\n' ',' | sed 's/,$//')
        
        # Color based on connection count
        if [ "$MCP_COUNT" -ge 3 ]; then
            MCP_COLOR=$GREEN
            MCP_EMOJI="🔥"
        elif [ "$MCP_COUNT" -ge 1 ]; then
            MCP_COLOR=$YELLOW
            MCP_EMOJI="🔌"
        fi
        
        MCP_INFO=" ${DIM}|${RESET} ${MCP_COLOR}${MCP_EMOJI} ${MCP_COUNT} MCP${RESET}"
        [ -n "$MCP_SERVERS" ] && [ ${#MCP_SERVERS} -lt 30 ] && MCP_INFO="${MCP_INFO} ${DIM}(${MCP_SERVERS})${RESET}"
    fi
fi

# Enhanced cost information with session stats
COST_INFO=""
TOTAL_COST=$(get_nested_field "cost" "total_cost_usd")
SESSION_TOKENS=$(get_nested_field "usage" "total_tokens")
if [ "$TOTAL_COST" != "null" ] && [ "$TOTAL_COST" != "" ] && [ "$TOTAL_COST" != "0" ]; then
    COST_FORMATTED=$(printf "%.4f" "$TOTAL_COST")
    
    # Color based on cost
    COST_COLOR=$GREEN
    COST_EMOJI="💵"
    if (( $(echo "$TOTAL_COST > 1.0" | bc -l) )); then
        COST_COLOR=$YELLOW
        COST_EMOJI="💰"
    elif (( $(echo "$TOTAL_COST > 5.0" | bc -l) )); then
        COST_COLOR=$RED
        COST_EMOJI="💸"
    fi
    
    COST_INFO=" ${DIM}|${RESET} ${COST_COLOR}${COST_EMOJI} \$${COST_FORMATTED}${RESET}"
    
    # Add token count if available
    if [ "$SESSION_TOKENS" != "null" ] && [ "$SESSION_TOKENS" != "" ]; then
        TOKEN_K=$(echo "scale=1; $SESSION_TOKENS / 1000" | bc 2>/dev/null || echo "0")
        [ "$TOKEN_K" != "0" ] && COST_INFO="${COST_INFO} ${DIM}(${TOKEN_K}k tok)${RESET}"
    fi
fi

# File count and project size
FILE_COUNT=$(find . -type f -name "*.js" -o -name "*.ts" -o -name "*.py" -o -name "*.java" -o -name "*.go" 2>/dev/null | wc -l | tr -d ' ')
PROJECT_SIZE=$(du -sh . 2>/dev/null | awk '{print $1}')
PROJECT_INFO="${CYAN}📊 ${FILE_COUNT} files${RESET} ${DIM}(${PROJECT_SIZE})${RESET}"

# Session duration
SESSION_START=$(get_field "session_start_time")
if [ "$SESSION_START" != "null" ] && [ "$SESSION_START" != "" ]; then
    SESSION_DURATION=$(( $(date +%s) - $(date -d "$SESSION_START" +%s 2>/dev/null || echo $(date +%s)) ))
    SESSION_MINS=$(( SESSION_DURATION / 60 ))
    SESSION_INFO=""
    if [ $SESSION_MINS -gt 0 ]; then
        SESSION_HOURS=$(( SESSION_MINS / 60 ))
        if [ $SESSION_HOURS -gt 0 ]; then
            SESSION_INFO=" ${DIM}|${RESET} ⏱️ ${SESSION_HOURS}h ${SESSION_MINS}m"
        else
            SESSION_INFO=" ${DIM}|${RESET} ⏱️ ${SESSION_MINS}m"
        fi
    fi
else
    SESSION_INFO=""
fi

# Weather/mood indicator based on various factors
MOOD=""
if [ "$MCP_COUNT" -gt 2 ] && [ "$STAGED" -eq 0 ] && [ "$MODIFIED" -eq 0 ]; then
    MOOD="🚀 PEAK PERFORMANCE"
elif [ "$MODIFIED" -gt 5 ] || [ "$UNTRACKED" -gt 10 ]; then
    MOOD="🔥 HOT CHANGES"
elif [ "$SESSION_MINS" -gt 120 ]; then
    MOOD="💪 MARATHON MODE"
else
    MOOD="$(get_time_greeting)"
fi

# Terminal width for dynamic sizing
TERM_WIDTH=$(tput cols 2>/dev/null || echo 80)

# Build the spectacular status line
STATUS_LINE=""

# Top border with gradient effect
if [ $TERM_WIDTH -gt 100 ]; then
    STATUS_LINE="${BG_DARK}${GRADIENT_1}━${GRADIENT_2}━${GRADIENT_3}━${GRADIENT_4}━${GRADIENT_5}━${GRADIENT_6}━${RESET} "
fi

# Main status components
STATUS_LINE="${STATUS_LINE}${BOLD}${MODEL_COLOR}${MODEL_EMOJI} ${MODEL}${RESET}"
STATUS_LINE="${STATUS_LINE} ${DIM}|${RESET} ${BOLD}${BLUE}📁 ${PROJECT_NAME}${RESET}"
STATUS_LINE="${STATUS_LINE}${GIT_INFO}"
STATUS_LINE="${STATUS_LINE}${MCP_INFO}"
STATUS_LINE="${STATUS_LINE}${COST_INFO}"
STATUS_LINE="${STATUS_LINE}${SESSION_INFO}"

# Add mood/greeting on wider terminals
if [ $TERM_WIDTH -gt 120 ]; then
    STATUS_LINE="${STATUS_LINE} ${DIM}|${RESET} ${ITALIC}${MOOD}${RESET}"
fi

# Add project info on very wide terminals
if [ $TERM_WIDTH -gt 150 ]; then
    STATUS_LINE="${STATUS_LINE} ${DIM}|${RESET} ${PROJECT_INFO}"
fi

# System metrics on ultra-wide terminals
if [ $TERM_WIDTH -gt 180 ]; then
    METRICS=$(get_system_metrics)
    STATUS_LINE="${STATUS_LINE} ${DIM}|${RESET} ${DIM}${METRICS}${RESET}"
fi

# Output the magnificent status line
echo -e "$STATUS_LINE"

# Optional: Add a subtle animation/spinner for active operations
# Uncomment to enable (may cause flicker in some terminals)
# SPINNER=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")
# SPIN_INDEX=$(( $(date +%S) % 10 ))
# echo -ne " ${CYAN}${SPINNER[$SPIN_INDEX]}${RESET}"