---
description: Check Serena installation and connection status
---

Check the complete status of your Serena integration with MCF.

This will verify:
- Serena MCP server configuration in Claude Code
- Serena configuration file
- Project initialization status  
- Connection health

**🔍 Checking Serena Status...**

```bash
echo "🧪 MCF Serena Integration Status Check"
echo "======================================"
echo ""

# 1. Check Claude Code MCP configuration
echo "1️⃣ Claude Code MCP Configuration:"
if claude mcp list 2>/dev/null | grep -q "serena"; then
    echo "   ✅ Serena MCP server is configured"
else
    echo "   ❌ Serena MCP server not found"
    echo "   💡 Run '/serena:install' to add Serena MCP server"
fi
echo ""

# 2. Check Serena accessibility
echo "2️⃣ Serena Accessibility:"
if uvx --from git+https://github.com/oraios/serena serena --help > /dev/null 2>&1; then
    echo "   ✅ Serena is accessible via uvx"
else
    echo "   ❌ Cannot access Serena"
    echo "   💡 Check if uv is installed: curl -LsSf https://astral.sh/uv/install.sh | sh"
fi
echo ""

# 3. Check Serena configuration
echo "3️⃣ Serena Configuration:"
if [ -f ~/.serena/serena_config.yml ]; then
    echo "   ✅ Configuration file exists: ~/.serena/serena_config.yml"
    
    # Check key settings
    if grep -q "record_tool_usage_stats: true" ~/.serena/serena_config.yml; then
        echo "   ✅ Usage analytics enabled"
    else
        echo "   ⚠️  Usage analytics disabled"
    fi
    
    if grep -q "web_dashboard: true" ~/.serena/serena_config.yml; then
        echo "   ✅ Web dashboard enabled"
    else
        echo "   ⚠️  Web dashboard disabled"
    fi
    
    if grep -q "default_max_tool_answer_chars: 200000" ~/.serena/serena_config.yml; then
        echo "   ✅ MCF optimized response size"
    else
        echo "   ⚠️  Default response size (may need optimization)"
    fi
else
    echo "   ❌ Configuration file missing"
    echo "   💡 Run '/serena:config update' to create optimized config"
fi
echo ""

# 4. Check current project initialization
echo "4️⃣ Current Project Status:"
if [ -f .serena/project.yml ]; then
    echo "   ✅ Project initialized for Serena"
    if grep -q "MCF" .serena/project.yml; then
        echo "   ✅ MCF-specific configuration detected"
    else
        echo "   ⚠️  Basic configuration (consider running /serena:init)"
    fi
    
    # Check for memories
    if [ -d .serena/memories ] && [ "$(ls -A .serena/memories)" ]; then
        MEMORY_COUNT=$(ls .serena/memories | wc -l)
        echo "   ✅ $MEMORY_COUNT memory files found"
    else
        echo "   ⚠️  No memory files found"
    fi
else
    echo "   ❌ Project not initialized"
    echo "   💡 Run '/serena:init' to initialize current project"
fi
echo ""

# 5. Check for code compatibility
echo "5️⃣ Code Compatibility:"
CODE_TYPES=0
for ext in py md sh yml yaml json js ts; do
    if find . -name "*.$ext" -not -path "./.git/*" -not -path "./node_modules/*" | head -1 | read; then
        CODE_TYPES=$((CODE_TYPES + 1))
    fi
done

if [ $CODE_TYPES -gt 0 ]; then
    echo "   ✅ Found $CODE_TYPES supported file types for semantic analysis"
else
    echo "   ⚠️  No supported code files found"
fi
echo ""

# 6. Connection test (if possible)
echo "6️⃣ Connection Test:"
echo "   💡 To test live connection:"
echo "   • Use '/mcp' command to see active servers"
echo "   • Try '/serena:overview' for semantic analysis" 
echo "   • Monitor at: http://localhost:24282/dashboard/"
echo ""

echo "📋 Quick Actions:"
echo "• Install: /serena:install"
echo "• Initialize: /serena:init" 
echo "• Configure: /serena:config"
echo "• Test: /serena:overview"
echo ""
echo "🚀 MCF + Serena = Semantic Code Superpowers!"
```

**📊 What This Status Check Reveals:**

1. **MCP Integration** - Whether Serena is connected to Claude Code
2. **Configuration Health** - If settings are optimized for MCF workflow  
3. **Project Readiness** - Whether current project is set up for semantic analysis
4. **Code Compatibility** - If your files can benefit from Serena's tools
5. **Quick Diagnostics** - Common issues and their solutions

**🛠️ Common Issues & Solutions:**

- **Serena not found in MCP list** → Run `/serena:install`
- **Configuration missing** → Run `/serena:config update`
- **Project not initialized** → Run `/serena:init`
- **Connection issues** → Restart Claude Code and check `/mcp`

This gives you a complete health check of your MCF + Serena integration! 🎯
