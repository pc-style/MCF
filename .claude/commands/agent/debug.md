---
description: Use specialized agents to debug complex issues with comprehensive analysis
argument-hint: <issue-description or error-pattern>
allowed-tools: Task, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
---

## Multi-Agent Debugging System

**Debug Target:** $ARGUMENTS

## Gemini MCP Capabilities

**Use Gemini MCP tools for advanced debugging analysis and creative problem solving:**

- **ask-gemini**: Advanced debugging strategy analysis, complex issue evaluation, and structured debugging recommendations with changeMode
- **brainstorm**: Creative debugging solutions, innovative troubleshooting strategies, and root cause analysis techniques
- Perfect for analyzing complex debugging requirements, generating efficient troubleshooting strategies, and exploring creative debugging approaches
- Use changeMode parameter with ask-gemini for structured debugging implementation and issue resolution suggestions
- These tools can save context usage by handling complex debugging analysis and problem-solving decisions efficiently

Execute comprehensive debugging using specialized agents for rapid issue resolution:

### 1. **Debug Team Assembly**

**Primary Debug Agents:**

- **micro-analyzer**: Code analysis and issue pattern identification
- **semantic-navigator**: Codebase exploration and context understanding
- **perf-optimizer**: Performance issue analysis and bottleneck detection

**Specialized Debug Agents:**

- **security-auditor**: Security vulnerability and exploit analysis
- **api-designer**: API and interface issue diagnosis
- **devops-engineer**: Infrastructure and deployment problem analysis

### 2. **Debug Methodology**

**Phase 1: Issue Analysis & Reproduction**

- Error pattern identification and classification
- Issue reproduction steps and environment analysis
- Impact assessment and severity classification
- Context gathering and related code analysis

**Phase 2: Root Cause Investigation**

- Code path analysis using Serena semantic tools
- Dependency and data flow tracing
- Performance profiling and bottleneck identification
- Security vulnerability assessment
- Configuration and environment validation

**Phase 3: Solution Development & Testing**

- Fix implementation and validation
- Regression testing and impact assessment
- Performance validation and optimization
- Security review and vulnerability closure
- Documentation and prevention measures

### 3. **Debug Focus Areas**

**Code-Level Issues:**

- Logic errors and algorithm problems
- Data type and validation issues
- Memory leaks and resource management
- Exception handling and error propagation
- Race conditions and concurrency issues

**System-Level Issues:**

- Performance bottlenecks and scalability problems
- Integration and API communication issues
- Database and data persistence problems
- Security vulnerabilities and access control issues
- Infrastructure and deployment problems

**User Experience Issues:**

- Interface bugs and usability problems
- Workflow and business logic issues
- Data consistency and validation problems
- Accessibility and internationalization issues

### 4. **Debug Tools & Techniques**

**Serena Semantic Analysis:**

- find_symbol for locating relevant code elements
- get_symbol_info for detailed function/class analysis
- find_referencing_symbols for usage pattern analysis
- get_project_structure for architectural context

**Traditional Debugging:**

- Log analysis and error pattern recognition
- Code review and static analysis
- Performance profiling and benchmarking
- Security scanning and vulnerability assessment
- Integration testing and validation

### 5. **Debug Outputs**

- **Issue Summary**: Clear problem description and impact
- **Root Cause Analysis**: Detailed technical explanation
- **Solution Implementation**: Working fix with explanation
- **Testing Results**: Validation and regression testing
- **Prevention Measures**: How to avoid similar issues
- **Documentation Updates**: Code comments and guides

### 6. **Debug Workflow**

1. **Issue Triage**: Classify and prioritize the problem
2. **Context Gathering**: Collect relevant code and environment info
3. **Pattern Analysis**: Identify common debugging patterns
4. **Solution Development**: Implement and test the fix
5. **Validation**: Ensure the fix resolves the issue completely
6. **Prevention**: Document and implement preventive measures

### 7. **Advanced Debug Scenarios**

**Performance Debugging:**

- perf-optimizer + micro-analyzer for bottleneck analysis
- Focus: CPU, memory, I/O, and network optimization

**Security Debugging:**

- security-auditor + micro-analyzer for vulnerability analysis
- Focus: Input validation, access control, data exposure

**Integration Debugging:**

- api-designer + devops-engineer for system integration issues
- Focus: API contracts, data flow, deployment problems

**Data Debugging:**

- micro-analyzer + semantic-navigator for data flow analysis
- Focus: Data validation, transformation, and persistence

### 8. **Debug Best Practices**

- **Systematic Approach**: Follow structured debugging methodology
- **Evidence-Based**: Base conclusions on data and analysis
- **Minimal Changes**: Make smallest possible changes to fix issues
- **Comprehensive Testing**: Validate fixes don't introduce new problems
- **Documentation**: Record solutions for future reference
- **Prevention**: Implement measures to avoid similar issues

**Execute multi-agent debugging session now.**
