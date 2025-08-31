---
allowed-tools: Bash, Read, Edit, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
argument-hint: [package-manager]
description: Smart dependency updates with safety checks
---

Update project dependencies safely:

1. Detect package manager (npm, yarn, pip, cargo, etc.) or use $1

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.

2. List outdated dependencies
3. Categorize updates by risk level (patch/minor/major)
4. Suggest update strategy
5. Create backup/branch if requested
6. Perform updates with testing
