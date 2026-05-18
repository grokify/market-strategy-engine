# mse analyze

Run gap analysis on a single market analysis file.

## Usage

```bash
mse analyze [file] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `file` | Path to market analysis JSON file (required) |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `summary` | Output format: summary, json, priorities, html |
| `--output` | `-o` | | Output file (required for html format) |
| `--skip-validate` | | `false` | Skip input validation before analysis |

## Output Formats

### summary (default)

Human-readable text output:

```bash
mse analyze analysis.json
```

```
═══════════════════════════════════════════════════════════════
  MARKET STRATEGY ANALYSIS: Sales CRM Market
═══════════════════════════════════════════════════════════════

Focus Vendor: Acme CRM
Market: Sales CRM

───────────────────────────────────────────────────────────────
  SEGMENT READINESS SCORES
───────────────────────────────────────────────────────────────

  ● SMB: 85/100
     Ready to compete effectively in SMB segment

  ◐ Mid-Market: 72/100
     Adjacent to Mid-Market segment; minor improvements needed
```

### json

Full analysis results as JSON:

```bash
mse analyze analysis.json --format json
```

```json
{
  "id": "security-analysis",
  "name": "Security Market",
  "gaps": [...],
  "readiness": [...],
  "priorities": [...],
  "comparisons": [...]
}
```

### priorities

Tabular priority list:

```bash
mse analyze analysis.json --format priorities
```

```
Rank | Capability           | Gap | Target Segment | Effort    | Gap Type
-----|---------------------|-----|----------------|-----------|----------
   1 | Telemetry Access    | 5.5 | Enterprise     | very_high | structural
   2 | Response Automation | 3.0 | Enterprise     | medium    | product
```

### html

Dark-themed HTML report:

```bash
mse analyze analysis.json --format html --output report.html
```

!!! note "Output Required"
    HTML format requires the `--output` flag.

## Examples

Basic analysis:

```bash
mse analyze crm-market.json
```

Generate HTML report:

```bash
mse analyze crm-market.json -f html -o report.html
```

Export full JSON:

```bash
mse analyze crm-market.json -f json > results.json
```

Skip validation (for debugging):

```bash
mse analyze broken-data.json --skip-validate
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (file not found, invalid JSON, validation failed) |
