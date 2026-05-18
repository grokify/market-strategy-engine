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
	analyzeOutputFile   string
	analyzeOutputFormat string
	analyzeSkipValidate bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Run gap analysis on an analysis file",
	Long: `Run gap analysis on an analysis JSON file.

This command:
  - Computes capability gaps for the focus vendor against market benchmarks
  - Calculates readiness scores for each market segment
  - Generates a prioritized list of strategic actions

The results can be output as a summary to stdout or saved to a file.

Example:
  mse analyze examples/security-market.json
  mse analyze examples/security-market.json --output results.json
  mse analyze examples/security-market.json --format summary
  mse analyze examples/security-market.json --format html --output report.html`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&analyzeOutputFile, "output", "o", "",
		"Output file for full results (JSON)")
	analyzeCmd.Flags().StringVarP(&analyzeOutputFormat, "format", "f", "summary",
		"Output format: summary, json, priorities, or html")
	analyzeCmd.Flags().BoolVar(&analyzeSkipValidate, "skip-validate", false,
		"Skip input validation before analysis")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Load analysis file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var analysis model.Analysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate input data
	if !analyzeSkipValidate {
		result := validate.Validate(&analysis)
		if !result.IsValid() {
			return fmt.Errorf("input validation failed:\n%s", result.Error())
		}
	}

	// Run analysis
	eng := engine.New(&analysis)
	eng.Run()

	// Output results
	switch analyzeOutputFormat {
	case "json":
		return outputJSON(&analysis)
	case "priorities":
		return outputPriorities(&analysis)
	case "html":
		return outputHTML(&analysis)
	case "summary":
		fallthrough
	default:
		return outputSummary(&analysis)
	}
}

func outputSummary(a *model.Analysis) error {
	// Header
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  MARKET STRATEGY ANALYSIS: %s\n", a.Name)
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	// Focus vendor
	focusVendor := a.GetVendor(a.FocusVendorID)
	if focusVendor != nil {
		fmt.Printf("Focus Vendor: %s\n", focusVendor.Name)
	}
	fmt.Printf("Market: %s\n\n", a.Market.Name)

	// Readiness scores
	fmt.Printf("───────────────────────────────────────────────────────────────\n")
	fmt.Printf("  SEGMENT READINESS SCORES\n")
	fmt.Printf("───────────────────────────────────────────────────────────────\n\n")

	for _, r := range a.Readiness {
		segment := a.GetSegment(r.SegmentID)
		segmentName := r.SegmentID
		if segment != nil {
			segmentName = segment.Name
		}

		level := model.GetReadinessLevel(r.Score)
		levelIcon := getReadinessIcon(level)

		fmt.Printf("  %s %s: %.0f/100\n", levelIcon, segmentName, r.Score)
		fmt.Printf("     %s\n\n", r.Interpretation)

		if len(r.CriticalGaps) > 0 {
			fmt.Printf("     Critical gaps: ")
			for i, capID := range r.CriticalGaps {
				cap := a.GetCapability(capID)
				name := capID
				if cap != nil {
					name = cap.Name
				}
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", name)
			}
			fmt.Printf("\n\n")
		}
	}

	// Top priorities
	fmt.Printf("───────────────────────────────────────────────────────────────\n")
	fmt.Printf("  TOP STRATEGIC PRIORITIES\n")
	fmt.Printf("───────────────────────────────────────────────────────────────\n\n")

	maxPriorities := 10
	if len(a.Priorities) < maxPriorities {
		maxPriorities = len(a.Priorities)
	}

	for _, p := range a.Priorities[:maxPriorities] {
		effortIcon := getEffortIcon(p.Effort)
		gapTypeIcon := getGapTypeIcon(p.GapType)

		fmt.Printf("  %d. %s [%s %s]\n", p.Rank, p.CapabilityName, gapTypeIcon, p.GapType)
		fmt.Printf("     Score: %.1f → %.1f (gap: %.1f)\n", p.CurrentScore, p.TargetScore, p.TargetScore-p.CurrentScore)
		fmt.Printf("     Target: %s | Effort: %s %s\n", p.SegmentName, effortIcon, p.Effort)
		fmt.Printf("     %s\n\n", p.Impact)
	}

	// Save to file if requested
	if analyzeOutputFile != "" {
		return saveAnalysis(a, analyzeOutputFile)
	}

	return nil
}

func outputPriorities(a *model.Analysis) error {
	fmt.Printf("Rank | Capability | Gap | Target Segment | Effort | Gap Type\n")
	fmt.Printf("-----|------------|-----|----------------|--------|----------\n")

	for _, p := range a.Priorities {
		fmt.Printf("%4d | %-20s | %.1f | %-14s | %-9s | %s\n",
			p.Rank, truncate(p.CapabilityName, 20),
			p.TargetScore-p.CurrentScore,
			p.SegmentName, p.Effort, p.GapType)
	}

	if analyzeOutputFile != "" {
		return saveAnalysis(a, analyzeOutputFile)
	}

	return nil
}

func outputHTML(a *model.Analysis) error {
	if analyzeOutputFile == "" {
		return fmt.Errorf("HTML output requires --output flag to specify the output file")
	}

	// Ensure directory exists
	dir := filepath.Dir(analyzeOutputFile)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	f, err := os.Create(analyzeOutputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	htmlReport := report.NewHTMLReport(a)
	if err := htmlReport.Write(f); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	fmt.Printf("HTML report saved to: %s\n", analyzeOutputFile)
	return nil
}

func outputJSON(a *model.Analysis) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func saveAnalysis(a *model.Analysis, filename string) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	// Ensure directory exists
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

func getReadinessIcon(level model.ReadinessLevel) string {
	switch level {
	case model.ReadinessLevelReady:
		return "●"
	case model.ReadinessLevelAdjacent:
		return "◐"
	case model.ReadinessLevelPartial:
		return "○"
	case model.ReadinessLevelDistant:
		return "◌"
	default:
		return "?"
	}
}

func getEffortIcon(effort model.EffortLevel) string {
	switch effort {
	case model.EffortLevelLow:
		return "▪"
	case model.EffortLevelMedium:
		return "▪▪"
	case model.EffortLevelHigh:
		return "▪▪▪"
	case model.EffortLevelVeryHigh:
		return "▪▪▪▪"
	default:
		return "?"
	}
}

func getGapTypeIcon(gapType model.GapType) string {
	switch gapType {
	case model.GapTypeProduct:
		return "⚙"
	case model.GapTypeStructural:
		return "⚡"
	case model.GapTypePerception:
		return "◉"
	default:
		return "?"
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
