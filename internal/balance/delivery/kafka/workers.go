package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"

	"github.com/shaolinjehzu/go-queue-microservice/balance_service/internal/models"
)

// updateBalanceWorker processes the message, updates the balance and sends the response to one of the queues.
func (bcg *BalancesConsumerGroup) updateBalanceWorker(
	ctx context.Context,
	r *kafka.Reader,
	w *kafka.Writer,
	we *kafka.Writer,
	wg *sync.WaitGroup,
	m kafka.Message,
	workerID int,
) {
	defer wg.Done()

	bcg.log.Infof(
		"WORKER: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n",
		workerID,
		m.Topic,
		m.Partition,
		m.Offset,
		string(m.Key),
		string(m.Value),
	)

	var err error
	var task models.Task
	if err = json.Unmarshal(m.Value, &task); err != nil {
		bcg.log.Errorf("json.Unmarshal", err)
		return
	}

	if err = bcg.validate.StructCtx(ctx, task); err != nil {
		bcg.log.Errorf("validate.StructCtx", err)
		return
	}

	var amount decimal.Decimal

	switch task.Type {
	case "decrease":
		amount, err = decimal.NewFromString("-" + task.Value)
		if err != nil {
			bcg.log.Errorf("decimal.NewFromString", err)
			return
		}
	case "increase":
		amount, err = decimal.NewFromString(task.Value)
		if err != nil {
			bcg.log.Errorf("decimal.NewFromString", err)
			return
		}
	default:
		return
	}

	if err = bcg.balanceUC.Update(ctx, amount, task.Id); err != nil {
		if err := bcg.publishErrorMessage(ctx, we, m, err); err != nil {
			bcg.log.Errorf("publishErrorMessage", err)
		}
		bcg.log.Errorf("balanceUC.Update.publishErrorMessage", err)
		return
	}

	if err := r.CommitMessages(ctx, m); err != nil {
		bcg.log.Errorf("CommitMessages", err)
		return
	}

	if err = bcg.publishSuccessMessage(ctx, w, m, task); err != nil {
		bcg.log.Errorf("publishSuccessMessage", err)
		return
	}
}
