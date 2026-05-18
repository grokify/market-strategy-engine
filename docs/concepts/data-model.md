# Data Model Reference

Complete reference for all data types in Market Strategy Engine.

## Analysis

The top-level container for a market analysis.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Human-readable name |
| `description` | string | No | Context about the analysis |
| `version` | string | No | Version string (e.g., "1.0.0") |
| `createdAt` | datetime | No | Creation timestamp |
| `updatedAt` | datetime | No | Last update timestamp |
| `focusVendorId` | string | No | Primary vendor being analyzed |
| `market` | Market | Yes | Market definition |
| `segments` | Segment[] | Yes | Market segments |
| `vendors` | Vendor[] | Yes | Competitors |
| `weights` | SegmentWeight[] | Yes | Capability weights per segment |
| `scores` | CapabilityScore[] | Yes | Vendor scores |
| `sources` | DataSource[] | No | Data provenance |
| `tags` | string[] | No | Classification tags |

## Market

Defines the market structure.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Display name |
| `description` | string | No | Market description |
| `categories` | ProductCategory[] | Yes | Product categories |

## ProductCategory

Groups related capabilities.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Display name |
| `description` | string | No | Category description |
| `capabilities` | Capability[] | Yes | Capabilities in this category |

## Capability

A specific feature or attribute that can be scored.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Display name |
| `description` | string | No | What this capability represents |
| `gapType` | GapType | No | How gaps should be addressed |
| `scoreRubric` | string | No | Scoring guidance |

### GapType Values

| Value | Description |
|-------|-------------|
| `product` | Feature gap - can be built |
| `structural` | Architecture gap - requires platform changes |
| `perception` | Market perception gap - marketing problem |

## Segment

A target market segment.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Display name |
| `description` | string | No | Segment description |
| `minEndpoints` | int | No | Minimum size (for sizing-based segments) |
| `maxEndpoints` | int | No | Maximum size |
| `characteristics` | string[] | No | Key buyer characteristics |

## Vendor

A company competing in the market.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Company name |
| `description` | string | No | Company description |
| `color` | string | No | Hex color for visualizations |
| `website` | string | No | Company website |
| `targetSegments` | string[] | No | Primary segments served |
| `founded` | int | No | Year founded |
| `headquarters` | string | No | HQ location |
| `isPublic` | bool | No | Publicly traded? |
| `ticker` | string | No | Stock ticker symbol |

## SegmentWeight

Defines capability importance per segment.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `segmentId` | string | Yes | Segment reference |
| `capabilityId` | string | Yes | Capability reference |
| `weight` | float | Yes | Importance (0-1) |
| `rationale` | string | No | Why this weight |

!!! warning "Weights Must Sum to 1.0"
    For each segment, all capability weights must sum to exactly 1.0.

## CapabilityScore

A vendor's score on a capability.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `vendorId` | string | Yes | Vendor reference |
| `capabilityId` | string | Yes | Capability reference |
| `score` | float | Yes | Score (0-10) |
| `rawValue` | string | No | Original value before normalization |
| `sourceId` | string | No | Data source reference |
| `asOfDate` | datetime | No | When score was measured |
| `notes` | string | No | Score justification |

## DataSource

Tracks where data comes from.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Source name |
| `type` | SourceType | No | Type of source |
| `url` | string | No | Source URL |
| `credibility` | int | No | Credibility score (1-10) |

### SourceType Values

| Value | Description |
|-------|-------------|
| `analyst` | Analyst report (Gartner, Forrester) |
| `review` | Review platform (G2, Capterra) |
| `internal` | Internal assessment |
| `customer` | Customer feedback |

## Computed Types

These are populated by the engine after running analysis.

### GapAnalysis

| Field | Type | Description |
|-------|------|-------------|
| `vendorId` | string | Vendor analyzed |
| `capabilityId` | string | Capability with gap |
| `segmentId` | string | Target segment |
| `benchmarkScore` | float | Best-in-class score |
| `benchmarkVendorId` | string | Who has best score |
| `vendorScore` | float | Focus vendor's score |
| `gap` | float | Raw gap (benchmark - vendor) |
| `segmentWeight` | float | Capability weight |
| `weightedGap` | float | Gap × Weight |
| `gapType` | GapType | How to address |
| `severity` | GapSeverity | How critical |

### ReadinessScore

| Field | Type | Description |
|-------|------|-------------|
| `vendorId` | string | Vendor |
| `segmentId` | string | Segment |
| `score` | float | Readiness (0-100) |
| `totalWeightedGap` | float | Sum of weighted gaps |
| `criticalGaps` | string[] | Critical gap IDs |
| `interpretation` | string | Human-readable assessment |

### PriorityAction

| Field | Type | Description |
|-------|------|-------------|
| `rank` | int | Priority order (1 = highest) |
| `capabilityId` | string | Capability to improve |
| `segmentId` | string | Target segment |
| `priority` | float | Priority score |
| `gapType` | GapType | How to address |
| `currentScore` | float | Current score |
| `targetScore` | float | Target score |
| `impact` | string | What improving enables |
| `effort` | EffortLevel | Implementation effort |
