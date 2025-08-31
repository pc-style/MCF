---
description: Use specialized subagents with isolated context to handle complex tasks efficiently
argument-hint: <task-description>
allowed-tools: Task, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
---

## Efficient Multi-Agent Task Execution

**Task:** $ARGUMENTS

## Gemini MCP Capabilities

**Use Gemini MCP tools for advanced task coordination and creative problem solving:**

- **ask-gemini**: Advanced task decomposition analysis, complex agent coordination evaluation, and structured workflow recommendations with changeMode
- **brainstorm**: Creative multi-agent solutions, innovative task distribution strategies, and optimization techniques
- Perfect for analyzing complex task requirements, generating efficient agent workflows, and exploring creative coordination patterns
- Use changeMode parameter with ask-gemini for structured task orchestration and agent coordination suggestions
- These tools can save context usage by handling complex workflow analysis and coordination decisions efficiently

Execute this task using specialized subagents with minimal context to optimize token usage:

### 1. **Task Analysis & Agent Selection**

First, analyze the task and determine which specialized agents to use:

**Available Context-Efficient Agents:**

- **docs-researcher**: Documentation research using Context7 MCP
- **security-auditor**: Security analysis and vulnerability assessment
- **perf-optimizer**: Performance analysis and optimization
- **devops-engineer**: Infrastructure, deployment, and automation
- **api-designer**: API design and implementation
- **micro-researcher**: Lightweight quick research tasks
- **micro-analyzer**: Focused code analysis with minimal context
- **micro-security**: Targeted security checks
- **micro-executor**: Implementation tasks with isolated context

### 2. **Parallel Execution Strategy**

Identify subtasks that can be executed in parallel:

- **Independent Research**: Documentation, best practices, examples
- **Analysis Tasks**: Code review, security audit, performance check
- **Implementation Work**: Code changes, configuration updates
- **Validation Tasks**: Testing, verification, quality checks

### 3. **Context Isolation Plan**

For each subagent:

- Provide **only essential context** for their specific task
- Avoid duplicating information across agents
- Use **minimal background** necessary for task completion
- Focus on **specific deliverables** from each agent

### 4. **Execution Protocol**

Execute the following workflow:

1. **Decompose** the main task into agent-specific subtasks
2. **Assign** each subtask to the most appropriate specialist agent
3. **Execute** agents in parallel where possible (using Task tool)
4. **Collect** results from each agent efficiently
5. **Synthesize** results into unified solution
6. **Integrate** all outcomes into actionable next steps

### 5. **Efficiency Metrics**

Target outcomes:

- **50-70% token savings** through context isolation
- **Parallel execution** for faster completion
- **Specialized expertise** for higher quality results
- **Focused outputs** with clear actionable items

### 6. **Result Integration**

Combine agent outputs to provide:

- **Unified solution** addressing all aspects of the task
- **Priority-ordered action items** from all agents
- **Cross-agent insights** and recommendations
- **Next steps** with clear ownership and timeline

**Execute this multi-agent approach now.**
