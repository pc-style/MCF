---
description: Activate current project for Serena semantic analysis
---

Activate the current project with Serena for semantic code analysis.

This command will:
1. Ensure the project is properly initialized
2. Activate semantic analysis for the current directory
3. Trigger Serena's onboarding process if needed

**Project Activation:**

```bash
echo "🚀 Activating MCF project for Serena semantic analysis..."
echo "📍 Project Location: $(pwd)"
echo ""

# Check if project is initialized
if [ ! -f .serena/project.yml ]; then
    echo "⚠️  Project not initialized. Initializing now..."
    # Initialize first
    uvx --from git+https://github.com/oraios/serena serena project generate-yml
    echo "✅ Basic project configuration created"
else
    echo "✅ Project configuration found"
fi

echo ""
echo "🔍 Project Details:"
if [ -f .serena/project.yml ]; then
    echo "   Configuration: .serena/project.yml"
    if grep -q "name:" .serena/project.yml; then
        PROJECT_NAME=$(grep "name:" .serena/project.yml | head -1 | cut -d':' -f2 | xargs)
        echo "   Name: $PROJECT_NAME"
    fi
fi

if [ -d .serena/memories ]; then
    MEMORY_COUNT=$(ls .serena/memories | wc -l 2>/dev/null || echo 0)
    echo "   Memories: $MEMORY_COUNT files"
fi

echo ""
echo "📋 To complete activation, ask me:"
echo '   "Activate the current project for Serena semantic analysis"'
echo ""
echo "Or more specifically:"
echo '   "Use Serena to activate this project for semantic code understanding"'
```

**What Happens During Activation:**

1. **Project Recognition** - Serena scans your codebase structure
2. **Language Detection** - Identifies programming languages and frameworks
3. **Semantic Indexing** - Creates symbol tables and cross-references  
4. **Memory Loading** - Reads project-specific memories and context
5. **Tool Preparation** - Makes semantic tools available for this project

**After Activation, You Can Use:**

- 🔍 **find_symbol** - Locate any function, class, or variable by name
- 📊 **get_symbol_info** - Get detailed information about code symbols
- 🔗 **find_referencing_symbols** - Find all references to a symbol
- 🏗️ **get_project_structure** - Get comprehensive project overview
- ✏️ **Symbol editing** - Precise code insertion and modification

**Verify Activation:**

```bash
# Check if project appears in Serena's project list
echo "After activation, your project should appear when you ask Serena to list projects"
echo "You can also test with semantic commands like:"
echo "• /serena:overview"
echo "• /serena:find <symbol_name>"
```

**Pro Tips:**

- **First Time Setup**: Serena will perform an onboarding process to understand your project
- **Memory Creation**: Let Serena create memories about your codebase architecture
- **Context Switch**: Consider starting a new conversation after onboarding to free up context
- **Dashboard Monitoring**: Watch the process at http://localhost:24282/dashboard/

Your MCF project is ready for semantic superpowers! 🎯
