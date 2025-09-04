// internal/modules/order/dependency.go
package order

import (
	"movie-app-go/internal/jobs"
	"movie-app-go/internal/middleware"
	notificationRepositories "movie-app-go/internal/modules/notification/repositories"
	notificationServices "movie-app-go/internal/modules/notification/services"
	"movie-app-go/internal/modules/order/controllers"
	"movie-app-go/internal/modules/order/repositories"
	"movie-app-go/internal/modules/order/services"
	promoServices "movie-app-go/internal/modules/promo/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderModule struct {
	TransactionController *controllers.TransactionController
	TicketController      *controllers.TicketController
	TicketService         *services.TicketService
	PromoService          *promoServices.PromoService
}

func NewOrderModule(db *gorm.DB, queueService *jobs.QueueService, promoService *promoServices.PromoService) *OrderModule {
	transactionRepo := repositories.NewTransactionRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)

	notificationRepo := notificationRepositories.NewNotificationRepository(db)
	notificationService := notificationServices.NewNotificationService(notificationRepo)

	transactionService := services.NewTransactionService(
		transactionRepo,
		queueService,
		promoService,
		notificationService,
	)
	ticketService := services.NewTicketService(ticketRepo)

	return &OrderModule{
		TransactionController: controllers.NewTransactionController(transactionService),
		TicketController:      controllers.NewTicketController(ticketService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *OrderModule, mf *middleware.Factory) {
	// Transaction routes
	rg.POST("/transactions", mf.Auth(), module.TransactionController.Create)
	rg.GET("/transactions/my", mf.Auth(), module.TransactionController.GetMyTransactions)
	rg.GET("/transactions", mf.Auth(), mf.RequirePermission("orders.read"), module.TransactionController.GetAll)
	rg.GET("/transactions/:id", mf.Auth(), module.TransactionController.GetByID)
	rg.POST("/transactions/:id/payment", mf.Auth(), mf.RequirePermission("orders.update"), module.TransactionController.ProcessPayment)

	// Ticket routes
	rg.GET("/tickets/my", mf.Auth(), module.TicketController.GetMyTickets)
	rg.GET("/tickets", mf.Auth(), mf.RequirePermission("orders.read"), module.TicketController.GetAll)
	rg.GET("/tickets/:id", mf.Auth(), module.TicketController.GetByID)
	rg.GET("/tickets/by-schedule/:schedule_id", mf.Auth(), mf.RequirePermission("orders.read"), module.TicketController.GetBySchedule)
	rg.POST("/tickets/:id/scan", mf.Auth(), module.TicketController.ScanTicket)
}
