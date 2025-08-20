// internal/modules/order/dependency.go
package order

import (
	"movie-app-go/internal/jobs"
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/order/controllers"
	"movie-app-go/internal/modules/order/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderModule struct {
	TransactionController *controllers.TransactionController
	TicketController      *controllers.TicketController
}

func NewOrderModule(db *gorm.DB, queueService *jobs.QueueService) *OrderModule {
	transactionService := services.NewTransactionService(db, queueService)
	ticketService := services.NewTicketService(db)

	return &OrderModule{
		TransactionController: controllers.NewTransactionController(transactionService),
		TicketController:      controllers.NewTicketController(ticketService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *OrderModule) {
	// Transaction routes
	rg.POST("/transactions", middleware.Auth(), module.TransactionController.Create)
	rg.GET("/transactions/my", middleware.Auth(), module.TransactionController.GetMyTransactions)
	rg.GET("/transactions", middleware.AdminOnly(), module.TransactionController.GetAll)
	rg.GET("/transactions/:id", middleware.Auth(), module.TransactionController.GetByID)
	rg.PUT("/transactions/:id/status", middleware.AdminOnly(), module.TransactionController.UpdateStatus)

	// Ticket routes
	rg.GET("/tickets/my", middleware.Auth(), module.TicketController.GetMyTickets)
	rg.GET("/tickets", middleware.AdminOnly(), module.TicketController.GetAll)
	rg.GET("/tickets/:id", middleware.Auth(), module.TicketController.GetByID)
	rg.PUT("/tickets/:id/status", middleware.Auth(), module.TicketController.UpdateStatus)
	rg.GET("/schedules/:schedule_id/tickets", middleware.AdminOnly(), module.TicketController.GetBySchedule)
}
