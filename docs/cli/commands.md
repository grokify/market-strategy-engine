# CLI Commands

Market Strategy Engine provides the `mse` command-line tool for running analyses.

## Overview

```bash
mse [command] [flags]
```

## Available Commands

| Command | Description |
|---------|-------------|
| `analyze` | Run single-market gap analysis |
| `multimarket` | Run cross-market analysis |
| `validate` | Validate analysis JSON files |
| `generate-schema` | Generate JSON Schema files |
| `version` | Show version information |
| `help` | Show help for any command |

## Global Flags

```bash
-h, --help   Show help for any command
```

## Quick Examples

```bash
# Validate data
mse validate analysis.json

# Run analysis with text output
mse analyze analysis.json

# Generate HTML report
mse analyze analysis.json --format html --output report.html

# Multi-market analysis
mse multimarket --focus huntress edr.json mdr.json --format html -o portfolio.html

# Generate JSON schemas
mse generate-schema --output ./schemas

# Show version
mse version
```

## Command Details

- [analyze](analyze.md) - Single-market gap analysis
- [multimarket](multimarket.md) - Cross-market portfolio analysis
- [validate](validate.md) - Input validation and data integrity checks
