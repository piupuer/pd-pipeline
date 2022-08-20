package rec

import "github.com/shopspring/decimal"

type Res struct {
	Text   string                `json:"text"`
	Acc    decimal.Decimal       `json:"acc"`
	Points [4][2]decimal.Decimal `json:"points"`
}
