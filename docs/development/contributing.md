# Contributing

Contributions to Market Strategy Engine are welcome.

## Development Setup

```bash
# Clone the repository
git clone https://github.com/grokify/market-strategy-engine.git
cd market-strategy-engine

# Install dependencies
go mod download

# Build
go build -o mse ./cmd/mse

# Run tests
go test ./...

# Run linting
golangci-lint run
```

## Code Style

- Use `gofmt` for formatting
- Follow standard Go conventions
- Run `golangci-lint run` before submitting

## Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Common types:

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting changes |
| `refactor` | Code restructuring |
| `test` | Adding tests |
| `chore` | Maintenance tasks |

Examples:

```
feat(engine): add segment-level priority ranking
fix: handle nil weights in validation
docs: update CLI reference
```

## Project Structure

```
market-strategy-engine/
├── cmd/mse/           # CLI application
│   └── cmd/           # Cobra commands
├── model/             # Data types (source of truth)
├── engine/            # Analysis computation
├── report/            # Output generation (HTML)
├── validate/          # Input validation
├── colors/            # Color palettes
├── schema/            # Embedded JSON schemas
└── docs/              # Documentation
```

## Adding Features

1. **Data Model Changes**: Start in `model/` - Go structs are the source of truth
2. **Schema Generation**: Run `mse generate-schema` after model changes
3. **Validation**: Add validation rules in `validate/validate.go`
4. **Engine**: Implement computation logic in `engine/`
5. **Reports**: Update HTML generation in `report/`
6. **CLI**: Add/update commands in `cmd/mse/cmd/`
7. **Tests**: Add tests alongside your changes
8. **Docs**: Update documentation

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./engine/
go test -v ./validate/
```

## Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Ensure tests pass and linting is clean
5. Submit a pull request

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include steps to reproduce for bugs
- Include example data when possible
