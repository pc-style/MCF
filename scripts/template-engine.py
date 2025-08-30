#!/usr/bin/env python3
"""
MCF Template Engine
Handles template processing, variable substitution, and project generation.
"""
import json
import os
import sys
import re
import subprocess
from pathlib import Path
from typing import Dict, List, Any


class TemplateEngine:
    def __init__(self):
        self.mcf_dir = Path.home() / "mcf"
        self.templates_dir = self.mcf_dir / "templates"
        self.variables = {}

    def list_templates(self) -> List[Dict[str, Any]]:
        """List all available templates."""
        if not self.templates_dir.exists():
            return []
        
        templates = []
        for template_file in self.templates_dir.glob("*.json"):
            try:
                with open(template_file, 'r') as f:
                    template = json.load(f)
                    template['filename'] = template_file.stem
                    templates.append(template)
            except (json.JSONDecodeError, KeyError) as e:
                print(f"Warning: Invalid template file {template_file}: {e}")
        
        return sorted(templates, key=lambda x: x.get('name', ''))

    def load_template(self, template_name: str) -> Dict[str, Any]:
        """Load a specific template by name."""
        template_file = self.templates_dir / f"{template_name}.json"
        
        if not template_file.exists():
            raise FileNotFoundError(f"Template '{template_name}' not found")
        
        with open(template_file, 'r') as f:
            return json.load(f)

    def collect_variables(self, template: Dict[str, Any]) -> Dict[str, str]:
        """Collect variable values from user input."""
        variables = {}
        
        if 'variables' not in template:
            return variables
        
        print(f"\n🔧 Configuring template: {template['name']}")
        print(f"📝 {template.get('description', '')}\n")
        
        for var_config in template['variables']:
            name = var_config['name']
            prompt = var_config.get('prompt', f"Value for {name}")
            default = var_config.get('default', '')
            options = var_config.get('options', [])
            
            if options:
                prompt += f" ({'/'.join(options)})"
            if default:
                prompt += f" [default: {default}]"
            
            while True:
                value = input(f"{prompt}: ").strip()
                
                if not value and default:
                    value = default
                
                if not value:
                    print("❌ This field is required")
                    continue
                
                if options and value not in options:
                    print(f"❌ Must be one of: {', '.join(options)}")
                    continue
                
                # Validate if pattern provided
                if 'validation' in var_config:
                    pattern = var_config['validation']
                    if not re.match(pattern, value):
                        print(f"❌ Invalid format (expected: {pattern})")
                        continue
                
                variables[name] = value
                break
        
        return variables

    def substitute_variables(self, text: str, variables: Dict[str, str]) -> str:
        """Replace {{variable}} placeholders with actual values."""
        result = text
        for name, value in variables.items():
            result = result.replace(f"{{{{{name}}}}}", value)
        return result

    def execute_step(self, step: Dict[str, Any], variables: Dict[str, str], current_dir: Path) -> bool:
        """Execute a single template step."""
        step_type = step.get('type', 'command')
        description = step.get('description', 'Executing step')
        
        print(f"  ▶️ {description}")
        
        if step_type == 'command':
            command = self.substitute_variables(step['command'], variables)
            
            try:
                result = subprocess.run(
                    command,
                    shell=True,
                    cwd=current_dir,
                    capture_output=True,
                    text=True,
                    timeout=300  # 5 minute timeout
                )
                
                if result.returncode != 0:
                    print(f"❌ Command failed: {command}")
                    print(f"Error: {result.stderr}")
                    return False
                
                if result.stdout.strip():
                    print(f"   {result.stdout.strip()}")
                
            except subprocess.TimeoutExpired:
                print(f"❌ Command timed out: {command}")
                return False
            except Exception as e:
                print(f"❌ Error executing command: {e}")
                return False
        
        elif step_type == 'directory':
            path = self.substitute_variables(step['path'], variables)
            dir_path = Path(path)
            
            if not dir_path.is_absolute():
                dir_path = current_dir / dir_path
            
            try:
                dir_path.mkdir(parents=True, exist_ok=True)
                print(f"   📁 Created directory: {dir_path}")
            except Exception as e:
                print(f"❌ Error creating directory: {e}")
                return False
        
        elif step_type == 'file':
            file_path = self.substitute_variables(step['path'], variables)
            content = self.substitute_variables(step.get('content', ''), variables)
            
            full_path = Path(file_path)
            if not full_path.is_absolute():
                full_path = current_dir / full_path
            
            try:
                full_path.parent.mkdir(parents=True, exist_ok=True)
                with open(full_path, 'w') as f:
                    f.write(content)
                print(f"   📄 Created file: {full_path}")
            except Exception as e:
                print(f"❌ Error creating file: {e}")
                return False
        
        else:
            print(f"⚠️ Unknown step type: {step_type}")
        
        return True

    def execute_template(self, template_name: str, target_dir: Path = None) -> bool:
        """Execute a complete template."""
        try:
            template = self.load_template(template_name)
        except FileNotFoundError as e:
            print(f"❌ {e}")
            return False
        
        # Collect variables
        variables = self.collect_variables(template)
        
        # Set working directory
        if target_dir is None:
            target_dir = Path.cwd()
        
        print(f"\n🚀 Initializing project from template: {template['name']}")
        
        # Check prerequisites
        if 'prerequisites' in template:
            print("🔍 Checking prerequisites...")
            for prereq in template['prerequisites']:
                if not self.check_prerequisite(prereq):
                    print(f"❌ Missing prerequisite: {prereq}")
                    return False
            print("✅ All prerequisites satisfied")
        
        # Execute steps
        print("\n📋 Executing template steps:")
        for i, step in enumerate(template.get('steps', []), 1):
            print(f"\n{i}. {step.get('description', 'Step ' + str(i))}")
            if not self.execute_step(step, variables, target_dir):
                print(f"❌ Template execution failed at step {i}")
                return False
        
        # Execute post-install commands
        if 'postInstall' in template:
            print("\n🎯 Running post-installation steps:")
            for command in template['postInstall']:
                cmd = self.substitute_variables(command, variables)
                print(f"  ▶️ {cmd}")
                os.system(cmd)
        
        # Show next steps
        if 'documentation' in template and 'nextSteps' in template['documentation']:
            print("\n📚 Suggested next steps:")
            for step in template['documentation']['nextSteps']:
                print(f"  • {step}")
        
        print(f"\n✅ Template '{template['name']}' executed successfully!")
        return True

    def check_prerequisite(self, prereq: str) -> bool:
        """Check if a prerequisite is available."""
        try:
            result = subprocess.run(
                f"command -v {prereq}",
                shell=True,
                capture_output=True,
                text=True
            )
            return result.returncode == 0
        except:
            return False

    def save_template(self, template_name: str, template_data: Dict[str, Any]) -> bool:
        """Save a template to disk."""
        try:
            self.templates_dir.mkdir(parents=True, exist_ok=True)
            template_file = self.templates_dir / f"{template_name}.json"
            
            with open(template_file, 'w') as f:
                json.dump(template_data, f, indent=2)
            
            print(f"✅ Template saved: {template_file}")
            return True
        except Exception as e:
            print(f"❌ Error saving template: {e}")
            return False


def main():
    """Command line interface for the template engine."""
    if len(sys.argv) < 2:
        print("Usage: template-engine.py <command> [args...]")
        print("Commands:")
        print("  list                    - List available templates")
        print("  init <template-name>    - Initialize project from template")
        print("  info <template-name>    - Show template information")
        return
    
    engine = TemplateEngine()
    command = sys.argv[1]
    
    if command == 'list':
        templates = engine.list_templates()
        if not templates:
            print("📭 No templates found. Use '/template:add-template' to create some!")
            return
        
        print("📁 Available Templates:")
        print()
        for template in templates:
            print(f"  {template['name']:<15} - {template.get('description', 'No description')}")
        print()
        print(f"Total: {len(templates)} templates")
    
    elif command == 'init':
        if len(sys.argv) < 3:
            print("❌ Please specify a template name")
            return
        
        template_name = sys.argv[2]
        engine.execute_template(template_name)
    
    elif command == 'info':
        if len(sys.argv) < 3:
            print("❌ Please specify a template name")
            return
        
        template_name = sys.argv[2]
        try:
            template = engine.load_template(template_name)
            print(f"📋 Template: {template['name']}")
            print(f"📝 Description: {template.get('description', 'No description')}")
            print(f"🏷️  Category: {template.get('category', 'Uncategorized')}")
            
            if 'prerequisites' in template:
                print(f"⚙️  Prerequisites: {', '.join(template['prerequisites'])}")
            
            if 'variables' in template:
                print("🔧 Variables:")
                for var in template['variables']:
                    print(f"  • {var['name']}: {var.get('prompt', 'No description')}")
            
            print(f"📋 Steps: {len(template.get('steps', []))}")
            
        except FileNotFoundError:
            print(f"❌ Template '{template_name}' not found")
    
    else:
        print(f"❌ Unknown command: {command}")


if __name__ == "__main__":
    main()