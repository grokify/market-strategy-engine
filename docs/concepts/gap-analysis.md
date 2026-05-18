# Gap Analysis

Gap analysis is the core computation that identifies where your product falls short of competitors.

## How Gaps Are Computed

### Step 1: Find Benchmarks

For each capability, the engine finds the **benchmark** - the highest score among all vendors:

```
Capability: Pipeline Management
├── Vendor A: 7.0
├── Vendor B: 9.0  ← Benchmark
└── Vendor C: 8.5
```

### Step 2: Calculate Raw Gaps

For the focus vendor, compute the gap from benchmark:

```
Gap = BenchmarkScore - VendorScore
Gap = 9.0 - 7.0 = 2.0
```

A positive gap means the vendor is behind the benchmark.

### Step 3: Apply Segment Weights

Different segments care about capabilities differently. Weight the gap:

```
WeightedGap = Gap × SegmentWeight

SMB (weight 0.3):        2.0 × 0.3 = 0.6
Enterprise (weight 0.7): 2.0 × 0.7 = 1.4
```

The same 2-point gap matters more in Enterprise because Pipeline Management is more important there.

## Gap Severity

Weighted gaps map to severity levels:

| Severity | Weighted Gap | Meaning |
|----------|--------------|---------|
| None | ≤ 0 | Meets or exceeds benchmark |
| Low | < 0.5 | Minor improvement opportunity |
| Medium | 0.5 - 1.5 | Notable gap |
| High | 1.5 - 2.5 | Significant disadvantage |
| Critical | ≥ 2.5 | Blocking gap - must address |

## Gap Types

Capabilities have different gap types that indicate how to address them:

### Product Gaps

```json
{"id": "reporting", "name": "Reporting", "gapType": "product"}
```

- Can be built with engineering effort
- Typically lower effort to close
- Example: Missing dashboard feature

### Structural Gaps

```json
{"id": "scalability", "name": "Scalability", "gapType": "structural"}
```

- Requires architecture changes
- Higher effort, potentially quarters of work
- Example: Database doesn't scale horizontally

### Perception Gaps

```json
{"id": "brand-recognition", "name": "Brand Recognition", "gapType": "perception"}
```

- Market perception problem
- Marketing/positioning solution
- Example: Not seen as "enterprise-ready"

## Priority Ranking

Gaps are ranked by weighted gap (descending). This prioritizes:

1. Large gaps in important capabilities
2. Over smaller gaps in less important capabilities

```
Rank  Capability        Weighted Gap  Severity
──────────────────────────────────────────────
#1    Integrations      1.8           High
#2    Forecasting       1.2           Medium
#3    Mobile App        0.6           Medium
#4    Reporting         0.3           Low
```

## Effort Estimation

The engine estimates effort based on gap type and gap size:

| Gap Type | Small Gap (<2) | Medium Gap (2-4) | Large Gap (>4) |
|----------|---------------|------------------|----------------|
| Product | Low | Medium | High |
| Structural | Medium | High | Very High |
| Perception | Medium | High | High |

## Example Analysis

Given:

- Capability: Integrations
- Benchmark (Competitor): 9.0
- Your Score: 5.0
- Enterprise Weight: 0.35

Calculation:

```
Gap = 9.0 - 5.0 = 4.0
WeightedGap = 4.0 × 0.35 = 1.4
Severity = Medium (0.5-1.5 range)
```

Interpretation: "Notable gap affecting Enterprise competitiveness"
