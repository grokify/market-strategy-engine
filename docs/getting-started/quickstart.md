# Quick Start

This guide walks you through running your first analysis.

## 1. Create a Minimal Analysis File

Create `analysis.json`:

```json
{
  "id": "demo",
  "name": "Demo Analysis",
  "focusVendorId": "us",
  "market": {
    "id": "widgets",
    "name": "Widget Market",
    "categories": [
      {
        "id": "core",
        "name": "Core Features",
        "capabilities": [
          {"id": "speed", "name": "Speed", "gapType": "product"},
          {"id": "reliability", "name": "Reliability", "gapType": "structural"}
        ]
      }
    ]
  },
  "segments": [
    {"id": "smb", "name": "SMB"},
    {"id": "enterprise", "name": "Enterprise"}
  ],
  "vendors": [
    {"id": "us", "name": "Our Company", "color": "#8b5cf6"},
    {"id": "them", "name": "Competitor", "color": "#ef4444"}
  ],
  "weights": [
    {"segmentId": "smb", "capabilityId": "speed", "weight": 0.7},
    {"segmentId": "smb", "capabilityId": "reliability", "weight": 0.3},
    {"segmentId": "enterprise", "capabilityId": "speed", "weight": 0.3},
    {"segmentId": "enterprise", "capabilityId": "reliability", "weight": 0.7}
  ],
  "scores": [
    {"vendorId": "us", "capabilityId": "speed", "score": 9.0},
    {"vendorId": "us", "capabilityId": "reliability", "score": 6.0},
    {"vendorId": "them", "capabilityId": "speed", "score": 7.0},
    {"vendorId": "them", "capabilityId": "reliability", "score": 9.0}
  ]
}
```

## 2. Validate the Data

```bash
mse validate analysis.json
```

Expected output:

```
Validating: analysis.json

Analysis: Demo Analysis
Market: Widget Market
Vendors: 2
Segments: 2
Capabilities: 2
Scores: 4
Weights: 4

✓ Validation passed
```

## 3. Run Analysis

```bash
mse analyze analysis.json
```

This outputs a text summary showing:

- Readiness scores per segment
- Top strategic priorities
- Gap analysis details

## 4. Generate HTML Report

```bash
mse analyze analysis.json --format html --output report.html
```

Open `report.html` in your browser to see the visual report.

## 5. What's Next?

- [Your First Analysis](first-analysis.md) - Build a complete real-world analysis
- [Data Model](../concepts/data-model.md) - Understand all the fields
- [CLI Reference](../cli/commands.md) - See all command options
