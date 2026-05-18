package model

// Segment represents a customer segment defined by organization size,
// maturity, or other characteristics. Common segments: SMB, Mid-Market, Enterprise.
//
// Different segments value different capabilities. SMB prioritizes simplicity
// and cost; Enterprise prioritizes depth and customization.
type Segment struct {
	// ID is a unique identifier for the segment.
	ID string `json:"id"`

	// Name is the display name (e.g., "SMB", "Mid-Market", "Enterprise").
	Name string `json:"name"`

	// Description explains the segment characteristics.
	Description string `json:"description,omitempty"`

	// MinEndpoints is the minimum endpoint count for this segment (optional).
	MinEndpoints int `json:"minEndpoints,omitempty"`

	// MaxEndpoints is the maximum endpoint count for this segment (optional).
	// Use 0 or omit for "unlimited" (e.g., Enterprise has no upper bound).
	MaxEndpoints int `json:"maxEndpoints,omitempty"`

	// Characteristics describes what buyers in this segment prioritize.
	Characteristics []string `json:"characteristics,omitempty"`
}

// SegmentWeight defines the importance of a capability for a specific segment.
// Weights should sum to 1.0 across all capabilities for a given segment.
type SegmentWeight struct {
	// SegmentID references the segment this weight applies to.
	SegmentID string `json:"segmentId"`

	// CapabilityID references the capability being weighted.
	CapabilityID string `json:"capabilityId"`

	// Weight is the importance factor (0.0 to 1.0).
	// Higher weight means this capability matters more for this segment.
	Weight float64 `json:"weight"`

	// Rationale explains why this weight was chosen.
	Rationale string `json:"rationale,omitempty"`
}

// StandardSegments returns the common three-tier segment model.
func StandardSegments() []Segment {
	return []Segment{
		{
			ID:           "smb",
			Name:         "SMB",
			Description:  "Small and medium businesses with limited or no dedicated security staff",
			MinEndpoints: 1,
			MaxEndpoints: 500,
			Characteristics: []string{
				"Optimizes for simplicity and cost",
				"No internal SOC or security engineering",
				"Prefers managed/outsourced security",
				"Values fast deployment and low maintenance",
			},
		},
		{
			ID:           "mid-market",
			Name:         "Mid-Market",
			Description:  "Growing organizations with some security maturity but not full enterprise scale",
			MinEndpoints: 500,
			MaxEndpoints: 10000,
			Characteristics: []string{
				"Balance of automation and control",
				"Some internal security staff",
				"Needs integration with existing IT stack",
				"Values coverage breadth and operational efficiency",
			},
		},
		{
			ID:           "enterprise",
			Name:         "Enterprise",
			Description:  "Large organizations with mature security programs and dedicated teams",
			MinEndpoints: 10000,
			MaxEndpoints: 0, // unlimited
			Characteristics: []string{
				"Optimizes for depth and customization",
				"Has internal security engineering teams",
				"Requires deep telemetry and data access",
				"Values threat hunting and custom detection",
				"Needs global scale and compliance",
			},
		},
	}
}
