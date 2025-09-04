package notification

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/notification/controllers"
	"movie-app-go/internal/modules/notification/repositories"
	"movie-app-go/internal/modules/notification/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationModule struct {
	NotificationController *controllers.NotificationController
	NotificationService    *services.NotificationService
}

func NewNotificationModule(db *gorm.DB) *NotificationModule {
	notificationRepo := repositories.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationController := controllers.NewNotificationController(notificationService)

	return &NotificationModule{
		NotificationController: notificationController,
		NotificationService:    notificationService,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *NotificationModule) {
	userNotifications := rg.Group("/notifications")
	userNotifications.Use(middleware.Auth())
	{
		userNotifications.GET("", module.NotificationController.GetUserNotifications)
		userNotifications.GET("/stats", module.NotificationController.GetNotificationStats)
		userNotifications.GET("/:id", module.NotificationController.GetNotificationByID)
		userNotifications.PUT("/:id/read", module.NotificationController.MarkAsRead)
		userNotifications.PUT("/:id/unread", module.NotificationController.MarkAsUnread)
		userNotifications.PUT("/mark-all-read", module.NotificationController.MarkAllAsRead)
		userNotifications.PUT("/bulk-read", module.NotificationController.BulkMarkAsRead)
		userNotifications.DELETE("/:id", module.NotificationController.DeleteNotification)
	}

	adminNotifications := rg.Group("/admin/notifications")
	adminNotifications.Use(middleware.AdminOnly())
	{
		adminNotifications.POST("", module.NotificationController.CreateNotification)
		adminNotifications.POST("/system", module.NotificationController.CreateSystemNotification)
	}
}
