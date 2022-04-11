package usecase

import (
	"context"

	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/balance"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/pkg/logger"
	"github.com/shopspring/decimal"
)

// balanceUC.
type balanceUC struct {
	repo balance.Repository
	log  logger.Logger
}

// NewBalanceUC constructor.
func NewBalanceUC(repo balance.Repository, log logger.Logger) *balanceUC {
	return &balanceUC{
		repo: repo,
		log:  log,
	}
}

// Update single balance.
func (b *balanceUC) Update(ctx context.Context, amount decimal.Decimal, userId int64) error {
	return b.repo.Update(ctx, amount, userId)
}