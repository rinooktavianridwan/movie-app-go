package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/notification/options"
	"movie-app-go/internal/modules/notification/requests"
	"movie-app-go/internal/modules/notification/responses"
	"movie-app-go/internal/modules/notification/services"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: notificationService,
	}
}

func (c *NotificationController) getUserID(ctx *gin.Context) (uint, error) {
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: fmt.Errorf("user not authenticated"), Type: gin.ErrorTypePublic}
	}

	switch v := userIDInterface.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	default:
		return 0, gin.Error{Err: fmt.Errorf("invalid user ID type"), Type: gin.ErrorTypePublic}
	}
}

func (c *NotificationController) GetUserNotifications(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	opts, err := options.ParseNotificationOptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.NotificationService.GetUserNotifications(userID, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notificationResponses := responses.ToNotificationResponses(result.Data)
	paginatedResponse := responses.PaginatedNotificationResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      notificationResponses,
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

func (c *NotificationController) GetNotificationByID(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	notification, err := c.NotificationService.GetNotificationByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	if notification.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	notificationResponse := responses.ToNotificationResponse(*notification)
	ctx.JSON(http.StatusOK, notificationResponse)
}

func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	notification, err := c.NotificationService.MarkAsRead(uint(id), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	notificationResponse := responses.ToNotificationResponse(*notification)
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "notification marked as read",
		"notification": notificationResponse,
	})
}

func (c *NotificationController) MarkAsUnread(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	notification, err := c.NotificationService.MarkAsUnread(uint(id), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	notificationResponse := responses.ToNotificationResponse(*notification)
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "notification marked as unread",
		"notification": notificationResponse,
	})
}

func (c *NotificationController) BulkMarkAsRead(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req requests.BulkMarkReadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.NotificationService.BulkMarkAsRead(&req, userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "notifications marked as read"})
}

func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := c.NotificationService.MarkAllAsRead(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

func (c *NotificationController) GetNotificationStats(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	stats, err := c.NotificationService.GetNotificationStats(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := c.NotificationService.DeleteNotification(uint(id), userID); err != nil {
		if err.Error() == "notification not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "notification deleted successfully"})
}

func (c *NotificationController) CreateNotification(ctx *gin.Context) {
	var req requests.CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := c.NotificationService.CreateNotification(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notificationResponse := responses.ToNotificationResponse(*notification)
	ctx.JSON(http.StatusCreated, notificationResponse)
}

func (c *NotificationController) CreateSystemNotification(ctx *gin.Context) {
	var req requests.CreateSystemNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userIDs []uint
	if err := c.NotificationService.DB.Model(&models.User{}).
		Where("is_admin = ?", false).
		Pluck("id", &userIDs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	if len(userIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no users found"})
		return
	}

	err := c.NotificationService.CreateBulkNotifications(
		userIDs,
		req.Title,
		req.Message,
		constants.NotificationTypeSystem,
		req.Data,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    fmt.Sprintf("system notification sent to %d users", len(userIDs)),
		"user_count": len(userIDs),
	})
}
