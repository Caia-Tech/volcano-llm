# Calculator Example

The simplest demonstration of Volcano LLM's deterministic execution.

## Overview

This example shows:
- Natural language math processing
- Sub-25ms response times
- 100% deterministic results
- No GPU required

## Sample Requests

### Simple Addition
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d @request-add.json
```

### Complex Expression
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d @request-complex.json
```

### Using Previous Result
```bash
# First request stores result in session
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "Calculate 100 + 50", "sessionId": "calc-session"}'

# Second request uses "it" to reference previous result
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "Multiply it by 2", "sessionId": "calc-session"}'
```

## Performance

Average response times:
- Simple operations: 20-25ms
- Complex expressions: 30-40ms
- Session-based: 25-30ms

## How It Works

1. **Classification**: Identifies math domain
2. **Entity Extraction**: Extracts numbers and operators
3. **Decomposition**: Builds calculation workflow
4. **Execution**: Runs calculator tool
5. **Result**: Returns deterministic answer

No neural networks, no randomness, just pure deterministic logic!