package rec

import "github.com/shopspring/decimal"

type Res struct {
	Text string             `json:"text"`
	Acc  decimal.Decimal    `json:"acc"`
	P0   [2]decimal.Decimal `json:"p0"` // left-top point
	P1   [2]decimal.Decimal `json:"p1"` // right-top point
	P2   [2]decimal.Decimal `json:"p2"` // right-bottom point
	P3   [2]decimal.Decimal `json:"p3"` // left-bottom point
}
