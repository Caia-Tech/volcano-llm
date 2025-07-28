# Data Pipeline Example

Enterprise-grade data processing with Temporal workflows.

## Overview

This example demonstrates:
- Multi-stage data pipeline (Extract → Transform → Load)
- Temporal workflow orchestration
- Error handling and retries
- Progress tracking via Temporal UI

## Workflow Definition

```yaml
name: CustomerDataPipeline
stages:
  - extract:    # Pull from multiple sources
  - validate:   # Check data quality
  - transform:  # Apply business rules
  - load:       # Write to destination
  - notify:     # Send completion notice
```

## Running the Pipeline

### Start Pipeline
```bash
./run.sh start-pipeline
```

### Monitor Progress
Open http://localhost:8088 to see:
- Workflow execution timeline
- Individual activity status
- Retry attempts (if any)
- Final results

### Simulate Failure
```bash
# This will trigger retry logic
./run.sh simulate-failure
```

## Hot Reload Demo

Watch how workflow updates without restart:

```bash
# Terminal 1: Start monitoring
./run.sh monitor

# Terminal 2: Update workflow
./run.sh update-workflow

# See immediate changes in Terminal 1!
```

## Key Features

1. **Durability**: Survives server crashes
2. **Visibility**: Full execution history
3. **Flexibility**: Hot-reload configuration
4. **Reliability**: Automatic retries