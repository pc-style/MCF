#!/usr/bin/env python3
"""
Agent Command Suggestions Hook
Analyzes user prompts and suggests using agent commands for complex development tasks
"""
import json
import sys
import re

# Development task patterns that benefit from agent commands
AGENT_PATTERNS = [
    # Feature development
    (r'\b(create|build|implement|develop|add)\s+(a\s+)?(new\s+)?(feature|functionality|system|component)\b',
     'Use /agent:feature for comprehensive feature development with multiple specialized agents'),
    
    # System redesign
    (r'\b(redesign|refactor|restructure|reorganize|modernize)\s+(system|architecture|component|module)\b',
     'Use /agent:redesign for system redesign with architectural analysis and modern patterns'),
    
    # Debugging complex issues
    (r'\b(debug|fix|resolve|troubleshoot)\s+(complex|difficult|persistent|intermittent)\s+(issue|problem|bug|error)\b',
     'Use /agent:debug for comprehensive debugging with multiple specialized agents'),
    
    # Performance optimization
    (r'\b(optimize|improve|enhance|boost)\s+(performance|speed|efficiency|scalability)\b',
     'Use /agent:optimize for performance analysis and optimization with specialized agents'),
    
    # Security and compliance
    (r'\b(audit|review|check|validate)\s+(security|compliance|vulnerabilities|standards)\b',
     'Use /agent:audit for comprehensive security and compliance auditing'),
    
    # Testing strategies
    (r'\b(create|design|implement|set\s+up)\s+(testing|test\s+strategy|test\s+suite|test\s+coverage)\b',
     'Use /agent:test for comprehensive testing strategy and implementation'),
    
    # Deployment and CI/CD
    (r'\b(set\s+up|create|implement|configure)\s+(deployment|ci/cd|pipeline|infrastructure)\b',
     'Use /agent:deploy for deployment pipeline design and CI/CD implementation'),
    
    # Brainstorming and ideation
    (r'\b(brainstorm|generate\s+ideas|explore\s+options|find\s+solutions)\b',
     'Use /agent:brainstorm for collaborative ideation with multiple specialized agents'),
    
    # Agent generation
    (r'\b(create|generate|build)\s+(an?\s+)?(agent|specialist|expert)\b',
     'Use /agent:generate-ag to create project-specific specialized agents'),
    
    # Fully automatic development
    (r'\b(build|create|implement|develop)\s+(a\s+)?(complete|full|entire|whole)\s+(system|application|platform|solution)\b',
     'Use /agent:auto for fully automatic development with a complete team of specialized agents'),
    
    # Complex project development
    (r'\b(project|system|application|platform)\s+(development|creation|implementation|building)\b',
     'Use /agent:auto for complete project development with automated team coordination'),
]

# Keywords that suggest complex development tasks
COMPLEX_TASK_KEYWORDS = [
    'architecture', 'system design', 'scalability', 'performance', 'security',
    'compliance', 'testing strategy', 'deployment', 'ci/cd', 'infrastructure',
    'optimization', 'refactoring', 'redesign', 'modernization', 'migration'
]

def analyze_for_agent_suggestions(prompt):
    """Analyze prompt for opportunities to use agent commands."""
    suggestions = []
    prompt_lower = prompt.lower()
    
    # Check for agent pattern matches
    for pattern, suggestion in AGENT_PATTERNS:
        if re.search(pattern, prompt_lower):
            suggestions.append(suggestion)
    
    # Check for complex task keywords
    has_complex_keywords = any(keyword in prompt_lower for keyword in COMPLEX_TASK_KEYWORDS)
    
    # Check for question patterns that suggest complex tasks
    question_patterns = [
        r'\bhow\s+(do|can|should)\s+(we|I|you)\s+(design|implement|optimize|secure)\b',
        r'\bwhat\s+(is|are)\s+(the\s+)?(best\s+)?(way|approach|strategy|method)\b',
        r'\bcan\s+(you|we)\s+(help|assist|guide)\s+(with|on)\s+(complex|difficult|challenging)\b',
        r'\bneed\s+(help|guidance|assistance)\s+(with|on)\s+(architecture|design|optimization)\b',
    ]
    
    has_complex_question = any(re.search(pattern, prompt_lower) for pattern in question_patterns)
    
    # Check for large scope indicators
    scope_indicators = [
        r'\b(entire|whole|complete|full)\s+(system|application|codebase|project)\b',
        r'\b(major|significant|comprehensive|extensive)\s+(refactor|redesign|overhaul)\b',
        r'\b(multiple|several|various)\s+(components|modules|systems|areas)\b',
    ]
    
    has_large_scope = any(re.search(pattern, prompt_lower) for pattern in scope_indicators)
    
    # Suggest agent commands for complex tasks
    if has_complex_keywords or has_complex_question or has_large_scope:
        if not suggestions:  # Only add if no specific pattern matched
            suggestions.append('Consider using agent commands for complex development tasks that benefit from multiple specialized perspectives')
    
    return list(set(suggestions))  # Remove duplicates

def main():
    try:
        input_data = json.load(sys.stdin)
        prompt = input_data.get('prompt', '')
        
        if not prompt.strip():
            sys.exit(0)
        
        suggestions = analyze_for_agent_suggestions(prompt)
        
        if suggestions:
            context_message = "\n🚀 **Agent Command Suggestions**:\n"
            for i, suggestion in enumerate(suggestions[:2], 1):  # Limit to 2 suggestions
                context_message += f"{i}. {suggestion}\n"
            
            context_message += "\n💡 Agent commands provide specialized expertise and comprehensive analysis for complex development tasks."
            context_message += "\n📚 Use /agent:help to see all available agent commands."
            
            print(context_message)
    
    except Exception:
        # Silent fail - don't interrupt workflow
        pass

if __name__ == "__main__":
    main()
