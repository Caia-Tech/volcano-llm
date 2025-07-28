# 🏗️ Volcano LLM Architecture

## Overview

Volcano LLM is built on a revolutionary architecture that treats Git as a runtime database, enabling fluid software that evolves without restarts.

## Core Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Client Applications                          │
│                    (REST API, SDK, CLI, UI)                         │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                      API Gateway (Gin)                               │
│  • Authentication           • Rate Limiting                          │
│  • Request Routing          • Session Management                     │
│  • Tenant Isolation         • Metrics Collection                     │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                    Intelligent Router                                │
│                                                                      │
│  if (simple_query) {        │  if (complex_workflow) {              │
│    route → Fast Path        │    route → Temporal Path              │
│  }                          │  }                                     │
└──────────┬──────────────────┴─────────────────────┬─────────────────┘
           │                                        │
┌──────────▼──────────┐                  ┌─────────▼─────────┐
│     Fast Path       │                  │  Temporal Path    │
│   (<100ms SLA)      │                  │   (Durable)       │
├─────────────────────┤                  ├───────────────────┤
│ • Classifier        │                  │ • Workflows       │
│ • Extractor         │                  │ • Activities      │
│ • Decomposer        │                  │ • Sagas          │
│ • Simple Tools      │                  │ • State Machines │
└──────────┬──────────┘                  └─────────┬─────────┘
           │                                        │
           └────────────────┬───────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────┐
│                    Git-Native Runtime Layer                          │
│                                                                      │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Tools     │  │  Workflows   │  │   Configs    │              │
│  │  Registry   │  │  Definitions │  │   Manager    │              │
│  └──────┬──────┘  └──────┬───────┘  └──────┬───────┘              │
│         │                │                  │                        │
│  ┌──────▼────────────────▼──────────────────▼──────┐               │
│  │              Git Repository                      │               │
│  │  • File-based configuration                     │               │
│  │  • Branch-based multi-tenancy                   │               │
│  │  • Commit-triggered hot reload                  │               │
│  │  • Complete audit trail                         │               │
│  └──────────────────────────────────────────────────┘               │
└──────────────────────────────────────────────────────────────────────┘
```

## Component Deep Dive

### 1. API Gateway Layer

**Purpose**: Entry point for all requests

**Components**:
- **Gin Router**: High-performance HTTP router
- **Middleware Stack**: Auth, rate limiting, logging
- **Session Manager**: Maintains conversation context
- **Metrics Collector**: Tracks performance and usage

**Key Code**: `/pkg/api/server.go`

### 2. Intelligent Router

**Purpose**: Decides execution path based on query complexity

**Decision Logic**:
```go
func routeRequest(req ExecuteRequest) ExecutionPath {
    classification := classify(req.Text)
    
    if requiresDurability(classification) {
        return TemporalPath
    }
    
    if isSimpleCalculation(classification) {
        return FastPath
    }
    
    if hasMultipleStages(classification) {
        return TemporalPath
    }
    
    return FastPath // Default
}
```

### 3. Fast Path Engine

**Purpose**: Ultra-fast execution for simple queries

**Components**:
- **Pattern Classifier**: Regex-based domain identification
- **Entity Extractor**: Pulls out numbers, dates, names
- **Task Decomposer**: Breaks query into executable steps
- **Tool Executor**: Runs tools with <25ms latency

**Performance Optimizations**:
- Pre-compiled regex patterns
- Zero-allocation string processing
- Memory-pooled objects
- Minimal abstraction layers

### 4. Temporal Integration

**Purpose**: Durable execution for complex workflows

**Architecture**:
```
┌────────────────┐     ┌─────────────────┐     ┌────────────────┐
│  Workflow      │────▶│    Temporal     │────▶│   Activities   │
│  Definition    │     │    Server       │     │                │
└────────────────┘     └─────────────────┘     └────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │  (Persistence)  │
                       └─────────────────┘
```

**Features**:
- Workflow versioning
- Activity retries with exponential backoff
- Signal handling (pause/resume)
- Complete execution history

### 5. Git-Native Runtime

**Purpose**: The revolutionary core - Git as runtime database

**Components**:

#### File Watcher
```go
type GitWatcher struct {
    repoPath string
    handlers map[string]ChangeHandler
}

func (w *GitWatcher) Watch() {
    for event := range w.events {
        if handler, ok := w.handlers[event.Type]; ok {
            handler.Handle(event)
        }
    }
}
```

#### Hot Reload Engine
```go
type HotReloader struct {
    toolRegistry     *ToolRegistry
    workflowRegistry *WorkflowRegistry
    configManager    *ConfigManager
}

func (r *HotReloader) Reload(change GitChange) error {
    switch change.FileType {
    case "tool":
        return r.reloadTool(change)
    case "workflow":
        return r.reloadWorkflow(change)
    case "config":
        return r.reloadConfig(change)
    }
}
```

#### Branch Manager (Multi-tenancy)
```go
type BranchManager struct {
    tenantBranches map[string]string // tenant -> branch
}

func (m *BranchManager) GetTenantBranch(tenantID string) string {
    return fmt.Sprintf("customer/%s", tenantID)
}
```

## Data Flow

### Simple Request Flow (Fast Path)

```
1. Client Request
   POST /api/v1/execute
   {"text": "Calculate 42 + 58"}
   ↓
2. API Gateway
   - Validate request
   - Create/load session
   ↓
3. Intelligent Router
   - Classify as "simple_math"
   - Route to Fast Path
   ↓
4. Fast Path Processing
   - Extract: [42, +, 58]
   - Decompose: SingleCalculation
   - Execute: CalculatorTool
   ↓
5. Response (24ms)
   {"result": 100, "duration": "24ms"}
```

### Complex Workflow Flow (Temporal Path)

```
1. Client Request
   POST /api/v1/execute
   {"text": "Process customer data and generate report"}
   ↓
2. API Gateway
   - Validate request
   - Check permissions
   ↓
3. Intelligent Router
   - Classify as "complex_workflow"
   - Route to Temporal Path
   ↓
4. Temporal Workflow
   - Start workflow execution
   - Run activities in sequence:
     • ExtractData (500ms)
     • ValidateData (300ms)
     • TransformData (800ms)
     • GenerateReport (1.2s)
   ↓
5. Response (3.8s)
   {"workflow_id": "wf-123", "status": "completed"}
```

## Git Integration Architecture

### Repository Structure
```
repos/
├── main/                    # Base configuration
│   ├── tools/              # Core tools
│   ├── workflows/          # Standard workflows
│   └── configs/            # Default settings
│
├── customer/acme-corp/     # ACME tenant
│   ├── tools/              # ACME-specific tools
│   ├── workflows/          # Custom workflows
│   └── configs/            # Custom settings
│
└── customer/globex-inc/    # Globex tenant
    ├── tools/              # Globex tools
    ├── workflows/          # Globex workflows
    └── configs/            # Globex settings
```

### Hot Reload Sequence
```
1. File Change Detected
   git commit -m "Update calculator precision"
   ↓
2. Git Hook Triggered
   post-commit: notify-volcano-runtime
   ↓
3. File Watcher Event
   Event{Type: "modified", Path: "tools/calculator.json"}
   ↓
4. Hot Reload Engine
   - Parse new tool definition
   - Validate configuration
   - Update tool registry
   ↓
5. Instant Availability
   Next request uses updated tool
   (Zero downtime!)
```

## Security Architecture

### Multi-Layer Security
```
┌─────────────────────────────────────┐
│         API Gateway                 │
│  • API Key validation              │
│  • Rate limiting per key           │
│  • Request sanitization            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Tenant Isolation               │
│  • Branch-based separation         │
│  • Task queue isolation            │
│  • Resource quotas                 │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Git Security                  │
│  • Signed commits                  │
│  • Branch protection               │
│  • Audit logging                   │
└─────────────────────────────────────┘
```

## Scalability Architecture

### Horizontal Scaling
```
                    Load Balancer
                         │
        ┌────────────────┼────────────────┐
        │                │                │
   API Server 1     API Server 2     API Server 3
        │                │                │
        └────────────────┼────────────────┘
                         │
                 Shared Git Repo
                 (Source of Truth)
```

### Temporal Scaling
```
   Temporal Frontend → Temporal Workers (1..N)
                            │
                      ┌─────┴─────┐
                      │           │
                 Worker 1    Worker 2 ... Worker N
                      │           │
                      └─────┬─────┘
                            │
                      Git Repository
```

## Performance Characteristics

### Fast Path Performance
- **Latency**: 20-30ms average
- **Throughput**: 10,000+ requests/second
- **CPU Usage**: <5% per request
- **Memory**: 10MB baseline + 1KB per request

### Temporal Path Performance
- **Workflow Start**: <1 second
- **Activity Execution**: Depends on complexity
- **Durability**: 100% (survives crashes)
- **Scalability**: Limited by worker count

### Git Operations Performance
- **Hot Reload**: <100ms
- **Branch Switch**: <50ms
- **File Watch**: Instant (inotify)
- **Commit Hook**: <10ms

## Deployment Architecture

### Docker Compose (Development)
```yaml
services:
  api:
    image: volcano-llm:latest
    depends_on: [temporal, postgres]
    
  temporal:
    image: temporalio/auto-setup:latest
    
  postgres:
    image: postgres:13
    
  git-sync:
    image: volcano-llm-git-sync:latest
    volumes:
      - ./repos:/repos
```

### Kubernetes (Production)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: volcano-llm
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: api
        image: volcano-llm:latest
      - name: git-sync
        image: k8s.gcr.io/git-sync:latest
```

## Monitoring & Observability

### Metrics Collection
```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   OpenTelemetry   │   Prometheus  │     │  Grafana   │
│   (Traces)   │────│   (Metrics)   │────│(Dashboards)│
└──────┬──────┘     └──────┬───────┘     └────────────┘
       │                   │
       └───────┬───────────┘
               │
         Volcano LLM
```

### Key Metrics
- Request latency (p50, p95, p99)
- Tool execution time
- Workflow completion rate
- Git sync latency
- Hot reload frequency

## Future Architecture Enhancements

### 1. Distributed Git
- Multiple Git replicas for HA
- Eventual consistency model
- Conflict resolution strategies

### 2. Edge Deployment
- Run Volcano LLM at edge locations
- Git sync via satellite/limited connectivity
- Offline-first architecture

### 3. Federation
- Connect multiple Volcano LLM instances
- Cross-region workflow execution
- Global tool registry

---

This architecture enables Volcano LLM to be the world's first truly fluid software runtime, where code and configuration flow like lava, constantly reshaping the application without ever cooling into a fixed form.