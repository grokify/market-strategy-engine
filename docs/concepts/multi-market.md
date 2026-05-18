# Multi-Market Analysis

Multi-market analysis lets you evaluate a vendor's competitive position across multiple product markets simultaneously.

## When to Use Multi-Market

Use multi-market analysis when:

- Your company competes in multiple markets (e.g., Sales CRM, Marketing Automation, Support)
- You want to identify your strongest/weakest markets
- You need to understand segment performance across markets
- You're planning portfolio-level strategy

## How It Works

Multi-market analysis:

1. **Runs each single-market analysis** individually
2. **Aggregates results** across markets
3. **Identifies patterns** in strengths and weaknesses
4. **Generates portfolio view** for strategic planning

## Multi-Market Structure

```
MultiMarketAnalysis
├── FocusVendor: "Acme CRM"
├── Markets[]
│   ├── Sales CRM Analysis (computed)
│   ├── Marketing Automation Analysis (computed)
│   └── Customer Support Analysis (computed)
├── MarketSummaries[]
│   ├── Sales CRM: Rank #2, Readiness 75%
│   ├── Marketing: Rank #1, Readiness 85%
│   └── Support: Rank #4, Readiness 45%
└── CrossMarketComparison
    ├── StrongestMarket: "Marketing Automation"
    ├── WeakestMarket: "Customer Support"
    ├── StrongestSegment: "SMB"
    └── WeakestSegment: "Enterprise"
```

## Running Multi-Market Analysis

### CLI

```bash
mse multimarket \
    --focus acme-crm \
    --name "Acme CRM Portfolio" \
    sales-crm.json marketing.json support.json \
    --format html \
    --output portfolio.html
```

### Programmatic

```go
analyses := []*model.Analysis{salesCRM, marketing, support}

eng := engine.NewMultiMarket(analyses)
mma := eng.ComputeMultiMarket("portfolio", "Acme CRM Portfolio", "acme-crm")
```

## Cross-Market Insights

### Market Scores

Shows overall performance in each market:

| Market | Rank | Score | Level |
|--------|------|-------|-------|
| Marketing | #1 of 4 | 85% | Ready |
| Sales CRM | #2 of 5 | 72% | Adjacent |
| Support | #4 of 6 | 45% | Partial |

### Segment Patterns

Aggregates segment performance across markets:

| Segment | Avg Score | Best Market | Worst Market |
|---------|-----------|-------------|--------------|
| SMB | 82% | Marketing (90%) | Support (65%) |
| Mid-Market | 68% | Sales CRM (75%) | Support (55%) |
| Enterprise | 42% | Marketing (50%) | Support (30%) |

This reveals you're consistently strong in SMB and weak in Enterprise across all markets.

## Strategic Questions Answered

### "Where should we invest?"

Look at weak markets with high strategic value:

- Support readiness is low (45%)
- If Support is strategic, prioritize closing gaps there

### "What are our consistent strengths?"

Look at segment patterns:

- SMB is strong across all markets
- Leverage this for go-to-market messaging

### "Where are we bleeding?"

Look at weak segments in strong markets:

- Marketing is your best market, but Enterprise is only 50%
- Opportunity to expand upmarket in Marketing Automation

## Example Output

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

## Data Requirements

Each market analysis file must:

- Use the same `focusVendorId`
- Use consistent segment IDs (e.g., "smb", "enterprise")
- Be independently valid (pass `mse validate`)

The engine handles different vendors and capabilities per market.
