---
description: Configure Serena settings and check current configuration
argument-hint: [check|update|reset]
---

Manage Serena configuration settings.

**Usage:**
- `/serena:config` or `/serena:config check` - Display current configuration
- `/serena:config update` - Update configuration with MCF optimizations
- `/serena:config reset` - Reset to default configuration

**Current Configuration Location:** `~/.serena/serena_config.yml`

Let me check your current Serena configuration:

```bash
# Check if Serena config exists
if [ -f ~/.serena/serena_config.yml ]; then
    echo "✅ Serena configuration found"
    echo "📍 Location: ~/.serena/serena_config.yml"
    echo ""
    echo "📋 Current settings:"
    cat ~/.serena/serena_config.yml
else
    echo "❌ Serena configuration not found"
    echo "💡 Run '/serena:install' to create initial configuration"
fi
```

$ARGUMENTS contains: $ARGUMENTS

If you want to update the configuration with MCF optimizations:

```bash
# Create MCF-optimized configuration
mkdir -p ~/.serena

cat > ~/.serena/serena_config.yml << 'EOF'
# Optimized Serena Configuration for MCF Development Workflow
gui_log_window: false
web_dashboard: true
web_dashboard_open_on_launch: false
log_level: 20
trace_lsp_communication: false
tool_timeout: 300
excluded_tools: []
included_optional_tools: []
jetbrains: false
default_max_tool_answer_chars: 200000
record_tool_usage_stats: true
token_count_estimator: TIKTOKEN_GPT4O

# Projects will be managed automatically
projects: []
EOF

echo "✅ Updated Serena configuration with MCF optimizations"
```

**Key MCF Optimizations:**
- 📊 **Usage Analytics Enabled** - Track which tools are most valuable
- ⏱️ **Extended Timeout** - 5 minutes for complex analysis
- 📝 **Larger Response Size** - 200k chars for detailed code info
- 🌐 **Web Dashboard** - Monitor without auto-opening browser
- 🔍 **Info Level Logging** - Balanced detail without spam

**Next Steps:**
1. Restart Claude Code if you updated configuration
2. Test with `/mcp` to verify Serena connection
3. Monitor usage at http://localhost:24282/dashboard/
