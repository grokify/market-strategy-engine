// Package model defines the core data types for the Market Strategy Engine.
//
// The Market Strategy Engine is a strategic capability gap analysis system that:
//   - Maps competitors across product lines and market segments
//   - Scores capabilities with segment-specific weightings
//   - Computes weighted gaps to identify strategic priorities
//   - Generates prioritized roadmaps for market expansion
//
// # Core Concepts
//
// Market: An industry vertical (e.g., Security, CRM, Fintech)
//
// Segment: A customer segment defined by size/maturity (SMB, Mid-Market, Enterprise)
//
// ProductCategory: A functional domain within a market (e.g., EDR, MDR, SIEM for Security)
//
// Capability: A measurable product attribute that can be scored and compared
//
// Vendor: A company being analyzed within a market
//
// # Data Flow
//
// Reference Data → Scores → Weights → Gap Analysis → Readiness → Priorities
//
// External data sources (Gartner, G2, internal assessments) feed into capability
// scores. Segment-specific weights determine importance. The engine computes
// gaps and generates prioritized actions for market expansion.
package model
