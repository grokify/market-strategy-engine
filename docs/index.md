# Market Strategy Engine

A strategic capability gap analysis system for competitive market intelligence.

## What It Does

Market Strategy Engine helps product and strategy teams:

- **Map the competitive landscape** across markets and segments
- **Identify capability gaps** between your product and competitors
- **Prioritize investments** based on weighted segment importance
- **Track readiness** for entering new market segments
- **Generate visual reports** for stakeholder communication

## Key Features

- **Segment-Weighted Scoring** - Different segments value capabilities differently
- **Multi-Market Analysis** - Analyze across EDR, MDR, SIEM, etc. in one view
- **Gap Type Classification** - Product, structural, and perception gaps
- **Readiness Levels** - Actionable tiers from "Ready" to "Distant"
- **Dark-Themed HTML Reports** - Professional visualizations
- **JSON Schema Validation** - Catch data errors early

## Quick Example

```bash
# Validate your analysis data
mse validate security-market.json

# Generate an HTML report
mse analyze security-market.json --format html --output report.html

# Analyze across multiple markets
mse multimarket --focus huntress edr.json mdr.json siem.json \
    --format html --output portfolio.html
```

## How It Works

1. **Define your market** - Categories, capabilities, and gap types
2. **Define segments** - SMB, Mid-Market, Enterprise, etc.
3. **Add vendors** - Your company and competitors
4. **Set weights** - How important is each capability per segment?
5. **Score vendors** - Rate each vendor on each capability (0-10)
6. **Run analysis** - Engine computes gaps, readiness, and priorities

## Next Steps

- [Installation](getting-started/installation.md) - Get the CLI installed
- [Quick Start](getting-started/quickstart.md) - Run your first analysis
- [Concepts](concepts/overview.md) - Understand the data model
