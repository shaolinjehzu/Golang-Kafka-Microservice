package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/config"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/balance/delivery/kafka"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/balance/repository"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/balance/usecase"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/pkg/logger"
)

const (
	kafkaGroupID = "balances_group"
)

// server.
type server struct {
	log logger.Logger
	cfg *config.Config
	db  *sqlx.DB
}

// NewServer constructor.
func NewServer(log logger.Logger, cfg *config.Config, db *sqlx.DB) *server {
	return &server{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

// Run start server
func (s *server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	validate := validator.New()

	repo := repository.NewBalancePostgresRepo(s.db)

	balanceUC := usecase.NewBalanceUC(repo, s.log)

	balancesCG := kafka.NewBalancesConsumerGroup(s.cfg.Kafka.Brokers, kafkaGroupID, s.log, s.cfg, balanceUC, validate, wg)

	balancesCG.RunConsumers(ctx, cancel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.log.Errorf("ctx.Done: %v", done)
	}

	wg.Wait()
}
