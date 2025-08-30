---
allowed-tools: Bash, Read, Write, Edit, MultiEdit
argument-hint: <template-name>
description: Initialize a new project from a saved template
---

# Project Template Initialization

You are initializing a project using the template: **$1**

## Process:

1. **Check if template exists** in `~/mcf/templates/$1.json`
2. **Read the template configuration** to understand the setup steps
3. **Execute each step** in the template, asking for user input when needed
4. **Create the project structure** as defined in the template
5. **Run any initialization commands** specified in the template
6. **Provide setup completion summary**

## Template Structure Expected:

Templates are JSON files with this structure:
```json
{
  "name": "template-name",
  "description": "Template description",
  "variables": [
    {
      "name": "project_name",
      "prompt": "What is your project name?",
      "default": "my-app"
    }
  ],
  "steps": [
    {
      "type": "command",
      "description": "Initialize project",
      "command": "npm create vite@latest {{project_name}} -- --template react"
    },
    {
      "type": "file",
      "path": "{{project_name}}/.env.example",
      "content": "# Environment variables\nVITE_API_URL=http://localhost:3000"
    }
  ]
}
```

## Your Task:

1. Check if the template `$1` exists
2. If it doesn't exist, list available templates
3. If it exists, follow the template steps and create the project
4. Handle variable substitution using the `{{variable_name}}` syntax
5. Provide clear feedback on each step

## Implementation Steps:

### 1. Check Template Existence
```bash
if [ ! -f "~/mcf/templates/$1.json" ]; then
  echo "❌ Template '$1' not found. Available templates:"
  ls ~/mcf/templates/*.json | sed 's/.*\///' | sed 's/\.json$//'
  exit 1
fi
```

### 2. Load Template Configuration
```bash
# Read the JSON template file
TEMPLATE_CONTENT=$(cat ~/mcf/templates/$1.json)
```

### 3. Process Variables
```bash
# Extract and process variables from template
VARIABLES=$(echo "$TEMPLATE_CONTENT" | jq -r '.variables[]? | @base64')

# Collect user input for each variable
for var in $VARIABLES; do
  VAR_DATA=$(echo "$var" | base64 -d)
  NAME=$(echo "$VAR_DATA" | jq -r '.name')
  PROMPT=$(echo "$VAR_DATA" | jq -r '.prompt')
  DEFAULT=$(echo "$VAR_DATA" | jq -r '.default // ""')

  # Prompt user for input
  read -p "$PROMPT [$DEFAULT]: " VALUE
  VALUE=${VALUE:-$DEFAULT}

  # Store variable for substitution
  export $NAME="$VALUE"
done
```

### 4. Execute Template Steps
```bash
# Extract steps from template
STEPS=$(echo "$TEMPLATE_CONTENT" | jq -r '.steps[]? | @base64')

# Execute each step
for step in $STEPS; do
  STEP_DATA=$(echo "$step" | base64 -d)
  STEP_TYPE=$(echo "$STEP_DATA" | jq -r '.type')
  DESCRIPTION=$(echo "$STEP_DATA" | jq -r '.description')

  echo "▶️ $DESCRIPTION"

  case $STEP_TYPE in
    "command")
      COMMAND=$(echo "$STEP_DATA" | jq -r '.command')
      # Substitute variables in command
      COMMAND=$(envsubst <<< "$COMMAND")
      eval "$COMMAND"
      ;;
    "directory")
      DIR_PATH=$(echo "$STEP_DATA" | jq -r '.path')
      DIR_PATH=$(envsubst <<< "$DIR_PATH")
      mkdir -p "$DIR_PATH"
      ;;
    "file")
      FILE_PATH=$(echo "$STEP_DATA" | jq -r '.path')
      FILE_CONTENT=$(echo "$STEP_DATA" | jq -r '.content')
      FILE_PATH=$(envsubst <<< "$FILE_PATH")
      FILE_CONTENT=$(envsubst <<< "$FILE_CONTENT")
      mkdir -p "$(dirname "$FILE_PATH")"
      echo "$FILE_CONTENT" > "$FILE_PATH"
      ;;
  esac
done
```

### 5. Run Post-Installation Steps
```bash
# Execute post-install commands
POST_INSTALL=$(echo "$TEMPLATE_CONTENT" | jq -r '.postInstall[]?')
for cmd in $POST_INSTALL; do
  cmd=$(envsubst <<< "$cmd")
  eval "$cmd"
done
```

### 6. Show Completion Summary
```bash
echo "✅ Project initialized successfully from template '$1'!"
echo ""
echo "📚 Next Steps:"
echo "$TEMPLATE_CONTENT" | jq -r '.documentation.nextSteps[]?' 2>/dev/null || echo "No specific next steps provided."
```
