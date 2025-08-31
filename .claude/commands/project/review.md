---
allowed-tools: Bash, Read, Grep, Glob, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
argument-hint: [commit-hash or file-pattern]
description: Automated code review with best practices check
---

!`git diff HEAD~1..HEAD --name-only`

Perform comprehensive code review:

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.


Target: $ARGUMENTS (if provided, otherwise recent changes)

Focus areas:
- Code quality and readability
- Security vulnerabilities
- Performance implications
- Best practices adherence
- Test coverage
- Documentation completeness

Provide actionable feedback with examples.
