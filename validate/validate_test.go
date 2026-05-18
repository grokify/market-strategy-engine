package validate

import (
	"strings"
	"testing"

	"github.com/grokify/market-strategy-engine/model"
)

// validAnalysis returns a minimal but valid analysis for testing.
func validAnalysis() *model.Analysis {
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
						{ID: "cap-1", Name: "Capability 1"},
						{ID: "cap-2", Name: "Capability 2"},
					},
				},
			},
		},
		Segments: []model.Segment{
			{ID: "seg-1", Name: "Segment 1"},
		},
		Vendors: []model.Vendor{
			{ID: "vendor-a", Name: "Vendor A"},
			{ID: "vendor-b", Name: "Vendor B"},
		},
		Weights: []model.SegmentWeight{
			{SegmentID: "seg-1", CapabilityID: "cap-1", Weight: 0.6},
			{SegmentID: "seg-1", CapabilityID: "cap-2", Weight: 0.4},
		},
		Scores: []model.CapabilityScore{
			{VendorID: "vendor-a", CapabilityID: "cap-1", Score: 7.0},
			{VendorID: "vendor-a", CapabilityID: "cap-2", Score: 5.0},
			{VendorID: "vendor-b", CapabilityID: "cap-1", Score: 9.0},
			{VendorID: "vendor-b", CapabilityID: "cap-2", Score: 8.0},
		},
	}
}

func TestValidAnalysis(t *testing.T) {
	a := validAnalysis()
	result := Validate(a)

	if !result.IsValid() {
		t.Errorf("expected valid analysis, got errors: %s", result.Error())
	}
}

func TestMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*model.Analysis)
		errText string
	}{
		{
			name:    "missing id",
			modify:  func(a *model.Analysis) { a.ID = "" },
			errText: "id",
		},
		{
			name:    "missing name",
			modify:  func(a *model.Analysis) { a.Name = "" },
			errText: "name",
		},
		{
			name:    "missing market.id",
			modify:  func(a *model.Analysis) { a.Market.ID = "" },
			errText: "market.id",
		},
		{
			name:    "missing market.name",
			modify:  func(a *model.Analysis) { a.Market.Name = "" },
			errText: "market.name",
		},
		{
			name:    "no segments",
			modify:  func(a *model.Analysis) { a.Segments = nil },
			errText: "segments",
		},
		{
			name:    "no vendors",
			modify:  func(a *model.Analysis) { a.Vendors = nil },
			errText: "vendors",
		},
		{
			name:    "no categories",
			modify:  func(a *model.Analysis) { a.Market.Categories = nil },
			errText: "categories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := validAnalysis()
			tt.modify(a)
			result := Validate(a)

			if result.IsValid() {
				t.Error("expected validation error")
			}

			found := false
			for _, e := range result.Errors {
				if strings.Contains(e.Field, tt.errText) || strings.Contains(e.Message, tt.errText) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error about %q, got: %v", tt.errText, result.Errors)
			}
		})
	}
}

func TestDuplicateIDs(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*model.Analysis)
	}{
		{
			name: "duplicate segment ID",
			modify: func(a *model.Analysis) {
				a.Segments = append(a.Segments, model.Segment{ID: "seg-1", Name: "Duplicate"})
			},
		},
		{
			name: "duplicate vendor ID",
			modify: func(a *model.Analysis) {
				a.Vendors = append(a.Vendors, model.Vendor{ID: "vendor-a", Name: "Duplicate"})
			},
		},
		{
			name: "duplicate capability ID",
			modify: func(a *model.Analysis) {
				a.Market.Categories[0].Capabilities = append(
					a.Market.Categories[0].Capabilities,
					model.Capability{ID: "cap-1", Name: "Duplicate"},
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := validAnalysis()
			tt.modify(a)
			result := Validate(a)

			if result.IsValid() {
				t.Error("expected validation error for duplicate ID")
			}

			found := false
			for _, e := range result.Errors {
				if strings.Contains(e.Message, "duplicate") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected duplicate error, got: %v", result.Errors)
			}
		})
	}
}

func TestInvalidReferences(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*model.Analysis)
		errMsg string
	}{
		{
			name: "invalid focusVendorId",
			modify: func(a *model.Analysis) {
				a.FocusVendorID = "nonexistent"
			},
			errMsg: "unknown vendor",
		},
		{
			name: "weight references invalid segment",
			modify: func(a *model.Analysis) {
				a.Weights = append(a.Weights, model.SegmentWeight{
					SegmentID:    "nonexistent",
					CapabilityID: "cap-1",
					Weight:       0.5,
				})
			},
			errMsg: "unknown segment",
		},
		{
			name: "weight references invalid capability",
			modify: func(a *model.Analysis) {
				a.Weights = append(a.Weights, model.SegmentWeight{
					SegmentID:    "seg-1",
					CapabilityID: "nonexistent",
					Weight:       0.5,
				})
			},
			errMsg: "unknown capability",
		},
		{
			name: "score references invalid vendor",
			modify: func(a *model.Analysis) {
				a.Scores = append(a.Scores, model.CapabilityScore{
					VendorID:     "nonexistent",
					CapabilityID: "cap-1",
					Score:        5.0,
				})
			},
			errMsg: "unknown vendor",
		},
		{
			name: "score references invalid capability",
			modify: func(a *model.Analysis) {
				a.Scores = append(a.Scores, model.CapabilityScore{
					VendorID:     "vendor-a",
					CapabilityID: "nonexistent",
					Score:        5.0,
				})
			},
			errMsg: "unknown capability",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := validAnalysis()
			tt.modify(a)
			result := Validate(a)

			if result.IsValid() {
				t.Error("expected validation error")
			}

			found := false
			for _, e := range result.Errors {
				if strings.Contains(e.Message, tt.errMsg) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected %q error, got: %v", tt.errMsg, result.Errors)
			}
		})
	}
}

func TestWeightValidation(t *testing.T) {
	t.Run("weights don't sum to 1.0", func(t *testing.T) {
		a := validAnalysis()
		// Change weights to sum to 0.9
		a.Weights[0].Weight = 0.5
		a.Weights[1].Weight = 0.4

		result := Validate(a)

		if result.IsValid() {
			t.Error("expected validation error for incorrect weight sum")
		}

		found := false
		for _, e := range result.Errors {
			if strings.Contains(e.Message, "sum to 1.0") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected weight sum error, got: %v", result.Errors)
		}
	})

	t.Run("weight out of range", func(t *testing.T) {
		a := validAnalysis()
		a.Weights[0].Weight = 1.5 // Invalid
		a.Weights[1].Weight = -0.5

		result := Validate(a)

		if result.IsValid() {
			t.Error("expected validation error for out-of-range weight")
		}
	})

	t.Run("missing weight for capability", func(t *testing.T) {
		a := validAnalysis()
		// Remove one weight
		a.Weights = a.Weights[:1]

		result := Validate(a)

		if result.IsValid() {
			t.Error("expected validation error for missing weight")
		}

		found := false
		for _, e := range result.Errors {
			if strings.Contains(e.Message, "missing weight") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected missing weight error, got: %v", result.Errors)
		}
	})
}

func TestScoreValidation(t *testing.T) {
	t.Run("score out of range (too high)", func(t *testing.T) {
		a := validAnalysis()
		a.Scores[0].Score = 11.0

		result := Validate(a)

		if result.IsValid() {
			t.Error("expected validation error for out-of-range score")
		}

		found := false
		for _, e := range result.Errors {
			if strings.Contains(e.Message, "between 0 and 10") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected score range error, got: %v", result.Errors)
		}
	})

	t.Run("score out of range (negative)", func(t *testing.T) {
		a := validAnalysis()
		a.Scores[0].Score = -1.0

		result := Validate(a)

		if result.IsValid() {
			t.Error("expected validation error for negative score")
		}
	})

	t.Run("missing score for vendor/capability", func(t *testing.T) {
		a := validAnalysis()
		// Remove scores for vendor-b
		a.Scores = a.Scores[:2]

		result := Validate(a)

		if result.IsValid() {
			t.Error("expected validation error for missing scores")
		}

		found := false
		for _, e := range result.Errors {
			// Could be "missing score" or "no scores defined"
			if strings.Contains(e.Message, "missing score") || strings.Contains(e.Message, "no scores defined") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected missing score error, got: %v", result.Errors)
		}
	})
}

func TestValidationResultError(t *testing.T) {
	result := &ValidationResult{}
	result.Add("field1", "error 1")
	result.Add("field2", "error 2")

	errStr := result.Error()

	if !strings.Contains(errStr, "2 error(s)") {
		t.Errorf("expected error count in message, got: %s", errStr)
	}
	if !strings.Contains(errStr, "field1") || !strings.Contains(errStr, "error 1") {
		t.Errorf("expected field1 error in message, got: %s", errStr)
	}
	if !strings.Contains(errStr, "field2") || !strings.Contains(errStr, "error 2") {
		t.Errorf("expected field2 error in message, got: %s", errStr)
	}
}

func TestValidationResultIsValid(t *testing.T) {
	result := &ValidationResult{}

	if !result.IsValid() {
		t.Error("empty result should be valid")
	}

	result.Add("field", "error")

	if result.IsValid() {
		t.Error("result with errors should not be valid")
	}
}
