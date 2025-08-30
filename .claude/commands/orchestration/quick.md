---
name: orchestration:quick
description: Quick orchestration for simple development tasks
parameters:
  type: Type of task (feature, bugfix, refactor, research, docs)
  description: Brief description of the task
  priority: Task priority (low, medium, high)
---

# Quick Orchestration

Fast-track orchestration for simple development tasks that don't require the full workflow.

## Usage

```
/orchestration:quick type=feature description="Add user authentication" priority=high
```

## Task Types

### **feature**
New feature development with basic testing
- **Duration**: 1-3 days
- **Team**: Backend Developer + Test Engineer
- **Deliverables**: Working feature + basic tests

### **bugfix**
Bug identification and resolution
- **Duration**: 4-12 hours
- **Team**: Backend/Frontend Developer + Test Engineer
- **Deliverables**: Fixed bug + regression tests

### **refactor**
Code refactoring and optimization
- **Duration**: 1-2 days
- **Team**: Backend/Frontend Developer + Test Engineer
- **Deliverables**: Refactored code + validation tests

### **research**
Technical research and feasibility analysis
- **Duration**: 4-8 hours
- **Team**: System Architect + Docs Researcher
- **Deliverables**: Research report + recommendations

### **docs**
Documentation creation and updates
- **Duration**: 2-6 hours
- **Team**: Docs Researcher + Semantic Navigator
- **Deliverables**: Updated documentation + examples

## Parameters

### **type** (required)
- `feature`: New feature implementation
- `bugfix`: Bug resolution
- `refactor`: Code improvement
- `research`: Technical investigation
- `docs`: Documentation work

### **description** (required)
Brief description of the task (max 200 characters)

### **priority** (optional, default: medium)
- `low`: Standard timeline
- `medium`: Within 24 hours
- `high`: Within 4 hours

## What Happens

1. **Automatic Team Selection**: Based on task type
2. **Quick Planning**: 15-minute planning session
3. **Parallel Execution**: Team works simultaneously
4. **Rapid Delivery**: Focused on speed and quality
5. **Immediate Results**: Progress updates every 2 hours

## Example Commands

### Feature Development
```
/orchestration:quick type=feature description="Implement password reset functionality" priority=high
```

### Bug Fix
```
/orchestration:quick type=bugfix description="Fix login form validation error" priority=medium
```

### Code Refactor
```
/orchestration:quick type=refactor description="Extract common database utilities" priority=low
```

### Technical Research
```
/orchestration:quick type=research description="Evaluate React vs Vue for new component" priority=medium
```

### Documentation
```
/orchestration:quick type=docs description="API documentation for user management" priority=low
```

## Team Composition by Type

### **Feature Tasks**
- **Lead**: Backend/Frontend Developer (based on feature type)
- **Support**: Test Engineer for validation
- **Review**: System Architect for design validation

### **Bug Fix Tasks**
- **Lead**: Backend/Frontend Developer (based on bug location)
- **Support**: Test Engineer for reproduction and validation
- **Analysis**: Perf Optimizer if performance-related

### **Refactor Tasks**
- **Lead**: Backend/Frontend Developer
- **Support**: Test Engineer for regression testing
- **Review**: Semantic Navigator for impact analysis

### **Research Tasks**
- **Lead**: System Architect for technical evaluation
- **Support**: Docs Researcher for information gathering
- **Analysis**: API Designer if API-related research

### **Documentation Tasks**
- **Lead**: Docs Researcher for content creation
- **Support**: Semantic Navigator for code examples
- **Review**: API Designer for API documentation

## Progress Tracking

### **Real-time Updates**
- Progress notifications every 30 minutes for high priority
- Hourly updates for medium priority
- Daily summaries for low priority

### **Status Indicators**
- 🚀 **Planning**: Team assembly and task breakdown
- 🔄 **Active**: Implementation in progress
- ✅ **Review**: Code review and testing
- 🎯 **Complete**: Task finished and delivered

## Quality Assurance

### **Automated Checks**
- Code linting and formatting
- Unit test execution
- Integration test validation
- Security scanning (for applicable tasks)

### **Review Process**
- Peer code review for all changes
- Functional testing validation
- Documentation review (for docs tasks)
- Performance impact assessment

## Success Criteria

### **Feature Tasks**
- ✅ Feature implemented and functional
- ✅ Basic unit tests written and passing
- ✅ No breaking changes to existing functionality
- ✅ Code follows project standards

### **Bug Fix Tasks**
- ✅ Bug reproduced and root cause identified
- ✅ Fix implemented and tested
- ✅ Regression tests added
- ✅ No new bugs introduced

### **Refactor Tasks**
- ✅ Code improved and more maintainable
- ✅ All existing tests still pass
- ✅ No functional changes to behavior
- ✅ Performance maintained or improved

### **Research Tasks**
- ✅ Research question answered thoroughly
- ✅ Multiple options evaluated with pros/cons
- ✅ Clear recommendation provided
- ✅ Implementation guidance included

### **Documentation Tasks**
- ✅ Documentation complete and accurate
- ✅ Examples provided and tested
- ✅ Consistent with existing documentation style
- ✅ Accessible to target audience

## Cost Estimation

### **By Priority**
- **High**: 4-12 agent hours
- **Medium**: 2-6 agent hours
- **Low**: 1-3 agent hours

### **By Type**
- **Feature**: 6-20 hours (based on complexity)
- **Bugfix**: 2-8 hours
- **Refactor**: 4-12 hours
- **Research**: 2-6 hours
- **Docs**: 1-4 hours

## When to Use Quick Orchestration

### **Perfect For**
- Small, well-defined tasks
- Urgent fixes and patches
- Proof-of-concept work
- Documentation updates
- Technical research questions

### **Not Ideal For**
- Complex multi-week projects
- System-wide architectural changes
- High-risk production deployments
- Tasks requiring extensive stakeholder involvement
- Projects with unclear requirements

## Integration

Quick orchestration integrates with:
- **MCF Agents**: Leverages existing specialized agents
- **Version Control**: Automatic commit and PR creation
- **CI/CD**: Triggers automated testing and deployment
- **Documentation**: Updates project documentation automatically
- **Notifications**: Progress updates via preferred channels
