package postgres

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/config"
)

const (
	retryAttempts = 5
	retryDelay    = 2 * time.Second
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

func NewPostgresDB(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	dbUrl := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PostgresDB.Host,
		cfg.PostgresDB.Port,
		cfg.PostgresDB.DB,
		cfg.PostgresDB.User,
		cfg.PostgresDB.Password,
		cfg.PostgresDB.SSL,
	)

	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}


	if err = retry.Do(func() error {
		return db.PingContext(ctx)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		return nil, err
	}

	return db, nil
}
