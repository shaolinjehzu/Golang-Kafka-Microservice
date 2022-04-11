package balance

import (
	"context"

	"github.com/shopspring/decimal"
)

// Repository Balance.
type Repository interface {
	Update(ctx context.Context, amount decimal.Decimal, userId int64) error
}
