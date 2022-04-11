package main

import (
	"context"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/config"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/server"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/pkg/kafka"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/pkg/logger"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/pkg/postgres"
	"log"

)

func main(){

	log.Println("Starting balance microservice")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	postgresDB, err := postgres.NewPostgresDB(ctx, cfg)
	if err != nil {
		appLogger.Fatal("cannot connect postgres", err.Error())
	}

	defer func() {
		if err := postgresDB.Close(); err != nil {
			appLogger.Fatal("Postgres.Disconnect", err)
		}
	}()

	appLogger.Info("Postgres connected")

	conn, err := kafka.NewKafkaConn(cfg)
	if err != nil {
		appLogger.Fatal("NewKafkaConn", err)
	}
	defer conn.Close()

	brokers, err := conn.Brokers()
	if err != nil {
		appLogger.Fatal("conn.Brokers", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)

	s := server.NewServer(appLogger, cfg, postgresDB)
	s.Run()

}
