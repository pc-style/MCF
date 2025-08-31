---
description: Use multiple specialized agents to create a new feature from concept to implementation
argument-hint: <feature-description>
allowed-tools: Task, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
---

## Multi-Agent Feature Development

**Feature Request:** $ARGUMENTS

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.

Execute comprehensive feature development using specialized agents for maximum efficiency and quality:

### 1. **Feature Analysis & Planning**

**Agent Team:**

- **micro-analyzer**: Analyze current codebase structure and identify integration points
- **api-designer**: Design API endpoints and data models if needed
- **semantic-navigator**: Understand existing architecture and dependencies

**Objectives:**

- Break down the feature request into technical requirements
- Identify existing code that can be leveraged or modified
- Design the feature architecture and integration points
- Create implementation roadmap and task breakdown

### 2. **API & Data Design**

**Agent Team:**

- **api-designer**: RESTful API endpoint design and specification
- **micro-analyzer**: Database schema analysis and modifications

**Deliverables:**

- API endpoint specifications (OpenAPI/Swagger)
- Database schema changes and migration scripts
- Data validation and serialization logic
- API documentation and usage examples

### 3. **Feature Implementation**

**Agent Team:**

- **backend-developer**: Server-side logic and business rules
- **frontend-developer**: User interface and user experience
- **fullstack-developer**: Full-stack integration and coordination

**Implementation Phases:**

- **Backend Development**: Core business logic, API endpoints, data models
- **Frontend Development**: UI components, state management, API integration
- **Integration**: End-to-end feature integration and validation

### 4. **Testing & Quality Assurance**

**Agent Team:**

- **test-engineer**: Comprehensive testing strategy and implementation
- **security-auditor**: Security validation and vulnerability assessment

**Testing Strategy:**

- **Unit Tests**: Component-level testing and validation
- **Integration Tests**: API and database integration testing
- **End-to-End Tests**: Complete user workflow validation
- **Security Tests**: Authentication, authorization, and data protection
- **Performance Tests**: Load testing and optimization

### 5. **Performance Optimization**

**Agent Team:**

- **perf-optimizer**: Performance analysis and optimization
- **micro-analyzer**: Code efficiency and bottleneck identification

**Optimization Areas:**

- **Database Performance**: Query optimization and indexing
- **API Performance**: Response time and throughput optimization
- **Frontend Performance**: Load time and user experience optimization
- **Caching Strategy**: Implementation of appropriate caching mechanisms

### 6. **Documentation & Deployment**

**Agent Team:**

- **docs-researcher**: Comprehensive feature documentation
- **devops-engineer**: Deployment strategy and CI/CD integration

**Documentation:**

- **Technical Documentation**: Implementation details and architecture
- **User Documentation**: Feature usage and configuration guides
- **API Documentation**: Endpoint specifications and examples
- **Deployment Documentation**: Setup and configuration procedures

### 7. **Feature Deployment**

**Agent Team:**

- **devops-engineer**: Deployment pipeline and infrastructure
- **test-engineer**: Deployment validation and monitoring

**Deployment Strategy:**

- **Staging Deployment**: Feature testing in production-like environment
- **Progressive Rollout**: Gradual feature activation and monitoring
- **Monitoring Setup**: Performance and error tracking implementation
- **Rollback Strategy**: Quick rollback procedures if issues arise

### **Common Feature Development Patterns**

**User Authentication Feature:**

```
1. api-designer → security-auditor → backend-developer
2. frontend-developer → test-engineer → devops-engineer
```

**Data Dashboard Feature:**

```
1. micro-analyzer → api-designer → backend-developer
2. frontend-developer → perf-optimizer → test-engineer
```

**Integration Feature:**

```
1. semantic-navigator → api-designer → backend-developer
2. test-engineer → security-auditor → devops-engineer
```

### **Success Metrics**

- **Development Velocity**: Time from concept to production
- **Code Quality**: Test coverage, maintainability, performance
- **User Experience**: Feature usability and adoption rates
- **System Impact**: Performance impact and system stability
- **Security Posture**: Vulnerability assessment and compliance

**Execute comprehensive multi-agent feature development now.**
