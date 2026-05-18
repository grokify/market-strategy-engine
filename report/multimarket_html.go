package report

import (
	"fmt"
	"html/template"
	"io"

	"github.com/grokify/market-strategy-engine/model"
)

// MultiMarketHTMLReport generates an HTML report for multi-market analysis.
type MultiMarketHTMLReport struct {
	Analysis *model.MultiMarketAnalysis
}

// NewMultiMarketHTMLReport creates a new multi-market HTML report generator.
func NewMultiMarketHTMLReport(a *model.MultiMarketAnalysis) *MultiMarketHTMLReport {
	return &MultiMarketHTMLReport{Analysis: a}
}

// Write writes the HTML report to the given writer.
func (r *MultiMarketHTMLReport) Write(w io.Writer) error {
	tmpl, err := template.New("multimarket").Funcs(template.FuncMap{
		"readinessColor": readinessColor,
		"readinessLevel": readinessLevelName,
		"formatScore":    formatScore,
		"percentBar":     percentBar,
		"marketBar":      marketBar,
	}).Parse(multiMarketHTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl.Execute(w, r.Analysis)
}

func marketBar(score float64, color string) template.HTML {
	if color == "" {
		color = readinessColor(score)
	}
	return template.HTML(fmt.Sprintf(
		`<div class="market-bar"><div class="market-fill" style="width: %.0f%%; background-color: %s;"></div></div>`,
		score, color))
}

const multiMarketHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Multi-Market Analysis: {{.Name}}</title>
    <style>
        :root {
            --bg-primary: #0f172a;
            --bg-secondary: #1e293b;
            --bg-tertiary: #334155;
            --text-primary: #f8fafc;
            --text-secondary: #94a3b8;
            --border-color: #475569;
            --accent: #3b82f6;
            --purple: #8b5cf6;
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
            max-width: 1400px;
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

        h3 {
            font-size: 1.25rem;
            margin-bottom: 1rem;
            color: var(--text-secondary);
        }

        /* Cross-Market Overview */
        .overview-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
        }

        .overview-card {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            border: 1px solid var(--border-color);
        }

        .overview-card.highlight {
            border-color: var(--purple);
            background: linear-gradient(135deg, rgba(139, 92, 246, 0.1), rgba(59, 130, 246, 0.1));
        }

        .overview-card h4 {
            font-size: 0.9rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--text-secondary);
            margin-bottom: 0.5rem;
        }

        .overview-value {
            font-size: 2rem;
            font-weight: 700;
        }

        .overview-detail {
            font-size: 0.9rem;
            color: var(--text-secondary);
            margin-top: 0.5rem;
        }

        /* Market Matrix */
        .market-matrix {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            border: 1px solid var(--border-color);
            overflow-x: auto;
        }

        .matrix-table {
            width: 100%;
            border-collapse: collapse;
        }

        .matrix-table th,
        .matrix-table td {
            padding: 1rem;
            text-align: left;
            border-bottom: 1px solid var(--border-color);
        }

        .matrix-table th {
            background: var(--bg-tertiary);
            font-weight: 600;
            font-size: 0.85rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--text-secondary);
        }

        .matrix-table tr:last-child td {
            border-bottom: none;
        }

        .matrix-table tr:hover {
            background: var(--bg-tertiary);
        }

        .market-name {
            font-weight: 600;
            color: var(--text-primary);
        }

        .rank-badge {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 28px;
            height: 28px;
            border-radius: 50%;
            font-weight: 700;
            font-size: 0.9rem;
        }

        .rank-1 { background: #fbbf24; color: #000; }
        .rank-2 { background: #9ca3af; color: #000; }
        .rank-3 { background: #cd7f32; color: #000; }
        .rank-other { background: var(--bg-tertiary); color: var(--text-secondary); }

        .market-bar {
            height: 24px;
            background: var(--bg-tertiary);
            border-radius: 4px;
            overflow: hidden;
            min-width: 100px;
        }

        .market-fill {
            height: 100%;
            border-radius: 4px;
            transition: width 0.3s ease;
        }

        .score-cell {
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .score-value {
            min-width: 50px;
            font-weight: 600;
        }

        /* Segment Heatmap */
        .heatmap {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            border: 1px solid var(--border-color);
        }

        .heatmap-grid {
            display: grid;
            gap: 1rem;
        }

        .heatmap-row {
            display: grid;
            grid-template-columns: 120px repeat(auto-fit, minmax(100px, 1fr));
            gap: 0.5rem;
            align-items: center;
        }

        .heatmap-header {
            font-weight: 600;
            font-size: 0.85rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--text-secondary);
        }

        .heatmap-cell {
            padding: 0.75rem;
            border-radius: 8px;
            text-align: center;
            font-weight: 600;
            font-size: 0.9rem;
        }

        .heatmap-label {
            color: var(--text-secondary);
            font-size: 0.9rem;
        }

        /* Strengths/Weaknesses */
        .insights-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1.5rem;
        }

        .insight-card {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            border: 1px solid var(--border-color);
        }

        .insight-card h4 {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            margin-bottom: 1rem;
            font-size: 1rem;
        }

        .insight-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 0.75rem 0;
            border-bottom: 1px solid var(--border-color);
        }

        .insight-item:last-child {
            border-bottom: none;
        }

        .insight-label {
            color: var(--text-secondary);
        }

        .insight-value {
            font-weight: 600;
        }

        .strength { color: #22c55e; }
        .weakness { color: #ef4444; }

        /* Market Cards */
        .market-cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 1.5rem;
        }

        .market-card {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 1.5rem;
            border: 1px solid var(--border-color);
        }

        .market-card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
        }

        .market-card-title {
            font-size: 1.1rem;
            font-weight: 600;
        }

        .market-card-score {
            font-size: 1.5rem;
            font-weight: 700;
            padding: 0.25rem 0.75rem;
            border-radius: 8px;
        }

        .segment-list {
            margin-top: 1rem;
        }

        .segment-item {
            display: flex;
            align-items: center;
            gap: 1rem;
            padding: 0.5rem 0;
        }

        .segment-name {
            width: 100px;
            font-size: 0.9rem;
            color: var(--text-secondary);
        }

        .segment-bar-container {
            flex: 1;
        }

        .segment-score {
            width: 50px;
            text-align: right;
            font-weight: 600;
            font-size: 0.9rem;
        }

        footer {
            text-align: center;
            margin-top: 3rem;
            padding-top: 2rem;
            border-top: 1px solid var(--border-color);
            color: var(--text-secondary);
            font-size: 0.9rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>{{.Name}}</h1>
            <p class="subtitle">Multi-Market Competitive Analysis</p>
            <div class="meta">
                <div class="meta-item">
                    <span class="meta-label">Focus:</span>
                    <span>{{.FocusVendorName}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Markets:</span>
                    <span>{{len .Markets}}</span>
                </div>
            </div>
        </header>

        {{if .CrossMarketComparison}}
        <section>
            <h2>Cross-Market Overview</h2>
            <div class="overview-grid">
                <div class="overview-card highlight">
                    <h4>Strongest Market</h4>
                    <div class="overview-value">{{.CrossMarketComparison.StrongestMarket}}</div>
                    <div class="overview-detail">Best overall competitive position</div>
                </div>
                <div class="overview-card">
                    <h4>Weakest Market</h4>
                    <div class="overview-value">{{.CrossMarketComparison.WeakestMarket}}</div>
                    <div class="overview-detail">Most opportunity for improvement</div>
                </div>
                <div class="overview-card highlight">
                    <h4>Strongest Segment</h4>
                    <div class="overview-value">{{.CrossMarketComparison.StrongestSegment}}</div>
                    <div class="overview-detail">Consistent strength across markets</div>
                </div>
                <div class="overview-card">
                    <h4>Weakest Segment</h4>
                    <div class="overview-value">{{.CrossMarketComparison.WeakestSegment}}</div>
                    <div class="overview-detail">Focus area for expansion</div>
                </div>
            </div>
        </section>

        <section>
            <h2>Market Performance Matrix</h2>
            <div class="market-matrix">
                <table class="matrix-table">
                    <thead>
                        <tr>
                            <th>Market</th>
                            <th>Rank</th>
                            <th>Overall Score</th>
                            <th>Readiness</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .CrossMarketComparison.MarketScores}}
                        <tr>
                            <td class="market-name">{{.MarketName}}</td>
                            <td>
                                <span class="rank-badge {{if eq .Rank 1}}rank-1{{else if eq .Rank 2}}rank-2{{else if eq .Rank 3}}rank-3{{else}}rank-other{{end}}">
                                    #{{.Rank}}
                                </span>
                                <span style="color: var(--text-secondary); margin-left: 0.5rem;">of {{.TotalVendors}}</span>
                            </td>
                            <td>
                                <div class="score-cell">
                                    <span class="score-value">{{formatScore .OverallScore}}%</span>
                                    {{marketBar .OverallScore ""}}
                                </div>
                            </td>
                            <td>
                                <span style="color: {{readinessColor .OverallScore}}">●</span>
                                {{readinessLevel .OverallScore}}
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </section>

        <section>
            <h2>Segment Performance Across Markets</h2>
            <div class="insights-grid">
                {{range .CrossMarketComparison.SegmentScores}}
                <div class="insight-card">
                    <h4>{{.SegmentName}}</h4>
                    <div class="insight-item">
                        <span class="insight-label">Average Score</span>
                        <span class="insight-value" style="color: {{readinessColor .AverageScore}}">{{formatScore .AverageScore}}%</span>
                    </div>
                    <div class="insight-item">
                        <span class="insight-label">Markets Analyzed</span>
                        <span class="insight-value">{{.MarketCount}}</span>
                    </div>
                    {{if .BestMarket}}
                    <div class="insight-item">
                        <span class="insight-label">Best In</span>
                        <span class="insight-value strength">{{.BestMarket}} ({{formatScore .BestMarketScore}}%)</span>
                    </div>
                    {{end}}
                    {{if .WorstMarket}}
                    <div class="insight-item">
                        <span class="insight-label">Weakest In</span>
                        <span class="insight-value weakness">{{.WorstMarket}} ({{formatScore .WorstMarketScore}}%)</span>
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </section>
        {{end}}

        <section>
            <h2>Individual Market Details</h2>
            <div class="market-cards">
                {{range .MarketSummaries}}
                <div class="market-card">
                    <div class="market-card-header">
                        <span class="market-card-title">{{.MarketName}}</span>
                        <span class="market-card-score" style="background-color: {{readinessColor .OverallReadiness}}; color: white;">
                            {{formatScore .OverallReadiness}}%
                        </span>
                    </div>
                    <div class="insight-item">
                        <span class="insight-label">Rank</span>
                        <span class="insight-value">#{{.FocusVendorRank}} of {{.TotalVendors}}</span>
                    </div>
                    {{if .BestSegment}}
                    <div class="insight-item">
                        <span class="insight-label">Best Segment</span>
                        <span class="insight-value strength">{{.BestSegment}} ({{formatScore .BestSegmentScore}}%)</span>
                    </div>
                    {{end}}
                    {{if .WorstSegment}}
                    <div class="insight-item">
                        <span class="insight-label">Weakest Segment</span>
                        <span class="insight-value weakness">{{.WorstSegment}} ({{formatScore .WorstSegmentScore}}%)</span>
                    </div>
                    {{end}}
                    {{if .TopCompetitors}}
                    <div class="insight-item">
                        <span class="insight-label">Top Competitors</span>
                        <span class="insight-value">{{range $i, $c := .TopCompetitors}}{{if $i}}, {{end}}{{$c}}{{end}}</span>
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </section>

        <footer>
            <p>Analysis: {{.Name}}</p>
            <p style="margin-top: 0.5rem; font-size: 0.8rem;">Generated by Market Strategy Engine</p>
        </footer>
    </div>
</body>
</html>`
