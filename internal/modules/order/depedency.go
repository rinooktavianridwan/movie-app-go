// internal/modules/order/dependency.go
package order

import (
	"movie-app-go/internal/jobs"
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/order/controllers"
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
	transactionService := services.NewTransactionService(db, queueService, promoService)
	ticketService := services.NewTicketService(db)

	return &OrderModule{
		TransactionController: controllers.NewTransactionController(transactionService),
		TicketController:      controllers.NewTicketController(ticketService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *OrderModule) {
	// Transaction
	rg.POST("/transactions", middleware.Auth(), module.TransactionController.Create)
	rg.GET("/transactions/my", middleware.Auth(), module.TransactionController.GetMyTransactions)
	rg.GET("/transactions", middleware.AdminOnly(), module.TransactionController.GetAll)
	rg.GET("/transactions/:id", middleware.Auth(), module.TransactionController.GetByID)
	rg.POST("/transactions/:id/payment", middleware.AdminOnly(), module.TransactionController.ProcessPayment)

	// Ticket
	rg.GET("/tickets/my", middleware.Auth(), module.TicketController.GetMyTickets)
	rg.GET("/tickets", middleware.AdminOnly(), module.TicketController.GetAll)
	rg.GET("/tickets/:id", middleware.Auth(), module.TicketController.GetByID)
	rg.GET("/tickets/by-schedule/:schedule_id", middleware.AdminOnly(), module.TicketController.GetBySchedule)
	rg.POST("/tickets/:id/scan", middleware.Auth(), module.TicketController.ScanTicket)
}
