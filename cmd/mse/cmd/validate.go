package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/market-strategy-engine/model"
	"github.com/grokify/market-strategy-engine/validate"
	"github.com/spf13/cobra"
)

var (
	validateStrict bool
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate an analysis file",
	Long: `Validate an analysis JSON file against the data model.

This command performs comprehensive validation including:
  - Required fields are present
  - No duplicate IDs
  - All references point to existing entities
  - Weights sum to 1.0 per segment
  - Scores are within valid range (0-10)
  - All vendors have scores for all capabilities

Example:
  mse validate examples/security-market.json
  mse validate --strict examples/security-market.json`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().BoolVar(&validateStrict, "strict", false,
		"Treat all issues as errors (exit non-zero on any issue)")
}

func runValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var analysis model.Analysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Run comprehensive validation
	result := validate.Validate(&analysis)

	// Print summary
	fmt.Printf("Validating: %s\n\n", filename)

	summary := analysis.Summary()
	fmt.Printf("Analysis: %s\n", summary.Name)
	fmt.Printf("Market: %s\n", summary.Market)
	fmt.Printf("Vendors: %d\n", summary.VendorCount)
	fmt.Printf("Segments: %d\n", summary.SegmentCount)
	fmt.Printf("Capabilities: %d\n", summary.CapabilityCount)
	fmt.Printf("Scores: %d\n", len(analysis.Scores))
	fmt.Printf("Weights: %d\n", len(analysis.Weights))
	fmt.Println()

	// Print validation results
	if result.IsValid() {
		fmt.Println("✓ Validation passed")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", len(result.Errors))
	for _, e := range result.Errors {
		fmt.Printf("  ✗ %s\n", e.Error())
	}
	fmt.Println()

	if validateStrict || len(result.Errors) > 0 {
		return fmt.Errorf("validation failed with %d issue(s)", len(result.Errors))
	}

	return nil
}
