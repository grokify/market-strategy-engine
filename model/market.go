package model

// Market represents an industry vertical (e.g., Security, CRM, Fintech).
// A market contains multiple product categories and defines the competitive
// landscape for analysis.
type Market struct {
	// ID is a unique identifier for the market.
	ID string `json:"id"`

	// Name is the display name of the market.
	Name string `json:"name"`

	// Description provides additional context about the market.
	Description string `json:"description,omitempty"`

	// Categories are the product categories within this market.
	Categories []ProductCategory `json:"categories,omitempty"`
}

// ProductCategory represents a functional domain within a market.
// For example, in Security: EDR, MDR, SIEM, ITDR, SAT.
type ProductCategory struct {
	// ID is a unique identifier for the category.
	ID string `json:"id"`

	// Name is the display name (e.g., "EDR", "MDR").
	Name string `json:"name"`

	// Description explains what this category covers.
	Description string `json:"description,omitempty"`

	// Capabilities are the measurable attributes for this category.
	Capabilities []Capability `json:"capabilities,omitempty"`
}

// Capability represents a measurable product attribute that can be
// scored and compared across vendors. Examples: detection depth,
// ease of deployment, SOC quality, telemetry access.
type Capability struct {
	// ID is a unique identifier for the capability.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Description explains what this capability measures.
	Description string `json:"description,omitempty"`

	// GapType categorizes the nature of gaps in this capability.
	// This affects how gaps should be addressed (build vs. architecture vs. marketing).
	GapType GapType `json:"gapType,omitempty"`

	// ScoreRubric describes how to score this capability (0-10 scale).
	ScoreRubric string `json:"scoreRubric,omitempty"`
}

// GapType categorizes the nature of a capability gap.
// Different gap types have different cost/time profiles for remediation.
type GapType string

const (
	// GapTypeProduct indicates a buildable feature gap.
	// Can typically be addressed in quarters through product development.
	// Examples: missing integrations, UI improvements, reporting features.
	GapTypeProduct GapType = "product"

	// GapTypeStructural indicates an architectural limitation.
	// Requires significant platform changes to address.
	// Examples: telemetry pipeline, data model limitations, scalability architecture.
	GapTypeStructural GapType = "structural"

	// GapTypePerception indicates a market perception or brand gap.
	// Addressed through marketing, sales, and market presence rather than product.
	// Examples: brand credibility, analyst positioning, enterprise trust.
	GapTypePerception GapType = "perception"
)

// GapTypes returns all valid gap types.
func GapTypes() []GapType {
	return []GapType{GapTypeProduct, GapTypeStructural, GapTypePerception}
}

// IsValid returns true if the gap type is a known valid value.
func (g GapType) IsValid() bool {
	switch g {
	case GapTypeProduct, GapTypeStructural, GapTypePerception:
		return true
	default:
		return false
	}
}
