# CRM Market Example

This example demonstrates analyzing a CRM software market with three vendors competing across SMB, Mid-Market, and Enterprise segments.

## Complete Analysis File

```json
{
  "id": "crm-analysis",
  "name": "CRM Market Analysis",
  "description": "Competitive analysis of CRM software market",
  "focusVendorId": "acme-crm",
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
            "description": "Predict revenue based on pipeline data",
            "gapType": "product"
          },
          {
            "id": "lead-scoring",
            "name": "Lead Scoring",
            "description": "Automated lead qualification",
            "gapType": "product"
          }
        ]
      },
      {
        "id": "platform",
        "name": "Platform & Integration",
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
            "description": "Native mobile apps for iOS and Android",
            "gapType": "product"
          },
          {
            "id": "api",
            "name": "API & Extensibility",
            "description": "Developer platform and APIs",
            "gapType": "structural"
          }
        ]
      },
      {
        "id": "analytics",
        "name": "Analytics & Reporting",
        "capabilities": [
          {
            "id": "reporting",
            "name": "Reporting",
            "description": "Standard and custom reports",
            "gapType": "product"
          },
          {
            "id": "dashboards",
            "name": "Dashboards",
            "description": "Real-time visual dashboards",
            "gapType": "product"
          }
        ]
      }
    ]
  },
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
  ],
  "vendors": [
    {
      "id": "acme-crm",
      "name": "Acme CRM",
      "description": "SMB-focused CRM with excellent mobile experience",
      "color": "#8b5cf6",
      "targetSegments": ["smb", "mid-market"]
    },
    {
      "id": "salesforce",
      "name": "Salesforce",
      "description": "Enterprise CRM leader",
      "color": "#00a1e0",
      "targetSegments": ["mid-market", "enterprise"],
      "isPublic": true,
      "ticker": "CRM"
    },
    {
      "id": "hubspot",
      "name": "HubSpot",
      "description": "Inbound marketing and sales platform",
      "color": "#ff7a59",
      "targetSegments": ["smb", "mid-market"],
      "isPublic": true,
      "ticker": "HUBS"
    }
  ],
  "weights": [
    {"segmentId": "smb", "capabilityId": "pipeline-mgmt", "weight": 0.15},
    {"segmentId": "smb", "capabilityId": "forecasting", "weight": 0.05},
    {"segmentId": "smb", "capabilityId": "lead-scoring", "weight": 0.10},
    {"segmentId": "smb", "capabilityId": "integrations", "weight": 0.15},
    {"segmentId": "smb", "capabilityId": "mobile", "weight": 0.25},
    {"segmentId": "smb", "capabilityId": "api", "weight": 0.05},
    {"segmentId": "smb", "capabilityId": "reporting", "weight": 0.15},
    {"segmentId": "smb", "capabilityId": "dashboards", "weight": 0.10},

    {"segmentId": "mid-market", "capabilityId": "pipeline-mgmt", "weight": 0.15},
    {"segmentId": "mid-market", "capabilityId": "forecasting", "weight": 0.15},
    {"segmentId": "mid-market", "capabilityId": "lead-scoring", "weight": 0.10},
    {"segmentId": "mid-market", "capabilityId": "integrations", "weight": 0.20},
    {"segmentId": "mid-market", "capabilityId": "mobile", "weight": 0.10},
    {"segmentId": "mid-market", "capabilityId": "api", "weight": 0.10},
    {"segmentId": "mid-market", "capabilityId": "reporting", "weight": 0.10},
    {"segmentId": "mid-market", "capabilityId": "dashboards", "weight": 0.10},

    {"segmentId": "enterprise", "capabilityId": "pipeline-mgmt", "weight": 0.10},
    {"segmentId": "enterprise", "capabilityId": "forecasting", "weight": 0.15},
    {"segmentId": "enterprise", "capabilityId": "lead-scoring", "weight": 0.10},
    {"segmentId": "enterprise", "capabilityId": "integrations", "weight": 0.20},
    {"segmentId": "enterprise", "capabilityId": "mobile", "weight": 0.05},
    {"segmentId": "enterprise", "capabilityId": "api", "weight": 0.15},
    {"segmentId": "enterprise", "capabilityId": "reporting", "weight": 0.15},
    {"segmentId": "enterprise", "capabilityId": "dashboards", "weight": 0.10}
  ],
  "scores": [
    {"vendorId": "acme-crm", "capabilityId": "pipeline-mgmt", "score": 8.5},
    {"vendorId": "acme-crm", "capabilityId": "forecasting", "score": 6.0},
    {"vendorId": "acme-crm", "capabilityId": "lead-scoring", "score": 7.0},
    {"vendorId": "acme-crm", "capabilityId": "integrations", "score": 5.0},
    {"vendorId": "acme-crm", "capabilityId": "mobile", "score": 9.5},
    {"vendorId": "acme-crm", "capabilityId": "api", "score": 5.5},
    {"vendorId": "acme-crm", "capabilityId": "reporting", "score": 7.0},
    {"vendorId": "acme-crm", "capabilityId": "dashboards", "score": 7.5},

    {"vendorId": "salesforce", "capabilityId": "pipeline-mgmt", "score": 9.0},
    {"vendorId": "salesforce", "capabilityId": "forecasting", "score": 9.5},
    {"vendorId": "salesforce", "capabilityId": "lead-scoring", "score": 9.0},
    {"vendorId": "salesforce", "capabilityId": "integrations", "score": 9.0},
    {"vendorId": "salesforce", "capabilityId": "mobile", "score": 7.0},
    {"vendorId": "salesforce", "capabilityId": "api", "score": 9.5},
    {"vendorId": "salesforce", "capabilityId": "reporting", "score": 9.0},
    {"vendorId": "salesforce", "capabilityId": "dashboards", "score": 8.5},

    {"vendorId": "hubspot", "capabilityId": "pipeline-mgmt", "score": 8.0},
    {"vendorId": "hubspot", "capabilityId": "forecasting", "score": 7.5},
    {"vendorId": "hubspot", "capabilityId": "lead-scoring", "score": 8.5},
    {"vendorId": "hubspot", "capabilityId": "integrations", "score": 8.5},
    {"vendorId": "hubspot", "capabilityId": "mobile", "score": 8.0},
    {"vendorId": "hubspot", "capabilityId": "api", "score": 7.5},
    {"vendorId": "hubspot", "capabilityId": "reporting", "score": 8.0},
    {"vendorId": "hubspot", "capabilityId": "dashboards", "score": 8.5}
  ],
  "sources": [
    {
      "id": "g2",
      "name": "G2 Reviews",
      "type": "review",
      "url": "https://www.g2.com",
      "credibility": 8
    },
    {
      "id": "internal",
      "name": "Internal Assessment",
      "type": "internal",
      "credibility": 9
    }
  ]
}
```

## Running the Analysis

### Validate

```bash
mse validate crm-analysis.json
```

### Run Analysis

```bash
mse analyze crm-analysis.json
```

### Generate HTML Report

```bash
mse analyze crm-analysis.json --format html --output crm-report.html
```

## Expected Results

### Readiness Scores

| Segment | Score | Level | Interpretation |
|---------|-------|-------|----------------|
| SMB | 85 | Ready | Ready to compete effectively in SMB segment |
| Mid-Market | 68 | Adjacent | Adjacent to Mid-Market; minor improvements needed |
| Enterprise | 45 | Partial | Partial fit for Enterprise; significant work required |

### Top Strategic Priorities

1. **Integrations** (Enterprise) - Structural gap, high effort
2. **API & Extensibility** (Enterprise) - Structural gap, high effort
3. **Forecasting** (Mid-Market) - Product gap, medium effort
4. **Integrations** (Mid-Market) - Structural gap, high effort

### Key Insights

- **Strong in SMB**: Acme CRM's mobile experience (9.5) and ease of use make it competitive
- **Adjacent to Mid-Market**: Integration gaps (5.0 vs 9.0 benchmark) hold back expansion
- **Distant from Enterprise**: Structural gaps in API and integrations require major investment

## Multi-Market Extension

To analyze across multiple CRM-related markets (Sales, Marketing, Support), create separate analysis files and run:

```bash
mse multimarket \
    --focus acme-crm \
    --name "Acme CRM Portfolio" \
    sales-crm.json marketing-automation.json customer-support.json \
    --format html \
    --output portfolio.html
```
