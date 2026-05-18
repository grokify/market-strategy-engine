package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/market-strategy-engine/engine"
	"github.com/grokify/market-strategy-engine/model"
	"github.com/grokify/market-strategy-engine/report"
	"github.com/grokify/market-strategy-engine/validate"
	"github.com/spf13/cobra"
)

var (
	multimarketName         string
	multimarketFocusVendor  string
	multimarketOutputFile   string
	multimarketOutputFormat string
	multimarketSkipValidate bool
)

var multimarketCmd = &cobra.Command{
	Use:   "multimarket [files...]",
	Short: "Analyze a vendor across multiple markets",
	Long: `Run competitive analysis across multiple markets.

This command combines multiple single-market analysis files to produce
a cross-market view showing:
  - Overall market performance comparison
  - Segment-level performance across markets
  - Strengths and weaknesses by market and segment

Example:
  mse multimarket --focus huntress --name "Huntress Portfolio" \
      edr.json mdr.json siem.json --format html --output portfolio.html`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMultimarket,
}

func init() {
	rootCmd.AddCommand(multimarketCmd)

	multimarketCmd.Flags().StringVar(&multimarketName, "name", "Multi-Market Analysis",
		"Name for the combined analysis")
	multimarketCmd.Flags().StringVar(&multimarketFocusVendor, "focus", "",
		"Focus vendor ID (required)")
	multimarketCmd.Flags().StringVarP(&multimarketOutputFile, "output", "o", "",
		"Output file for results")
	multimarketCmd.Flags().StringVarP(&multimarketOutputFormat, "format", "f", "summary",
		"Output format: summary, json, or html")
	multimarketCmd.Flags().BoolVar(&multimarketSkipValidate, "skip-validate", false,
		"Skip input validation")

	_ = multimarketCmd.MarkFlagRequired("focus")
}

func runMultimarket(cmd *cobra.Command, args []string) error {
	// Load all analysis files
	analyses := make([]*model.Analysis, 0, len(args))

	for _, filename := range args {
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filename, err)
		}

		var analysis model.Analysis
		if err := json.Unmarshal(data, &analysis); err != nil {
			return fmt.Errorf("failed to parse %s: %w", filename, err)
		}

		// Validate if not skipped
		if !multimarketSkipValidate {
			result := validate.Validate(&analysis)
			if !result.IsValid() {
				return fmt.Errorf("validation failed for %s:\n%s", filename, result.Error())
			}
		}

		analyses = append(analyses, &analysis)
	}

	// Run multi-market analysis
	eng := engine.NewMultiMarket(analyses)
	mma := eng.ComputeMultiMarket("multimarket", multimarketName, multimarketFocusVendor)

	// Output results
	switch multimarketOutputFormat {
	case "json":
		return outputMultimarketJSON(mma)
	case "html":
		return outputMultimarketHTML(mma)
	case "summary":
		fallthrough
	default:
		return outputMultimarketSummary(mma)
	}
}

func outputMultimarketSummary(mma *model.MultiMarketAnalysis) error {
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  MULTI-MARKET ANALYSIS: %s\n", mma.Name)
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	fmt.Printf("Focus Vendor: %s\n", mma.FocusVendorName)
	fmt.Printf("Markets Analyzed: %d\n\n", len(mma.Markets))

	if mma.CrossMarketComparison != nil {
		cmc := mma.CrossMarketComparison

		fmt.Printf("───────────────────────────────────────────────────────────────\n")
		fmt.Printf("  CROSS-MARKET INSIGHTS\n")
		fmt.Printf("───────────────────────────────────────────────────────────────\n\n")

		fmt.Printf("  Strongest Market:  %s\n", cmc.StrongestMarket)
		fmt.Printf("  Weakest Market:    %s\n", cmc.WeakestMarket)
		fmt.Printf("  Strongest Segment: %s\n", cmc.StrongestSegment)
		fmt.Printf("  Weakest Segment:   %s\n\n", cmc.WeakestSegment)

		fmt.Printf("───────────────────────────────────────────────────────────────\n")
		fmt.Printf("  MARKET PERFORMANCE\n")
		fmt.Printf("───────────────────────────────────────────────────────────────\n\n")

		for _, ms := range cmc.MarketScores {
			level := model.GetReadinessLevel(ms.OverallScore)
			fmt.Printf("  %-20s  #%d of %d  %.0f%%  [%s]\n",
				ms.MarketName, ms.Rank, ms.TotalVendors, ms.OverallScore, level)
		}
		fmt.Println()

		fmt.Printf("───────────────────────────────────────────────────────────────\n")
		fmt.Printf("  SEGMENT PERFORMANCE (ACROSS ALL MARKETS)\n")
		fmt.Printf("───────────────────────────────────────────────────────────────\n\n")

		for _, ss := range cmc.SegmentScores {
			fmt.Printf("  %-15s  Avg: %.0f%%  (Best: %s %.0f%%, Worst: %s %.0f%%)\n",
				ss.SegmentName, ss.AverageScore,
				ss.BestMarket, ss.BestMarketScore,
				ss.WorstMarket, ss.WorstMarketScore)
		}
		fmt.Println()
	}

	if multimarketOutputFile != "" {
		return saveMultimarketJSON(mma, multimarketOutputFile)
	}

	return nil
}

func outputMultimarketJSON(mma *model.MultiMarketAnalysis) error {
	data, err := json.MarshalIndent(mma, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputMultimarketHTML(mma *model.MultiMarketAnalysis) error {
	if multimarketOutputFile == "" {
		return fmt.Errorf("HTML output requires --output flag to specify the output file")
	}

	// Ensure directory exists
	dir := filepath.Dir(multimarketOutputFile)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	f, err := os.Create(multimarketOutputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	htmlReport := report.NewMultiMarketHTMLReport(mma)
	if err := htmlReport.Write(f); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	fmt.Printf("HTML report saved to: %s\n", multimarketOutputFile)
	return nil
}

func saveMultimarketJSON(mma *model.MultiMarketAnalysis, filename string) error {
	data, err := json.MarshalIndent(mma, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("\nResults saved to: %s\n", filename)
	return nil
}
