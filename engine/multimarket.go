package engine

import (
	"sort"

	"github.com/grokify/market-strategy-engine/model"
)

// MultiMarketEngine computes cross-market analysis from multiple single-market analyses.
type MultiMarketEngine struct {
	analyses []*model.Analysis
}

// NewMultiMarket creates a new multi-market engine from a slice of analyses.
func NewMultiMarket(analyses []*model.Analysis) *MultiMarketEngine {
	return &MultiMarketEngine{analyses: analyses}
}

// ComputeMultiMarket generates a complete multi-market analysis.
func (e *MultiMarketEngine) ComputeMultiMarket(id, name, focusVendorID string) *model.MultiMarketAnalysis {
	mma := &model.MultiMarketAnalysis{
		ID:            id,
		Name:          name,
		FocusVendorID: focusVendorID,
		Markets:       make([]model.Analysis, 0, len(e.analyses)),
	}

	// Run analysis on each market
	for _, a := range e.analyses {
		// Set focus vendor if not already set
		if a.FocusVendorID == "" {
			a.FocusVendorID = focusVendorID
		}

		// Run single-market analysis
		eng := New(a)
		eng.Run()

		mma.Markets = append(mma.Markets, *a)

		// Get focus vendor name
		if mma.FocusVendorName == "" {
			if v := a.GetVendor(focusVendorID); v != nil {
				mma.FocusVendorName = v.Name
			}
		}
	}

	// Compute market summaries
	mma.MarketSummaries = e.computeMarketSummaries(mma)

	// Compute cross-market comparison
	mma.CrossMarketComparison = e.computeCrossMarketComparison(mma)

	return mma
}

func (e *MultiMarketEngine) computeMarketSummaries(mma *model.MultiMarketAnalysis) []model.MarketSummary {
	summaries := make([]model.MarketSummary, 0, len(mma.Markets))

	for _, market := range mma.Markets {
		summary := model.MarketSummary{
			MarketID:    market.Market.ID,
			MarketName:  market.Market.Name,
			VendorCount: len(market.Vendors),
		}

		// Calculate overall readiness (average across segments)
		var totalReadiness float64
		var bestScore, worstScore float64 = -1, 101
		var bestSegment, worstSegment string

		for _, r := range market.Readiness {
			if r.VendorID == mma.FocusVendorID {
				totalReadiness += r.Score

				if r.Score > bestScore {
					bestScore = r.Score
					if seg := market.GetSegment(r.SegmentID); seg != nil {
						bestSegment = seg.Name
					}
				}
				if r.Score < worstScore {
					worstScore = r.Score
					if seg := market.GetSegment(r.SegmentID); seg != nil {
						worstSegment = seg.Name
					}
				}
			}
		}

		focusReadinessCount := 0
		for _, r := range market.Readiness {
			if r.VendorID == mma.FocusVendorID {
				focusReadinessCount++
			}
		}
		if focusReadinessCount > 0 {
			summary.OverallReadiness = totalReadiness / float64(focusReadinessCount)
		}

		summary.BestSegment = bestSegment
		summary.BestSegmentScore = bestScore
		summary.WorstSegment = worstSegment
		summary.WorstSegmentScore = worstScore

		// Find focus vendor rank from comparisons
		if len(market.Comparisons) > 0 {
			// Use first segment's comparison for overall rank
			for _, vs := range market.Comparisons[0].VendorScores {
				if vs.VendorID == mma.FocusVendorID {
					summary.FocusVendorRank = vs.Rank
					break
				}
			}
			summary.TotalVendors = len(market.Comparisons[0].VendorScores)

			// Top competitors (vendors ranked higher than focus)
			for _, vs := range market.Comparisons[0].VendorScores {
				if vs.VendorID != mma.FocusVendorID && vs.Rank < summary.FocusVendorRank {
					summary.TopCompetitors = append(summary.TopCompetitors, vs.VendorName)
				}
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries
}

func (e *MultiMarketEngine) computeCrossMarketComparison(mma *model.MultiMarketAnalysis) *model.CrossMarketComparison {
	cmc := &model.CrossMarketComparison{
		FocusVendorID:   mma.FocusVendorID,
		FocusVendorName: mma.FocusVendorName,
	}

	// Compute market scores
	var bestMarketScore, worstMarketScore float64 = -1, 101
	var bestMarket, worstMarket string

	for _, market := range mma.Markets {
		ms := model.MarketScore{
			MarketID:     market.Market.ID,
			MarketName:   market.Market.Name,
			TotalVendors: len(market.Vendors),
		}

		// Calculate overall score from comparisons
		var totalScore float64
		var segmentCount int

		for _, comp := range market.Comparisons {
			for _, vs := range comp.VendorScores {
				if vs.VendorID == mma.FocusVendorID {
					totalScore += vs.NormalizedScore

					bd := model.SegmentBreakdown{
						SegmentID:   comp.SegmentID,
						SegmentName: comp.SegmentName,
						Score:       vs.NormalizedScore,
						Rank:        vs.Rank,
					}
					ms.SegmentBreakdown = append(ms.SegmentBreakdown, bd)

					if ms.Rank == 0 || vs.Rank < ms.Rank {
						ms.Rank = vs.Rank
					}
					segmentCount++
					break
				}
			}
		}

		if segmentCount > 0 {
			ms.OverallScore = totalScore / float64(segmentCount)
			ms.ReadinessLevel = model.GetReadinessLevel(ms.OverallScore)
		}

		// Track best/worst markets
		if ms.OverallScore > bestMarketScore {
			bestMarketScore = ms.OverallScore
			bestMarket = ms.MarketName
		}
		if ms.OverallScore < worstMarketScore && ms.OverallScore > 0 {
			worstMarketScore = ms.OverallScore
			worstMarket = ms.MarketName
		}

		cmc.MarketScores = append(cmc.MarketScores, ms)
	}

	cmc.StrongestMarket = bestMarket
	cmc.WeakestMarket = worstMarket

	// Compute segment scores (aggregate across markets)
	segmentData := make(map[string]*model.SegmentScore)

	for _, market := range mma.Markets {
		for _, comp := range market.Comparisons {
			for _, vs := range comp.VendorScores {
				if vs.VendorID != mma.FocusVendorID {
					continue
				}

				ss, ok := segmentData[comp.SegmentID]
				if !ok {
					ss = &model.SegmentScore{
						SegmentID:       comp.SegmentID,
						SegmentName:     comp.SegmentName,
						BestMarketScore: -1,
						WorstMarketScore: 101,
					}
					segmentData[comp.SegmentID] = ss
				}

				ss.AverageScore += vs.NormalizedScore
				ss.MarketCount++

				if vs.NormalizedScore > ss.BestMarketScore {
					ss.BestMarketScore = vs.NormalizedScore
					ss.BestMarket = market.Market.Name
				}
				if vs.NormalizedScore < ss.WorstMarketScore {
					ss.WorstMarketScore = vs.NormalizedScore
					ss.WorstMarket = market.Market.Name
				}
			}
		}
	}

	// Finalize segment scores
	var bestSegmentScore, worstSegmentScore float64 = -1, 101
	var bestSegment, worstSegment string

	for _, ss := range segmentData {
		if ss.MarketCount > 0 {
			ss.AverageScore = ss.AverageScore / float64(ss.MarketCount)
		}

		if ss.AverageScore > bestSegmentScore {
			bestSegmentScore = ss.AverageScore
			bestSegment = ss.SegmentName
		}
		if ss.AverageScore < worstSegmentScore && ss.AverageScore > 0 {
			worstSegmentScore = ss.AverageScore
			worstSegment = ss.SegmentName
		}

		cmc.SegmentScores = append(cmc.SegmentScores, *ss)
	}

	cmc.StrongestSegment = bestSegment
	cmc.WeakestSegment = worstSegment

	// Sort segment scores by average score descending
	sort.Slice(cmc.SegmentScores, func(i, j int) bool {
		return cmc.SegmentScores[i].AverageScore > cmc.SegmentScores[j].AverageScore
	})

	// Sort market scores by overall score descending
	sort.Slice(cmc.MarketScores, func(i, j int) bool {
		return cmc.MarketScores[i].OverallScore > cmc.MarketScores[j].OverallScore
	})

	return cmc
}
