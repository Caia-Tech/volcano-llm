# 🚀 Volcano LLM Quick Start Guide

Get Volcano LLM running in under 5 minutes!

## Prerequisites

- Docker & Docker Compose installed
- Git installed
- 8GB RAM minimum
- Ports 8080, 7233, and 8088 available

## 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/volcano-llm
cd volcano-llm

# Create necessary directories
mkdir -p repos/{tools,workflows,configs}
```

## 2. Start with Docker Compose

```bash
# Start all services
docker-compose up -d

# Check service health
docker-compose ps

# You should see:
# - volcano-api     (port 8080) - Main API
# - temporal        (port 7233) - Workflow engine  
# - temporal-ui     (port 8088) - Workflow UI
# - postgres        (port 5432) - Temporal storage
# - git-sync        (internal)  - Git watcher
```

## 3. Your First Request

### Simple Math (Fast Path - 24ms)
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "What is 42 plus 58?"}'

# Response:
{
  "success": true,
  "result": 100,
  "sessionId": "session-1234567890",
  "duration": "24.5ms",
  "deterministic": true
}
```

### Complex Workflow (Temporal Path)
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "Process sales data for customer ACME and generate monthly report"}'

# Response:
{
  "success": true,
  "result": {
    "workflow_id": "data-pipeline-1234",
    "status": "completed",
    "report_url": "/reports/acme-2025-01.pdf"
  },
  "duration": "4.5s",
  "temporal_workflow": true
}
```

## 4. Git-Native Hot Reload Demo

### Add a New Tool (Zero Downtime!)
```bash
# Create a new tool
cat > repos/tools/greeting-tool.json << EOF
{
  "name": "GreetingTool",
  "version": "1.0.0",
  "category": "text",
  "description": "Generates personalized greetings",
  "patterns": ["hello", "greet", "welcome"],
  "config": {
    "templates": [
      "Hello, {name}! Welcome to Volcano LLM!",
      "Greetings, {name}! How can I help you today?"
    ]
  }
}
EOF

# Commit to Git
cd repos
git add tools/greeting-tool.json
git commit -m "Add greeting tool"

# Tool is INSTANTLY available - no restart!
cd ..
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "Say hello to John"}'
```

### Modify a Workflow
```bash
# Update workflow configuration
cat > repos/workflows/data-pipeline.yaml << EOF
name: DataPipelineWorkflow
version: 1.1.0
max_duration: 2h  # Changed from 1h
retry_policy:
  maximum_attempts: 5  # Changed from 3
stages:
  - extract
  - validate
  - transform
  - load
  - report
EOF

# Commit changes
cd repos
git add workflows/data-pipeline.yaml
git commit -m "Increase pipeline timeout and retries"

# New workflows will use updated config immediately!
```

## 5. Multi-Tenant Setup

### Create Customer Branch
```bash
# Create isolated environment for customer
cd repos
git checkout -b customer/acme-corp

# Add customer-specific tool
cat > tools/acme-validator.json << EOF
{
  "name": "ACMEValidator",
  "version": "1.0.0",
  "description": "ACME Corp specific validation rules",
  "rules": ["must_have_po_number", "max_amount_1million"]
}
EOF

git add tools/acme-validator.json
git commit -m "Add ACME custom validator"

# ACME's requests will use their branch
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: acme-corp" \
  -d '{"text": "Validate purchase order PO-12345"}'
```

## 6. Monitor with Temporal UI

Open http://localhost:8088 to see:
- Running workflows
- Execution history
- Retry attempts
- Workflow timelines

## 7. Check System Status

```bash
# API health check
curl http://localhost:8080/health

# System status
curl http://localhost:8080/api/v1/status

# List available tools
curl http://localhost:8080/api/v1/tools
```

## 8. Development Workflow

### Local Development
```bash
# Run in development mode
go run cmd/api-server/main.go

# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./pkg/...
```

### Adding New Tools
1. Create tool definition in `repos/tools/`
2. Implement tool in `pkg/tools/`
3. Register in `pkg/tools/registry.go`
4. Commit to Git - instantly available!

### Adding Workflows
1. Define workflow in `repos/workflows/`
2. Implement activities in `pkg/temporal/`
3. Commit to Git - hot reloaded!

## 9. Configuration

### Environment Variables
```bash
# API Configuration
PORT=8080
API_KEYS=key1,key2  # Optional API keys
ENABLE_AUTH=false   # Enable authentication
ENABLE_TRACE=true   # Enable execution tracing

# Temporal Configuration
TEMPORAL_SERVER=localhost:7233
TEMPORAL_NAMESPACE=default
TEMPORAL_TASK_QUEUE=volcano-workflows

# Git Configuration
REPOS_DIR=/app/repos
GIT_SYNC_INTERVAL=5s
```

### Docker Compose Override
```yaml
# docker-compose.override.yml
version: '3.8'
services:
  api:
    environment:
      - ENABLE_AUTH=true
      - API_KEYS=prod-key-1,prod-key-2
```

## 10. Troubleshooting

### Services Won't Start
```bash
# Check logs
docker-compose logs -f

# Ensure ports are free
lsof -i :8080
lsof -i :7233
lsof -i :8088
```

### Git Sync Issues
```bash
# Check git-sync logs
docker-compose logs git-sync

# Manually trigger sync
docker-compose exec git-sync /app/sync.sh
```

### Temporal Connection Failed
```bash
# Ensure Temporal is running
docker-compose ps temporal

# Test connection
docker-compose exec api curl temporal:7233
```

## Next Steps

1. **Explore Examples**: Check the `examples/` directory
2. **Read Architecture**: See `docs/architecture.md`
3. **Contribute**: Read `CONTRIBUTING.md`
4. **Deploy Production**: See `docs/production.md`

## Get Help

- GitHub Issues: https://github.com/yourusername/volcano-llm/issues
- Documentation: https://volcano-llm.dev/docs
- Community Slack: https://volcano-llm.slack.com

---

**Congratulations!** 🎉 You're now running the world's first fluid software runtime!