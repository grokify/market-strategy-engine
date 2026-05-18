package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/market-strategy-engine/model"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
)

var (
	schemaOutputDir string
)

var generateSchemaCmd = &cobra.Command{
	Use:   "generate-schema",
	Short: "Generate JSON Schema files from Go types",
	Long: `Generate JSON Schema files from the Go type definitions.

This command generates JSON Schema files in the specified output directory
(default: schema/). These schemas can be used for validation and documentation.

Example:
  mse generate-schema
  mse generate-schema --output ./my-schemas`,
	RunE: runGenerateSchema,
}

func init() {
	rootCmd.AddCommand(generateSchemaCmd)

	generateSchemaCmd.Flags().StringVarP(&schemaOutputDir, "output", "o", "schema",
		"Output directory for schema files")
}

func runGenerateSchema(cmd *cobra.Command, args []string) error {
	// Ensure schema directory exists
	if err := os.MkdirAll(schemaOutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create schema directory: %w", err)
	}

	// Define schemas to generate
	schemas := []struct {
		filename string
		typ      any
	}{
		{"analysis.schema.json", model.Analysis{}},
		{"market.schema.json", model.Market{}},
		{"segment.schema.json", model.Segment{}},
		{"vendor.schema.json", model.Vendor{}},
		{"vendor_product.schema.json", model.VendorProduct{}},
		{"capability.schema.json", model.Capability{}},
		{"product_category.schema.json", model.ProductCategory{}},
		{"gap_analysis.schema.json", model.GapAnalysis{}},
		{"readiness_score.schema.json", model.ReadinessScore{}},
		{"priority_action.schema.json", model.PriorityAction{}},
		{"capability_score.schema.json", model.CapabilityScore{}},
		{"segment_weight.schema.json", model.SegmentWeight{}},
		{"segment_comparison.schema.json", model.SegmentComparison{}},
		{"vendor_segment_score.schema.json", model.VendorSegmentScore{}},
		{"capability_comparison.schema.json", model.CapabilityComparison{}},
		{"data_source.schema.json", model.DataSource{}},
	}

	reflector := &jsonschema.Reflector{
		DoNotReference:             false,
		ExpandedStruct:             false,
		RequiredFromJSONSchemaTags: true,
	}

	for _, s := range schemas {
		schema := reflector.Reflect(s.typ)

		// Add schema metadata
		schema.ID = jsonschema.ID(fmt.Sprintf(
			"https://github.com/grokify/market-strategy-engine/schema/%s", s.filename))
		schema.Version = "https://json-schema.org/draft/2020-12/schema"

		data, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal schema for %s: %w", s.filename, err)
		}

		path := filepath.Join(schemaOutputDir, s.filename)
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}

		fmt.Printf("Generated %s\n", path)
	}

	fmt.Printf("\nGenerated %d schema files in %s/\n", len(schemas), schemaOutputDir)
	return nil
}
