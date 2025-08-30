---
name: orchestration:team
description: Activate the full orchestrator team for complex development tasks
parameters:
  workflow: The workflow template to use (new-feature-development, system-modernization, performance-optimization)
  scope: The scope of work (small, medium, large, enterprise)
  priority: Priority level (low, medium, high, critical)
---

# Orchestration Team Activation

This command activates the complete MCF Orchestrator Team for complex development initiatives.

## Usage

```
/orchestration:team workflow=new-feature-development scope=medium priority=high
```

## Available Workflows

### **new-feature-development**
Complete workflow for developing new features from concept to production
- **Duration**: 2-4 weeks
- **Team**: Orchestrator, System Architect, API Architect, Backend Developer, Frontend Developer, Test Engineer
- **Best for**: New feature development, product enhancements, user story implementation

### **system-modernization**
Comprehensive workflow for modernizing legacy systems
- **Duration**: 4-8 weeks
- **Team**: Full orchestrator team including DevOps Engineer
- **Best for**: Legacy system modernization, technology stack updates, architectural refactoring

### **performance-optimization**
Systematic approach to performance bottleneck identification and resolution
- **Duration**: 2-3 weeks
- **Team**: Orchestrator, Performance Engineer, Backend Developer, Frontend Developer, DevOps Engineer
- **Best for**: Performance issues, scalability problems, user experience optimization

## Parameters

### **workflow** (required)
- `new-feature-development`: Complete feature development lifecycle
- `system-modernization`: Legacy system modernization
- `performance-optimization`: Performance bottleneck resolution

### **scope** (optional, default: medium)
- `small`: < 2 weeks, 1-2 developers
- `medium`: 2-4 weeks, 2-4 developers
- `large`: 1-2 months, 4-8 developers
- `enterprise`: 2-6 months, 8+ developers

### **priority** (optional, default: medium)
- `low`: Standard timeline, normal resources
- `medium`: Accelerated timeline, additional resources
- `high`: Fast-tracked with dedicated resources
- `critical`: Emergency response with maximum resources

## What Happens When You Run This Command

1. **Orchestrator Activation**: The main orchestrator agent is activated
2. **Team Assembly**: Required specialized agents are identified and prepared
3. **Workflow Loading**: The specified workflow template is loaded and customized
4. **Kickoff Meeting**: Virtual kickoff with all team members to align on objectives
5. **Task Distribution**: Initial tasks are assigned to appropriate team members
6. **Progress Tracking**: Automated progress monitoring and reporting begins

## Example Commands

### New Feature Development
```
/orchestration:team workflow=new-feature-development scope=medium priority=high
```

### System Modernization
```
/orchestration:team workflow=system-modernization scope=large priority=critical
```

### Performance Optimization
```
/orchestration:team workflow=performance-optimization scope=small priority=medium
```

## Monitoring Progress

Once activated, you can monitor progress using:

```
/orchestration:status
```

This provides real-time updates on:
- Current phase and progress
- Active tasks and assignments
- Completed deliverables
- Upcoming milestones
- Any blockers or issues

## Communication

The orchestrator team will:
- Provide daily progress updates
- Send milestone completion notifications
- Alert you to any blockers requiring your input
- Deliver final results with comprehensive documentation

## Cost and Resource Considerations

- **Small scope**: ~20-40 agent hours
- **Medium scope**: ~40-80 agent hours
- **Large scope**: ~80-160 agent hours
- **Enterprise scope**: ~160-400+ agent hours

Resources are allocated based on priority level and automatically scaled as needed.

## Emergency Stop

If you need to pause or stop the orchestration:

```
/orchestration:stop reason="Project requirements changed"
```

This gracefully stops all activities and provides a final status report.
