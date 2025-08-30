---
description: Create intelligent bookmarks of current conversation state for later restoration
argument-hint: <bookmark-name>
---

## Context Bookmark Creation

**Bookmark Name:** $ARGUMENTS

Create an intelligent bookmark capturing the current state of our conversation:

### 1. **Session State Summary**
- Current project objectives and active goals
- Key decisions made in this session
- Major milestones and achievements reached
- Current development focus and priorities

### 2. **Technical State Capture**
- Files modified or created in this session
- Configuration changes and settings updates
- New commands, hooks, or agents created
- MCP server configurations and integrations
- Database schemas or API changes

### 3. **Code Context**
- Important code implementations and patterns
- Architecture decisions and design choices
- Security considerations and implementations
- Performance optimizations applied
- Testing strategies and coverage

### 4. **Project Knowledge**
- Lessons learned and best practices identified
- Problems solved and solutions implemented
- Patterns and approaches that worked well
- Mistakes to avoid and pitfalls discovered

### 5. **Development Environment**
- Tool configurations and customizations
- Dependencies and version information
- Environment variables and secrets setup
- Build and deployment configurations

### 6. **Future Planning**
- Next steps and planned improvements
- Technical debt identified but not addressed
- Future feature requirements and considerations
- Refactoring opportunities and architectural improvements

### 7. **Context Restoration Guide**
Create a restoration guide including:
- Key files to review when restoring this context
- Commands to re-run to verify state
- Configurations to check for consistency
- Tests to run to validate current functionality

### 8. **Bookmark Metadata**
- Creation timestamp and session info
- Claude Code version and model used
- Project state and git branch information
- Token usage estimates at bookmark time

**Save this comprehensive bookmark to `.claude/bookmarks/$ARGUMENTS.md` for future reference.**

Execute bookmark creation now.