#!/usr/bin/env python3
"""
Context7 Documentation Reminder Hook
Analyzes user prompts and suggests using Context7 MCP for documentation when external frameworks/libraries are detected.
"""
import json
import sys
import re

# Common frameworks, libraries, and technologies that benefit from Context7
TECH_PATTERNS = [
    # JavaScript/Node.js
    r'\b(react|vue|angular|svelte|next\.?js|nuxt|express|fastify|nest\.?js)\b',
    r'\b(webpack|vite|rollup|parcel|babel|typescript|eslint|prettier)\b',
    r'\b(lodash|axios|moment|dayjs|socket\.io|prisma|mongoose)\b',

    # Python
    r'\b(django|flask|fastapi|pyramid|tornado|celery|sqlalchemy)\b',
    r'\b(pandas|numpy|matplotlib|seaborn|sklearn|tensorflow|pytorch)\b',
    r'\b(pytest|unittest|requests|beautifulsoup|scrapy)\b',

    # Databases
    r'\b(postgresql|mysql|mongodb|redis|elasticsearch|cassandra|supabase)\b',

    # Cloud/DevOps
    r'\b(aws|azure|gcp|docker|kubernetes|terraform|ansible)\b',
    r'\b(github\s+actions|gitlab\s+ci|jenkins|circleci)\b',

    # Other popular tools
    r'\b(git|graphql|apollo|prisma|strapi|sanity|contentful)\b',
    r'\b(tailwind|bootstrap|material\-ui|chakra|ant\s+design)\b',
]

# Keywords that suggest documentation needs
DOC_KEYWORDS = [
    'how to', 'tutorial', 'guide', 'documentation', 'docs', 'api reference',
    'getting started', 'setup', 'install', 'configure', 'examples',
    'best practices', 'patterns', 'implementation', 'integration'
]

def analyze_prompt(prompt):
    """Analyze prompt for technology mentions and documentation needs."""
    prompt_lower = prompt.lower()

    # Check for technology patterns
    detected_techs = []
    for pattern in TECH_PATTERNS:
        matches = re.findall(pattern, prompt_lower, re.IGNORECASE)
        detected_techs.extend(matches)

    # Check for documentation-related keywords
    has_doc_keywords = any(keyword in prompt_lower for keyword in DOC_KEYWORDS)

    # Check for question patterns that suggest needing docs
    question_patterns = [
        r'\bhow\s+(do|can|to)\b',
        r'\bwhat\s+is\b',
        r'\bwhere\s+(do|can|to)\b',
        r'\bwhy\s+(do|does)\b',
        r'\bwhich\s+way\b',
        r'\bcan\s+you\s+(help|show|explain)\b',
        r'\bneed\s+(help|to\s+understand|examples?)\b',
    ]

    has_question_pattern = any(re.search(pattern, prompt_lower) for pattern in question_patterns)

    return {
        'detected_techs': list(set(detected_techs)),
        'has_doc_keywords': has_doc_keywords,
        'has_question_pattern': has_question_pattern,
        'should_suggest_context7': bool(detected_techs) and (has_doc_keywords or has_question_pattern)
    }

def main():
    try:
        # Read JSON input from stdin
        input_data = json.load(sys.stdin)

        prompt = input_data.get('prompt', '')
        if not prompt.strip():
            sys.exit(0)

        analysis = analyze_prompt(prompt)

        if analysis['should_suggest_context7']:
            techs_str = ', '.join(analysis['detected_techs'][:3])  # Limit to first 3
            if len(analysis['detected_techs']) > 3:
                techs_str += f" (and {len(analysis['detected_techs']) - 3} more)"

            context_message = f"""💡 **Documentation Suggestion**: I detected you're asking about {techs_str}.

Consider using Context7 MCP for the most up-to-date documentation:
- Use Context7 to resolve library IDs: "Can you use Context7 to get the latest docs for {analysis['detected_techs'][0]}?"
- Get specific documentation: "Use Context7 to fetch documentation about [specific topic] in {analysis['detected_techs'][0]}"

This will give you the most current and comprehensive information available."""

            # Output the additional context
            print(context_message)

    except Exception as e:
        # Silent fail - don't interrupt the workflow
        pass

if __name__ == "__main__":
    main()
