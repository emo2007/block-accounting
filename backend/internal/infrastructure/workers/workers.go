package workers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/infrastructure/queue"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	id  string
	log *slog.Logger

	rmqc *amqp.Connection

	queueName string
}

func (w *Worker) Run() {
	w.log = w.log.With(slog.String("worker-id", w.id), slog.String("worker-queue", w.queueName))

	defer func() {
		if p := recover(); p != nil {
			w.log.Error(
				"worker paniced!",
				slog.String("worker id", w.id),
				slog.Any("panic", p),
			)
		} else {
			w.log.Info("worker stoped. bye bye 0w0", slog.String("worker id", w.id))
		}
	}()

	channel, err := w.rmqc.Channel()
	if err != nil {
		w.log.Error("error create rmq channel", logger.Err(err))
		return
	}

	delivery, err := channel.Consume(
		w.queueName,
		w.id,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		w.log.Error("error consume from rmq channel", logger.Err(err))
		return
	}

	w.handleJobs(delivery)
}

func (w *Worker) handleJobs(ch <-chan amqp.Delivery) {
	for msg := range ch {
		w.log.Debug("job received", slog.Any("job", msg.MessageId))

		var job queue.Job

		if err := json.Unmarshal(msg.Body, &job); err != nil {
			w.log.Error("error parse message body. %w", err)
			continue
		}

		// TODO check job.IdempotentKey for duplicate

		// TODO dispatch job
	}
}

func (w *Worker) handleDeployMultisig(
	ctx context.Context,
)
