package colors

import (
	"strings"
	"testing"
)

func TestPaletteGet(t *testing.T) {
	p := CompetitiveAnalysis

	// Test normal access
	color := p.Get(0)
	if color != "#8b5cf6" {
		t.Errorf("Get(0): expected #8b5cf6, got %s", color)
	}

	// Test wrap-around
	color = p.Get(10) // Should wrap to index 0
	if color != "#8b5cf6" {
		t.Errorf("Get(10): expected #8b5cf6 (wrap), got %s", color)
	}

	// Test wrap-around with larger index
	color = p.Get(15) // 15 % 10 = 5
	if color != "#06b6d4" {
		t.Errorf("Get(15): expected #06b6d4, got %s", color)
	}
}

func TestPaletteGetEmpty(t *testing.T) {
	p := Palette{Name: "Empty", Colors: []string{}}

	color := p.Get(0)
	if color != "#6b7280" {
		t.Errorf("Get on empty palette: expected default gray, got %s", color)
	}
}

func TestPaletteLen(t *testing.T) {
	if CompetitiveAnalysis.Len() != 10 {
		t.Errorf("CompetitiveAnalysis.Len(): expected 10, got %d", CompetitiveAnalysis.Len())
	}

	if Categorical.Len() != 10 {
		t.Errorf("Categorical.Len(): expected 10, got %d", Categorical.Len())
	}

	if Diverging.Len() != 7 {
		t.Errorf("Diverging.Len(): expected 7, got %d", Diverging.Len())
	}

	if Sequential.Len() != 7 {
		t.Errorf("Sequential.Len(): expected 7, got %d", Sequential.Len())
	}
}

func TestCompetitiveAnalysisPalette(t *testing.T) {
	p := CompetitiveAnalysis

	// First color should be purple (focus vendor)
	if p.Get(0) != "#8b5cf6" {
		t.Errorf("first color should be purple for focus vendor")
	}

	// All colors should be valid hex
	for i, c := range p.Colors {
		if !strings.HasPrefix(c, "#") || len(c) != 7 {
			t.Errorf("color %d (%s) is not valid hex", i, c)
		}
	}
}

func TestCategoricalPalette(t *testing.T) {
	p := Categorical

	// All colors should be valid hex
	for i, c := range p.Colors {
		if !strings.HasPrefix(c, "#") || len(c) != 7 {
			t.Errorf("color %d (%s) is not valid hex", i, c)
		}
	}
}

func TestGetReadinessColor(t *testing.T) {
	tests := []struct {
		score    float64
		expected string
	}{
		{100, "#22c55e"}, // ready (green)
		{80, "#22c55e"},  // ready (green)
		{79, "#eab308"},  // adjacent (yellow)
		{60, "#eab308"},  // adjacent (yellow)
		{59, "#f97316"},  // partial (orange)
		{40, "#f97316"},  // partial (orange)
		{39, "#ef4444"},  // distant (red)
		{0, "#ef4444"},   // distant (red)
	}

	for _, tt := range tests {
		got := GetReadinessColor(tt.score)
		if got != tt.expected {
			t.Errorf("GetReadinessColor(%.0f): expected %s, got %s",
				tt.score, tt.expected, got)
		}
	}
}

func TestSegmentColors(t *testing.T) {
	expected := map[string]string{
		"smb":        "#22c55e",
		"mid-market": "#3b82f6",
		"enterprise": "#8b5cf6",
	}

	for segment, color := range expected {
		if SegmentColors[segment] != color {
			t.Errorf("SegmentColors[%s]: expected %s, got %s",
				segment, color, SegmentColors[segment])
		}
	}
}

func TestGapTypeColors(t *testing.T) {
	expected := map[string]string{
		"product":    "#3b82f6",
		"structural": "#ef4444",
		"perception": "#eab308",
	}

	for gapType, color := range expected {
		if GapTypeColors[gapType] != color {
			t.Errorf("GapTypeColors[%s]: expected %s, got %s",
				gapType, color, GapTypeColors[gapType])
		}
	}
}

func TestSeverityColors(t *testing.T) {
	expected := map[string]string{
		"critical": "#ef4444",
		"high":     "#f97316",
		"medium":   "#eab308",
		"low":      "#22c55e",
		"none":     "#6b7280",
	}

	for severity, color := range expected {
		if SeverityColors[severity] != color {
			t.Errorf("SeverityColors[%s]: expected %s, got %s",
				severity, color, SeverityColors[severity])
		}
	}
}

func TestDefaultPalettes(t *testing.T) {
	palettes := DefaultPalettes()

	if len(palettes) != 4 {
		t.Errorf("expected 4 default palettes, got %d", len(palettes))
	}

	names := make(map[string]bool)
	for _, p := range palettes {
		names[p.Name] = true
	}

	expectedNames := []string{"CompetitiveAnalysis", "Categorical", "Diverging", "Sequential"}
	for _, name := range expectedNames {
		if !names[name] {
			t.Errorf("expected palette %s in defaults", name)
		}
	}
}
