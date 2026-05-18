// Package schema provides embedded JSON Schema files for the Market Strategy Engine
// data model. Schemas are generated from Go types and embedded for runtime access.
package schema

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed *.schema.json
var schemaFS embed.FS

// SchemaName identifies available schemas.
type SchemaName string

const (
	SchemaAnalysis             SchemaName = "analysis.schema.json"
	SchemaMarket               SchemaName = "market.schema.json"
	SchemaSegment              SchemaName = "segment.schema.json"
	SchemaVendor               SchemaName = "vendor.schema.json"
	SchemaVendorProduct        SchemaName = "vendor_product.schema.json"
	SchemaCapability           SchemaName = "capability.schema.json"
	SchemaProductCategory      SchemaName = "product_category.schema.json"
	SchemaGapAnalysis          SchemaName = "gap_analysis.schema.json"
	SchemaReadinessScore       SchemaName = "readiness_score.schema.json"
	SchemaPriorityAction       SchemaName = "priority_action.schema.json"
	SchemaCapabilityScore      SchemaName = "capability_score.schema.json"
	SchemaSegmentWeight        SchemaName = "segment_weight.schema.json"
	SchemaSegmentComparison    SchemaName = "segment_comparison.schema.json"
	SchemaVendorSegmentScore   SchemaName = "vendor_segment_score.schema.json"
	SchemaCapabilityComparison SchemaName = "capability_comparison.schema.json"
	SchemaDataSource           SchemaName = "data_source.schema.json"
)

// Get returns the JSON Schema for the given schema name.
func Get(name SchemaName) ([]byte, error) {
	return schemaFS.ReadFile(string(name))
}

// GetAsMap returns the JSON Schema as a map for programmatic access.
func GetAsMap(name SchemaName) (map[string]any, error) {
	data, err := Get(name)
	if err != nil {
		return nil, err
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema %s: %w", name, err)
	}

	return schema, nil
}

// List returns all available schema names.
func List() []SchemaName {
	return []SchemaName{
		SchemaAnalysis,
		SchemaMarket,
		SchemaSegment,
		SchemaVendor,
		SchemaVendorProduct,
		SchemaCapability,
		SchemaProductCategory,
		SchemaGapAnalysis,
		SchemaReadinessScore,
		SchemaPriorityAction,
		SchemaCapabilityScore,
		SchemaSegmentWeight,
		SchemaSegmentComparison,
		SchemaVendorSegmentScore,
		SchemaCapabilityComparison,
		SchemaDataSource,
	}
}
