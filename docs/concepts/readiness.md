# Readiness Scoring

Readiness scores indicate how prepared a vendor is to compete in a market segment.

## How Readiness Is Computed

Readiness starts at 100 and is reduced based on weighted gaps:

```
Readiness = 100 - (TotalWeightedGap × 10)
```

The score is clamped to 0-100.

### Example

Given weighted gaps for SMB segment:

| Capability | Weighted Gap |
|------------|--------------|
| Ease of Use | 0.2 |
| Pricing | 0.4 |
| Support | 0.3 |
| **Total** | **0.9** |

Calculation:

```
Readiness = 100 - (0.9 × 10) = 91
```

## Readiness Levels

Scores map to actionable levels:

| Level | Score Range | Meaning |
|-------|-------------|---------|
| **Ready** | 80-100 | Competitive in this segment |
| **Adjacent** | 60-79 | Minor improvements needed |
| **Partial** | 40-59 | Significant work required |
| **Distant** | 0-39 | Major investment needed |

## Interpretations

The engine generates human-readable interpretations:

| Level | Example Interpretation |
|-------|----------------------|
| Ready | "Ready to compete effectively in SMB segment" |
| Adjacent | "Adjacent to Enterprise segment; minor improvements needed" |
| Partial | "Partial fit for Mid-Market segment; 2 critical gap(s) must be addressed" |
| Distant | "Distant from Enterprise segment; major investment required" |

## Critical Gaps

Gaps with severity "critical" are flagged:

```json
{
  "score": 45,
  "criticalGaps": ["integrations", "scalability"],
  "interpretation": "Partial fit for Enterprise; 2 critical gap(s) must be addressed"
}
```

Critical gaps often block segment entry entirely.

## Using Readiness Scores

### Market Entry Decisions

| Readiness | Recommendation |
|-----------|----------------|
| Ready | Launch/expand in segment |
| Adjacent | Address gaps, then expand |
| Partial | Prioritize gaps, build roadmap |
| Distant | Consider if segment is strategic |

### Resource Allocation

Higher readiness = lower investment needed:

```
Adjacent: 1-2 quarters of focused work
Partial:  2-4 quarters, possibly new hires
Distant:  Major initiative, possibly acquisitions
```

### Competitive Positioning

Compare readiness across segments:

```
Segment      Readiness  Level
─────────────────────────────
SMB          85         Ready
Mid-Market   62         Adjacent
Enterprise   38         Distant
```

This shows you're strongest in SMB and should focus expansion there.

## Top Gaps

Each readiness score includes the top gaps (by weighted gap):

```json
{
  "score": 72,
  "topGaps": [
    {"capabilityId": "integrations", "weightedGap": 1.2},
    {"capabilityId": "forecasting", "weightedGap": 0.8},
    {"capabilityId": "mobile", "weightedGap": 0.5}
  ]
}
```

Addressing these top gaps will most improve readiness.
