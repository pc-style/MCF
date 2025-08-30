#!/usr/bin/env python3
"""
Auto-formatting hook that runs appropriate formatters after file edits.
"""
import json
import sys
import os
import subprocess

def run_formatter(file_path):
    """Run appropriate formatter based on file extension."""
    _, ext = os.path.splitext(file_path)

    formatters = {
        '.js': ['npx', 'prettier', '--write'],
        '.ts': ['npx', 'prettier', '--write'],
        '.jsx': ['npx', 'prettier', '--write'],
        '.tsx': ['npx', 'prettier', '--write'],
        '.css': ['npx', 'prettier', '--write'],
        '.scss': ['npx', 'prettier', '--write'],
        '.json': ['npx', 'prettier', '--write'],
        '.py': ['black', '--quiet'],
        '.rs': ['rustfmt'],
        '.go': ['gofmt', '-w'],
        '.md': ['npx', 'prettier', '--write'],
        '.yaml': ['npx', 'prettier', '--write'],
        '.yml': ['npx', 'prettier', '--write'],
    }

    if ext in formatters:
        try:
            cmd = formatters[ext] + [file_path]
            subprocess.run(cmd, capture_output=True, timeout=30)
            print(f"✨ Formatted {os.path.basename(file_path)}")
        except:
            pass  # Silent fail

def main():
    try:
        input_data = json.load(sys.stdin)
        file_path = input_data.get('tool_input', {}).get('file_path', '')

        if file_path and os.path.exists(file_path):
            run_formatter(file_path)
    except:
        pass

if __name__ == "__main__":
    main()
