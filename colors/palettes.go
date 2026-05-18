// Package colors provides reusable color palettes for visualizations.
package colors

// Palette represents a named collection of colors.
type Palette struct {
	Name        string
	Description string
	Colors      []string // hex colors
}

// Get returns a color by index, wrapping around if needed.
func (p Palette) Get(index int) string {
	if len(p.Colors) == 0 {
		return "#6b7280" // default gray
	}
	return p.Colors[index%len(p.Colors)]
}

// Len returns the number of colors in the palette.
func (p Palette) Len() int {
	return len(p.Colors)
}

// CompetitiveAnalysis is a palette designed for competitive analysis charts.
// Colors are chosen to be distinct and work well on dark backgrounds.
// First color (purple) is intended for the focus/target company.
var CompetitiveAnalysis = Palette{
	Name:        "CompetitiveAnalysis",
	Description: "Optimized for competitive analysis visualizations with focus vendor emphasis",
	Colors: []string{
		"#8b5cf6", // purple - focus vendor
		"#ef4444", // red - primary competitor
		"#3b82f6", // blue - secondary competitor
		"#22c55e", // green - tertiary competitor
		"#f97316", // orange
		"#06b6d4", // cyan
		"#ec4899", // pink
		"#eab308", // yellow
		"#6366f1", // indigo
		"#14b8a6", // teal
	},
}

// Categorical is a general-purpose categorical palette.
// Based on Tableau 10 with modifications for dark backgrounds.
var Categorical = Palette{
	Name:        "Categorical",
	Description: "General-purpose categorical palette for charts",
	Colors: []string{
		"#4e79a7", // blue
		"#f28e2c", // orange
		"#e15759", // red
		"#76b7b2", // teal
		"#59a14f", // green
		"#edc949", // yellow
		"#af7aa1", // purple
		"#ff9da7", // pink
		"#9c755f", // brown
		"#bab0ab", // gray
	},
}

// Diverging is a palette for showing positive/negative or above/below comparisons.
var Diverging = Palette{
	Name:        "Diverging",
	Description: "For showing positive/negative divergence from a center point",
	Colors: []string{
		"#ef4444", // negative (red)
		"#f97316", // slightly negative
		"#eab308", // neutral negative
		"#6b7280", // neutral (gray)
		"#84cc16", // neutral positive
		"#22c55e", // slightly positive
		"#10b981", // positive (green)
	},
}

// Sequential is a palette for showing ordered/ranked data.
var Sequential = Palette{
	Name:        "Sequential",
	Description: "For showing ordered or ranked data",
	Colors: []string{
		"#1e3a5f", // darkest
		"#2563eb",
		"#3b82f6",
		"#60a5fa",
		"#93c5fd",
		"#bfdbfe",
		"#dbeafe", // lightest
	},
}

// SegmentColors provides consistent colors for market segments.
var SegmentColors = map[string]string{
	"smb":        "#22c55e", // green - accessible
	"mid-market": "#3b82f6", // blue - growing
	"enterprise": "#8b5cf6", // purple - premium
}

// GapTypeColors provides consistent colors for gap types.
var GapTypeColors = map[string]string{
	"product":    "#3b82f6", // blue - buildable
	"structural": "#ef4444", // red - hard
	"perception": "#eab308", // yellow - marketing
}

// SeverityColors provides consistent colors for severity levels.
var SeverityColors = map[string]string{
	"critical": "#ef4444", // red
	"high":     "#f97316", // orange
	"medium":   "#eab308", // yellow
	"low":      "#22c55e", // green
	"none":     "#6b7280", // gray
}

// ReadinessColors provides colors for readiness score ranges.
var ReadinessColors = map[string]string{
	"ready":    "#22c55e", // green (80-100)
	"adjacent": "#eab308", // yellow (60-79)
	"partial":  "#f97316", // orange (40-59)
	"distant":  "#ef4444", // red (0-39)
}

// GetReadinessColor returns a color for a readiness score.
func GetReadinessColor(score float64) string {
	switch {
	case score >= 80:
		return ReadinessColors["ready"]
	case score >= 60:
		return ReadinessColors["adjacent"]
	case score >= 40:
		return ReadinessColors["partial"]
	default:
		return ReadinessColors["distant"]
	}
}

// DefaultPalettes returns all available palettes.
func DefaultPalettes() []Palette {
	return []Palette{
		CompetitiveAnalysis,
		Categorical,
		Diverging,
		Sequential,
	}
}
