package jobs

import (
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypePaymentTimeout = "payment:timeout"
)

type QueueService struct {
	client *asynq.Client
}

func NewQueueService(redisAddr string) *QueueService {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &QueueService{client: client}
}

func (q *QueueService) SchedulePaymentTimeout(transactionID uint, delay time.Duration) error {
	payload, err := CreatePaymentTimeoutPayload(transactionID)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypePaymentTimeout, payload)

	// Schedule task to run after delay
	_, err = q.client.Enqueue(task, asynq.ProcessIn(delay))
	return err
}

func (q *QueueService) Close() error {
	return q.client.Close()
}
