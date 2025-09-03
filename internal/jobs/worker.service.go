package jobs

import (
	"log"

	"movie-app-go/internal/constants"
	notificationJobs "movie-app-go/internal/modules/notification/jobs"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type WorkerService struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	DB     *gorm.DB
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
	ws := &WorkerService{
		server: server,
		mux:    mux,
		DB:     db,
	}

	ws.registerHandlers()

	return ws
}

func (w *WorkerService) Start() error {
	log.Println("Starting Asynq worker...")
	return w.server.Start(w.mux)
}

func (w *WorkerService) Shutdown() {
	log.Println("Shutting down Asynq worker...")
	w.server.Shutdown()
}

func (w *WorkerService) registerHandlers() {
	paymentHandler := NewPaymentJobHandler(w.DB)
	w.mux.HandleFunc(constants.TypePaymentTimeout, paymentHandler.HandlePaymentTimeout)

	notificationHandler := notificationJobs.NewNotificationJobHandler(w.DB)
	w.mux.HandleFunc(constants.JobTypeMovieReminder, notificationHandler.HandleMovieReminder)
	w.mux.HandleFunc(constants.JobTypePromoNotification, notificationHandler.HandlePromoNotification)
	w.mux.HandleFunc(constants.JobTypeBookingConfirm, notificationHandler.HandleBookingConfirmation)
}
