---
allowed-tools: Bash(git*), Bash(gh*)
argument-hint: [natural language instruction]
description: Perform Git/GitHub operations based on natural language
---

Perform Git operation: "$ARGUMENTS"

Repository context:
!`git status --porcelain`
!`git branch --show-current`
!`git log --oneline -5`
!`git remote -v`

Analyze the instruction "$ARGUMENTS" and perform the appropriate Git operation.

Examples I can handle:
- "create a new branch for login feature" → `git checkout -b feature/login`
- "undo my last commit" → `git reset HEAD~1`
- "merge main into current branch" → `git pull origin main`
- "create pull request" → `gh pr create`
- "switch to main branch" → `git checkout main`
- "delete branch feature-x" → `git branch -d feature-x`
- "tag this as v1.0.0" → `git tag v1.0.0 && git push origin v1.0.0`
- "stash my changes" → `git stash push`
- "apply my stash" → `git stash pop`

Execute the appropriate Git commands safely based on the instruction.
