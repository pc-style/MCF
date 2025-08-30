---
description: Display comprehensive guide and current status for all context management features
---

## 🧠 Context Management System Overview

### **Available Context Commands**

| Command | Purpose | Example Usage |
|---------|---------|---------------|
| `/context:purge [areas]` | Smart context cleanup | `/context:purge security,api` |
| `/context:agent <task>` | Use specialized subagents | `/context:agent "optimize this API"` |
| `/context:state` | Current usage analysis | `/context:state` |
| `/context:bookmark <name>` | Save conversation state | `/context:bookmark stable-v1` |
| `/context:load <name>` | Restore saved state | `/context:load stable-v1` |
| `/context:split [topics]` | Split by topics | `/context:split frontend,backend` |
| `/context:merge <contexts>` | Merge conversations | `/context:merge api,frontend` |

### **Specialized Micro-Agents for Efficiency**

| Agent | Context Usage | Best For |
|-------|---------------|----------|
| `micro-researcher` | Minimal | Quick documentation lookup |
| `micro-analyzer` | Low | Targeted code review |
| `micro-executor` | Medium | Focused implementation tasks |
| `docs-researcher` | Context7 focused | Deep documentation research |
| `security-auditor` | Security focused | Comprehensive security analysis |
| `perf-optimizer` | Performance focused | Speed and optimization work |

### **Context Health Indicators**

**Current Monitoring:**
- 🟢 **Excellent**: <200 tokens/message average
- 🟡 **Good**: 200-400 tokens/message  
- 🟠 **Fair**: 400-600 tokens/message
- 🔴 **Poor**: >600 tokens/message

### **Smart Usage Patterns**

**For Research Tasks:**
```
> Use micro-researcher to find the latest Next.js authentication patterns
```

**For Complex Multi-Step Tasks:**
```
> /context:agent "analyze this API for security issues and performance bottlenecks, then suggest improvements"
```

**For Context Cleanup:**
```
> /context:purge authentication,database
```

**For Strategic Bookmarking:**
```
> /context:bookmark pre-deployment-state
```

### **Automatic Context Monitoring**

Your setup includes intelligent monitoring that will:
- Track token usage automatically
- Suggest optimization at key thresholds
- Recommend context management actions
- Monitor conversation efficiency

**Thresholds:**
- 8k+ tokens: Suggests analysis with `/context:state`
- 12k+ tokens: Recommends `/context:purge` 
- 20k+ tokens: Urgent cleanup needed
- 40+ messages: Consider topic splitting
- 60+ messages: Strongly recommend organization

### **Best Practices Summary**

1. **Proactive Management**: Use `/context:state` every 20-30 messages
2. **Strategic Bookmarking**: Save state before major changes
3. **Agent Efficiency**: Use micro-agents for focused tasks
4. **Topic Organization**: Split conversations that cover multiple areas
5. **Regular Cleanup**: Purge context when approaching token limits

### **Token Savings Potential**

- **Context Purging**: 50-70% reduction in conversation size
- **Micro-Agent Usage**: 30-50% savings on complex tasks
- **Strategic Splitting**: Prevents context bloat in focused discussions
- **Smart Bookmarking**: Enables efficient context restoration

**Your context management system is now active and monitoring!** 🎯