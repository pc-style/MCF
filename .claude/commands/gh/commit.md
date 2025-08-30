---
allowed-tools: Bash(git*)
argument-hint: [commit-message]
description: Add all changes and commit with provided message
---

Add all changes and commit with message: "$ARGUMENTS"

Current changes:
!`git status --porcelain`
!`git diff --stat`

Please:
1. Add all changes with `git add .`
2. Commit with the message: "$ARGUMENTS" (or generate a good message if none provided)
