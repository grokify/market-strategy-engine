# Your First Analysis

This guide walks through building a complete competitive analysis from scratch.

## Step 1: Define Your Market

Start by identifying the market structure:

```json
{
  "market": {
    "id": "crm",
    "name": "CRM Software",
    "description": "Customer relationship management platforms",
    "categories": [
      {
        "id": "sales",
        "name": "Sales Automation",
        "capabilities": [
          {
            "id": "pipeline-mgmt",
            "name": "Pipeline Management",
            "description": "Track deals through sales stages",
            "gapType": "product"
          },
          {
            "id": "forecasting",
            "name": "Sales Forecasting",
            "description": "Predict revenue based on pipeline",
            "gapType": "product"
          }
        ]
      },
      {
        "id": "platform",
        "name": "Platform",
        "capabilities": [
          {
            "id": "integrations",
            "name": "Integrations",
            "description": "Connect to other business tools",
            "gapType": "structural"
          },
          {
            "id": "mobile",
            "name": "Mobile Experience",
            "description": "Native mobile apps",
            "gapType": "product"
          }
        ]
      }
    ]
  }
}
```

!!! tip "Gap Types"
    - **product**: Can be built with engineering effort
    - **structural**: Requires architecture changes
    - **perception**: Marketing/positioning problem

## Step 2: Define Segments

Identify who buys in this market:

```json
{
  "segments": [
    {
      "id": "smb",
      "name": "SMB",
      "description": "Small businesses with 1-100 employees",
      "characteristics": [
        "Price sensitive",
        "Values ease of use",
        "Limited IT resources"
      ]
    },
    {
      "id": "mid-market",
      "name": "Mid-Market",
      "description": "Growing companies with 100-1000 employees",
      "characteristics": [
        "Needs integrations",
        "Has some IT support",
        "Wants scalability"
      ]
    },
    {
      "id": "enterprise",
      "name": "Enterprise",
      "description": "Large organizations with 1000+ employees",
      "characteristics": [
        "Complex requirements",
        "Dedicated IT teams",
        "Compliance needs"
      ]
    }
  ]
}
```

## Step 3: Add Vendors

Include your company and key competitors:

```json
{
  "vendors": [
    {
      "id": "acme-crm",
      "name": "Acme CRM",
      "description": "Our product",
      "color": "#8b5cf6",
      "website": "https://acme-crm.com",
      "targetSegments": ["smb", "mid-market"]
    },
    {
      "id": "salesforce",
      "name": "Salesforce",
      "color": "#00a1e0",
      "targetSegments": ["mid-market", "enterprise"],
      "isPublic": true,
      "ticker": "CRM"
    },
    {
      "id": "hubspot",
      "name": "HubSpot",
      "color": "#ff7a59",
      "targetSegments": ["smb", "mid-market"],
      "isPublic": true,
      "ticker": "HUBS"
    }
  ]
}
```

## Step 4: Set Segment Weights

Define how important each capability is per segment:

```json
{
  "weights": [
    {"segmentId": "smb", "capabilityId": "pipeline-mgmt", "weight": 0.30},
    {"segmentId": "smb", "capabilityId": "forecasting", "weight": 0.15},
    {"segmentId": "smb", "capabilityId": "integrations", "weight": 0.20},
    {"segmentId": "smb", "capabilityId": "mobile", "weight": 0.35},

    {"segmentId": "mid-market", "capabilityId": "pipeline-mgmt", "weight": 0.25},
    {"segmentId": "mid-market", "capabilityId": "forecasting", "weight": 0.25},
    {"segmentId": "mid-market", "capabilityId": "integrations", "weight": 0.30},
    {"segmentId": "mid-market", "capabilityId": "mobile", "weight": 0.20},

    {"segmentId": "enterprise", "capabilityId": "pipeline-mgmt", "weight": 0.20},
    {"segmentId": "enterprise", "capabilityId": "forecasting", "weight": 0.30},
    {"segmentId": "enterprise", "capabilityId": "integrations", "weight": 0.35},
    {"segmentId": "enterprise", "capabilityId": "mobile", "weight": 0.15}
  ]
}
```

!!! warning "Weights Must Sum to 1.0"
    For each segment, all capability weights must sum to exactly 1.0.

## Step 5: Score Vendors

Rate each vendor on each capability (0-10 scale):

```json
{
  "scores": [
    {"vendorId": "acme-crm", "capabilityId": "pipeline-mgmt", "score": 8.5},
    {"vendorId": "acme-crm", "capabilityId": "forecasting", "score": 6.0},
    {"vendorId": "acme-crm", "capabilityId": "integrations", "score": 5.0},
    {"vendorId": "acme-crm", "capabilityId": "mobile", "score": 9.0},

    {"vendorId": "salesforce", "capabilityId": "pipeline-mgmt", "score": 9.0},
    {"vendorId": "salesforce", "capabilityId": "forecasting", "score": 9.5},
    {"vendorId": "salesforce", "capabilityId": "integrations", "score": 9.0},
    {"vendorId": "salesforce", "capabilityId": "mobile", "score": 7.0},

    {"vendorId": "hubspot", "capabilityId": "pipeline-mgmt", "score": 8.0},
    {"vendorId": "hubspot", "capabilityId": "forecasting", "score": 7.5},
    {"vendorId": "hubspot", "capabilityId": "integrations", "score": 8.5},
    {"vendorId": "hubspot", "capabilityId": "mobile", "score": 8.0}
  ]
}
```

## Step 6: Run Analysis

Combine everything into one file and run:

```bash
mse validate crm-analysis.json
mse analyze crm-analysis.json --format html --output crm-report.html
```

## Interpreting Results

The report shows:

1. **Readiness Scores** - How ready you are to compete in each segment
2. **Strategic Priorities** - Which capabilities to improve first
3. **Competitive Comparison** - How you rank against competitors
4. **Gap Details** - Specific gaps with severity levels

## Next Steps

- Add more capabilities as you discover them
- Update scores as products evolve
- Track changes over time with versioned files
- Run multi-market analysis across product lines
