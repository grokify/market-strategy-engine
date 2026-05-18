package model

import "time"

// CapabilityScore represents a vendor's score on a specific capability.
// Scores are normalized to a 0-10 scale for consistency across sources.
type CapabilityScore struct {
	// VendorID references the vendor being scored.
	VendorID string `json:"vendorId"`

	// CapabilityID references the capability being measured.
	CapabilityID string `json:"capabilityId"`

	// Score is the normalized score (0.0 to 10.0).
	// 0 = not present/worst, 10 = best in class.
	Score float64 `json:"score"`

	// RawValue is the original value from the source before normalization.
	// Examples: "4.7/5", "Leader", "87%"
	RawValue string `json:"rawValue,omitempty"`

	// SourceID references the data source for this score.
	SourceID string `json:"sourceId"`

	// Confidence indicates how confident we are in this score (0.0 to 1.0).
	// Higher confidence for direct measurements, lower for estimates.
	Confidence float64 `json:"confidence,omitempty"`

	// AsOfDate is when this score was valid.
	AsOfDate time.Time `json:"asOfDate"`

	// Notes provides additional context about the score.
	Notes string `json:"notes,omitempty"`
}

// ScoreNormalization defines how to convert raw scores to the 0-10 scale.
type ScoreNormalization struct {
	// SourceID references the source this normalization applies to.
	SourceID string `json:"sourceId"`

	// CapabilityID references the capability (optional, for capability-specific rules).
	CapabilityID string `json:"capabilityId,omitempty"`

	// InputType describes the input format.
	InputType ScoreInputType `json:"inputType"`

	// InputMin is the minimum value in the source scale.
	InputMin float64 `json:"inputMin,omitempty"`

	// InputMax is the maximum value in the source scale.
	InputMax float64 `json:"inputMax,omitempty"`

	// CategoryMapping maps categorical values to scores (for categorical inputs).
	CategoryMapping map[string]float64 `json:"categoryMapping,omitempty"`

	// Description explains the normalization logic.
	Description string `json:"description,omitempty"`
}

// ScoreInputType describes the format of raw score inputs.
type ScoreInputType string

const (
	// ScoreInputTypeNumeric indicates a numeric scale (e.g., 1-5, 0-100).
	ScoreInputTypeNumeric ScoreInputType = "numeric"

	// ScoreInputTypeCategorical indicates categorical values (e.g., Leader/Challenger).
	ScoreInputTypeCategorical ScoreInputType = "categorical"

	// ScoreInputTypePercentage indicates a percentage (0-100%).
	ScoreInputTypePercentage ScoreInputType = "percentage"

	// ScoreInputTypeBoolean indicates presence/absence (yes/no).
	ScoreInputTypeBoolean ScoreInputType = "boolean"
)

// StandardNormalizations returns common normalization rules.
func StandardNormalizations() []ScoreNormalization {
	return []ScoreNormalization{
		{
			SourceID:    "g2",
			InputType:   ScoreInputTypeNumeric,
			InputMin:    0,
			InputMax:    5,
			Description: "G2 uses 0-5 star ratings, multiply by 2 for 0-10 scale",
		},
		{
			SourceID:  "gartner-mq",
			InputType: ScoreInputTypeCategorical,
			CategoryMapping: map[string]float64{
				"Leader":     9.0,
				"Challenger": 7.0,
				"Visionary":  6.0,
				"Niche":      4.0,
			},
			Description: "Gartner MQ positions mapped to approximate scores",
		},
		{
			SourceID:    "internal",
			InputType:   ScoreInputTypeNumeric,
			InputMin:    0,
			InputMax:    10,
			Description: "Internal assessments use 0-10 scale directly",
		},
	}
}

// Normalize converts a raw score to the 0-10 scale using the given normalization.
func (n ScoreNormalization) Normalize(rawValue string, numericValue float64) float64 {
	switch n.InputType {
	case ScoreInputTypeNumeric:
		if n.InputMax == n.InputMin {
			return 5.0 // avoid division by zero
		}
		return ((numericValue - n.InputMin) / (n.InputMax - n.InputMin)) * 10.0

	case ScoreInputTypeCategorical:
		if score, ok := n.CategoryMapping[rawValue]; ok {
			return score
		}
		return 5.0 // default for unknown categories

	case ScoreInputTypePercentage:
		return numericValue / 10.0 // 0-100% -> 0-10

	case ScoreInputTypeBoolean:
		if numericValue > 0 {
			return 10.0
		}
		return 0.0

	default:
		return numericValue
	}
}
