---
description: Show all available Serena commands and usage guide
---

Complete guide to MCF's Serena integration commands.

## 🚀 **MCF Serena Command Suite**

Your MCF setup includes comprehensive Serena semantic code analysis integration. Here are all available commands:

### **🛠️ Setup & Management Commands**

**`/serena:install`** - Install and configure Serena MCP server
- Adds Serena to Claude Code configuration  
- Creates optimized settings for MCF workflow
- Sets up proper analytics and monitoring

**`/serena:config`** - Manage Serena configuration  
- Check current settings
- Update with MCF optimizations
- Reset to defaults

**`/serena:status`** - Complete integration health check
- Verify MCP server connection
- Check configuration status  
- Validate project setup
- Diagnose common issues

### **📂 Project Commands**

**`/serena:init`** - Initialize current project for semantic analysis
- Generates project configuration using Serena's built-in tools
- Creates MCF-specific memories and settings
- Sets up language-specific patterns

**`/serena:activate`** - Activate project for Serena analysis  
- Ensures proper project initialization
- Triggers semantic indexing process
- Prepares tools for current codebase

### **🔍 Analysis & Navigation Commands**

**`/serena:overview`** - Get semantic project structure overview
- High-level project organization
- Key modules and entry points
- Architectural patterns

**`/serena:find [symbol]`** - Find symbols using semantic search
- Locate functions, classes, variables by name
- Type-aware searching
- Cross-reference discovery

**`/serena:analyze [symbol]`** - Deep analysis of code symbols
- Detailed symbol information
- Dependencies and relationships
- Documentation and context

**`/serena:refs [symbol]`** - Find all references to a symbol
- Usage locations throughout codebase  
- Import/export tracking
- Call chain analysis

**`/serena:help`** - This comprehensive guide

## 🎯 **Quick Start Workflow**

If you're new to Serena with MCF:

```bash
# 1. Install Serena MCP server
/serena:install

# 2. Initialize current project  
/serena:init

# 3. Check everything is working
/serena:status

# 4. Activate project for analysis
/serena:activate

# 5. Start exploring your code
/serena:overview
```

## 🛠️ **Troubleshooting Guide**

**Connection Issues:**
- Run `/serena:status` for complete diagnostics
- Check `/mcp` for active MCP servers
- Restart Claude Code if needed

**Configuration Problems:**
- Use `/serena:config check` to see current settings
- Run `/serena:config update` for MCF optimizations
- Check `~/.serena/serena_config.yml` manually

**Project Issues:**
- Ensure `/serena:init` completed successfully
- Verify `.serena/project.yml` exists
- Re-run `/serena:activate` if needed

## 🌟 **Advanced Features**

### **Semantic Analysis Tools Available:**
- `find_symbol` - IDE-like symbol search
- `get_symbol_info` - Detailed symbol information  
- `find_referencing_symbols` - Complete reference tracking
- `get_project_structure` - Architectural overview
- `insert_before_symbol` / `insert_after_symbol` - Precise editing

### **MCF Integration Benefits:**
- **Token Efficiency** - Work at symbol level instead of full files
- **Agent Enhancement** - All 8+ MCF agents get semantic superpowers
- **Smart Suggestions** - Context-aware recommendations via hooks
- **Usage Analytics** - Track which tools provide most value

### **Monitoring & Analytics:**
- **Web Dashboard** - http://localhost:24282/dashboard/
- **Usage Statistics** - Track tool effectiveness  
- **Performance Metrics** - Monitor analysis speed
- **Token Savings** - See efficiency improvements

## 🎨 **Integration with MCF Agents**

Your specialized agents are enhanced with Serena:

- **micro-analyzer** - Symbol-level code analysis
- **security-auditor** - Data flow tracing and vulnerability detection
- **perf-optimizer** - Function-level performance bottleneck identification  
- **api-designer** - Endpoint-to-implementation mapping
- **semantic-navigator** - Pure semantic code exploration

## 📊 **Configuration Location**

**Global Config:** `~/.serena/serena_config.yml`  
**Project Config:** `.serena/project.yml`  
**Project Memories:** `.serena/memories/`

## 🚀 **Next Steps**

1. Run `/serena:status` to check your current setup
2. Use semantic commands to explore your codebase
3. Monitor usage at the web dashboard
4. Let your MCF agents leverage semantic superpowers!

**Your MCF + Serena integration provides IDE-like semantic code understanding directly in Claude Code!** 🎯
