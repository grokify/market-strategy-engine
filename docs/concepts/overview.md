# Concepts Overview

Market Strategy Engine is built around a few core concepts that work together to provide competitive intelligence.

## The Big Picture

```
┌─────────────────────────────────────────────────────────────┐
│                        ANALYSIS                              │
│  ┌─────────┐  ┌──────────┐  ┌─────────┐  ┌──────────────┐  │
│  │ MARKET  │  │ SEGMENTS │  │ VENDORS │  │   WEIGHTS    │  │
│  │         │  │          │  │         │  │ (per segment)│  │
│  │Categories│  │ SMB      │  │ You     │  │              │  │
│  │  └─Caps  │  │ Mid-Mkt  │  │ Comp A  │  │ Capability   │  │
│  │         │  │ Enterpr. │  │ Comp B  │  │ importance   │  │
│  └─────────┘  └──────────┘  └─────────┘  └──────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                      SCORES                           │   │
│  │         Vendor × Capability matrix (0-10)             │   │
│  └──────────────────────────────────────────────────────┘   │
│                           ↓                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                   COMPUTED RESULTS                    │   │
│  │  ┌──────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐ │   │
│  │  │ GAPS │ │ READINESS │ │ PRIORITIES│ │COMPARISONS│ │   │
│  │  └──────┘ └───────────┘ └───────────┘ └───────────┘ │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Core Entities

### Market

A market contains **categories** which contain **capabilities**. This hierarchy lets you organize complex markets:

```
CRM Market
├── Sales Automation
│   ├── Pipeline Management
│   ├── Forecasting
│   └── Lead Scoring
├── Marketing
│   ├── Email Campaigns
│   └── Landing Pages
└── Platform
    ├── Integrations
    └── API
```

### Segments

Segments are buyer groups with different needs. Common segments:

- **SMB** - Small businesses, price-sensitive, ease-of-use focused
- **Mid-Market** - Growing companies, need scalability
- **Enterprise** - Large orgs, complex requirements, compliance needs

### Vendors

Companies competing in the market. One vendor is the "focus vendor" - typically your company - which the analysis centers on.

### Weights

How important is each capability for each segment? Weights are 0-1 and must sum to 1.0 per segment.

Example: SMB cares more about ease-of-use than advanced features.

### Scores

How well does each vendor perform on each capability? Scores are 0-10.

## Computed Results

### Gaps

The difference between benchmark (best score) and vendor score, weighted by segment importance:

```
WeightedGap = (BenchmarkScore - VendorScore) × SegmentWeight
```

### Readiness

Overall score (0-100) indicating how ready a vendor is to compete in a segment:

```
Readiness = 100 - (TotalWeightedGaps × 10)
```

### Priorities

Ranked list of capabilities to improve, sorted by weighted gap (highest first).

### Comparisons

Per-segment rankings showing how vendors compare overall and on each capability.

## Key Formulas

| Metric | Formula |
|--------|---------|
| Gap | `BenchmarkScore - VendorScore` |
| Weighted Gap | `Gap × SegmentWeight` |
| Readiness | `100 - (Σ WeightedGaps × 10)` |
| Normalized Score | `(WeightedScore / MaxPossible) × 100` |

## Next Steps

- [Data Model](data-model.md) - Detailed field reference
- [Gap Analysis](gap-analysis.md) - How gaps are computed
- [Readiness Scoring](readiness.md) - Understanding readiness levels
