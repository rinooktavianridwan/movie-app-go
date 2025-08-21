package jobs

import (
	"log"

	"movie-app-go/internal/constants"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type WorkerService struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewWorkerService(redisAddr string, db *gorm.DB) *WorkerService {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	paymentHandler := NewPaymentJobHandler(db)

	mux.HandleFunc(constants.TypePaymentTimeout, paymentHandler.HandlePaymentTimeout)

	return &WorkerService{
		server: server,
		mux:    mux,
	}
}

func (w *WorkerService) Start() error {
	log.Println("Starting Asynq worker...")
	return w.server.Start(w.mux)
}

func (w *WorkerService) Shutdown() {
	log.Println("Shutting down Asynq worker...")
	w.server.Shutdown()
}
