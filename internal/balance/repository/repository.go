package repository

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/models"
	"github.com/shopspring/decimal"
)

// balancePostgresRepo.
type balancePostgresRepo struct {
	db *sqlx.DB
}

// NewBalancePostgresRepo constructor.
func NewBalancePostgresRepo(db *sqlx.DB) *balancePostgresRepo {
	return &balancePostgresRepo{db: db}
}

// Update single balance.
func (b *balancePostgresRepo) Update(ctx context.Context, amount decimal.Decimal, userId int64) error {
	var hasCommitted = false

	for ok := true; ok; ok = !hasCommitted {
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}

		defer tx.Rollback()

		_, err = tx.Exec(`set transaction isolation level repeatable read`)
		if err != nil {
			return err
		}

		balance := &models.Balance{}
		if err = tx.QueryRow(
			"SELECT balance FROM balance WHERE user_id = $1",
			userId,
		).Scan(
			&balance.Balance,
		); err != nil {
			return err
		}

		newBalance := balance.Balance.Add(amount)

		sqlStatement := `
		UPDATE balance
		SET balance = $1
		WHERE user_id = $2;`

		_, err = tx.Exec(sqlStatement, newBalance, userId)
		if err != nil {
			if strings.Contains(err.Error(), "could not serialize access due to concurrent update") || strings.Contains(err.Error(), "не удалось сериализовать доступ из-за параллельного изменения") {
				continue
			} else {
				return err
			}
		}

		tx.Commit()
		hasCommitted = true
	}

	return nil
}
