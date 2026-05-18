package model

// MultiMarketAnalysis represents competitive analysis across multiple markets.
// This allows tracking a vendor's position across different product categories
// (e.g., EDR, MDR, ITDR, SIEM, SAT).
type MultiMarketAnalysis struct {
	// ID is a unique identifier for this multi-market analysis.
	ID string `json:"id"`

	// Name is a human-readable name for this analysis.
	Name string `json:"name"`

	// Description provides context about the analysis scope.
	Description string `json:"description,omitempty"`

	// FocusVendorID is the primary vendor being analyzed.
	FocusVendorID string `json:"focusVendorId"`

	// FocusVendorName is cached for convenience.
	FocusVendorName string `json:"focusVendorName,omitempty"`

	// Markets contains the individual market analyses.
	Markets []Analysis `json:"markets"`

	// MarketSummaries provides an overview of each market.
	MarketSummaries []MarketSummary `json:"marketSummaries,omitempty"`

	// CrossMarketComparison shows how the focus vendor compares across all markets.
	CrossMarketComparison *CrossMarketComparison `json:"crossMarketComparison,omitempty"`
}

// MarketSummary provides a high-level view of a single market.
type MarketSummary struct {
	// MarketID is the market identifier.
	MarketID string `json:"marketId"`

	// MarketName is the display name.
	MarketName string `json:"marketName"`

	// VendorCount is how many vendors compete in this market.
	VendorCount int `json:"vendorCount"`

	// FocusVendorRank is where the focus vendor ranks overall in this market.
	FocusVendorRank int `json:"focusVendorRank"`

	// TotalVendors is the total number of vendors analyzed.
	TotalVendors int `json:"totalVendors"`

	// BestSegment is the segment where focus vendor performs best.
	BestSegment string `json:"bestSegment,omitempty"`

	// BestSegmentScore is the score in the best segment.
	BestSegmentScore float64 `json:"bestSegmentScore,omitempty"`

	// WorstSegment is the segment where focus vendor needs most improvement.
	WorstSegment string `json:"worstSegment,omitempty"`

	// WorstSegmentScore is the score in the worst segment.
	WorstSegmentScore float64 `json:"worstSegmentScore,omitempty"`

	// OverallReadiness is the average readiness across all segments.
	OverallReadiness float64 `json:"overallReadiness"`

	// TopCompetitors lists the main competitors in this market.
	TopCompetitors []string `json:"topCompetitors,omitempty"`
}

// CrossMarketComparison shows how a vendor performs across all markets.
type CrossMarketComparison struct {
	// FocusVendorID is the vendor being analyzed.
	FocusVendorID string `json:"focusVendorId"`

	// FocusVendorName is the display name.
	FocusVendorName string `json:"focusVendorName"`

	// MarketScores shows performance in each market.
	MarketScores []MarketScore `json:"marketScores"`

	// SegmentScores shows performance in each segment across all markets.
	SegmentScores []SegmentScore `json:"segmentScores"`

	// StrongestMarket is where the vendor performs best overall.
	StrongestMarket string `json:"strongestMarket"`

	// WeakestMarket is where the vendor needs most improvement.
	WeakestMarket string `json:"weakestMarket"`

	// StrongestSegment is the segment type where vendor excels.
	StrongestSegment string `json:"strongestSegment"`

	// WeakestSegment is the segment type needing most work.
	WeakestSegment string `json:"weakestSegment"`
}

// MarketScore represents a vendor's overall score in a specific market.
type MarketScore struct {
	// MarketID is the market identifier.
	MarketID string `json:"marketId"`

	// MarketName is the display name.
	MarketName string `json:"marketName"`

	// OverallScore is the average normalized score across segments (0-100).
	OverallScore float64 `json:"overallScore"`

	// Rank is position among all vendors in this market.
	Rank int `json:"rank"`

	// TotalVendors is how many vendors compete.
	TotalVendors int `json:"totalVendors"`

	// ReadinessLevel is the overall readiness (ready, adjacent, partial, distant).
	ReadinessLevel ReadinessLevel `json:"readinessLevel"`

	// SegmentBreakdown shows performance by segment.
	SegmentBreakdown []SegmentBreakdown `json:"segmentBreakdown,omitempty"`
}

// SegmentBreakdown shows performance in a single segment within a market.
type SegmentBreakdown struct {
	// SegmentID is the segment identifier.
	SegmentID string `json:"segmentId"`

	// SegmentName is the display name.
	SegmentName string `json:"segmentName"`

	// Score is the normalized score (0-100).
	Score float64 `json:"score"`

	// Rank is position among vendors in this segment.
	Rank int `json:"rank"`
}

// SegmentScore aggregates vendor performance across all markets for a segment type.
// This helps identify whether a vendor is consistently strong/weak in certain segments.
type SegmentScore struct {
	// SegmentID is the segment identifier (e.g., "smb", "enterprise").
	SegmentID string `json:"segmentId"`

	// SegmentName is the display name.
	SegmentName string `json:"segmentName"`

	// AverageScore is the average score across all markets for this segment.
	AverageScore float64 `json:"averageScore"`

	// MarketCount is how many markets have this segment.
	MarketCount int `json:"marketCount"`

	// BestMarket is where the vendor does best in this segment.
	BestMarket string `json:"bestMarket,omitempty"`

	// BestMarketScore is the score in the best market.
	BestMarketScore float64 `json:"bestMarketScore,omitempty"`

	// WorstMarket is where the vendor does worst in this segment.
	WorstMarket string `json:"worstMarket,omitempty"`

	// WorstMarketScore is the score in the worst market.
	WorstMarketScore float64 `json:"worstMarketScore,omitempty"`
}
