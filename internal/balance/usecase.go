package balance

import (
	"context"

	"github.com/shopspring/decimal"
)

// UseCase Balance.
type UseCase interface {
	Update(context.Context, decimal.Decimal, int64) error
}
