---
name: orchestration:status
description: Check the status of active orchestration workflows
parameters:
  workflow: Specific workflow to check (optional - shows all if not specified)
  detailed: Show detailed status including task breakdown (true/false, default: false)
---

# Orchestration Status

Monitor the progress and status of active orchestration workflows.

## Usage

```
/orchestration:status
/orchestration:status workflow=new-feature-development
/orchestration:status detailed=true
```

## Status Information

When you run this command, you'll receive:

### **Workflow Overview**
- Current phase and overall progress
- Start date and estimated completion
- Team members and their current status
- Key milestones and deadlines

### **Active Tasks** (when detailed=true)
- Task name and description
- Assigned agent and current status
- Start date and estimated completion
- Dependencies and blockers
- Progress percentage

### **Completed Work**
- Recently completed tasks
- Quality validation results
- Performance metrics achieved
- Documentation delivered

### **Upcoming Work**
- Next phase planning
- Resource requirements
- Risk mitigation activities
- Stakeholder deliverables

### **Issues & Blockers**
- Current problems requiring attention
- Resolution plans and timelines
- Impact assessment
- Escalation requirements

## Status Indicators

### **Workflow Status**
- 🔄 **Active**: Workflow is currently executing
- ⏸️ **Paused**: Workflow temporarily stopped
- ✅ **Completed**: Workflow finished successfully
- ❌ **Failed**: Workflow encountered critical issues
- 🟡 **At Risk**: Workflow behind schedule or over budget

### **Task Status**
- 🎯 **Assigned**: Task assigned but not started
- 🚀 **In Progress**: Task currently being worked on
- ✅ **Completed**: Task finished successfully
- ❌ **Failed**: Task encountered errors
- ⏳ **Blocked**: Task waiting on dependencies
- 🔄 **Review**: Task completed, awaiting review

### **Quality Gates**
- ✅ **Passed**: Quality requirements met
- ❌ **Failed**: Quality requirements not met
- 🔄 **In Review**: Quality validation in progress
- ⏭️ **Pending**: Quality validation not yet started

## Example Output

```
🎭 MCF ORCHESTRATOR STATUS REPORT
══════════════════════════════════════════════

📊 WORKFLOW OVERVIEW
├── Workflow: new-feature-development
├── Status: 🔄 Active (Phase 2/5)
├── Progress: ████████░░░░ 65%
├── Started: 2023-12-01
├── ETA: 2023-12-15
└── Team: 6 active agents

👥 TEAM STATUS
├── 🎯 system-architect: ✅ Phase 1 complete
├── 🚀 api-architect: In progress (85%)
├── ⏳ backend-developer: Blocked (waiting for API spec)
├── 🎯 frontend-developer: Assigned
├── 🎯 test-engineer: Pending
└── 🎯 devops-engineer: Pending

📋 CURRENT PHASE: Development & Integration
├── Duration: Days 6-10
├── Progress: ████░░░░░░ 40%
├── Active Tasks: 3
└── Completed Tasks: 7

⚠️ BLOCKERS & ISSUES
└── API specification review pending
    ├── Impact: Medium
    ├── Resolution: Stakeholder review scheduled
    └── ETA: 2023-12-03

🎯 UPCOMING MILESTONES
├── Phase 3 completion: 2023-12-08
├── Quality gate: 2023-12-10
├── Deployment: 2023-12-12
└── Launch: 2023-12-15

💰 RESOURCE UTILIZATION
├── Budget Used: $2,340 / $5,000 (47%)
├── Agent Hours: 156 / 320 (49%)
└── Risk Level: 🟢 Low
```

## Alert Types

### **Immediate Attention Required**
- 🚨 Critical blockers stopping progress
- ❌ Failed quality gates
- ⏰ Missed deadlines
- 💰 Budget overruns

### **Monitor Closely**
- 🟡 Tasks at risk of delay
- ⚠️ Quality concerns identified
- 📉 Performance below targets
- 🔄 Dependencies not met

### **Information Only**
- ✅ Task completions
- 📊 Progress updates
- 📋 Documentation delivered
- 🎯 Milestone achievements

## Automated Reporting

The status command provides:

### **Daily Summary** (9 AM)
- Previous day accomplishments
- Current day priorities
- Any new issues or blockers

### **Phase Completion Alerts**
- Automatic notification when phases complete
- Quality gate results
- Go/no-go recommendations

### **Milestone Notifications**
- Key deliverable completions
- Stakeholder review requirements
- Deployment readiness confirmations

## Filtering Options

### **By Workflow**
```
/orchestration:status workflow=new-feature-development
```
Shows status for specific workflow only

### **Detailed View**
```
/orchestration:status detailed=true
```
Includes task-level breakdown and technical details

### **Team Focus**
```
/orchestration:status team=backend-developer
```
Shows tasks and status for specific team member

## Integration with MCF

This command integrates with:

- **MCF Dashboard**: Visual progress tracking
- **Notification System**: Email/Slack alerts for important events
- **Documentation System**: Automatic generation of status reports
- **Resource Management**: Budget and time tracking
- **Quality Management**: Automated quality gate reporting

## Troubleshooting

### **No Active Workflows**
If you see "No active orchestration workflows":
- Check if workflows were properly started with `/orchestration:team`
- Verify workflow parameters were correct
- Confirm all required agents are available

### **Missing Status Information**
If status information seems incomplete:
- Wait a few minutes for status to update
- Check agent connectivity
- Run `/orchestration:status detailed=true` for more information

### **Performance Issues**
If status updates are slow:
- Status is generated in real-time
- Large workflows may take longer to process
- Consider using summary view for faster results
