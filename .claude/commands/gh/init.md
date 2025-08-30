---
allowed-tools: Bash(git*), Bash(gh*), Bash(pwd*), Bash(ls*)
description: Initialize git repo and create private GitHub repository
---

Initialize a Git repository and create a GitHub repo.

Current directory status:
!`pwd`
!`ls -la`

Check if git repo exists:
!`git status 2>/dev/null || echo "Not a git repository"`

Please:
1. Check if this is already a git repository
2. If not, initialize with `git init`
3. Create .gitignore and README.md if they don't exist
4. Make initial commit
5. Create private GitHub repository using `gh repo create --private`
6. Add remote and push initial commit
