// Package engine implements the gap analysis and readiness scoring logic.
package engine

import (
	"fmt"
	"sort"

	"github.com/grokify/market-strategy-engine/model"
)

// Engine computes gap analysis and readiness scores from an Analysis.
type Engine struct {
	analysis *model.Analysis
}

// New creates a new Engine from an Analysis.
func New(a *model.Analysis) *Engine {
	return &Engine{analysis: a}
}

// ComputeGaps calculates capability gaps for the focus vendor against benchmarks.
// If no focus vendor is set, it computes gaps for all vendors.
func (e *Engine) ComputeGaps() []model.GapAnalysis {
	var gaps []model.GapAnalysis

	// Build score lookup: vendorID -> capabilityID -> score
	scoreMap := make(map[string]map[string]float64)
	for _, s := range e.analysis.Scores {
		if scoreMap[s.VendorID] == nil {
			scoreMap[s.VendorID] = make(map[string]float64)
		}
		scoreMap[s.VendorID][s.CapabilityID] = s.Score
	}

	// Find benchmark (max score) for each capability
	benchmarks := make(map[string]struct {
		score    float64
		vendorID string
	})
	for vendorID, caps := range scoreMap {
		for capID, score := range caps {
			if b, ok := benchmarks[capID]; !ok || score > b.score {
				benchmarks[capID] = struct {
					score    float64
					vendorID string
				}{score, vendorID}
			}
		}
	}

	// Get all capabilities
	var capabilities []model.Capability
	for _, cat := range e.analysis.Market.Categories {
		capabilities = append(capabilities, cat.Capabilities...)
	}

	// Determine which vendors to analyze
	vendorsToAnalyze := []string{}
	if e.analysis.FocusVendorID != "" {
		vendorsToAnalyze = append(vendorsToAnalyze, e.analysis.FocusVendorID)
	} else {
		for _, v := range e.analysis.Vendors {
			vendorsToAnalyze = append(vendorsToAnalyze, v.ID)
		}
	}

	// Compute gaps for each vendor, segment, and capability
	for _, vendorID := range vendorsToAnalyze {
		vendorScores := scoreMap[vendorID]
		if vendorScores == nil {
			continue
		}

		for _, segment := range e.analysis.Segments {
			for _, cap := range capabilities {
				vendorScore := vendorScores[cap.ID]
				benchmark := benchmarks[cap.ID]
				weight := e.analysis.GetWeight(segment.ID, cap.ID)

				gap := benchmark.score - vendorScore
				weightedGap := gap * weight

				gapAnalysis := model.GapAnalysis{
					VendorID:          vendorID,
					CapabilityID:      cap.ID,
					SegmentID:         segment.ID,
					BenchmarkScore:    benchmark.score,
					BenchmarkVendorID: benchmark.vendorID,
					VendorScore:       vendorScore,
					Gap:               gap,
					SegmentWeight:     weight,
					WeightedGap:       weightedGap,
					GapType:           cap.GapType,
					Severity:          model.ComputeSeverity(weightedGap),
				}

				gaps = append(gaps, gapAnalysis)
			}
		}
	}

	return gaps
}

// ComputeReadiness calculates readiness scores for the focus vendor per segment.
func (e *Engine) ComputeReadiness(gaps []model.GapAnalysis) []model.ReadinessScore {
	// Group gaps by vendor and segment
	type key struct {
		vendorID  string
		segmentID string
	}
	grouped := make(map[key][]model.GapAnalysis)
	for _, g := range gaps {
		k := key{g.VendorID, g.SegmentID}
		grouped[k] = append(grouped[k], g)
	}

	var scores []model.ReadinessScore
	for k, gapList := range grouped {
		var totalWeightedGap float64
		var criticalGaps []string
		var maxWeightedGap float64

		for _, g := range gapList {
			totalWeightedGap += g.WeightedGap
			if g.WeightedGap > maxWeightedGap {
				maxWeightedGap = g.WeightedGap
			}
			if g.Severity == model.GapSeverityCritical {
				criticalGaps = append(criticalGaps, g.CapabilityID)
			}
		}

		// Readiness score: 100 - (total weighted gap * 10)
		// Clamped to 0-100
		readiness := 100.0 - (totalWeightedGap * 10.0)
		if readiness < 0 {
			readiness = 0
		}
		if readiness > 100 {
			readiness = 100
		}

		// Sort gaps by weighted gap descending for top gaps
		sort.Slice(gapList, func(i, j int) bool {
			return gapList[i].WeightedGap > gapList[j].WeightedGap
		})

		// Top 5 gaps
		topGaps := gapList
		if len(topGaps) > 5 {
			topGaps = topGaps[:5]
		}

		// Generate interpretation
		level := model.GetReadinessLevel(readiness)
		interpretation := e.interpretReadiness(level, k.segmentID, len(criticalGaps))

		scores = append(scores, model.ReadinessScore{
			VendorID:         k.vendorID,
			SegmentID:        k.segmentID,
			Score:            readiness,
			MaxPossible:      100,
			TotalWeightedGap: totalWeightedGap,
			CriticalGaps:     criticalGaps,
			TopGaps:          topGaps,
			Interpretation:   interpretation,
		})
	}

	// Sort by segment for consistent output
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].SegmentID < scores[j].SegmentID
	})

	return scores
}

func (e *Engine) interpretReadiness(level model.ReadinessLevel, segmentID string, criticalCount int) string {
	segment := e.analysis.GetSegment(segmentID)
	segmentName := segmentID
	if segment != nil {
		segmentName = segment.Name
	}

	switch level {
	case model.ReadinessLevelReady:
		return fmt.Sprintf("Ready to compete effectively in %s segment", segmentName)
	case model.ReadinessLevelAdjacent:
		return fmt.Sprintf("Adjacent to %s segment; minor improvements needed", segmentName)
	case model.ReadinessLevelPartial:
		if criticalCount > 0 {
			return fmt.Sprintf("Partial fit for %s segment; %d critical gap(s) must be addressed", segmentName, criticalCount)
		}
		return fmt.Sprintf("Partial fit for %s segment; significant work required", segmentName)
	case model.ReadinessLevelDistant:
		return fmt.Sprintf("Distant from %s segment; major investment required", segmentName)
	default:
		return ""
	}
}

// ComputePriorities generates a prioritized list of actions from gaps.
func (e *Engine) ComputePriorities(gaps []model.GapAnalysis) []model.PriorityAction {
	// Filter to positive gaps only and dedupe by capability (take highest weighted gap)
	capGaps := make(map[string]model.GapAnalysis)
	for _, g := range gaps {
		if g.WeightedGap <= 0 {
			continue
		}
		existing, ok := capGaps[g.CapabilityID]
		if !ok || g.WeightedGap > existing.WeightedGap {
			capGaps[g.CapabilityID] = g
		}
	}

	// Convert to slice and sort by weighted gap descending
	var sorted []model.GapAnalysis
	for _, g := range capGaps {
		sorted = append(sorted, g)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].WeightedGap > sorted[j].WeightedGap
	})

	// Generate priority actions
	var actions []model.PriorityAction
	for i, g := range sorted {
		cap := e.analysis.GetCapability(g.CapabilityID)
		capName := g.CapabilityID
		if cap != nil {
			capName = cap.Name
		}

		segment := e.analysis.GetSegment(g.SegmentID)
		segmentName := g.SegmentID
		if segment != nil {
			segmentName = segment.Name
		}

		action := model.PriorityAction{
			Rank:           i + 1,
			CapabilityID:   g.CapabilityID,
			CapabilityName: capName,
			SegmentID:      g.SegmentID,
			SegmentName:    segmentName,
			Priority:       g.WeightedGap,
			GapType:        g.GapType,
			CurrentScore:   g.VendorScore,
			TargetScore:    g.BenchmarkScore,
			Impact:         e.describeImpact(g),
			Effort:         model.GetEffortLevel(g.GapType, g.Gap),
		}
		actions = append(actions, action)
	}

	return actions
}

func (e *Engine) describeImpact(g model.GapAnalysis) string {
	segment := e.analysis.GetSegment(g.SegmentID)
	segmentName := g.SegmentID
	if segment != nil {
		segmentName = segment.Name
	}

	cap := e.analysis.GetCapability(g.CapabilityID)
	capName := g.CapabilityID
	if cap != nil {
		capName = cap.Name
	}

	switch g.Severity {
	case model.GapSeverityCritical:
		return fmt.Sprintf("Critical blocker for %s entry; improving %s is essential", segmentName, capName)
	case model.GapSeverityHigh:
		return fmt.Sprintf("Significant competitive disadvantage in %s; high priority", segmentName)
	case model.GapSeverityMedium:
		return fmt.Sprintf("Notable gap affecting %s competitiveness", segmentName)
	case model.GapSeverityLow:
		return fmt.Sprintf("Minor improvement opportunity for %s", segmentName)
	default:
		return ""
	}
}

// ComputeComparisons generates segment-level competitive comparisons.
func (e *Engine) ComputeComparisons() []model.SegmentComparison {
	var comparisons []model.SegmentComparison

	// Build score lookup: vendorID -> capabilityID -> score
	scoreMap := make(map[string]map[string]float64)
	for _, s := range e.analysis.Scores {
		if scoreMap[s.VendorID] == nil {
			scoreMap[s.VendorID] = make(map[string]float64)
		}
		scoreMap[s.VendorID][s.CapabilityID] = s.Score
	}

	// Get all capabilities
	var capabilities []model.Capability
	for _, cat := range e.analysis.Market.Categories {
		capabilities = append(capabilities, cat.Capabilities...)
	}

	for _, segment := range e.analysis.Segments {
		comparison := model.SegmentComparison{
			SegmentID:     segment.ID,
			SegmentName:   segment.Name,
			FocusVendorID: e.analysis.FocusVendorID,
		}

		// Compute per-vendor weighted scores
		var vendorScores []model.VendorSegmentScore
		for _, vendor := range e.analysis.Vendors {
			vendorCaps := scoreMap[vendor.ID]
			if vendorCaps == nil {
				continue
			}

			var weightedScore float64
			for _, cap := range capabilities {
				weight := e.analysis.GetWeight(segment.ID, cap.ID)
				score := vendorCaps[cap.ID]
				weightedScore += score * weight
			}

			// Max possible is 10 (max score) × 1.0 (sum of weights) = 10
			maxPossible := 10.0
			normalized := (weightedScore / maxPossible) * 100

			vendorScores = append(vendorScores, model.VendorSegmentScore{
				VendorID:         vendor.ID,
				VendorName:       vendor.Name,
				VendorColor:      vendor.Color,
				WeightedScore:    weightedScore,
				MaxPossibleScore: maxPossible,
				NormalizedScore:  normalized,
				IsFocusVendor:    vendor.ID == e.analysis.FocusVendorID,
			})
		}

		// Sort by weighted score descending and assign ranks
		sort.Slice(vendorScores, func(i, j int) bool {
			return vendorScores[i].WeightedScore > vendorScores[j].WeightedScore
		})
		for i := range vendorScores {
			vendorScores[i].Rank = i + 1
		}

		comparison.VendorScores = vendorScores

		// Compute per-capability comparisons
		var capComparisons []model.CapabilityComparison
		for _, cap := range capabilities {
			weight := e.analysis.GetWeight(segment.ID, cap.ID)
			if weight == 0 {
				continue
			}

			capComp := model.CapabilityComparison{
				CapabilityID:   cap.ID,
				CapabilityName: cap.Name,
				SegmentWeight:  weight,
				VendorScores:   make(map[string]float64),
			}

			for _, vendor := range e.analysis.Vendors {
				if vendorCaps := scoreMap[vendor.ID]; vendorCaps != nil {
					score := vendorCaps[cap.ID]
					capComp.VendorScores[vendor.ID] = score
					if score > capComp.BestScore {
						capComp.BestScore = score
						capComp.BestVendorID = vendor.ID
					}
				}
			}

			capComparisons = append(capComparisons, capComp)
		}

		// Sort by weight descending
		sort.Slice(capComparisons, func(i, j int) bool {
			return capComparisons[i].SegmentWeight > capComparisons[j].SegmentWeight
		})

		comparison.CapabilityComparisons = capComparisons
		comparisons = append(comparisons, comparison)
	}

	return comparisons
}

// Run executes the full analysis and updates the Analysis struct.
func (e *Engine) Run() {
	gaps := e.ComputeGaps()
	readiness := e.ComputeReadiness(gaps)
	priorities := e.ComputePriorities(gaps)
	comparisons := e.ComputeComparisons()

	e.analysis.Gaps = gaps
	e.analysis.Readiness = readiness
	e.analysis.Priorities = priorities
	e.analysis.Comparisons = comparisons
}
