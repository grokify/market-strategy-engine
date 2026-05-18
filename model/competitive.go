package model

// SegmentComparison represents a competitive comparison within a segment.
type SegmentComparison struct {
	// SegmentID is the segment being compared.
	SegmentID string `json:"segmentId"`

	// SegmentName is the display name.
	SegmentName string `json:"segmentName"`

	// FocusVendorID is the primary vendor being analyzed.
	FocusVendorID string `json:"focusVendorId,omitempty"`

	// VendorScores contains the weighted scores per vendor.
	VendorScores []VendorSegmentScore `json:"vendorScores"`

	// CapabilityComparisons shows per-capability breakdowns.
	CapabilityComparisons []CapabilityComparison `json:"capabilityComparisons"`
}

// VendorSegmentScore represents a vendor's overall score in a segment.
type VendorSegmentScore struct {
	// VendorID is the vendor identifier.
	VendorID string `json:"vendorId"`

	// VendorName is the display name.
	VendorName string `json:"vendorName"`

	// VendorColor is the hex color for this vendor.
	VendorColor string `json:"vendorColor,omitempty"`

	// WeightedScore is the sum of (capability score × segment weight).
	WeightedScore float64 `json:"weightedScore"`

	// MaxPossibleScore is the theoretical maximum (10 × sum of weights = 10).
	MaxPossibleScore float64 `json:"maxPossibleScore"`

	// NormalizedScore is WeightedScore as a percentage of max (0-100).
	NormalizedScore float64 `json:"normalizedScore"`

	// IsFocusVendor indicates if this is the vendor being analyzed.
	IsFocusVendor bool `json:"isFocusVendor,omitempty"`

	// Rank is the position in this segment (1 = highest score).
	Rank int `json:"rank"`
}

// CapabilityComparison shows how vendors compare on a single capability.
type CapabilityComparison struct {
	// CapabilityID is the capability identifier.
	CapabilityID string `json:"capabilityId"`

	// CapabilityName is the display name.
	CapabilityName string `json:"capabilityName"`

	// SegmentWeight is how important this capability is for the segment.
	SegmentWeight float64 `json:"segmentWeight"`

	// VendorScores maps vendor ID to their score (0-10).
	VendorScores map[string]float64 `json:"vendorScores"`

	// BestScore is the highest score among all vendors.
	BestScore float64 `json:"bestScore"`

	// BestVendorID is the vendor with the best score.
	BestVendorID string `json:"bestVendorId"`
}
