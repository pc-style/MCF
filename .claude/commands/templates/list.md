---
allowed-tools: Bash, Read, Glob, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
description: List all available project templates with descriptions
---

# List Available Project Templates

Display all saved project templates from `~/mcf/templates/` with their descriptions and usage examples.

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.


## Process:

1. **Check if templates directory exists** (`~/mcf/templates/`)
2. **Read all `.json` files** in the templates directory
3. **Parse each template** to extract name and description
4. **Display in a formatted table** with:
   - Template name
   - Description
   - Variables it accepts
   - Example usage

## Output Format:

```
📁 Available Project Templates:

┌─────────────────┬─────────────────────────────────────────────┬──────────────────────────┐
│ Template Name   │ Description                                 │ Usage Example            │
├─────────────────┼─────────────────────────────────────────────┼──────────────────────────┤
│ shadcn          │ Next.js + shadcn/ui + Tailwind CSS setup   │ /template:init shadcn    │
│ python-uv       │ Python project with UV package manager     │ /template:init python-uv │
│ vite-react      │ Vite + React + TypeScript starter          │ /template:init vite-react│
└─────────────────┴─────────────────────────────────────────────┴──────────────────────────┘

Use '/template:init <name>' to initialize a project from any template.
Use '/template:save <name>' to save your current project as a template.
```

## Implementation Steps:

### 1. Check Templates Directory
```bash
TEMPLATES_DIR="$HOME/mcf/templates"

if [ ! -d "$TEMPLATES_DIR" ]; then
  echo "❌ Templates directory not found: $TEMPLATES_DIR"
  echo "Run setup script first: ~/setup-templates.sh"
  exit 1
fi

# Count available templates
TEMPLATE_COUNT=$(ls "$TEMPLATES_DIR"/*.json 2>/dev/null | wc -l)

if [ "$TEMPLATE_COUNT" -eq 0 ]; then
  echo "📭 No templates found."
  echo ""
  echo "💡 To create templates:"
  echo "   /template:add <technology> - Auto-create from documentation"
  echo "   /template:save <name>      - Save current project as template"
  exit 0
fi
```

### 2. Generate Table Header
```bash
echo "📁 Available Project Templates:"
echo ""
printf "┌─────────────────┬─────────────────────────────────────────────┬──────────────────────────┐\n"
printf "│ Template Name   │ Description                                 │ Usage Example            │\n"
printf "├─────────────────┼─────────────────────────────────────────────┼──────────────────────────┤\n"
```

### 3. Process Each Template
```bash
for template_file in "$TEMPLATES_DIR"/*.json; do
  if [ ! -f "$template_file" ]; then continue; fi

  # Extract template name from filename
  TEMPLATE_NAME=$(basename "$template_file" .json)

  # Parse template JSON
  if ! TEMPLATE_DATA=$(cat "$template_file" 2>/dev/null); then
    printf "│ %-15s │ %-43s │ %-24s │\n" "$TEMPLATE_NAME" "Invalid JSON" "N/A"
    continue
  fi

  # Extract description
  DESCRIPTION=$(echo "$TEMPLATE_DATA" | jq -r '.description // "No description"' 2>/dev/null || echo "No description")

  # Truncate description if too long
  if [ ${#DESCRIPTION} -gt 43 ]; then
    DESCRIPTION="${DESCRIPTION:0:40}..."
  fi

  # Format usage example
  USAGE="/template:init $TEMPLATE_NAME"
  if [ ${#USAGE} -gt 24 ]; then
    USAGE="${USAGE:0:21}..."
  fi

  printf "│ %-15s │ %-43s │ %-24s │\n" "$TEMPLATE_NAME" "$DESCRIPTION" "$USAGE"
done
```

### 4. Generate Table Footer
```bash
printf "└─────────────────┴─────────────────────────────────────────────┴──────────────────────────┘\n"
echo ""
echo "📊 Total: $TEMPLATE_COUNT templates"
echo ""
echo "🚀 Usage:"
echo "   /template:init <name>     - Initialize project from template"
echo "   /template:save <name>     - Save current project as template"
echo "   /template:add <tech>      - Auto-create template from docs"
```

### 5. Show Additional Template Info
```bash
echo ""
echo "📋 Template Details:"
echo ""

for template_file in "$TEMPLATES_DIR"/*.json; do
  if [ ! -f "$template_file" ]; then continue; fi

  TEMPLATE_NAME=$(basename "$template_file" .json)
  TEMPLATE_DATA=$(cat "$template_file" 2>/dev/null)

  echo "🔹 $TEMPLATE_NAME"

  # Show category
  CATEGORY=$(echo "$TEMPLATE_DATA" | jq -r '.category // "Uncategorized"' 2>/dev/null)
  echo "   Category: $CATEGORY"

  # Show prerequisites
  PREREQ_COUNT=$(echo "$TEMPLATE_DATA" | jq '.prerequisites | length' 2>/dev/null || echo 0)
  if [ "$PREREQ_COUNT" -gt 0 ]; then
    PREREQS=$(echo "$TEMPLATE_DATA" | jq -r '.prerequisites | join(", ")' 2>/dev/null)
    echo "   Prerequisites: $PREREQS"
  fi

  # Show variable count
  VAR_COUNT=$(echo "$TEMPLATE_DATA" | jq '.variables | length' 2>/dev/null || echo 0)
  echo "   Variables: $VAR_COUNT configurable"

  # Show step count
  STEP_COUNT=$(echo "$TEMPLATE_DATA" | jq '.steps | length' 2>/dev/null || echo 0)
  echo "   Steps: $STEP_COUNT setup steps"

  echo ""
done
```

### 6. Show Next Steps
```bash
echo "💡 Quick Start:"
echo "1. Choose a template from the list above"
echo "2. Run: /template:init <template-name>"
echo "3. Follow the prompts to configure your project"
echo "4. Your project will be created and ready to use!"
```
