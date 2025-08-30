---
allowed-tools: Bash(git*)
argument-hint: [branch-name]
description: Create feature branch, commit changes, and push
---

Create feature branch and push changes. Branch name: "$ARGUMENTS"

Current status:
!`git status --porcelain`
!`git branch --show-current`

Please:
1. Generate a branch name like "feature/beta-YYYYMMDD" if "$ARGUMENTS" is empty
2. Create and checkout the new branch
3. Add all changes with `git add .`
4. Create a descriptive commit message
5. Commit the changes
6. Push the new branch to remote
7. Suggest next steps (like creating a PR)
