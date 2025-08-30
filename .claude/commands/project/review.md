---
allowed-tools: Bash, Read, Grep, Glob
argument-hint: [commit-hash or file-pattern]
description: Automated code review with best practices check
---

!`git diff HEAD~1..HEAD --name-only`

Perform comprehensive code review:

Target: $ARGUMENTS (if provided, otherwise recent changes)

Focus areas:
- Code quality and readability
- Security vulnerabilities
- Performance implications
- Best practices adherence
- Test coverage
- Documentation completeness

Provide actionable feedback with examples.
