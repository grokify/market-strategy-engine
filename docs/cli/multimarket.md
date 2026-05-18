# mse multimarket

Run cross-market analysis to evaluate a vendor across multiple product markets.

## Usage

```bash
mse multimarket [files...] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `files` | Paths to market analysis JSON files (one or more required) |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--focus` | | | Focus vendor ID (required) |
| `--name` | | | Name for the combined analysis |
| `--format` | `-f` | `summary` | Output format: summary, json, html |
| `--output` | `-o` | | Output file (required for html format) |

## Output Formats

### summary (default)

Human-readable text output:

```bash
mse multimarket --focus acme-crm sales.json marketing.json support.json
```

```
═══════════════════════════════════════════════════════════════
  MULTI-MARKET ANALYSIS: Acme CRM Portfolio
═══════════════════════════════════════════════════════════════

Focus Vendor: Acme CRM
Markets Analyzed: 3

───────────────────────────────────────────────────────────────
  CROSS-MARKET INSIGHTS
───────────────────────────────────────────────────────────────

  Strongest Market:  Marketing Automation
  Weakest Market:    Customer Support
  Strongest Segment: SMB
  Weakest Segment:   Enterprise

───────────────────────────────────────────────────────────────
  MARKET PERFORMANCE
───────────────────────────────────────────────────────────────

  Marketing             #1 of 4  85%  [ready]
  Sales CRM             #2 of 5  72%  [adjacent]
  Support               #4 of 6  45%  [partial]
```

### json

Full multi-market analysis as JSON:

```bash
mse multimarket --focus acme-crm sales.json marketing.json -f json
```

```json
{
  "id": "portfolio",
  "name": "Acme CRM Portfolio",
  "focusVendorId": "acme-crm",
  "markets": [...],
  "marketSummaries": [...],
  "crossMarketComparison": {
    "strongestMarketId": "marketing",
    "weakestMarketId": "support",
    "strongestSegmentId": "smb",
    "weakestSegmentId": "enterprise"
  }
}
```

### html

Dark-themed HTML report with interactive visualizations:

```bash
mse multimarket --focus acme-crm sales.json marketing.json support.json \
    --format html --output portfolio.html
```

!!! note "Output Required"
    HTML format requires the `--output` flag.

## Examples

Basic multi-market analysis:

```bash
mse multimarket --focus acme-crm sales.json marketing.json support.json
```

Generate HTML portfolio report:

```bash
mse multimarket \
    --focus acme-crm \
    --name "Acme CRM Portfolio Analysis" \
    sales-crm.json marketing.json support.json \
    --format html \
    --output portfolio.html
```

Export full JSON for processing:

```bash
mse multimarket --focus acme-crm *.json -f json > portfolio.json
```

## Data Requirements

Each market analysis file must:

- Use the same `focusVendorId` across all files
- Use consistent segment IDs (e.g., "smb", "enterprise")
- Be independently valid (pass `mse validate`)

Different vendors and capabilities per market are handled automatically.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (missing focus vendor, file not found, invalid JSON) |
