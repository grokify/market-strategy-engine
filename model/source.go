package model

import "time"

// DataSource represents where capability scores originate.
// Tracking sources enables auditability and helps assess data quality.
type DataSource struct {
	// ID is a unique identifier for the source.
	ID string `json:"id"`

	// Name is the display name (e.g., "Gartner Magic Quadrant", "G2 Grid").
	Name string `json:"name"`

	// Type categorizes the source.
	Type SourceType `json:"type"`

	// Description provides context about the source.
	Description string `json:"description,omitempty"`

	// URL is a reference link to the source (optional).
	URL string `json:"url,omitempty"`

	// PublishedAt is when the source was published.
	PublishedAt time.Time `json:"publishedAt,omitempty"`

	// RetrievedAt is when the data was retrieved from this source.
	RetrievedAt time.Time `json:"retrievedAt,omitempty"`

	// Credibility is a subjective rating of source reliability (1-10).
	Credibility int `json:"credibility,omitempty"`

	// Notes provides additional context about the source.
	Notes string `json:"notes,omitempty"`
}

// SourceType categorizes the type of data source.
type SourceType string

const (
	// SourceTypeAnalyst indicates analyst firm reports (Gartner, Forrester, IDC).
	SourceTypeAnalyst SourceType = "analyst"

	// SourceTypeReview indicates peer review platforms (G2, TrustRadius, PeerSpot).
	SourceTypeReview SourceType = "review"

	// SourceTypeInternal indicates internal company assessments.
	SourceTypeInternal SourceType = "internal"

	// SourceTypePublic indicates publicly available data (SEC filings, press releases).
	SourceTypePublic SourceType = "public"

	// SourceTypeCommunity indicates community sources (Reddit, forums, social media).
	SourceTypeCommunity SourceType = "community"
)

// SourceTypes returns all valid source types.
func SourceTypes() []SourceType {
	return []SourceType{
		SourceTypeAnalyst,
		SourceTypeReview,
		SourceTypeInternal,
		SourceTypePublic,
		SourceTypeCommunity,
	}
}

// IsValid returns true if the source type is a known valid value.
func (s SourceType) IsValid() bool {
	switch s {
	case SourceTypeAnalyst, SourceTypeReview, SourceTypeInternal, SourceTypePublic, SourceTypeCommunity:
		return true
	default:
		return false
	}
}

// StandardSources returns commonly used data sources.
func StandardSources() []DataSource {
	return []DataSource{
		{
			ID:          "gartner-mq",
			Name:        "Gartner Magic Quadrant",
			Type:        SourceTypeAnalyst,
			Description: "Gartner Magic Quadrant positioning (Leader, Challenger, Visionary, Niche)",
			Credibility: 9,
		},
		{
			ID:          "forrester-wave",
			Name:        "Forrester Wave",
			Type:        SourceTypeAnalyst,
			Description: "Forrester Wave positioning and scores",
			Credibility: 9,
		},
		{
			ID:          "g2",
			Name:        "G2",
			Type:        SourceTypeReview,
			Description: "G2 peer review ratings and grid positioning",
			URL:         "https://www.g2.com",
			Credibility: 7,
		},
		{
			ID:          "trustradius",
			Name:        "TrustRadius",
			Type:        SourceTypeReview,
			Description: "TrustRadius peer review ratings",
			URL:         "https://www.trustradius.com",
			Credibility: 7,
		},
		{
			ID:          "internal",
			Name:        "Internal Assessment",
			Type:        SourceTypeInternal,
			Description: "Internal company evaluation and testing",
			Credibility: 6,
		},
	}
}
