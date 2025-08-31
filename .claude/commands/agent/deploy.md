---
description: Use specialized agents to design, implement, and manage deployment pipelines and CI/CD workflows
argument-hint: <deployment-target or pipeline-type>
allowed-tools: Task, mcp__gemini-cli__ask-gemini, mcp__gemini-cli__brainstorm
---

## Multi-Agent Deployment & CI/CD System

**Deployment Target:** $ARGUMENTS

**Gemini MCP Support**: This command has access to Gemini MCP tools for advanced analysis, brainstorming, and efficient problem solving to save context usage.

Execute comprehensive deployment and CI/CD implementation using specialized agents for maximum reliability and efficiency:

### 1. **Deployment Team Assembly**

**Core Deployment Agents:**

- **devops-engineer**: Infrastructure and deployment pipeline design
- **security-auditor**: Security validation and compliance checking
- **perf-optimizer**: Performance testing and optimization
- **test-engineer**: Automated testing integration

### 2. **Infrastructure Assessment**

**Team:** devops-engineer + security-auditor

**Objectives:**

- Analyze current infrastructure requirements
- Design scalable deployment architecture
- Implement security best practices
- Setup monitoring and alerting

**Deliverables:**

- Infrastructure-as-Code templates
- Security configuration guidelines
- Monitoring dashboard setup
- Backup and disaster recovery plan

### 3. **CI/CD Pipeline Design**

**Team:** devops-engineer + test-engineer

**Pipeline Components:**

- **Source Control Integration**: Git webhooks and branch strategies
- **Build Automation**: Compilation, packaging, and artifact generation
- **Testing Integration**: Unit, integration, and end-to-end tests
- **Security Scanning**: SAST, DAST, and dependency vulnerability checks
- **Deployment Automation**: Blue-green, canary, or rolling deployments

**Key Features:**

- Automated rollback mechanisms
- Environment-specific configurations
- Approval workflows for production
- Deployment metrics and health checks

### 4. **Security & Compliance Integration**

**Team:** security-auditor + devops-engineer

**Security Measures:**

- Container scanning and hardening
- Secrets management (HashiCorp Vault, AWS Secrets Manager)
- Access control and RBAC implementation
- Compliance monitoring (SOC2, GDPR, HIPAA)
- Audit logging and traceability

### 5. **Performance & Monitoring Setup**

**Team:** perf-optimizer + devops-engineer

**Monitoring Stack:**

- Application Performance Monitoring (APM)
- Infrastructure monitoring (CPU, memory, disk, network)
- Log aggregation and analysis
- Error tracking and alerting
- User experience monitoring

**Performance Optimization:**

- Load testing automation
- Performance regression detection
- Auto-scaling configuration
- Resource optimization recommendations

### 6. **Testing Strategy Integration**

**Team:** test-engineer + devops-engineer

**Test Automation:**

- Unit test execution in CI
- Integration test environments
- End-to-end test automation
- Performance test integration
- Security test automation

**Quality Gates:**

- Code coverage requirements
- Security scan thresholds
- Performance benchmarks
- Manual approval processes

### 7. **Deployment Environments**

**Environment Strategy:**

- **Development**: Continuous deployment from feature branches
- **Staging**: Production-like environment for final testing
- **Production**: Controlled deployment with monitoring and rollback

**Configuration Management:**

- Environment-specific variables
- Feature flag management
- Database migration strategies
- Service dependency management

### 8. **Monitoring & Observability**

**Team:** devops-engineer + perf-optimizer

**Observability Stack:**

- Metrics collection (Prometheus, DataDog, CloudWatch)
- Distributed tracing (Jaeger, Zipkin)
- Log management (ELK, Splunk, Loki)
- Alerting and notification systems
- Dashboards and visualization

### 9. **Disaster Recovery & Business Continuity**

**Team:** devops-engineer + security-auditor

**DR Strategy:**

- Backup automation and testing
- Multi-region deployment strategies
- Data replication and consistency
- Recovery time and point objectives (RTO/RPO)
- Incident response procedures

### 10. **Documentation & Knowledge Transfer**

**Team:** devops-engineer + test-engineer

**Documentation:**

- Deployment runbooks and procedures
- Troubleshooting guides
- Architecture diagrams and documentation
- Team onboarding materials
- Incident post-mortem templates

### **Execution Plan**

1. **Infrastructure Setup** (devops-engineer)
2. **Security Framework** (security-auditor)
3. **CI/CD Pipeline Implementation** (devops-engineer + test-engineer)
4. **Monitoring Setup** (perf-optimizer + devops-engineer)
5. **Testing Integration** (test-engineer)
6. **Documentation** (devops-engineer)
7. **Production Deployment** (Full team coordination)

### **Success Metrics**

- Deployment frequency and success rate
- Mean time to deployment (MTTD)
- Mean time to recovery (MTTR)
- Security vulnerability detection and resolution time
- System uptime and availability
- Performance metrics and user experience scores

**Agent Coordination:** Each agent brings specialized expertise while maintaining awareness of the overall deployment ecosystem for seamless integration and optimal results.
