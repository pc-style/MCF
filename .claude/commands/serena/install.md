---
description: Install and configure Serena MCP server for Claude Code
---

Install and configure Serena MCP server for this project.

This command will:
1. Add Serena MCP server to Claude Code configuration
2. Set up optimal configuration for MCF workflow
3. Create project-specific settings and memories

**Installation Steps:**

First, let me add the Serena MCP server to Claude Code:

```bash
claude mcp add serena --scope user -- uvx --from git+https://github.com/oraios/serena serena start-mcp-server
```

Now I'll create an optimized Serena configuration:

```bash
# Create config directory
mkdir -p ~/.serena

# Create optimized config for MCF
cat > ~/.serena/serena_config.yml << 'EOF'
# Optimized Serena Configuration for MCF
gui_log_window: false
web_dashboard: true
web_dashboard_open_on_launch: false
log_level: 20
tool_timeout: 300
default_max_tool_answer_chars: 200000
record_tool_usage_stats: true
token_count_estimator: TIKTOKEN_GPT4O
projects: []
EOF
```

✅ **Serena MCP server installed successfully!**

**Next Steps:**
1. Restart Claude Code to load the new MCP server
2. Run `/serena:init` to initialize the project
3. Check connection with `/mcp` command
4. Monitor at http://localhost:24282/dashboard/ when active

**Verify Installation:**
```bash
# Check if Serena is configured in Claude Code
claude mcp list | grep serena

# Check if Serena config exists
ls -la ~/.serena/serena_config.yml

# Test Serena accessibility
uvx --from git+https://github.com/oraios/serena serena --help
```

**Configuration Location:** `~/.serena/serena_config.yml`

Your MCF workflow now has IDE-like semantic code capabilities!

💡 **Pro Tip:** Use `/serena:config` to check and manage your Serena settings anytime.
