package model

// GapAnalysis represents a computed capability gap for a vendor
// relative to the market benchmark for a specific segment.
//
// The core formula is:
//
//	WeightedGap = (BenchmarkScore - VendorScore) × SegmentWeight
//
// Higher weighted gaps indicate more critical areas for improvement.
type GapAnalysis struct {
	// VendorID references the vendor being analyzed.
	VendorID string `json:"vendorId"`

	// CapabilityID references the capability with the gap.
	CapabilityID string `json:"capabilityId"`

	// SegmentID references the target segment for this analysis.
	SegmentID string `json:"segmentId"`

	// BenchmarkScore is the market leader's score (or target score).
	BenchmarkScore float64 `json:"benchmarkScore"`

	// BenchmarkVendorID is the vendor used as benchmark (optional).
	BenchmarkVendorID string `json:"benchmarkVendorId,omitempty"`

	// VendorScore is the analyzed vendor's score.
	VendorScore float64 `json:"vendorScore"`

	// Gap is the raw gap (BenchmarkScore - VendorScore).
	// Positive values indicate the vendor is behind the benchmark.
	Gap float64 `json:"gap"`

	// SegmentWeight is the importance of this capability for the segment.
	SegmentWeight float64 `json:"segmentWeight"`

	// WeightedGap is Gap × SegmentWeight.
	// This is the primary metric for prioritization.
	WeightedGap float64 `json:"weightedGap"`

	// GapType categorizes the nature of this gap.
	GapType GapType `json:"gapType,omitempty"`

	// Severity categorizes how critical this gap is.
	Severity GapSeverity `json:"severity,omitempty"`
}

// GapSeverity indicates how critical a gap is for market success.
type GapSeverity string

const (
	// GapSeverityCritical indicates a blocking gap that prevents segment entry.
	GapSeverityCritical GapSeverity = "critical"

	// GapSeverityHigh indicates a significant competitive disadvantage.
	GapSeverityHigh GapSeverity = "high"

	// GapSeverityMedium indicates a notable but manageable gap.
	GapSeverityMedium GapSeverity = "medium"

	// GapSeverityLow indicates a minor gap with limited impact.
	GapSeverityLow GapSeverity = "low"

	// GapSeverityNone indicates no gap (vendor meets or exceeds benchmark).
	GapSeverityNone GapSeverity = "none"
)

// ComputeSeverity determines gap severity based on weighted gap score.
func ComputeSeverity(weightedGap float64) GapSeverity {
	switch {
	case weightedGap <= 0:
		return GapSeverityNone
	case weightedGap < 0.5:
		return GapSeverityLow
	case weightedGap < 1.5:
		return GapSeverityMedium
	case weightedGap < 2.5:
		return GapSeverityHigh
	default:
		return GapSeverityCritical
	}
}

// ReadinessScore represents a vendor's overall readiness to compete
// in a specific market segment. Computed from capability gaps.
type ReadinessScore struct {
	// VendorID references the vendor.
	VendorID string `json:"vendorId"`

	// SegmentID references the target segment.
	SegmentID string `json:"segmentId"`

	// Score is the readiness percentage (0-100).
	// 100 = fully ready, meets all benchmarks.
	// 0 = not ready, significant gaps across all capabilities.
	Score float64 `json:"score"`

	// MaxPossible is the maximum possible score (typically 100).
	MaxPossible float64 `json:"maxPossible"`

	// TotalWeightedGap is the sum of all weighted gaps.
	TotalWeightedGap float64 `json:"totalWeightedGap"`

	// CriticalGaps lists gaps with critical severity.
	CriticalGaps []string `json:"criticalGaps,omitempty"`

	// TopGaps lists the most impactful gaps (by weighted gap).
	TopGaps []GapAnalysis `json:"topGaps,omitempty"`

	// Interpretation provides a human-readable assessment.
	Interpretation string `json:"interpretation,omitempty"`
}

// ReadinessLevel categorizes readiness into actionable tiers.
type ReadinessLevel string

const (
	// ReadinessLevelReady indicates the vendor can compete effectively.
	ReadinessLevelReady ReadinessLevel = "ready"

	// ReadinessLevelAdjacent indicates minor improvements needed.
	ReadinessLevelAdjacent ReadinessLevel = "adjacent"

	// ReadinessLevelPartial indicates significant work required.
	ReadinessLevelPartial ReadinessLevel = "partial"

	// ReadinessLevelDistant indicates major investment required.
	ReadinessLevelDistant ReadinessLevel = "distant"
)

// GetReadinessLevel returns the readiness level for a score.
func GetReadinessLevel(score float64) ReadinessLevel {
	switch {
	case score >= 80:
		return ReadinessLevelReady
	case score >= 60:
		return ReadinessLevelAdjacent
	case score >= 40:
		return ReadinessLevelPartial
	default:
		return ReadinessLevelDistant
	}
}

// PriorityAction represents a prioritized strategic action
// derived from gap analysis.
type PriorityAction struct {
	// Rank is the priority order (1 = highest priority).
	Rank int `json:"rank"`

	// CapabilityID references the capability to improve.
	CapabilityID string `json:"capabilityId"`

	// CapabilityName is included for convenience.
	CapabilityName string `json:"capabilityName,omitempty"`

	// SegmentID references the target segment this unlocks.
	SegmentID string `json:"segmentId"`

	// SegmentName is included for convenience.
	SegmentName string `json:"segmentName,omitempty"`

	// Priority is the computed priority score (higher = more important).
	Priority float64 `json:"priority"`

	// GapType indicates how this gap should be addressed.
	GapType GapType `json:"gapType,omitempty"`

	// CurrentScore is the vendor's current score.
	CurrentScore float64 `json:"currentScore"`

	// TargetScore is the benchmark/target score.
	TargetScore float64 `json:"targetScore"`

	// Impact describes what improving this capability enables.
	Impact string `json:"impact,omitempty"`

	// Effort provides a rough estimate of implementation effort.
	Effort EffortLevel `json:"effort,omitempty"`
}

// EffortLevel estimates the effort required to close a gap.
type EffortLevel string

const (
	// EffortLevelLow indicates weeks of effort.
	EffortLevelLow EffortLevel = "low"

	// EffortLevelMedium indicates months of effort.
	EffortLevelMedium EffortLevel = "medium"

	// EffortLevelHigh indicates quarters of effort.
	EffortLevelHigh EffortLevel = "high"

	// EffortLevelVeryHigh indicates years or major platform changes.
	EffortLevelVeryHigh EffortLevel = "very_high"
)

// GetEffortLevel estimates effort based on gap type and gap size.
func GetEffortLevel(gapType GapType, gap float64) EffortLevel {
	switch gapType {
	case GapTypePerception:
		if gap < 2 {
			return EffortLevelMedium
		}
		return EffortLevelHigh

	case GapTypeProduct:
		if gap < 2 {
			return EffortLevelLow
		} else if gap < 4 {
			return EffortLevelMedium
		}
		return EffortLevelHigh

	case GapTypeStructural:
		if gap < 2 {
			return EffortLevelMedium
		} else if gap < 4 {
			return EffortLevelHigh
		}
		return EffortLevelVeryHigh

	default:
		return EffortLevelMedium
	}
}
