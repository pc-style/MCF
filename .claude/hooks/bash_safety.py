#!/usr/bin/env python3
"""
Bash command safety checker that warns about potentially dangerous operations.
"""
import json
import sys
import re

# Dangerous command patterns
DANGEROUS_PATTERNS = [
    (r'\brm\s+-rf\s+/[^/\s]*\s*$', 'Dangerous: rm -rf on root-level directory'),
    (r'\bsudo\s+rm\s+-rf', 'Dangerous: sudo rm -rf command'),
    (r'\bchmod\s+-R\s+777', 'Security risk: chmod 777 -R'),
    (r'\bcurl\s+[^|]*\|\s*bash', 'Security risk: piping curl to bash'),
    (r'\bwget\s+[^|]*\|\s*bash', 'Security risk: piping wget to bash'),
    (r'>\s*/dev/sd[a-z]', 'Dangerous: writing to raw disk device'),
    (r'\bdd\s+.*of=/dev/', 'Dangerous: dd command to device'),
]

def check_command_safety(command):
    """Check if command contains dangerous patterns."""
    warnings = []
    for pattern, message in DANGEROUS_PATTERNS:
        if re.search(pattern, command, re.IGNORECASE):
            warnings.append(message)
    return warnings

def main():
    try:
        input_data = json.load(sys.stdin)
        command = input_data.get('tool_input', {}).get('command', '')

        if command:
            warnings = check_command_safety(command)
            if warnings:
                for warning in warnings:
                    print(f"⚠️  {warning}", file=sys.stderr)
                print(f"Command: {command}", file=sys.stderr)
                print("Please confirm this is intentional.", file=sys.stderr)
                sys.exit(2)  # Block the command
    except:
        pass

if __name__ == "__main__":
    main()
