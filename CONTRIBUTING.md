# Contributing to Volcano LLM

Thank you for your interest in contributing to Volcano LLM! We're building the world's first fluid software runtime, and we need your help to make it even better.

## 🌋 Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:
- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Respect differing viewpoints

## 🚀 Getting Started

### Prerequisites
- Go 1.23 or higher
- Docker & Docker Compose
- Git
- Basic understanding of Temporal workflows (optional but helpful)

### Development Setup
```bash
# Clone the repository
git clone https://github.com/yourusername/volcano-llm
cd volcano-llm

# Install dependencies
go mod download

# Run tests
go test ./...

# Start development environment
docker-compose up -d

# Run the API server locally
go run cmd/api-server/main.go
```

## 🔥 Areas We Need Help

### 1. Core Features
- [ ] Additional tool implementations
- [ ] More workflow patterns
- [ ] Performance optimizations
- [ ] Multi-language support

### 2. Git-Native Runtime
- [ ] Advanced hot-reload strategies
- [ ] Distributed Git synchronization
- [ ] Conflict resolution algorithms
- [ ] Branch merge strategies

### 3. Temporal Integration
- [ ] Custom activity implementations
- [ ] Workflow templates
- [ ] Saga pattern examples
- [ ] Performance benchmarks

### 4. Documentation
- [ ] Tutorial videos
- [ ] Architecture diagrams
- [ ] API documentation
- [ ] Example applications

## 📝 How to Contribute

### 1. Find an Issue
Browse our [issues](https://github.com/yourusername/volcano-llm/issues) labeled:
- `good first issue` - Perfect for newcomers
- `help wanted` - We need your expertise
- `enhancement` - New features
- `bug` - Something needs fixing

### 2. Fork & Branch
```bash
# Fork the repository on GitHub

# Clone your fork
git clone https://github.com/yourusername/volcano-llm
cd volcano-llm

# Add upstream remote
git remote add upstream https://github.com/original/volcano-llm

# Create a feature branch
git checkout -b feature/your-feature-name
```

### 3. Make Your Changes

#### Adding a New Tool
1. Create tool definition in `repos/tools/`
2. Implement tool in `pkg/tools/`
3. Add tests in `pkg/tools/your_tool_test.go`
4. Register in `pkg/tools/registry.go`

Example:
```go
// pkg/tools/weather_tool.go
type WeatherTool struct {
    BaseTools
}

func (t *WeatherTool) Execute(ctx context.Context, input ToolInput) (ToolOutput, error) {
    // Your implementation
}
```

#### Adding a Workflow
1. Define workflow in `repos/workflows/`
2. Implement activities in `pkg/temporal/activities/`
3. Add workflow logic in `pkg/temporal/workflows/`
4. Write integration tests

### 4. Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/tools/...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./pkg/...

# Run integration tests
docker-compose up -d
go test -tags=integration ./test/...
```

### 5. Code Style

We use standard Go formatting:
```bash
# Format your code
go fmt ./...

# Run linter
golangci-lint run

# Check for issues
go vet ./...
```

#### Style Guidelines
- Keep functions small and focused
- Use meaningful variable names
- Add comments for exported functions
- Follow Go idioms and best practices
- Performance matters - benchmark critical paths

### 6. Commit Messages

Follow conventional commits:
```
feat: add weather tool for forecast queries
fix: correct entity extraction for dates
docs: update architecture diagram
test: add benchmarks for calculator tool
refactor: simplify workflow execution logic
```

### 7. Submit Pull Request

1. Push your branch to your fork
2. Create a Pull Request against `main`
3. Fill out the PR template
4. Wait for review and address feedback

#### PR Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation
- [ ] Performance improvement

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No hardcoded secrets
```

## 🏗️ Architecture Decisions

### When Adding Tools
- Tools must be deterministic
- Execution time should be <100ms
- No external API calls in Fast Path tools
- Use Temporal activities for long-running operations

### When Adding Workflows
- Keep workflows idempotent
- Use Temporal's built-in retry mechanisms
- Document failure scenarios
- Add monitoring/metrics

### Git-Native Principles
- All configuration via Git files
- Hot-reload must work without restart
- Multi-tenant isolation via branches
- No hardcoded values

## 🐛 Reporting Issues

### Bug Reports
Include:
- Volcano LLM version
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs

### Feature Requests
Describe:
- The problem you're solving
- Proposed solution
- Alternative approaches considered
- Impact on existing features

## 📊 Performance Contributions

If you're optimizing performance:
1. Include benchmark before/after
2. Document the optimization
3. Ensure no regression in functionality
4. Consider memory allocation

Example:
```go
// Before: 45ms, 1000 allocations
// After: 25ms, 10 allocations
func BenchmarkCalculator(b *testing.B) {
    // Your benchmark
}
```

## 🎓 Learning Resources

- [Temporal Documentation](https://docs.temporal.io)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Git Workflows](https://www.atlassian.com/git/tutorials/comparing-workflows)
- [Our Architecture Docs](docs/architecture.md)

## 💬 Communication

- **GitHub Issues**: Bug reports and feature requests
- **Discussions**: General questions and ideas
- **Slack**: Join our community at volcano-llm.slack.com
- **Twitter**: Follow @VolcanoLLM for updates

## 🙏 Recognition

Contributors will be:
- Listed in our CONTRIBUTORS.md file
- Mentioned in release notes
- Invited to our contributor Slack channel
- Eligible for Volcano LLM swag!

## Questions?

Don't hesitate to ask! Open an issue with the `question` label or reach out on Slack.

---

Thank you for helping make Volcano LLM the future of fluid software! 🌋