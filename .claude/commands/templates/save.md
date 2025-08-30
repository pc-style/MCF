---
allowed-tools: Bash, Read, Write, Edit, MultiEdit, Glob
argument-hint: <template-name>
description: Save the current project setup as a reusable template
---

# Save Current Project as Template

You are saving the current project structure as a template named: **$1**

## Process:

1. **Analyze current directory structure** to understand the project layout
2. **Identify key configuration files** (package.json, requirements.txt, etc.)
3. **Detect the project type** and main technologies used
4. **Ask user for template details**:
   - Description of the template
   - Variables that should be configurable (project name, database type, etc.)
   - Which files should be templated vs copied exactly
5. **Generate template JSON** with setup steps
6. **Save to** `~/mcf/templates/$1.json`

## Template Generation Logic:

- **Detect package managers**: npm, yarn, pnpm, uv, pip, cargo, etc.
- **Identify frameworks**: React, Vue, Django, FastAPI, etc.
- **Find config files**: .env.example, docker-compose.yml, etc.
- **Capture initialization commands** based on detected stack
- **Create parameterized steps** for project creation

## Variables to Consider:

Common variables that should be configurable:
- `project_name`: The name of the new project
- `database_type`: PostgreSQL, MySQL, SQLite, etc.
- `auth_provider`: Auth0, Supabase, custom, etc.
- `css_framework`: Tailwind, Bootstrap, etc.
- `deployment_target`: Vercel, Netlify, AWS, etc.

## Implementation Steps:

### 1. Analyze Current Project
```bash
echo "🔍 Analyzing current project structure..."

# Detect package managers and frameworks
if [ -f "package.json" ]; then
  echo "📦 Detected Node.js project"
  PACKAGE_MANAGER="npm"
  if [ -f "yarn.lock" ]; then PACKAGE_MANAGER="yarn"; fi
  if [ -f "pnpm-lock.yaml" ]; then PACKAGE_MANAGER="pnpm"; fi

  FRAMEWORK=$(jq -r '.dependencies | keys[]' package.json | grep -E "(react|vue|angular|next|svelte)" | head -1 || echo "unknown")
elif [ -f "pyproject.toml" ] || [ -f "requirements.txt" ] || [ -f "setup.py" ]; then
  echo "🐍 Detected Python project"
  PACKAGE_MANAGER="pip"
  if [ -f "uv.lock" ]; then PACKAGE_MANAGER="uv"; fi
  if [ -f "poetry.lock" ]; then PACKAGE_MANAGER="poetry"; fi

  FRAMEWORK=$(grep -r "fastapi\|flask\|django" . --include="*.py" | head -1 | cut -d: -f2 || echo "unknown")
elif [ -f "Cargo.toml" ]; then
  echo "🦀 Detected Rust project"
  PACKAGE_MANAGER="cargo"
  FRAMEWORK="rust"
else
  echo "❓ Unknown project type"
  FRAMEWORK="unknown"
fi

# Detect other configuration files
CONFIG_FILES=""
if [ -f ".env.example" ]; then CONFIG_FILES="$CONFIG_FILES .env.example"; fi
if [ -f "docker-compose.yml" ]; then CONFIG_FILES="$CONFIG_FILES docker-compose.yml"; fi
if [ -f "Dockerfile" ]; then CONFIG_FILES="$CONFIG_FILES Dockerfile"; fi
```

### 2. Collect Template Information
```bash
echo ""
echo "📝 Template Configuration for '$1'"
echo ""

# Get template description
read -p "Description of this template: " DESCRIPTION

# Get category
echo "Category options: frontend, backend, fullstack, mobile, desktop"
read -p "Template category: " CATEGORY

# Detect prerequisites
PREREQUISITES="[\"git\"]"
if [ "$PACKAGE_MANAGER" = "npm" ] || [ "$PACKAGE_MANAGER" = "yarn" ] || [ "$PACKAGE_MANAGER" = "pnpm" ]; then
  PREREQUISITES="$PREREQUISITES,\"node\""
elif [ "$PACKAGE_MANAGER" = "pip" ] || [ "$PACKAGE_MANAGER" = "uv" ] || [ "$PACKAGE_MANAGER" = "poetry" ]; then
  PREREQUISITES="$PREREQUISITES,\"python\""
elif [ "$PACKAGE_MANAGER" = "cargo" ]; then
  PREREQUISITES="$PREREQUISITES,\"rust\""
fi
```

### 3. Define Variables
```bash
echo ""
echo "🔧 Configurable Variables"

# Always include project_name
VARIABLES="[
  {
    \"name\": \"project_name\",
    \"prompt\": \"What is your project name?\",
    \"default\": \"my-$1-app\",
    \"validation\": \"^[a-zA-Z0-9-_]+$\"
  }
]"

# Ask about additional variables
read -p "Add custom variables? (y/N): " ADD_VARS
if [[ $ADD_VARS =~ ^[Yy]$ ]]; then
  echo "Enter variables as: name,prompt,default (empty line to finish)"
  while true; do
    read -p "Variable: " VAR_INPUT
    if [ -z "$VAR_INPUT" ]; then break; fi

    NAME=$(echo "$VAR_INPUT" | cut -d, -f1)
    PROMPT=$(echo "$VAR_INPUT" | cut -d, -f2)
    DEFAULT=$(echo "$VAR_INPUT" | cut -d, -f3)

    VARIABLES=$(echo "$VARIABLES" | jq ". + [{\"name\":\"$NAME\",\"prompt\":\"$PROMPT\",\"default\":\"$DEFAULT\"}]")
  done
fi
```

### 4. Generate Setup Steps
```bash
echo ""
echo "📋 Generating setup steps..."

STEPS="["

# Add project creation step
if [ "$PACKAGE_MANAGER" = "npm" ] || [ "$PACKAGE_MANAGER" = "yarn" ] || [ "$PACKAGE_MANAGER" = "pnpm" ]; then
  STEPS="$STEPS{\"type\":\"command\",\"description\":\"Create project\",\"command\":\"mkdir {{project_name}} && cd {{project_name}}\"}"
  if [ -f "package.json" ]; then
    STEPS="$STEPS,{\"type\":\"file\",\"path\":\"{{project_name}}/package.json\",\"content\":$(jq -c . package.json | sed 's/"/\\"/g' | sed 's/\n/\\n/g')}"
  fi
elif [ "$PACKAGE_MANAGER" = "uv" ]; then
  STEPS="$STEPS{\"type\":\"command\",\"description\":\"Initialize UV project\",\"command\":\"uv init {{project_name}} --python 3.12\"}"
elif [ "$PACKAGE_MANAGER" = "cargo" ]; then
  STEPS="$STEPS{\"type\":\"command\",\"description\":\"Create Rust project\",\"command\":\"cargo new {{project_name}}\"}"
fi

# Add configuration files
for config in $CONFIG_FILES; do
  if [ -f "$config" ]; then
    content=$(cat "$config" | jq -R -s . | sed 's/"/\\"/g' | sed 's/\n/\\n/g')
    STEPS="$STEPS,{\"type\":\"file\",\"path\":\"{{project_name}}/$config\",\"content\":$content}"
  fi
done

STEPS="$STEPS]"
```

### 5. Create Template JSON
```bash
TEMPLATE_JSON="{
  \"name\": \"$1\",
  \"description\": \"$DESCRIPTION\",
  \"category\": \"$CATEGORY\",
  \"prerequisites\": [$PREREQUISITES],
  \"variables\": $VARIABLES,
  \"steps\": $STEPS,
  \"postInstall\": [
    \"echo '✅ $1 project created successfully!'\",
    \"echo '📁 Navigate to your project: cd {{project_name}}'\",
    \"echo '🚀 Run: $PACKAGE_MANAGER install && $PACKAGE_MANAGER run dev'\"
  ],
  \"documentation\": {
    \"nextSteps\": [
      \"Install dependencies\",
      \"Configure environment variables\",
      \"Start development server\"
    ]
  }
}"

echo "$TEMPLATE_JSON" | jq . > ~/mcf/templates/$1.json
```

### 6. Verification
```bash
echo "✅ Template '$1' saved to ~/mcf/templates/$1.json"
echo ""
echo "📋 Template Summary:"
echo "Name: $1"
echo "Description: $DESCRIPTION"
echo "Category: $CATEGORY"
echo "Prerequisites: $PREREQUISITES"
echo "Variables: $(echo "$VARIABLES" | jq length)"
echo "Steps: $(echo "$STEPS" | jq length)"
echo ""
echo "Test your template with: /template:init $1"
```
