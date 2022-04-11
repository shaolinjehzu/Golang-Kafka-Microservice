package models

import (
	"github.com/shopspring/decimal"
)

// Balance model.
type Balance struct {
	UserId  int64
	Balance decimal.Decimal
}
