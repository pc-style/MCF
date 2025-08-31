---
allowed-tools: Bash(git*), mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
description: Create a commit with proper message formatting
---

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.

!`git status`
!`git diff --cached`

Create a commit with proper message formatting.

1. **Review staged changes** to understand what will be committed
2. **Write a clear commit message** that explains the "why" not just the "what"
3. **Follow conventional commit format** (feat:, fix:, docs:, etc.)
4. **Keep first line under 50 characters**
5. **Add detailed description if needed** after a blank line

Commit message format:

```
type(scope): brief description

Detailed explanation of what and why, not how.
Include any breaking changes or special notes.
```

Example commit messages:

- `feat(auth): add OAuth2 integration for Google login`
- `fix(api): resolve race condition in user registration`
- `docs(readme): update installation instructions for Windows`
- `refactor(utils): extract validation logic to separate module`

Create the commit now.
