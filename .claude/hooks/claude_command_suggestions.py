#!/usr/bin/env python3
"""
Claude Code Command Suggestion Hook
Suggests relevant Claude Code slash commands based on context and repository state
"""

import json
import sys
import subprocess
import os
import random
from datetime import datetime, timedelta

def run_git_command(cmd):
    """Run a git command and return output safely"""
    try:
        result = subprocess.run(
            f"git {cmd}",
            shell=True,
            capture_output=True,
            text=True,
            timeout=5
        )
        if result.returncode == 0:
            return result.stdout.strip()
        return None
    except:
        return None

def get_git_status():
    """Get current git repository status"""
    status = run_git_command("status --porcelain")
    if not status:
        return {"clean": True, "staged": 0, "unstaged": 0, "untracked": 0}
    
    lines = status.split('\n') if status else []
    staged = len([l for l in lines if l[0] in 'MARCDT'])
    unstaged = len([l for l in lines if l[1] in 'MARCDT'])
    untracked = len([l for l in lines if l.startswith('??')])
    
    return {
        "clean": len(lines) == 0,
        "staged": staged,
        "unstaged": unstaged,
        "untracked": untracked,
        "total_changes": len(lines)
    }

def get_branch_info():
    """Get current branch and remote tracking info"""
    current_branch = run_git_command("rev-parse --abbrev-ref HEAD")
    if not current_branch:
        return {}
    
    # Check if branch has upstream
    upstream = run_git_command(f"rev-parse --abbrev-ref {current_branch}@{{upstream}}")
    
    # Get commit counts
    ahead = run_git_command(f"rev-list --count {current_branch}@{{upstream}}..HEAD") if upstream else None
    behind = run_git_command(f"rev-list --count HEAD..{current_branch}@{{upstream}}") if upstream else None
    
    return {
        "branch": current_branch,
        "has_upstream": upstream is not None,
        "ahead": int(ahead) if ahead and ahead.isdigit() else 0,
        "behind": int(behind) if behind and behind.isdigit() else 0
    }

def should_suggest(hook_event):
    """Determine if we should suggest based on frequency and context"""
    # Create a state file to track suggestions
    state_file = os.path.expanduser("~/.claude/claude-cmd-suggestions-state.json")
    
    try:
        if os.path.exists(state_file):
            with open(state_file, 'r') as f:
                state = json.load(f)
        else:
            state = {"last_suggestion": None, "suggestion_count": 0}
    except:
        state = {"last_suggestion": None, "suggestion_count": 0}
    
    now = datetime.now().isoformat()
    last_suggestion = state.get("last_suggestion")
    
    # Don't suggest too frequently
    if last_suggestion:
        try:
            last_time = datetime.fromisoformat(last_suggestion)
            if datetime.now() - last_time < timedelta(minutes=3):
                return False
        except:
            pass
    
    # Suggest based on different triggers and randomness
    suggest_probability = {
        "UserPromptSubmit": 0.20,  # 20% chance on user prompts
        "Stop": 0.30,              # 30% chance when Claude stops
        "PostToolUse": 0.15        # 15% chance after file operations
    }
    
    should_suggest_now = random.random() < suggest_probability.get(hook_event, 0.1)
    
    if should_suggest_now:
        state["last_suggestion"] = now
        state["suggestion_count"] = state.get("suggestion_count", 0) + 1
        
        # Save state
        try:
            os.makedirs(os.path.dirname(state_file), exist_ok=True)
            with open(state_file, 'w') as f:
                json.dump(state, f)
        except:
            pass
    
    return should_suggest_now

def generate_claude_command_suggestions():
    """Generate contextual Claude Code command suggestions"""
    
    # Check if we're in a git repository
    is_git_repo = run_git_command("rev-parse --git-dir") is not None
    
    suggestions = []
    
    if is_git_repo:
        status = get_git_status()
        branch_info = get_branch_info()
        
        # Git workflow suggestions
        if status["total_changes"] > 0:
            if status["unstaged"] > 0 and status["staged"] == 0:
                suggestions.extend([
                    "🚀 **Claude Tip**: Got changes to commit? Try `/gh:commit \"Your commit message\"` to add and commit quickly!",
                    "⚡ **Claude Tip**: Use `/gh:push` to add, commit with smart message, and push in one command!",
                ])
            elif status["staged"] > 0:
                suggestions.extend([
                    "📝 **Claude Tip**: Ready to commit staged changes? Use `/gh:commit \"Your message\"` or just `/gh:push`!",
                    "🎯 **Claude Tip**: Try `/gh:push` to commit and push your staged changes automatically!",
                ])
        
        if branch_info.get("ahead", 0) > 0:
            suggestions.append("⬆️ **Claude Tip**: Branch is ahead! Use `/gh:push` to push your commits to remote.")
        
        # General git suggestions
        if status["total_changes"] > 3:
            suggestions.append("🔄 **Claude Tip**: Lots of changes? Try `/gh:auto \"create feature branch\"` for branch management!")
        
        # Advanced git operations
        git_advanced_suggestions = [
            "🌿 **Claude Tip**: Need to undo something? Use `/gh:revert` to safely revert commits!",
            "🎪 **Claude Tip**: Try `/gh:auto \"merge main into current branch\"` for natural language git operations!",
            "🔄 **Claude Tip**: Use `/gh:auto \"stash my changes\"` when you need to switch contexts quickly!",
        ]
        suggestions.extend(random.sample(git_advanced_suggestions, min(1, len(git_advanced_suggestions))))
    
    else:
        # Not a git repo - suggest initialization
        suggestions.extend([
            "🎬 **Claude Tip**: Not a git repo yet? Try `/gh:init` to initialize git and create a GitHub repository!",
            "📚 **Claude Tip**: Start tracking your project! Use `/gh:init` to set up git and GitHub integration.",
        ])
    
    # General Claude Code workflow suggestions
    general_suggestions = [
        "💡 **Claude Tip**: Use `/help` to see all available commands and shortcuts!",
        "🧠 **Claude Tip**: Try `/memory` to edit your CLAUDE.md file and give me project context!",
        "🔍 **Claude Tip**: Use `/review` to get a comprehensive code review of your changes!",
        "📊 **Claude Tip**: Check your token usage with `/cost` to track your API consumption!",
        "🎛️ **Claude Tip**: Use `/config` to customize your Claude Code settings and preferences!",
        "📝 **Claude Tip**: Try `/compact` to compress our conversation when it gets too long!",
        "🤖 **Claude Tip**: Use `/agents` to create specialized AI subagents for specific tasks!",
        "🔌 **Claude Tip**: Check `/mcp` to see available integrations and connect new tools!",
        "🎯 **Claude Tip**: Use `/init` to create a CLAUDE.md guide for your project!",
        "⚡ **Claude Tip**: Press `Ctrl+R` to see our conversation transcript with all the details!",
    ]
    
    # Add some general suggestions
    if len(suggestions) < 2:
        suggestions.extend(random.sample(general_suggestions, min(2, len(general_suggestions))))
    
    return suggestions[:2]  # Return max 2 suggestions

def main():
    try:
        # Read input from stdin
        input_data = json.load(sys.stdin)
        hook_event = input_data.get("hook_event_name", "")
        
        # Only suggest sometimes to avoid being annoying
        if not should_suggest(hook_event):
            sys.exit(0)
        
        # Generate suggestions
        suggestions = generate_claude_command_suggestions()
        
        if suggestions:
            # Pick one random suggestion to show
            suggestion = random.choice(suggestions)
            
            # Format the output based on the hook type
            if hook_event == "UserPromptSubmit":
                # For UserPromptSubmit, we can add context that Claude will see
                output = {
                    "hookSpecificOutput": {
                        "hookEventName": "UserPromptSubmit",
                        "additionalContext": f"\n{suggestion}\n"
                    }
                }
                print(json.dumps(output))
            else:
                # For other events, just print the suggestion
                print(suggestion)
        
        sys.exit(0)
        
    except Exception as e:
        # Fail silently to avoid disrupting Claude Code
        sys.exit(0)

if __name__ == "__main__":
    main()