---
allowed-tools: Bash, Read, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
argument-hint: [environment]
description: Smart deployment with pre-flight checks
---

Deploy to $1 environment (default: staging):

1. Run pre-deployment checks:

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.

   - Tests passing
   - No uncommitted changes
   - Dependencies up to date
   - Security scan
2. Build optimization
3. Deploy with rollback plan
4. Post-deployment verification
5. Notify relevant stakeholders
