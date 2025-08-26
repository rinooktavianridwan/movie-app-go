package controllers

import (
    "net/http"
    "strconv"

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

func (c *NotificationController) GetUserNotifications(ctx *gin.Context) {
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
        return
    }

    if err := c.NotificationService.MarkAllAsRead(userID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

func (c *NotificationController) GetNotificationStats(ctx *gin.Context) {
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
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