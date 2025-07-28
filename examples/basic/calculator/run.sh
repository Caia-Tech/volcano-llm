#!/bin/bash

echo "🧮 Volcano LLM Calculator Example"
echo "================================="
echo ""

# Simple addition
echo "1. Simple Addition: 42 + 58"
echo "Request:"
cat request-add.json | jq .
echo ""
echo "Response:"
curl -s -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d @request-add.json | jq .
echo ""

# Complex expression
echo "2. Complex Expression: (100 + 50) * 2 - 25"
echo "Request:"
cat request-complex.json | jq .
echo ""
echo "Response:"
curl -s -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d @request-complex.json | jq .
echo ""

# Session-based calculation
echo "3. Session-Based Calculation"
echo "First: Calculate 100 + 50"
RESULT1=$(curl -s -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "Calculate 100 + 50", "sessionId": "calc-demo"}')
echo "$RESULT1" | jq .
echo ""

echo "Then: Multiply it by 2"
RESULT2=$(curl -s -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"text": "Multiply it by 2", "sessionId": "calc-demo"}')
echo "$RESULT2" | jq .
echo ""

echo "✅ All calculations completed with deterministic results!"