# Multi-Tenant Example

Complete customer isolation using Git branches.

## Overview

This example shows how Volcano LLM provides:
- Complete tenant isolation
- Custom tools per customer
- Independent configurations
- Separate audit trails

## Architecture

```
main branch
├── core tools
├── default workflows
└── base configuration

customer/acme-corp branch
├── inherits from main
├── custom ACME tools
└── ACME-specific workflows

customer/globex-inc branch
├── inherits from main
├── custom Globex tools
└── Globex compliance rules
```

## Setup Customer Environments

### 1. Create ACME Corp Environment
```bash
./setup-tenant.sh acme-corp
```

This creates:
- Git branch: `customer/acme-corp`
- Custom namespace: `acme-corp`
- Isolated task queue: `volcano-acme-corp`

### 2. Create Globex Inc Environment
```bash
./setup-tenant.sh globex-inc
```

### 3. Test Isolation
```bash
# ACME request - uses ACME branch
curl -X POST http://localhost:8080/api/v1/execute \
  -H "X-Tenant-ID: acme-corp" \
  -d '{"text": "Process ACME order 12345"}'

# Globex request - uses Globex branch
curl -X POST http://localhost:8080/api/v1/execute \
  -H "X-Tenant-ID: globex-inc" \
  -d '{"text": "Generate Globex compliance report"}'
```

## Customer-Specific Customization

### Add ACME-Specific Tool
```bash
# Switch to ACME branch
cd repos
git checkout customer/acme-corp

# Add custom validation
cat > tools/acme-po-validator.json << EOF
{
  "name": "ACMEPOValidator",
  "rules": [
    "po_number_format: ^PO-[0-9]{5}$",
    "max_amount: 1000000",
    "required_fields: [department, cost_center]"
  ]
}
EOF

git add tools/acme-po-validator.json
git commit -m "Add ACME PO validation rules"
```

### Add Globex Compliance Workflow
```bash
# Switch to Globex branch
git checkout customer/globex-inc

# Add compliance workflow
cat > workflows/globex-sox-compliance.yaml << EOF
name: GlobexSOXCompliance
stages:
  - verify_approvals
  - check_thresholds
  - generate_audit_trail
  - file_with_sec
EOF

git add workflows/globex-sox-compliance.yaml
git commit -m "Add SOX compliance workflow"
```

## Benefits

1. **Perfect Isolation**: No data leakage between customers
2. **Easy Customization**: Each customer gets their own Git branch
3. **Simple Billing**: Count commits/executions per branch
4. **Instant Provisioning**: New customer = new branch (seconds)
5. **Complete Audit Trail**: Git history per customer

## Monitoring

### Per-Customer Metrics
```bash
# ACME metrics
./metrics.sh acme-corp

# Globex metrics
./metrics.sh globex-inc
```

### Temporal UI
- ACME workflows: Filter by task queue `volcano-acme-corp`
- Globex workflows: Filter by task queue `volcano-globex-inc`

## Upgrades

### Apply Core Updates to All Customers
```bash
# Update main branch
git checkout main
# ... make changes ...
git commit -m "Core security update"

# Merge to each customer
for customer in acme-corp globex-inc; do
  git checkout customer/$customer
  git merge main
  git push
done
```

### Customer-Specific Updates
```bash
# Just commit to their branch
git checkout customer/acme-corp
# ... make ACME-specific changes ...
git commit -m "Update ACME validation rules"
```

## Security

- Each branch has separate access controls
- Customers can't see each other's branches
- API validates X-Tenant-ID header
- Temporal task queues are isolated
- Git commits are signed and audited