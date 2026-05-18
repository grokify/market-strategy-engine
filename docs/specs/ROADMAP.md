# Market Strategy Engine Roadmap

This document consolidates the project plan, active tasks, and feature roadmap.

## Overview

Market Strategy Engine is a strategic capability gap analysis system that:

- Maps competitors across product lines and market segments
- Scores capabilities with segment-specific weightings
- Computes weighted gaps to identify strategic priorities
- Generates prioritized roadmaps for market expansion

## Completed

### v0.1.0 - Foundation

- [x] **Data Model** - Go struct-first data model with JSON tags
  - Market, ProductCategory, Capability, GapType
  - Segment, SegmentWeight
  - Vendor, VendorProduct
  - DataSource, SourceType
  - CapabilityScore, ScoreNormalization
  - GapAnalysis, ReadinessScore, PriorityAction
  - Analysis (top-level container)
  - SegmentComparison, VendorSegmentScore, CapabilityComparison

- [x] **Cobra CLI** - Command-line interface
  - `mse generate-schema` - Generate JSON Schema from Go types
  - `mse validate` - Validate analysis JSON files
  - `mse analyze` - Run gap analysis with multiple output formats
  - `mse version` - Show version information

- [x] **Analysis Engine** - Core computation logic
  - `ComputeGaps()` - Calculate capability gaps vs benchmarks
  - `ComputeReadiness()` - Calculate segment readiness scores
  - `ComputePriorities()` - Generate prioritized action list
  - `ComputeComparisons()` - Generate competitive comparisons per segment

- [x] **HTML Report** - Dark-themed visualization
  - Segment readiness cards with scores and interpretations
  - Strategic priorities table with effort indicators
  - Detailed gap analysis per segment
  - Competitive comparison with vendor rankings
  - Capability breakdown with mini-bars

- [x] **Color System** - Reusable color palettes
  - `colors/palettes.go` with CompetitiveAnalysis, Categorical, Diverging, Sequential
  - Vendor color support (Color field on Vendor model)
  - Vendor legend in HTML reports
  - Consistent colors across all visualizations

### v0.2.0 - High-Value Features

- [x] **JSON Schema Generation** - Schema files generated from Go types
  - 16 schema files generated for all model types
  - Schemas embedded via `//go:embed` in schema package
  - `mse generate-schema` command outputs to schema/ directory
  - Includes competitive comparison types

- [x] **Input Validation** - Comprehensive data validation
  - `validate` package with detailed error reporting
  - Validates weights sum to 1.0 per segment
  - Validates all vendors have scores for all capabilities
  - Validates all referenced IDs exist
  - Validates scores in valid range (0-10)
  - Validation runs automatically before analysis
  - `mse validate` command with `--strict` flag

- [x] **Unit Tests** - Comprehensive test coverage
  - `engine/engine_test.go` - Tests for gap, readiness, priority, comparison computations
  - `validate/validate_test.go` - Tests for all validation rules
  - `colors/palettes_test.go` - Tests for color palette functions
  - All tests passing

- [x] **Multi-Market Visualization** - Cross-market analysis
  - `model/multimarket.go` - MultiMarketAnalysis, MarketSummary, CrossMarketComparison
  - `engine/multimarket.go` - Multi-market computation engine
  - `report/multimarket_html.go` - Multi-market HTML report
  - `mse multimarket` command for analyzing across multiple markets
  - Cross-market overview showing strongest/weakest markets and segments
  - Market performance matrix with rankings
  - Segment performance aggregated across markets

## In Progress

_No items currently in progress._

## Backlog

### Medium Priority

#### Additional Export Formats

Support multiple output formats beyond HTML and JSON.

**Formats to add:**

- [ ] Markdown report (for GitHub/GitLab wikis)
- [ ] CSV export (for spreadsheet analysis)
- [ ] PDF report (for executive presentations)

#### Radar/Spider Charts

Visualize vendor capabilities on a radar chart per segment.

**Requirements:**

- Show all capabilities as spokes
- Overlay multiple vendors on same chart
- Use vendor colors for consistency
- Include in HTML report

#### Trend Tracking

Compare scores across time periods.

**Requirements:**

- Support multiple snapshots with timestamps
- Calculate score deltas between periods
- Visualize trends (improving/declining)
- Alert on significant changes

#### Market Positioning Quadrant

Gartner-style 2D positioning visualization.

**Requirements:**

- Configurable X and Y axes (e.g., Vision vs Execution)
- Plot vendors as bubbles (size = market share or similar)
- Quadrant labels (Leaders, Challengers, Visionaries, Niche)
- Interactive hover for vendor details

### Low Priority

#### Documentation

Comprehensive project documentation.

- [ ] README with quick start guide
- [ ] CLI usage examples
- [ ] Data model reference
- [ ] Contribution guidelines

#### Configuration System

Support for config files and environment variables.

- [ ] YAML/JSON config file support
- [ ] Environment variable overrides
- [ ] Default config locations (~/.mse/config.yaml)

#### Data Import

Import data from external sources.

**Potential sources:**

- [ ] G2 Crowd ratings
- [ ] Gartner Magic Quadrant data
- [ ] CSV/Excel import
- [ ] API integrations

#### Interactive HTML Features

Add JavaScript interactivity to reports.

- [ ] Filter by segment
- [ ] Sort tables by column
- [ ] Drill-down into capabilities
- [ ] Toggle vendor visibility

## Design Decisions

### Go Struct-First Data Model

We use Go structs as the source of truth for data types, generating JSON Schema from them rather than the reverse. This ensures type safety and makes the Go code authoritative.

### Segment-Weighted Scoring

Capabilities have different importance in different segments. The weighting system allows the same capability scores to produce different gap analyses depending on the target segment.

### Three Gap Types

1. **Product** - Gaps that can be closed by building features
2. **Structural** - Gaps requiring architectural changes (harder)
3. **Perception** - Gaps in market perception (marketing problem)

### Readiness Levels

- **Ready** (80-100): Competitive in segment
- **Adjacent** (60-79): Minor improvements needed
- **Partial** (40-59): Significant work required
- **Distant** (0-39): Major investment needed

### Multi-Market Analysis

Cross-market analysis aggregates single-market analyses to show:

- Overall competitive position across all markets
- Strongest/weakest markets for the focus vendor
- Segment performance patterns (e.g., consistently strong in SMB)
- Market-specific insights and competitor identification

## Architecture

```
market-strategy-engine/
├── cmd/mse/           # CLI application
│   ├── main.go
│   └── cmd/           # Cobra commands
│       ├── root.go
│       ├── analyze.go
│       ├── validate.go
│       ├── multimarket.go
│       ├── generate_schema.go
│       └── version.go
├── model/             # Data types (source of truth)
│   ├── analysis.go
│   ├── market.go
│   ├── segment.go
│   ├── vendor.go
│   ├── score.go
│   ├── gap.go
│   ├── source.go
│   ├── competitive.go
│   └── multimarket.go
├── engine/            # Analysis computation
│   ├── engine.go
│   └── multimarket.go
├── report/            # Output generation
│   ├── html.go
│   └── multimarket_html.go
├── validate/          # Input validation
│   └── validate.go
├── colors/            # Color palettes
│   └── palettes.go
├── schema/            # Embedded JSON schemas
│   ├── schema.go
│   └── *.schema.json
└── docs/
    └── specs/         # Specifications and roadmap
        └── ROADMAP.md
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `mse analyze <file>` | Run single-market gap analysis |
| `mse multimarket <files...>` | Run cross-market analysis |
| `mse validate <file>` | Validate analysis JSON |
| `mse generate-schema` | Generate JSON Schema files |
| `mse version` | Show version information |

## References

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Structured Changelog](https://github.com/grokify/structured-changelog)
- [JSON Schema](https://json-schema.org/)
