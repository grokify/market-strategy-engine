// Package validate provides validation for analysis data.
package validate

import (
	"fmt"
	"math"
	"strings"

	"github.com/grokify/market-strategy-engine/model"
)

// ValidationError represents a validation issue.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidationResult contains all validation errors found.
type ValidationResult struct {
	Errors []ValidationError
}

// IsValid returns true if there are no validation errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// Add adds a validation error.
func (r *ValidationResult) Add(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// Error returns a combined error message.
func (r *ValidationResult) Error() string {
	if r.IsValid() {
		return ""
	}
	var msgs []string
	for _, e := range r.Errors {
		msgs = append(msgs, e.Error())
	}
	return fmt.Sprintf("validation failed with %d error(s):\n  - %s",
		len(r.Errors), strings.Join(msgs, "\n  - "))
}

// Validate performs comprehensive validation on an Analysis.
func Validate(a *model.Analysis) *ValidationResult {
	result := &ValidationResult{}

	validateRequired(a, result)
	validateIDs(a, result)
	validateReferences(a, result)
	validateWeights(a, result)
	validateScores(a, result)
	validateScoreCoverage(a, result)

	return result
}

// validateRequired checks that required fields are present.
func validateRequired(a *model.Analysis, result *ValidationResult) {
	if a.ID == "" {
		result.Add("id", "required field is missing")
	}
	if a.Name == "" {
		result.Add("name", "required field is missing")
	}
	if a.Market.ID == "" {
		result.Add("market.id", "required field is missing")
	}
	if a.Market.Name == "" {
		result.Add("market.name", "required field is missing")
	}
	if len(a.Segments) == 0 {
		result.Add("segments", "at least one segment is required")
	}
	if len(a.Vendors) == 0 {
		result.Add("vendors", "at least one vendor is required")
	}
	if len(a.Market.Categories) == 0 {
		result.Add("market.categories", "at least one category is required")
	}

	// Check segments have required fields
	for i, seg := range a.Segments {
		if seg.ID == "" {
			result.Add(fmt.Sprintf("segments[%d].id", i), "required field is missing")
		}
		if seg.Name == "" {
			result.Add(fmt.Sprintf("segments[%d].name", i), "required field is missing")
		}
	}

	// Check vendors have required fields
	for i, v := range a.Vendors {
		if v.ID == "" {
			result.Add(fmt.Sprintf("vendors[%d].id", i), "required field is missing")
		}
		if v.Name == "" {
			result.Add(fmt.Sprintf("vendors[%d].name", i), "required field is missing")
		}
	}

	// Check categories and capabilities have required fields
	for i, cat := range a.Market.Categories {
		if cat.ID == "" {
			result.Add(fmt.Sprintf("market.categories[%d].id", i), "required field is missing")
		}
		if cat.Name == "" {
			result.Add(fmt.Sprintf("market.categories[%d].name", i), "required field is missing")
		}
		for j, cap := range cat.Capabilities {
			if cap.ID == "" {
				result.Add(fmt.Sprintf("market.categories[%d].capabilities[%d].id", i, j), "required field is missing")
			}
			if cap.Name == "" {
				result.Add(fmt.Sprintf("market.categories[%d].capabilities[%d].name", i, j), "required field is missing")
			}
		}
	}
}

// validateIDs checks for duplicate IDs.
func validateIDs(a *model.Analysis, result *ValidationResult) {
	// Check segment IDs
	segmentIDs := make(map[string]bool)
	for i, seg := range a.Segments {
		if segmentIDs[seg.ID] {
			result.Add(fmt.Sprintf("segments[%d].id", i), fmt.Sprintf("duplicate segment ID: %s", seg.ID))
		}
		segmentIDs[seg.ID] = true
	}

	// Check vendor IDs
	vendorIDs := make(map[string]bool)
	for i, v := range a.Vendors {
		if vendorIDs[v.ID] {
			result.Add(fmt.Sprintf("vendors[%d].id", i), fmt.Sprintf("duplicate vendor ID: %s", v.ID))
		}
		vendorIDs[v.ID] = true
	}

	// Check capability IDs
	capabilityIDs := make(map[string]bool)
	for i, cat := range a.Market.Categories {
		for j, cap := range cat.Capabilities {
			if capabilityIDs[cap.ID] {
				result.Add(fmt.Sprintf("market.categories[%d].capabilities[%d].id", i, j),
					fmt.Sprintf("duplicate capability ID: %s", cap.ID))
			}
			capabilityIDs[cap.ID] = true
		}
	}

	// Check category IDs
	categoryIDs := make(map[string]bool)
	for i, cat := range a.Market.Categories {
		if categoryIDs[cat.ID] {
			result.Add(fmt.Sprintf("market.categories[%d].id", i), fmt.Sprintf("duplicate category ID: %s", cat.ID))
		}
		categoryIDs[cat.ID] = true
	}
}

// validateReferences checks that all referenced IDs exist.
func validateReferences(a *model.Analysis, result *ValidationResult) {
	// Build lookup sets
	segmentIDs := make(map[string]bool)
	for _, seg := range a.Segments {
		segmentIDs[seg.ID] = true
	}

	vendorIDs := make(map[string]bool)
	for _, v := range a.Vendors {
		vendorIDs[v.ID] = true
	}

	capabilityIDs := make(map[string]bool)
	for _, cat := range a.Market.Categories {
		for _, cap := range cat.Capabilities {
			capabilityIDs[cap.ID] = true
		}
	}

	// Check focusVendorId
	if a.FocusVendorID != "" && !vendorIDs[a.FocusVendorID] {
		result.Add("focusVendorId", fmt.Sprintf("references unknown vendor: %s", a.FocusVendorID))
	}

	// Check weights reference valid segments and capabilities
	for i, w := range a.Weights {
		if !segmentIDs[w.SegmentID] {
			result.Add(fmt.Sprintf("weights[%d].segmentId", i),
				fmt.Sprintf("references unknown segment: %s", w.SegmentID))
		}
		if !capabilityIDs[w.CapabilityID] {
			result.Add(fmt.Sprintf("weights[%d].capabilityId", i),
				fmt.Sprintf("references unknown capability: %s", w.CapabilityID))
		}
	}

	// Check scores reference valid vendors and capabilities
	for i, s := range a.Scores {
		if !vendorIDs[s.VendorID] {
			result.Add(fmt.Sprintf("scores[%d].vendorId", i),
				fmt.Sprintf("references unknown vendor: %s", s.VendorID))
		}
		if !capabilityIDs[s.CapabilityID] {
			result.Add(fmt.Sprintf("scores[%d].capabilityId", i),
				fmt.Sprintf("references unknown capability: %s", s.CapabilityID))
		}
	}

	// Check vendor target segments reference valid segments
	for i, v := range a.Vendors {
		for j, segID := range v.TargetSegments {
			if !segmentIDs[segID] {
				result.Add(fmt.Sprintf("vendors[%d].targetSegments[%d]", i, j),
					fmt.Sprintf("references unknown segment: %s", segID))
			}
		}
	}
}

// validateWeights checks that weights sum to 1.0 per segment.
func validateWeights(a *model.Analysis, result *ValidationResult) {
	// Get all capability IDs
	var capabilityIDs []string
	for _, cat := range a.Market.Categories {
		for _, cap := range cat.Capabilities {
			capabilityIDs = append(capabilityIDs, cap.ID)
		}
	}

	// Build weight lookup: segmentID -> capabilityID -> weight
	weightMap := make(map[string]map[string]float64)
	for _, w := range a.Weights {
		if weightMap[w.SegmentID] == nil {
			weightMap[w.SegmentID] = make(map[string]float64)
		}
		weightMap[w.SegmentID][w.CapabilityID] = w.Weight
	}

	// Check each segment
	for _, seg := range a.Segments {
		segWeights := weightMap[seg.ID]
		if segWeights == nil {
			result.Add(fmt.Sprintf("weights[%s]", seg.ID),
				"no weights defined for segment")
			continue
		}

		// Check weight values are valid (0-1)
		for capID, weight := range segWeights {
			if weight < 0 || weight > 1 {
				result.Add(fmt.Sprintf("weights[%s][%s]", seg.ID, capID),
					fmt.Sprintf("weight must be between 0 and 1, got %.2f", weight))
			}
		}

		// Sum weights for this segment
		var sum float64
		for _, capID := range capabilityIDs {
			sum += segWeights[capID]
		}

		// Allow small floating point tolerance
		if math.Abs(sum-1.0) > 0.001 {
			result.Add(fmt.Sprintf("weights[%s]", seg.ID),
				fmt.Sprintf("weights must sum to 1.0, got %.3f", sum))
		}

		// Check for missing weights
		for _, capID := range capabilityIDs {
			if _, ok := segWeights[capID]; !ok {
				result.Add(fmt.Sprintf("weights[%s][%s]", seg.ID, capID),
					"missing weight for capability")
			}
		}
	}
}

// validateScores checks that scores are within valid range.
func validateScores(a *model.Analysis, result *ValidationResult) {
	for i, s := range a.Scores {
		if s.Score < 0 || s.Score > 10 {
			result.Add(fmt.Sprintf("scores[%d].score", i),
				fmt.Sprintf("score must be between 0 and 10, got %.1f", s.Score))
		}
	}
}

// validateScoreCoverage checks that all vendors have scores for all capabilities.
func validateScoreCoverage(a *model.Analysis, result *ValidationResult) {
	// Get all capability IDs
	var capabilityIDs []string
	for _, cat := range a.Market.Categories {
		for _, cap := range cat.Capabilities {
			capabilityIDs = append(capabilityIDs, cap.ID)
		}
	}

	// Build score lookup: vendorID -> capabilityID -> exists
	scoreMap := make(map[string]map[string]bool)
	for _, s := range a.Scores {
		if scoreMap[s.VendorID] == nil {
			scoreMap[s.VendorID] = make(map[string]bool)
		}
		scoreMap[s.VendorID][s.CapabilityID] = true
	}

	// Check each vendor has scores for all capabilities
	for _, v := range a.Vendors {
		vendorScores := scoreMap[v.ID]
		if vendorScores == nil {
			result.Add(fmt.Sprintf("scores[%s]", v.ID),
				"no scores defined for vendor")
			continue
		}

		for _, capID := range capabilityIDs {
			if !vendorScores[capID] {
				result.Add(fmt.Sprintf("scores[%s][%s]", v.ID, capID),
					"missing score for capability")
			}
		}
	}
}
