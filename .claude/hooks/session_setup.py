#!/usr/bin/env python3
"""
Enhanced Session Setup Hook with Serena Integration
Provides context about available tools, MCPs, Serena capabilities, and project-specific information at session start.
"""
import json
import sys
import os
import subprocess
from pathlib import Path


def get_git_info():
    """Get current git branch and recent commits."""
    try:
        # Get current branch
        branch = subprocess.run(
            ["git", "branch", "--show-current"],
            capture_output=True,
            text=True,
            check=True,
        )
        current_branch = branch.stdout.strip()

        # Get recent commits
        commits = subprocess.run(
            ["git", "log", "--oneline", "-5"],
            capture_output=True,
            text=True,
            check=True,
        )
        recent_commits = commits.stdout.strip()

        return {"branch": current_branch, "recent_commits": recent_commits}
    except:
        return None


def get_project_info():
    """Extract project information from common config files."""
    info = {}

    # Check package.json
    if os.path.exists("package.json"):
        try:
            with open("package.json", "r") as f:
                pkg = json.load(f)
                info["name"] = pkg.get("name", "")
                info["version"] = pkg.get("version", "")
                info["description"] = pkg.get("description", "")
                info["main_deps"] = list(pkg.get("dependencies", {}).keys())[:5]
        except:
            pass

    # Check pyproject.toml or requirements.txt for Python projects
    elif os.path.exists("pyproject.toml"):
        info["type"] = "Python (pyproject.toml)"
    elif os.path.exists("requirements.txt"):
        info["type"] = "Python (requirements.txt)"

    # Check for other common project indicators
    elif os.path.exists("Cargo.toml"):
        info["type"] = "Rust"
    elif os.path.exists("go.mod"):
        info["type"] = "Go"
    elif os.path.exists("pom.xml"):
        info["type"] = "Java (Maven)"
    elif os.path.exists("build.gradle"):
        info["type"] = "Java/Kotlin (Gradle)"

    return info


def check_serena_status():
    """Check Serena MCP integration status."""
    serena_info = {
        "mcp_configured": False,
        "global_config": False,
        "project_initialized": False,
        "memories_count": 0,
        "compatible_files": 0,
    }

    # Check if Serena MCP is configured
    try:
        result = subprocess.run(
            ["claude", "mcp", "list"], capture_output=True, text=True, check=True
        )
        if "serena" in result.stdout:
            serena_info["mcp_configured"] = True
    except:
        pass

    # Check global Serena config
    serena_config_path = Path.home() / ".serena" / "serena_config.yml"
    if serena_config_path.exists():
        serena_info["global_config"] = True

    # Check project initialization
    project_config = Path(".serena/project.yml")
    if project_config.exists():
        serena_info["project_initialized"] = True

    # Count memories
    memories_dir = Path(".serena/memories")
    if memories_dir.exists():
        serena_info["memories_count"] = len(list(memories_dir.glob("*.md")))

    # Count compatible files
    compatible_extensions = {
        ".py",
        ".js",
        ".ts",
        ".tsx",
        ".jsx",
        ".md",
        ".sh",
        ".yml",
        ".yaml",
        ".json",
    }
    for ext in compatible_extensions:
        if list(Path(".").rglob(f"*{ext}")):
            serena_info["compatible_files"] += 1
            if serena_info["compatible_files"] >= 5:  # Stop counting after 5 types
                break

    return serena_info


def main():
    try:
        # Read JSON input
        input_data = json.load(sys.stdin)

        # Get project, git, and Serena information
        project_info = get_project_info()
        git_info = get_git_info()
        serena_info = check_serena_status()

        # Build context message
        context_parts = []

        # MCP Tools section (enhanced with Serena)
        mcp_section = """🔧 **Available MCP Tools**:"""

        if serena_info["mcp_configured"]:
            mcp_section += """
                - **Serena**: Semantic code analysis with IDE-like capabilities ✨
  • find_symbol, get_symbol_info, find_referencing_symbols
  • Symbol-level editing and project structure analysis"""

        mcp_section += """
                - **Context7**: Up-to-date documentation on any framework/library
                - **Google Drive**: Access your Drive files and documents
                - **Exa**: Advanced web search capabilities"""

        context_parts.append(mcp_section)

        # Serena-specific status and recommendations
        if serena_info["compatible_files"] > 0:
            serena_section = f"🧠 **Serena Status**:"

            if serena_info["mcp_configured"] and serena_info["project_initialized"]:
                serena_section += f"""
                - ✅ Ready for semantic analysis ({serena_info['compatible_files']} file types supported)"""

                if serena_info["memories_count"] > 0:
                    serena_section += f"""
                - ✅ {serena_info['memories_count']} project memories loaded"""

                serena_section += """
                - 💡 Try: `/serena:overview` or `/serena:find <symbol>`"""

            elif (
                serena_info["mcp_configured"] and not serena_info["project_initialized"]
            ):
                serena_section += """
                - ⚠️  MCP configured but project not initialized
                - 💡 Run: `/serena:init` to enable semantic analysis"""

            elif not serena_info["mcp_configured"]:
                serena_section += """
                - 💡 Semantic code analysis available! Run: `/serena:install`"""

            context_parts.append(serena_section)

        # Project information
        if project_info:
            context_parts.append("📋 **Project Info**:")
            if "name" in project_info:
                context_parts.append(f"- Name: {project_info['name']}")
            if "description" in project_info:
                context_parts.append(f"- Description: {project_info['description']}")
            if "main_deps" in project_info:
                deps = ", ".join(project_info["main_deps"])
                context_parts.append(f"- Key dependencies: {deps}")
            if "type" in project_info:
                context_parts.append(f"- Project type: {project_info['type']}")

        # Git information
        if git_info:
            context_parts.append(f"🌿 **Git Status**:")
            context_parts.append(f"- Current branch: {git_info['branch']}")
            if git_info["recent_commits"]:
                context_parts.append("- Recent commits:")
                for line in git_info["recent_commits"].split("\n")[:3]:
                    context_parts.append(f"  • {line}")

        # Custom commands section (enhanced with Serena commands)
        commands_section = """⚡ **Your Custom Commands**:
                - `/gh:push`: Quick add, commit, push
                - `/gh:commit [msg]`: Add and commit with message
                - `/gh:auto <instruction>`: Natural language Git operations"""

        if serena_info["mcp_configured"]:
            commands_section += """
                - `/serena:overview`: Get project structure overview
                - `/serena:find <symbol>`: Locate functions/classes/variables
                - `/serena:status`: Check Serena integration health"""

        commands_section += """
                - Use `/help` to see all available commands"""

        context_parts.append(commands_section)

        # Pro tips based on setup
        pro_tips = "💡 **Pro Tips**:"

        if serena_info["mcp_configured"] and serena_info["project_initialized"]:
            pro_tips += """
                - Use semantic analysis to understand code faster and save tokens
                - Let Serena find symbols before reading entire files"""

        pro_tips += """
                - Ask me to use Context7 when working with external frameworks
                - Monitor Serena usage at http://localhost:24282/dashboard/ (when active)"""

        context_parts.append(pro_tips)

        # Output the context
        print("\n".join(context_parts))

    except Exception as e:
        # Silent fail to not interrupt workflow
        pass


if __name__ == "__main__":
    main()
