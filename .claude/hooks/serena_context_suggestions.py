#!/usr/bin/env python3
"""
Serena Context Suggestions Hook
Suggests using Serena tools for code-related tasks to improve efficiency
"""
import json
import sys
import re

# Code-related patterns that benefit from Serena
CODE_PATTERNS = [
    # Function/method operations
    (r'\b(function|method|def|class)\b.*\b(find|search|locate|get)\b', 
     'Use find_symbol to locate functions/classes efficiently'),
    
    # Refactoring operations
    (r'\b(refactor|rename|move|extract|inline)\b.*\b(function|method|class|variable)\b',
     'Use Serena\'s symbol editing tools for precise refactoring'),
    
    # Code analysis
    (r'\b(analyze|review|inspect|examine)\b.*\b(code|function|class|method)\b',
     'Use get_symbol_info for detailed code analysis'),
    
    # Finding references/usage
    (r'\b(where|find|show).*\b(used|called|referenced|imported)\b',
     'Use find_referencing_symbols to trace usage efficiently'),
    
    # Code editing at symbol level
    (r'\b(add|insert|modify|update).*\b(before|after|in).*\b(function|method|class)\b',
     'Use insert_before_symbol/insert_after_symbol for precise editing'),
]

def analyze_for_serena_suggestions(prompt):
    """Analyze prompt for opportunities to use Serena tools."""
    suggestions = []
    prompt_lower = prompt.lower()
    
    for pattern, suggestion in CODE_PATTERNS:
        if re.search(pattern, prompt_lower):
            suggestions.append(suggestion)
    
    # Check for large file operations that could benefit from symbol-level work
    if re.search(r'\b(large|big|huge|entire)\s+(file|codebase|project)\b', prompt_lower):
        suggestions.append('Consider using Serena tools to work at symbol level instead of entire files')
    
    return list(set(suggestions))  # Remove duplicates

def main():
    try:
        input_data = json.load(sys.stdin)
        prompt = input_data.get('prompt', '')
        
        if not prompt.strip():
            sys.exit(0)
        
        suggestions = analyze_for_serena_suggestions(prompt)
        
        if suggestions:
            context_message = "\n🔍 **Serena Suggestions**:\n"
            for i, suggestion in enumerate(suggestions[:2], 1):  # Limit to 2 suggestions
                context_message += f"{i}. {suggestion}\n"
            
            context_message += "\n💡 Serena provides IDE-like semantic code tools."
            print(context_message)
    
    except Exception:
        # Silent fail - don't interrupt workflow
        pass

if __name__ == "__main__":
    main()
