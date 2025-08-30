---
description: Intelligently purge conversation context while preserving essential information
argument-hint: [focus-areas]
---

## Context Purge Request

**Focus Areas:** $ARGUMENTS

Please perform an intelligent context purge by compacting this conversation while preserving all essential information:

### 1. **Preserve Critical Context**
- Current project objectives and active tasks
- Recent key decisions and their reasoning
- Important file changes and modifications  
- Configuration settings and custom commands
- MCP server configurations and authentication
- Subagent definitions and custom hooks
- Active debugging or development state

### 2. **Compress Verbose Content**
- Lengthy explanations of basic concepts
- Repetitive code examples and patterns
- Multiple iterations of similar discussions
- Redundant confirmations and acknowledgments
- Excessive debugging output and logs
- Verbose installation or setup instructions

### 3. **Remove Low-Value Content**
- Outdated information superseded by newer decisions
- Exploratory discussions that didn't lead to action
- Tangential conversations not related to current goals
- Error messages that have been resolved
- Duplicate information across multiple messages

### 4. **Maintain Context Flow**
- Keep logical progression of development work
- Preserve cause-and-effect relationships between decisions
- Maintain references between related code changes
- Keep context necessary for understanding current state

### 5. **Special Preservation**
If focus areas specified, give extra attention to preserving context related to: $ARGUMENTS

**Goal:** Create a clean, focused conversation that maintains essential project context while reducing token usage by 50-70%.

Execute the purge now, providing a brief summary of what was preserved vs. compressed.