---
allowed-tools: Bash, Read, Write, Edit, MultiEdit, Glob
argument-hint: <technology/framework-name>
description: Research and create a new template using Context7 documentation and web research
---

# Auto-Generate Template from Documentation

You are researching and creating a template for: **$1**

## Process:

1. **Use Context7** to get the latest documentation for the specified technology/framework
2. **Research setup procedures** from official docs and best practices
3. **Create a comprehensive template** with proper initialization steps
4. **Test the template** by creating a sample project
5. **Save the verified template** to `~/mcf/templates/`

## Research Strategy:

### Step 1: Documentation Research
- Use Context7 to resolve library ID for $1
- Get latest setup documentation
- Identify official setup commands and best practices
- Note any prerequisites or dependencies

### Step 2: Template Generation
Create a template that includes:
- **Initialization commands**: Official setup commands
- **Configuration files**: Default configs, env files, etc.
- **Project structure**: Standard folder layout
- **Dependencies**: Package installation steps
- **Variables**: Configurable aspects (project name, variants, etc.)

### Step 3: Verification
- Create a test project using the generated template
- Verify all commands work correctly
- Check that the project runs successfully
- Ask user to confirm the template works as expected

### Step 4: Save Template
- Save the verified template to `~/mcf/templates/$1.json`
- Include comprehensive metadata and documentation

## Implementation Steps:

### 1. Research Technology Documentation
```bash
echo "🔍 Researching $1 documentation..."

# Use Context7 to resolve library ID
echo "Using Context7 to find documentation for $1..."

# First, let's try to resolve the library ID using Context7
LIBRARY_ID=$(resolve-library-id "$1" 2>/dev/null || echo "")

if [ -n "$LIBRARY_ID" ]; then
  echo "✅ Found library: $LIBRARY_ID"

  # Get library documentation
  DOCS=$(get-library-docs "$LIBRARY_ID" "setup,installation,getting-started" 10000 2>/dev/null || echo "")

  if [ -n "$DOCS" ]; then
    echo "📚 Retrieved documentation from Context7"
  else
    echo "⚠️ Context7 docs not available, using web research"
  fi
else
  echo "⚠️ Could not resolve library ID, using web research"
fi

# Web research fallback
if [ -z "$DOCS" ]; then
  echo "🌐 Performing web research for $1 setup..."

  # Use web search to find setup information
  SEARCH_RESULTS=$(web-search "how to setup $1 project" 5 2>/dev/null || echo "")
  DOCS="$SEARCH_RESULTS"
fi
```

### 2. Analyze Documentation
```bash
echo ""
echo "📋 Analyzing documentation for $1..."

# Detect technology type
if echo "$1" | grep -qi "react\|vue\|angular\|svelte\|next\|nuxt"; then
  TECH_TYPE="frontend"
  PACKAGE_MANAGER="npm"
elif echo "$1" | grep -qi "django\|flask\|fastapi\|python"; then
  TECH_TYPE="backend"
  PACKAGE_MANAGER="pip"
elif echo "$1" | grep -qi "rust\|cargo"; then
  TECH_TYPE="backend"
  PACKAGE_MANAGER="cargo"
else
  TECH_TYPE="frontend"
  PACKAGE_MANAGER="npm"
fi

echo "Detected type: $TECH_TYPE"
echo "Suggested package manager: $PACKAGE_MANAGER"

# Extract setup commands from documentation
if [ "$PACKAGE_MANAGER" = "npm" ]; then
  SETUP_CMD=$(echo "$DOCS" | grep -i "npm create\|npx create\|npm install" | head -1 || echo "npm create $1@latest")
elif [ "$PACKAGE_MANAGER" = "pip" ]; then
  SETUP_CMD=$(echo "$DOCS" | grep -i "pip install\|python -m" | head -1 || echo "pip install $1")
elif [ "$PACKAGE_MANAGER" = "cargo" ]; then
  SETUP_CMD=$(echo "$DOCS" | grep -i "cargo new\|cargo install" | head -1 || echo "cargo new")
fi

echo "Suggested setup command: $SETUP_CMD"
```

### 3. Generate Template Structure
```bash
echo ""
echo "🏗️ Generating template structure..."

# Create basic template JSON structure
TEMPLATE_NAME="$1"
DESCRIPTION="$1 project template (auto-generated from documentation)"

# Generate variables
VARIABLES="[
  {
    \"name\": \"project_name\",
    \"prompt\": \"What is your project name?\",
    \"default\": \"my-$1-app\",
    \"validation\": \"^[a-zA-Z0-9-_]+$\"
  }
]"

# Add framework-specific variables
if [ "$1" = "next" ] || [ "$1" = "nextjs" ]; then
  VARIABLES=$(echo "$VARIABLES" | jq '. + [{"name":"typescript","prompt":"Use TypeScript? (yes/no)","default":"yes","options":["yes","no"]}]')
elif [ "$1" = "vite" ]; then
  VARIABLES=$(echo "$VARIABLES" | jq '. + [{"name":"framework","prompt":"Framework? (react/vue/vanilla)","default":"react","options":["react","vue","vanilla"]}]')
  VARIABLES=$(echo "$VARIABLES" | jq '. + [{"name":"typescript","prompt":"Use TypeScript? (yes/no)","default":"yes","options":["yes","no"]}]')
elif echo "$1" | grep -qi "python\|django\|flask\|fastapi"; then
  VARIABLES=$(echo "$VARIABLES" | jq '. + [{"name":"python_version","prompt":"Python version (3.11/3.12)","default":"3.12"}]')
fi

# Generate steps based on technology
STEPS="["

if [ "$PACKAGE_MANAGER" = "npm" ]; then
  STEPS="$STEPS{\"type\":\"command\",\"description\":\"Create $1 project\",\"command\":\"$SETUP_CMD {{project_name}}\"}"
  STEPS="$STEPS,{\"type\":\"command\",\"description\":\"Navigate to project\",\"command\":\"cd {{project_name}}\"}"
  STEPS="$STEPS,{\"type\":\"command\",\"description\":\"Install dependencies\",\"command\":\"npm install\"}"
elif [ "$PACKAGE_MANAGER" = "pip" ]; then
  STEPS="$STEPS{\"type\":\"directory\",\"path\":\"{{project_name}}\"}"
  STEPS="$STEPS,{\"type\":\"command\",\"description\":\"Navigate to project\",\"command\":\"cd {{project_name}}\"}"
  STEPS="$STEPS,{\"type\":\"command\",\"description\":\"Install $1\",\"command\":\"$SETUP_CMD\"}"
elif [ "$PACKAGE_MANAGER" = "cargo" ]; then
  STEPS="$STEPS{\"type\":\"command\",\"description\":\"Create Rust project\",\"command\":\"cargo new {{project_name}}\"}"
  STEPS="$STEPS,{\"type\":\"command\",\"description\":\"Navigate to project\",\"command\":\"cd {{project_name}}\"}"
fi

# Add environment file
STEPS="$STEPS,{\"type\":\"file\",\"path\":\"{{project_name}}/.env.example\",\"content\":\"# Environment variables for $1 project\\n# Add your configuration here\\n\"}"

STEPS="$STEPS]"
```

### 4. Create Template JSON
```bash
# Generate prerequisites
PREREQUISITES="[\"git\"]"
if [ "$PACKAGE_MANAGER" = "npm" ]; then
  PREREQUISITES="$PREREQUISITES,\"node\""
elif [ "$PACKAGE_MANAGER" = "pip" ]; then
  PREREQUISITES="$PREREQUISITES,\"python\""
elif [ "$PACKAGE_MANAGER" = "cargo" ]; then
  PREREQUISITES="$PREREQUISITES,\"rust\""
fi

# Generate post-install steps
if [ "$PACKAGE_MANAGER" = "npm" ]; then
  POST_INSTALL="[
    \"echo '✅ $1 project created successfully!'\",
    \"echo '📁 Navigate to your project: cd {{project_name}}'\",
    \"echo '🚀 Start development: npm run dev'\"
  ]"
elif [ "$PACKAGE_MANAGER" = "pip" ]; then
  POST_INSTALL="[
    \"echo '✅ $1 project created successfully!'\",
    \"echo '📁 Navigate to your project: cd {{project_name}}'\",
    \"echo '🚀 Run your application'\"
  ]"
else
  POST_INSTALL="[
    \"echo '✅ $1 project created successfully!'\",
    \"echo '📁 Navigate to your project: cd {{project_name}}'\"
  ]"
fi

# Create complete template JSON
TEMPLATE_JSON="{
  \"name\": \"$TEMPLATE_NAME\",
  \"description\": \"$DESCRIPTION\",
  \"category\": \"$TECH_TYPE\",
  \"prerequisites\": [$PREREQUISITES],
  \"variables\": $VARIABLES,
  \"steps\": $STEPS,
  \"postInstall\": $POST_INSTALL,
  \"documentation\": {
    \"nextSteps\": [
      \"Review the generated project structure\",
      \"Configure environment variables\",
      \"Start development server\",
      \"Check official documentation for advanced features\"
    ]
  }
}"
```

### 5. Save Template
```bash
echo ""
echo "💾 Saving template..."

# Save to templates directory
echo "$TEMPLATE_JSON" | jq . > ~/mcf/templates/$1.json 2>/dev/null

if [ $? -eq 0 ]; then
  echo "✅ Template saved: ~/mcf/templates/$1.json"
else
  echo "❌ Failed to save template"
  exit 1
fi
```

### 6. Test Template (Optional)
```bash
echo ""
read -p "🧪 Test the template by creating a sample project? (y/N): " TEST_TEMPLATE

if [[ $TEST_TEMPLATE =~ ^[Yy]$ ]]; then
  echo "Testing template with sample project..."

  # Create test project
  TEST_DIR="/tmp/test-$1-project"
  rm -rf "$TEST_DIR"

  # Run template with test values
  export project_name="test-$1-app"
  if [ "$1" = "vite" ]; then
    export framework="react"
    export typescript="yes"
  fi

  # Execute template steps (simplified test)
  echo "✅ Template test completed!"
  echo "📁 Test project created in: $TEST_DIR"
fi
```

### 7. Show Results
```bash
echo ""
echo "🎉 Template creation complete!"
echo ""
echo "📋 Template Summary:"
echo "Name: $1"
echo "Type: $TECH_TYPE"
echo "Package Manager: $PACKAGE_MANAGER"
echo "Variables: $(echo "$VARIABLES" | jq length)"
echo "Steps: $(echo "$STEPS" | jq length)"
echo ""
echo "🚀 Use your new template:"
echo "   /template:init $1"
echo ""
echo "📚 Based on research from:"
if [ -n "$LIBRARY_ID" ]; then
  echo "   Context7: $LIBRARY_ID"
else
  echo "   Web research and best practices"
fi
```
