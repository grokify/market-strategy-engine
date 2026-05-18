// Package report generates reports from analysis data.
package report

import (
	"fmt"
	"html/template"
	"io"

	"github.com/grokify/market-strategy-engine/model"
)

// HTMLReport generates an HTML report from an Analysis.
type HTMLReport struct {
	Analysis *model.Analysis
}

// NewHTMLReport creates a new HTML report generator.
func NewHTMLReport(a *model.Analysis) *HTMLReport {
	return &HTMLReport{Analysis: a}
}

// Write writes the HTML report to the given writer.
func (r *HTMLReport) Write(w io.Writer) error {
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"readinessColor":    readinessColor,
		"readinessLevel":    readinessLevelName,
		"severityColor":     severityColor,
		"gapTypeIcon":       gapTypeIcon,
		"effortBars":        effortBars,
		"formatScore":       formatScore,
		"percentBar":        percentBar,
		"competitorBar":     competitorBar,
		"vendorColor":       vendorColor,
		"sub":               func(a, b float64) float64 { return a - b },
		"getCapability":     r.getCapability,
		"getSegment":        r.getSegment,
		"getVendor":         r.getVendor,
		"topGaps":           topGaps,
		"getVendorScore":    getVendorScore,
		"formatWeight":      formatWeight,
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl.Execute(w, r.Analysis)
}

func (r *HTMLReport) getCapability(id string) string {
	if cap := r.Analysis.GetCapability(id); cap != nil {
		return cap.Name
	}
	return id
}

func (r *HTMLReport) getSegment(id string) string {
	if seg := r.Analysis.GetSegment(id); seg != nil {
		return seg.Name
	}
	return id
}

func (r *HTMLReport) getVendor(id string) string {
	if v := r.Analysis.GetVendor(id); v != nil {
		return v.Name
	}
	return id
}

func readinessColor(score float64) string {
	switch {
	case score >= 80:
		return "#22c55e" // green
	case score >= 60:
		return "#eab308" // yellow
	case score >= 40:
		return "#f97316" // orange
	default:
		return "#ef4444" // red
	}
}

func readinessLevelName(score float64) string {
	level := model.GetReadinessLevel(score)
	switch level {
	case model.ReadinessLevelReady:
		return "Ready"
	case model.ReadinessLevelAdjacent:
		return "Adjacent"
	case model.ReadinessLevelPartial:
		return "Partial"
	case model.ReadinessLevelDistant:
		return "Distant"
	default:
		return "Unknown"
	}
}

func severityColor(severity model.GapSeverity) string {
	switch severity {
	case model.GapSeverityCritical:
		return "#ef4444"
	case model.GapSeverityHigh:
		return "#f97316"
	case model.GapSeverityMedium:
		return "#eab308"
	case model.GapSeverityLow:
		return "#22c55e"
	default:
		return "#6b7280"
	}
}

func gapTypeIcon(gapType model.GapType) string {
	switch gapType {
	case model.GapTypeProduct:
		return "🔧"
	case model.GapTypeStructural:
		return "⚡"
	case model.GapTypePerception:
		return "👁"
	default:
		return "?"
	}
}

func effortBars(effort model.EffortLevel) string {
	switch effort {
	case model.EffortLevelLow:
		return "▮"
	case model.EffortLevelMedium:
		return "▮▮"
	case model.EffortLevelHigh:
		return "▮▮▮"
	case model.EffortLevelVeryHigh:
		return "▮▮▮▮"
	default:
		return "?"
	}
}

func formatScore(score float64) string {
	return fmt.Sprintf("%.1f", score)
}

func percentBar(score, max float64) template.HTML {
	pct := (score / max) * 100
	if pct > 100 {
		pct = 100
	}
	color := readinessColor(pct)
	return template.HTML(fmt.Sprintf(
		`<div class="progress-bar"><div class="progress-fill" style="width: %.0f%%; background-color: %s;"></div></div>`,
		pct, color))
}

func topGaps(gaps []model.GapAnalysis, n int) []model.GapAnalysis {
	if len(gaps) <= n {
		return gaps
	}
	return gaps[:n]
}

func competitorBar(score float64, isFocus bool) template.HTML {
	pct := score * 10 // score is 0-10, convert to 0-100%
	if pct > 100 {
		pct = 100
	}
	color := "#3b82f6" // default blue
	if isFocus {
		color = "#8b5cf6" // purple for focus vendor
	}
	return template.HTML(fmt.Sprintf(
		`<div class="competitor-bar"><div class="competitor-fill" style="width: %.0f%%; background-color: %s;"></div></div>`,
		pct, color))
}

func vendorColor(color string, isFocus bool, rank int) string {
	// Use explicit color if provided
	if color != "" {
		return color
	}
	// Fall back to computed colors
	if isFocus {
		return "#8b5cf6" // purple
	}
	colors := []string{"#3b82f6", "#22c55e", "#f97316", "#ef4444", "#6b7280"}
	if rank > 0 && rank <= len(colors) {
		return colors[rank-1]
	}
	return "#6b7280"
}

func getVendorScore(scores map[string]float64, vendorID string) float64 {
	if score, ok := scores[vendorID]; ok {
		return score
	}
	return 0
}

func formatWeight(weight float64) string {
	return fmt.Sprintf("%.0f%%", weight*100)
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Market Strategy Analysis: {{.Name}}</title>
    <style>
        :root {
            --bg-primary: #0f172a;
            --bg-secondary: #1e293b;
            --bg-tertiary: #334155;
            --text-primary: #f8fafc;
            --text-secondary: #94a3b8;
            --border-color: #475569;
            --accent: #3b82f6;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.6;
            padding: 2rem;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            margin-bottom: 3rem;
            padding-bottom: 2rem;
            border-bottom: 1px solid var(--border-color);
        }

        h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            background: linear-gradient(135deg, #3b82f6, #8b5cf6);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 1.1rem;
        }

        .meta {
            display: flex;
            justify-content: center;
            gap: 2rem;
            margin-top: 1rem;
            color: var(--text-secondary);
        }

        .meta-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .meta-label {
            font-weight: 600;
            color: var(--text-primary);
        }

        section {
            margin-bottom: 3rem;
        }

        h2 {
            font-size: 1.5rem;
            margin-bottom: 1.5rem;
            padding-bottom: 0.5rem;
            border-bottom: 2px solid var(--accent);
            display: inline-block;
        }

        .cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1.5rem;
        }

        .card {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            border: 1px solid var(--border-color);
        }

        .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
        }

        .card-title {
            font-size: 1.25rem;
            font-weight: 600;
        }

        .score-badge {
            font-size: 1.5rem;
            font-weight: 700;
            padding: 0.25rem 0.75rem;
            border-radius: 8px;
            color: white;
        }

        .level-badge {
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            background: var(--bg-tertiary);
            color: var(--text-secondary);
        }

        .interpretation {
            color: var(--text-secondary);
            font-size: 0.9rem;
            margin-top: 0.5rem;
        }

        .progress-bar {
            height: 8px;
            background: var(--bg-tertiary);
            border-radius: 4px;
            overflow: hidden;
            margin-top: 1rem;
        }

        .progress-fill {
            height: 100%;
            border-radius: 4px;
            transition: width 0.3s ease;
        }

        .critical-gaps {
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid var(--border-color);
        }

        .critical-gaps-label {
            font-size: 0.8rem;
            color: #ef4444;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            margin-bottom: 0.5rem;
        }

        .critical-gap-item {
            font-size: 0.9rem;
            color: var(--text-secondary);
        }

        table {
            width: 100%;
            border-collapse: collapse;
            background: var(--bg-secondary);
            border-radius: 12px;
            overflow: hidden;
        }

        th, td {
            padding: 1rem;
            text-align: left;
            border-bottom: 1px solid var(--border-color);
        }

        th {
            background: var(--bg-tertiary);
            font-weight: 600;
            font-size: 0.85rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--text-secondary);
        }

        tr:last-child td {
            border-bottom: none;
        }

        tr:hover {
            background: var(--bg-tertiary);
        }

        .rank {
            font-weight: 700;
            color: var(--accent);
            font-size: 1.1rem;
        }

        .gap-type {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            background: var(--bg-tertiary);
            font-size: 0.85rem;
        }

        .score-change {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .score-arrow {
            color: var(--text-secondary);
        }

        .effort {
            font-family: monospace;
            letter-spacing: -0.1em;
        }

        .impact {
            font-size: 0.85rem;
            color: var(--text-secondary);
            max-width: 300px;
        }

        /* Competitive Comparison Styles */
        .segment-comparison {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 2rem;
            border: 1px solid var(--border-color);
        }

        .segment-comparison h3 {
            font-size: 1.25rem;
            margin-bottom: 1.5rem;
            color: var(--text-primary);
        }

        .vendor-ranking {
            margin-bottom: 2rem;
        }

        .vendor-row {
            display: flex;
            align-items: center;
            margin-bottom: 0.75rem;
            padding: 0.5rem;
            border-radius: 8px;
        }

        .vendor-row.focus {
            background: rgba(139, 92, 246, 0.1);
            border: 1px solid rgba(139, 92, 246, 0.3);
        }

        .vendor-rank {
            width: 30px;
            font-weight: 700;
            color: var(--accent);
        }

        .vendor-name {
            width: 150px;
            font-weight: 500;
        }

        .vendor-bar-container {
            flex: 1;
            margin: 0 1rem;
        }

        .competitor-bar {
            height: 24px;
            background: var(--bg-tertiary);
            border-radius: 4px;
            overflow: hidden;
        }

        .competitor-fill {
            height: 100%;
            border-radius: 4px;
            transition: width 0.3s ease;
            display: flex;
            align-items: center;
            justify-content: flex-end;
            padding-right: 8px;
        }

        .vendor-score {
            width: 60px;
            text-align: right;
            font-weight: 600;
        }

        .capability-breakdown {
            margin-top: 1.5rem;
            padding-top: 1.5rem;
            border-top: 1px solid var(--border-color);
        }

        .capability-breakdown h4 {
            font-size: 1rem;
            margin-bottom: 1rem;
            color: var(--text-secondary);
        }

        .capability-row {
            display: grid;
            grid-template-columns: 180px 60px 1fr;
            gap: 1rem;
            align-items: center;
            margin-bottom: 0.5rem;
            font-size: 0.9rem;
        }

        .capability-name {
            color: var(--text-secondary);
        }

        .capability-weight {
            color: var(--text-secondary);
            font-size: 0.8rem;
        }

        .capability-bars {
            display: flex;
            gap: 4px;
            align-items: center;
        }

        .mini-bar {
            height: 16px;
            border-radius: 2px;
            min-width: 4px;
        }

        .mini-bar-label {
            font-size: 0.7rem;
            color: var(--text-secondary);
            margin-left: 2px;
        }

        .vendor-legend {
            display: flex;
            flex-wrap: wrap;
            gap: 1.5rem;
            margin-bottom: 1.5rem;
            padding: 1rem;
            background: var(--bg-secondary);
            border-radius: 8px;
            border: 1px solid var(--border-color);
        }

        .legend-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .legend-color {
            width: 16px;
            height: 16px;
            border-radius: 4px;
        }

        .legend-name {
            font-size: 0.9rem;
            color: var(--text-primary);
        }

        footer {
            text-align: center;
            margin-top: 3rem;
            padding-top: 2rem;
            border-top: 1px solid var(--border-color);
            color: var(--text-secondary);
            font-size: 0.9rem;
        }

        .generated-by {
            margin-top: 0.5rem;
            font-size: 0.8rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>{{.Name}}</h1>
            <p class="subtitle">Market Strategy Analysis Report</p>
            <div class="meta">
                {{if .FocusVendorID}}
                <div class="meta-item">
                    <span class="meta-label">Focus:</span>
                    <span>{{getVendor .FocusVendorID}}</span>
                </div>
                {{end}}
                <div class="meta-item">
                    <span class="meta-label">Market:</span>
                    <span>{{.Market.Name}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Vendors:</span>
                    <span>{{len .Vendors}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Segments:</span>
                    <span>{{len .Segments}}</span>
                </div>
            </div>
        </header>

        <section>
            <h2>Segment Readiness</h2>
            <div class="cards">
                {{range .Readiness}}
                <div class="card">
                    <div class="card-header">
                        <span class="card-title">{{getSegment .SegmentID}}</span>
                        <span class="score-badge" style="background-color: {{readinessColor .Score}}">
                            {{formatScore .Score}}
                        </span>
                    </div>
                    <span class="level-badge">{{readinessLevel .Score}}</span>
                    {{percentBar .Score 100}}
                    <p class="interpretation">{{.Interpretation}}</p>
                    {{if .CriticalGaps}}
                    <div class="critical-gaps">
                        <div class="critical-gaps-label">⚠ Critical Gaps</div>
                        {{range .CriticalGaps}}
                        <div class="critical-gap-item">• {{getCapability .}}</div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </section>

        <section>
            <h2>Competitive Comparison by Segment</h2>
            {{if .Comparisons}}
            <div class="vendor-legend">
                {{$first := index .Comparisons 0}}
                {{range $first.VendorScores}}
                <div class="legend-item">
                    <span class="legend-color" style="background-color: {{vendorColor .VendorColor .IsFocusVendor .Rank}};"></span>
                    <span class="legend-name">{{.VendorName}}</span>
                </div>
                {{end}}
            </div>
            {{end}}
            {{range .Comparisons}}
            <div class="segment-comparison">
                <h3>{{.SegmentName}} Segment</h3>

                <div class="vendor-ranking">
                    {{range .VendorScores}}
                    <div class="vendor-row {{if .IsFocusVendor}}focus{{end}}">
                        <span class="vendor-rank">#{{.Rank}}</span>
                        <span class="vendor-name">{{.VendorName}}</span>
                        <div class="vendor-bar-container">
                            <div class="competitor-bar">
                                <div class="competitor-fill" style="width: {{formatScore .NormalizedScore}}%; background-color: {{vendorColor .VendorColor .IsFocusVendor .Rank}};"></div>
                            </div>
                        </div>
                        <span class="vendor-score">{{formatScore .NormalizedScore}}%</span>
                    </div>
                    {{end}}
                </div>

                {{$vendorList := .VendorScores}}
                <div class="capability-breakdown">
                    <h4>Capability Scores (by importance)</h4>
                    {{range .CapabilityComparisons}}
                    <div class="capability-row">
                        <span class="capability-name">{{.CapabilityName}}</span>
                        <span class="capability-weight">{{formatWeight .SegmentWeight}}</span>
                        <div class="capability-bars">
                            {{$capScores := .VendorScores}}
                            {{range $vendorList}}
                            <div class="mini-bar" style="width: {{formatScore (getVendorScore $capScores .VendorID)}}0px; background-color: {{vendorColor .VendorColor .IsFocusVendor .Rank}};" title="{{.VendorName}}: {{formatScore (getVendorScore $capScores .VendorID)}}"></div>
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
        </section>

        <section>
            <h2>Strategic Priorities</h2>
            <table>
                <thead>
                    <tr>
                        <th>Rank</th>
                        <th>Capability</th>
                        <th>Type</th>
                        <th>Score Gap</th>
                        <th>Target Segment</th>
                        <th>Effort</th>
                        <th>Impact</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Priorities}}
                    <tr>
                        <td><span class="rank">#{{.Rank}}</span></td>
                        <td><strong>{{.CapabilityName}}</strong></td>
                        <td>
                            <span class="gap-type">
                                {{gapTypeIcon .GapType}} {{.GapType}}
                            </span>
                        </td>
                        <td>
                            <div class="score-change">
                                <span>{{formatScore .CurrentScore}}</span>
                                <span class="score-arrow">→</span>
                                <span>{{formatScore .TargetScore}}</span>
                            </div>
                        </td>
                        <td>{{.SegmentName}}</td>
                        <td>
                            <span class="effort" title="{{.Effort}}">{{effortBars .Effort}}</span>
                        </td>
                        <td><span class="impact">{{.Impact}}</span></td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </section>

        <section>
            <h2>Detailed Gap Analysis</h2>
            {{range .Readiness}}
            <h3 style="margin: 1.5rem 0 1rem; color: var(--text-secondary);">{{getSegment .SegmentID}} Segment</h3>
            <table>
                <thead>
                    <tr>
                        <th>Capability</th>
                        <th>Your Score</th>
                        <th>Benchmark</th>
                        <th>Gap</th>
                        <th>Weight</th>
                        <th>Weighted Gap</th>
                        <th>Severity</th>
                    </tr>
                </thead>
                <tbody>
                    {{range topGaps .TopGaps 10}}
                    <tr>
                        <td><strong>{{getCapability .CapabilityID}}</strong></td>
                        <td>{{formatScore .VendorScore}}</td>
                        <td>{{formatScore .BenchmarkScore}}</td>
                        <td>{{formatScore .Gap}}</td>
                        <td>{{formatScore .SegmentWeight}}</td>
                        <td><strong>{{formatScore .WeightedGap}}</strong></td>
                        <td>
                            <span style="color: {{severityColor .Severity}}">●</span>
                            {{.Severity}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{end}}
        </section>

        <footer>
            <p>Analysis: {{.Name}}</p>
            <p class="generated-by">Generated by Market Strategy Engine</p>
        </footer>
    </div>
</body>
</html>`
