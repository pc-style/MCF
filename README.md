# MCF - Multi Component Framework

A sophisticated development automation platform built around Claude Code integration, featuring intelligent hooks, custom commands, comprehensive documentation, and **semantic code analysis via Serena integration**.

## 📁 Project Structure

```
MCF/
├── .claude/                          # Claude Code configuration hub
│   ├── agents/                       # Specialized AI agents (9 total)
│   │   ├── api-designer.md           # RESTful/GraphQL API design specialist
│   │   ├── devops-engineer.md        # Infrastructure & deployment specialist
│   │   ├── docs-researcher.md        # Documentation research specialist
│   │   ├── micro-analyzer.md         # Lightweight code analysis specialist
│   │   ├── micro-executor.md         # Focused task execution specialist
│   │   ├── micro-researcher.md       # Quick documentation research
│   │   ├── perf-optimizer.md         # Performance analysis specialist
│   │   ├── security-auditor.md       # Security vulnerability specialist
│   │   └── semantic-navigator.md     # Serena semantic code navigation specialist
│   │
│   ├── commands/                     # Custom slash commands
│   │   ├── context/                  # Context management commands
│   │   │   ├── agent.md              # Launch specialized agents
│   │   │   ├── bookmark.md           # Save conversation bookmarks
│   │   │   ├── help.md               # Context-specific help
│   │   │   ├── load.md               # Load saved contexts
│   │   │   ├── merge.md              # Merge conversation contexts
│   │   │   ├── purge.md              # Clean up context data
│   │   │   ├── split.md              # Split conversations
│   │   │   └── state.md              # View context state
│   │   │
│   │   ├── gh/                       # Git/GitHub workflow commands
│   │   │   ├── auto.md               # Natural language git operations
│   │   │   ├── commit.md             # Quick add + commit workflow
│   │   │   ├── init.md               # Initialize git + GitHub repo
│   │   │   ├── push.md               # Add + commit + push workflow
│   │   │   ├── push-beta.md          # Beta push workflow
│   │   │   └── revert.md             # Safe commit reverting
│   │   │
│   │   ├── project/                  # Project management commands
│   │   │   ├── analyze.md            # Comprehensive project analysis
│   │   │   ├── debug.md              # Debug project issues
│   │   │   ├── deploy.md             # Deployment workflows
│   │   │   ├── deps-update.md        # Update project dependencies
│   │   │   ├── docs.md               # Generate documentation
│   │   │   ├── review.md             # Code review workflows
│   │   │   └── security.md           # Security analysis
│   │   │
│   │   ├── serena/                   # ⭐ NEW: Serena semantic analysis commands
│   │   │   ├── activate.md           # Activate project for semantic analysis
│   │   │   ├── analyze.md            # Deep symbol analysis
│   │   │   ├── config.md             # Serena configuration management
│   │   │   ├── find.md               # Symbol search and discovery
│   │   │   ├── help.md               # Complete Serena command guide
│   │   │   ├── init.md               # Initialize project for Serena
│   │   │   ├── install.md            # Install and setup Serena MCP
│   │   │   ├── overview.md           # Project structure analysis
│   │   │   ├── refs.md               # Reference tracking
│   │   │   └── status.md             # Integration health check
│   │   │
│   │   └── templates/                # Template system commands
│   │       ├── add.md                # Add new templates
│   │       ├── init.md               # Initialize from templates
│   │       ├── list.md               # List available templates
│   │       └── save.md               # Save current project as template
│   │
│   ├── hooks/                        # Event-driven automation system
│   │   ├── auto_format.py            # Automatic code formatting
│   │   ├── bash_safety.py            # Bash command validation
│   │   ├── claude_command_suggestions.py  # Contextual command suggestions
│   │   ├── context7_reminder.py      # Context7 MCP integration reminders
│   │   ├── context_monitor.py        # Conversation context monitoring
│   │   ├── enhanced_statusline.sh    # Enhanced status line display
│   │   ├── git_suggestions.py        # Git workflow suggestions (legacy)
│   │   ├── git_work_tracker.py       # File modification tracking
│   │   ├── serena_context_suggestions.py  # ⭐ NEW: Serena usage suggestions
│   │   ├── serena_project_init.py    # ⭐ NEW: Serena project initialization
│   │   └── session_setup.py          # Session initialization (Serena-aware)
│   │
│   ├── settings.json                 # Core Claude Code configuration
│   └── settings.local.json           # Local project settings
│
├── .serena/                          # ⭐ NEW: Serena semantic analysis config
│   ├── memories/                     # Project-specific memories
│   │   ├── mcf_architecture.md       # MCF system architecture
│   │   └── mcf_workflow_patterns.md  # Common workflow patterns
│   └── project.yml                   # Serena project configuration
│
├── .gitignore                        # Git ignore rules
└── README.md                         # This file
```

## 🚀 Features

### **🧠 Semantic Code Analysis (NEW!)**

- **Serena Integration**: IDE-like semantic code understanding and navigation
- **Symbol-Level Operations**: Work with functions, classes, and variables directly
- **Token Efficiency**: Massive token savings through precise symbol targeting
- **Cross-Reference Analysis**: Trace data flow and dependencies throughout codebase
- **10 Serena Commands**: Complete semantic analysis toolkit (`/serena:*`)

### **🤖 AI Agent System**

- **9 Specialized Agents**: Each optimized for specific development tasks (including new semantic-navigator)
- **Semantic Enhancement**: All agents upgraded with Serena semantic capabilities
- **Micro Agents**: Lightweight, focused execution with minimal context usage
- **Domain Experts**: API design, DevOps, security, performance optimization
- **Auto-Discovery**: Agents are automatically suggested based on task context

### **⚡ Custom Command System**

- **50+ Custom Commands**: Organized by functionality (gh/, project/, context/, templates/, serena/)
- **Natural Language Git**: `/gh:auto "create feature branch"` translates to appropriate git commands
- **Workflow Automation**: `/gh:push` does add + commit + push in one command
- **Template Management**: Complete project scaffolding system
- **Semantic Operations**: `/serena:find`, `/serena:analyze`, `/serena:refs` for code navigation

### **🔧 Intelligent Hook System**

- **Event-Driven**: Hooks respond to file changes, user input, and session events
- **Smart Suggestions**: Contextual command recommendations based on repository state
- **Serena Awareness**: Hooks suggest semantic operations when appropriate
- **Safety Mechanisms**: Bash command validation and input sanitization
- **Auto-Formatting**: Automatic code formatting on file saves

### **📚 Comprehensive Documentation**

- **Complete Claude Code Reference**: All features and configurations documented
- **Serena Integration Guide**: Setup, usage, and best practices
- **Best Practices**: Security, performance, and workflow guidelines
- **Integration Guides**: MCP servers, IDEs, and third-party tools
- **Troubleshooting**: Common issues and solutions

## 🎯 Key Components

### **Serena Semantic Engine** (`.claude/commands/serena/` + `.serena/`)

- **Symbol Discovery**: `find_symbol` for locating functions, classes, variables
- **Deep Analysis**: `get_symbol_info` for detailed code understanding
- **Reference Tracking**: `find_referencing_symbols` for usage analysis
- **Project Structure**: `get_project_structure` for architectural insights
- **Precise Editing**: Symbol-level code insertion and modification

### **Template Engine** (`workflow/template-engine.py`)

- **Project Scaffolding**: Initialize new projects from templates
- **Variable Substitution**: Dynamic template customization
- **Post-Install Scripts**: Automated setup after template application
- **Built-in Templates**: Common project types (React, Python, etc.)

### **Hook System** (`.claude/hooks/`)

- **Command Suggestions**: Contextual `/gh:*`, `/serena:*` and workflow recommendations
- **Safety Validation**: Prevents dangerous bash operations
- **Auto-Formatting**: Code formatting on file operations
- **Context Monitoring**: Tracks conversation context and suggests optimizations
- **Serena Integration**: Smart suggestions for semantic operations

### **Slash Commands** (`.claude/commands/`)

- **Git Workflows**: Natural language git operations and quick workflows
- **Project Management**: Analysis, debugging, deployment, and reviews
- **Context Management**: Save, load, merge, and split conversation contexts
- **Template Operations**: Create, list, and apply project templates
- **Semantic Analysis**: Complete suite of code understanding tools

## 🛡️ Security Features

- **Input Validation**: All user inputs are sanitized and validated
- **Command Whitelisting**: Only approved bash commands are allowed
- **Path Traversal Protection**: Prevents unauthorized file access
- **Hook Safety**: All hooks have timeout limits and error handling
- **Audit Logging**: Track all automated operations
- **Semantic Security**: Code vulnerability detection via semantic analysis

## 🚀 Quick Start

### **Standard Setup**

1. **Clone the repository**:

   ```bash
   git clone <your-repo-url>
   cd MCF
   ```

2. **Initialize Claude Code**:

   ```bash
   claude --project .
   ```

### **⭐ NEW: Serena Semantic Setup**

3. **Install Serena integration**:

   ```
   /serena:install
   ```

4. **Initialize project for semantic analysis**:

   ```
   /serena:init
   ```

5. **Verify everything works**:

   ```
   /serena:status
   ```

### **Start Using**

6. **Try semantic code analysis**:

   ```
   /serena:overview                    # See project structure
   /serena:find MyFunction             # Find specific symbols
   ```

7. **Try workflow commands**:

   ```
   /gh:push                           # Git workflow
   ```

8. **Get contextual suggestions**:
   - The hooks will automatically suggest relevant commands based on your work

9. **Explore available commands**:
   ```
   /help
   ```

## 📖 Usage Examples

### **⭐ NEW: Semantic Code Analysis**

```bash
/serena:overview                     # Get high-level project structure
/serena:find UserService             # Find UserService class/function
/serena:analyze authenticate         # Deep analysis of authenticate method
/serena:refs getCurrentUser          # Find all references to getCurrentUser
/serena:config                       # Check Serena configuration
```

### **Git Workflows**

```bash
/gh:commit "Add new feature"          # Quick add + commit
/gh:push                              # Add + commit + push with smart message
/gh:auto "create feature branch"      # Natural language git operations
/gh:revert                           # Safe commit reverting
```

### **Project Management**

```bash
/project:analyze                     # Comprehensive project analysis (now with Serena!)
/project:review                      # Code review workflow (semantic-enhanced)
/project:security                    # Security analysis (with data flow tracing)
```

### **Template System**

```bash
/templates:list                      # Show available templates
/templates:init react-app            # Initialize from template
/templates:save my-template          # Save current project as template
```

## 🔧 Configuration

The system is configured through:

- **`.claude/settings.json`**: Core Claude Code configuration with hooks and status line
- **`.claude/settings.local.json`**: Local project-specific settings
- **`~/.serena/serena_config.yml`**: Global Serena semantic analysis configuration
- **`.serena/project.yml`**: Project-specific Serena settings
- **Individual command files**: Each slash command is customizable
- **Hook scripts**: Event-driven automation can be modified or extended

## 🌟 What Makes MCF Special

### **Semantic Superpowers**

- **10x Token Efficiency**: Work at symbol-level instead of reading entire files
- **IDE-like Navigation**: Find, analyze, and modify code with surgical precision  
- **Agent Enhancement**: All 9 agents get semantic code understanding
- **Smart Suggestions**: Context-aware recommendations for semantic operations

### **Complete Portability**

- **One-Folder Setup**: Just copy `.claude/` folder to any new PC
- **Self-Contained**: All intelligence, commands, hooks, and agents included
- **Auto-Configuration**: Serena and other integrations set up automatically
- **Cross-Platform**: Works on any system with Claude Code and uv

### **Professional Automation**

- **50+ Commands**: Comprehensive workflow automation
- **Event-Driven Hooks**: Intelligent responses to development events
- **Template System**: Rapid project scaffolding and standardization
- **Safety First**: Input validation, command whitelisting, audit logging

## 🤝 Contributing

1. **Add new commands**: Create `.md` files in `.claude/commands/`
2. **Create hooks**: Add Python/shell scripts to `.claude/hooks/`
3. **Extend agents**: Modify agent definitions in `.claude/agents/`
4. **Add templates**: Update `workflow/builtin-templates.json`
5. **Enhance Serena integration**: Extend `.serena/memories/` or create new semantic workflows

## 📊 Performance Benefits

With Serena integration, MCF provides:

- **Massive Token Savings**: Symbol-level operations vs full file reading
- **Faster Development**: IDE-like code navigation and understanding
- **Better Code Quality**: Semantic analysis catches issues regular text analysis misses
- **Enhanced Debugging**: Trace data flow and dependencies precisely
- **Smarter Refactoring**: Understand impact before making changes

## 📄 License

This project contains configuration and documentation for Claude Code, an AI development assistant by Anthropic.

---

**Built with Claude Code + Serena** - The ultimate development productivity platform combining AI automation with semantic code intelligence. 🚀