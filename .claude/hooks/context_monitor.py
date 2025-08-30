#!/usr/bin/env python3
"""
Context Usage Monitor Hook
Tracks context usage and provides intelligent suggestions for optimization.
"""
import json
import sys
import os
from pathlib import Path

def estimate_context_usage(transcript_path):
    """Estimate current context usage from transcript file."""
    try:
        if not os.path.exists(transcript_path):
            return {"estimated_tokens": 0, "message_count": 0, "chars": 0}
        
        with open(transcript_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Count messages (rough estimate)
        message_count = content.count('"role":')
        
        # Estimate tokens (rough: 4 chars per token average)
        total_chars = len(content)
        estimated_tokens = total_chars // 4
        
        return {
            "estimated_tokens": estimated_tokens,
            "message_count": message_count,
            "chars": total_chars
        }
    except Exception as e:
        return {"estimated_tokens": 0, "message_count": 0, "chars": 0}

def analyze_context_health(usage_stats):
    """Analyze context health and provide recommendations."""
    tokens = usage_stats.get("estimated_tokens", 0)
    messages = usage_stats.get("message_count", 0)
    
    recommendations = []
    
    # Token-based recommendations
    if tokens > 20000:
        recommendations.append("🔥 High token usage detected - consider `/context:purge` immediately")
    elif tokens > 12000:
        recommendations.append("⚠️  Growing context - `/context:purge` recommended soon")
    elif tokens > 8000:
        recommendations.append("💡 Consider `/context:state` to analyze usage patterns")
    
    # Message-based recommendations  
    if messages > 60:
        recommendations.append("📚 Long conversation - consider `/context:split` for topic organization")
    elif messages > 40:
        recommendations.append("🎯 Large conversation - `/context:agent` for complex tasks recommended")
    
    # Combined recommendations
    if tokens > 15000 and messages > 50:
        recommendations.append("🧠 Complex conversation - create `/context:bookmark` before major changes")
    
    return recommendations

def get_context_efficiency_score(usage_stats):
    """Calculate context efficiency score."""
    tokens = usage_stats.get("estimated_tokens", 1)
    messages = usage_stats.get("message_count", 1)
    
    # Higher score = more efficient (less tokens per message)
    avg_tokens_per_message = tokens / max(messages, 1)
    
    if avg_tokens_per_message < 200:
        return "🟢 Excellent"
    elif avg_tokens_per_message < 400:
        return "🟡 Good"
    elif avg_tokens_per_message < 600:
        return "🟠 Fair"
    else:
        return "🔴 Poor"

def main():
    try:
        input_data = json.load(sys.stdin)
        transcript_path = input_data.get('transcript_path', '')
        hook_event = input_data.get('hook_event_name', '')
        
        # Only run monitoring on Stop events to avoid spam
        if hook_event != 'Stop' or not transcript_path:
            sys.exit(0)
        
        usage_stats = estimate_context_usage(transcript_path)
        recommendations = analyze_context_health(usage_stats)
        efficiency_score = get_context_efficiency_score(usage_stats)
        
        # Only show output if there are recommendations
        if recommendations or usage_stats['estimated_tokens'] > 8000:
            print(f"🧠 **Context Monitor** ({efficiency_score} efficiency)")
            print(f"📊 ~{usage_stats['estimated_tokens']:,} tokens, {usage_stats['message_count']} messages")
            
            if recommendations:
                print("💡 **Recommendations:**")
                for rec in recommendations[:2]:  # Limit to 2 recommendations
                    print(f"  • {rec}")
    
    except Exception:
        # Silent fail to not interrupt workflow
        pass

if __name__ == "__main__":
    main()