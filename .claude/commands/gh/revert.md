---
allowed-tools: Bash(git*)
argument-hint: [commit-hash]
description: Show recent commits and revert the specified one
---

Revert commit: "$ARGUMENTS"

Recent commits:
!`git log --oneline -10`

Please:
1. If no commit hash provided in "$ARGUMENTS", show the above commits and ask which one to revert
2. Show details of the commit to be reverted
3. Use `git revert` to safely revert the commit
4. Handle any merge conflicts if they occur
