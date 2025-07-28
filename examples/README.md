# 🌋 Volcano LLM Examples

This directory contains practical examples demonstrating the power of fluid software.

## Directory Structure

```
examples/
├── basic/         # Simple examples for getting started
├── advanced/      # Complex workflows and integrations  
└── enterprise/    # Multi-tenant and production patterns
```

## Quick Tour

### 1. Basic Examples
- **Calculator**: Simple math operations (24ms responses)
- **Text Processing**: Entity extraction and classification
- **Hot Reload**: Live configuration updates

### 2. Advanced Examples
- **Data Pipeline**: Multi-stage ETL with Temporal
- **Report Generation**: Complex document workflows
- **API Integration**: External service orchestration

### 3. Enterprise Examples
- **Multi-Tenant**: Branch-based customer isolation
- **GitOps Workflow**: Full CI/CD via Git commits
- **Compliance**: Audit trail and regulatory reporting

## Running Examples

Each example includes:
- `README.md` - Detailed explanation
- `request.json` - Sample API request
- `workflow.yaml` - Workflow definition (if applicable)
- `tools/` - Custom tool definitions

To run any example:
```bash
cd examples/basic/calculator
./run.sh
```

Start with the basic examples and work your way up to enterprise patterns!