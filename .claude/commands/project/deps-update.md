---
allowed-tools: Bash, Read, Edit
argument-hint: [package-manager]
description: Smart dependency updates with safety checks
---

Update project dependencies safely:

1. Detect package manager (npm, yarn, pip, cargo, etc.) or use $1
2. List outdated dependencies
3. Categorize updates by risk level (patch/minor/major)
4. Suggest update strategy
5. Create backup/branch if requested
6. Perform updates with testing
