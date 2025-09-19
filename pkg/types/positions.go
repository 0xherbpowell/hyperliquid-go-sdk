package types

import (
	"encoding/json"
	"strconv"
)

// PositionType represents the type of position
type PositionType string

const (
	PositionTypeOneWay PositionType = "oneWay"
)

// LeverageType represents the type of leverage
type LeverageType string

const (
	LeverageTypeCross    LeverageType = "cross"
	LeverageTypeIsolated LeverageType = "isolated"
)

// DetailedPosition represents a detailed trading position with all fields
type DetailedPosition struct {
	Coin           string         `json:"coin"`
	EntryPx        *FloatString   `json:"entryPx"`
	Leverage       LeverageInfo   `json:"leverage"`
	LiquidationPx  *FloatString   `json:"liquidationPx"`
	MarginUsed     FloatString    `json:"marginUsed"`
	MaxTradeSzs    [2]FloatString `json:"maxTradeSzs"`
	PositionValue  FloatString    `json:"positionValue"`
	ReturnOnEquity FloatString    `json:"returnOnEquity"`
	Szi            FloatString    `json:"szi"` // Signed size (positive for long, negative for short)
	UnrealizedPnl  FloatString    `json:"unrealizedPnl"`
}

// LeverageInfo represents leverage information that can be either cross or isolated
type LeverageInfo struct {
	Type   LeverageType `json:"type"`
	Value  int          `json:"value"`
	RawUsd *FloatString `json:"rawUsd,omitempty"` // Only for isolated margin
}

// UnmarshalJSON implements custom JSON unmarshaling for LeverageInfo
func (l *LeverageInfo) UnmarshalJSON(data []byte) error {
	var temp struct {
		Type   LeverageType `json:"type"`
		Value  int          `json:"value"`
		RawUsd *FloatString `json:"rawUsd,omitempty"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	l.Type = temp.Type
	l.Value = temp.Value
	l.RawUsd = temp.RawUsd

	return nil
}

// MarshalJSON implements custom JSON marshaling for LeverageInfo
func (l LeverageInfo) MarshalJSON() ([]byte, error) {
	temp := struct {
		Type   LeverageType `json:"type"`
		Value  int          `json:"value"`
		RawUsd *FloatString `json:"rawUsd,omitempty"`
	}{
		Type:   l.Type,
		Value:  l.Value,
		RawUsd: l.RawUsd,
	}

	return json.Marshal(temp)
}

// IsCross returns true if the leverage is cross margin
func (l LeverageInfo) IsCross() bool {
	return l.Type == LeverageTypeCross
}

// IsIsolated returns true if the leverage is isolated margin
func (l LeverageInfo) IsIsolated() bool {
	return l.Type == LeverageTypeIsolated
}

// GetLeverageValue returns the leverage multiplier
func (l LeverageInfo) GetLeverageValue() int {
	return l.Value
}

// GetRawUsd returns the raw USD amount for isolated margin (nil for cross margin)
func (l LeverageInfo) GetRawUsd() *FloatString {
	return l.RawUsd
}

// DetailedAssetPosition represents an asset position with type information
type DetailedAssetPosition struct {
	Position DetailedPosition `json:"position"`
	Type     PositionType     `json:"type"`
}

// PositionSummary represents a summary of position information
type PositionSummary struct {
	Coin              string       `json:"coin"`
	Size              FloatString  `json:"size"` // Absolute size
	Side              string       `json:"side"` // "long", "short", or "none"
	EntryPrice        *FloatString `json:"entryPrice"`
	MarkPrice         FloatString  `json:"markPrice"`
	UnrealizedPnl     FloatString  `json:"unrealizedPnl"`
	UnrealizedPnlPerc FloatString  `json:"unrealizedPnlPerc"`
	MarginUsed        FloatString  `json:"marginUsed"`
	LeverageType      LeverageType `json:"leverageType"`
	LeverageValue     int          `json:"leverageValue"`
}

// GetSide returns the position side based on the signed size
func (p DetailedPosition) GetSide() string {
	szi, err := strconv.ParseFloat(string(p.Szi), 64)
	if err != nil {
		return "none"
	}

	if szi > 0 {
		return "long"
	} else if szi < 0 {
		return "short"
	}
	return "none"
}

// GetAbsoluteSize returns the absolute size of the position
func (p DetailedPosition) GetAbsoluteSize() FloatString {
	szi, err := strconv.ParseFloat(string(p.Szi), 64)
	if err != nil {
		return FloatString("0")
	}

	if szi < 0 {
		szi = -szi
	}

	return FloatString(strconv.FormatFloat(szi, 'f', -1, 64))
}

// IsLong returns true if the position is long
func (p DetailedPosition) IsLong() bool {
	szi, err := strconv.ParseFloat(string(p.Szi), 64)
	if err != nil {
		return false
	}
	return szi > 0
}

// IsShort returns true if the position is short
func (p DetailedPosition) IsShort() bool {
	szi, err := strconv.ParseFloat(string(p.Szi), 64)
	if err != nil {
		return false
	}
	return szi < 0
}

// HasPosition returns true if there is an active position
func (p DetailedPosition) HasPosition() bool {
	szi, err := strconv.ParseFloat(string(p.Szi), 64)
	if err != nil {
		return false
	}
	return szi != 0
}

// ToSummary converts a DetailedPosition to a PositionSummary
func (p DetailedPosition) ToSummary() PositionSummary {
	return PositionSummary{
		Coin:              p.Coin,
		Size:              p.GetAbsoluteSize(),
		Side:              p.GetSide(),
		EntryPrice:        p.EntryPx,
		MarkPrice:         FloatString("0"), // Would need to be populated from market data
		UnrealizedPnl:     p.UnrealizedPnl,
		UnrealizedPnlPerc: FloatString("0"), // Would need to be calculated
		MarginUsed:        p.MarginUsed,
		LeverageType:      p.Leverage.Type,
		LeverageValue:     p.Leverage.Value,
	}
}

// PositionRisk represents position risk metrics
type PositionRisk struct {
	Coin              string       `json:"coin"`
	Size              FloatString  `json:"size"`
	NotionalValue     FloatString  `json:"notionalValue"`
	MarginUsed        FloatString  `json:"marginUsed"`
	MaintenanceMargin FloatString  `json:"maintenanceMargin"`
	LiquidationPrice  *FloatString `json:"liquidationPrice"`
	UnrealizedPnl     FloatString  `json:"unrealizedPnl"`
	MarginRatio       FloatString  `json:"marginRatio"`
	LeverageUsed      FloatString  `json:"leverageUsed"`
	RiskLevel         string       `json:"riskLevel"` // "low", "medium", "high", "critical"
}

// GetRiskLevel calculates the risk level based on margin ratio
func (pr PositionRisk) GetRiskLevel() string {
	marginRatio, err := strconv.ParseFloat(string(pr.MarginRatio), 64)
	if err != nil {
		return "unknown"
	}

	if marginRatio > 0.5 {
		return "low"
	} else if marginRatio > 0.25 {
		return "medium"
	} else if marginRatio > 0.1 {
		return "high"
	}
	return "critical"
}

// PositionHistory represents historical position data
type PositionHistory struct {
	Coin       string       `json:"coin"`
	OpenTime   int64        `json:"openTime"`
	CloseTime  *int64       `json:"closeTime,omitempty"`
	EntryPrice FloatString  `json:"entryPrice"`
	ExitPrice  *FloatString `json:"exitPrice,omitempty"`
	Size       FloatString  `json:"size"`
	Side       string       `json:"side"`
	Pnl        *FloatString `json:"pnl,omitempty"`
	Fee        FloatString  `json:"fee"`
	Duration   *int64       `json:"duration,omitempty"` // in milliseconds
}

// IsOpen returns true if the position is still open
func (ph PositionHistory) IsOpen() bool {
	return ph.CloseTime == nil
}

// GetDuration returns the duration of the position in milliseconds
func (ph PositionHistory) GetDuration() int64 {
	if ph.Duration != nil {
		return *ph.Duration
	}

	if ph.CloseTime != nil {
		return *ph.CloseTime - ph.OpenTime
	}

	// Position is still open, return current duration
	return GetTimestampMs() - ph.OpenTime
}

// PositionUpdate represents a real-time position update
type PositionUpdate struct {
	Coin          string       `json:"coin"`
	Size          FloatString  `json:"size"`
	Side          string       `json:"side"`
	MarkPrice     FloatString  `json:"markPrice"`
	UnrealizedPnl FloatString  `json:"unrealizedPnl"`
	MarginUsed    FloatString  `json:"marginUsed"`
	Timestamp     int64        `json:"timestamp"`
	LiquidationPx *FloatString `json:"liquidationPx,omitempty"`
}

// PositionFilter represents filters for querying positions
type PositionFilter struct {
	Coins     []string     `json:"coins,omitempty"`     // Filter by specific coins
	MinSize   *FloatString `json:"minSize,omitempty"`   // Minimum position size
	HasPnl    *bool        `json:"hasPnl,omitempty"`    // Filter positions with/without PnL
	Side      *string      `json:"side,omitempty"`      // Filter by side ("long", "short")
	RiskLevel *string      `json:"riskLevel,omitempty"` // Filter by risk level
}

// PositionMetrics represents aggregated position metrics
type PositionMetrics struct {
	TotalPositions     int         `json:"totalPositions"`
	TotalNotional      FloatString `json:"totalNotional"`
	TotalUnrealizedPnl FloatString `json:"totalUnrealizedPnl"`
	TotalMarginUsed    FloatString `json:"totalMarginUsed"`
	LongPositions      int         `json:"longPositions"`
	ShortPositions     int         `json:"shortPositions"`
	AvgLeverage        FloatString `json:"avgLeverage"`
	LargestPosition    *string     `json:"largestPosition,omitempty"` // Coin with largest notional
}

// CalculateMetrics calculates position metrics from a list of positions
func CalculateMetrics(positions []DetailedPosition) PositionMetrics {
	metrics := PositionMetrics{}

	var totalNotional, totalPnl, totalMargin, totalLeverage float64
	var longCount, shortCount int
	var largestNotional float64
	var largestCoin string

	for _, pos := range positions {
		if !pos.HasPosition() {
			continue
		}

		metrics.TotalPositions++

		// Parse values
		notional, _ := strconv.ParseFloat(string(pos.PositionValue), 64)
		pnl, _ := strconv.ParseFloat(string(pos.UnrealizedPnl), 64)
		margin, _ := strconv.ParseFloat(string(pos.MarginUsed), 64)
		leverage := float64(pos.Leverage.Value)

		totalNotional += notional
		totalPnl += pnl
		totalMargin += margin
		totalLeverage += leverage

		if notional > largestNotional {
			largestNotional = notional
			largestCoin = pos.Coin
		}

		if pos.IsLong() {
			longCount++
		} else if pos.IsShort() {
			shortCount++
		}
	}

	metrics.TotalNotional = FloatString(strconv.FormatFloat(totalNotional, 'f', -1, 64))
	metrics.TotalUnrealizedPnl = FloatString(strconv.FormatFloat(totalPnl, 'f', -1, 64))
	metrics.TotalMarginUsed = FloatString(strconv.FormatFloat(totalMargin, 'f', -1, 64))
	metrics.LongPositions = longCount
	metrics.ShortPositions = shortCount

	if metrics.TotalPositions > 0 {
		avgLev := totalLeverage / float64(metrics.TotalPositions)
		metrics.AvgLeverage = FloatString(strconv.FormatFloat(avgLev, 'f', 2, 64))
	}

	if largestCoin != "" {
		metrics.LargestPosition = &largestCoin
	}

	return metrics
}
