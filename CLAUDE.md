# CLAUDE.md

Project-specific instructions for Market Strategy Engine.

## Project Overview

Market Strategy Engine is a strategic capability gap analysis system for competitive market intelligence. It helps analyze competitive positioning across markets and segments.

**Current Version:** v0.1.0 (not yet tagged)

## Architecture

```
market-strategy-engine/
├── cmd/mse/           # CLI application (Cobra)
│   └── cmd/           # Subcommands: analyze, multimarket, validate, generate-schema
├── model/             # Data types (source of truth for JSON Schema)
├── engine/            # Gap analysis computation
│   ├── engine.go      # Single-market analysis
│   └── multimarket.go # Multi-market portfolio analysis
├── report/            # Output generation
│   ├── html.go        # Single-market HTML reports
│   └── multimarket_html.go  # Multi-market HTML reports
├── validate/          # Input validation
├── colors/            # Color palette system
├── schema/            # Embedded JSON schemas (generated from model/)
└── docs/              # MkDocs documentation site
```

## Key Concepts

### Gap Analysis Formula

```
Gap = BenchmarkScore - VendorScore
WeightedGap = Gap × SegmentWeight
```

### Readiness Scoring

```
Readiness = 100 - (TotalWeightedGap × 10)
```

Levels:
- **Ready** (80-100): Competitive in segment
- **Adjacent** (60-79): Minor improvements needed
- **Partial** (40-59): Significant work required
- **Distant** (0-39): Major investment needed

### Gap Severity

| Severity | Weighted Gap |
|----------|--------------|
| Critical | ≥ 2.5 |
| High | 1.5 - 2.5 |
| Medium | 0.5 - 1.5 |
| Low | < 0.5 |

### Gap Types

- **product**: Feature gaps (can be built)
- **structural**: Architecture gaps (platform changes)
- **perception**: Market perception (marketing)

## CLI Commands

```bash
mse analyze <file>           # Single-market analysis
mse multimarket <files...>   # Cross-market portfolio analysis
mse validate <file>          # Validate JSON input
mse generate-schema          # Generate JSON Schema files
```

## Development Commands

```bash
# Build
go build -o mse ./cmd/mse

# Test
go test ./...

# Lint
golangci-lint run

# Generate schemas (after model changes)
./mse generate-schema --output ./schema

# Build docs locally
mkdocs serve
```

## Example Data

Use CRM market examples in documentation:
- Markets: Sales CRM, Marketing Automation, Customer Support
- Vendors: Acme CRM, Salesforce, HubSpot
- Segments: SMB, Mid-Market, Enterprise

**Do NOT use security examples** (EDR, MDR, SIEM, Huntress) in documentation.

## Changelog Management

Uses structured-changelog format:
- `CHANGELOG.json` - Source of truth
- `CHANGELOG.md` - Generated via `schangelog generate CHANGELOG.json -o CHANGELOG.md`

## Roadmap

See `docs/specs/ROADMAP.md` for planned features including:
- CSV/Excel export
- Trend analysis over time
- API server mode
- Scenario planning ("what-if" analysis)

## Testing

Tests are in `*_test.go` files alongside source:
- `engine/engine_test.go` - Gap computation, readiness scoring
- `validate/validate_test.go` - Validation rules
- `colors/palettes_test.go` - Color functions

## Lint Fixes

Project uses nolint comments for gosec false positives:
- G203: template.HTML with numeric values and hardcoded colors
- G602: Slice bounds checks not recognized by static analysis

See `~/go/src/github.com/grokify/mogo/lintfix/` for common patterns.
