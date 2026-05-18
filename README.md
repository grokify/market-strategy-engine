# Market Strategy Engine

A strategic capability gap analysis system for competitive market intelligence.

[![Go Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

## Overview

Market Strategy Engine helps you analyze competitive positioning across markets and segments by:

- **Mapping competitors** across product lines and market segments
- **Scoring capabilities** with segment-specific weightings
- **Computing weighted gaps** to identify strategic priorities
- **Generating prioritized roadmaps** for market expansion
- **Visualizing competitive position** with dark-themed HTML reports

## Installation

```bash
go install github.com/grokify/market-strategy-engine/cmd/mse@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/market-strategy-engine.git
cd market-strategy-engine
go build -o mse ./cmd/mse
```

## Quick Start

### 1. Create an Analysis File

Create a JSON file defining your market, segments, vendors, capabilities, and scores:

```json
{
  "id": "my-analysis",
  "name": "My Market Analysis",
  "focusVendorId": "my-company",
  "market": {
    "id": "my-market",
    "name": "My Market",
    "categories": [
      {
        "id": "core",
        "name": "Core Features",
        "capabilities": [
          {"id": "cap-1", "name": "Feature A", "gapType": "product"},
          {"id": "cap-2", "name": "Feature B", "gapType": "structural"}
        ]
      }
    ]
  },
  "segments": [
    {"id": "smb", "name": "SMB"},
    {"id": "enterprise", "name": "Enterprise"}
  ],
  "vendors": [
    {"id": "my-company", "name": "My Company", "color": "#8b5cf6"},
    {"id": "competitor", "name": "Competitor", "color": "#ef4444"}
  ],
  "weights": [
    {"segmentId": "smb", "capabilityId": "cap-1", "weight": 0.6},
    {"segmentId": "smb", "capabilityId": "cap-2", "weight": 0.4},
    {"segmentId": "enterprise", "capabilityId": "cap-1", "weight": 0.3},
    {"segmentId": "enterprise", "capabilityId": "cap-2", "weight": 0.7}
  ],
  "scores": [
    {"vendorId": "my-company", "capabilityId": "cap-1", "score": 8.0},
    {"vendorId": "my-company", "capabilityId": "cap-2", "score": 5.0},
    {"vendorId": "competitor", "capabilityId": "cap-1", "score": 7.0},
    {"vendorId": "competitor", "capabilityId": "cap-2", "score": 9.0}
  ]
}
```

### 2. Validate Your Data

```bash
mse validate analysis.json
```

### 3. Run Analysis

```bash
# Summary output
mse analyze analysis.json

# HTML report
mse analyze analysis.json --format html --output report.html

# JSON output
mse analyze analysis.json --format json
```

### 4. Multi-Market Analysis

Analyze across multiple markets:

```bash
mse multimarket --focus my-company --name "Portfolio Analysis" \
    market1.json market2.json market3.json \
    --format html --output portfolio.html
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `mse analyze <file>` | Run single-market gap analysis |
| `mse multimarket <files...>` | Run cross-market analysis |
| `mse validate <file>` | Validate analysis JSON |
| `mse generate-schema` | Generate JSON Schema files |
| `mse version` | Show version information |

### Analyze Options

```bash
mse analyze analysis.json [flags]

Flags:
  -f, --format string   Output format: summary, json, priorities, html (default "summary")
  -o, --output string   Output file (required for html format)
      --skip-validate   Skip input validation
```

### Multimarket Options

```bash
mse multimarket [files...] [flags]

Flags:
      --focus string    Focus vendor ID (required)
      --name string     Name for the combined analysis
  -f, --format string   Output format: summary, json, html (default "summary")
  -o, --output string   Output file
```

## Core Concepts

### Gap Types

Capabilities are categorized by how gaps should be addressed:

| Type | Description | Effort |
|------|-------------|--------|
| `product` | Feature gaps that can be built | Lower |
| `structural` | Architecture gaps requiring platform changes | Higher |
| `perception` | Market perception gaps (marketing) | Variable |

### Segment Weights

Different segments value capabilities differently. Weights must sum to 1.0 per segment:

```json
{
  "weights": [
    {"segmentId": "smb", "capabilityId": "ease-of-use", "weight": 0.4},
    {"segmentId": "smb", "capabilityId": "price", "weight": 0.4},
    {"segmentId": "smb", "capabilityId": "features", "weight": 0.2},

    {"segmentId": "enterprise", "capabilityId": "ease-of-use", "weight": 0.1},
    {"segmentId": "enterprise", "capabilityId": "price", "weight": 0.2},
    {"segmentId": "enterprise", "capabilityId": "features", "weight": 0.7}
  ]
}
```

### Readiness Levels

Computed readiness scores map to actionable levels:

| Level | Score | Meaning |
|-------|-------|---------|
| Ready | 80-100 | Competitive in segment |
| Adjacent | 60-79 | Minor improvements needed |
| Partial | 40-59 | Significant work required |
| Distant | 0-39 | Major investment needed |

### Gap Severity

Weighted gaps determine priority:

| Severity | Weighted Gap | Action |
|----------|--------------|--------|
| Critical | ≥2.5 | Blocking - must address |
| High | 1.5-2.5 | Significant disadvantage |
| Medium | 0.5-1.5 | Notable gap |
| Low | <0.5 | Minor improvement |

## Data Model

### Analysis Structure

```
Analysis
├── Market
│   └── Categories[]
│       └── Capabilities[]
├── Segments[]
├── Vendors[]
├── Weights[]           (segment × capability weights)
├── Scores[]            (vendor × capability scores)
├── Sources[]           (data provenance)
└── Computed Results
    ├── Gaps[]
    ├── Readiness[]
    ├── Priorities[]
    └── Comparisons[]
```

### JSON Schema

JSON Schema files are available for validation:

```bash
mse generate-schema --output ./schemas
```

Or access embedded schemas programmatically:

```go
import "github.com/grokify/market-strategy-engine/schema"

schemaBytes, err := schema.Get(schema.SchemaAnalysis)
```

## HTML Reports

Reports feature a dark theme optimized for readability:

- **Segment Readiness Cards** - Score, level, and interpretation per segment
- **Strategic Priorities Table** - Ranked actions with effort indicators
- **Competitive Comparison** - Vendor rankings with capability breakdowns
- **Gap Analysis Details** - Per-segment gap breakdown with severity

### Vendor Colors

Assign colors to vendors for consistent visualization:

```json
{
  "vendors": [
    {"id": "my-company", "name": "My Company", "color": "#8b5cf6"},
    {"id": "competitor-a", "name": "Competitor A", "color": "#ef4444"},
    {"id": "competitor-b", "name": "Competitor B", "color": "#3b82f6"}
  ]
}
```

## Programmatic Usage

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/grokify/market-strategy-engine/engine"
    "github.com/grokify/market-strategy-engine/model"
    "github.com/grokify/market-strategy-engine/report"
    "github.com/grokify/market-strategy-engine/validate"
)

func main() {
    // Load analysis
    data, _ := os.ReadFile("analysis.json")
    var analysis model.Analysis
    json.Unmarshal(data, &analysis)

    // Validate
    result := validate.Validate(&analysis)
    if !result.IsValid() {
        panic(result.Error())
    }

    // Run analysis
    eng := engine.New(&analysis)
    eng.Run()

    // Generate HTML report
    f, _ := os.Create("report.html")
    defer f.Close()

    htmlReport := report.NewHTMLReport(&analysis)
    htmlReport.Write(f)
}
```

## Project Structure

```
market-strategy-engine/
├── cmd/mse/           # CLI application
├── model/             # Data types (source of truth)
├── engine/            # Analysis computation
├── report/            # Output generation (HTML)
├── validate/          # Input validation
├── colors/            # Color palettes
├── schema/            # Embedded JSON schemas
└── docs/              # Documentation
```

## Contributing

Contributions are welcome. Please ensure:

1. Code passes `golangci-lint run`
2. Tests pass with `go test ./...`
3. Commits follow [Conventional Commits](https://www.conventionalcommits.org/)

## License

MIT License - see [LICENSE](LICENSE) for details.

 [build-status-svg]: https://github.com/grokify/market-strategy-engine/actions/workflows/go-ci.yaml/badge.svg
 [build-status-url]: https://github.com/grokify/market-strategy-engine/actions/workflows/go-ci.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/market-strategy-engine
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/market-strategy-engine
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/market-strategy-engine
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/market-strategy-engine
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/market-strategy-engine/blob/main/LICENSE
