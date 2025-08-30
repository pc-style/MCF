---
description: Restore conversation context from a previously saved bookmark
argument-hint: <bookmark-name>
allowed-tools: Read
---

## Context Restoration from Bookmark

**Bookmark to Load:** $ARGUMENTS

Restore and integrate context from the specified bookmark:

### 1. **Load Bookmark Data**
Read the bookmark file from `.claude/bookmarks/$ARGUMENTS.md` and extract:
- Session state and objectives at time of bookmark
- Technical configurations and settings
- Code implementations and architecture decisions
- Project knowledge and lessons learned
- Development environment state

### 2. **Current State Analysis**
Compare bookmark state with current reality:
- Files that have changed since bookmark
- Configurations that may have been updated
- New features or changes implemented
- Dependencies or tools that may have changed

### 3. **Context Integration**
Integrate the bookmarked context into our current session:
- Restore key project objectives and priorities
- Bring back important technical decisions and reasoning
- Re-establish development patterns and approaches
- Integrate lessons learned and best practices

### 4. **State Verification**
Verify that the restored context is still valid:
- Check that referenced files and configurations exist
- Validate that saved commands and tools still work
- Confirm that project structure matches expectations
- Test that integrations and dependencies are functional

### 5. **Gap Analysis**
Identify what has changed since the bookmark:
- New features or functionality added
- Configuration changes or updates made
- Problems solved or new issues discovered
- Technical debt addressed or accumulated

### 6. **Context Synchronization**
Update the restored context with current information:
- Merge recent changes with bookmarked state
- Resolve any conflicts between old and new information
- Update outdated information with current facts
- Integrate new knowledge with restored context

### 7. **Restoration Summary**
Provide a comprehensive summary of:
- What context was successfully restored
- What information needed to be updated
- Any inconsistencies or conflicts found
- Current project state after restoration

### 8. **Next Steps**
Based on the restored context, suggest:
- Immediate actions to take
- Areas that need attention or review
- Opportunities for continued development
- Context management strategies going forward

**Execute context restoration now.**