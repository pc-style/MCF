---
allowed-tools: Bash(git*), Bash(gh*), mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
description: Add all changes, commit with smart message, and push
---

Quick workflow to add all changes, commit, and push to current branch.

Current status:

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.

!`git status --porcelain`
!`git branch --show-current`

Please:

1. Add all changes with `git add .`
2. Create a meaningful commit message based on the changes shown above
3. Commit the changes
4. Push to the current branch (create upstream if needed)
