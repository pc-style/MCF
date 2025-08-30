#!/usr/bin/env python3
"""
Serena Project Initialization Hook
Provides project info when Serena-compatible code is detected
"""
import json
import sys
import os

def check_serena_compatibility():
    """Check if project is compatible with Serena semantic analysis."""
    supported_extensions = {
        '.py', '.js', '.ts', '.tsx', '.jsx', '.php', '.go', '.rs', 
        '.cpp', '.c', '.hpp', '.h', '.zig', '.cs', '.rb', '.swift',
        '.kt', '.java', '.clj', '.dart', '.sh', '.bash', '.lua', 
        '.nix', '.ex', '.exs', '.erl'
    }
    
    code_files = []
    languages_found = set()
    
    for root, dirs, files in os.walk('.'):
        # Skip hidden directories and common ignore patterns
        dirs[:] = [d for d in dirs if not d.startswith('.') and d not in {'node_modules', '__pycache__', 'target', 'build', 'dist'}]
        
        for file in files:
            ext = os.path.splitext(file)[1].lower()
            if ext in supported_extensions:
                code_files.append(os.path.join(root, file))
                languages_found.add(ext)
                if len(code_files) >= 20:  # Stop after finding enough files
                    break
        
        if len(code_files) >= 20:
            break
    
    return len(code_files) > 0, len(code_files), languages_found

def main():
    try:
        input_data = json.load(sys.stdin)
        session_start_source = input_data.get('source', '')
        
        # Only run on session startup
        if session_start_source != 'startup':
            sys.exit(0)
        
        is_compatible, file_count, languages = check_serena_compatibility()
        
        if is_compatible and file_count >= 3:  # Only show for meaningful projects
            lang_list = ', '.join(sorted(languages)[:5])  # Show up to 5 extensions
            context_message = f"""🏗️ **Serena Ready**: Found {file_count}+ code files ({lang_list})

Available Serena tools:
• find_symbol - Locate functions/classes by name  
• get_symbol_info - Get detailed symbol information
• find_referencing_symbols - Find all symbol references
• get_project_structure - Project organization overview

💡 Use these for efficient, token-saving code operations."""
            
            print(context_message)
    
    except Exception:
        # Silent fail
        pass

if __name__ == "__main__":
    main()
