package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"


	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/balance"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

	"github.com/shaolinjehzu/go-queue-microservice/balance_service/config"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/models"
	"github.com/shaolinjehzu/go-queue-microservice/balance_service/pkg/logger"
)

// BalancesConsumerGroup struct.
type BalancesConsumerGroup struct {
	Brokers   	[]string
	GroupID   	string
	workers   	int
	log       	logger.Logger
	cfg       	*config.Config
	balanceUC 	balance.UseCase
	validate  	*validator.Validate
	wg			*sync.WaitGroup
}

// NewBalancesConsumerGroup constructor.
func NewBalancesConsumerGroup(
	brokers []string,
	groupID string,
	log logger.Logger,
	cfg *config.Config,
	balanceUC balance.UseCase,
	validate *validator.Validate,
	wg *sync.WaitGroup,
) *BalancesConsumerGroup {
	return &BalancesConsumerGroup{
		Brokers:   	brokers,
		GroupID:   	groupID,
		workers:   	0,
		log:       	log,
		cfg:       	cfg,
		balanceUC: 	balanceUC,
		validate:  	validate,
		wg: 		wg,
	}
}

// getNewKafkaReader Create new kafka reader.
func (bcg *BalancesConsumerGroup) getNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(bcg.log.Debugf),
		ErrorLogger:            kafka.LoggerFunc(bcg.log.Errorf),
		MaxAttempts:            maxAttempts,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}

// getNewKafkaWriter Create new kafka writer.
func (bcg *BalancesConsumerGroup) getNewKafkaWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(bcg.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(bcg.log.Debugf),
		ErrorLogger:  kafka.LoggerFunc(bcg.log.Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
	return w
}

// consumeUpdateBalance run consumer for update-balance topic and spawn.
func (bcg *BalancesConsumerGroup) consumeUpdateBalance(
	ctx context.Context,
	cancel context.CancelFunc,
	groupID string,
	topic string,
) {
	r := bcg.getNewKafkaReader(bcg.Brokers, topic, groupID)
	defer cancel()
	defer func() {
		if err := r.Close(); err != nil {
			bcg.log.Errorf("r.Close", err)
			cancel()
		}
	}()

	w := bcg.getNewKafkaWriter(successLetterQueueTopic)
	defer func() {
		if err := w.Close(); err != nil {
			bcg.log.Errorf("w.Close", err)
			cancel()
		}
	}()

	we := bcg.getNewKafkaWriter(deadLetterQueueTopic)
	defer func() {
		if err := we.Close(); err != nil {
			bcg.log.Errorf("w.Close", err)
			cancel()
		}
	}()

	bcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			fmt.Println("ERROR", err)
			continue
		}
		bcg.wg.Add(1)
		bcg.workers++
		go bcg.updateBalanceWorker(ctx, r, w, we, bcg.wg, m, bcg.workers)
	}

}

// publishErrorMessage publish messages to dead-letter-queue topic.
func (bcg *BalancesConsumerGroup) publishErrorMessage(ctx context.Context, w *kafka.Writer, m kafka.Message, err error) error {
	errMsg := &models.ErrorMessage{
		Offset:    m.Offset,
		Error:     err.Error(),
		Time:      m.Time.UTC(),
		Partition: m.Partition,
		Topic:     m.Topic,
	}

	errMsgBytes, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}

	return w.WriteMessages(ctx, kafka.Message{
		Value: errMsgBytes,
	})
}

// publishSuccessMessage publish messages to success-letter-queue topic.
func (bcg *BalancesConsumerGroup) publishSuccessMessage(ctx context.Context, w *kafka.Writer, m kafka.Message, task models.Task) error {
	msg := &models.SuccessMessage{
		Offset:    m.Offset,
		Time:      m.Time.UTC(),
		Partition: m.Partition,
		Task:      task,
		Topic:     m.Topic,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return w.WriteMessages(ctx, kafka.Message{
		Value: msgBytes,
	})
}

// RunConsumers run kafka consumers.
func (bcg *BalancesConsumerGroup) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	go bcg.consumeUpdateBalance(ctx, cancel, balanceGroupID, updateBalanceTopic)
}
