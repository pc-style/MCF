#!/usr/bin/env python3
"""
Git Work Session Tracker - Suggests git commands based on work patterns
Tracks file modifications and suggests appropriate git commands
"""

import json
import sys
import subprocess
import os
import time
from pathlib import Path

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

def track_file_modifications(tool_input):
    """Track file modifications and suggest git commands"""
    file_path = tool_input.get("file_path", "")
    if not file_path or not os.path.exists(file_path):
        return []
    
    # Check if file is in git repository
    if not run_git_command("rev-parse --git-dir"):
        return []
    
    # Check if this file is tracked by git
    git_status = run_git_command(f"status --porcelain '{file_path}'")
    
    suggestions = []
    
    if git_status:
        status_code = git_status[:2] if len(git_status) >= 2 else ""
        
        if status_code.startswith("??"):
            suggestions.append(f"🆕 **Git Tip**: New file `{os.path.basename(file_path)}` created! Add it with `git add \"{file_path}\"`")
        elif status_code[1] == "M":  # Modified in working directory
            suggestions.append(f"✏️ **Git Tip**: File `{os.path.basename(file_path)}` was modified. Stage it with `git add \"{file_path}\"`")
        elif status_code[0] == "M":  # Modified and staged
            suggestions.append(f"📝 **Git Tip**: `{os.path.basename(file_path)}` is staged and ready to commit!")
    
    return suggestions

def get_work_session_suggestions():
    """Provide suggestions based on overall work session patterns"""
    
    # Get git status
    status = run_git_command("status --porcelain")
    if not status:
        return ["✅ **Git Tip**: Working directory is clean! Good time to start a new feature or pull latest changes."]
    
    lines = status.split('\n')
    staged_files = len([l for l in lines if l[0] in 'MARCDT'])
    unstaged_files = len([l for l in lines if l[1] in 'MARCDT'])
    untracked_files = len([l for l in lines if l.startswith('??')])
    
    suggestions = []
    
    # Suggest based on number of changes
    total_changes = staged_files + unstaged_files + untracked_files
    
    if total_changes >= 5:
        suggestions.append("📚 **Git Tip**: Lots of changes detected! Consider committing work in smaller, logical chunks.")
    
    if staged_files > 0 and unstaged_files > 0:
        suggestions.append("🎯 **Git Tip**: You have both staged and unstaged changes. Review with `git status` and commit when ready.")
    
    if unstaged_files >= 3:
        suggestions.append("📁 **Git Tip**: Multiple modified files detected. Use `git add .` to stage all, or `git add -p` for selective staging.")
    
    if untracked_files >= 2:
        suggestions.append("🔍 **Git Tip**: Several new files found. Review them with `git status` and add the ones you want to track.")
    
    return suggestions[:1]  # Return just one suggestion

def main():
    try:
        # Read input from stdin
        input_data = json.load(sys.stdin)
        
        tool_name = input_data.get("tool_name", "")
        tool_input = input_data.get("tool_input", {})
        hook_event = input_data.get("hook_event_name", "")
        
        suggestions = []
        
        # Handle PostToolUse events for file operations
        if hook_event == "PostToolUse" and tool_name in ["Write", "Edit", "MultiEdit"]:
            suggestions.extend(track_file_modifications(tool_input))
        
        # Handle Stop events - suggest based on work session
        elif hook_event == "Stop":
            suggestions.extend(get_work_session_suggestions())
        
        # Show suggestion if we have any
        if suggestions:
            import random
            suggestion = random.choice(suggestions)
            print(suggestion)
        
        sys.exit(0)
        
    except Exception:
        # Fail silently to avoid disrupting Claude Code
        sys.exit(0)

if __name__ == "__main__":
    main()