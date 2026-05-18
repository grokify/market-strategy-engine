package engine

import (
	"testing"

	"github.com/grokify/market-strategy-engine/model"
)

// testAnalysis returns a minimal but complete analysis for testing.
func testAnalysis() *model.Analysis {
	return &model.Analysis{
		ID:            "test-analysis",
		Name:          "Test Analysis",
		FocusVendorID: "vendor-a",
		Market: model.Market{
			ID:   "test-market",
			Name: "Test Market",
			Categories: []model.ProductCategory{
				{
					ID:   "cat-1",
					Name: "Category 1",
					Capabilities: []model.Capability{
						{ID: "cap-1", Name: "Capability 1", GapType: model.GapTypeProduct},
						{ID: "cap-2", Name: "Capability 2", GapType: model.GapTypeStructural},
					},
				},
			},
		},
		Segments: []model.Segment{
			{ID: "seg-1", Name: "Segment 1"},
			{ID: "seg-2", Name: "Segment 2"},
		},
		Vendors: []model.Vendor{
			{ID: "vendor-a", Name: "Vendor A", Color: "#8b5cf6"},
			{ID: "vendor-b", Name: "Vendor B", Color: "#ef4444"},
		},
		Weights: []model.SegmentWeight{
			{SegmentID: "seg-1", CapabilityID: "cap-1", Weight: 0.6},
			{SegmentID: "seg-1", CapabilityID: "cap-2", Weight: 0.4},
			{SegmentID: "seg-2", CapabilityID: "cap-1", Weight: 0.3},
			{SegmentID: "seg-2", CapabilityID: "cap-2", Weight: 0.7},
		},
		Scores: []model.CapabilityScore{
			{VendorID: "vendor-a", CapabilityID: "cap-1", Score: 7.0},
			{VendorID: "vendor-a", CapabilityID: "cap-2", Score: 5.0},
			{VendorID: "vendor-b", CapabilityID: "cap-1", Score: 9.0},
			{VendorID: "vendor-b", CapabilityID: "cap-2", Score: 8.0},
		},
	}
}

func TestComputeGaps(t *testing.T) {
	a := testAnalysis()
	eng := New(a)

	gaps := eng.ComputeGaps()

	// With focus vendor set, should only compute gaps for vendor-a
	// vendor-a has 2 capabilities × 2 segments = 4 gaps
	if len(gaps) != 4 {
		t.Errorf("expected 4 gaps, got %d", len(gaps))
	}

	// Check specific gap calculations
	// For seg-1, cap-1: benchmark=9.0 (vendor-b), vendor-a=7.0, gap=2.0, weight=0.6
	// weighted gap = 2.0 × 0.6 = 1.2
	for _, g := range gaps {
		if g.SegmentID == "seg-1" && g.CapabilityID == "cap-1" {
			if g.Gap != 2.0 {
				t.Errorf("seg-1/cap-1 gap: expected 2.0, got %.1f", g.Gap)
			}
			if g.WeightedGap != 1.2 {
				t.Errorf("seg-1/cap-1 weighted gap: expected 1.2, got %.1f", g.WeightedGap)
			}
			if g.BenchmarkVendorID != "vendor-b" {
				t.Errorf("seg-1/cap-1 benchmark vendor: expected vendor-b, got %s", g.BenchmarkVendorID)
			}
		}
	}
}

func TestComputeGapsAllVendors(t *testing.T) {
	a := testAnalysis()
	a.FocusVendorID = "" // No focus vendor
	eng := New(a)

	gaps := eng.ComputeGaps()

	// Without focus vendor, should compute gaps for all vendors
	// 2 vendors × 2 capabilities × 2 segments = 8 gaps
	if len(gaps) != 8 {
		t.Errorf("expected 8 gaps, got %d", len(gaps))
	}
}

func TestComputeReadiness(t *testing.T) {
	a := testAnalysis()
	eng := New(a)

	gaps := eng.ComputeGaps()
	readiness := eng.ComputeReadiness(gaps)

	// Should have readiness scores for focus vendor in each segment
	if len(readiness) != 2 {
		t.Errorf("expected 2 readiness scores, got %d", len(readiness))
	}

	// Readiness formula: 100 - (totalWeightedGap × 10)
	for _, r := range readiness {
		if r.Score < 0 || r.Score > 100 {
			t.Errorf("readiness score out of range: %.1f", r.Score)
		}
		if r.VendorID != "vendor-a" {
			t.Errorf("expected vendor-a, got %s", r.VendorID)
		}
	}
}

func TestComputePriorities(t *testing.T) {
	a := testAnalysis()
	eng := New(a)

	gaps := eng.ComputeGaps()
	priorities := eng.ComputePriorities(gaps)

	// Should have priorities for capabilities with positive gaps
	if len(priorities) == 0 {
		t.Error("expected at least one priority action")
	}

	// Priorities should be sorted by weighted gap (descending)
	for i := 1; i < len(priorities); i++ {
		if priorities[i].Priority > priorities[i-1].Priority {
			t.Error("priorities not sorted by weighted gap descending")
		}
	}

	// Ranks should be sequential
	for i, p := range priorities {
		if p.Rank != i+1 {
			t.Errorf("expected rank %d, got %d", i+1, p.Rank)
		}
	}
}

func TestComputeComparisons(t *testing.T) {
	a := testAnalysis()
	eng := New(a)

	comparisons := eng.ComputeComparisons()

	// Should have one comparison per segment
	if len(comparisons) != 2 {
		t.Errorf("expected 2 comparisons, got %d", len(comparisons))
	}

	for _, comp := range comparisons {
		// Each comparison should have scores for all vendors
		if len(comp.VendorScores) != 2 {
			t.Errorf("segment %s: expected 2 vendor scores, got %d",
				comp.SegmentID, len(comp.VendorScores))
		}

		// Vendor scores should be ranked
		for i, vs := range comp.VendorScores {
			if vs.Rank != i+1 {
				t.Errorf("vendor %s rank: expected %d, got %d",
					vs.VendorID, i+1, vs.Rank)
			}
		}

		// Focus vendor should be marked
		foundFocus := false
		for _, vs := range comp.VendorScores {
			if vs.IsFocusVendor {
				foundFocus = true
				if vs.VendorID != "vendor-a" {
					t.Error("wrong vendor marked as focus")
				}
			}
		}
		if !foundFocus {
			t.Error("no vendor marked as focus")
		}

		// Should have capability comparisons
		if len(comp.CapabilityComparisons) == 0 {
			t.Error("expected capability comparisons")
		}
	}
}

func TestComputeComparisonsVendorColors(t *testing.T) {
	a := testAnalysis()
	eng := New(a)

	comparisons := eng.ComputeComparisons()

	for _, comp := range comparisons {
		for _, vs := range comp.VendorScores {
			if vs.VendorID == "vendor-a" && vs.VendorColor != "#8b5cf6" {
				t.Errorf("vendor-a color: expected #8b5cf6, got %s", vs.VendorColor)
			}
			if vs.VendorID == "vendor-b" && vs.VendorColor != "#ef4444" {
				t.Errorf("vendor-b color: expected #ef4444, got %s", vs.VendorColor)
			}
		}
	}
}

func TestRun(t *testing.T) {
	a := testAnalysis()
	eng := New(a)

	eng.Run()

	// After Run(), analysis should have computed fields populated
	if len(a.Gaps) == 0 {
		t.Error("gaps not populated after Run()")
	}
	if len(a.Readiness) == 0 {
		t.Error("readiness not populated after Run()")
	}
	if len(a.Priorities) == 0 {
		t.Error("priorities not populated after Run()")
	}
	if len(a.Comparisons) == 0 {
		t.Error("comparisons not populated after Run()")
	}
}

func TestSeverityComputation(t *testing.T) {
	tests := []struct {
		weightedGap float64
		expected    model.GapSeverity
	}{
		{0.0, model.GapSeverityNone},
		{0.4, model.GapSeverityLow},
		{0.5, model.GapSeverityMedium},  // 0.5 is not < 0.5, so medium
		{1.0, model.GapSeverityMedium},
		{1.5, model.GapSeverityHigh},    // 1.5 is not < 1.5, so high
		{2.0, model.GapSeverityHigh},
		{2.5, model.GapSeverityCritical}, // 2.5 is not < 2.5, so critical
	}

	for _, tt := range tests {
		got := model.ComputeSeverity(tt.weightedGap)
		if got != tt.expected {
			t.Errorf("ComputeSeverity(%.1f): expected %s, got %s",
				tt.weightedGap, tt.expected, got)
		}
	}
}

func TestReadinessLevel(t *testing.T) {
	tests := []struct {
		score    float64
		expected model.ReadinessLevel
	}{
		{100, model.ReadinessLevelReady},
		{80, model.ReadinessLevelReady},
		{79, model.ReadinessLevelAdjacent},
		{60, model.ReadinessLevelAdjacent},
		{59, model.ReadinessLevelPartial},
		{40, model.ReadinessLevelPartial},
		{39, model.ReadinessLevelDistant},
		{0, model.ReadinessLevelDistant},
	}

	for _, tt := range tests {
		got := model.GetReadinessLevel(tt.score)
		if got != tt.expected {
			t.Errorf("GetReadinessLevel(%.0f): expected %s, got %s",
				tt.score, tt.expected, got)
		}
	}
}
