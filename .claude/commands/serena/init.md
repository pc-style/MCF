---
description: Initialize current project for Serena semantic analysis
---

Initialize the current project for Serena semantic code analysis and create project-specific configuration.

This command will:
1. Generate project configuration using Serena's built-in generator
2. Create MCF-specific memories and settings
3. Set up semantic analysis for your codebase

**Project Initialization:**

Let me initialize this project with Serena:

```bash
# Generate project configuration using Serena's built-in generator
uvx --from git+https://github.com/oraios/serena serena project generate-yml
```

Now I'll enhance the generated configuration for MCF workflow:

```bash
# Create .serena directory if it doesn't exist
mkdir -p .serena/memories

# Enhance project.yml with MCF-specific settings
cat >> .serena/project.yml << 'EOF'

# MCF-Specific Enhancements
description: "My Claude Flow - Advanced Claude Code development automation with semantic analysis"

# Language-specific optimizations for MCF
languages:
  python:
    enabled: true
    include_patterns:
      - ".claude/hooks/*.py"
      - "workflow/*.py"
  
  markdown:
    enabled: true
    include_patterns:
      - "**/*.md"
      - ".claude/commands/**/*.md"
      - ".claude/agents/**/*.md"
  
  shell:
    enabled: true
    include_patterns:
      - "**/*.sh"
      - ".claude/hooks/*.sh"
  
  yaml:
    enabled: true
    include_patterns:
      - "**/*.yml"
      - "**/*.yaml"
  
  json:
    enabled: true
    include_patterns:
      - "**/*.json"

# MCF-specific ignore patterns
ignore_patterns:
  - ".git"
  - "node_modules"
  - "__pycache__"
  - "*.log"
  - ".claude/settings.local.json"
  - "cc-docs/not-relevant"

# Enable comprehensive analysis for MCF
analysis:
  deep_analysis: true
  include_docs: true
  cross_reference_analysis: true
EOF

# Create MCF architecture memory
cat > .serena/memories/mcf_architecture.md << 'EOF'
# MCF Architecture Memory

## Key Components

### Agent System (.claude/agents/)
- 8+ specialized agents with semantic analysis capabilities
- micro-analyzer, security-auditor, perf-optimizer, api-designer
- Each agent uses Serena tools for efficient code operations

### Hook System (.claude/hooks/)
- Event-driven automation with Python/shell scripts
- Context monitoring, command suggestions, auto-formatting
- Integration with Serena for semantic awareness

### Command System (.claude/commands/)
- 40+ custom slash commands organized by functionality
- gh/ (git operations), project/ (management), context/ (state)
- Serena-specific commands for semantic operations

### Template System (workflow/)
- Dynamic project scaffolding with template engine
- Built-in templates for common project types
- Integration with semantic code understanding

## Semantic Development Patterns

### Symbol-First Approach
- Use find_symbol before reading files
- Leverage get_symbol_info for code understanding
- Apply find_referencing_symbols for impact analysis

### Agent Enhancement
- Agents use semantic tools for precise operations
- Token efficiency through symbol-level work
- Context-aware automation and suggestions
EOF

# Create workflow patterns memory
cat > .serena/memories/mcf_workflow_patterns.md << 'EOF'
# MCF Workflow Patterns

## Serena Integration Patterns

### Code Analysis Flow
1. /serena:overview - Get project structure
2. /serena:find - Locate specific symbols
3. /serena:analyze - Get detailed symbol information
4. /serena:refs - Understand usage patterns

### Agent-Enhanced Development
- micro-analyzer: Symbol-level code analysis
- security-auditor: Data flow tracing with find_referencing_symbols
- perf-optimizer: Bottleneck identification at function level
- api-designer: Endpoint-to-implementation mapping

### MCF + Serena Best Practices
1. Start with semantic tools before file operations
2. Use specialized agents with semantic capabilities
3. Leverage symbol-level operations for token efficiency
4. Combine MCF commands with Serena understanding
EOF
```

✅ **Project initialized for Serena semantic analysis!**

Now I'll activate this project with Serena:

```bash
# The project should now be ready for semantic analysis
echo "Project initialized with Serena semantic capabilities"
echo "Location: $(pwd)"
echo "Configuration: .serena/project.yml"
echo "Memories: .serena/memories/"
```

**Project Activation:**
To activate this project in the current session, ask me to:
"Activate the current project for Serena semantic analysis"

**What's Available Now:**
- 🔍 **find_symbol** - Locate functions, classes, variables by name
- 📊 **get_symbol_info** - Detailed symbol information and documentation  
- 🔗 **find_referencing_symbols** - Find all references to a symbol
- 🏗️ **get_project_structure** - Comprehensive project overview
- ✏️ **Symbol editing tools** - Precise code insertion and modification

**Next Steps:**
1. Test with `/serena:overview` to see project structure
2. Use `/serena:find <symbol>` to locate code elements
3. Monitor usage at http://localhost:24282/dashboard/

Your MCF project now has full semantic code understanding! 🚀
